package drouter

import (
	"github.com/andersfylling/disgord"
	"log"
	"strings"
)

func findPluginCommand(plugin *Plugin, receivedCommand string) *Command {
	for _, cmdObj := range plugin.Commands {
		if cmdObj.IsMatching(receivedCommand) {
			return cmdObj
		}
	}

	return nil
}

// Find looks for a matching command. Returns the
// matched prefix and command, if found.
func (router *RouterDefinition) Find(args ...string) (string, *Command) {
	argCount := len(args)

	if argCount < 1 {
		return "", nil
	}

	// The raw command, with prefix (if any)
	prefixedCommand := args[0]

	for _, plugin := range router.Plugins {
		// Skip plugins that are not active
		if !plugin.IsReady {
			continue
		}

		// Skip plugins that don't match the message prefix
		if !strings.HasPrefix(prefixedCommand, plugin.Prefix) {
			continue
		}

		// Skip commands that only have a prefix
		if len(plugin.Prefix) == len(prefixedCommand) {
			continue
		}

		// Remove the command prefix
		receivedCommand := prefixedCommand[len(plugin.Prefix):]

		// Skip plugins that don't match the command
		if !plugin.RootCommand.IsMatching(receivedCommand) {
			continue
		}

		// Try to Find a sub-command if a second argument was supplied
		if len(args) > 1 {
			// Return the sub-command if the second argument matches a sub-command
			if subCommand := findPluginCommand(plugin, args[1]); subCommand != nil {
				return plugin.Prefix, subCommand
			}
		}

		return plugin.Prefix, &plugin.RootCommand
	}

	return "", nil
}

func (router *RouterDefinition) dispatchMessage(ctx *Context) {
	var err error

	// Run wrappers
	for _, wrappingFunc := range ctx.Command.Wrappers {
		err = wrappingFunc(ctx)

		if err != nil {
			break
		}
	}

	// If the wrappers did not report any error, invoke the command
	if err == nil {
		err = ctx.Command.HandlerFunc(ctx)
	}

	// If the command returned an error, report it to the user
	// (as is, it's expected to be already user formatted errors)
	if err != nil {
		err = ctx.Say(err.Error())
	}

	// If we failed to communicate the error (as reply), there may
	// be a permission issue or else. Thus, log the error.
	if err != nil {
		log.Printf(
			"failed to communicate %s%s error: %s",
			ctx.MatchedPrefix, ctx.Command.Names, err)
	}
}

func parseMessage(message *disgord.Message) []string {
	return strings.Fields(message.Content)
}
