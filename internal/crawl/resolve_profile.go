package crawl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/yourname/tool-facebook/internal/actiontest"
)

// resolveProfileID chuyển username/vanity URL sang numeric profile ID
func (h *CrawlHandler) resolveProfileID(target, cookie string) string {
	isNum := true
	for _, ch := range target {
		if ch < '0' || ch > '9' { isNum = false; break }
	}
	if isNum { return target }

	fmt.Printf("[RESOLVE] Resolving '%s' to numeric ID...\n", target)

	// Chiến lược 1: Dùng GraphQL CometUserLookupQuery để tìm ID từ vanity URL
	if id := h.resolveViaGraphQL(target, cookie); id != "" {
		fmt.Printf("[RESOLVE] GraphQL lookup succeeded: %s → %s\n", target, id)
		return id
	}

	// Chiến lược 2: Scrape www.facebook.com và tìm trong JSON data
	if id := h.resolveViaHTMLScrape(target, cookie); id != "" {
		fmt.Printf("[RESOLVE] HTML scrape succeeded: %s → %s\n", target, id)
		return id
	}

	fmt.Printf("[RESOLVE] All strategies failed, using original: %s\n", target)
	return target
}

// resolveViaGraphQL dùng Facebook GraphQL để lookup user/page bằng vanity URL
func (h *CrawlHandler) resolveViaGraphQL(vanity, cookie string) string {
	sd, err := actiontest.FetchSessionData(cookie)
	if err != nil {
		return ""
	}

	// Query để lookup profile/page bằng vanity URL
	// Dùng UserProfileCometRootQuery hoặc CometUserByVanityURLRootQuery
	vars := map[string]interface{}{
		"vanityUrl": vanity,
		"scale":     1,
	}
	varsJSON, _ := json.Marshal(vars)
	actorID := actiontest.ActorIDFromCookie(cookie)

	fd := url.Values{}
	fd.Set("av", actorID)
	fd.Set("__user", actorID)
	fd.Set("__a", "1")
	fd.Set("fb_dtsg", sd.DTSG)
	fd.Set("fb_api_caller_class", "RelayModern")
	fd.Set("fb_api_req_friendly_name", "CometUserByVanityURLRootQuery")
	fd.Set("variables", string(varsJSON))
	fd.Set("doc_id", "4381972608565505") // CometUserByVanityURLRootQuery or similar
	fd.Set("lsd", sd.LSD)
	fd.Set("server_timestamps", "true")

	req, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(fd.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", actiontest.SanitizeCookie(cookie))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-FB-LSD", sd.LSD)
	req.Header.Set("X-ASBD-ID", "129477")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Referer", "https://www.facebook.com/")

	client := &http.Client{Timeout: 15 * http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout}
	_ = client
	resp, err := (&http.Client{Timeout: 15 * 1e9}).Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("[RESOLVE] GraphQL response size: %d bytes\n", len(body))

	// Parse multi-line response
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimPrefix(strings.TrimSpace(line), "for (;;);")
		if line == "" { continue }
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil { continue }

		// Traverse tìm user ID
		if id := extractIDFromGraphQL(data); id != "" {
			return id
		}
	}
	return ""
}

func extractIDFromGraphQL(data interface{}) string {
	switch v := data.(type) {
	case map[string]interface{}:
		// Tìm các keys đặc trưng cho profile/page node
		for _, key := range []string{"user", "page", "profile", "node", "actor"} {
			if node, ok := v[key].(map[string]interface{}); ok {
				if id, ok := node["id"].(string); ok && len(id) >= 6 {
					return id
				}
			}
		}
		for _, val := range v {
			if id := extractIDFromGraphQL(val); id != "" {
				return id
			}
		}
	case []interface{}:
		for _, item := range v {
			if id := extractIDFromGraphQL(item); id != "" {
				return id
			}
		}
	}
	return ""
}

// resolveViaHTMLScrape scrape www.facebook.com với full browser headers
func (h *CrawlHandler) resolveViaHTMLScrape(vanity, cookie string) string {
	profileURL := "https://www.facebook.com/" + vanity
	fmt.Printf("[RESOLVE] HTML scrape: %s\n", profileURL)

	req, _ := http.NewRequest("GET", profileURL, nil)
	req.Header.Set("Cookie", actiontest.SanitizeCookie(cookie))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", "https://www.facebook.com/")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := (&http.Client{Timeout: 20 * 1e9}).Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	// Lưu debug
	limit := len(body)
	if limit > 200000 { limit = 200000 }
	_ = os.WriteFile("data/last_resolve_debug.html", bodyBytes[:limit], 0644)

	// Các pattern tìm trong JSON inline của Facebook desktop
	// Lưu ý: tránh match ID của viewer (61572157625165)
	viewerID := ""
	if m := regexp.MustCompile(`"USER_ID":"(\d+)"`).FindStringSubmatch(body); len(m) > 1 {
		viewerID = m[1]
	}
	if viewerID == "" {
		if m := regexp.MustCompile(`"userID":"(\d+)"`).FindStringSubmatch(body); len(m) > 1 {
			viewerID = m[1]
		}
	}
	fmt.Printf("[RESOLVE] Viewer ID detected as: %s\n", viewerID)

	// Tìm patterns đặc trưng cho profile được xem
	patterns := []struct{ label, regex string }{
		{"entity_id (first)",  `"entity_id":"(\d+)"`},
		{"timeline_owner",     `timeline_owner.*?"id":"(\d+)"`},
		{"subject_id",         `"subject_id":"(\d+)"`},
		{"profile node id",    `"__typename":"(?:User|UserProfile|Profile)"[^}]*"id":"(\d+)"`},
		{"page node id",       `"__typename":"Page"[^}]*"id":"(\d+)"`},
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p.regex)
		matches := re.FindAllStringSubmatch(body, -1)
		for _, m := range matches {
			if len(m) > 1 && len(m[1]) >= 6 && m[1] != viewerID {
				fmt.Printf("[RESOLVE][html] Found via '%s': %s\n", p.label, m[1])
				return m[1]
			}
		}
	}
	return ""
}
