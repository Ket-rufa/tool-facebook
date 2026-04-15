<script setup lang="ts">
import TargetInputPanel from '../components/crawl/TargetInputPanel.vue'
import CrawlConfigPanel from '../components/crawl/CrawlConfigPanel.vue'
import RuntimeStatusPanel from '../components/crawl/RuntimeStatusPanel.vue'
import RuntimeLogPanel from '../components/crawl/RuntimeLogPanel.vue'
import EntityResultViewer from '../components/crawl/EntityResultViewer.vue'
import { crawlMockData } from '../data/crawlMockData'
</script>

<template>
  <div class="workspace-page">
    <div class="page-header">
      <div class="title-area">
        <h1>Quét thông tin người dùng</h1>
        <p>Public User Info Crawling Workspace</p>
      </div>
    </div>

    <!-- Top Section: Operations -->
    <div class="operations-grid">
      <!-- Zone A: Target Input -->
      <TargetInputPanel :targets="crawlMockData.targets" />
      
      <!-- Zone B: Configurations -->
      <CrawlConfigPanel :config="crawlMockData.config" />
      
      <!-- Zone C: Runtime Monitoring -->
      <RuntimeStatusPanel :status="crawlMockData.runtime" />
      
      <!-- Zone D: Runtime Logs -->
      <RuntimeLogPanel :logs="crawlMockData.logs" />
    </div>

    <!-- Bottom Section: Results Container -->
    <div class="results-section">
      <div class="section-title">Kết quả trích xuất Entity</div>
      <EntityResultViewer :entities="crawlMockData.entities" :logs="crawlMockData.logs" />
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
