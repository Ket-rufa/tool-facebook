<script setup lang="ts">
import { ref } from 'vue'
import ProfileTabs from './ProfileTabs.vue'
import { icons } from '../utils/icons'

defineProps<{
  detail: any
}>()

const tabs = ['Tổng quan', 'Lịch sử', 'Bài viết', 'Media']
const activeTab = ref(tabs[0])
</script>

<template>
  <div class="profile-detail-card">
    <div class="cover-photo" :style="{ backgroundImage: `url(${detail.coverImg})` }">
      <button class="btn-edit-cover" aria-label="Edit Cover">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
      </button>
    </div>

    <div class="profile-header-info">
      <div class="avatar-wrap">
        <img :src="detail.avatarImg" alt="avatar" />
      </div>
      
      <div class="title-row">
        <h2>{{ detail.name }}</h2>
        <span class="badge-status">{{ detail.statusBadge }}</span>
      </div>
      <a href="#" class="handle-link">{{ detail.handle }} <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg></a>
    </div>

    <div class="card-body">
      <ProfileTabs :tabs="tabs" :active-tab="activeTab" @change="t => activeTab = t" />

      <div v-if="activeTab === 'Tổng quan'" class="tab-content">
        <div class="metrics-grid">
          <div class="metric-box">
            <span class="m-label">FOLLOWERS</span>
            <span class="m-val">{{ detail.followers }}</span>
          </div>
          <div class="metric-box">
            <span class="m-label">ENGAGEMENT</span>
            <span class="m-val">{{ detail.engagement }}</span>
          </div>
        </div>

        <div class="info-list">
          <div class="info-item">
            <div class="icon-wrap"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg></div>
            <div class="info-text">
              <span class="i-title">MÔ TẢ</span>
              <p>{{ detail.description }}</p>
            </div>
          </div>
          
          <div class="info-item">
            <div class="icon-wrap" v-html="icons.clock"></div>
            <div class="info-text">
              <span class="i-title">LỊCH QUÉT ĐỊNH KỲ</span>
              <p>{{ detail.schedule }}</p>
            </div>
          </div>
        </div>

        <div class="action-buttons">
          <button class="btn-primary w-full"><span v-html="icons.repeat"></span> Quét lại ngay bây giờ</button>
          <div class="flex-row gap-8">
            <button class="btn-outline w-full">Tải dữ liệu (.csv)</button>
            <button class="btn-outline btn-danger-icon" aria-label="Delete">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
            </button>
          </div>
        </div>
      </div>
      <div v-else class="tab-placeholder">
        Nội dung tab "{{ activeTab }}" đang cập nhật...
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-detail-card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.cover-photo {
  height: 120px;
  background-size: cover;
  background-position: center;
  position: relative;
  background-color: #1e3a8a; /* Fallback */
}

.btn-edit-cover {
  position: absolute;
  top: 16px;
  right: 16px;
  background: rgba(255,255,255,0.2);
  border: none;
  color: white;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  backdrop-filter: blur(4px);
}

.btn-edit-cover:hover { background: rgba(255,255,255,0.3); }

.profile-header-info {
  padding: 0 24px;
  position: relative;
  margin-top: -40px; /* Overlap avatar */
  margin-bottom: 20px;
}

.avatar-wrap {
  width: 80px;
  height: 80px;
  border: 4px solid var(--c-surface);
  border-radius: 12px;
  overflow: hidden;
  background: white;
  margin-bottom: 12px;
  box-shadow: var(--shadow-sm);
}

.avatar-wrap img { width: 100%; height: 100%; object-fit: cover; }

.title-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 4px;
}

.title-row h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: var(--c-text-title);
}

.badge-status {
  background: #dcfce7;
  color: #166534;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 700;
}

.handle-link {
  font-size: 13px;
  color: var(--c-primary);
  text-decoration: none;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.handle-link:hover { text-decoration: underline; }

.card-body { padding: 0 24px 24px 24px; }

.metrics-grid {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
}

.metric-box {
  flex: 1;
  background: #f9fafb;
  border-radius: var(--radius-sm);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.m-label { font-size: 11px; font-weight: 600; color: var(--c-text-muted); letter-spacing: 0.5px; }
.m-val { font-size: 20px; font-weight: 700; color: var(--c-text-title); }

.info-list { display: flex; flex-direction: column; gap: 20px; margin-bottom: 32px; }
.info-item { display: flex; gap: 12px; align-items: flex-start; }
.icon-wrap { color: var(--c-text-muted); flex-shrink: 0; margin-top: 2px; }
.i-title { font-size: 11px; font-weight: 700; color: var(--c-text-title); display: block; margin-bottom: 4px; letter-spacing: 0.5px; }
.info-text p { margin: 0; font-size: 13px; color: var(--c-text-main); line-height: 1.5; }

.action-buttons { display: flex; flex-direction: column; gap: 12px; }
.w-full { width: 100%; }
.flex-row { display: flex; width: 100%; }
.gap-8 { gap: 8px; }

.btn-primary {
  background: var(--c-primary);
  color: white;
  border: none;
  padding: 12px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: opacity 0.2s;
}
.btn-primary:hover { opacity: 0.9; }

.btn-outline {
  background: white;
  border: 1px solid var(--c-border);
  color: var(--c-text-title);
  padding: 12px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-outline:hover { background: #f9fafb; }
.btn-danger-icon { color: var(--c-danger); padding: 12px 14px; flex-shrink: 0; }
.btn-danger-icon:hover { background: #fee2e2; border-color: #fca5a5; }

.tab-placeholder {
  padding: 40px 0;
  text-align: center;
  color: var(--c-text-muted);
  font-size: 14px;
  background: #f9fafb;
  border-radius: var(--radius-sm);
}
</style>
