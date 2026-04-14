<script setup lang="ts">
defineProps<{
  tabs: string[]
  activeTab: string
}>()

const emit = defineEmits(['change'])
</script>

<template>
  <div class="tabs-container">
    <button 
      v-for="tab in tabs" 
      :key="tab"
      class="tab-item"
      :class="{ 'active': activeTab === tab }"
      @click="emit('change', tab)"
    >
      {{ tab }}
    </button>
  </div>
</template>

<style scoped>
.tabs-container {
  display: flex;
  border-bottom: 1px solid var(--c-border);
  margin-bottom: 24px;
}

.tab-item {
  background: transparent;
  border: none;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-muted);
  padding: 12px 16px;
  cursor: pointer;
  position: relative;
  transition: color 0.2s;
}

.tab-item:hover {
  color: var(--c-text-main);
}

.tab-item.active {
  color: var(--c-primary);
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: var(--c-primary);
}
</style>
