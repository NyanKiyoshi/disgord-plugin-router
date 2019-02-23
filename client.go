package drouter

import "github.com/andersfylling/disgord"

// routerClient defines the expected client type.
type routerClient interface {
	On(event string, inputs ...interface{}) error
	CreateChannelMessage(
		channelID disgord.Snowflake,
		params *disgord.CreateMessageParams) (ret *disgord.Message, err error)
}
