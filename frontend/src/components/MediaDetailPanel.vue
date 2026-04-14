<script setup lang="ts">
import { icons } from '../utils/icons'
import RelatedPostCard from './RelatedPostCard.vue'

defineProps<{
  detail: any
}>()

const emit = defineEmits(['close'])
</script>

<template>
  <div class="media-detail-panel">
    <div class="panel-header">
      <h3 class="panel-title">Chi tiết phương tiện</h3>
      <button class="btn-close" @click="emit('close')" v-html="icons.x"></button>
    </div>

    <div class="panel-body">
      <div class="preview-wrap">
        <!-- Mock assuming it's an image for simplicity, video will have play icon if needed via item data -->
        <img :src="detail.url" class="p-img" alt="preview" />
      </div>

      <div class="detail-block">
        <label>TÊN TỆP</label>
        <div class="d-val">{{ detail.fileName }}</div>
      </div>

      <div class="detail-row">
        <div class="detail-block">
          <label>KÍCH THƯỚC</label>
          <div class="d-val">{{ detail.size }}</div>
        </div>
        <div class="detail-block">
          <label>ĐỘ PHÂN GIẢI</label>
          <div class="d-val">{{ detail.resolution }}</div>
        </div>
      </div>

      <div class="detail-block">
        <label>NGÀY THU THẬP</label>
        <div class="d-val">{{ detail.date }}</div>
      </div>

      <div class="section-divider"></div>

      <div class="section-block">
        <h4 class="s-title">HỒ SƠ GỐC</h4>
        <div class="source-card">
          <img :src="detail.source.avatar" class="s-avatar" alt="avatar" />
          <div class="s-info">
            <div class="s-name">{{ detail.source.name }}</div>
            <div class="s-handle">{{ detail.source.handle }}</div>
          </div>
          <a href="#" class="s-link" v-html="icons.externalLink"></a>
        </div>
      </div>

      <div class="section-block">
        <h4 class="s-title">BÀI VIẾT LIÊN QUAN</h4>
        <RelatedPostCard :post="detail.post" />
      </div>
    </div>

    <div class="panel-footer">
      <button class="btn-action btn-download"><span v-html="icons.download"></span> Tải về</button>
      <button class="btn-action btn-remove"><span v-html="icons.trash"></span> Gỡ khỏi danh sách</button>
      <div class="footer-note">
        Hành động này sẽ gỡ tệp khỏi trình quản lý nhưng không ảnh hưởng đến dữ liệu gốc trên Facebook.
      </div>
    </div>
  </div>
</template>

<style scoped>
.media-detail-panel {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  height: 100%; /* assume it fills col-right */
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--c-border);
}

.panel-title {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
  color: var(--c-text-title);
}

.btn-close {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--c-text-muted);
  display: flex;
}
.btn-close:hover { color: var(--c-text-main); }

.panel-body {
  padding: 24px;
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.preview-wrap {
  background: #f3f4f6;
  border-radius: var(--radius-sm);
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
}

.p-img {
  max-width: 100%;
  max-height: 250px;
  object-fit: contain;
  box-shadow: var(--shadow-sm);
  border-radius: 4px;
}

.detail-row {
  display: flex;
  gap: 24px;
}

.detail-block {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-block label {
  font-size: 10px;
  font-weight: 700;
  color: var(--c-text-muted);
  text-transform: uppercase;
}

.d-val {
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
  word-break: break-all;
}

.section-divider {
  height: 1px;
  background: var(--c-border);
  margin: 8px 0;
}

.section-block {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.s-title {
  margin: 0;
  font-size: 10px;
  font-weight: 700;
  color: var(--c-text-muted);
  text-transform: uppercase;
}

.source-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: #f9fafb;
  padding: 12px;
  border-radius: var(--radius-sm);
}

.s-avatar { width: 32px; height: 32px; border-radius: 50%; object-fit: cover; }
.s-info { flex: 1; display: flex; flex-direction: column; }
.s-name { font-size: 13px; font-weight: 700; color: var(--c-text-title); }
.s-handle { font-size: 11px; color: var(--c-text-muted); }
.s-link { color: var(--c-text-muted); display: flex; }
.s-link:hover { color: var(--c-text-title); }

.panel-footer {
  padding: 24px;
  border-top: 1px solid var(--c-border);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.btn-action {
  padding: 10px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.btn-download { background: #e5e7eb; color: var(--c-text-title); }
.btn-download:hover { background: #d1d5db; }

.btn-remove { background: #fee2e2; color: #dc2626; letter-spacing: 0.5px; text-transform: uppercase;}
.btn-remove:hover { background: #fecaca; }

.footer-note {
  font-size: 10px;
  color: var(--c-text-muted);
  text-align: center;
  font-style: italic;
  margin-top: 4px;
  line-height: 1.4;
}
</style>
