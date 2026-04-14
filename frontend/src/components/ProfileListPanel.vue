<script setup lang="ts">
import ProfileListItem from './ProfileListItem.vue'

defineProps<{
  profiles: any[]
  selectedId: number | string
}>()

const emit = defineEmits(['select'])
</script>

<template>
  <div class="profile-list-panel">
    <div class="panel-header">
      <h3 class="title">DANH SÁCH HỒ SƠ</h3>
      <a href="#" class="btn-export">Xuất báo cáo</a>
    </div>

    <div class="list-header">
      <div class="col-h">HỒ SƠ / TRANG</div>
      <div class="col-h">SỐ NGƯỜI THEO DÕI</div>
      <div class="col-h">LẦN QUÉT CUỐI</div>
    </div>

    <div class="list-body">
      <ProfileListItem 
        v-for="p in profiles" 
        :key="p.id" 
        :profile="p" 
        :is-selected="p.id === selectedId"
        @click="emit('select', p)"
      />
    </div>
  </div>
</template>

<style scoped>
.profile-list-panel {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid var(--c-border);
}

.title {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-title);
  letter-spacing: 0.5px;
}

.btn-export {
  font-size: 12px;
  font-weight: 600;
  color: var(--c-primary);
  text-decoration: none;
}
.btn-export:hover { text-decoration: underline; }

.list-header {
  display: flex;
  padding: 16px 24px 16px 48px; /* account for border-left space */
  background: #f9fafb;
  border-bottom: 1px solid var(--c-border);
}

.col-h {
  font-size: 11px;
  font-weight: 600;
  color: var(--c-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.col-h:first-child { flex: 2; }
.col-h:nth-child(2) { flex: 1; padding-left: 16px;} /* Align with items */
.col-h:last-child { flex: 1; padding-left: 16px; }

.list-body {
  display: flex;
  flex-direction: column;
}
</style>
