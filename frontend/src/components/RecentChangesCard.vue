<script setup lang="ts">
defineProps<{
  changes: any[]
}>()
</script>

<template>
  <div class="card bg-gray-50">
    <div class="card-title">Lịch sử thay đổi gần nhất</div>
    
    <div class="timeline">
      <div v-for="(ch, idx) in changes" :key="idx" class="t-item">
        <div class="t-line">
          <div class="t-dot"></div>
          <div v-if="idx < changes.length - 1" class="t-tail"></div>
        </div>
        <div class="t-content">
          <div class="t-text">{{ ch.text }}</div>
          <div class="t-time">{{ ch.time }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: #f9fafb;
  border-radius: var(--radius-md);
  padding: 24px;
}

.card-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--c-text-title);
  margin-bottom: 20px;
}

.timeline {
  display: flex;
  flex-direction: column;
}

.t-item {
  display: flex;
  gap: 12px;
}

.t-line {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 12px;
}

.t-dot {
  width: 6px;
  height: 6px;
  background: var(--c-primary);
  border-radius: 50%;
  margin-top: 6px;
}

.t-tail {
  flex: 1;
  width: 2px;
  background: #e5e7eb;
  margin-top: 4px;
  margin-bottom: 4px;
}

.t-content {
  padding-bottom: 16px;
}
.t-item:last-child .t-content { padding-bottom: 0; }

.t-text {
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
  margin-bottom: 2px;
}

.t-time {
  font-size: 11px;
  color: var(--c-text-muted);
}
</style>
