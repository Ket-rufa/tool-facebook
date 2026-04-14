<script setup lang="ts">
defineProps<{
  source: any
  content: string
}>()
</script>

<template>
  <div class="post-content-card">
    <div class="pc-header">
      <div 
        class="pc-avatar" 
        :style="{ backgroundColor: source.bgColor, color: source.textColor }"
      >
        {{ source.avatar }}
      </div>
      <div class="pc-meta">
        <div class="pc-name">{{ source.name }}</div>
        <div class="pc-time">{{ source.time }}</div>
      </div>
    </div>
    <div class="pc-body">
      <p class="pc-text">{{ content }}</p>
    </div>
  </div>
</template>

<style scoped>
.post-content-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.pc-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.pc-avatar {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
}

.pc-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.pc-name {
  font-size: 16px;
  font-weight: 700;
  color: var(--c-text-title);
}

.pc-time {
  font-size: 12px;
  color: var(--c-text-muted);
}

.pc-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: var(--c-text-main);
  white-space: pre-wrap;
}
</style>
