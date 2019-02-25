package drouter

// DefaultPrefix defines the default command prefix of plugins.
var DefaultPrefix = "/"

// Plugin defines the structure of a disgord plugin.
type Plugin struct {
	// ImportName defines the import name of the plugin.
	ImportName string

	// Prefix defines the commands prefix.
	Prefix string

	// RootCommand defines the base plugin's root/ base command
	RootCommand Command

	// Listeners defines the different event handlers
	// of the plugin (see https://godoc.org/github.com/andersfylling/disgord/event).
	Listeners map[string][]interface{}

	// Wrappers Contains the registered sub-commands of the plugin.
	Commands []*Command

	// IsLoaded is true is the module was loaded and installed into the client.
	IsLoaded bool
}

// Use appends given callbacks to a plugin to call
// whenever a command is being invoked.
func (plugin *Plugin) Use(callbackFuncs ...callbackFunc) *Plugin {
	// FIXME: we should put them as global wrappers in Plugin
	//        instead of the root command.
	plugin.RootCommand.Use(callbackFuncs...)
	return plugin
}

// SetPrefix sets the plugin commands prefix (can be empty for no prefix).
func (plugin *Plugin) SetPrefix(prefix string) *Plugin {
	plugin.Prefix = prefix
	return plugin
}

// Handler defines the function to invoke whenever the plugin command
// is being invoked.
func (plugin *Plugin) Handler(callbackFunc callbackFunc) *Plugin {
	plugin.RootCommand.Handler(callbackFunc)
	return plugin
}

// On registers given handlers to be invoked whenever the event is fired.
func (plugin *Plugin) On(eventName string, inputs ...interface{}) *Plugin {
	existingEvents := plugin.Listeners[eventName]

	if existingEvents != nil {
		plugin.Listeners[eventName] = append(existingEvents, inputs...)
	} else {
		plugin.Listeners[eventName] = inputs
	}

	return plugin
}

// Command creates a new sub-command for the plugin.
func (plugin *Plugin) Command(names ...string) *Command {
	command := newCommand(names...)
	plugin.Commands = append(plugin.Commands, &command)
	return &command
}

// Help sets the help text of a command. The first line is
// the short and straightforward documentation. The whole text
// is the long and descriptive documentation.
func (plugin *Plugin) Help(helpText string) *Plugin {
	plugin.RootCommand.Help(helpText)
	return plugin
}

// activate marks a plugin as ready.
func (plugin *Plugin) activate() {
	// TODO: we should dispatch setUp(...)
	plugin.IsLoaded = true
}
