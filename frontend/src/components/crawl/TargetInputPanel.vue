<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  targets: any[]
}>()

const emit = defineEmits(['add', 'resolve', 'clear'])
const rawInput = ref('')

const handleAdd = () => {
  if (!rawInput.value.trim()) return
  emit('add', rawInput.value)
  rawInput.value = ''
}
</script>

<template>
  <div class="crawl-panel target-input-panel">
    <div class="panel-header">
      <h3>Zone A: Mục tiêu quét</h3>
      <span class="badge">{{ targets.length }} URLs</span>
    </div>
    
    <div class="panel-body">
      <textarea 
        class="target-textarea" 
        v-model="rawInput"
        placeholder="Dán danh sách Facebook profile/page URL vào đây (mỗi dòng 1 dạng)..."
      ></textarea>
      
      <div class="action-row">
        <button class="btn btn-primary" @click="handleAdd">Thêm mục tiêu</button>
        <button class="btn btn-outline" @click="emit('resolve')">Giải mã UID</button>
        <button class="btn btn-text text-danger" @click="emit('clear')">Xoá tất cả</button>
      </div>

      <div class="target-list">
        <div v-for="target in targets" :key="target.id" class="target-item">
          <div class="target-avatar">
            <img v-if="target.avatar" :src="target.avatar" alt="avatar" />
            <div v-else class="avatar-placeholder"></div>
          </div>
          <div class="target-info">
            <div class="target-name">
              {{ target.displayName || '—' }}
              <span class="type-badge">{{ target.type }}</span>
            </div>
            <div class="target-url">{{ target.url }}</div>
          </div>
          <div class="target-status">
            <span class="uid-badge" v-if="target.uid">UID: {{ target.uid }}</span>
            <span :class="['status-badge', target.resolveStatus.toLowerCase()]">
              {{ target.resolveStatus === 'Resolved' ? 'Đã giải mã' : target.resolveStatus === 'Failed' ? 'Thất bại' : 'Chờ xử lý' }}
            </span>
            <span :class="['access-badge', target.publicAccess.toLowerCase()]">
              {{ target.publicAccess === 'Accessible' ? 'Truy cập được' : target.publicAccess === 'Inaccessible' ? 'Bị chặn' : 'Không xác định' }}
            </span>
          </div>
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f9fafb;
}

.panel-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.badge {
  background: #e5e7eb;
  color: #4b5563;
  padding: 2px 8px;
  border-radius: 99px;
  font-size: 12px;
  font-weight: 500;
}

.panel-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex: 1;
}

.target-textarea {
  width: 100%;
  height: 80px;
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  resize: none;
  font-family: inherit;
  font-size: 13px;
}

.target-textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.action-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.btn {
  padding: 6px 12px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 6px;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn-primary { background: #2563eb; color: white; }
.btn-primary:hover { background: #1d4ed8; }
.btn-outline { background: white; border-color: #d1d5db; color: #374151; }
.btn-outline:hover { background: #f3f4f6; }
.btn-text { background: transparent; }
.btn-text.text-danger { color: #dc2626; }
.btn-text:hover { text-decoration: underline; }

.target-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
  max-height: 200px;
}

.target-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px;
  border: 1px solid #f3f4f6;
  border-radius: 6px;
  background: white;
}

.target-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
  background: #f3f4f6;
}

.target-avatar img { width: 100%; height: 100%; object-fit: cover; }
.avatar-placeholder { width: 100%; height: 100%; }

.target-info { flex: 1; min-width: 0; }
.target-name { 
  font-size: 13px; font-weight: 600; color: #111827; 
  display: flex; align-items: center; gap: 8px;
}
.type-badge { font-size: 10px; background: #f3f4f6; color: #6b7280; padding: 2px 6px; border-radius: 4px; }
.target-url { font-size: 12px; color: #6b7280; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.target-status {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.uid-badge { font-size: 11px; font-family: monospace; color: #374151; background: #e5e7eb; padding: 2px 6px; border-radius: 4px; }
.status-badge, .access-badge { font-size: 10px; font-weight: 500; }
.status-badge.resolved { color: #059669; }
.status-badge.failed { color: #dc2626; }
.status-badge.pending { color: #6b7280; }
.access-badge.accessible { color: #059669; }
.access-badge.partial { color: #d97706; }
.access-badge.inaccessible { color: #dc2626; }
.access-badge.unknown { color: #6b7280; }
</style>
