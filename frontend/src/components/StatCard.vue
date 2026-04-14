<script setup lang="ts">
import { computed } from 'vue'
import { icons } from '../utils/icons'

const props = defineProps<{
  id?: string
  label: string
  value: string | number
  subValue: string
  subType: string | 'success' | 'danger' | 'neutral'
  iconRight?: string
}>()

const statusIcon = computed(() => {
  if (props.subType === 'success') return icons.trendingUp
  if (props.subType === 'danger') return icons.alertCircle
  if (props.iconRight === 'dot-blue') return '<span class="status-dot blue"></span>'
  if (props.iconRight === 'check') return icons.checkCircle
  if (props.iconRight === 'alert') return icons.alertCircle
  return icons.media
})
</script>

<template>
  <div class="stat-card">
    <div class="label">{{ label }}</div>
    <div class="value-row">
      <div class="value" :class="{ 'text-primary': iconRight === 'dot-blue' || id === 'profiles', 'text-danger': subType === 'danger' }">{{ value }}</div>
      <div v-if="iconRight === 'dot-blue'" class="dot blue"></div>
    </div>
    
    <div class="sub-info" :class="subType">
      <span class="sub-icon" v-html="statusIcon"></span>
      <span>{{ subValue }}</span>
    </div>
  </div>
</template>

<style scoped>
.stat-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-shadow: var(--shadow-sm);
}

.label {
  font-size: 12px;
  font-weight: 600;
  color: var(--c-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.value-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.value {
  font-size: 36px;
  font-weight: 700;
  color: var(--c-text-title);
  line-height: 1;
}

.value.text-primary {
  color: var(--c-primary);
}

.value.text-danger {
  color: var(--c-danger);
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.dot.blue { background-color: var(--c-primary); }

.sub-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  margin-top: auto;
}

.sub-info.success { color: var(--c-text-main); }
.sub-info.success .sub-icon { color: var(--c-text-main); }

.sub-info.danger { color: var(--c-danger); }
.sub-info.danger .sub-icon { color: var(--c-danger); }

.sub-info.neutral { color: var(--c-text-muted); }
.sub-info.neutral .sub-icon { color: var(--c-text-muted); }

.sub-icon {
  display: flex;
  align-items: center;
}
</style>
