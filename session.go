package drouter

import "github.com/andersfylling/disgord"

type clientSession interface {
	SendMsg(
		channelID disgord.Snowflake,
		message *disgord.Message) (msg *disgord.Message, err error)
	SendMsgString(
		channelID disgord.Snowflake,
		content string) (msg *disgord.Message, err error)
}
