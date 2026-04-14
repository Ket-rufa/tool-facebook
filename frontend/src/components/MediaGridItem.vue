<script setup lang="ts">
defineProps<{
  item: any
  isSelected: boolean
}>()

const emit = defineEmits(['select'])
</script>

<template>
  <div class="media-grid-item" :class="{ selected: isSelected }" @click="emit('select')">
    <img :src="item.url" class="m-thumb" alt="thumb" />
    
    <div class="m-overlay"></div>
    
    <div v-if="item.type === 'video'" class="m-video-meta">
      <span class="v-icon"><svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor" stroke="none"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg></span>
      {{ item.duration }}
    </div>

    <div v-if="isSelected" class="m-check">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
    </div>
  </div>
</template>

<style scoped>
.media-grid-item {
  aspect-ratio: 1;
  border-radius: var(--radius-md);
  position: relative;
  overflow: hidden;
  cursor: pointer;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  border: 2px solid transparent; /* reserved for selected state */
  transition: all 0.2s;
}

.media-grid-item:hover .m-overlay {
  opacity: 1;
}

.media-grid-item.selected {
  border-color: var(--c-primary);
}

.m-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.media-grid-item:hover .m-thumb {
  transform: scale(1.05);
}

.m-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.1);
  opacity: 0;
  transition: opacity 0.2s;
}

.m-video-meta {
  position: absolute;
  bottom: 8px;
  right: 8px;
  background: rgba(0,0,0,0.7);
  color: white;
  font-size: 11px;
  font-weight: 600;
  padding: 4px 6px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  z-index: 2;
}

.m-check {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 20px;
  height: 20px;
  background: var(--c-primary);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
  box-shadow: 0 1px 4px rgba(0,0,0,0.2);
}
</style>
