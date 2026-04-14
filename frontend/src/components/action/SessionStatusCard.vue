<script setup lang="ts">
import { computed } from 'vue'
import type { SessionStatus } from '../../types/action'
import { icons } from '../../utils/icons'

const props = defineProps<{
  session: SessionStatus | null
}>()

const isOk = computed(() => props.session?.is_active ?? false)
</script>

<template>
  <div class="card">
    <div class="card-header">
      <h3 class="card-title">Trạng thái phiên đăng nhập</h3>
      <span class="provider-badge" v-if="session">{{ session.provider }}</span>
    </div>
    
    <div class="card-body">
      <div v-if="!session" class="loading-state">
        Đang tải thông tin...
      </div>
      <div v-else class="status-box" :class="isOk ? 'is-success' : 'is-error'">
        <div class="status-icon" v-html="isOk ? icons.checkCircle : icons.alertCircle"></div>
        <div class="status-text">{{ session.message }}</div>
        <div class="status-dot">
          <span class="pulse-dot" :class="isOk ? 'pulse-green' : 'pulse-red'"></span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
}

.card-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--c-border-light);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
  color: var(--c-text-title);
}

.provider-badge {
  font-size: 11px;
  font-weight: 600;
  background: var(--c-border-light);
  padding: 2px 8px;
  border-radius: 12px;
  color: var(--c-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.card-body {
  padding: 20px;
}

.loading-state {
  font-size: 13px;
  color: var(--c-text-muted);
  font-style: italic;
}

.status-box {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  border: 1px solid transparent;
}

.status-icon {
  display: flex;
}
.status-icon :deep(svg) { width: 18px; height: 18px; }

.status-text {
  flex: 1;
  font-size: 14px;
  font-weight: 600;
}

/* Success State */
.status-box.is-success {
  background: #ecfdf5;
  border-color: #a7f3d0;
  color: #065f46;
}
.status-box.is-success .status-icon { color: #10b981; }

/* Error State */
.status-box.is-error {
  background: #fef2f2;
  border-color: #fecaca;
  color: #991b1b;
}
.status-box.is-error .status-icon { color: #ef4444; }


/* Pulser */
.status-dot {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 12px;
  height: 12px;
}

.pulse-dot {
  position: absolute;
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.pulse-green {
  background-color: #10b981;
  box-shadow: 0 0 0 2px rgba(16,185,129,0.3);
}

.pulse-red {
  background-color: #ef4444;
  box-shadow: 0 0 0 2px rgba(239,68,68,0.3);
}
</style>
