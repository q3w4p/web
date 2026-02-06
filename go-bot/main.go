package main

import (
	"bytes"
	"discord-selfbot/discord"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

const TOKEN = ""

type Selfbot struct {
	client     *discord.DiscordClient
	prefix     string
	hostedBots map[string]*discord.DiscordClient
	gateway    *discord.WebSocketClient

	streaming        bool
	streamDetails    string
	streamState      string
	streamName       string
	streamImageKey   string
	streamImageText  string
	streamMode       string
	startTime        time.Time
	reconnectAttempt int
	gatewayActive    bool
	heartbeatTicker  *time.Ticker
	lastGatewayEvent time.Time

	// Reaction tracking
	reactTargets map[string]string // userID -> emoji
	reactMutex   sync.RWMutex
}

func NewSelfbot(token string) *Selfbot {
	return &Selfbot{
		client:           discord.NewDiscordClient(token),
		prefix:           ">",
		hostedBots:       make(map[string]*discord.DiscordClient),
		reactTargets:     make(map[string]string),
		streamMode:       "triple",
		startTime:        time.Now(),
		lastGatewayEvent: time.Now(),
		reconnectAttempt: 0,
		gatewayActive:    false,
	}
}

func loadASCII() string {
	data, err := os.ReadFile("ascii.txt")
	if err != nil {
		return "=== Discord Selfbot ==="
	}
	return string(data)
}

func (sb *Selfbot) Start() error {
	user, err := sb.client.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	fmt.Println("=== Discord Selfbot ===")
	fmt.Print("\033[35m")
	fmt.Println(loadASCII())
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("Logged in as: %s#%s\n", user.Username, user.Discriminator)
	fmt.Printf("Commands Available: 14\n")
	fmt.Printf("Prefix: %s\n", sb.prefix)
	fmt.Printf("Gateway: Connecting...\n")
	fmt.Printf("Stream Mode: %s\n", sb.streamMode)
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Print("\033[0m")

	if err := sb.connectGateway(); err != nil {
		fmt.Println("[!] Gateway connection failed, using polling mode only")
	} else {
		fmt.Println("[+] Gateway mode active")
		sb.gatewayActive = true
	}

	go sb.pollMessages()
	go sb.monitorGateway()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	<-sigch

	fmt.Println("\n[*] Shutting down...")
	if sb.heartbeatTicker != nil {
		sb.heartbeatTicker.Stop()
	}
	if sb.gateway != nil {
		sb.gateway.Close()
	}
	sb.client.Close()

	return nil
}

func (sb *Selfbot) monitorGateway() {
	checkInterval := 5 * time.Second
	activityTimeout := 120 * time.Second
	
	for {
		time.Sleep(checkInterval)

		isConnected := sb.gateway != nil && sb.gateway.IsConnected()
		
		timeSinceLastEvent := time.Since(sb.lastGatewayEvent)
		if isConnected && sb.gatewayActive && timeSinceLastEvent > activityTimeout {
			fmt.Printf("[!] No gateway activity for %v - forcing reconnect\n", timeSinceLastEvent)
			isConnected = false
			sb.gatewayActive = false
		}
		
		if !isConnected {
			sb.gatewayActive = false
			sb.reconnectAttempt++
			
			delay := min(2<<uint(min(sb.reconnectAttempt, 5)), 30)
			
			fmt.Printf("[!] Gateway disconnected, reconnecting in %ds (attempt %d)\n", delay, sb.reconnectAttempt)
			time.Sleep(time.Duration(delay) * time.Second)

			if sb.gateway != nil {
				sb.gateway.Close()
				sb.gateway = nil
			}
			
			if sb.heartbeatTicker != nil {
				sb.heartbeatTicker.Stop()
				sb.heartbeatTicker = nil
			}

			fmt.Println("[+] Attempting to reconnect...")
			if err := sb.connectGateway(); err != nil {
				fmt.Printf("[!] Reconnection failed: %v - will retry in %ds\n", err, checkInterval/time.Second)
			} else {
				fmt.Println("[+] Gateway reconnected successfully")
				sb.lastGatewayEvent = time.Now()
			}
			
			if sb.reconnectAttempt >= 999 {
				sb.reconnectAttempt = 0
				fmt.Println("[!] Reset reconnection counter")
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (sb *Selfbot) connectGateway() error {
	sb.lastGatewayEvent = time.Now()
	
	headers := map[string]string{
		"User-Agent":      sb.client.HeaderSpoofer.Profile.UserAgent,
		"Origin":          "https://discord.com",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "en-US,en;q=0.9",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
	}
	sb.gateway = discord.NewWebSocketClient("wss://gateway.discord.gg/?v=9&encoding=json", headers)

	if err := sb.gateway.Connect(); err != nil {
		return err
	}

	sb.gateway.On("MESSAGE_CREATE", func(data map[string]interface{}) {
		sb.lastGatewayEvent = time.Now()
		sb.handleGatewayMessage(data)
	})

	sb.gateway.On("READY", func(data map[string]interface{}) {
		fmt.Println("[+] Gateway READY - instant responses enabled")
		sb.gatewayActive = true
		sb.reconnectAttempt = 0
		sb.lastGatewayEvent = time.Now()
		
		// Reapply streaming if it was active
		if sb.streaming {
			time.Sleep(2 * time.Second)
			sb.reapplyStreaming()
		}
	})

	sb.gateway.OnClose(func() {
		fmt.Println("[!] Gateway connection closed - will attempt reconnect")
		sb.gatewayActive = false
		sb.streaming = false
	})

	sb.gateway.On("ERROR", func(data map[string]interface{}) {
		fmt.Printf("[!] Gateway error: %v\n", data)
	})

	time.Sleep(500 * time.Millisecond)
	sb.sendIdentify()

	sb.startHeartbeat()

	return nil
}

func (sb *Selfbot) startHeartbeat() {
	if sb.heartbeatTicker != nil {
		sb.heartbeatTicker.Stop()
	}

	sb.heartbeatTicker = time.NewTicker(41250 * time.Millisecond)
	go func() {
		for range sb.heartbeatTicker.C {
			if sb.gateway != nil && sb.gateway.IsConnected() {
				heartbeat := map[string]interface{}{
					"op": 1,
					"d":  nil,
				}
				if err := sb.gateway.SendJSON(heartbeat); err != nil {
					fmt.Printf("[!] Heartbeat failed: %v\n", err)
					sb.gatewayActive = false
				}
			} else {
				sb.gatewayActive = false
			}
		}
	}()
}

func (sb *Selfbot) sendIdentify() {
	identify := map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"token": sb.client.Token,
			"properties": map[string]interface{}{
				"$os":              "linux",
				"$browser":         "Chrome",
				"$device":          "desktop",
				"$referrer":        "",
				"$referring_domain": "",
			},
			"compress":        false,
			"large_threshold": 250,
			"intents":         3276799, // All intents for selfbot
			"presence": map[string]interface{}{
				"status":     "online",
				"since":      0,
				"activities": []interface{}{},
				"afk":        false,
			},
		},
	}
	sb.gateway.SendJSON(identify)
}

func (sb *Selfbot) reapplyStreaming() {
	if !sb.streaming || !sb.gatewayActive {
		return
	}

	activity := map[string]interface{}{
		"name": "Custom Stream",
		"type": 1,
		"url":  "https://twitch.tv/toolers",
	}

	if sb.streamMode == "single" {
		if sb.streamName != "" && sb.streamName != "skip" && sb.streamName != "‎" {
			activity["name"] = sb.streamName
		}
	} else {
		if sb.streamDetails != "" && sb.streamDetails != "skip" && sb.streamDetails != "‎" {
			activity["details"] = sb.streamDetails
		}
		if sb.streamState != "" && sb.streamState != "skip" && sb.streamState != "‎" {
			activity["state"] = sb.streamState
		}
		activity["timestamps"] = map[string]interface{}{
			"start": time.Now().Unix(),
		}
	}

	if sb.streamImageKey != "" && sb.streamImageKey != "skip" && sb.streamImageKey != "‎" {
		assets := map[string]interface{}{
			"large_image": sb.streamImageKey,
		}
		if sb.streamMode == "triple" {
			if sb.streamImageText != "" && sb.streamImageText != "skip" && sb.streamImageText != "‎" {
				assets["large_text"] = sb.streamImageText
			} else {
				assets["large_text"] = "Streaming"
			}
		}
		activity["assets"] = assets
	}

	payload := map[string]interface{}{
		"op": 3,
		"d": map[string]interface{}{
			"since":      time.Now().Unix() * 1000,
			"activities": []map[string]interface{}{activity},
			"status":     "online",
			"afk":        false,
		},
	}

	sb.gateway.SendJSON(payload)
	fmt.Println("[+] Reapplied streaming status after reconnection")
}

func (sb *Selfbot) handleGatewayMessage(data map[string]interface{}) {
	// Handle commands
	author, ok := data["author"].(map[string]interface{})
	if !ok {
		return
	}

	authorID, ok := author["id"].(string)
	if !ok {
		return
	}

	// Check if message is from a react target
	sb.reactMutex.RLock()
	emoji, isTarget := sb.reactTargets[authorID]
	sb.reactMutex.RUnlock()

	if isTarget {
		channelID, _ := data["channel_id"].(string)
		messageID, _ := data["id"].(string)
		if channelID != "" && messageID != "" {
			go sb.addReaction(channelID, messageID, emoji)
		}
	}

	// Handle own commands
	if authorID != sb.client.CurrentUser.ID {
		return
	}

	content, ok := data["content"].(string)
	if !ok || !strings.HasPrefix(content, sb.prefix) {
		return
	}

	channelID, _ := data["channel_id"].(string)
	messageID, _ := data["id"].(string)

	cmdContent := strings.TrimPrefix(content, sb.prefix)
	parts := strings.Fields(cmdContent)

	if len(parts) == 0 {
		return
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	go sb.handleCommand(cmd, args, channelID, messageID)
}

func (sb *Selfbot) addReaction(channelID, messageID, emoji string) {
	sb.client.AddReaction(channelID, messageID, emoji)
}

func (sb *Selfbot) pollMessages() {
	checked := make(map[string]bool)
	failedAttempts := 0

	for {
		time.Sleep(2000 * time.Millisecond)

		if sb.gatewayActive && failedAttempts == 0 {
			continue
		}

		guilds, err := sb.client.GetGuilds(100)
		if err != nil {
			failedAttempts++
			if failedAttempts > 3 {
				fmt.Printf("[!] Polling failed %d times: %v\n", failedAttempts, err)
				time.Sleep(10 * time.Second)
			}
			continue
		}

		failedAttempts = 0

		for _, guild := range guilds {
			channels, err := sb.client.GetGuildChannels(guild.ID)
			if err != nil {
				continue
			}

			for _, channel := range channels {
				if channel.Type != 0 {
					continue
				}

				messages, err := sb.client.GetMessages(channel.ID, 15)
				if err != nil {
					continue
				}

				for _, msg := range messages {
					if checked[msg.ID] {
						continue
					}

					checked[msg.ID] = true

					// Check react targets
					sb.reactMutex.RLock()
					emoji, isTarget := sb.reactTargets[msg.Author.ID]
					sb.reactMutex.RUnlock()

					if isTarget {
						go sb.addReaction(msg.ChannelID, msg.ID, emoji)
					}

					// Check own commands
					if msg.Author.ID == sb.client.CurrentUser.ID && strings.HasPrefix(msg.Content, sb.prefix) {
						cmdContent := strings.TrimPrefix(msg.Content, sb.prefix)
						parts := strings.Fields(cmdContent)
						if len(parts) > 0 {
							go sb.handleCommand(strings.ToLower(parts[0]), parts[1:], msg.ChannelID, msg.ID)
						}
					}
				}
			}
		}

		if len(checked) > 5000 {
			checked = make(map[string]bool)
		}
	}
}

func (sb *Selfbot) handleCommand(cmd string, args []string, channelID, messageID string) {
	fmt.Printf("[+] Processing command: %s %v\n", cmd, args)

	switch cmd {
	case "info":
		sb.cmdInfo(channelID, messageID)
	case "ping":
		sb.cmdPing(channelID, messageID)
	case "purge":
		if len(args) > 0 {
			sb.cmdPurge(channelID, messageID, args[0])
		}
	case "spam":
		if len(args) >= 2 {
			sb.cmdSpam(channelID, messageID, args)
		}
	case "userinfo":
		if len(args) > 0 {
			sb.cmdUserInfo(channelID, messageID, args[0])
		}
	case "host":
		if len(args) > 0 {
			sb.cmdHost(channelID, messageID, args[0])
		}
	case "hosted":
		sb.cmdHosted(channelID, messageID)
	case "stream":
		sb.cmdStream(channelID, messageID, args)
	case "mode":
		sb.cmdMode(channelID, messageID, args)
	case "react":
		sb.cmdReact(channelID, messageID, args)
	case "end":
		sb.cmdEnd(channelID, messageID, args)
	case "massdm":
		sb.cmdMassDM(channelID, messageID, args)
	case "massgc":
		sb.cmdMassGC(channelID, messageID, args)
	case "report":
		sb.cmdReport(channelID, messageID, args)
	}
}

func (sb *Selfbot) sendTemp(channelID, content, originalMsgID string) {
	msg, err := sb.client.SendMessage(channelID, content)
	if err != nil {
		fmt.Printf("[!] Failed to send message: %v\n", err)
		return
	}
	go func() {
		time.Sleep(10 * time.Second)
		sb.client.DeleteMessage(channelID, msg.ID)
		sb.client.DeleteMessage(channelID, originalMsgID)
	}()
}

func (sb *Selfbot) cmdInfo(channelID, messageID string) {
	info := "```asciidoc\n= Selfbot Commands =\n\n" +
		"> info        :: Show commands\n" +
		"> ping        :: Check latency\n" +
		"> purge <n>   :: Delete messages\n" +
		"> spam <n> <msg> :: Spam\n" +
		"> userinfo <id> :: User info\n" +
		"> host <token> :: Host account\n" +
		"> hosted      :: List hosted\n" +
		"> stream      :: Control RPC streaming\n" +
		"> mode <1/3>  :: Switch stream mode\n" +
		"> react @user <emoji> :: Auto-react to user\n" +
		"> end react @user :: Stop reacting\n" +
		"> massdm <1|2|3> <msg> :: Mass DM\n" +
		"> massgc <msg> :: Mass GC spam\n" +
		"> report @user :: Auto-report user (100 msgs)\n```"

	msg, _ := sb.client.SendMessage(channelID, info)
	go func() {
		time.Sleep(10 * time.Second)
		sb.client.DeleteMessage(channelID, msg.ID)
		sb.client.DeleteMessage(channelID, messageID)
	}()
}

func (sb *Selfbot) cmdReact(channelID, messageID string, args []string) {
	if len(args) < 2 {
		sb.sendTemp(channelID, "```diff\n- Usage: >react @user <emoji>\nExample: >react @user 👍\nExample: >react @user <:customemoji:123456>\n```", messageID)
		return
	}

	userMention := args[0]
	emoji := strings.Join(args[1:], " ")

	// Extract user ID from mention
	userID := strings.TrimPrefix(userMention, "<@")
	userID = strings.TrimPrefix(userID, "!")
	userID = strings.TrimSuffix(userID, ">")

	// Validate emoji format
	emojiToUse := emoji
	if strings.HasPrefix(emoji, "<") && strings.HasSuffix(emoji, ">") {
		// Custom emoji format: <:name:id> or <a:name:id>
		emojiToUse = emoji
	} else {
		// Regular unicode emoji
		emojiToUse = emoji
	}

	sb.reactMutex.Lock()
	sb.reactTargets[userID] = emojiToUse
	sb.reactMutex.Unlock()

	user, _ := sb.client.GetUser(userID)
	username := "User"
	if user != nil {
		username = user.Username
	}

	sb.sendTemp(channelID, fmt.Sprintf("```diff\n+ Now reacting to all messages from %s with %s\n```", username, emoji), messageID)
	fmt.Printf("[+] Auto-react enabled for user %s with emoji %s\n", userID, emojiToUse)
}

func (sb *Selfbot) cmdEnd(channelID, messageID string, args []string) {
	if len(args) < 2 || strings.ToLower(args[0]) != "react" {
		sb.sendTemp(channelID, "```diff\n- Usage: >end react @user\n```", messageID)
		return
	}

	userMention := args[1]
	userID := strings.TrimPrefix(userMention, "<@")
	userID = strings.TrimPrefix(userID, "!")
	userID = strings.TrimSuffix(userID, ">")

	sb.reactMutex.Lock()
	_, existed := sb.reactTargets[userID]
	delete(sb.reactTargets, userID)
	sb.reactMutex.Unlock()

	if existed {
		user, _ := sb.client.GetUser(userID)
		username := "User"
		if user != nil {
			username = user.Username
		}
		sb.sendTemp(channelID, fmt.Sprintf("```diff\n+ Stopped reacting to %s\n```", username), messageID)
		fmt.Printf("[+] Auto-react disabled for user %s\n", userID)
	} else {
		sb.sendTemp(channelID, "```diff\n- Not currently reacting to that user\n```", messageID)
	}
}

func (sb *Selfbot) cmdMassDM(channelID, messageID string, args []string) {
	if len(args) < 2 {
		info := "```asciidoc\n= Mass DM Options =\n\n" +
			"1 :: All DM history\n" +
			"2 :: Friends with existing DMs\n" +
			"3 :: Both (history + friends)\n\n" +
			"Usage: >massdm <1|2|3> <message>\n" +
			"Example: >massdm 1 Hello everyone!\n```"
		sb.sendTemp(channelID, info, messageID)
		return
	}

	option := args[0]
	message := strings.Join(args[1:], " ")

	statusMsg, _ := sb.client.SendMessage(channelID, "```asciidoc\n= Mass DM\nStatus :: Initializing...\n```")

	go func() {
		sb.client.DeleteMessage(channelID, messageID)

		targets := make(map[string]struct {
			channelID string
			username  string
			dmType    string
		})

		// Fetch DM history using TLS client with proper headers
		headers := sb.client.HeaderSpoofer.GetHeaders(sb.client.TLSClient, map[string]string{})

		resp, err := sb.client.TLSClient.Get("https://discord.com/api/v9/users/@me/channels", headers)
		if err != nil {
			sb.client.EditMessage(channelID, statusMsg.ID, "```diff\n- Failed to fetch DMs\n```")
			return
		}

		if resp.StatusCode != 200 {
			sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```diff\n- Failed to fetch DMs (status %d)\n```", resp.StatusCode))
			return
		}

		var dms []map[string]interface{}
		json.Unmarshal(resp.Body, &dms)

		// Process DMs
		existingDMChannels := make(map[string]struct {
			channelID string
			username  string
		})

		for _, dm := range dms {
			dmType, _ := dm["type"].(float64)
			if dmType == 1 { // DM channel
				lastMsgID, _ := dm["last_message_id"].(string)
				recipients, _ := dm["recipients"].([]interface{})
				if lastMsgID != "" && len(recipients) > 0 {
					recipient := recipients[0].(map[string]interface{})
					userID, _ := recipient["id"].(string)
					username, _ := recipient["username"].(string)
					dmChannelID, _ := dm["id"].(string)

					if userID != "" {
						existingDMChannels[userID] = struct {
							channelID string
							username  string
						}{dmChannelID, username}

						if option == "1" || option == "3" {
							targets[userID] = struct {
								channelID string
								username  string
								dmType    string
							}{dmChannelID, username, "dm_history"}
						}
					}
				}
			}
		}

		// Fetch friends if option 2 or 3
		if option == "2" || option == "3" {
			friendHeaders := sb.client.HeaderSpoofer.GetHeaders(sb.client.TLSClient, map[string]string{
				"x-discord-timezone": "America/Los_Angeles",
			})
			delete(friendHeaders, "Content-Type")

			friendResp, err := sb.client.TLSClient.Get("https://discord.com/api/v9/users/@me/relationships", friendHeaders)
			if err == nil && friendResp.StatusCode == 200 {
				var relationships []map[string]interface{}
				json.Unmarshal(friendResp.Body, &relationships)

				for _, rel := range relationships {
					relType, _ := rel["type"].(float64)
					if relType == 1 { // Friend
						user, _ := rel["user"].(map[string]interface{})
						userID, _ := user["id"].(string)
						username, _ := user["username"].(string)

						// Check if already has DM
						if dmInfo, exists := existingDMChannels[userID]; exists {
							if _, alreadyAdded := targets[userID]; !alreadyAdded {
								targets[userID] = struct {
									channelID string
									username  string
									dmType    string
								}{dmInfo.channelID, username, "friend_existing_dm"}
							}
						}
					}
				}
			} else if option == "2" {
				sb.client.EditMessage(channelID, statusMsg.ID, "```diff\n- Failed to fetch friends\n```")
				return
			}
		}

		total := len(targets)
		if total == 0 {
			sb.client.EditMessage(channelID, statusMsg.ID, "```diff\n- No targets found (existing DMs only)\n```")
			return
		}

		sent := 0
		failed := 0

		sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```asciidoc\n= Mass DM\nTargets :: %d (safe - existing DMs)\nStatus :: Sending...\nSent :: 0/%d\n```", total, total))

		count := 0
		for _, target := range targets {
			count++
			
			// Send using TLS client with retries
			payload := map[string]interface{}{
				"content": message,
				"nonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
				"tts":     false,
				"flags":   0,
			}

			success := false
			for attempt := 0; attempt < 3; attempt++ {
				msgResp, err := sb.client.TLSClient.Post(
					fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", target.channelID),
					headers,
					payload,
				)

				if err == nil && msgResp.StatusCode == 200 {
					sent++
					success = true
					break
				} else if msgResp != nil && msgResp.StatusCode == 429 {
					var rlData map[string]interface{}
					json.Unmarshal(msgResp.Body, &rlData)
					retryAfter, _ := rlData["retry_after"].(float64)
					time.Sleep(time.Duration((retryAfter+rand.Float64())*1000) * time.Millisecond)
					continue
				} else if msgResp != nil && (msgResp.StatusCode == 400 || msgResp.StatusCode == 403 || msgResp.StatusCode == 401) {
					break
				} else if attempt < 2 {
					time.Sleep(2 * time.Second)
				}
			}

			if !success {
				failed++
			}

			if count%5 == 0 || count == total {
				sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```asciidoc\n= Mass DM\nMode :: %s\nTargets :: %d (safe - existing DMs only)\nStatus :: Sending...\nSent :: %d/%d\nFailed :: %d\nLast :: %s (%s)\n```",
					map[string]string{"1": "History", "2": "Friends", "3": "Both"}[option], total, sent, total, failed, target.username, target.dmType))
			}

			time.Sleep(time.Duration(2500+rand.Intn(1500)) * time.Millisecond)
		}

		sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```asciidoc\n= Mass DM\nMode :: %s\nStatus :: Completed SAFELY\nSent :: %d/%d\nFailed :: %d\nNote :: Only sent to existing DMs (no captcha risk)\nTime :: %s\n```",
			map[string]string{"1": "History", "2": "Friends", "3": "Both"}[option], sent, total, failed, time.Now().Format("03:04 PM MST")))

		time.Sleep(10 * time.Second)
		sb.client.DeleteMessage(channelID, statusMsg.ID)
	}()
}

func (sb *Selfbot) cmdMassGC(channelID, messageID string, args []string) {
	if len(args) == 0 {
		sb.sendTemp(channelID, "```diff\n- Usage: >massgc <message>\n```", messageID)
		return
	}

	message := strings.Join(args, " ")
	statusMsg, _ := sb.client.SendMessage(channelID, "```asciidoc\n= Mass GC\nStatus :: Fetching group chats...\n```")

	go func() {
		sb.client.DeleteMessage(channelID, messageID)

		headers := sb.client.HeaderSpoofer.GetHeaders(sb.client.TLSClient, map[string]string{})

		resp, err := sb.client.TLSClient.Get("https://discord.com/api/v9/users/@me/channels", headers)
		if err != nil || resp.StatusCode != 200 {
			sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```diff\n- Failed to fetch channels (%d)\n```", resp.StatusCode))
			return
		}

		var channels []map[string]interface{}
		json.Unmarshal(resp.Body, &channels)

		var groupChats []map[string]interface{}
		for _, ch := range channels {
			chType, _ := ch["type"].(float64)
			if chType == 3 { // Group chat
				groupChats = append(groupChats, ch)
			}
		}

		total := len(groupChats)
		if total == 0 {
			sb.client.EditMessage(channelID, statusMsg.ID, "```diff\n- No group chats found\n```")
			return
		}

		sent := 0
		failed := 0

		sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```asciidoc\n= Mass GC\nFound :: %d groupchats\nStatus :: Sending...\nProgress :: 0/%d\n```", total, total))

		for _, gc := range groupChats {
			gcID, _ := gc["id"].(string)
			gcName, _ := gc["name"].(string)
			if gcName == "" {
				gcName = "Unnamed Group"
			}

			recipients, _ := gc["recipients"].([]interface{})
			memberCount := len(recipients)
			displayName := fmt.Sprintf("%s (%d members)", gcName, memberCount)

			payload := map[string]interface{}{
				"content": message,
				"nonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
				"tts":     false,
				"flags":   0,
			}

			success := false
			for attempt := 0; attempt < 8; attempt++ {
				msgResp, err := sb.client.TLSClient.Post(
					fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", gcID),
					headers,
					payload,
				)

				if err == nil && msgResp.StatusCode == 200 {
					sent++
					success = true
					fmt.Printf("[SENT] %s\n", displayName)
					break
				} else if msgResp != nil && msgResp.StatusCode == 429 {
					var rlData map[string]interface{}
					json.Unmarshal(msgResp.Body, &rlData)
					retryAfter, _ := rlData["retry_after"].(float64)
					wait := retryAfter + rand.Float64()*3 + 2
					fmt.Printf("[RATELIMIT] %s → sleeping %.2fs\n", displayName, wait)
					time.Sleep(time.Duration(wait*1000) * time.Millisecond)
					continue
				} else if msgResp != nil && (msgResp.StatusCode == 10003 || msgResp.StatusCode == 50001 || msgResp.StatusCode == 50013) {
					fmt.Printf("[KICKED/LEFT] %s\n", displayName)
					failed++
					break
				} else if attempt < 7 {
					time.Sleep(3 * time.Second)
				}
			}

			if !success && failed == 0 {
				failed++
			}

			sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```asciidoc\n= Mass GC\nFound :: %d groupchats\nProgress :: %d/%d\nSent :: %d\nFailed/Kicked :: %d\nCurrent :: %s\n```", total, sent+failed, total, sent, failed, displayName))

			time.Sleep(time.Duration(3500+rand.Intn(3000)) * time.Millisecond)
		}

		sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```asciidoc\n= Mass GC :: DONE\nSuccessfully sent to :: %d/%d\nFailed/Kicked :: %d\nFinished at :: %s\n```", sent, total, failed, time.Now().Format("15:04:05")))

		time.Sleep(10 * time.Second)
		sb.client.DeleteMessage(channelID, statusMsg.ID)
	}()
}

func (sb *Selfbot) cmdMode(channelID, messageID string, args []string) {
	if len(args) == 0 {
		sb.sendTemp(channelID, fmt.Sprintf("```diff\n+ Current stream mode: %s\nUse '>mode 1' for single text (name only)\nUse '>mode 3' for triple text (details, state, small_text)\n```", sb.streamMode), messageID)
		return
	}

	mode := args[0]
	if mode == "1" {
		sb.streamMode = "single"
		fmt.Println("[+] Set to single text mode (name only)")
		sb.sendTemp(channelID, "```diff\n+ Stream mode set to SINGLE (name only)\nDetails, state and small_text will be omitted\n```", messageID)
	} else if mode == "3" {
		sb.streamMode = "triple"
		fmt.Println("[+] Set to triple text mode (details, state, small_text)")
		sb.sendTemp(channelID, "```diff\n+ Stream mode set to TRIPLE (details, state, small_text)\n```", messageID)
	} else {
		sb.sendTemp(channelID, "```diff\n- Invalid mode. Use '1' or '3'\n```", messageID)
	}
}

func (sb *Selfbot) cmdReport(channelID, messageID string, args []string) {
	if len(args) < 1 {
		sb.sendTemp(channelID, "```diff\n- Usage: >report @user\nExample: >report @username\nThis will scan and report their last 100 messages\n```", messageID)
		return
	}

	userMention := args[0]
	userID := strings.TrimPrefix(userMention, "<@")
	userID = strings.TrimPrefix(userID, "!")
	userID = strings.TrimSuffix(userID, ">")

	statusMsg, _ := sb.client.SendMessage(channelID, "```asciidoc\n= Mass Report\nStatus :: Scanning messages...\n```")

	go func() {
		sb.client.DeleteMessage(channelID, messageID)

		// Get user info
		user, err := sb.client.GetUser(userID)
		if err != nil {
			sb.client.EditMessage(channelID, statusMsg.ID, "```diff\n- Failed to fetch user info\n```")
			return
		}

		username := user.Username

		// Fetch messages from current channel
		messages, err := sb.client.GetMessages(channelID, 100)
		if err != nil {
			sb.client.EditMessage(channelID, statusMsg.ID, "```diff\n- Failed to fetch messages\n```")
			return
		}

		// Filter messages from target user
		var targetMessages []string
		for _, msg := range messages {
			if msg.Author.ID == userID {
				targetMessages = append(targetMessages, msg.ID)
			}
		}

		if len(targetMessages) == 0 {
			sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```diff\n- No messages found from %s in this channel\n```", username))
			time.Sleep(5 * time.Second)
			sb.client.DeleteMessage(channelID, statusMsg.ID)
			return
		}

		sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```asciidoc\n= Mass Report\nTarget :: %s\nFound :: %d messages\nStatus :: Reporting...\nProgress :: 0/%d\n```", username, len(targetMessages), len(targetMessages)))

		// Report each message
		reported := 0
		failed := 0

		headers := sb.client.HeaderSpoofer.GetHeaders(sb.client.TLSClient, map[string]string{})

		reasons := []string{
			"Illegal content involving minors",
			"Harassment or bullying",
			"Spam or phishing links",
			"Self-harm or suicide content",
			"Violent threats or hate speech",
		}

		for i, msgID := range targetMessages {
			reason := reasons[rand.Intn(len(reasons))]

			payload := map[string]interface{}{
				"channel_id": channelID,
				"message_id": msgID,
				"reason":     reason,
				"guild_id":   nil,
			}

			success := false
			for attempt := 0; attempt < 3; attempt++ {
				reportResp, err := sb.client.TLSClient.Post(
					"https://discord.com/api/v9/report",
					headers,
					payload,
				)

				if err == nil && (reportResp.StatusCode == 200 || reportResp.StatusCode == 201) {
					reported++
					success = true
					break
				} else if reportResp != nil && reportResp.StatusCode == 429 {
					var rlData map[string]interface{}
					json.Unmarshal(reportResp.Body, &rlData)
					retryAfter, _ := rlData["retry_after"].(float64)
					time.Sleep(time.Duration((retryAfter+1)*1000) * time.Millisecond)
					continue
				} else if attempt < 2 {
					time.Sleep(2 * time.Second)
				}
			}

			if !success {
				failed++
			}

			if (i+1)%10 == 0 || i+1 == len(targetMessages) {
				sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```asciidoc\n= Mass Report\nTarget :: %s\nFound :: %d messages\nStatus :: Reporting...\nProgress :: %d/%d\nReported :: %d\nFailed :: %d\n```", username, len(targetMessages), i+1, len(targetMessages), reported, failed))
			}

			time.Sleep(time.Duration(1000+rand.Intn(2000)) * time.Millisecond)
		}

		sb.client.EditMessage(channelID, statusMsg.ID, fmt.Sprintf("```diff\n+ Mass Report Complete\n+ Target: %s\n+ Reported: %d/%d messages\n- Failed: %d\n```", username, reported, len(targetMessages), failed))

		time.Sleep(10 * time.Second)
		sb.client.DeleteMessage(channelID, statusMsg.ID)
	}()
}

func (sb *Selfbot) cmdPing(channelID, messageID string) {
	start := time.Now()
	msg, _ := sb.client.SendMessage(channelID, "```diff\n+ Pong!\n```")
	latency := time.Since(start).Milliseconds()
	sb.client.EditMessage(channelID, msg.ID, fmt.Sprintf("```diff\n+ Pong! %dms\n```", latency))

	go func() {
		time.Sleep(10 * time.Second)
		sb.client.DeleteMessage(channelID, msg.ID)
		sb.client.DeleteMessage(channelID, messageID)
	}()
}

func (sb *Selfbot) cmdPurge(channelID, messageID, amountStr string) {
	var amount int
	fmt.Sscanf(amountStr, "%d", &amount)

	if amount < 1 || amount > 100 {
		return
	}

	messages, _ := sb.client.GetMessages(channelID, 100)

	deleted := 0
	for _, msg := range messages {
		if msg.Author.ID == sb.client.CurrentUser.ID && deleted < amount {
			sb.client.DeleteMessage(channelID, msg.ID)
			deleted++
			time.Sleep(300 * time.Millisecond)
		}
	}

	sb.sendTemp(channelID, fmt.Sprintf("```diff\n+ Deleted %d\n```", deleted), messageID)
}

func (sb *Selfbot) cmdSpam(channelID, messageID string, args []string) {
	var count int
	fmt.Sscanf(args[0], "%d", &count)

	if count < 1 || count > 20 {
		return
	}

	message := strings.Join(args[1:], " ")
	sb.client.DeleteMessage(channelID, messageID)

	for i := 0; i < count; i++ {
		sb.client.SendMessage(channelID, message)
		time.Sleep(500 * time.Millisecond)
	}
}

func (sb *Selfbot) cmdUserInfo(channelID, messageID, userID string) {
	user, err := sb.client.GetUser(userID)
	if err != nil {
		return
	}

	info := fmt.Sprintf("```asciidoc\n= User =\n\nName :: %s#%s\nID   :: %s\nBot  :: %t\n```",
		user.Username, user.Discriminator, user.ID, user.Bot)

	msg, _ := sb.client.SendMessage(channelID, info)
	go func() {
		time.Sleep(10 * time.Second)
		sb.client.DeleteMessage(channelID, msg.ID)
		sb.client.DeleteMessage(channelID, messageID)
	}()
}

func (sb *Selfbot) cmdHost(channelID, messageID, token string) {
	newClient := discord.NewDiscordClient(token)
	user, err := newClient.GetCurrentUser()
	if err != nil {
		sb.sendTemp(channelID, "```diff\n- Invalid token\n```", messageID)
		return
	}

	if _, exists := sb.hostedBots[user.ID]; exists {
		sb.sendTemp(channelID, "```diff\n- Already hosting\n```", messageID)
		return
	}

	newClient.CurrentUser = user
	sb.hostedBots[user.ID] = newClient

	fmt.Printf("[+] Now hosting %s#%s (ID: %s)\n", user.Username, user.Discriminator, user.ID)

	sb.sendTemp(channelID, fmt.Sprintf("```diff\n+ Hosting %s#%s\n```", user.Username, user.Discriminator), messageID)

	go sb.pollHostedBot(newClient, user.ID)
}

func (sb *Selfbot) cmdHosted(channelID, messageID string) {
	if len(sb.hostedBots) == 0 {
		sb.sendTemp(channelID, "```diff\n- No hosted accounts\n```", messageID)
		return
	}

	list := "```asciidoc\n= Hosted =\n\n"
	for _, bot := range sb.hostedBots {
		if bot.CurrentUser != nil {
			list += fmt.Sprintf("+ %s#%s\n", bot.CurrentUser.Username, bot.CurrentUser.Discriminator)
		}
	}
	list += "```"

	msg, _ := sb.client.SendMessage(channelID, list)
	go func() {
		time.Sleep(10 * time.Second)
		sb.client.DeleteMessage(channelID, msg.ID)
		sb.client.DeleteMessage(channelID, messageID)
	}()
}

func (sb *Selfbot) cmdStream(channelID, messageID string, args []string) {
	if len(args) == 0 {
		info := "```asciidoc\n= Stream Command =\n\n" +
			"> stream on/off        :: Start/stop streaming\n" +
			"> stream set <fields>  :: Set stream fields\n" +
			"> stream image <url>   :: Set stream image\n" +
			"> mode <1/3>           :: Switch stream mode\n" +
			"\nCurrent Mode: " + sb.streamMode + "\n"

		if sb.streamMode == "single" {
			info += "\nSingle Mode Example:\n>stream set name \"Playing game\"\n>stream image https://cdn.discordapp.com/attachments/...\n>stream on\n"
		} else {
			info += "\nTriple Mode Example:\n>stream set details \"Playing game\" state \"Online\" small_text \"Level 99\"\n>stream image https://cdn.discordapp.com/attachments/...\n>stream on\n"
		}

		info += "\nUse 'skip' to omit a field\n```"

		msg, _ := sb.client.SendMessage(channelID, info)
		go func() {
			time.Sleep(10 * time.Second)
			sb.client.DeleteMessage(channelID, msg.ID)
			sb.client.DeleteMessage(channelID, messageID)
		}()
		return
	}

	subCmd := strings.ToLower(args[0])

	switch subCmd {
	case "on":
		sb.startStreaming(channelID, messageID)
	case "off":
		sb.stopStreaming(channelID, messageID)
	case "set":
		sb.setStreamDetails(channelID, messageID, args[1:])
	case "image":
		sb.setStreamImage(channelID, messageID, args[1:])
	default:
		sb.sendTemp(channelID, "```diff\n- Invalid stream command\n```", messageID)
	}
}

func (sb *Selfbot) startStreaming(channelID, messageID string) {
	if sb.streaming {
		sb.sendTemp(channelID, "```diff\n- Already streaming\n```", messageID)
		return
	}

	if !sb.gatewayActive {
		sb.sendTemp(channelID, "```diff\n- Gateway not connected\n```", messageID)
		return
	}

	activity := map[string]interface{}{
		"name": "Custom Stream",
		"type": 1,
		"url":  "https://twitch.tv/toolers",
	}

	if sb.streamMode == "single" {
		if sb.streamName != "" && sb.streamName != "skip" && sb.streamName != "‎" {
			activity["name"] = sb.streamName
		}
	} else {
		if sb.streamDetails != "" && sb.streamDetails != "skip" && sb.streamDetails != "‎" {
			activity["details"] = sb.streamDetails
		}

		if sb.streamState != "" && sb.streamState != "skip" && sb.streamState != "‎" {
			activity["state"] = sb.streamState
		}
		activity["timestamps"] = map[string]interface{}{
			"start": time.Now().Unix(),
		}
	}

	if sb.streamImageKey != "" && sb.streamImageKey != "skip" && sb.streamImageKey != "‎" {
		assets := map[string]interface{}{
			"large_image": sb.streamImageKey,
		}
		if sb.streamMode == "triple" {
			if sb.streamImageText != "" && sb.streamImageText != "skip" && sb.streamImageText != "‎" {
				assets["large_text"] = sb.streamImageText
			} else {
				assets["large_text"] = "Streaming"
			}
		}
		activity["assets"] = assets
	}

	payload := map[string]interface{}{
		"op": 3,
		"d": map[string]interface{}{
			"since":      time.Now().Unix() * 1000,
			"activities": []map[string]interface{}{activity},
			"status":     "online",
			"afk":        false,
		},
	}

	sb.gateway.SendJSON(payload)
	sb.streaming = true

	statusMsg := "```diff\n+ Streaming started!\n"
	if sb.streamMode == "single" {
		if sb.streamName != "" && sb.streamName != "skip" && sb.streamName != "‎" {
			statusMsg += fmt.Sprintf("Name: %s\n", sb.streamName)
		}
	} else {
		if sb.streamDetails != "" && sb.streamDetails != "skip" && sb.streamDetails != "‎" {
			statusMsg += fmt.Sprintf("Details: %s\n", sb.streamDetails)
		}
		if sb.streamState != "" && sb.streamState != "skip" && sb.streamState != "‎" {
			statusMsg += fmt.Sprintf("State: %s\n", sb.streamState)
		}
		if sb.streamImageText != "" && sb.streamImageText != "skip" && sb.streamImageText != "‎" {
			statusMsg += fmt.Sprintf("Small Text: %s\n", sb.streamImageText)
		}
	}
	if sb.streamImageKey != "" && sb.streamImageKey != "skip" && sb.streamImageKey != "‎" {
		statusMsg += "Image: Set\n"
	}
	statusMsg += "```"

	sb.sendTemp(channelID, statusMsg, messageID)
	fmt.Printf("[+] Started streaming in %s mode\n", sb.streamMode)
}

func (sb *Selfbot) stopStreaming(channelID, messageID string) {
	if !sb.streaming {
		sb.sendTemp(channelID, "```diff\n- Not streaming\n```", messageID)
		return
	}

	if !sb.gatewayActive {
		sb.sendTemp(channelID, "```diff\n- Gateway not connected\n```", messageID)
		return
	}

	payload := map[string]interface{}{
		"op": 3,
		"d": map[string]interface{}{
			"since":      nil,
			"activities": []map[string]interface{}{},
			"status":     "online",
			"afk":        false,
		},
	}

	sb.gateway.SendJSON(payload)
	sb.streaming = false
	sb.sendTemp(channelID, "```diff\n+ Streaming stopped\n```", messageID)
	fmt.Println("[-] Stopped streaming")
}

func (sb *Selfbot) setStreamDetails(channelID, messageID string, args []string) {
	var name, details, state, smallText string

	if sb.streamMode == "single" {
		for i := 0; i < len(args); i++ {
			if strings.ToLower(args[i]) == "name" && i+1 < len(args) {
				i++
				if args[i] == "skip" {
					name = "‎"
				} else {
					name = args[i]
				}
			}
		}
	} else {
		for i := 0; i < len(args); i++ {
			switch strings.ToLower(args[i]) {
			case "details":
				if i+1 < len(args) {
					i++
					if args[i] == "skip" {
						details = "‎"
					} else {
						details = args[i]
					}
				}
			case "state":
				if i+1 < len(args) {
					i++
					if args[i] == "skip" {
						state = "‎"
					} else {
						state = args[i]
					}
				}
			case "small_text":
				if i+1 < len(args) {
					i++
					if args[i] == "skip" {
						smallText = "‎"
					} else {
						smallText = args[i]
					}
				}
			}
		}
	}

	if name == "" && details == "" && state == "" && smallText == "" {
		sb.sendTemp(channelID, "```diff\n- No changes specified\n```", messageID)
		return
	}

	if name != "" {
		sb.streamName = name
	}
	if details != "" {
		sb.streamDetails = details
	}
	if state != "" {
		sb.streamState = state
	}
	if smallText != "" {
		sb.streamImageText = smallText
	}

	statusMsg := "```diff\n+ Stream settings updated\n"
	if sb.streamMode == "single" {
		if name != "" && name != "‎" {
			statusMsg += "Name: " + sb.streamName + "\n"
		}
		statusMsg += "Mode: Single (name only)\n"
	} else {
		if details != "" && details != "‎" {
			statusMsg += "Details: " + sb.streamDetails + "\n"
		}
		if state != "" && state != "‎" {
			statusMsg += "State: " + sb.streamState + "\n"
		}
		if smallText != "" && smallText != "‎" {
			statusMsg += "Small Text: " + sb.streamImageText + "\n"
		}
		statusMsg += "Mode: Triple\n"
	}
	statusMsg += "Use '>stream on' to start streaming\n```"

	sb.sendTemp(channelID, statusMsg, messageID)
}

func getAssetKeyFromCDN(url string) string {
	re := regexp.MustCompile(`https?://(?:cdn\.discordapp\.com|media\.discordapp\.net)/attachments/(\d+)/(\d+)/(.+)`)
	matches := re.FindStringSubmatch(url)
	if matches != nil {
		channelID := matches[1]
		attachmentID := matches[2]
		filename := matches[3]
		return fmt.Sprintf("mp:attachments/%s/%s/%s", channelID, attachmentID, filename)
	}
	return ""
}

func (sb *Selfbot) uploadAndGetAssetKey(imageUrl string) (string, error) {
	assetKey := getAssetKeyFromCDN(imageUrl)
	if assetKey != "" {
		return assetKey, nil
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to download image")
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("URL is not an image")
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	filename := filepath.Base(imageUrl)
	if strings.Contains(filename, "?") {
		filename = strings.Split(filename, "?")[0]
	}

	ext, err := mime.ExtensionsByType(contentType)
	if err != nil || len(ext) == 0 {
		ext = []string{".png"}
	}

	if !strings.Contains(filename, ".") || len(filename) > 50 {
		if strings.Contains(contentType, "gif") {
			filename = "asset.gif"
		} else {
			filename = "asset" + ext[0]
		}
	}

	dmChannel, err := sb.client.CreateDM(sb.client.CurrentUser.ID)
	if err != nil {
		return "", err
	}

	msg, err := sb.client.SendMessage(dmChannel.ID, "Uploading image for stream...")
	if err != nil {
		return "", err
	}

	time.Sleep(2 * time.Second)
	sb.client.DeleteMessage(dmChannel.ID, msg.ID)

	headers := map[string]string{
		"Authorization": sb.client.Token,
		"Content-Type":  contentType,
		"User-Agent":    sb.client.HeaderSpoofer.Profile.UserAgent,
	}

	url := fmt.Sprintf("%s/channels/%s/messages", discord.APIBase, dmChannel.ID)

	req, err := http.NewRequest("POST", url, bytes.NewReader(imageBytes))
	if err != nil {
		return "", err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp2, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		return "", fmt.Errorf("failed to upload image to Discord")
	}

	body, err := io.ReadAll(resp2.Body)
	if err != nil {
		return "", err
	}

	var messageData map[string]interface{}
	if err := json.Unmarshal(body, &messageData); err != nil {
		return "", err
	}

	attachments, ok := messageData["attachments"].([]interface{})
	if !ok || len(attachments) == 0 {
		return "", fmt.Errorf("no attachments in response")
	}

	attachment := attachments[0].(map[string]interface{})
	attachmentURL, ok := attachment["url"].(string)
	if !ok {
		return "", fmt.Errorf("no URL in attachment")
	}

	assetKey = getAssetKeyFromCDN(attachmentURL)
	if assetKey == "" {
		return "", fmt.Errorf("failed to extract asset key from uploaded image")
	}

	return assetKey, nil
}

func (sb *Selfbot) setStreamImage(channelID, messageID string, args []string) {
	if len(args) == 0 {
		sb.sendTemp(channelID, "```diff\n- Provide image URL or 'attachment'\n```", messageID)
		return
	}

	if args[0] == "attachment" {
		sb.sendTemp(channelID, "```diff\n+ Reply to an image with '>stream image attachment'\n```", messageID)
		return
	}

	url := args[0]

	if args[0] == "skip" {
		sb.streamImageKey = ""
		sb.sendTemp(channelID, "```diff\n+ Image cleared\n```", messageID)
		return
	}

	if !strings.HasPrefix(url, "http") {
		sb.sendTemp(channelID, "```diff\n- Invalid URL\n```", messageID)
		return
	}

	assetKey, err := sb.uploadAndGetAssetKey(url)
	if err != nil {
		sb.sendTemp(channelID, fmt.Sprintf("```diff\n- Failed to set image: %v\n```", err), messageID)
		return
	}

	sb.streamImageKey = assetKey
	sb.sendTemp(channelID, fmt.Sprintf("```diff\n+ Image set successfully!\nAsset Key: %s\n```", assetKey), messageID)
}

func (sb *Selfbot) pollHostedBot(client *discord.DiscordClient, userID string) {
	checked := make(map[string]bool)

	fmt.Printf("[*] Hosted bot polling started for %s\n", client.CurrentUser.Username)

	for {
		if _, exists := sb.hostedBots[userID]; !exists {
			fmt.Printf("[-] Stopped hosting %s\n", client.CurrentUser.Username)
			return
		}

		time.Sleep(2000 * time.Millisecond)

		guilds, err := client.GetGuilds(100)
		if err != nil {
			continue
		}

		for _, guild := range guilds {
			channels, err := client.GetGuildChannels(guild.ID)
			if err != nil {
				continue
			}

			for _, channel := range channels {
				if channel.Type != 0 {
					continue
				}

				messages, err := client.GetMessages(channel.ID, 15)
				if err != nil {
					continue
				}

				for _, msg := range messages {
					if checked[msg.ID] {
						continue
					}

					checked[msg.ID] = true

					if msg.Author.ID == client.CurrentUser.ID {
						if strings.HasPrefix(msg.Content, sb.prefix) {
							cmdContent := strings.TrimPrefix(msg.Content, sb.prefix)
							parts := strings.Fields(cmdContent)
							if len(parts) > 0 {
								cmd := strings.ToLower(parts[0])
								args := parts[1:]

								// Hosted bots have all commands except host
								go sb.handleHostedCommand(client, cmd, args, msg.ChannelID, msg.ID)
							}
						}
					}
				}
			}
		}

		if len(checked) > 5000 {
			checked = make(map[string]bool)
		}
	}
}

func (sb *Selfbot) handleHostedCommand(client *discord.DiscordClient, cmd string, args []string, channelID, messageID string) {
	// Hosted bots can use ALL commands except host
	if cmd == "host" {
		msg, _ := client.SendMessage(channelID, "```diff\n- Hosted accounts cannot host\n```")
		go func() {
			time.Sleep(10 * time.Second)
			client.DeleteMessage(channelID, msg.ID)
			client.DeleteMessage(channelID, messageID)
		}()
		return
	}

	// Use main command handler for everything else
	// Create temporary selfbot instance for hosted bot
	hostedSB := &Selfbot{
		client:       client,
		prefix:       sb.prefix,
		hostedBots:   sb.hostedBots,   // Share hosted bots map
		reactTargets: sb.reactTargets, // Share react targets
		streamMode:   sb.streamMode,
	}

	hostedSB.handleCommand(cmd, args, channelID, messageID)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=== Discord Selfbot ===")

	bot := NewSelfbot(TOKEN)

	if err := bot.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}