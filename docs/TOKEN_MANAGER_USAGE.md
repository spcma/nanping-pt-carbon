# Token 管理器使用指南

## 概述

当前项目支持两种 Token 管理模式：

1. **JWT 模式（TokenTypeJWT）**：传统 JWT，所有用户信息都编码在 token 中
2. **雪花 ID + 缓存模式（TokenTypeSnowflake）**：token 只是雪花 ID，用户信息存储在 Redis 缓存中

## 两种模式的区别

### JWT 模式
- **优点**：
  - 无状态，不依赖 Redis
  - 性能好，无需查询数据库/缓存
  - 部署简单
  
- **缺点**：
  - Token 较大（包含用户信息）
  - 无法主动撤销（只能等过期或维护黑名单）
  - 用户信息变更不会反映到已发放的 token

- **适用场景**：
  - 小型应用
  - 对性能要求高
  - 不需要频繁撤销 token

### 雪花 ID + 缓存模式
- **优点**：
  - Token 短小精悍（只是一个 ID）
  - 可以主动撤销 token
  - 用户信息变更实时生效
  - 支持踢用户下线
  
- **缺点**：
  - 依赖 Redis
  - 每次请求需要查询缓存
  - 部署复杂度高

- **适用场景**：
  - 中大型应用
  - 需要灵活的 token 管理
  - 安全性要求高

## 使用方式

### 1. 自动模式（推荐）

系统会根据 Redis 连接情况自动选择模式：

```go
// internal/transport/http/router.go
var tokenManager token.Manager

if redisClient != nil {
    // 优先使用雪花 ID + 缓存模式
    tokenManager, err = token.NewManager(token.ConfigEx{
        Type:        token.TokenTypeSnowflake,
        RedisClient: redisClient.GetClient(),
        ExpireTime:  24 * time.Hour,
    })
} else {
    // 降级到 JWT 模式
    tokenManager = initJWTManager()
}
```

### 2. 手动指定 JWT 模式

```go
tokenManager, err := token.NewManager(token.ConfigEx{
    Type: token.TokenTypeJWT,
    JWTConfig: token.Config{
        SecretKey:     "your-secret-key",
        ExpireTime:    24 * time.Hour,
        RefreshTime:   7 * 24 * time.Hour,
        Issuer:        "your-app",
        BlacklistTime: 30 * 24 * time.Hour,
    },
})
```

### 3. 手动指定雪花 ID 模式

```go
tokenManager, err := token.NewManager(token.ConfigEx{
    Type:        token.TokenTypeSnowflake,
    RedisClient: redisClient,
    UserCache:   cache.NewUserCache(redisClient, 1*time.Hour),
    ExpireTime:  24 * time.Hour,
})
```

## 接口说明

### Manager 接口

```go
type Manager interface {
    // 生成 Token（单角色）
    GenerateToken(userID int64, username, role string) (string, error)
    
    // 生成 Token（多角色）
    GenerateTokenWithRoles(userID int64, username string, roles []string) (string, error)
    
    // 验证 Token，返回 Claims
    ValidateToken(tokenString string) (*Claims, error)
    
    // 刷新 Token
    RefreshToken(refreshToken string) (string, error)
    
    // 添加到黑名单
    AddToBlacklist(tokenString string, expireTime time.Time) error
    
    // 检查是否在黑名单
    IsInBlacklist(tokenString string) (bool, error)
    
    // 获取用户信息（雪花 ID 模式专用）
    GetUserInfo(tokenString string) (*UserInfo, error)
}
```

### 雪花 ID 模式特有方法

```go
// SnowflakeTokenManager 特有的管理方法
type SnowflakeTokenManager struct {
    // ...
}

// 撤销指定 Token
func (m *SnowflakeTokenManager) RevokeToken(token string) error

// 撤销用户的所有 Token
func (m *SnowflakeTokenManager) RevokeAllUserTokens(userID int64) error

// 踢用户下线
func (m *SnowflakeTokenManager) KickoutUser(userID int64) error

// 获取用户的活跃 Token 数量
func (m *SnowflakeTokenManager) GetTokenCount(userID int64) (int64, error)
```

## 使用示例

### 基础使用

```go
// 登录时生成 Token
tokenString, err := tokenManager.GenerateToken(user.ID, user.Username, "USER")

// 中间件自动验证 Token
claims, err := tokenManager.ValidateToken(tokenString)
if err != nil {
    // Token 无效
    return
}

// 获取用户信息
userInfo, err := tokenManager.GetUserInfo(tokenString)
```

### 踢用户下线（雪花 ID 模式）

```go
// 转换为 SnowflakeTokenManager
if sfm, ok := tokenManager.(*token.SnowflakeTokenManager); ok {
    // 踢用户下线
    err := sfm.KickoutUser(userID)
    
    // 或者撤销指定 Token
    err = sfm.RevokeToken(tokenString)
}
```

## 配置建议

### 开发环境
- 使用 JWT 模式（简化部署）
- Token 过期时间：24 小时

### 生产环境
- 使用雪花 ID + 缓存模式（更好的安全性和管理性）
- Token 过期时间：根据业务需求（建议 8-24 小时）
- Redis 缓存过期时间：1 小时（会自动续期）

## 注意事项

1. **接口统一性**：两种模式都实现 `token.Manager` 接口，中间件和业务代码无需关心具体实现
2. **降级策略**：如果 Redis 不可用，系统会自动降级到 JWT 模式
3. **兼容性**：现有代码无需修改，`ValidateToken` 统一返回 `*Claims`
4. **扩展性**：可以通过类型断言访问雪花 ID 模式的特有功能
