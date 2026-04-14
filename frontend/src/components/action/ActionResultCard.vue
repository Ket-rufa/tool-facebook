<script setup lang="ts">
import { computed } from 'vue'
import type { ActionResponse } from '../../types/action'
import { icons } from '../../utils/icons'

const props = defineProps<{
  result: ActionResponse | null
}>()

const uiState = computed(() => {
  if (!props.result) return null
  switch (props.result.status) {
    case 'success':
      return { cssClass: 'state-success', icon: icons.checkCircle, title: 'Thực thi thành công' }
    case 'validation_error':
    case 'session_invalid':
    case 'failed':
      return { cssClass: 'state-error', icon: icons.alertCircle, title: 'Thực thi thất bại' }
    case 'not_enabled':
      return { cssClass: 'state-warning', icon: icons.info, title: 'Từ chối thực thi' }
    default:
      return { cssClass: 'state-info', icon: icons.info, title: 'Thông tin bổ sung' }
  }
})

const formatTime = (timeStr: string) => {
  if (!timeStr) return ''
  return new Date(timeStr).toLocaleString('vi-VN')
}
</script>

<template>
  <div v-if="result && uiState" class="result-alert" :class="uiState.cssClass">
    <div class="ra-icon" v-html="uiState.icon"></div>
    <div class="ra-content">
      <h4 class="ra-title">{{ uiState.title }}</h4>
      <p class="ra-msg">{{ result.message }}</p>
      
      <div class="ra-meta">
        <span class="meta-item">
          <span v-html="icons.clock"></span>
          {{ formatTime(result.executed_at) }}
        </span>
        <span class="meta-item" v-if="result.account_id">
          <span v-html="icons.users"></span>
          {{ result.account_id }}
        </span>
        <span class="meta-item capitalize" v-if="result.reaction_type">
          {{ result.reaction_type }}
        </span>
        <span class="meta-badge" :class="result.dry_run ? 'bd-dry' : 'bd-real'">
          {{ result.dry_run ? 'Dry Run' : 'Real Run' }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.result-alert {
  display: flex;
  gap: 16px;
  padding: 16px;
  border-radius: var(--radius-md);
  border: 1px solid transparent;
}

.ra-icon {
  display: flex;
  flex-shrink: 0;
  margin-top: 2px;
}
.ra-icon :deep(svg) { width: 20px; height: 20px; }

.ra-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.ra-title {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
}

.ra-msg {
  margin: 0;
  font-size: 13px;
  line-height: 1.5;
  opacity: 0.9;
}

.ra-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 8px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  opacity: 0.8;
  font-weight: 500;
}
.meta-item :deep(svg) { width: 12px; height: 12px; }

.meta-badge {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid currentColor;
}
.bd-dry { opacity: 0.8; }
.bd-real { border-style: dashed; }

/* Theming per state */
.state-success {
  background: #f0fdf4;
  border-color: #bbf7d0;
  color: #166534;
}
.state-success .ra-icon { color: #15803d; }

.state-error {
  background: #fef2f2;
  border-color: #fecaca;
  color: #991b1b;
}
.state-error .ra-icon { color: #b91c1c; }

.state-warning {
  background: #fffbeb;
  border-color: #fde68a;
  color: #92400e;
}
.state-warning .ra-icon { color: #d97706; }

.state-info {
  background: #eff6ff;
  border-color: #bfdbfe;
  color: #1e40af;
}
.state-info .ra-icon { color: #1d4ed8; }
</style>
