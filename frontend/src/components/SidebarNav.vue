<script setup lang="ts">
import { inject, Ref } from 'vue'
import { icons } from '../utils/icons'
import StorageUsageCard from './StorageUsageCard.vue'

const currentRoute = inject('currentRoute') as Ref<string>
const navigate = inject('navigate') as (route: string) => void

const navItems = [
  { id: 'overview', label: 'Tổng quan', icon: icons.grid },
  { id: 'tasks', label: 'Công việc quét', icon: icons.repeat },
  { id: 'profiles', label: 'Hồ sơ', icon: icons.users },
  { id: 'posts', label: 'Bài viết', icon: icons.fileText },
  { id: 'media', label: 'Phương tiện', icon: icons.image },
  { id: 'logs', label: 'Nhật ký', icon: icons.clock },
  { id: 'integrations', label: 'Tích hợp', icon: icons.plug },
  { id: 'action-like-test', label: 'Thử nghiệm Like', icon: icons.checkBadge },
  { id: 'settings', label: 'Cài đặt', icon: icons.settings },
]
</script>

<template>
  <aside class="sidebar">
    <div class="logo-area">
      <div class="logo-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"></path></svg>
      </div>
      <div class="logo-text">
        <h1>Facebook<br/>Manager</h1>
        <span>KIẾN TRÚC SỰ CHÍNH XÁC</span>
      </div>
    </div>

    <nav class="nav-menu">
      <a 
        v-for="item in navItems" 
        :key="item.id" 
        href="#" 
        class="nav-item"
        :class="{ active: currentRoute === item.id || (item.id === 'posts' && currentRoute === 'post-detail') }"
        @click.prevent="navigate(item.id)"
      >
        <span class="icon" v-html="item.icon"></span>
        <span class="label">{{ item.label }}</span>
      </a>
    </nav>
    
    <StorageUsageCard />

    <div class="user-profile">
      <img src="https://ui-avatars.com/api/?name=Admin+User&background=eff6ff&color=0f62fe" alt="Avatar" class="avatar" />
      <div class="user-info">
        <span class="name">Admin User</span>
        <span class="plan">Enterprise Plan</span>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 260px;
  background-color: var(--c-sidebar);
  border-right: 1px solid var(--c-border);
  display: flex;
  flex-direction: column;
}

.logo-area {
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background-color: var(--c-primary);
  color: white;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text h1 {
  font-size: 16px;
  font-weight: 700;
  line-height: 1.2;
  margin: 0;
  color: var(--c-primary);
}

.logo-text span {
  font-size: 10px;
  color: var(--c-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.nav-menu {
  flex: 1;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  text-decoration: none;
  color: #4b5563;
  border-radius: var(--radius-sm);
  font-weight: 500;
  font-size: 14px;
  transition: all 0.2s ease;
}

.nav-item:hover {
  background-color: #f3f4f6;
  color: var(--c-text-main);
}

.nav-item.active {
  background-color: #eff6ff;
  color: var(--c-primary);
  border-left: 3px solid var(--c-primary);
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
}

.icon {
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.7;
}

.nav-item.active .icon {
  opacity: 1;
}

.user-profile {
  margin: 20px;
  padding: 12px;
  background-color: #f3f4f6;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.name {
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
}

.plan {
  font-size: 11px;
  color: var(--c-text-muted);
}
</style>
