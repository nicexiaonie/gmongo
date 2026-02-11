# GMongo - MongoDB 工具库

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

GMongo 是一个基于 `go.mongodb.org/mongo-driver/mongo` 的现代化 MongoDB 工具库，提供全面、稳定、高性能的 MongoDB 操作封装。

## 特性

- ✨ **开箱即用** - 简洁的 API 设计，快速上手
- 🔗 **链式调用** - 流畅的查询构建器，优雅的代码风格
- 🌍 **全局/独立** - 支持全局默认客户端和多数据库实例
- 🎯 **便捷方法** - 提供大量便捷方法，简化常用操作
- 💾 **事务支持** - 完善的事务操作，支持多种隔离级别
- 🪝 **钩子系统** - 灵活的回调机制，支持操作前后钩子
- 🔒 **并发安全** - 线程安全的设计，适合高并发场景
- 📊 **聚合管道** - 强大的聚合操作构建器
- 🔍 **索引管理** - 完善的索引创建和管理功能
- ⚡ **高性能** - 基于官方 MongoDB Driver，性能优异
- 🎨 **现代化设计** - 符合 Go 语言最佳实践

## 安装

```bash
go get -u github.com/nicexiaonie/gmongo
```

## 快速开始

### 基础使用

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/yourusername/go-sophon/gmongo"
    "go.mongodb.org/mongo-driver/bson"
)

type User struct {
    ID    string `bson:"_id,omitempty"`
    Name  string `bson:"name"`
    Email string `bson:"email"`
    Age   int    `bson:"age"`
}

func main() {
    // 1. 连接数据库（全局模式）
    config := gmongo.DefaultConfig()
    config.Host = "localhost"
    config.Port = 27017
    config.Database = "myapp"

    if err := gmongo.Connect(config); err != nil {
        log.Fatal(err)
    }
    defer gmongo.Close()

    ctx := context.Background()

    // 2. 获取集合
    users := gmongo.Collection("users")

    // 3. 插入文档
    user := &User{
        Name:  "张三",
        Email: "zhangsan@example.com",
        Age:   25,
    }
    result, err := users.InsertOne(ctx, user)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("插入成功，ID: %v\n", result.InsertedID)

    // 4. 查询文档
    var foundUser User
    err = users.Query().
        WhereEq("name", "张三").
        FindOne(&foundUser)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("查询结果: %+v\n", foundUser)

    // 5. 更新文档
    update := bson.M{"$set": bson.M{"age": 26}}
    _, err = users.UpdateOne(ctx, bson.M{"name": "张三"}, update)
    if err != nil {
        log.Fatal(err)
    }

    // 6. 删除文档
    _, err = users.DeleteOne(ctx, bson.M{"name": "张三"})
    if err != nil {
        log.Fatal(err)
    }
}
```

### 独立客户端模式

```go
// 创建独立客户端
config := gmongo.DefaultConfig()
config.Database = "myapp"

client, err := gmongo.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 使用独立客户端
users := client.Collection("users")
```

## 核心功能

### 1. 链式查询构建器

```go
var users []User

// 复杂查询
err := gmongo.Collection("users").Query().
    WhereGt("age", 18).                    // 年龄大于 18
    WhereLt("age", 60).                    // 年龄小于 60
    WhereIn("status", []string{"active"}). // 状态为 active
    WhereRegex("name", "^张", "i").        // 名字以"张"开头（不区分大小写）
    Select("name", "email", "age").        // 只选择指定字段
    Sort("-age", "name").                  // 按年龄降序，名字升序
    Skip(10).                              // 跳过 10 条
    Limit(20).                             // 限制 20 条
    Find(&users)

if err != nil {
    log.Fatal(err)
}
```

### 2. 便捷查询方法

```go
ctx := context.Background()
users := gmongo.Collection("users")

// 根据 ID 查询
var user User
err := users.FindOneByID(ctx, "user_id", &user)

// 检查是否存在
exists, err := users.Exists(ctx, bson.M{"email": "test@example.com"})

// 统计数量
count, err := users.Count(ctx, bson.M{"age": bson.M{"$gt": 18}})

// 分页查询
var results []User
total, err := users.Paginate(ctx, bson.M{}, 1, 10, &results)
fmt.Printf("总数: %d, 当前页: %d 条\n", total, len(results))

// Upsert（更新或插入）
update := bson.M{"$set": bson.M{"name": "李四", "age": 30}}
result, err := users.Upsert(ctx, bson.M{"email": "lisi@example.com"}, update)
```

### 3. 聚合管道

```go
// 复杂聚合查询
type AgeGroup struct {
    ID    string `bson:"_id"`
    Count int    `bson:"count"`
    Avg   float64 `bson:"avgAge"`
}

var results []AgeGroup
err := gmongo.NewAggregation(users).
    Match(bson.M{"status": "active"}).                    // 过滤活跃用户
    Group("$city", bson.M{                                // 按城市分组
        "count":  gmongo.Sum(1),                          // 统计数量
        "avgAge": gmongo.Avg("$age"),                     // 平均年龄
    }).
    Sort(bson.D{{Key: "count", Value: -1}}).              // 按数量降序
    Limit(10).                                            // 前 10 个
    Execute(&results)

if err != nil {
    log.Fatal(err)
}

// 关联查询（类似 SQL JOIN）
err = gmongo.NewAggregation(orders).
    Lookup("users", "user_id", "_id", "user").            // 关联用户表
    Unwind("$user", true).                                // 展开用户数组
    Match(bson.M{"status": "completed"}).                 // 过滤已完成订单
    Project(bson.M{                                       // 投影字段
        "order_id":   1,
        "user_name":  "$user.name",
        "user_email": "$user.email",
        "total":      1,
    }).
    Execute(&results)
```

### 4. 事务支持

```go
ctx := context.Background()

// 简单事务
err := gmongo.Tx(ctx, func(sessCtx mongo.SessionContext) error {
    users := gmongo.Collection("users")
    accounts := gmongo.Collection("accounts")

    // 在事务中执行多个操作
    _, err := users.InsertOne(sessCtx, user)
    if err != nil {
        return err // 自动回滚
    }

    _, err = accounts.InsertOne(sessCtx, account)
    if err != nil {
        return err // 自动回滚
    }

    return nil // 自动提交
})

// 带选项的事务
opts := gmongo.NewTxOptions().
    ReadConcernMajority().
    WriteConcernMajority().
    Build()

err = gmongo.TxWithOptions(ctx, opts, func(sessCtx mongo.SessionContext) error {
    // 事务操作
    return nil
})

// 使用隔离级别
err = gmongo.WithIsolationLevel(ctx, gmongo.IsolationLevelSnapshot, func(sessCtx mongo.SessionContext) error {
    // 事务操作
    return nil
})
```

### 5. 索引管理

```go
ctx := context.Background()
users := gmongo.Collection("users")
indexes := users.Indexes()

// 创建唯一索引
indexName, err := indexes.CreateUniqueIndex(ctx, bson.D{
    {Key: "email", Value: 1},
})

// 创建复合索引
indexName, err = indexes.CreateCompoundIndex(ctx, map[string]int{
    "name": 1,
    "age":  -1,
}, false)

// 创建文本索引
indexName, err = indexes.CreateTextIndex(ctx, "description")

// 创建 TTL 索引（自动过期）
indexName, err = indexes.CreateTTLIndex(ctx, "created_at", 24*time.Hour)

// 创建地理空间索引
indexName, err = indexes.CreateGeoIndex(ctx, "location", "2dsphere")

// 使用索引构建器
indexName, err = gmongo.NewIndexBuilder().
    AddAscending("name").
    AddDescending("age").
    Unique(true).
    Name("idx_name_age").
    Sparse(true).
    Create(ctx, indexes)

// 列出所有索引
specs, err := indexes.ListSpecifications(ctx)
for _, spec := range specs {
    fmt.Printf("索引: %s\n", spec.Name)
}

// 删除索引
err = indexes.DropByName(ctx, "idx_name_age")
```

### 6. 钩子系统

```go
// 注册全局钩子
gmongo.RegisterBeforeInsert(func(ctx context.Context, operation string, args ...interface{}) error {
    fmt.Printf("准备插入: %v\n", args)
    return nil
})

gmongo.RegisterAfterInsert(func(ctx context.Context, operation string, args ...interface{}) error {
    fmt.Printf("插入完成: %v\n", args)
    return nil
})

// 集合级别钩子
users := gmongo.Collection("users")
hooksRegistry := gmongo.NewHookRegistry()
usersWithHooks := gmongo.NewCollectionWithHooks(users, hooksRegistry)

usersWithHooks.RegisterBefore(gmongo.OpInsert, func(ctx context.Context, operation string, args ...interface{}) error {
    // 插入前验证
    return nil
})

usersWithHooks.RegisterAfter(gmongo.OpUpdate, func(ctx context.Context, operation string, args ...interface{}) error {
    // 更新后清除缓存
    return nil
})

// 使用便捷钩子
gmongo.RegisterBeforeInsert(gmongo.LoggingHook(func(format string, args ...interface{}) {
    log.Printf(format, args...)
}))
```

### 7. 批量操作

```go
ctx := context.Background()
users := gmongo.Collection("users")

// 批量插入
documents := []interface{}{
    User{Name: "用户1", Age: 20},
    User{Name: "用户2", Age: 25},
    User{Name: "用户3", Age: 30},
}
result, err := users.InsertMany(ctx, documents)
fmt.Printf("插入了 %d 条记录\n", len(result.InsertedIDs))

// 批量写入（混合操作）
models := []mongo.WriteModel{
    mongo.NewInsertOneModel().SetDocument(User{Name: "新用户", Age: 22}),
    mongo.NewUpdateOneModel().
        SetFilter(bson.M{"name": "用户1"}).
        SetUpdate(bson.M{"$set": bson.M{"age": 21}}),
    mongo.NewDeleteOneModel().SetFilter(bson.M{"name": "用户3"}),
}
bulkResult, err := users.BulkWrite(ctx, models)
```

### 8. 变更流（Change Streams）

```go
ctx := context.Background()
users := gmongo.Collection("users")

// 监听集合变更
pipeline := []bson.M{
    {"$match": bson.M{"operationType": bson.M{"$in": []string{"insert", "update", "delete"}}}},
}

stream, err := users.Watch(ctx, pipeline)
if err != nil {
    log.Fatal(err)
}
defer stream.Close(ctx)

// 处理变更事件
for stream.Next(ctx) {
    var changeEvent bson.M
    if err := stream.Decode(&changeEvent); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("变更类型: %v\n", changeEvent["operationType"])
    fmt.Printf("文档: %v\n", changeEvent["fullDocument"])
}
```

## 配置选项

### 完整配置示例

```go
config := &gmongo.Config{
    // 基础连接
    Host:     "localhost",
    Port:     27017,
    Username: "admin",
    Password: "password",
    Database: "myapp",

    // 认证
    AuthSource:    "admin",
    AuthMechanism: "SCRAM-SHA-256",
    ReplicaSet:    "rs0",

    // 连接池
    MaxPoolSize:     100,
    MinPoolSize:     10,
    MaxConnIdleTime: 10 * time.Minute,
    MaxConnecting:   10,

    // 超时
    ConnectTimeout:      10 * time.Second,
    SocketTimeout:       30 * time.Second,
    ServerSelectTimeout: 30 * time.Second,
    HeartbeatInterval:   10 * time.Second,

    // 读写配置
    ReadPreference: "primary",
    ReadConcern:    "majority",
    WriteConcern:   "majority",
    WTimeout:       5000,

    // 压缩
    Compressors: []string{"snappy", "zlib", "zstd"},
    ZlibLevel:   6,
    ZstdLevel:   6,

    // TLS
    TLS:         true,
    TLSInsecure: false,

    // 应用
    AppName: "MyApp",
}

// 或使用 URI
config := &gmongo.Config{
    URI:      "mongodb://user:pass@localhost:27017/myapp?authSource=admin",
    Database: "myapp",
}
```

### 使用默认配置

```go
config := gmongo.DefaultConfig()
config.Database = "myapp"
// 修改需要的配置项
```

## 高级特性

### 1. 地理空间查询

```go
// 创建地理空间索引
indexes := gmongo.Collection("places").Indexes()
_, err := indexes.CreateGeoIndex(ctx, "location", "2dsphere")

// 查询附近的地点
pipeline := []bson.M{
    {
        "$geoNear": bson.M{
            "near": bson.M{
                "type":        "Point",
                "coordinates": []float64{116.404, 39.915}, // 经度, 纬度
            },
            "distanceField": "distance",
            "maxDistance":   5000, // 5公里
            "spherical":     true,
        },
    },
}

var results []Place
err = gmongo.NewAggregation(places).
    AddStage(pipeline[0]).
    Execute(&results)
```

### 2. 全文搜索

```go
// 创建文本索引
indexes := gmongo.Collection("articles").Indexes()
_, err := indexes.CreateTextIndex(ctx, "content")

// 全文搜索
var articles []Article
err = gmongo.Collection("articles").Query().
    Filter(bson.M{"$text": bson.M{"$search": "MongoDB 教程"}}).
    Find(&articles)
```

### 3. 数组操作

```go
// 查询数组包含特定元素
err := gmongo.Collection("users").Query().
    WhereIn("tags", []string{"golang", "mongodb"}).
    Find(&users)

// 查询数组包含所有元素
err = gmongo.Collection("users").Query().
    WhereAll("tags", []string{"golang", "mongodb"}).
    Find(&users)

// 查询数组大小
err = gmongo.Collection("users").Query().
    WhereSize("tags", 3).
    Find(&users)

// 数组元素匹配
err = gmongo.Collection("users").Query().
    WhereElemMatch("scores", bson.M{"$gte": 80, "$lte": 100}).
    Find(&users)
```

### 4. 投影和排除字段

```go
// 只选择特定字段
var users []User
err := gmongo.Collection("users").Query().
    Select("name", "email").
    Find(&users)

// 排除特定字段
err = gmongo.Collection("users").Query().
    Omit("password", "secret_key").
    Find(&users)
```

## 性能优化建议

1. **使用索引**: 为常用查询字段创建索引
2. **批量操作**: 使用 `InsertMany` 和 `BulkWrite` 进行批量操作
3. **投影**: 只查询需要的字段，减少网络传输
4. **连接池**: 合理配置连接池大小
5. **读偏好**: 根据场景选择合适的读偏好
6. **聚合管道**: 在数据库端进行数据处理，减少应用层计算
7. **使用 Hint**: 为复杂查询指定使用的索引

```go
// 使用索引提示
err := gmongo.Collection("users").Query().
    WhereGt("age", 18).
    Hint(bson.D{{Key: "age", Value: 1}}).
    Find(&users)

// 允许使用磁盘（大数据量聚合）
err = gmongo.NewAggregation(users).
    AllowDiskUse(true).
    Group("$city", bson.M{"count": gmongo.Sum(1)}).
    Execute(&results)
```

## 错误处理

```go
import (
    "errors"
    "go.mongodb.org/mongo-driver/mongo"
)

// 检查文档不存在
err := users.FindOne(ctx, filter).Decode(&user)
if errors.Is(err, mongo.ErrNoDocuments) {
    fmt.Println("文档不存在")
}

// 检查重复键错误
_, err = users.InsertOne(ctx, user)
if mongo.IsDuplicateKeyError(err) {
    fmt.Println("文档已存在")
}
```

## 最佳实践

1. **使用上下文**: 始终传递 context，支持超时和取消
2. **关闭资源**: 使用 defer 关闭游标和客户端
3. **错误处理**: 妥善处理所有错误
4. **事务使用**: 只在必要时使用事务，避免长事务
5. **索引策略**: 定期分析和优化索引
6. **监控**: 使用钩子系统实现日志和监控

## 示例项目

查看 `examples/` 目录获取更多示例：

- [基础 CRUD 操作](examples/basic/main.go)
- [事务处理](examples/transaction/main.go)
- [聚合管道](examples/aggregation/main.go)
- [索引管理](examples/index/main.go)
- [变更流](examples/changestream/main.go)

## 与 GRDS 的对比

GMongo 借鉴了 GRDS 的优秀设计理念：

| 特性 | GRDS (MySQL) | GMongo (MongoDB) |
|------|--------------|------------------|
| 全局/独立模式 | ✅ | ✅ |
| 链式查询 | ✅ | ✅ |
| 事务支持 | ✅ | ✅ |
| 钩子系统 | ✅ | ✅ |
| 便捷方法 | ✅ | ✅ |
| 聚合操作 | SQL | Pipeline |
| 索引管理 | ✅ | ✅ |
| 变更监听 | ❌ | ✅ (Change Streams) |

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 相关链接

- [MongoDB 官方文档](https://docs.mongodb.com/)
- [Go MongoDB Driver](https://github.com/mongodb/mongo-go-driver)
- [GRDS 工具库](../grds/README.md)
