<script setup lang="ts">
import { icons } from '../utils/icons'
import PostFilterChip from './PostFilterChip.vue'
</script>

<template>
  <div class="post-filter-container">
    <div class="filter-row-top">
      <div class="filter-group col-keyword">
        <label>TỪ KHÓA</label>
        <div class="input-wrap">
          <input type="text" placeholder="Nhập từ khóa..." />
        </div>
      </div>
      
      <div class="filter-group col-date">
        <label>KHOẢNG NGÀY</label>
        <select class="filter-select custom-arrow">
          <option>7 ngày qua</option>
        </select>
      </div>
      
      <div class="filter-group col-source">
        <label>NGUỒN</label>
        <select class="filter-select custom-arrow">
          <option>Tất cả nguồn</option>
        </select>
      </div>
      
      <div class="filter-actions">
        <button class="btn-add-filter">
          <span v-html="icons.plus"></span> Thêm bộ lọc
        </button>
        <a href="#" class="btn-reset">Đặt lại</a>
      </div>
    </div>
    
    <div class="filter-row-bottom">
      <PostFilterChip label="Chỉ bài có phương tiện" icon="image" />
      <PostFilterChip label="Tương tác cao (>1k)" icon="trendingUp" />
      <button class="btn-lib">+ THƯ VIỆN LỌC</button>
    </div>
  </div>
</template>

<style scoped>
.post-filter-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
}

.filter-row-top {
  display: flex;
  align-items: flex-end;
  gap: 16px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group label {
  font-size: 11px;
  font-weight: 700;
  color: var(--c-text-muted);
}

.col-keyword { width: 300px; }
.col-date { width: 160px; }
.col-source { width: 160px; }

.input-wrap {
  background: #f3f4f6;
  border-radius: var(--radius-sm);
  padding: 0 12px;
  height: 38px;
  display: flex;
  align-items: center;
}

.input-wrap input {
  border: none;
  background: transparent;
  width: 100%;
  font-size: 13px;
  outline: none;
}

.filter-select {
  background: #f3f4f6;
  border: none;
  border-radius: var(--radius-sm);
  height: 38px;
  padding: 0 12px;
  font-size: 13px;
  outline: none;
  cursor: pointer;
  appearance: none;
}

.custom-arrow {
  padding-right: 32px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
}

.filter-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-left: 8px;
  height: 38px;
}

.btn-add-filter {
  display: flex;
  align-items: center;
  gap: 6px;
  background: white;
  border: 1px solid var(--c-border);
  height: 38px;
  padding: 0 16px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
}
.btn-add-filter:hover { background: #f9fafb; }

.btn-reset {
  font-size: 13px;
  font-weight: 600;
  color: var(--c-primary);
  text-decoration: none;
}
.btn-reset:hover { text-decoration: underline; }

.filter-row-bottom {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid var(--c-border);
}

.btn-lib {
  background: transparent;
  border: none;
  font-size: 12px;
  font-weight: 700;
  color: var(--c-primary);
  cursor: pointer;
  padding: 4px 8px;
}
</style>
