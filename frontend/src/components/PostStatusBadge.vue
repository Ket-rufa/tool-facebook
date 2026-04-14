<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
}>()

const statusStyle = computed(() => {
  if (props.status === 'HOÀN TẤT') return 's-success'
  if (props.status === 'ĐANG QUÉT') return 's-primary'
  if (props.status === 'LỖI QUÉT') return 's-danger'
  return 's-default'
})
</script>

<template>
  <span class="post-badge" :class="statusStyle">
    <span class="text-inner">{{ status }}</span>
  </span>
</template>

<style scoped>
.post-badge {
  display: inline-flex;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: pre-wrap;
  text-align: center;
  max-width: 60px;
  line-height: 1.2;
}

.s-success {
  background-color: #dcfce7;
  color: #166534;
}

.s-primary {
  background-color: #dbeafe;
  color: #1e40af;
}

.s-danger {
  background-color: #fee2e2;
  color: #991b1b;
}

.s-default {
  background-color: #f3f4f6;
  color: #4b5563;
}
</style>
