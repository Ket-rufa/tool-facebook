package crawl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/yourname/tool-facebook/internal/accounts"
	"github.com/yourname/tool-facebook/internal/actiontest"
	"github.com/yourname/tool-facebook/internal/fbdata"
)

type CrawlHandler struct {
	accStore *accounts.JSONStore
	fbStore  *fbdata.JSONStore
}

func NewCrawlHandler(accStore *accounts.JSONStore, fbStore *fbdata.JSONStore) *CrawlHandler {
	return &CrawlHandler{
		accStore: accStore,
		fbStore:  fbStore,
	}
}

// Request and Response types
type CrawlRequest struct {
	TargetID  string `json:"target_id"`
	Type      string `json:"type"` // "profile" or "group"
	Limit     int    `json:"limit"`
	AccountID string `json:"account_id"`
	Cursor    string `json:"cursor"` // For pagination
}

type CrawlPostEntity struct {
	PostID    string   `json:"post_id"`
	Author    string   `json:"author"`
	Content   string   `json:"content"`
	Time      string   `json:"time"`
	Reactions string   `json:"reactions"`
	Comments  string   `json:"comments"`
	MediaURLs []string `json:"media_urls"`
}

type CrawlResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []CrawlPostEntity `json:"data"`
}

// fetchPostsViaGraphQL uses Facebook's Comet API to fetch timeline posts
func (h *CrawlHandler) fetchPostsViaGraphQL(crawlReq CrawlRequest, cookie string) ([]CrawlPostEntity, string, error) {
	targetID := crawlReq.TargetID
	fmt.Printf("[CRAWL] Fetching posts via GraphQL for ID: %s (Cursor: %s)\n", targetID, crawlReq.Cursor)

	sd, err := actiontest.FetchSessionData(cookie)
	if err != nil {
		fmt.Printf("[CRAWL] Failed to fetch session tokens: %v\n", err)
		return nil, "", err
	}

	fbDtsg := sd.DTSG
	lsd := sd.LSD

	// Use doc_id for ProfileCometTimelineFeedRefetchQuery
	docID := "26445493535117508"

	vars := map[string]interface{}{
		"count":                         3,
		"cursor":                        nil,
		"feedLocation":                  "TIMELINE",
		"feedbackSource":                0,
		"focusCommentID":                nil,
		"memorializedSplitTimeFilter":   nil,
		"omitPinnedPost":                true,
		"postedBy":                      nil,
		"privacy":                       nil,
		"privacySelectorRenderLocation": "COMET_STREAM",
		"referringStoryRenderLocation":  nil,
		"renderLocation":                "timeline",
		"scale":                         1,
		"stream_count":                  1,
		"taggedInOnly":                  nil,
		"trackingCode":                  nil,
		"useDefaultActor":               false,
		"id":                            targetID,
		"__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider":            true,
		"__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider":   true,
		"__relay_internal__pv__CometFeedStory_enable_post_permalink_white_space_clickrelayprovider": false,
		"__relay_internal__pv__CometUFICommentActionLinksRewriteEnabledrelayprovider": false,
		"__relay_internal__pv__CometUFICommentAvatarStickerAnimatedImagerelayprovider": false,
		"__relay_internal__pv__IsWorkUserrelayprovider":                             false,
		"__relay_internal__pv__TestPilotShouldIncludeDemoAdUseCaserelayprovider":            false,
		"__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider": true,
		"__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider":      true,
		"__relay_internal__pv__CometImmersivePhotoCanUserDisable3DMotionrelayprovider":     false,
		"__relay_internal__pv__WorkCometIsEmployeeGKProviderrelayprovider":           false,
		"__relay_internal__pv__IsMergQAPollsrelayprovider":                             false,
		"__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider": true,
		"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider":          false,
		"__relay_internal__pv__CometUFICommentAutoTranslationTyperelayprovider":        "ORIGINAL",
		"__relay_internal__pv__CometUFIShareActionMigrationrelayprovider":            true,
		"__relay_internal__pv__CometUFISingleLineUFIrelayprovider":                   true,
		"__relay_internal__pv__CometUFI_dedicated_comment_routable_dialog_gkrelayprovider": true,
		"__relay_internal__pv__FBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider":    true,
		"__relay_internal__pv__GroupsCometGYSJFeedItemHeightrelayprovider":               206,
		"__relay_internal__pv__ShouldEnableBakedInTextStoriesrelayprovider":             false,
		"__relay_internal__pv__StoriesShouldIncludeFbNotesrelayprovider":                 false,
	}
	if crawlReq.Cursor != "" {
		vars["cursor"] = crawlReq.Cursor
	}

	variablesJSON, _ := json.Marshal(vars)
	actorID := actiontest.ActorIDFromCookie(cookie)

	fd := url.Values{}
	fd.Set("av", actorID)
	fd.Set("__user", actorID)
	fd.Set("__a", "1")
	fd.Set("fb_dtsg", fbDtsg)
	fd.Set("fb_api_caller_class", "RelayModern")
	fd.Set("fb_api_req_friendly_name", "ProfileCometTimelineFeedRefetchQuery")
	fd.Set("variables", string(variablesJSON))
	fd.Set("doc_id", docID)

	hReq, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(fd.Encode()))

	cleanCookie := actiontest.SanitizeCookie(cookie)
	hReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	hReq.Header.Set("Cookie", cleanCookie)
	hReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	hReq.Header.Set("X-FB-LSD", lsd)
	hReq.Header.Set("X-ASBD-ID", "129477")
	hReq.Header.Set("X-FB-Friendly-Name", "ProfileCometTimelineFeedRefetchQuery")
	hReq.Header.Set("Sec-Fetch-Dest", "empty")
	hReq.Header.Set("Sec-Fetch-Mode", "cors")
	hReq.Header.Set("Sec-Fetch-Site", "same-origin")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(hReq)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	_ = os.WriteFile("data/last_graphql_debug.json", body, 0644)

	bodyStr := string(body)
	return h.parseGraphQLBody(bodyStr)
}

// parseGraphQLBody là hàm dùng chung để parse response GraphQL của Facebook
// cho cả profile lẫn page — cấu trúc JSON hoàn toàn giống nhau
func (h *CrawlHandler) parseGraphQLBody(bodyStr string) ([]CrawlPostEntity, string, error) {
	lines := strings.Split(bodyStr, "\n")

	getMap := func(m interface{}) map[string]interface{} {
		if val, ok := m.(map[string]interface{}); ok { return val }
		return nil
	}
	getList := func(m interface{}) []interface{} {
		if val, ok := m.([]interface{}); ok { return val }
		return nil
	}

	var dataNodes []map[string]interface{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" { continue }
		line = strings.TrimPrefix(line, "for (;;);")
		
		var chunk interface{}
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		if m := getMap(chunk); m != nil {
			if d := getMap(m["data"]); d != nil {
				dataNodes = append(dataNodes, d)
			} else if p := getMap(m["payload"]); p != nil {
				if d := getMap(p["data"]); d != nil {
					dataNodes = append(dataNodes, d)
				}
			}
		} else if l := getList(chunk); l != nil {
			for _, item := range l {
				if m := getMap(item); m != nil {
					if d := getMap(m["data"]); d != nil {
						dataNodes = append(dataNodes, d)
					}
				}
			}
		}
	}

	fmt.Printf("[DEBUG] Parsed %d chunks, found %d valid data nodes\n", len(lines), len(dataNodes))

	var results []CrawlPostEntity
	processedIDs := make(map[string]bool)

	var timelineCursor string
	var findEdges func(m interface{}) []interface{}
	findEdges = func(m interface{}) []interface{} {
		var edges []interface{}
		switch v := m.(type) {
		case map[string]interface{}:
			keys := []string{"timeline_list_feed_units", "timeline_feed_units", "timeline_feed_units_by_location", "news_feed", "edge_comet_viewer_main_feed"}
			for _, k := range keys {
				if u := getMap(v[k]); u != nil {
					if e := getList(u["edges"]); e != nil { edges = append(edges, e...) }
					if pi := getMap(u["page_info"]); pi != nil {
						if hasNext, ok := pi["has_next_page"].(bool); ok && hasNext {
							if c, ok := pi["end_cursor"].(string); ok && c != "" {
								timelineCursor = c
							}
						}
					}
				}
			}
			for _, child := range v {
				if e := findEdges(child); e != nil { edges = append(edges, e...) }
			}
		case []interface{}:
			for _, child := range v {
				if e := findEdges(child); e != nil { edges = append(edges, e...) }
			}
		}
		return edges
	}

	var findCreationTime func(m interface{}) float64
	findCreationTime = func(m interface{}) float64 {
		switch v := m.(type) {
		case map[string]interface{}:
			if ct, ok := v["creation_time"].(float64); ok { return ct }
			for _, child := range v {
				if f := findCreationTime(child); f > 0 { return f }
			}
		case []interface{}:
			for _, child := range v {
				if f := findCreationTime(child); f > 0 { return f }
			}
		}
		return 0
	}

	var findReactionCount func(m interface{}) float64
	findReactionCount = func(m interface{}) float64 {
		switch v := m.(type) {
		case map[string]interface{}:
			if rc, ok := v["reaction_count"]; ok {
				if rcMap, ok2 := rc.(map[string]interface{}); ok2 {
					if count, ok3 := rcMap["count"].(float64); ok3 { return count }
				}
				if rcFl, ok2 := rc.(float64); ok2 { return rcFl }
			}
			for _, child := range v {
				if f := findReactionCount(child); f >= 0 { return f }
			}
		case []interface{}:
			for _, child := range v {
				if f := findReactionCount(child); f >= 0 { return f }
			}
		}
		return -1
	}

	var findCommentCount func(m interface{}) float64
	findCommentCount = func(m interface{}) float64 {
		switch v := m.(type) {
		case map[string]interface{}:
			if cri, ok := v["comment_rendering_instance"]; ok {
				if criMap, ok2 := cri.(map[string]interface{}); ok2 {
					if commentsMap := getMap(criMap["comments"]); commentsMap != nil {
						if tc, ok3 := commentsMap["total_count"].(float64); ok3 { return tc }
					}
				}
			}
			for _, child := range v {
				if f := findCommentCount(child); f >= 0 { return f }
			}
		case []interface{}:
			for _, child := range v {
				if f := findCommentCount(child); f >= 0 { return f }
			}
		}
		return -1
	}

	var findText func(m interface{}) string
	findText = func(m interface{}) string {
		switch v := m.(type) {
		case map[string]interface{}:
			if msg, ok := v["message"]; ok {
				if msgMap, ok2 := msg.(map[string]interface{}); ok2 {
					if t, ok3 := msgMap["text"].(string); ok3 && t != "" { return t }
				}
			}
			for _, child := range v {
				if s := findText(child); s != "" { return s }
			}
		case []interface{}:
			for _, child := range v {
				if s := findText(child); s != "" { return s }
			}
		}
		return ""
	}

	var findCursor func(m interface{}) string
	findCursor = func(m interface{}) string {
		switch v := m.(type) {
		case map[string]interface{}:
			if pi, ok := v["page_info"]; ok {
				if piMap, ok2 := pi.(map[string]interface{}); ok2 {
					if hasNext, ok3 := piMap["has_next_page"].(bool); ok3 && hasNext {
						if c, ok3 := piMap["end_cursor"].(string); ok3 && c != "" { 
							return c 
						}
					}
				}
			}
			if c, ok := v["cursor"].(string); ok && c != "" && len(c) > 10 { 
				return c 
			}
			
			for _, child := range v {
				if c := findCursor(child); c != "" { return c }
			}
		case []interface{}:
			for i := len(v) - 1; i >= 0; i-- {
				if c := findCursor(v[i]); c != "" { return c }
			}
		}
		return ""
	}

	nextCursor := ""
	for _, data := range dataNodes {
		edges := findEdges(data)
		fmt.Printf("[DEBUG] Found %d edges\n", len(edges))
		
		if timelineCursor != "" {
			nextCursor = timelineCursor
		} else if cur := findCursor(data); cur != "" && nextCursor == "" {
			nextCursor = cur
		}
		
		for _, edgeItem := range edges {
			edgeMap := getMap(edgeItem)
			if edgeMap == nil { continue }
			edgeNode := getMap(edgeMap["node"])
			if edgeNode == nil { continue }

			postID, _ := edgeNode["post_id"].(string)
			if postID == "" { postID, _ = edgeNode["id"].(string) }
			if postID == "" || processedIDs[postID] { continue }

			text := findText(edgeNode)
			ct := findCreationTime(edgeNode)
			rc := findReactionCount(edgeNode)
			cc := findCommentCount(edgeNode)

			reactions, comments := "0", "0"
			if rc >= 0 { reactions = fmt.Sprintf("%.0f", rc) }
			if cc >= 0 { comments = fmt.Sprintf("%.0f", cc) }

			if text != "" || postID != "" {
				processedIDs[postID] = true
				if ct > 0 {
					fmt.Printf("[DEBUG] >> Found Post: %s (len:%d, R:%s, C:%s)\n", postID, len(text), reactions, comments)
					p := CrawlPostEntity{
						PostID:    postID,
						Content:   text,
						Time:      time.Unix(int64(ct), 0).Format("02/01/2006 15:04"),
						Reactions: reactions,
						Comments:  comments,
						MediaURLs: h.ExtractMedia(edgeNode),
					}
					if p.Content == "" && len(p.MediaURLs) > 0 {
						p.Content = fmt.Sprintf("(%d Hình ảnh)", len(p.MediaURLs))
					}
					results = append(results, p)
				}
			}
		}
	}
	fmt.Printf("[DEBUG] Extracted Next Cursor: %s\n", nextCursor)
	return results, nextCursor, nil
}

// fetchPostsViaGraphQLPage fetches posts from a Fanpage using CometProfilePostsTabFeedPaginationQuery
func (h *CrawlHandler) fetchPostsViaGraphQLPage(crawlReq CrawlRequest, cookie string) ([]CrawlPostEntity, string, error) {
	targetID := crawlReq.TargetID
	fmt.Printf("[CRAWL] Fetching PAGE posts for Page ID: %s (Cursor: %s)\n", targetID, crawlReq.Cursor)

	sd, err := actiontest.FetchSessionData(cookie)
	if err != nil {
		return nil, "", err
	}

	docID := "8731949750167478"

	vars := map[string]interface{}{
		"count":           3,
		"cursor":          nil,
		"feedLocation":    "PAGE_TIMELINE",
		"feedbackSource":  0,
		"focusCommentID":  nil,
		"renderLocation":  "timeline",
		"scale":           1,
		"useDefaultActor": false,
		"id":              targetID,
		"__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider":          true,
		"__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider": true,
		"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider":      false,
		"__relay_internal__pv__CometUFIShareActionMigrationrelayprovider":          true,
		"__relay_internal__pv__IsWorkUserrelayprovider":                            false,
	}
	if crawlReq.Cursor != "" {
		vars["cursor"] = crawlReq.Cursor
	}

	variablesJSON, _ := json.Marshal(vars)
	actorID := actiontest.ActorIDFromCookie(cookie)

	fd := url.Values{}
	fd.Set("av", actorID)
	fd.Set("__user", actorID)
	fd.Set("__a", "1")
	fd.Set("fb_dtsg", sd.DTSG)
	fd.Set("fb_api_caller_class", "RelayModern")
	fd.Set("fb_api_req_friendly_name", "CometProfilePostsTabFeedPaginationQuery")
	fd.Set("variables", string(variablesJSON))
	fd.Set("doc_id", docID)
	fd.Set("lsd", sd.LSD)
	fd.Set("server_timestamps", "true")

	hReq, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(fd.Encode()))
	hReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	hReq.Header.Set("Cookie", actiontest.SanitizeCookie(cookie))
	hReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	hReq.Header.Set("X-FB-LSD", sd.LSD)
	hReq.Header.Set("X-ASBD-ID", "129477")
	hReq.Header.Set("X-FB-Friendly-Name", "CometProfilePostsTabFeedPaginationQuery")
	hReq.Header.Set("Sec-Fetch-Dest", "empty")
	hReq.Header.Set("Sec-Fetch-Mode", "cors")
	hReq.Header.Set("Sec-Fetch-Site", "same-origin")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(hReq)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	_ = os.WriteFile("data/last_graphql_page_debug.json", body, 0644)
	fmt.Printf("[PAGE] Response status: %d, size: %d bytes\n", resp.StatusCode, len(body))
	if len(body) < 1000 {
		fmt.Printf("[PAGE] Raw response: %s\n", string(body))
	}

	return h.parseGraphQLBody(string(body))
}

func (h *CrawlHandler) ExtractTextFromComet(cometSections map[string]interface{}) string {
	if cometSections == nil { return "" }
	content := h.getMap(cometSections["content"])
	story := h.getMap(content["story"])
	if story == nil { return "" }
	message := h.getMap(story["message"])
	if message == nil { message = h.getMap(h.getMap(story["comet_sections"])["message"]) }
	text, _ := message["text"].(string)
	return text
}

func (h *CrawlHandler) ExtractMedia(m interface{}) []string {
	var urls []string
	seen := make(map[string]bool)

	blacklist := []string{"p40x40", "p32x32", "p24x24", "p50x50", "p16x16", "s50x50", "s32x32", "s40x40", "_s50x", "_s40x", "_s32x", "t39.30808"}

	addURL := func(uri string) {
		if !seen[uri] && strings.Contains(uri, "scontent") {
			for _, b := range blacklist {
				if strings.Contains(uri, b) {
					return
				}
			}
			urls = append(urls, uri)
			seen[uri] = true
		}
	}

	var findInMedia func(data interface{})
	findInMedia = func(data interface{}) {
		switch v := data.(type) {
		case map[string]interface{}:
			if uri, ok := v["uri"].(string); ok {
				addURL(uri)
			}
			// Chỉ đi sâu vào các key có thể chứa media thật sự
			for _, key := range []string{"image", "large_image", "photo_image", "full_image", "media", "thumbnail", "attachments", "style_infos", "story", "attached_story"} {
				if child, ok := v[key]; ok {
					findInMedia(child)
				}
			}
		case []interface{}:
			for _, val := range v {
				findInMedia(val)
			}
		}
	}

	// Chỉ bắt đầu từ các key an toàn ở root của edge node
	if node, ok := m.(map[string]interface{}); ok {
		for _, key := range []string{"attachments", "comet_sections", "attached_story", "media"} {
			if child, ok := node[key]; ok {
				findInMedia(child)
			}
		}
	}

	return urls
}

func (h *CrawlHandler) getMap(m interface{}) map[string]interface{} {
	if val, ok := m.(map[string]interface{}); ok { return val }
	return nil
}

func (h *CrawlHandler) getList(m interface{}) []interface{} {
	if val, ok := m.([]interface{}); ok { return val }
	return nil
}

// RunCrawl is the entry point for the scraping task
func (h *CrawlHandler) RunCrawl(req CrawlRequest) CrawlResponse {
	if req.AccountID == "" {
		return CrawlResponse{Success: false, Message: "Cần chọn tài khoản."}
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	profile, err := h.accStore.Get(req.AccountID)
	if err != nil || profile == nil {
		return CrawlResponse{Success: false, Message: "Tài khoản không tồn tại."}
	}

	cookie := actiontest.SanitizeCookie(profile.Cookie)
	if cookie == "" {
		fbInfo, _ := h.fbStore.Get(req.AccountID)
		if fbInfo != nil {
			cookie = actiontest.SanitizeCookie(fbInfo.Info.Cookie)
		}
	}

	if cookie == "" {
		return CrawlResponse{Success: false, Message: "Lỗi: Không tìm thấy Cookie. Vui lòng cập nhật Cookie."}
	}

	finalTarget := req.TargetID
	if finalTarget == "" {
		finalTarget = actiontest.ActorIDFromCookie(cookie)
		if finalTarget == "" {
			finalTarget = profile.AccountID
		}
	}

	// Nếu target là username (không phải số) → resolve sang numeric ID
	// vì GraphQL của Facebook cần numeric ID
	if finalTarget != "" {
		resolved := h.resolveProfileID(finalTarget, cookie)
		if resolved != "" && resolved != finalTarget {
			fmt.Printf("[CRAWL] Resolved '%s' → numeric ID: %s\n", finalTarget, resolved)
			finalTarget = resolved
		}
	}

	crawlReq := CrawlRequest{
		TargetID: finalTarget,
		Limit:    req.Limit,
		Type:     req.Type,
	}

	var allResults []CrawlPostEntity
	seenIDs := make(map[string]bool)
	cursor := ""
	for len(allResults) < req.Limit {
		crawlReq.Cursor = cursor
		var gPosts []CrawlPostEntity
		var nextCur string
		var err error
		
		if req.Type == "page" {
			gPosts, nextCur, err = h.fetchPostsViaGraphQLPage(crawlReq, cookie)
		} else {
			gPosts, nextCur, err = h.fetchPostsViaGraphQL(crawlReq, cookie)
		}
		if err != nil {
			fmt.Printf("[CRAWL] GraphQL Error: %v\n", err)
			break
		}
		
		addedInThisChunk := 0
		for _, p := range gPosts {
			if !seenIDs[p.PostID] && p.PostID != "" {
				seenIDs[p.PostID] = true
				allResults = append(allResults, p)
				addedInThisChunk++
			}
		}
		
		if nextCur == "" || nextCur == cursor || addedInThisChunk == 0 {
			break
		}
		cursor = nextCur
	}

	if len(allResults) > req.Limit {
		allResults = allResults[:req.Limit]
	}

	if len(allResults) == 0 {
		fmt.Printf("[CRAWL] GraphQL provided 0 results. Falling back to HTML...\n")
		htmlPosts, _ := h.fetchFromHTMLStrategies(crawlReq, cookie)
		for _, hp := range htmlPosts {
			if len(allResults) >= req.Limit {
				break
			}
			if !seenIDs[hp.PostID] && hp.PostID != "" {
				seenIDs[hp.PostID] = true
				allResults = append(allResults, hp)
			}
		}
	}

	if len(allResults) == 0 {
		return CrawlResponse{Success: false, Message: "Không tìm thấy bài viết nào."}
	}

	return CrawlResponse{
		Success: true,
		Message: fmt.Sprintf("✅ Thành công! Đã tìm thấy %d bài viết.", len(allResults)),
		Data:    allResults,
	}
}

func (h *CrawlHandler) fetchFromHTMLStrategies(req CrawlRequest, cookie string) ([]CrawlPostEntity, error) {
	fmt.Println("[CRAWL] Fetching from mbasic Fallback (Safe Mode)...")
	targetURL := ""
	isNumeric := true
	for _, ch := range req.TargetID {
		if ch < '0' || ch > '9' {
			isNumeric = false
			break
		}
	}

	if isNumeric {
		targetURL = fmt.Sprintf("https://mbasic.facebook.com/profile.php?id=%s&v=timeline", req.TargetID)
	} else {
		targetURL = fmt.Sprintf("https://mbasic.facebook.com/%s?v=timeline", req.TargetID)
	}

	time.Sleep(1500 * time.Millisecond)

	operaMiniUA := "Opera/9.80 (J2ME/MIDP; Opera Mini/9.80 (S60; SymbOS; Opera Mobi/23.348; U; en) Presto/2.5.25 Version/10.54"
	status, body, _, err := h.fetchHTML(targetURL, cookie, operaMiniUA)
	if err != nil {
		return nil, err
	}

	_ = os.WriteFile("data/last_crawl_debug.html", []byte(body), 0644)
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d", status)
	}

	return h.ExtractPostsFromHTML(body), nil
}

func (h *CrawlHandler) fetchHTML(targetURL, cookie, ua string) (int, string, string, error) {
	req, _ := http.NewRequest("GET", targetURL, nil)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", ua)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), resp.Request.URL.String(), nil
}


