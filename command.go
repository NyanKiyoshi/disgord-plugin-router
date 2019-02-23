package discplugins

import (
	"regexp"
	"strings"
)

type matcherFunc func(input string) bool

// Command defines the structure of a plugin sub-command.
type Command struct {
	// Names lists the different aliases of the sub-command.
	Names []string

	// MatchFunc is called whenever a new message
	// (starting with the correct prefix) is received by the plugin
	// and return true or false if it should be handled by this command or not.
	MatchFunc matcherFunc

	// HandlerFunc contains the function to invoke
	// whenever the command is requested.
	HandlerFunc callbackFunc

	// Wrappers contains the functions to invoke before the command.
	Wrappers []callbackFunc

	// ShortHelp contains the short straightforward command help.
	ShortHelp string

	// LongHelp contains the long descriptive command documentation.
	LongHelp string
}

// Match sets the matching function from a given function.
func (cmd *Command) Match(matcherFunc matcherFunc) *Command {
	cmd.MatchFunc = matcherFunc
	return cmd
}

// MatchRE defines a matching function from a given regex.
func (cmd *Command) MatchRE(regex string) *Command {
	matcher := regexp.MustCompile(regex)

	cmd.MatchFunc = func(input string) bool {
		return matcher.MatchString(input)
	}

	return cmd
}

// Handler defines the function to invoke whenever the command
// is being invoked.
func (cmd *Command) Handler(callbackFunc callbackFunc) *Command {
	cmd.HandlerFunc = callbackFunc
	return cmd
}

// Use appends given callbacks to a command to call
// whenever a command is being invoked.
func (cmd *Command) Use(callbackFuncs ...callbackFunc) *Command {
	cmd.Wrappers = append(cmd.Wrappers, callbackFuncs...)
	return cmd
}

// Help sets the help text of a command. The first line is
// the short and straightforward documentation. The whole text
// is the long and descriptive documentation.
func (cmd *Command) Help(helpText string) *Command {
	if startPos := strings.Index(helpText, "\n"); startPos > -1 {
		cmd.ShortHelp = helpText[:startPos]
	} else {
		cmd.ShortHelp = helpText
	}
	cmd.LongHelp = helpText

	return cmd
}
