<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { CreatePostRequest } from '../../types/action'
import { icons } from '../../utils/icons'
import { ListAccounts } from '../../../wailsjs/go/accounts/AccountService'
import { PickMediaFiles } from '../../../wailsjs/go/main/App'
import { GetCreatePostDocID, UpdateCreatePostDocID } from '../../../wailsjs/go/actiontest/ActionHandler'
import type { accounts } from '../../../wailsjs/go/models'

type LocalMediaFile = {
  path: string
  name: string
}

const props = defineProps<{
  isLoading: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', payloads: CreatePostRequest[]): void
  (e: 'updateAccount', accountId: string): void
}>()

// State
const accountsList = ref<accounts.AccountProfile[]>([])
const accountsLoading = ref(false)
const selectedAccounts = ref<string[]>([])
const postList = ref('')
const mediaInput = ref('')
const mediaPlaceholders = ref<string[]>([])
const mediaFiles = ref<LocalMediaFile[]>([])
const dryRun = ref(true)
const validationError = ref('')

// Doc ID
const createPostDocID = ref('')
const docIDInput = ref('')
const docIDSaving = ref(false)
const docIDMessage = ref('')

// Load accounts
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
    console.error('Loi tai tai khoan:', e)
    accountsList.value = []
  } finally {
    accountsLoading.value = false
  }
}

// Load doc_id
async function loadCreatePostDocID() {
  try {
    const id = await GetCreatePostDocID()
    createPostDocID.value = id || ''
    docIDInput.value = id || ''
  } catch (e) {
    console.error('Loi tai doc_id:', e)
  }
}

// Save doc_id
async function saveCreatePostDocID() {
  docIDSaving.value = true
  docIDMessage.value = ''
  try {
    const saved = await UpdateCreatePostDocID(docIDInput.value.trim())
    createPostDocID.value = saved
    docIDMessage.value = saved ? 'Da luu doc_id: ' + saved : 'Da xoa doc_id - dung mac dinh'
  } catch (e: any) {
    docIDMessage.value = 'Loi: ' + (e?.message || String(e))
  } finally {
    docIDSaving.value = false
    setTimeout(() => { docIDMessage.value = '' }, 4000)
  }
}

onMounted(() => {
  loadAccounts()
  loadCreatePostDocID()
})

watch(selectedAccounts, (newVal) => {
  if (newVal && newVal.length > 0) {
    emit('updateAccount', newVal[newVal.length - 1])
  }
}, { deep: true, immediate: true })

// Media
const parsedMediaPlaceholders = computed(() => {
  return mediaInput.value
    .split('\n')
    .map(l => l.trim())
    .filter(l => l.startsWith('{') && l.endsWith('}'))
})

async function openMediaPicker() {
  try {
    const files: string[] = await PickMediaFiles()
    if (!files || files.length === 0) return
    const newFiles = files.map(p => {
      const parts = p.replace(/\\/g, '/').split('/')
      return { path: p, name: parts[parts.length - 1] }
    })
    mediaFiles.value = [...mediaFiles.value, ...newFiles]
  } catch (e) {
    console.error('Loi picker:', e)
  }
}

function removeMedia(index: number) {
  mediaFiles.value.splice(index, 1)
}

// Computed
const hasAccounts = computed(() => accountsList.value.length > 0)
const isAccountSelectable = (acc: accounts.AccountProfile) =>
  acc.sessionStatus === 'active' || acc.sessionStatus === 'unchecked'

const parsedPosts = computed(() =>
  postList.value.split('\n').map(c => c.trim()).filter(c => c.length > 0)
)

const isReady = computed(() => {
  const hasAccount = selectedAccounts.value.length > 0
  const hasContent = parsedPosts.value.length > 0 || mediaFiles.value.length > 0
  return hasAccount && hasContent
})

const getAccountName = computed(() => {
  if (selectedAccounts.value.length === 0) return 'Chua chon'
  if (selectedAccounts.value.length === 1) {
    const acc = accountsList.value.find(a => a.id === selectedAccounts.value[0])
    return acc ? acc.displayName : 'Chua chon'
  }
  return selectedAccounts.value.length + ' tai khoan da chon'
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
    active: 'Hoat dong', expired: 'Het han',
    error: 'Loi', invalid: 'Khong hop le', unchecked: 'Chua kiem tra'
  }
  return map[status] || status
}

// Submit
const handleSubmit = () => {
  validationError.value = ''
  if (!hasAccounts.value) {
    validationError.value = 'Chua co tai khoan nao. Vui long them tai khoan truoc.'
    return
  }
  if (selectedAccounts.value.length === 0) {
    validationError.value = 'Vui long chon it nhat mot tai khoan.'
    return
  }
  if (parsedPosts.value.length === 0 && mediaFiles.value.length === 0) {
    validationError.value = 'Vui long nhap noi dung hoac chon anh truoc khi dang.'
    return
  }

  const payloads: CreatePostRequest[] = selectedAccounts.value.map((id) => {
    let randomContent = ''
    if (parsedPosts.value.length > 0) {
      const idx = Math.floor(Math.random() * parsedPosts.value.length)
      randomContent = parsedPosts.value[idx]
    }
    const mediaList = mediaFiles.value.map(f => f.path)
    return {
      account_id: id,
      post_text: randomContent,
      image_paths: mediaList,
      dry_run: dryRun.value,
      actor_source: 'manual_test'
    }
  })

  emit('submit', payloads)
}
</script>

<template>
  <div class="test-form-wrapper">

    <!-- 1. Card Tai khoan -->
    <div class="card">
      <h2 class="card-title">Tai khoan thuc thi</h2>

      <div v-if="accountsLoading" class="accounts-loading">
        <span class="spinner-sm"></span> Dang tai danh sach tai khoan...
      </div>

      <div v-else-if="!hasAccounts" class="empty-accounts">
        <span v-html="icons.users" class="ea-icon"></span>
        <p>Chua co tai khoan nao duoc gan phien.</p>
        <p class="ea-hint">Vui long them tai khoan o man <strong>Tai khoan</strong></p>
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
              <span v-if="acc.isDefault" class="default-tag">Mac dinh</span>
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

    <!-- 2. Card Noi dung bai viet -->
    <div class="card">
      <h2 class="card-title">Noi dung bai viet (Status)</h2>

      <div class="form-group mt-2">
        <p class="text-muted text-xs mb-2">Nhap moi Status tren 1 dong. Tool se boc ngau nhien 1 noi dung cho moi tai khoan.</p>
        <textarea
          v-model="postList"
          class="form-input"
          rows="5"
          placeholder="Hom nay troi dep qua!&#10;Minh vua ra mat san pham moi.&#10;Trua nay an gi nhi?"
          :disabled="isLoading"
        ></textarea>
        <div class="comment-stats" v-if="parsedPosts.length > 0">
          Da nhan dien <strong>{{ parsedPosts.length }}</strong> kich ban noi dung khac nhau.
        </div>
      </div>
    </div>

    <!-- 3. Card Media -->
    <div class="card">
      <h2 class="card-title">Dinh kem Hinh anh / Video</h2>
      <p class="text-muted text-xs mb-2">Chon anh/video de dang kem bai viet. Tool se upload va dinh kem tu dong.</p>

      <div class="media-input-row">
        <button
          type="button"
          class="btn btn-outline"
          :disabled="isLoading"
          @click="openMediaPicker"
        >
          <span v-html="icons.image" class="inline-icon"></span>
          Chon anh/video tu may tinh
        </button>
      </div>

      <div v-if="mediaFiles.length > 0" class="media-file-list mt-2">
        <div v-for="(f, i) in mediaFiles" :key="i" class="media-file-item">
          <span v-html="icons.image" class="inline-icon file-icon"></span>
          <span class="file-name">{{ f.name }}</span>
          <button type="button" class="btn-remove" @click="removeMedia(i)" :disabled="isLoading">
            <span v-html="icons.x"></span>
          </button>
        </div>
        <div class="comment-stats">Tong file: <strong>{{ mediaFiles.length }}</strong></div>
      </div>
    </div>

    <!-- 4. Card CreatePost Doc ID -->
    <div class="card docid-card">
      <h2 class="card-title">⚙ CreatePost Doc ID</h2>
      <p class="text-muted text-xs mb-2">
        Neu gap loi 1357010 field_exception, Facebook da cap nhat schema.
        Vao F12 - Network - ComposerStoryCreateMutation - Request Payload de lay doc_id moi.
      </p>
      <div class="docid-current" :class="createPostDocID ? 'docid-set' : 'docid-default'">
        {{ createPostDocID ? 'Doc ID hien tai: ' + createPostDocID : 'Dang dung doc_id mac dinh (chua cau hinh)' }}
      </div>
      <div class="media-input-row mt-2">
        <input
          v-model="docIDInput"
          type="text"
          class="form-input"
          placeholder="VD: 27581837698072404"
          :disabled="isLoading || docIDSaving"
          @keydown.enter.prevent="saveCreatePostDocID"
        />
        <button type="button" class="btn btn-outline" :disabled="isLoading || docIDSaving" @click="saveCreatePostDocID">
          {{ docIDSaving ? 'Dang luu...' : 'Luu' }}
        </button>
        <button type="button" class="btn btn-outline btn-reset" title="Reset ve mac dinh"
          :disabled="isLoading || docIDSaving"
          @click="docIDInput = ''; saveCreatePostDocID()">
          X
        </button>
      </div>
      <div v-if="docIDMessage" class="docid-msg" :class="docIDMessage.startsWith('Loi') ? 'msg-error' : 'msg-ok'">
        {{ docIDMessage }}
      </div>
    </div>

    <!-- 5. Card Che do -->
    <div class="card">
      <div class="flex-heading">
        <h2 class="card-title m-0">Che do chay</h2>
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
        <p v-if="dryRun" class="text-muted"><span v-html="icons.info" class="inline-icon"></span> Mo phong hanh dong, khong dang that len Profile.</p>
        <p v-else class="text-warning"><span v-html="icons.alertCircle" class="inline-icon"></span> Dang thang len tuong Profile (yeu cau GraphQL doc_id).</p>
      </div>
    </div>

    <!-- 6. Tom tat & Submit -->
    <div class="card summary-card">
      <h3 class="summary-title">Tom tat hanh dong</h3>
      <div class="summary-grid">
        <div class="sum-row"><span>Tai khoan:</span> <strong>{{ getAccountName }}</strong></div>
        <div class="sum-row"><span>Dang vao:</span> <strong>Profile Ca Nhan</strong></div>
        <div class="sum-row"><span>Noi dung:</span> <strong>{{ parsedPosts.length }} cau</strong></div>
        <div class="sum-row"><span>Media:</span> <strong>{{ mediaFiles.length }} file</strong></div>
        <div class="sum-row"><span>Che do:</span> <strong>{{ dryRun ? 'Dry Run' : 'Real Run' }}</strong></div>
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
          {{ isLoading ? 'Dang gui...' : 'Dang Bai Hang Loat' }}
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
.acc-check { color: var(--c-primary); display: flex; }
.acc-check :deep(svg) { width: 20px; height: 20px; }

/* Form */
.form-group { margin-bottom: 16px; }
.form-label { display: block; font-size: 13px; font-weight: 600; color: var(--c-text-main); margin-bottom: 8px; }
.form-input { width: 100%; padding: 10px 12px; border: 1px solid var(--c-border); border-radius: var(--radius-sm); font-size: 14px; color: var(--c-text-title); background: var(--c-bg); outline: none; font-family: monospace; transition: border-color 0.2s; box-sizing: border-box; }
.form-input:focus { border-color: var(--c-primary); background: var(--c-surface); }
textarea.form-input { resize: vertical; }
.mt-2 { margin-top: 10px; }
.mt-4 { margin-top: 24px; }
.mb-2 { margin-bottom: 8px; }
.text-xs { font-size: 12px; }
.comment-stats { margin-top: 8px; font-size: 12px; color: var(--c-primary); }

/* Media */
.media-input-row { display: flex; gap: 8px; align-items: center; }
.media-file-list { display: flex; flex-direction: column; gap: 6px; }
.media-file-item { display: flex; align-items: center; gap: 8px; background: #f9fafb; border: 1px solid var(--c-border-light); border-radius: var(--radius-sm); padding: 6px 10px; }
.file-name { flex: 1; font-size: 12px; font-family: monospace; color: var(--c-text-main); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.file-icon :deep(svg) { width: 14px; height: 14px; color: var(--c-primary); }
.btn-remove { border: none; background: transparent; cursor: pointer; color: var(--c-text-muted); display: flex; padding: 2px; }
.btn-remove :deep(svg) { width: 14px; height: 14px; }
.btn-remove:hover { color: var(--c-danger); }

/* Doc ID Card */
.docid-card { border-left: 3px solid #f59e0b; }
.docid-current { font-size: 12px; font-family: monospace; padding: 8px 12px; border-radius: var(--radius-sm); margin-bottom: 4px; }
.docid-set { background: #eff6ff; color: var(--c-primary); border: 1px solid #bfdbfe; }
.docid-default { background: #f9fafb; color: var(--c-text-muted); border: 1px dashed var(--c-border); }
.docid-msg { margin-top: 8px; padding: 8px 12px; border-radius: var(--radius-sm); font-size: 12px; font-weight: 600; }
.msg-ok { background: #ecfdf5; color: #065f46; border: 1px solid #a7f3d0; }
.msg-error { background: #fee2e2; color: #991b1b; border: 1px solid #fca5a5; }
.btn-reset { padding: 10px 12px !important; min-width: auto; }

/* Segmented Control */
.segmented-control { display: inline-flex; position: relative; background: var(--c-border-light); padding: 4px; border-radius: var(--radius-sm); width: max-content; }
.segment-bg { position: absolute; top: 4px; left: 4px; bottom: 4px; width: calc(50% - 4px); background: var(--c-surface); border-radius: calc(var(--radius-sm) - 2px); box-shadow: var(--shadow-sm); transition: transform 0.2s ease; }
.segment-bg.is-right { transform: translateX(100%); }
.segment-btn { position: relative; z-index: 1; padding: 6px 16px; min-width: 100px; border: none; background: transparent; font-size: 13px; font-weight: 600; cursor: pointer; border-radius: calc(var(--radius-sm) - 2px); }
.text-primary { color: var(--c-primary); }
.text-muted { color: var(--c-text-muted); font-weight: 500; }
.text-warning { color: var(--c-warning); font-weight: 600; }
.text-title { color: var(--c-text-title); font-weight: 700; }
.mode-info { margin-top: 16px; font-size: 13px; }
.inline-icon { display: inline-flex; vertical-align: middle; margin-top: -2px; }
.inline-icon :deep(svg) { width: 16px; height: 16px; }

/* Buttons */
.btn { padding: 10px 20px; font-size: 14px; font-weight: 600; border-radius: var(--radius-sm); cursor: pointer; border: none; display: inline-flex; align-items: center; justify-content: center; gap: 8px; transition: all 0.2s; }
.btn-lg { padding: 14px 24px; font-size: 15px; }
.w-full { width: 100%; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-primary { background: var(--c-primary); color: white; }
.btn-primary:not(:disabled):hover { background: var(--c-primary-hover); }
.btn-outline { background: var(--c-surface); border: 1px solid var(--c-border); color: var(--c-text-main); }
.btn-outline:not(:disabled):hover { border-color: var(--c-primary); color: var(--c-primary); }

/* Summary */
.summary-card { background: #fbfbfc; }
.summary-title { font-size: 13px; font-weight: 600; color: var(--c-text-muted); text-transform: uppercase; margin: 0 0 16px 0; letter-spacing: 0.5px; }
.summary-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; background: white; padding: 16px; border: 1px solid var(--c-border-light); border-radius: var(--radius-sm); margin-bottom: 24px; }
.sum-row { display: flex; flex-direction: column; gap: 4px; font-size: 13px; }
.sum-row span { color: var(--c-text-muted); }
.sum-row strong { color: var(--c-text-title); font-weight: 600; }

/* Spinner */
.spinner { display: inline-block; width: 16px; height: 16px; border: 2px solid rgba(255,255,255,0.3); border-radius: 50%; border-top-color: white; animation: spin 0.8s linear infinite; }
.spinner-sm { display: inline-block; width: 14px; height: 14px; border: 2px solid #e5e7eb; border-top-color: var(--c-primary); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.error-msg { font-size: 13px; font-weight: 600; text-align: center; margin-bottom: 16px; }
.text-danger { color: var(--c-danger); }
.submit-area { margin-top: 8px; }

/* Empty & loading */
.accounts-loading { display: flex; align-items: center; gap: 10px; font-size: 13px; color: var(--c-text-muted); padding: 12px 0; }
.empty-accounts { display: flex; flex-direction: column; align-items: center; padding: 32px 16px; text-align: center; gap: 4px; }
.empty-accounts p { margin: 0; font-size: 13px; color: var(--c-text-main); }
.ea-hint { color: var(--c-text-muted) !important; font-size: 12px !important; margin-top: 4px; }
.ea-icon { color: var(--c-text-muted); margin-bottom: 8px; }
.ea-icon :deep(svg) { width: 40px; height: 40px; stroke-width: 1.5; }
.st-error { background: #fee2e2; color: #991b1b; }
.st-unknown { background: #f3f4f6; color: #6b7280; }
.default-tag { font-size: 9px; font-weight: 800; color: var(--c-primary); border: 1px solid #bfdbfe; padding: 1px 5px; border-radius: 4px; text-transform: uppercase; vertical-align: middle; margin-left: 4px; }
</style>
