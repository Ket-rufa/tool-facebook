package main

import (
	"context"
	"fmt"
	"strings"
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
