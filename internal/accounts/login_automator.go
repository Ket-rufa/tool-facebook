package accounts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/pquerna/otp/totp"
)

// LoginAutomator xử lý việc điều khiển trình duyệt để lấy session
type LoginAutomator struct {
	// Có thể thêm cấu hình timeout, v.v.
}

func NewLoginAutomator() *LoginAutomator {
	return &LoginAutomator{}
}

// Result là kết quả sau khi đăng nhập thành công
type AutomatorResult struct {
	UID    string
	Name   string
	Cookie string
	Dtsg   string
}

// StartLogin opens a real Chrome window and waits for the user to log in.
// Once logged in (c_user cookie found), it extracts dtsg and other info.
func (a *LoginAutomator) StartLogin(ctx context.Context) (*AutomatorResult, error) {
	// Tắt headless để người dùng thấy trình duyệt
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Tạo context con cho chromedp
	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Timeout tổng cộng là 5 phút cho việc đăng nhập
	timeoutCtx, cancel := context.WithTimeout(browserCtx, 5*time.Minute)
	defer cancel()

	result := &AutomatorResult{}
	var cookies []*network.Cookie
	var dtsg string
	var userName string

	fmt.Println("[LOGIN-AUTO] Đang mở trình duyệt...")

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate("https://www.facebook.com/login"),
		// Vòng lặp chờ cho đến khi thấy cookie c_user
		chromedp.ActionFunc(func(ctx context.Context) error {
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					// Kiểm tra cookie
					currentCookies, err := network.GetCookies().Do(ctx)
					if err != nil {
						continue
					}

					var uid string
					for _, c := range currentCookies {
						if c.Name == "c_user" {
							uid = c.Value
							break
						}
					}

					// Nếu có c_user, nghĩa là đã đăng nhập thành công
					if uid != "" {
						result.UID = uid
						cookies = currentCookies
						fmt.Printf("[LOGIN-AUTO] Đăng nhập thành công! UID: %s\n", uid)
						return nil
					}
				}
			}
		}),
		// Đợi thêm một chút để trang load ổn định
		chromedp.Sleep(2*time.Second),
		// Lấy fb_dtsg qua Javascript
		chromedp.Evaluate(`(function() {
			try {
				// Mẫu 1: DTSGInitData
				if (window.DTSGInitData && window.DTSGInitData.token) return window.DTSGInitData.token;
				// Mẫu 2: DTSGInitialData
				if (window.DTSGInitialData && window.DTSGInitialData.token) return window.DTSGInitialData.token;
				// Mẫu 3: Bốc từ mã nguồn
				const m = document.documentElement.innerHTML.match(/"DTSGInitData"\s*,\s*\[\s*\]\s*,\s*\{\s*"token"\s*:\s*"([^"]+)"/);
				if (m) return m[1];
				return "";
			} catch(e) { return ""; }
		})()`, &dtsg),
		// Lấy tên người dùng
		chromedp.Evaluate(`document.title.replace(/\s*\|\s*Facebook/i, "").trim()`, &userName),
	)

	if err != nil {
		return nil, fmt.Errorf("lỗi trong quá trình automation: %w", err)
	}

	result.Dtsg = dtsg
	result.Name = userName
	if result.Name == "" || result.Name == "Facebook" || result.Name == "Log in" {
		result.Name = "New Account (" + result.UID + ")"
	}

	// Gộp cookie thành chuỗi
	var cookieParts []string
	for _, c := range cookies {
		cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}
	result.Cookie = strings.Join(cookieParts, "; ")

	return result, nil
}

// AutomatedLogin performs a fully automated login using credentials.
// It opens a visible browser window so the user can watch the process.
func (a *LoginAutomator) AutomatedLogin(ctx context.Context, email, pass, twoFactorKey string) (*AutomatorResult, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Hiện cửa sổ trình duyệt
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Timeout 2 phút cho quá trình tự động
	timeoutCtx, cancel := context.WithTimeout(browserCtx, 2*time.Minute)
	defer cancel()

	result := &AutomatorResult{}
	var cookies []*network.Cookie
	var dtsg string
	var userName string

	fmt.Printf("[LOGIN-AUTO] Bắt đầu đăng nhập tự động cho: %s\n", email)

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate("https://www.facebook.com"),
		chromedp.Sleep(2*time.Second),
		
		// 1. Xử lý bảng chấp nhận Cookie (nếu có)
		chromedp.ActionFunc(func(ctx context.Context) error {
			var nodes []*cdp.Node
			// Thử tìm nút Chấp nhận/Cho phép cookie
			// Các selector phổ biến: [data-testid="cookie-policy-manage-dialog-accept-button"], button[title*="Allow"], button[title*="Chấp nhận"]
			err := chromedp.Nodes(`button[data-cookiebanner="accept_button"], button[data-testid*="accept"], button[title*="Allow"], button[title*="Chấp nhận"]`, &nodes, chromedp.AtLeast(0)).Do(ctx)
			if err == nil && len(nodes) > 0 {
				fmt.Println("[LOGIN-AUTO] Đang nhấn chấp nhận Cookie...")
				return chromedp.Click(`button[data-cookiebanner="accept_button"], button[data-testid*="accept"], button[title*="Allow"], button[title*="Chấp nhận"]`, chromedp.ByQuery).Do(ctx)
			}
			return nil
		}),
		chromedp.Sleep(1*time.Second),

		// 2. Điền thông tin với selector linh hoạt và gõ phím như người thật
		chromedp.WaitVisible(`input[name="email"], #email`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("[LOGIN-AUTO] Đang nhập Email...")
			// Gõ từng phím với độ trễ ngẫu nhiên để giống người thật
			return chromedp.SendKeys(`input[name="email"], #email`, email, chromedp.ByQuery).Do(ctx)
		}),
		chromedp.Sleep(1*time.Second),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("[LOGIN-AUTO] Đang nhập Mật khẩu...")
			// Nhấn Enter sau khi gõ pass để gửi form
			return chromedp.SendKeys(`input[name="pass"], #pass`, pass + "\n", chromedp.ByQuery).Do(ctx)
		}),
		chromedp.Sleep(2*time.Second),
		
		// 3. Xử lý trang 2FA / Checkpoint
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Chờ tối đa 5 giây để xem có trang 2FA hiện ra không
			subCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			var nodes []*cdp.Node
			err := chromedp.Nodes("input[name='approvals_code'], #approvals_code", &nodes, chromedp.AtLeast(0)).Do(subCtx)
			if err == nil && len(nodes) > 0 {
				fmt.Println("[LOGIN-AUTO] Phát hiện yêu cầu mã 2FA...")
				if twoFactorKey != "" {
					fmt.Println("[LOGIN-AUTO] Đang tự động tạo và điền mã OTP...")
					code, _ := totp.GenerateCode(strings.ReplaceAll(twoFactorKey, " ", ""), time.Now())
					return chromedp.Run(ctx,
						chromedp.SendKeys("input[name='approvals_code'], #approvals_code", code, chromedp.ByQuery),
						chromedp.Sleep(1*time.Second),
						// Click nút gửi mã (thường là checkpointSubmitButton)
						chromedp.Click("#checkpointSubmitButton, button[type='submit']", chromedp.ByQuery),
						chromedp.Sleep(3*time.Second),
						// Xử lý trang "Lưu trình duyệt" (Save Browser)
						chromedp.ActionFunc(func(ctx context.Context) error {
							var saveNodes []*cdp.Node
							chromedp.Nodes("#checkpointSubmitButton, button[value='continue']", &saveNodes, chromedp.AtLeast(0)).Do(ctx)
							if len(saveNodes) > 0 {
								fmt.Println("[LOGIN-AUTO] Đang nhấn 'Tiếp tục' qua trang Lưu trình duyệt...")
								return chromedp.Click("#checkpointSubmitButton, button[value='continue']", chromedp.ByQuery).Do(ctx)
							}
							return nil
						}),
						chromedp.Sleep(2*time.Second),
					)
				} else {
					fmt.Println("[WARN] Tài khoản yêu cầu 2FA nhưng bạn chưa nhập khóa bí mật (Secret Key). Vui lòng tự nhập mã OTP trên trình duyệt.")
				}
			}
			return nil
		}),

		// 4. Vòng lặp chờ c_user (xác nhận đăng nhập thành công)
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("[LOGIN-AUTO] Chờ xác nhận phiên làm việc từ Facebook...")
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return fmt.Errorf("hết thời gian chờ đăng nhập thành công. Có thể bạn cần xác minh thêm trên trình duyệt.")
				case <-ticker.C:
					currentCookies, _ := network.GetCookies().Do(ctx)
					for _, c := range currentCookies {
						if c.Name == "c_user" {
							result.UID = c.Value
							cookies = currentCookies
							fmt.Println("[LOGIN-AUTO] Đã tìm thấy phiên (c_user). Đăng nhập thành công!")
							return nil
						}
					}
				}
			}
		}),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`(function() {
			try {
				// Cách 1: Tìm trong các biến toàn cục (Chuẩn nhất)
				if (window.DTSGInitData && window.DTSGInitData.token) return window.DTSGInitData.token;
				if (window.DTSGInitialData && window.DTSGInitialData.token) return window.DTSGInitialData.token;
				
				// Cách 2: Tìm trong form (Dành cho bản mobile/mbasic)
				const input = document.querySelector('input[name="fb_dtsg"]');
				if (input) return input.value;
				
				// Cách 3: Quét chuỗi trong toàn bộ HTML (Dự phòng cuối cùng)
				const html = document.documentElement.innerHTML;
				const patterns = [
					/\"DTSGInitData\"\,\[\]\,\{\"token\"\:\"([^\"]+)\"/,
					/\"dtsg\"\:\{\"token\"\:\"([^\"]+)\"/,
					/name=\"fb_dtsg\" value=\"([^\"]+)\"/
				];
				for (const p of patterns) {
					const m = html.match(p);
					if (m && m[1]) return m[1];
				}
				return "";
			} catch(e) { return ""; }
		})()`, &dtsg),
		chromedp.Evaluate(`document.title.replace(/\s*\|\s*Facebook/i, "").trim()`, &userName),
	)

	if err != nil {
		return nil, err
	}

	result.Dtsg = dtsg
	result.Name = userName
	if result.Name == "" || result.Name == "Facebook" {
		result.Name = "Account " + result.UID
	}

	var cookieParts []string
	for _, c := range cookies {
		cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}
	result.Cookie = strings.Join(cookieParts, "; ")

	return result, nil
}
