<script setup lang="ts">
import { ref } from 'vue'
import { logsData } from '../data/logsMockData'
import FilterBar from '../components/FilterBar.vue'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import DetailPanel from '../components/DetailPanel.vue'
import CodeViewer from '../components/CodeViewer.vue'
import { icons } from '../utils/icons'

const columns = [
  { key: 'time', label: 'Thời gian', width: '180px' },
  { key: 'level', label: 'Mức độ', width: '120px' },
  { key: 'module', label: 'Module', width: '160px' },
  { key: 'job', label: 'Công việc', width: '160px' },
  { key: 'desc', label: 'Thông tin bổ sung' }
]

const selectedLog = ref<any>(null)

function selectLog(log: any) {
  selectedLog.value = log
}

function closePanel() {
  selectedLog.value = null
}
</script>

<template>
  <div class="logs-page" :class="{ 'has-panel': selectedLog }">
    <div class="main-content">
      <div class="search-top">
        <div class="search-input-wrap">
          <span class="s-icon" v-html="icons.search"></span>
          <input type="text" placeholder="Tìm kiếm log..." />
        </div>
      </div>
      
      <FilterBar>
        <button class="filter-chip active">Mọi mức độ <span v-html="icons.trendingUp"></span></button>
        <button class="filter-chip">Module: Crawler &times;</button>
        <button class="filter-chip">Ngày: Hôm nay <span v-html="icons.clock"></span></button>
        <select class="filter-select outline bg-transparent"><option>Hồ sơ liên quan</option></select>
        
        <template #actions>
          <button class="icon-btn" v-html="icons.fileText"></button>
          <button class="icon-btn" v-html="icons.repeat"></button>
        </template>
      </FilterBar>

      <DataTable 
        :columns="columns" 
        :data="logsData.list"
        :selected-id="selectedLog?.id"
        @row-click="selectLog"
      >
        <template #time="{ value }"><span class="log-time">{{ value }}</span></template>
        <template #level="{ value }"><StatusBadge :type="value" /></template>
        <template #module="{ value }"><span class="log-module">{{ value }}</span></template>
        <template #job="{ value }"><span class="log-job">{{ value }}</span></template>
        <template #desc="{ value }"><span class="log-desc">{{ value }}</span></template>
      </DataTable>

      <div class="pagination-bar">
        <span>Hiển thị 1-50 của 2,451 logs</span>
        <div class="page-controls">
          <button class="page-btn">&lt;</button>
          <button class="page-btn">&gt;</button>
        </div>
      </div>
    </div>

    <!-- Detail Panel -->
    <div class="side-panel dark-panel" v-if="selectedLog">
      <div class="panel-header-dark">
        <h3 class="panel-title-dark"><span v-html="icons.box"></span> CHI TIẾT NHẬT KÝ THÔ</h3>
        <button class="btn-dark-copy">Sao chép</button>
      </div>
      
      <div class="panel-body-dark">
        <CodeViewer :code="logsData.detail.rawLog" :show-copy="false" />
      </div>

      <div class="panel-footer-dark">
        <div class="meta-row-dark">
          <span class="meta-label">Mã định danh Worker</span>
          <span class="meta-value">{{ logsData.detail.workerId }}</span>
        </div>
        <div class="meta-row-dark">
          <span class="meta-label">Mã định danh liên kết</span>
          <span class="meta-value">{{ logsData.detail.correlationId }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.logs-page {
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

.search-top {
  margin-bottom: 24px;
  display: flex;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background-color: #f3f4f6;
  border-radius: var(--radius-md);
  padding: 0 16px;
  width: 100%;
  max-width: 600px;
  height: 44px;
}

.search-input-wrap input {
  border: none;
  background: transparent;
  width: 100%;
  height: 100%;
  outline: none;
  font-size: 14px;
  color: var(--c-text-main);
  padding-left: 10px;
}

.s-icon { color: var(--c-text-muted); display: flex; }

.icon-btn {
  background: transparent; border: none; cursor: pointer; color: var(--c-text-muted); padding: 4px; border-radius: 4px;
}
.icon-btn:hover { background: #f3f4f6; color: var(--c-text-main); }
.outline { border: 1px solid var(--c-border); background-color: transparent !important; }

.log-time { font-family: 'Consolas', monospace; font-size: 13px; color: var(--c-text-muted); }
.log-module { font-size: 13px; font-weight: 600; }
.log-job { font-size: 13px; color: var(--c-text-muted); }
.log-desc { font-size: 13px; color: var(--c-text-main); }

.pagination-bar {
  display: flex; justify-content: space-between; align-items: center; padding: 16px;
  border: 1px solid var(--c-border); border-top: none; font-size: 13px; color: var(--c-text-muted); background: white;
}
.page-controls { display: flex; gap: 4px; }
.page-btn { background: transparent; border: none; padding: 4px 8px; cursor: pointer; color: var(--c-text-main); }
.page-btn:hover { background: #f3f4f6; }

/* Dark Panel overrrides just for Logs strictly following the instruction "CodeViewer/raw panel được nền tối" */
.side-panel.dark-panel {
  width: 450px; flex-shrink: 0; background: #111827; height: 100vh;
  margin: -32px -32px -32px 0; display: flex; flex-direction: column;
}

.panel-header-dark {
  display: flex; justify-content: space-between; align-items: center; padding: 20px 24px;
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.panel-title-dark { margin: 0; font-size: 13px; color: white; display: flex; align-items: center; gap: 8px; }

.btn-dark-copy {
  background: rgba(255,255,255,0.1); border: none; color: #d1d5db; padding: 6px 12px; border-radius: 4px; font-size: 12px; cursor: pointer;
}
.btn-dark-copy:hover { background: rgba(255,255,255,0.2); color: white; }

.panel-body-dark { flex: 1; overflow-y: auto; padding: 24px; /* code viewer background will naturally blend here since it is #111827 too */ }

.panel-footer-dark { padding: 16px 24px; border-top: 1px solid rgba(255,255,255,0.1); background: #1f2937; }
.meta-row-dark { display: flex; justify-content: space-between; font-size: 11px; font-family: monospace; color: #9ca3af; margin-bottom: 8px; }
.meta-row-dark:last-child { margin-bottom: 0; }
</style>
