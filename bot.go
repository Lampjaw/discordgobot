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
	Commands        map[string]*CommandDefinition
	Config          *GobotConf
	messageChannels []chan Message
	State           interface{}
}

// Open loads plugin data and starts listening for discord messages
func (b *Gobot) Open() {
	var invalidPlugin = false
	var invalidCommand = false

	for _, plugin := range b.Plugins {
		if !validatePlugin(plugin) {
			invalidPlugin = true
		}
	}

	if invalidPlugin {
		log.Printf("A misconfigured plugin was found.")
		return
	}

	for _, command := range b.Commands {
		if !validateCommand(command) {
			invalidCommand = true
		}
	}

	if invalidCommand {
		log.Printf("A misconfigured command was found.")
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

// RegisterCommand registers a command
func (b *Gobot) RegisterCommand(trigger string, description string, callback func(bot *Gobot, client *DiscordClient, payload CommandPayload)) {
	b.RegisterPrefixCommand("", trigger, description, callback)
}

// RegisterPrefixCommand registers a command with a static prefix
func (b *Gobot) RegisterPrefixCommand(prefix string, trigger string, description string, callback func(bot *Gobot, client *DiscordClient, payload CommandPayload)) {
	def := &CommandDefinition{
		CommandID:   fmt.Sprintf("gobot-cmd-%s", trigger),
		Description: description,
		Triggers: []string{
			trigger,
		},
		CommandPrefix: prefix,
		Callback:      callback,
	}
	b.Commands[def.CommandID] = def
}

// RegisterCommandDefinition registers a command definition
func (b *Gobot) RegisterCommandDefinition(cmdDef *CommandDefinition) {
	if b.Commands[cmdDef.CommandID] != nil {
		log.Println("Command with that id is already registered", cmdDef.CommandID)
	}
	b.Commands[cmdDef.CommandID] = cmdDef
}

// GetCommandPrefix returns the prefix as configured in the GobotConf or the default if none is available
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

		if isCommandsRequest(b.Client, commandPrefix, message) {
			go handleCommandsRequest(b, message, commandPrefix)
			continue
		}

		messageParts := strings.Fields(message.RawMessage())

		for _, command := range b.Commands {
			if !b.Client.IsMe(message) {
				go findCommandDefinitionCommandMatch(b, command, message, commandPrefix, messageParts)
			}
		}

		for _, plugin := range b.Plugins {
			go plugin.Message(b, b.Client, message)
			if !b.Client.IsMe(message) {
				go findPluginCommandMatch(b, plugin, message, commandPrefix, messageParts)
			}
		}
	}
}

func findPluginCommandMatch(b *Gobot, plugin IPlugin, message Message, commandPrefix string, parts []string) {
	if plugin.Commands() == nil || message.Message() == "" {
		return
	}

	for _, commandDefinition := range plugin.Commands() {
		findCommandDefinitionCommandMatch(b, commandDefinition, message, commandPrefix, parts)
	}
}

func findCommandDefinitionCommandMatch(b *Gobot, commandDefinition *CommandDefinition, message Message, commandPrefix string, parts []string) {
	if message.Message() == "" || !validateCommandAccess(b.Client, commandDefinition, message) {
		return
	}

	definitionPrefix := getPrefixFromCommand(b, b.Client, commandDefinition, message)

	if definitionPrefix == "" {
		definitionPrefix = commandPrefix
	}

	for _, trigger := range commandDefinition.Triggers {
		if isTriggerMatch, triggerMatch := findTriggerMatch(commandDefinition, trigger, definitionPrefix, parts, message); isTriggerMatch {
			if isArgumentMatch, parsedArgs := extractCommandArguments(message, triggerMatch, commandDefinition.Arguments); isArgumentMatch {
				log.Printf("<%s> %s: %s\n", message.Channel(), message.UserName(), message.RawMessage())

				payload := CommandPayload{
					Trigger:   trigger,
					Arguments: parsedArgs,
					Message:   message,
				}

				go commandDefinition.Callback(b, b.Client, payload)
			}
		}
	}
}

func findTriggerMatch(commandDefinition *CommandDefinition, commandTrigger string, definitionPrefix string, messageParts []string, message Message) (bool, string) {
	if messageParts[0] == definitionPrefix+commandTrigger {
		return true, messageParts[0]
	}

	if !commandDefinition.DisableTriggerOnMention && len(messageParts) > 1 {
		return message.IsMentionTrigger(commandTrigger)
	}

	return false, ""
}

func validateCommandAccess(client *DiscordClient, commandDefinition *CommandDefinition, message Message) bool {
	if commandDefinition.ExposureLevel > 0 {
		switch commandDefinition.ExposureLevel {
		case EXPOSURE_PRIVATE:
			if !client.IsPrivate(message) {
				return false
			}
		case EXPOSURE_PUBLIC:
			if client.IsPrivate(message) {
				return false
			}
		}
	}

	return validateCommandAccessPermission(client, commandDefinition.PermissionLevel, message)
}

func validateCommandAccessPermission(client *DiscordClient, permissionLevel PermissionLevel, message Message) bool {
	if permissionLevel <= 0 {
		return true
	}

	switch permissionLevel {
	case PERMISSION_USER:
		return true
	case PERMISSION_MODERATOR:
		if client.IsModerator(message) {
			return true
		}
		fallthrough
	case PERMISSION_ADMIN:
		if client.IsChannelOwner(message) {
			return true
		}
		fallthrough
	case PERMISSION_OWNER:
		if client.IsBotOwner(message) {
			return true
		}
	}

	return false
}

func extractCommandArguments(message Message, trigger string, arguments []CommandDefinitionArgument) (bool, map[string]string) {
	parsedArgs := make(map[string]string)

	if arguments == nil || len(arguments) == 0 {
		return true, parsedArgs
	}

	var argPatterns []string

	for i, argument := range arguments {
		pattern := ""

		if i == 0 {
			pattern = fmt.Sprintf("(?P<%s>%s)", argument.Alias, argument.Pattern)
		} else {
			pattern = fmt.Sprintf("(?:\\s+(?P<%s>%s))", argument.Alias, argument.Pattern)
		}

		if argument.Optional {
			pattern += "?"
		}

		argPatterns = append(argPatterns, pattern)
	}
	var pattern = fmt.Sprintf("^%s$", strings.Join(argPatterns, ""))

	var trimmedContent = strings.TrimSpace(strings.TrimPrefix(message.RawMessage(), fmt.Sprintf("%s", trigger)))
	pat := regexp.MustCompile(pattern)
	argsMatch := pat.FindStringSubmatch(trimmedContent)

	if len(argsMatch) == len(arguments)-1 && arguments[len(arguments)-1].Optional {
		argsMatch = append(argsMatch, "")
	}

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

func handleCommandsRequest(b *Gobot, message Message, commandPrefix string) {
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

	for _, commandDefinition := range b.Commands {
		if commandDefinition.Unlisted {
			continue
		}

		definitionPrefix := getPrefixFromCommand(b, b.Client, commandDefinition, message)

		if definitionPrefix == "" {
			definitionPrefix = commandPrefix
		}

		help = append(help, commandDefinition.Help(b.Client, definitionPrefix))
	}

	sort.Strings(help)

	if len(help) == 0 {
		help = []string{"No commands found"}
	}

	b.Client.SendMessage(message.Channel(), strings.Join(help, "\n"))
}

func isCommandsRequest(client *DiscordClient, commandPrefix string, message Message) bool {
	triggerTerm := "commands"
	commandTrigger := commandPrefix + triggerTerm

	triggered, _ := message.IsMentionTrigger(triggerTerm)

	return triggered || strings.HasPrefix(message.RawMessage(), commandTrigger)
}

func getPrefixFromCommand(bot *Gobot, client *DiscordClient, command *CommandDefinition, message Message) string {
	if command.CommandPrefixFunc != nil {
		return command.CommandPrefixFunc(bot, client, message)
	}

	if command.CommandPrefix != "" {
		return command.CommandPrefix
	}

	return ""
}
