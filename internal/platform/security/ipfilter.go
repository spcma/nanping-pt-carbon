package security

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

// IPFilter IP过滤器接口
type IPFilter interface {
	// IsBlocked 检查IP是否被阻止
	IsBlocked(ctx context.Context, ip string) (bool, error)
	// IsAllowed 检查IP是否被允许
	IsAllowed(ctx context.Context, ip string) (bool, error)
	// BlockIP 阻止IP
	BlockIP(ctx context.Context, ip string, duration time.Duration, reason string) error
	// AllowIP 允许IP
	AllowIP(ctx context.Context, ip string) error
	// UnblockIP 解除阻止IP
	UnblockIP(ctx context.Context, ip string) error
	// GetBlockedIPs 获取被阻止的IP列表
	GetBlockedIPs(ctx context.Context) ([]*BlockedIP, error)
}

// BlockedIP 被阻止的IP信息
type BlockedIP struct {
	IP        string    `json:"ip"`
	Reason    string    `json:"reason"`
	BlockedAt time.Time `json:"blocked_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RedisIPFilter Redis IP过滤器实现
type RedisIPFilter struct {
	client *redis.Client
	prefix string
}

// NewRedisIPFilter 创建Redis IP过滤器
func NewRedisIPFilter(client *redis.Client, prefix string) *RedisIPFilter {
	return &RedisIPFilter{
		client: client,
		prefix: prefix,
	}
}

// IsBlocked 检查IP是否被阻止
func (r *RedisIPFilter) IsBlocked(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("%s:blocked:%s", r.prefix, ip)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// IsAllowed 检查IP是否被允许
func (r *RedisIPFilter) IsAllowed(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("%s:allowed:%s", r.prefix, ip)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// BlockIP 阻止IP
func (r *RedisIPFilter) BlockIP(ctx context.Context, ip string, duration time.Duration, reason string) error {
	key := fmt.Sprintf("%s:blocked:%s", r.prefix, ip)

	// 存储阻止信息
	blockInfo := map[string]interface{}{
		"ip":         ip,
		"reason":     reason,
		"blocked_at": time.Now().Unix(),
		"expires_at": time.Now().Add(duration).Unix(),
	}

	err := r.client.HMSet(ctx, key, blockInfo).Err()
	if err != nil {
		return err
	}

	// 设置过期时间
	return r.client.Expire(ctx, key, duration).Err()
}

// AllowIP 允许IP
func (r *RedisIPFilter) AllowIP(ctx context.Context, ip string) error {
	key := fmt.Sprintf("%s:allowed:%s", r.prefix, ip)
	return r.client.Set(ctx, key, "allowed", 0).Err()
}

// UnblockIP 解除阻止IP
func (r *RedisIPFilter) UnblockIP(ctx context.Context, ip string) error {
	key := fmt.Sprintf("%s:blocked:%s", r.prefix, ip)
	return r.client.Del(ctx, key).Err()
}

// GetBlockedIPs 获取被阻止的IP列表
func (r *RedisIPFilter) GetBlockedIPs(ctx context.Context) ([]*BlockedIP, error) {
	pattern := fmt.Sprintf("%s:blocked:*", r.prefix)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var blockedIPs []*BlockedIP
	for _, key := range keys {
		ipInfo, err := r.client.HGetAll(ctx, key).Result()
		if err != nil {
			continue
		}

		if len(ipInfo) > 0 {
			blockedIP := &BlockedIP{
				IP:        ipInfo["ip"],
				Reason:    ipInfo["reason"],
				BlockedAt: time.Unix(parseInt64(ipInfo["blocked_at"]), 0),
				ExpiresAt: time.Unix(parseInt64(ipInfo["expires_at"]), 0),
			}
			blockedIPs = append(blockedIPs, blockedIP)
		}
	}

	return blockedIPs, nil
}

// MemoryIPFilter 内存IP过滤器实现（备用方案）
type MemoryIPFilter struct {
	blocked map[string]*BlockedIP
	allowed map[string]bool
	prefix  string
}

// NewMemoryIPFilter 创建内存IP过滤器
func NewMemoryIPFilter(prefix string) *MemoryIPFilter {
	return &MemoryIPFilter{
		blocked: make(map[string]*BlockedIP),
		allowed: make(map[string]bool),
		prefix:  prefix,
	}
}

// IsBlocked 检查IP是否被阻止
func (m *MemoryIPFilter) IsBlocked(ctx context.Context, ip string) (bool, error) {
	if blockedIP, exists := m.blocked[ip]; exists {
		// 检查是否过期
		if time.Now().After(blockedIP.ExpiresAt) {
			delete(m.blocked, ip)
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

// IsAllowed 检查IP是否被允许
func (m *MemoryIPFilter) IsAllowed(ctx context.Context, ip string) (bool, error) {
	return m.allowed[ip], nil
}

// BlockIP 阻止IP
func (m *MemoryIPFilter) BlockIP(ctx context.Context, ip string, duration time.Duration, reason string) error {
	m.blocked[ip] = &BlockedIP{
		IP:        ip,
		Reason:    reason,
		BlockedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return nil
}

// AllowIP 允许IP
func (m *MemoryIPFilter) AllowIP(ctx context.Context, ip string) error {
	m.allowed[ip] = true
	return nil
}

// UnblockIP 解除阻止IP
func (m *MemoryIPFilter) UnblockIP(ctx context.Context, ip string) error {
	delete(m.blocked, ip)
	return nil
}

// GetBlockedIPs 获取被阻止的IP列表
func (m *MemoryIPFilter) GetBlockedIPs(ctx context.Context) ([]*BlockedIP, error) {
	var blockedIPs []*BlockedIP
	now := time.Now()

	for ip, blockedIP := range m.blocked {
		// 清理过期的阻止记录
		if now.After(blockedIP.ExpiresAt) {
			delete(m.blocked, ip)
			continue
		}
		blockedIPs = append(blockedIPs, blockedIP)
	}

	return blockedIPs, nil
}

// parseInt64 安全解析int64
func parseInt64(s string) int64 {
	var result int64
	fmt.Sscanf(s, "%d", &result)
	return result
}

// IsValidIP 验证IP地址格式
func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsPrivateIP 检查是否为私有IP地址
func IsPrivateIP(ip string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}

	// 私有IP地址范围
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
	}

	for _, cidr := range privateRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ipAddr) {
			return true
		}
	}

	return false
}

// NormalizeIP 标准化IP地址（处理IPv4映射的IPv6地址）
func NormalizeIP(ip string) string {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return ip
	}

	// 如果是IPv4映射的IPv6地址，转换为IPv4
	if ipAddr.To4() != nil {
		return ipAddr.To4().String()
	}

	return ipAddr.String()
}
