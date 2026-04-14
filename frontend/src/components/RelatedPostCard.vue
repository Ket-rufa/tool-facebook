<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  post: {
    preview: string
    link: string
  }
}>()
</script>

<template>
  <div class="related-post-card">
    <p class="rp-text">{{ post.preview }}</p>
    <div class="rp-actions">
      <button class="btn-rp btn-rp-solid"><span v-html="icons.externalLink"></span> Xem bài gốc</button>
      <button class="btn-rp btn-rp-outline"><span v-html="icons.share2"></span> Mở trên Facebook</button>
    </div>
  </div>
</template>

<style scoped>
.related-post-card {
  background: #f9fafb;
  border-radius: var(--radius-sm);
  padding: 16px;
}

.rp-text {
  margin: 0 0 16px;
  font-size: 13px;
  color: var(--c-text-main);
  line-height: 1.5;
  font-style: italic;
}

.rp-actions {
  display: flex;
  gap: 8px;
}

.btn-rp {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  border: 1px solid transparent;
}

.btn-rp-solid {
  background: #1e40af; /* Darker blue for contrast like picture */
  color: white;
}
.btn-rp-solid:hover { background: #1e3a8a; }

.btn-rp-outline {
  background: white;
  border-color: var(--c-border);
  color: var(--c-text-title);
}
.btn-rp-outline:hover { background: #f3f4f6; }

.btn-rp :deep(svg) { width: 12px; height: 12px; }
</style>
