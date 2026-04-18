package crawl

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/yourname/tool-facebook/internal/actiontest"
)

var (
	profileHTMLScriptRe     = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script>`)
	profileHTMLStyleRe      = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style>`)
	profileHTMLCommentRe    = regexp.MustCompile(`(?is)<!--.*?-->`)
	profileHTMLTagRe        = regexp.MustCompile(`(?is)<[^>]+>`)
	profileHTMLSpaceRe      = regexp.MustCompile(`\s+`)
	profileBirthdayViRe     = regexp.MustCompile(`^\d{1,2}\s+tháng\s+\d{1,2},\s+\d{4}$`)
	profileBirthdayEnRe     = regexp.MustCompile(`^[A-Za-z]+\s+\d{1,2},\s+\d{4}$`)
	profileFollowersLineRe  = regexp.MustCompile(`(?i)([\d\.,kKmM]+)\s+(?:người theo dõi|followers?)`)
	profileFollowersValueRe = regexp.MustCompile(`[\d\.,]+`)
)

func (h *CrawlHandler) needsHTMLProfileFallback(info *ProfileInfo) bool {
	return info.Birthday == "" ||
		info.CurrentCity == "" ||
		info.Hometown == "" ||
		info.Relationship == "" ||
		len(info.Education) == 0 ||
		len(info.Work) == 0 ||
		info.Followers == "" ||
		info.Bio == ""
}

func (h *CrawlHandler) hydrateProfileFromHTML(uid, cookie string, info *ProfileInfo) {
	attempts := []struct {
		label string
		url   string
		ua    string
	}{
		{
			label: "desktop_timeline",
			url:   fmt.Sprintf("https://www.facebook.com/profile.php?id=%s", uid),
			ua:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		},
		{
			label: "desktop_about",
			url:   fmt.Sprintf("https://www.facebook.com/profile.php?id=%s&sk=about", uid),
			ua:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		},
		{
			label: "mbasic_info",
			url:   fmt.Sprintf("https://mbasic.facebook.com/profile.php?id=%s&v=info", uid),
			ua:    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Mobile/15E148 Safari/604.1",
		},
	}

	for _, attempt := range attempts {
		status, body, finalURL, err := h.fetchProfileHTMLPage(attempt.url, cookie, attempt.ua)
		if err != nil {
			continue
		}
		if status != http.StatusOK || h.isUnsupportedProfileHTML(body) {
			continue
		}

		h.saveProfileHTMLDebug(attempt.label, body)
		h.parseProfileInfoFromHTML(body, info)

		if info.Name == "" && finalURL != "" {
			if derived := h.deriveProfileFallbackName(uid, profileAboutTokens{AboutURL: finalURL, InfoAllURL: finalURL}); derived != "" && derived != uid {
				info.Name = derived
			}
		}

		if !h.needsHTMLProfileFallback(info) {
			return
		}
	}
}

func (h *CrawlHandler) fetchProfileHTMLPage(targetURL, cookie, ua string) (int, string, string, error) {
	req, _ := http.NewRequest("GET", targetURL, nil)
	req.Header.Set("Cookie", actiontest.SanitizeCookie(cookie))
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "vi-VN,vi;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", "https://www.facebook.com/")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", "", err
	}

	finalURL := targetURL
	if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}
	return resp.StatusCode, string(body), finalURL, nil
}

func (h *CrawlHandler) isUnsupportedProfileHTML(body string) bool {
	return strings.Contains(body, "Trình duyệt này không hỗ trợ Facebook") ||
		strings.Contains(body, "This browser is not supported") ||
		strings.Contains(body, "unsupported-interstitial")
}

func (h *CrawlHandler) saveProfileHTMLDebug(label, body string) {
	limit := len(body)
	if limit > 250000 {
		limit = 250000
	}
	_ = os.WriteFile(fmt.Sprintf("data/last_profile_html_%s_debug.html", label), []byte(body[:limit]), 0644)
}

func (h *CrawlHandler) parseProfileInfoFromHTML(body string, info *ProfileInfo) {
	lines := extractProfileHTMLLines(body)
	if len(lines) == 0 {
		return
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		norm := normalizeLabel(line)

		// Xử lý bio placeholder
		if norm == "bạn đang nghĩ gì?" || norm == "what's on your mind?" {
			if info.Bio == "" {
				info.Bio = line
				fmt.Printf("[SCAN][HTML] >> Đã chuyển bio placeholder sang Tiểu sử: %s\n", line)
			}
			continue
		}

		switch {
		case info.CurrentCity == "" && hasAnyPrefix(norm, "sống ở ", "sống tại ", "lives in "):
			value := trimKnownPrefixes(line, "Sống ở ", "Sống tại ", "Lives in ")
			if value != "" {
				info.CurrentCity = value
				fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Thành phố hiện tại: %s\n", value)
			}
		case info.Hometown == "" && hasAnyPrefix(norm, "từ ", "đến từ ", "from "):
			value := trimKnownPrefixes(line, "Từ ", "Đến từ ", "From ")
			if value != "" {
				info.Hometown = value
				fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Quê quán: %s\n", value)
			}
		case info.Birthday == "" && isProfileBirthdayLine(line):
			info.Birthday = line
			fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Ngày sinh: %s\n", line)
		case info.Relationship == "" && isProfileRelationshipLine(line):
			info.Relationship = line
			fmt.Printf("[SCAN][HTML] >> Đã tìm thấy MQH: %s\n", line)
		case info.Gender == "" && isProfileGenderLine(line):
			info.Gender = normalizeProfileGender(line)
			fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Giới tính: %s\n", info.Gender)
		case info.Followers == "" && strings.Contains(norm, "người theo dõi"):
			if value := parseFollowersText(line); value != "" {
				info.Followers = value
				fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Followers: %s\n", value)
			}
		case info.Followers == "" && strings.Contains(norm, "followers"):
			if value := parseFollowersText(line); value != "" {
				info.Followers = value
				fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Followers: %s\n", value)
			}
		case hasAnyPrefix(norm, "học tại ", "từng học tại ", "studied at "):
			value := trimKnownPrefixes(line, "Học tại ", "Từng học tại ", "Studied at ")
			if value != "" {
				info.Education = appendUnique(info.Education, value)
				fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Học vấn/Giáo dục: %s\n", value)
			}
		case hasAnyPrefix(norm, "làm việc tại ", "works at "):
			value := trimKnownPrefixes(line, "Làm việc tại ", "Works at ")
			if value != "" {
				info.Work = appendUnique(info.Work, value)
				fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Công việc: %s\n", value)
			}
		case norm == "tiểu sử" || norm == "bio":
			if info.Bio == "" {
				if value := h.nextProfileHTMLValue(lines, i+1); value != "" {
					info.Bio = value
					fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Tiểu sử: %s\n", value)
				}
			}
		case norm == "giáo dục" || norm == "education" || norm == "học vấn":
			for _, value := range h.collectProfileHTMLSectionValues(lines, i+1) {
				if !isProfileHeadingLine(value) {
					info.Education = appendUnique(info.Education, value)
					fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Học vấn/Giáo dục: %s\n", value)
				}
			}
		case norm == "công việc" || norm == "work":
			for _, value := range h.collectProfileHTMLSectionValues(lines, i+1) {
				if !isProfileHeadingLine(value) {
					info.Work = appendUnique(info.Work, value)
					fmt.Printf("[SCAN][HTML] >> Đã tìm thấy Công việc: %s\n", value)
				}
			}
		}
	}
}

func extractProfileHTMLLines(body string) []string {
	body = profileHTMLScriptRe.ReplaceAllString(body, "\n")
	body = profileHTMLStyleRe.ReplaceAllString(body, "\n")
	body = profileHTMLCommentRe.ReplaceAllString(body, "\n")
	body = strings.ReplaceAll(body, "<br>", "\n")
	body = strings.ReplaceAll(body, "<br/>", "\n")
	body = strings.ReplaceAll(body, "<br />", "\n")
	body = profileHTMLTagRe.ReplaceAllString(body, "\n")
	body = html.UnescapeString(body)
	body = strings.ReplaceAll(body, "\u00a0", " ")

	rawLines := strings.Split(body, "\n")
	lines := make([]string, 0, len(rawLines))
	prev := ""
	for _, raw := range rawLines {
		line := strings.TrimSpace(profileHTMLSpaceRe.ReplaceAllString(raw, " "))
		if !isUsefulProfileHTMLLine(line) {
			continue
		}
		// Bỏ qua các chuỗi placeholder chung
		norm := normalizeLabel(line)
		if norm == "bạn đang nghĩ gì?" || norm == "what's on your mind?" {
			continue
		}
		if line == prev {
			continue
		}
		lines = append(lines, line)
		prev = line
	}
	return lines
}

func isUsefulProfileHTMLLine(line string) bool {
	if line == "" || len(line) < 2 {
		return false
	}
	if strings.Contains(line, "http://") || strings.Contains(line, "https://") {
		return false
	}
	if strings.Contains(line, "function(") || strings.Contains(line, "return ") {
		return false
	}

	switch normalizeLabel(line) {
	case "facebook", "trang chủ", "tin nhắn", "thông báo", "chat", "tìm bạn bè", "menu", "tìm kiếm",
		"quay lại đầu trang", "chỉnh sửa trang cá nhân", "edit profile", "help", "trợ giúp":
		return false
	}
	return true
}

func isProfileBirthdayLine(line string) bool {
	line = strings.TrimSpace(line)
	return profileBirthdayViRe.MatchString(line) || profileBirthdayEnRe.MatchString(line)
}

func isProfileRelationshipLine(line string) bool {
	switch normalizeLabel(line) {
	case "độc thân", "single", "đang hẹn hò", "in a relationship", "đã kết hôn", "married",
		"đã đính hôn", "engaged", "ly hôn", "divorced", "góa", "widowed", "phức tạp", "it's complicated":
		return true
	default:
		return false
	}
}

func isProfileGenderLine(line string) bool {
	switch normalizeLabel(line) {
	case "nữ", "female", "nam", "male":
		return true
	default:
		return false
	}
}

func normalizeProfileGender(line string) string {
	switch normalizeLabel(line) {
	case "nữ", "female":
		return "FEMALE"
	case "nam", "male":
		return "MALE"
	default:
		return strings.TrimSpace(line)
	}
}

func parseFollowersText(line string) string {
	if m := profileFollowersLineRe.FindStringSubmatch(line); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	if m := profileFollowersValueRe.FindString(line); m != "" {
		return strings.TrimSpace(m)
	}
	return ""
}

func (h *CrawlHandler) nextProfileHTMLValue(lines []string, start int) string {
	for i := start; i < len(lines) && i < start+4; i++ {
		if isProfileHeadingLine(lines[i]) {
			break
		}
		if !isProfileMetadataLine(lines[i]) {
			return lines[i]
		}
	}
	return ""
}

func (h *CrawlHandler) collectProfileHTMLSectionValues(lines []string, start int) []string {
	values := make([]string, 0, 2)
	for i := start; i < len(lines) && i < start+6; i++ {
		line := lines[i]
		if isProfileHeadingLine(line) {
			break
		}
		if isProfileMetadataLine(line) {
			continue
		}
		// Bỏ qua bio placeholder trong lúc thu thập value của section
		norm := normalizeLabel(line)
		if norm == "bạn đang nghĩ gì?" || norm == "what's on your mind?" {
			continue
		}
		values = append(values, line)
		if len(values) >= 2 {
			break
		}
	}
	return values
}

func isProfileHeadingLine(line string) bool {
	switch normalizeLabel(line) {
	case "thông tin cá nhân", "personal information", "giáo dục", "education", "học vấn",
		"công việc", "work", "tiểu sử", "bio", "bài viết", "posts", "giới thiệu", "about":
		return true
	default:
		return false
	}
}

func isProfileMetadataLine(line string) bool {
	norm := normalizeLabel(line)
	if norm == "" {
		return true
	}
	if strings.HasPrefix(norm, "xem ") || strings.HasPrefix(norm, "see ") || strings.HasPrefix(norm, "chỉnh sửa") {
		return true
	}
	if isProfileBirthdayLine(line) || isProfileRelationshipLine(line) || isProfileGenderLine(line) {
		return true
	}
	if hasAnyPrefix(norm, "sống ở ", "sống tại ", "lives in ", "từ ", "đến từ ", "from ", "học tại ", "từng học tại ", "studied at ", "làm việc tại ", "works at ") {
		return true
	}
	return strings.Contains(norm, "người theo dõi") || strings.Contains(norm, "followers")
}

func hasAnyPrefix(line string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(line, normalizeLabel(prefix)) {
			return true
		}
	}
	return false
}

func trimKnownPrefixes(line string, prefixes ...string) string {
	trimmed := strings.TrimSpace(line)
	for _, prefix := range prefixes {
		if strings.HasPrefix(normalizeLabel(trimmed), normalizeLabel(prefix)) {
			return strings.TrimSpace(trimmed[len(prefix):])
		}
	}
	return trimmed
}
