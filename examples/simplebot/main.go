package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/lampjaw/discordgobot"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string

func main() {
	if token == "" {
		fmt.Println("No token provided. Please run: simplebot -t <bot token>")
		return
	}

	q := make(chan bool)

	config := &discordgobot.GobotConf{
		CommandPrefix: "?",
	}

	b, err := discordgobot.NewBot(token, config, nil)

	if err != nil {
		log.Println(err)
	}

	b.RegisterPlugin(NewExamplePlugin())

	b.RegisterCommand("cmd",
		"this was registered with RegisterCommand!",
		func(b *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
			client.SendMessage(payload.Message.Channel(), "A RegisterCommand response!")
		},
	)

	b.RegisterPrefixCommand("??",
		"pcmd",
		"this was registered with RegisterPrefixCommand!",
		func(b *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
			client.SendMessage(payload.Message.Channel(), "A RegisterPrefixCommand response!")
		},
	)

	b.RegisterCommandDefinition(&discordgobot.CommandDefinition{
		CommandID: "command-definition-command",
		Triggers: []string{
			"def",
		},
		Description: "this was registered with RegisterCommandDefinition!",
		Callback: func(b *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
			client.SendMessage(payload.Message.Channel(), "A RegisterCommandDefinition response!")
		},
	})

	b.Open()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

out:
	for {
		select {
		case <-q:
			break out
		case <-c:
			break out
		}
	}

	b.Save()
}
