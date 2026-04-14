<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { accountsData as data } from '../data/accountsMockData'
import AccountToolbar from '../components/AccountToolbar.vue'
import AccountListPanel from '../components/AccountListPanel.vue'
import AccountDetailCard from '../components/AccountDetailCard.vue'
import AddAccountModal from '../components/AddAccountModal.vue'
import { icons } from '../utils/icons'

// Import Wails Backend bindings - AccountService
import {
  ListAccounts,
  CheckSessionStatus,
  RemoveAccount,
  SetDefaultAccount,
  AddAccountProfile,
  StartAccountAttachFlow,
  GetAttachFlowStatus,
  CompleteAccountAttachFlow,
  CancelAccountAttachFlow,
  CheckAllSessionStatuses
} from '../../wailsjs/go/accounts/AccountService'
import type { accounts } from '../../wailsjs/go/models'

const accountsList = ref<accounts.AccountProfile[]>([])
const isModalOpen = ref(false)
const selectedId = ref<string | null>(null)

const selectedAccount = computed(() => {
  return accountsList.value.find((a: any) => a.id === selectedId.value) || null
})

async function fetchAccounts() {
  const res = await ListAccounts()
  accountsList.value = res || []
  if (accountsList.value.length > 0 && !selectedId.value) {
    selectedId.value = accountsList.value[0].id
  } else if (accountsList.value.length === 0) {
    selectedId.value = null
  }
}

onMounted(() => {
  fetchAccounts()
})

function handleSelect(id: string) {
  selectedId.value = id
}

async function handleAddAccountComplete(updatedList: accounts.AccountProfile[]) {
  // Backend đã lưu rồi, chỉ cần cập nhật UI từ danh sách trả về
  if (updatedList && updatedList.length > 0) {
    accountsList.value = updatedList
    selectedId.value = updatedList[0].id
  } else {
    await fetchAccounts()
  }
  isModalOpen.value = false
}

async function handleCheckSession(id: string) {
  const res = await CheckSessionStatus(id)
  accountsList.value = res
}

async function handleSetDefault(id: string) {
  const res = await SetDefaultAccount(id)
  accountsList.value = res
}

async function handleRemove(id: string) {
  const res = await RemoveAccount(id)
  accountsList.value = res
  if (selectedId.value === id) {
    selectedId.value = res.length > 0 ? res[0].id : null
  }
}
</script>

<template>
  <div class="accounts-page">
    <div class="page-header">
      <h1 class="page-title">{{ data.header.title }}</h1>
      <p class="page-desc">{{ data.header.description }}</p>
    </div>

    <AccountToolbar @add="isModalOpen = true" />

    <div class="page-layout">
      <!-- Left Column: List -->
      <div class="layout-left">
        <AccountListPanel 
          :accounts="accountsList" 
          :selectedId="selectedId" 
          @selectItem="handleSelect" 
        />
      </div>

      <!-- Right Column: Details + Instructions -->
      <div class="layout-right">
        <AccountDetailCard 
          :account="selectedAccount" 
          @setDefault="handleSetDefault"
          @checkSession="handleCheckSession"
          @remove="handleRemove"
        />
        
        <!-- Instruction Card -->
        <div class="instruction-card">
          <div class="i-header">
            <span class="i-icon text-primary" v-html="icons.info"></span>
            Hướng dẫn sử dụng
          </div>
          <ul class="i-list">
            <li>Thêm tài khoản lưu trữ phiên qua cửa sổ đăng nhập riêng biệt.</li>
            <li>Hệ thống <strong>không lưu mật khẩu</strong> trực tiếp trong danh sách này.</li>
            <li>Chỉ lưu và quản lý Phiên (Session/Cookie) đã được xác thực thành công.</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Modal Workflow Thêm Account -->
    <AddAccountModal 
      v-if="isModalOpen" 
      @close="isModalOpen = false" 
      @complete="handleAddAccountComplete" 
    />
  </div>
</template>

<style scoped>
.accounts-page {
  display: flex;
  flex-direction: column;
  max-width: 1200px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 800;
  color: var(--c-text-title);
}

.page-desc {
  margin: 0;
  font-size: 15px;
  color: var(--c-text-muted);
}

.page-layout {
  display: grid;
  grid-template-columns: 400px 1fr;
  gap: 32px;
}

.layout-right {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.instruction-card {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: var(--radius-md);
  padding: 20px;
}

.i-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 700;
  color: #1e3a8a;
  margin-bottom: 12px;
}

.i-icon { display: flex; }
.i-icon :deep(svg) { width: 16px; height: 16px; }

.i-list {
  margin: 0;
  padding-left: 20px;
  font-size: 13px;
  color: #1e40af;
  line-height: 1.6;
}

.i-list li {
  margin-bottom: 4px;
}
</style>
