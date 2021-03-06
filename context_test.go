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
var selfUserID = disgord.NewSnowflake(456)

func createDummyMessage() *disgord.Message {
	return &disgord.Message{
		ChannelID: channelID,
		Author: &disgord.User{
			ID: selfUserID,
		},
		Content: "hello",
	}
}

func createDummyContext(session *mocked_disgord.MockrouterSession) *drouter.Context {
	return &drouter.Context{
		Message: createDummyMessage(),
		Session: session,
	}
}

func TestContext_Say(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	const messageToSend = "Hello world"

	mockedSession := mocked_disgord.NewMockrouterSession(mockCtrl)
	mockedSession.
		EXPECT().
		SendMsgString(channelID, messageToSend).
		Return(nil, nil)

	ctx := createDummyContext(mockedSession)
	assert.Nil(t, ctx.Say("Hello world"))
}

func TestContext_SayComplex(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var messageToSend = &disgord.Message{}

	mockedSession := mocked_disgord.NewMockrouterSession(mockCtrl)
	mockedSession.
		EXPECT().
		SendMsg(channelID, messageToSend).
		Return(nil, nil)

	ctx := createDummyContext(mockedSession)
	assert.Nil(t, ctx.SayComplex(messageToSend))
}
