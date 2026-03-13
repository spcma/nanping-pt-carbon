package valueobject

import (
	"errors"
	"fmt"
)

// License license plate value object
type License struct {
	Prefix    string `json:"prefix"`
	Number    string `json:"number"`
	FullPlate string `json:"fullPlate"`
}

// NewLicense creates a license plate
func NewLicense(prefix, number string) (*License, error) {
	if prefix == "" {
		return nil, errors.New("license prefix cannot be empty")
	}
	if number == "" {
		return nil, errors.New("license number cannot be empty")
	}

	// 简单的车牌号验证（可根据需要扩展）
	fullPlate := fmt.Sprintf("%s%s", prefix, number)
	if len(fullPlate) < 7 {
		return nil, errors.New("invalid license plate format")
	}

	return &License{
		Prefix:    prefix,
		Number:    number,
		FullPlate: fullPlate,
	}, nil
}

// String string representation
func (l *License) String() string {
	return l.FullPlate
}

// Equal checks if equal
func (l *License) Equal(other *License) bool {
	if other == nil {
		return false
	}
	return l.Prefix == other.Prefix && l.Number == other.Number
}
