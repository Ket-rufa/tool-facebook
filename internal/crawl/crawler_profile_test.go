package crawl

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractAboutTokensFromDebug(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "data", "last_ProfileCometAboutAppSectionQuery_debug.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("debug fixture not found, skipping")
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read debug fixture: %v", err)
	}

	h := &CrawlHandler{}
	tokens := h.extractAboutTokens(body)

	if tokens.SectionToken == "" {
		t.Log("Note: SectionToken is empty in this fixture")
	}
}

func TestParseProfileDataSupportsVietnameseCollections(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"node": map[string]interface{}{
				"id":         "999", // Distractor ID
				"name":       "OpenAI Distractor",
				"__typename": "Page",
			},
			"user": map[string]interface{}{
				"id":     "100014491935602",
				"name":   "Hai Yen",
				"gender": "FEMALE",
				"all_collections": map[string]interface{}{
					"nodes": []interface{}{
						map[string]interface{}{
							"title": map[string]interface{}{"text": "Công việc"},
							"style_renderer": map[string]interface{}{
								"renderer": map[string]interface{}{
									"items": []interface{}{
										map[string]interface{}{
											"title": map[string]interface{}{"text": "OpenAI"},
										},
									},
								},
							},
						},
						map[string]interface{}{
							"title": map[string]interface{}{"text": "Học vấn"},
							"style_renderer": map[string]interface{}{
								"renderer": map[string]interface{}{
									"items": []interface{}{
										map[string]interface{}{"title": map[string]interface{}{"text": "Đại học A"}},
									},
								},
							},
						},
						map[string]interface{}{
							"title": map[string]interface{}{"text": "Thành phố hiện tại"},
							"style_renderer": map[string]interface{}{
								"renderer": map[string]interface{}{
									"items": []interface{}{
										map[string]interface{}{"title": map[string]interface{}{"text": "Hà Nội"}},
									},
								},
							},
						},
						map[string]interface{}{
							"title": map[string]interface{}{"text": "Quê quán"},
							"style_renderer": map[string]interface{}{
								"renderer": map[string]interface{}{
									"items": []interface{}{
										map[string]interface{}{"title": map[string]interface{}{"text": "Nam Định"}},
									},
								},
							},
						},
						map[string]interface{}{
							"title": map[string]interface{}{"text": "Tình trạng mối quan hệ"},
							"style_renderer": map[string]interface{}{
								"renderer": map[string]interface{}{
									"items": []interface{}{
										map[string]interface{}{"title": map[string]interface{}{"text": "Độc thân"}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	h := &CrawlHandler{}
	info := &ProfileInfo{UID: "100014491935602"}
	h.parseProfileData(data, info)

	if info.Name != "Hai Yen" {
		t.Fatalf("expected name, got %q", info.Name)
	}
	if info.Gender != "FEMALE" {
		t.Fatalf("expected gender, got %q", info.Gender)
	}
	if len(info.Work) != 1 || info.Work[0] != "OpenAI" {
		t.Fatalf("expected work to be parsed, got %#v", info.Work)
	}
	if len(info.Education) != 1 || info.Education[0] != "Đại học A" {
		t.Fatalf("expected education to be parsed, got %#v", info.Education)
	}
	if info.CurrentCity != "Hà Nội" {
		t.Fatalf("expected current city, got %q", info.CurrentCity)
	}
	if info.Hometown != "Nam Định" {
		t.Fatalf("expected hometown, got %q", info.Hometown)
	}
	if info.Relationship != "Độc thân" {
		t.Fatalf("expected relationship, got %q", info.Relationship)
	}
}
