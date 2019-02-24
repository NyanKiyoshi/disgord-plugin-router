package drouter

import "github.com/andersfylling/disgord"

// Context defines callbacks invocation context.
type Context struct {
	// Message Contains the received message
	Message *disgord.Message

	// Session Contains the received discord session
	Session clientSession
}

// Say replies to a message with a given string.
func (ctx *Context) Say(message string) error {
	_, err := ctx.Session.SendMsgString(ctx.Message.ChannelID, message)
	return err
}

// SayComplex replies to a message with a given message object.
func (ctx *Context) SayComplex(message *disgord.Message) error {
	_, err := ctx.Session.SendMsg(ctx.Message.ChannelID, message)
	return err
}
