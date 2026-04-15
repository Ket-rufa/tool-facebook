<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { ActionRequest } from '../../types/action'
import { icons } from '../../utils/icons'
import { ListAccounts } from '../../../wailsjs/go/accounts/AccountService'
import { GetReactDocID, UpdateReactDocID } from '../../../wailsjs/go/actiontest/ActionHandler'
import type { accounts } from '../../../wailsjs/go/models'

const props = defineProps<{
  isLoading: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', payloads: ActionRequest[]): void
  (e: 'updateAccount', accountId: string): void
}>()

// ── State ─────────────────────────────────────────────────────────────
const accountsList = ref<accounts.AccountProfile[]>([])
const accountsLoading = ref(false)

const reactions = [
  { value: 'like', label: 'Like', icon: '👍', fb_id: '1635855486666999' },
  { value: 'love', label: 'Tim', icon: '❤️', fb_id: '1678524932434102' },
  { value: 'care', label: 'Thương thương', icon: '🥰', fb_id: '613557422527858' },
  { value: 'haha', label: 'Haha', icon: '😆', fb_id: '115940658764963' },
  { value: 'wow', label: 'Wow', icon: '😮', fb_id: '478547315650144' },
  { value: 'sad', label: 'Buồn', icon: '😢', fb_id: '908563459236466' },
  { value: 'angry', label: 'Phẫn nộ', icon: '😡', fb_id: '444813342392137' }
]

const selectedAccounts = ref<string[]>([])
const postId = ref('')
const selectedReaction = ref('like')
const dryRun = ref(true)
const validationError = ref('')
const reactDocID = ref('')
const docIdSaving = ref(false)
const docIdError = ref('')
const docIdNotice = ref('')
// doc_id đã bị xóa khỏi UI — backend tự resolve từ graphqlDocID constant nội bộ

// ── Load accounts thật từ backend ─────────────────────────────────────
async function loadAccounts() {
  accountsLoading.value = true
  try {
    const list = await ListAccounts()
    accountsList.value = list || []
    // Auto-select default account hoặc account đầu tiên
    if (accountsList.value.length > 0) {
      const def = accountsList.value.find(a => a.isDefault)
      if (def) {
        selectedAccounts.value = [def.id]
      } else {
        selectedAccounts.value = [accountsList.value[0].id]
      }
    }
  } catch (e) {
    console.error('Không thể tải danh sách tài khoản:', e)
    accountsList.value = []
  } finally {
    accountsLoading.value = false
  }
}

const isNumericDocID = (value: string) => /^\d+$/.test(value)

async function loadReactDocID() {
  try {
    const v = await GetReactDocID()
    reactDocID.value = (v || '').trim()
  } catch (e) {
    console.error('Khong the tai react doc_id:', e)
  }
}

async function persistReactDocID(showNotice = true): Promise<boolean> {
  docIdError.value = ''
  const value = reactDocID.value.trim()

  if (value && !isNumericDocID(value)) {
    docIdNotice.value = ''
    docIdError.value = 'Doc ID chỉ được chứa chữ số.'
    return false
  }

  try {
    docIdSaving.value = true
    const saved = await UpdateReactDocID(value)
    reactDocID.value = (saved || '').trim()
    if (showNotice) {
      docIdNotice.value = reactDocID.value ? 'Đã lưu doc_id cục bộ.' : 'Đã xóa doc_id cục bộ.'
    }
    return true
  } catch (e: any) {
    docIdNotice.value = ''
    docIdError.value = 'Không lưu được doc_id: ' + (e?.message || e?.toString() || 'Unknown error')
    return false
  } finally {
    docIdSaving.value = false
  }
}

const handleSaveDocID = async () => {
  await persistReactDocID(true)
}

const handleClearDocID = async () => {
  reactDocID.value = ''
  await persistReactDocID(true)
}

onMounted(() => {
  loadAccounts()
  loadReactDocID()
})

watch(selectedAccounts, (newVal) => {
  if (newVal && newVal.length > 0) {
    emit('updateAccount', newVal[newVal.length - 1]) // Lấy session status của user click cuối cùng
  }
}, { deep: true, immediate: true })

// ── Computed ──────────────────────────────────────────────────────────
const hasAccounts = computed(() => accountsList.value.length > 0)

const isAccountSelectable = (acc: accounts.AccountProfile) => {
  return acc.sessionStatus === 'active' || acc.sessionStatus === 'unchecked'
}

const isReady = computed(() => {
  return selectedAccounts.value.length > 0 && postId.value.trim() && selectedReaction.value
})

const getAccountName = computed(() => {
  if (selectedAccounts.value.length === 0) return 'Chưa chọn'
  if (selectedAccounts.value.length === 1) {
    const acc = accountsList.value.find(a => a.id === selectedAccounts.value[0])
    return acc ? acc.displayName : 'Chưa chọn'
  }
  return `${selectedAccounts.value.length} tài khoản đã chọn`
})

const getReactionLabel = computed(() => {
  const r = reactions.find(rx => rx.value === selectedReaction.value)
  return r ? `${r.icon} ${r.label}` : 'Chưa chọn'
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
const handleSubmit = async () => {
  validationError.value = ''

  if (!hasAccounts.value) {
    validationError.value = 'Chưa có tài khoản nào. Vui lòng thêm tài khoản ở màn "Tài khoản" trước.'
    return
  }
  if (selectedAccounts.value.length === 0) {
    validationError.value = 'Vui lòng chọn ít nhất một tài khoản thực thi'
    return
  }
  if (!postId.value.trim()) {
    validationError.value = 'Vui lòng nhập Feedback ID (base64)'
    return
  }
	if (!selectedReaction.value) {
		validationError.value = 'Vui lòng chọn loại cảm xúc'
		return
	}

  if (!dryRun.value && !reactDocID.value.trim()) {
    validationError.value = 'Vui lòng nhập GraphQL doc_id trước khi chạy Real Run.'
    return
  }

	// 1) Validation Semantic ID 
	if (!dryRun.value && (!postId.value.trim().startsWith("ZmVlZ") && !postId.value.trim().includes("=="))) {
		validationError.value = 'ID không đúng Semantic. Meta GraphQL yêu cầu Feedback ID dạng Base64 (feedback:<id> -> ZmVlZGJhY2s6...). Không truyền pfbid hoặc ID bài viết tại đây.'
		return
	}

  if (!dryRun.value) {
    const ok = await persistReactDocID(false)
    if (!ok) {
      validationError.value = docIdError.value || 'Không lưu được GraphQL doc_id.'
      return
    }
  }

	const selectedReactionObj = reactions.find(r => r.value === selectedReaction.value)

  // doc_id KHÔNG được gửi từ frontend — backend tự resolve
  console.debug('[DocID-Frontend] Submitting payload WITHOUT doc_id — backend resolves internally')
  const payloads: ActionRequest[] = selectedAccounts.value.map((id) => ({
    post_id: postId.value.trim(),
    account_id: id,
    reaction_type: selectedReaction.value,
    reaction_id: selectedReactionObj ? selectedReactionObj.fb_id : '',
    dry_run: dryRun.value,
    actor_source: 'manual_test'
  }))

  emit('submit', payloads)
}
</script>

<template>
  <div class="test-form-wrapper">
    
    <!-- 1. Card Tài khoản -->
    <div class="card">
      <h2 class="card-title">Tài khoản thực thi</h2>
      
      <!-- Loading state -->
      <div v-if="accountsLoading" class="accounts-loading">
        <span class="spinner-sm"></span> Đang tải danh sách tài khoản...
      </div>
      
      <!-- Empty state -->
      <div v-else-if="!hasAccounts" class="empty-accounts">
        <span v-html="icons.users" class="ea-icon"></span>
        <p>Chưa có tài khoản nào được gắn phiên.</p>
        <p class="ea-hint">Vui lòng thêm tài khoản ở màn <strong>Tài khoản</strong> trước khi thử nghiệm.</p>
      </div>
      
      <!-- Account list từ backend -->
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
          <img :src="acc.avatar || `https://ui-avatars.com/api/?name=${encodeURIComponent(acc.displayName)}&background=eff6ff&color=0f62fe`" class="acc-avatar" alt="avatar"/>
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

    <!-- 2. Card Cấu hình Cảm xúc -->
    <div class="card">
      <h2 class="card-title">Cấu hình cảm xúc</h2>
      
      <div class="form-group">
        <label class="form-label">Feedback ID (Base64)</label>
        <input
          v-model="postId"
          type="text"
          class="form-input"
          placeholder="VD: ZmVlZGJhY2s6MTUw..."
          :disabled="isLoading"
        />
      </div>

      <div class="form-group">
        <label class="form-label">GraphQL Doc ID (Real Run)</label>
        <div class="docid-row">
          <input
            v-model="reactDocID"
            type="text"
            class="form-input"
            placeholder="VD: 26477330531933156"
            :disabled="isLoading || docIdSaving"
          />
          <button
            type="button"
            class="btn btn-outline btn-docid"
            :disabled="isLoading || docIdSaving"
            @click="handleSaveDocID"
          >
            {{ docIdSaving ? 'Đang lưu...' : 'Lưu' }}
          </button>
          <button
            type="button"
            class="btn btn-link btn-docid"
            :disabled="isLoading || docIdSaving"
            @click="handleClearDocID"
          >
            Xóa
          </button>
        </div>
        <div class="docid-hint">Lưu cục bộ để Real Run dùng tự động, không cần nhập terminal mỗi lần.</div>
        <div v-if="docIdError" class="docid-error">{{ docIdError }}</div>
        <div v-else-if="docIdNotice" class="docid-note">{{ docIdNotice }}</div>
      </div>

      <div class="form-group mt-4">
        <label class="form-label">Loại cảm xúc</label>
        <div class="reaction-grid">
          <label v-for="r in reactions" :key="r.value" class="reaction-btn" :class="{'is-active': selectedReaction === r.value}">
            <input type="radio" :value="r.value" v-model="selectedReaction" :disabled="isLoading" class="sr-only" />
            <span class="r-icon">{{ r.icon }}</span>
            <span class="r-label">{{ r.label }}</span>
          </label>
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
        <p v-if="dryRun" class="text-muted"><span v-html="icons.info" class="inline-icon"></span> Mô phỏng hành động, không gửi request thật lên hệ thống đích.</p>
        <p v-else class="text-warning"><span v-html="icons.alertCircle" class="inline-icon"></span> Sẽ gửi request thật đẩy tương tác lên Facebook nếu hệ thống hỗ trợ.</p>
      </div>
    </div>

    <!-- 4. Tóm tắt & Submit -->
    <div class="card summary-card">
      <h3 class="summary-title">Tóm tắt hành động</h3>
      <div class="summary-grid">
        <div class="sum-row"><span>Tài khoản:</span> <strong>{{ getAccountName }}</strong></div>
        <div class="sum-row"><span>Feedback ID:</span> <strong>{{ postId || '—' }}</strong></div>
        <div class="sum-row"><span>Cảm xúc:</span> <strong>{{ getReactionLabel }}</strong></div>
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
          {{ isLoading ? 'Đang gửi...' : 'Chạy tương tác' }}
        </button>
      </div>
    </div>

  </div>
</template>

<style scoped>
.test-form-wrapper {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 20px 24px;
  box-shadow: var(--shadow-sm);
}

.card-title {
  margin: 0 0 16px;
  font-size: 15px;
  font-weight: 700;
  color: var(--c-text-title);
}

.m-0 { margin: 0; }

.flex-heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* Account List */
.account-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.account-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.account-item:not(.is-disabled):hover {
  background: #f9fafb;
}

.account-item.is-selected {
  border-color: var(--c-primary);
  background: #eff6ff;
}

.account-item.is-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.sr-only {
  position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border-width: 0;
}

.acc-avatar {
  width: 40px; height: 40px; border-radius: 50%;
}

.acc-info {
  flex: 1;
}

.acc-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--c-text-title);
}

.acc-type {
  font-size: 12px;
  color: var(--c-text-muted);
}

.acc-status {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 12px;
}
.st-active { background: #d1fae5; color: #065f46; }
.st-expired { background: #fee2e2; color: #991b1b; }

.acc-check { color: var(--c-primary); display: flex;}
.acc-check :deep(svg) { width: 20px; height: 20px; }


/* Form Base */
.form-group {
  margin-bottom: 16px;
}

.form-label {
  display: block; font-size: 13px; font-weight: 600; color: var(--c-text-main); margin-bottom: 8px;
}

.form-input {
  width: 100%; padding: 10px 12px; border: 1px solid var(--c-border); border-radius: var(--radius-sm);
  font-size: 14px; color: var(--c-text-title); background: var(--c-bg); outline: none; font-family: monospace; transition: border-color 0.2s;
}
.form-input:focus { border-color: var(--c-primary); background: var(--c-surface); }
.mt-4 { margin-top: 24px; }

.docid-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-docid {
  height: 40px;
  white-space: nowrap;
}

.docid-hint {
  margin-top: 8px;
  font-size: 12px;
  color: var(--c-text-muted);
}

.docid-error {
  margin-top: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--c-danger);
}

.docid-note {
  margin-top: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--c-success, #059669);
}


/* Reaction Grid */
.reaction-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
}

.reaction-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 12px 8px;
  border: 1px solid var(--c-border-light);
  background: var(--c-bg);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.reaction-btn:hover { background: #f3f4f6; border-color: #d1d5db; }

.reaction-btn.is-active {
  background: #eff6ff; border-color: #bfdbfe;
  box-shadow: 0 0 0 1px var(--c-primary);
}

.r-icon { font-size: 24px; }
.r-label { font-size: 12px; font-weight: 600; color: var(--c-text-main); }


/* Segmented Control */
.segmented-control {
  display: inline-flex; position: relative; background: var(--c-border-light); padding: 4px; border-radius: var(--radius-sm); width: max-content;
}

.segment-bg {
  position: absolute; top: 4px; left: 4px; bottom: 4px; width: calc(50% - 4px); background: var(--c-surface); border-radius: calc(var(--radius-sm) - 2px); box-shadow: var(--shadow-sm); transition: transform 0.2s ease;
}
.segment-bg.is-right { transform: translateX(100%); }

.segment-btn {
  position: relative; z-index: 1; padding: 6px 16px; min-width: 100px; border: none; background: transparent; font-size: 13px; font-weight: 600; cursor: pointer; border-radius: calc(var(--radius-sm) - 2px);
}
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

.summary-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  background: white;
  padding: 16px;
  border: 1px solid var(--c-border-light);
  border-radius: var(--radius-sm);
  margin-bottom: 24px;
}

.sum-row { display: flex; flex-direction: column; gap: 4px; font-size: 13px; }
.sum-row span { color: var(--c-text-muted); }
.sum-row strong { color: var(--c-text-title); font-weight: 600;}

.btn {
  padding: 10px 20px; font-size: 14px; font-weight: 600; border-radius: var(--radius-sm); cursor: pointer; border: none; display: inline-flex; align-items: center; justify-content: center; gap: 8px; transition: all 0.2s;
}
.btn-lg { padding: 14px 24px; font-size: 15px; }
.w-full { width: 100%; }

.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-primary { background: var(--c-primary); color: white; }
.btn-primary:not(:disabled):hover { background: var(--c-primary-hover); }

/* Spinner */
.spinner {
  display: inline-block; width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-radius: 50%; border-top-color: white; animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.error-msg { font-size: 13px; font-weight: 600; text-align: center; margin-bottom: 16px; }
.text-danger { color: var(--c-danger); }

/* Empty & loading accounts */
.accounts-loading {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: var(--c-text-muted);
  padding: 12px 0;
}

.empty-accounts {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 16px;
  text-align: center;
  gap: 4px;
}
.empty-accounts p { margin: 0; font-size: 13px; color: var(--c-text-main); }
.ea-hint { color: var(--c-text-muted) !important; font-size: 12px !important; margin-top: 4px; }
.ea-icon { color: var(--c-text-muted); margin-bottom: 8px; }
.ea-icon :deep(svg) { width: 40px; height: 40px; stroke-width: 1.5; }

/* Additional session status badges */
.st-error { background: #fee2e2; color: #991b1b; }
.st-unknown { background: #f3f4f6; color: #6b7280; }

/* Default tag inline */
.default-tag {
  font-size: 9px;
  font-weight: 800;
  color: var(--c-primary);
  border: 1px solid #bfdbfe;
  padding: 1px 5px;
  border-radius: 4px;
  text-transform: uppercase;
  vertical-align: middle;
  margin-left: 4px;
}

/* Spinner small */
.spinner-sm {
  display: inline-block;
  width: 14px; height: 14px;
  border: 2px solid #e5e7eb;
  border-top-color: var(--c-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
</style>
