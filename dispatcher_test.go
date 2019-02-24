package drouter_test

import (
	"bytes"
	"fmt"
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/NyanKiyoshi/disgord-plugin-router/mocks/mocked_disgord"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var routerDefinitionFindTests = []struct {
	name                    string
	in                      []string
	shouldReturnRootCommand bool
	shouldReturnSubCommand  bool
}{
	{"valid root command", []string{"?color"}, true, false},
	{"valid root command (alias)", []string{"?colour"}, true, false},
	{"invalid root command", []string{"?no"}, false, false},

	{"valid sub command", []string{"?color", "blue"}, false, true},
	{"invalid sub command => should invoke root", []string{"?color", "black"}, true, false},

	{"no arguments", []string{}, false, false},
	{"whitespace only", []string{" "}, false, false},
	{"prefix only", []string{"?"}, false, false},
	{"invalid prefix", []string{"!color"}, false, false},
}

func TestRouterDefinition_Find(t *testing.T) {
	var (
		foundCommand *drouter.Command
		prefix       string
	)

	router := createTestRouter()
	plugin := router.Plugin(_myModuleInternalType{}, "color", "colour").SetPrefix("?")

	// Should return nil as the plugin is not yet enabled
	prefix, foundCommand = router.Find("color")
	assert.Nil(t, foundCommand, "plugin should not be enabled")
	assert.Empty(t, prefix, "prefix of a non found command should be empty")

	// Add dummy commands
	plugin.Command("red")
	subCommand := plugin.Command("blue")

	// Enable the plugin
	plugin.Activate()

	for _, tt := range routerDefinitionFindTests {
		t.Run(fmt.Sprintf("%s: %s", tt.name, tt.in), func(t *testing.T) {
			prefix, foundCommand = router.Find(tt.in...)

			if tt.shouldReturnRootCommand {
				assert.Equal(t, &plugin.RootCommand, foundCommand, "expected root command")
				assert.Equal(t, "?", prefix)
			} else if tt.shouldReturnSubCommand {
				assert.Equal(t, subCommand, foundCommand, "expected sub command")
				assert.Equal(t, "?", prefix)
			} else {
				assert.Nil(t, foundCommand, "expected nil, command should'nt because found")
				assert.Empty(t, prefix, "prefix of a non found command should be empty")
			}
		})
	}
}

func TestDispatchMessage_EverythingValid(t *testing.T) {
	var (
		WrapperCalled        bool
		CommandHandlerCalled bool
	)

	// Create a dummy command
	command := createDummyCommand()

	// Add a non-errored wrapper that we will check if it gets called
	command.Use(func(ctx *drouter.Context) error {
		WrapperCalled = true
		return nil
	})

	// Add a non-errored command handler that we will check if it gets called
	command.Handler(func(ctx *drouter.Context) error {
		CommandHandlerCalled = true
		return nil
	})

	// Execute the dispatcher, it is expected to succeed
	isSuccessChan := make(chan bool, 1)
	drouter.DispatchMessage(&drouter.Context{
		Command: command,
	}, isSuccessChan)

	// Check if it was successful, as we expect
	assert.True(t, <-isSuccessChan)

	// Check if the wrapper and command handler were called
	assert.True(t, WrapperCalled, "the command wrapper was not called")
	assert.True(t, CommandHandlerCalled, "the command handler was not called")
}

func TestDispatchMessage_HandlesFuncErrors(t *testing.T) {
	// Create a dummy command
	command := createDummyCommand("ping")

	var checkResultFunc = func(checkT *testing.T, discordAPIError error) {
		// Create mocked session
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockedSession := mocked_disgord.NewMockrouterSession(mockCtrl)

		// Create the test context
		ctx := createDummyContext(mockedSession)
		ctx.Command = command
		ctx.MatchedPrefix = "!"

		// Create the result channel
		isSuccessChan := make(chan bool, 1)

		// Prepare the mock
		mockedSession.
			EXPECT().
			SendMsgString(channelID, successError.Error()).
			Return(ctx.Message, discordAPIError)

		// Execute the dispatcher, it is expected to succeed
		drouter.DispatchMessage(ctx, isSuccessChan)

		// Check if it was successful, as we expect
		assert.False(checkT, <-isSuccessChan)
	}

	// Test wrapper errors are handled
	t.Run("against wrapper", func(t *testing.T) {
		oldWrappers := command.Wrappers

		// Add a wrapper that will return an error
		command.Use(func(ctx *drouter.Context) error {
			return successError
		})

		checkResultFunc(t, nil)
		command.Wrappers = oldWrappers
	})

	// Test command's handler errors are handled
	t.Run("against command handler", func(t *testing.T) {
		// Add a wrapper that will return an error
		command.Handler(func(ctx *drouter.Context) error {
			return successError
		})

		checkResultFunc(t, nil)
	})

	// Test if it failed to report an error to the user
	// (e.g.: not enough permission), it reports to the logger instead.
	t.Run("against discord API", func(t *testing.T) {
		var dummyBuffer bytes.Buffer
		log.SetOutput(&dummyBuffer)
		defer log.SetOutput(os.Stderr)

		checkResultFunc(t, successError)

		// Check the error was sent
		assert.Contains(
			t, dummyBuffer.String(),
			"failed to communicate prefix: ! / names: map[ping:{}]: success",
			"wrong message error/ not called",
		)
	})
}

var parseMessageTests = []struct {
	in  string
	out []string
}{
	{"ping", []string{"ping"}},
	{"ping   pong", []string{"ping", "pong"}},
	{"    ", []string{}},
	{"", []string{}},
}

func TestParseMessage(t *testing.T) {
	for _, tt := range parseMessageTests {
		t.Run(tt.in, func(t *testing.T) {
			assert.EqualValues(t, tt.out, drouter.ParseMessage(tt.in))
		})
	}
}
