<script setup lang="ts">
import { postDetailData as data } from '../data/postDetailMockData'
import { icons } from '../utils/icons'
import ParserWarningBanner from '../components/ParserWarningBanner.vue'
import PostContentCard from '../components/PostContentCard.vue'
import PublicMetricsCard from '../components/PublicMetricsCard.vue'
import PostMetadataCard from '../components/PostMetadataCard.vue'
import RecentChangesCard from '../components/RecentChangesCard.vue'
import PostMediaGallery from '../components/PostMediaGallery.vue'
import CrawlSummaryCard from '../components/CrawlSummaryCard.vue'
import JsonTabs from '../components/JsonTabs.vue'
</script>

<template>
  <div class="post-detail-page">
    <div class="page-header">
      <div class="title-area">
        <div class="breadcrumb">
          <span>Dữ liệu bài viết</span>
          <span class="sep">&rsaquo;</span>
          <span class="current text-primary">ID: {{ data.id }}</span>
        </div>
        <h1>Chi tiết bài viết</h1>
      </div>

      <div class="header-actions">
        <button class="btn btn-danger-soft"><span v-html="icons.trash"></span> Xóa dữ liệu</button>
        <button class="btn btn-outline"><span v-html="icons.download"></span> Tải JSON</button>
        <button class="btn btn-primary"><span v-html="icons.repeat"></span> Quét lại ngay</button>
      </div>
    </div>

    <ParserWarningBanner :text="data.warning" v-if="data.warning" />

    <div class="detail-grid">
      <div class="col-main">
        <PostContentCard :source="data.source" :content="data.content" />
        <PostMediaGallery :media="data.media" />
        <CrawlSummaryCard :summary="data.crawlSummary" />
        <JsonTabs :json-content="data.rawJson" />
      </div>

      <div class="col-side">
        <PublicMetricsCard :metrics="data.metrics" />
        <PostMetadataCard :metadata="data.metadata" />
        <RecentChangesCard :changes="data.changes" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.post-detail-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.breadcrumb {
  font-size: 11px;
  font-weight: 500;
  color: var(--c-text-muted);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.current.text-primary { color: var(--c-primary); font-weight: 600; }

.title-area h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: var(--c-text-title);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: none;
}

.btn :deep(svg) { width: 14px; height: 14px; }

.btn-danger-soft { background: #fee2e2; color: #dc2626; }
.btn-danger-soft:hover { background: #fecaca; }

.btn-outline { background: #f3f4f6; color: var(--c-text-title); }
.btn-outline:hover { background: #e5e7eb; }

.btn-primary { background: var(--c-primary); color: white; }
.btn-primary:hover { background: var(--c-primary-hover); }

.detail-grid {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

.col-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
  min-width: 0;
}

.col-side {
  width: 320px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
</style>
