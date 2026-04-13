package tests

import (
	"testing"

	app "github.com/yourname/tool-facebook/internal/app"
)

func TestNormalizeProfileInput(t *testing.T) {
	result := app.NormalizeProfileInput(app.ProfileInput{URL: " https://facebook.com/example ", Language: ""})
	if result.Language != "vi" {
		t.Fatalf("expected default language vi, got %s", result.Language)
	}
	if result.Status != "ready" {
		t.Fatalf("expected status ready, got %s", result.Status)
	}
}
