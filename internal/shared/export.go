package shared

import (
	"app/internal/shared/cache"
	"app/internal/shared/crypto"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	err "app/internal/shared/error"
	"app/internal/shared/event"
	"app/internal/shared/eventbus"
	"app/internal/shared/http"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/logger"
	"app/internal/shared/persistence"
	"app/internal/shared/timeutil"
	"app/internal/shared/token"
	"app/internal/shared/transaction"
	"app/internal/shared/valueobject"
)

// 类型导出，方便外部使用
type (
	// Cache
	Cache         = cache.Cache
	InMemoryCache = cache.InMemoryCache
	RedisClient   = cache.RedisClient

	// Crypto
	CryptoSalt = string
	CryptoHash = string

	// DB
	DBConfig = db.Config

	// Entity
	BaseEntity = entity.BaseEntity
	// PaginationResult 需要根据实际使用场景指定类型参数

	// Error
	Error = err.Error

	// Event
	EventBus         = event.EventBus
	EventHandler     = event.EventHandler
	DomainEvent      = event.DomainEvent
	BaseEvent        = event.BaseEvent
	EventType        = event.EventType
	InMemoryEventBus = event.InMemoryEventBus

	// EventBus (RabbitMQ)
	RabbitMQEventBus = eventbus.RabbitMQEventBus

	// HTTP
	RouteHandler     = http.RouteHandler
	RouteGroupConfig = http.RouteGroupConfig
	AuthType         = http.AuthType

	// ID Generator
	IdgenGenerateFactory = idgen.IdgenGenerateFactory

	// Logger
	LogConfig      = logger.LogConfig
	CategoryConfig = logger.CategoryConfig
	RotationConfig = logger.RotationConfig
	OtherConfig    = logger.OtherConfig

	// Persistence
	PersistenceConfig = persistence.Config
	BaseRepository    = persistence.BaseRepository

	// Time
	Time = timeutil.Time

	// Token
	TokenConfig = token.Config
	Claims      = token.Claims
	JWTManager  = token.JWTManager

	// Transaction
	TransactionManager = transaction.TransactionManager

	// Value Objects
	Coordinate = valueobject.Coordinate
	License    = valueobject.License
	TimeRange  = valueobject.TimeRange
)

// 函数导出，方便外部使用
var (
	// Cache
	NewInMemoryCache = cache.NewInMemoryCache
	NewRedisClient   = cache.NewRedisClient

	// Crypto
	GenerateSalt   = crypto.GenerateSalt
	HashPassword   = crypto.HashPassword
	VerifyPassword = crypto.VerifyPassword

	// DB
	NewGormDB   = db.NewGormDB
	WithContext = db.WithContext

	// Entity
	NewPaginationResult = entity.NewPaginationResult

	// Error
	NewError = err.NewError

	// Event
	NewInMemoryEventBus = event.NewInMemoryEventBus
	GetEventType        = event.GetEventType

	// EventBus
	NewRabbitMQEventBus = eventbus.NewRabbitMQEventBus

	// HTTP
	ValidateRouteConfigs = http.ValidateRouteConfigs

	// Logger
	Initialize       = logger.Initialize
	NewLoggerFactory = logger.NewLoggerFactory
	Close            = logger.Close
	Get              = logger.Get
	Default          = logger.Default
	DebugFunc        = logger.Debug
	InfoFunc         = logger.Info
	WarnFunc         = logger.Warn
	ErrorFunc        = logger.Error

	// Persistence
	NewGORM           = persistence.NewGORM
	AutoMigrate       = persistence.AutoMigrate
	NewBaseRepository = persistence.NewBaseRepository

	// Time
	NewTime    = timeutil.New
	FromTime   = timeutil.FromTime
	FromGoTime = timeutil.FromGoTime

	// Token
	NewJWTManager = token.NewJWTManager

	// Transaction
	NewTransactionManager = transaction.NewTransactionManager
	RunInTransaction      = transaction.RunInTransaction

	// Value Objects
	NewCoordinate = valueobject.NewCoordinate
	NewLicense    = valueobject.NewLicense
	NewTimeRange  = valueobject.NewTimeRange
)
