<script setup lang="ts">
import { ref } from 'vue'

const tabs = ['Tất cả', 'Hình ảnh', 'Video']
const activeTab = ref(tabs[0])
</script>

<template>
  <div class="media-filter-tabs">
    <button 
      v-for="t in tabs" 
      :key="t"
      class="m-tab"
      :class="{ 'active': activeTab === t }"
      @click="activeTab = t"
    >
      {{ t }}
    </button>
  </div>
</template>

<style scoped>
.media-filter-tabs {
  display: inline-flex;
  background: #f3f4f6;
  border-radius: var(--radius-sm);
  padding: 4px;
}

.m-tab {
  background: transparent;
  border: none;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-muted);
  padding: 6px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.m-tab:hover {
  color: var(--c-text-main);
}

.m-tab.active {
  background: white;
  color: var(--c-text-title);
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}
</style>
