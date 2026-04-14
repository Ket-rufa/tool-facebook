export const integrationsData = {
  header: {
    title: 'Tích hợp & Mở rộng',
    description: 'Mở rộng khả năng của hệ thống một cách an toàn thông qua các module nội bộ và tích hợp có kiểm soát.'
  },
  mainCards: [
    {
      id: 'local_session',
      icon: 'sliders',
      title: 'Quản lý phiên cục bộ',
      badge: { text: 'ĐANG PHÁT TRIỂN', color: 'warning' },
      description: 'Quản lý trạng thái phiên làm việc cục bộ của ứng dụng, khôi phục phiên trước đó, kiểm soát kết nối app-engine.',
      cta: 'Thiết lập module'
    },
    {
      id: 'account_tasks',
      icon: 'repeat',
      title: 'Tác vụ tài khoản (Sắp ra mắt)',
      badge: { text: 'SẮP RA MẮT', color: 'info' },
      description: 'Module tương lai tích hợp quyền kiểm soát, xem xét truy cập và các hành động thủ công có xác nhận.',
      cta: 'Đăng ký nhận thông báo'
    }
  ],
  smallCards: [
    {
      id: 'export_report',
      icon: 'fileText',
      title: 'Xuất dữ liệu & Báo cáo',
      badge: { text: 'CẦN THIẾT LẬP', color: 'warning' },
      description: 'Xuất dữ liệu sang các định dạng CSV, JSON, Excel hoặc tạo báo cáo PDF chuyên sâu.',
      cta: 'BẮT ĐẦU'
    },
    {
      id: 'webhook_api',
      icon: 'box',
      title: 'Webhook & API nội bộ',
      badge: { text: 'CHƯA KÍCH HOẠT', color: 'default' },
      description: 'Kết nối và đẩy dữ liệu sang các hệ thống nội bộ khác thông qua API an toàn.',
      cta: 'XEM TÀI LIỆU'
    },
    {
      id: 'analytics',
      icon: 'barChart',
      title: 'Phân tích dữ liệu',
      badge: { text: 'SẮP RA MẮT', color: 'info' },
      description: 'Sử dụng các công cụ phân tích để đánh giá xu hướng và chất lượng dữ liệu thu thập.',
      cta: 'ĐĂNG KÝ SỚM'
    }
  ],
  banner: {
    title: 'Kiến trúc mở rộng an toàn',
    description: 'Hệ thống được thiết kế theo kiến trúc module, cho phép kích hoạt có kiểm soát, phân quyền chi tiết và lưu nhật ký kiểm soát toàn diện.',
    cta: 'Khám phá tài liệu kỹ thuật'
  },
  toast: {
    title: 'Yêu cầu quyền truy cập',
    description: 'Vui lòng liên hệ quản trị viên để đăng ký tham gia thử nghiệm các module Beta.'
  }
}
