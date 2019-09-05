package main

import (
	"github.com/lampjaw/discordgobot"
)

type ExamplePlugin struct {
	discordgobot.Plugin
}

func NewExamplePlugin() discordgobot.IPlugin {
	return &ExamplePlugin{}
}

func (p *ExamplePlugin) Name() string {
	return "TestPlugin"
}

// Commands defines how we want to listen for things to execute on this plugin
func (p *ExamplePlugin) Commands() []*discordgobot.CommandDefinition {
	return []*discordgobot.CommandDefinition{
		&discordgobot.CommandDefinition{
			CommandID: "hello-command",
			Triggers: []string{
				"hello",
			},
			Description: "Displays hello world",
			Callback:    p.hellocallback,
		},
	}
}

func (p *ExamplePlugin) hellocallback(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload discordgobot.CommandPayload) {
	p.RLock()

	client.SendMessage(payload.Message.Channel(), "Hello, World!")

	p.RUnlock()
}
