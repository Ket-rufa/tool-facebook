<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  isLoading: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', payload: { postId: string, dryRun: boolean }): void
}>()

const postId = ref('')
const dryRun = ref(true)
const validationError = ref('')

const handleSubmit = () => {
  validationError.value = ''
  
  if (!postId.value.trim()) {
    validationError.value = 'Vui lòng nhập mã bài viết (Post ID)'
    return
  }

  emit('submit', {
    postId: postId.value.trim(),
    dryRun: dryRun.value
  })
}
</script>

<template>
  <div class="card">
    <h2 class="card-title">Cấu hình giả lập Like</h2>
    
    <div class="form-group">
      <label class="form-label">Mã bài viết (Post ID)</label>
      <input
        v-model="postId"
        type="text"
        class="form-input"
        placeholder="VD: FB_982129"
        :disabled="isLoading"
        @keyup.enter="handleSubmit"
      />
      <div v-if="validationError" class="error-msg">{{ validationError }}</div>
    </div>

    <div class="form-group mode-selection">
      <label class="form-label">Chế độ chạy</label>
      <div class="segmented-control">
        <div 
          class="segment-bg" 
          :class="{ 'is-right': !dryRun }"
        ></div>
        <button
          type="button"
          class="segment-btn"
          :class="{ 'active text-primary': dryRun, 'text-muted': !dryRun }"
          @click="dryRun = true"
          :disabled="isLoading"
        >
          Dry Run
        </button>
        <button
          type="button"
          class="segment-btn"
          :class="{ 'active font-bold': !dryRun, 'text-muted': dryRun }"
          @click="dryRun = false"
          :disabled="isLoading"
        >
          Real Run
        </button>
      </div>
      <p class="mode-desc" v-if="dryRun">
        * Dry Run: Mô phỏng hành động Like, KHÔNG gửi request thật lên hệ thống đích.
      </p>
      <p class="mode-desc warning" v-else>
        * Real Run: Cảnh báo! Sẽ gửi request đẩy Like thật nếu ứng dụng hỗ trợ.
      </p>
    </div>

    <div class="card-footer">
      <button
        type="button"
        @click="handleSubmit"
        :disabled="isLoading"
        class="btn btn-primary"
      >
        <span class="spinner" v-if="isLoading"></span>
        {{ isLoading ? 'Đang gửi...' : 'Test Like' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: var(--c-surface);
  border: 1px solid var(--c-border);
  border-radius: var(--radius-md);
  padding: 24px;
  box-shadow: var(--shadow-sm);
}

.card-title {
  margin: 0 0 20px;
  font-size: 16px;
  font-weight: 700;
  color: var(--c-text-title);
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--c-text-main);
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  font-size: 14px;
  color: var(--c-text-title);
  background: var(--c-bg);
  outline: none;
  font-family: monospace;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: var(--c-primary);
  background: var(--c-surface);
}

.form-input:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error-msg {
  margin-top: 6px;
  font-size: 12px;
  color: var(--c-danger);
}

.mode-selection {
  margin-top: 24px;
}

.segmented-control {
  display: inline-flex;
  position: relative;
  background: var(--c-border-light);
  padding: 4px;
  border-radius: var(--radius-sm);
  width: max-content;
}

.segment-bg {
  position: absolute;
  top: 4px;
  left: 4px;
  bottom: 4px;
  width: calc(50% - 4px);
  background: var(--c-surface);
  border-radius: calc(var(--radius-sm) - 2px);
  box-shadow: var(--shadow-sm);
  transition: transform 0.2s ease;
}

.segment-bg.is-right {
  transform: translateX(100%);
}

.segment-btn {
  position: relative;
  z-index: 1;
  padding: 6px 16px;
  min-width: 100px;
  border: none;
  background: transparent;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border-radius: calc(var(--radius-sm) - 2px);
}

.text-primary { color: var(--c-primary); }
.text-muted { color: var(--c-text-muted); font-weight: 500;}
.font-bold { color: var(--c-text-title); font-weight: 700; }

.mode-desc {
  margin-top: 10px;
  font-size: 12px;
  color: var(--c-text-muted);
}

.mode-desc.warning {
  color: var(--c-warning);
  font-weight: 600;
}

.card-footer {
  margin-top: 28px;
  padding-top: 20px;
  border-top: 1px solid var(--c-border-light);
  display: flex;
  justify-content: flex-end;
}

.btn {
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 600;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: none;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary { 
  background: var(--c-primary); 
  color: white; 
}

.btn-primary:not(:disabled):hover { 
  background: var(--c-primary-hover); 
}

/* Spinner */
.spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
