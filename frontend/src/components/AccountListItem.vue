<script setup lang="ts">
import SessionStatusBadge from './SessionStatusBadge.vue'

defineProps<{
  account: any
  isSelected: boolean
}>()

const emit = defineEmits(['select'])
</script>

<template>
  <div 
    class="account-list-item" 
    :class="{ selected: isSelected }"
    @click="emit('select')"
  >
    <img :src="account.avatar" class="a-avatar" alt="avatar" />
    
    <div class="a-main">
      <div class="a-name">
        {{ account.displayName }}
        <span v-if="account.isDefault" class="default-badge">Mặc định</span>
      </div>
      <div class="a-meta">
        Loại: <span class="fw-600">{{ account.accountType }}</span>
        <span class="sep">&bull;</span>
        Kiểm tra: {{ account.lastCheckedAt }}
      </div>
    </div>
    
    <div class="a-status">
      <SessionStatusBadge :status="account.sessionStatus" />
    </div>
  </div>
</template>

<style scoped>
.account-list-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background: white;
  border-radius: var(--radius-sm);
  border: 1px solid var(--c-border);
  cursor: pointer;
  transition: all 0.2s;
}

.account-list-item:hover {
  background: #f9fafb;
}

.account-list-item.selected {
  border-color: var(--c-primary);
  background: #eff6ff;
  box-shadow: 0 0 0 1px var(--c-primary);
}

.a-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 16px;
  border: 1px solid var(--c-border);
}

.a-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.a-name {
  font-size: 14px;
  font-weight: 700;
  color: var(--c-text-title);
  display: flex;
  align-items: center;
  gap: 8px;
}

.default-badge {
  font-size: 10px;
  font-weight: 700;
  color: var(--c-primary);
  background: #dbeafe;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
}

.a-meta {
  font-size: 12px;
  color: var(--c-text-muted);
}

.fw-600 { font-weight: 600; color: var(--c-text-main); }
.sep { padding: 0 6px; color: #d1d5db; }

.a-status {
  flex-shrink: 0;
}
</style>
