package trace

import (
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	defaultTracerName = "github.com/nicexiaonie/gmongo/tracing/trace"
	defaultDBSystem   = "mongodb"
)

// Config trace 配置
type Config struct {
	Enabled bool

	TracerProvider oteltrace.TracerProvider
	TracerName     string

	SpanNameFormatter func(CommandStartedInfo) string
	Attributes        func(CommandStartedInfo) []attribute.KeyValue
	RecordCommand     bool
	CommandSanitizer  func(CommandStartedInfo) []attribute.KeyValue

	PoolEventHandler   func(PoolEventInfo)
	ServerEventHandler func(ServerEventInfo)
}

// CommandStartedInfo 命令开始事件信息
type CommandStartedInfo struct {
	DatabaseName string
	CommandName  string
	Collection   string
	RequestID    int64
	ConnectionID string
}

// PoolEventInfo 连接池事件信息
type PoolEventInfo struct {
	Type         string
	Address      string
	ConnectionID uint64
	Reason       string
	Error        error
}

// ServerEventInfo 服务发现和心跳事件信息
type ServerEventInfo struct {
	Type         string
	ConnectionID string
	Failure      error
}

// DefaultConfig 返回默认 trace 配置
func DefaultConfig() Config {
	return Config{
		Enabled:    true,
		TracerName: defaultTracerName,
	}
}
