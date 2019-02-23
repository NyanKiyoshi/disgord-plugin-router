package drouter

import (
	"github.com/andersfylling/disgord"
	"strings"
)

type callbackFunc func(ctx *Context) error

func (router *RouterDefinition) onMessageReceived(session disgord.Session, event *disgord.MessageCreate) {
	args := strings.Fields(event.Message.Content)
	router.find(args...)

	// TODO: dispatch
	//ctx := &Context{
	//	Session: session,
	//	Message: event.Message,
	//}
}
