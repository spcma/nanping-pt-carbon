package infrastructure

import (
	"app/internal/module/carbonreport/domain"
	"app/internal/module/ipfs/rpc"
	"context"
	"time"
)

type npfsSessionRepository struct {
	client  *rpc.LApiStub
	session string
}

// NewSessionRepository 创建会话仓储
func NewSessionRepository(client *rpc.LApiStub, session string) domain.SessionRepository {
	return &npfsSessionRepository{
		client:  client,
		session: session,
	}
}

func (r *npfsSessionRepository) GetCurrentSession(ctx context.Context) (*domain.NpfsSession, error) {
	return &domain.NpfsSession{
		SessionID: r.session,
		ClientURL: "127.0.0.1:4080",
		LoginTime: time.Now(),
	}, nil
}

func (r *npfsSessionRepository) RefreshSession(ctx context.Context) error {
	// TODO: 实现会话刷新逻辑
	return nil
}

func (r *npfsSessionRepository) CloseSession(ctx context.Context) error {
	if r.client != nil {
		return r.client.Logout(r.session, "")
	}
	return nil
}
