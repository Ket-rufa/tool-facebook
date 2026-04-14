export const settingsData = {
  header: {
    title: 'Cài đặt hệ thống',
    description: 'Quản trị và tối ưu hóa các thông số kỹ thuật cho hệ thống Facebook Manager.'
  },
  sectionsNav: [
    { id: 'crawl', label: 'Cấu hình Crawl', icon: 'sliders' },
    { id: 'storage', label: 'Lưu trữ & Dữ liệu', icon: 'server' },
    { id: 'system', label: 'Hệ thống & UI', icon: 'layout' },
    { id: 'about', label: 'Giới thiệu', icon: 'info' }
  ],
  crawl: {
    retryPolicy: 'Thử lại 3 lần (Mặc định)',
    useExponentialBackoff: true,
    concurrentThreads: 12
  },
  storage: {
    path: '/var/www/ffb_manager/data/raw_crawls',
    retention: '30 Ngày',
    retentionOptions: ['30 Ngày', '90 Ngày', '1 Năm', 'Vĩnh viễn']
  },
  logLevel: 'info',
  uiOptions: {
    density: 'Cao',
    fontSize: 'Mặc định (14px)',
    pureWhite: true
  },
  version: {
    current: '2.4.0-stable',
    lastUpdate: '12/10/2023'
  }
}
