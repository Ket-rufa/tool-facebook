<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  type: string // 'running', 'success', 'error' | 'INFO', 'WARNING', 'ERROR'
  text?: string
}>()

const statusConfig = computed(() => {
  switch (props.type) {
    case 'running': return { class: 'bg-warning-light text-warning', dot: 'bg-warning', label: props.text || 'Đang chạy' }
    case 'success': return { class: 'bg-success-light text-success', icon: 'check', label: props.text || 'Hoàn thành' }
    case 'error': return { class: 'bg-danger-light text-danger', icon: 'alert', label: props.text || 'Lỗi' }
    
    // Log levels
    case 'INFO': return { class: 'bg-primary-light text-primary', dot: 'bg-primary', label: 'INFO', isLog: true }
    case 'WARNING': return { class: 'bg-warning-light text-warning', icon: 'alert', label: 'WARNING', isLog: true }
    case 'ERROR': return { class: 'bg-danger-light text-danger', icon: 'alert', label: 'ERROR', isLog: true }
    default: return { class: 'bg-neutral-light text-neutral', label: props.type }
  }
})
</script>

<template>
  <span 
    class="status-badge" 
    :class="[statusConfig.class, { 'log-badge': statusConfig.isLog }]"
  >
    <span v-if="statusConfig.dot" class="dot" :class="statusConfig.dot"></span>
    <span v-if="statusConfig.icon === 'check'" class="icon">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
    </span>
    <span v-if="statusConfig.icon === 'alert'" class="icon">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
    </span>
    {{ statusConfig.label }}
  </span>
</template>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
}

.log-badge {
  border-radius: 4px;
  padding: 2px 8px;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.icon {
  display: flex;
  align-items: center;
}

/* Colors */
.bg-success-light { background-color: #d1fae5; }
.text-success { color: #059669; }

.bg-danger-light { background-color: #fee2e2; }
.text-danger { color: #dc2626; }

.bg-warning-light { background-color: #fef3c7; }
.text-warning { color: #d97706; }
.bg-warning { background-color: #d97706; }

.bg-primary-light { background-color: #eff6ff; }
.text-primary { color: #0f62fe; }
.bg-primary { background-color: #0f62fe; }
</style>
