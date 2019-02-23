package discplugins_test

import (
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createDummyCommand() *discplugins.Command {
	return &discplugins.Command{
		Names: []string{"test"},
	}
}

func TestCommand_Match(t *testing.T) {
	// Create a testing command and a dummy handler
	cmd := createDummyCommand()

	// Ensure the default matching function is nil
	assert.Nil(t, cmd.MatchFunc)

	// Register the custom matcher
	cmd.Match(func(input string) bool {
		return true
	})

	// Check the custom matcher was set
	assert.NotNil(t, cmd.MatchFunc)
	assert.Equal(t, true, cmd.MatchFunc(""))
}

func TestCommand_MatchRE(t *testing.T) {
	// Create a testing command and a dummy handler
	cmd := createDummyCommand()

	// Ensure the default matching function is nil
	assert.Nil(t, cmd.MatchFunc)

	// Register the custom matcher
	cmd.MatchRE("^p([a-z]+)ch$")

	// Check the custom matcher was set
	assert.NotNil(t, cmd.MatchFunc)

	// Check the matcher regex is working as expected
	assert.Equal(t, true, cmd.MatchFunc("punch"))
	assert.Equal(t, true, cmd.MatchFunc("peach"))
	assert.Equal(t, false, cmd.MatchFunc("hello"))
}

func TestCommand_Handler(t *testing.T) {
	// Create a testing command and a dummy handler
	cmd := createDummyCommand()
	handler := func(ctx *discplugins.Context) error {
		return successError
	}

	// Ensure no handler is registered by default
	assert.Nil(t, cmd.HandlerFunc)

	// Register an handler to the plugin
	cmd.Handler(handler)

	// Check it was correctly set
	assert.NotNil(t, cmd.HandlerFunc)
	assert.Equal(t, successError, cmd.HandlerFunc(nil))
}

func TestCommand_Use(t *testing.T) {
	// Create a testing plugin
	cmd := createDummyCommand()

	// Create the dummy callback
	callback := func(ctx *discplugins.Context) error {
		return successError
	}

	// Ensure there are not wrappers on an empty plugin
	assert.Empty(t, cmd.Wrappers)

	// Add the callback to the module
	cmd.Use(callback)

	// Ensure it was added
	assert.Len(t, cmd.Wrappers, 1)
	assert.Equal(t, successError, cmd.Wrappers[0](nil))
}

var helpCommandTests = []struct {
	in                string
	expectedShortText string
}{
	{in: "Two\nLines:)", expectedShortText: "Two"},
	{in: "Only one line", expectedShortText: "Only one line"},
	{in: "", expectedShortText: ""},
	{in: " ", expectedShortText: " "},
}

func TestCommand_Help(t *testing.T) {
	// Create a testing plugin
	cmd := createDummyCommand()

	// Ensure the help texts are empty by default
	assert.Empty(t, cmd.ShortHelp)
	assert.Empty(t, cmd.LongHelp)

	for _, tt := range helpCommandTests {
		t.Run(tt.in, func(t *testing.T) {
			// Set the new help
			cmd.Help(tt.in)

			// Check the results
			assert.Equal(t, tt.expectedShortText, cmd.ShortHelp)
			assert.Equal(t, tt.in, cmd.LongHelp)
		})
	}
}