package accounts

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pquerna/otp/totp"
	"regexp"
)

// RequestAutomator xử lý đăng nhập qua HTTP Requests
type RequestAutomator struct {
	client *http.Client
}

func NewRequestAutomator() *RequestAutomator {
	jar, _ := cookiejar.New(nil)
	return &RequestAutomator{
		client: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil // Cho phép auto redirect
			},
		},
	}
}

// RequestLoginResult chứa kết quả đăng nhập
type RequestLoginResult struct {
	Success bool
	Message string
	UID     string
	Name    string
	Cookie  string
	Dtsg    string
	Need2FA bool
}

// Login performs the multi-step login process via mbasic.facebook.com
func (a *RequestAutomator) Login(email, pass, twoFactorKey string) (*RequestLoginResult, error) {
	// 1. GET trang login để lấy tokens (lsd, jazoest)
	reqInit, _ := http.NewRequest("GET", "https://mbasic.facebook.com/login.php", nil)
	a.setCommonHeaders(reqInit)
	resp, err := a.client.Do(reqInit)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Facebook: %v", err)
	}
	defer resp.Body.Close()

	// Đọc nội dung trang đầu để debug nếu cần
	initBody, _ := io.ReadAll(resp.Body)
	initHTML := string(initBody)

	// Nếu bị đá sang trang chọn ngôn ngữ, hãy follow nó
	currURL := resp.Request.URL.String()
	if strings.Contains(currURL, "/intl/save_locale/") {
		fmt.Printf("[LOGIN-REQ] Redirected to locale: %s\n", currURL)
		respLocale, errLocale := a.client.Get(currURL)
		if errLocale == nil {
			defer respLocale.Body.Close()
			b, _ := io.ReadAll(respLocale.Body)
			initHTML = string(b)
		}
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(initHTML))
	if err != nil {
		return nil, fmt.Errorf("không thể phân tích trang login: %v", err)
	}

	// Tìm chính xác bảng chứa ô nhập password (ô login chính)
	form := doc.Find("form").Has("input[type='password']").First()
	if form.Length() == 0 {
		form = doc.Find("form").Has("input[name='email']").First()
	}
	if form.Length() == 0 {
		form = doc.Find("form").First()
	}

	action, _ := form.Attr("action")
	if action == "" || action == "#" {
		action = "https://mbasic.facebook.com/login.php"
	}
	if !strings.HasPrefix(action, "http") {
		action = "https://mbasic.facebook.com" + action
	}

	fmt.Printf("[LOGIN-REQ] Found form with action: %s\n", action)

	var postDataParts []string
	
	// Thu thập các trường theo thứ tự xuất hiện trong HTML (rất quan trọng)
	form.Find("input").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		val, _ := s.Attr("value")
		typ, _ := s.Attr("type")
		
		if name == "" { return }

		finalVal := val
		if typ == "password" {
			finalVal = pass
		} else if name == "email" {
			finalVal = email
		}

		postDataParts = append(postDataParts, fmt.Sprintf("%s=%s", url.QueryEscape(name), url.QueryEscape(finalVal)))
	})

	// Đảm bảo có nút login
	hasLogin := false
	for _, p := range postDataParts {
		if strings.HasPrefix(p, "login=") { hasLogin = true; break }
	}
	if !hasLogin {
		postDataParts = append(postDataParts, "login="+url.QueryEscape("Log In"))
	}

	bodyStr := strings.Join(postDataParts, "&")

	// DEBUG: Xem toàn bộ dữ liệu gửi đi (ẩn mật khẩu)
	debugBody := bodyStr
	rePass := regexp.MustCompile(`(?i)(pass(?:word)?|prefilled)=[^&]+`)
	debugBody = rePass.ReplaceAllString(debugBody, "$1=********")
	fmt.Printf("[LOGIN-REQ] Sending POST Body: %s\n", debugBody)

	// 2. Thực hiện POST login
	req, err := http.NewRequest("POST", action, strings.NewReader(bodyStr))
	if err != nil {
		return nil, err
	}
	a.setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://mbasic.facebook.com/login.php")
	req.Header.Set("Origin", "https://mbasic.facebook.com")
	
	// Thêm các header bảo mật của Facebook
	if lsd, ok := a.extractValue(bodyStr, "lsd"); ok {
		req.Header.Set("X-FB-LSD", lsd)
	}

	resp, err = a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 3. Kiểm tra kết quả
	finalURL := resp.Request.URL.String()
	fmt.Printf("[LOGIN-REQ] Final URL: %s\n", finalURL)

	// Kiểm tra các mã lỗi phổ biến
	if strings.Contains(finalURL, "e=1348131") {
		return &RequestLoginResult{Success: false, Message: "Facebook báo sai mật khẩu hoặc từ chối do bảo mật."}, nil
	}
	if strings.Contains(finalURL, "login.php") && !strings.Contains(finalURL, "checkpoint") {
		// Thử bóc lỗi từ page
		return a.handleFailure(resp)
	}

	// Trường hợp: Yêu cầu 2FA
	if strings.Contains(finalURL, "/checkpoint/") || strings.Contains(finalURL, "next=https%3A%2F%2Fmbasic.facebook.com%2Fcheckpoint") {
		// Nếu có key 2FA, thử tự động nhập luôn
		if twoFactorKey != "" {
			return a.handle2FA(twoFactorKey)
		}
		return &RequestLoginResult{Success: false, Need2FA: true, Message: "Tài khoản yêu cầu mã 2FA."}, nil
	}

	// Trường hợp: Thành công
	if strings.Contains(finalURL, "facebook.com/home.php") || 
	   strings.Contains(finalURL, "facebook.com/mbasic/home") || 
	   strings.Contains(finalURL, "facebook.com/m.facebook.com") ||
	   a.hasCUser() {
		return a.extractPostLoginData()
	}

	// Trường hợp đặc biệt: Facebook hỏi "One-tap login"
	if strings.Contains(finalURL, "/login/device-based/update-nonce") {
		return a.extractPostLoginData()
	}

	// Trường hợp: Thất bại (Bóc tách lỗi từ HTML nếu có)
	errorMsg := "Đăng nhập thất bại. Vui lòng kiểm tra lại tài khoản/mật khẩu."
	respBody, _ := io.ReadAll(resp.Body)
	htmlContent := string(respBody)
	
	docErr, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	fbError := docErr.Find("#login_error, [role='alert']").Text()
	if fbError == "" {
		fbError = docErr.Find("div[style*='background: #fa3e3e']").Text()
	}

	if fbError != "" {
		errorMsg = "Facebook báo lỗi: " + fbError
	}

	return &RequestLoginResult{Success: false, Message: errorMsg}, nil
}

// handle2FA xử lý bước nhập mã OTP
func (a *RequestAutomator) handle2FA(key string) (*RequestLoginResult, error) {
	code, err := totp.GenerateCode(strings.ReplaceAll(key, " ", ""), time.Now())
	if err != nil {
		return &RequestLoginResult{Success: false, Message: "Key 2FA không hợp lệ: " + err.Error()}, nil
	}

	resp, err := a.client.Get("https://mbasic.facebook.com/checkpoint/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	doc.Find("form input[type='hidden']").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		val, _ := s.Attr("value")
		formData.Set(name, val)
	})
	formData.Set("approvals_code", code)
	formData.Set("submit[Submit Code]", "Submit Code")

	action, _ := doc.Find("form").Attr("action")
	if !strings.HasPrefix(action, "http") {
		action = "https://mbasic.facebook.com" + action
	}

	req, err := http.NewRequest("POST", action, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err = a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return a.handleSaveBrowser()
}

func (a *RequestAutomator) handleSaveBrowser() (*RequestLoginResult, error) {
	resp, err := a.client.Get("https://mbasic.facebook.com/checkpoint/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	doc.Find("form input[type='hidden']").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		val, _ := s.Attr("value")
		formData.Set(name, val)
	})
	
	if doc.Find("input[name='submit[Continue]']").Length() > 0 {
		formData.Set("submit[Continue]", "Continue")
	} else if doc.Find("button[name='submit[Continue]']").Length() > 0 {
		formData.Set("submit[Continue]", "Continue")
	}

	action, _ := doc.Find("form").Attr("action")
	if action != "" {
		if !strings.HasPrefix(action, "http") {
			action = "https://mbasic.facebook.com" + action
		}
		a.client.PostForm(action, formData)
	}

	return a.extractPostLoginData()
}

func (a *RequestAutomator) hasCUser() bool {
	u, _ := url.Parse("https://facebook.com")
	for _, c := range a.client.Jar.Cookies(u) {
		if c.Name == "c_user" {
			return true
		}
	}
	return false
}

func (a *RequestAutomator) extractPostLoginData() (*RequestLoginResult, error) {
	u, _ := url.Parse("https://mbasic.facebook.com")
	var cookieParts []string
	var uid string
	
	// Lấy tất cả cookies từ Jar cho domain facebook
	cookies := a.client.Jar.Cookies(u)
	fmt.Printf("[LOGIN-REQ] Đã lấy được %d cookies từ Jar\n", len(cookies))
	
	for _, c := range cookies {
		cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", c.Name, c.Value))
		if c.Name == "c_user" {
			uid = c.Value
		}
	}
	fullCookie := strings.Join(cookieParts, "; ")

	// Fetch www.facebook.com để lấy tiếp fb_dtsg (Dùng bản Desktop để dtsg chuẩn hơn)
	resp, err := a.client.Get("https://www.facebook.com/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	html := string(body)

	dtsg := ""
	re := regexp.MustCompile(`"DTSGInitData"\s*,\s*\[\s*\]\s*,\s*\{\s*"token"\s*:\s*"([^"]+)"`)
	m := re.FindStringSubmatch(html)
	if len(m) > 1 {
		dtsg = m[1]
	} else {
		re2 := regexp.MustCompile(`"dtsg":\{"token":"([^"]+)"\}`)
		m2 := re2.FindStringSubmatch(html)
		if len(m2) > 1 {
			dtsg = m2[1]
		}
	}

	name := "Account " + uid
	reName := regexp.MustCompile(`<title>(.*?)<\/title>`)
	mName := reName.FindStringSubmatch(html)
	if len(mName) > 1 {
		name = strings.TrimSuffix(mName[1], " | Facebook")
	}

	return &RequestLoginResult{
		Success: true,
		UID:     uid,
		Name:    name,
		Cookie:  fullCookie,
		Dtsg:    dtsg,
	}, nil
}

func (a *RequestAutomator) setCommonHeaders(req *http.Request) {
	// Giả lập Chrome 120 trên Windows (Vân tay trình duyệt hiện đại)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
}

func (a *RequestAutomator) extractValue(body, key string) (string, bool) {
	re := regexp.MustCompile(fmt.Sprintf(`%s=([^&]+)`, regexp.QuoteMeta(key)))
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		val, _ := url.QueryUnescape(match[1])
		return val, true
	}
	return "", false
}

func (a *RequestAutomator) handleFailure(resp *http.Response) (*RequestLoginResult, error) {
	// Re-get current body if needed, but since we already read it in Login we'd need to pass it
	// For simplicity, let's just use a generic message or re-read if possible (not possible with closing body)
	return &RequestLoginResult{Success: false, Message: "Đăng nhập thất bại. Facebook từ chối do bảo mật hoặc sai thông tin."}, nil
}
