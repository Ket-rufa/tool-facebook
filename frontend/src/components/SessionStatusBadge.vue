<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
}>()

const statusMap: Record<string, { label: string, type: string }> = {
  active: { label: 'Hoạt động', type: 'success' },
  expired: { label: 'Hết hạn', type: 'warning' },
  error: { label: 'Lỗi', type: 'danger' },
  unchecked: { label: 'Chưa kiểm tra', type: 'neutral' }
}

const badgeBase = computed(() => {
  return statusMap[props.status] || { label: 'Unknown', type: 'neutral' }
})
</script>

<template>
  <span class="session-badge" :class="`badge-${badgeBase.type}`">
    <span class="badge-dot"></span>
    {{ badgeBase.label }}
  </span>
</template>

<style scoped>
.session-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 12px;
}

.badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.badge-success { background: #dcfce7; color: #166534; }
.badge-success .badge-dot { background: #22c55e; }

.badge-warning { background: #fef9c3; color: #854d0e; }
.badge-warning .badge-dot { background: #eab308; }

.badge-danger { background: #fee2e2; color: #991b1b; }
.badge-danger .badge-dot { background: #ef4444; }

.badge-neutral { background: #f3f4f6; color: #4b5563; }
.badge-neutral .badge-dot { background: #9ca3af; }
</style>
