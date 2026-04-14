package actiontest

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/yourname/tool-facebook/internal/accounts"
	"github.com/yourname/tool-facebook/internal/fbdata"
)

// jsonStr trả về chuỗi được JSON encode an toàn (có xử lý ký tự đặc biệt)
func jsonStr(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// generateToken tạo một chuỗi token ngẫu nhiên dạng UUID đơn giản
func generateToken() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		rand.Uint32(),
		rand.Uint32()&0xffff,
		rand.Uint32()&0xffff,
		rand.Uint32()&0xffff,
		rand.Uint64()&0xffffffffffff,
	)
}

// actorIDFromCookie trích xuất UID từ cookie string
func actorIDFromCookie(cookie string) string {
	for _, part := range strings.Split(cookie, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "c_user=") {
			return strings.TrimPrefix(part, "c_user=")
		}
	}
	return ""
}

// extractDtsgFromError trích xuất dtsgToken từ phản hồi lỗi 1357004 của Facebook
// Facebook “tặng” token hợp lệ trong body lỗi, dùng cho lần request tiếp theo
func extractDtsgFromError(body string) string {
	re := regexp.MustCompile(`"dtsgToken"\s*:\s*"([^"]+)"`)
	m := re.FindStringSubmatch(body)
	if len(m) > 1 {
		token := m[1] // GIỮ NGUYÊN toàn bộ, kể cả :3:timestamp
		display := token
		if len(display) > 20 {
			display = display[:20]
		}
		fmt.Printf("[INFO] Trích xuất dtsgToken từ lỗi Facebook: %s\n", display)
		return token
	}
	return ""
}

// SessionData chứa các thông số phiên làm việc lấy từ trang chủ
type SessionData struct {
	LSD   string
	DTSG  string
	Rev   string
	HSI   string
	SpinR string
	SpinT string
}

// fetchSessionData truy cập trang chủ Facebook ngầm để lấy lsd, fb_dtsg và các thông số kỹ thuật khác
func fetchSessionData(cookie string) (data SessionData, err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://www.facebook.com/", nil)
	
	// Headers siêu tàng hình
	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Ch-Ua", `"Not-A.Brand";v="99", "Chromium";v="124", "Google Chrome";v="124"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	// Lấy Revision (__rev)
	revRegex := regexp.MustCompile(`"client_revision":(\d+)`)
	revMatches := revRegex.FindStringSubmatch(bodyStr)
	if len(revMatches) > 1 {
		data.Rev = revMatches[1]
	}

	// Lấy HSI (__hsi)
	hsiRegex := regexp.MustCompile(`"hsi":"(\d+)"`)
	hsiMatches := hsiRegex.FindStringSubmatch(bodyStr)
	if len(hsiMatches) > 1 {
		data.HSI = hsiMatches[1]
	}

	// Lấy Spin R (__spin_r)
	spinRRegex := regexp.MustCompile(`"spin_r":(\d+)`)
	spinRMatches := spinRRegex.FindStringSubmatch(bodyStr)
	if len(spinRMatches) > 1 {
		data.SpinR = spinRMatches[1]
	}

	// Lấy Spin T (__spin_t)
	spinTRegex := regexp.MustCompile(`"spin_t":(\d+)`)
	spinTMatches := spinTRegex.FindStringSubmatch(bodyStr)
	if len(spinTMatches) > 1 {
		data.SpinT = spinTMatches[1]
	}

	// Kiểm tra xem có thực sự đang ở trạng thái đăng nhập không
	isLoggedIn := strings.Contains(bodyStr, `"USER_ID":"`) || strings.Contains(bodyStr, "c_user=") || strings.Contains(bodyStr, actorIDFromCookie(cookie))
	
	if !isLoggedIn {
		fmt.Printf("[WARN] Facebook yêu cầu đăng nhập (hoặc Cookie đã hết hạn) khi lấy LSD.\n")
		// In thử tiêu đề trang để biết đang ở đâu
		titleRegex := regexp.MustCompile(`<title>(.*?)</title>`)
		titleMatch := titleRegex.FindStringSubmatch(bodyStr)
		if len(titleMatch) > 1 {
			fmt.Printf("[DEBUG] Tiêu đề trang nhận được: %s\n", titleMatch[1])
		}
	}

	// Tìm LSD token (Cải tiến Regex linh hoạt hơn)
	lsdRegex := regexp.MustCompile(`"LSD"\s*,\s*\[\s*\]\s*,\s*\{\s*"token"\s*:\s*"(.*?)"`)
	lsdMatches := lsdRegex.FindStringSubmatch(bodyStr)
	if len(lsdMatches) > 1 {
		data.LSD = lsdMatches[1]
	} else {
		// Fallback 1: lsd":"..."
		lsdRegexFB1 := regexp.MustCompile(`lsd":"(.*?)"`)
		lsdMatchesFB1 := lsdRegexFB1.FindStringSubmatch(bodyStr)
		if len(lsdMatchesFB1) > 1 {
			data.LSD = lsdMatchesFB1[1]
		} else {
			// Fallback 2: type="hidden" name="lsd" value="..."
			lsdRegexFB2 := regexp.MustCompile(`name="lsd"\s*value="(.*?)"`)
			lsdMatchesFB2 := lsdRegexFB2.FindStringSubmatch(bodyStr)
			if len(lsdMatchesFB2) > 1 {
				data.LSD = lsdMatchesFB2[1]
			}
		}
	}

	// Tìm fb_dtsg - thử nhiều pattern
	// Debug: In context xung quanh "DTSG" để xem format thực tế
	dtsgPatterns := []string{
		`"(NAfv[a-zA-Z0-9_\-\:]+)"`, // Ưu tiên hàng hiệu NAfv (Chìa khóa vạn năng cho mutations)
		`"(NAfu[a-zA-Z0-9_\-\:]+)"`, // Ưu tiên mã NAfu
		`"DTSGInitialData"\s*,\s*\[\s*\]\s*,\s*\{\s*"token"\s*:\s*"(.*?)"`,
		`"DTSGInitData"\s*,\s*\[\s*\]\s*,\s*\{\s*"token"\s*:\s*"(.*?)"`,
		`"fb_dtsg"\s*:\s*"(.*?)"`,
		`name="fb_dtsg"\s*value="(.*?)"`,
	}
	for _, pattern := range dtsgPatterns {
		r := regexp.MustCompile(pattern)
		m := r.FindStringSubmatch(bodyStr)
		if len(m) > 1 {
			token := m[1]
			// Chuẩn hóa ngay lập tức: bỏ phần :3:timestamp nếu có
			if idx := strings.Index(token, ":"); idx > 0 {
				token = token[:idx]
			}
			data.DTSG = token
			fmt.Printf("[INFO] Đã lấy được fb_dtsg mới (pattern: %s): %s\n", pattern[:min(20, len(pattern))], data.DTSG[:min(20, len(data.DTSG))])
			break
		}
	}

	if data.DTSG == "" {
		fmt.Printf("[WARN] Không tìm được fb_dtsg mới. Sẽ dùng fb_dtsg cũ trong file json.\n")
	}

	if data.LSD == "" {
		snippet := bodyStr
		if len(snippet) > 1000 {
			snippet = snippet[:1000]
		}
		fmt.Printf("[DEBUG] Nội dung trang chủ (1000 ký tự đầu): \n%s\n", snippet)
		return data, fmt.Errorf("không tìm thấy LSD token trên trang chủ (đã thử 3 mẫu regex)")
	}

	return data, nil
}

// min trả về giá trị nhỏ hơn của 2 số
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// calcJazoest tính jazoest checksum từ fb_dtsg
// Công thức: "2" + tổng ASCII code của tất cả ký tự trong dtsg
func calcJazoest(dtsg string) string {
	sum := 0
	for _, c := range dtsg {
		sum += int(c)
	}
	return fmt.Sprintf("2%d", sum)
}

type ActionHandler struct {
	sessionHandler *SessionHandler
	accountStore   accounts.Store
	fbDataStore    fbdata.Store
}

// NewActionHandler tạo handler không có account store (backward compat)
func NewActionHandler(sessionHandler *SessionHandler) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   nil,
		fbDataStore:    nil,
	}
}

// NewActionHandlerWithStore tạo handler có account store thật
func NewActionHandlerWithStore(sessionHandler *SessionHandler, store accounts.Store, fbStore fbdata.Store) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   store,
		fbDataStore:    fbStore,
	}
}

func (h *ActionHandler) LikePost(req LikePostRequest) LikePostResponse {
	now := time.Now()

	postID := strings.TrimSpace(req.PostID)
	if postID == "" {
		resp := LikePostResponse{
			Success: false, Status: StatusValidationError,
			Action: req.ReactionType, PostID: "",
			AccountID: req.AccountID, ReactionType: req.ReactionType,
			DryRun: req.DryRun, Message: "Post ID không được để trống.", ExecutedAt: now,
		}
		return resp
	}

	reactionType := strings.ToLower(strings.TrimSpace(req.ReactionType))
	if reactionType == "" {
		reactionType = ReactionLike
	}
	if !ValidReactions[reactionType] {
		resp := LikePostResponse{
			Success: false, Status: StatusValidationError,
			Action: reactionType, PostID: postID,
			AccountID: req.AccountID, ReactionType: reactionType,
			DryRun: req.DryRun, Message: fmt.Sprintf("Loại cảm xúc không hợp lệ: %s", reactionType), ExecutedAt: now,
		}
		return resp
	}

	accountID := strings.TrimSpace(req.AccountID)
	displayName := accountID
	if accountID == "" {
		resp := LikePostResponse{
			Success: false, Status: StatusValidationError,
			Action: reactionType, PostID: postID,
			AccountID: "", ReactionType: reactionType,
			DryRun: req.DryRun, Message: "Chưa chọn tài khoản thực thi.", ExecutedAt: now,
		}
		return resp
	}

	sessionValid := true
	sessionStatusStr := "active"
	if h.accountStore != nil {
		acc, err := h.accountStore.Get(accountID)
		if err != nil || acc == nil {
			resp := LikePostResponse{
				Success: false, Status: StatusValidationError,
				Action: reactionType, PostID: postID,
				AccountID: accountID, ReactionType: reactionType,
				DryRun: req.DryRun,
				Message: fmt.Sprintf("Tài khoản '%s' không tìm thấy trong local store.", accountID),
				ExecutedAt: now,
			}
			AddLog(toLog(resp, "Tài khoản không tồn tại"))
			return resp
		}
		displayName = acc.DisplayName
		sessionStatusStr = acc.SessionStatus

		if acc.SessionStatus == string(accounts.SessionInvalid) ||
			acc.SessionStatus == string(accounts.SessionExpired) ||
			acc.SessionStatus == "error" {
			sessionValid = false
		}
	} else {
		status := h.sessionHandler.GetSessionStatus()
		if !status.IsActive {
			sessionValid = false
		}
	}

	if !sessionValid {
		resp := LikePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: reactionType, PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
			DryRun: req.DryRun,
			Message: fmt.Sprintf("Phiên của tài khoản '%s' không hợp lệ (trạng thái: %s).", displayName, sessionStatusStr),
			ExecutedAt: now,
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	if req.DryRun {
		sleepMs := 500 + rand.Intn(701)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		resp := LikePostResponse{
			Success: true, Status: StatusSuccess,
			Action: reactionType, PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
			DryRun: true,
			Message: fmt.Sprintf("[DRY RUN] Mô phỏng gửi '%s' thành công — Tài khoản: %s — Bài viết: %s", reactionType, displayName, postID),
			ExecutedAt: time.Now(),
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	resp := LikePostResponse{
		Success: false, Status: StatusNotEnabled,
		Action: reactionType, PostID: postID,
		AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
		DryRun: false,
		Message: "Chế độ Real Run chưa được kích hoạt cho Action Like.",
		ExecutedAt: time.Now(),
	}
	AddLog(toLog(resp, resp.Message))
	return resp
}

func (h *ActionHandler) CommentPost(req CommentPostRequest) CommentPostResponse {
	now := time.Now()

	postID := strings.TrimSpace(req.PostID)
	if postID == "" {
		resp := CommentPostResponse{
			Success: false, Status: StatusValidationError,
			Action: "comment", PostID: "",
			AccountID: req.AccountID, CommentText: req.CommentText,
			DryRun: req.DryRun, Message: "Post ID không được để trống.", ExecutedAt: now,
		}
		return resp
	}

	commentText := strings.TrimSpace(req.CommentText)
	if commentText == "" {
		resp := CommentPostResponse{
			Success: false, Status: StatusValidationError,
			Action: "comment", PostID: postID,
			AccountID: req.AccountID, CommentText: commentText,
			DryRun: req.DryRun, Message: "Nội dung bình luận không được để trống.", ExecutedAt: now,
		}
		return resp
	}

	accountID := strings.TrimSpace(req.AccountID)
	displayName := accountID
	if accountID == "" {
		resp := CommentPostResponse{
			Success: false, Status: StatusValidationError,
			Action: "comment", PostID: postID,
			AccountID: "", CommentText: commentText,
			DryRun: req.DryRun, Message: "Chưa chọn tài khoản thực thi.", ExecutedAt: now,
		}
		return resp
	}

	sessionValid := true
	sessionStatusStr := "active"
	if h.accountStore != nil {
		acc, err := h.accountStore.Get(accountID)
		if err != nil || acc == nil {
			resp := CommentPostResponse{
				Success: false, Status: StatusValidationError,
				Action: "comment", PostID: postID,
				AccountID: accountID, CommentText: commentText,
				DryRun: req.DryRun,
				Message: fmt.Sprintf("Tài khoản '%s' không tìm thấy trong local store.", accountID),
				ExecutedAt: now,
			}
			AddLog(logFromComment(resp, "Tài khoản không tồn tại"))
			return resp
		}
		displayName = acc.DisplayName
		sessionStatusStr = acc.SessionStatus

		if acc.SessionStatus == string(accounts.SessionInvalid) ||
			acc.SessionStatus == string(accounts.SessionExpired) ||
			acc.SessionStatus == "error" {
			sessionValid = false
		}
	}

	if !sessionValid {
		resp := CommentPostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
			DryRun: req.DryRun,
			Message: fmt.Sprintf("Phiên của tài khoản '%s' không hợp lệ (trạng thái: %s).", displayName, sessionStatusStr),
			ExecutedAt: now,
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	if req.DryRun {
		sleepMs := 500 + rand.Intn(701)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		snippet := commentText
		if len(snippet) > 20 {
			snippet = snippet[:20] + "..."
		}

		resp := CommentPostResponse{
			Success: true, Status: StatusSuccess,
			Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
			DryRun: true,
			Message: fmt.Sprintf("[DRY RUN] Đã bình luận: \"%s\" bằng tài khoản %s — Bài viết: %s", snippet, displayName, postID),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	// ── 6. Real Run — Nối ống Executor bắt đầu ───────────────────────────────
	if h.fbDataStore == nil {
		resp := CommentPostResponse{
			Success: false, Status: StatusNotEnabled,
			Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
			DryRun: false,
			Message: "Lỗi hệ thống: fbDataStore chưa được khởi tạo.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	fbData, err := h.fbDataStore.Get(accountID)
	if err != nil || fbData == nil || fbData.Info.Cookie == "" || fbData.Info.FbDtsg == "" {
		resp := CommentPostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
			DryRun: false,
			Message: "Lỗi: Không tìm thấy Cookie hoặc fb_dtsg trong FB Data của tài khoản này.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}
	
	// Delay chống Rate Limit — nghỉ ngẫu nhiên 3-8 giây để qua mặt bot-detector FB
	sleepMs := 3000 + rand.Intn(5001) // 3000ms → 8000ms
	fmt.Printf("[DELAY] Chờ %dms trước khi gửi request...\n", sleepMs)
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)

	client := &http.Client{Timeout: 15 * time.Second}
	reqUrl := "https://www.facebook.com/api/graphql/"

	// doc_id mới nhất cho Comment Create Mutation do user dâng hiến
	docID := "27321223354145000"

	// Lọc lấy c_user (UID thật) từ Cookie để làm actor_id
	actorID := ""
	for _, part := range strings.Split(fbData.Info.Cookie, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "c_user=") {
			actorID = strings.TrimPrefix(part, "c_user=")
			break
		}
	}
	if actorID == "" {
		actorID = fbData.UID // Dự phòng
	}

	variables := fmt.Sprintf(`{"input":{"client_mutation_id":"%d","actor_id":"%s","feedback_id":"%s","message":{"ranges":[],"text":"%s"}}}`,
		time.Now().UnixNano(), actorID, postID, commentText)

	data := url.Values{}
	data.Set("doc_id", docID)
	data.Set("fb_dtsg", fbData.Info.FbDtsg)
	data.Set("variables", variables)

	httpReq, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		resp := CommentPostResponse{
			Success: false, Status: "error", Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
			Message: "Lỗi cấu trúc HTTP request", ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Cookie", fbData.Info.Cookie)
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	httpReq.Header.Set("Sec-Fetch-Site", "same-origin")
	httpReq.Header.Set("X-FB-Friendly-Name", "UFI2CommentCreateMutation")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		resp := CommentPostResponse{
			Success: false, Status: "error", Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
			Message: fmt.Sprintf("Lỗi mạng: %v", err), ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}
	defer httpResp.Body.Close()

	bodyBytes, _ := io.ReadAll(httpResp.Body)
	bodyStr := string(bodyBytes)
	fmt.Printf("\n=== GRAPHQL RAW RESP ===\n%s\n====================\n\n", bodyStr)

	if httpResp.StatusCode != 200 {
		resp := CommentPostResponse{
			Success: false, Status: "error", Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
			Message: fmt.Sprintf("GraphQL Server trả về mã lỗi %d", httpResp.StatusCode),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	isSuccess := strings.Contains(bodyStr, "\"comment_create\":{")
	
	if !isSuccess && (strings.Contains(bodyStr, "\"errors\":[{") || strings.Contains(bodyStr, "Exception")) {
		snippet := bodyStr
		if len(snippet) > 80 {
			snippet = snippet[:80] + "..."
		}
		resp := CommentPostResponse{
			Success: false, Status: "error", Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
			Message:   fmt.Sprintf("FB GraphQL Error: %s", snippet),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	resp := CommentPostResponse{
		Success: true, Status: StatusSuccess,
		Action: "comment", PostID: postID,
		AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
		DryRun: false,
		Message: "Real Run: Gửi payload GraphQl Comment thành công!",
		ExecutedAt: time.Now(),
	}
	AddLog(logFromComment(resp, resp.Message))
	return resp
}

func (h *ActionHandler) GetActionLogs() []ActionLog {
	return GetLogs()
}

func toLog(resp LikePostResponse, msg string) ActionLog {
	return ActionLog{
		ActionType:         resp.Action,
		AccountID:          resp.AccountID,
		AccountDisplayName: resp.AccountDisplayName,
		PostID:             resp.PostID,
		ReactionType:       resp.ReactionType,
		Status:             resp.Status,
		DryRun:             resp.DryRun,
		Message:            msg,
		ExecutedAt:         resp.ExecutedAt,
	}
}

func logFromComment(resp CommentPostResponse, msg string) ActionLog {
	return ActionLog{
		ActionType:         resp.Action,
		AccountID:          resp.AccountID,
		AccountDisplayName: resp.AccountDisplayName,
		PostID:             resp.PostID,
		CommentText:        resp.CommentText,
		Status:             resp.Status,
		DryRun:             resp.DryRun,
		Message:            msg,
		ExecutedAt:         resp.ExecutedAt,
	}
}

func logFromPost(resp CreatePostResponse, msg string) ActionLog {
	return ActionLog{
		ActionType:         resp.Action,
		AccountID:          resp.AccountID,
		AccountDisplayName: resp.AccountDisplayName,
		PostText:           resp.PostText,
		Status:             resp.Status,
		DryRun:             resp.DryRun,
		Message:            msg,
		ExecutedAt:         resp.ExecutedAt,
	}
}

func (h *ActionHandler) CreatePost(req CreatePostRequest) CreatePostResponse {
	now := time.Now()

	postText := strings.TrimSpace(req.PostText)
	if postText == "" {
		resp := CreatePostResponse{
			Success: false, Status: StatusValidationError,
			Action: "create_post", AccountID: req.AccountID, PostText: postText,
			DryRun: req.DryRun, Message: "Nội dung đăng bài không được để trống.", ExecutedAt: now,
		}
		return resp
	}

	accountID := strings.TrimSpace(req.AccountID)
	displayName := accountID
	if accountID == "" {
		resp := CreatePostResponse{
			Success: false, Status: StatusValidationError,
			Action: "create_post", AccountID: "", PostText: postText,
			DryRun: req.DryRun, Message: "Chưa chọn tài khoản thực thi.", ExecutedAt: now,
		}
		return resp
	}

	sessionValid := true
	sessionStatusStr := "active"
	if h.accountStore != nil {
		acc, err := h.accountStore.Get(accountID)
		if err != nil || acc == nil {
			resp := CreatePostResponse{
				Success: false, Status: StatusValidationError,
				Action: "create_post", AccountID: accountID, PostText: postText,
				DryRun: req.DryRun,
				Message: fmt.Sprintf("Tài khoản '%s' không tìm thấy trong local store.", accountID),
				ExecutedAt: now,
			}
			AddLog(logFromPost(resp, "Tài khoản không tồn tại"))
			return resp
		}
		displayName = acc.DisplayName
		sessionStatusStr = acc.SessionStatus

		if acc.SessionStatus == string(accounts.SessionInvalid) ||
			acc.SessionStatus == string(accounts.SessionExpired) ||
			acc.SessionStatus == "error" {
			sessionValid = false
		}
	}

	if !sessionValid {
		resp := CreatePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: req.DryRun,
			Message: fmt.Sprintf("Phiên của tài khoản '%s' không hợp lệ (trạng thái: %s).", displayName, sessionStatusStr),
			ExecutedAt: now,
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	if req.DryRun {
		sleepMs := 500 + rand.Intn(701)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		snippet := postText
		if len(snippet) > 20 {
			snippet = snippet[:20] + "..."
		}

		resp := CreatePostResponse{
			Success: true, Status: StatusSuccess,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: true,
			Message: fmt.Sprintf("[DRY RUN] Đã đăng bài nháp: \"%s\" bằng tài khoản %s", snippet, displayName),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	// ── Real Run ────────────────────────────────────────────────────────

	// Lấy thông tin fb_dtsg và Cookie từ fbdata store đã inject
	if h.fbDataStore == nil {
		resp := CreatePostResponse{
			Success: false, Status: StatusFailed,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Lỗi: fbDataStore chưa được khởi tạo trong backend.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	// Lấy actor_id (c_user) từ cookie
	actorID := ""
	fbInfo, err := h.fbDataStore.Get(accountID)
	if err != nil || fbInfo == nil {
		resp := CreatePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false,
			Message: fmt.Sprintf("Lỗi: %s - Không tìm thấy dữ liệu FB (Cookie/fb_dtsg). Cập nhật vào facebook_data.json!", err),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}
	for _, part := range strings.Split(fbInfo.Info.Cookie, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "c_user=") {
			actorID = strings.TrimPrefix(part, "c_user=")
			actorID = strings.TrimSpace(actorID)
			break
		}
	}
	if actorID == "" {
		resp := CreatePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Không tìm thấy c_user (UID) trong Cookie. Cập nhật Cookie mới vào facebook_data.json!",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	// Delay chống Rate Limit
	sleepMs := 3000 + rand.Intn(5001)
	fmt.Printf("[DELAY] Nghỉ %dms trước khi đăng bài...\n", sleepMs)
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)

	// ── Two-Step DTSG Strategy ──────────────────────────────────────────────
	// Facebook trả về token hợp lệ trong lỗi 1357004 → học rồi thử lại ngay

	docID := "26695477506728717"

	// Helper: xây dựng variables JSON CHÍNH XÁC theo browser (captured từ F12)
	buildVariables := func(_ string) string {
		rawToken := generateToken()
		idempotenceToken := fmt.Sprintf("%s_FEED", rawToken)
		composerSessionID := rawToken // Cùng UUID, không có _FEED
		return fmt.Sprintf(`{`+
			`"input":{`+
			`"composer_entry_point":"inline_composer",`+
			`"composer_source_surface":"timeline",`+
			`"idempotence_token":%s,`+
			`"source":"WWW",`+
			`"attachments":[],`+
			`"audience":{"privacy":{"allow":[],"base_state":"FRIENDS","deny":[],"tag_expansion_state":"UNSPECIFIED"}},`+
			`"message":{"ranges":[],"text":%s},`+
			`"with_tags_ids":null,`+
			`"inline_activities":[],`+
			`"text_format_preset_id":"0",`+
			`"publishing_flow":{"supported_flows":["ASYNC_SILENT","ASYNC_NOTIF","FALLBACK"]},`+
			`"logging":{"composer_session_id":%s},`+
			`"navigation_data":{"attribution_id_v2":"ProfileCometTimelineListViewRoot.react,comet.profile.timeline.list,via_cold_start,1776157041276,201193,190055527696468,,"},`+
			`"tracking":[null],`+
			`"event_share_metadata":{"surface":"timeline"},`+
			`"actor_id":%s,`+
			`"client_mutation_id":"1"`+
			`},`+
			`"feedLocation":"TIMELINE",`+
			`"feedbackSource":0,`+
			`"focusCommentID":null,`+
			`"gridMediaWidth":230,`+
			`"groupID":null,`+
			`"scale":1,`+
			`"privacySelectorRenderLocation":"COMET_STREAM",`+
			`"checkPhotosToReelsUpsellEligibility":true,`+
			`"referringStoryRenderLocation":null,`+
			`"renderLocation":"timeline",`+
			`"useDefaultActor":false,`+
			`"inviteShortLinkKey":null,`+
			`"isFeed":false,`+
			`"isFundraiser":false,`+
			`"isFunFactPost":false,`+
			`"isGroup":false,`+
			`"isEvent":false,`+
			`"isTimeline":true,`+
			`"isSocialLearning":false,`+
			`"isPageNewsFeed":false,`+
			`"isProfileReviews":false,`+
			`"isWorkSharedDraft":false,`+
			`"canUserManageOffers":false,`+
			`"__relay_internal__pv__CometUFIShareActionMigrationrelayprovider":true,`+
			`"__relay_internal__pv__GHLShouldChangeSponsoredDataFieldNamerelayprovider":true,`+
			`"__relay_internal__pv__GHLShouldChangeAdIdFieldNamerelayprovider":true,`+
			`"__relay_internal__pv__CometUFI_dedicated_comment_routable_dialog_gkrelayprovider":true,`+
			`"__relay_internal__pv__CometUFICommentAutoTranslationTyperelayprovider":"ORIGINAL",`+
			`"__relay_internal__pv__CometUFICommentAvatarStickerAnimatedImagerelayprovider":false,`+
			`"__relay_internal__pv__CometUFICommentActionLinksRewriteEnabledrelayprovider":false,`+
			`"__relay_internal__pv__IsWorkUserrelayprovider":false,`+
			`"__relay_internal__pv__CometUFIReactionsEnableShortNamerelayprovider":false,`+
			`"__relay_internal__pv__CometUFISingleLineUFIrelayprovider":false,`+
			`"__relay_internal__pv__CometFeedStory_enable_post_permalink_white_space_clickrelayprovider":false,`+
			`"__relay_internal__pv__TestPilotShouldIncludeDemoAdUseCaserelayprovider":false,`+
			`"__relay_internal__pv__FBReels_deprecate_short_form_video_context_gkrelayprovider":true,`+
			`"__relay_internal__pv__FBReels_enable_view_dubbed_audio_type_gkrelayprovider":true,`+
			`"__relay_internal__pv__CometImmersivePhotoCanUserDisable3DMotionrelayprovider":false,`+
			`"__relay_internal__pv__WorkCometIsEmployeeGKProviderrelayprovider":false,`+
			`"__relay_internal__pv__IsMergQAPollsrelayprovider":false,`+
			`"__relay_internal__pv__FBReelsMediaFooter_comet_enable_reels_ads_gkrelayprovider":true,`+
			`"__relay_internal__pv__FBReelsIFUTileContent_reelsIFUPlayOnHoverrelayprovider":true,`+
			`"__relay_internal__pv__GroupsCometGYSJFeedItemHeightrelayprovider":206,`+
			`"__relay_internal__pv__ShouldEnableBakedInTextStoriesrelayprovider":false,`+
			`"__relay_internal__pv__StoriesShouldIncludeFbNotesrelayprovider":false,`+
			`"__relay_internal__pv__groups_comet_use_glvrelayprovider":false,`+
			`"__relay_internal__pv__GHLShouldChangeSponsoredAuctionDistanceFieldNamerelayprovider":true,`+
			`"__relay_internal__pv__GHLShouldUseSponsoredAuctionLabelFieldNameV1relayprovider":true,`+
			`"__relay_internal__pv__GHLShouldUseSponsoredAuctionLabelFieldNameV2relayprovider":false`+
			`}`,
			jsonStr(idempotenceToken), jsonStr(postText), jsonStr(composerSessionID), jsonStr(actorID))
	}

	// Lấy LSD token (cần thiết cho Comet GraphQL endpoint)
	lsd := ""
	sessionData, sdErr := fetchSessionData(fbInfo.Info.Cookie)
	if sdErr == nil && sessionData.LSD != "" {
		lsd = sessionData.LSD
		fmt.Printf("[INFO] LSD token: %s\n", lsd)
	} else {
		fmt.Printf("[WARN] Không lấy được LSD: %v\n", sdErr)
	}

	// Helper: thực hiện một lần POST lên GraphQL
	doGraphQLPost := func(fbDtsg string) (int, string, error) {
		vars := buildVariables(fbDtsg)
		fd := url.Values{}
		fd.Set("fb_dtsg", fbDtsg) // Gửi FULL token (kể cả :3:timestamp)
		// Jazoest tính từ phần token trước dấu ":" đầu tiên
		jazoestBase := fbDtsg
		if idx := strings.Index(jazoestBase, ":"); idx > 0 {
			jazoestBase = jazoestBase[:idx]
		}
		fd.Set("jazoest", calcJazoest(jazoestBase))
		fd.Set("lsd", lsd)
		fd.Set("doc_id", docID)
		fd.Set("variables", vars)

		req, err := http.NewRequest("POST", "https://www.facebook.com/api/graphql/", strings.NewReader(fd.Encode()))
		if err != nil {
			return 0, "", err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", fbInfo.Info.Cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
		req.Header.Set("Origin", "https://www.facebook.com")
		req.Header.Set("Referer", "https://www.facebook.com/")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("X-FB-Friendly-Name", "ComposerStoryCreateMutation")
		req.Header.Set("X-FB-LSD", lsd)
		req.Header.Set("X-ASBD-ID", "129477")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")

		// Client với redirect handler để giữ nguyên headers qua redirect
		cl := &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Giữ lại Cookie và headers qua các lần redirect
				if len(via) > 0 {
					fmt.Printf("[DEBUG] Redirect: %s -> %s\n", via[len(via)-1].URL, req.URL)
					req.Header.Set("Cookie", fbInfo.Info.Cookie)
					req.Header.Set("User-Agent", via[0].Header.Get("User-Agent"))
				}
				if len(via) >= 3 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		}
		httpResp, err := cl.Do(req)
		if err != nil {
			return 0, "", err
		}
		defer httpResp.Body.Close()
		// Debug: log response headers
		fmt.Printf("[DEBUG] Resp headers: Content-Type=%s, Content-Length=%s, Location=%s\n",
			httpResp.Header.Get("Content-Type"),
			httpResp.Header.Get("Content-Length"),
			httpResp.Header.Get("Location"))
		b, _ := io.ReadAll(httpResp.Body)
		return httpResp.StatusCode, string(b), nil
	}

	// ── Lần 1: Thử với fb_dtsg ĐẦY ĐỦ (bao gồm :3:timestamp) ────────────────────
	// QUAN TRỌẠNG: KHÔNG cắt bỏ :version:ts vì nó là một phần của token!
	currentDtsg := fbInfo.Info.FbDtsg
	shortDtsg := currentDtsg
	if len(shortDtsg) > 20 {
		shortDtsg = shortDtsg[:20]
	}
	fmt.Printf("[INFO] Lần 1: Gửi với fb_dtsg=%s...\n", shortDtsg)

	status1, body1, err := doGraphQLPost(currentDtsg)
	if err != nil {
		resp := CreatePostResponse{
			Success: false, Status: StatusFailed,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Lỗi HTTP lần 1: " + err.Error(), ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	snippet1 := body1
	if len(snippet1) > 300 {
		snippet1 = snippet1[:300]
	}
	fmt.Printf("=== RESP Lần 1 (status=%d, len=%d) ===\n%s\n====================\n", status1, len(body1), snippet1)

	// Empty body = silent rejection (thiếu LSD hoặc header khác)
	if body1 == "" {
		resp := CreatePostResponse{
			Success: false, Status: StatusFailed,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: fmt.Sprintf("Facebook trả về body trống (HTTP %d) - LSD: %s", status1, lsd),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	// Kiểm tra nếu lần 1 đã thành công (không có lỗi)
	if !strings.Contains(body1, `"error":`) && !strings.Contains(body1, `"errors":[`) {
		resp := CreatePostResponse{
			Success: true, Status: StatusSuccess,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: fmt.Sprintf("Đăng bài thành công lần 1 — Tài khoản: %s!", displayName),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	// ── Lần 2: Học token từ lỗi rồi thử lại ────────────────────────────────
	freshDtsg := extractDtsgFromError(body1)
	if freshDtsg == "" {
		resp := CreatePostResponse{
			Success: false, Status: StatusFailed,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Không trích xuất được dtsgToken từ lỗi. Body: " + snippet1,
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	fmt.Printf("[INFO] Lần 2: Thử lại với fresh dtsg=%s...\n", freshDtsg[:min(20, len(freshDtsg))])
	time.Sleep(1 * time.Second) // Nghỉ 1 giây giữa 2 lần

	_, body2, err := doGraphQLPost(freshDtsg)
	if err != nil {
		resp := CreatePostResponse{
			Success: false, Status: StatusFailed,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Lỗi HTTP lần 2: " + err.Error(), ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	snippet2 := body2
	if len(snippet2) > 400 {
		snippet2 = snippet2[:400]
	}
	fmt.Printf("=== RESP Lần 2 ===\n%s\n====================\n", snippet2)

	if strings.Contains(body2, `"error":`) || strings.Contains(body2, `"errors":[`) {
		resp := CreatePostResponse{
			Success: false, Status: StatusFailed,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Vẫn lỗi ở lần 2: " + snippet2, ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	resp := CreatePostResponse{
		Success: true, Status: StatusSuccess,
		Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
		DryRun: false, Message: fmt.Sprintf("Đăng bài thành công (Two-Step DTSG) — Tài khoản: %s!", displayName),
		ExecutedAt: time.Now(),
	}
	AddLog(logFromPost(resp, resp.Message))
	return resp
}

