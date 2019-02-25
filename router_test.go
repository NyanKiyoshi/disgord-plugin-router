package drouter_test

import (
	"fmt"
	"github.com/NyanKiyoshi/disgord-plugin-router"
	"github.com/NyanKiyoshi/disgord-plugin-router/mocks/mocked_disgord"
	"github.com/andersfylling/disgord"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type _myModuleInternalType struct{}

func createTestRouter() *drouter.RouterDefinition {
	return &drouter.RouterDefinition{}
}

func Example_newPlugin() {
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

func TestRouterDefinition_Plugin(t *testing.T) {
	router := createTestRouter()

	// Ensure there are no plugins by default
	assert.Empty(t, router.Plugins)

	// Register a new plugin
	plugin := router.Plugin(_myModuleInternalType{})

	// Ensure the plugin was registerd
	assert.NotEmpty(t, router.Plugins)
	assert.NotNil(t, router.Plugins[0])
	assert.Equal(t, plugin, router.Plugins[0])
}

func TestRouterDefinition_ShouldUse(t *testing.T) {
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

func TestRouterDefinition_ShouldNotUseRE(t *testing.T) {
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

// Test .Configure(...) without registering any plugin to configure.
func TestRouterDefinition_Configure_OnlyInstallInternalEvents(t *testing.T) {
	// Backup log.Fatal before mocking it
	oldLogFatal := drouter.LogFatalf
	defer func() { drouter.LogFatalf = oldLogFatal }()

	runTest := func(subT *testing.T, customErr error) {
		// Mock log.Fatal
		var receivedLog string
		drouter.LogFatalf = func(format string, v ...interface{}) {
			receivedLog = fmt.Sprintf(format, v...)
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockedClient := mocked_disgord.NewMockrouterClient(mockCtrl)
		mockedClient.
			EXPECT().
			On(disgord.EventMessageCreate, gomock.Any()).
			Return(customErr)

		router := createTestRouter()
		router.Configure(mockedClient)

		// Check if log.Fatal was called if an error was set
		if customErr != nil {
			assert.Equal(
				subT,
				"failed to register router's internal MessageCreate event: successful test",
				receivedLog,
			)
		}
	}

	t.Run("with proper configuration", func(t *testing.T) {
		runTest(t, nil)
	})

	t.Run("with errored configuration", func(t *testing.T) {
		runTest(t, errors.New("successful test"))
	})
}

func TestRouterDefinition_Configure_WithRegisteredPlugins(t *testing.T) {
	// Backup log.Fatal before mocking it
	oldLogFatal := drouter.LogFatalf
	defer func() { drouter.LogFatalf = oldLogFatal }()

	router := createTestRouter().ShouldUse(func(plugin *drouter.Plugin) bool {
		return plugin.Prefix != "disabled"
	})
	plugin := router.Plugin(_myModuleInternalType{}, "ping").
		On(disgord.EventMessageCreate, nil) // Register a dummy event
	disabledPlugin := router.Plugin(_myModuleInternalType{}).SetPrefix("disabled")

	// Ensure it is false by default
	assert.False(t, plugin.IsLoaded)

	runTest := func(subT *testing.T, customErr error) {
		// Mock log.Fatal
		var receivedLog string
		drouter.LogFatalf = func(format string, v ...interface{}) {
			receivedLog = fmt.Sprintf(format, v...)
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockedClient := mocked_disgord.NewMockrouterClient(mockCtrl)
		mockedClient.
			EXPECT().
			On(disgord.EventMessageCreate, gomock.Any()).
			Return(customErr).
			Times(2) // should be called twice: internal event & dummy custom event

		// Run configure
		router.Configure(mockedClient)

		// Ensure the plugin was enabled and the disabled one was ignore
		assert.True(t, plugin.IsLoaded)
		assert.False(t, disabledPlugin.IsLoaded)

		// Check if log.Fatal was called if an error was set
		if customErr != nil {
			assert.Equal(
				subT,
				"failed to register event MESSAGE_CREATE "+
					"for plugin github.com/NyanKiyoshi/disgord-plugin-router_test: successful test",
				receivedLog,
			)
		}
	}

	t.Run("with proper event", func(t *testing.T) {
		runTest(t, nil)
	})

	t.Run("with invalid event error", func(t *testing.T) {
		runTest(t, errors.New("successful test"))
	})
}
