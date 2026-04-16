<script setup lang="ts">
// Chức năng Quét bạn bè — đang phát triển
</script>

<template>
  <div class="friend-scraper-page">

    <div class="page-header">
      <div class="title-area">
        <h1>Quét bạn bè</h1>
        <p>Friend Scraping Workspace — Thu thập danh sách bạn bè từ tài khoản Facebook.</p>
      </div>
    </div>

    <!-- Coming Soon Card -->
    <div class="coming-soon-card">
      <div class="cs-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" viewBox="0 0 24 24" fill="none"
          stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
          <circle cx="9" cy="7" r="4" />
          <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
          <path d="M16 3.13a4 4 0 0 1 0 7.75" />
        </svg>
      </div>
      <div class="cs-content">
        <h2>Chức năng đang được phát triển</h2>
        <p>Tính năng <strong>Quét bạn bè</strong> sẽ cho phép bạn thu thập toàn bộ danh sách bạn bè từ một tài khoản Facebook và xuất ra file. Hiện tại chức năng đang được hoàn thiện.</p>
        <div class="feature-tags">
          <span class="tag">Thu thập UID</span>
          <span class="tag">Tên hiển thị</span>
          <span class="tag">Xuất CSV / JSON</span>
          <span class="tag">Lọc & Sắp xếp</span>
        </div>
      </div>
      <div class="cs-badge">Sắp ra mắt</div>
    </div>

    <!-- Placeholder config skeleton -->
    <div class="card skeleton-card" aria-hidden="true">
      <div class="card-header">Thiết lập Quét bạn bè</div>
      <div class="card-body">
        <div class="skeleton-row">
          <div class="skeleton w-40"></div>
          <div class="skeleton w-full"></div>
        </div>
        <div class="skeleton-row">
          <div class="skeleton w-40"></div>
          <div class="skeleton w-60"></div>
        </div>
        <div class="skeleton-row">
          <div class="skeleton btn-shape"></div>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
.friend-scraper-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  background-color: var(--c-background, #f8fafc);
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 800;
  color: var(--c-text-title, #111827);
}

.page-header p {
  color: var(--c-text-muted, #6b7280);
  margin-top: 4px;
  font-size: 14px;
}

/* ── Coming Soon Card ─────────────────────────── */
.coming-soon-card {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 24px;
  background: linear-gradient(135deg, #eff6ff 0%, #f0fdf4 100%);
  border: 1px solid #bfdbfe;
  border-radius: 12px;
  padding: 32px 28px;
  overflow: hidden;
}

.coming-soon-card::before {
  content: '';
  position: absolute;
  top: -40px; right: -40px;
  width: 180px; height: 180px;
  background: radial-gradient(circle, #bfdbfe 0%, transparent 70%);
  opacity: 0.4;
  pointer-events: none;
}

.cs-icon {
  flex-shrink: 0;
  color: var(--c-primary, #2563eb);
  opacity: 0.85;
}

.cs-content {
  flex: 1;
}

.cs-content h2 {
  margin: 0 0 10px;
  font-size: 18px;
  font-weight: 700;
  color: var(--c-text-title, #1e3a5f);
}

.cs-content p {
  margin: 0 0 16px;
  font-size: 14px;
  color: var(--c-text-main, #374151);
  line-height: 1.6;
}

.feature-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  padding: 4px 12px;
  background: white;
  border: 1px solid #bfdbfe;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  color: var(--c-primary, #2563eb);
}

.cs-badge {
  position: absolute;
  top: 16px; right: 20px;
  background: #f59e0b;
  color: white;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 4px 10px;
  border-radius: 20px;
}

/* ── Skeleton Placeholder ─────────────────────── */
.card {
  background: #fff;
  border: 1px solid var(--c-border, #e5e7eb);
  border-radius: 8px;
  overflow: hidden;
}

.card-header {
  padding: 14px 16px;
  border-bottom: 1px solid var(--c-border, #e5e7eb);
  font-weight: 600;
  font-size: 14px;
  background: #f9fafb;
  color: var(--c-text-muted, #6b7280);
}

.card-body {
  padding: 20px 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.skeleton-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.skeleton {
  background: linear-gradient(90deg, #f3f4f6 25%, #e5e7eb 50%, #f3f4f6 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 6px;
  height: 36px;
}

.w-40  { width: 120px; flex-shrink: 0; }
.w-60  { width: 60%; }
.w-full { flex: 1; }
.btn-shape { width: 140px; height: 40px; border-radius: 8px; }

.skeleton-card { opacity: 0.55; }

@keyframes shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
