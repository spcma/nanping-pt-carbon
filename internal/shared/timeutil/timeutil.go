package timeutil

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// 自定义时间类型
type Time struct {
	time.Time
}

// Now 创建时间, 默认为当前参数, 可接收一个参数[time.Time]创建指定时间
func Now(newTime ...time.Time) Time {
	if len(newTime) == 0 {
		return Time{time.Now()}
	}

	return Time{newTime[0]}
}

// FromString 从字符串创建时间
func FromString(s string) (Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return Time{}, err
	}
	return Time{t}, nil
}

// FromUnix 从Unix时间戳创建时间
func FromUnix(sec int64) Time {
	return Time{time.Unix(sec, 0)}
}

// IsZero 判断是否为零值
func (t Time) IsZero() bool {
	return t.Time.IsZero()
}

// 实现JSON序列化接口
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format("2006-01-02 15:04:05"))
}

// 实现JSON反序列化接口
func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}

	t.Time = parsedTime
	return nil
}

// 实现数据库驱动接口
func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t *Time) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		t.Time = v
	case string:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		t.Time = parsedTime
	case nil:
		t.Time = time.Time{}
	default:
		return fmt.Errorf("cannot scan %T into custom Time", value)
	}
	return nil
}

func (t Time) ToTime() time.Time {
	return t.Time
}

// FromGoTime 从 time.Time 创建 Time
func FromGoTime(t time.Time) Time {
	return Time{t}
}

// FromTime 是 FromGoTime 的别名
func FromTime(t time.Time) Time {
	return Time{t}
}
