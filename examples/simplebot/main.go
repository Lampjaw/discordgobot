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

	b, err := discordgobot.NewBot(token, config)

	if err != nil {
		log.Println(err)
	}

	b.RegisterPlugin(NewExamplePlugin())

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
