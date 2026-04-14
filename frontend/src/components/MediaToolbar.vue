<script setup lang="ts">
import MediaFilterTabs from './MediaFilterTabs.vue'
import { icons } from '../utils/icons'
</script>

<template>
  <div class="media-toolbar">
    <div class="toolbar-left">
      <MediaFilterTabs />
      
      <div class="divider"></div>
      
      <div class="tools-group">
        <span class="t-icon" v-html="icons.users"></span>
        <select class="t-select">
          <option>Nguồn hồ sơ</option>
        </select>
      </div>

      <div class="tools-group">
        <span class="t-icon" v-html="icons.clock"></span>
        <select class="t-select">
          <option>Khoảng ngày</option>
        </select>
      </div>
    </div>
    
    <div class="toolbar-right">
      <span class="t-count">Hiển thị 1 - 24 trên 1,420</span>
      <div class="view-icons">
        <button class="v-btn active" v-html="icons.grid"></button>
        <button class="v-btn" v-html="icons.list"></button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.media-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 12px 24px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.divider {
  width: 1px;
  height: 24px;
  background: var(--c-border);
  margin: 0 4px;
}

.tools-group {
  display: flex;
  align-items: center;
  gap: 6px;
}

.t-icon {
  color: var(--c-text-muted);
  display: flex;
}

.t-select {
  border: none;
  background: transparent;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
  outline: none;
  cursor: pointer;
  appearance: none;
  padding-right: 16px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right center;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.t-count {
  font-size: 13px;
  color: var(--c-text-muted);
}

.view-icons {
  display: flex;
  gap: 4px;
}

.v-btn {
  background: transparent;
  border: none;
  color: var(--c-text-muted);
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}
.v-btn:hover { background: #f3f4f6; color: var(--c-text-main); }
.v-btn.active { color: var(--c-primary); }
</style>
