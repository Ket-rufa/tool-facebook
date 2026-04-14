<script setup lang="ts">
import { ref } from 'vue'
import { profilesData } from '../data/profilesMockData'
import ProfileStatsRow from '../components/ProfileStatsRow.vue'
import ProfileListPanel from '../components/ProfileListPanel.vue'
import ProfileDetailCard from '../components/ProfileDetailCard.vue'
import LatestPostsCard from '../components/LatestPostsCard.vue'

const selectedId = ref(profilesData.list[0].id)
const currentTabFilter = ref('Tất cả')

function handleSelectProfile(p: any) {
  selectedId.value = p.id
}
</script>

<template>
  <div class="profiles-page">
    <div class="page-header">
      <div class="title-area">
        <div class="breadcrumb">
          <span>Quản lý</span>
          <span class="sep">&rsaquo;</span>
          <span class="current">Hồ sơ đã quét</span>
        </div>
        <h1>Danh mục Hồ sơ</h1>
        <p>Quản lý và giám sát {{ profilesData.stats.total.value }} thực thể đang được quét dữ liệu.</p>
      </div>

      <div class="filter-tabs">
        <button 
          v-for="tab in ['Tất cả', 'Trang', 'Cá nhân']" 
          :key="tab"
          class="f-tab"
          :class="{ active: currentTabFilter === tab }"
          @click="currentTabFilter = tab"
        >
          {{ tab }}
        </button>
      </div>
    </div>

    <!-- Top KPI Row -->
    <ProfileStatsRow :stats="profilesData.stats" />

    <!-- 2 Cols Layout -->
    <div class="content-grid">
      <div class="col-left">
        <ProfileListPanel 
          :profiles="profilesData.list"
          :selected-id="selectedId"
          @select="handleSelectProfile"
        />
      </div>
      <div class="col-right">
        <!-- Render detail for selected ID, using mock detail for all rows temporarily -->
        <ProfileDetailCard :detail="profilesData.detail" />
        <LatestPostsCard :posts="profilesData.latestPosts" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.profiles-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 8px;
}

.breadcrumb {
  font-size: 11px;
  font-weight: 500;
  color: var(--c-text-muted);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.breadcrumb .current {
  color: var(--c-primary);
  font-weight: 600;
}

.title-area h1 {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 700;
  color: var(--c-text-title);
}

.title-area p {
  margin: 0;
  color: var(--c-text-muted);
  font-size: 14px;
}

.filter-tabs {
  display: flex;
  background: #f3f4f6;
  padding: 4px;
  border-radius: var(--radius-sm);
  gap: 4px;
}

.f-tab {
  background: transparent;
  border: none;
  font-size: 12px;
  font-weight: 600;
  color: var(--c-text-muted);
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.f-tab:hover {
  color: var(--c-text-main);
}

.f-tab.active {
  background: white;
  color: var(--c-text-title);
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 24px;
  align-items: flex-start;
}

.col-left {
  display: flex;
  flex-direction: column;
}

.col-right {
  display: flex;
  flex-direction: column;
  gap: 24px;
}
</style>
