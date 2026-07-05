package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"lets_config/config"
)

//go:embed all:frontend/dist
var frontendAssets embed.FS

//go:embed all:assets
var bundledAssets embed.FS

//go:embed wails.json
var wailsJSON []byte

// AppVersion is the current application version, read from wails.json.
var AppVersion string

func init() {
	var v struct {
		Info struct {
			ProductVersion string `json:"productVersion"`
		} `json:"Info"`
	}
	if err := json.Unmarshal(wailsJSON, &v); err == nil {
		AppVersion = v.Info.ProductVersion
	}
	if AppVersion == "" {
		AppVersion = "0.0.0"
	}
}

// AssetsDir returns a usable assets directory path.
// In dev mode (local assets/ exists) it returns "assets" directly.
// In production it extracts embedded assets to a persistent cache dir.
func AssetsDir() string {
	local := "assets"
	if info, err := os.Stat(local); err == nil && info.IsDir() {
		return local
	}

	// Production: extract to %APPDATA%\unieditdept\foil\assets\
	cacheDir := filepath.Join(os.Getenv("APPDATA"), "unieditdept", "foil", "assets")

	// Check if already fully extracted by verifying key files
	keyFiles := []string{
		"foil-example.apk",
		"apktool.jar",
		filepath.Join("jre-minimal", "bin", "java.exe"),
	}
	allExist := true
	for _, kf := range keyFiles {
		if _, err := os.Stat(filepath.Join(cacheDir, kf)); err != nil {
			allExist = false
			break
		}
	}
	if allExist {
		return cacheDir
	}

	// Extract embedded assets
	os.RemoveAll(cacheDir)
	if err := extractEmbedFS(bundledAssets, "assets", cacheDir); err != nil {
		// Fallback: try to use local assets/ anyway
		return local
	}
	return cacheDir
}

// extractEmbedFS copies all files from an embedded filesystem path to destDir.
func extractEmbedFS(fsys embed.FS, embedPath, destDir string) error {
	return fs.WalkDir(fsys, embedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(embedPath, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}

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
		Width:     780,
		Height:    720,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: frontendAssets,
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
		os.Exit(1)
	}
}
