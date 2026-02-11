package gmongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Version 版本号
const Version = "1.0.0"

// 全局默认客户端
var defaultClient *Client

// Connect 连接 MongoDB 并设置为默认客户端
func Connect(config *Config) error {
	client, err := NewClient(config)
	if err != nil {
		return err
	}
	defaultClient = client
	return nil
}

// MustConnect 连接 MongoDB，失败则 panic
func MustConnect(config *Config) {
	if err := Connect(config); err != nil {
		panic(err)
	}
}

// SetDefaultClient 设置默认客户端
func SetDefaultClient(client *Client) {
	defaultClient = client
}

// GetDefaultClient 获取默认客户端
func GetDefaultClient() *Client {
	if defaultClient == nil {
		panic("default client not initialized, call Connect() first")
	}
	return defaultClient
}

// Close 关闭默认客户端
func Close() error {
	if defaultClient != nil {
		return defaultClient.Close()
	}
	return nil
}

// MongoClient 获取默认客户端的原生 mongo.Client
func MongoClient() *mongo.Client {
	return GetDefaultClient().Client()
}

// DB 使用默认客户端获取数据库
func DB(name ...string) *Database {
	return GetDefaultClient().Database(name...)
}

// Coll 使用默认客户端获取集合
func Coll(name string) *Collection {
	return GetDefaultClient().Collection(name)
}

// Ping 测试连接（使用默认客户端）
func Ping(ctx context.Context) error {
	return GetDefaultClient().Ping(ctx)
}

// HealthCheck 健康检查（使用默认客户端）
func HealthCheck() error {
	return GetDefaultClient().HealthCheck()
}

// Stats 获取统计信息（使用默认客户端）
func Stats() string {
	return GetDefaultClient().Stats()
}

// StartSession 开始会话（使用默认客户端）
func StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return GetDefaultClient().StartSession(opts...)
}

// UseSession 使用会话执行操作（使用默认客户端）
func UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return GetDefaultClient().UseSession(ctx, fn)
}

// UseSessionWithOptions 使用会话和选项执行操作（使用默认客户端）
func UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	return GetDefaultClient().UseSessionWithOptions(ctx, opts, fn)
}

// ListDatabases 列出所有数据库（使用默认客户端）
func ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error) {
	return GetDefaultClient().ListDatabases(ctx, filter, opts...)
}

// ListDatabaseNames 列出所有数据库名称（使用默认客户端）
func ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	return GetDefaultClient().ListDatabaseNames(ctx, filter, opts...)
}

// Watch 监听变更流（使用默认客户端）
func Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return GetDefaultClient().Watch(ctx, pipeline, opts...)
}

// Transaction 执行事务（使用默认客户端）
func Transaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) error, opts ...*options.TransactionOptions) error {
	return GetDefaultClient().UseSessionWithOptions(ctx, options.Session(), func(sessCtx mongo.SessionContext) error {
		if err := sessCtx.StartTransaction(opts...); err != nil {
			return err
		}
		// 执行事务函数
		if err := fn(sessCtx); err != nil {
			_ = sessCtx.AbortTransaction(ctx)
			return err
		}
		return sessCtx.CommitTransaction(ctx)
	})
}

// WithTransaction 使用事务执行操作（使用默认客户端）
func WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error) {
	session, err := GetDefaultClient().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, fn, opts...)
}

// Tx 简化的事务执行（使用默认客户端）
func Tx(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	session, err := GetDefaultClient().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})
	return err
}

// TxWithOptions 带选项的事务执行（使用默认客户端）
func TxWithOptions(ctx context.Context, opts *options.TransactionOptions, fn func(sessCtx mongo.SessionContext) error) error {
	session, err := GetDefaultClient().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	}, opts)
	return err
}
