package drouter

import (
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

func (router *RouterDefinition) Find(args ...string) *Command {
	argCount := len(args)

	if argCount < 1 {
		return nil
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
				return subCommand
			}
		}

		return &plugin.RootCommand
	}

	return nil
}
