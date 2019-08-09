package discordgobot

import (
	"errors"
	"fmt"
)

// VERSION of discordgobot
const VERSION = "0.1.1"

// NewBot creates a new Gobot
func NewBot(token string, prefix string, clientID string, ownerUserID string) (b *Gobot, err error) {
	if token == "" {
		fmt.Println("No token provided. Please run: mutterblack -t <bot token>")
		return nil, errors.New("Missing discord token")
	}

	if prefix == "" {
		prefix = "?"
	}

	args := []interface{}{("Bot " + token)}

	bot := &Gobot{
		Client: &DiscordClient{
			args:                args,
			messageChan:         make(chan Message, 200),
			ApplicationClientID: clientID,
			OwnerUserID:         ownerUserID,
			prefix:              prefix,
		},
		Plugins: make(map[string]IPlugin, 0),
	}

	return bot, nil
}
