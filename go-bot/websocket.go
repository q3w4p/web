package discord

import (
	"bytes"
	"compress/zlib"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"sync"
	"time"
)

const (
	// Gateway opcodes
	OpDispatch            = 0
	OpHeartbeat           = 1
	OpIdentify            = 2
	OpPresenceUpdate      = 3
	OpVoiceStateUpdate    = 4
	OpResume              = 6
	OpReconnect           = 7
	OpRequestGuildMembers = 8
	OpInvalidSession      = 9
	OpHello               = 10
	OpHeartbeatACK        = 11
)

// WebSocketClient handles WebSocket connections to Discord gateway
type WebSocketClient struct {
	conn           net.Conn
	url            string
	headers        map[string]string
	isConnected    bool
	mutex          sync.RWMutex
	heartbeatMutex sync.Mutex

	// Gateway state
	sessionID         string
	sequence          int64
	heartbeatInterval int
	lastHeartbeatACK  int64
	identified        bool
	resumed           bool

	// Handlers
	messageHandlers map[string][]func(map[string]interface{})
	closeHandler    func()

	// Control
	stopHeartbeat chan bool
	zlibReader    io.ReadCloser
}

// GatewayPayload represents a Discord gateway payload
type GatewayPayload struct {
	Op int         `json:"op"`
	D  interface{} `json:"d"`
	S  *int64      `json:"s,omitempty"`
	T  *string     `json:"t,omitempty"`
}

// NewWebSocketClient creates a new WebSocket client
func NewWebSocketClient(gatewayURL string, headers map[string]string) *WebSocketClient {
	return &WebSocketClient{
		url:             gatewayURL,
		headers:         headers,
		messageHandlers: make(map[string][]func(map[string]interface{})),
		stopHeartbeat:   make(chan bool),
	}
}

// Connect establishes WebSocket connection
func (ws *WebSocketClient) Connect() error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	u, err := url.Parse(ws.url)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	host := u.Host
	if u.Port() == "" {
		if u.Scheme == "wss" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	// Connect with TLS
	tlsConfig := &tls.Config{
		ServerName:         u.Hostname(),
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", host, tlsConfig)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	ws.conn = conn

	// Perform WebSocket handshake
	if err := ws.performHandshake(u); err != nil {
		conn.Close()
		return fmt.Errorf("handshake failed: %w", err)
	}

	ws.isConnected = true

	// Start message reader
	go ws.readMessages()

	return nil
}

// performHandshake performs WebSocket handshake
func (ws *WebSocketClient) performHandshake(u *url.URL) error {
	key := make([]byte, 16)
	rand.Read(key)
	wsKey := base64.StdEncoding.EncodeToString(key)

	path := u.Path
	if path == "" {
		path = "/"
	}
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	handshake := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
	handshake += fmt.Sprintf("Host: %s\r\n", u.Host)
	handshake += "Upgrade: websocket\r\n"
	handshake += "Connection: Upgrade\r\n"
	handshake += fmt.Sprintf("Sec-WebSocket-Key: %s\r\n", wsKey)
	handshake += "Sec-WebSocket-Version: 13\r\n"

	// Add custom headers
	for k, v := range ws.headers {
		handshake += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	handshake += "\r\n"

	if _, err := ws.conn.Write([]byte(handshake)); err != nil {
		return err
	}

	// Read handshake response
	buf := make([]byte, 4096)
	n, err := ws.conn.Read(buf)
	if err != nil {
		return err
	}

	response := string(buf[:n])
	if !bytes.Contains([]byte(response), []byte("HTTP/1.1 101")) {
		return fmt.Errorf("handshake failed: %s", response)
	}

	return nil
}

// readMessages reads and processes WebSocket messages
func (ws *WebSocketClient) readMessages() {
	defer func() {
		ws.mutex.Lock()
		ws.isConnected = false
		ws.mutex.Unlock()

		if ws.closeHandler != nil {
			ws.closeHandler()
		}
	}()

	for {
		ws.mutex.RLock()
		conn := ws.conn
		ws.mutex.RUnlock()

		if conn == nil {
			break
		}

		conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		// Read WebSocket frame
		data, err := ws.readFrame()
		if err != nil {
			break
		}

		if data == nil {
			continue
		}

		// Handle compressed data
		var payload []byte
		if ws.zlibReader != nil {
			// Zlib compressed
			payload, err = ws.decompressZlib(data)
			if err != nil {
				fmt.Printf("Zlib decompression error: %v\n", err)
				continue
			}
		} else {
			payload = data
		}

		// Parse payload
		var gatewayPayload GatewayPayload
		if err := json.Unmarshal(payload, &gatewayPayload); err != nil {
			fmt.Printf("JSON unmarshal error: %v\n", err)
			continue
		}

		// Handle opcode
		ws.handleOpcode(&gatewayPayload)
	}
}

// readFrame reads a WebSocket frame
func (ws *WebSocketClient) readFrame() ([]byte, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(ws.conn, header); err != nil {
		return nil, err
	}

	fin := (header[0] & 0x80) != 0
	opcode := header[0] & 0x0F
	masked := (header[1] & 0x80) != 0
	payloadLen := int64(header[1] & 0x7F)

	// Handle different payload lengths
	if payloadLen == 126 {
		lenBytes := make([]byte, 2)
		if _, err := io.ReadFull(ws.conn, lenBytes); err != nil {
			return nil, err
		}
		payloadLen = int64(binary.BigEndian.Uint16(lenBytes))
	} else if payloadLen == 127 {
		lenBytes := make([]byte, 8)
		if _, err := io.ReadFull(ws.conn, lenBytes); err != nil {
			return nil, err
		}
		payloadLen = int64(binary.BigEndian.Uint64(lenBytes))
	}

	// Read mask key if masked
	var maskKey []byte
	if masked {
		maskKey = make([]byte, 4)
		if _, err := io.ReadFull(ws.conn, maskKey); err != nil {
			return nil, err
		}
	}

	// Read payload
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(ws.conn, payload); err != nil {
		return nil, err
	}

	// Unmask if needed
	if masked {
		for i := range payload {
			payload[i] ^= maskKey[i%4]
		}
	}

	// Handle opcode
	switch opcode {
	case 0x1: // Text frame
		if !fin {
			// TODO: Handle fragmented messages
			return payload, nil
		}
		return payload, nil
	case 0x2: // Binary frame
		return payload, nil
	case 0x8: // Close frame
		return nil, fmt.Errorf("connection closed by server")
	case 0x9: // Ping frame
		ws.sendPong(payload)
		return nil, nil
	case 0xA: // Pong frame
		return nil, nil
	default:
		return payload, nil
	}
}

// writeFrame writes a WebSocket frame
func (ws *WebSocketClient) writeFrame(opcode byte, data []byte) error {
	ws.mutex.RLock()
	conn := ws.conn
	ws.mutex.RUnlock()

	if conn == nil {
		return fmt.Errorf("connection closed")
	}

	header := []byte{0x80 | opcode} // FIN bit + opcode

	payloadLen := len(data)
	if payloadLen < 126 {
		header = append(header, 0x80|byte(payloadLen)) // Mask bit + length
	} else if payloadLen < 65536 {
		header = append(header, 0x80|126)
		header = append(header, byte(payloadLen>>8), byte(payloadLen))
	} else {
		header = append(header, 0x80|127)
		for i := 7; i >= 0; i-- {
			header = append(header, byte(payloadLen>>(i*8)))
		}
	}

	// Generate mask key
	maskKey := make([]byte, 4)
	rand.Read(maskKey)
	header = append(header, maskKey...)

	// Mask data
	maskedData := make([]byte, len(data))
	for i := range data {
		maskedData[i] = data[i] ^ maskKey[i%4]
	}

	// Write frame
	frame := append(header, maskedData...)
	_, err := conn.Write(frame)
	return err
}

// SendJSON sends a JSON payload
func (ws *WebSocketClient) SendJSON(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ws.mutex.RLock()
	conn := ws.conn
	ws.mutex.RUnlock()

	if conn == nil {
		return fmt.Errorf("not connected")
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	return ws.writeFrame(0x1, data) // Text frame
}

// sendPong sends a pong frame
func (ws *WebSocketClient) sendPong(data []byte) error {
	return ws.writeFrame(0xA, data)
}

// handleOpcode handles Discord gateway opcodes
func (ws *WebSocketClient) handleOpcode(payload *GatewayPayload) {
	switch payload.Op {
	case OpHello:
		// Start heartbeating
		if data, ok := payload.D.(map[string]interface{}); ok {
			if interval, ok := data["heartbeat_interval"].(float64); ok {
				ws.heartbeatInterval = int(interval)
				go ws.startHeartbeat()
			}
		}

		// Send identify or resume
		if ws.sessionID != "" && ws.sequence > 0 {
			ws.sendResume()
		} else {
			// Identify handled by client
		}

	case OpHeartbeatACK:
		ws.heartbeatMutex.Lock()
		ws.lastHeartbeatACK = time.Now().UnixNano() / 1000000
		ws.heartbeatMutex.Unlock()

	case OpDispatch:
		if payload.S != nil {
			ws.sequence = *payload.S
		}

		if payload.T != nil {
			eventName := *payload.T

			// Handle READY event
			if eventName == "READY" {
				if data, ok := payload.D.(map[string]interface{}); ok {
					if sid, ok := data["session_id"].(string); ok {
						ws.sessionID = sid
					}
				}
			}

			// Call event handlers
			if handlers, ok := ws.messageHandlers[eventName]; ok {
				if data, ok := payload.D.(map[string]interface{}); ok {
					for _, handler := range handlers {
						go handler(data)
					}
				}
			}

			// Call wildcard handlers
			if handlers, ok := ws.messageHandlers["*"]; ok {
				if data, ok := payload.D.(map[string]interface{}); ok {
					for _, handler := range handlers {
						go handler(data)
					}
				}
			}
		}

	case OpInvalidSession:
		// Session invalidated
		canResume := false
		if b, ok := payload.D.(bool); ok {
			canResume = b
		}

		if !canResume {
			ws.sessionID = ""
			ws.sequence = 0
		}

		// Wait a bit then reconnect
		time.Sleep(time.Duration(1+ws.randomInt(5)) * time.Second)

	case OpReconnect:
		// Server requested reconnect
		ws.Close()
		// Reconnection handled by client
	}
}

// startHeartbeat starts the heartbeat goroutine
func (ws *WebSocketClient) startHeartbeat() {
	// Initial random jitter
	time.Sleep(time.Duration(ws.randomInt(ws.heartbeatInterval)) * time.Millisecond)

	ticker := time.NewTicker(time.Duration(ws.heartbeatInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ws.sendHeartbeat()

			// Check if we missed a heartbeat ACK
			ws.heartbeatMutex.Lock()
			lastACK := ws.lastHeartbeatACK
			ws.heartbeatMutex.Unlock()

			if lastACK > 0 && time.Now().UnixNano()/1000000-lastACK > int64(ws.heartbeatInterval)*2 {
				// Missed heartbeat, reconnect
				ws.Close()
				return
			}

		case <-ws.stopHeartbeat:
			return
		}
	}
}

// sendHeartbeat sends a heartbeat
func (ws *WebSocketClient) sendHeartbeat() error {
	payload := &GatewayPayload{
		Op: OpHeartbeat,
		D:  ws.sequence,
	}

	return ws.SendJSON(payload)
}

// sendResume sends a resume payload
func (ws *WebSocketClient) sendResume() error {
	// Resume data set by client
	return nil
}

// decompressZlib decompresses zlib data
func (ws *WebSocketClient) decompressZlib(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// On registers an event handler
func (ws *WebSocketClient) On(event string, handler func(map[string]interface{})) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	ws.messageHandlers[event] = append(ws.messageHandlers[event], handler)
}

// OnClose registers a close handler
func (ws *WebSocketClient) OnClose(handler func()) {
	ws.closeHandler = handler
}

// Close closes the WebSocket connection
func (ws *WebSocketClient) Close() error {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if ws.stopHeartbeat != nil {
		close(ws.stopHeartbeat)
	}

	if ws.conn != nil {
		// Send close frame
		closeFrame := []byte{0x03, 0xE8} // Code 1000 (normal closure)
		ws.writeFrame(0x8, closeFrame)

		ws.conn.Close()
		ws.conn = nil
	}

	ws.isConnected = false
	return nil
}

// IsConnected returns connection status
func (ws *WebSocketClient) IsConnected() bool {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	return ws.isConnected
}

// GetSessionID returns the session ID
func (ws *WebSocketClient) GetSessionID() string {
	return ws.sessionID
}

// GetSequence returns the sequence number
func (ws *WebSocketClient) GetSequence() int64 {
	return ws.sequence
}

// randomInt generates a random integer
func (ws *WebSocketClient) randomInt(max int) int {
	b := make([]byte, 4)
	rand.Read(b)
	return int(binary.BigEndian.Uint32(b)) % max
}
