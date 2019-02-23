package drouter_test

import (
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/NyanKiyoshi/disgord-plugin-router/mocks/mocked_disgord"
	"github.com/andersfylling/disgord"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var channelID = disgord.NewSnowflake(123)

func createDummyMessage() *disgord.Message {
	return &disgord.Message{
		ChannelID: channelID,
	}
}

func createDummyContext(session *mocked_disgord.MockclientSession) *drouter.Context {
	return &drouter.Context{
		Message: createDummyMessage(),
		Session: session,
	}
}

func TestContext_Reply(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	const messageToSend = "Hello world"

	mockedSession := mocked_disgord.NewMockclientSession(mockCtrl)
	mockedSession.
		EXPECT().
		SendMsgString(channelID, messageToSend).
		Return(nil, nil)

	ctx := createDummyContext(mockedSession)
	assert.Nil(t, ctx.Reply("Hello world"))
}
