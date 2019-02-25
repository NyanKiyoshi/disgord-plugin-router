package drouter

import (
	"github.com/andersfylling/disgord/event"
	"reflect"
	"regexp"
)

// shouldEnablePluginFunc defines the type of a module matcher.
// Returns false to disable a given plugin. Or true to allow it
// and proceed to the next matcher function.
type shouldEnablePluginFunc func(plugin *Plugin) bool

// RouterDefinition defines the structure of a bot plugins routing.
type RouterDefinition struct {
	// Plugins Contains every registered plugin.
	Plugins []*Plugin

	// ShouldEnablePluginFuncs Contains a list of matcher to test
	// against package paths. If it returns false, the plugin must disabled.
	ShouldEnablePluginFuncs []shouldEnablePluginFunc
}

// Router is the global plugin routing object.
// It should be used by plugins to register themselves
// during their initialization or else.
var Router = &RouterDefinition{}

// Plugin creates a new module from a given type that will generate
// the proper Plugin.ImportName value. And takes a human readable plugin name.
//
// Parameters:
//
//     pluginType   Any object defined in the plugin's package that will be used
//                  to extract the module's package path, such as "my-modules/stats".
//
//     name         The plugin's human readable name, such as "bot-statistics".
func (router *RouterDefinition) Plugin(pluginType interface{}, names ...string) *Plugin {
	newPlugin := &Plugin{
		ImportName:  reflect.TypeOf(pluginType).PkgPath(),
		RootCommand: newCommand(names...),
		Prefix:      DefaultPrefix,
		Listeners:   map[string][]interface{}{},
	}

	router.Plugins = append(router.Plugins, newPlugin)
	return newPlugin
}

// ShouldUse register matchers to test against plugins to keep enabled or to disable.
func (router *RouterDefinition) ShouldUse(pluginFuncs ...shouldEnablePluginFunc) *RouterDefinition {
	router.ShouldEnablePluginFuncs = append(router.ShouldEnablePluginFuncs, pluginFuncs...)
	return router
}

// ShouldNotUseRE registers a regex to be tested against plugins that should be disabled.
func (router *RouterDefinition) ShouldNotUseRE(regex string) *RouterDefinition {
	re := regexp.MustCompile(regex)
	router.ShouldUse(func(plugin *Plugin) bool {
		return !re.MatchString(plugin.ImportName)
	})
	return router
}

// isPluginEnabled check whether a given plugin should be enabled or not.
func (router *RouterDefinition) isPluginEnabled(plugin *Plugin) bool {
	for _, matcher := range router.ShouldEnablePluginFuncs {
		if !matcher(plugin) {
			return false
		}
	}

	return true
}

// Configure configures the bot's client with the router and plugins.
// 1. It hooks the internal events of the router;
// 2. It configures and enables every plugin that are enabled.
func (router *RouterDefinition) Configure(client routerClient) {
	// Add router events
	router.installInternalEvents(client)

	// Add plugins events
	for _, plugin := range router.Plugins {
		// Skip disabled modules
		if !router.isPluginEnabled(plugin) {
			continue
		}

		// Setup the plugin
		configurePlugin(plugin, client)
	}
}

// installInternalEvents installs the router internal events
// into a given client.
func (router *RouterDefinition) installInternalEvents(client routerClient) {
	if err := client.On(event.MessageCreate, router.OnMessageReceived); err != nil {
		LogFatalf("failed to register router's internal MessageCreate event: %s", err)
	}
}

// configurePlugin configures a client for a given plugin.
func configurePlugin(plugin *Plugin, client routerClient) {
	for eventName, handlers := range plugin.Listeners {
		if err := client.On(eventName, handlers...); err != nil {
			LogFatalf(
				"failed to register event %s for plugin %s: %s",
				eventName, plugin.ImportName, err)
		}
	}

	// Set the module as ready and enabled
	plugin.Activate()
}
