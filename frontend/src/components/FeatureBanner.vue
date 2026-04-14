<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  banner: {
    title: string
    description: string
    cta: string
  }
}>()
</script>

<template>
  <div class="feature-banner">
    <div class="fb-badge">
      <span class="fb-dot"></span>
      TIÊU CHUẨN HỆ THỐNG
    </div>
    <h2 class="fb-title">{{ banner.title }}</h2>
    <p class="fb-desc">{{ banner.description }}</p>
    <button class="btn btn-primary">
      {{ banner.cta }}
      <span class="icon" v-html="icons.arrowRight"></span>
    </button>
  </div>
</template>

<style scoped>
.feature-banner {
  background: linear-gradient(180deg, #f9fafb 0%, #ffffff 100%);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.fb-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: white;
  border: 1px solid var(--c-border);
  padding: 4px 12px;
  border-radius: 16px;
  font-size: 10px;
  font-weight: 800;
  color: var(--c-text-muted);
  margin-bottom: 24px;
}

.fb-dot {
  width: 6px;
  height: 6px;
  background: var(--c-primary);
  border-radius: 50%;
}

.fb-title {
  font-size: 28px;
  font-weight: 800;
  color: var(--c-text-title);
  margin: 0 0 16px;
}

.fb-desc {
  font-size: 15px;
  color: var(--c-text-muted);
  max-width: 600px;
  line-height: 1.6;
  margin: 0 0 32px;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: #0f62fe; /* solid blue from mockup */
  color: white;
  padding: 12px 24px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 700;
  border: none;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-primary:hover { background: #0050e6; }
.icon { display: flex; }
.icon :deep(svg) { width: 16px; height: 16px; }
</style>
