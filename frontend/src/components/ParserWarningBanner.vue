<script setup lang="ts">
defineProps<{
  text: string
}>()
</script>

<template>
  <div class="parser-warning-banner">
    <div class="warn-icon">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
    </div>
    <div class="warn-content">
      <h4 class="warn-title">Cảnh báo Parser (Cơ chế phân tích dữ liệu)</h4>
      <p class="warn-desc">{{ text }}</p>
    </div>
  </div>
</template>

<style scoped>
.parser-warning-banner {
  display: flex;
  gap: 16px;
  background: #fee2e2;
  border: 1px solid #fecaca;
  padding: 16px 20px;
  border-radius: var(--radius-md);
  color: #991b1b;
}

.warn-icon {
  margin-top: 2px;
  flex-shrink: 0;
}

.warn-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.warn-title {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
}

.warn-desc {
  margin: 0;
  font-size: 13px;
  line-height: 1.5;
  color: #b91c1c;
}
</style>
