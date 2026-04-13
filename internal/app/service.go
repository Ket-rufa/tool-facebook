package app

import "strings"

type ProfileInput struct {
	URL      string
	Language string
}

type ProfileResult struct {
	NormalizedURL string `json:"normalized_url"`
	Language      string `json:"language"`
	Status        string `json:"status"`
}

func NormalizeProfileInput(input ProfileInput) ProfileResult {
	url := strings.TrimSpace(input.URL)
	lang := strings.TrimSpace(input.Language)
	if lang == "" {
		lang = "vi"
	}

	return ProfileResult{
		NormalizedURL: url,
		Language:      lang,
		Status:        "ready",
	}
}
