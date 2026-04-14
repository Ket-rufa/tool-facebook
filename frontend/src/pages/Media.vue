<script setup lang="ts">
import { ref } from 'vue'
import { mediaData } from '../data/mediaMockData'
import MediaToolbar from '../components/MediaToolbar.vue'
import MediaGrid from '../components/MediaGrid.vue'
import MediaDetailPanel from '../components/MediaDetailPanel.vue'

const selectedId = ref<number | null>(1)

function handleSelect(id: number) {
  selectedId.value = id
}
</script>

<template>
  <div class="media-page">
    <div class="page-header">
      <MediaToolbar />
    </div>

    <!-- 2 Cols Layout -->
    <div class="content-grid" :class="{ 'with-panel': selectedId }">
      <div class="col-left">
        <MediaGrid 
          :items="mediaData.items"
          :selected-id="selectedId"
          @select="handleSelect"
        />
      </div>
      <div class="col-right" v-if="selectedId">
        <div class="panel-sticky">
          <MediaDetailPanel 
            :detail="mediaData.detail" 
            @close="selectedId = null" 
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.media-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
  height: 100%;
}

.page-header {
  margin-bottom: 8px;
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
  align-items: flex-start;
  flex: 1;
}

.content-grid.with-panel {
  grid-template-columns: 1fr 340px;
}

.col-left {
  display: flex;
  flex-direction: column;
}

.col-right {
  display: flex;
  flex-direction: column;
  height: 100%;
  position: relative;
}

/* Simulate sticky panel inside scrollable area or relative area */
.panel-sticky {
  position: sticky;
  top: 0;
  height: calc(100vh - 180px); /* estimate based on header height */
}
</style>
