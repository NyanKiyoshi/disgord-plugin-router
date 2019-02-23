package drouter_test

import (
	"fmt"
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/NyanKiyoshi/disgord-plugin-router/mocks/mocked_disgord"
	"github.com/andersfylling/disgord"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type _myModuleInternalType struct{}

func createTestRouter() *drouter.RouterDefinition {
	return &drouter.RouterDefinition{}
}

func createTestPlugin() *drouter.Plugin {
	return createTestRouter().Plugin(_myModuleInternalType{}, "my-module")
}

func ExampleNew_newPlugin() {
	myPlugin := drouter.Router.Plugin(_myModuleInternalType{}, "my-plugin")
	fmt.Printf("ImportName: %s\nNames: %s", myPlugin.ImportName, myPlugin.RootCommand.Names.Keys())
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
	router.ShouldUse(func(plugin *drouter.Plugin) bool {
		return true
	})

	// Ensure the matcher was added as expected
	assert.Len(t, router.ShouldEnablePluginFuncs, 1)
	assert.True(t, router.ShouldEnablePluginFuncs[0](nil))
}

func TestRouter_ShouldNotUseRE(t *testing.T) {
	// Create dummy router and test plugins
	router := createTestRouter()
	enabledPlugin := &drouter.Plugin{ImportName: "enabled!!"}
	disabledPlugin := &drouter.Plugin{ImportName: "peach"}

	// Ensure there are not matcher by default
	assert.Len(t, router.ShouldEnablePluginFuncs, 0)

	// Add dummy matcher
	router.ShouldNotUseRE("^peach$")

	// Ensure the matcher was added and is behaving as expected
	assert.Len(t, router.ShouldEnablePluginFuncs, 1)
	assert.True(t, router.ShouldEnablePluginFuncs[0](enabledPlugin))
	assert.False(t, router.ShouldEnablePluginFuncs[0](disabledPlugin))
}

func TestRouter_Configure(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockedClient := mocked_disgord.NewMockrouterClient(mockCtrl)
	mockedClient.
		EXPECT().
		On(disgord.EventMessageCreate, gomock.Any()).
		Return(nil)

	router := createTestRouter()
	router.Configure(mockedClient)
}
