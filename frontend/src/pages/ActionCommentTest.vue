<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ActionCommentPanel from '../components/action/ActionCommentPanel.vue'
import ActionLogList from '../components/action/ActionLogList.vue'
import { actionsService } from '../services/actionsService'
import { sessionService } from '../services/sessionService'
import type { ActionResponse, SessionStatus, ActionLog, CommentRequest } from '../types/action'

const isLoading = ref(false)
const sessionInfo = ref<SessionStatus | null>(null)
const latestResult = ref<ActionResponse | null>(null)

const actionLogs = ref<ActionLog[]>([])

const loadData = async () => {
  try {
    actionLogs.value = await actionsService.getActionLogs()
  } catch (e) {
    console.error(e)
  }
}

const handleTestComment = async (payloads: CommentRequest[]) => {
  isLoading.value = true
  latestResult.value = null
  
  try {
    for (const payload of payloads) {
      const res = await actionsService.commentPost(payload)
      latestResult.value = res
    }
  } catch (error) {
    console.error('Test Action error', error)
  } finally {
    isLoading.value = false
    loadData() // Refresh logs
  }
}

const handleAccountUpdate = async (accountId: string) => {
  try {
    sessionInfo.value = null
    sessionInfo.value = await sessionService.getSessionStatus(accountId)
  } catch (error) {
    console.error('Error fetching session for account', error)
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="test-page">
    <div class="page-header">
      <h1 class="page-title">Thử nghiệm Bình luận</h1>
      <p class="page-desc">Giả lập và kiểm tra hành vi auto comment (random nội dung) lên bài viết.</p>
    </div>

    <div class="test-layout">
      <!-- Cột trái: Form cấu hình -->
      <div class="col-left">
        <ActionCommentPanel 
          :isLoading="isLoading" 
          @submit="handleTestComment"
          @updateAccount="handleAccountUpdate"
        />
        
        <div v-show="latestResult" class="result-wrapper">
          <div class="result-header">Kết quả trả về gần nhất</div>
          <pre class="result-box" :class="{ 'is-success': latestResult?.success }">{{ JSON.stringify(latestResult, null, 2) }}</pre>
        </div>
      </div>

      <!-- Cột phải: Thông tin Session & Lịch sử -->
      <div class="col-right">
        <!-- Dashboard Session Status -->
        <div class="card bg-white mb-4">
          <div class="card-header pb-0">
            <h3 class="card-title text-sm m-0">Trạng thái phiên đăng nhập</h3>
            <span class="badge" v-if="sessionInfo?.provider">{{ sessionInfo.provider.toUpperCase() }}</span>
          </div>
          <div class="card-body">
            <div v-if="sessionInfo" class="session-box">
              <div class="status-indicator" :class="{
                'bg-success': sessionInfo.is_active,
                'bg-danger': !sessionInfo.is_active
              }">
                <span class="status-dot"></span>
                <span>{{ sessionInfo.is_active ? 'Phiên hoạt động' : 'Hết hạn / Lỗi' }}</span>
              </div>
              <p class="session-msg mt-2 text-xs text-muted" v-if="sessionInfo.message">
                {{ sessionInfo.message }}
              </p>
            </div>
            <div v-else class="text-muted text-xs">
              Đang tải dữ liệu phiên... (Vui lòng chọn tài khoản)
            </div>
          </div>
        </div>

        <!-- Bảng Logs -->
        <ActionLogList :logs="actionLogs" :loading="isLoading" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.test-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  margin-bottom: 8px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--c-text-title);
  margin: 0 0 8px 0;
}

.page-desc {
  font-size: 14px;
  color: var(--c-text-muted);
  margin: 0;
}

.test-layout {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 24px;
  align-items: start;
}

.result-wrapper {
  margin-top: 24px;
}

.result-header {
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.result-box {
  background: var(--c-bg);
  border: 1px solid var(--c-border);
  padding: 16px;
  border-radius: var(--radius-md);
  font-family: monospace;
  font-size: 13px;
  color: var(--c-text-main);
  overflow-x: auto;
  white-space: pre-wrap;
}

.result-box.is-success {
  border-left: 4px solid var(--c-success);
}

.card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}
.mb-4 { margin-bottom: 24px; }
.card-header {
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--c-border-light);
}
.pb-0 { padding-bottom: 0; border-bottom: none; }
.card-body { padding: 16px; }
.text-sm { font-size: 14px; }
.m-0 { margin: 0; }
.badge { font-size: 10px; font-weight: 700; background: var(--c-border-light); padding: 2px 6px; border-radius: 4px; color: var(--c-text-muted); }
.text-xs { font-size: 12px; }
.text-muted { color: var(--c-text-muted); }
.mt-2 { margin-top: 8px; }

.session-box {
  background: #f8fafc;
  border: 1px dashed var(--c-border);
  padding: 12px;
  border-radius: var(--radius-sm);
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 20px;
}
.status-dot {
  width: 8px; height: 8px; border-radius: 50%; background: currentColor;
}
.bg-success { background: #dcfce7; color: #166534; }
.bg-danger { background: #fee2e2; color: #991b1b; }
</style>
