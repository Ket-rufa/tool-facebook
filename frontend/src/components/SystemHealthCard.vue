<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  status: string
  modules: {
    id: string
    name: string
    status: string
    icon: string
  }[]
}>()

function getIcon(name: string) {
  return (icons as any)[name] || icons.box
}
</script>

<template>
  <div class="health-card">
    <div class="card-header">
      <h3>SỨC KHỎE HỆ THỐNG</h3>
      <span class="status-badge">{{ status }}</span>
    </div>

    <div class="modules-grid">
      <div v-for="mod in modules" :key="mod.id" class="module-item">
        <div class="module-icon" v-html="getIcon(mod.icon)"></div>
        <div class="module-info">
          <div class="module-name">{{ mod.name }}</div>
          <div class="module-status">{{ mod.status }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.health-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.card-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--c-text-title);
  letter-spacing: 0.5px;
}

.status-badge {
  font-size: 11px;
  font-weight: 700;
  color: var(--c-text-title);
  letter-spacing: 0.5px;
}

.modules-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.module-item {
  background-color: #f8fafc;
  border-radius: var(--radius-sm);
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.module-icon {
  color: var(--c-text-title);
  display: flex;
  align-items: center;
  justify-content: center;
}

.module-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.module-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--c-text-main);
}

.module-status {
  font-size: 12px;
  color: var(--c-text-muted);
}
</style>
