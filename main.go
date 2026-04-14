package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/yourname/tool-facebook/internal/accounts"
	"github.com/yourname/tool-facebook/internal/actiontest"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	// Tạo thư mục data cục bộ cạnh file thực thi
	exePath, _ := os.Executable()
	dataDir := filepath.Join(filepath.Dir(exePath), "data")

	// Khởi tạo store trước để chia sẻ giữa các service
	store, err := accounts.NewJSONStore(dataDir)
	if err != nil {
		log.Fatalf("Không thể khởi tạo account store: %v", err)
	}

	// Tạo AccountService dùng store đã tạo
	accountService := accounts.NewAccountServiceWithStore(store)

	// Khởi tạo action handler với cùng store để validate session
	sessionHandler := actiontest.NewSessionHandler()
	actionHandler := actiontest.NewActionHandlerWithStore(sessionHandler, store)

	err = wails.Run(&options.App{
		Title:  "Tool Facebook",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			sessionHandler,
			actionHandler,
			accountService,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
