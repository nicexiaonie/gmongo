# GMongo 设计文档

## 设计理念

GMongo 是一个现代化的 MongoDB 工具库，借鉴了 GRDS 的优秀设计理念，专为 Go 开发者打造。

### 核心设计原则

1. **简洁易用** - API 设计简洁直观，降低学习成本
2. **链式调用** - 支持流畅的链式调用，代码更优雅
3. **类型安全** - 充分利用 Go 的类型系统，减少运行时错误
4. **高性能** - 基于官方 MongoDB Driver，性能优异
5. **功能全面** - 覆盖 MongoDB 的所有核心功能

## 架构设计

```
gmongo/
├── config.go          # 配置管理
├── client.go          # 客户端封装
├── collection.go      # 集合操作
├── query.go           # 查询构建器
├── aggregation.go     # 聚合管道
├── index.go           # 索引管理
├── transaction.go     # 事务支持
├── hooks.go           # 钩子系统
├── utils.go           # 工具函数
├── gmongo.go          # 全局实例
└── README.md          # 文档
```

## 核心组件

### 1. Config - 配置管理

提供完整的 MongoDB 连接配置选项：

- 基础连接配置（Host, Port, URI）
- 认证配置（Username, Password, AuthSource）
- 连接池配置（MaxPoolSize, MinPoolSize）
- 超时配置（ConnectTimeout, SocketTimeout）
- 读写配置（ReadPreference, WriteConcern）
- TLS 配置
- 压缩配置

**设计亮点**：
- 提供 `DefaultConfig()` 快速开始
- 支持 URI 和参数两种配置方式
- 自动验证配置有效性
- 智能转换为 MongoDB Driver 选项

### 2. Client - 客户端封装

封装 `mongo.Client`，提供线程安全的客户端管理：

- 连接管理（Connect, Close, Ping）
- 数据库访问（Database, Collection）
- 会话管理（StartSession, UseSession）
- 健康检查（HealthCheck, Stats）

**设计亮点**：
- 使用 `sync.RWMutex` 保证并发安全
- 自动管理连接生命周期
- 支持多数据库实例
- 提供全局默认客户端

### 3. Collection - 集合操作

封装 `mongo.Collection`，提供便捷的 CRUD 操作：

- 基础操作（Insert, Find, Update, Delete）
- 便捷方法（FindByID, Exists, Count, Paginate）
- 批量操作（InsertMany, BulkWrite）
- 原子操作（FindOneAndUpdate, Upsert）

**设计亮点**：
- 自动处理 context 超时
- 提供便捷的分页查询
- 支持链式查询构建器
- 完善的错误处理

### 4. QueryBuilder - 查询构建器

流畅的链式查询 API，类似 GRDS 的设计：

```go
users.Query().
    WhereGt("age", 18).
    WhereLt("age", 60).
    WhereIn("status", []string{"active"}).
    Select("name", "email").
    Sort("-age").
    Limit(10).
    Find(&results)
```

**支持的查询条件**：
- 比较操作：Eq, Ne, Gt, Gte, Lt, Lte
- 范围操作：In, Nin, Between
- 逻辑操作：Or, And, Nor, Not
- 字段操作：Exists, Type, Size
- 数组操作：All, ElemMatch
- 文本操作：Regex

**设计亮点**：
- 链式调用，代码优雅
- 类型安全的查询构建
- 自动处理复杂查询条件
- 支持投影、排序、分页

### 5. AggregationBuilder - 聚合管道

强大的聚合操作构建器：

```go
gmongo.NewAggregation(collection).
    Match(filter).
    Group("$city", bson.M{"count": gmongo.Sum(1)}).
    Sort(bson.D{{Key: "count", Value: -1}}).
    Limit(10).
    Execute(&results)
```

**支持的聚合阶段**：
- Match, Project, Group, Sort, Limit, Skip
- Unwind, Lookup, AddFields, ReplaceRoot
- Sample, Count, Bucket, BucketAuto
- Facet, GraphLookup, GeoNear
- Out, Merge, Redact

**设计亮点**：
- 链式构建聚合管道
- 提供便捷的聚合函数（Sum, Avg, Min, Max）
- 支持复杂的关联查询（Lookup）
- 完善的选项配置

### 6. IndexManager - 索引管理

完善的索引创建和管理功能：

```go
indexes := collection.Indexes()

// 创建唯一索引
indexes.CreateUniqueIndex(ctx, bson.D{{Key: "email", Value: 1}})

// 使用构建器
gmongo.NewIndexBuilder().
    AddAscending("name").
    AddDescending("age").
    Unique(true).
    Create(ctx, indexes)
```

**支持的索引类型**：
- 单字段索引
- 复合索引
- 唯一索引
- 文本索引
- TTL 索引
- 地理空间索引（2d, 2dsphere）

**设计亮点**：
- 提供索引构建器
- 支持索引存在性检查
- 便捷的索引管理方法
- 完善的索引选项配置

### 7. TransactionManager - 事务支持

简化的事务操作，支持多种隔离级别：

```go
// 简单事务
gmongo.Tx(ctx, func(sessCtx mongo.SessionContext) error {
    // 事务操作
    return nil
})

// 带选项的事务
opts := gmongo.NewTxOptions().
    ReadConcernMajority().
    WriteConcernMajority().
    Build()
gmongo.TxWithOptions(ctx, opts, fn)

// 使用隔离级别
gmongo.WithIsolationLevel(ctx, gmongo.IsolationLevelSnapshot, fn)
```

**支持的隔离级别**：
- Snapshot（快照隔离）
- Majority（多数派）
- Local（本地）

**设计亮点**：
- 自动管理会话生命周期
- 自动提交/回滚
- 支持嵌套事务
- 提供事务选项构建器

### 8. HookRegistry - 钩子系统

灵活的回调机制，支持操作前后钩子：

```go
// 全局钩子
gmongo.RegisterBeforeInsert(func(ctx context.Context, operation string, args ...interface{}) error {
    // 插入前验证
    return nil
})

// 集合级别钩子
usersWithHooks := gmongo.NewCollectionWithHooks(users, hooksRegistry)
usersWithHooks.RegisterBefore(gmongo.OpInsert, validationHook)
```

**支持的钩子类型**：
- BeforeInsert / AfterInsert
- BeforeUpdate / AfterUpdate
- BeforeDelete / AfterDelete
- BeforeFind / AfterFind
- BeforeAggregate / AfterAggregate

**设计亮点**：
- 支持全局和集合级别钩子
- 线程安全的钩子注册
- 提供便捷的钩子函数（Logging, Timing, Validation）
- 灵活的钩子管理

### 9. Utils - 工具函数

丰富的工具函数，简化常用操作：

- BSON 操作：ToDoc, FromDoc, PrintDoc
- 查询构建：BuildFilter, BuildUpdate, BuildSort
- 条件构建：Eq, Gt, In, Between, Regex
- 更新操作：Set, Inc, Push, Pull
- 分页辅助：PageInfo, PageResult, PaginateQuery
- 时间范围：Today, Yesterday, ThisWeek, ThisMonth
- ObjectID 操作：NewObjectID, ObjectIDFromHex

**设计亮点**：
- 类型别名简化使用（M, D, A, E）
- 便捷的 BSON 构建函数
- 完善的分页支持
- 实用的时间范围查询

## 与 GRDS 的对比

| 特性 | GRDS (MySQL) | GMongo (MongoDB) |
|------|--------------|------------------|
| 全局/独立模式 | ✅ | ✅ |
| 链式查询 | ✅ QueryBuilder | ✅ QueryBuilder |
| 便捷方法 | WhereEq, WhereGt | WhereEq, WhereGt |
| 事务支持 | ✅ SQL 事务 | ✅ MongoDB 事务 |
| 钩子系统 | ✅ GORM Callbacks | ✅ HookRegistry |
| 连接池 | ✅ database/sql | ✅ mongo-driver |
| 聚合操作 | SQL 查询 | ✅ Aggregation Pipeline |
| 索引管理 | ✅ Migrator | ✅ IndexManager |
| 变更监听 | ❌ | ✅ Change Streams |
| 模型生成 | ✅ Generator | 计划中 |

## 性能优化

1. **连接池管理**
   - 合理配置 MaxPoolSize 和 MinPoolSize
   - 设置连接空闲时间和最大生命周期

2. **查询优化**
   - 使用索引加速查询
   - 使用投影减少数据传输
   - 使用 Hint 指定索引

3. **批量操作**
   - 使用 InsertMany 批量插入
   - 使用 BulkWrite 批量写入
   - 使用聚合管道处理大数据

4. **并发控制**
   - 使用 context 控制超时
   - 合理使用事务
   - 避免长时间持有连接

## 最佳实践

1. **配置管理**
   - 使用环境变量管理敏感信息
   - 根据环境调整连接池大小
   - 启用压缩减少网络传输

2. **错误处理**
   - 检查 ErrNoDocuments
   - 处理重复键错误
   - 使用 context 超时

3. **索引策略**
   - 为常用查询创建索引
   - 定期分析索引使用情况
   - 删除未使用的索引

4. **事务使用**
   - 只在必要时使用事务
   - 避免长事务
   - 合理选择隔离级别

5. **监控和日志**
   - 使用钩子系统记录操作
   - 监控连接池状态
   - 记录慢查询

## 未来规划

1. **模型生成器** - 从 MongoDB 集合生成 Go 结构体
2. **迁移工具** - 数据库迁移和版本管理
3. **性能分析** - 查询性能分析和优化建议
4. **缓存集成** - 集成 Redis 等缓存系统
5. **分片支持** - 更好的分片集群支持
6. **监控面板** - 可视化监控和管理

## 总结

GMongo 是一个功能全面、性能优异、易于使用的 MongoDB 工具库。它借鉴了 GRDS 的优秀设计理念，结合 MongoDB 的特性，为 Go 开发者提供了一个现代化的数据库操作解决方案。

通过链式查询、聚合管道、索引管理、事务支持、钩子系统等核心功能，GMongo 大大简化了 MongoDB 的使用，提高了开发效率，同时保持了高性能和灵活性。
