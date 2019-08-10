package discordgobot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// DEFAULT_COMMAND_PREFIX is the default prefix character if none is configured
const DEFAULT_COMMAND_PREFIX = "?"

// GobotConf defines operational parameters to be used with NewBot
type GobotConf struct {
	// CommandPrefix is a string thats prefixed to every command trigger.
	CommandPrefix string
	// CommandPrefixFunc allows for a CommandPrefix to be dynamically chosen during processing.
	CommandPrefixFunc func(bot *Gobot, client *DiscordClient, message Message) string
	// ClientID sets the known client id of the bot. Potentially useful in some plugins.
	ClientID string
	// OwnerUserID is the OwnerUserId. Needed for processing commands restricted to the owner permission.
	OwnerUserID string
	// CommandLookupDisabled allows for the ?commands command to be disabled
	CommandLookupDisabled bool
}

// Gobot handles bot related functionality
type Gobot struct {
	Client          *DiscordClient
	Plugins         map[string]IPlugin
	Config          *GobotConf
	messageChannels []chan Message
}

// Open loads plugin data and starts listening for discord messages
func (b *Gobot) Open() {
	var invalidPlugin = false

	for _, plugin := range b.Plugins {
		if !validatePlugin(plugin) {
			invalidPlugin = true
		}
	}

	if invalidPlugin {
		log.Printf("A misconfigured plugin was found.")
		return
	}

	if messageChan, err := b.Client.Open(); err == nil {
		for _, plugin := range b.Plugins {
			plugin.Load(b.Client)
		}
		go b.listen(messageChan)
	} else {
		log.Printf("Error creating discord service: %v\n", err)
	}
}

// Save writes all plugin data to disk
func (b *Gobot) Save() {
	for _, plugin := range b.Plugins {
		plugin.Save()
	}
}

// RegisterPlugin registers a plugin to process messages or commands
func (b *Gobot) RegisterPlugin(plugin IPlugin) {
	if b.Plugins[plugin.Name()] != nil {
		log.Println("Plugin with that name already registered", plugin.Name())
	}
	b.Plugins[plugin.Name()] = plugin
}

// GetCommandPrefix returns the prefix as configured in the GobotConf or the default if non is available
func (b *Gobot) GetCommandPrefix(message Message) string {
	if b.Config != nil {
		if b.Config.CommandPrefixFunc != nil {
			return b.Config.CommandPrefixFunc(b, b.Client, message)
		}

		if b.Config.CommandPrefix != "" {
			return b.Config.CommandPrefix
		}
	}

	return DEFAULT_COMMAND_PREFIX
}

func (b *Gobot) getData(plugin IPlugin) []byte {
	fileName := "data/" + plugin.Name()

	if b, err := ioutil.ReadFile(fileName); err == nil {
		return b
	}

	return nil
}

func (b *Gobot) listen(messageChan <-chan Message) {
	log.Printf("Listening")
	for {
		message := <-messageChan

		commandPrefix := b.GetCommandPrefix(message)

		if handleCommandsRequest(b, message, commandPrefix) {
			continue
		}

		plugins := b.Plugins
		for _, plugin := range plugins {
			go plugin.Message(b, b.Client, message)
			if !b.Client.IsMe(message) {
				go findCommandMatch(b, plugin, message, commandPrefix)
			}
		}
	}
}

func findCommandMatch(b *Gobot, plugin IPlugin, message Message, commandPrefix string) {
	if plugin.Commands() == nil || message.Message() == "" {
		return
	}

	for _, commandDefinition := range plugin.Commands() {
		if commandDefinition.ExposureLevel > 0 {
			switch commandDefinition.ExposureLevel {
			case EXPOSURE_PRIVATE:
				if !b.Client.IsPrivate(message) {
					return
				}
			case EXPOSURE_PUBLIC:
				if b.Client.IsPrivate(message) {
					return
				}
			}
		}

		if commandDefinition.PermissionLevel > 0 {
			switch commandDefinition.PermissionLevel {
			case PERMISSION_MODERATOR:
				if !b.Client.IsModerator(message) {
					return
				}
				fallthrough
			case PERMISSION_ADMIN:
				if !b.Client.IsChannelOwner(message) {
					return
				}
				fallthrough
			case PERMISSION_OWNER:
				if !b.Client.IsBotOwner(message) {
					return
				}
			}
		}

		definitionPrefix := getPrefixFromCommand(b, b.Client, commandDefinition, message)

		if definitionPrefix == "" {
			definitionPrefix = commandPrefix
		}

		for _, trigger := range commandDefinition.Triggers {
			var trig = definitionPrefix + trigger
			var parts = strings.Split(message.Message(), " ")

			if parts[0] == trig {
				log.Printf("<%s> %s: %s\n", message.Channel(), message.UserName(), message.Message())

				if isMatch, parsedArgs := extractCommandArguments(message, trig, commandDefinition.Arguments); isMatch {
					commandDefinition.Callback(b, b.Client, message, parsedArgs, trigger)
					return
				}
			}
		}
	}
}

func extractCommandArguments(message Message, trigger string, arguments []CommandDefinitionArgument) (bool, map[string]string) {
	parsedArgs := make(map[string]string)

	if arguments == nil || len(arguments) == 0 {
		return true, parsedArgs
	}

	var argPatterns []string

	for _, argument := range arguments {
		argPatterns = append(argPatterns, fmt.Sprintf("(?P<%s>%s)", argument.Alias, argument.Pattern))
	}
	var pattern = fmt.Sprintf("^%s$", strings.Join(argPatterns, " "))

	var trimmedContent = strings.TrimPrefix(message.Message(), fmt.Sprintf("%s ", trigger))
	pat := regexp.MustCompile(pattern)
	argsMatch := pat.FindStringSubmatch(trimmedContent)

	if argsMatch == nil || len(argsMatch) == 1 {
		return false, nil
	}

	for i := 1; i < len(argsMatch); i++ {
		parsedArgs[pat.SubexpNames()[i]] = argsMatch[i]
	}

	if len(parsedArgs) != len(arguments) {
		return false, nil
	}

	return true, parsedArgs
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func handleCommandsRequest(b *Gobot, message Message, commandPrefix string) bool {
	if b.Config.CommandLookupDisabled {
		return false
	}

	var trig = commandPrefix + "commands"

	var parts = strings.Split(message.Message(), " ")

	if parts[0] != trig {
		return false
	}

	help := []string{}

	for _, plugin := range b.Plugins {
		var h []string

		helpResult := plugin.Help(b, b.Client, message, false)

		if helpResult != nil {
			h = helpResult
		} else if plugin.Commands() != nil {
			for _, commandDefinition := range plugin.Commands() {
				if commandDefinition.Unlisted {
					continue
				}

				definitionPrefix := getPrefixFromCommand(b, b.Client, commandDefinition, message)

				if definitionPrefix == "" {
					definitionPrefix = commandPrefix
				}

				h = append(h, commandDefinition.Help(b.Client, definitionPrefix))
			}
		}

		if h != nil && len(h) > 0 {
			help = append(help, h...)
		}
	}

	sort.Strings(help)

	if len(help) == 0 {
		help = []string{"No commands found"}
	}

	b.Client.SendMessage(message.Channel(), strings.Join(help, "\n"))

	return true
}

func getPrefixFromCommand(bot *Gobot, client *DiscordClient, command CommandDefinition, message Message) string {
	if command.CommandPrefixFunc != nil {
		return command.CommandPrefixFunc(bot, client, message)
	}

	if command.CommandPrefix != "" {
		return command.CommandPrefix
	}

	return ""
}
