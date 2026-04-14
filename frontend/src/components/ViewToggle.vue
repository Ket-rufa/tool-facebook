<script setup lang="ts">
import { ref } from 'vue'
import { icons } from '../utils/icons'

const view = ref('table') // 'table' or 'grid'
</script>

<template>
  <div class="view-toggle">
    <button 
      class="toggle-btn" 
      :class="{ active: view === 'table' }"
      @click="view = 'table'"
    >
      <span class="ic" v-html="icons.layout"></span>
      Bảng
    </button>
    <button 
      class="toggle-btn" 
      :class="{ active: view === 'grid' }"
      @click="view = 'grid'"
    >
      <span class="ic" v-html="icons.grid"></span>
      Lưới
    </button>
  </div>
</template>

<style scoped>
.view-toggle {
  display: inline-flex;
  background: #f3f4f6;
  border-radius: var(--radius-sm);
  padding: 4px;
}

.toggle-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: transparent;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-muted);
  cursor: pointer;
  transition: all 0.2s;
}

.toggle-btn:hover {
  color: var(--c-text-main);
}

.toggle-btn.active {
  background: white;
  color: var(--c-primary);
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

.ic {
  display: flex;
  align-items: center;
}
</style>
