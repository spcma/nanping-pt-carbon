package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端包装
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(addr, password string, db, poolSize, minIdleConns int, timeout, idleTimeout time.Duration) (*RedisClient, error) {

	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		PoolSize:        poolSize,
		MinIdleConns:    minIdleConns,
		DialTimeout:     timeout,
		ReadTimeout:     timeout,
		WriteTimeout:    timeout,
		ConnMaxIdleTime: idleTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

// Close 关闭Redis连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Get 获取缓存值
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set 设置缓存值
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// SetNX 设置缓存值（仅当key不存在时）
func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// Del 删除缓存键
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire 设置键的过期时间
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}

// TTL 获取键的剩余生存时间
func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// HSet 设置哈希字段
func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段值
func (r *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希所有字段
func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// SAdd 添加集合成员
func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// SRem 移除集合成员
func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

// SIsMember 检查成员是否在集合中
func (r *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

// Pipeline 管道操作
func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// TxPipeline 事务管道操作
func (r *RedisClient) TxPipeline() redis.Pipeliner {
	return r.client.TxPipeline()
}

// FlushDB 清空当前数据库
func (r *RedisClient) FlushDB(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

// Info 获取Redis服务器信息
func (r *RedisClient) Info(ctx context.Context, section ...string) *redis.StringCmd {
	return r.client.Info(ctx, section...)
}

// Incr 自增计数器
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Decr 自减计数器
func (r *RedisClient) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// IncrBy 按指定值自增
func (r *RedisClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// DecrBy 按指定值自减
func (r *RedisClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.DecrBy(ctx, key, value).Result()
}

// GetClient 获取底层的redis.Client
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}
