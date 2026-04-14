<script setup lang="ts">
import { mediaData } from '../data/mediaMockData'
</script>

<template>
  <div class="storage-usage-card">
    <div class="s-title">Dung lượng lưu trữ</div>
    <div class="s-bar">
      <div class="s-progress" :style="{ width: mediaData.storage.percent + '%' }"></div>
    </div>
    <div class="s-text">
      {{ mediaData.storage.usedGB }} GB / {{ mediaData.storage.totalGB }} GB ({{ mediaData.storage.percent }}%)
    </div>
  </div>
</template>

<style scoped>
.storage-usage-card {
  margin: 24px 16px 16px;
  background: #eff6ff;
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.s-title {
  font-size: 11px;
  font-weight: 700;
  color: var(--c-text-title);
  letter-spacing: 0.5px;
}

.s-bar {
  width: 100%;
  height: 4px;
  background: #dbeafe;
  border-radius: 2px;
  overflow: hidden;
  margin: 4px 0;
}

.s-progress {
  height: 100%;
  background: var(--c-primary);
  border-radius: 2px;
}

.s-text {
  font-size: 10px;
  font-weight: 600;
  color: var(--c-primary);
}
</style>
