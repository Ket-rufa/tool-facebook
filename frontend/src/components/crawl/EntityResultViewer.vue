<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  entities: any[],
  logs: any[]
}>()

const activeEntityId = ref(props.entities.length > 0 ? props.entities[0].entity_id : null)
const activeTab = ref('overview')

const activeEntity = computed(() => props.entities.find(e => e.entity_id === activeEntityId.value))
const entityLogs = computed(() => props.logs.filter(l => l.uid === activeEntity.value?.uid))

const selectEntity = (id: string) => {
  activeEntityId.value = id
  activeTab.value = 'overview'
}

const getMetricValue = (text: any, value: any) => {
  if (text) return text
  if (value !== null && value !== undefined) return value
  return '—'
}

const tabs = [
  { id: 'overview', label: 'Tổng quan' },
  { id: 'fields', label: 'Dữ liệu' },
  { id: 'raw', label: 'JSON Gốc' },
  { id: 'snapshots', label: 'Lịch sử' },
  { id: 'logs', label: 'Nhật ký' }
]
</script>

<template>
  <div class="crawl-panel entity-viewer-panel">
    <div class="viewer-layout">
      <!-- Left Pane: Entity List -->
      <div class="entity-list-pane">
        <div class="pane-header">DANH SÁCH ENTITY</div>
        <div class="entity-list">
          <div 
            v-for="entity in entities" 
            :key="entity.entity_id"
            :class="['entity-item', { active: activeEntityId === entity.entity_id }]"
            @click="selectEntity(entity.entity_id)"
          >
            <div class="entity-avatar">
              <img v-if="entity.avatar" :src="entity.avatar" alt="avatar" />
            </div>
            <div class="entity-info">
              <div class="entity-name">{{ entity.identity?.display_name || '—' }}</div>
              <div class="entity-uid">{{ entity.uid ? `UID: ${entity.uid}` : 'No UID' }}</div>
            </div>
            <div class="entity-status">
              <span :class="['status-dot', entity.scan_status]"></span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Pane: Detail Viewer -->
      <div class="entity-detail-pane" v-if="activeEntity">
        <div class="detail-header">
          <div class="header-info">
            <h2>{{ activeEntity.identity?.display_name || 'Thực thể không xác định' }}</h2>
            <div class="badges">
              <span class="badge">{{ activeEntity.entity_type }}</span>
              <span class="badge" v-if="activeEntity.uid">UID: {{ activeEntity.uid }}</span>
              <span :class="['badge status-type', activeEntity.scan_status]">{{ activeEntity.scan_status === 'success' ? 'Thành công' : activeEntity.scan_status }}</span>
            </div>
          </div>
        </div>

        <div class="tabs-bar">
          <button 
            v-for="tab in tabs" 
            :key="tab.id"
            :class="['tab-btn', { active: activeTab === tab.id }]"
            @click="activeTab = tab.id"
          >
            {{ tab.label }}
          </button>
        </div>

        <div class="tab-content">
          <!-- OVERVIEW TAB -->
          <div v-if="activeTab === 'overview'" class="overview-grid">
            <div class="info-card">
              <h4>Thông tin cơ bản</h4>
              <div class="info-row"><span class="label">Tên hiển thị:</span> <span class="val">{{ activeEntity.identity?.display_name || '—' }}</span></div>
              <div class="info-row"><span class="label">Tên khác:</span> <span class="val">{{ activeEntity.identity?.alternate_name || '—' }}</span></div>
              <div class="info-row"><span class="label">Tiểu sử:</span> <span class="val">{{ activeEntity.basic_info?.bio || '—' }}</span></div>
              <div class="info-row"><span class="label">Ngày sinh:</span> <span class="val">{{ activeEntity.basic_info?.birthday_text || '—' }}</span></div>
            </div>
            
            <div class="info-card">
              <h4>Vị trí & Học vấn</h4>
              <div class="info-row"><span class="label">Thành phố hiện tại:</span> <span class="val">{{ activeEntity.location?.current_city || '—' }}</span></div>
              <div class="info-row"><span class="label">Quê quán:</span> <span class="val">{{ activeEntity.location?.hometown || '—' }}</span></div>
              <div class="info-row"><span class="label">Học vấn:</span> 
                <span class="val">
                  <ul class="val-list" v-if="activeEntity.education?.length">
                    <li v-for="edu in activeEntity.education" :key="edu">{{ edu }}</li>
                  </ul>
                  <span v-else>—</span>
                </span>
              </div>
            </div>

            <div class="info-card">
              <h4>Chỉ số (Công khai)</h4>
              <div class="info-row"><span class="label">Người theo dõi:</span> <span class="val">{{ getMetricValue(activeEntity.metrics?.followers_text, activeEntity.metrics?.followers_value) }}</span></div>
              <div class="info-row"><span class="label">Đang theo dõi:</span> <span class="val">{{ getMetricValue(activeEntity.metrics?.following_text, activeEntity.metrics?.following_value) }}</span></div>
              <div class="info-row"><span class="label">Bạn bè:</span> <span class="val">{{ getMetricValue(activeEntity.metrics?.friends_count_text, activeEntity.metrics?.friends_count_value) }}</span></div>
            </div>

            <div class="info-card">
              <h4>Hệ thống</h4>
              <div class="info-row"><span class="label">Trạng thái dữ liệu:</span> <span class="val">{{ activeEntity.field_status || '—' }}</span></div>
              <div class="info-row"><span class="label">Lần quét cuối:</span> <span class="val">{{ activeEntity.last_updated || '—' }}</span></div>
              <div class="info-row"><span class="label">Canonical URL:</span> <span class="val"><a :href="activeEntity.canonical_url" target="_blank">{{ activeEntity.canonical_url || '—' }}</a></span></div>
            </div>
          </div>

          <!-- FIELDS TAB -->
          <div v-else-if="activeTab === 'fields'" class="fields-table-container">
             <table class="fields-table">
               <thead>
                 <tr>
                   <th>Trường</th>
                   <th>Nhóm</th>
                   <th>Giá trị</th>
                   <th>Trạng thái</th>
                 </tr>
               </thead>
               <tbody>
                 <tr><td>display_name</td><td>Identity</td><td>{{ activeEntity.identity?.display_name || '—' }}</td><td><span class="status-yes">Có</span></td></tr>
                 <tr><td>bio</td><td>Basic info</td><td>{{ activeEntity.basic_info?.bio || '—' }}</td><td><span v-if="activeEntity.basic_info?.bio" class="status-yes">Có</span><span v-else class="status-no">Rỗng</span></td></tr>
                 <tr><td>followers</td><td>Metrics</td><td>{{ getMetricValue(activeEntity.metrics?.followers_text, null) }}</td><td><span v-if="activeEntity.metrics?.followers_text" class="status-yes">Có</span><span v-else class="status-no">Rỗng</span></td></tr>
               </tbody>
             </table>
          </div>

          <!-- RAW JSON TAB -->
          <div v-else-if="activeTab === 'raw'" class="raw-json-container">
            <pre><code>{{ JSON.stringify(activeEntity, null, 2) }}</code></pre>
          </div>

          <!-- SNAPSHOTS TAB -->
          <div v-else-if="activeTab === 'snapshots'" class="snapshots-container">
            <div class="empty-state">Chưa có snapshot nào được lưu thay đổi.</div>
          </div>

          <!-- LOGS TAB -->
          <div v-else-if="activeTab === 'logs'" class="logs-container">
            <div v-if="entityLogs.length === 0" class="empty-state">Không có nhật ký nào cho entity này.</div>
            <div v-else class="entity-log-line" v-for="log in entityLogs" :key="log.id">
              <span class="time">[{{ log.time }}]</span>
              <span :class="['level', log.level]">{{ log.level.toUpperCase() }}</span>
              <span>{{ log.message }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="entity-detail-pane empty" v-else>
        <p>Chọn một Entity để xem cấu trúc chi tiết.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.crawl-panel {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  flex: 1;
}

.viewer-layout {
  display: flex;
  height: 100%;
}

.entity-list-pane {
  width: 280px;
  border-right: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  background: #f9fafb;
}

.pane-header {
  padding: 12px 16px;
  font-size: 11px;
  font-weight: 700;
  color: #6b7280;
  text-transform: uppercase;
  border-bottom: 1px solid #e5e7eb;
}

.entity-list {
  flex: 1;
  overflow-y: auto;
}

.entity-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid #f3f4f6;
  cursor: pointer;
  transition: background 0.2s;
}
.entity-item:hover { background: #f3f4f6; }
.entity-item.active { background: #eff6ff; border-left: 3px solid #3b82f6; padding-left: 13px; }

.entity-avatar { width: 32px; height: 32px; border-radius: 50%; overflow: hidden; background: #e5e7eb; }
.entity-avatar img { width: 100%; height: 100%; object-fit: cover; }

.entity-info { flex: 1; min-width: 0; }
.entity-name { font-size: 13px; font-weight: 600; color: #111827; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.entity-uid { font-size: 11px; color: #6b7280; font-family: monospace; }

.status-dot { width: 8px; height: 8px; border-radius: 50%; display: block; }
.status-dot.queued { background: #9ca3af; }
.status-dot.running { background: #3b82f6; }
.status-dot.success { background: #10b981; }
.status-dot.failed { background: #ef4444; }

.entity-detail-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: white;
}
.entity-detail-pane.empty { align-items: center; justify-content: center; color: #9ca3af; font-size: 14px; }

.detail-header { padding: 20px; border-bottom: 1px solid #e5e7eb; display: flex; justify-content: space-between; }
.header-info h2 { margin: 0 0 8px 0; font-size: 18px; color: #111827; }
.badges { display: flex; gap: 8px; }
.badge { font-size: 11px; background: #f3f4f6; color: #4b5563; padding: 2px 8px; border-radius: 4px; font-weight: 500; }
.badge.status-type.running { background: #dbeafe; color: #1e40af; }
.badge.status-type.queued { background: #f3f4f6; color: #4b5563; }

.tabs-bar { display: flex; border-bottom: 1px solid #e5e7eb; background: #f9fafb; padding: 0 12px; }
.tab-btn { background: none; border: none; padding: 12px 16px; font-size: 13px; font-weight: 500; color: #6b7280; cursor: pointer; border-bottom: 2px solid transparent; margin-bottom: -1px; }
.tab-btn:hover { color: #111827; }
.tab-btn.active { color: #2563eb; border-bottom-color: #2563eb; }

.tab-content { padding: 20px; flex: 1; overflow-y: auto; }

.overview-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20px; }
.info-card { background: #f9fafb; border: 1px solid #f3f4f6; border-radius: 8px; padding: 16px; }
.info-card h4 { margin: 0 0 12px 0; font-size: 13px; color: #374151; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px; }
.info-row { display: flex; margin-bottom: 8px; font-size: 13px; }
.info-row:last-child { margin-bottom: 0; }
.info-row .label { width: 120px; color: #6b7280; flex-shrink: 0; }
.info-row .val { color: #111827; font-weight: 500; }
.val-list { margin: 0; padding-left: 20px; }
.val-list li { margin-bottom: 4px; }
.info-row a { color: #2563eb; text-decoration: none; word-break: break-all; }
.info-row a:hover { text-decoration: underline; }

.fields-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.fields-table th, .fields-table td { padding: 12px; text-align: left; border-bottom: 1px solid #e5e7eb; }
.fields-table th { background: #f9fafb; font-weight: 600; color: #374151; }
.status-yes { color: #059669; font-weight: 500; }
.status-no { color: #9ca3af; }

.raw-json-container pre { background: #111827; color: #e5e7eb; padding: 16px; border-radius: 8px; font-size: 12px; overflow-x: auto; margin: 0; }
.empty-state { color: #6b7280; font-size: 13px; text-align: center; padding: 40px; }

.entity-log-line { font-family: monospace; font-size: 12px; margin-bottom: 8px; color: #374151; }
.entity-log-line .time { color: #6b7280; margin-right: 8px; }
.entity-log-line .level { font-weight: 700; margin-right: 8px; }
.entity-log-line .level.info { color: #2563eb; }
.entity-log-line .level.warning { color: #d97706; }
.entity-log-line .level.error { color: #dc2626; }
</style>
