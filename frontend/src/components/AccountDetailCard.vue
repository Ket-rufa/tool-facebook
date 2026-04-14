<script setup lang="ts">
import SessionStatusBadge from './SessionStatusBadge.vue'
import { icons } from '../utils/icons'

defineProps<{
  account: any | null
}>()

defineEmits(['setDefault', 'checkSession', 'remove'])
</script>

<template>
  <div class="account-detail-card">
    <template v-if="account">
      <div class="d-header">
        <img :src="account.avatar" class="d-avatar" alt="Avatar" />
        <div class="d-summary">
          <h2 class="d-name">{{ account.displayName }}</h2>
          <div class="d-meta-badges">
            <span class="meta-tag">{{ account.accountType }}</span>
            <span class="meta-tag">{{ account.provider }}</span>
            <span v-if="account.isDefault" class="meta-tag tag-primary">Mặc định</span>
          </div>
        </div>
      </div>
      
      <div class="d-body">
        <div class="info-row">
          <label>Account ID</label>
          <span class="val fw-code">{{ account.accountId }}</span>
        </div>
        
        <div class="info-row">
          <label>Trạng thái Phiên</label>
          <SessionStatusBadge :status="account.sessionStatus" />
        </div>
        
        <div class="info-row">
          <label>Kiểm tra lần cuối</label>
          <span class="val">{{ account.lastCheckedAt }}</span>
        </div>
        
        <div class="info-row col">
          <label>Ghi chú nội bộ</label>
          <div class="val-box">{{ account.note }}</div>
        </div>
      </div>
      
      <div class="d-footer">
        <button v-if="!account.isDefault" class="btn btn-outline" @click="$emit('setDefault', account.id)"><span class="icon" v-html="icons.checkCircle"></span> Đặt làm mặc định</button>
        <button class="btn btn-outline" @click="$emit('checkSession', account.id)"><span class="icon" v-html="icons.refreshCcw"></span> Kiểm tra lại phiên</button>
        <button class="btn btn-danger" @click="$emit('remove', account.id)"><span class="icon" v-html="icons.trash"></span> Gỡ khỏi danh sách</button>
      </div>
    </template>
    
    <div v-else class="empty-state">
      Chọn một tài khoản từ danh sách bên trái để xem chi tiết.
    </div>
  </div>
</template>

<style scoped>
.account-detail-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  padding: 24px;
}

.empty-state {
  display: flex;
  height: 200px;
  align-items: center;
  justify-content: center;
  color: var(--c-text-muted);
  font-size: 13px;
  font-style: italic;
}

.d-header {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--c-border);
}

.d-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #eff6ff;
}

.d-summary {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.d-name {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
  color: var(--c-text-title);
}

.d-meta-badges {
  display: flex;
  gap: 8px;
}

.meta-tag {
  font-size: 11px;
  background: #f3f4f6;
  color: var(--c-text-main);
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 600;
}

.tag-primary {
  background: #dbeafe;
  color: var(--c-primary);
}

.d-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
}

.info-row.col {
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.info-row label {
  color: var(--c-text-muted);
  font-weight: 600;
  width: 140px;
}

.val {
  color: var(--c-text-main);
  font-weight: 500;
  text-align: right;
}

.fw-code {
  font-family: monospace;
  background: #f9fafb;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.val-box {
  width: 100%;
  padding: 12px;
  background: #f9fafb;
  border: 1px dashed var(--c-border);
  border-radius: var(--radius-sm);
  color: var(--c-text-main);
  line-height: 1.5;
}

.d-footer {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn-outline { background: white; border-color: var(--c-border); color: var(--c-text-title); }
.btn-outline:hover { background: #f9fafb; }

.btn-danger { background: white; border-color: #fecaca; color: #dc2626; }
.btn-danger:hover { background: #fef2f2; border-color: #f87171; }

.icon { display: flex; }
.icon :deep(svg) { width: 14px; height: 14px; }
</style>
