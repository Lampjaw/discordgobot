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

### Overwritable plugin functions
* `func (p *Plugin) Name() string` - (Required) Returns the name of the plugin
* `func (p *Plugin) Load(*discordgobot.DiscordClient) error` - Loads plugin state
* `func (p *Plugin) Save() error` - Saves plugin state
* `func (p *Plugin) Help(*discordgobot.DiscordClient, discordgobot.Message, bool) []string` - Returns a help message for `?commands` calls
* `func (p *Plugin) Message(*discordgobot.Gobot, *discordgobot.DiscordClient, discordgobot.Message) error` - If your plugin looks at all messages and isn't triggered by commands use this to process every message.
* `func (p *Plugin) Commands() []discordgobot.CommandDefinition` - Returns an array of CommandDefinitions to listen for
