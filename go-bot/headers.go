package discord

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

// BrowserProfile represents a consistent browser fingerprint
type BrowserProfile struct {
	UserAgent           string
	OS                  string
	Browser             string
	BrowserVersion      string
	OSVersion           string
	Locale              string
	Timezone            string
	ScreenResolution    string
	HardwareConcurrency int
	DeviceMemory        int
	Fonts               []string
	Plugins             []string
}

// HeaderSpoofer handles advanced header spoofing and fingerprinting
type HeaderSpoofer struct {
	Token          string
	Fingerprint    string
	Cookies        string
	CacheTime      int64
	SessionID      string
	BrowserSession string
	Profile        *BrowserProfile
	BuildNumber    int
	XSPHash        string
}

// Location represents a timezone/locale pair
type Location struct {
	Timezone string
	Locale   string
}

// ChromeVersion represents Chrome version info
type ChromeVersion struct {
	Major string
	Full  string
}

// NewHeaderSpoofer creates a new header spoofing instance
func NewHeaderSpoofer(token string) *HeaderSpoofer {
	hs := &HeaderSpoofer{
		Token:          token,
		SessionID:      fmt.Sprintf("%d", time.Now().UnixNano()/1000000),
		BrowserSession: fmt.Sprintf("session_%d", randomInt(100000, 999999)),
		BuildNumber:    284054,
	}

	hs.Profile = hs.createConsistentProfile()
	hs.XSPHash = hs.generateXSPHash()

	return hs
}

// createConsistentProfile creates a consistent browser profile
func (hs *HeaderSpoofer) createConsistentProfile() *BrowserProfile {
	timestamp := time.Now().Unix()

	locations := []Location{
		{"America/New_York", "en-US"},
		{"America/Chicago", "en-US"},
		{"America/Los_Angeles", "en-US"},
		{"Europe/London", "en-GB"},
		{"Europe/Paris", "fr-FR"},
		{"Asia/Tokyo", "ja-JP"},
		{"Australia/Sydney", "en-AU"},
	}

	location := locations[timestamp%int64(len(locations))]

	chromeVersions := []ChromeVersion{
		{"125", "125.0.6422.113"},
		{"124", "124.0.6367.207"},
		{"123", "123.0.6312.122"},
		{"126", "126.0.6478.126"},
		{"127", "127.0.6533.88"},
	}

	chrome := chromeVersions[(timestamp/3600)%int64(len(chromeVersions))]

	osVersions := map[string][]string{
		"Windows": {"10", "11"},
		"Mac":     {"13_6", "14_4"},
		"Linux":   {"x86_64"},
	}

	osType := "Windows"
	osVersion := osVersions[osType][timestamp%int64(len(osVersions[osType]))]

	resolutions := []string{"1920x1080", "2560x1440", "3840x2160", "1366x768", "1536x864"}
	resolution := resolutions[(timestamp/1000)%int64(len(resolutions))]

	hwConcurrency := []int{8, 12, 16}
	deviceMemory := []int{8, 16, 32}

	return &BrowserProfile{
		UserAgent:           fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", chrome.Full),
		OS:                  osType,
		Browser:             "Chrome",
		BrowserVersion:      chrome.Full,
		OSVersion:           osVersion,
		Locale:              location.Locale,
		Timezone:            location.Timezone,
		ScreenResolution:    resolution,
		HardwareConcurrency: hwConcurrency[timestamp%int64(len(hwConcurrency))],
		DeviceMemory:        deviceMemory[timestamp%int64(len(deviceMemory))],
		Fonts: []string{
			"Arial", "Helvetica", "Times New Roman", "Verdana",
			"Georgia", "Courier New", "Comic Sans MS", "Trebuchet MS",
		},
		Plugins: []string{
			"Chrome PDF Plugin",
			"Chrome PDF Viewer",
			"Native Client",
			"Widevine Content Decryption Module",
		},
	}
}

// generateXSPHash generates X-Super-Properties hash
func (hs *HeaderSpoofer) generateXSPHash() string {
	data := fmt.Sprintf("%s%s%s%d",
		hs.Profile.OS,
		hs.Profile.Browser,
		hs.Profile.Locale,
		hs.BuildNumber,
	)

	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])[:8]
}

// FetchFingerprint fetches Discord fingerprint
func (hs *HeaderSpoofer) FetchFingerprint(client *TLSClient) (string, string) {
	// Use cache if valid
	if time.Now().Unix()-hs.CacheTime < 3600 && hs.Fingerprint != "" {
		return hs.Fingerprint, hs.Cookies
	}

	headers := map[string]string{
		"User-Agent":      hs.Profile.UserAgent,
		"Accept":          "application/json",
		"Accept-Language": hs.Profile.Locale,
	}

	resp, err := client.Get("https://discord.com/api/v9/experiments", headers)
	if err == nil && resp.StatusCode == 200 {
		var data map[string]interface{}
		if err := json.Unmarshal(resp.Body, &data); err == nil {
			if fp, ok := data["fingerprint"].(string); ok {
				hs.Fingerprint = fp
			} else {
				hs.Fingerprint = hs.fallbackFingerprint()
			}

			// Build cookies from response
			cookieParts := []string{}
			for _, cookie := range resp.Cookies {
				cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
			}
			cookieParts = append(cookieParts, fmt.Sprintf("locale=%s", hs.Profile.Locale))
			hs.Cookies = joinStrings(cookieParts, "; ")
			hs.CacheTime = time.Now().Unix()
		} else {
			hs.Fingerprint = hs.fallbackFingerprint()
			hs.Cookies = hs.defaultCookies()
		}
	} else {
		hs.Fingerprint = hs.fallbackFingerprint()
		hs.Cookies = hs.defaultCookies()
	}

	return hs.Fingerprint, hs.Cookies
}

// fallbackFingerprint generates fallback fingerprint
func (hs *HeaderSpoofer) fallbackFingerprint() string {
	base := time.Now().UnixNano() / 1000000
	randNum := randomInt64(100000000000000000, 999999999999999999)
	return fmt.Sprintf("%d.%d", base, randNum)
}

// defaultCookies generates default cookies
func (hs *HeaderSpoofer) defaultCookies() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("__dcfduid=%dabcdef; __sdcfduid=%dghijkl; locale=%s",
		timestamp, timestamp, hs.Profile.Locale)
}

// GenerateSuperProperties generates X-Super-Properties header
func (hs *HeaderSpoofer) GenerateSuperProperties() string {
	props := map[string]interface{}{
		"os":                  hs.Profile.OS,
		"browser":             hs.Profile.Browser,
		"device":              "",
		"system_locale":       hs.Profile.Locale,
		"browser_user_agent":  hs.Profile.UserAgent,
		"browser_version":     hs.Profile.BrowserVersion,
		"os_version":          hs.Profile.OSVersion,
		"referrer":            "",
		"referring_domain":    "",
		"release_channel":     "stable",
		"client_build_number": hs.BuildNumber,
		"client_event_source": nil,
		"design_id":           0,
	}

	jsonData, _ := json.Marshal(props)
	return base64.StdEncoding.EncodeToString(jsonData)
}

// GenerateSecChUA generates Sec-CH-UA header
func (hs *HeaderSpoofer) GenerateSecChUA() string {
	majorVersion := hs.Profile.BrowserVersion[:3]
	return fmt.Sprintf(`"Chromium";v="%s", "Google Chrome";v="%s", "Not=A?Brand";v="99"`,
		majorVersion, majorVersion)
}

// GetHeaders returns complete headers with fingerprinting
func (hs *HeaderSpoofer) GetHeaders(client *TLSClient, additional map[string]string) map[string]string {
	hs.RotateIfNeeded()

	fingerprint, cookies := hs.FetchFingerprint(client)

	headers := map[string]string{
		"Authorization":             hs.Token,
		"User-Agent":                hs.Profile.UserAgent,
		"Content-Type":              "application/json",
		"Accept":                    "*/*",
		"Accept-Language":           fmt.Sprintf("%s,en;q=0.9", hs.Profile.Locale),
		"Accept-Encoding":           "gzip, deflate, br",
		"Origin":                    "https://discord.com",
		"Referer":                   "https://discord.com/channels/@me",
		"Sec-Ch-Ua":                 hs.GenerateSecChUA(),
		"Sec-Ch-Ua-Mobile":          "?0",
		"Sec-Ch-Ua-Platform":        fmt.Sprintf(`"%s"`, hs.Profile.OS),
		"Sec-Fetch-Dest":            "empty",
		"Sec-Fetch-Mode":            "cors",
		"Sec-Fetch-Site":            "same-origin",
		"Dnt":                       "1",
		"Upgrade-Insecure-Requests": "1",
		"X-Debug-Options":           "bugReporterEnabled",
		"X-Discord-Locale":          hs.Profile.Locale,
		"X-Discord-Timezone":        hs.Profile.Timezone,
		"X-Super-Properties":        hs.GenerateSuperProperties(),
		"X-Fingerprint":             fingerprint,
		"Cookie":                    cookies,
		"X-Track":                   generateTrackHash(),
		"X-Super-Properties-Hash":   hs.XSPHash,
	}

	// Merge additional headers
	for k, v := range additional {
		headers[k] = v
	}

	return headers
}

// GetWebSocketHeaders returns WebSocket-specific headers
func (hs *HeaderSpoofer) GetWebSocketHeaders() map[string]string {
	hs.RotateIfNeeded()

	wsKey := make([]byte, 16)
	rand.Read(wsKey)

	return map[string]string{
		"User-Agent":               hs.Profile.UserAgent,
		"Accept-Encoding":          "gzip, deflate, br",
		"Accept-Language":          hs.Profile.Locale,
		"Cache-Control":            "no-cache",
		"Pragma":                   "no-cache",
		"Sec-WebSocket-Extensions": "permessage-deflate; client_max_window_bits",
		"Sec-WebSocket-Key":        base64.StdEncoding.EncodeToString(wsKey),
		"Sec-WebSocket-Version":    "13",
		"Upgrade":                  "websocket",
		"Connection":               "Upgrade",
		"Origin":                   "https://discord.com",
		"Sec-WebSocket-Protocol":   "json",
	}
}

// RotateProfile rotates to a new browser profile
func (hs *HeaderSpoofer) RotateProfile() *BrowserProfile {
	hs.Profile = hs.createConsistentProfile()
	hs.XSPHash = hs.generateXSPHash()
	hs.CacheTime = 0
	return hs.Profile
}

// RotateIfNeeded auto-rotates profile every 6 hours
func (hs *HeaderSpoofer) RotateIfNeeded() bool {
	if time.Now().Unix()-hs.CacheTime > 21600 { // 6 hours
		hs.RotateProfile()
		return true
	}
	return false
}

// Helper functions

func randomInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return int(n.Int64()) + min
}

func randomInt64(min, max int64) int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return n.Int64() + min
}

func generateTrackHash() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return hex.EncodeToString(hash[:])
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
