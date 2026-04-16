package actiontest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yourname/tool-facebook/internal/accounts"
	"github.com/yourname/tool-facebook/internal/fbdata"
)

// JsonStr trả về chuỗi được JSON encode an toàn (có xử lý ký tự đặc biệt)
func JsonStr(s string) string {
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

// ActorIDFromCookie trích xuất UID từ cookie string
func ActorIDFromCookie(cookie string) string {
	for _, part := range strings.Split(cookie, ";") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "c_user=") {
			return strings.TrimPrefix(part, "c_user=")
		}
	}
	return ""
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func isLikelyLocalMediaInput(val string) bool {
	v := strings.TrimSpace(strings.ToLower(val))
	if v == "" {
		return false
	}
	if strings.HasPrefix(v, "file://") {
		return true
	}
	if strings.Contains(v, `\`) || strings.HasPrefix(v, "/") || strings.HasPrefix(v, "./") || strings.HasPrefix(v, "../") {
		return true
	}
	if strings.Contains(v, "://") {
		return false
	}
	exts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".mp4", ".mov", ".avi", ".mkv", ".webm"}
	for _, ext := range exts {
		if strings.HasSuffix(v, ext) {
			return true
		}
	}
	return false
}

// extractPhotoIDsFromImagePaths chuyen danh sach image_paths thanh photo IDs dung cho GraphQL attachments.
// Ho tro cac dinh dang:
// - "122100411308738587"
// - "media_fbid_122100411308738587"
// - URL co query fbid=...
func extractPhotoIDsFromImagePaths(imagePaths []string) ([]string, []string, []string) {
	seen := make(map[string]struct{})
	ids := make([]string, 0, len(imagePaths))
	localInputs := make([]string, 0)
	ignored := make([]string, 0)

	for _, raw := range imagePaths {
		val := strings.TrimSpace(raw)
		if val == "" {
			continue
		}

		candidate := ""
		switch {
		case strings.HasPrefix(val, "media_fbid_"):
			candidate = strings.TrimPrefix(val, "media_fbid_")
		case isAllDigits(val):
			candidate = val
		default:
			if u, err := url.Parse(val); err == nil {
				if fbid := strings.TrimSpace(u.Query().Get("fbid")); isAllDigits(fbid) {
					candidate = fbid
				}
			}
		}

		if !isAllDigits(candidate) {
			if isLikelyLocalMediaInput(val) {
				localInputs = append(localInputs, val)
				continue
			}
			ignored = append(ignored, val)
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		ids = append(ids, candidate)
	}

	return ids, localInputs, ignored
}

func normalizeLocalMediaPath(input string) string {
	path := strings.TrimSpace(input)
	if path == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(path), "file://") {
		if u, err := url.Parse(path); err == nil {
			unescapedPath, _ := url.PathUnescape(u.Path)
			if len(unescapedPath) >= 3 && unescapedPath[0] == '/' && unescapedPath[2] == ':' {
				unescapedPath = unescapedPath[1:]
			}
			if u.Host != "" && !strings.Contains(unescapedPath, ":") {
				path = `\\` + u.Host + strings.ReplaceAll(unescapedPath, "/", `\`)
			} else {
				path = strings.ReplaceAll(unescapedPath, "/", string(os.PathSeparator))
			}
		}
	}
	return filepath.Clean(path)
}

func parseIDByPatterns(text string, patterns []string) string {
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			continue
		}
		if m := re.FindStringSubmatch(text); len(m) > 1 && isAllDigits(m[1]) {
			return m[1]
		}
	}
	return ""
}

func parseUploadedPhotoIDFromBody(body string) (string, string) {
	cleanBody := strings.TrimSpace(strings.TrimPrefix(body, "for (;;);"))
	preferredPatterns := []struct {
		key      string
		patterns []string
	}{
		{key: "photo_id", patterns: []string{`"photo_id"\s*:\s*"?(\d{8,})"?`, `"photoID"\s*:\s*"?(\d{8,})"?`}},
		{key: "legacy_fbid", patterns: []string{`"legacy_fbid"\s*:\s*"?(\d{8,})"?`}},
		// metadata[].fbid thường đáng tin hơn image_id/media_id.
		{key: "metadata_fbid", patterns: []string{`(?s)"metadata"\s*:\s*\[.*?"fbid"\s*:\s*"?(\d{8,})"?`}},
		{key: "fbid", patterns: []string{`"fbid"\s*:\s*"?(\d{8,})"?`}},
		{key: "image_id", patterns: []string{`"image_id"\s*:\s*"?(\d{8,})"?`}},
		{key: "media_id", patterns: []string{`"media_id"\s*:\s*"?(\d{8,})"?`}},
		{key: "id_122", patterns: []string{`\b(122\d{12,})\b`}},
		{key: "id_generic", patterns: []string{`"id"\s*:\s*"?(\d{15,})"?`}},
	}
	variants := []string{
		cleanBody,
		strings.ReplaceAll(cleanBody, `\"`, `"`),
	}
	for _, variant := range variants {
		for _, rule := range preferredPatterns {
			if id := parseIDByPatterns(variant, rule.patterns); id != "" {
				return id, rule.key
			}
		}
	}
	return "", ""
}

func parseUploadedFileIDFromBody(body string) string {
	cleanBody := strings.TrimSpace(strings.TrimPrefix(body, "for (;;);"))
	filePatterns := []string{
		`"file_id"\s*:\s*"?(\d{8,})"?`,
	}
	variants := []string{
		cleanBody,
		strings.ReplaceAll(cleanBody, `\"`, `"`),
	}
	for _, variant := range variants {
		if id := parseIDByPatterns(variant, filePatterns); id != "" {
			return id
		}
	}
	return ""
}

func summarizeUploadIDCandidates(body string) string {
	cleanBody := strings.TrimSpace(strings.TrimPrefix(body, "for (;;);"))
	variants := []string{
		cleanBody,
		strings.ReplaceAll(cleanBody, `\"`, `"`),
	}
	keys := []string{"photo_id", "photoID", "legacy_fbid", "image_id", "media_id", "file_id", "fbid"}
	seen := make(map[string]struct{})
	out := make([]string, 0, 8)
	for _, variant := range variants {
		for _, key := range keys {
			pattern := fmt.Sprintf(`"?%s"?\s*:\s*"?(\d{8,})"?`, regexp.QuoteMeta(key))
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}
			matches := re.FindAllStringSubmatch(variant, 5)
			for _, m := range matches {
				if len(m) < 2 || !isAllDigits(m[1]) {
					continue
				}
				pair := fmt.Sprintf("%s=%s", key, m[1])
				if _, ok := seen[pair]; ok {
					continue
				}
				seen[pair] = struct{}{}
				out = append(out, pair)
			}
		}
	}
	if len(out) == 0 {
		return "none"
	}
	return strings.Join(out, ", ")
}

// detectMimeType trả về Content-Type đúng dựa theo extension file.
// Cực kỳ quan trọng: nếu upload với Content-Type sai (application/octet-stream),
// Facebook sẽ không nhận ra đây là ảnh và không trả về photo_id/fbid.
func detectMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".mkv":
		return "video/x-matroska"
	case ".webm":
		return "video/webm"
	default:
		return "application/octet-stream"
	}
}

func buildAttachmentsJSONWithKey(ids []string, key string) string {
	if len(ids) == 0 || strings.TrimSpace(key) == "" {
		return "[]"
	}
	attachments := make([]map[string]map[string]string, 0, len(ids))
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		attachments = append(attachments, map[string]map[string]string{
			key: {"id": trimmed},
		})
	}
	if len(attachments) == 0 {
		return "[]"
	}
	if b, err := json.Marshal(attachments); err == nil {
		return string(b)
	}
	return "[]"
}

type attachmentVariant struct {
	Name string
	JSON string
}

func buildAttachmentVariants(photoIDs []string) []attachmentVariant {
	if len(photoIDs) == 0 {
		return []attachmentVariant{{Name: "none", JSON: "[]"}}
	}
	variants := []attachmentVariant{
		{Name: "photo", JSON: buildAttachmentsJSONWithKey(photoIDs, "photo")},
		{Name: "existing_photo", JSON: buildAttachmentsJSONWithKey(photoIDs, "existing_photo")},
		{Name: "media", JSON: buildAttachmentsJSONWithKey(photoIDs, "media")},
	}
	out := make([]attachmentVariant, 0, len(variants))
	seen := make(map[string]struct{})
	for _, v := range variants {
		if strings.TrimSpace(v.JSON) == "" || v.JSON == "[]" {
			continue
		}
		if _, ok := seen[v.JSON]; ok {
			continue
		}
		seen[v.JSON] = struct{}{}
		out = append(out, v)
	}
	if len(out) == 0 {
		return []attachmentVariant{{Name: "none", JSON: "[]"}}
	}
	return out
}

func hasGraphQLErrorCode(body string, code string) bool {
	needle := `"code":` + strings.TrimSpace(code)
	needleAlt := `"api_error_code":` + strings.TrimSpace(code)
	return strings.Contains(body, needle) || strings.Contains(body, needleAlt)
}

func uploadOneLocalMedia(cookie, actorID, fbDtsg, lsd string, sd SessionData, localInput string, composerSessionID string) (string, error) {
	localPath := normalizeLocalMediaPath(localInput)
	if localPath == "" {
		return "", fmt.Errorf("local media path is empty")
	}
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return "", fmt.Errorf("cannot access file '%s': %w", localPath, err)
	}
	if fileInfo.IsDir() {
		return "", fmt.Errorf("'%s' is a directory, expected file", localPath)
	}

	type uploadStrategy struct {
		endpoint  string
		fileField string
	}
	strategies := []uploadStrategy{
		{endpoint: "https://upload.facebook.com/ajax/react_composer/attachments/photo/upload", fileField: "farr"},
		{endpoint: "https://www.facebook.com/ajax/react_composer/attachments/photo/upload", fileField: "farr"},
		{endpoint: "https://upload.facebook.com/ajax/react_composer/attachments/photo/upload", fileField: "source"},
		{endpoint: "https://upload.facebook.com/ajax/composerx/attachment/media/upload", fileField: "farr"},
		{endpoint: "https://www.facebook.com/ajax/composerx/attachment/media/upload", fileField: "farr"},
		{endpoint: "https://upload.facebook.com/ajax/composerx/attachment/media/upload", fileField: "source"},
		{endpoint: "https://upload.facebook.com/ajax/browser/graphql_media_upload", fileField: "farr"}, // Modern comet endpoint
		{endpoint: "https://www.facebook.com/ajax/browser/graphql_media_upload", fileField: "farr"},
		{endpoint: "https://upload.facebook.com/ajax/mercury/upload.php", fileField: "upload_1024"},
		{endpoint: "https://www.facebook.com/ajax/mercury/upload.php", fileField: "upload_1024"},
	}

	var lastErr error
	for _, strategy := range strategies {
		fmt.Printf("[INFO] Uploading local media via %s (field=%s): %s\n", strategy.endpoint, strategy.fileField, localPath)
		file, err := os.Open(localPath)
		if err != nil {
			return "", fmt.Errorf("cannot open file '%s': %w", localPath, err)
		}

		var reqBody bytes.Buffer
		writer := multipart.NewWriter(&reqBody)
		_ = writer.WriteField("av", actorID)
		_ = writer.WriteField("__user", actorID)
		_ = writer.WriteField("__a", "1")
		_ = writer.WriteField("__req", "2b")
		_ = writer.WriteField("__comet_req", "15")
		if strings.TrimSpace(composerSessionID) != "" {
			_ = writer.WriteField("composer_session_id", strings.TrimSpace(composerSessionID))
		}
		
		// BUG FIX: Essential CSRF bypass
		_ = writer.WriteField("jazoest", CalcJazoest(fbDtsg))
		_ = writer.WriteField("lsd", lsd)

		_ = writer.WriteField("fb_dtsg", fbDtsg)
		_ = writer.WriteField("__ccg", "EXCELLENT")
		if strings.TrimSpace(sd.HS) != "" {
			_ = writer.WriteField("__hs", strings.TrimSpace(sd.HS))
		}
		if strings.TrimSpace(sd.S) != "" {
			_ = writer.WriteField("__s", strings.TrimSpace(sd.S))
		}
		if strings.TrimSpace(sd.Dyn) != "" {
			_ = writer.WriteField("__dyn", strings.TrimSpace(sd.Dyn))
		}
		if strings.TrimSpace(sd.CSR) != "" {
			_ = writer.WriteField("__csr", strings.TrimSpace(sd.CSR))
		}
		if strings.TrimSpace(sd.Rev) != "" {
			_ = writer.WriteField("__rev", strings.TrimSpace(sd.Rev))
		}
		if strings.TrimSpace(sd.HSI) != "" {
			_ = writer.WriteField("__hsi", strings.TrimSpace(sd.HSI))
		}
		if strings.TrimSpace(sd.SpinR) != "" {
			_ = writer.WriteField("__spin_r", strings.TrimSpace(sd.SpinR))
		}
		_ = writer.WriteField("__spin_b", "trunk")
		if strings.TrimSpace(sd.SpinT) != "" {
			_ = writer.WriteField("__spin_t", strings.TrimSpace(sd.SpinT))
		}
		_ = writer.WriteField("fb_dtsg", fbDtsg)
		_ = writer.WriteField("jazoest", CalcJazoest(fbDtsg))
		if strings.TrimSpace(lsd) != "" {
			_ = writer.WriteField("lsd", strings.TrimSpace(lsd))
		}
		_ = writer.WriteField("profile_id", actorID)
		_ = writer.WriteField("target_id", actorID)
		_ = writer.WriteField("upload_source", "composer")
		if strings.TrimSpace(composerSessionID) != "" {
			_ = writer.WriteField("composer_session_id", strings.TrimSpace(composerSessionID))
		}
		_ = writer.WriteField("composer_entry_point", "inline_composer")
		_ = writer.WriteField("composer_source_surface", "timeline")
		_ = writer.WriteField("waterfallxapp", "comet")
		_ = writer.WriteField("voice_clip", "false")
		_ = writer.WriteField("story_attachment", "true")
		_ = writer.WriteField("profile_id", actorID)
		_ = writer.WriteField("target_id", actorID)
		_ = writer.WriteField("source", "8")

		// BUG FIX: Dùng CreatePart với Content-Type chính xác thay vì CreateFormFile.
		// CreateFormFile luôn set "application/octet-stream" → Facebook không nhận ra là ảnh
		// → không trả về photo_id/fbid → upload thất bại.
		mimeType := detectMimeType(localPath)
		partHeader := make(textproto.MIMEHeader)
		partHeader.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="%s"; filename="%s"`, strategy.fileField, filepath.Base(localPath)))
		partHeader.Set("Content-Type", mimeType)
		part, err := writer.CreatePart(partHeader)
		if err == nil {
			_, err = io.Copy(part, file)
		}
		_ = file.Close()
		if err != nil {
			lastErr = fmt.Errorf("cannot append media file to multipart payload: %w", err)
			_ = writer.Close()
			continue
		}
		_ = writer.Close()

		req, err := http.NewRequest("POST", strategy.endpoint, &reqBody)
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Cookie", cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
		req.Header.Set("Origin", "https://www.facebook.com")
		req.Header.Set("Referer", "https://www.facebook.com/")
		req.Header.Set("Accept", "*/*")

		if strings.Contains(strategy.endpoint, "upload.facebook.com") {
			req.Header.Set("Sec-Fetch-Site", "same-site")
		} else {
			req.Header.Set("Sec-Fetch-Site", "same-origin")
		}

		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("X-FB-LSD", lsd)
		req.Header.Set("X-ASBD-ID", "129477")

		httpClient := &http.Client{Timeout: 60 * time.Second}
		httpResp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		respBytes, _ := io.ReadAll(httpResp.Body)
		_ = httpResp.Body.Close()
		respText := strings.TrimSpace(string(respBytes))

		// DEBUG: In toàn bộ body nếu không phải là mercury để xem FB trả lỗi gì
		if !strings.Contains(strategy.endpoint, "mercury") {
			fmt.Printf("[DEBUG] Endpoint %s response (status=%d): %s\n", strategy.endpoint, httpResp.StatusCode, respText)
		}

		if httpResp.StatusCode >= 400 {
			snippet := respText
			if len(snippet) > 180 {
				snippet = snippet[:180] + "..."
			}
			lastErr = fmt.Errorf("upload endpoint '%s' returned status=%d body=%s", strategy.endpoint, httpResp.StatusCode, snippet)
			continue
		}
		if strings.Contains(respText, `"error":1357001`) || strings.Contains(respText, `"errorSummary":"\u0110\u0103ng nh\u1eadp \u0111\u1ec3 ti\u1ebfp t\u1ee5c"`) {
			lastErr = fmt.Errorf("upload endpoint '%s' yêu cầu đăng nhập lại (error 1357001)", strategy.endpoint)
			continue
		}

		photoID, idSource := parseUploadedPhotoIDFromBody(respText)
		if photoID == "" {
			candidates := summarizeUploadIDCandidates(respText)
			if candidates != "none" {
				fmt.Printf("[WARN] Upload endpoint %s khong parse duoc photo_id. Candidates: %s\n", strategy.endpoint, candidates)
			}
			if fileID := parseUploadedFileIDFromBody(respText); fileID != "" {
				lastErr = fmt.Errorf("upload endpoint '%s' returned file_id=%s but no usable photo_id/fbid", strategy.endpoint, fileID)
				continue
			}
			snippet := respText
			if len(snippet) > 220 {
				snippet = snippet[:220] + "..."
			}
			lastErr = fmt.Errorf("upload endpoint '%s' did not return photo id, body=%s", strategy.endpoint, snippet)
			continue
		}
		fmt.Printf("[INFO] Upload candidates: %s\n", summarizeUploadIDCandidates(respText))
		fmt.Printf("[INFO] Uploaded local media '%s' -> photo_id=%s (source=%s) via %s\n", localPath, photoID, idSource, strategy.endpoint)
		return photoID, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("all upload strategies failed for '%s'", localPath)
	}
	return "", lastErr
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
	HS    string
	S     string
	Dyn   string
	CSR   string
	SpinR string
	SpinT string
}

// FetchSessionData truy cập trang chủ Facebook ngầm để lấy lsd, fb_dtsg và các thông số kỹ thuật khác
func FetchSessionData(cookie string) (data SessionData, err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://www.facebook.com/", nil)

	// Làm sạch cookie triệt để để tránh lỗi 400 Bad Request
	cleanCookie := strings.ReplaceAll(cookie, "\n", "")
	cleanCookie = strings.ReplaceAll(cleanCookie, "\r", "")
	cleanCookie = strings.TrimSpace(cleanCookie)

	// Headers siêu tàng hình
	req.Header.Set("Cookie", cleanCookie)
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

	// Lay __hs
	hsRegex := regexp.MustCompile(`"__hs":"([^"]+)"`)
	hsMatches := hsRegex.FindStringSubmatch(bodyStr)
	if len(hsMatches) > 1 {
		data.HS = hsMatches[1]
	}

	// Lay __s
	sRegex := regexp.MustCompile(`"__s":"([^"]+)"`)
	sMatches := sRegex.FindStringSubmatch(bodyStr)
	if len(sMatches) > 1 {
		data.S = sMatches[1]
	}

	// Lay __dyn
	dynRegex := regexp.MustCompile(`"__dyn":"([^"]+)"`)
	dynMatches := dynRegex.FindStringSubmatch(bodyStr)
	if len(dynMatches) > 1 {
		data.Dyn = dynMatches[1]
	}

	// Lay __csr
	csrRegex := regexp.MustCompile(`"__csr":"([^"]+)"`)
	csrMatches := csrRegex.FindStringSubmatch(bodyStr)
	if len(csrMatches) > 1 {
		data.CSR = csrMatches[1]
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
	isLoggedIn := strings.Contains(bodyStr, `"USER_ID":"`) || strings.Contains(bodyStr, "c_user=") || strings.Contains(bodyStr, ActorIDFromCookie(cookie))

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
		`"DTSGInitialData"\s*,\s*\[\s*\]\s*,\s*\{\s*"token"\s*:\s*"(.+?)"`,
		`"DTSGInitData"\s*,\s*\[\s*\]\s*,\s*\{\s*"token"\s*:\s*"(.+?)"`,
		`"fb_dtsg"\s*:\s*"(.+?)"`,
		`name="fb_dtsg"\s*value="(.+?)"`,
	}
	for _, pattern := range dtsgPatterns {
		r := regexp.MustCompile(pattern)
		m := r.FindStringSubmatch(bodyStr)
		if len(m) > 1 {
			token := strings.TrimSpace(m[1])
			if token == "" {
				continue
			}
			// Giữ nguyên token full (kể cả :n:timestamp nếu có).
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

// CalcJazoest tính jazoest checksum từ fb_dtsg
// Công thức: "2" + tổng ASCII code của tất cả ký tự trong dtsg
func CalcJazoest(fbDtsg string) string {
	sum := 0
	for _, c := range fbDtsg {
		sum += int(c)
	}
	return "2" + fmt.Sprintf("%d", sum)
}

type ActionHandler struct {
	sessionHandler *SessionHandler
	accountStore   accounts.Store
	fbDataStore    fbdata.Store
	configStore    *ConfigStore
}

// NewActionHandler tạo handler không có account store (backward compat)
func NewActionHandler(sessionHandler *SessionHandler) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   nil,
		fbDataStore:    nil,
		configStore:    nil,
	}
}

// NewActionHandlerWithStore tạo handler có account store thật
func NewActionHandlerWithStore(sessionHandler *SessionHandler, store accounts.Store, fbStore fbdata.Store) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   store,
		fbDataStore:    fbStore,
		configStore:    nil,
	}
}

// NewActionHandlerWithStoreAndConfig tao handler co account store va config store.
func NewActionHandlerWithStoreAndConfig(sessionHandler *SessionHandler, store accounts.Store, fbStore fbdata.Store, cfg *ConfigStore) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   store,
		fbDataStore:    fbStore,
		configStore:    cfg,
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
			ReactionID: req.ReactionID,
			DryRun:     req.DryRun, Message: "Post ID không được để trống.", ExecutedAt: now,
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
			ReactionID: req.ReactionID,
			DryRun:     req.DryRun, Message: fmt.Sprintf("Loại cảm xúc không hợp lệ: %s", reactionType), ExecutedAt: now,
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
			ReactionID: req.ReactionID,
			DryRun:     req.DryRun, Message: "Chưa chọn tài khoản thực thi.", ExecutedAt: now,
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
				ReactionID: req.ReactionID,
				DryRun:     req.DryRun,
				Message:    fmt.Sprintf("Tài khoản '%s' không tìm thấy trong local store.", accountID),
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
			ReactionID: req.ReactionID,
			DryRun:     req.DryRun,
			Message:    fmt.Sprintf("Phiên của tài khoản '%s' không hợp lệ (trạng thái: %s). Vui lòng kiểm tra lại phiên.", displayName, sessionStatusStr),
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
			ReactionID: req.ReactionID,
			DryRun:     true,
			Message:    fmt.Sprintf("[DRY RUN] Mô phỏng gửi '%s' (ID: %s) thành công — Tài khoản: %s — Bài viết: %s", reactionType, req.ReactionID, displayName, postID),
			ExecutedAt: time.Now(),
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	// ── 6. Real Run ───────────────────────────────────────────────────
	var cookie, fbDtsg, actorID string
	// Uu tien fbDataStore vi luong nay dang duoc dung on dinh cho comment/post.
	if h.fbDataStore != nil {
		if fbInfo, err := h.fbDataStore.Get(accountID); err == nil && fbInfo != nil {
			cookie = strings.TrimSpace(fbInfo.Info.Cookie)
			fbDtsg = strings.TrimSpace(fbInfo.Info.FbDtsg)
			actorID = strings.TrimSpace(fbInfo.UID)
		}
	}
	// Fallback sang accountStore neu fbDataStore khong co du lieu.
	if h.accountStore != nil {
		if acc, err := h.accountStore.Get(accountID); err == nil && acc != nil {
			if cookie == "" {
				cookie = strings.TrimSpace(acc.Cookie)
			}
			if fbDtsg == "" {
				fbDtsg = strings.TrimSpace(acc.FbDtsg)
			}
		}
	}
	// actor_id that phai la c_user trong cookie.
	if cookie != "" {
		if uid := strings.TrimSpace(ActorIDFromCookie(cookie)); uid != "" {
			actorID = uid
		}
	}
	if cookie == "" || fbDtsg == "" {
		resp := LikePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: reactionType, PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
			ReactionID: req.ReactionID,
			DryRun:     false,
			Message:    "Thieu Cookie hoac fb_dtsg cho tai khoan nay (nguon account/fbdata).",
			ExecutedAt: time.Now(),
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}
	if actorID == "" {
		resp := LikePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: reactionType, PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
			ReactionID: req.ReactionID,
			DryRun:     false,
			Message:    "Khong tim thay actor_id (c_user) trong cookie.",
			ExecutedAt: time.Now(),
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	// ── 6a. Pre-flight log — strategy is resolved inside executor ────────
	log.Printf("[RealRun] === Pre-flight ===")
	log.Printf("[RealRun]   strategy   = auto (comet/doc_id fallback)")
	log.Printf("[RealRun]   postID     = %q", postID)
	log.Printf("[RealRun]   reactionID = %q", req.ReactionID)
	log.Printf("[RealRun]   actorID    = %q", actorID)
	configuredDocID := h.GetReactDocID()
	if configuredDocID != "" {
		log.Printf("[RealRun]   doc_id     = configured in app settings")
	} else {
		log.Printf("[RealRun]   doc_id     = from env FB_REACT_DOC_ID (if set), else inline query")
	}
	fbResp, errCode, err := ExecuteFacebookReaction(cookie, fbDtsg, actorID, postID, req.ReactionID, configuredDocID)

	if err != nil {
		errorStatus := StatusValidationError
		if errCode != "" {
			errorStatus = string(errCode)
		}

		resp := LikePostResponse{
			Success: false, Status: errorStatus,
			Action: reactionType, PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
			ReactionID: req.ReactionID,
			DryRun:     false,
			ErrorCode:  errCode,
			Message:    fmt.Sprintf("[%s] %v", errCode, err),
			ExecutedAt: time.Now(),
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	resp := LikePostResponse{
		Success: true, Status: StatusSuccess,
		Action: reactionType, PostID: postID,
		AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
		ReactionID: req.ReactionID,
		DryRun:     false,
		Message:    "Tha cam xuc that thanh cong. Phan hoi FB: " + fbResp,
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
	if err := validateFeedbackID(postID); err != nil {
		resp := CommentPostResponse{
			Success: false, Status: StatusValidationError,
			Action: "comment", PostID: postID,
			AccountID: req.AccountID, CommentText: req.CommentText,
			DryRun:     req.DryRun,
			Message:    fmt.Sprintf("Feedback ID không hợp lệ: %v (gợi ý: dùng ID dạng base64 bắt đầu 'ZmVlZGJhY2s6').", err),
			ExecutedAt: now,
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
				DryRun:     req.DryRun,
				Message:    fmt.Sprintf("Tài khoản '%s' không tìm thấy trong local store.", accountID),
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
			DryRun:     req.DryRun,
			Message:    fmt.Sprintf("Phiên của tài khoản '%s' không hợp lệ (trạng thái: %s).", displayName, sessionStatusStr),
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
			DryRun:     true,
			Message:    fmt.Sprintf("[DRY RUN] Đã bình luận: \"%s\" bằng tài khoản %s — Bài viết: %s", snippet, displayName, postID),
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
			DryRun:     false,
			Message:    "Lỗi hệ thống: fbDataStore chưa được khởi tạo.",
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
			DryRun:     false,
			Message:    "Lỗi: Không tìm thấy Cookie hoặc fb_dtsg trong FB Data của tài khoản này.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	// Delay chống Rate Limit: comment thường bị FB siết mạnh hơn post.
	sleepMs := 5000 + rand.Intn(7001) // 5-12 giây
	fmt.Printf("[DELAY] Chờ %dms trước khi gửi request comment...\n", sleepMs)
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)

	reqURL := "https://www.facebook.com/api/graphql/"
	docID := "25720979764242405" // default
	if h.configStore != nil {
		if override := h.configStore.GetCommentDocID(); override != "" {
			docID = override
			fmt.Printf("[INFO] Comment dùng doc_id từ config: %s\n", docID)
		}
	}

	actorID := ActorIDFromCookie(fbData.Info.Cookie)
	if actorID == "" {
		actorID = strings.TrimSpace(fbData.UID) // fallback
	}
	if actorID == "" {
		resp := CommentPostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
			DryRun: false, Message: "Không tìm thấy actor_id (c_user) trong cookie.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	sessionData, sdErr := FetchSessionData(fbData.Info.Cookie)
	lsd := ""
	if sdErr == nil && strings.TrimSpace(sessionData.LSD) != "" {
		lsd = strings.TrimSpace(sessionData.LSD)
		fmt.Printf("[INFO] Comment LSD token: %s\n", lsd)
	} else {
		fmt.Printf("[WARN] Comment không lấy được LSD: %v\n", sdErr)
	}

	currentDtsg := strings.TrimSpace(fbData.Info.FbDtsg)
	if sdErr == nil && strings.TrimSpace(sessionData.DTSG) != "" {
		currentDtsg = strings.TrimSpace(sessionData.DTSG)
	}
	if currentDtsg == "" {
		resp := CommentPostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
			DryRun: false, Message: "Không có fb_dtsg hợp lệ để gửi comment.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	doCommentGraphQL := func(fbDtsg, lsdToken string) (int, string, error) {
		clientMutationID := fmt.Sprintf("%d", time.Now().UnixNano()%100000)
		idempotenceToken := "client:" + generateToken()
		sessionID := generateToken()
		attribution := fmt.Sprintf(
			"CometSinglePostDialogRoot.react,comet.post.single_dialog,via_cold_start,%d,%d,,,",
			time.Now().UnixMilli(),
			time.Now().UnixNano()%1000000,
		)

		variablesMap := map[string]interface{}{
			"feedLocation":   "POST_PERMALINK_DIALOG",
			"feedbackSource": 2,
			"groupID":        nil,
			"input": map[string]interface{}{
				"actor_id":           actorID,
				"client_mutation_id": clientMutationID,
				"attachments":        nil,
				"feedback_id":        postID,
				"formatting_style":   nil,
				"message": map[string]interface{}{
					"ranges": []interface{}{},
					"text":   commentText,
				},
				"reply_target_clicked":  false,
				"attribution_id_v2":     attribution,
				"vod_video_timestamp":   nil,
				"feedback_referrer":     "/",
				"is_tracking_encrypted": true,
				"tracking":              []string{},
				"feedback_source":       "OBJECT",
				"idempotence_token":     idempotenceToken,
				"session_id":            sessionID,
			},
			"inviteShortLinkKey": nil,
			"renderLocation":     nil,
			"scale":              1,
			"useDefaultActor":    false,
			"focusCommentID":     nil,
			"__relay_internal__pv__groups_comet_use_glvrelayprovider":                      false,
			"__relay_internal__pv__CometUFICommentActionLinksRewriteEnabledrelayprovider":  false,
			"__relay_internal__pv__CometUFICommentAvatarStickerAnimatedImagerelayprovider": false,
			"__relay_internal__pv__IsWorkUserrelayprovider":                                false,
			"__relay_internal__pv__CometUFICommentAutoTranslationTyperelayprovider":        "ORIGINAL",
		}
		variablesBytes, _ := json.Marshal(variablesMap)
		variables := string(variablesBytes)

		fd := url.Values{}
		fd.Set("doc_id", docID)
		fd.Set("av", actorID)
		fd.Set("__aaid", "0")
		fd.Set("__user", actorID)
		fd.Set("__a", "1")
		fd.Set("__req", "4d")
		fd.Set("dpr", "1")
		fd.Set("__ccg", "EXCELLENT")
		if strings.TrimSpace(sessionData.Rev) != "" {
			fd.Set("__rev", strings.TrimSpace(sessionData.Rev))
		}
		if strings.TrimSpace(sessionData.HSI) != "" {
			fd.Set("__hsi", strings.TrimSpace(sessionData.HSI))
		}
		fd.Set("fb_dtsg", fbDtsg)
		fd.Set("jazoest", CalcJazoest(fbDtsg))
		fd.Set("variables", variables)
		if lsdToken != "" {
			fd.Set("lsd", lsdToken)
		}
		if strings.TrimSpace(sessionData.SpinR) != "" {
			fd.Set("__spin_r", strings.TrimSpace(sessionData.SpinR))
		}
		fd.Set("__spin_b", "trunk")
		if strings.TrimSpace(sessionData.SpinT) != "" {
			fd.Set("__spin_t", strings.TrimSpace(sessionData.SpinT))
		}
		fd.Set("fb_api_caller_class", "RelayModern")
		fd.Set("fb_api_req_friendly_name", "useCometUFICreateCommentMutation")
		fd.Set("server_timestamps", "true")
		fd.Set("__comet_req", "15")

		httpReq, err := http.NewRequest("POST", reqURL, strings.NewReader(fd.Encode()))
		if err != nil {
			return 0, "", err
		}
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		httpReq.Header.Set("Cookie", fbData.Info.Cookie)
		httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/124.0.0.0")
		httpReq.Header.Set("Origin", "https://www.facebook.com")
		httpReq.Header.Set("Referer", "https://www.facebook.com/")
		httpReq.Header.Set("Accept", "*/*")
		httpReq.Header.Set("Sec-Fetch-Site", "same-origin")
		httpReq.Header.Set("Sec-Fetch-Mode", "cors")
		httpReq.Header.Set("Sec-Fetch-Dest", "empty")
		httpReq.Header.Set("X-FB-Friendly-Name", "useCometUFICreateCommentMutation")
		if lsdToken != "" {
			httpReq.Header.Set("X-FB-LSD", lsdToken)
		}
		httpReq.Header.Set("X-ASBD-ID", "129477")

		client := &http.Client{Timeout: 30 * time.Second}
		httpResp, err := client.Do(httpReq)
		if err != nil {
			return 0, "", err
		}
		defer httpResp.Body.Close()
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return httpResp.StatusCode, string(bodyBytes), nil
	}

	status1, body1, err := doCommentGraphQL(currentDtsg, lsd)
	if err != nil {
		resp := CommentPostResponse{
			Success: false, Status: "error", Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
			Message: fmt.Sprintf("Lỗi mạng khi gửi comment: %v", err), ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	fmt.Printf("\n=== GRAPHQL RAW RESP (Comment Attempt 1) ===\n%s\n====================\n\n", body1)

	if status1 != 200 {
		resp := CommentPostResponse{
			Success: false, Status: "error", Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
			Message:    fmt.Sprintf("GraphQL Server trả về mã lỗi %d", status1),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}
	if strings.Contains(body1, `"error":1357001`) || strings.Contains(body1, `"error":1357004`) {
		resp := CommentPostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
			DryRun:     false,
			Message:    "Facebook trả về lỗi đăng nhập (1357001/1357004). Cookie hoặc fb_dtsg không còn hợp lệ/không đồng bộ với phiên hiện tại.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	if strings.Contains(body1, `"comment_create":{`) {
		resp := CommentPostResponse{
			Success: true, Status: StatusSuccess,
			Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
			DryRun:     false,
			Message:    "Real Run: Gửi payload GraphQL comment thành công!",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	if strings.Contains(body1, `"code":1675004`) {
		backoffMs := 15000 + rand.Intn(15001) // 15-30 giây
		fmt.Printf("[DELAY] Rate limit comment (1675004), chờ %dms rồi thử lại...\n", backoffMs)
		time.Sleep(time.Duration(backoffMs) * time.Millisecond)

		// Làm mới token trước khi retry.
		if refreshed, rErr := FetchSessionData(fbData.Info.Cookie); rErr == nil {
			if strings.TrimSpace(refreshed.LSD) != "" {
				lsd = strings.TrimSpace(refreshed.LSD)
			}
			if strings.TrimSpace(refreshed.DTSG) != "" {
				currentDtsg = strings.TrimSpace(refreshed.DTSG)
			}
		}

		_, body2, err2 := doCommentGraphQL(currentDtsg, lsd)
		if err2 != nil {
			resp := CommentPostResponse{
				Success: false, Status: "error", Action: "comment", PostID: postID,
				AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
				Message:    fmt.Sprintf("Rate limit và retry lỗi mạng: %v", err2),
				ExecutedAt: time.Now(),
			}
			AddLog(logFromComment(resp, resp.Message))
			return resp
		}

		fmt.Printf("\n=== GRAPHQL RAW RESP (Comment Attempt 2) ===\n%s\n====================\n\n", body2)
		if strings.Contains(body2, `"comment_create":{`) {
			resp := CommentPostResponse{
				Success: true, Status: StatusSuccess,
				Action: "comment", PostID: postID,
				AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
				DryRun:     false,
				Message:    "Retry sau rate limit: comment thành công.",
				ExecutedAt: time.Now(),
			}
			AddLog(logFromComment(resp, resp.Message))
			return resp
		}

		resp := CommentPostResponse{
			Success: false, Status: "error", Action: "comment", PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
			Message:    "Facebook đang giới hạn bình luận (code 1675004). Hãy tăng delay và thử lại sau.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromComment(resp, resp.Message))
		return resp
	}

	// Nếu FB trả về dtsgToken mới trong lỗi, thử lại 1 lần.
	freshDtsg := extractDtsgFromError(body1)
	if freshDtsg != "" {
		time.Sleep(1200 * time.Millisecond)
		_, body2, err2 := doCommentGraphQL(freshDtsg, lsd)
		if err2 == nil && strings.Contains(body2, `"comment_create":{`) {
			resp := CommentPostResponse{
				Success: true, Status: StatusSuccess,
				Action: "comment", PostID: postID,
				AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText,
				DryRun:     false,
				Message:    "Retry với dtsgToken mới: comment thành công.",
				ExecutedAt: time.Now(),
			}
			AddLog(logFromComment(resp, resp.Message))
			return resp
		}
	}

	snippet := body1
	if len(snippet) > 220 {
		snippet = snippet[:220] + "..."
	}
	resp := CommentPostResponse{
		Success: false, Status: "error", Action: "comment", PostID: postID,
		AccountID: accountID, AccountDisplayName: displayName, CommentText: commentText, DryRun: false,
		Message:    fmt.Sprintf("FB GraphQL Error: %s", snippet),
		ExecutedAt: time.Now(),
	}
	AddLog(logFromComment(resp, resp.Message))
	return resp
}

func (h *ActionHandler) GetActionLogs() []ActionLog {
	return GetLogs()
}

// GetReactDocID lay doc_id tu local config, fallback ve env.
func (h *ActionHandler) GetReactDocID() string {
	if h.configStore != nil {
		if v := h.configStore.GetReactDocID(); v != "" {
			return v
		}
	}
	return strings.TrimSpace(os.Getenv("FB_REACT_DOC_ID"))
}

// UpdateReactDocID cap nhat doc_id local.
// Truyen chuoi rong de xoa doc_id da luu.
func (h *ActionHandler) UpdateReactDocID(docID string) (string, error) {
	if h.configStore == nil {
		return "", errors.New("config store chua duoc khoi tao")
	}
	if err := h.configStore.SetReactDocID(docID); err != nil {
		return "", err
	}
	return h.configStore.GetReactDocID(), nil
}

// GetCreatePostDocID lay doc_id cho ComposerStoryCreateMutation tu config.
func (h *ActionHandler) GetCreatePostDocID() string {
	if h.configStore != nil {
		if v := h.configStore.GetCreatePostDocID(); v != "" {
			return v
		}
	}
	return ""
}

// UpdateCreatePostDocID cap nhat create_post doc_id.
// Lay doc_id moi tu F12 > Network > ComposerStoryCreateMutation > Request Payload > doc_id.
// Truyen chuoi rong de reset ve default.
func (h *ActionHandler) UpdateCreatePostDocID(docID string) (string, error) {
	if h.configStore == nil {
		return "", errors.New("config store chua duoc khoi tao")
	}
	if err := h.configStore.SetCreatePostDocID(docID); err != nil {
		return "", err
	}
	return h.configStore.GetCreatePostDocID(), nil
}

// GetCommentDocID lay doc_id cho useCometUFICreateCommentMutation tu config.
func (h *ActionHandler) GetCommentDocID() string {
	if h.configStore != nil {
		if v := h.configStore.GetCommentDocID(); v != "" {
			return v
		}
	}
	return ""
}

// UpdateCommentDocID cap nhat comment doc_id.
// Lay doc_id moi tu F12 > Network > useCometUFICreateCommentMutation > Request Payload > doc_id.
// Truyen chuoi rong de reset ve default.
func (h *ActionHandler) UpdateCommentDocID(docID string) (string, error) {
	if h.configStore == nil {
		return "", errors.New("config store chua duoc khoi tao")
	}
	if err := h.configStore.SetCommentDocID(docID); err != nil {
		return "", err
	}
	return h.configStore.GetCommentDocID(), nil
}

// GetScrapeDocID lay doc_id cho scrape/crawl post tu config.
func (h *ActionHandler) GetScrapeDocID() string {
	if h.configStore != nil {
		if v := h.configStore.GetScrapeDocID(); v != "" {
			return v
		}
	}
	return ""
}

// UpdateScrapeDocID cap nhat scrape doc_id.
// Lay doc_id moi tu F12 > Network > ProfileCometTimelineFeed > Request Payload > doc_id.
// Truyen chuoi rong de reset ve default.
func (h *ActionHandler) UpdateScrapeDocID(docID string) (string, error) {
	if h.configStore == nil {
		return "", errors.New("config store chua duoc khoi tao")
	}
	if err := h.configStore.SetScrapeDocID(docID); err != nil {
		return "", err
	}
	return h.configStore.GetScrapeDocID(), nil
}


// toLog chuyển response thành ActionLog entry
func toLog(resp LikePostResponse, msg string) ActionLog {
	return ActionLog{
		ActionType:         resp.Action,
		AccountID:          resp.AccountID,
		AccountDisplayName: resp.AccountDisplayName,
		PostID:             resp.PostID,
		ReactionType:       resp.ReactionType,
		ReactionID:         resp.ReactionID,
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
	if len(req.ImagePaths) == 0 {
		fmt.Printf("[INFO] CreatePost payload: image_paths=0 (frontend khong gui media)\n")
	} else {
		fmt.Printf("[INFO] CreatePost payload: image_paths=%d -> %v\n", len(req.ImagePaths), req.ImagePaths)
	}

	postText := strings.TrimSpace(req.PostText)
	photoIDs, localMediaInputs, ignoredMediaInputs := extractPhotoIDsFromImagePaths(req.ImagePaths)
	hasMediaAttachments := len(photoIDs) > 0
	if postText == "" && !hasMediaAttachments && len(localMediaInputs) == 0 {
		msg := "Nội dung đăng bài không được để trống."
		if len(req.ImagePaths) > 0 {
			msg = "Bai dang can co text hoac media id hop le (photo_id/media_fbid)."
		}
		resp := CreatePostResponse{
			Success: false, Status: StatusValidationError,
			Action: "create_post", AccountID: req.AccountID, PostText: postText,
			DryRun: req.DryRun, Message: msg, ExecutedAt: now,
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
				DryRun:     req.DryRun,
				Message:    fmt.Sprintf("Tài khoản '%s' không tìm thấy trong local store.", accountID),
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
			DryRun:     req.DryRun,
			Message:    fmt.Sprintf("Phiên của tài khoản '%s' không hợp lệ (trạng thái: %s).", displayName, sessionStatusStr),
			ExecutedAt: now,
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}
	if len(ignoredMediaInputs) > 0 {
		fmt.Printf("[WARN] Bỏ qua %d media input không parse được photo_id: %v\n", len(ignoredMediaInputs), ignoredMediaInputs)
	}
	if len(localMediaInputs) > 0 {
		fmt.Printf("[INFO] Có %d media local sẽ upload để lấy photo_id: %v\n", len(localMediaInputs), localMediaInputs)
	}

	if req.DryRun {
		sleepMs := 500 + rand.Intn(701)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		snippet := "[MEDIA ONLY]"
		if postText != "" {
			snippet = postText
			if len(snippet) > 20 {
				snippet = snippet[:20] + "..."
			}
		}
		mediaInfo := ""
		if hasMediaAttachments {
			mediaInfo = fmt.Sprintf(" | media_id=%d", len(photoIDs))
		}

		resp := CreatePostResponse{
			Success: true, Status: StatusSuccess,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun:     true,
			Message:    fmt.Sprintf("[DRY RUN] Đã đăng bài nháp: \"%s\" bằng tài khoản %s%s", snippet, displayName, mediaInfo),
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
			DryRun:     false,
			Message:    fmt.Sprintf("Lỗi: %s - Không tìm thấy dữ liệu FB (Cookie/fb_dtsg). Cập nhật vào facebook_data.json!", err),
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

	// Uu tien doc_id moi (user capture thanh cong), sau do den doc_id cu de fallback.
	// Neu user set doc_id trong settings thi van duoc uu tien.
	preferredCreateDocIDs := []string{
		"8508448842550728",  // doc_id mới nhất (2025-04)
		"7993695614062380",  // doc_id mới (2025-Q1)
		"26400547339631027", // doc_id user đã capture thành công (cũ hơn)
		"27581837698072404", // fallback cũ
	}
	docIDCandidates := make([]string, 0, 3)
	seenDocID := make(map[string]struct{})
	appendDocID := func(id string) {
		id = strings.TrimSpace(id)
		if id == "" {
			return
		}
		if _, ok := seenDocID[id]; ok {
			return
		}
		seenDocID[id] = struct{}{}
		docIDCandidates = append(docIDCandidates, id)
	}
	if h.configStore != nil {
		appendDocID(h.configStore.GetCreatePostDocID())
	}
	for _, id := range preferredCreateDocIDs {
		appendDocID(id)
	}
	if len(docIDCandidates) == 0 {
		docIDCandidates = append(docIDCandidates, "26400547339631027")
	}
	currentDocID := docIDCandidates[0]
	fmt.Printf("[INFO] CreatePost doc_id đang dùng: %s\n", currentDocID)

	attachmentsJSON := "[]"
	mediaNote := ""
	attachmentVariants := []attachmentVariant{{Name: "none", JSON: "[]"}}
	currentAttachmentVariantIdx := 0
	attachmentMode := "none"
	composerSessionID := generateToken()
	idempotenceToken := fmt.Sprintf("%s_FEED", composerSessionID)
	messageTextForMutation := postText

	// Helper: xây dựng variables JSON CHÍNH XÁC theo browser (captured từ F12)
	buildVariables := func() string {
		// Facebook ComposerStoryCreateMutation expects MessageInput object (non-nullable).
		// null → field_exception 1357010.
		// Omit → field_exception 1357010.
		messageField := `"message":{"ranges":[],"text":""},`
		if strings.TrimSpace(messageTextForMutation) != "" {
			messageField = fmt.Sprintf(`"message":{"ranges":[],"text":%s},`, JsonStr(messageTextForMutation))
		}

		return fmt.Sprintf(`{`+
			`"input":{`+
			`"composer_entry_point":"inline_composer",`+
			`"composer_source_surface":"timeline",`+
			`"idempotence_token":%s,`+
			`"source":"WWW",`+
			`"attachments":%s,`+
			`"audience":{"privacy":{"allow":[],"base_state":"EVERYONE","deny":[],"tag_expansion_state":"UNSPECIFIED"}},`+
			`%s`+ // messageField — conditional: có text thì gửi, không có thì bỏ qua
			`"with_tags_ids":null,`+
			`"inline_activities":[],`+
			`"text_format_preset_id":"0",`+
			`"publishing_flow":{"supported_flows":["ASYNC_SILENT","ASYNC_NOTIF","FALLBACK"]},`+
			`"post_publish_story_data":{"reshare_post_as_sticker":"DISABLED"},`+
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
			`"__relay_internal__pv__CometUFISingleLineUFIrelayprovider":true,`+
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
			`"__relay_internal__pv__GHLShouldChangeSponsoredAuctionDistanceFieldNamerelayprovider":false,`+
			`"__relay_internal__pv__GHLShouldUseSponsoredAuctionLabelFieldNameV1relayprovider":false,`+
			`"__relay_internal__pv__GHLShouldUseSponsoredAuctionLabelFieldNameV2relayprovider":false`+
			`}`,
			JsonStr(idempotenceToken), attachmentsJSON, messageField, JsonStr(composerSessionID), JsonStr(actorID))
	}

	// Lấy LSD token (cần thiết cho Comet GraphQL endpoint)
	lsd := ""
	sessionData, sdErr := FetchSessionData(fbInfo.Info.Cookie)
	if sdErr == nil && sessionData.LSD != "" {
		lsd = sessionData.LSD
		fmt.Printf("[INFO] LSD token: %s\n", lsd)
	} else {
		fmt.Printf("[WARN] Không lấy được LSD: %v\n", sdErr)
	}

	// Helper: thực hiện một lần POST lên GraphQL

	currentDtsg := strings.TrimSpace(fbInfo.Info.FbDtsg)
	if sdErr == nil && strings.TrimSpace(sessionData.DTSG) != "" {
		currentDtsg = strings.TrimSpace(sessionData.DTSG)
	}
	if currentDtsg == "" {
		resp := CreatePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Khong co fb_dtsg hop le de dang bai.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	if len(localMediaInputs) > 0 {
		for _, localInput := range localMediaInputs {
			uploadedPhotoID, upErr := uploadOneLocalMedia(
				fbInfo.Info.Cookie,
				actorID,
				currentDtsg,
				lsd,
				sessionData,
				localInput,
				composerSessionID,
			)
			if upErr != nil {
				resp := CreatePostResponse{
					Success: false, Status: StatusFailed,
					Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
					DryRun: false, Message: fmt.Sprintf("Upload media local that bai (%s): %v", localInput, upErr),
					ExecutedAt: time.Now(),
				}
				AddLog(logFromPost(resp, resp.Message))
				return resp
			}
			photoIDs = append(photoIDs, uploadedPhotoID)
		}
	}

	hasMediaAttachments = len(photoIDs) > 0
	if postText == "" && !hasMediaAttachments {
		resp := CreatePostResponse{
			Success: false, Status: StatusValidationError,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Khong lay duoc photo_id tu media da chon.",
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	if hasMediaAttachments {
		attachmentVariants = buildAttachmentVariants(photoIDs)
		currentAttachmentVariantIdx = 0
		attachmentsJSON = attachmentVariants[currentAttachmentVariantIdx].JSON
		attachmentMode = attachmentVariants[currentAttachmentVariantIdx].Name
		fmt.Printf("[INFO] Attachment mode=%s\n", attachmentMode)
	}
	if hasMediaAttachments {
		mediaNote += fmt.Sprintf(" | media_id=%d", len(photoIDs))
	}
	if len(ignoredMediaInputs) > 0 {
		mediaNote += fmt.Sprintf(" | bo_qua_media=%d", len(ignoredMediaInputs))
	}
	if len(localMediaInputs) > 0 {
		mediaNote += fmt.Sprintf(" | local_media=%d", len(localMediaInputs))
	}
	attachmentModeSuffix := func() string {
		if !hasMediaAttachments {
			return ""
		}
		return " | attach_mode=" + attachmentMode
	}
	doGraphQLPost := func(fbDtsg string) (int, string, error) {
		vars := buildVariables()
		fd := url.Values{}
		fd.Set("av", actorID)
		fd.Set("__aaid", "0")
		fd.Set("__user", actorID)
		fd.Set("__a", "1")
		fd.Set("__req", "28")
		if strings.TrimSpace(sessionData.HS) != "" {
			fd.Set("__hs", strings.TrimSpace(sessionData.HS))
		}
		if strings.TrimSpace(sessionData.S) != "" {
			fd.Set("__s", strings.TrimSpace(sessionData.S))
		}
		if strings.TrimSpace(sessionData.Dyn) != "" {
			fd.Set("__dyn", strings.TrimSpace(sessionData.Dyn))
		}
		if strings.TrimSpace(sessionData.CSR) != "" {
			fd.Set("__csr", strings.TrimSpace(sessionData.CSR))
		}
		fd.Set("dpr", "1")
		fd.Set("__ccg", "EXCELLENT")
		if strings.TrimSpace(sessionData.Rev) != "" {
			fd.Set("__rev", strings.TrimSpace(sessionData.Rev))
		}
		if strings.TrimSpace(sessionData.HSI) != "" {
			fd.Set("__hsi", strings.TrimSpace(sessionData.HSI))
		}
		fd.Set("fb_dtsg", fbDtsg) // Gửi FULL token (kể cả :3:timestamp)
		fd.Set("jazoest", CalcJazoest(fbDtsg))
		if strings.TrimSpace(sessionData.SpinR) != "" {
			fd.Set("__spin_r", strings.TrimSpace(sessionData.SpinR))
		}
		fd.Set("__spin_b", "trunk")
		if strings.TrimSpace(sessionData.SpinT) != "" {
			fd.Set("__spin_t", strings.TrimSpace(sessionData.SpinT))
		}
		fd.Set("__comet_req", "15")
		fd.Set("__crn", "comet.fbweb.CometProfileTimelineListViewRoute")
		fd.Set("lsd", lsd)
		fd.Set("doc_id", currentDocID)
		fd.Set("fb_api_caller_class", "RelayModern")
		fd.Set("fb_api_req_friendly_name", "ComposerStoryCreateMutation")
		fd.Set("server_timestamps", "true")
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
	currentDtsg = strings.TrimSpace(currentDtsg)
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
			DryRun: false, Message: fmt.Sprintf("Đăng bài thành công lần 1 — Tài khoản: %s!%s%s", displayName, mediaNote, attachmentModeSuffix()),
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	// Fallback attachment shape cho lỗi field_exception (1357010).
	if hasMediaAttachments && hasGraphQLErrorCode(body1, "1357010") && len(attachmentVariants) > 1 {
		for idx := currentAttachmentVariantIdx + 1; idx < len(attachmentVariants); idx++ {
			currentAttachmentVariantIdx = idx
			attachmentsJSON = attachmentVariants[currentAttachmentVariantIdx].JSON
			attachmentMode = attachmentVariants[currentAttachmentVariantIdx].Name
			fmt.Printf("[WARN] GraphQL code 1357010. Retry with attachment mode=%s\n", attachmentMode)

			statusAlt, bodyAlt, errAlt := doGraphQLPost(currentDtsg)
			if errAlt != nil {
				fmt.Printf("[WARN] Retry mode=%s failed by HTTP error: %v\n", attachmentMode, errAlt)
				continue
			}
			snippetAlt := bodyAlt
			if len(snippetAlt) > 300 {
				snippetAlt = snippetAlt[:300]
			}
			fmt.Printf("=== RESP Retry mode=%s (status=%d, len=%d) ===\n%s\n====================\n",
				attachmentMode, statusAlt, len(bodyAlt), snippetAlt)

			if bodyAlt != "" && !strings.Contains(bodyAlt, `"error":`) && !strings.Contains(bodyAlt, `"errors":[`) {
				resp := CreatePostResponse{
					Success: true, Status: StatusSuccess,
					Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
					DryRun: false, Message: fmt.Sprintf("Đăng bài thành công (fallback attachment) — Tài khoản: %s!%s%s", displayName, mediaNote, attachmentModeSuffix()),
					ExecutedAt: time.Now(),
				}
				AddLog(logFromPost(resp, resp.Message))
				return resp
			}

			status1 = statusAlt
			body1 = bodyAlt
			snippet1 = snippetAlt
			if !hasGraphQLErrorCode(body1, "1357010") {
				break
			}
		}
	}

	// Neu van 1357010, thu doc_id khac (payload user capture thanh cong la 26400547339631027).
	if (hasGraphQLErrorCode(body1, "1357010") || hasGraphQLErrorCode(body1, "1675030")) && len(docIDCandidates) > 1 {
		for d := 1; d < len(docIDCandidates); d++ {
			currentDocID = docIDCandidates[d]
			fmt.Printf("[WARN] GraphQL code 1357010. Retry with another doc_id=%s\n", currentDocID)

			// Reset state ve mode photo + message goc de bam sat payload tay.
			messageTextForMutation = postText
			if hasMediaAttachments && len(attachmentVariants) > 0 {
				currentAttachmentVariantIdx = 0
				attachmentsJSON = attachmentVariants[currentAttachmentVariantIdx].JSON
				attachmentMode = attachmentVariants[currentAttachmentVariantIdx].Name
			}

			statusDoc, bodyDoc, errDoc := doGraphQLPost(currentDtsg)
			if errDoc != nil {
				fmt.Printf("[WARN] Retry doc_id=%s failed by HTTP error: %v\n", currentDocID, errDoc)
				continue
			}
			snippetDoc := bodyDoc
			if len(snippetDoc) > 300 {
				snippetDoc = snippetDoc[:300]
			}
			fmt.Printf("=== RESP Retry doc_id=%s (status=%d, len=%d) ===\n%s\n====================\n",
				currentDocID, statusDoc, len(bodyDoc), snippetDoc)

			if bodyDoc != "" && !strings.Contains(bodyDoc, `"error":`) && !strings.Contains(bodyDoc, `"errors":[`) {
				resp := CreatePostResponse{
					Success: true, Status: StatusSuccess,
					Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
					DryRun: false, Message: fmt.Sprintf("Đăng bài thành công (doc_id fallback) — Tài khoản: %s!%s%s", displayName, mediaNote, attachmentModeSuffix()),
					ExecutedAt: time.Now(),
				}
				AddLog(logFromPost(resp, resp.Message))
				return resp
			}

			// Thu attachment variants tren doc_id moi.
			if hasMediaAttachments && hasGraphQLErrorCode(bodyDoc, "1357010") && len(attachmentVariants) > 1 {
				for idx := 1; idx < len(attachmentVariants); idx++ {
					currentAttachmentVariantIdx = idx
					attachmentsJSON = attachmentVariants[currentAttachmentVariantIdx].JSON
					attachmentMode = attachmentVariants[currentAttachmentVariantIdx].Name
					fmt.Printf("[WARN] doc_id=%s retry with attachment mode=%s\n", currentDocID, attachmentMode)

					statusAlt, bodyAlt, errAlt := doGraphQLPost(currentDtsg)
					if errAlt != nil {
						fmt.Printf("[WARN] doc_id=%s mode=%s failed by HTTP error: %v\n", currentDocID, attachmentMode, errAlt)
						continue
					}
					snippetAlt := bodyAlt
					if len(snippetAlt) > 300 {
						snippetAlt = snippetAlt[:300]
					}
					fmt.Printf("=== RESP Retry doc_id=%s mode=%s (status=%d, len=%d) ===\n%s\n====================\n",
						currentDocID, attachmentMode, statusAlt, len(bodyAlt), snippetAlt)

					if bodyAlt != "" && !strings.Contains(bodyAlt, `"error":`) && !strings.Contains(bodyAlt, `"errors":[`) {
						resp := CreatePostResponse{
							Success: true, Status: StatusSuccess,
							Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
							DryRun: false, Message: fmt.Sprintf("Đăng bài thành công (doc_id + attachment fallback) — Tài khoản: %s!%s%s", displayName, mediaNote, attachmentModeSuffix()),
							ExecutedAt: time.Now(),
						}
						AddLog(logFromPost(resp, resp.Message))
						return resp
					}

					statusDoc = statusAlt
					bodyDoc = bodyAlt
					snippetDoc = snippetAlt
					if !hasGraphQLErrorCode(bodyDoc, "1357010") {
						break
					}
				}
			}

			status1 = statusDoc
			body1 = bodyDoc
			snippet1 = snippetDoc
			if !hasGraphQLErrorCode(body1, "1357010") && !hasGraphQLErrorCode(body1, "1675030") {
				break
			}
		}
	}

	// ── Lần 2: Chỉ retry khi lỗi DTSG/session hết hạn (1357004, 1357032) ──────
	// Các lỗi khác (1357010=field_exception, v.v.) là lỗi THẬT từ FB — báo thẳng.
	isDtsgExpiry := strings.Contains(body1, `"error":1357004`) ||
		strings.Contains(body1, `"error":1357032`)

	if !isDtsgExpiry {
		// Lỗi thật từ Facebook — trả về ngay, không retry vô nghĩa
		resp := CreatePostResponse{
			Success: false, Status: StatusFailed,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Facebook từ chối bài đăng: " + snippet1,
			ExecutedAt: time.Now(),
		}
		AddLog(logFromPost(resp, resp.Message))
		return resp
	}

	// Chỉ đến đây khi lỗi là DTSG hết hạn — thử lấy token mới từ body lỗi
	freshDtsg := extractDtsgFromError(body1)
	if freshDtsg == "" {
		resp := CreatePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: "create_post", AccountID: accountID, AccountDisplayName: displayName, PostText: postText,
			DryRun: false, Message: "Phiên DTSG đã hết hạn nhưng không lấy được token mới. Cập nhật lại fb_dtsg trong facebook_data.json. Body: " + snippet1,
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
		DryRun: false, Message: fmt.Sprintf("Đăng bài thành công (Two-Step DTSG) — Tài khoản: %s!%s%s", displayName, mediaNote, attachmentModeSuffix()),
		ExecutedAt: time.Now(),
	}
	AddLog(logFromPost(resp, resp.Message))
	return resp
}

// SanitizeCookie loại bỏ các ký tự xuống dòng và khoảng trắng thừa trong cookie
func SanitizeCookie(cookie string) string {
	cookie = strings.ReplaceAll(cookie, "\n", "")
	cookie = strings.ReplaceAll(cookie, "\r", "")
	return strings.TrimSpace(cookie)
}
