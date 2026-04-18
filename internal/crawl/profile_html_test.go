package crawl

import "testing"

func TestParseProfileInfoFromHTMLVietnameseCard(t *testing.T) {
	t.Parallel()

	body := `
	<html>
	<body>
		<div>Sống ở Thành phố Hồ Chí Minh</div>
		<div>Từ Đắk Lắk</div>
		<div>8 tháng 12, 1996</div>
		<div>Độc thân</div>
		<div>Nữ</div>
		<div>Giáo dục</div>
		<div>Trường Đại Học Tây Nguyên</div>
		<div>Công việc</div>
		<div>OpenAI</div>
		<div>Tiểu sử</div>
		<div>xin chao</div>
		<div>1.2K người theo dõi</div>
	</body>
	</html>`

	h := &CrawlHandler{}
	info := &ProfileInfo{UID: "61572157625165"}
	h.parseProfileInfoFromHTML(body, info)

	if info.CurrentCity != "Thành phố Hồ Chí Minh" {
		t.Fatalf("expected current city, got %q", info.CurrentCity)
	}
	if info.Hometown != "Đắk Lắk" {
		t.Fatalf("expected hometown, got %q", info.Hometown)
	}
	if info.Birthday != "8 tháng 12, 1996" {
		t.Fatalf("expected birthday, got %q", info.Birthday)
	}
	if info.Relationship != "Độc thân" {
		t.Fatalf("expected relationship, got %q", info.Relationship)
	}
	if info.Gender != "FEMALE" {
		t.Fatalf("expected gender FEMALE, got %q", info.Gender)
	}
	if len(info.Education) != 1 || info.Education[0] != "Trường Đại Học Tây Nguyên" {
		t.Fatalf("expected education, got %#v", info.Education)
	}
	if len(info.Work) != 1 || info.Work[0] != "OpenAI" {
		t.Fatalf("expected work, got %#v", info.Work)
	}
	if info.Bio != "xin chao" {
		t.Fatalf("expected bio, got %q", info.Bio)
	}
	if info.Followers != "1.2K" {
		t.Fatalf("expected followers, got %q", info.Followers)
	}
}
