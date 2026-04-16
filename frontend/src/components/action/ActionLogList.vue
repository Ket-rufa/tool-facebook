<script setup lang="ts">
import type { ActionLog } from '../../types/action'

defineProps<{
  logs: ActionLog[]
}>()

const formatTime = (timeStr: string) => {
  if (!timeStr) return ''
  return new Date(timeStr).toLocaleString('vi-VN', {
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'success': return 'st-success'
    case 'not_enabled': return 'st-warning'
    case 'validation_error':
    case 'session_invalid':
    case 'failed':
      return 'st-error'
    default: return 'st-info'
  }
}
</script>

<template>
  <div class="card">
    <div class="card-header">
      <h3 class="card-title">Action Logs</h3>
      <span class="count-badge">{{ logs.length }}</span>
    </div>
    
    <div class="card-body">
      <div v-if="logs.length === 0" class="empty-state">
        Chưa có bản ghi thử nghiệm nào.
      </div>
      
      <div v-else class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Thời gian</th>
              <th>Tài khoản</th>
              <th>Mã Bài Viết</th>
              <th>Hành Động</th>
              <th>Trạng Thái</th>
              <th>Message</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in logs" :key="log.id" class="log-row">
              <td class="cell-time">{{ formatTime(log.executed_at) }}</td>
              <td class="cell-account">
                <div class="acc-name-cell">{{ log.account_display_name || log.account_id || '—' }}</div>
                <div class="acc-id-cell">{{ log.account_id }}</div>
              </td>
              <td class="cell-id">
                <span class="id-text">{{ log.post_id || '-' }}</span>
              </td>
              <td class="cell-reaction capitalize">
                {{ log.action_type === 'comment' ? 'Comment' : (log.reaction_type || log.action_type) }}
              </td>
              <td>
                <div class="status-wrapper">
                  <span class="status-pill" :class="getStatusClass(log.status)">
                    {{ log.status }}
                  </span>
                  <span v-if="log.dry_run" class="dry-pill">Dry</span>
                </div>
              </td>
              <td class="cell-msg" :title="log.message">{{ log.message }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.card-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--c-border-light);
  background: #fbfbfc;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title {
  margin: 0;
  font-size: 14px;
  font-weight: 700;
  color: var(--c-text-title);
}

.count-badge {
  background: #e5e7eb;
  color: #4b5563;
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 12px;
}

.card-body {
  padding: 0;
}

.empty-state {
  padding: 40px;
  text-align: center;
  font-size: 13px;
  color: var(--c-text-muted);
  font-style: italic;
}

.table-container {
  max-height: 400px;
  overflow-y: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.data-table th {
  position: sticky;
  top: 0;
  background: #fbfbfc;
  padding: 12px 16px;
  font-size: 11px;
  text-transform: uppercase;
  font-weight: 600;
  color: var(--c-text-muted);
  border-bottom: 1px solid var(--c-border);
  z-index: 10;
}

.data-table td {
  padding: 12px 16px;
  font-size: 13px;
  border-bottom: 1px solid var(--c-border-light);
  color: var(--c-text-main);
}

.log-row {
  transition: background-color 0.2s;
}

.log-row:hover {
  background-color: #f9fafb;
}

.log-row:last-child td {
  border-bottom: none;
}

.cell-time {
  font-size: 12px;
  color: var(--c-text-muted);
  white-space: nowrap;
}

.id-text {
  font-family: monospace;
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  color: #374151;
  font-size: 12px;
}

.status-wrapper {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-pill {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 10px;
  text-transform: capitalize;
}

.st-success { background: #d1fae5; color: #065f46; }
.st-warning { background: #fef3c7; color: #92400e; }
.st-error { background: #fee2e2; color: #991b1b; }
.st-info { background: #e0e7ff; color: #3730a3; }

.dry-pill {
  font-size: 9px;
  font-weight: 800;
  color: var(--c-primary);
  border: 1px solid var(--c-primary);
  padding: 1px 4px;
  border-radius: 4px;
  text-transform: uppercase;
}

.cell-msg {
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--c-text-muted);
}

.acc-text {
  font-weight: 600;
  color: var(--c-primary);
  font-size: 12px;
}

.acc-name-cell {
  font-weight: 600;
  color: var(--c-text-title);
  font-size: 13px;
}

.acc-id-cell {
  font-size: 10px;
  color: var(--c-text-muted);
  font-family: monospace;
}

.capitalize {
  text-transform: capitalize;
}
</style>
