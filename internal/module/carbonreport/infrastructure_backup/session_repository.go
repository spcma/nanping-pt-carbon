package infrastructure

import (
	"app/internal/module/carbonreport/domain"
	"app/internal/rpc"
	"time"
)

// NpfsSessionRepository NPFS 会话仓储实现
type NpfsSessionRepository struct {
	client  *rpc.LApiStub
	session string
}

// NewNpfsSessionRepository 创建 NPFS 会话仓储
func NewNpfsSessionRepository(client *rpc.LApiStub, session string) *NpfsSessionRepository {
	return &NpfsSessionRepository{
		client:  client,
		session: session,
	}
}

// GetCurrentSession 获取当前会话
func (r *NpfsSessionRepository) GetCurrentSession() (*domain.NpfsSession, error) {
	return &domain.NpfsSession{
		SessionID: r.session,
		ClientURL: "127.0.0.1:4080",
		LoginTime: time.Now(),
	}, nil
}

// RefreshSession 刷新会话
func (r *NpfsSessionRepository) RefreshSession() error {
	// TODO: 实现会话刷新逻辑
	return nil
}

// CloseSession 关闭会话
func (r *NpfsSessionRepository) CloseSession() error {
	if r.client != nil {
		return r.client.Logout(r.session, "")
	}
	return nil
}
