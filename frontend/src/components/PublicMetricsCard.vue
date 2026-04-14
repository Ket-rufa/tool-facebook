<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  metrics: {
    likes: string
    comments: string
    shares: string
  }
}>()
</script>

<template>
  <div class="card bg-white">
    <div class="card-title">CHỈ SỐ CÔNG KHAI</div>
    
    <div class="metrics-list">
      <div class="m-row bg-gray">
        <span class="m-label"><span class="m-icon text-primary"><svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" stroke="none"><path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path></svg></span> Cảm xúc</span>
        <span class="m-val">{{ metrics.likes }}</span>
      </div>
      
      <div class="m-row bg-gray">
        <span class="m-label"><span class="m-icon text-muted"><svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" stroke="none"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg></span> Bình luận</span>
        <span class="m-val">{{ metrics.comments }}</span>
      </div>
      
      <div class="m-row bg-gray">
        <span class="m-label"><span class="m-icon text-danger"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="18" cy="5" r="3"></circle><circle cx="6" cy="12" r="3"></circle><circle cx="18" cy="19" r="3"></circle><line x1="8.59" y1="13.51" x2="15.42" y2="17.49"></line><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"></line></svg></span> Chia sẻ</span>
        <span class="m-val">{{ metrics.shares }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.card-title {
  font-size: 11px;
  font-weight: 700;
  color: var(--c-text-muted);
  letter-spacing: 0.5px;
  margin-bottom: 16px;
}

.metrics-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.m-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-radius: var(--radius-sm);
}

.bg-gray {
  background: #f9fafb;
}

.m-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
  color: var(--c-text-main);
}

.m-icon {
  display: flex;
  align-items: center;
}

.text-primary { color: var(--c-primary); }
.text-danger { color: #dc2626; }
.text-muted { color: var(--c-text-muted); }

.m-val {
  font-size: 16px;
  font-weight: 700;
  color: var(--c-text-title);
}
</style>
