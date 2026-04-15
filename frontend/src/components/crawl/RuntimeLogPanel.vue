<script setup lang="ts">
defineProps<{
  logs: any[]
}>()
</script>

<template>
  <div class="crawl-panel runtime-log-panel">
    <div class="panel-header">
      <h3>Zone D: Nhật ký kỹ thuật</h3>
    </div>
    
    <div class="panel-body">
      <div class="log-container">
        <div v-for="log in logs" :key="log.id" class="log-line">
          <span class="log-time">[{{ log.time }}]</span>
          <span :class="['log-level', log.level]">{{ log.level.toUpperCase() }}</span>
          <span class="log-msg">{{ log.message }}</span>
          <span v-if="log.uid" class="log-meta">uid:{{ log.uid }}</span>
          <span v-if="log.step" class="log-meta">[step:{{ log.step }}]</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.crawl-panel {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

.panel-header {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
}
.panel-header h3 { margin: 0; font-size: 14px; font-weight: 600; color: #374151; }

.panel-body {
  padding: 0;
  display: flex;
  flex-direction: column;
  flex: 1;
  background: #111827;
}

.log-container {
  padding: 12px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #d1d5db;
  overflow-y: auto;
  max-height: 320px;
}

.log-line {
  margin-bottom: 4px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.log-time { color: #6b7280; }
.log-level { font-weight: 700; min-width: 45px; }
.log-level.info { color: #3b82f6; }
.log-level.warning { color: #f59e0b; }
.log-level.error { color: #ef4444; }

.log-msg { color: #f3f4f6; }
.log-meta { color: #9ca3af; font-size: 11px; }
</style>
