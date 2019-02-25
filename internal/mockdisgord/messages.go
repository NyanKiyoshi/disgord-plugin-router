package mockdisgord

import "github.com/andersfylling/disgord"

// ChannelID represents a dummy message's ChannelID
var ChannelID = disgord.NewSnowflake(123)

// SelfUserID represents a dummy message's author ID
var SelfUserID = disgord.NewSnowflake(456)

// CreateDummyMessage creates a disgord test message.
func CreateDummyMessage() *disgord.Message {
	return &disgord.Message{
		ChannelID: ChannelID,
		Author: &disgord.User{
			ID: SelfUserID,
		},
		Content: "hello",
	}
}
