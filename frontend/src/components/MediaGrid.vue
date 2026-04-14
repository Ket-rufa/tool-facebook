<script setup lang="ts">
import MediaGridItem from './MediaGridItem.vue'

defineProps<{
  items: any[]
  selectedId: number | null
}>()

const emit = defineEmits(['select'])
</script>

<template>
  <div class="media-grid-container">
    <MediaGridItem 
      v-for="item in items" 
      :key="item.id" 
      :item="item"
      :is-selected="item.id === selectedId"
      @select="emit('select', item.id)"
    />
  </div>
</template>

<style scoped>
.media-grid-container {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  align-content: flex-start;
}
</style>
