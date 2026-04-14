<script setup lang="ts">
defineProps<{
  activities: {
    id: number | string
    type: string
    title: string
    description: string
    timeago: string
  }[]
}>()
</script>

<template>
  <div class="activity-card">
    <div class="card-header">
      <h3>HOẠT ĐỘNG GẦN ĐÂY</h3>
      <a href="#" class="view-all">Xem tất cả nhật ký</a>
    </div>

    <div class="activity-list">
      <div v-for="act in activities" :key="act.id" class="activity-item">
        <div class="status-dot-wrapper">
          <div class="status-dot" :class="act.type"></div>
        </div>
        <div class="activity-content">
          <div class="activity-header">
            <span class="activity-title" v-html="act.title"></span>
            <span class="activity-time">{{ act.timeago }}</span>
          </div>
          <p class="activity-desc">{{ act.description }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.activity-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.card-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--c-text-title);
  letter-spacing: 0.5px;
}

.view-all {
  font-size: 13px;
  font-weight: 600;
  color: var(--c-primary);
  text-decoration: none;
}

.view-all:hover {
  text-decoration: underline;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.activity-item {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.status-dot-wrapper {
  padding-top: 6px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.info, .status-dot.success {
  background-color: var(--c-primary);
}

.status-dot.danger {
  background-color: var(--c-danger);
}

.activity-content {
  flex: 1;
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 6px;
}

.activity-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--c-text-title);
}

.activity-time {
  font-size: 12px;
  color: var(--c-text-muted);
}

.activity-desc {
  margin: 0;
  font-size: 13px;
  color: var(--c-text-muted);
  line-height: 1.5;
}
</style>
