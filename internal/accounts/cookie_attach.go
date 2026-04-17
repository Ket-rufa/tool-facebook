package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	cUserRegex   = regexp.MustCompile(`c_user=([0-9]+)`)
	fbDtsgRegex  = regexp.MustCompile(`name="fb_dtsg" value="(.*?)"`)
	fbDtsgRegex2 = regexp.MustCompile(`fb_dtsg" value="(.*?)"`)
	fbDtsgRegex3 = regexp.MustCompile(`"DTSGInitialData",\[\],\{"token":"(.*?)"\}`)
	fbDtsgRegex4 = regexp.MustCompile(`"fb_dtsg":"(.*?)"`)
)

// ProcessCookieAttachFlow nhận cookie, lấy thử uid và cào fb_dtsg, trả về AttachFlowStatusResponse
func (s *AccountService) ProcessCookieAttachFlow(cookieValue string) (AttachFlowStatusResponse, error) {
	cookieValue = strings.TrimSpace(cookieValue)
	if cookieValue == "" {
		return AttachFlowStatusResponse{}, errors.New("cookie không được để trống")
	}

	// Hỗ trợ parse JSON của J2TEAM Cookie
	if strings.HasPrefix(cookieValue, "[") || strings.HasPrefix(cookieValue, "{") {
		var j2Cookies []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}
		
		// Thử parse dạng Array trực tiếp: [{"name":"...", "value":"..."}]
		err := json.Unmarshal([]byte(cookieValue), &j2Cookies)
		
		if err != nil {
			// Thử parse dạng Object (bản update của J2TEAM): {"url":"...", "cookies":[{"name":"...", "value":"..."}]}
			var j2Object struct {
				Cookies []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"cookies"`
			}
			if errObj := json.Unmarshal([]byte(cookieValue), &j2Object); errObj == nil {
				j2Cookies = j2Object.Cookies
			}
		}

		if len(j2Cookies) > 0 {
			var builder strings.Builder
			for _, c := range j2Cookies {
				builder.WriteString(fmt.Sprintf("%s=%s; ", c.Name, c.Value))
			}
			cookieValue = builder.String()
		}
	}

	// Dọn dẹp newline nếu người dùng copy dư
	cookieValue = strings.ReplaceAll(cookieValue, "\n", "")
	cookieValue = strings.ReplaceAll(cookieValue, "\r", "")

	// 1. Phân tích c_user (AccountID)
	matches := cUserRegex.FindStringSubmatch(cookieValue)
	if len(matches) < 2 {
		return AttachFlowStatusResponse{}, errors.New("không tìm thấy c_user trong Cookie (Cookie không hợp lệ hoặc bị thiếu)")
	}
	uidExtracted := matches[1]

	// 2. Fetch www.facebook.com để lấy fb_dtsg (Dùng bản Desktop vì User-Agent đang là Desktop)
	req, err := http.NewRequest("GET", "https://www.facebook.com/", nil)
	if err != nil {
		return AttachFlowStatusResponse{}, fmt.Errorf("lỗi tạo request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookieValue)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return AttachFlowStatusResponse{}, fmt.Errorf("lỗi mạng khi cào facebook: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return AttachFlowStatusResponse{}, fmt.Errorf("lỗi đọc nội dung: %v", err)
	}
	bodyStr := string(bodyBytes)

	// Có thể parse luôn URL hiện tại nếu là trang đăng nhập
	if strings.Contains(resp.Request.URL.String(), "login") || strings.Contains(resp.Request.URL.String(), "checkpoint") {
		return AttachFlowStatusResponse{}, errors.New("cookie bị từ chối / đã chết (bị ép về trang đăng nhập hoặc Checkpoint)")
	}

	// 3. Regex lấy fb_dtsg
	var fbDtsg string
	if match := fbDtsgRegex.FindStringSubmatch(bodyStr); len(match) >= 2 {
		fbDtsg = match[1]
	} else if match := fbDtsgRegex2.FindStringSubmatch(bodyStr); len(match) >= 2 {
		fbDtsg = match[1]
	} else if match := fbDtsgRegex3.FindStringSubmatch(bodyStr); len(match) >= 2 {
		fbDtsg = match[1]
	} else if match := fbDtsgRegex4.FindStringSubmatch(bodyStr); len(match) >= 2 {
		fbDtsg = match[1]
	} else {
		// Parse thử cái Title để xem Facebook trả về trang gì
		titleRegex := regexp.MustCompile(`<title>(.*?)</title>`)
		titleMatch := titleRegex.FindStringSubmatch(bodyStr)
		title := "Không rõ"
		if len(titleMatch) >= 2 {
			title = titleMatch[1]
		}
		return AttachFlowStatusResponse{}, fmt.Errorf("không trích xuất được fb_dtsg. FB trả về trang: [%s]", title)
	}

	// 4. Tạo Flow & Mock Preview thành công
	flow, _ := s.flowManager.Start()
	flowID := flow.FlowID

	existing, _ := s.store.List()
	cloneIndex := len(existing) + 1
	displayName := fmt.Sprintf("Clone %d", cloneIndex)

	preview := &AccountProfile{
		ID:            "acc_" + uidExtracted,
		AccountID:     uidExtracted,
		DisplayName:   displayName,
		Avatar:        fmt.Sprintf("https://graph.facebook.com/%s/picture?type=large", uidExtracted), // Avatar chuẩn FB
		AccountType:   string(AccountTypeProfile),
		Provider:      "Facebook",
		SessionID:     "sess_" + uidExtracted,
		SessionStatus: string(SessionActive),
		AttachedAt:    time.Now().Format(time.RFC3339),
		LastCheckedAt: time.Now().Format("2006-01-02 15:04:05"),
		IsDefault:     false,
		Note:          "Thêm qua Cookie thủ công. Có chứa fb_dtsg: " + fbDtsg[:10] + "...",
		Cookie:        cookieValue,
		FbDtsg:        fbDtsg,
	}

	s.flowManager.SetAuthenticated(flowID, preview)

	return AttachFlowStatusResponse{
		FlowID:         flowID,
		State:          string(FlowAuthenticated),
		AccountPreview: preview,
	}, nil
}
