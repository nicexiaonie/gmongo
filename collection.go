package gmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection 集合封装
type Collection struct {
	client     *Client
	db         *Database
	collection *mongo.Collection
	name       string
}

// Name 获取集合名称
func (c *Collection) Name() string {
	return c.name
}

// Collection 获取原生 mongo.Collection
func (c *Collection) Collection() *mongo.Collection {
	return c.collection
}

// Query 创建查询构建器
func (c *Collection) Query() *QueryBuilder {
	return &QueryBuilder{
		collection: c,
		filter:     bson.M{},
		ctx:        context.Background(),
	}
}

// InsertOne 插入单个文档
func (c *Collection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.InsertOne(ctx, document, opts...)
}

// InsertMany 插入多个文档
func (c *Collection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.InsertMany(ctx, documents, opts...)
}

// FindOne 查询单个文档
func (c *Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.FindOne(ctx, filter, opts...)
}

// Find 查询多个文档
func (c *Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.Find(ctx, filter, opts...)
}

// UpdateOne 更新单个文档
func (c *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.UpdateOne(ctx, filter, update, opts...)
}

// UpdateMany 更新多个文档
func (c *Collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.UpdateMany(ctx, filter, update, opts...)
}

// ReplaceOne 替换单个文档
func (c *Collection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.ReplaceOne(ctx, filter, replacement, opts...)
}

// DeleteOne 删除单个文档
func (c *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.DeleteOne(ctx, filter, opts...)
}

// DeleteMany 删除多个文档
func (c *Collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.DeleteMany(ctx, filter, opts...)
}

// CountDocuments 统计文档数量
func (c *Collection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.CountDocuments(ctx, filter, opts...)
}

// EstimatedDocumentCount 估算文档数量（快速）
func (c *Collection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.EstimatedDocumentCount(ctx, opts...)
}

// Distinct 获取不重复的值
func (c *Collection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.Distinct(ctx, fieldName, filter, opts...)
}

// Aggregate 执行聚合操作
func (c *Collection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.Aggregate(ctx, pipeline, opts...)
}

// BulkWrite 批量写入操作
func (c *Collection) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.BulkWrite(ctx, models, opts...)
}

// Watch 监听集合变更流
func (c *Collection) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return c.collection.Watch(ctx, pipeline, opts...)
}

// Drop 删除集合
func (c *Collection) Drop(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.Drop(ctx)
}

// Clone 克隆集合（使用新选项）
func (c *Collection) Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	return c.collection.Clone(opts...)
}

// FindOneAndUpdate 查找并更新单个文档
func (c *Collection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.FindOneAndUpdate(ctx, filter, update, opts...)
}

// FindOneAndReplace 查找并替换单个文档
func (c *Collection) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.FindOneAndReplace(ctx, filter, replacement, opts...)
}

// FindOneAndDelete 查找并删除单个文档
func (c *Collection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}
	return c.collection.FindOneAndDelete(ctx, filter, opts...)
}

// Indexes 获取索引管理器
func (c *Collection) Indexes() *IndexManager {
	return &IndexManager{
		collection: c,
		view:       c.collection.Indexes(),
	}
}

// 便捷方法

// Insert 插入文档（自动判断单个或多个）
func (c *Collection) Insert(ctx context.Context, documents interface{}, opts ...interface{}) (interface{}, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// 判断是否为切片
	switch docs := documents.(type) {
	case []interface{}:
		return c.InsertMany(ctx, docs)
	default:
		return c.InsertOne(ctx, documents)
	}
}

// Update 更新文档（便捷方法）
func (c *Collection) Update(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.UpdateOne(ctx, filter, update, opts...)
}

// Delete 删除文档（便捷方法）
func (c *Collection) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.DeleteOne(ctx, filter, opts...)
}

// Count 统计文档数量（便捷方法）
func (c *Collection) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return c.CountDocuments(ctx, filter, opts...)
}

// Exists 检查文档是否存在
func (c *Collection) Exists(ctx context.Context, filter interface{}) (bool, error) {
	count, err := c.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll 查询所有文档并解码到切片
func (c *Collection) FindAll(ctx context.Context, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	cursor, err := c.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// FindOneByID 根据 ID 查询单个文档
func (c *Collection) FindOneByID(ctx context.Context, id interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	filter := bson.M{"_id": id}
	return c.FindOne(ctx, filter, opts...).Decode(result)
}

// UpdateByID 根据 ID 更新文档
func (c *Collection) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	filter := bson.M{"_id": id}
	return c.UpdateOne(ctx, filter, update, opts...)
}

// DeleteByID 根据 ID 删除文档
func (c *Collection) DeleteByID(ctx context.Context, id interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	filter := bson.M{"_id": id}
	return c.DeleteOne(ctx, filter, opts...)
}

// Upsert 更新或插入文档
func (c *Collection) Upsert(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	opts := options.Update().SetUpsert(true)
	return c.UpdateOne(ctx, filter, update, opts)
}

// Paginate 分页查询
func (c *Collection) Paginate(ctx context.Context, filter interface{}, page, pageSize int64, results interface{}, opts ...*options.FindOptions) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 计算总数
	total, err := c.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// 计算跳过数量
	skip := (page - 1) * pageSize

	// 设置分页选项
	findOpts := options.Find().SetSkip(skip).SetLimit(pageSize)
	if len(opts) > 0 {
		findOpts = opts[0].SetSkip(skip).SetLimit(pageSize)
	}

	// 查询数据
	cursor, err := c.Find(ctx, filter, findOpts)
	if err != nil {
		return 0, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, results); err != nil {
		return 0, fmt.Errorf("failed to decode documents: %w", err)
	}

	return total, nil
}
