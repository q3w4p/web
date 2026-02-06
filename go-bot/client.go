package discord

import (
	"encoding/json"
	"fmt"
)

const (
	APIBase = "https://discord.com/api/v10"
	CDNBase = "https://cdn.discordapp.com"
)

// DiscordClient is the main client for Discord API
type DiscordClient struct {
	Token         string
	TLSClient     *TLSClient
	HeaderSpoofer *HeaderSpoofer
	Gateway       *WebSocketClient

	// Cache
	CurrentUser *User
	Guilds      map[string]*Guild
	Channels    map[string]*Channel
}

// NewDiscordClient creates a new Discord client
func NewDiscordClient(token string) *DiscordClient {
	tlsClient := NewTLSClient()
	headerSpoofer := NewHeaderSpoofer(token)

	return &DiscordClient{
		Token:         token,
		TLSClient:     tlsClient,
		HeaderSpoofer: headerSpoofer,
		Guilds:        make(map[string]*Guild),
		Channels:      make(map[string]*Channel),
	}
}

// Request makes an API request
func (c *DiscordClient) Request(method, endpoint string, body interface{}) (*HTTPResponse, error) {
	url := APIBase + endpoint
	headers := c.HeaderSpoofer.GetHeaders(c.TLSClient, nil)

	switch method {
	case "GET":
		return c.TLSClient.Get(url, headers)
	case "POST":
		return c.TLSClient.Post(url, headers, body)
	case "PATCH":
		return c.TLSClient.Patch(url, headers, body)
	case "PUT":
		return c.TLSClient.Put(url, headers, body)
	case "DELETE":
		return c.TLSClient.Delete(url, headers)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}

// GetCurrentUser gets the current user
func (c *DiscordClient) GetCurrentUser() (*User, error) {
	resp, err := c.Request("GET", "/users/@me", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	var user User
	if err := json.Unmarshal(resp.Body, &user); err != nil {
		return nil, err
	}

	c.CurrentUser = &user
	return &user, nil
}

// ModifyCurrentUser modifies the current user
func (c *DiscordClient) ModifyCurrentUser(username *string, avatar *string, bio *string) (*User, error) {
	payload := make(map[string]interface{})
	if username != nil {
		payload["username"] = *username
	}
	if avatar != nil {
		payload["avatar"] = *avatar
	}
	if bio != nil {
		payload["bio"] = *bio
	}

	resp, err := c.Request("PATCH", "/users/@me", payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	var user User
	if err := json.Unmarshal(resp.Body, &user); err != nil {
		return nil, err
	}

	c.CurrentUser = &user
	return &user, nil
}

// GetUser gets a user by ID
func (c *DiscordClient) GetUser(userID string) (*User, error) {
	resp, err := c.Request("GET", "/users/"+userID, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var user User
	if err := json.Unmarshal(resp.Body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetGuilds gets current user's guilds
func (c *DiscordClient) GetGuilds(limit int) ([]Guild, error) {
	endpoint := fmt.Sprintf("/users/@me/guilds?limit=%d", limit)
	resp, err := c.Request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var guilds []Guild
	if err := json.Unmarshal(resp.Body, &guilds); err != nil {
		return nil, err
	}

	return guilds, nil
}

// GetGuild gets a guild by ID
func (c *DiscordClient) GetGuild(guildID string) (*Guild, error) {
	resp, err := c.Request("GET", "/guilds/"+guildID, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var guild Guild
	if err := json.Unmarshal(resp.Body, &guild); err != nil {
		return nil, err
	}

	c.Guilds[guildID] = &guild
	return &guild, nil
}

// GetChannel gets a channel by ID
func (c *DiscordClient) GetChannel(channelID string) (*Channel, error) {
	resp, err := c.Request("GET", "/channels/"+channelID, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var channel Channel
	if err := json.Unmarshal(resp.Body, &channel); err != nil {
		return nil, err
	}

	c.Channels[channelID] = &channel
	return &channel, nil
}

// GetMessages gets messages from a channel
func (c *DiscordClient) GetMessages(channelID string, limit int) ([]Message, error) {
	endpoint := fmt.Sprintf("/channels/%s/messages?limit=%d", channelID, limit)
	resp, err := c.Request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var messages []Message
	if err := json.Unmarshal(resp.Body, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// SendMessage sends a message to a channel
func (c *DiscordClient) SendMessage(channelID, content string) (*Message, error) {
	payload := map[string]interface{}{
		"content": content,
	}

	resp, err := c.Request("POST", "/channels/"+channelID+"/messages", payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	var message Message
	if err := json.Unmarshal(resp.Body, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// EditMessage edits a message
func (c *DiscordClient) EditMessage(channelID, messageID, content string) (*Message, error) {
	payload := map[string]interface{}{
		"content": content,
	}

	resp, err := c.Request("PATCH", fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID), payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var message Message
	if err := json.Unmarshal(resp.Body, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// DeleteMessage deletes a message
func (c *DiscordClient) DeleteMessage(channelID, messageID string) error {
	resp, err := c.Request("DELETE", fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// BulkDeleteMessages deletes multiple messages
func (c *DiscordClient) BulkDeleteMessages(channelID string, messageIDs []string) error {
	payload := map[string]interface{}{
		"messages": messageIDs,
	}

	resp, err := c.Request("POST", fmt.Sprintf("/channels/%s/messages/bulk-delete", channelID), payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// AddReaction adds a reaction to a message
func (c *DiscordClient) AddReaction(channelID, messageID, emoji string) error {
	resp, err := c.Request("PUT", fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/@me", channelID, messageID, emoji), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// RemoveReaction removes a reaction
func (c *DiscordClient) RemoveReaction(channelID, messageID, emoji string) error {
	resp, err := c.Request("DELETE", fmt.Sprintf("/channels/%s/messages/%s/reactions/%s/@me", channelID, messageID, emoji), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// CreateDM creates a DM channel
func (c *DiscordClient) CreateDM(recipientID string) (*Channel, error) {
	payload := map[string]interface{}{
		"recipient_id": recipientID,
	}

	resp, err := c.Request("POST", "/users/@me/channels", payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var channel Channel
	if err := json.Unmarshal(resp.Body, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

// TriggerTyping triggers typing indicator
func (c *DiscordClient) TriggerTyping(channelID string) error {
	resp, err := c.Request("POST", fmt.Sprintf("/channels/%s/typing", channelID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// GetGuildMember gets a guild member
func (c *DiscordClient) GetGuildMember(guildID, userID string) (*GuildMember, error) {
	resp, err := c.Request("GET", fmt.Sprintf("/guilds/%s/members/%s", guildID, userID), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var member GuildMember
	if err := json.Unmarshal(resp.Body, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// ModifyGuildMember modifies a guild member
func (c *DiscordClient) ModifyGuildMember(guildID, userID string, nick *string, roles *[]string, mute, deaf *bool) (*GuildMember, error) {
	payload := make(map[string]interface{})
	if nick != nil {
		payload["nick"] = *nick
	}
	if roles != nil {
		payload["roles"] = *roles
	}
	if mute != nil {
		payload["mute"] = *mute
	}
	if deaf != nil {
		payload["deaf"] = *deaf
	}

	resp, err := c.Request("PATCH", fmt.Sprintf("/guilds/%s/members/%s", guildID, userID), payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var member GuildMember
	if err := json.Unmarshal(resp.Body, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// CreateGuildRole creates a guild role
func (c *DiscordClient) CreateGuildRole(guildID, name string, color int, permissions string) (*Role, error) {
	payload := map[string]interface{}{
		"name":        name,
		"color":       color,
		"permissions": permissions,
	}

	resp, err := c.Request("POST", fmt.Sprintf("/guilds/%s/roles", guildID), payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var role Role
	if err := json.Unmarshal(resp.Body, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// DeleteGuildRole deletes a guild role
func (c *DiscordClient) DeleteGuildRole(guildID, roleID string) error {
	resp, err := c.Request("DELETE", fmt.Sprintf("/guilds/%s/roles/%s", guildID, roleID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// CreateGuildBan bans a user from a guild
func (c *DiscordClient) CreateGuildBan(guildID, userID string, deleteMessageDays int, reason string) error {
	payload := map[string]interface{}{
		"delete_message_days": deleteMessageDays,
	}

	headers := c.HeaderSpoofer.GetHeaders(c.TLSClient, map[string]string{
		"X-Audit-Log-Reason": reason,
	})

	url := fmt.Sprintf("%s/guilds/%s/bans/%s", APIBase, guildID, userID)
	resp, err := c.TLSClient.Put(url, headers, payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// RemoveGuildBan removes a ban
func (c *DiscordClient) RemoveGuildBan(guildID, userID string) error {
	resp, err := c.Request("DELETE", fmt.Sprintf("/guilds/%s/bans/%s", guildID, userID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// KickGuildMember kicks a member from a guild
func (c *DiscordClient) KickGuildMember(guildID, userID string) error {
	resp, err := c.Request("DELETE", fmt.Sprintf("/guilds/%s/members/%s", guildID, userID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// GetGuildChannels gets all channels in a guild
func (c *DiscordClient) GetGuildChannels(guildID string) ([]Channel, error) {
	resp, err := c.Request("GET", fmt.Sprintf("/guilds/%s/channels", guildID), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var channels []Channel
	if err := json.Unmarshal(resp.Body, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}

// CreateGuildChannel creates a channel in a guild
func (c *DiscordClient) CreateGuildChannel(guildID, name string, channelType int) (*Channel, error) {
	payload := map[string]interface{}{
		"name": name,
		"type": channelType,
	}

	resp, err := c.Request("POST", fmt.Sprintf("/guilds/%s/channels", guildID), payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var channel Channel
	if err := json.Unmarshal(resp.Body, &channel); err != nil {
		return nil, err
	}

	return &channel, nil
}

// DeleteChannel deletes a channel
func (c *DiscordClient) DeleteChannel(channelID string) error {
	resp, err := c.Request("DELETE", fmt.Sprintf("/channels/%s", channelID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// GetGuildRoles gets all roles in a guild
func (c *DiscordClient) GetGuildRoles(guildID string) ([]Role, error) {
	resp, err := c.Request("GET", fmt.Sprintf("/guilds/%s/roles", guildID), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var roles []Role
	if err := json.Unmarshal(resp.Body, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// AddGuildMemberRole adds a role to a member
func (c *DiscordClient) AddGuildMemberRole(guildID, userID, roleID string) error {
	resp, err := c.Request("PUT", fmt.Sprintf("/guilds/%s/members/%s/roles/%s", guildID, userID, roleID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// RemoveGuildMemberRole removes a role from a member
func (c *DiscordClient) RemoveGuildMemberRole(guildID, userID, roleID string) error {
	resp, err := c.Request("DELETE", fmt.Sprintf("/guilds/%s/members/%s/roles/%s", guildID, userID, roleID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// LeaveGuild leaves a guild
func (c *DiscordClient) LeaveGuild(guildID string) error {
	resp, err := c.Request("DELETE", fmt.Sprintf("/users/@me/guilds/%s", guildID), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// Close closes the client
func (c *DiscordClient) Close() {
	if c.Gateway != nil {
		c.Gateway.Close()
	}
	if c.TLSClient != nil {
		c.TLSClient.Close()
	}
}
