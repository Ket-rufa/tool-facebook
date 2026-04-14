<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  label: string
  icon?: string
}>()
</script>

<template>
  <div class="filter-chip">
    <span v-if="icon" class="ic" v-html="(icons as any)[icon]"></span>
    <span class="label">{{ label }}</span>
    <button class="btn-close">&times;</button>
  </div>
</template>

<style scoped>
.filter-chip {
  display: inline-flex;
  align-items: center;
  background: #f3f4f6;
  border-radius: 4px;
  padding: 6px 12px;
  gap: 8px;
}

.ic {
  display: flex;
  align-items: center;
  color: var(--c-text-muted);
}

.label {
  font-size: 13px;
  color: var(--c-text-main);
}

.btn-close {
  background: transparent;
  border: none;
  font-size: 16px;
  color: var(--c-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  line-height: 1;
}
.btn-close:hover { color: var(--c-text-title); }
</style>
