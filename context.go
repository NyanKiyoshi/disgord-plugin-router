package drouter

import "github.com/andersfylling/disgord"

// Context defines callbacks invocation context.
type Context struct {
	// Message contains the received message
	Message *disgord.Message

	// Session contains the received discord session
	Session routerSession

	// Command is the matched command
	Command *Command

	// MatchedPrefix is the command (plugin's) prefix
	MatchedPrefix string

	// Args contains the received arguments
	// Deprecated: it will totally change in a future release
	Args []string
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
