<script setup lang="ts">
import { ref } from 'vue'
import { settingsData as data } from '../data/settingsMockData'
import SettingsSection from '../components/SettingsSection.vue'
import OptionRow from '../components/OptionRow.vue'
import SegmentedControl from '../components/SegmentedControl.vue'
import { icons } from '../utils/icons'

// Mocks local states for UI toggling
const retention = ref(data.storage.retention)
const pureWhite = ref(data.uiOptions.pureWhite)
const checkedLogLevel = ref(data.logLevel)
</script>

<template>
  <div class="settings-page">
    <div class="page-header">
      <div class="title-area">
        <h1 class="page-title">{{ data.header.title }}</h1>
        <p class="page-desc">{{ data.header.description }}</p>
      </div>

      <div class="header-actions">
        <!-- Mock profile avatar like picture -->
        <div class="avatar-mock">
          <img src="https://ui-avatars.com/api/?name=Admin+User&background=eff6ff&color=0f62fe" class="avatar-img" />
        </div>
      </div>
    </div>

    <div class="settings-layout">
      <!-- Mini Sidebar -->
      <div class="s-sidebar">
        <nav class="s-nav">
          <a v-for="nav in data.sectionsNav" :key="nav.id" href="#" class="s-nav-item" :class="{ active: nav.id === 'crawl' }">
            <span class="s-icon" v-html="icons[nav.icon as keyof typeof icons]"></span>
            {{ nav.label }}
          </a>
        </nav>
      </div>

      <!-- Main Config Area -->
      <div class="s-content">
        <!-- Section 1 -->
        <SettingsSection title="CẤU HÌNH CRAWL" icon="checkBadge">
          <div class="row-2-col">
            <OptionRow label="Chính sách Retry" description="Xác định số lần thử lại và chiến lược trì hoãn khi gặp lỗi kết nối.">
              <select class="form-select">
                <option>{{ data.crawl.retryPolicy }}</option>
              </select>
              <label class="form-check">
                <input type="checkbox" :checked="data.crawl.useExponentialBackoff" />
                Sử dụng chiến lược Exponential Backoff
              </label>
            </OptionRow>

            <OptionRow label="Số luồng đồng thời" description="Cấu hình song song hóa để tối ưu hiệu suất xử lý dữ liệu.">
              <div class="slider-control">
                <input type="range" min="1" max="32" :value="data.crawl.concurrentThreads" class="form-range" />
                <span class="range-val">{{ data.crawl.concurrentThreads }}</span>
              </div>
              <div class="slider-hint">Khuyến nghị 8-16 luồng để tránh bị giới hạn băng thông.</div>
            </OptionRow>
          </div>
        </SettingsSection>

        <!-- Section 2 -->
        <SettingsSection title="LƯU TRỮ & DỮ LIỆU" icon="list">
          <OptionRow label="Vị trí lưu trữ dữ liệu" description="Đường dẫn thư mục cục bộ dùng để lưu trữ dữ liệu thô và tệp tin phương tiện." layout="vertical">
            <div class="input-group">
              <div class="path-input">
                <input type="text" :value="data.storage.path" readonly />
                <span class="folder-icon" v-html="icons.folder"></span>
              </div>
              <button class="btn btn-outline" v-html="icons.folder + ' Duyệt...'"></button>
            </div>
          </OptionRow>

          <OptionRow label="Thời gian lưu trữ dữ liệu" description="Chính sách tự động dọn dẹp (Retention Policy) để tối ưu dung lượng ổ đĩa." layout="vertical">
            <SegmentedControl 
              :options="data.storage.retentionOptions" 
              v-model="retention" 
              style="margin-top: 12px; max-width: 500px;"
            />
          </OptionRow>
        </SettingsSection>

        <!-- 2 Cols Grid for Section 3 & 4 -->
        <div class="row-2-col-wide">
          <SettingsSection title="MỨC ĐỘ LOG" icon="list">
            <div class="radio-list">
              <label class="radio-item">
                <input type="radio" value="error" name="loglevel" :checked="checkedLogLevel === 'error'" @change="checkedLogLevel = 'error'" />
                <div class="r-content">
                  <div class="r-title">Chỉ lỗi (Error)</div>
                  <div class="r-desc">Chỉ ghi nhận các lỗi làm dừng hệ thống</div>
                </div>
              </label>
              
              <label class="radio-item" :class="{ selected: checkedLogLevel === 'info' }">
                <input type="radio" value="info" name="loglevel" :checked="checkedLogLevel === 'info'" @change="checkedLogLevel = 'info'" />
                <div class="r-content">
                  <div class="r-title">Thông tin (Info)</div>
                  <div class="r-desc">Ghi nhận các sự kiện hoạt động chính</div>
                </div>
              </label>

              <label class="radio-item">
                <input type="radio" value="debug" name="loglevel" :checked="checkedLogLevel === 'debug'" @change="checkedLogLevel = 'debug'" />
                <div class="r-content">
                  <div class="r-title">Chi tiết (Debug/Verbose)</div>
                  <div class="r-desc">Mô tả chi tiết từng bước (Tốn tài nguyên)</div>
                </div>
              </label>
            </div>
          </SettingsSection>

          <SettingsSection title="TÙY CHỌN HIỂN THỊ" icon="checkBadge">
            <OptionRow label="Mật độ hiển thị" description="Điều chỉnh khoảng cách các thành phần" layout="horizontal">
              <select class="form-select-sm">
                <option>{{ data.uiOptions.density }}</option>
              </select>
            </OptionRow>

            <OptionRow label="Kích thước phông chữ" description="Kích thước chữ mặc định hệ thống" layout="horizontal" style="margin-top: 24px;">
              <select class="form-select-sm">
                <option>{{ data.uiOptions.fontSize }}</option>
              </select>
            </OptionRow>

            <OptionRow label="Giao diện sạch (Pure White)" description="Luôn sử dụng tông màu sáng tối giản" layout="horizontal" style="margin-top: 24px;">
              <!-- Simple Switch UI -->
              <label class="toggle-switch">
                <input type="checkbox" v-model="pureWhite" />
                <span class="slider round"></span>
              </label>
            </OptionRow>
          </SettingsSection>
        </div>

        <!-- Section 5 -->
        <div class="version-banner">
          <div class="vb-icon text-primary" v-html="icons.checkBadge"></div>
          <div class="vb-content">
            <div class="vb-title">Thông tin phiên bản</div>
            <div class="vb-desc">Hệ thống đang chạy trên phiên bản: <strong>{{ data.version.current }}</strong></div>
          </div>
          <div class="vb-action">
            <div class="vb-time">Cập nhật lần cuối: {{ data.version.lastUpdate }}</div>
            <button class="btn btn-outline float-right"><span v-html="icons.refreshCcw"></span> Kiểm tra cập nhật</button>
          </div>
        </div>

        <!-- Action Footer -->
        <div class="page-footer">
          <button class="btn btn-link">Khôi phục mặc định</button>
          <button class="btn btn-primary">Lưu thay đổi</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  gap: 32px;
  max-width: 1200px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--c-border);
}

.page-title {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 800;
  color: var(--c-text-title);
}

.page-desc {
  margin: 0;
  font-size: 15px;
  color: var(--c-text-muted);
}

.avatar-mock img {
  width: 40px; height: 40px; border-radius: 50%;
}

.settings-layout {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 32px;
}

.s-nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.s-nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  text-decoration: none;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
  border-radius: var(--radius-sm);
  transition: all 0.2s;
}

.s-nav-item:hover {
  background: #f3f4f6;
}

.s-nav-item.active {
  background: #eff6ff;
  color: var(--c-primary);
}

.s-icon { display: flex; }
.s-icon :deep(svg) { width: 16px; height: 16px; }

.s-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.row-2-col {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
}

.row-2-col-wide {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  align-items: stretch;
}

/* Forms */
.form-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  background: #f9fafb;
  font-size: 13px;
  color: var(--c-text-main);
  outline: none;
  margin-top: 12px;
  margin-bottom: 12px;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
}

.form-select-sm {
  padding: 6px 32px 6px 12px;
  border: 1px solid var(--c-border);
  border-radius: 4px;
  background: white;
  font-size: 13px;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%236b7280' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 8px center;
}

.form-check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
}

/* Custom Slider Input */
.slider-control {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
  margin-bottom: 8px;
}
.form-range {
  flex: 1;
  accent-color: var(--c-primary);
}
.range-val {
  background: var(--c-primary);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 700;
  min-width: 32px;
  text-align: center;
}
.slider-hint { font-size: 11px; color: var(--c-text-muted); font-style: italic; }


/* Input Group */
.input-group {
  display: flex;
  gap: 12px;
  margin-top: 12px;
  width: 100%;
}
.path-input {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
}
.path-input input {
  width: 100%;
  padding: 10px 12px 10px 40px;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  background: #f9fafb;
  font-size: 13px;
  font-family: monospace;
  outline: none;
  color: var(--c-text-main);
}
.folder-icon {
  position: absolute;
  left: 12px;
  color: var(--c-text-muted);
  display: flex;
}
.folder-icon :deep(svg) { width: 16px; height: 16px; }

/* Radio List */
.radio-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.radio-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px;
  border: 1px solid transparent; /* default */
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;
}
.radio-item.selected {
  background: #eff6ff;
  border-color: #bfdbfe;
}
.radio-item input[type="radio"] {
  margin-top: 2px;
  accent-color: var(--c-primary);
}
.r-content { display: flex; flex-direction: column; gap: 4px; }
.r-title { font-size: 13px; font-weight: 700; color: var(--c-text-title); }
.r-desc { font-size: 12px; color: var(--c-text-muted); }


/* Toggle Switch */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}

.toggle-switch input { 
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: #ccc;
  transition: .4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .4s;
}

input:checked + .slider {
  background-color: var(--c-primary);
}

input:checked + .slider:before {
  transform: translateX(20px);
}

.slider.round {
  border-radius: 24px;
}

.slider.round:before {
  border-radius: 50%;
}


/* Version Banner */
.version-banner {
  background: #f9fafb;
  border-radius: var(--radius-md);
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  border: 1px solid var(--c-border);
}

.vb-icon {
  width: 48px;
  height: 48px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-sm);
}

.vb-content { flex: 1; display: flex; flex-direction: column; gap: 4px; }
.vb-title { font-size: 13px; font-weight: 700; color: var(--c-text-title); }
.vb-desc { font-size: 13px; color: var(--c-text-main); }

.vb-action { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; }
.vb-time { font-size: 11px; color: var(--c-text-muted); font-style: italic; }


/* Buttons */
.btn {
  padding: 10px 20px;
  font-size: 13px;
  font-weight: 700;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: 1px solid transparent;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
}

.btn-primary { background: #0f62fe; color: white; }
.btn-primary:hover { background: #0050e6; }

.btn-outline { background: white; border-color: var(--c-border); color: var(--c-text-title); }
.btn-outline:hover { background: #f9fafb; }

.btn-link { background: transparent; color: var(--c-text-muted); font-weight: 600; }
.btn-link:hover { color: var(--c-text-main); }

.btn :deep(svg) { width: 14px; height: 14px; }

/* Action Footer */
.page-footer {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 16px;
}
</style>
