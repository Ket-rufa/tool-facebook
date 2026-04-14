<script setup lang="ts">
defineProps<{
  statusText: string
  steps: { name: string; time: string; status: string }[]
}>()
</script>

<template>
  <div class="progress-card">
    <div class="progress-header">
      <div class="titles">
        <h4>Tiến trình trích xuất</h4>
        <span class="status-text">{{ statusText }}</span>
      </div>
      <button class="btn-retry">CHẠY LẠI</button>
    </div>

    <div class="steps-list">
      <div v-for="(step, idx) in steps" :key="idx" class="step-item">
        <div class="step-icon" :class="step.status">
          <svg v-if="step.status === 'success'" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
          <div v-else-if="step.status === 'running'" class="spinner"></div>
          <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
        </div>
        <div class="step-name">{{ step.name }}</div>
        <div class="step-time">{{ step.time }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.progress-card {
  background: #f8fafc;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  padding: 16px;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.titles h4 {
  margin: 0 0 4px;
  font-size: 14px;
  font-weight: 600;
  color: var(--c-text-title);
}

.status-text {
  font-size: 12px;
  color: var(--c-text-muted);
}

.btn-retry {
  background: var(--c-primary);
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
}

.btn-retry:hover {
  background: var(--c-primary-hover);
}

.steps-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.step-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.step-icon {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.step-icon.success { background-color: #10b981; color: white; }
.step-icon.error { background-color: #ef4444; color: white; }
.step-icon.running { border: 2px solid #e5e7eb; border-top-color: var(--c-primary); background: transparent; }

.spinner {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 2px solid transparent;
  border-top-color: inherit;
  animation: spin 1s linear infinite;
}

@keyframes spin { 100% { transform: rotate(360deg); } }

.step-name {
  flex: 1;
  font-size: 13px;
  color: var(--c-text-main);
}

.step-time {
  font-size: 12px;
  color: var(--c-text-muted);
}
</style>
