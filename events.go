package drouter

import (
	"github.com/andersfylling/disgord"
	"log"
	"strings"
)

type callbackFunc func(ctx *Context) error

func (router *RouterDefinition) onMessageReceived(session disgord.Session, event *disgord.MessageCreate) {
	args := strings.Fields(event.Message.Content)
	foundCommand := router.Find(args...)

	if foundCommand == nil {
		return
	}

	ctx := &Context{
		Session: session,
		Message: event.Message,
	}
	go func() {
		var err error

		// Run wrappers
		for _, wrappingFunc := range foundCommand.Wrappers {
			err = wrappingFunc(ctx)

			if err != nil {
				break
			}
		}

		// If the wrappers did not report any error, invoke the command
		if err == nil {
			err = foundCommand.HandlerFunc(ctx)
		}

		// If the command returned an error, report it to the user
		// (as is, it's expected to be already user formatted errors)
		if err != nil {
			err = ctx.Say(err.Error())
		}

		// If we failed to communicate the error (as reply), there may
		// be a permission issue or else. Thus, log the error.
		if err != nil {
			log.Printf("failed to communicate %s error: %s", foundCommand.Names, err)
		}
	}()
}
