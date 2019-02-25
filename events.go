package drouter

import (
	"github.com/andersfylling/disgord"
)

type callbackFunc func(ctx *Context) error

// onMessageReceived is invoked whenever a new message is created,
// it will dispatch instructions if the message is a command,
// and if the message is not from the bot itself.
func (router *RouterDefinition) onMessageReceived(session disgord.Session, event *disgord.MessageCreate) {
	myself, err := session.Myself()
	if err != nil || event.Message.Author.ID == myself.ID {
		return
	}

	args := parseMessage(event.Message.Content)
	matchedPrefix, foundCommand := router.Find(args...)

	if foundCommand == nil {
		return
	}

	go dispatchMessage(&Context{
		Session:       session,
		Message:       event.Message,
		Command:       foundCommand,
		MatchedPrefix: matchedPrefix,
		Args:          args,
	}, make(chan bool, 0))
}
