package drouter_test

import (
	"fmt"
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/stretchr/testify/assert"
	"testing"
)

var routerDefinitionFindTests = []struct {
	name                    string
	in                      []string
	shouldReturnRootCommand bool
	shouldReturnSubCommand  bool
}{
	{"valid root command", []string{"?color"}, true, false},
	{"valid root command (alias)", []string{"?colour"}, true, false},
	{"invalid root command", []string{"?no"}, false, false},

	{"valid sub command", []string{"?color", "blue"}, false, true},
	{"invalid sub command => should invoke root", []string{"?color", "black"}, true, false},

	{"no arguments", []string{}, false, false},
	{"whitespace only", []string{" "}, false, false},
	{"prefix only", []string{"?"}, false, false},
	{"invalid prefix", []string{"!color"}, false, false},
}

func TestRouterDefinition_Find(t *testing.T) {
	var (
		foundCommand *drouter.Command
		prefix       string
	)

	router := createTestRouter()
	plugin := router.Plugin(_myModuleInternalType{}, "color", "colour").SetPrefix("?")

	// Should return nil as the plugin is not yet enabled
	prefix, foundCommand = router.Find("color")
	assert.Nil(t, foundCommand, "plugin should not be enabled")
	assert.Empty(t, prefix, "prefix of a non found command should be empty")

	// Add dummy commands
	plugin.Command("red")
	subCommand := plugin.Command("blue")

	// Enable the plugin
	plugin.Activate()

	for _, tt := range routerDefinitionFindTests {
		t.Run(fmt.Sprintf("%s: %s", tt.name, tt.in), func(t *testing.T) {
			prefix, foundCommand = router.Find(tt.in...)

			if tt.shouldReturnRootCommand {
				assert.Equal(t, &plugin.RootCommand, foundCommand, "expected root command")
				assert.Equal(t, "?", prefix)
			} else if tt.shouldReturnSubCommand {
				assert.Equal(t, subCommand, foundCommand, "expected sub command")
				assert.Equal(t, "?", prefix)
			} else {
				assert.Nil(t, foundCommand, "expected nil, command should'nt because found")
				assert.Empty(t, prefix, "prefix of a non found command should be empty")
			}
		})
	}
}
