# discordgobot

discordgobot is a wrapper for the [DiscordGo](https://github.com/bwmarrin/discordgo) framework that makes it easy to subscribe to commands with a plugin framework

### Usage

Import this package into your project

```go
import "github.com/lampjaw/discordgobot"
```

Construct a new bot

```go
bot, err := discordgobot.NewBot(token, commandPrefix, clientId, ownerUserId)
```

Register plugins on the bot

```go
bot.RegisterPlugin(NewMyAwesomePlugin())    
bot.RegisterPlugin(NewMyAwesomePlugin())
```

Start listening

```go
bot.Open()
```

Check out the Examples to see how everything is tied together and how to make a plugin.