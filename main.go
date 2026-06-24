package main

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"lets_config/config"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Config file lives alongside the WebView user data directory
	configPath := filepath.Join(os.Getenv("APPDATA"), "unieditdept", "foil", "config.toml")
	cfgManager, err := config.NewManager(configPath)
	if err != nil {
		println("Error: failed to init config manager:", err.Error())
		os.Exit(1)
	}

	// Create an instance of the app structure
	app := NewApp(cfgManager)

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "Foil",
		Width:     1280,
		Height:    720,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Acrylic,

			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               filepath.Join(os.Getenv("APPDATA"), "unieditdept", "foil"),
			Theme:                             windows.SystemDefault,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
