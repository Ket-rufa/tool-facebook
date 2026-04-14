<script setup lang="ts">
import { useSlots } from 'vue'

const slots = useSlots()
</script>

<template>
  <div class="filter-bar">
    <div v-if="slots.label" class="filter-label">
      <slot name="label"></slot>
    </div>
    
    <div class="filter-items">
      <slot></slot>
    </div>

    <div class="filter-actions" v-if="slots.actions">
      <slot name="actions"></slot>
    </div>
  </div>
</template>

<style scoped>
.filter-bar {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  background-color: var(--c-surface);
  border: 1px solid var(--c-border);
  border-bottom: none;
  border-radius: var(--radius-md) var(--radius-md) 0 0;
  gap: 16px;
}

.filter-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--c-text-muted);
  text-transform: uppercase;
  margin-right: 8px;
}

.filter-items {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.filter-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Common filter select styles to be used via slots */
:deep(.filter-select) {
  padding: 8px 12px;
  background-color: #f3f4f6;
  border: 1px solid transparent;
  border-radius: var(--radius-sm);
  font-size: 13px;
  color: var(--c-text-main);
  outline: none;
  cursor: pointer;
  appearance: none;
  padding-right: 32px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
}

:deep(.filter-select:hover) {
  background-color: #e5e7eb;
}

:deep(.filter-chip) {
  padding: 6px 12px;
  border: 1px solid var(--c-border);
  border-radius: 16px;
  font-size: 13px;
  background: white;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
:deep(.filter-chip.active) {
  background-color: var(--c-primary);
  color: white;
  border-color: var(--c-primary);
}
</style>
