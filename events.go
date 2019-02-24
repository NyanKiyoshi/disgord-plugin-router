package drouter

import (
	"github.com/andersfylling/disgord"
)

type callbackFunc func(ctx *Context) error

// OnMessageReceived is invoked whenever a new message is created,
// it will dispatch instructions if the message is a command,
// and if the message is not from the bot itself.
func (router *RouterDefinition) OnMessageReceived(session routerSession, event *disgord.MessageCreate) chan bool {
	myself, err := session.Myself()
	if err != nil || event.Message.Author.ID == myself.ID {
		return nil
	}

	args := ParseMessage(event.Message.Content)
	matchedPrefix, foundCommand := router.Find(args...)

	if foundCommand == nil {
		return nil
	}

	isSuccessChan := make(chan bool, 1)

	go DispatchMessage(&Context{
		Session:       session,
		Message:       event.Message,
		Command:       foundCommand,
		MatchedPrefix: matchedPrefix,
		Args:          args,
	}, isSuccessChan)

	return isSuccessChan
}
