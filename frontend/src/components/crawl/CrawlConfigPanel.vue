<script setup lang="ts">
defineProps<{
  config: any
}>()

const availableGroups = ['Identity', 'Basic info', 'Location', 'Education', 'Work', 'Metrics', 'External links', 'Public posts']
</script>

<template>
  <div class="crawl-panel crawl-config-panel">
    <div class="panel-header">
      <h3>Zone B: Cấu hình quét</h3>
    </div>
    
    <div class="panel-body">
      <div class="config-section">
        <label>Preset cấu hình</label>
        <select class="form-select" :value="config.preset">
          <option value="Basic">Basic (Thông tin cơ bản)</option>
          <option value="Standard">Standard (Tiêu chuẩn)</option>
          <option value="Deep">Deep Public Scan (Chuyên sâu)</option>
        </select>
      </div>

      <div class="config-section">
        <label>Nhóm dữ liệu công khai cần quét</label>
        <div class="checkbox-grid">
          <label v-for="group in availableGroups" :key="group" class="checkbox-label">
            <input type="checkbox" :checked="config.fieldGroups.includes(group)" />
            <span>{{ group }}</span>
          </label>
        </div>
      </div>

      <div class="config-section runner-config">
        <div class="config-item">
          <label>Concurrency</label>
          <input type="number" class="form-input" :value="config.concurrency" />
        </div>
        <div class="config-item">
          <label>Delay (ms)</label>
          <input type="number" class="form-input" :value="config.delay" />
        </div>
        <div class="config-item">
          <label>Retry limit</label>
          <input type="number" class="form-input" :value="config.retry" />
        </div>
        <div class="config-item">
          <label>Max Posts</label>
          <input type="number" class="form-input" :value="config.maxPosts" />
        </div>
      </div>

      <div class="actions-footer">
        <button class="btn btn-run">Bắt đầu quét</button>
        <button class="btn btn-pause">Tạm dừng</button>
        <button class="btn btn-stop">Dừng</button>
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

.config-section label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.form-select, .form-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 13px;
  color: #111827;
}

.checkbox-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.checkbox-label {
  display: flex !important;
  align-items: center;
  gap: 8px;
  font-size: 13px !important;
  font-weight: 400 !important;
  color: #374151 !important;
  text-transform: none !important;
  letter-spacing: normal !important;
  margin: 0 !important;
}

.runner-config {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}
.config-item label { margin-bottom: 4px; }

.actions-footer {
  margin-top: auto;
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #f3f4f6;
}

.btn {
  padding: 8px 16px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 6px;
  cursor: pointer;
  border: transparent;
  flex: 1;
  text-align: center;
  transition: opacity 0.2s;
}
.btn:hover { opacity: 0.9; }

.btn-run { background: #10b981; color: white; }
.btn-pause { background: #f59e0b; color: white; }
.btn-stop { background: #ef4444; color: white; }
</style>
