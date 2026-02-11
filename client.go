package gmongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client MongoDB 客户端
type Client struct {
	client *mongo.Client
	config *Config
	mu     sync.RWMutex
	closed bool
}

// NewClient 创建客户端
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 创建客户端选项
	clientOpts := config.ToClientOptions()

	// 连接 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Client{
		client: client,
		config: config,
	}, nil
}

// Client 获取原生 mongo.Client
func (c *Client) Client() *mongo.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.client
}

// Config 获取配置
func (c *Client) Config() *Config {
	return c.config
}

// Database 获取数据库
func (c *Client) Database(name ...string) *Database {
	dbName := c.config.Database
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}

	return &Database{
		client: c,
		db:     c.client.Database(dbName),
		name:   dbName,
	}
}

// Collection 获取集合（使用默认数据库）
func (c *Client) Collection(name string) *Collection {
	return c.Database().Collection(name)
}

// Close 关闭连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c.closed = true
	return c.client.Disconnect(ctx)
}

// IsClosed 是否已关闭
func (c *Client) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// Ping 测试连接
func (c *Client) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("client is closed")
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	return c.client.Ping(ctx, nil)
}

// StartSession 开始会话
func (c *Client) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	return c.client.StartSession(opts...)
}

// UseSession 使用会话执行操作
func (c *Client) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return c.client.UseSession(ctx, fn)
}

// UseSessionWithOptions 使用会话和选项执行操作
func (c *Client) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	return c.client.UseSessionWithOptions(ctx, opts, fn)
}

// ListDatabases 列出所有数据库
func (c *Client) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return mongo.ListDatabasesResult{}, fmt.Errorf("client is closed")
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return c.client.ListDatabases(ctx, filter, opts...)
}

// ListDatabaseNames 列出所有数据库名称
func (c *Client) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return c.client.ListDatabaseNames(ctx, filter, opts...)
}

// Watch 监听变更流
func (c *Client) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("client is closed")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return c.client.Watch(ctx, pipeline, opts...)
}

// HealthCheck 健康检查
func (c *Client) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.Ping(ctx)
}

// Stats 获取连接统计信息
func (c *Client) Stats() string {
	// MongoDB driver 不直接提供连接池统计，这里返回基本信息
	return fmt.Sprintf("Database: %s, Closed: %v", c.config.Database, c.closed)
}

// Database 数据库封装
type Database struct {
	client *Client
	db     *mongo.Database
	name   string
}

// Name 获取数据库名称
func (d *Database) Name() string {
	return d.name
}

// DB 获取原生 mongo.Database
func (d *Database) DB() *mongo.Database {
	return d.db
}

// Collection 获取集合
func (d *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return &Collection{
		client:     d.client,
		db:         d,
		collection: d.db.Collection(name, opts...),
		name:       name,
	}
}

// RunCommand 执行数据库命令
func (d *Database) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *mongo.SingleResult {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return d.db.RunCommand(ctx, runCommand, opts...)
}

// Drop 删除数据库
func (d *Database) Drop(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return d.db.Drop(ctx)
}

// ListCollections 列出所有集合
func (d *Database) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*mongo.Cursor, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return d.db.ListCollections(ctx, filter, opts...)
}

// ListCollectionNames 列出所有集合名称
func (d *Database) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return d.db.ListCollectionNames(ctx, filter, opts...)
}

// CreateCollection 创建集合
func (d *Database) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return d.db.CreateCollection(ctx, name, opts...)
}

// Watch 监听数据库变更流
func (d *Database) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return d.db.Watch(ctx, pipeline, opts...)
}

// Aggregate 执行聚合操作
func (d *Database) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return d.db.Aggregate(ctx, pipeline, opts...)
}
