package gmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexManager 索引管理器
type IndexManager struct {
	collection *Collection
	view       mongo.IndexView
}

// CreateOne 创建单个索引
func (im *IndexManager) CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return im.view.CreateOne(ctx, model, opts...)
}

// CreateMany 创建多个索引
func (im *IndexManager) CreateMany(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return im.view.CreateMany(ctx, models, opts...)
}

// DropOne 删除单个索引
func (im *IndexManager) DropOne(ctx context.Context, name string, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return im.view.DropOne(ctx, name, opts...)
}

// DropAll 删除所有索引
func (im *IndexManager) DropAll(ctx context.Context, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return im.view.DropAll(ctx, opts...)
}

// List 列出所有索引
func (im *IndexManager) List(ctx context.Context, opts ...*options.ListIndexesOptions) (*mongo.Cursor, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return im.view.List(ctx, opts...)
}

// ListSpecifications 列出所有索引规范
func (im *IndexManager) ListSpecifications(ctx context.Context, opts ...*options.ListIndexesOptions) ([]*mongo.IndexSpecification, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return im.view.ListSpecifications(ctx, opts...)
}

// 便捷方法

// CreateIndex 创建索引（便捷方法）
func (im *IndexManager) CreateIndex(ctx context.Context, keys bson.D, unique bool, opts ...*options.IndexOptions) (string, error) {
	indexOpts := options.Index()
	if unique {
		indexOpts.SetUnique(true)
	}
	if len(opts) > 0 {
		indexOpts = opts[0]
	}

	model := mongo.IndexModel{
		Keys:    keys,
		Options: indexOpts,
	}

	return im.CreateOne(ctx, model)
}

// CreateUniqueIndex 创建唯一索引
func (im *IndexManager) CreateUniqueIndex(ctx context.Context, keys bson.D, opts ...*options.IndexOptions) (string, error) {
	return im.CreateIndex(ctx, keys, true, opts...)
}

// CreateTextIndex 创建文本索引
func (im *IndexManager) CreateTextIndex(ctx context.Context, field string, opts ...*options.IndexOptions) (string, error) {
	keys := bson.D{{Key: field, Value: "text"}}
	return im.CreateIndex(ctx, keys, false, opts...)
}

// CreateTTLIndex 创建 TTL 索引
func (im *IndexManager) CreateTTLIndex(ctx context.Context, field string, expireAfter time.Duration, opts ...*options.IndexOptions) (string, error) {
	keys := bson.D{{Key: field, Value: 1}}
	indexOpts := options.Index().SetExpireAfterSeconds(int32(expireAfter.Seconds()))
	if len(opts) > 0 {
		indexOpts = opts[0].SetExpireAfterSeconds(int32(expireAfter.Seconds()))
	}
	return im.CreateIndex(ctx, keys, false, indexOpts)
}

// CreateCompoundIndex 创建复合索引
func (im *IndexManager) CreateCompoundIndex(ctx context.Context, fields map[string]int, unique bool, opts ...*options.IndexOptions) (string, error) {
	keys := bson.D{}
	for field, order := range fields {
		keys = append(keys, bson.E{Key: field, Value: order})
	}
	return im.CreateIndex(ctx, keys, unique, opts...)
}

// CreateGeoIndex 创建地理空间索引
func (im *IndexManager) CreateGeoIndex(ctx context.Context, field string, indexType string, opts ...*options.IndexOptions) (string, error) {
	keys := bson.D{{Key: field, Value: indexType}} // "2d" or "2dsphere"
	return im.CreateIndex(ctx, keys, false, opts...)
}

// DropByName 根据名称删除索引
func (im *IndexManager) DropByName(ctx context.Context, name string) error {
	_, err := im.DropOne(ctx, name)
	return err
}

// EnsureIndex 确保索引存在（如果不存在则创建）
func (im *IndexManager) EnsureIndex(ctx context.Context, keys bson.D, unique bool, opts ...*options.IndexOptions) (string, error) {
	// 检查索引是否存在
	specs, err := im.ListSpecifications(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list indexes: %w", err)
	}

	// 生成索引名称
	indexName := generateIndexName(keys)
	if len(opts) > 0 && opts[0].Name != nil {
		indexName = *opts[0].Name
	}

	// 检查是否已存在
	for _, spec := range specs {
		if spec.Name == indexName {
			return indexName, nil // 索引已存在
		}
	}

	// 创建索引
	return im.CreateIndex(ctx, keys, unique, opts...)
}

// GetIndexNames 获取所有索引名称
func (im *IndexManager) GetIndexNames(ctx context.Context) ([]string, error) {
	specs, err := im.ListSpecifications(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(specs))
	for _, spec := range specs {
		names = append(names, spec.Name)
	}
	return names, nil
}

// IndexExists 检查索引是否存在
func (im *IndexManager) IndexExists(ctx context.Context, name string) (bool, error) {
	names, err := im.GetIndexNames(ctx)
	if err != nil {
		return false, err
	}

	for _, n := range names {
		if n == name {
			return true, nil
		}
	}
	return false, nil
}

// generateIndexName 生成索引名称
func generateIndexName(keys bson.D) string {
	name := ""
	for i, key := range keys {
		if i > 0 {
			name += "_"
		}
		name += fmt.Sprintf("%s_%v", key.Key, key.Value)
	}
	return name
}

// IndexBuilder 索引构建器
type IndexBuilder struct {
	keys    bson.D
	options *options.IndexOptions
}

// NewIndexBuilder 创建索引构建器
func NewIndexBuilder() *IndexBuilder {
	return &IndexBuilder{
		keys:    bson.D{},
		options: options.Index(),
	}
}

// AddField 添加字段（1 升序，-1 降序）
func (ib *IndexBuilder) AddField(field string, order int) *IndexBuilder {
	ib.keys = append(ib.keys, bson.E{Key: field, Value: order})
	return ib
}

// AddAscending 添加升序字段
func (ib *IndexBuilder) AddAscending(field string) *IndexBuilder {
	return ib.AddField(field, 1)
}

// AddDescending 添加降序字段
func (ib *IndexBuilder) AddDescending(field string) *IndexBuilder {
	return ib.AddField(field, -1)
}

// AddText 添加文本字段
func (ib *IndexBuilder) AddText(field string) *IndexBuilder {
	ib.keys = append(ib.keys, bson.E{Key: field, Value: "text"})
	return ib
}

// AddGeo2D 添加 2D 地理空间字段
func (ib *IndexBuilder) AddGeo2D(field string) *IndexBuilder {
	ib.keys = append(ib.keys, bson.E{Key: field, Value: "2d"})
	return ib
}

// AddGeo2DSphere 添加 2DSphere 地理空间字段
func (ib *IndexBuilder) AddGeo2DSphere(field string) *IndexBuilder {
	ib.keys = append(ib.keys, bson.E{Key: field, Value: "2dsphere"})
	return ib
}

// Unique 设置唯一索引
func (ib *IndexBuilder) Unique(unique bool) *IndexBuilder {
	ib.options.SetUnique(unique)
	return ib
}

// Name 设置索引名称
func (ib *IndexBuilder) Name(name string) *IndexBuilder {
	ib.options.SetName(name)
	return ib
}

// Background 设置后台创建
func (ib *IndexBuilder) Background(background bool) *IndexBuilder {
	ib.options.SetBackground(background)
	return ib
}

// Sparse 设置稀疏索引
func (ib *IndexBuilder) Sparse(sparse bool) *IndexBuilder {
	ib.options.SetSparse(sparse)
	return ib
}

// ExpireAfterSeconds 设置 TTL（秒）
func (ib *IndexBuilder) ExpireAfterSeconds(seconds int32) *IndexBuilder {
	ib.options.SetExpireAfterSeconds(seconds)
	return ib
}

// ExpireAfter 设置 TTL（时间段）
func (ib *IndexBuilder) ExpireAfter(duration time.Duration) *IndexBuilder {
	ib.options.SetExpireAfterSeconds(int32(duration.Seconds()))
	return ib
}

// PartialFilterExpression 设置部分过滤表达式
func (ib *IndexBuilder) PartialFilterExpression(filter interface{}) *IndexBuilder {
	ib.options.SetPartialFilterExpression(filter)
	return ib
}

// Collation 设置排序规则
func (ib *IndexBuilder) Collation(collation *options.Collation) *IndexBuilder {
	ib.options.SetCollation(collation)
	return ib
}

// Build 构建索引模型
func (ib *IndexBuilder) Build() mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    ib.keys,
		Options: ib.options,
	}
}

// Create 创建索引
func (ib *IndexBuilder) Create(ctx context.Context, im *IndexManager) (string, error) {
	return im.CreateOne(ctx, ib.Build())
}
