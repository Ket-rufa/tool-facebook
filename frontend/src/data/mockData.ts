export const dashboardData = {
  lastUpdated: '09:41 AM, 24/05/2024',
  kpiCards: [
    { id: 'profiles', label: 'TỔNG HỒ SƠ', value: '1,240', subValue: '+12%', subType: 'success' },
    { id: 'posts', label: 'TỔNG BÀI VIẾT', value: '45,800', subValue: '+5.2k', subType: 'success' },
    { id: 'media', label: 'PHƯƠNG TIỆN', value: '12,300', subValue: 'Media Files', subType: 'neutral' },
    { id: 'scanning', label: 'ĐANG QUÉT', value: '8', subValue: 'Tasks Active', subType: 'neutral', iconRight: 'dot-blue' },
    { id: 'parser', label: 'TỶ LỆ PARSER', value: '98.5%', subValue: 'Ổn định', subType: 'success', iconRight: 'check' },
    { id: 'failed', label: 'TASK THẤT BẠI', value: '12', subValue: 'Cần kiểm tra', subType: 'danger', iconRight: 'alert' }
  ],
  systemHealth: {
    status: 'ALL SYSTEMS NOMINAL',
    modules: [
      { id: 'crawler', name: 'Crawler Module', status: 'Hoạt động (99%)', icon: 'cpu' },
      { id: 'proxy', name: 'Proxy Pool', status: '450/450 Live', icon: 'globe' },
      { id: 'parser', name: 'Data Parser', status: 'Phản hồi 45ms', icon: 'git-merge' }
    ]
  },
  recentActivities: [
    {
      id: 1,
      type: 'info',
      title: 'Bắt đầu quét ID: <span class="text-primary">fb_882910</span>',
      description: 'Tiến hành thu thập dữ liệu bài viết từ fanpage "Tech Community VN". Dự kiến hoàn thành sau 15 phút.',
      timeago: '2 phút trước'
    },
    {
      id: 2,
      type: 'success',
      title: 'Hoàn tất Parser: <span class="text-primary">task_77301</span>',
      description: 'Đã chuyển đổi 1,450 bình luận thành cấu trúc JSON thành công. Tỷ lệ lỗi 0%.',
      timeago: '12 phút trước'
    },
    {
      id: 3,
      type: 'danger',
      title: 'Cảnh báo Proxy: <span class="text-danger">node_asia_4</span>',
      description: 'Phát hiện checkpoint trên 3 tài khoản ảo. Hệ thống đã tự động chuyển sang pool dự phòng.',
      timeago: '45 phút trước'
    },
    {
      id: 4,
      type: 'info',
      title: 'Khởi tạo Backup',
      description: 'Dữ liệu hàng tuần đã được đồng bộ hóa với Amazon S3.',
      timeago: '1 giờ trước'
    }
  ],
  topProfiles: [
    { id: 1, name: 'VinFast Official', value: 14200, displayValue: '14.2k', progress: 100, color: 'bg-gradient-1' },
    { id: 2, name: 'Tech Review VN', value: 11800, displayValue: '11.8k', progress: 83, color: 'bg-gradient-2' },
    { id: 3, name: 'FPT Shop', value: 9400, displayValue: '9.4k', progress: 66, color: 'bg-gradient-3' },
    { id: 4, name: 'Shopee Vietnam', value: 8100, displayValue: '8.1k', progress: 57, color: 'bg-gradient-4' },
    { id: 5, name: 'Cộng đồng IT', value: 7600, displayValue: '7.6k', progress: 54, color: 'bg-gradient-5' },
    { id: 6, name: 'Tin tức 24h', value: 6200, displayValue: '6.2k', progress: 43, color: 'bg-gradient-6' }
  ]
}
