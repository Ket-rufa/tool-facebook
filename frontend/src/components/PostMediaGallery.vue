<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  media: any[]
}>()
</script>

<template>
  <div class="post-media-gallery">
    <div class="gallery-title">
      <span class="icon" v-html="icons.image"></span>
      THƯ VIỆN PHƯƠNG TIỆN ({{ media.length }})
    </div>
    
    <div class="media-grid">
      <div v-for="item in media" :key="item.id" class="media-item">
        <template v-if="item.type === 'image'">
          <img :src="item.url" class="m-img" alt="media" />
        </template>
        <template v-else-if="item.type === 'pdf'">
          <div class="m-pdf">
            <span class="pdf-icon"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg></span>
            <span class="pdf-title">{{ item.title }}</span>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.post-media-gallery {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.gallery-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
  font-weight: 700;
  color: var(--c-text-muted);
  letter-spacing: 0.5px;
  margin-bottom: 16px;
}

.icon {
  display: flex;
  align-items: center;
}

.media-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.media-item {
  aspect-ratio: 16 / 9;
  border-radius: var(--radius-sm);
  overflow: hidden;
  background: #f3f4f6;
  position: relative;
  cursor: pointer;
}

.media-item::after {
  content: '';
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0);
  transition: background 0.2s;
}

.media-item:hover::after {
  background: rgba(0,0,0,0.05);
}

.m-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.m-pdf {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #e5e7eb;
  color: var(--c-text-muted);
  gap: 8px;
  padding: 16px;
  text-align: center;
}

.pdf-title {
  font-size: 11px;
  font-weight: 600;
}
</style>
