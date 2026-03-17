package logger

import (
	"context"
	"fmt"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"os"
	"strings"
	"sync"
	"time"
)

// 日志级别
const (
	DebugLevel string = "debug"
	InfoLevel  string = "info"
	WarnLevel  string = "warn"
	ErrorLevel string = "error"
)

type LogConfig struct {
	Level       string           `mapstructure:"level" label:"日志等级"`
	Path        string           `mapstructure:"path" label:"日志目录"`
	Development bool             `mapstructure:"development"`
	Rotation    RotationConfig   `mapstructure:"rotation"`
	Other       OtherConfig      `mapstructure:"other"`
	Categories  []CategoryConfig `mapstructure:"categories"`
}

type RotationConfig struct {
	Enabled    bool `mapstructure:"enabled"`
	MaxSize    int  `mapstructure:"max_size"`    // MB
	MaxBackups int  `mapstructure:"max_backups"` // 文件数量
	MaxAge     int  `mapstructure:"max_age"`     // 天数
	Compress   bool `mapstructure:"compress"`    // 轮转的日志是否压缩
}

type OtherConfig struct {
	Path  string `mapstructure:"path"`
	Level string `mapstructure:"level"`
}

type CategoryConfig struct {
	Path     string `mapstructure:"path"`
	Name     string `mapstructure:"name"`
	Filename string `mapstructure:"filename"`
	Level    string `mapstructure:"level"`
}

type Logger struct {
	mu       sync.RWMutex
	config   *LogConfig
	loggers  map[string]*zap.Logger
	writers  map[string]zapcore.WriteSyncer
	initLock map[string]*sync.Mutex
	initOnce map[string]*sync.Once
}

// 全局日志实例
var (
	globalLogger *Logger
	once         sync.Once
)

func Initialize(config *LogConfig) error {
	var initErr error
	once.Do(func() {
		globalLogger, initErr = NewLoggerFactory(config)
		if initErr == nil {
			InitGlobalLoggers()
		}
	})
	return initErr
}

func NewLoggerFactory(config *LogConfig) (*Logger, error) {
	pl := &Logger{
		config:   config,
		loggers:  make(map[string]*zap.Logger),
		writers:  make(map[string]zapcore.WriteSyncer),
		initLock: make(map[string]*sync.Mutex),
		initOnce: make(map[string]*sync.Once),
	}

	if config.Path == "" {
		config.Path = "./logs/"
	}

	if config.Other.Path == "" {
		config.Other.Path = config.Path + ".//other"
	}

	if config.Other.Level == "" {
		config.Other.Level = InfoLevel
	}

	// 创建日志目录
	if err := os.MkdirAll(config.Path, 0755); err != nil {
		return nil, fmt.Errorf("create log directory failed: %w", err)
	}

	return pl, nil
}

func (pl *Logger) getLoggerSimple(category string) *zap.Logger {
	// fast path：已存在
	pl.mu.RLock()
	logger, ok := pl.loggers[category]
	pl.mu.RUnlock()
	if ok {
		return logger
	}

	// 获取 or 创建 onceInstance
	pl.mu.Lock()
	onceInstance, ok := pl.initOnce[category]
	if !ok {
		onceInstance = &sync.Once{}
		pl.initOnce[category] = onceInstance
	}
	pl.mu.Unlock()

	// 确保只初始化一次
	onceInstance.Do(func() {
		if err := pl.initCategoryLogger(category); err != nil {
			fmt.Printf("init logger failed: %s %v\n", category, err)
		}
	})

	// 再次获取
	pl.mu.RLock()
	logger = pl.loggers[category]
	pl.mu.RUnlock()

	if logger == nil {
		return pl.createFallbackLogger()
	}

	return logger
}

// 获取指定分类的logger
func (pl *Logger) getLogger(category string) *zap.Logger {
	return pl.getLoggerSimple(category)
}

func (pl *Logger) getLoggerFull(category string) *zap.Logger {
	pl.mu.RLock()
	if logger, exists := pl.loggers[category]; exists {
		pl.mu.RUnlock()
		return logger
	}

	// 获取该 category 的初始化锁
	initMu, exists := pl.initLock[category]
	if !exists {
		pl.mu.RUnlock()
		pl.mu.Lock()
		if initMu, exists = pl.initLock[category]; !exists {
			initMu = &sync.Mutex{}
			pl.initLock[category] = initMu
		}
		pl.mu.Unlock()
	} else {
		pl.mu.RUnlock()
	}

	// 设置超时时间，避免长时间阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 使用带超时的锁获取
	done := make(chan struct{})

	go func() {
		initMu.Lock()
		defer initMu.Unlock()
		defer close(done)

		// 检查 context 是否已取消
		select {
		case <-ctx.Done():
			// 超时或取消，不执行初始化
			return
		default:
		}

		// 双重检查
		pl.mu.RLock()
		if _, exists = pl.loggers[category]; exists {
			pl.mu.RUnlock()
			return
		}
		pl.mu.RUnlock()

		// 初始化
		if err := pl.initCategoryLogger(category); err != nil {
			fmt.Println(fmt.Sprintf("Failed to initialize logger for category: %s %v", category, err))
			return
		}
	}()

	select {
	case <-done:
		pl.mu.RLock()
		defer pl.mu.RUnlock()
		if logger, exists := pl.loggers[category]; exists {
			return logger
		}
		return pl.fallbackLogger()
	case <-ctx.Done():
		// 超时，返回降级 logger
		return pl.fallbackLogger()
	}
}

var fallbackOnce sync.Once
var fallback *zap.Logger

func (pl *Logger) fallbackLogger() *zap.Logger {
	fallbackOnce.Do(func() {
		fallback = pl.createFallbackLogger()
	})
	return fallback
}

// 创建降级logger（当配置有问题时使用）
func (pl *Logger) createFallbackLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     pl.customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.Lock(os.Stderr),
		zapcore.InfoLevel,
	)

	return zap.New(core)
}

func (pl *Logger) initCategoryLogger(category string) error {
	// 查找配置
	var categoryConfig *CategoryConfig
	for _, cfg := range pl.config.Categories {
		if cfg.Name == category {
			categoryConfig = &cfg
			break
		}
	}

	if categoryConfig == nil {
		categoryConfig = &CategoryConfig{
			Name:     category,
			Filename: fmt.Sprintf("%s/%s/log.log", pl.config.Other.Path, category),
			Level:    pl.config.Other.Level,
		}
	} else {
		categoryConfig.Filename = fmt.Sprintf("%s/%s/log.log", pl.config.Path, category)
	}

	// 创建日志目录
	if err := os.MkdirAll(pl.config.Other.Path+"/"+category, 0755); err != nil {
		return err
	}

	// 创建写入器
	writer, err := pl.createWriter(categoryConfig)
	if err != nil {
		return err
	}

	// 创建编码器
	encoder := pl.createEncoder()

	// 设置日志级别
	level := pl.parseLevel(categoryConfig.Level)
	if level == zapcore.Level(-1) {
		level = pl.parseLevel(pl.config.Level)
	}

	core := zapcore.NewCore(encoder, writer, level)

	// 创建logger
	logger := zap.New(core, pl.buildOptions()...)

	pl.writers[category] = writer
	pl.loggers[category] = logger

	return nil
}

func (pl *Logger) createEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     pl.customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if pl.config.Development {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}

func (pl *Logger) customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func (pl *Logger) createWriter(config *CategoryConfig) (zapcore.WriteSyncer, error) {

	if pl.config.Rotation.Enabled {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    pl.config.Rotation.MaxSize,
			MaxBackups: pl.config.Rotation.MaxBackups,
			MaxAge:     pl.config.Rotation.MaxAge,
			Compress:   pl.config.Rotation.Compress,
			LocalTime:  true,
		}
		return zapcore.AddSync(lumberjackLogger), nil
	}

	file, err := os.OpenFile(config.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file failed: %w", err)
	}

	return zapcore.AddSync(file), nil
}

func (pl *Logger) buildOptions() []zap.Option {
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	}

	if pl.config.Development {
		options = append(options, zap.Development())
	}

	return options
}

func (pl *Logger) parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// 基础日志方法
func (pl *Logger) Debug(category string, msg string, fields ...zap.Field) {
	if logger := pl.getLogger(category); logger != nil {
		logger.Debug(msg, fields...)
	}
}

func (pl *Logger) Info(category string, msg string, fields ...zap.Field) {
	if logger := pl.getLogger(category); logger != nil {
		logger.Info(msg, fields...)
	}
}

func (pl *Logger) Warn(category string, msg string, fields ...zap.Field) {
	if logger := pl.getLogger(category); logger != nil {
		logger.Warn(msg, fields...)
	}
}

func (pl *Logger) Error(category string, msg string, fields ...zap.Field) {
	if logger := pl.getLogger(category); logger != nil {
		logger.Error(msg, fields...)
	}
}

// 关闭所有logger，确保日志写入完成
func (pl *Logger) Close() error {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	var errs []error

	for category, logger := range pl.loggers {
		if err := logger.Sync(); err != nil {
			// 忽略sync的标准错误（如同步到stdout/stderr时的错误）
			if !isSyncError(err) {
				errs = append(errs, fmt.Errorf("sync logger %s failed: %w", category, err))
			}
		}
	}

	for _, writer := range pl.writers {
		if err := writer.Sync(); err != nil {
			errs = append(errs, fmt.Errorf("sync writer failed: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close logger with errors: %v", errs)
	}

	return nil
}

// 检查是否是可忽略的sync错误
func isSyncError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := err.Error()
	// 这些是zap同步时的常见可忽略错误
	ignoreErrors := []string{
		"sync /dev/stdout: invalid argument",
		"sync /dev/stderr: invalid argument",
		"invalid argument",
	}

	for _, ignore := range ignoreErrors {
		if strings.Contains(errorStr, ignore) {
			return true
		}
	}

	return false
}

// 全局关闭方法
func Close() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}

func Get(category string) *zap.Logger {
	return globalLogger.getLogger(category)
}

func Default() *zap.Logger {
	return globalLogger.getLogger("runtime")
}

// 全局便捷方法 - 结构化日志
func Debug(category string, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Debug(category, msg, fields...)
	}
}

func Info(category string, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Info(category, msg, fields...)
	}
}

func Warn(category string, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Warn(category, msg, fields...)
	}
}

func Error(category string, msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Error(category, msg, fields...)
	}
}
