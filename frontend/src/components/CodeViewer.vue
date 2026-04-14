<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  code: string
  showCopy?: boolean
}>()

const copied = ref(false)

function copyCode() {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(props.code)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}
</script>

<template>
  <div class="code-viewer-wrapper">
    <div v-if="showCopy" class="code-header">
      <button class="btn-copy" @click="copyCode">
        <svg v-if="copied" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>
        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
        {{ copied ? 'Đã sao chép' : 'SAO CHÉP' }}
      </button>
    </div>
    <pre class="code-viewer"><code>{{ code }}</code></pre>
  </div>
</template>

<style scoped>
.code-viewer-wrapper {
  background-color: #111827;
  border-radius: var(--radius-sm);
  overflow: hidden;
  position: relative;
}

.code-header {
  display: flex;
  justify-content: flex-end;
  padding: 8px 12px;
  background-color: rgba(255,255,255,0.05);
  border-bottom: 1px solid rgba(255,255,255,0.1);
}

.btn-copy {
  display: flex;
  align-items: center;
  gap: 6px;
  background: transparent;
  border: none;
  color: #9ca3af;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
}

.btn-copy:hover {
  background-color: rgba(255,255,255,0.1);
  color: white;
}

.code-viewer {
  margin: 0;
  padding: 16px;
  color: #e5e7eb;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 12px;
  line-height: 1.5;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
