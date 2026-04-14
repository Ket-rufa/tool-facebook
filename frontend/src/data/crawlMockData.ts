export const crawlMockData = {
  // Zone A: Target input mock
  targets: [
    { id: '1', url: 'https://fb.com/zuck', displayName: 'Mark Zuckerberg', avatar: 'https://ui-avatars.com/api/?name=Mark', type: 'Profile', uid: '4', resolveStatus: 'Resolved', publicAccess: 'Accessible' },
    { id: '2', url: 'https://fb.com/groups/vuejs', displayName: 'Vue.js Developers', avatar: 'https://ui-avatars.com/api/?name=Vue', type: 'Group', uid: '123456', resolveStatus: 'Resolved', publicAccess: 'Partial' },
    { id: '3', url: 'https://fb.com/invalid123', displayName: null, avatar: null, type: 'Unknown', uid: null, resolveStatus: 'Failed', publicAccess: 'Inaccessible' }
  ],
  // Zone B: Config defaults
  config: {
    preset: 'Standard',
    fieldGroups: ['Identity', 'Basic info', 'Location', 'Work', 'Metrics'],
    concurrency: 5,
    delay: 2000,
    retry: 3,
    maxPosts: 10
  },
  // Zone C: Runtime summary
  runtime: {
    total: 3,
    queued: 1,
    running: 1,
    success: 0,
    failed: 1,
    partial: 0,
    eta: '02:15',
    progress: 35,
    currentEntity: 'Mark Zuckerberg',
    currentUid: '4',
    currentStep: 'đang đọc intro/about'
  },
  // Zone D: Logs
  logs: [
    { id: 1, time: '10:00:01', level: 'info', message: 'Resolve target success', uid: '4', step: 'resolve' },
    { id: 2, time: '10:00:05', level: 'info', message: 'Public intro extracted', uid: '4', step: 'intro' },
    { id: 3, time: '10:00:08', level: 'warning', message: 'Metrics partially available', uid: '4', step: 'metrics' },
    { id: 4, time: '10:01:00', level: 'error', message: 'Target inaccessible', uid: null, step: 'resolve' }
  ],
  // Bottom: Entity results
  entities: [
    {
      entity_id: 'e_1',
      uid: '4',
      entity_type: 'Profile',
      canonical_url: 'https://facebook.com/zuck',
      scan_status: 'running',
      last_updated: '2023-10-25T10:00:00Z',
      avatar: 'https://ui-avatars.com/api/?name=Mark',
      identity: {
        display_name: 'Mark Zuckerberg',
        alternate_name: null,
      },
      basic_info: {
        bio: 'Bringing the world closer together.',
        birthday_text: 'May 14, 1984',
      },
      location: {
        current_city: 'Palo Alto, California',
        hometown: 'Dobbs Ferry, New York',
      },
      education: [
        'Harvard University'
      ],
      work: [
        'Founder and CEO at Meta'
      ],
      metrics: {
        followers_text: '119M followers',
        followers_value: 119000000,
        following_text: null,
        following_value: null,
        friends_count_text: null,
        friends_count_value: null,
      },
      links: ['https://instagram.com/zuck'],
      field_status: 'partial'
    },
    {
      entity_id: 'e_2',
      uid: '123456',
      entity_type: 'Group',
      canonical_url: 'https://facebook.com/groups/vuejs',
      scan_status: 'queued',
      last_updated: null,
      avatar: 'https://ui-avatars.com/api/?name=Vue',
      identity: {
        display_name: 'Vue.js Developers',
        alternate_name: null,
      },
      basic_info: {
        bio: null,
        birthday_text: null,
      },
      location: {
        current_city: null,
        hometown: null,
      },
      education: [],
      work: [],
      metrics: {
        followers_text: null,
        followers_value: null,
        following_text: null,
        following_value: null,
        friends_count_text: null,
        friends_count_value: null,
      },
      links: [],
      field_status: 'pending'
    }
  ]
}
