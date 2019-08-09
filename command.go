package discordgobot

import "fmt"

// CommandDefinition is the basic type for defining plugin commands
type CommandDefinition struct {
	// Description is a summary of the command that's returned in Help text.
	Description string
	// CommandID is an internal id used for internal tracking.
	CommandID string
	// Triggers are an array of strings used to determine if the commands has been called.
	Triggers []string
	// Arguments are an array of CommandDefinitionArgument types that define how to parse a message.
	Arguments []CommandDefinitionArgument
	// PermissionLevel is the minimum level of command access. Default is PERMISSION_USER.
	PermissionLevel PermissionLevel
	// ExposureLevel restricts commands from being processed in either public, private, or both settings. Default is EXPOSURE_EVERYWHERE.
	ExposureLevel ExposureLevel
	// Callback is a function reference that's called when a message meets trigger and argument requirements.
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

// PermissionLevel access required to execute command
type PermissionLevel int

const (
	PERMISSION_OWNER PermissionLevel = 1 + iota
	PERMISSION_ADMIN
	PERMISSION_MODERATOR
	PERMISSION_USER
)

// ExposureLevel expose commands to private messages, public, or both
type ExposureLevel int

const (
	EXPOSURE_EVERYWHERE ExposureLevel = 1 + iota
	EXPOSURE_PUBLIC
	EXPOSURE_PRIVATE
)

// IsValid determines if the command definition is configured correctly
func (c *CommandDefinition) IsValid() (bool, []string) {
	errors := make([]string, 0)

	if c.CommandID == "" {
		errors = append(errors, "No CommandID provided for CommandDefinition")
	}

	if c.Triggers == nil || len(c.Triggers) == 0 {
		errors = append(errors, "No triggers provided for CommandDefinition")
	}

	if c.Callback == nil {
		errors = append(errors, "No callback provided for CommandDefinition")
	}

	if c.Arguments != nil && len(c.Arguments) > 0 {
		for _, argument := range c.Arguments {
			if isValid, argErrors := argument.IsValid(); !isValid {
				errors = append(errors, argErrors...)
			}
		}
	}

	return len(errors) > 0, errors
}

// Help generates a help string from a CommandDefinition
func (c *CommandDefinition) Help(client *DiscordClient) string {
	var arguments []string

	if c.Arguments != nil && len(c.Arguments) > 0 {
		arguments = make([]string, len(c.Arguments))
		for i, argument := range c.Arguments {
			arguments[i] = argument.Alias
		}
	}

	return CommandHelp(client, c.Triggers[0], arguments, c.Description)
}

// IsValid determines if the command definition argument is configured correctly
func (c *CommandDefinitionArgument) IsValid() (bool, []string) {
	errors := make([]string, 0)

	if c.Pattern == "" {
		errors = append(errors, "No regex pattern provided for CommandDefinitionArgument")
	}

	if c.Alias == "" {
		errors = append(errors, "No argument alias provided for CommandDefinitionArgument")
	}

	return len(errors) > 0, errors
}

// CommandHelp is a helper message that creates help text for a command.
func CommandHelp(client *DiscordClient, command string, arguments []string, description string) string {
	commandString := fmt.Sprintf("%s%s", client.CommandPrefix(), command)

	if arguments != nil && len(arguments) > 0 {
		for _, argument := range arguments {
			commandString = fmt.Sprintf("%s <%s>", commandString, argument)
		}
	}

	return fmt.Sprintf("`%s` - %s", commandString, description)
}
