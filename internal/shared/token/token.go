package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明
type Claims struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// Config Token 配置
type Config struct {
	SecretKey     string        `mapstructure:"secret_key"`
	ExpireTime    time.Duration `mapstructure:"expire_time"`
	RefreshTime   time.Duration `mapstructure:"refresh_time"`
	Issuer        string        `mapstructure:"issuer"`
	BlacklistTime time.Duration `mapstructure:"blacklist_time"`
}

// Manager Token 管理器接口
type Manager interface {
	GenerateToken(userID int64, username string, roles []string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshToken(refreshToken string) (string, error)
	AddToBlacklist(tokenString string, expireTime time.Time) error
	IsInBlacklist(tokenString string) (bool, error)
	GetUserInfo(tokenString string) (*UserInfo, error)
}

// UserInfo 用户信息（用于缓存模式）
type UserInfo struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Status   string   `json:"status"` // normal, frozen, canceled
}

// JWTManager JWT Token 管理器实现
type JWTManager struct {
	config Config
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(config Config) *JWTManager {
	if config.ExpireTime == 0 {
		config.ExpireTime = 24 * time.Hour
	}
	if config.RefreshTime == 0 {
		config.RefreshTime = 7 * 24 * time.Hour
	}
	if config.Issuer == "" {
		config.Issuer = "iam"
	}
	if config.BlacklistTime == 0 {
		config.BlacklistTime = 30 * 24 * time.Hour
	}

	return &JWTManager{
		config: config,
	}
}

func (jm *JWTManager) GenerateToken(userID int64, username string, roles []string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.config.ExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jm.config.Issuer,
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.SecretKey))
}

// ValidateToken 验证 Token
func (jm *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jm.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		inBlacklist, err := jm.IsInBlacklist(tokenString)
		if err != nil {
			return nil, err
		}
		if inBlacklist {
			return nil, errors.New("token is in blacklist")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Token
func (jm *JWTManager) RefreshToken(refreshToken string) (string, error) {
	claims, err := jm.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	newToken, err := jm.GenerateToken(claims.UserID, claims.Username, claims.Roles)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// AddToBlacklist 添加 Token 到黑名单
func (jm *JWTManager) AddToBlacklist(tokenString string, expireTime time.Time) error {
	// TODO: 实现黑名单逻辑
	return nil
}

// IsInBlacklist 检查 Token 是否在黑名单中
func (jm *JWTManager) IsInBlacklist(tokenString string) (bool, error) {
	// TODO: 实现黑名单检查逻辑
	return false, nil
}

// ParseToken 解析 Token（不验证签名）
func (jm *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("failed to parse token claims")
}

// GetTokenTTL 获取 Token 剩余有效期
func (jm *JWTManager) GetTokenTTL(tokenString string) (time.Duration, error) {
	claims, err := jm.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	if claims.ExpiresAt == nil {
		return 0, errors.New("token has no expiration time")
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl < 0 {
		return 0, errors.New("token has expired")
	}

	return ttl, nil
}

// HasRole 检查用户是否有指定角色
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否有任一指定角色
func (c *Claims) HasAnyRole(roles []string) bool {
	for _, requiredRole := range roles {
		if c.HasRole(requiredRole) {
			return true
		}
	}
	return false
}

// GetRoles 获取用户的所有角色
func (c *Claims) GetRoles() []string {
	if len(c.Roles) > 0 {
		return c.Roles
	}

	return []string{}
}

// GetUserInfo 获取用户信息（JWT 模式下从 Claims 转换）
func (jm *JWTManager) GetUserInfo(tokenString string) (*UserInfo, error) {
	claims, err := jm.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		UserID:   claims.UserID,
		Username: claims.Username,
		Roles:    claims.GetRoles(),
		Status:   "normal",
	}, nil
}
