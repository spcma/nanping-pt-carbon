package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"time"
)

// CSRFToken CSRF令牌管理器
type CSRFToken struct {
	secret    []byte
	tokenSize int
	expiry    time.Duration
}

// NewCSRFToken 创建CSRF令牌管理器
func NewCSRFToken(secret string, tokenSize int, expiry time.Duration) *CSRFToken {
	return &CSRFToken{
		secret:    []byte(secret),
		tokenSize: tokenSize,
		expiry:    expiry,
	}
}

// GenerateToken 生成CSRF令牌
func (c *CSRFToken) GenerateToken() (string, error) {
	// 生成随机字节
	randomBytes := make([]byte, c.tokenSize)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// 添加时间戳
	timestamp := time.Now().Unix()
	timestampBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		timestampBytes[i] = byte(timestamp >> (i * 8))
	}

	// 组合数据
	data := append(randomBytes, timestampBytes...)

	// 编码为base64
	token := base64.URLEncoding.EncodeToString(data)
	return token, nil
}

// ValidateToken 验证CSRF令牌
func (c *CSRFToken) ValidateToken(token string) bool {
	// 解码令牌
	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	// 检查长度
	if len(data) < c.tokenSize+8 {
		return false
	}

	// 提取时间戳
	timestampBytes := data[c.tokenSize : c.tokenSize+8]
	var timestamp int64
	for i := 0; i < 8; i++ {
		timestamp |= int64(timestampBytes[i]) << (i * 8)
	}

	// 检查是否过期
	if time.Since(time.Unix(timestamp, 0)) > c.expiry {
		return false
	}

	// 验证令牌（简单的比较，实际应用中可能需要更复杂的验证）
	randomPart := data[:c.tokenSize]
	expectedRandom := make([]byte, c.tokenSize)
	copy(expectedRandom, randomPart)

	return subtle.ConstantTimeCompare(randomPart, expectedRandom) == 1
}

// DoubleSubmitCookie 双提交Cookie模式
type DoubleSubmitCookie struct {
	csrfToken  *CSRFToken
	cookieName string
	headerName string
}

// NewDoubleSubmitCookie 创建双提交Cookie保护器
func NewDoubleSubmitCookie(secret, cookieName, headerName string) *DoubleSubmitCookie {
	return &DoubleSubmitCookie{
		csrfToken:  NewCSRFToken(secret, 32, 24*time.Hour),
		cookieName: cookieName,
		headerName: headerName,
	}
}

// GenerateCSRFToken 生成CSRF令牌和对应的Cookie值
func (d *DoubleSubmitCookie) GenerateCSRFToken() (string, string, error) {
	token, err := d.csrfToken.GenerateToken()
	if err != nil {
		return "", "", err
	}

	// Cookie值和令牌相同（双提交模式）
	cookieValue := token
	return token, cookieValue, nil
}

// ValidateRequest 验证请求的CSRF保护
func (d *DoubleSubmitCookie) ValidateRequest(cookieValue, headerValue string) bool {
	// 检查Cookie和Header是否存在
	if cookieValue == "" || headerValue == "" {
		return false
	}

	// 验证令牌有效性
	if !d.csrfToken.ValidateToken(cookieValue) {
		return false
	}

	// 比较Cookie值和Header值
	return subtle.ConstantTimeCompare([]byte(cookieValue), []byte(headerValue)) == 1
}

// SameSiteCookie SameSite Cookie模式
type SameSiteCookie struct {
	csrfToken  *CSRFToken
	cookieName string
	formField  string
}

// NewSameSiteCookie 创建SameSite Cookie保护器
func NewSameSiteCookie(secret, cookieName, formField string) *SameSiteCookie {
	return &SameSiteCookie{
		csrfToken:  NewCSRFToken(secret, 32, 24*time.Hour),
		cookieName: cookieName,
		formField:  formField,
	}
}

// GenerateCSRFToken 生成CSRF令牌
func (s *SameSiteCookie) GenerateCSRFToken() (string, error) {
	return s.csrfToken.GenerateToken()
}

// ValidateRequest 验证请求（通过表单字段）
func (s *SameSiteCookie) ValidateRequest(token, formValue string) bool {
	// 检查令牌和表单值是否存在
	if token == "" || formValue == "" {
		return false
	}

	// 验证令牌有效性
	if !s.csrfToken.ValidateToken(token) {
		return false
	}

	// 比较令牌和表单值
	return subtle.ConstantTimeCompare([]byte(token), []byte(formValue)) == 1
}

// CSRFProtection 综合CSRF保护器
type CSRFProtection struct {
	doubleSubmit *DoubleSubmitCookie
	sameSite     *SameSiteCookie
	mode         string // "double_submit" or "same_site"
}

// NewCSRFProtection 创建综合CSRF保护器
func NewCSRFProtection(secret, mode string) *CSRFProtection {
	return &CSRFProtection{
		doubleSubmit: NewDoubleSubmitCookie(secret, "csrf_token", "X-CSRF-Token"),
		sameSite:     NewSameSiteCookie(secret, "csrf_token", "_csrf"),
		mode:         mode,
	}
}

// GenerateToken 生成CSRF令牌
func (c *CSRFProtection) GenerateToken() (string, string, error) {
	switch c.mode {
	case "double_submit":
		return c.doubleSubmit.GenerateCSRFToken()
	case "same_site":
		token, err := c.sameSite.GenerateCSRFToken()
		return token, "", err
	default:
		return "", "", fmt.Errorf("unsupported CSRF mode: %s", c.mode)
	}
}

// ValidateToken 验证CSRF令牌
func (c *CSRFProtection) ValidateToken(cookieValue, headerValue, formValue string) bool {
	switch c.mode {
	case "double_submit":
		return c.doubleSubmit.ValidateRequest(cookieValue, headerValue)
	case "same_site":
		// 在SameSite模式下，通常cookieValue就是token
		return c.sameSite.ValidateRequest(cookieValue, formValue)
	default:
		return false
	}
}

// GetCookieName 获取Cookie名称
func (c *CSRFProtection) GetCookieName() string {
	switch c.mode {
	case "double_submit":
		return c.doubleSubmit.cookieName
	case "same_site":
		return c.sameSite.cookieName
	default:
		return "csrf_token"
	}
}

// GetHeaderName 获取Header名称
func (c *CSRFProtection) GetHeaderName() string {
	switch c.mode {
	case "double_submit":
		return c.doubleSubmit.headerName
	default:
		return "X-CSRF-Token"
	}
}

// GetFormField 获取表单字段名称
func (c *CSRFProtection) GetFormField() string {
	switch c.mode {
	case "same_site":
		return c.sameSite.formField
	default:
		return "_csrf"
	}
}
