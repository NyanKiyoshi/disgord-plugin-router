package drouter

import (
	"github.com/NyanKiyoshi/disgord-plugin-router/internal/stringset"
	"regexp"
	"strings"
)

type matcherFunc func(input string) bool

// command defines the structure of a plugin sub-command.
type command struct {
	// Names lists the different aliases of the sub-command.
	Names stringset.StringSet

	// MatchFunc is called whenever a new message
	// (starting with the correct prefix) is received by the plugin
	// and return true or false if it should be handled by this command or not.
	MatchFunc matcherFunc

	// HandlerFunc Contains the function to invoke
	// whenever the command is requested.
	HandlerFunc callbackFunc

	// Wrappers Contains the functions to invoke before the command.
	Wrappers []callbackFunc

	// ShortHelp Contains the short straightforward command help.
	ShortHelp string

	// LongHelp Contains the long descriptive command documentation.
	LongHelp string
}

// newCommand creates and initialize a new command object.
func newCommand(names ...string) command {
	return command{
		Names: stringset.NewStringSet(names...),
	}
}

// Match sets the matching function from a given function.
func (cmd *command) Match(matcherFunc matcherFunc) *command {
	cmd.MatchFunc = matcherFunc
	return cmd
}

// MatchRE defines a matching function from a given regex.
func (cmd *command) MatchRE(regex string) *command {
	matcher := regexp.MustCompile(regex)

	cmd.MatchFunc = func(input string) bool {
		return matcher.MatchString(input)
	}

	return cmd
}

// Handler defines the function to invoke whenever the command
// is being invoked.
func (cmd *command) Handler(callbackFunc callbackFunc) *command {
	cmd.HandlerFunc = callbackFunc
	return cmd
}

// Use appends given callbacks to a command to call
// whenever a command is being invoked.
func (cmd *command) Use(callbackFuncs ...callbackFunc) *command {
	cmd.Wrappers = append(cmd.Wrappers, callbackFuncs...)
	return cmd
}

// Help sets the help text of a command. The first line is
// the short and straightforward documentation. The whole text
// is the long and descriptive documentation.
func (cmd *command) Help(helpText string) *command {
	if startPos := strings.Index(helpText, "\n"); startPos > -1 {
		cmd.ShortHelp = helpText[:startPos]
	} else {
		cmd.ShortHelp = helpText
	}
	cmd.LongHelp = helpText

	return cmd
}

// IsMatching returns true if the command name exists or
// if it matches the matching function, if provided.
func (cmd *command) IsMatching(targetCommand string) bool {
	return cmd.Names.Contains(targetCommand) || (cmd.MatchFunc != nil && cmd.MatchFunc(targetCommand))
}
