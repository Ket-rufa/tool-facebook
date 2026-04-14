<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  profile: any
  isSelected: boolean
}>()

const statusClass = computed(() => {
  if (props.profile.statusType === 'success') return 'text-success'
  if (props.profile.statusType === 'danger') return 'text-danger'
  return 'text-muted'
})

// Generate a dummy avatar color based on ID
const avatarUrl = computed(() => {
  const bg = ['e0f2fe', 'f3f4f6', 'fef3c7'][props.profile.id % 3]
  const tr = ['0369a1', '1f2937', 'b45309'][props.profile.id % 3]
  const abbr = props.profile.name.substring(0,2).toUpperCase()
  return `https://ui-avatars.com/api/?name=${abbr}&background=${bg}&color=${tr}&size=48`
})
</script>

<template>
  <div class="profile-item" :class="{ 'selected': isSelected }">
    <div class="col-main">
      <img :src="avatarUrl" class="avatar" alt="avatar" />
      <div class="info">
        <div class="name">{{ profile.name }}</div>
        <div class="handle">{{ profile.handle }}</div>
      </div>
    </div>
    <div class="col-followers">
      {{ profile.followers }}
    </div>
    <div class="col-scan">
      <div class="scan-time">{{ profile.lastScanTime }}</div>
      <div class="scan-status" :class="statusClass">{{ profile.lastScanStatus }}</div>
    </div>
  </div>
</template>

<style scoped>
.profile-item {
  display: flex;
  padding: 24px;
  border-bottom: 1px solid var(--c-border);
  cursor: pointer;
  background: var(--c-surface);
  transition: all 0.2s;
  border-left: 3px solid transparent;
}

.profile-item:hover {
  background: #f9fafb;
}

.profile-item.selected {
  background: #eff6ff;
  border-left-color: var(--c-primary);
}

.col-main {
  flex: 2;
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  object-fit: cover;
}

.info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name {
  font-size: 14px;
  font-weight: 700;
  color: var(--c-text-title);
}

.handle {
  font-size: 12px;
  color: var(--c-text-muted);
}

.col-followers {
  flex: 1;
  display: flex;
  align-items: center;
  font-size: 14px;
  color: var(--c-text-main);
  font-weight: 500;
}

.col-scan {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 6px;
}

.scan-time {
  font-size: 13px;
  color: var(--c-text-main);
  white-space: pre-line;
  line-height: 1.4;
}

.scan-status {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
}

.text-success { color: var(--c-success); }
.text-danger { color: var(--c-danger); }
.text-muted { color: var(--c-text-muted); }
</style>
