package gmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// TxFunc 事务函数类型
type TxFunc func(sessCtx mongo.SessionContext) error

// TxFuncWithResult 带返回值的事务函数类型
type TxFuncWithResult func(sessCtx mongo.SessionContext) (interface{}, error)

// TransactionManager 事务管理器
type TransactionManager struct {
	client *Client
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(client *Client) *TransactionManager {
	return &TransactionManager{client: client}
}

// Execute 执行事务
func (tm *TransactionManager) Execute(ctx context.Context, fn TxFunc, opts ...*options.TransactionOptions) error {
	session, err := tm.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	}, opts...)
	return err
}

// ExecuteWithResult 执行事务并返回结果
func (tm *TransactionManager) ExecuteWithResult(ctx context.Context, fn TxFuncWithResult, opts ...*options.TransactionOptions) (interface{}, error) {
	session, err := tm.client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, fn, opts...)
}

// WithSession 使用会话执行操作
func (tm *TransactionManager) WithSession(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	return tm.client.UseSession(ctx, fn)
}

// WithSessionAndOptions 使用会话和选项执行操作
func (tm *TransactionManager) WithSessionAndOptions(ctx context.Context, opts *options.SessionOptions, fn func(sessCtx mongo.SessionContext) error) error {
	return tm.client.UseSessionWithOptions(ctx, opts, fn)
}

// 事务选项构建器

// TxOptionsBuilder 事务选项构建器
type TxOptionsBuilder struct {
	opts *options.TransactionOptions
}

// NewTxOptions 创建事务选项构建器
func NewTxOptions() *TxOptionsBuilder {
	return &TxOptionsBuilder{
		opts: options.Transaction(),
	}
}

// ReadConcern 设置读关注
func (tb *TxOptionsBuilder) ReadConcern(rc *readconcern.ReadConcern) *TxOptionsBuilder {
	tb.opts.SetReadConcern(rc)
	return tb
}

// ReadConcernLocal 设置读关注为 local
func (tb *TxOptionsBuilder) ReadConcernLocal() *TxOptionsBuilder {
	return tb.ReadConcern(readconcern.Local())
}

// ReadConcernMajority 设置读关注为 majority
func (tb *TxOptionsBuilder) ReadConcernMajority() *TxOptionsBuilder {
	return tb.ReadConcern(readconcern.Majority())
}

// ReadConcernSnapshot 设置读关注为 snapshot
func (tb *TxOptionsBuilder) ReadConcernSnapshot() *TxOptionsBuilder {
	return tb.ReadConcern(readconcern.Snapshot())
}

// WriteConcern 设置写关注
func (tb *TxOptionsBuilder) WriteConcern(wc *writeconcern.WriteConcern) *TxOptionsBuilder {
	tb.opts.SetWriteConcern(wc)
	return tb
}

// WriteConcernMajority 设置写关注为 majority
func (tb *TxOptionsBuilder) WriteConcernMajority() *TxOptionsBuilder {
	return tb.WriteConcern(writeconcern.Majority())
}

// WriteConcernW1 设置写关注为 w1
func (tb *TxOptionsBuilder) WriteConcernW1() *TxOptionsBuilder {
	return tb.WriteConcern(writeconcern.W1())
}

// MaxCommitTime 设置最大提交时间
func (tb *TxOptionsBuilder) MaxCommitTime(d interface{}) *TxOptionsBuilder {
	if duration, ok := d.(*time.Duration); ok {
		tb.opts.SetMaxCommitTime(duration)
	}
	return tb
}

// Build 构建选项
func (tb *TxOptionsBuilder) Build() *options.TransactionOptions {
	return tb.opts
}

// 便捷事务函数

// RunInTransaction 在事务中运行函数（使用默认客户端）
func RunInTransaction(ctx context.Context, fn TxFunc, opts ...*options.TransactionOptions) error {
	return GetDefaultClient().UseSessionWithOptions(ctx, options.Session(), func(sessCtx mongo.SessionContext) error {
		if err := sessCtx.StartTransaction(opts...); err != nil {
			return err
		}

		if err := fn(sessCtx); err != nil {
			_ = sessCtx.AbortTransaction(ctx)
			return err
		}

		return sessCtx.CommitTransaction(ctx)
	})
}

// RunInTransactionWithResult 在事务中运行函数并返回结果（使用默认客户端）
func RunInTransactionWithResult(ctx context.Context, fn TxFuncWithResult, opts ...*options.TransactionOptions) (interface{}, error) {
	session, err := GetDefaultClient().StartSession()
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, fn, opts...)
}

// WithReadConcernMajority 使用 majority 读关注执行事务
func WithReadConcernMajority(ctx context.Context, fn TxFunc) error {
	opts := NewTxOptions().ReadConcernMajority().Build()
	return RunInTransaction(ctx, fn, opts)
}

// WithWriteConcernMajority 使用 majority 写关注执行事务
func WithWriteConcernMajority(ctx context.Context, fn TxFunc) error {
	opts := NewTxOptions().WriteConcernMajority().Build()
	return RunInTransaction(ctx, fn, opts)
}

// WithSnapshot 使用 snapshot 读关注执行事务
func WithSnapshot(ctx context.Context, fn TxFunc) error {
	opts := NewTxOptions().ReadConcernSnapshot().Build()
	return RunInTransaction(ctx, fn, opts)
}

// 事务辅助函数

// InTransaction 检查是否在事务中
func InTransaction(ctx context.Context) bool {
	if sessCtx, ok := ctx.(mongo.SessionContext); ok {
		return sessCtx.Client() != nil
	}
	return false
}

// GetSessionFromContext 从上下文获取会话
func GetSessionFromContext(ctx context.Context) (mongo.SessionContext, bool) {
	sessCtx, ok := ctx.(mongo.SessionContext)
	return sessCtx, ok
}

// AbortTransaction 中止事务
func AbortTransaction(ctx context.Context) error {
	if sessCtx, ok := GetSessionFromContext(ctx); ok {
		return sessCtx.AbortTransaction(ctx)
	}
	return fmt.Errorf("not in transaction context")
}

// CommitTransaction 提交事务
func CommitTransaction(ctx context.Context) error {
	if sessCtx, ok := GetSessionFromContext(ctx); ok {
		return sessCtx.CommitTransaction(ctx)
	}
	return fmt.Errorf("not in transaction context")
}

// 事务隔离级别常量

// TransactionIsolationLevel 事务隔离级别
type TransactionIsolationLevel string

const (
	// IsolationLevelSnapshot Snapshot 隔离级别
	IsolationLevelSnapshot TransactionIsolationLevel = "snapshot"
	// IsolationLevelMajority Majority 隔离级别
	IsolationLevelMajority TransactionIsolationLevel = "majority"
	// IsolationLevelLocal Local 隔离级别
	IsolationLevelLocal TransactionIsolationLevel = "local"
)

// GetTransactionOptions 根据隔离级别获取事务选项
func GetTransactionOptions(level TransactionIsolationLevel) *options.TransactionOptions {
	switch level {
	case IsolationLevelSnapshot:
		return NewTxOptions().ReadConcernSnapshot().WriteConcernMajority().Build()
	case IsolationLevelMajority:
		return NewTxOptions().ReadConcernMajority().WriteConcernMajority().Build()
	case IsolationLevelLocal:
		return NewTxOptions().ReadConcernLocal().WriteConcernW1().Build()
	default:
		return options.Transaction()
	}
}

// WithIsolationLevel 使用指定隔离级别执行事务
func WithIsolationLevel(ctx context.Context, level TransactionIsolationLevel, fn TxFunc) error {
	opts := GetTransactionOptions(level)
	return RunInTransaction(ctx, fn, opts)
}
