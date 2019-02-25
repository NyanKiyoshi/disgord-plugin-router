package mockdisgord

import "github.com/andersfylling/disgord"

var ChannelID = disgord.NewSnowflake(123)
var SelfUserID = disgord.NewSnowflake(456)

func CreateDummyMessage() *disgord.Message {
	return &disgord.Message{
		ChannelID: ChannelID,
		Author: &disgord.User{
			ID: SelfUserID,
		},
		Content: "hello",
	}
}
