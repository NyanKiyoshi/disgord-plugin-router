package discplugins

import (
	"reflect"
	"regexp"
)

// routerClient defines the expected client type.
type routerClient interface {
	On(event string, inputs ...interface{}) error
}

// shouldEnablePluginFunc defines the type of a module matcher.
// Returns false to disable a given plugin. Or true to allow it
// and proceed to the next matcher function.
type shouldEnablePluginFunc func(plugin *Plugin) bool

// Router defines the structure of a bot plugins routing.
type Router struct {
	// Plugins contains every registered plugin.
	Plugins []Plugin

	// ShouldEnablePluginFuncs contains a list of matcher to test
	// against package paths. If it returns false, the plugin must disabled.
	ShouldEnablePluginFuncs []shouldEnablePluginFunc
}

// New creates a new plugin router.
// It manages every plugin and the routing mechanism.
func New() *Router {
	return &Router{}
}

// Plugin creates a new module from a given type that will generate
// the proper Plugin.ImportName value. And takes a human readable plugin name.
//
// Parameters:
//
//     pluginType   Any object defined in the plugin's package that will be used
//                  to extract the module's package path, such as "my-modules/stats".
//
//     name         The plugin's human readable name, such as "bot-statistics".
func (router *Router) Plugin(pluginType interface{}, names ...string) *Plugin {
	newPlugin := &Plugin{
		ImportName: reflect.TypeOf(pluginType).PkgPath(),
		RootCommand: Command{
			Names: names,
		},
		Prefix:    DefaultPrefix,
		Listeners: map[string][]interface{}{},
	}

	return newPlugin
}

// ShouldUse register matchers to test against plugins to keep enabled or to disable.
func (router *Router) ShouldUse(pluginFuncs ...shouldEnablePluginFunc) *Router {
	router.ShouldEnablePluginFuncs = append(router.ShouldEnablePluginFuncs, pluginFuncs...)
	return router
}

// ShouldNotUseRE registers a regex to be tested against plugins that should be disabled.
func (router *Router) ShouldNotUseRE(regex string) *Router {
	re := regexp.MustCompile(regex)
	router.ShouldUse(func(plugin *Plugin) bool {
		return !re.MatchString(plugin.ImportName)
	})
	return router
}

// Load loads the router and plugins into the bot client.
func (router *Router) Load(client routerClient) {
	// FIXME: Not implemented
}
