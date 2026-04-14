<script setup lang="ts">
import AccountListItem from './AccountListItem.vue'

defineProps<{
  accounts: any[]
  selectedId: string | null
}>()

const emit = defineEmits(['selectItem'])
</script>

<template>
  <div class="account-list-panel">
    <div class="panel-header">
      <h3 class="panel-title">Tài khoản đã lưu ({{ accounts.length }})</h3>
      <div class="panel-filter">
        <!-- Mock filter -->
        <select class="f-select">
          <option>Tất cả trạng thái</option>
          <option>Hoạt động</option>
          <option>Lỗi</option>
        </select>
      </div>
    </div>
    
    <div class="panel-body">
      <AccountListItem
        v-for="acc in accounts"
        :key="acc.id"
        :account="acc"
        :isSelected="acc.id === selectedId"
        @select="emit('selectItem', acc.id)"
      />
    </div>
  </div>
</template>

<style scoped>
.account-list-panel {
  display: flex;
  flex-direction: column;
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  height: calc(100vh - 200px);
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #f9fafb;
  border-bottom: 1px solid var(--c-border);
}

.panel-title {
  margin: 0;
  font-size: 13px;
  font-weight: 800;
  text-transform: uppercase;
  color: var(--c-text-title);
  letter-spacing: 0.5px;
}

.f-select {
  border: 1px solid var(--c-border);
  background: white;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 12px;
  color: var(--c-text-main);
  outline: none;
}

.panel-body {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
