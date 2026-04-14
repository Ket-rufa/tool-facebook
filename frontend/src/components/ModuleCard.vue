<script setup lang="ts">
import { icons } from '../utils/icons'
import StatusPill from './StatusPill.vue'

defineProps<{
  card: any
  size?: 'large' | 'small'
}>()
</script>

<template>
  <div class="module-card" :class="`size-${size || 'large'}`">
    <div class="mc-header">
      <div v-if="size === 'large'" class="mc-icon" v-html="icons[card.icon as keyof typeof icons]"></div>
      <div v-else class="mc-icon-small" v-html="icons[card.icon as keyof typeof icons]"></div>
      <StatusPill :text="card.badge.text" :color="card.badge.color" />
    </div>
    <div class="mc-body">
      <h3 class="mc-title">{{ card.title }}</h3>
      <p class="mc-desc">{{ card.description }}</p>
    </div>
    <div class="mc-footer">
      <a href="#" class="mc-cta">
        {{ card.cta }} 
        <span v-if="size === 'large'" class="cta-icon" v-html="icons.arrowRight"></span>
        <span v-else-if="card.icon === 'fileText'" class="cta-icon" v-html="icons.externalLink"></span>
      </a>
    </div>
  </div>
</template>

<style scoped>
.module-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-sm);
  transition: box-shadow 0.2s, transform 0.2s;
  height: 100%;
}

.module-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.mc-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.mc-icon {
  width: 40px;
  height: 40px;
  background: #eff6ff;
  color: var(--c-primary);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.mc-icon :deep(svg) { width: 24px; height: 24px; stroke-width: 2.5; }

.mc-icon-small {
  width: 32px;
  height: 32px;
  background: #f3f4f6;
  color: var(--c-text-main);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.mc-body {
  flex: 1;
}

.size-large .mc-title { font-size: 18px; margin: 0 0 12px; font-weight: 700; color: var(--c-text-title); }
.size-small .mc-title { font-size: 15px; margin: 0 0 8px; font-weight: 700; color: var(--c-text-title); }

.mc-desc {
  font-size: 13px;
  line-height: 1.5;
  color: var(--c-text-muted);
  margin: 0;
}

.mc-footer {
  margin-top: 24px;
}

.mc-cta {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 700;
  text-decoration: none;
}

.size-large .mc-cta { color: var(--c-primary); }
.size-small .mc-cta { color: var(--c-primary); text-transform: uppercase; font-size: 11px; }

.mc-cta:hover { text-decoration: underline; }

.cta-icon { display: flex; }
.cta-icon :deep(svg) { width: 14px; height: 14px; }
</style>
