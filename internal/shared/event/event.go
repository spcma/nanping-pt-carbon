package event

import "time"

// BaseEvent 领域事件基类
type BaseEvent struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// DomainEvent 领域事件接口
type DomainEvent interface {
	// GetID 获取事件 ID
	GetID() int64
	// GetTimestamp 获取事件时间戳
	GetTimestamp() time.Time
	// SetID 设置事件 ID
	SetID(id int64)
	// SetTimestamp 设置事件时间戳
	SetTimestamp(time time.Time)
}

// SetID 设置事件 ID
func (e *BaseEvent) SetID(id int64) {
	e.ID = id
}

// GetID 获取事件 ID
func (e *BaseEvent) GetID() int64 {
	return e.ID
}

// SetTimestamp 设置事件时间戳
func (e *BaseEvent) SetTimestamp(t time.Time) {
	e.Timestamp = t
}

// GetTimestamp 获取事件时间戳
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}
