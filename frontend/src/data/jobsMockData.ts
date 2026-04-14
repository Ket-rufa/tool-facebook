export const jobsData = {
  list: [
    {
      id: '#CR-8921',
      url: 'facebook.com/vinfast.official',
      subText: 'VinFast Vietnam',
      type: 'Page',
      status: 'running',
      time: '14:23:45\n22/10/2023'
    },
    {
      id: '#CR-8920',
      url: 'facebook.com/groups/techr...',
      subText: 'Tech Review Community',
      type: 'Group',
      status: 'success',
      time: '12:15:32\n22/10/2023'
    },
    {
      id: '#CR-8919',
      url: 'facebook.com/ronaldo',
      subText: 'Cristiano Ronaldo',
      type: 'Profile',
      status: 'error',
      time: '11:05:10\n22/10/2023'
    },
    {
      id: '#CR-8918',
      url: 'facebook.com/nike',
      subText: 'Nike Global',
      type: 'Page',
      status: 'success',
      time: '09:30:00\n22/10/2023'
    }
  ],
  detail: {
    id: '#CR-8920',
    progress: {
      statusText: 'Đã hoàn thành lúc 12:15:32',
      steps: [
        { name: 'Kiểm tra URL', time: '0.2s', status: 'success' },
        { name: 'Trình phân tích (DOM)', time: '2.1s', status: 'success' },
        { name: 'Tải hình ảnh/media', time: '5.4s', status: 'success' },
        { name: 'Lưu cơ sở dữ liệu', time: '1.5s', status: 'success' }
      ]
    },
    warning: 'Đã bỏ qua 4 bài viết do lỗi mã hóa ký tự (UTF-8 mismatch). Kiểm tra lại cấu hình parser.',
    metadata: `{\n  "job_id": "CR-8920",\n  "engine_version": "v4.2-stable",\n  "user_agent": "Mozilla/5.0 (Windows NT 10.0...)",\n  "proxy_pool": "ASIA_VIP_04",\n  "retries": 0,\n  "payload": {\n    "target": "group_activity",\n    "depth": 50,\n    "scrape_comments": true\n  }\n}`,
    executor: 'Hệ thống (Cron)',
    priority: 'High (P0)'
  }
}
