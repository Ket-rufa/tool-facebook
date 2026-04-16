<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { crawlService } from '../services/crawlService'
import { ListAccounts } from '../../wailsjs/go/accounts/AccountService'
import { crawl, accounts } from '../../wailsjs/go/models'

const targetType = ref('profile')
const targetId = ref('')
const accountId = ref('')
const limit = ref(10)
const scanSelf = ref(true) // Mặc định là quét chính nick đang chọn
const loading = ref(false)
const message = ref('')
const results = ref<any[]>([]) // Dùng any để tránh lỗi TS khi model chưa cập nhật kịp
const availableAccounts = ref<accounts.AccountProfile[]>([])

onMounted(async () => {
  try {
    const list = await ListAccounts()
    availableAccounts.value = list || []
    if (list && list.length > 0) {
      accountId.value = list[0].id // Auto-select first account
    }
  } catch (error) {
    console.error("Lỗi khi tải danh sách tài khoản", error)
  }
})

const handleCrawl = async () => {
  let finalTarget = targetId.value.trim()
  
  // TỰ ĐỘNG XỬ LÝ NẾU NGƯỜI DÙNG DÁN CẢ LINK
  if (finalTarget.includes('facebook.com')) {
    try {
      const url = new URL(finalTarget.startsWith('http') ? finalTarget : `https://${finalTarget}`)
      if (url.pathname.includes('/groups/')) {
        const parts = url.pathname.split('/')
        const idx = parts.indexOf('groups')
        if (idx !== -1 && parts[idx + 1]) {
          finalTarget = parts[idx + 1]
          targetType.value = 'group'
        }
      } else if (url.pathname.includes('profile.php')) {
        finalTarget = url.searchParams.get('id') || finalTarget
        targetType.value = 'profile'
      } else {
        const paths = url.pathname.split('/').filter(p => p && p !== 'profile.php')
        if (paths.length > 0) {
          finalTarget = paths[paths.length - 1]
        }
      }
    } catch (e) {
      console.error("Lỗi parse URL", e)
    }
  }

  if (!scanSelf.value && !finalTarget) {
    message.value = 'Vui lòng nhập ID (hoặc dán link) mục tiêu.'
    return
  }
  
  if (!accountId.value) {
    message.value = 'Vui lòng chọn Tài khoản thực hiện.'
    return
  }

  loading.value = true
  message.value = '⚡ Đang khởi tạo bộ máy quét Hybrid...'
  results.value = [] // Reset cũ
  
  try {
    const req = new crawl.CrawlRequest()
    req.target_id = scanSelf.value ? "" : finalTarget
    ;(req as any).type = scanSelf.value ? 'profile' : targetType.value
    req.account_id = accountId.value
    req.limit = Number(limit.value)
    // req.cursor = "" // Có thể mở rộng sau nếu cần nút "Tải thêm"

    const response = await crawlService.runCrawl(req)
    if (response.success) {
      results.value = response.data || []
      message.value = `✅ ${response.message}`
    } else {
      message.value = `❌ Lỗi: ${response.message}`
    }
  } catch (error: any) {
    message.value = `❌ Có lỗi hệ thống: ${error.message || error}`
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="workspace-page">
    <div class="page-header">
      <div class="title-area">
        <h1>Quét bài viết</h1>
        <p>Post Crawling / Scraping Workspace</p>
      </div>
    </div>

    <!-- Cấu hình và Chạy -->
    <div class="card config-card">
      <div class="card-header">Thiết lập Quét</div>
      <div class="card-body">
        
        <div class="form-row">
          <div class="form-group flex-2">
            <label>Chọn Tài khoản (Của bạn):</label>
            <select v-model="accountId" class="select-box">
              <option disabled value="">-- Vui lòng chọn tài khoản --</option>
              <option v-for="acc in availableAccounts" :key="acc.id" :value="acc.id">
                {{ acc.displayName || 'Không tên' }} ({{ acc.accountId || 'No ID' }})
              </option>
            </select>
          </div>
          <div class="form-group flex-1" style="justify-content: flex-end; padding-bottom: 12px;">
            <label class="checkbox-label">
              <input type="checkbox" v-model="scanSelf" />
              Quét chính nick này
            </label>
          </div>
        </div>

        <div class="form-row animate-if" v-if="!scanSelf">
          <div class="form-group flex-1">
            <label>Nơi muốn quét:</label>
            <select v-model="targetType" class="select-box">
              <option value="profile">Trang cá nhân / Fanpage khác</option>
              <option value="group">Trong Hội nhóm (Group)</option>
            </select>
          </div>
          <div class="form-group flex-2">
            <label>Link hoặc ID mục tiêu:</label>
            <input type="text" v-model="targetId" placeholder="Dán link hoặc ID người khác vào đây..." />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group flex-1">
            <label>Số lượng bài viết:</label>
            <input type="number" v-model="limit" min="1" max="100" />
          </div>
          <div class="form-group flex-2" style="justify-content: flex-end;">
            <div class="action-footer">
              <button class="btn-primary" @click="handleCrawl" :disabled="loading" style="min-width: 150px;">
                {{ loading ? 'Đang thực thi...' : 'Bắt đầu Quét' }}
              </button>
            </div>
          </div>
        </div>
        
        <p v-if="message" class="status-msg" :class="{ error: message.includes('Lỗi') || message.includes('Vui lòng') }">{{ message }}</p>
      </div>
    </div>

    <!-- Kết quả -->
    <div class="card results-card" v-if="results.length > 0">
      <div class="card-header">Kết quả quét ({{ results.length }})</div>
      <div class="card-body p-0">
        <table class="results-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Thời gian</th>
              <th>Nội dung</th>
              <th>Cảm xúc</th>
              <th>Bình luận</th>
              <th>Media</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="post in results" :key="post.post_id">
              <td>{{ post.post_id || 'N/A' }}</td>
              <td>{{ post.time }}</td>
              <td class="content-cell" :title="post.content">{{ post.content }}</td>
              <td><span class="badge reactions">{{ post.reactions || '0' }}</span></td>
              <td><span class="badge comments">{{ post.comments || '0' }}</span></td>
              <td>
                <div v-if="post.media_urls && post.media_urls.length > 0" class="image-gallery">
                  <img v-for="(url, idx) in post.media_urls.slice(0, 3)" :key="idx" :src="url" class="post-thumb" />
                  <span v-if="post.media_urls.length > 3" class="more-images">+{{ post.media_urls.length - 3 }}</span>
                </div>
                <span v-else class="text-muted">Không có</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.workspace-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  background-color: var(--c-background);
}
.page-header h1 {
  margin: 0;
  font-size: 24px;
  color: var(--c-text-main);
}
.page-header p {
  color: var(--c-text-muted);
  margin-top: 4px;
}
.card {
  background: #fff;
  border: 1px solid var(--c-border);
  border-radius: 8px;
  overflow: hidden;
}
.card-header {
  padding: 16px;
  border-bottom: 1px solid var(--c-border);
  font-weight: 600;
  background: #f9fafb;
}
.card-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.card-body.p-0 {
  padding: 0;
}
.form-row {
  display: flex;
  gap: 16px;
  width: 100%;
}
.flex-1 { flex: 1; }
.flex-2 { flex: 2; }
.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.form-group label {
  font-size: 13px;
  font-weight: 500;
}
.form-group input, .select-box {
  padding: 10px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  color: #111827;
  font-size: 14px;
}
.action-footer {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 8px;
}
.btn-primary {
  padding: 12px;
  background: var(--c-primary);
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
}
.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}
.status-msg {
  font-size: 13px;
  color: #059669;
}
.status-msg.error {
  color: #dc2626;
}
.results-table {
  width: 100%;
  border-collapse: collapse;
}
.results-table th, .results-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
  font-size: 13px;
}
.results-table th {
  background: #f9fafb;
  font-weight: 600;
  color: #4b5563;
}
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  color: var(--c-primary);
}
.checkbox-label input {
  width: 18px;
  height: 18px;
}
.badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: 600;
  font-size: 12px;
}
.badge.reactions { background: #e0f2fe; color: #0369a1; }
.badge.comments { background: #fef3c7; color: #92400e; }
.animate-if {
  animation: fadeIn 0.3s ease;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}
.image-gallery {
  display: flex;
  gap: 4px;
  align-items: center;
}
.post-thumb {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
}
.more-images {
  font-size: 11px;
  font-weight: bold;
  color: #6b7280;
  background: #f3f4f6;
  padding: 4px 6px;
  border-radius: 4px;
}
.text-muted {
  color: #9ca3af;
  font-style: italic;
}
</style>
