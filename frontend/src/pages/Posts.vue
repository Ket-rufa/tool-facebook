<script setup lang="ts">
import { ref } from 'vue'
import { postsData } from '../data/postsMockData'
import ViewToggle from '../components/ViewToggle.vue'
import PostFilterBar from '../components/PostFilterBar.vue'
import BulkActionBar from '../components/BulkActionBar.vue'
import PostTable from '../components/PostTable.vue'

const selectedIds = ref<string[]>(['FB_982129']) // Giả lập chọn mặc định 1 dòng

function handleToggleSelect(id: string) {
  const index = selectedIds.value.indexOf(id)
  if (index > -1) {
    selectedIds.value.splice(index, 1)
  } else {
    selectedIds.value.push(id)
  }
}
</script>

<template>
  <div class="posts-page">
    <div class="page-header">
      <div class="title-area">
        <h1>Quản lý bài viết</h1>
        <p>Theo dõi và phân tích {{ postsData.total }} bài viết đã thu thập từ các nguồn.</p>
      </div>
      <div class="header-actions">
        <ViewToggle />
      </div>
    </div>

    <div class="page-content">
      <PostFilterBar class="no-bottom-radius" />
      
      <BulkActionBar 
        :selected-count="selectedIds.length"
        :total-count="postsData.total"
        @export="() => {}"
        @delete="() => {}"
      />
      
      <PostTable 
        :posts="postsData.list"
        :selected-ids="selectedIds"
        @toggle-select="handleToggleSelect"
        :style="{ borderTop: selectedIds.length ? 'none' : '' }"
      />

      <div class="pagination-footer">
        <div class="page-size">
          Số dòng mỗi trang: <strong>15</strong>
        </div>
        <div class="page-controls">
          <button class="page-btn"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="11 17 6 12 11 7"></polyline><polyline points="18 17 13 12 18 7"></polyline></svg></button>
          <button class="page-btn"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"></polyline></svg></button>
          
          <button class="page-btn active">1</button>
          <button class="page-btn">2</button>
          <button class="page-btn">3</button>
          <span class="page-dots">...</span>
          <button class="page-btn">856</button>

          <button class="page-btn"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"></polyline></svg></button>
          <button class="page-btn"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="13 17 18 12 13 7"></polyline><polyline points="6 17 11 12 6 7"></polyline></svg></button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.posts-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.title-area h1 {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 700;
  color: var(--c-text-title);
}

.title-area p {
  margin: 0;
  color: var(--c-text-muted);
  font-size: 14px;
}

.page-content {
  display: flex;
  flex-direction: column;
}

/* Connect the filter bar to the table seamlessly */
:deep(.no-bottom-radius) {
  border-bottom-left-radius: 0;
  border-bottom-right-radius: 0;
  border-bottom: none;
}

.pagination-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: #f9fafb;
  border: 1px solid var(--c-border);
  border-top: none;
  border-radius: 0 0 var(--radius-md) var(--radius-md);
}

.page-size {
  font-size: 13px;
  color: var(--c-text-muted);
}
.page-size strong {
  color: var(--c-text-main);
}

.page-controls {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  color: var(--c-text-main);
  cursor: pointer;
  transition: all 0.2s;
}

.page-btn:hover { background: #e5e7eb; }
.page-btn.active {
  background: var(--c-primary);
  color: white;
  font-weight: 700;
}

.page-dots {
  font-size: 14px;
  color: var(--c-text-muted);
  padding: 0 4px;
}
</style>
