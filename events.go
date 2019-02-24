package drouter

import (
	"github.com/andersfylling/disgord"
)

type callbackFunc func(ctx *Context) error

func (router *RouterDefinition) onMessageReceived(session disgord.Session, event *disgord.MessageCreate) {
	myself, err := session.Myself()
	if err != nil || event.Message.Author.ID == myself.ID {
		return
	}

	args := parseMessage(event.Message)
	matchedPrefix, foundCommand := router.Find(args...)

	if foundCommand == nil {
		return
	}

	go router.dispatchMessage(&Context{
		Session:       session,
		Message:       event.Message,
		Command:       foundCommand,
		MatchedPrefix: matchedPrefix,
		Args:          args,
	})
}
