package gmongo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

// HookFunc 钩子函数类型
type HookFunc func(ctx context.Context, operation string, args ...interface{}) error

// HookRegistry 钩子注册器
type HookRegistry struct {
	beforeHooks map[string][]HookFunc
	afterHooks  map[string][]HookFunc
	mu          sync.RWMutex
}

// NewHookRegistry 创建钩子注册器
func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		beforeHooks: make(map[string][]HookFunc),
		afterHooks:  make(map[string][]HookFunc),
	}
}

// RegisterBefore 注册前置钩子
func (hr *HookRegistry) RegisterBefore(operation string, hook HookFunc) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	hr.beforeHooks[operation] = append(hr.beforeHooks[operation], hook)
}

// RegisterAfter 注册后置钩子
func (hr *HookRegistry) RegisterAfter(operation string, hook HookFunc) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	hr.afterHooks[operation] = append(hr.afterHooks[operation], hook)
}

// ExecuteBefore 执行前置钩子
func (hr *HookRegistry) ExecuteBefore(ctx context.Context, operation string, args ...interface{}) error {
	hr.mu.RLock()
	hooks := hr.beforeHooks[operation]
	hr.mu.RUnlock()

	for _, hook := range hooks {
		if err := hook(ctx, operation, args...); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteAfter 执行后置钩子
func (hr *HookRegistry) ExecuteAfter(ctx context.Context, operation string, args ...interface{}) error {
	hr.mu.RLock()
	hooks := hr.afterHooks[operation]
	hr.mu.RUnlock()

	for _, hook := range hooks {
		if err := hook(ctx, operation, args...); err != nil {
			return err
		}
	}
	return nil
}

// RemoveBefore 移除前置钩子
func (hr *HookRegistry) RemoveBefore(operation string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	delete(hr.beforeHooks, operation)
}

// RemoveAfter 移除后置钩子
func (hr *HookRegistry) RemoveAfter(operation string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	delete(hr.afterHooks, operation)
}

// Clear 清除所有钩子
func (hr *HookRegistry) Clear() {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	hr.beforeHooks = make(map[string][]HookFunc)
	hr.afterHooks = make(map[string][]HookFunc)
}

// 全局钩子注册器
var globalHookRegistry = NewHookRegistry()

// RegisterBeforeInsert 注册插入前钩子
func RegisterBeforeInsert(hook HookFunc) {
	globalHookRegistry.RegisterBefore("insert", hook)
}

// RegisterAfterInsert 注册插入后钩子
func RegisterAfterInsert(hook HookFunc) {
	globalHookRegistry.RegisterAfter("insert", hook)
}

// RegisterBeforeUpdate 注册更新前钩子
func RegisterBeforeUpdate(hook HookFunc) {
	globalHookRegistry.RegisterBefore("update", hook)
}

// RegisterAfterUpdate 注册更新后钩子
func RegisterAfterUpdate(hook HookFunc) {
	globalHookRegistry.RegisterAfter("update", hook)
}

// RegisterBeforeDelete 注册删除前钩子
func RegisterBeforeDelete(hook HookFunc) {
	globalHookRegistry.RegisterBefore("delete", hook)
}

// RegisterAfterDelete 注册删除后钩子
func RegisterAfterDelete(hook HookFunc) {
	globalHookRegistry.RegisterAfter("delete", hook)
}

// RegisterBeforeFind 注册查询前钩子
func RegisterBeforeFind(hook HookFunc) {
	globalHookRegistry.RegisterBefore("find", hook)
}

// RegisterAfterFind 注册查询后钩子
func RegisterAfterFind(hook HookFunc) {
	globalHookRegistry.RegisterAfter("find", hook)
}

// RegisterBeforeAggregate 注册聚合前钩子
func RegisterBeforeAggregate(hook HookFunc) {
	globalHookRegistry.RegisterBefore("aggregate", hook)
}

// RegisterAfterAggregate 注册聚合后钩子
func RegisterAfterAggregate(hook HookFunc) {
	globalHookRegistry.RegisterAfter("aggregate", hook)
}

// GetGlobalHookRegistry 获取全局钩子注册器
func GetGlobalHookRegistry() *HookRegistry {
	return globalHookRegistry
}

// 钩子操作常量
const (
	OpInsert    = "insert"
	OpUpdate    = "update"
	OpDelete    = "delete"
	OpFind      = "find"
	OpAggregate = "aggregate"
	OpCount     = "count"
	OpDistinct  = "distinct"
)

// CollectionWithHooks 带钩子的集合封装
type CollectionWithHooks struct {
	*Collection
	hooks *HookRegistry
}

// NewCollectionWithHooks 创建带钩子的集合
func NewCollectionWithHooks(collection *Collection, hooks *HookRegistry) *CollectionWithHooks {
	if hooks == nil {
		hooks = NewHookRegistry()
	}
	return &CollectionWithHooks{
		Collection: collection,
		hooks:      hooks,
	}
}

// InsertOne 插入单个文档（带钩子）
func (c *CollectionWithHooks) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if err := c.hooks.ExecuteBefore(ctx, OpInsert, document); err != nil {
		return nil, err
	}

	result, err := c.Collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	if err := c.hooks.ExecuteAfter(ctx, OpInsert, document, result); err != nil {
		return result, err
	}

	return result, nil
}

// UpdateOne 更新单个文档（带钩子）
func (c *CollectionWithHooks) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	if err := c.hooks.ExecuteBefore(ctx, OpUpdate, filter, update); err != nil {
		return nil, err
	}

	result, err := c.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if err := c.hooks.ExecuteAfter(ctx, OpUpdate, filter, update, result); err != nil {
		return result, err
	}

	return result, nil
}

// DeleteOne 删除单个文档（带钩子）
func (c *CollectionWithHooks) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	if err := c.hooks.ExecuteBefore(ctx, OpDelete, filter); err != nil {
		return nil, err
	}

	result, err := c.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := c.hooks.ExecuteAfter(ctx, OpDelete, filter, result); err != nil {
		return result, err
	}

	return result, nil
}

// FindOne 查询单个文档（带钩子）
func (c *CollectionWithHooks) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	_ = c.hooks.ExecuteBefore(ctx, OpFind, filter)
	result := c.Collection.FindOne(ctx, filter)
	_ = c.hooks.ExecuteAfter(ctx, OpFind, filter, result)
	return result
}

// RegisterBefore 注册前置钩子
func (c *CollectionWithHooks) RegisterBefore(operation string, hook HookFunc) {
	c.hooks.RegisterBefore(operation, hook)
}

// RegisterAfter 注册后置钩子
func (c *CollectionWithHooks) RegisterAfter(operation string, hook HookFunc) {
	c.hooks.RegisterAfter(operation, hook)
}

// 便捷钩子函数

// LoggingHook 日志钩子
func LoggingHook(logger func(format string, args ...interface{})) HookFunc {
	return func(ctx context.Context, operation string, args ...interface{}) error {
		logger("Operation: %s, Args: %v", operation, args)
		return nil
	}
}

// TimingHook 计时钩子
func TimingHook(onComplete func(operation string, duration int64)) HookFunc {
	return func(ctx context.Context, operation string, args ...interface{}) error {
		// 这里需要配合 after hook 使用来计算时间
		return nil
	}
}

// ValidationHook 验证钩子
func ValidationHook(validator func(ctx context.Context, operation string, args ...interface{}) error) HookFunc {
	return validator
}

// CacheInvalidationHook 缓存失效钩子
func CacheInvalidationHook(invalidate func(ctx context.Context, operation string, args ...interface{}) error) HookFunc {
	return invalidate
}
