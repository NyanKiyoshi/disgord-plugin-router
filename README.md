<div align='center'>
  <h1>disgord-plugin-router</h1>
  <p>
    <a href='https://travis-ci.org/NyanKiyoshi/disgord-plugin-router'>
      <img src='https://travis-ci.org/NyanKiyoshi/disgord-plugin-router.svg?branch=master'
           alt='Build Status' />
    </a>
    <a href='https://codecov.io/gh/NyanKiyoshi/disgord-plugin-router'>
      <img src='https://codecov.io/gh/NyanKiyoshi/disgord-plugin-router/branch/master/graph/badge.svg'
           alt='Code coverage' />
    </a>
    <a href='https://codeclimate.com/github/NyanKiyoshi/disgord-plugin-router/maintainability'>
      <img src='https://api.codeclimate.com/v1/badges/cabe385f12ad3b20d336/maintainability'
           alt='Maintainability' />
    </a>
  </p>
  <p>
    <a href='https://godoc.org/github.com/NyanKiyoshi/disgord-plugin-router'>
      <img src='https://godoc.org/github.com/NyanKiyoshi/disgord-plugin-router?status.svg'
           alt='Godoc' />
    </a>
  </p>
</div>

A plugin management and routing mechanism for
[Disgord](https://github.com/andersfylling/disgord).

**WARNING: this repository is still in active development and is not yet ready 
for the proper expected user experience. See the [roadmap](#roadmap) for more details.**


## Design Approach
The approach of this library is that the bot developer will create a plugin 
for each command, each plugin can have a set of sub-commands.

### Example Design
```
- Roles -> lists the roles when invoked (e.g.: /roles)
    - add {role-name}               -> appends a given role
    - remove {role-name}            -> pops a given role

- Colors [color-name] -> lists the roles or set the user's color when invoked 
                         (e.g.: /colors or /colors blue)
    - add {color-hex} {color-name}  -> registers a new color
    - remove {color-name}           -> removes a given color
```

## Features
<!-- TODO: update links -->
- Can define or not a command to be executed on root commands
  ([read](#) or [example](#))
- Can define sub-commands (one level max by design, see above section)
  ([read](#) or [example](#))
- Can define [Disgord events](https://godoc.org/github.com/andersfylling/disgord/event) 
  handlers per plugin ([read](#) or [example](#))
- Can define typed arguments with auto parsing and error management ([read](#) or [example](#))
- Can blacklist/ disable pluging using patterns or matchers ([read](#) or [example](#))
- Can define plugin `setUp` and `tearDown` handlers ([read](#) or [example](#))
- Command handlers/ callbacks are context based ([read](#) or [example](#))
- Plugins can have separate or the same command prefix everywhere ([read](#) or [example](#))
- Errors can be returned by handlers and directly reported to the user ([read](#) or [example](#))
- Modular, modules can be imported from anywhere using go modules ([read](#) or [example](#))
- ...It's open source, fully tested and made with love! 🚀


## Usage
```go
package main

import (
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/andersfylling/disgord"
)

type _internal struct {}

func main() {
	client := disgord.New(&disgord.Config{BotToken: "YOUR BOT TOKEN"})
	
	drouter.Router.Plugin(_internal{}, "ping").
		Handler(func(ctx *drouter.Context) error {
			return ctx.Say("pong!")
		}).
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
```


## Roadmap
- [ ] Implement arguments parsing with dynamic typing (instead of strings).
- [ ] Context arguments should not contain the dispatched root command or sub command.
- [ ] Drop plugins pointer array in the router.
- [ ] Implement `SetUp(...)` and `TearDown(...)` functions on the plugins.
- [ ] Add final documentation and complete examples for each feature.
