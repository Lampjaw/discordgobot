package discordgobot

import "fmt"

// CommandDefinition is the basic type for defining plugin commands
type CommandDefinition struct {
	// Description is a summary of the command that's returned in Help text
	Description string
	// CommandID is an internal id used for internal tracking
	CommandID string
	// Triggers are an array of strings used to determine if the commands has been called
	Triggers []string
	// Arguments are an array of CommandDefinitionArgument types that define how to parse a message
	Arguments []CommandDefinitionArgument
	// Callback is a function reference that's called when a message meets trigger and argument requirements
	Callback func(bot *Gobot, client *DiscordClient, message Message, args map[string]string, trigger string)
}

// CommandDefinitionArgument defines parameters to parse from message text
type CommandDefinitionArgument struct {
	// Optional determines if this argument is required to process the command
	Optional bool
	// Pattern holds a regex to match message parts
	Pattern string
	// Alias is the name of the parameter to return when the argument map is sent to the CommandDefinition Callback
	Alias string
}

// Help generates a help string from a CommandDefinition
func (c *CommandDefinition) Help(client *DiscordClient) string {
	commandString := fmt.Sprintf("%s%s", client.CommandPrefix(), c.Triggers[0])

	if len(c.Arguments) > 0 {
		for _, argument := range c.Arguments {
			commandString = fmt.Sprintf("%s <%s>", commandString, argument.Alias)
		}
	}

	return fmt.Sprintf("`%s` - %s", commandString, c.Description)
}
