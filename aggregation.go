package gmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AggregationBuilder 聚合管道构建器
type AggregationBuilder struct {
	collection *Collection
	pipeline   []bson.M
	opts       *options.AggregateOptions
	ctx        context.Context
}

// NewAggregation 创建聚合构建器
func NewAggregation(collection *Collection) *AggregationBuilder {
	return &AggregationBuilder{
		collection: collection,
		pipeline:   []bson.M{},
		opts:       options.Aggregate(),
		ctx:        context.Background(),
	}
}

// NewAggregationWithContext 使用上下文创建聚合构建器
func NewAggregationWithContext(ctx context.Context, collection *Collection) *AggregationBuilder {
	return NewAggregation(collection).Context(ctx)
}

// Context 设置上下文
func (ab *AggregationBuilder) Context(ctx context.Context) *AggregationBuilder {
	ab.ctx = ctx
	return ab
}

// Pipeline 设置完整管道
func (ab *AggregationBuilder) Pipeline(pipeline []bson.M) *AggregationBuilder {
	ab.pipeline = pipeline
	return ab
}

// AddStage 添加自定义阶段
func (ab *AggregationBuilder) AddStage(stage bson.M) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, stage)
	return ab
}

// Match 匹配阶段（过滤文档）
func (ab *AggregationBuilder) Match(filter bson.M) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$match": filter})
	return ab
}

// Project 投影阶段（选择字段）
func (ab *AggregationBuilder) Project(projection bson.M) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$project": projection})
	return ab
}

// Group 分组阶段
func (ab *AggregationBuilder) Group(id interface{}, fields bson.M) *AggregationBuilder {
	group := bson.M{"_id": id}
	for k, v := range fields {
		group[k] = v
	}
	ab.pipeline = append(ab.pipeline, bson.M{"$group": group})
	return ab
}

// Sort 排序阶段
func (ab *AggregationBuilder) Sort(sort bson.D) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$sort": sort})
	return ab
}

// SortBy 排序阶段（便捷方法）
func (ab *AggregationBuilder) SortBy(fields ...string) *AggregationBuilder {
	sort := bson.D{}
	for _, field := range fields {
		if len(field) > 0 && field[0] == '-' {
			sort = append(sort, bson.E{Key: field[1:], Value: -1})
		} else {
			sort = append(sort, bson.E{Key: field, Value: 1})
		}
	}
	return ab.Sort(sort)
}

// Limit 限制阶段
func (ab *AggregationBuilder) Limit(limit int64) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$limit": limit})
	return ab
}

// Skip 跳过阶段
func (ab *AggregationBuilder) Skip(skip int64) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$skip": skip})
	return ab
}

// Unwind 展开数组阶段
func (ab *AggregationBuilder) Unwind(path string, preserveNullAndEmptyArrays ...bool) *AggregationBuilder {
	unwind := bson.M{"path": path}
	if len(preserveNullAndEmptyArrays) > 0 && preserveNullAndEmptyArrays[0] {
		unwind["preserveNullAndEmptyArrays"] = true
	}
	ab.pipeline = append(ab.pipeline, bson.M{"$unwind": unwind})
	return ab
}

// Lookup 关联查询阶段（类似 SQL JOIN）
func (ab *AggregationBuilder) Lookup(from, localField, foreignField, as string) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{
		"$lookup": bson.M{
			"from":         from,
			"localField":   localField,
			"foreignField": foreignField,
			"as":           as,
		},
	})
	return ab
}

// LookupPipeline 关联查询阶段（使用管道）
func (ab *AggregationBuilder) LookupPipeline(from string, let bson.M, pipeline []bson.M, as string) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{
		"$lookup": bson.M{
			"from":     from,
			"let":      let,
			"pipeline": pipeline,
			"as":       as,
		},
	})
	return ab
}

// AddFields 添加字段阶段
func (ab *AggregationBuilder) AddFields(fields bson.M) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$addFields": fields})
	return ab
}

// ReplaceRoot 替换根文档阶段
func (ab *AggregationBuilder) ReplaceRoot(newRoot interface{}) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$replaceRoot": bson.M{"newRoot": newRoot}})
	return ab
}

// Sample 随机采样阶段
func (ab *AggregationBuilder) Sample(size int64) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$sample": bson.M{"size": size}})
	return ab
}

// Count 计数阶段
func (ab *AggregationBuilder) Count(field string) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$count": field})
	return ab
}

// Bucket 分桶阶段
func (ab *AggregationBuilder) Bucket(groupBy interface{}, boundaries []interface{}, defaultBucket interface{}, output bson.M) *AggregationBuilder {
	bucket := bson.M{
		"groupBy":    groupBy,
		"boundaries": boundaries,
	}
	if defaultBucket != nil {
		bucket["default"] = defaultBucket
	}
	if output != nil {
		bucket["output"] = output
	}
	ab.pipeline = append(ab.pipeline, bson.M{"$bucket": bucket})
	return ab
}

// BucketAuto 自动分桶阶段
func (ab *AggregationBuilder) BucketAuto(groupBy interface{}, buckets int, output bson.M) *AggregationBuilder {
	bucketAuto := bson.M{
		"groupBy": groupBy,
		"buckets": buckets,
	}
	if output != nil {
		bucketAuto["output"] = output
	}
	ab.pipeline = append(ab.pipeline, bson.M{"$bucketAuto": bucketAuto})
	return ab
}

// Facet 多面搜索阶段
func (ab *AggregationBuilder) Facet(facets bson.M) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$facet": facets})
	return ab
}

// GraphLookup 图查找阶段
func (ab *AggregationBuilder) GraphLookup(from, startWith, connectFromField, connectToField, as string, maxDepth ...int) *AggregationBuilder {
	graphLookup := bson.M{
		"from":             from,
		"startWith":        startWith,
		"connectFromField": connectFromField,
		"connectToField":   connectToField,
		"as":               as,
	}
	if len(maxDepth) > 0 {
		graphLookup["maxDepth"] = maxDepth[0]
	}
	ab.pipeline = append(ab.pipeline, bson.M{"$graphLookup": graphLookup})
	return ab
}

// Out 输出到集合阶段
func (ab *AggregationBuilder) Out(collection string) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$out": collection})
	return ab
}

// Merge 合并到集合阶段
func (ab *AggregationBuilder) Merge(into interface{}, on interface{}, whenMatched, whenNotMatched string) *AggregationBuilder {
	merge := bson.M{
		"into":           into,
		"whenMatched":    whenMatched,
		"whenNotMatched": whenNotMatched,
	}
	if on != nil {
		merge["on"] = on
	}
	ab.pipeline = append(ab.pipeline, bson.M{"$merge": merge})
	return ab
}

// Redact 编辑阶段
func (ab *AggregationBuilder) Redact(expression interface{}) *AggregationBuilder {
	ab.pipeline = append(ab.pipeline, bson.M{"$redact": expression})
	return ab
}

// GeoNear 地理位置查询阶段
func (ab *AggregationBuilder) GeoNear(near interface{}, distanceField string, spherical bool, opts bson.M) *AggregationBuilder {
	geoNear := bson.M{
		"near":          near,
		"distanceField": distanceField,
		"spherical":     spherical,
	}
	for k, v := range opts {
		geoNear[k] = v
	}
	ab.pipeline = append(ab.pipeline, bson.M{"$geoNear": geoNear})
	return ab
}

// 选项设置

// AllowDiskUse 允许使用磁盘
func (ab *AggregationBuilder) AllowDiskUse(allow bool) *AggregationBuilder {
	ab.opts.SetAllowDiskUse(allow)
	return ab
}

// BatchSize 设置批次大小
func (ab *AggregationBuilder) BatchSize(size int32) *AggregationBuilder {
	ab.opts.SetBatchSize(size)
	return ab
}

// MaxTime 设置最大执行时间
func (ab *AggregationBuilder) MaxTime(d time.Duration) *AggregationBuilder {
	ab.opts.SetMaxTime(d)
	return ab
}

// Collation 设置排序规则
func (ab *AggregationBuilder) Collation(collation *options.Collation) *AggregationBuilder {
	ab.opts.SetCollation(collation)
	return ab
}

// Hint 使用索引提示
func (ab *AggregationBuilder) Hint(hint interface{}) *AggregationBuilder {
	ab.opts.SetHint(hint)
	return ab
}

// Comment 添加注释
func (ab *AggregationBuilder) Comment(comment string) *AggregationBuilder {
	ab.opts.SetComment(comment)
	return ab
}

// 执行方法

// Execute 执行聚合
func (ab *AggregationBuilder) Execute(results interface{}) error {
	ctx := ab.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	cursor, err := ab.collection.Aggregate(ctx, ab.pipeline, ab.opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// Cursor 获取游标
func (ab *AggregationBuilder) Cursor() (*mongo.Cursor, error) {
	ctx := ab.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	return ab.collection.Aggregate(ctx, ab.pipeline, ab.opts)
}

// One 执行聚合并获取单个结果
func (ab *AggregationBuilder) One(result interface{}) error {
	ctx := ab.ctx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	cursor, err := ab.collection.Aggregate(ctx, ab.pipeline, ab.opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		return cursor.Decode(result)
	}

	return mongo.ErrNoDocuments
}

// GetPipeline 获取管道
func (ab *AggregationBuilder) GetPipeline() []bson.M {
	return ab.pipeline
}

// 便捷聚合函数

// Sum 求和
func Sum(field string) bson.M {
	return bson.M{"$sum": field}
}

// Avg 平均值
func Avg(field string) bson.M {
	return bson.M{"$avg": field}
}

// Min 最小值
func Min(field string) bson.M {
	return bson.M{"$min": field}
}

// Max 最大值
func Max(field string) bson.M {
	return bson.M{"$max": field}
}

// First 第一个值
func First(field string) bson.M {
	return bson.M{"$first": field}
}

// Last 最后一个值
func Last(field string) bson.M {
	return bson.M{"$last": field}
}

// Push 添加到数组
func Push(field string) bson.M {
	return bson.M{"$push": field}
}

// AddToSet 添加到集合（去重）
func AddToSet(field string) bson.M {
	return bson.M{"$addToSet": field}
}

// StdDevPop 总体标准差
func StdDevPop(field string) bson.M {
	return bson.M{"$stdDevPop": field}
}

// StdDevSamp 样本标准差
func StdDevSamp(field string) bson.M {
	return bson.M{"$stdDevSamp": field}
}
