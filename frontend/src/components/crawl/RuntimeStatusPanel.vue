<script setup lang="ts">
defineProps<{
  status: any
}>()
</script>

<template>
  <div class="crawl-panel runtime-status-panel">
    <div class="panel-header">
      <h3>Zone C: Tiến trình Runtime</h3>
    </div>
    
    <div class="panel-body">
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-label">Tổng mục tiêu</div>
          <div class="stat-value">{{ status.total }}</div>
        </div>
        <div class="stat-card">
          <div class="stat-label">Queued</div>
          <div class="stat-value">{{ status.queued }}</div>
        </div>
        <div class="stat-card running">
          <div class="stat-label">Running</div>
          <div class="stat-value">{{ status.running }}</div>
        </div>
        <div class="stat-card success">
          <div class="stat-label">Success</div>
          <div class="stat-value">{{ status.success }}</div>
        </div>
        <div class="stat-card failed">
          <div class="stat-label">Failed</div>
          <div class="stat-value">{{ status.failed }}</div>
        </div>
        <div class="stat-card partial">
          <div class="stat-label">Partial</div>
          <div class="stat-value">{{ status.partial }}</div>
        </div>
      </div>

      <div class="progress-section">
        <div class="progress-header">
          <span>Tiến độ tổng hợp</span>
          <span class="eta">ETA: {{ status.eta }}</span>
        </div>
        <div class="progress-bar-container">
          <div class="progress-fill" :style="`width: ${status.progress}%`"></div>
        </div>
        <div class="progress-text">{{ status.progress }}% hoàn thành</div>
      </div>

      <div class="current-target-card">
        <div class="card-title">TIẾN TRÌNH HIỆN TẠI</div>
        <div class="current-info target-name">Entity: <strong>{{ status.currentEntity }}</strong></div>
        <div class="current-info target-uid">UID: <code>{{ status.currentUid }}</code></div>
        <div class="current-info target-step">Bước hiện tại: <span class="step-badge">{{ status.currentStep }}</span></div>
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
}

.panel-header {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
}
.panel-header h3 { margin: 0; font-size: 14px; font-weight: 600; color: #374151; }

.panel-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex: 1;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.stat-card {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  padding: 12px;
  border-radius: 6px;
  text-align: center;
}
.stat-card.running .stat-value { color: #3b82f6; }
.stat-card.success .stat-value { color: #10b981; }
.stat-card.failed .stat-value { color: #ef4444; }
.stat-card.partial .stat-value { color: #f59e0b; }

.stat-label { font-size: 11px; color: #6b7280; font-weight: 500; text-transform: uppercase; margin-bottom: 4px; }
.stat-value { font-size: 20px; font-weight: 700; color: #111827; }

.progress-section { margin-top: 4px; }
.progress-header { display: flex; justify-content: space-between; font-size: 12px; font-weight: 600; color: #374151; margin-bottom: 8px; }
.eta { color: #6b7280; font-weight: 500; }
.progress-bar-container { background: #e5e7eb; height: 8px; border-radius: 4px; overflow: hidden; }
.progress-fill { background: #3b82f6; height: 100%; transition: width 0.3s; }
.progress-text { text-align: right; font-size: 11px; color: #6b7280; margin-top: 4px; }

.current-target-card {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  padding: 12px;
  border-radius: 6px;
  margin-top: auto;
}
.card-title { font-size: 11px; font-weight: 700; color: #1d4ed8; margin-bottom: 8px; letter-spacing: 0.5px; }
.current-info { font-size: 13px; color: #1e3a8a; margin-bottom: 4px; }
.current-info strong { font-weight: 600; }
.current-info code { background: white; padding: 2px 6px; border-radius: 4px; font-family: monospace; font-size: 12px;}
.step-badge {
  display: inline-block;
  background: #3b82f6;
  color: white;
  padding: 2px 8px;
  border-radius: 99px;
  font-size: 11px;
  font-weight: 500;
  margin-left: 4px;
}
</style>
