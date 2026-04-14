<script setup lang="ts">
defineProps<{
  options: string[]
  modelValue: string
}>()

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div class="segmented-control">
    <button
      v-for="opt in options"
      :key="opt"
      class="seg-btn"
      :class="{ active: modelValue === opt }"
      @click="emit('update:modelValue', opt)"
    >
      {{ opt }}
    </button>
  </div>
</template>

<style scoped>
.segmented-control {
  display: flex;
}

.seg-btn {
  flex: 1;
  background: white;
  border: 1px solid var(--c-border);
  color: var(--c-text-muted);
  font-size: 13px;
  font-weight: 600;
  padding: 8px 0;
  cursor: pointer;
  margin-left: -1px; /* border collapse */
  transition: all 0.2s;
}

.seg-btn:first-child {
  margin-left: 0;
  border-top-left-radius: var(--radius-sm);
  border-bottom-left-radius: var(--radius-sm);
}

.seg-btn:last-child {
  border-top-right-radius: var(--radius-sm);
  border-bottom-right-radius: var(--radius-sm);
}

.seg-btn:hover {
  background: #f9fafb;
  color: var(--c-text-main);
}

.seg-btn.active {
  background: #eff6ff;
  border-color: var(--c-primary);
  color: var(--c-primary);
  z-index: 2; /* overlap the adjacent borders */
}
</style>
