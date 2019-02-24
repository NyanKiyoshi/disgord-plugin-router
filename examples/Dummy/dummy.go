package main

import (
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/andersfylling/disgord"
	"os"
)

type _internal struct{}

func main() {
	client := disgord.New(&disgord.Config{BotToken: os.Getenv("DISCORD_TOKEN")})

	plugin := drouter.Router.Plugin(_internal{}, "ping").
		Handler(func(ctx *drouter.Context) error {
			return ctx.Say("pong!")
		})
	plugin.
		Command("miss").Handler(func(ctx *drouter.Context) error {
		return ctx.Say("I missed.")
	})

	// Setup the client from the router
	drouter.Router.Configure(client)

	// Connect to the discord gateway to receive events
	if err := client.Connect(); err != nil {
		panic(err)
	}

	// Wait for ever for interrupt. Then, disconnect.
	if err := client.DisconnectOnInterrupt(); err != nil {
		panic(err)
	}
}
