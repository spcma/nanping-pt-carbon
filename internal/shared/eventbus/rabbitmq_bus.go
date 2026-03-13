package eventbus

import (
	"app/internal/shared/event"
	"context"
)

// RabbitMQEventBus RabbitMQ 事件总线实现
type RabbitMQEventBus struct {
	// TODO: 添加 RabbitMQ 连接字段
}

// NewRabbitMQEventBus 创建 RabbitMQ 事件总线
func NewRabbitMQEventBus() *RabbitMQEventBus {
	return &RabbitMQEventBus{}
}

// Publish 发布事件到 RabbitMQ
func (b *RabbitMQEventBus) Publish(ctx context.Context, evt event.DomainEvent) error {
	// TODO: 实现 RabbitMQ 发布逻辑
	return nil
}

// Subscribe 订阅事件
func (b *RabbitMQEventBus) Subscribe(eventType string, handler event.EventHandler) {
	// TODO: 实现 RabbitMQ 订阅逻辑
}

// Unsubscribe 取消订阅
func (b *RabbitMQEventBus) Unsubscribe(eventType string, handler event.EventHandler) {
	// TODO: 实现取消订阅逻辑
}

// 确保实现了 event.EventBus 接口
var _ event.EventBus = (*RabbitMQEventBus)(nil)
