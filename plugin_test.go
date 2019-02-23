package drouter_test

import (
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlugin_Use(t *testing.T) {
	// Create a testing plugin
	plugin := createTestPlugin()

	// Create the dummy callback
	callback := func(ctx *drouter.Context) error {
		return successError
	}

	// Ensure there are not wrappers on an empty plugin
	assert.Empty(t, plugin.RootCommand.Wrappers)

	// Add the callback to the module
	plugin.Use(callback)

	// Ensure it was added
	assert.Len(t, plugin.RootCommand.Wrappers, 1)
	assert.Equal(t, successError, plugin.RootCommand.Wrappers[0](nil))
}

func TestPlugin_SetPrefix(t *testing.T) {
	// Create a testing plugin
	plugin := createTestPlugin()

	// Ensure the default prefix is used by default
	assert.Equal(t, drouter.DefaultPrefix, plugin.Prefix)

	// Set a new prefix
	plugin.SetPrefix("??")

	// Check if the prefix was correctly set
	assert.Equal(t, "??", plugin.Prefix)
}

func TestPlugin_Handler(t *testing.T) {
	// Create a testing plugin and a dummy handler
	plugin := createTestPlugin()
	handler := func(ctx *drouter.Context) error {
		return successError
	}

	// Ensure no handler is registered by default
	assert.Nil(t, plugin.RootCommand.HandlerFunc)

	// Register an handler to the plugin
	plugin.Handler(handler)

	// Check it was correctly set
	assert.NotNil(t, plugin.RootCommand.HandlerFunc)
	assert.Equal(t, successError, plugin.RootCommand.HandlerFunc(nil))
}

func TestPlugin_On(t *testing.T) {
	// Create a testing plugin and a dummy handler
	plugin := createTestPlugin()

	// Dummy handlers
	handler1 := struct {
		Handler1 bool
	}{}
	handler2 := struct {
		Handler2 bool
	}{}

	// Add a "hello" event
	plugin.On("hello", handler1)

	// Ensure it was correctly set
	handlers := plugin.Listeners["hello"]
	assert.NotNil(t, handlers)
	assert.Len(t, handlers, 1)
	assert.Equal(t, handler1, handlers[0])

	// Append another "hello" event
	plugin.On("hello", handler2)

	// Ensure it was appended
	handlers = plugin.Listeners["hello"]
	assert.Len(t, handlers, 2)
	assert.Equal(t, handler2, handlers[1])

	// Append an inexisting event
	plugin.On("notfound", handler1)

	// Ensure it was appended
	handlers = plugin.Listeners["notfound"]
	assert.Len(t, handlers, 1)
	assert.Equal(t, handler1, handlers[0])
}

func TestPlugin_Command(t *testing.T) {
	// Create a testing plugin
	plugin := createTestPlugin()

	// Ensure no commands are registered on empty plugin
	assert.Len(t, plugin.Commands, 0)

	// Register a new command
	cmd := plugin.Command("color", "colour")

	// Check if the command was correctly registered
	assert.EqualValues(t, []*drouter.Command{cmd}, plugin.Commands)

	// Check if the registered command is correct
	assert.ElementsMatch(t, []string{"color", "colour"}, cmd.Names.Keys())
	assert.Nil(t, cmd.HandlerFunc)
	assert.Empty(t, cmd.Wrappers)
}

func TestPlugin_Help(t *testing.T) {
	// Create a testing plugin
	plugin := createTestPlugin()

	// Ensure no help is set by default
	assert.Empty(t, plugin.RootCommand.ShortHelp)
	assert.Empty(t, plugin.RootCommand.LongHelp)

	// Set custom help
	plugin.Help("hello\nworld")

	// Check help is set
	assert.Equal(t, "hello", plugin.RootCommand.ShortHelp)
	assert.Equal(t, "hello\nworld", plugin.RootCommand.LongHelp)
}
