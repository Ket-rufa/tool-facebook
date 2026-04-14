<script setup lang="ts">
import { icons } from '../utils/icons'

defineProps<{
  profiles: {
    id: number | string
    name: string
    displayValue: string
    progress: number
    color: string
  }[]
}>()

function getAvatarUrl(name: string) {
  // Use UI avatars as a dummy thumbnail logic, or place dummy icons.
  // Using generic images or generated avatars depending on requirements.
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}&background=111827&color=fff&size=40`
}
</script>

<template>
  <div class="top-profiles-card">
    <div class="card-header">
      <h3>HỒ SƠ NHIỀU BÀI NHẤT</h3>
      <span class="icon" v-html="icons.barChart"></span>
    </div>

    <div class="profiles-list">
      <div v-for="profile in profiles" :key="profile.id" class="profile-item">
        <img :src="getAvatarUrl(profile.name)" alt="avatar" class="profile-avatar" />
        <div class="profile-info">
          <div class="profile-name">{{ profile.name }}</div>
          <div class="progress-track">
            <div class="progress-bar" :class="profile.color" :style="{ width: profile.progress + '%' }"></div>
          </div>
        </div>
        <div class="profile-value">{{ profile.displayValue }}</div>
      </div>
    </div>

    <button class="btn-manage-all">QUẢN LÝ TẤT CẢ HỒ SƠ</button>
  </div>
</template>

<style scoped>
.top-profiles-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
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

.icon {
  color: var(--c-text-title);
  display: flex;
  align-items: center;
}

.profiles-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex: 1;
}

.profile-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.profile-avatar {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm);
  object-fit: cover;
  background-color: #1f2937;
}

.profile-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.profile-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--c-text-title);
}

.progress-track {
  height: 4px;
  background-color: #f3f4f6;
  border-radius: 2px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  border-radius: 2px;
}

.profile-value {
  font-size: 13px;
  font-weight: 500;
  color: var(--c-text-main);
  min-width: 40px;
  text-align: right;
}

.btn-manage-all {
  margin-top: 24px;
  width: 100%;
  padding: 12px;
  background-color: #f8fafc;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  color: var(--c-text-title);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
  font-family: inherit;
}

.btn-manage-all:hover {
  background-color: #f1f5f9;
}
</style>
