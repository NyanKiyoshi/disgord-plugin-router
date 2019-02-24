package drouter

import "github.com/andersfylling/disgord"

type routerSession interface {
	Myself() (*disgord.User, error)

	SendMsg(
		channelID disgord.Snowflake,
		message *disgord.Message) (msg *disgord.Message, err error)

	SendMsgString(
		channelID disgord.Snowflake,
		content string) (msg *disgord.Message, err error)
}
