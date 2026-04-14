<script setup lang="ts">
import { ref } from 'vue'
import { integrationsData as data } from '../data/integrationsMockData'
import ModuleCard from '../components/ModuleCard.vue'
import FeatureBanner from '../components/FeatureBanner.vue'
import { icons } from '../utils/icons'

const showToast = ref(true)
</script>

<template>
  <div class="integrations-page">
    <div class="page-header">
      <h1 class="page-title">{{ data.header.title }}</h1>
      <p class="page-desc">{{ data.header.description }}</p>
    </div>

    <div class="main-cards-grid">
      <div v-for="card in data.mainCards" :key="card.id">
        <ModuleCard :card="card" size="large" />
      </div>
    </div>

    <div class="small-cards-grid">
      <div v-for="card in data.smallCards" :key="card.id">
        <ModuleCard :card="card" size="small" />
      </div>
    </div>

    <div class="banner-section">
      <FeatureBanner :banner="data.banner" />
    </div>

    <!-- Floating Toast Mock -->
    <div v-if="showToast" class="access-toast">
      <div class="t-icon text-primary" v-html="icons.info"></div>
      <div class="t-content">
        <div class="t-title">{{ data.toast.title }}</div>
        <div class="t-desc">{{ data.toast.description }}</div>
      </div>
      <button class="t-close" @click="showToast = false" v-html="icons.x"></button>
    </div>
  </div>
</template>

<style scoped>
.integrations-page {
  display: flex;
  flex-direction: column;
  gap: 32px;
  max-width: 1000px;
  position: relative;
}

.page-header {
  margin-bottom: 8px;
}

.page-title {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 800;
  color: var(--c-text-title);
}

.page-desc {
  margin: 0;
  font-size: 15px;
  color: var(--c-text-muted);
}

.main-cards-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.small-cards-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

/* Toast */
.access-toast {
  position: fixed;
  bottom: 24px;
  right: 24px;
  background: white;
  border: 1px solid var(--c-border);
  box-shadow: var(--shadow-md);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  gap: 12px;
  width: 320px;
  z-index: 50;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from { transform: translateX(100%); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}

.t-icon {
  background: #eff6ff;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.text-primary { color: var(--c-primary); }

.t-content {
  flex: 1;
}

.t-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--c-text-title);
  margin-bottom: 4px;
}

.t-desc {
  font-size: 11px;
  color: var(--c-text-muted);
  line-height: 1.4;
}

.t-close {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--c-text-muted);
  padding: 0;
  display: flex;
  height: max-content;
}
.t-close:hover { color: var(--c-text-main); }
.t-close :deep(svg) { width: 14px; height: 14px; }
</style>
