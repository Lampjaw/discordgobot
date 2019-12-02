package discordgobot

import (
	"errors"
	"fmt"

	"github.com/lampjaw/discordclient"
)

// VERSION of discordgobot
const VERSION = "0.4.0"

// NewBot creates a new Gobot
func NewBot(token string, config *GobotConf, state interface{}) (b *Gobot, err error) {
	if token == "" {
		fmt.Println("No token provided.")
		return nil, errors.New("Missing discord token")
	}

	bot := &Gobot{
		Client: &DiscordClient{
			DiscordClient: discordclient.NewDiscordClient(token, config.OwnerUserID, config.ClientID),
		},
		Plugins:  make(map[string]IPlugin, 0),
		Commands: make(map[string]*CommandDefinition, 0),
		Config:   config,
		State:    state,
	}

	return bot, nil
}
