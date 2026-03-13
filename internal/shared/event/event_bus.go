package event

import (
	"context"
	"reflect"
	"strings"
	"sync"
)

// EventHandler 事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event DomainEvent) error
	// EventType 返回能处理的事件类型
	EventType() string
}

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event DomainEvent) error
	// Subscribe 订阅事件
	Subscribe(eventType string, handler EventHandler)
	// Unsubscribe 取消订阅
	Unsubscribe(eventType string, handler EventHandler)
}

// EventType 事件类型
type EventType string

const (
	// BusCreated 车辆创建事件
	BusCreated EventType = "bus.created"
	// BusAssigned 车辆分配线路事件
	BusAssigned EventType = "bus.assigned"
	// BusStatusChanged 车辆状态变更事件
	BusStatusChanged EventType = "bus.status_changed"

	// RouteCreated 线路创建事件
	RouteCreated EventType = "route.created"
	// RouteChanged 线路变更事件
	RouteChanged EventType = "route.changed"

	// StationCreated 站点创建事件
	StationCreated EventType = "station.created"
	// StationChanged 站点变更事件
	StationChanged EventType = "station.changed"
)

// InMemoryEventBus 内存事件总线实现
type InMemoryEventBus struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

// NewInMemoryEventBus 创建内存事件总线
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish 发布事件
func (b *InMemoryEventBus) Publish(ctx context.Context, event DomainEvent) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	eventType := string(GetEventType(event))
	if handlers, ok := b.handlers[eventType]; ok {
		for _, handler := range handlers {
			if err := handler.Handle(ctx, event); err != nil {
				return err
			}
		}
	}
	return nil
}

// Subscribe 订阅事件
func (b *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Unsubscribe 取消订阅
func (b *InMemoryEventBus) Unsubscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventType]
	for i, h := range handlers {
		if h == handler {
			b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// GetEventType 获取事件类型 (使用反射实现)
func GetEventType(event DomainEvent) EventType {
	if event == nil {
		return "unknown"
	}

	// 使用反射获取类型名称
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 将类型名称转换为 kebab-case 格式
	typeName := t.Name()
	if typeName == "" {
		return "unknown"
	}

	// 简单的转换逻辑：大写字母前加连字符，转小写
	var result strings.Builder
	for i, r := range typeName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('.')
			result.WriteRune(r + 32) // 转为小写
		} else {
			result.WriteRune(r)
		}
	}

	return EventType(strings.ToLower(result.String()))
}
