package discplugins

// Context defines callbacks invocation context.
type Context struct {

}

// Reply replies to a message (Can be from a guild channel or DM).
func (ctx *Context) Reply(message string) error {
	return nil
}
