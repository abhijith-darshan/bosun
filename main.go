package main

import (
	"bosun/pkg"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:dashboards/dist/apps/bosun-ui
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := pkg.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "bosun",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Start,
		Bind: []interface{}{
			app,
		},
		WindowStartState: options.Maximised,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
