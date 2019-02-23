package discplugins_test

import (
	"fmt"
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/stretchr/testify/assert"
	"testing"
)

type _myModuleInternalType struct{}

func createTestRouter() *discplugins.Router {
	return discplugins.New()
}

func createTestPlugin() *discplugins.Plugin {
	return createTestRouter().Plugin(_myModuleInternalType{}, "my-module")
}

func ExampleNew_newPlugin() {
	myRouter := discplugins.New()
	myPlugin := myRouter.Plugin(_myModuleInternalType{}, "my-plugin")
	fmt.Printf("ImportName: %s\nNames: %s", myPlugin.ImportName, myPlugin.RootCommand.Names)
	// Output:
	// ImportName: github.com/NyanKiyoshi/disgord-plugin-router_test
	// Names: [my-plugin]
}

func TestNew(t *testing.T) {
	router := createTestRouter()
	assert.Len(t, router.Plugins, 0)
}

func TestRouter_ShouldUse(t *testing.T) {
	// Create dummy router
	router := createTestRouter()

	// Ensure there are not matcher by default
	assert.Len(t, router.ShouldEnablePluginFuncs, 0)

	// Add dummy matcher
	router.ShouldUse(func(plugin *discplugins.Plugin) bool {
		return true
	})

	// Ensure the matcher was added as expected
	assert.Len(t, router.ShouldEnablePluginFuncs, 1)
	assert.True(t, router.ShouldEnablePluginFuncs[0](nil))
}


func TestRouter_ShouldNotUseRE(t *testing.T) {
	// Create dummy router and test plugins
	router := createTestRouter()
	enabledPlugin := &discplugins.Plugin{ImportName: "enabled!!"}
	disabledPlugin := &discplugins.Plugin{ImportName: "peach"}

	// Ensure there are not matcher by default
	assert.Len(t, router.ShouldEnablePluginFuncs, 0)

	// Add dummy matcher
	router.ShouldNotUseRE("^peach$")

	// Ensure the matcher was added and is behaving as expected
	assert.Len(t, router.ShouldEnablePluginFuncs, 1)
	assert.True(t, router.ShouldEnablePluginFuncs[0](enabledPlugin))
	assert.False(t, router.ShouldEnablePluginFuncs[0](disabledPlugin))
}

func TestRouter_Load(t *testing.T) {
	router := createTestRouter()
	router.Load(nil)
}
