<script setup lang="ts">
defineProps<{
  text: string
  color: 'warning' | 'info' | 'default'
}>()
</script>

<template>
  <span class="status-pill" :class="`pill-${color}`">
    {{ text }}
  </span>
</template>

<style scoped>
.status-pill {
  font-size: 9px;
  font-weight: 800;
  padding: 4px 8px;
  border-radius: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.pill-warning { background: #fef3c7; color: #b45309; }
.pill-info { background: #eff6ff; color: #1d4ed8; }
.pill-default { background: #f3f4f6; color: #4b5563; }
</style>
