<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  metadata: {
    status: string
    postType: string
    visibility: string
    sourceWorker: string
    lastUpdate: string
  }
}>()
</script>

<template>
  <div class="card bg-white">
    <div class="card-title">METADATA CHI TIẾT</div>
    
    <div class="meta-list">
      <div class="meta-item">
        <label>TRẠNG THÁI BẢN GHI</label>
        <div class="m-val text-primary">
          <span class="m-icon"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg></span>
          {{ metadata.status }}
        </div>
      </div>

      <div class="meta-item">
        <label>KIỂU BÀI VIẾT</label>
        <div class="m-val">
          <span class="m-icon" v-html="icons.link"></span>
          {{ metadata.postType }}
        </div>
      </div>

      <div class="meta-item">
        <label>CHẾ ĐỘ HIỂN THỊ</label>
        <div class="m-val">
          <span class="m-icon" v-html="icons.globe"></span>
          {{ metadata.visibility }}
        </div>
      </div>

      <div class="meta-item">
        <label>NGUỒN QUÉT</label>
        <div class="m-val">
          <span class="m-icon" v-html="icons.server"></span>
          {{ metadata.sourceWorker }}
        </div>
      </div>

      <div class="meta-item">
        <label>THỜI GIAN CẬP NHẬT CUỐI</label>
        <div class="m-val text-primary font-bold">
          {{ metadata.lastUpdate }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.card-title {
  font-size: 11px;
  font-weight: 700;
  color: var(--c-text-muted);
  letter-spacing: 0.5px;
  margin-bottom: 20px;
}

.meta-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.meta-item label {
  font-size: 10px;
  font-weight: 600;
  color: var(--c-text-muted);
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.m-val {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--c-text-main);
  font-weight: 500;
}

.m-icon {
  display: flex;
  align-items: center;
  color: var(--c-text-muted);
}

.text-primary { color: var(--c-primary); }
.text-primary .m-icon { color: var(--c-primary); }
.font-bold { font-weight: 600; }
</style>
