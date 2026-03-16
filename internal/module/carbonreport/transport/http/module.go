package http

import (
	"app/internal/module/ipfs/rpc"
)

// Module carbonreport 模块
type Module struct {
	DDD      *FileDDD
	Handlers *Handlers
}

// NewModule 创建 carbonreport 模块
func NewModule(client *rpc.LApiStub, session string) (*Module, error) {
	// 初始化 DDD 组件
	ddd := InitFileDDD(client, session)

	// 初始化 HTTP handlers
	handlers := &Handlers{
		FileHandler: NewFileHandler(ddd.AppService),
	}

	return &Module{
		DDD:      ddd,
		Handlers: handlers,
	}, nil
}

// Close 关闭模块
func (m *Module) Close() error {
	if m.DDD.SessionRepo != nil {
		// 可以在此处关闭会话
	}
	return nil
}
