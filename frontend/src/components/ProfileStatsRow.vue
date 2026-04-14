<script setup lang="ts">
defineProps<{
  stats: {
    total: { value: string; sub: string }
    scan24h: { value: string; sub: string }
    newData: { value: string; unit: string; sub: string }
    status: { value: string; isGood: boolean }
  }
}>()
</script>

<template>
  <div class="stats-row">
    <!-- Card 1 -->
    <div class="stat-box">
      <div class="s-label">TỔNG HỒ SƠ</div>
      <div class="s-val-row">
        <span class="s-val">{{ stats.total.value }}</span>
        <span class="s-sub text-success">{{ stats.total.sub }}</span>
      </div>
    </div>
    
    <!-- Card 2 -->
    <div class="stat-box">
      <div class="s-label">LƯỢT QUÉT 24H</div>
      <div class="s-val-row">
        <span class="s-val">{{ stats.scan24h.value }}</span>
        <span class="s-sub text-primary">{{ stats.scan24h.sub }}</span>
      </div>
    </div>

    <!-- Card 3 -->
    <div class="stat-box">
      <div class="s-label">DỮ LIỆU MỚI</div>
      <div class="s-val-row">
        <span class="s-val">{{ stats.newData.value }}</span>
        <span class="s-unit">{{ stats.newData.unit }}</span>
        <span class="s-sub default">{{ stats.newData.sub }}</span>
      </div>
    </div>

    <!-- Card 4 (Blue primary card) -->
    <div class="stat-box box-primary">
      <div class="s-label-light">TRẠNG THÁI HỆ THỐNG</div>
      <div class="s-status">
        <div class="dot-green" v-if="stats.status.isGood"></div>
        <span>{{ stats.status.value }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-box {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-sm);
}

.stat-box.box-primary {
  background: var(--c-primary);
  border-color: var(--c-primary);
  color: white;
}

.s-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--c-text-title);
  margin-bottom: 16px;
  letter-spacing: 0.5px;
}

.s-label-light {
  font-size: 12px;
  font-weight: 600;
  color: rgba(255,255,255,0.8);
  margin-bottom: 16px;
  letter-spacing: 0.5px;
}

.s-val-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-top: auto;
}

.s-val {
  font-size: 32px;
  font-weight: 700;
  color: var(--c-text-title);
  line-height: 1;
}

.s-unit {
  font-size: 18px;
  font-weight: 600;
  color: var(--c-text-title);
}

.s-sub {
  font-size: 13px;
  font-weight: 600;
}
.s-sub.text-success { color: var(--c-success); }
.s-sub.text-primary { color: var(--c-primary); }
.s-sub.default { color: var(--c-text-muted); font-weight: 500; }

.s-status {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: auto;
  font-size: 20px;
  font-weight: 700;
}

.dot-green {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: #34d399;
}
</style>
