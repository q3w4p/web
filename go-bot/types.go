package discord

// User represents a Discord user account
type User struct {
	ID            string  `json:"id"`
	Username      string  `json:"username"`
	Discriminator string  `json:"discriminator"`
	Avatar        *string `json:"avatar"`
	Bot           bool    `json:"bot,omitempty"`
	System        bool    `json:"system,omitempty"`
	MFAEnabled    bool    `json:"mfa_enabled,omitempty"`
	Banner        *string `json:"banner"`
	AccentColor   *int    `json:"accent_color"`
	Locale        string  `json:"locale,omitempty"`
	Verified      bool    `json:"verified,omitempty"`
	Email         *string `json:"email"`
	Flags         int     `json:"flags,omitempty"`
	PremiumType   int     `json:"premium_type,omitempty"`
	PublicFlags   int     `json:"public_flags,omitempty"`
	Bio           string  `json:"bio,omitempty"`
}

// Guild represents a Discord server/guild
type Guild struct {
	ID                          string                `json:"id"`
	Name                        string                `json:"name"`
	Icon                        *string               `json:"icon"`
	IconHash                    *string               `json:"icon_hash"`
	Splash                      *string               `json:"splash"`
	DiscoverySplash             *string               `json:"discovery_splash"`
	Owner                       bool                  `json:"owner,omitempty"`
	OwnerID                     string                `json:"owner_id"`
	Permissions                 *string               `json:"permissions"`
	Region                      string                `json:"region"`
	AFKChannelID                *string               `json:"afk_channel_id"`
	AFKTimeout                  int                   `json:"afk_timeout"`
	WidgetEnabled               bool                  `json:"widget_enabled,omitempty"`
	WidgetChannelID             *string               `json:"widget_channel_id"`
	VerificationLevel           int                   `json:"verification_level"`
	DefaultMessageNotifications int                   `json:"default_message_notifications"`
	ExplicitContentFilter       int                   `json:"explicit_content_filter"`
	Roles                       []Role                `json:"roles"`
	Emojis                      []Emoji               `json:"emojis"`
	Features                    []string              `json:"features"`
	MFALevel                    int                   `json:"mfa_level"`
	ApplicationID               *string               `json:"application_id"`
	SystemChannelID             *string               `json:"system_channel_id"`
	SystemChannelFlags          int                   `json:"system_channel_flags"`
	RulesChannelID              *string               `json:"rules_channel_id"`
	JoinedAt                    *string               `json:"joined_at"`
	Large                       bool                  `json:"large,omitempty"`
	Unavailable                 bool                  `json:"unavailable,omitempty"`
	MemberCount                 int                   `json:"member_count,omitempty"`
	VoiceStates                 []VoiceState          `json:"voice_states,omitempty"`
	Members                     []GuildMember         `json:"members,omitempty"`
	Channels                    []Channel             `json:"channels,omitempty"`
	Threads                     []Channel             `json:"threads,omitempty"`
	Presences                   []PresenceUpdate      `json:"presences,omitempty"`
	MaxPresences                int                   `json:"max_presences,omitempty"`
	MaxMembers                  int                   `json:"max_members,omitempty"`
	VanityURLCode               *string               `json:"vanity_url_code"`
	Description                 *string               `json:"description"`
	Banner                      *string               `json:"banner"`
	PremiumTier                 int                   `json:"premium_tier"`
	PremiumSubscriptionCount    int                   `json:"premium_subscription_count,omitempty"`
	PreferredLocale             string                `json:"preferred_locale"`
	PublicUpdatesChannelID      *string               `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                   `json:"max_video_channel_users,omitempty"`
	ApproximateMemberCount      int                   `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    int                   `json:"approximate_presence_count,omitempty"`
	WelcomeScreen               *WelcomeScreen        `json:"welcome_screen"`
	NSFWLevel                   int                   `json:"nsfw_level"`
	StageInstances              []StageInstance       `json:"stage_instances,omitempty"`
	Stickers                    []Sticker             `json:"stickers,omitempty"`
	GuildScheduledEvents        []GuildScheduledEvent `json:"guild_scheduled_events,omitempty"`
	PremiumProgressBarEnabled   bool                  `json:"premium_progress_bar_enabled"`
}

// Channel represents any Discord channel type
type Channel struct {
	ID                   string                `json:"id"`
	Type                 int                   `json:"type"`
	GuildID              *string               `json:"guild_id"`
	Position             int                   `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`
	Name                 *string               `json:"name"`
	Topic                *string               `json:"topic"`
	NSFW                 bool                  `json:"nsfw,omitempty"`
	LastMessageID        *string               `json:"last_message_id"`
	Bitrate              int                   `json:"bitrate,omitempty"`
	UserLimit            int                   `json:"user_limit,omitempty"`
	RateLimitPerUser     int                   `json:"rate_limit_per_user,omitempty"`
	Recipients           []User                `json:"recipients,omitempty"`
	Icon                 *string               `json:"icon"`
	OwnerID              *string               `json:"owner_id"`
	ApplicationID        *string               `json:"application_id"`
	ParentID             *string               `json:"parent_id"`
	LastPinTimestamp     *string               `json:"last_pin_timestamp"`
	RTCRegion            *string               `json:"rtc_region"`
	VideoQualityMode     int                   `json:"video_quality_mode,omitempty"`
	MessageCount         int                   `json:"message_count,omitempty"`
	MemberCount          int                   `json:"member_count,omitempty"`
	ThreadMetadata       *ThreadMetadata       `json:"thread_metadata"`
	Member               *ThreadMember         `json:"member"`
	DefaultAutoArchive   int                   `json:"default_auto_archive_duration,omitempty"`
	Permissions          *string               `json:"permissions"`
	Flags                int                   `json:"flags,omitempty"`
}

// Message represents a Discord message
type Message struct {
	ID                string              `json:"id"`
	ChannelID         string              `json:"channel_id"`
	GuildID           *string             `json:"guild_id"`
	Author            User                `json:"author"`
	Member            *GuildMember        `json:"member"`
	Content           string              `json:"content"`
	Timestamp         string              `json:"timestamp"`
	EditedTimestamp   *string             `json:"edited_timestamp"`
	TTS               bool                `json:"tts"`
	MentionEveryone   bool                `json:"mention_everyone"`
	Mentions          []User              `json:"mentions"`
	MentionRoles      []string            `json:"mention_roles"`
	MentionChannels   []ChannelMention    `json:"mention_channels,omitempty"`
	Attachments       []Attachment        `json:"attachments"`
	Embeds            []Embed             `json:"embeds"`
	Reactions         []Reaction          `json:"reactions,omitempty"`
	Nonce             interface{}         `json:"nonce,omitempty"`
	Pinned            bool                `json:"pinned"`
	WebhookID         *string             `json:"webhook_id"`
	Type              int                 `json:"type"`
	Activity          *MessageActivity    `json:"activity"`
	Application       *Application        `json:"application"`
	ApplicationID     *string             `json:"application_id"`
	MessageReference  *MessageReference   `json:"message_reference"`
	Flags             int                 `json:"flags,omitempty"`
	ReferencedMessage *Message            `json:"referenced_message"`
	Interaction       *MessageInteraction `json:"interaction"`
	Thread            *Channel            `json:"thread"`
	Components        []Component         `json:"components,omitempty"`
	StickerItems      []StickerItem       `json:"sticker_items,omitempty"`
	Stickers          []Sticker           `json:"stickers,omitempty"`
	Position          int                 `json:"position,omitempty"`
}

// Role represents a Discord guild role
type Role struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Color        int       `json:"color"`
	Hoist        bool      `json:"hoist"`
	Icon         *string   `json:"icon"`
	UnicodeEmoji *string   `json:"unicode_emoji"`
	Position     int       `json:"position"`
	Permissions  string    `json:"permissions"`
	Managed      bool      `json:"managed"`
	Mentionable  bool      `json:"mentionable"`
	Tags         *RoleTags `json:"tags"`
}

// RoleTags contains additional role metadata
type RoleTags struct {
	BotID             *string `json:"bot_id"`
	IntegrationID     *string `json:"integration_id"`
	PremiumSubscriber bool    `json:"premium_subscriber,omitempty"`
}

// GuildMember represents a user's membership in a guild
type GuildMember struct {
	User                       *User    `json:"user"`
	Nick                       *string  `json:"nick"`
	Avatar                     *string  `json:"avatar"`
	Roles                      []string `json:"roles"`
	JoinedAt                   string   `json:"joined_at"`
	PremiumSince               *string  `json:"premium_since"`
	Deaf                       bool     `json:"deaf"`
	Mute                       bool     `json:"mute"`
	Pending                    bool     `json:"pending,omitempty"`
	Permissions                *string  `json:"permissions"`
	CommunicationDisabledUntil *string  `json:"communication_disabled_until"`
}

// Emoji represents a Discord emoji
type Emoji struct {
	ID            *string  `json:"id"`
	Name          *string  `json:"name"`
	Roles         []string `json:"roles,omitempty"`
	User          *User    `json:"user,omitempty"`
	RequireColons bool     `json:"require_colons,omitempty"`
	Managed       bool     `json:"managed,omitempty"`
	Animated      bool     `json:"animated,omitempty"`
	Available     bool     `json:"available,omitempty"`
}

// Sticker represents a Discord sticker
type Sticker struct {
	ID          string  `json:"id"`
	PackID      *string `json:"pack_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Tags        string  `json:"tags"`
	Asset       string  `json:"asset,omitempty"`
	Type        int     `json:"type"`
	FormatType  int     `json:"format_type"`
	Available   bool    `json:"available,omitempty"`
	GuildID     *string `json:"guild_id"`
	User        *User   `json:"user,omitempty"`
	SortValue   int     `json:"sort_value,omitempty"`
}

// StickerItem represents a partial sticker object
type StickerItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FormatType int    `json:"format_type"`
}

// Attachment represents a message attachment
type Attachment struct {
	ID          string  `json:"id"`
	Filename    string  `json:"filename"`
	Description *string `json:"description"`
	ContentType *string `json:"content_type"`
	Size        int     `json:"size"`
	URL         string  `json:"url"`
	ProxyURL    string  `json:"proxy_url"`
	Height      *int    `json:"height"`
	Width       *int    `json:"width"`
	Ephemeral   bool    `json:"ephemeral,omitempty"`
}

// Embed represents a message embed
type Embed struct {
	Title       *string         `json:"title"`
	Type        *string         `json:"type"`
	Description *string         `json:"description"`
	URL         *string         `json:"url"`
	Timestamp   *string         `json:"timestamp"`
	Color       *int            `json:"color"`
	Footer      *EmbedFooter    `json:"footer"`
	Image       *EmbedImage     `json:"image"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail"`
	Video       *EmbedVideo     `json:"video"`
	Provider    *EmbedProvider  `json:"provider"`
	Author      *EmbedAuthor    `json:"author"`
	Fields      []EmbedField    `json:"fields"`
}

// EmbedFooter represents embed footer
type EmbedFooter struct {
	Text         string  `json:"text"`
	IconURL      *string `json:"icon_url"`
	ProxyIconURL *string `json:"proxy_icon_url"`
}

// EmbedImage represents embed image
type EmbedImage struct {
	URL      string  `json:"url"`
	ProxyURL *string `json:"proxy_url"`
	Height   *int    `json:"height"`
	Width    *int    `json:"width"`
}

// EmbedThumbnail represents embed thumbnail
type EmbedThumbnail struct {
	URL      string  `json:"url"`
	ProxyURL *string `json:"proxy_url"`
	Height   *int    `json:"height"`
	Width    *int    `json:"width"`
}

// EmbedVideo represents embed video
type EmbedVideo struct {
	URL      *string `json:"url"`
	ProxyURL *string `json:"proxy_url"`
	Height   *int    `json:"height"`
	Width    *int    `json:"width"`
}

// EmbedProvider represents embed provider
type EmbedProvider struct {
	Name *string `json:"name"`
	URL  *string `json:"url"`
}

// EmbedAuthor represents embed author
type EmbedAuthor struct {
	Name         string  `json:"name"`
	URL          *string `json:"url"`
	IconURL      *string `json:"icon_url"`
	ProxyIconURL *string `json:"proxy_icon_url"`
}

// EmbedField represents an embed field
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// Reaction represents a message reaction
type Reaction struct {
	Count int   `json:"count"`
	Me    bool  `json:"me"`
	Emoji Emoji `json:"emoji"`
}

// ChannelMention represents a channel mentioned in a message
type ChannelMention struct {
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Type    int    `json:"type"`
	Name    string `json:"name"`
}

// MessageActivity represents message activity data
type MessageActivity struct {
	Type    int     `json:"type"`
	PartyID *string `json:"party_id"`
}

// MessageReference represents a reference to another message
type MessageReference struct {
	MessageID       *string `json:"message_id"`
	ChannelID       *string `json:"channel_id"`
	GuildID         *string `json:"guild_id"`
	FailIfNotExists bool    `json:"fail_if_not_exists,omitempty"`
}

// MessageInteraction represents the interaction that triggered the message
type MessageInteraction struct {
	ID     string       `json:"id"`
	Type   int          `json:"type"`
	Name   string       `json:"name"`
	User   User         `json:"user"`
	Member *GuildMember `json:"member"`
}

// Component represents a message component
type Component struct {
	Type        int            `json:"type"`
	CustomID    *string        `json:"custom_id"`
	Disabled    bool           `json:"disabled,omitempty"`
	Style       int            `json:"style,omitempty"`
	Label       *string        `json:"label"`
	Emoji       *Emoji         `json:"emoji"`
	URL         *string        `json:"url"`
	Options     []SelectOption `json:"options,omitempty"`
	Placeholder *string        `json:"placeholder"`
	MinValues   int            `json:"min_values,omitempty"`
	MaxValues   int            `json:"max_values,omitempty"`
	Components  []Component    `json:"components,omitempty"`
}

// SelectOption represents a select menu option
type SelectOption struct {
	Label       string  `json:"label"`
	Value       string  `json:"value"`
	Description *string `json:"description"`
	Emoji       *Emoji  `json:"emoji"`
	Default     bool    `json:"default,omitempty"`
}

// Application represents a Discord application
type Application struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Icon                *string  `json:"icon"`
	Description         string   `json:"description"`
	RPCOrigins          []string `json:"rpc_origins,omitempty"`
	BotPublic           bool     `json:"bot_public"`
	BotRequireCodeGrant bool     `json:"bot_require_code_grant"`
	TermsOfServiceURL   *string  `json:"terms_of_service_url"`
	PrivacyPolicyURL    *string  `json:"privacy_policy_url"`
	Owner               *User    `json:"owner"`
	Summary             string   `json:"summary"`
	VerifyKey           string   `json:"verify_key"`
	Team                *Team    `json:"team"`
	GuildID             *string  `json:"guild_id"`
	PrimarySKUID        *string  `json:"primary_sku_id"`
	Slug                *string  `json:"slug"`
	CoverImage          *string  `json:"cover_image"`
	Flags               int      `json:"flags,omitempty"`
}

// Team represents an application team
type Team struct {
	Icon        *string      `json:"icon"`
	ID          string       `json:"id"`
	Members     []TeamMember `json:"members"`
	Name        string       `json:"name"`
	OwnerUserID string       `json:"owner_user_id"`
}

// TeamMember represents a team member
type TeamMember struct {
	MembershipState int      `json:"membership_state"`
	Permissions     []string `json:"permissions"`
	TeamID          string   `json:"team_id"`
	User            User     `json:"user"`
}

// PermissionOverwrite represents a channel permission overwrite
type PermissionOverwrite struct {
	ID    string `json:"id"`
	Type  int    `json:"type"`
	Allow string `json:"allow"`
	Deny  string `json:"deny"`
}

// ThreadMetadata represents thread metadata
type ThreadMetadata struct {
	Archived            bool    `json:"archived"`
	AutoArchiveDuration int     `json:"auto_archive_duration"`
	ArchiveTimestamp    string  `json:"archive_timestamp"`
	Locked              bool    `json:"locked"`
	Invitable           bool    `json:"invitable,omitempty"`
	CreateTimestamp     *string `json:"create_timestamp"`
}

// ThreadMember represents a thread member
type ThreadMember struct {
	ID            *string `json:"id"`
	UserID        *string `json:"user_id"`
	JoinTimestamp string  `json:"join_timestamp"`
	Flags         int     `json:"flags"`
}

// VoiceState represents a user's voice connection state
type VoiceState struct {
	GuildID                 *string      `json:"guild_id"`
	ChannelID               *string      `json:"channel_id"`
	UserID                  string       `json:"user_id"`
	Member                  *GuildMember `json:"member"`
	SessionID               string       `json:"session_id"`
	Deaf                    bool         `json:"deaf"`
	Mute                    bool         `json:"mute"`
	SelfDeaf                bool         `json:"self_deaf"`
	SelfMute                bool         `json:"self_mute"`
	SelfStream              bool         `json:"self_stream,omitempty"`
	SelfVideo               bool         `json:"self_video"`
	Suppress                bool         `json:"suppress"`
	RequestToSpeakTimestamp *string      `json:"request_to_speak_timestamp"`
}

// PresenceUpdate represents a user's presence
type PresenceUpdate struct {
	User         User         `json:"user"`
	GuildID      string       `json:"guild_id"`
	Status       string       `json:"status"`
	Activities   []Activity   `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
}

// Activity represents a user activity
type Activity struct {
	Name          string              `json:"name"`
	Type          int                 `json:"type"`
	URL           *string             `json:"url"`
	CreatedAt     int64               `json:"created_at"`
	Timestamps    *ActivityTimestamps `json:"timestamps"`
	ApplicationID *string             `json:"application_id"`
	Details       *string             `json:"details"`
	State         *string             `json:"state"`
	Emoji         *Emoji              `json:"emoji"`
	Party         *ActivityParty      `json:"party"`
	Assets        *ActivityAssets     `json:"assets"`
	Secrets       *ActivitySecrets    `json:"secrets"`
	Instance      bool                `json:"instance,omitempty"`
	Flags         int                 `json:"flags,omitempty"`
	Buttons       []ActivityButton    `json:"buttons,omitempty"`
}

// ActivityTimestamps represents activity timestamps
type ActivityTimestamps struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

// ActivityParty represents activity party info
type ActivityParty struct {
	ID   *string `json:"id"`
	Size []int   `json:"size,omitempty"`
}

// ActivityAssets represents activity assets
type ActivityAssets struct {
	LargeImage *string `json:"large_image"`
	LargeText  *string `json:"large_text"`
	SmallImage *string `json:"small_image"`
	SmallText  *string `json:"small_text"`
}

// ActivitySecrets represents activity secrets
type ActivitySecrets struct {
	Join     *string `json:"join"`
	Spectate *string `json:"spectate"`
	Match    *string `json:"match"`
}

// ActivityButton represents an activity button
type ActivityButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// ClientStatus represents client platform status
type ClientStatus struct {
	Desktop *string `json:"desktop"`
	Mobile  *string `json:"mobile"`
	Web     *string `json:"web"`
}

// WelcomeScreen represents a guild welcome screen
type WelcomeScreen struct {
	Description     *string                `json:"description"`
	WelcomeChannels []WelcomeScreenChannel `json:"welcome_channels"`
}

// WelcomeScreenChannel represents a welcome screen channel
type WelcomeScreenChannel struct {
	ChannelID   string  `json:"channel_id"`
	Description string  `json:"description"`
	EmojiID     *string `json:"emoji_id"`
	EmojiName   *string `json:"emoji_name"`
}

// StageInstance represents a stage channel instance
type StageInstance struct {
	ID                    string  `json:"id"`
	GuildID               string  `json:"guild_id"`
	ChannelID             string  `json:"channel_id"`
	Topic                 string  `json:"topic"`
	PrivacyLevel          int     `json:"privacy_level"`
	DiscoverableDisabled  bool    `json:"discoverable_disabled"`
	GuildScheduledEventID *string `json:"guild_scheduled_event_id"`
}

// GuildScheduledEvent represents a scheduled guild event
type GuildScheduledEvent struct {
	ID                 string                       `json:"id"`
	GuildID            string                       `json:"guild_id"`
	ChannelID          *string                      `json:"channel_id"`
	CreatorID          *string                      `json:"creator_id"`
	Name               string                       `json:"name"`
	Description        *string                      `json:"description"`
	ScheduledStartTime string                       `json:"scheduled_start_time"`
	ScheduledEndTime   *string                      `json:"scheduled_end_time"`
	PrivacyLevel       int                          `json:"privacy_level"`
	Status             int                          `json:"status"`
	EntityType         int                          `json:"entity_type"`
	EntityID           *string                      `json:"entity_id"`
	EntityMetadata     *GuildScheduledEventMetadata `json:"entity_metadata"`
	Creator            *User                        `json:"creator"`
	UserCount          int                          `json:"user_count,omitempty"`
	Image              *string                      `json:"image"`
}

// GuildScheduledEventMetadata represents event metadata
type GuildScheduledEventMetadata struct {
	Location *string `json:"location"`
}

// Invite represents a Discord invite
type Invite struct {
	Code                     string               `json:"code"`
	Guild                    *Guild               `json:"guild"`
	Channel                  *Channel             `json:"channel"`
	Inviter                  *User                `json:"inviter"`
	TargetType               int                  `json:"target_type,omitempty"`
	TargetUser               *User                `json:"target_user"`
	TargetApplication        *Application         `json:"target_application"`
	ApproximatePresenceCount int                  `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   int                  `json:"approximate_member_count,omitempty"`
	ExpiresAt                *string              `json:"expires_at"`
	GuildScheduledEvent      *GuildScheduledEvent `json:"guild_scheduled_event"`
}

// Integration represents a guild integration
type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           bool               `json:"enabled"`
	Syncing           bool               `json:"syncing,omitempty"`
	RoleID            *string            `json:"role_id"`
	EnableEmoticons   bool               `json:"enable_emoticons,omitempty"`
	ExpireBehavior    int                `json:"expire_behavior,omitempty"`
	ExpireGracePeriod int                `json:"expire_grace_period,omitempty"`
	User              *User              `json:"user"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          *string            `json:"synced_at"`
	SubscriberCount   int                `json:"subscriber_count,omitempty"`
	Revoked           bool               `json:"revoked,omitempty"`
	Application       *Application       `json:"application"`
}

// IntegrationAccount represents an integration account
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Ban represents a guild ban
type Ban struct {
	Reason *string `json:"reason"`
	User   User    `json:"user"`
}

// Webhook represents a Discord webhook
type Webhook struct {
	ID            string   `json:"id"`
	Type          int      `json:"type"`
	GuildID       *string  `json:"guild_id"`
	ChannelID     *string  `json:"channel_id"`
	User          *User    `json:"user"`
	Name          *string  `json:"name"`
	Avatar        *string  `json:"avatar"`
	Token         *string  `json:"token"`
	ApplicationID *string  `json:"application_id"`
	SourceGuild   *Guild   `json:"source_guild"`
	SourceChannel *Channel `json:"source_channel"`
	URL           *string  `json:"url"`
}

// AuditLog represents a guild audit log
type AuditLog struct {
	AuditLogEntries      []AuditLogEntry       `json:"audit_log_entries"`
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`
	Integrations         []Integration         `json:"integrations"`
	Threads              []Channel             `json:"threads"`
	Users                []User                `json:"users"`
	Webhooks             []Webhook             `json:"webhooks"`
}

// AuditLogEntry represents an audit log entry
type AuditLogEntry struct {
	TargetID   *string          `json:"target_id"`
	Changes    []AuditLogChange `json:"changes,omitempty"`
	UserID     *string          `json:"user_id"`
	ID         string           `json:"id"`
	ActionType int              `json:"action_type"`
	Options    *AuditEntryInfo  `json:"options"`
	Reason     *string          `json:"reason"`
}

// AuditLogChange represents a change in an audit log entry
type AuditLogChange struct {
	NewValue interface{} `json:"new_value,omitempty"`
	OldValue interface{} `json:"old_value,omitempty"`
	Key      string      `json:"key"`
}

// AuditEntryInfo represents optional audit entry information
type AuditEntryInfo struct {
	DeleteMemberDays string  `json:"delete_member_days,omitempty"`
	MembersRemoved   string  `json:"members_removed,omitempty"`
	ChannelID        *string `json:"channel_id"`
	MessageID        *string `json:"message_id"`
	Count            string  `json:"count,omitempty"`
	ID               *string `json:"id"`
	Type             *string `json:"type"`
	RoleName         *string `json:"role_name"`
}

// VoiceRegion represents a voice region
type VoiceRegion struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Optimal    bool   `json:"optimal"`
	Deprecated bool   `json:"deprecated"`
	Custom     bool   `json:"custom"`
}

// GuildPreview represents a guild preview
type GuildPreview struct {
	ID                       string    `json:"id"`
	Name                     string    `json:"name"`
	Icon                     *string   `json:"icon"`
	Splash                   *string   `json:"splash"`
	DiscoverySplash          *string   `json:"discovery_splash"`
	Emojis                   []Emoji   `json:"emojis"`
	Features                 []string  `json:"features"`
	ApproximateMemberCount   int       `json:"approximate_member_count"`
	ApproximatePresenceCount int       `json:"approximate_presence_count"`
	Description              *string   `json:"description"`
	Stickers                 []Sticker `json:"stickers"`
}

// GuildWidget represents a guild widget
type GuildWidget struct {
	Enabled   bool    `json:"enabled"`
	ChannelID *string `json:"channel_id"`
}

// GuildWidgetSettings represents guild widget settings
type GuildWidgetSettings struct {
	Enabled   bool    `json:"enabled"`
	ChannelID *string `json:"channel_id"`
}

// Connection represents a user connection
type Connection struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Revoked      bool          `json:"revoked,omitempty"`
	Integrations []Integration `json:"integrations,omitempty"`
	Verified     bool          `json:"verified"`
	FriendSync   bool          `json:"friend_sync"`
	ShowActivity bool          `json:"show_activity"`
	Visibility   int           `json:"visibility"`
}

// GuildTemplate represents a guild template
type GuildTemplate struct {
	Code                  string  `json:"code"`
	Name                  string  `json:"name"`
	Description           *string `json:"description"`
	UsageCount            int     `json:"usage_count"`
	CreatorID             string  `json:"creator_id"`
	Creator               User    `json:"creator"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
	SourceGuildID         string  `json:"source_guild_id"`
	SerializedSourceGuild Guild   `json:"serialized_source_guild"`
	IsDirty               *bool   `json:"is_dirty"`
}
