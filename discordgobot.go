package discordgobot

import (
	"errors"
	"fmt"
)

// VERSION of discordgobot
const VERSION = "0.2.0"

// NewBot creates a new Gobot
func NewBot(token string, config *GobotConf) (b *Gobot, err error) {
	if token == "" {
		fmt.Println("No token provided. Please run: mutterblack -t <bot token>")
		return nil, errors.New("Missing discord token")
	}

	args := []interface{}{("Bot " + token)}

	bot := &Gobot{
		Client: &DiscordClient{
			args:        args,
			messageChan: make(chan Message, 200),
			OwnerUserID: config.OwnerUserID,
		},
		Plugins: make(map[string]IPlugin, 0),
		Config:  config,
	}

	return bot, nil
}
