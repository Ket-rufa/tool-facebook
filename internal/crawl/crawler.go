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

type ProfileInfo struct {
	UID            string   `json:"uid"`
	Name           string   `json:"name"`
	Gender         string   `json:"gender"`
	Birthday       string   `json:"birthday"`
	CurrentCity    string   `json:"current_city"`
	Hometown       string   `json:"hometown"`
	Relationship   string   `json:"relationship"`
	Education      []string `json:"education"`
	Work           []string `json:"work"`
	Followers      string   `json:"followers"`
	Bio            string   `json:"bio"`
	ProfilePicture string   `json:"profile_picture"`
}

type profileAboutTokens struct {
	SectionToken      string
	RawSectionToken   string
	AppSectionFeedKey string
	AboutURL          string
	InfoAllURL        string
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
		"__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider":                           true,
		"__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider":                  true,
		"__relay_internal__pv__CometFeedStory_enable_post_permalink_white_space_clickrelayprovider": false,
		"__relay_internal__pv__CometUFICommentActionLinksRewriteEnabledrelayprovider":               false,
		"__relay_internal__pv__CometUFICommentAvatarStickerAnimatedImagerelayprovider":              false,
		"__relay_internal__pv__IsWorkUserrelayprovider":                                             false,
		"__relay_internal__pv__TestPilotShouldIncludeDemoAdUseCaserelayprovider":                    false,
		"__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider":          true,
		"__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider":               true,
		"__relay_internal__pv__CometImmersivePhotoCanUserDisable3DMotionrelayprovider":              false,
		"__relay_internal__pv__WorkCometIsEmployeeGKProviderrelayprovider":                          false,
		"__relay_internal__pv__IsMergQAPollsrelayprovider":                                          false,
		"__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider":           true,
		"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider":                       false,
		"__relay_internal__pv__CometUFICommentAutoTranslationTyperelayprovider":                     "ORIGINAL",
		"__relay_internal__pv__CometUFIShareActionMigrationrelayprovider":                           true,
		"__relay_internal__pv__CometUFISingleLineUFIrelayprovider":                                  true,
		"__relay_internal__pv__CometUFI_dedicated_comment_routable_dialog_gkrelayprovider":          true,
		"__relay_internal__pv__FBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider":              true,
		"__relay_internal__pv__GroupsCometGYSJFeedItemHeightrelayprovider":                          206,
		"__relay_internal__pv__ShouldEnableBakedInTextStoriesrelayprovider":                         false,
		"__relay_internal__pv__StoriesShouldIncludeFbNotesrelayprovider":                            false,
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
		if val, ok := m.(map[string]interface{}); ok {
			return val
		}
		return nil
	}
	getList := func(m interface{}) []interface{} {
		if val, ok := m.([]interface{}); ok {
			return val
		}
		return nil
	}

	var dataNodes []map[string]interface{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
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
					if e := getList(u["edges"]); e != nil {
						edges = append(edges, e...)
					}
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
				if e := findEdges(child); e != nil {
					edges = append(edges, e...)
				}
			}
		case []interface{}:
			for _, child := range v {
				if e := findEdges(child); e != nil {
					edges = append(edges, e...)
				}
			}
		}
		return edges
	}

	var findCreationTime func(m interface{}) float64
	findCreationTime = func(m interface{}) float64 {
		switch v := m.(type) {
		case map[string]interface{}:
			if ct, ok := v["creation_time"].(float64); ok {
				return ct
			}
			for _, child := range v {
				if f := findCreationTime(child); f > 0 {
					return f
				}
			}
		case []interface{}:
			for _, child := range v {
				if f := findCreationTime(child); f > 0 {
					return f
				}
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
					if count, ok3 := rcMap["count"].(float64); ok3 {
						return count
					}
				}
				if rcFl, ok2 := rc.(float64); ok2 {
					return rcFl
				}
			}
			for _, child := range v {
				if f := findReactionCount(child); f >= 0 {
					return f
				}
			}
		case []interface{}:
			for _, child := range v {
				if f := findReactionCount(child); f >= 0 {
					return f
				}
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
						if tc, ok3 := commentsMap["total_count"].(float64); ok3 {
							return tc
						}
					}
				}
			}
			for _, child := range v {
				if f := findCommentCount(child); f >= 0 {
					return f
				}
			}
		case []interface{}:
			for _, child := range v {
				if f := findCommentCount(child); f >= 0 {
					return f
				}
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
					if t, ok3 := msgMap["text"].(string); ok3 && t != "" {
						return t
					}
				}
			}
			for _, child := range v {
				if s := findText(child); s != "" {
					return s
				}
			}
		case []interface{}:
			for _, child := range v {
				if s := findText(child); s != "" {
					return s
				}
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
				if c := findCursor(child); c != "" {
					return c
				}
			}
		case []interface{}:
			for i := len(v) - 1; i >= 0; i-- {
				if c := findCursor(v[i]); c != "" {
					return c
				}
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
			if edgeMap == nil {
				continue
			}
			edgeNode := getMap(edgeMap["node"])
			if edgeNode == nil {
				continue
			}

			postID, _ := edgeNode["post_id"].(string)
			if postID == "" {
				postID, _ = edgeNode["id"].(string)
			}
			if postID == "" || processedIDs[postID] {
				continue
			}

			text := findText(edgeNode)
			ct := findCreationTime(edgeNode)
			rc := findReactionCount(edgeNode)
			cc := findCommentCount(edgeNode)

			reactions, comments := "0", "0"
			if rc >= 0 {
				reactions = fmt.Sprintf("%.0f", rc)
			}
			if cc >= 0 {
				comments = fmt.Sprintf("%.0f", cc)
			}

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
	if cometSections == nil {
		return ""
	}
	content := h.getMap(cometSections["content"])
	story := h.getMap(content["story"])
	if story == nil {
		return ""
	}
	message := h.getMap(story["message"])
	if message == nil {
		message = h.getMap(h.getMap(story["comet_sections"])["message"])
	}
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
		if depth > 20 {
			return
		} // tránh đệ quy vô hạn
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
				"vote_attachments":                 true,
				"feedback":                         true, "story_ufi_container": true,
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
	if val, ok := m.(map[string]interface{}); ok {
		return val
	}
	return nil
}

func (h *CrawlHandler) getList(m interface{}) []interface{} {
	if val, ok := m.([]interface{}); ok {
		return val
	}
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

// ExportedResolveUID is a wrapper for resolveProfileID to be called from the frontend
func (h *CrawlHandler) ExportedResolveUID(target, accountID string) string {
	profile, err := h.accStore.Get(accountID)
	if err != nil || profile == nil {
		return target
	}
	cookie := actiontest.SanitizeCookie(profile.Cookie)
	return h.resolveProfileID(target, cookie)
}

// ScanProfileInfo fetches detailed profile information using GraphQL doc_id
func (h *CrawlHandler) ScanProfileInfo(uid, cookie string) (*ProfileInfo, error) {
	fmt.Printf("[SCAN] Scanning profile info for UID: %s\n", uid)

	sd, err := actiontest.FetchSessionData(cookie)
	if err != nil {
		return nil, err
	}

	actorID := actiontest.ActorIDFromCookie(cookie)
	info := &ProfileInfo{UID: uid}
	return h.scanProfileInfoV2(uid, sd, actorID, cookie, info)

	// BƯỚC 1: Thử ProfileCometAboutAppSectionQuery (Chi tiết) - Cập nhật theo doc_id mới nhất
	docID := "24801675826197206"
	vars := map[string]interface{}{
		"appSectionFeedKey": fmt.Sprintf("ProfileCometAppSectionFeed_timeline_nav_app_sections__%s", uid),
		"collectionToken":   nil,
		"pageID":            uid,
		"rawSectionToken":   nil,
		"scale":             1,
		"sectionToken":      nil,
		"showReactions":     true,
		"userID":            uid,
		"__relay_internal__pv__FBProfile_enable_perf_improv_gkrelayprovider":                                               true,
		"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider":                                              false,
		"__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider":                                 true,
		"__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider":                                  true,
		"__relay_internal__pv__FBUnifiedVideoMediaContentContainer_comet_reels_video_footer_defer_loading_gkrelayprovider": false,
		"__relay_internal__pv__ShouldEnableBakedInTextUnifiedVideorelayprovider":                                           false,
		"__relay_internal__pv__FBUnifiedVideoMediaFooter_comet_enable_reels_ads_gkrelayprovider":                           true,
		"__relay_internal__pv__FBUnifiedVideoMediaFooter_enable_meta_ai_pill_gkrelayprovider":                              true,
		"__relay_internal__pv__FBUnifiedVideoMediaFooter_enable_group_character_ai_info_pill_gkrelayprovider":              true,
		"__relay_internal__pv__FBUnifiedVideoMediaFooter_enable_video_augment_pills_gkrelayprovider":                       false,
		"__relay_internal__pv__FBUnifiedVideoFeedbackBar_comet_reels_save_button_gkrelayprovider":                          false,
		"__relay_internal__pv__usePushPipEngagementCounts_comet_video_document_picture_in_picture_gkrelayprovider":         false,
		"__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider":                                      true,
		"__relay_internal__pv__FBUnifiedVideoMenu_fb_reels_ranking_debug_tool_gkrelayprovider":                             false,
	}

	body, err := h.executeGraphQL(docID, "ProfileCometAboutAppSectionQuery", vars, sd, actorID, cookie)
	if err == nil && !strings.Contains(string(body), `"error":1357054`) {
		h.parseProfileLines(body, info)
	}

	// BƯỚC 2: Nếu chưa có tên (do truy vấn trên lỗi), thử ProfileCometHeaderQuery (Cơ bản nhưng ổn định)
	if info.Name == "" {
		fmt.Printf("[SCAN] About query failed or empty, falling back to Header query for UID: %s\n", uid)
		headerDocID := "5380494448574168" // ProfileCometHeaderQuery
		headerVars := map[string]interface{}{
			"id":    uid,
			"scale": 1,
		}
		headerBody, err := h.executeGraphQL(headerDocID, "ProfileCometHeaderQuery", headerVars, sd, actorID, cookie)
		if err == nil {
			h.parseProfileLines(headerBody, info)
		}
	}

	if info.Name != "" {
		h.SaveProfileToText(info)
		return info, nil
	}

	return nil, fmt.Errorf("không thể lấy thông tin profile qua GraphQL")
}

func (h *CrawlHandler) scanProfileInfoV2(uid string, sd actiontest.SessionData, actorID, cookie string, info *ProfileInfo) (*ProfileInfo, error) {
	// Bước 1: probe About query để lấy đúng sectionToken/rawSectionToken từ response mới.
	docID := "24801675826197206"
	probeBody, err := h.executeGraphQL(docID, "ProfileCometAboutAppSectionQuery", h.buildProfileAboutVars(uid, profileAboutTokens{}), sd, actorID, cookie)
	tokens := profileAboutTokens{}
	if err == nil && !strings.Contains(string(probeBody), `"error":1357054`) {
		tokens = h.extractAboutTokens(probeBody)
		h.parseProfileLines(probeBody, info)
	}

	// Bước 2: khi đã có token thật, gọi lại About query với appSectionFeedKey/rawSectionToken đúng.
	if tokens.SectionToken != "" || tokens.RawSectionToken != "" {
		fullBody, err := h.executeGraphQL(docID, "ProfileCometAboutAppSectionQuery", h.buildProfileAboutVars(uid, tokens), sd, actorID, cookie)
		if err == nil && !strings.Contains(string(fullBody), `"error":1357054`) {
			h.parseProfileLines(fullBody, info)
			if refreshed := h.extractAboutTokens(fullBody); refreshed.SectionToken != "" || refreshed.RawSectionToken != "" {
				tokens = refreshed
			}
		}
	}

	// Bước 3: Header query hiện tại đã stale; dùng Timeline query ổn định để bổ sung name/avatar/gender.
	if info.Name == "" || info.ProfilePicture == "" || info.Gender == "" {
		fmt.Printf("[SCAN] About query incomplete, hydrating identity from Timeline query for UID: %s\n", uid)
		h.hydrateProfileFromTimeline(uid, sd, actorID, cookie, info)
	}

	if info.Name == "" && h.hasMeaningfulProfileInfo(info) {
		info.Name = h.deriveProfileFallbackName(uid, tokens)
	}

	h.saveLastProfileScanDebug(info)
	if h.hasMeaningfulProfileInfo(info) {
		h.SaveProfileToText(info)
		return info, nil
	}

	return nil, fmt.Errorf("không thể lấy thông tin profile qua GraphQL")
}

func (h *CrawlHandler) executeGraphQL(docID, friendlyName string, vars map[string]interface{}, sd actiontest.SessionData, actorID, cookie string) ([]byte, error) {
	variablesJSON, _ := json.Marshal(vars)
	fd := url.Values{}
	fd.Set("av", actorID)
	fd.Set("__user", actorID)
	fd.Set("__a", "1")
	fd.Set("fb_dtsg", sd.DTSG)
	fd.Set("jazoest", actiontest.CalcJazoest(sd.DTSG))
	fd.Set("fb_api_caller_class", "RelayModern")
	fd.Set("fb_api_req_friendly_name", friendlyName)
	fd.Set("variables", string(variablesJSON))
	fd.Set("doc_id", docID)
	fd.Set("lsd", sd.LSD)
	fd.Set("server_timestamps", "true")

	req, _ := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(fd.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", actiontest.SanitizeCookie(cookie))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("X-FB-LSD", sd.LSD)
	req.Header.Set("X-ASBD-ID", "129477")
	req.Header.Set("X-FB-Friendly-Name", friendlyName)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Referer", "https://www.facebook.com/")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		_ = os.WriteFile(fmt.Sprintf("data/last_%s_debug.json", friendlyName), body, 0644)
	}
	return body, err
}

func (h *CrawlHandler) parseProfileLines(body []byte, info *ProfileInfo) {
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		line = strings.TrimPrefix(strings.TrimSpace(line), "for (;;);")
		if line == "" {
			continue
		}
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}
		h.parseProfileData(data, info)
	}
}

func (h *CrawlHandler) parseProfileData(data map[string]interface{}, info *ProfileInfo) {
	roots := []map[string]interface{}{
		h.getMap(data["data"]),
		h.getMap(h.getMap(data["payload"])["data"]),
		h.getMap(data["node"]),
	}
	for _, root := range roots {
		if root == nil {
			continue
		}
		h.walkProfileValue(root, info, 0)
	}
	return

	node := h.getMap(data["data"])
	if node == nil {
		node = h.getMap(h.getMap(data["payload"])["data"])
	}
	if node == nil {
		node = h.getMap(data["node"])
	}
	if node == nil {
		return
	}

	user := h.getMap(node["user"])
	if user == nil {
		// Trong Header query, node thường chính là user hoặc có field user trực tiếp
		if _, ok := node["name"].(string); ok {
			user = node
		} else {
			// Thử tìm sâu hơn
			for _, v := range node {
				if m := h.getMap(v); m != nil {
					if _, ok := m["name"].(string); ok {
						user = m
						break
					}
				}
			}
		}
	}

	if user == nil {
		return
	}

	if name, ok := user["name"].(string); ok && info.Name == "" {
		info.Name = name
	}
	if gender, ok := user["gender"].(string); ok && info.Gender == "" {
		info.Gender = gender
	}

	if bday := h.getMap(user["birth_date"]); bday != nil {
		day := bday["day"]
		month := bday["month"]
		year := bday["year"]
		info.Birthday = fmt.Sprintf("%v/%v/%v", day, month, year)
	}

	collections := h.getMap(user["all_collections"])
	if collections != nil {
		nodes := h.getList(collections["nodes"])
		for _, nItem := range nodes {
			n := h.getMap(nItem)
			if n == nil {
				continue
			}
			title := h.getMap(n["title"])
			if title == nil {
				continue
			}
			titleText, _ := title["text"].(string)

			style := h.getMap(n["style_renderer"])
			if style == nil {
				continue
			}
			renderer := h.getMap(style["renderer"])
			if renderer == nil {
				continue
			}

			items := h.getList(renderer["items"])
			for _, iItem := range items {
				item := h.getMap(iItem)
				if item == nil {
					continue
				}
				titleObj := h.getMap(item["title"])
				if titleObj == nil {
					continue
				}
				text, _ := titleObj["text"].(string)

				switch titleText {
				case "Work":
					info.Work = append(info.Work, text)
				case "Education":
					info.Education = append(info.Education, text)
				case "Current City":
					info.CurrentCity = text
				case "Hometown":
					info.Hometown = text
				case "Relationship":
					info.Relationship = text
				}
			}
		}
	}

	if followers := h.getMap(user["followers"]); followers != nil {
		if count, ok := followers["count"].(float64); ok {
			info.Followers = fmt.Sprintf("%.0f", count)
		}
	}

	if bio := h.getMap(user["bio_text"]); bio != nil {
		info.Bio, _ = bio["text"].(string)
	}
}

func (h *CrawlHandler) buildProfileAboutVars(uid string, tokens profileAboutTokens) map[string]interface{} {
	appSectionFeedKey := fmt.Sprintf("ProfileCometAppSectionFeed_timeline_nav_app_sections__%s", uid)
	var rawSectionToken interface{}
	var sectionToken interface{}

	if tokens.RawSectionToken != "" {
		rawSectionToken = tokens.RawSectionToken
		appSectionFeedKey = fmt.Sprintf("ProfileCometAppSectionFeed_timeline_nav_app_sections__%s", tokens.RawSectionToken)
	}
	if tokens.AppSectionFeedKey != "" {
		appSectionFeedKey = tokens.AppSectionFeedKey
	}
	if tokens.SectionToken != "" {
		sectionToken = tokens.SectionToken
	}

	return map[string]interface{}{
		"appSectionFeedKey": appSectionFeedKey,
		"collectionToken":   nil,
		"pageID":            uid,
		"rawSectionToken":   rawSectionToken,
		"scale":             1,
		"sectionToken":      sectionToken,
		"showReactions":     true,
		"userID":            uid,
		"__relay_internal__pv__FBProfile_enable_perf_improv_gkrelayprovider":                                               true,
		"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider":                                              false,
		"__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider":                                 true,
		"__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider":                                  true,
		"__relay_internal__pv__FBUnifiedVideoMediaContentContainer_comet_reels_video_footer_defer_loading_gkrelayprovider": false,
		"__relay_internal__pv__ShouldEnableBakedInTextUnifiedVideorelayprovider":                                           false,
		"__relay_internal__pv__FBUnifiedVideoMediaFooter_comet_enable_reels_ads_gkrelayprovider":                           true,
		"__relay_internal__pv__FBUnifiedVideoMediaFooter_enable_meta_ai_pill_gkrelayprovider":                              true,
		"__relay_internal__pv__FBUnifiedVideoMediaFooter_enable_group_character_ai_info_pill_gkrelayprovider":              true,
		"__relay_internal__pv__FBUnifiedVideoMediaFooter_enable_video_augment_pills_gkrelayprovider":                       false,
		"__relay_internal__pv__FBUnifiedVideoFeedbackBar_comet_reels_save_button_gkrelayprovider":                          false,
		"__relay_internal__pv__usePushPipEngagementCounts_comet_video_document_picture_in_picture_gkrelayprovider":         false,
		"__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider":                                      true,
		"__relay_internal__pv__FBUnifiedVideoMenu_fb_reels_ranking_debug_tool_gkrelayprovider":                             false,
	}
}

func (h *CrawlHandler) extractAboutTokens(body []byte) profileAboutTokens {
	var tokens profileAboutTokens

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		line = strings.TrimPrefix(strings.TrimSpace(line), "for (;;);")
		if line == "" {
			continue
		}

		var chunk map[string]interface{}
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		if root := h.getMap(chunk["data"]); root != nil {
			h.extractAboutTokensFromRoot(root, &tokens)
		}
		if payload := h.getMap(chunk["payload"]); payload != nil {
			if root := h.getMap(payload["data"]); root != nil {
				h.extractAboutTokensFromRoot(root, &tokens)
			}
		}
	}

	if tokens.RawSectionToken != "" && tokens.AppSectionFeedKey == "" {
		tokens.AppSectionFeedKey = fmt.Sprintf("ProfileCometAppSectionFeed_timeline_nav_app_sections__%s", tokens.RawSectionToken)
	}

	return tokens
}

func (h *CrawlHandler) extractAboutTokensFromRoot(root map[string]interface{}, tokens *profileAboutTokens) {
	if user := h.getMap(root["user"]); user != nil {
		if sections := h.getMap(user["about_app_sections"]); sections != nil {
			h.extractAboutTokensFromSectionContainer(sections, tokens, false)
		}
	}
	if sections := h.getMap(root["about_app_sections"]); sections != nil {
		h.extractAboutTokensFromSectionContainer(sections, tokens, false)
	}
	if sections := h.getMap(root["timeline_nav_app_sections"]); sections != nil {
		h.extractAboutTokensFromSectionContainer(sections, tokens, true)
	}
}

func (h *CrawlHandler) extractAboutTokensFromSectionContainer(container map[string]interface{}, tokens *profileAboutTokens, wantCursor bool) {
	for _, section := range h.connectionNodes(container) {
		if !h.isAboutSection(section) {
			continue
		}
		if tokens.SectionToken == "" {
			if id, _ := section["id"].(string); id != "" {
				tokens.SectionToken = id
			}
		}
		if tokens.AboutURL == "" {
			if u, _ := section["url"].(string); u != "" {
				tokens.AboutURL = u
			}
		}
		if tokens.InfoAllURL == "" {
			if u := h.firstCollectionURL(h.getMap(section["nav_collections"])); u != "" {
				tokens.InfoAllURL = u
			}
		}
		if tokens.InfoAllURL == "" {
			if u := h.firstCollectionURL(h.getMap(section["all_collections"])); u != "" {
				tokens.InfoAllURL = u
			}
		}
	}

	if wantCursor {
		for _, edgeItem := range h.getList(container["edges"]) {
			edge := h.getMap(edgeItem)
			if edge == nil {
				continue
			}
			section := h.getMap(edge["node"])
			if section == nil || !h.isAboutSection(section) {
				continue
			}
			if tokens.RawSectionToken == "" {
				if cursor, _ := edge["cursor"].(string); cursor != "" {
					tokens.RawSectionToken = cursor
				}
			}
			break
		}
	}
}

func (h *CrawlHandler) hydrateProfileFromTimeline(uid string, sd actiontest.SessionData, actorID, cookie string, info *ProfileInfo) {
	vars := map[string]interface{}{
		"count":                         1,
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
		"id":                            uid,
		"__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider":                           true,
		"__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider":                  true,
		"__relay_internal__pv__CometFeedStory_enable_post_permalink_white_space_clickrelayprovider": false,
		"__relay_internal__pv__CometUFICommentActionLinksRewriteEnabledrelayprovider":               false,
		"__relay_internal__pv__CometUFICommentAvatarStickerAnimatedImagerelayprovider":              false,
		"__relay_internal__pv__IsWorkUserrelayprovider":                                             false,
		"__relay_internal__pv__TestPilotShouldIncludeDemoAdUseCaserelayprovider":                    false,
		"__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider":          true,
		"__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider":               true,
		"__relay_internal__pv__CometImmersivePhotoCanUserDisable3DMotionrelayprovider":              false,
		"__relay_internal__pv__WorkCometIsEmployeeGKProviderrelayprovider":                          false,
		"__relay_internal__pv__IsMergQAPollsrelayprovider":                                          false,
		"__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider":           true,
		"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider":                       false,
		"__relay_internal__pv__CometUFICommentAutoTranslationTyperelayprovider":                     "ORIGINAL",
		"__relay_internal__pv__CometUFIShareActionMigrationrelayprovider":                           true,
		"__relay_internal__pv__CometUFISingleLineUFIrelayprovider":                                  true,
		"__relay_internal__pv__CometUFI_dedicated_comment_routable_dialog_gkrelayprovider":          true,
		"__relay_internal__pv__FBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider":              true,
		"__relay_internal__pv__GroupsCometGYSJFeedItemHeightrelayprovider":                          206,
		"__relay_internal__pv__ShouldEnableBakedInTextStoriesrelayprovider":                         false,
		"__relay_internal__pv__StoriesShouldIncludeFbNotesrelayprovider":                            false,
	}

	body, err := h.executeGraphQL("26445493535117508", "ProfileCometTimelineFeedRefetchQuery", vars, sd, actorID, cookie)
	if err == nil {
		h.parseProfileLines(body, info)
	}
}

func (h *CrawlHandler) walkProfileValue(value interface{}, info *ProfileInfo, depth int) {
	if value == nil || depth > 20 {
		return
	}

	switch v := value.(type) {
	case map[string]interface{}:
		h.applyProfileEntity(v, info)
		if sections := h.getMap(v["about_app_sections"]); sections != nil {
			h.parseSectionContainer(sections, info)
		}
		if sections := h.getMap(v["timeline_nav_app_sections"]); sections != nil {
			h.parseSectionContainer(sections, info)
		}
		if collections := h.getMap(v["all_collections"]); collections != nil {
			h.parseCollectionsContainer(collections, info)
		}
		for _, child := range v {
			h.walkProfileValue(child, info, depth+1)
		}
	case []interface{}:
		for _, child := range v {
			h.walkProfileValue(child, info, depth+1)
		}
	}
}

func (h *CrawlHandler) applyProfileEntity(node map[string]interface{}, info *ProfileInfo) {
	id, _ := node["id"].(string)
	if info.UID != "" && id != "" && id != info.UID {
		return
	}
	if info.UID != "" && id == "" {
		return
	}

	if name, ok := node["name"].(string); ok && info.Name == "" && !h.isAboutSection(map[string]interface{}{"name": name}) {
		info.Name = name
	}
	if gender, ok := node["gender"].(string); ok && info.Gender == "" {
		info.Gender = gender
	}
	if bday := h.getMap(node["birth_date"]); bday != nil && info.Birthday == "" {
		day := fmt.Sprintf("%v", bday["day"])
		month := fmt.Sprintf("%v", bday["month"])
		year := fmt.Sprintf("%v", bday["year"])
		if day != "<nil>" && month != "<nil>" && year != "<nil>" {
			info.Birthday = fmt.Sprintf("%s/%s/%s", day, month, year)
		}
	}
	if followers := h.getMap(node["followers"]); followers != nil && info.Followers == "" {
		if count, ok := followers["count"].(float64); ok {
			info.Followers = fmt.Sprintf("%.0f", count)
		}
	}
	if bio := h.getMap(node["bio_text"]); bio != nil && info.Bio == "" {
		info.Bio = h.findText(bio)
	}
	if info.CurrentCity == "" {
		info.CurrentCity = h.findText(h.getMap(node["current_city"]))
	}
	if info.Hometown == "" {
		info.Hometown = h.findText(h.getMap(node["hometown"]))
	}
	if info.Relationship == "" {
		info.Relationship = h.findText(h.getMap(node["relationship_status"]))
	}
	if info.ProfilePicture == "" {
		for _, key := range []string{"profile_picture", "profile_picture_depth_0", "profile_picture_depth_1", "profile_photo"} {
			if uri := h.extractURI(h.getMap(node[key])); uri != "" {
				info.ProfilePicture = uri
				break
			}
		}
	}
}

func (h *CrawlHandler) parseSectionContainer(container map[string]interface{}, info *ProfileInfo) {
	for _, section := range h.connectionNodes(container) {
		if collections := h.getMap(section["all_collections"]); collections != nil {
			h.parseCollectionsContainer(collections, info)
		}
	}
}

func (h *CrawlHandler) parseCollectionsContainer(container map[string]interface{}, info *ProfileInfo) {
	for _, collectionNode := range h.connectionNodes(container) {
		h.parseCollectionNode(collectionNode, info)
	}
}

func (h *CrawlHandler) parseCollectionNode(collectionNode map[string]interface{}, info *ProfileInfo) {
	collectionTitle := firstNonEmpty(h.findText(collectionNode["title"]), h.findText(collectionNode["name"]))

	if style := h.getMap(collectionNode["style_renderer"]); style != nil {
		if collection := h.getMap(style["collection"]); collection != nil {
			if collectionTitle == "" {
				collectionTitle = firstNonEmpty(h.findText(collection["title"]), h.findText(collection["name"]))
			}
			h.parseCollectionItems(collectionTitle, collection["pageItems"], info)
			h.parseCollectionItems(collectionTitle, collection["items"], info)
		}
		if renderer := h.getMap(style["renderer"]); renderer != nil {
			h.parseCollectionItems(collectionTitle, renderer["items"], info)
		}
	}

	h.parseCollectionItems(collectionTitle, collectionNode["items"], info)
	h.parseCollectionItems(collectionTitle, collectionNode["pageItems"], info)
}

func (h *CrawlHandler) parseCollectionItems(collectionTitle string, items interface{}, info *ProfileInfo) {
	switch v := items.(type) {
	case []interface{}:
		for _, item := range v {
			h.applyCollectionItem(collectionTitle, h.getMap(item), info)
		}
	case map[string]interface{}:
		for _, edge := range h.getList(v["edges"]) {
			edgeMap := h.getMap(edge)
			if edgeMap == nil {
				continue
			}
			h.applyCollectionItem(collectionTitle, h.getMap(edgeMap["node"]), info)
		}
		for _, node := range h.getList(v["nodes"]) {
			h.applyCollectionItem(collectionTitle, h.getMap(node), info)
		}
	}
}

func (h *CrawlHandler) applyCollectionItem(collectionTitle string, item map[string]interface{}, info *ProfileInfo) {
	if item == nil {
		return
	}

	normalizedCollection := normalizeLabel(collectionTitle)
	label := firstNonEmpty(h.findText(item["title"]), h.findText(item["label"]))
	value := firstNonEmpty(
		h.findText(item["subtitle"]),
		h.findText(item["description"]),
		h.extractPrimaryText(item, collectionTitle),
	)
	if value == "" {
		switch normalizedCollection {
		case "work", "công việc", "education", "học vấn", "current city", "thành phố hiện tại", "hometown", "quê quán", "relationship", "tình trạng mối quan hệ":
			value = label
		}
	}
	if value == "" {
		return
	}

	switch normalizedCollection {
	case "work", "công việc":
		info.Work = appendUnique(info.Work, value)
	case "education", "học vấn":
		info.Education = appendUnique(info.Education, value)
	case "current city", "thành phố hiện tại":
		if info.CurrentCity == "" {
			info.CurrentCity = value
		}
	case "hometown", "quê quán":
		if info.Hometown == "" {
			info.Hometown = value
		}
	case "relationship", "tình trạng mối quan hệ":
		if info.Relationship == "" {
			info.Relationship = value
		}
	case "places lived", "nơi từng sống":
		switch normalizeLabel(label) {
		case "current city", "thành phố hiện tại":
			if info.CurrentCity == "" {
				info.CurrentCity = value
			}
		case "hometown", "quê quán":
			if info.Hometown == "" {
				info.Hometown = value
			}
		}
	case "contact and basic info", "thông tin liên hệ và cơ bản":
		switch normalizeLabel(label) {
		case "birthday", "ngày sinh":
			if info.Birthday == "" {
				info.Birthday = value
			}
		case "relationship", "tình trạng mối quan hệ":
			if info.Relationship == "" {
				info.Relationship = value
			}
		}
	}
}

func (h *CrawlHandler) extractPrimaryText(value interface{}, excludes ...string) string {
	exclude := map[string]struct{}{}
	for _, item := range excludes {
		if item = normalizeLabel(item); item != "" {
			exclude[item] = struct{}{}
		}
	}

	var texts []string
	var walk func(interface{})
	walk = func(node interface{}) {
		switch v := node.(type) {
		case map[string]interface{}:
			for key, child := range v {
				switch key {
				case "text", "name":
					if s, ok := child.(string); ok {
						texts = append(texts, s)
					}
				default:
					walk(child)
				}
			}
		case []interface{}:
			for _, child := range v {
				walk(child)
			}
		}
	}
	walk(value)

	for _, text := range texts {
		norm := normalizeLabel(text)
		if norm == "" {
			continue
		}
		if _, found := exclude[norm]; found {
			continue
		}
		if norm == "about" || norm == "giới thiệu" {
			continue
		}
		return strings.TrimSpace(text)
	}
	return ""
}

func (h *CrawlHandler) extractURI(node map[string]interface{}) string {
	if node == nil {
		return ""
	}
	if uri, _ := node["uri"].(string); uri != "" {
		return uri
	}
	for _, child := range node {
		if nested := h.getMap(child); nested != nil {
			if uri := h.extractURI(nested); uri != "" {
				return uri
			}
		}
	}
	return ""
}

func (h *CrawlHandler) findText(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case map[string]interface{}:
		if s, ok := v["text"].(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
		if s, ok := v["name"].(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
		for _, child := range v {
			if text := h.findText(child); text != "" {
				return text
			}
		}
	case []interface{}:
		for _, child := range v {
			if text := h.findText(child); text != "" {
				return text
			}
		}
	}
	return ""
}

func (h *CrawlHandler) isAboutSection(section map[string]interface{}) bool {
	sectionType, _ := section["section_type"].(string)
	if strings.EqualFold(sectionType, "ABOUT") {
		return true
	}
	name, _ := section["name"].(string)
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "about", "giới thiệu":
		return true
	default:
		return false
	}
}

func (h *CrawlHandler) connectionNodes(container map[string]interface{}) []map[string]interface{} {
	if container == nil {
		return nil
	}

	var nodes []map[string]interface{}
	for _, item := range h.getList(container["nodes"]) {
		if m := h.getMap(item); m != nil {
			nodes = append(nodes, m)
		}
	}
	for _, item := range h.getList(container["edges"]) {
		if edge := h.getMap(item); edge != nil {
			if node := h.getMap(edge["node"]); node != nil {
				nodes = append(nodes, node)
			}
		}
	}
	return nodes
}

func (h *CrawlHandler) firstCollectionURL(container map[string]interface{}) string {
	for _, item := range h.connectionNodes(container) {
		if u, _ := item["url"].(string); u != "" {
			return u
		}
		if style := h.getMap(item["style_renderer"]); style != nil {
			if collection := h.getMap(style["collection"]); collection != nil {
				if u, _ := collection["url"].(string); u != "" {
					return u
				}
			}
		}
	}
	return ""
}

func (h *CrawlHandler) hasMeaningfulProfileInfo(info *ProfileInfo) bool {
	return info.Name != "" ||
		info.Gender != "" ||
		info.Birthday != "" ||
		info.CurrentCity != "" ||
		info.Hometown != "" ||
		info.Relationship != "" ||
		len(info.Work) > 0 ||
		len(info.Education) > 0 ||
		info.Followers != "" ||
		info.Bio != "" ||
		info.ProfilePicture != ""
}

func (h *CrawlHandler) deriveProfileFallbackName(uid string, tokens profileAboutTokens) string {
	for _, raw := range []string{tokens.InfoAllURL, tokens.AboutURL} {
		if raw == "" {
			continue
		}
		parsed, err := url.Parse(raw)
		if err != nil {
			continue
		}
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) == 0 || parts[0] == "" || parts[0] == "profile.php" {
			continue
		}
		return parts[0]
	}
	return uid
}

func (h *CrawlHandler) saveLastProfileScanDebug(info *ProfileInfo) {
	body, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile("data/last_profile_scan_debug.json", body, 0644)
}

func normalizeLabel(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func appendUnique(list []string, value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return list
	}
	for _, existing := range list {
		if existing == value {
			return list
		}
	}
	return append(list, value)
}

func (h *CrawlHandler) SaveProfileToText(info *ProfileInfo) {
	dataDir := "DATA"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		_ = os.MkdirAll(dataDir, 0755)
	}

	filename := fmt.Sprintf("%s/%s.txt", dataDir, info.UID)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UID: %s\n", info.UID))
	sb.WriteString(fmt.Sprintf("Họ tên: %s\n", info.Name))
	sb.WriteString(fmt.Sprintf("Giới tính: %s\n", info.Gender))
	sb.WriteString(fmt.Sprintf("Ngày sinh: %s\n", info.Birthday))
	sb.WriteString(fmt.Sprintf("Thành phố hiện tại: %s\n", info.CurrentCity))
	sb.WriteString(fmt.Sprintf("Quê quán: %s\n", info.Hometown))
	sb.WriteString(fmt.Sprintf("Tình trạng mối quan hệ: %s\n", info.Relationship))

	sb.WriteString("\nCông việc:\n")
	for _, w := range info.Work {
		sb.WriteString(fmt.Sprintf("- %s\n", w))
	}

	sb.WriteString("\nGiáo dục:\n")
	for _, e := range info.Education {
		sb.WriteString(fmt.Sprintf("- %s\n", e))
	}

	sb.WriteString(fmt.Sprintf("\nNgười theo dõi: %s\n", info.Followers))
	sb.WriteString(fmt.Sprintf("Tiểu sử: %s\n", info.Bio))

	_ = os.WriteFile(filename, []byte(sb.String()), 0644)
	fmt.Printf("[SCAN] Saved profile data to %s\n", filename)
}
