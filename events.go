package drouter

import (
	"github.com/andersfylling/disgord"
)

type callbackFunc func(ctx *Context) error

// OnMessageReceived is invoked whenever a new message is created,
// it will dispatch instructions if the message is a command,
// and if the message is not from the bot itself.
func (router *RouterDefinition) OnMessageReceived(session disgord.Session, event *disgord.MessageCreate) {
	myself, err := session.Myself()
	if err != nil || event.Message.Author.ID == myself.ID {
		return
	}

	args := ParseMessage(event.Message.Content)
	matchedPrefix, foundCommand := router.Find(args...)

	if foundCommand == nil {
		return
	}

	go DispatchMessage(&Context{
		Session:       session,
		Message:       event.Message,
		Command:       foundCommand,
		MatchedPrefix: matchedPrefix,
		Args:          args,
	}, make(chan bool, 0))
}
