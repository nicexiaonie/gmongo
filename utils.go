package gmongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// M 是 bson.M 的别名，便于使用
type M = bson.M

// D 是 bson.D 的别名，便于使用
type D = bson.D

// A 是 bson.A 的别名，便于使用
type A = bson.A

// E 是 bson.E 的别名，便于使用
type E = bson.E

// ObjectID 是 primitive.ObjectID 的别名
type ObjectID = primitive.ObjectID

// NewObjectID 创建新的 ObjectID
func NewObjectID() ObjectID {
	return primitive.NewObjectID()
}

// ObjectIDFromHex 从十六进制字符串创建 ObjectID
func ObjectIDFromHex(hex string) (ObjectID, error) {
	return primitive.ObjectIDFromHex(hex)
}

// IsValidObjectID 检查字符串是否为有效的 ObjectID
func IsValidObjectID(hex string) bool {
	_, err := primitive.ObjectIDFromHex(hex)
	return err == nil
}

// TimeToObjectID 将时间转换为 ObjectID（用于时间范围查询）
func TimeToObjectID(t time.Time) ObjectID {
	return primitive.NewObjectIDFromTimestamp(t)
}

// 便捷的 BSON 构建函数

// BuildFilter 构建过滤条件
func BuildFilter(conditions ...M) M {
	if len(conditions) == 0 {
		return M{}
	}
	if len(conditions) == 1 {
		return conditions[0]
	}
	return M{"$and": conditions}
}

// BuildUpdate 构建更新操作
func BuildUpdate(set M, unset ...string) M {
	update := M{}
	if len(set) > 0 {
		update["$set"] = set
	}
	if len(unset) > 0 {
		unsetMap := M{}
		for _, field := range unset {
			unsetMap[field] = ""
		}
		update["$unset"] = unsetMap
	}
	return update
}

// BuildSort 构建排序
func BuildSort(fields ...string) D {
	sort := D{}
	for _, field := range fields {
		if len(field) > 0 && field[0] == '-' {
			sort = append(sort, E{Key: field[1:], Value: -1})
		} else {
			sort = append(sort, E{Key: field, Value: 1})
		}
	}
	return sort
}

// BuildProjection 构建投影
func BuildProjection(include []string, exclude []string) M {
	projection := M{}
	for _, field := range include {
		projection[field] = 1
	}
	for _, field := range exclude {
		projection[field] = 0
	}
	return projection
}

// 查询条件构建器

// Eq 等于
func Eq(value interface{}) M {
	return M{"$eq": value}
}

// Ne 不等于
func Ne(value interface{}) M {
	return M{"$ne": value}
}

// Gt 大于
func Gt(value interface{}) M {
	return M{"$gt": value}
}

// Gte 大于等于
func Gte(value interface{}) M {
	return M{"$gte": value}
}

// Lt 小于
func Lt(value interface{}) M {
	return M{"$lt": value}
}

// Lte 小于等于
func Lte(value interface{}) M {
	return M{"$lte": value}
}

// In 在列表中
func In(values ...interface{}) M {
	return M{"$in": values}
}

// Nin 不在列表中
func Nin(values ...interface{}) M {
	return M{"$nin": values}
}

// Between 在范围内
func Between(min, max interface{}) M {
	return M{"$gte": min, "$lte": max}
}

// Exists 字段存在
func Exists(exists bool) M {
	return M{"$exists": exists}
}

// Regex 正则表达式
func Regex(pattern string, options ...string) M {
	regex := M{"$regex": pattern}
	if len(options) > 0 {
		regex["$options"] = options[0]
	}
	return regex
}

// Type 字段类型
func Type(bsonType interface{}) M {
	return M{"$type": bsonType}
}

// Size 数组大小
func Size(size int) M {
	return M{"$size": size}
}

// All 数组包含所有元素
func All(values ...interface{}) M {
	return M{"$all": values}
}

// ElemMatch 数组元素匹配
func ElemMatch(condition M) M {
	return M{"$elemMatch": condition}
}

// Or 或条件
func Or(conditions ...M) M {
	return M{"$or": conditions}
}

// And 与条件
func And(conditions ...M) M {
	return M{"$and": conditions}
}

// Nor 非或条件
func Nor(conditions ...M) M {
	return M{"$nor": conditions}
}

// Not 非条件
func Not(condition M) M {
	return M{"$not": condition}
}

// 更新操作构建器

// Set 设置字段值
func Set(fields M) M {
	return M{"$set": fields}
}

// Unset 删除字段
func Unset(fields ...string) M {
	unset := M{}
	for _, field := range fields {
		unset[field] = ""
	}
	return M{"$unset": unset}
}

// Inc 增加数值
func Inc(fields M) M {
	return M{"$inc": fields}
}

// Mul 乘以数值
func Mul(fields M) M {
	return M{"$mul": fields}
}

// Rename 重命名字段
func Rename(fields M) M {
	return M{"$rename": fields}
}

// SetOnInsert 插入时设置（仅 upsert）
func SetOnInsert(fields M) M {
	return M{"$setOnInsert": fields}
}

// CurrentDate 设置为当前日期
func CurrentDate(fields ...string) M {
	current := M{}
	for _, field := range fields {
		current[field] = true
	}
	return M{"$currentDate": current}
}

// 数组更新操作

// AddToSetOp 添加到集合（去重）
func AddToSetOp(field string, values ...interface{}) M {
	if len(values) == 1 {
		return M{"$addToSet": M{field: values[0]}}
	}
	return M{"$addToSet": M{field: M{"$each": values}}}
}

// PushOp 添加到数组
func PushOp(field string, values ...interface{}) M {
	if len(values) == 1 {
		return M{"$push": M{field: values[0]}}
	}
	return M{"$push": M{field: M{"$each": values}}}
}

// Pull 从数组中删除
func Pull(field string, condition interface{}) M {
	return M{"$pull": M{field: condition}}
}

// PullAll 从数组中删除多个值
func PullAll(field string, values ...interface{}) M {
	return M{"$pullAll": M{field: values}}
}

// Pop 删除数组第一个或最后一个元素（1 最后，-1 第一个）
func Pop(field string, position int) M {
	return M{"$pop": M{field: position}}
}

// 工具函数

// ToDoc 将结构体转换为 bson.M
func ToDoc(v interface{}) (M, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}
	var doc M
	err = bson.Unmarshal(data, &doc)
	return doc, err
}

// FromDoc 将 bson.M 转换为结构体
func FromDoc(doc M, v interface{}) error {
	data, err := bson.Marshal(doc)
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, v)
}

// MustToDoc 将结构体转换为 bson.M（panic on error）
func MustToDoc(v interface{}) M {
	doc, err := ToDoc(v)
	if err != nil {
		panic(err)
	}
	return doc
}

// PrintDoc 打印文档（用于调试）
func PrintDoc(doc interface{}) {
	data, err := bson.MarshalExtJSON(doc, true, true)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// 分页辅助

// PageInfo 分页信息
type PageInfo struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int64 `json:"pages"`
}

// NewPageInfo 创建分页信息
func NewPageInfo(page, pageSize, total int64) *PageInfo {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	pages := total / pageSize
	if total%pageSize > 0 {
		pages++
	}
	return &PageInfo{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Pages:    pages,
	}
}

// PageResult 分页结果
type PageResult struct {
	Data     interface{} `json:"data"`
	PageInfo *PageInfo   `json:"page_info"`
}

// NewPageResult 创建分页结果
func NewPageResult(data interface{}, page, pageSize, total int64) *PageResult {
	return &PageResult{
		Data:     data,
		PageInfo: NewPageInfo(page, pageSize, total),
	}
}

// PaginateQuery 分页查询辅助函数
func PaginateQuery(ctx context.Context, collection *Collection, filter interface{}, page, pageSize int64, results interface{}) (*PageResult, error) {
	total, err := collection.Paginate(ctx, filter, page, pageSize, results)
	if err != nil {
		return nil, err
	}
	return NewPageResult(results, page, pageSize, total), nil
}

// 时间范围查询

// TimeRange 时间范围
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// NewTimeRange 创建时间范围
func NewTimeRange(start, end time.Time) *TimeRange {
	return &TimeRange{Start: start, End: end}
}

// Today 今天
func Today() *TimeRange {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	return &TimeRange{Start: start, End: end}
}

// Yesterday 昨天
func Yesterday() *TimeRange {
	now := time.Now().Add(-24 * time.Hour)
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	return &TimeRange{Start: start, End: end}
}

// ThisWeek 本周
func ThisWeek() *TimeRange {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := now.Add(-time.Duration(weekday-1) * 24 * time.Hour)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end := start.Add(7 * 24 * time.Hour)
	return &TimeRange{Start: start, End: end}
}

// ThisMonth 本月
func ThisMonth() *TimeRange {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0)
	return &TimeRange{Start: start, End: end}
}

// ToFilter 转换为过滤条件
func (tr *TimeRange) ToFilter(field string) M {
	return M{field: M{"$gte": tr.Start, "$lt": tr.End}}
}
