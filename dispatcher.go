package drouter

import (
	"strings"
)

func (router *RouterDefinition) find(args ...string) *Command {
	argCount := len(args)

	if argCount < 1 {
		return nil
	}

	// The raw command, with prefix (if any)
	prefixedCommand := args[0]

	for _, plugin := range router.Plugins {
		// Skip modules that are not active
		if !plugin.IsReady {
			continue
		}

		// Skip modules that don't match the message prefix
		if !strings.HasPrefix(prefixedCommand, plugin.Prefix) {
			continue
		}
	}

	return nil
}
