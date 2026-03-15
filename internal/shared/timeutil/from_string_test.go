package timeutil

import (
	"testing"
	"time"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantHour  int
		wantMin   int
		wantSec   int
	}{
		{
			name:      "正常日期时间格式",
			input:     "2026-03-01 23:15:30",
			wantErr:   false,
			wantYear:  2026,
			wantMonth: time.March,
			wantDay:   1,
			wantHour:  23,
			wantMin:   15,
			wantSec:   30,
		},
		{
			name:      "零点时间",
			input:     "2026-03-01 00:00:00",
			wantErr:   false,
			wantYear:  2026,
			wantMonth: time.March,
			wantDay:   1,
			wantHour:  0,
			wantMin:   0,
			wantSec:   0,
		},
		{
			name:    "错误格式-缺少秒",
			input:   "2026-03-01 23:15",
			wantErr: true,
		},
		{
			name:    "错误格式-月份超出范围",
			input:   "2026-13-01 23:15:30",
			wantErr: true,
		},
		{
			name:    "错误格式-日期超出范围",
			input:   "2026-02-30 23:15:30",
			wantErr: true,
		},
		{
			name:    "错误格式-小时超出范围",
			input:   "2026-03-01 25:15:30",
			wantErr: true,
		},
		{
			name:    "错误格式-分钟超出范围",
			input:   "2026-03-01 23:61:30",
			wantErr: true,
		},
		{
			name:    "错误格式-秒超出范围",
			input:   "2026-03-01 23:15:61",
			wantErr: true,
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
		{
			name:      "边界值-最小有效日期",
			input:     "0001-01-01 00:00:00",
			wantErr:   false,
			wantYear:  1,
			wantMonth: time.January,
			wantDay:   1,
			wantHour:  0,
			wantMin:   0,
			wantSec:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.input)

			// 检查错误预期
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			// 如果期望错误，直接返回
			if tt.wantErr {
				return
			}

			// 检查解析后的值是否正确
			if got.Year() != tt.wantYear {
				t.Errorf("FromString(%q) year = %v, want %v", tt.input, got.Year(), tt.wantYear)
			}

			if got.Month() != tt.wantMonth {
				t.Errorf("FromString(%q) month = %v, want %v", tt.input, got.Month(), tt.wantMonth)
			}

			if got.Day() != tt.wantDay {
				t.Errorf("FromString(%q) day = %v, want %v", tt.input, got.Day(), tt.wantDay)
			}

			if got.Hour() != tt.wantHour {
				t.Errorf("FromString(%q) hour = %v, want %v", tt.input, got.Hour(), tt.wantHour)
			}

			if got.Minute() != tt.wantMin {
				t.Errorf("FromString(%q) minute = %v, want %v", tt.input, got.Minute(), tt.wantMin)
			}

			if got.Second() != tt.wantSec {
				t.Errorf("FromString(%q) second = %v, want %v", tt.input, got.Second(), tt.wantSec)
			}

			// 验证转换回time.Time的方法
			convertedTime := got.ToTime()
			if convertedTime.Year() != tt.wantYear ||
				convertedTime.Month() != tt.wantMonth ||
				convertedTime.Day() != tt.wantDay {
				t.Errorf("ToTime() conversion failed for %q", tt.input)
			}
		})
	}
}

func TestFromStringIntegration(t *testing.T) {
	// 集成测试：验证FromString与其它方法的配合使用

	// 测试正常的解析流程
	dateStr := "2026-12-25 15:30:45"
	parsedTime, err := FromString(dateStr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 验证IsZero方法
	if parsedTime.IsZero() {
		t.Error("Parsed time should not be zero")
	}

	// 验证格式化输出
	expectedFormat := "2026-12-25 15:30:45"
	actualFormat := parsedTime.Format("2006-01-02 15:04:05")
	if actualFormat != expectedFormat {
		t.Errorf("Format mismatch: got %q, want %q", actualFormat, expectedFormat)
	}

	// 测试Unix时间戳转换
	unixTime := parsedTime.Unix()
	expectedUnix := time.Date(2026, 12, 25, 15, 30, 45, 0, time.UTC).Unix()
	if unixTime != expectedUnix {
		t.Errorf("Unix timestamp mismatch: got %d, want %d", unixTime, expectedUnix)
	}

	// 测试从Unix时间戳重建
	rebuiltTime := FromUnix(unixTime)
	if rebuiltTime.Unix() != parsedTime.Unix() {
		t.Error("Rebuilt time from Unix timestamp doesn't match original")
	}
}
