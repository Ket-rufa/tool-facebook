<script setup lang="ts">
import { ref } from 'vue'
import { jobsData } from '../data/jobsMockData'
import FilterBar from '../components/FilterBar.vue'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import DetailPanel from '../components/DetailPanel.vue'
import JobProgressCard from '../components/JobProgressCard.vue'
import WarningBox from '../components/WarningBox.vue'
import CodeViewer from '../components/CodeViewer.vue'
import { icons } from '../utils/icons'

const columns = [
  { key: 'id', label: 'Mã công việc', width: '120px' },
  { key: 'url', label: 'URL Mục tiêu', width: '35%' },
  { key: 'type', label: 'Loại', width: '120px' },
  { key: 'status', label: 'Trạng thái', width: '140px' },
  { key: 'time', label: 'Bắt đầu / Kết thúc', width: '180px' }
]

const selectedJob = ref<any>(null)

function selectJob(job: any) {
  selectedJob.value = job
}

function closePanel() {
  selectedJob.value = null
}
</script>

<template>
  <div class="jobs-page" :class="{ 'has-panel': selectedJob }">
    <div class="main-content">
      <div class="page-header">
        <div class="title-area">
          <h1>Công việc quét</h1>
          <p>Quản lý và giám sát tiến trình trích xuất dữ liệu mạng xã hội.</p>
        </div>
        <div class="header-actions">
          <button class="icon-btn" aria-label="Refresh" v-html="icons.repeat"></button>
          <button class="icon-btn" aria-label="Download" v-html="icons.fileText"></button>
        </div>
      </div>

      <FilterBar>
        <template #label>LỌC THEO:</template>
        <select class="filter-select"><option>Tất cả trạng thái</option></select>
        <select class="filter-select"><option>Mọi loại mục tiêu</option></select>
        <button class="filter-select date-btn">15/10/2023 - 22/10/2023</button>
      </FilterBar>

      <DataTable 
        :columns="columns" 
        :data="jobsData.list"
        :selected-id="selectedJob?.id"
        @row-click="selectJob"
      >
        <template #id="{ value }"><span class="job-id text-primary font-bold">{{ value }}</span></template>
        
        <template #url="{ row }">
          <div class="url-cell">
            <div class="url-text font-bold">{{ row.url }}</div>
            <div class="sub-text">{{ row.subText }}</div>
          </div>
        </template>
        
        <template #type="{ value }"><span class="type-badge">{{ value }}</span></template>
        
        <template #status="{ value }"><StatusBadge :type="value" /></template>
        
        <template #time="{ value }">
          <div class="time-cell">{{ value }}</div>
        </template>
      </DataTable>

      <div class="pagination-bar">
        <span>Đang hiển thị 1 - 20 trong tổng số 1,429 công việc</span>
        <div class="page-controls">
          <button class="page-btn">&lt;</button>
          <button class="page-btn active">1</button>
          <button class="page-btn">2</button>
          <button class="page-btn">3</button>
          <button class="page-btn">&gt;</button>
        </div>
      </div>
    </div>

    <!-- Detail Panel -->
    <div class="side-panel" v-if="selectedJob">
      <DetailPanel :title="'CHI TIẾT: ' + selectedJob.id" @close="closePanel">
        <JobProgressCard 
          :status-text="jobsData.detail.progress.statusText"
          :steps="jobsData.detail.progress.steps"
        />
        
        <div class="block-section">
          <h5>CẢNH BÁO HỆ THỐNG</h5>
          <WarningBox type="warning" :text="jobsData.detail.warning" />
        </div>

        <div class="block-section">
          <div class="section-header-flex">
            <h5>METADATA THÔ</h5>
            <button class="btn-text" @click="() => {}">
              SAO CHÉP JSON
            </button>
          </div>
          <CodeViewer :code="jobsData.detail.metadata" />
        </div>

        <template #footer>
          <div class="meta-info-row">
            <span class="meta-label">Người thực hiện:</span>
            <span class="meta-value">{{ jobsData.detail.executor }}</span>
          </div>
          <div class="meta-info-row">
            <span class="meta-label">Priority:</span>
            <span class="meta-value text-primary">{{ jobsData.detail.priority }}</span>
          </div>
        </template>
      </DetailPanel>
    </div>
  </div>
</template>

<style scoped>
.jobs-page {
  display: flex;
  height: 100%;
  gap: 24px;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.side-panel {
  width: 400px;
  flex-shrink: 0;
  height: 100vh;
  margin: -32px -32px -32px 0; /* Negate page padding to make panel full height */
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.title-area h1 { margin: 0 0 8px; font-size: 24px; font-weight: 700; }
.title-area p { margin: 0; color: var(--c-text-muted); font-size: 14px; }

.header-actions {
  display: flex;
  gap: 12px;
}

.icon-btn {
  background: white;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  width: 36px; height: 36px;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer;
  color: var(--c-text-main);
}
.icon-btn:hover { background: #f3f4f6; }

.date-btn {
  background-image: none !important;
  display: flex;
  align-items: center;
  gap: 8px;
  padding-right: 12px !important;
}

.job-id { font-size: 13px; }
.font-bold { font-weight: 600; }
.url-cell { display: flex; flex-direction: column; gap: 4px; }
.url-text { font-size: 14px; color: var(--c-text-title); }
.sub-text { font-size: 12px; color: var(--c-text-muted); }

.type-badge {
  background: #f3f4f6;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.time-cell { font-size: 12px; color: var(--c-text-muted); white-space: pre-wrap; }

.pagination-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f9fafb;
  border: 1px solid var(--c-border);
  border-top: none;
  font-size: 13px;
  color: var(--c-text-muted);
}

.page-controls { display: flex; gap: 4px; }
.page-btn {
  background: transparent; border: none; padding: 4px 8px; min-width: 28px;
  cursor: pointer; font-size: 13px; color: var(--c-text-main); border-radius: 4px;
}
.page-btn:hover { background: #e5e7eb; }
.page-btn.active { background: var(--c-primary); color: white; }

.block-section h5 { margin: 0 0 12px; font-size: 12px; letter-spacing: 0.5px; color: var(--c-text-muted); }
.section-header-flex { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 12px; }
.section-header-flex h5 { margin: 0; }
.btn-text { background: transparent; border: none; color: var(--c-primary); font-size: 11px; font-weight: 700; cursor: pointer; }
.btn-text:hover { text-decoration: underline; }

.meta-info-row { display: flex; justify-content: space-between; font-size: 13px; margin-bottom: 8px; }
.meta-info-row:last-child { margin-bottom: 0; }
.meta-label { color: var(--c-text-muted); }
.meta-value { font-weight: 600; color: var(--c-text-main); }
</style>
