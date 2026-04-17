package accounts

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AttachFlowManager quản lý các phiên gắn tài khoản đang chạy
type AttachFlowManager struct {
	mu    sync.Mutex
	flows map[string]*AttachFlow
}

func NewAttachFlowManager() *AttachFlowManager {
	return &AttachFlowManager{
		flows: make(map[string]*AttachFlow),
	}
}

// Start tạo flow mới, trả về flowId
func (m *AttachFlowManager) Start() (*AttachFlow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	flow := &AttachFlow{
		FlowID:    "flow_" + uuid.New().String()[:8],
		State:     FlowCreated,
		Message:   "Flow đã khởi tạo. Đang chờ mở cửa sổ đăng nhập.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.flows[flow.FlowID] = flow
	return flow, nil
}

// SetWaitingAuth chuyển trạng thái sang chờ đăng nhập
func (m *AttachFlowManager) SetWaitingAuth(flowID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	flow, ok := m.flows[flowID]
	if !ok {
		return fmt.Errorf("không tìm thấy flow: %s", flowID)
	}
	flow.State = FlowWaitingAuth
	flow.Message = "Cửa sổ đăng nhập đang mở. Đang chờ xác thực..."
	flow.UpdatedAt = time.Now()
	return nil
}

// SetAuthenticated chuyển trạng thái sau khi đã đăng nhập thành công (có preview)
func (m *AttachFlowManager) SetAuthenticated(flowID string, preview *AccountProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	flow, ok := m.flows[flowID]
	if !ok {
		return fmt.Errorf("không tìm thấy flow: %s", flowID)
	}
	flow.State = FlowAuthenticated
	flow.Message = "Xác thực thành công. Đang chờ xác nhận từ người dùng."
	flow.AccountPreview = preview
	flow.UpdatedAt = time.Now()
	return nil
}

// SetFailed chuyển trạng thái khi flow gặp lỗi
func (m *AttachFlowManager) SetFailed(flowID string, message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	flow, ok := m.flows[flowID]
	if !ok {
		return fmt.Errorf("không tìm thấy flow: %s", flowID)
	}
	flow.State = FlowFailed
	flow.Message = message
	flow.UpdatedAt = time.Now()
	return nil
}

// Complete hoàn tất flow
func (m *AttachFlowManager) Complete(flowID string) (*AttachFlow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	flow, ok := m.flows[flowID]
	if !ok {
		return nil, fmt.Errorf("không tìm thấy flow: %s", flowID)
	}
	if flow.State != FlowAuthenticated {
		return nil, fmt.Errorf("flow chưa ở trạng thái authenticated: %s", flow.State)
	}
	flow.State = FlowCompleted
	flow.Message = "Gắn phiên hoàn tất."
	flow.UpdatedAt = time.Now()
	return flow, nil
}

// Cancel hủy flow
func (m *AttachFlowManager) Cancel(flowID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	flow, ok := m.flows[flowID]
	if !ok {
		return nil // không có cũng không cần báo lỗi khi cancel
	}
	flow.State = FlowCancelled
	flow.Message = "Flow đã bị hủy."
	flow.UpdatedAt = time.Now()
	return nil
}

// GetStatus trả về trạng thái của flow
func (m *AttachFlowManager) GetStatus(flowID string) (*AttachFlowStatusResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	flow, ok := m.flows[flowID]
	if !ok {
		return nil, fmt.Errorf("không tìm thấy flow: %s", flowID)
	}
	return &AttachFlowStatusResponse{
		FlowID:         flow.FlowID,
		State:          string(flow.State),
		Message:        flow.Message,
		AccountPreview: flow.AccountPreview,
		UpdatedAt:      flow.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// Cleanup xóa flow đã xong (Completed/Cancelled/Failed)
func (m *AttachFlowManager) Cleanup(flowID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if f, ok := m.flows[flowID]; ok {
		if f.State == FlowCompleted || f.State == FlowCancelled || f.State == FlowFailed {
			delete(m.flows, flowID)
		}
	}
}
