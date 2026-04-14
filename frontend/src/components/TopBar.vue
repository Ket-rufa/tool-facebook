<script setup lang="ts">
import { icons } from '../utils/icons'
</script>

<template>
  <header class="topbar">
    <div class="search-box">
      <span class="search-icon" v-html="icons.search"></span>
      <input type="text" placeholder="Tìm kiếm hồ sơ hoặc công việc..." />
    </div>

    <div class="actions">
      <button class="btn-create">
        <span class="icon" v-html="icons.plus"></span>
        Tạo công việc quét mới
      </button>
      
      <div class="icon-btn" v-html="icons.bell"></div>
      <div class="icon-btn" v-html="icons.wifi"></div>
    </div>
  </header>
</template>

<style scoped>
.topbar {
  height: 72px;
  background-color: var(--c-surface);
  border-bottom: 1px solid var(--c-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
}

.search-box {
  display: flex;
  align-items: center;
  background-color: #f3f4f6;
  border-radius: var(--radius-md);
  padding: 0 16px;
  width: 360px;
  height: 40px;
}

.search-icon {
  color: var(--c-text-muted);
  display: flex;
  align-items: center;
  margin-right: 8px;
}

.search-box input {
  border: none;
  background: transparent;
  width: 100%;
  height: 100%;
  font-family: inherit;
  font-size: 14px;
  color: var(--c-text-main);
  outline: none;
}

.search-box input::placeholder {
  color: #9ca3af;
}

.actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.btn-create {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: var(--c-primary);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  padding: 10px 16px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  font-family: inherit;
}

.btn-create:hover {
  background-color: var(--c-primary-hover);
}

.icon {
  display: flex;
  align-items: center;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #4b5563;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-btn:hover {
  background-color: #f3f4f6;
  color: var(--c-text-main);
}
</style>
