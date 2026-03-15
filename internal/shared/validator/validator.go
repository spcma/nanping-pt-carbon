package validator

import (
	"errors"
	"regexp"
)

// Validator 验证器接口
type Validator interface {
	// Validate 验证对象
	Validate(obj interface{}) error

	// ValidateField 验证字段
	ValidateField(fieldName string, value interface{}) error
}

// DefaultValidator 默认验证器
type DefaultValidator struct{}

// NewDefaultValidator 创建默认验证器
func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{}
}

// Validate 验证对象
func (v *DefaultValidator) Validate(obj interface{}) error {
	// TODO: 使用反射验证所有字段
	return nil
}

// ValidateField 验证字段
func (v *DefaultValidator) ValidateField(fieldName string, value interface{}) error {
	switch fieldName {
	case "license":
		return v.validateLicense(value)
	case "phone":
		return v.validatePhone(value)
	case "email":
		return v.validateEmail(value)
	default:
		return nil
	}
}

// validateLicense 验证车牌号
func (v *DefaultValidator) validateLicense(value interface{}) error {
	license, ok := value.(string)
	if !ok {
		return errors.New("车牌号必须是字符串")
	}

	if len(license) < 7 {
		return errors.New("车牌号长度不足")
	}

	return nil
}

// validatePhone 验证手机号
func (v *DefaultValidator) validatePhone(value interface{}) error {
	phone, ok := value.(string)
	if !ok {
		return errors.New("手机号必须是字符串")
	}

	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	if !matched {
		return errors.New("手机号格式不正确")
	}

	return nil
}

// validateEmail 验证邮箱
func (v *DefaultValidator) validateEmail(value interface{}) error {
	email, ok := value.(string)
	if !ok {
		return errors.New("邮箱必须是字符串")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return errors.New("邮箱格式不正确")
	}

	return nil
}
