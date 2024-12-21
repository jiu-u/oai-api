package log

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jiu-u/oai-api/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const ctxLoggerKey = "zapLogger"

type Logger struct {
	*zap.Logger
}

type fileOutputOption struct {
	Filename   string // 日志文件路径
	MaxSize    int    // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups int    // 日志文件最多保存多少个备份
	MaxAge     int    // 文件最多保存多少天
	Compress   bool   // 是否压缩
}

func NewLogger(cfg *config.Config) *Logger {
	option := fileOutputOption{
		Filename:   cfg.Log.LogPath + "/" + cfg.Log.FileName,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
	}
	var level zapcore.Level
	//debug<info<warn<error<fatal<panic
	switch cfg.Log.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	coreList := make([]zapcore.Core, 0)
	// 根据config配置编码器
	consoleEncoder := getConsoleEncoder(getConsoleEncoderConfig())
	jsonEncoder := getJSONEncoder(getJSONEncoderConfig())
	// 输出到终端
	coreList = append(coreList, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	// 输出到文件
	coreList = append(coreList, zapcore.NewCore(jsonEncoder, getLogWriter(&option), level))
	// 输出到错误文件
	option.Filename = cfg.Log.LogPath + "/" + cfg.Log.ErrorFileName
	coreList = append(coreList, zapcore.NewCore(jsonEncoder, getLogWriter(&option), zap.ErrorLevel))
	logger := zap.New(zapcore.NewTee(coreList...), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(0))
	return &Logger{logger}
}

func getConsoleEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg
}

func getJSONEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg
}

func getConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return zapcore.NewConsoleEncoder(cfg)
}

func getJSONEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return zapcore.NewJSONEncoder(cfg)
}

func getLogWriter(option *fileOutputOption) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   option.Filename,   // 日志文件路径
		MaxSize:    option.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: option.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     option.MaxAge,     // 文件最多保存多少天
		Compress:   option.Compress,   // 是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// WithValue Adds a field to the specified context
// WithValue withValue会创建一个新的zapLogger上下文，并添加一些字段,然后存储在context.Context中。
// 将 zap 日志记录器及其附加的字段存储在 context.Context 中，以便在后续的处理中使用。
func (l *Logger) WithValue(ctx context.Context, fields ...zapcore.Field) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...)))
		return c
	}
	return context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
// WithContext 从 context.Context 中获取存储的 zap 日志记录器。
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	zl := ctx.Value(ctxLoggerKey)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &Logger{ctxLogger}
	}
	return l
}
