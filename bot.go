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

// Gobot handles bot related functionality
type Gobot struct {
	Client          *DiscordClient
	Plugins         map[string]IPlugin
	messageChannels []chan Message
}

// Open loads plugin data and starts listening for discord messages
func (b *Gobot) Open() {
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

		if handleCommandsRequest(b, message) {
			continue
		}

		plugins := b.Plugins
		for _, plugin := range plugins {
			go plugin.Message(b, b.Client, message)
			if !b.Client.IsMe(message) {
				go findCommandMatch(b, plugin, message)
			}
		}
	}
}

func findCommandMatch(b *Gobot, plugin IPlugin, message Message) {
	if plugin.Commands() == nil || message.Message() == "" {
		return
	}

	for _, commandDefinition := range plugin.Commands() {
		for _, trigger := range commandDefinition.Triggers {
			var trig = b.Client.CommandPrefix() + trigger
			var parts = strings.Split(message.Message(), " ")

			if parts[0] == trig {
				log.Printf("<%s> %s: %s\n", message.Channel(), message.UserName(), message.Message())

				if commandDefinition.Arguments == nil {
					commandDefinition.Callback(b, b.Client, message, nil, trigger)
					return
				}

				parsedArgs := extractCommandArguments(message, trig, commandDefinition.Arguments)

				if parsedArgs != nil {
					commandDefinition.Callback(b, b.Client, message, parsedArgs, trigger)
					return
				}
			}
		}
	}
}

func extractCommandArguments(message Message, trigger string, arguments []CommandDefinitionArgument) map[string]string {
	var argPatterns []string
	for _, argument := range arguments {
		argPatterns = append(argPatterns, fmt.Sprintf("(?P<%s>%s)", argument.Alias, argument.Pattern))
	}
	var pattern = fmt.Sprintf("^%s$", strings.Join(argPatterns, " "))

	var trimmedContent = strings.TrimPrefix(message.Message(), fmt.Sprintf("%s ", trigger))
	pat := regexp.MustCompile(pattern)
	argsMatch := pat.FindStringSubmatch(trimmedContent)

	parsedArgs := make(map[string]string)

	if argsMatch == nil || len(argsMatch) == 1 {
		return nil
	}

	for i := 1; i < len(argsMatch); i++ {
		parsedArgs[pat.SubexpNames()[i]] = argsMatch[i]
	}

	if len(parsedArgs) != len(arguments) {
		return nil
	}

	return parsedArgs
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func handleCommandsRequest(b *Gobot, message Message) bool {
	var trig = b.Client.CommandPrefix() + "commands"
	var parts = strings.Split(message.Message(), " ")

	if parts[0] != trig {
		return false
	}

	help := []string{}

	for _, plugin := range b.Plugins {
		var h []string

		if plugin.Commands() == nil {
			h = plugin.Help(b.Client, message, false)
		} else {
			for _, commandDefinition := range plugin.Commands() {
				h = append(h, commandDefinition.Help(b.Client))
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