<script setup lang="ts">
import { ref, onMounted } from 'vue'
import TargetInputPanel from '../components/crawl/TargetInputPanel.vue'
import CrawlConfigPanel from '../components/crawl/CrawlConfigPanel.vue'
import RuntimeStatusPanel from '../components/crawl/RuntimeStatusPanel.vue'
import RuntimeLogPanel from '../components/crawl/RuntimeLogPanel.vue'
import EntityResultViewer from '../components/crawl/EntityResultViewer.vue'

// Import Wails Backend bindings
import { 
  ScanProfileInfo, 
  ExportedResolveUID 
} from '../../wailsjs/go/crawl/CrawlHandler'
import { ListAccounts } from '../../wailsjs/go/accounts/AccountService'

// Reactive state
const targets = ref<any[]>([])
const config = ref({
  preset: 'Standard',
  fieldGroups: ['Danh tính', 'Thông tin cơ bản', 'Địa điểm', 'Công việc', 'Chỉ số'],
  concurrency: 1,
  delay: 2000,
  retry: 3,
  maxPosts: 10
})
const runtime = ref({
  total: 0,
  queued: 0,
  running: 0,
  success: 0,
  failed: 0,
  partial: 0,
  eta: '--:--',
  progress: 0,
  currentEntity: '',
  currentUid: '',
  currentStep: 'Sẵn sàng'
})
const logs = ref<any[]>([])
const entities = ref<any[]>([])
const selectedAccountID = ref<string>('')

onMounted(async () => {
  const accounts = await ListAccounts()
  if (accounts && accounts.length > 0) {
    selectedAccountID.value = accounts[0].id
  }
})

const addLog = (message: string, level: string = 'info', uid: string | null = null, step: string = 'process') => {
  const time = new Date().toLocaleTimeString()
  logs.value.unshift({ id: Date.now(), time, level, message, uid, step })
}

const handleAddTargets = (input: string) => {
  const lines = input.split('\n').map(l => l.trim()).filter(l => l)
  lines.forEach(line => {
    // Basic check if it's already in the list
    if (targets.value.some(t => t.url === line)) return
    
    targets.value.push({
      id: Date.now().toString() + Math.random(),
      url: line,
      displayName: null,
      avatar: null,
      type: 'Profile',
      uid: null,
      resolveStatus: 'Pending',
      publicAccess: 'Unknown'
    })
  })
  updateRuntimeSummary()
}

const handleClear = () => {
  targets.value = []
  updateRuntimeSummary()
}

const handleResolve = async () => {
  if (!selectedAccountID.value) {
    addLog('Cần chọn tài khoản trong phần Cài đặt/Tài khoản trước.', 'error')
    return
  }

  for (let i = 0; i < targets.value.length; i++) {
    const target = targets.value[i]
    if (target.uid) continue

    addLog(`Đang resolve UID cho: ${target.url}`, 'info')
    runtime.value.currentStep = `resolving ${target.url}`
    
    try {
      const resolved = await ExportedResolveUID(target.url, selectedAccountID.value)
      if (resolved && resolved !== target.url) {
        target.uid = resolved
        target.resolveStatus = 'Resolved'
        addLog(`Resolved: ${target.url} -> ${resolved}`, 'info', resolved)
      } else {
        target.resolveStatus = 'Failed'
        addLog(`Không thể resolve UID cho: ${target.url}`, 'warning')
      }
    } catch (e) {
      target.resolveStatus = 'Failed'
      addLog(`Lỗi khi resolve: ${e}`, 'error')
    }
  }
  updateRuntimeSummary()
}

const handleRun = async () => {
  if (!selectedAccountID.value) {
    addLog('Cần ít nhất một tài khoản để thực hiện quét.', 'error')
    return
  }

  if (targets.value.length === 0) {
    addLog('Chưa có mục tiêu nào để quét.', 'warning')
    return
  }

  // Ensure UIDs are resolved
  await handleResolve()

  let processed = 0
  runtime.value.total = targets.value.length
  runtime.value.running = 1
  runtime.value.queued = targets.value.length

  for (let i = 0; i < targets.value.length; i++) {
    const target = targets.value[i]
    if (!target.uid) {
      processed++
      runtime.value.failed++
      runtime.value.queued--
      continue
    }

    runtime.value.currentEntity = target.url
    runtime.value.currentUid = target.uid
    runtime.value.currentStep = 'fetching profile info...'
    addLog(`Bắt đầu quét profile UID: ${target.uid}`, 'info', target.uid)

    try {
      const accounts = await ListAccounts()
      const acc = accounts.find(a => a.id === selectedAccountID.value)
      if (!acc) throw new Error('Account not found')

      const result = await ScanProfileInfo(target.uid as string, acc.cookie as string)
      if (result) {
        target.displayName = result.name
        target.publicAccess = 'Accessible'
        
        // Add to entities list
        const entity = {
          entity_id: 'e_' + Date.now(),
          uid: target.uid,
          entity_type: 'Profile',
          canonical_url: `https://facebook.com/${target.uid}`,
          scan_status: 'success',
          last_updated: new Date().toISOString(),
          avatar: null,
          identity: { display_name: result.name, alternate_name: null },
          basic_info: { bio: result.bio, birthday_text: result.birthday },
          location: { current_city: result.current_city, hometown: result.hometown },
          education: result.education || [],
          work: result.work || [],
          metrics: { followers_text: `${result.followers} followers`, followers_value: parseInt(result.followers) || 0 },
          links: [],
          field_status: 'full'
        }
        entities.value.unshift(entity)
        
        runtime.value.success++
        addLog(`Đã quét xong và lưu vào DATA/${target.uid}.txt`, 'info', target.uid)
      } else {
        runtime.value.failed++
      }
    } catch (e) {
      console.error(e)
      runtime.value.failed++
      addLog(`Lỗi khi quét ${target.uid}: ${e}`, 'error', target.uid)
    }

    processed++
    runtime.value.queued--
    runtime.value.progress = Math.round((processed / targets.value.length) * 100)
    
    // Delay between targets
    if (i < targets.value.length - 1) {
      await new Promise(r => setTimeout(r, config.value.delay))
    }
  }

  runtime.value.running = 0
  runtime.value.currentStep = 'Hoàn tất'
  addLog('Hoàn tất tiến trình quét.', 'info')
}

const updateRuntimeSummary = () => {
  runtime.value.total = targets.value.length
  runtime.value.queued = targets.value.length
  runtime.value.progress = 0
}
</script>

<template>
  <div class="workspace-page">
    <div class="page-header">
      <div class="title-area">
        <h1>Quét thông tin người dùng</h1>
        <p>Khu vực làm việc Quét thông tin Công khai (GraphQL + Xuất dữ liệu)</p>
      </div>
    </div>

    <!-- Top Section: Operations -->
    <div class="operations-grid">
      <!-- Zone A: Target Input -->
      <TargetInputPanel 
        :targets="targets" 
        @add="handleAddTargets"
        @resolve="handleResolve"
        @clear="handleClear"
      />
      
      <!-- Zone B: Configurations -->
      <CrawlConfigPanel :config="config" @run="handleRun" />
      
      <!-- Zone C: Runtime Monitoring -->
      <RuntimeStatusPanel :status="runtime" />
      
      <!-- Zone D: Runtime Logs -->
      <RuntimeLogPanel :logs="logs" />
    </div>

    <!-- Bottom Section: Results Container -->
    <div class="results-section">
      <div class="section-title">Kết quả trích xuất Entity (Đã lưu vào thư mục DATA)</div>
      <EntityResultViewer :entities="entities" :logs="logs" />
    </div>
  </div>
</template>

<style scoped>
.workspace-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px); /* assuming some app header offsets */
  min-height: 800px;
}

.page-header {
  margin-bottom: 24px;
}
.title-area h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #111827;
}
.title-area p {
  margin: 4px 0 0;
  font-size: 14px;
  color: #6b7280;
}

.operations-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: auto auto;
  gap: 20px;
  margin-bottom: 24px;
}

.results-section {
  display: flex;
  flex-direction: column;
  min-height: 480px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 12px;
}
</style>
