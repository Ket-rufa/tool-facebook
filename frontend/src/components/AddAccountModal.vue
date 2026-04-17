<script setup lang="ts">
import { ref } from 'vue'
import { icons } from '../utils/icons'
import {
  StartAccountAttachFlow,
  GetAttachFlowStatus,
  CompleteAccountAttachFlow,
  CancelAccountAttachFlow,
  ProcessCookieAttachFlow,
  LoginByRequest,
  GetCookieFromCredentials
} from '../../wailsjs/go/accounts/AccountService'

const emit = defineEmits(['close', 'complete'])

const currentStep = ref(1)
const isOpeningWindow = ref(false)
const flowId = ref('')
const flowPreview = ref<any>(null)
const errorMsg = ref('')
const editableName = ref('')  // Tên do user nhập/sửa
const attachMethod = ref('browser')
const cookieString = ref('')
const isParsingCookie = ref(false)

// Dành cho Request Login
const email = ref('')
const password = ref('')
const twoFactorKey = ref('')
const isRequestLoggingIn = ref(false)
const isExtractingCookie = ref(false)
const extractedCookie = ref('')

function handleStartFlow() {
  if (attachMethod.value === 'browser') {
    currentStep.value = 2
  } else if (attachMethod.value === 'cookie') {
    currentStep.value = 2.5
  } else {
    currentStep.value = 2.7 // Step dành cho Request Login
  }
}

async function handleCookieSubmit() {
  isParsingCookie.value = true
  errorMsg.value = ''
  try {
    const status = await ProcessCookieAttachFlow(cookieString.value)
    flowPreview.value = status.accountPreview || flowPreview.value
    flowId.value = status.flowId
    editableName.value = flowPreview.value?.displayName || ''
    isParsingCookie.value = false
    currentStep.value = 3
  } catch (e: any) {
    console.error('Lỗi phân tích Cookie:', e)
    const msg = typeof e === 'string' ? e : (e?.message || JSON.stringify(e))
    errorMsg.value = msg || 'Có lỗi xảy ra khi phân tích Cookie.'
    isParsingCookie.value = false
  }
}

async function handleExtractCookie() {
  if (!email.value || !password.value) {
    errorMsg.value = 'Vui lòng nhập Email và Mật khẩu để lấy Cookie.'
    return
  }
  isExtractingCookie.value = true
  errorMsg.value = ''
  extractedCookie.value = ''
  try {
    const cookie = await GetCookieFromCredentials(email.value, password.value, twoFactorKey.value)
    extractedCookie.value = cookie
    isExtractingCookie.value = false
  } catch (e: any) {
    errorMsg.value = e || 'Lỗi khi lấy Cookie. Vui lòng kiểm tra lại tài khoản.'
    isExtractingCookie.value = false
  }
}

function handleUseExtractedCookie() {
  cookieString.value = extractedCookie.value
  currentStep.value = 2.5 // Chuyển sang tab nhập Cookie
}

async function handleRequestLogin() {
  if (!email.value || !password.value) {
    errorMsg.value = 'Vui lòng nhập đầy đủ Email và Mật khẩu.'
    return
  }
  isRequestLoggingIn.value = true
  errorMsg.value = ''
  try {
    // Chúng ta vẫn có thể dùng LoginByRequest (vốn dùng Chrome Auto) 
    // hoặc chuyển hẳn sang flow Extract -> Cookie. 
    // Ở đây tôi giữ LoginByRequest như một phương án "Auto hoàn toàn" cho user.
    const resp = await LoginByRequest(email.value, password.value, twoFactorKey.value)
    flowId.value = resp.flowId
    
    // Poll tương tự như browser login
    let attempts = 0
    const maxAttempts = 100
    const poll = async () => {
      if (attempts++ >= maxAttempts) {
        errorMsg.value = 'Hết thời gian chờ đăng nhập.'
        isRequestLoggingIn.value = false
        return
      }
      const status = await GetAttachFlowStatus(flowId.value)
      flowPreview.value = status.accountPreview || flowPreview.value
      if (status.state === 'authenticated') {
        editableName.value = flowPreview.value?.displayName || ''
        isRequestLoggingIn.value = false
        currentStep.value = 3
      } else if (status.state === 'failed') {
        errorMsg.value = status.message || 'Đăng nhập thất bại.'
        isRequestLoggingIn.value = false
      } else {
        setTimeout(poll, 1000)
      }
    }
    await poll()
  } catch (e: any) {
    errorMsg.value = e?.message || 'Lỗi kết nối hệ thống.'
    isRequestLoggingIn.value = false
  }
}

async function handleOpenLoginWindow() {
  isOpeningWindow.value = true
  errorMsg.value = ''
  try {
    // [THẬT] Gọi backend Go để khởi tạo flow và mở browser
    const resp = await StartAccountAttachFlow()
    flowId.value = resp.flowId
    flowPreview.value = resp.accountPreview || null

    // Poll trạng thái để đảm bảo flow đã Authenticated
    let attempts = 0
    const maxAttempts = 600 // Tăng lên 600 lần (~8-10 phút) để user thong thả đăng nhập/2FA
    const poll = async () => {
      if (attempts++ >= maxAttempts) {
        errorMsg.value = 'Hết thời gian chờ. Vui lòng thử lại.'
        isOpeningWindow.value = false
        return
      }
      const status = await GetAttachFlowStatus(flowId.value)
      flowPreview.value = status.accountPreview || flowPreview.value
      if (status.state === 'authenticated' || status.state === 'completed') {
        // Pre-fill tên có thể sửa được
        editableName.value = flowPreview.value?.displayName || ''
        isOpeningWindow.value = false
        currentStep.value = 3
      } else if (status.state === 'failed' || status.state === 'cancelled') {
        errorMsg.value = status.message || 'Đã có lỗi xảy ra.'
        isOpeningWindow.value = false
      } else {
        setTimeout(poll, 800)
      }
    }
    await poll()
  } catch (e: any) {
    errorMsg.value = e?.message || 'Lỗi kết nối backend.'
    isOpeningWindow.value = false
  }
}

async function handleComplete() {
  if (!flowId.value) {
    emit('close')
    return
  }
  // Truyền tên user nhập thẳng vào backend để lưu đúng
  const finalName = editableName.value.trim()
  const updatedList = await CompleteAccountAttachFlow(flowId.value, finalName)
  emit('complete', updatedList)
}

async function handleCancel() {
  if (flowId.value) {
    await CancelAccountAttachFlow(flowId.value)
  }
  emit('close')
}
</script>

<template>
  <div class="modal-overlay" @click.self="handleCancel">
    <div class="modal-card">
      <div class="m-header">
        <h2 class="m-title">Thêm tài khoản mới</h2>
        <button class="m-close" @click="handleCancel" v-html="icons.x"></button>
      </div>
      
      <div class="m-body">
        <!-- Step UI -->
        <div class="step-indicator">
          <div class="s-dot" :class="{ active: currentStep >= 1 }">1</div>
          <div class="s-line" :class="{ active: currentStep >= 2 }"></div>
          <div class="s-dot" :class="{ active: currentStep >= 2 }">2</div>
          <div class="s-line" :class="{ active: currentStep >= 3 }"></div>
          <div class="s-dot" :class="{ active: currentStep >= 3 }">3</div>
        </div>

        <div class="step-content">
          <!-- Step 1: Giới thiệu -->
          <div v-if="currentStep === 1" class="step-pane">
            <h3>Chọn phương thức đính kèm</h3>
            
            <div class="method-cards">
              <label class="m-card" :class="{ 'active': attachMethod === 'browser' }">
                <input type="radio" value="browser" v-model="attachMethod" class="sr-only" />
                <span class="mc-icon" v-html="icons.layout"></span>
                <div class="mc-content">
                  <div class="mc-title">Mở Cửa Sổ Trình Duyệt</div>
                  <div class="mc-desc">Đăng nhập trực quan qua cửa sổ mở riêng biệt. Ứng dụng không lưu mật khẩu.</div>
                </div>
              </label>
              
              <label class="m-card" :class="{ 'active': attachMethod === 'cookie' }">
                <input type="radio" value="cookie" v-model="attachMethod" class="sr-only" />
                <span class="mc-icon" v-html="icons.code"></span>
                <div class="mc-content">
                  <div class="mc-title">Nhập Cookie (Thủ công)</div>
                  <div class="mc-desc">Dán chuỗi Cookie từ trình duyệt.</div>
                </div>
              </label>

              <label class="m-card" :class="{ 'active': attachMethod === 'fast' }">
                <input type="radio" value="fast" v-model="attachMethod" class="sr-only" />
                <span class="mc-icon" v-html="icons.command"></span>
                <div class="mc-content">
                  <div class="mc-title">Đăng nhập bằng Request (Nhanh)</div>
                  <div class="mc-desc">Nhập User/Pass/2FA để tool tự động đăng nhập ngầm.</div>
                </div>
              </label>
            </div>

            <div class="m-actions">
              <button class="btn btn-outline" @click="handleCancel">Hủy</button>
              <button class="btn btn-primary" @click="handleStartFlow">Tiếp theo <span class="icon" v-html="icons.arrowRight"></span></button>
            </div>
          </div>

          <!-- Step 2: Mở cửa sổ & chờ -->
          <div v-if="currentStep === 2" class="step-pane centered">
            <template v-if="!isOpeningWindow">
              <div class="s2-icon text-primary" v-html="icons.layout"></div>
              <h3>Mở cửa sổ trình duyệt an toàn</h3>
              <p>Bấm nút bên dưới để bắt đầu đăng nhập.</p>
              <button class="btn btn-primary mt-4" @click="handleOpenLoginWindow">
                <span class="icon" v-html="icons.externalLink"></span> Mở cửa sổ đăng nhập
              </button>
            </template>
            
            <template v-else>
              <div class="loader-spinner"></div>
              <h3>Đang chờ xác thực...</h3>
              <p class="text-sm text-muted">Vui lòng đăng nhập ở cửa sổ vừa mở. Popup này sẽ tự chuyển trạng thái khi bạn đăng nhập thành công.</p>
              <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
            </template>
          </div>

          <!-- Step 2.5: Nhập Cookie -->
          <div v-if="currentStep === 2.5" class="step-pane">
            <h3>Nhập Cookie Facebook</h3>
            <p>Dán chuỗi Cookie đầy đủ chứa <code>c_user</code> và <code>xs</code> lấy từ tiện ích trình duyệt của bạn (VD: J2TEAM Cookie / GetCookie).</p>
            <textarea 
              v-model="cookieString" 
              class="cookie-input" 
              rows="5" 
              placeholder="Ví dụ: sb=123...; c_user=1000...; xs=...; fr=...;"
              :disabled="isParsingCookie"
            ></textarea>
            
            <div class="alert alert-info">
              <span class="a-icon text-primary" v-html="icons.info"></span>
              Hệ thống sẽ dùng Cookie để lấy fb_dtsg. Cookie của bạn bảo mật 100% trong máy tĩnh.
            </div>

            <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
            
            <div class="m-actions">
              <button class="btn btn-outline" @click="currentStep = 1">Quay lại</button>
              <button class="btn btn-primary" @click="handleCookieSubmit" :disabled="!cookieString.trim() || isParsingCookie">
                <span class="spinner" v-if="isParsingCookie"></span> Xác minh Cookie
              </button>
            </div>
          </div>

          <!-- Step 2.7: Request Login Form -->
          <div v-if="currentStep === 2.7" class="step-pane">
            <h3>Đăng nhập qua HTTP Request</h3>
            <p>Nhập thông tin tài khoản Facebook. Tool sẽ tự động xử lý bảo mật.</p>
            
            <div class="form-group-incard">
              <label>Email / Số điện thoại / UID</label>
              <input v-model="email" type="text" class="name-input" placeholder="Ví dụ: 1000123456789" :disabled="isRequestLoggingIn" />
            </div>

            <div class="form-group-incard">
              <label>Mật khẩu</label>
              <input v-model="password" type="password" class="name-input" placeholder="••••••••" :disabled="isRequestLoggingIn" />
            </div>

            <div class="form-group-incard">
              <label>Khóa 2FA (Tùy chọn)</label>
              <div class="hint">Nếu nick có 2FA, hãy dán mã Secret Key vào đây để tool tự giải mã OTP.</div>
              <input v-model="twoFactorKey" type="text" class="name-input" placeholder="J3XW..." :disabled="isRequestLoggingIn" />
            </div>

            <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>

            <!-- Kết quả trích xuất Cookie -->
            <div v-if="extractedCookie" class="extracted-cookie-box mt-4">
              <div class="ec-header">
                <span class="text-success" v-html="icons.checkBadge"></span>
                <strong>Đã lấy được Cookie!</strong>
              </div>
              <textarea readonly class="cookie-input mt-2" style="font-size: 11px; height: 60px;">{{ extractedCookie }}</textarea>
              <button class="btn btn-primary mt-2 w-full" @click="handleUseExtractedCookie">
                🚀 Dùng Cookie này để đăng nhập
              </button>
            </div>

            <div class="m-actions">
              <button class="btn btn-outline" @click="currentStep = 1">Quay lại</button>
              
              <button class="btn btn-outline-primary" @click="handleExtractCookie" :disabled="isExtractingCookie || isRequestLoggingIn">
                <span class="spinner" v-if="isExtractingCookie"></span>
                {{ isExtractingCookie ? 'Đang trích xuất...' : '🍪 Lấy Cookie' }}
              </button>

              <button class="btn btn-primary" @click="handleRequestLogin" :disabled="isRequestLoggingIn" title="Tự động đăng nhập qua trình duyệt">
                <span class="spinner" v-if="isRequestLoggingIn"></span>
                {{ isRequestLoggingIn ? 'Đang xử lý...' : '🚀 Đăng nhập Auto' }}
              </button>
            </div>
          </div>

          <!-- Step 3: Hoàn tất & Xác nhận -->
          <div v-if="currentStep === 3" class="step-pane">
            <div class="success-banner">
              <span class="sb-icon text-success" v-html="icons.checkBadge"></span>
              Xác thực thành công!
            </div>
            
            <div class="confirm-card">
              <div class="cc-header">
                <img :src="flowPreview?.avatar || 'https://ui-avatars.com/api/?name=FB&background=dbeafe&color=1d4ed8'" alt="" class="cc-avatar" />
                <div class="cc-info" style="flex:1">
                  <div class="c-name">{{ editableName || flowPreview?.displayName || 'Tài khoản mới' }}</div>
                  <div class="c-meta">{{ flowPreview?.provider || 'Facebook' }} &bull; {{ flowPreview?.accountType || 'Profile' }}</div>
                </div>
              </div>
              
              <!-- Input sửa tên -->
              <div class="name-edit-row">
                <label class="name-edit-label">Đặt tên tài khoản</label>
                <div class="name-edit-hint">Nhập tên hiển thị trên Facebook hoặc đặt tên riêng</div>
                <input 
                  v-model="editableName" 
                  type="text" 
                  class="name-input"
                  placeholder="VD: Victoria Thompson / Clone 1 / Nick phụ"
                />
              </div>

              <div class="cc-row">
                <span class="lbl">Account ID:</span>
                <span class="val fw-code">{{ flowPreview?.accountId || '—' }}</span>
              </div>
              <div class="cc-row">
                <span class="lbl">Trạng thái:</span>
                <span class="session-badge badge-success"><span class="badge-dot"></span>Hoạt động</span>
              </div>
            </div>

            <div class="m-actions">
              <button class="btn btn-primary" @click="handleComplete">Xác nhận <span class="icon" v-html="icons.check"></span></button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(17, 24, 39, 0.4);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-card {
  background: white;
  width: 540px;
  border-radius: var(--radius-md);
  box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1), 0 10px 10px -5px rgba(0,0,0,0.04);
  border: 1px solid var(--c-border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.m-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: #f9fafb;
  border-bottom: 1px solid var(--c-border);
}

.m-title {
  margin: 0;
  font-size: 16px;
  font-weight: 800;
  color: var(--c-text-title);
}

.m-close {
  background: transparent;
  border: none;
  color: var(--c-text-muted);
  cursor: pointer;
  padding: 4px;
}
.m-close:hover { color: var(--c-text-main); }
.m-close :deep(svg) { width: 18px; height: 18px; }

.m-body {
  padding: 24px;
}

/* Steps */
.step-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 32px;
}

.s-dot {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #f3f4f6;
  color: #9ca3af;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  transition: all 0.3s;
}

.s-dot.active {
  background: var(--c-primary);
  color: white;
  box-shadow: 0 0 0 3px #dbeafe;
}

.s-line {
  height: 2px;
  width: 40px;
  background: #e5e7eb;
  margin: 0 8px;
  transition: all 0.3s;
}

.s-line.active { background: var(--c-primary); }

/* Step Pane */
.step-pane {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.step-pane.centered {
  align-items: center;
  text-align: center;
  padding: 24px 0;
}

.step-pane h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--c-text-title);
}

.step-pane p {
  margin: 0;
  font-size: 14px;
  color: var(--c-text-main);
  line-height: 1.5;
}

.alert {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  padding: 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 13px;
  color: #1e40af;
  line-height: 1.4;
  margin-top: 8px;
}

.text-primary { color: var(--c-primary); }
.text-success { color: #22c55e; }
.text-muted { color: var(--c-text-muted) !important; }
.text-sm { font-size: 12px; }
.mt-4 { margin-top: 16px; }

.s2-icon {
  width: 56px;
  height: 56px;
  background: #eff6ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}
.s2-icon :deep(svg) { width: 28px; height: 28px; stroke-width: 2; }

/* Loader */
.loader-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top-color: var(--c-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}
@keyframes spin { to { transform: rotate(360deg); } }

/* Success Step */
.success-banner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 700;
  color: #166534;
  background: #dcfce7;
  padding: 12px;
  border-radius: var(--radius-sm);
  margin-bottom: 8px;
}
.sb-icon :deep(svg) { width: 20px; height: 20px; }

.confirm-card {
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: #f9fafb;
}

.cc-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--c-border);
}

.cc-avatar { width: 44px; height: 44px; border-radius: 50%; border: 1px solid var(--c-border); }
.c-name { font-weight: 700; color: var(--c-text-title); font-size: 15px; }
.c-meta { color: var(--c-text-muted); font-size: 12px; margin-top: 2px; }

.cc-row { display: flex; justify-content: space-between; align-items: center; font-size: 13px; }
.cc-row .lbl { color: var(--c-text-muted); font-weight: 600; }
.cc-row .val { color: var(--c-text-main); font-weight: 500; }
.fw-code { font-family: monospace; background: white; padding: 2px 6px; border: 1px solid var(--c-border); border-radius: 4px; }

/* Status Badge mini */
.session-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 700;
  padding: 4px 8px;
  border-radius: 12px;
  background: #dcfce7; color: #166534;
  text-transform: uppercase;
}
.badge-dot { width: 6px; height: 6px; border-radius: 50%; background: #22c55e; }

/* Actions */
.m-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 40px;
  padding: 0 20px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.btn-primary { background: var(--c-primary); color: white; }
.btn-primary:hover { background: var(--c-primary-hover); }

.btn-outline { background: white; border-color: var(--c-border); color: var(--c-text-title); }
.btn-outline:hover { background: #f9fafb; }

.icon { display: flex; }
.icon :deep(svg) { width: 14px; height: 14px; }

.error-msg {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  border-radius: 6px;
  padding: 8px 12px;
  font-size: 12px;
  margin-top: 8px;
}

/* Name edit */
.name-edit-row {
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: #f8faff;
  border: 1px solid #bfdbfe;
  border-radius: var(--radius-sm);
  padding: 14px 16px;
  margin-bottom: 8px;
}

.name-edit-label {
  font-size: 12px;
  font-weight: 700;
  color: var(--c-primary);
  text-transform: uppercase;
  letter-spacing: 0.4px;
}

.name-edit-hint {
  font-size: 11px;
  color: var(--c-text-muted);
}

.name-input, .cookie-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #93c5fd;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  color: var(--c-text-title);
  background: white;
  outline: none;
  box-sizing: border-box;
  transition: border-color 0.2s;
}

.cookie-input {
  font-family: monospace;
  font-weight: 500;
  font-size: 12px;
  resize: vertical;
  min-height: 80px;
}

.name-input:focus, .cookie-input:focus {
  border-color: var(--c-primary);
  box-shadow: 0 0 0 3px #dbeafe;
}

/* Method Cards */
.method-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 16px;
}

.m-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border: 1px solid var(--c-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  background: #f9fafb;
  transition: all 0.2s;
}

.m-card:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
}

.m-card.active {
  background: #eff6ff;
  border-color: var(--c-primary);
  box-shadow: 0 0 0 1px var(--c-primary);
}

.mc-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: white;
  border: 1px solid var(--c-border-light);
  color: var(--c-text-main);
}

.m-card.active .mc-icon {
  color: var(--c-primary);
  border-color: #bfdbfe;
  background: white;
}

.mc-content {
  flex: 1;
}

.mc-title {
  font-size: 14px;
  font-weight: 700;
  color: var(--c-text-title);
  margin-bottom: 4px;
}

.mc-desc {
  font-size: 12px;
  color: var(--c-text-muted);
}

.form-group-incard {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 4px;
}

.form-group-incard label {
  font-size: 12px;
  font-weight: 700;
  color: var(--c-text-main);
}

.form-group-incard .hint {
  font-size: 11px;
  color: var(--c-text-muted);
  margin-bottom: 2px;
}

.extracted-cookie-box {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: var(--radius-sm);
  padding: 12px;
}

.ec-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #166534;
}

.btn-outline-primary {
  background: white;
  border-color: var(--c-primary);
  color: var(--c-primary);
}
.btn-outline-primary:hover {
  background: #eff6ff;
}

.w-full {
  width: 100%;
}
</style>
