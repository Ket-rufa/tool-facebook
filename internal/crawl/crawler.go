package crawl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
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

	// Danh sách suffix/path của ảnh nhỏ (avatar, icon reaction, emoji)
	// QUAN TRỌNG: t39.30808-1 là avatar/profile pic, t39.30808-6 là ảnh bài viết!
	blacklist := []string{
		"p40x40", "p32x32", "p24x24", "p50x50", "p16x16",
		"s50x50", "s32x32", "s40x40", "s16x16",
		"_s50x", "_s40x", "_s32x", "_s16x",
		"t39.30808-1", // avatar/profile picture (KHÔNG phải t39.30808 để không block ảnh bài viết -6, -9...)
		"/emoji.php",  // emoji images
		"p180x180", "p100x100", "p75x75", "p60x60",
	}

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

	// Thực hiện DFS toàn bộ subtree — an toàn vì chỉ được gọi
	// từ các key attachment/media đã được cô lập khỏi header/feedback
	var fullDFS func(data interface{}, depth int)
	fullDFS = func(data interface{}, depth int) {
		if depth > 20 { return } // tránh đệ quy vô hạn
		switch v := data.(type) {
		case map[string]interface{}:
			// Nếu node này là User/Page/Actor thì bỏ qua (tránh avatar người comment)
			if typename, ok := v["__typename"].(string); ok {
				switch typename {
				case "User", "Page", "Group", "CometAnimatedReactionIconRenderer",
					"CometUFIActorPhotoType", "Actor":
					return
				}
			}
			if uri, ok := v["uri"].(string); ok {
				addURL(uri)
			}
			// Bỏ qua các key chứa thông tin người dùng / phản ứng
			skipKeys := map[string]bool{
				"actors": true, "actor": true,
				"profile_picture": true, "profile_picture_depth_0": true,
				"profile_picture_depth_1": true, "actor_photo": true,
				"icon_image": true, "extensions": true,
				"comet_ufi_reaction_icon_renderer": true,
				"vote_attachments": true,
				"feedback": true, "story_ufi_container": true,
			}
			for key, val := range v {
				if !skipKeys[key] {
					fullDFS(val, depth+1)
				}
			}
		case []interface{}:
			for _, val := range v {
				fullDFS(val, depth+1)
			}
		}
	}

	// Chỉ bắt đầu DFS từ các zone an toàn chứa attachment của bài viết
	// (không bắt đầu từ root để tránh lấy ảnh từ feedback/reactions)
	entryKeys := []string{
		"attachments", "comet_sections", "attached_story",
		"media", "photo", "photos", "subattachments",
	}
	if node, ok := m.(map[string]interface{}); ok {
		for _, key := range entryKeys {
			if child, ok := node[key]; ok {
				fullDFS(child, 0)
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

// ── Quét bạn bè ──────────────────────────────────────────────────────────────

// FriendEntity là thông tin một người bạn
type FriendEntity struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// FriendCrawlResponse là kết quả quét bạn bè
type FriendCrawlResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	TotalCount int            `json:"total_count"`
	Friends    []FriendEntity `json:"friends"`
}


// FetchFriends quét danh sách bạn bè. Ưu tiên GraphQL nếu có Doc ID, ngược lại dùng mbasic sau đó tự động thử các Doc ID phổ biến.
func (h *CrawlHandler) FetchFriends(accountID string, targetID string, docID string) FriendCrawlResponse {
	if accountID == "" {
		return FriendCrawlResponse{Success: false, Message: "Cần chọn tài khoản."}
	}

	profile, err := h.accStore.Get(accountID)
	if err != nil || profile == nil {
		return FriendCrawlResponse{Success: false, Message: "Tài khoản không tồn tại."}
	}

	cookie := actiontest.SanitizeCookie(profile.Cookie)
	if cookie == "" {
		fbInfo, _ := h.fbStore.Get(accountID)
		if fbInfo != nil {
			cookie = actiontest.SanitizeCookie(fbInfo.Info.Cookie)
		}
	}
	if cookie == "" {
		return FriendCrawlResponse{Success: false, Message: "Không tìm thấy Cookie."}
	}

	// Xác định viewerID (người thực hiện)
	viewerID := actiontest.ActorIDFromCookie(cookie)

	// Xác định target: Nếu trống thì là chính mình
	finalTarget := targetID
	isSelf := false
	if finalTarget == "" || finalTarget == viewerID || finalTarget == profile.AccountID {
		finalTarget = viewerID
		if finalTarget == "" { finalTarget = profile.AccountID }
		isSelf = true
	}

	// Resolve target (nếu là username/vanity URL)
	resolved := h.resolveProfileID(finalTarget, cookie)
	if resolved != "" {
		finalTarget = resolved
	}

	fmt.Printf("[FRIENDS] Bắt đầu quét cho Target: %s (qua tài khoản: %s)\n", finalTarget, profile.AccountID)

	// ── TRƯỜNG HỢP 1: CÓ DOC ID THỦ CÔNG ──────────────────────────────────────
	if docID != "" {
		return h.fetchFriendsGraphQL(finalTarget, cookie, docID)
	}

	// ── TRƯỜNG HỢP 2: THỬ MBASIC TRƯỚC ──
	mbasicFriends, totalCount := h.fetchFriendsMbasic(finalTarget, cookie)
	if len(mbasicFriends) > 0 {
		return FriendCrawlResponse{
			Success: true, TotalCount: totalCount, Friends: mbasicFriends,
			Message: fmt.Sprintf("✅ Quét thành công %d bạn bè qua mbasic.", len(mbasicFriends)),
		}
	}

	// ── TRƯỜNG HỢP 3: MBASIC THẤT BẠI -> TỰ ĐỘNG THỬ DANH SÁCH DOC ID PHỔ BIẾN ──
	fmt.Println("[FRIENDS] mbasic thất bại, bắt đầu thử các Doc ID phổ biến...")
	fallbacks := []string{
		"7114815461879024", // ProfileCometFriendsListQuery (Stable)
		"25686001222718131",
		"10156054341515257",
		"8856555194367018",
		"7086851601336040",
		"6793541504033754",
	}

	// Ưu tiên nạp ID phù hợp nhất lên đầu danh sách thử nghiệm
	if isSelf {
		// Nếu là chính mình, ưu tiên ID đã biết là hoạt động tốt cho self
		fallbacks = append([]string{"26206414195674994"}, fallbacks...)
	} else {
		// Nếu là người khác, ưu tiên ID Non-Self bạn vừa tìm được
		fallbacks = append([]string{"26565284836436780"}, fallbacks...)
	}

	for _, id := range fallbacks {
		fmt.Printf("[FRIENDS] Thử Fallback Doc ID: %s cho Target: %s\n", id, finalTarget)
		resp := h.fetchFriendsGraphQL(finalTarget, cookie, id)
		
		// Chỉ chấp nhận nếu tìm thấy ít nhất 1 bạn bè THẬT (không phải rác hệ thống)
		if resp.Success && len(resp.Friends) > 0 {
			resp.Message = fmt.Sprintf("✅ Quét thành công %d bạn bè (Tự động dùng ID: %s).", len(resp.Friends), id)
			return resp
		}
		fmt.Printf("[FRIENDS] Doc ID %s không trả về kết quả hợp lệ.\n", id)
	}

	return FriendCrawlResponse{
		Success: false,
		Message: "Không thể lấy danh sách bạn bè. Hãy thử F5 lại trang FB trên trình duyệt, bấm tab 'Tất cả bạn bè' rồi lấy Doc ID dán vào đây.",
	}
}

// fetchFriendsGraphQL thực hiện quét qua API GraphQL
func (h *CrawlHandler) fetchFriendsGraphQL(targetID, cookie, docID string) FriendCrawlResponse {
	sd, err := actiontest.FetchSessionData(cookie)
	if err != nil {
		return FriendCrawlResponse{Success: false, Message: "Lỗi session: " + err.Error()}
	}

	if sd.DTSG == "" {
		return FriendCrawlResponse{Success: false, Message: "Lỗi: Không lấy được token bảo mật (fb_dtsg). Có thể Cookie đã hết hạn hoặc bị Facebook chặn."}
	}

	var allFriends []FriendEntity
	seenUIDs := make(map[string]bool)
	cursor := ""
	totalCount := 0

	for page := 1; page <= 50; page++ {
		vars := map[string]interface{}{}

		// Nếu là mã "Non-Self" mới (2026), dùng bộ tham số đầy đủ
		if docID == "26565284836436780" {
			vars = map[string]interface{}{
				"id":      targetID,
				"count":   20,
				"cursor":  cursor,
				"scale":   1,
				"search":  nil,
				"__relay_internal__pv__FBProfile_enable_perf_improv_gkrelayprovider": true,
			}
		} else {
			// Với các mã cũ hoặc mã quét "Self", dùng bộ tham số tối giản và ổn định
			vars = map[string]interface{}{
				"id":     targetID,
				"count":  50,
				"scale":  1,
			}
			if cursor != "" {
				vars["cursor"] = cursor
			} else {
				vars["cursor"] = nil
			}
		}

		variablesJSON, _ := json.Marshal(vars)
		actorID := actiontest.ActorIDFromCookie(cookie)

		fd := url.Values{}
		fd.Set("av", actorID)
		fd.Set("__user", actorID)
		fd.Set("__a", "1")
		fd.Set("fb_dtsg", sd.DTSG)
		fd.Set("fb_api_caller_class", "RelayModern")
		fd.Set("fb_api_req_friendly_name", "ProfileCometAppCollectionNonSelfFriendsListRendererPaginationQuery")
		fd.Set("variables", string(variablesJSON))
		fd.Set("doc_id", docID)
		if sd.LSD != "" {
			fd.Set("lsd", sd.LSD)
		}

		req, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(fd.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
		
		if sd.LSD != "" {
			req.Header.Set("X-FB-LSD", sd.LSD)
		}
		req.Header.Set("X-ASBD-ID", "129477")
		req.Header.Set("X-FB-Friendly-Name", "ProfileCometAppCollectionNonSelfFriendsListRendererPaginationQuery")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil { break }
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		bodyStr := string(body)
		if strings.Contains(bodyStr, "\"errors\"") && !strings.Contains(bodyStr, "edges") {
			fmt.Printf("[FRIENDS-GQL] Facebook báo lỗi: %s\n", bodyStr)
			break
		}

		friends, nextCursor, count := h.parseFriendsGraphQLInternal(bodyStr, seenUIDs)
		if count > 0 && totalCount == 0 { totalCount = count }
		
		if len(friends) == 0 { break }

		allFriends = append(allFriends, friends...)

		fmt.Printf("[FRIENDS-GQL] Trang %d: +%d bạn, tổng=%d\n", page, len(friends), len(allFriends))
		if nextCursor == "" || nextCursor == cursor { break }
		cursor = nextCursor
		time.Sleep(300 * time.Millisecond)
	}

	if len(allFriends) == 0 {
		return FriendCrawlResponse{Success: false}
	}

	return FriendCrawlResponse{
		Success: true, TotalCount: totalCount, Friends: allFriends,
	}
}

// fetchFriendsMbasic thu thập danh sách bạn bè qua mbasic.facebook.com
func (h *CrawlHandler) fetchFriendsMbasic(uid, cookie string) ([]FriendEntity, int) {
	// Dùng User-Agent Chrome chuẩn để tránh bị chặn mbasic
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	seen := make(map[string]bool)
	var allFriends []FriendEntity
	totalCount := 0
	startIndex := 0

	for page := 1; page <= 100; page++ {
		var targetURL string
		if startIndex == 0 {
			targetURL = fmt.Sprintf("https://mbasic.facebook.com/profile.php?id=%s&v=friends", uid)
		} else {
			targetURL = fmt.Sprintf("https://mbasic.facebook.com/profile.php?id=%s&v=friends&startindex=%d", uid, startIndex)
		}

		req, _ := http.NewRequest("GET", targetURL, nil)
		req.Header.Set("Cookie", actiontest.SanitizeCookie(cookie))
		req.Header.Set("User-Agent", ua)
		req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8")

		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Do(req)
		if err != nil { break }
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		body := string(bodyBytes)

		if strings.Contains(body, "You must log in first") || strings.Contains(body, "login.php") {
			fmt.Println("[FRIENDS-MBASIC] Bị yêu cầu đăng nhập.")
			break
		}

		friends, nextIndex, count := h.parseFriendsMbasicHTML(body, seen)
		if count > 0 && totalCount == 0 { totalCount = count }
		
		if len(friends) == 0 { break }

		allFriends = append(allFriends, friends...)
		for _, f := range friends { seen[f.UID] = true }

		fmt.Printf("[FRIENDS-MBASIC] Trang %d: +%d bạn, tổng=%d\n", page, len(friends), len(allFriends))
		if nextIndex <= 0 || nextIndex == startIndex { break }
		startIndex = nextIndex
		time.Sleep(500 * time.Millisecond)
	}

	return allFriends, totalCount
}

func (h *CrawlHandler) parseFriendsGraphQLInternal(bodyStr string, seenUIDs map[string]bool) ([]FriendEntity, string, int) {
	var friends []FriendEntity
	nextCursor := ""
	totalCount := 0

	for _, line := range strings.Split(bodyStr, "\n") {
		line = strings.TrimPrefix(strings.TrimSpace(line), "for (;;);")
		if line == "" { continue }
		var chunk map[string]interface{}
		if err := json.Unmarshal([]byte(line), &chunk); err != nil { continue }
		
		var findEdges func(obj interface{})
		findEdges = func(obj interface{}) {
			m, ok := obj.(map[string]interface{})
			if !ok { return }
			
			// Hỗ trợ cả friends.edges và s_friends.edges
			for _, key := range []string{"edges", "friends", "all_friends", "node"} {
				if v, ok := m[key]; ok {
					if key == "edges" {
						if edges, ok := v.([]interface{}); ok {
							if tc, ok := m["total_count"].(float64); ok { totalCount = int(tc) }
							if pi, ok := m["page_info"].(map[string]interface{}); ok {
								if hasNext, _ := pi["has_next_page"].(bool); hasNext {
									if c, _ := pi["end_cursor"].(string); c != "" { nextCursor = c }
								}
							}
							for _, e := range edges {
								edge, ok := e.(map[string]interface{})
								if !ok { continue }
								node, ok := edge["node"].(map[string]interface{})
								if !ok { continue }
								uid, _ := node["id"].(string)
								name, _ := node["name"].(string)

								// BỘ LỌC THÔNG MINH:
								// 1. UID phải là dãy số và dài từ 6 ký tự trở lên.
								// 2. Không phải là các link hệ thống như profile.php, v.v.
								isNumericID := true
								if len(uid) < 6 { isNumericID = false }
								for _, ch := range uid {
									if ch < '0' || ch > '9' { isNumericID = false; break }
								}

								if !isNumericID || name == "" || seenUIDs[uid] { continue }
								
								// Loại bỏ các cụm từ rác
								lowerName := strings.ToLower(name)
								if strings.Contains(lowerName, "chỉnh sửa") || strings.Contains(lowerName, "edit profile") {
									continue
								}

								avatar := ""
								if pp, ok := node["profile_picture"].(map[string]interface{}); ok {
									avatar, _ = pp["uri"].(string)
								} else if pp, ok := node["profile_photo"].(map[string]interface{}); ok {
									// Thử thêm profile_photo cho một số cấu trúc khác
									avatar, _ = pp["uri"].(string)
								}

								friends = append(friends, FriendEntity{UID: uid, Name: name, Avatar: avatar})
								seenUIDs[uid] = true
							}
						}
					} else {
						findEdges(v)
					}
				}
			}
			for k, v := range m {
				if k != "edges" && k != "friends" { findEdges(v) }
			}
		}
		findEdges(chunk)
	}
	return friends, nextCursor, totalCount
}

func (h *CrawlHandler) parseFriendsMbasicHTML(body string, seen map[string]bool) ([]FriendEntity, int, int) {
	var friends []FriendEntity
	nextIndex := -1
	totalCount := 0

	reCount := regexp.MustCompile(`(\d[\d,\.]*)\s+(?:bạn bè|friends)`)
	if m := reCount.FindStringSubmatch(body); m != nil {
		numStr := strings.ReplaceAll(strings.ReplaceAll(m[1], ",", ""), ".", "")
		totalCount, _ = strconv.Atoi(numStr)
	}

	// Cải tiến regex để bắt được nhiều định dạng link mbasic hơn
	reFriend := regexp.MustCompile(`(?i)<img\s[^>]*src="([^"]+)"[^>]*>.*?<a\s+href="/([^"?&]+)(?:\?[^"]*)?"[^>]*>([^<]+)</a>`)
	matches := reFriend.FindAllStringSubmatch(body, -1)
	
	for _, m := range matches {
		img, href, name := m[1], m[2], strings.TrimSpace(m[3])
		uid := href
		if strings.Contains(href, "profile.php?id=") {
			uid = strings.Split(strings.Split(href, "id=")[1], "&")[0]
		}

		// BỘ LỌC THÔNG MINH CHO MBASIC:
		isNumericID := true
		if len(uid) < 6 { isNumericID = false }
		for _, ch := range uid {
			if ch < '0' || ch > '9' { isNumericID = false; break }
		}

		lowerName := strings.ToLower(name)
		if !isNumericID || name == "" || seen[uid] || 
		   strings.Contains(lowerName, "chỉnh sửa") || 
		   strings.Contains(lowerName, "edit profile") ||
		   strings.Contains(lowerName, "xem thêm") || 
		   strings.Contains(lowerName, "facebook") { 
			continue 
		}

		friends = append(friends, FriendEntity{UID: uid, Name: name, Avatar: img})
		seen[uid] = true
	}

	reNext := regexp.MustCompile(`(?i)startindex=(\d+)`)
	if m := reNext.FindStringSubmatch(body); m != nil {
		nextIndex, _ = strconv.Atoi(m[1])
	}

	return friends, nextIndex, totalCount
}




