export const profilesData = {
  stats: {
    total: { value: '124', sub: '+12%' },
    scan24h: { value: '1,402', sub: 'Ổn định' },
    newData: { value: '42.5', unit: 'GB', sub: 'Tháng này' },
    status: { value: 'Đang hoạt động', isGood: true }
  },
  list: [
    { 
      id: 1, 
      name: 'Công nghệ Xanh VN', 
      handle: 'fb.com/congnghexanhvn', 
      followers: '852,401', 
      lastScanTime: '15:30\nHôm nay', 
      lastScanStatus: 'THÀNH CÔNG', 
      statusType: 'success'
    },
    { 
      id: 2, 
      name: 'Alpha Media Group', 
      handle: 'fb.com/alphamedia', 
      followers: '1.2M', 
      lastScanTime: '10:15\nHôm qua', 
      lastScanStatus: 'CHỜ XỬ LÝ', 
      statusType: 'neutral'
    },
    { 
      id: 3, 
      name: 'Đời sống 24/7', 
      handle: 'fb.com/doisong247', 
      followers: '2.4M', 
      lastScanTime: '08:00\nHôm nay', 
      lastScanStatus: 'LỖI KÉT NỐI', 
      statusType: 'danger'
    }
  ],
  detail: {
    id: 1,
    name: 'Công nghệ Xanh VN',
    handle: 'fb.com/congnghexanhvn',
    statusBadge: 'ĐANG QUÉT',
    followers: '852,401',
    engagement: '3.2%',
    description: 'Cung cấp giải pháp công nghệ bền vững cho doanh nghiệp Việt Nam. Tập trung vào năng lượng tái tạo và AI.',
    schedule: 'Mỗi 4 tiếng (00:00, 04:00, 08:00...)',
    coverImg: 'https://images.unsplash.com/photo-1550751827-4bd374c3f58b?auto=format&fit=crop&q=80&w=800',
    avatarImg: 'https://ui-avatars.com/api/?name=C+N&background=e0f2fe&color=0369a1&size=80'
  },
  latestPosts: [
    { 
      id: 101, 
      title: 'Khai trương trung tâm dữ liệu mới tại Hà Nội...', 
      time: '2 giờ trước', 
      engagement: '1.2k Likes', 
      image: 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?auto=format&fit=crop&q=80&w=120' 
    },
    { 
      id: 102, 
      title: 'Giải pháp Cloud Hybrid cho doanh nghiệp vừa và nhỏ', 
      time: 'Hôm qua', 
      engagement: '850 Likes', 
      image: 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?auto=format&fit=crop&q=80&w=120' 
    }
  ]
}
