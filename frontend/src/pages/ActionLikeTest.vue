<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ActionTestPanel from '../components/action/ActionTestPanel.vue'
import SessionStatusCard from '../components/action/SessionStatusCard.vue'
import ActionResultCard from '../components/action/ActionResultCard.vue'
import ActionLogList from '../components/action/ActionLogList.vue'
import { actionsService } from '../services/actionsService'
import { sessionService } from '../services/sessionService'
import type { ActionResponse, SessionStatus, ActionLog, ActionRequest } from '../types/action'

const isLoading = ref(false)
const sessionInfo = ref<SessionStatus | null>(null)
const latestResult = ref<ActionResponse | null>(null)
const logs = ref<ActionLog[]>([])

const loadData = async () => {
  try {
    const [sessionStatus, recentLogs] = await Promise.all([
      sessionService.getSessionStatus(),
      actionsService.getActionLogs()
    ])
    sessionInfo.value = sessionStatus
    logs.value = recentLogs
  } catch (error) {
    console.error('Error loading test data', error)
  }
}

const handleTestLike = async (payloads: ActionRequest[]) => {
  isLoading.value = true
  latestResult.value = null
  
  try {
    for (const payload of payloads) {
      const res = await actionsService.likePost(payload)
      latestResult.value = res // Show the latest result conceptually, or update it per tick
    }
  } catch (error) {
    console.error('Test Action error', error)
  } finally {
    isLoading.value = false
    const recentLogs = await actionsService.getActionLogs()
    logs.value = recentLogs
  }
}

const handleAccountUpdate = async (accountId: string) => {
  try {
    sessionInfo.value = null // show loading state momentarily
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
      <h1 class="page-title">Thử nghiệm cảm xúc</h1>
      <p class="page-desc">Giả lập và kiểm tra hành vi thả cảm xúc lên bài viết trong môi trường an toàn.</p>
    </div>

    <div class="test-layout">
      <!-- Cột Trái -->
      <div class="col-left">
        <ActionTestPanel 
          :isLoading="isLoading" 
          @submit="handleTestLike"
          @updateAccount="handleAccountUpdate"
        />
        
        <div v-show="latestResult" class="result-wrapper">
          <ActionResultCard :result="latestResult" />
        </div>
      </div>

      <!-- Cột Phải -->
      <div class="col-right">
        <SessionStatusCard :session="sessionInfo" />
        <ActionLogList :logs="logs" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.test-page {
  display: flex;
  flex-direction: column;
  gap: 32px;
  max-width: 1200px;
}

.page-header {
  padding-bottom: 24px;
  border-bottom: 1px solid var(--c-border);
}

.page-title {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 800;
  color: var(--c-text-title);
}

.page-desc {
  margin: 0;
  font-size: 15px;
  color: var(--c-text-muted);
}

.test-layout {
  display: grid;
  grid-template-columns: 1.5fr 1fr;
  gap: 24px;
  align-items: start;
}

.col-left {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.col-right {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.result-wrapper {
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-8px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
