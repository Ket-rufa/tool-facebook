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
	"github.com/yourname/tool-facebook/internal/crawl"
	"github.com/yourname/tool-facebook/internal/fbdata"
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

	// Khởi tạo FB Data Store
	fbStore, err := fbdata.NewJSONStore(dataDir)
	if err != nil {
		log.Fatalf("Không thể khởi tạo fb data store: %v", err)
	}

	// Tạo AccountService dùng store đã tạo và truyền thêm fbStore để đồng bộ
	accountService := accounts.NewAccountServiceWithStore(store, fbStore)

	// Khởi tạo action handler với cùng store và fbStore
	sessionHandler := actiontest.NewSessionHandler()
	configStore, err := actiontest.NewConfigStore(dataDir)
	if err != nil {
		log.Fatalf("Khong the khoi tao action config store: %v", err)
	}
	actionHandler := actiontest.NewActionHandlerWithStoreAndConfig(sessionHandler, store, fbStore, configStore)

	crawlHandler := crawl.NewCrawlHandler(store, fbStore)

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
			crawlHandler,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
