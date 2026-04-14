<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  title: string
  icon?: string
}>()
</script>

<template>
  <div class="settings-section">
    <div class="ss-header">
      <span class="ss-icon text-primary" v-if="icon" v-html="icons[icon as keyof typeof icons]"></span>
      <h3 class="ss-title">{{ title }}</h3>
    </div>
    <div class="ss-body">
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
.settings-section {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  margin-bottom: 24px;
  overflow: hidden;
}

.ss-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  border-bottom: 1px solid var(--c-border);
  background: #f9fafb;
}

.ss-icon {
  display: flex;
  align-items: center;
}
.ss-icon :deep(svg) { width: 16px; height: 16px; stroke-width: 2.5; }
.text-primary { color: var(--c-primary); }

.ss-title {
  margin: 0;
  font-size: 11px;
  font-weight: 800;
  color: var(--c-text-title);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.ss-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
</style>
