package main

import (
	"context"
	"fmt"
	"strings"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		trimmed = "Developer"
	}
	return fmt.Sprintf("Xin chao %s, ung dung Go + Wails + Vue da san sang.", trimmed)
}

// PickMediaFiles opens a native file picker and returns absolute file paths.
func (a *App) PickMediaFiles() ([]string, error) {
	if a.ctx == nil {
		return nil, fmt.Errorf("app context is not initialized")
	}
	paths, err := wruntime.OpenMultipleFilesDialog(a.ctx, wruntime.OpenDialogOptions{
		Title: "Chon anh/video tu may",
		Filters: []wruntime.FileFilter{
			{
				DisplayName: "Media Files",
				Pattern:     "*.jpg;*.jpeg;*.png;*.gif;*.webp;*.bmp;*.mp4;*.mov;*.avi;*.mkv;*.webm",
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}
