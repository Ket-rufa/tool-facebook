<script setup lang="ts">
defineProps<{
  type?: 'warning' | 'error' | 'info'
  text: string
}>()
</script>

<template>
  <div class="warning-box" :class="type || 'warning'">
    <div class="icon-wrap">
      <svg v-if="type === 'error'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>
      <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
    </div>
    <div class="warning-text">{{ text }}</div>
  </div>
</template>

<style scoped>
.warning-box {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--radius-sm);
  background-color: #fef3c7;
  color: #b45309;
}

.warning-box.error {
  background-color: #fee2e2;
  color: #b91c1c;
}

.icon-wrap {
  margin-top: 2px;
  flex-shrink: 0;
}

.warning-text {
  font-size: 13px;
  line-height: 1.5;
}
</style>
