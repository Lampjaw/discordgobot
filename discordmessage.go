package discordgobot

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// DiscordMessage holds received message information
type DiscordMessage struct {
	Discord          *DiscordClient
	DiscordgoMessage *discordgo.Message
	MessageType      MessageType
	Nick             *string
	Content          *string
}

// Channel returns the message channel id
func (m *DiscordMessage) Channel() string {
	return m.DiscordgoMessage.ChannelID
}

// UserName returns the message username
func (m *DiscordMessage) UserName() string {
	me := m.DiscordgoMessage
	if me.Author == nil {
		return ""
	}

	if m.Nick == nil {
		n := m.Discord.NicknameForID(me.Author.ID, me.Author.Username, me.ChannelID)
		m.Nick = &n
	}
	return *m.Nick
}

// UserID returns the message userID
func (m *DiscordMessage) UserID() string {
	if m.DiscordgoMessage.Author == nil {
		return ""
	}

	return m.DiscordgoMessage.Author.ID
}

// UserAvatar returns the url to the message senders avatar
func (m *DiscordMessage) UserAvatar() string {
	if m.DiscordgoMessage.Author == nil {
		return ""
	}

	return discordgo.EndpointUserAvatar(m.DiscordgoMessage.Author.ID, m.DiscordgoMessage.Author.Avatar)
}

// Message returns the message in human readable text
func (m *DiscordMessage) Message() string {
	if m.Content == nil {
		c := m.DiscordgoMessage.ContentWithMentionsReplaced()
		c = m.Discord.replaceRoleNames(m.DiscordgoMessage, c)
		c = m.Discord.replaceChannelNames(m.DiscordgoMessage, c)

		m.Content = &c
	}
	return *m.Content
}

// RawMessage gets the raw message
func (m *DiscordMessage) RawMessage() string {
	return m.DiscordgoMessage.Content
}

// MessageID gets the ID of the message
func (m *DiscordMessage) MessageID() string {
	return m.DiscordgoMessage.ID
}

// Type gets the type of the message
func (m *DiscordMessage) Type() MessageType {
	return m.MessageType
}

// Timestamp gets the timestamp of the message
func (m *DiscordMessage) Timestamp() (time.Time, error) {
	return m.DiscordgoMessage.Timestamp.Parse()
}
