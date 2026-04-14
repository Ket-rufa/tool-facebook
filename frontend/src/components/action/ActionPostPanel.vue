<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { CreatePostRequest } from '../../types/action'
import { icons } from '../../utils/icons'
import { ListAccounts } from '../../../wailsjs/go/accounts/AccountService'
import type { accounts } from '../../../wailsjs/go/models'

const props = defineProps<{
  isLoading: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', payloads: CreatePostRequest[]): void
  (e: 'updateAccount', accountId: string): void
}>()

// ── State ─────────────────────────────────────────────────────────────
const accountsList = ref<accounts.AccountProfile[]>([])
const accountsLoading = ref(false)

const selectedAccounts = ref<string[]>([])
const postList = ref('')
const dryRun = ref(true)
const validationError = ref('')

// ── Load accounts thật từ backend ─────────────────────────────────────
async function loadAccounts() {
  accountsLoading.value = true
  try {
    const list = await ListAccounts()
    accountsList.value = list || []
    if (accountsList.value.length > 0) {
      const def = accountsList.value.find(a => a.isDefault)
      if (def) selectedAccounts.value = [def.id]
      else selectedAccounts.value = [accountsList.value[0].id]
    }
  } catch (e) {
    console.error('Không thể tải danh sách tài khoản:', e)
    accountsList.value = []
  } finally {
    accountsLoading.value = false
  }
}

onMounted(() => loadAccounts())

watch(selectedAccounts, (newVal) => {
  if (newVal && newVal.length > 0) {
    emit('updateAccount', newVal[newVal.length - 1])
  }
}, { deep: true, immediate: true })

// ── Computed ──────────────────────────────────────────────────────────
const hasAccounts = computed(() => accountsList.value.length > 0)

const isAccountSelectable = (acc: accounts.AccountProfile) => {
  return acc.sessionStatus === 'active' || acc.sessionStatus === 'unchecked'
}

const parsedPosts = computed(() => {
  return postList.value
    .split('\n')
    .map(c => c.trim())
    .filter(c => c.length > 0)
})

const isReady = computed(() => {
  return selectedAccounts.value.length > 0 && parsedPosts.value.length > 0
})

const getAccountName = computed(() => {
  if (selectedAccounts.value.length === 0) return 'Chưa chọn'
  if (selectedAccounts.value.length === 1) {
    const acc = accountsList.value.find(a => a.id === selectedAccounts.value[0])
    return acc ? acc.displayName : 'Chưa chọn'
  }
  return `${selectedAccounts.value.length} tài khoản đã chọn`
})

const sessionBadgeClass = (status: string) => {
  switch (status) {
    case 'active': return 'st-active'
    case 'expired': return 'st-expired'
    case 'error':
    case 'invalid': return 'st-error'
    default: return 'st-unknown'
  }
}

const sessionStatusLabel = (status: string) => {
  const map: Record<string, string> = {
    active: 'Hoạt động', expired: 'Hết hạn',
    error: 'Lỗi', invalid: 'Không hợp lệ', unchecked: 'Chưa kiểm tra'
  }
  return map[status] || status
}

// ── Submit ────────────────────────────────────────────────────────────
const handleSubmit = () => {
  validationError.value = ''

  if (!hasAccounts.value) {
    validationError.value = 'Chưa có tài khoản nào. Vui lòng thêm tài khoản ở màn "Tài khoản" trước.'
    return
  }
  if (selectedAccounts.value.length === 0) {
    validationError.value = 'Vui lòng chọn ít nhất một tài khoản thực thi'
    return
  }
  if (parsedPosts.value.length === 0) {
    validationError.value = 'Vui lòng nhập ít nhất một nội dung bài viết'
    return
  }

  const payloads: CreatePostRequest[] = selectedAccounts.value.map((id) => {
    const randomIndex = Math.floor(Math.random() * parsedPosts.value.length)
    const randomContent = parsedPosts.value[randomIndex]

    return {
      account_id: id,
      post_text: randomContent,
      dry_run: dryRun.value,
      actor_source: 'manual_test'
    }
  })

  emit('submit', payloads)
}
</script>

<template>
  <div class="test-form-wrapper">
    
    <!-- 1. Card Tài khoản -->
    <div class="card">
      <h2 class="card-title">Tài khoản thực thi</h2>
      
      <div v-if="accountsLoading" class="accounts-loading">
        <span class="spinner-sm"></span> Đang tải danh sách tài khoản...
      </div>
      
      <div v-else-if="!hasAccounts" class="empty-accounts">
        <span v-html="icons.users" class="ea-icon"></span>
        <p>Chưa có tài khoản nào được gắn phiên.</p>
        <p class="ea-hint">Vui lòng thêm tài khoản ở màn <strong>Tài khoản</strong></p>
      </div>
      
      <div v-else class="account-list">
        <label 
          v-for="acc in accountsList" 
          :key="acc.id" 
          class="account-item" 
          :class="{
            'is-selected': selectedAccounts.includes(acc.id),
            'is-disabled': !isAccountSelectable(acc)
          }"
        >
          <input 
            type="checkbox" :value="acc.id" v-model="selectedAccounts" 
            :disabled="!isAccountSelectable(acc) || isLoading" 
            class="sr-only" 
          />
          <img :src="`https://ui-avatars.com/api/?name=${encodeURIComponent(acc.displayName)}&background=eff6ff&color=0f62fe`" class="acc-avatar" alt="avatar"/>
          <div class="acc-info">
            <div class="acc-name">
              {{ acc.displayName }}
              <span v-if="acc.isDefault" class="default-tag">Mặc định</span>
            </div>
            <div class="acc-type">{{ acc.accountType }} · {{ acc.provider }}</div>
          </div>
          <div class="acc-status" :class="sessionBadgeClass(acc.sessionStatus)">
            {{ sessionStatusLabel(acc.sessionStatus) }}
          </div>
          <div class="acc-check" v-if="selectedAccounts.includes(acc.id)" v-html="icons.checkCircle"></div>
        </label>
      </div>
    </div>

    <!-- 2. Card Cấu hình Bài viết -->
    <div class="card">
      <h2 class="card-title">Nội dung bài viết (Status)</h2>
      
      <div class="form-group mt-2">
        <p class="text-muted text-xs mb-2">Nhập mỗi Status trên 1 dòng. Tool sẽ tự động bốc ngẫu nhiên (random) 1 nội dung cho mỗi tài khoản.</p>
        <textarea
          v-model="postList"
          class="form-input"
          rows="5"
          placeholder="Hôm nay trời đẹp quá!&#10;Mình vừa ra mắt sản phẩm mới.&#10;Trưa nay ăn gì nhỉ?"
          :disabled="isLoading"
        ></textarea>
        <div class="comment-stats" v-if="parsedPosts.length > 0">
          Đã nhận diện <strong>{{ parsedPosts.length }}</strong> kịch bản nội dung khác nhau.
        </div>
      </div>
      
      <div class="form-group mt-4">
        <label class="form-label text-muted">Đính kèm Hình ảnh / Video (Phát triển sau)</label>
        <div class="upload-placeholder">
           <span v-html="icons.image" class="inline-icon text-muted"></span> Tính năng up ảnh đang trong quá trình phát triển lõi Upload (MediaFBID).
        </div>
      </div>
    </div>

    <!-- 3. Card Chế độ -->
    <div class="card">
      <div class="flex-heading">
        <h2 class="card-title m-0">Chế độ chạy</h2>
        <div class="segmented-control">
          <div class="segment-bg" :class="{ 'is-right': !dryRun }"></div>
          <button type="button" class="segment-btn" :class="{ 'active text-primary': dryRun, 'text-muted': !dryRun }" @click="dryRun = true" :disabled="isLoading">
            Dry Run
          </button>
          <button type="button" class="segment-btn" :class="{ 'active font-bold text-title': !dryRun, 'text-muted': dryRun }" @click="dryRun = false" :disabled="isLoading">
            Real Run
          </button>
        </div>
      </div>
      <div class="mode-info">
        <p v-if="dryRun" class="text-muted"><span v-html="icons.info" class="inline-icon"></span> Mô phỏng hành động, không đăng thật lên Profile.</p>
        <p v-else class="text-warning"><span v-html="icons.alertCircle" class="inline-icon"></span> Đăng thẳng lên tường Profile (yêu cầu GraphQL doc_id).</p>
      </div>
    </div>

    <!-- 4. Tóm tắt & Submit -->
    <div class="card summary-card">
      <h3 class="summary-title">Tóm tắt hành động</h3>
      <div class="summary-grid">
        <div class="sum-row"><span>Tài khoản:</span> <strong>{{ getAccountName }}</strong></div>
        <div class="sum-row"><span>Đăng vào:</span> <strong>Profile Cá Nhân</strong></div>
        <div class="sum-row"><span>Nhật ký Status:</span> <strong>{{ parsedPosts.length }} câu</strong></div>
        <div class="sum-row"><span>Chế độ:</span> <strong>{{ dryRun ? 'Dry Run' : 'Real Run' }}</strong></div>
      </div>
      
      <div v-if="validationError" class="error-msg text-danger">{{ validationError }}</div>

      <div class="submit-area">
        <button
          type="button"
          @click="handleSubmit"
          :disabled="isLoading || !isReady"
          class="btn btn-primary btn-lg w-full"
        >
          <span class="spinner" v-if="isLoading"></span>
          {{ isLoading ? 'Đang gửi...' : 'Đăng Bài Hàng Loạt' }}
        </button>
      </div>
    </div>

  </div>
</template>

<style scoped>
.test-form-wrapper { display: flex; flex-direction: column; gap: 20px; }
.card { background: var(--c-surface); border: 1px solid var(--c-border); border-radius: var(--radius-md); padding: 20px 24px; box-shadow: var(--shadow-sm); }
.card-title { margin: 0 0 16px; font-size: 15px; font-weight: 700; color: var(--c-text-title); }
.m-0 { margin: 0; }
.flex-heading { display: flex; justify-content: space-between; align-items: center; }

/* Account List */
.account-list { display: flex; flex-direction: column; gap: 12px; }
.account-item { display: flex; align-items: center; gap: 16px; padding: 12px 16px; border: 1px solid var(--c-border); border-radius: var(--radius-sm); cursor: pointer; transition: all 0.2s; }
.account-item:not(.is-disabled):hover { background: #f9fafb; }
.account-item.is-selected { border-color: var(--c-primary); background: #eff6ff; }
.account-item.is-disabled { opacity: 0.5; cursor: not-allowed; }
.sr-only { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border-width: 0; }
.acc-avatar { width: 40px; height: 40px; border-radius: 50%; }
.acc-info { flex: 1; }
.acc-name { font-weight: 600; font-size: 14px; color: var(--c-text-title); }
.acc-type { font-size: 12px; color: var(--c-text-muted); }
.acc-status { font-size: 11px; font-weight: 600; padding: 2px 8px; border-radius: 12px; }
.st-active { background: #d1fae5; color: #065f46; }
.st-expired { background: #fee2e2; color: #991b1b; }
.acc-check { color: var(--c-primary); display: flex;}
.acc-check :deep(svg) { width: 20px; height: 20px; }

/* Form Base */
.form-group { margin-bottom: 16px; }
.form-label { display: block; font-size: 13px; font-weight: 600; color: var(--c-text-main); margin-bottom: 8px; }
.form-input { width: 100%; padding: 10px 12px; border: 1px solid var(--c-border); border-radius: var(--radius-sm); font-size: 14px; color: var(--c-text-title); background: var(--c-bg); outline: none; font-family: monospace; transition: border-color 0.2s; resize: vertical; }
.form-input:focus { border-color: var(--c-primary); background: var(--c-surface); }
.mt-4 { margin-top: 24px; }
.mb-2 { margin-bottom: 8px; }
.text-xs { font-size: 12px; }
.comment-stats { margin-top: 8px; font-size: 12px; color: var(--c-primary); }

.upload-placeholder { background: #f8fafc; border: 1px dashed var(--c-border-light); padding: 16px; border-radius: var(--radius-sm); text-align: center; font-size: 13px; color: var(--c-text-muted); }

/* Segmented Control */
.segmented-control { display: inline-flex; position: relative; background: var(--c-border-light); padding: 4px; border-radius: var(--radius-sm); width: max-content; }
.segment-bg { position: absolute; top: 4px; left: 4px; bottom: 4px; width: calc(50% - 4px); background: var(--c-surface); border-radius: calc(var(--radius-sm) - 2px); box-shadow: var(--shadow-sm); transition: transform 0.2s ease; }
.segment-bg.is-right { transform: translateX(100%); }
.segment-btn { position: relative; z-index: 1; padding: 6px 16px; min-width: 100px; border: none; background: transparent; font-size: 13px; font-weight: 600; cursor: pointer; border-radius: calc(var(--radius-sm) - 2px); }
.text-primary { color: var(--c-primary); }
.text-muted { color: var(--c-text-muted); font-weight: 500;}
.text-warning { color: var(--c-warning); font-weight: 600; }
.text-title { color: var(--c-text-title); font-weight: 700; }
.mode-info { margin-top: 16px; font-size: 13px; }
.inline-icon { display: inline-flex; vertical-align: middle; margin-top: -2px; }
.inline-icon :deep(svg) { width: 16px; height: 16px; }

/* Summary */
.summary-card { background: #fbfbfc; }
.summary-title { font-size: 13px; font-weight: 600; color: var(--c-text-muted); text-transform: uppercase; margin: 0 0 16px 0; letter-spacing: 0.5px; }
.summary-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; background: white; padding: 16px; border: 1px solid var(--c-border-light); border-radius: var(--radius-sm); margin-bottom: 24px; }
.sum-row { display: flex; flex-direction: column; gap: 4px; font-size: 13px; }
.sum-row span { color: var(--c-text-muted); }
.sum-row strong { color: var(--c-text-title); font-weight: 600;}

.btn { padding: 10px 20px; font-size: 14px; font-weight: 600; border-radius: var(--radius-sm); cursor: pointer; border: none; display: inline-flex; align-items: center; justify-content: center; gap: 8px; transition: all 0.2s; }
.btn-lg { padding: 14px 24px; font-size: 15px; }
.w-full { width: 100%; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-primary { background: var(--c-primary); color: white; }
.btn-primary:not(:disabled):hover { background: var(--c-primary-hover); }

/* Spinner */
.spinner { display: inline-block; width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-radius: 50%; border-top-color: white; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.error-msg { font-size: 13px; font-weight: 600; text-align: center; margin-bottom: 16px; }
.text-danger { color: var(--c-danger); }

/* Empty & loading accounts */
.accounts-loading { display: flex; align-items: center; gap: 10px; font-size: 13px; color: var(--c-text-muted); padding: 12px 0; }
.empty-accounts { display: flex; flex-direction: column; align-items: center; padding: 32px 16px; text-align: center; gap: 4px; }
.empty-accounts p { margin: 0; font-size: 13px; color: var(--c-text-main); }
.ea-hint { color: var(--c-text-muted) !important; font-size: 12px !important; margin-top: 4px; }
.ea-icon { color: var(--c-text-muted); margin-bottom: 8px; }
.ea-icon :deep(svg) { width: 40px; height: 40px; stroke-width: 1.5; }
.st-error { background: #fee2e2; color: #991b1b; }
.st-unknown { background: #f3f4f6; color: #6b7280; }
.default-tag { font-size: 9px; font-weight: 800; color: var(--c-primary); border: 1px solid #bfdbfe; padding: 1px 5px; border-radius: 4px; text-transform: uppercase; vertical-align: middle; margin-left: 4px; }
.spinner-sm { display: inline-block; width: 14px; height: 14px; border: 2px solid #e5e7eb; border-top-color: var(--c-primary); border-radius: 50%; animation: spin 0.8s linear infinite; }
</style>
