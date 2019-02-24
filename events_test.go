package drouter_test

import (
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/NyanKiyoshi/disgord-plugin-router/mocks/mocked_disgord"
	"github.com/andersfylling/disgord"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockMySelf(
	t *testing.T, senderID disgord.Snowflake,
	testFunc func(*testing.T, *mocked_disgord.MockrouterSession)) {

	// Start a new mocking controller
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Mock session and MySelf()
	mockedSession := mocked_disgord.NewMockrouterSession(mockCtrl)
	mockedSession.
		EXPECT().
		Myself().Return(&disgord.User{ID: senderID}, nil)

	// Run test
	testFunc(t, mockedSession)
}

// It ensures:
// - That it's ignoring messages that are the bot itself (returning nil).
// - And that it returns nil when the command is not found.
func TestRouterDefinition_OnMessageReceived_HandlesErrors(t *testing.T) {
	router := createTestRouter()

	t.Run("ignores self messages", func(t *testing.T) {
		mockMySelf(t, selfUserID, func(t *testing.T, session *mocked_disgord.MockrouterSession) {
			// It should return nil
			assert.Nil(t, router.OnMessageReceived(session, &disgord.MessageCreate{
				Message: createDummyMessage(),
			}))
		})
	})

	t.Run("returns nil when invalid command", func(t *testing.T) {
		mockMySelf(t, disgord.Snowflake(555), func(t *testing.T, session *mocked_disgord.MockrouterSession) {
			// It should return nil
			assert.Nil(t, router.OnMessageReceived(session, &disgord.MessageCreate{
				Message: createDummyMessage(),
			}))
		})
	})

	t.Run("dispatches command when valid request", func(t *testing.T) {
		// Will contain whether the ping command was invoked or not
		var pingHandlerWasCalled bool

		// Create a !ping module
		router.Plugin(_myModuleInternalType{}, "ping").Handler(func(ctx *drouter.Context) error {
			pingHandlerWasCalled = true
			return nil
		}).SetPrefix("!").Activate()

		// Create a valid command message
		message := createDummyMessage()
		message.Content = "!ping"

		mockMySelf(t, disgord.Snowflake(555), func(t *testing.T, session *mocked_disgord.MockrouterSession) {
			// Dispatch the event
			successChannel := router.OnMessageReceived(session, &disgord.MessageCreate{
				Message: message,
			})

			// Ensure it was a success
			assert.True(t, <-successChannel)
		})

		// Ensure the command was invoked
		assert.True(t, pingHandlerWasCalled)
	})
}
