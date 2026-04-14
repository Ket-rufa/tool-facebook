export const accountsData = {
  header: {
    title: 'Quản lý tài khoản',
    description: 'Quản lý các tài khoản đã gắn phiên đăng nhập cục bộ để dùng cho thao tác thủ công và thử nghiệm.'
  },
  list: [
    {
      id: 'acc_101',
      accountId: '100084323145',
      displayName: 'Nguyễn Văn A',
      avatar: 'https://ui-avatars.com/api/?name=NVA&background=eff6ff&color=0f62fe',
      accountType: 'Profile',
      provider: 'Facebook',
      sessionStatus: 'active',
      lastCheckedAt: '10 phút trước',
      isDefault: true,
      note: 'Dùng cho các chiến dịch bình luận mồi (Seeding A).'
    },
    {
      id: 'acc_102',
      accountId: '45910283011',
      displayName: 'Tech Startup VN',
      avatar: 'https://ui-avatars.com/api/?name=TS&background=fce7f3&color=be185d',
      accountType: 'Page',
      provider: 'Facebook',
      sessionStatus: 'expired',
      lastCheckedAt: '2 ngày trước',
      isDefault: false,
      note: 'Page phụ dùng để chia sẻ link bài viết.'
    },
    {
      id: 'acc_103',
      accountId: '100092837482',
      displayName: 'Test Account 01',
      avatar: 'https://ui-avatars.com/api/?name=T01&background=fef3c7&color=b45309',
      accountType: 'Test',
      provider: 'Facebook',
      sessionStatus: 'error',
      lastCheckedAt: 'Vừa xong',
      isDefault: false,
      note: 'Tài khoản sandbox, hiện báo lỗi đăng nhập (Check point).'
    },
    {
      id: 'acc_104',
      accountId: '29810398412',
      displayName: 'Trần B (Clone)',
      avatar: 'https://ui-avatars.com/api/?name=TB&background=f3f4f6&color=4b5563',
      accountType: 'Profile',
      provider: 'Facebook',
      sessionStatus: 'unchecked',
      lastCheckedAt: 'Chưa từng kiểm tra',
      isDefault: false,
      note: 'Mới import qua session cookie cục bộ.'
    }
  ]
}
