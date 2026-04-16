package crawl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ExtractPostsFromHTML parses Facebook HTML (both Comet JSON and mbasic HTML)
func (h *CrawlHandler) ExtractPostsFromHTML(html string) []CrawlPostEntity {
	var results []CrawlPostEntity
	seenIDs := make(map[string]struct{})

	// --- CHIẾN THUẬT 1: DÒ TÌM TRONG JSON (COMET/MOBILE APP) ---
	rePost := regexp.MustCompile(`(?s)[\\"]*post_id[\\"]*[:\s]+[\\"]*(\d+)[\\"]*.*?[\\"]*text[\\"]*[:\s]+[\\"]*(.*?)[\\"]*[,\}].*?[\\"]*creation_time[\\"]*[:\s]+(\d+)`)
	matches := rePost.FindAllStringSubmatch(html, -1)

	for _, m := range matches {
		id := m[1]
		if id == "" || strings.HasPrefix(id, "0") { continue }
		if _, ok := seenIDs[id]; ok { continue }
		seenIDs[id] = struct{}{}

		text := unescapeFBCore(m[2])
		ctime := m[3]
		var ts int64
		_, _ = fmt.Sscanf(ctime, "%d", &ts)

		results = append(results, CrawlPostEntity{
			PostID:    id,
			Content:   text,
			Time:      time.Unix(ts, 0).Format("02/01/2006 15:04"),
			Reactions: "0",
			Comments:  "0",
		})
	}

	// --- CHIẾN THUẬT 2: DÒ TÌM TRONG HTML THUẦN (MBASIC/TOUCH) ---
	if len(results) == 0 {
		// Tìm các khối bài viết (thường nằm trong thẻ <div> hoặc <article>)
		// mbasic: bài viết thường có link chứa /story.php?story_fbid=... hoặc /permalink.php?story_fbid=...
		reMBasic := regexp.MustCompile(`(?s)<div[^>]*>(.*?)</div>`)
		divs := reMBasic.FindAllString(html, -1)
		
		for _, div := range divs {
			// Tìm link Story ID
			idMatch := regexp.MustCompile(`(?:story_fbid|fbid)=(\d+)`).FindStringSubmatch(div)
			if len(idMatch) < 2 { continue }
			id := idMatch[1]
			
			if _, ok := seenIDs[id]; ok { continue }
			seenIDs[id] = struct{}{}

			// Tìm nội dung (thường nằm trong thẻ p hoặc div)
			content := ""
			// Loại bỏ tags để lấy text
			content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(div, "")
			content = strings.TrimSpace(content)
			
			if len(content) < 5 { continue }

			results = append(results, CrawlPostEntity{
				PostID:    id,
				Content:   content,
				Time:      "Vừa mới đây", // mbasic time parsing hơi khó, lấy tạm
				Reactions: "0",
				Comments:  "0",
			})
		}
	}

	// --- CHIẾN THUẬT 3: TÌM TƯƠNG TÁC ---
	reInteractions := []struct{
		Regex *regexp.Regexp
		Type  string
	}{
		{regexp.MustCompile(`(?s)[\\"]*post_id[\\"]*[:\s]+[\\"]*(\d+)[\\"]*.*?[\\"]*(?:reaction_count|reactors)[\\"]*[:\s]+\{.*?[\\"]*count[\\"]*[:\s]+(\d+)`), "react"},
		{regexp.MustCompile(`(?s)[\\"]*post_id[\\"]*[:\s]+[\\"]*(\d+)[\\"]*.*?[\\"]*comment_count[\\"]*[:\s]+\{.*?[\\"]*(?:total_count|count)[\\"]*[:\s]+(\d+)`), "comm"},
		// Support mbasic pattern: Like 123, Comment 456
		{regexp.MustCompile(`(\d+)\s+(?:Like|Thích|người thích)`), "react_mbasic"},
	}

	for _, ri := range reInteractions {
		mats := ri.Regex.FindAllStringSubmatch(html, -1)
		if ri.Type == "react_mbasic" && len(results) > 0 {
			// Gán tạm số Like cho bài gần nhất nếu là mbasic
			if len(mats) > 0 { results[0].Reactions = mats[0][1] }
			continue
		}
		for _, m := range mats {
			id := m[1]; val := m[2]
			for i := range results {
				if results[i].PostID == id {
					if ri.Type == "react" { results[i].Reactions = val }
					if ri.Type == "comm" { results[i].Comments = val }
					break
				}
			}
		}
	}

	return results
}

func unescapeFBCore(s string) string {
	var dummy struct{ S string }
	_ = json.Unmarshal([]byte(`{"S":"`+s+`"}`), &dummy)
	if dummy.S != "" { return dummy.S }
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "\\n", "\n")
	return s
}
