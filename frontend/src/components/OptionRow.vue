<script setup lang="ts">
defineProps<{
  label: string
  description?: string
  layout?: 'horizontal' | 'vertical'
}>()
</script>

<template>
  <div class="option-row" :class="layout || 'vertical'">
    <div class="or-info">
      <label class="or-label">{{ label }}</label>
      <div v-if="description" class="or-desc">{{ description }}</div>
    </div>
    <div class="or-control">
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
.option-row {
  display: flex;
  gap: 16px;
}

.option-row.horizontal {
  align-items: center;
  justify-content: space-between;
}

.option-row.vertical {
  flex-direction: column;
}

.or-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.or-label {
  font-size: 13px;
  font-weight: 700;
  color: var(--c-text-title);
}

.or-desc {
  font-size: 12px;
  color: var(--c-text-muted);
}

.or-control {
  /* Let inner component decide sizing */
}
</style>
