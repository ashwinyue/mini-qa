package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger 结构化日志接口
type Logger interface {
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, msg string, fields map[string]interface{})
	WithFields(fields map[string]interface{}) Logger
}

// logrusLogger logrus 实现
type logrusLogger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

// Config 日志配置
type Config struct {
	Level    string // debug, info, warn, error
	Format   string // json, text
	Output   string // stdout, file
	FilePath string // 日志文件路径
}

// New 创建新的日志实例
func New(config Config) (Logger, error) {
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	logger.SetLevel(level)

	// 设置日志格式
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		})
	default:
		return nil, fmt.Errorf("invalid log format: %s", config.Format)
	}

	// 设置输出目标
	var output io.Writer
	switch config.Output {
	case "stdout":
		output = os.Stdout
	case "file":
		if config.FilePath == "" {
			return nil, fmt.Errorf("file_path is required when output is file")
		}
		// 创建日志目录
		if err := os.MkdirAll(config.FilePath[:len(config.FilePath)-len("/app.log")], 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		output = file
	default:
		return nil, fmt.Errorf("invalid log output: %s", config.Output)
	}
	logger.SetOutput(output)

	return &logrusLogger{
		logger: logger,
		fields: logrus.Fields{},
	}, nil
}

// extractContextFields 从 context 中提取通用字段
func extractContextFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}

	// 提取 trace_id
	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields["trace_id"] = traceID
	}

	// 提取 tenant_id
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		fields["tenant_id"] = tenantID
	}

	// 提取 session_id
	if sessionID := ctx.Value("session_id"); sessionID != nil {
		fields["session_id"] = sessionID
	}

	// 提取 request_id
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields["request_id"] = requestID
	}

	return fields
}

// mergeFields 合并字段
func mergeFields(base, additional logrus.Fields) logrus.Fields {
	result := make(logrus.Fields)
	for k, v := range base {
		result[k] = v
	}
	for k, v := range additional {
		result[k] = v
	}
	return result
}

// Debug 记录 debug 级别日志
func (l *logrusLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	contextFields := extractContextFields(ctx)
	allFields := mergeFields(l.fields, contextFields)
	allFields = mergeFields(allFields, fields)
	l.logger.WithFields(allFields).Debug(msg)
}

// Info 记录 info 级别日志
func (l *logrusLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	contextFields := extractContextFields(ctx)
	allFields := mergeFields(l.fields, contextFields)
	allFields = mergeFields(allFields, fields)
	l.logger.WithFields(allFields).Info(msg)
}

// Warn 记录 warn 级别日志
func (l *logrusLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	contextFields := extractContextFields(ctx)
	allFields := mergeFields(l.fields, contextFields)
	allFields = mergeFields(allFields, fields)
	l.logger.WithFields(allFields).Warn(msg)
}

// Error 记录 error 级别日志
func (l *logrusLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	contextFields := extractContextFields(ctx)
	allFields := mergeFields(l.fields, contextFields)
	allFields = mergeFields(allFields, fields)
	l.logger.WithFields(allFields).Error(msg)
}

// WithFields 创建带有预设字段的新 logger
func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	newFields := mergeFields(l.fields, fields)
	return &logrusLogger{
		logger: l.logger,
		fields: newFields,
	}
}
