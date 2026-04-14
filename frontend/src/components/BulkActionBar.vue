<script setup lang="ts">
defineProps<{
  selectedCount: number
  totalCount: string
}>()

const emit = defineEmits(['export', 'delete'])
</script>

<template>
  <div class="bulk-action-bar" v-if="selectedCount > 0">
    <div class="left-actions">
      <span class="sel-text">Đã chọn: <span class="text-primary font-bold">{{ selectedCount }} bài viết</span></span>
      <button class="btn-bulk btn-export" @click="emit('export')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>
        Xuất CSV
      </button>
      <button class="btn-bulk btn-delete" @click="emit('delete')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
        Xóa đã chọn
      </button>
    </div>
    <div class="right-text">
      Hiển thị 1-15 của {{ totalCount }}
    </div>
  </div>
</template>

<style scoped>
.bulk-action-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-bottom: none;
  padding: 12px 24px;
}

.left-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.sel-text {
  font-size: 13px;
  color: var(--c-text-main);
}

.font-bold {
  font-weight: 700;
}

.btn-bulk {
  display: flex;
  align-items: center;
  gap: 8px;
  border: none;
  border-radius: 4px;
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.btn-export {
  background: #f3f4f6;
  color: var(--c-text-title);
}
.btn-export:hover { background: #e5e7eb; }

.btn-delete {
  background: #fee2e2;
  color: #dc2626;
}
.btn-delete:hover { background: #fecaca; }

.right-text {
  font-size: 13px;
  color: var(--c-text-muted);
}
</style>
