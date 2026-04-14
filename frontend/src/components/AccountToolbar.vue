<script setup lang="ts">
import { icons } from '../utils/icons'
</script>

<template>
  <div class="account-toolbar">
    <div class="t-left">
      <div class="search-box">
        <span class="s-icon" v-html="icons.search"></span>
        <input type="text" placeholder="Tìm kiếm tài khoản..." class="s-input" />
      </div>
    </div>
    
    <div class="t-right">
      <button class="btn btn-outline"><span class="b-icon" v-html="icons.refreshCcw"></span> Kiểm tra tất cả phiên</button>
      <button class="btn btn-primary" @click="$emit('add')"><span class="b-icon" v-html="icons.plus"></span> Thêm tài khoản</button>
    </div>
  </div>
</template>

<style scoped>
.account-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.t-left {
  flex: 1;
  max-width: 400px;
}

.search-box {
  display: flex;
  align-items: center;
  background: white;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  padding: 0 12px;
  height: 40px;
}

.s-icon {
  color: var(--c-text-muted);
  display: flex;
  margin-right: 8px;
}
.s-icon :deep(svg) { width: 16px; height: 16px; }

.s-input {
  border: none;
  background: transparent;
  outline: none;
  width: 100%;
  font-size: 13px;
  color: var(--c-text-main);
}
.s-input::placeholder { color: #9ca3af; }

.t-right {
  display: flex;
  gap: 12px;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 40px;
  padding: 0 16px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn-primary { background: var(--c-primary); color: white; }
.btn-primary:hover { background: var(--c-primary-hover); }

.btn-outline { background: white; border-color: var(--c-border); color: var(--c-text-title); }
.btn-outline:hover { background: #f9fafb; }

.b-icon { display: flex; }
.b-icon :deep(svg) { width: 14px; height: 14px; }
</style>
