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

Or you can register individual commands directly

```go
bot.RegisterCommand("6roll", "rolls a 6 sided die", callbackFunc)
bot.RegisterPrefixCommand("?", "20roll", "rolls a 20 sided die", callbackFunc)
bot.RegisterCommandDefinition(myCommandDefinition)
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
* `func (p *Plugin) Help(*discordgobot.Gobot, *discordgobot.DiscordClient, discordgobot.Message, bool) []string` - Returns a help message for `?commands` calls
* `func (p *Plugin) Message(*discordgobot.Gobot, *discordgobot.DiscordClient, discordgobot.Message) error` - If your plugin looks at all messages and isn't triggered by commands use this to process every message.
* `func (p *Plugin) Commands() []discordgobot.CommandDefinition` - Returns an array of CommandDefinitions to listen for

## Creating a command definition

A command definition is a built in way to tell a plugin when to run an action.

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

func (p *myCoolPlugin) runCommand(bot *discordgobot.Gobot, client *discordgobot.DiscordClient, payload CommandPayload) {
    myFirstArg := payload.Arguments["myFirstArg"]
    everythingElseArg := payload.Arguments["everythingElseArg"]

    client.SendMessage(payload.message.Channel(), "Hello!")
}
```

## Methods

`NewBot(token string, config GobotConf, state interface{}) (b *Gobot, err error)` 

`RegisterPlugin(plugin IPlugin) void` - Registers a plugin to process messages or commands

`RegisterCommand(trigger string, description string, callback func(bot *Gobot, client *DiscordClient, payload CommandPayload)) void` - Registers a command

`RegisterPrefixCommand(prefix string, trigger string, description string, callback func(bot *Gobot, client *DiscordClient, payload CommandPayload)) void` - Registers a command with a static prefix

`RegisterCommandDefinition(cmdDef *CommandDefinition) void` - Registers a command definition

`RemoveCommand(commandID string)` - Unregisters a command. Does not effect plugins.

`UpdateCommandDefinition(cmdDef *CommandDefinition)` - Updates a command definition or registers if it doesn't exist. Does not effect plugins.

`GetCommandPrefix(message Message) string` - Returns the prefix as configured in the GobotConf or the default if none is available

`Open() void` - Loads plugin data and starts listening for discord messages

`Save() void` - Writes all plugin data to disk


## Models

### [Model] GobotConf

`CommandPrefix string` - A string thats prefixed to every command trigger. Defaults to "?".

`CommandPrefixFunc func(bot *Gobot, client *DiscordClient, message Message) string` - If set, allows for a CommandPrefix to be dynamically chosen during processing.

`ClientID string` - The known client id of the bot. Potentially useful in some plugins.

`OwnerUserId string` - OwnerUserID is the OwnerUserId. Needed for processing commands restricted to the owner permission.

`CommandLookupDisabled bool` - Allows for the `?commands` command to be disabled

### [Model] CommandDefinition

`CommandID string` - (Required) A unique identifier for the command definition

`Triggers []string` - (Required) An array of activation terms. The command prefix is added automatically.

`Callback func(bot *discordgo.Gobot, client *discordgo.DiscordClient, payload CommandPayload)` - (Required) The callback function to use when a command is successfully called.

`Description string` - A short description used when the commands list is generated.

`Arguments []CommandDefinitionArgument` - Parsing rules for additional input arguments.

`PermissionLevel int` - An integer representing minimum permission required. Values are `PERMISSION_OWNER`, `PERMISSION_ADMIN`, `PERMISSION_MODERATOR`, and `PERMISSION_USER`. If no value is provided than `PERMISSION_USER` is used.

`ExposureLevel int` - An integer representing weather or not to allow commands to be restricted to private messages, guild channels, or both. Values are `EXPOSURE_EVERYWHERE`, `EXPOSURE_PUBLIC`, and `EXPOSURE_PRIVATE`, If no value is provided than `EXPOSURE_EVERYWHERE` is used.

`Unlisted bool` - Prevents the command from being displayed in the commands list lookup when set to true.

`DisableTriggerOnMention bool` - Prevents a command from being triggered when a user uses @BotName when set to true. example: `@BotName <trigger> <argument>`

`CommandPrefix string` - Allows a different prefix to be set compared to the rest of the bot commands.

`CommandPrefixFunc func(bot *Gobot, client *DiscordClient, message Message) string` - An optional function to be called if the command prefix should be assigned dynamically.

### [Model] CommandDefinitionArgument

`Alias` - (Required) The alias is the key used when returning the argument map to the callback function.

`Pattern string` - (Required) A regex pattern to validate and extract the argument from.

`Optional bool` - If an argument is optional than the command will execute even if the argument isn't provided in the input.

### [Model] CommandPayload

`CommandID` - The identifier of the command definition

`Message Message` - The entire message received that activated the command

`Arguments map[string]string` - Contain a hash of all configured CommandDefinitionArguments that could be parsed

`Trigger string` - The specific string that activated the command
	