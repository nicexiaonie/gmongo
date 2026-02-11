package gmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// QueryBuilder 查询构建器
type QueryBuilder struct {
	collection *Collection
	filter     bson.M
	opts       *options.FindOptions
	ctx        context.Context
}

// Context 设置上下文
func (qb *QueryBuilder) Context(ctx context.Context) *QueryBuilder {
	qb.ctx = ctx
	return qb
}

// Filter 设置过滤条件
func (qb *QueryBuilder) Filter(filter interface{}) *QueryBuilder {
	if f, ok := filter.(bson.M); ok {
		qb.filter = f
	}
	return qb
}

// Where 添加查询条件
func (qb *QueryBuilder) Where(key string, value interface{}) *QueryBuilder {
	qb.filter[key] = value
	return qb
}

// WhereEq 等于条件
func (qb *QueryBuilder) WhereEq(key string, value interface{}) *QueryBuilder {
	qb.filter[key] = value
	return qb
}

// WhereNe 不等于条件
func (qb *QueryBuilder) WhereNe(key string, value interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$ne": value}
	return qb
}

// WhereGt 大于条件
func (qb *QueryBuilder) WhereGt(key string, value interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$gt": value}
	return qb
}

// WhereGte 大于等于条件
func (qb *QueryBuilder) WhereGte(key string, value interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$gte": value}
	return qb
}

// WhereLt 小于条件
func (qb *QueryBuilder) WhereLt(key string, value interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$lt": value}
	return qb
}

// WhereLte 小于等于条件
func (qb *QueryBuilder) WhereLte(key string, value interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$lte": value}
	return qb
}

// WhereIn 在列表中
func (qb *QueryBuilder) WhereIn(key string, values interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$in": values}
	return qb
}

// WhereNin 不在列表中
func (qb *QueryBuilder) WhereNin(key string, values interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$nin": values}
	return qb
}

// WhereBetween 在范围内
func (qb *QueryBuilder) WhereBetween(key string, min, max interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$gte": min, "$lte": max}
	return qb
}

// WhereExists 字段存在
func (qb *QueryBuilder) WhereExists(key string, exists bool) *QueryBuilder {
	qb.filter[key] = bson.M{"$exists": exists}
	return qb
}

// WhereRegex 正则表达式匹配
func (qb *QueryBuilder) WhereRegex(key string, pattern string, options ...string) *QueryBuilder {
	regex := bson.M{"$regex": pattern}
	if len(options) > 0 {
		regex["$options"] = options[0]
	}
	qb.filter[key] = regex
	return qb
}

// WhereType 字段类型匹配
func (qb *QueryBuilder) WhereType(key string, bsonType interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$type": bsonType}
	return qb
}

// WhereSize 数组大小匹配
func (qb *QueryBuilder) WhereSize(key string, size int) *QueryBuilder {
	qb.filter[key] = bson.M{"$size": size}
	return qb
}

// WhereAll 数组包含所有元素
func (qb *QueryBuilder) WhereAll(key string, values interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$all": values}
	return qb
}

// WhereElemMatch 数组元素匹配
func (qb *QueryBuilder) WhereElemMatch(key string, condition interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$elemMatch": condition}
	return qb
}

// Or 或条件
func (qb *QueryBuilder) Or(conditions ...bson.M) *QueryBuilder {
	if len(conditions) > 0 {
		qb.filter["$or"] = conditions
	}
	return qb
}

// And 与条件
func (qb *QueryBuilder) And(conditions ...bson.M) *QueryBuilder {
	if len(conditions) > 0 {
		qb.filter["$and"] = conditions
	}
	return qb
}

// Nor 非或条件
func (qb *QueryBuilder) Nor(conditions ...bson.M) *QueryBuilder {
	if len(conditions) > 0 {
		qb.filter["$nor"] = conditions
	}
	return qb
}

// Not 非条件
func (qb *QueryBuilder) Not(key string, condition interface{}) *QueryBuilder {
	qb.filter[key] = bson.M{"$not": condition}
	return qb
}

// Select 选择字段
func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	projection := bson.M{}
	for _, field := range fields {
		projection[field] = 1
	}
	qb.opts.SetProjection(projection)
	return qb
}

// Omit 排除字段
func (qb *QueryBuilder) Omit(fields ...string) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	projection := bson.M{}
	for _, field := range fields {
		projection[field] = 0
	}
	qb.opts.SetProjection(projection)
	return qb
}

// Sort 排序
func (qb *QueryBuilder) Sort(fields ...string) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	sort := bson.D{}
	for _, field := range fields {
		if len(field) > 0 && field[0] == '-' {
			sort = append(sort, bson.E{Key: field[1:], Value: -1})
		} else {
			sort = append(sort, bson.E{Key: field, Value: 1})
		}
	}
	qb.opts.SetSort(sort)
	return qb
}

// Limit 限制数量
func (qb *QueryBuilder) Limit(limit int64) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	qb.opts.SetLimit(limit)
	return qb
}

// Skip 跳过数量
func (qb *QueryBuilder) Skip(skip int64) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	qb.opts.SetSkip(skip)
	return qb
}

// Page 分页（从 1 开始）
func (qb *QueryBuilder) Page(page, pageSize int64) *QueryBuilder {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	skip := (page - 1) * pageSize
	return qb.Skip(skip).Limit(pageSize)
}

// Hint 使用索引提示
func (qb *QueryBuilder) Hint(hint interface{}) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	qb.opts.SetHint(hint)
	return qb
}

// MaxTime 设置最大执行时间
func (qb *QueryBuilder) MaxTime(d time.Duration) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	qb.opts.SetMaxTime(d)
	return qb
}

// AllowDiskUse 允许使用磁盘
func (qb *QueryBuilder) AllowDiskUse(allow bool) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	qb.opts.SetAllowDiskUse(allow)
	return qb
}

// Collation 设置排序规则
func (qb *QueryBuilder) Collation(collation *options.Collation) *QueryBuilder {
	if qb.opts == nil {
		qb.opts = options.Find()
	}
	qb.opts.SetCollation(collation)
	return qb
}

// 执行方法

// Find 执行查询
func (qb *QueryBuilder) Find(results interface{}) error {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	cursor, err := qb.collection.Find(ctx, qb.filter, qb.opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// FindOne 查询单个文档
func (qb *QueryBuilder) FindOne(result interface{}) error {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	opts := options.FindOne()
	if qb.opts != nil {
		if qb.opts.Projection != nil {
			opts.SetProjection(qb.opts.Projection)
		}
		if qb.opts.Sort != nil {
			opts.SetSort(qb.opts.Sort)
		}
		if qb.opts.Skip != nil {
			opts.SetSkip(*qb.opts.Skip)
		}
		if qb.opts.Hint != nil {
			opts.SetHint(qb.opts.Hint)
		}
		if qb.opts.MaxTime != nil {
			opts.SetMaxTime(*qb.opts.MaxTime)
		}
		if qb.opts.Collation != nil {
			opts.SetCollation(qb.opts.Collation)
		}
	}

	return qb.collection.FindOne(ctx, qb.filter, opts).Decode(result)
}

// Count 统计数量
func (qb *QueryBuilder) Count() (int64, error) {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	opts := options.Count()
	if qb.opts != nil {
		if qb.opts.Limit != nil {
			opts.SetLimit(*qb.opts.Limit)
		}
		if qb.opts.Skip != nil {
			opts.SetSkip(*qb.opts.Skip)
		}
		if qb.opts.Hint != nil {
			opts.SetHint(qb.opts.Hint)
		}
		if qb.opts.MaxTime != nil {
			opts.SetMaxTime(*qb.opts.MaxTime)
		}
		if qb.opts.Collation != nil {
			opts.SetCollation(qb.opts.Collation)
		}
	}

	return qb.collection.CountDocuments(ctx, qb.filter, opts)
}

// Exists 检查是否存在
func (qb *QueryBuilder) Exists() (bool, error) {
	count, err := qb.Limit(1).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update 更新文档
func (qb *QueryBuilder) Update(update interface{}) (*mongo.UpdateResult, error) {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return qb.collection.UpdateOne(ctx, qb.filter, update)
}

// UpdateMany 更新多个文档
func (qb *QueryBuilder) UpdateMany(update interface{}) (*mongo.UpdateResult, error) {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return qb.collection.UpdateMany(ctx, qb.filter, update)
}

// Delete 删除文档
func (qb *QueryBuilder) Delete() (*mongo.DeleteResult, error) {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return qb.collection.DeleteOne(ctx, qb.filter)
}

// DeleteMany 删除多个文档
func (qb *QueryBuilder) DeleteMany() (*mongo.DeleteResult, error) {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return qb.collection.DeleteMany(ctx, qb.filter)
}

// Distinct 获取不重复的值
func (qb *QueryBuilder) Distinct(fieldName string) ([]interface{}, error) {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return qb.collection.Distinct(ctx, fieldName, qb.filter)
}

// Cursor 获取游标
func (qb *QueryBuilder) Cursor() (*mongo.Cursor, error) {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return qb.collection.Find(ctx, qb.filter, qb.opts)
}

// Paginate 分页查询（返回总数和数据）
func (qb *QueryBuilder) Paginate(page, pageSize int64, results interface{}) (int64, error) {
	ctx := qb.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// 计算总数
	total, err := qb.Count()
	if err != nil {
		return 0, err
	}

	// 分页查询
	qb.Page(page, pageSize)
	if err := qb.Find(results); err != nil {
		return 0, err
	}

	return total, nil
}
