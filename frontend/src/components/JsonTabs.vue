<script setup lang="ts">
import { ref } from 'vue'
import CodeViewer from './CodeViewer.vue'

defineProps<{
  jsonContent: string
}>()

const activeTab = ref('json')
</script>

<template>
  <div class="json-tabs-container">
    <div class="tabs-header">
      <button 
        class="tab-btn" 
        :class="{ active: activeTab === 'json' }" 
        @click="activeTab = 'json'"
      >
        Dữ liệu thô (JSON)
      </button>
      <button 
        class="tab-btn" 
        :class="{ active: activeTab === 'diff' }" 
        @click="activeTab = 'diff'"
      >
        Lịch sử quét (Diff)
      </button>
    </div>

    <div class="tab-content" v-if="activeTab === 'json'">
      <CodeViewer :code="jsonContent" />
    </div>
    
    <div class="tab-content" v-if="activeTab === 'diff'">
      <div class="diff-placeholder">
        <p>Tính năng xem Diff đang được phát triển...</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.json-tabs-container {
  display: flex;
  flex-direction: column;
}

.tabs-header {
  display: flex;
  border-bottom: 2px solid var(--c-border);
  margin-bottom: -1px; /* Align with border above code viewer manually or just let it gap */
}

.tab-btn {
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 13px;
  font-weight: 700;
  color: var(--c-text-muted);
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: -2px; /* Overlap border */
}

.tab-btn:hover { color: var(--c-text-main); }
.tab-btn.active {
  color: var(--c-primary);
  border-bottom-color: var(--c-primary);
}

.tab-content {
  margin-top: 16px;
}

.diff-placeholder {
  background: #111827;
  border-radius: var(--radius-sm);
  padding: 40px;
  text-align: center;
  color: #9ca3af;
  font-size: 13px;
  font-family: monospace;
}
</style>
