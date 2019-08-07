package discordgobot

import (
	"sync"
)

// IPlugin is the universal required interface for all plugins
type IPlugin interface {
	// Name returns the name of the plugin
	Name() string
	// Load retrieves state information for the plugin
	Load(*DiscordClient) error
	// Save stores state information for the plugin
	Save() error
	// Help returns an optional response for when ?commands is called
	Help(*DiscordClient, Message, bool) []string
	// Message is a callback from every incoming message. Setting Commands is recommended unless you need to see everything.
	Message(*Gobot, *DiscordClient, Message) error
	// Commands returns an array of CommandDefinitions
	Commands() []CommandDefinition
}

// Plugin is the basic model to build bot plugins off of
type Plugin struct {
	sync.RWMutex
}

// Commands returns an array of CommandDefinitions
func (p *Plugin) Commands() []CommandDefinition {
	return nil
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return ""
}

// Load retrieves state information for the plugin
func (p *Plugin) Load(client *DiscordClient) error {
	return nil
}

// Save stores state information for the plugin
func (p *Plugin) Save() error {
	return nil
}

// Help returns an optional response for when ?commands is called
func (p *Plugin) Help(client *DiscordClient, message Message, detailed bool) []string {
	return nil
}

// Message is a callback from every incoming message. Setting Commands is recommended unless you need to see everything.
func (p *Plugin) Message(bot *Gobot, client *DiscordClient, message Message) error {
	return nil
}
