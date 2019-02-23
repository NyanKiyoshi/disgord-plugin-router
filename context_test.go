package discplugins_test

import (
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createDummyContext() *discplugins.Context {
	return &discplugins.Context{

	}
}

func TestContext_Reply(t *testing.T) {
	ctx := createDummyContext()
	assert.Nil(t, ctx.Reply("Hello world"))
}
