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

## Overwritable plugin functions
* `func (p *Plugin) Name() string` - (Required) Returns the name of the plugin
* `func (p *Plugin) Load(*discordgobot.DiscordClient) error` - Loads plugin state
* `func (p *Plugin) Save() error` - Saves plugin state
* `func (p *Plugin) Help(*discordgobot.DiscordClient, discordgobot.Message, bool) []string` - Returns a help message for `?commands` calls
* `func (p *Plugin) Message(*discordgobot.Gobot, *discordgobot.DiscordClient, discordgobot.Message) error` - If your plugin looks at all messages and isn't triggered by commands use this to process every message.
* `func (p *Plugin) Commands() []discordgobot.CommandDefinition` - Returns an array of CommandDefinitions to listen for

## Creating a command definition

A command definition is a built in way to tell a plugin when to run an action.

### [Model] CommandDefinition

`CommandID string` - (Required) A unique identifier for the command definition

`Triggers []string` - (Required) An array of activation terms. The command prefix is added automatically.

`Callback func(bot *discordgo.Gobot, client *discordgo.DiscordClient, message discordgo.Message, args map[string]string, trigger string)` - (Required) The callback function to use when a command is successfully called.

`Description string` - A short description used when the commands list is generated.

`Arguments []CommandDefinitionArgument` - Parsing rules for additional input arguments.

`PermissionLevel int` - An integer representing minimum permission required. Values are `PERMISSION_OWNER`, `PERMISSION_ADMIN`, `PERMISSION_MODERATOR`, and `PERMISSION_USER`. If no value is provided than `PERMISSION_USER` is used.

`ExposureLevel int` - An integer representing weather or not to allow commands to be restricted to private messages, guild channels, or both. Values are `EXPOSURE_EVERYWHERE`, `EXPOSURE_PUBLIC`, and `EXPOSURE_PRIVATE`, If no value is provided than `EXPOSURE_EVERYWHERE` is used.

### [Model] CommandDefinitionArgument

`Alias` - (Required) The alias is the key used when returning the argument map to the callback function.

`Pattern string` - (Required) A regex pattern to validate and extract the argument from.

`Optional bool` - If an argument is optional than the command will execute even if the argument isn't provided in the input.

### Example

```go
discordgobot.CommandDefinition{
    CommandID: "example-command",
    Triggers: []string{
        "doathing",
    },
    Arguments: []discordgobot.CommandDefinitionArgument{
        discordgobot.CommandDefinitionArgument{
            Pattern: "[a-zA-Z0-9]*",
            Alias:   "myFirstArg",
        },
        discordgobot.CommandDefinitionArgument{
            Pattern: ".*",
            Alias:   "everythingElseArg",
        },
    },
    Description: "Say hello",
    ExposureLevel: discordgobot.EXPOSURE_EVERYWHERE,
    PermissionLevel: discordgobot.PERMISSION_MODERATOR,
    Callback:    p.runCommand,
}

func (p *myCoolPlugin) runCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, message discordgobot.Message, args map[string]string, trigger string) {
    myFirstArg := args["myFirstArg"]
    everythingElseArg := args["everythingElseArg"]

    client.SendMessage(message.Channel(), "Hello!")
}
```