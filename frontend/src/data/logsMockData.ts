export const logsData = {
  list: [
    {
      id: 1,
      time: '2023-11-24 14:32:01',
      level: 'ERROR',
      module: 'API_HANDLER',
      job: 'FB_Crawl_R...',
      desc: 'Request...'
    },
    {
      id: 2,
      time: '2023-11-24 14:31:55',
      level: 'WARNING',
      module: 'CRAWLER',
      job: 'Group_Mon...',
      desc: 'High...'
    },
    {
      id: 3,
      time: '2023-11-24 14:31:40',
      level: 'INFO',
      module: 'PARSER',
      job: 'Auto_Impor...',
      desc: 'Success...'
    },
    {
      id: 4,
      time: '2023-11-24 14:30:12',
      level: 'INFO',
      module: 'SCHEDULER',
      job: 'Daily_Report',
      desc: 'Job...'
    },
    {
      id: 5,
      time: '2023-11-24 14:28:45',
      level: 'ERROR',
      module: 'BROWSER_ENGINE',
      job: 'Profile_Logi...',
      desc: 'Pro...'
    },
    {
      id: 6,
      time: '2023-11-24 14:27:10',
      level: 'INFO',
      module: 'API_HANDLER',
      job: 'FB_Crawl_R...',
      desc: 'Token...'
    },
    {
      id: 7,
      time: '2023-11-24 14:25:33',
      level: 'WARNING',
      module: 'CRAWLER',
      job: 'Group_Mon...',
      desc: 'Rate...'
    }
  ],
  detail: {
    rawLog: `{\n  "timestamp": "2023-11-24T14:32:01.442Z",\n  "level": "ERROR",\n  "module": "API_HANDLER",\n  "context": {\n    "job_id": "job_88291_retail",\n    "endpoint": "/v16.0/1029384756/posts",\n    "method": "GET"\n  },\n  "error": {\n    "code": 403,\n    "type": "OAuthException",\n    "message": "The user hasn't authorized the application to perform this action",\n    "fbtrace_id": "Abx7Y2Jk9Lp0-22"\n  },\n  "stacktrace": [\n    "at FacebookApi.request (/src/api/handler.js:142:9)",\n    "at processTicksAndRejections (node:internal/process/task_queues:96:5)",\n    "at async Crawler.fetchPosts (/src/core/crawler.js:44:21)"\n  ],\n  "system_info": {\n    "node_id": "worker-delta-02",\n    "memory_usage": "242MB",\n    "uptime": "14d 2h 11m"\n  }\n}`,
    workerId: 'worker-delta-02',
    correlationId: '550e8400-e29b-41d4-a716-446655440000'
  }
}
