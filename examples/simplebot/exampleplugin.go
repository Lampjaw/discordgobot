package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

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
func (p *ExamplePlugin) Commands() []discordgobot.CommandDefinition {
	return []discordgobot.CommandDefinition{
		discordgobot.CommandDefinition{
			CommandID: "hello-command",
			Triggers: []string{
				"hello",
			},
			Description: "Displays hello world",
			Callback:    p.hellocallback,
		},
	}
}

// Optional Load override loads plugin data from disk
func (p *ExamplePlugin) Load(client *discordgobot.DiscordClient) error {
	fileName := "data/" + p.Name()

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Printf("Error reading plugin save data %s. %v", p.Name(), err)
		return err
	}

	if data != nil {
		if err := json.Unmarshal(data, p); err != nil {
			log.Println("Error loading data", err)
			return err
		}
	}

	return nil
}

// Optional Save override saves plugin data to disk
func (p *ExamplePlugin) Save() error {
	if err := os.Mkdir("data", os.ModePerm); err != nil {
		if !os.IsExist(err) {
			log.Println("Error creating service directory.")
			return err
		}
	}

	data, err := json.Marshal(p)

	if err != nil {
		log.Printf("Error marshaling plugin data %s. %v", p.Name(), err)
		return err
	}

	err = ioutil.WriteFile("data/"+p.Name(), data, os.ModePerm)

	if err != nil {
		log.Printf("Error saving plugin %s. %v", p.Name(), err)
		return err
	}

	return nil
}

func (p *ExamplePlugin) hellocallback(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, message discordgobot.Message, args map[string]string, trigger string) {
	p.RLock()

	client.SendMessage(message.Channel(), "Hello, World!")

	p.RUnlock()
}
