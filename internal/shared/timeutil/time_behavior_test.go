package timeutil

import (
	"testing"
)

func TestTimeValueBehavior(t *testing.T) {
	// 测试值类型的零值行为
	var zeroTime Time
	t.Logf("Zero Time: %+v", zeroTime)
	t.Logf("IsZero: %v", zeroTime.IsZero())
	t.Logf("Underlying time.Time IsZero: %v", zeroTime.Time.IsZero())

	// 创建一个具体的值
	now := Now()
	t.Logf("Now Time: %+v", now)
	t.Logf("IsZero: %v", now.IsZero())

	// 验证零值的时间戳
	zeroTimestamp := zeroTime.Unix()
	t.Logf("Zero Time Unix timestamp: %d", zeroTimestamp)

	// 验证具体值的时间戳
	nowTimestamp := now.Unix()
	t.Logf("Now Time Unix timestamp: %d", nowTimestamp)
}
