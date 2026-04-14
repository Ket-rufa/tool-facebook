<script setup lang="ts">
import StatCard from '../components/StatCard.vue'
import SystemHealthCard from '../components/SystemHealthCard.vue'
import ActivityList from '../components/ActivityList.vue'
import TopProfilesCard from '../components/TopProfilesCard.vue'
import { dashboardData } from '../data/mockData'
</script>

<template>
  <div class="overview-page">
    <div class="page-header">
      <h1>Tổng quan hệ thống</h1>
      <span class="last-updated">Cập nhật lúc {{ dashboardData.lastUpdated }}</span>
    </div>

    <div class="kpi-grid">
      <StatCard 
        v-for="kpi in dashboardData.kpiCards" 
        :key="kpi.id"
        :label="kpi.label"
        :value="kpi.value"
        :sub-value="kpi.subValue"
        :sub-type="kpi.subType"
        :icon-right="kpi.iconRight"
      />
    </div>

    <div class="section">
      <SystemHealthCard 
        :status="dashboardData.systemHealth.status" 
        :modules="dashboardData.systemHealth.modules" 
      />
    </div>

    <div class="bottom-grid">
      <div class="col-left">
        <ActivityList :activities="dashboardData.recentActivities" />
      </div>
      <div class="col-right">
        <TopProfilesCard :profiles="dashboardData.topProfiles" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.overview-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header h1 {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 700;
  color: var(--c-text-title);
}

.last-updated {
  font-size: 13px;
  color: var(--c-text-muted);
}

.kpi-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 20px;
}

.section {
  display: flex;
  flex-direction: column;
}

.bottom-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

/* Responsive adjust if necessary, though requested desktop only */
@media (max-width: 1200px) {
  .kpi-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>
