<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  summary: {
    time: string
    accuracy: string
    size: string
  }
}>()
</script>

<template>
  <div class="crawl-summary-card">
    <div class="cs-header">
      <span class="icon text-primary" v-html="icons.barChart"></span>
      <h3 class="cs-title">Tóm tắt kết quả trích xuất (Crawl Summary)</h3>
    </div>
    
    <div class="cs-grid">
      <div class="cs-box">
        <label>THỜI GIAN TRÍCH XUẤT</label>
        <div class="cs-val">{{ summary.time }}</div>
      </div>
      <div class="cs-box">
        <label>ĐỘ CHÍNH XÁC METADATA</label>
        <div class="cs-val text-danger">{{ summary.accuracy }}</div>
      </div>
      <div class="cs-box">
        <label>DUNG LƯỢNG BẢN GHI</label>
        <div class="cs-val">{{ summary.size }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.crawl-summary-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.cs-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.icon {
  display: flex;
  align-items: center;
}

.text-primary {
  color: var(--c-primary);
}

.cs-title {
  margin: 0;
  font-size: 13px;
  font-weight: 700;
  color: var(--c-text-title);
}

.cs-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.cs-box {
  background: #f9fafb;
  border-radius: var(--radius-sm);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cs-box label {
  font-size: 10px;
  font-weight: 700;
  color: var(--c-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.cs-val {
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
  white-space: pre-line;
  line-height: 1.4;
}

.text-danger { color: #dc2626; }
</style>
