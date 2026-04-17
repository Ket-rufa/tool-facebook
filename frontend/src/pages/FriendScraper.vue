<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { FetchFriends } from '../../wailsjs/go/crawl/CrawlHandler'
import { ListAccounts } from '../../wailsjs/go/accounts/AccountService'
import { accounts } from '../../wailsjs/go/models'

const availableAccounts = ref<accounts.AccountProfile[]>([])
const accountId = ref('')
const targetId = ref('')
const docId = ref('')
const loading = ref(false)
const message = ref('')
const totalCount = ref(0)
const friends = ref<{ uid: string; name: string; avatar: string }[]>([])
const searchQuery = ref('')

onMounted(async () => {
  try {
    const list = await ListAccounts()
    availableAccounts.value = list || []
    if (list && list.length > 0) {
      const def = list.find((a: accounts.AccountProfile) => a.isDefault)
      accountId.value = def ? def.id : list[0].id
    }
  } catch (e) {
    console.error('Lỗi tải tài khoản', e)
  }
})

const filteredFriends = () => {
  if (!searchQuery.value.trim()) return friends.value
  const q = searchQuery.value.toLowerCase()
  return friends.value.filter(f => f.name.toLowerCase().includes(q) || f.uid.includes(q))
}

const handleFetch = async () => {
  if (!accountId.value) {
    message.value = '⚠️ Vui lòng chọn tài khoản thực hiện.'
    return
  }
  loading.value = true
  message.value = '⚡ Đang phân tích mục tiêu và quét dữ liệu...'
  friends.value = []
  totalCount.value = 0
  searchQuery.value = ''

  try {
    // Lưu ý: FetchFriends giờ nhận 3 tham số: accountID, targetID, docID
    const resp = await FetchFriends(accountId.value, targetId.value, docId.value)
    if (resp.success) {
      friends.value = resp.friends || []
      totalCount.value = resp.total_count || friends.value.length
      message.value = resp.message
    } else {
      message.value = '❌ ' + (resp.message || 'Lỗi không xác định khi quét.')
    }
  } catch (e: any) {
    message.value = '❌ Lỗi hệ thống: ' + (e?.message || String(e))
  } finally {
    loading.value = false
  }
}

const exportCSV = () => {
  if (friends.value.length === 0) return
  
  const headers = "UID,Name,Avatar\n"
  const rows = friends.value.map(f => `"${f.uid}","${f.name.replace(/"/g, '""')}","${f.avatar}"`).join("\n")
  const csvContent = headers + rows
  
  // Thêm BOM (\uFEFF) để Excel nhận diện đúng bộ mã UTF-8 (Fix lỗi phông chữ tiếng Việt)
  const BOM = "\uFEFF";
  const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  
  const link = document.createElement("a")
  link.setAttribute("href", url)
  link.setAttribute("download", `friends_export_${targetId.value || 'self'}_${new Date().getTime()}.csv`)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  // Giải phóng bộ nhớ
  setTimeout(() => URL.revokeObjectURL(url), 100)
}

const avatarFallback = (name: string) =>
  `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}&background=6366f1&color=fff&size=100`

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text)
  message.value = '📋 Đã sao chép UID: ' + text
}
</script>

<template>
  <div class="friend-page">
    <div class="page-header">
      <div class="header-content">
        <div class="title-area">
          <h1>Quét Bạn bè Công khai</h1>
          <p>Thu thập thông tin UID, tên và avatar của đối tượng đích.</p>
        </div>
        <div class="header-actions" v-if="friends.length > 0">
          <button class="btn-export" @click="exportCSV">
            📥 Xuất CSV ({{ friends.length }})
          </button>
        </div>
      </div>
    </div>

    <!-- Main Entry Card -->
    <div class="card glass-card main-config">
      <div class="card-body">
        <div class="form-grid">
          <div class="form-group">
            <label>Tài khoản thực hiện:</label>
            <div class="select-wrapper">
              <select v-model="accountId" class="select-box" :disabled="loading">
                <option disabled value="">-- Chọn tài khoản --</option>
                <option v-for="acc in availableAccounts" :key="acc.id" :value="acc.id">
                  {{ acc.displayName || 'Tài khoản không tên' }}
                </option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label>Đối tượng quét (UID/Username):</label>
            <input 
              v-model="targetId" 
              type="text" 
              class="input-box" 
              placeholder="Để trống để quét chính mình..."
              :disabled="loading"
            />
          </div>

          <div class="form-group">
            <label>Doc ID (Nâng cao):</label>
            <input 
              v-model="docId" 
              type="text" 
              class="input-box" 
              placeholder="Tự động nếu để trống..."
              :disabled="loading"
            />
          </div>

          <div class="form-group action-group">
            <button
              class="btn-gradient-primary"
              @click="handleFetch"
              :disabled="loading || !accountId"
            >
              <span v-if="loading" class="spinner"></span>
              {{ loading ? 'Đang quét...' : '🚀 Bắt đầu Quét' }}
            </button>
          </div>
        </div>

        <div v-if="message" class="status-banner" :class="{ 'error': message.includes('❌') || message.includes('⚠️') }">
          <span class="status-icon">{{ message.includes('❌') ? '🚫' : (message.includes('✅') ? '✨' : 'ℹ️') }}</span>
          {{ message }}
        </div>
      </div>
    </div>

    <!-- Data Section -->
    <div class="data-section" v-if="friends.length > 0">
      <div class="stats-row">
        <div class="stat-bubble">
          <span class="bubble-val">{{ totalCount.toLocaleString() }}</span>
          <span class="bubble-lab">Tổng số bạn</span>
        </div>
        <div class="stat-bubble">
          <span class="bubble-val">{{ friends.length.toLocaleString() }}</span>
          <span class="bubble-lab">Đã lấy được</span>
        </div>
        <div class="search-wrap">
          <input
            v-model="searchQuery"
            type="text"
            class="search-bar"
            placeholder="🔎 Lọc theo tên hoặc ID..."
          />
        </div>
      </div>

      <div class="friends-container">
        <div
          v-for="friend in filteredFriends()"
          :key="friend.uid"
          class="friend-item-card"
        >
          <div class="avatar-wrap">
            <img
              :src="friend.avatar || avatarFallback(friend.name)"
              :alt="friend.name"
              class="friend-img"
              @error="($event.target as HTMLImageElement).src = avatarFallback(friend.name)"
            />
          </div>
          <div class="friend-meta">
            <div class="f-name" :title="friend.name">{{ friend.name }}</div>
            <div class="f-uid">{{ friend.uid }}</div>
          </div>
          <div class="item-actions">
             <button class="btn-icon-tiny" title="Copy UID" @click="copyToClipboard(friend.uid)">📋</button>
          </div>
        </div>
        
        <div v-if="filteredFriends().length === 0" class="no-results">
          Mục tiêu không khớp với từ khóa tìm kiếm.
        </div>
      </div>
    </div>
    
    <div class="empty-state" v-else-if="!loading">
       <div class="empty-icon">👥</div>
       <h3>Chưa có dữ liệu</h3>
       <p>Nhập UID đối tượng và nhấn nút Bắt đầu Quét để thu thập danh sách bạn bè.</p>
    </div>
  </div>
</template>

<style scoped>
.friend-page {
  padding: 30px;
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f4ff 0%, #ffffff 100%);
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-area h1 {
  font-size: 28px;
  font-weight: 850;
  background: linear-gradient(90deg, #1e293b, #4f46e5);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.title-area p {
  color: #64748b;
  margin-top: 5px;
  font-size: 15px;
}

/* Glass Card */
.glass-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.04);
}

.card-body {
  padding: 24px;
}

/* Form Grid */
.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  align-items: flex-end;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 13px;
  font-weight: 700;
  color: #475569;
  padding-left: 4px;
}

.select-box, .input-box {
  background: white;
  border: 1.5px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px 16px;
  font-size: 14px;
  transition: all 0.2s;
  outline: none;
}

.select-box:focus, .input-box:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.1);
}

/* Button Gradient */
.btn-gradient-primary {
  height: 48px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
}

.btn-gradient-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(79, 70, 229, 0.4);
}

.btn-gradient-primary:active {
  transform: translateY(0);
}

.btn-export {
  background: white;
  color: #4f46e5;
  border: 1.5px solid #4f46e5;
  padding: 10px 18px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-export:hover {
  background: #f5f3ff;
}

/* Status Banner */
.status-banner {
  margin-top: 20px;
  padding: 14px 18px;
  background: #f0fdf4;
  border: 1px solid #bcf0da;
  border-radius: 12px;
  color: #166534;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-banner.error {
  background: #fef2f2;
  border-color: #fee2e2;
  color: #991b1b;
}

/* Stats Row */
.stats-row {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
}

.stat-bubble {
  background: white;
  padding: 12px 20px;
  border-radius: 16px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.02);
  display: flex;
  flex-direction: column;
}

.bubble-val {
  font-size: 20px;
  font-weight: 800;
  color: #4f46e5;
}

.bubble-lab {
  font-size: 11px;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.search-wrap {
  flex: 1;
}

.search-bar {
  width: 100%;
  padding: 14px 20px;
  border-radius: 16px;
  border: 1.5px solid #e2e8f0;
  background: white;
  outline: none;
}

/* Friend Cards Grid */
.friends-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 15px;
}

.friend-item-card {
  display: flex;
  align-items: center;
  gap: 15px;
  background: white;
  padding: 15px;
  border-radius: 16px;
  border: 1px solid #f1f5f9;
  transition: all 0.2s;
}

.friend-item-card:hover {
  border-color: #6366f1;
  box-shadow: 0 4px 20px rgba(99, 102, 241, 0.08);
  transform: translateX(4px);
}

.friend-img {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  object-fit: cover;
}

.f-name {
  font-weight: 700;
  color: #1e293b;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 150px;
}

.f-uid {
  font-size: 11px;
  color: #94a3b8;
  font-family: monospace;
}

.item-actions {
  margin-left: auto;
}

.btn-icon-tiny {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 4px 8px;
  border-radius: 6px;
  cursor: pointer;
}

/* Empty State */
.empty-state {
  margin-top: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 60px;
  background: rgba(255, 255, 255, 0.4);
  border: 2px dashed #e2e8f0;
  border-radius: 30px;
}

.empty-icon {
  font-size: 60px;
  margin-bottom: 20px;
  opacity: 0.5;
}

.spinner {
  width: 20px; height: 20px;
  border: 3px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

</style>
