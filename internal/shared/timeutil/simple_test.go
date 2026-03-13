package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func ExampleFromString() {
	// 正常使用示例
	t, err := FromString("2026-03-01 23:15:30")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Parsed time: %s\n", t.Format("2006-01-02 15:04:05"))
	fmt.Printf("Year: %d, Month: %d, Day: %d\n", t.Year(), t.Month(), t.Day())

	// Output:
	// Parsed time: 2026-03-01 23:15:30
	// Year: 2026, Month: 3, Day: 1
}

func TestFromStringBasic(t *testing.T) {
	// 基本功能测试
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid datetime", "2026-03-01 23:15:30", false},
		{"invalid format", "2026-03-01", true},
		{"empty string", "", true},
		{"invalid month", "2026-13-01 23:15:30", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestFromStringValues(t *testing.T) {
	// 测试具体值的解析
	input := "2026-12-25 15:30:45"
	expected, _ := time.Parse("2006-01-02 15:04:05", input)

	result, err := FromString(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.Year() != expected.Year() {
		t.Errorf("Year mismatch: got %d, want %d", result.Year(), expected.Year())
	}

	if result.Month() != expected.Month() {
		t.Errorf("Month mismatch: got %d, want %d", result.Month(), expected.Month())
	}

	if result.Day() != expected.Day() {
		t.Errorf("Day mismatch: got %d, want %d", result.Day(), expected.Day())
	}

	// 测试零值判断
	var zeroTime Time
	if !zeroTime.IsZero() {
		t.Error("Zero time should be zero")
	}

	if result.IsZero() {
		t.Error("Valid time should not be zero")
	}
}
