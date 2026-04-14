package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/yourname/tool-facebook/internal/actiontest"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	sessionHandler := actiontest.NewSessionHandler()
	actionHandler := actiontest.NewActionHandler(sessionHandler)

	err := wails.Run(&options.App{
		Title:  "Tool Facebook",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 20, G: 24, B: 28, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			sessionHandler,
			actionHandler,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
