# GMongo

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://golang.org/)
[![MongoDB Driver](https://img.shields.io/badge/MongoDB%20Driver-v1.13.1-green)](https://github.com/mongodb/mongo-go-driver)
[![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-v1.24.0-blueviolet)](https://opentelemetry.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

GMongo 是一个基于 MongoDB 官方 Go Driver 的 MongoDB 工具库，提供连接配置、CRUD、链式查询、聚合管道、事务、索引、钩子、Change Streams 以及 OpenTelemetry trace 接入能力。

## 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [客户端模式](#客户端模式)
- [配置](#配置)
- [CRUD 操作](#crud-操作)
- [链式查询](#链式查询)
- [聚合管道](#聚合管道)
- [事务](#事务)
- [索引管理](#索引管理)
- [钩子系统](#钩子系统)
- [OpenTelemetry Trace](#opentelemetry-trace)
- [低层 Monitor 接入](#低层-monitor-接入)
- [批量操作](#批量操作)
- [Change Streams](#change-streams)
- [错误处理](#错误处理)
- [性能建议](#性能建议)
- [最佳实践](#最佳实践)

## 特性

- **官方 Driver 封装**：基于 `go.mongodb.org/mongo-driver/mongo`。
- **全局/独立客户端**：支持全局默认客户端，也支持多实例独立客户端。
- **完整连接配置**：支持 URI、认证、连接池、超时、读写关注、压缩、TLS、重试等配置。
- **链式查询**：提供 `QueryBuilder`，支持条件、排序、分页、投影、统计、更新、删除。
- **聚合管道**：提供 `AggregationBuilder`，支持常用 MongoDB aggregation stage。
- **事务支持**：支持 session、transaction、事务选项构建器和隔离级别辅助方法。
- **索引管理**：支持普通索引、唯一索引、复合索引、文本索引、TTL、地理空间索引。
- **钩子系统**：支持集合操作前后扩展逻辑。
- **Change Streams**：支持 client、database、collection 级变更监听。
- **OpenTelemetry Trace**：内置 MongoDB command span 创建、错误记录和 span 结束逻辑。
- **Driver Monitor 扩展**：支持 CommandMonitor、PoolMonitor、ServerMonitor 以及 `ClientOptionsHook`。

## 安装

```bash
go get github.com/nicexiaonie/gmongo
```

当前核心依赖：

| 依赖 | 版本 |
|---|---|
| `go.mongodb.org/mongo-driver` | `v1.13.1` |
| `go.opentelemetry.io/otel` | `v1.24.0` |
| `go.opentelemetry.io/otel/trace` | `v1.24.0` |
| `go.opentelemetry.io/otel/sdk` | `v1.24.0` |

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/nicexiaonie/gmongo"
    "go.mongodb.org/mongo-driver/bson"
)

type User struct {
    ID    string `bson:"_id,omitempty"`
    Name  string `bson:"name"`
    Email string `bson:"email"`
    Age   int    `bson:"age"`
}

func main() {
    config := gmongo.DefaultConfig()
    config.Host = "localhost"
    config.Port = 27017
    config.Database = "myapp"

    if err := gmongo.Connect(config); err != nil {
        log.Fatal(err)
    }
    defer gmongo.Close()

    ctx := context.Background()
    users := gmongo.Coll("users")

    result, err := users.InsertOne(ctx, User{
        Name:  "张三",
        Email: "zhangsan@example.com",
        Age:   25,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("inserted id: %v\n", result.InsertedID)

    var found User
    if err := users.QueryWithContext(ctx).
        WhereEq("name", "张三").
        FindOne(&found); err != nil {
        log.Fatal(err)
    }

    _, err = users.UpdateOne(ctx,
        bson.M{"name": "张三"},
        bson.M{"$set": bson.M{"age": 26}},
    )
    if err != nil {
        log.Fatal(err)
    }

    _, err = users.DeleteOne(ctx, bson.M{"name": "张三"})
    if err != nil {
        log.Fatal(err)
    }
}
```

## 客户端模式

### 全局默认客户端

全局模式适合应用只连接一个 MongoDB 实例的场景。

```go
config := gmongo.DefaultConfig()
config.Database = "myapp"

if err := gmongo.Connect(config); err != nil {
    log.Fatal(err)
}
defer gmongo.Close()

users := gmongo.Coll("users")
db := gmongo.DB()
```

常用全局方法：

| 方法 | 说明 |
|---|---|
| `Connect(config)` | 连接 MongoDB 并设置默认客户端 |
| `MustConnect(config)` | 连接失败时 panic |
| `Close()` | 关闭默认客户端 |
| `GetDefaultClient()` | 获取默认 `*Client` |
| `MongoClient()` | 获取原生 `*mongo.Client` |
| `DB(name ...string)` | 获取数据库封装 |
| `Coll(name string)` | 获取集合封装 |
| `Ping(ctx)` | 测试连接 |
| `HealthCheck()` | 健康检查 |

### 独立客户端

独立客户端适合多 MongoDB 实例、多租户或测试隔离场景。

```go
config := gmongo.DefaultConfig()
config.Database = "myapp"

client, err := gmongo.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

users := client.Collection("users")
db := client.Database()
```

### 获取原生 Driver 对象

```go
mongoClient := gmongo.MongoClient()
mongoDB := gmongo.DB().DB()
mongoCollection := gmongo.Coll("users").Collection()
```

## 配置

### 默认配置

```go
config := gmongo.DefaultConfig()
config.Database = "myapp"
```

默认配置：

| 配置项 | 默认值 | 说明 |
|---|---:|---|
| `Host` | `localhost` | MongoDB 主机 |
| `Port` | `27017` | MongoDB 端口 |
| `Database` | `test` | 默认数据库 |
| `AuthSource` | `admin` | 认证数据库 |
| `MaxPoolSize` | `100` | 最大连接池大小 |
| `MinPoolSize` | `10` | 最小连接池大小 |
| `MaxConnIdleTime` | `10m` | 连接最大空闲时间 |
| `MaxConnecting` | `10` | 最大并发建连数量 |
| `ConnectTimeout` | `10s` | 连接超时 |
| `SocketTimeout` | `30s` | Socket 超时 |
| `ServerSelectTimeout` | `30s` | 服务器选择超时 |
| `HeartbeatInterval` | `10s` | 心跳间隔 |
| `ReadPreference` | `primary` | 读偏好 |
| `ReadConcern` | `local` | 读关注 |
| `WriteConcern` | `majority` | 写关注 |
| `WTimeout` | `5000ms` | 写关注超时 |
| `RetryWrites` | `true` | 重试写 |
| `RetryReads` | `true` | 重试读 |
| `Journal` | `true` | 写入 journal |
| `Compressors` | `snappy,zlib,zstd` | 压缩算法 |
| `ZlibLevel` | `6` | zlib 压缩等级 |
| `ZstdLevel` | `6` | zstd 压缩等级 |

### 完整配置示例

```go
config := &gmongo.Config{
    URI:      "",
    Host:     "localhost",
    Port:     27017,
    Username: "admin",
    Password: "password",
    Database: "myapp",

    AuthSource:       "admin",
    AuthMechanism:    "SCRAM-SHA-256",
    ReplicaSet:       "rs0",
    DirectConnection: false,
    LoadBalanced:     false,

    MaxPoolSize:     100,
    MinPoolSize:     10,
    MaxConnIdleTime: 10 * time.Minute,
    MaxConnecting:   10,

    ConnectTimeout:      10 * time.Second,
    SocketTimeout:       30 * time.Second,
    ServerSelectTimeout: 30 * time.Second,
    HeartbeatInterval:   10 * time.Second,

    ReadPreference: "primary",
    ReadConcern:    "majority",
    WriteConcern:   "majority",
    WTimeout:       5000,

    Compressors: []string{"snappy", "zlib", "zstd"},
    ZlibLevel:   6,
    ZstdLevel:   6,

    AppName: "myapp",

    TLS:         true,
    TLSInsecure: false,
}
```

### URI 模式

如果设置 `URI`，连接地址优先使用 URI。

```go
config := &gmongo.Config{
    URI:      "mongodb://user:pass@localhost:27017/myapp?authSource=admin",
    Database: "myapp",
}
```

## CRUD 操作

```go
ctx := context.Background()
users := gmongo.Coll("users")

insertResult, err := users.InsertOne(ctx, User{Name: "张三", Age: 25})

insertManyResult, err := users.InsertMany(ctx, []interface{}{
    User{Name: "李四", Age: 30},
    User{Name: "王五", Age: 35},
})

var user User
err = users.FindOne(ctx, bson.M{"name": "张三"}).Decode(&user)

var list []User
err = users.FindAll(ctx, bson.M{"age": bson.M{"$gte": 18}}, &list)

updateResult, err := users.UpdateOne(ctx,
    bson.M{"name": "张三"},
    bson.M{"$set": bson.M{"age": 26}},
)

updateManyResult, err := users.UpdateMany(ctx,
    bson.M{"status": "inactive"},
    bson.M{"$set": bson.M{"status": "active"}},
)

deleteResult, err := users.DeleteOne(ctx, bson.M{"name": "张三"})

deleteManyResult, err := users.DeleteMany(ctx, bson.M{"status": "deleted"})
```

### 便捷方法

```go
ctx := context.Background()
users := gmongo.Coll("users")

var user User
err := users.FindOneByID(ctx, "user_id", &user)

exists, err := users.Exists(ctx, bson.M{"email": "test@example.com"})

count, err := users.Count(ctx, bson.M{"age": bson.M{"$gt": 18}})

var results []User
total, err := users.Paginate(ctx, bson.M{}, 1, 10, &results)

result, err := users.Upsert(ctx,
    bson.M{"email": "lisi@example.com"},
    bson.M{"$set": bson.M{"name": "李四", "age": 30}},
)
```

## 链式查询

### 基础查询

```go
ctx := context.Background()
var users []User

err := gmongo.Coll("users").QueryWithContext(ctx).
    WhereGt("age", 18).
    WhereLt("age", 60).
    WhereIn("status", []string{"active"}).
    Select("name", "email", "age").
    Sort("-age", "name").
    Skip(10).
    Limit(20).
    Find(&users)
```

### Context 入口

推荐使用 `QueryWithContext(ctx)`，确保查询能继承请求上下文、超时、取消信号和 trace parent span。

```go
query := gmongo.Coll("users").QueryWithContext(ctx)
```

也可以使用已有 builder 后再设置 context：

```go
query := gmongo.Coll("users").Query().Context(ctx)
```

如果直接使用 `Query()` 且不调用 `Context(ctx)`，builder 默认使用 `context.Background()`，不会自动关联上游 trace。

### 条件方法

| 方法 | 说明 |
|---|---|
| `Where` / `WhereEq` | 等于 |
| `WhereNe` | 不等于 |
| `WhereGt` / `WhereGte` | 大于 / 大于等于 |
| `WhereLt` / `WhereLte` | 小于 / 小于等于 |
| `WhereIn` / `WhereNin` | 在集合中 / 不在集合中 |
| `WhereBetween` | 范围查询 |
| `WhereExists` | 字段是否存在 |
| `WhereRegex` | 正则匹配 |
| `WhereType` | BSON 类型匹配 |
| `WhereSize` | 数组大小 |
| `WhereAll` | 数组包含所有元素 |
| `WhereElemMatch` | 数组元素匹配 |
| `Or` / `And` / `Nor` / `Not` | 逻辑条件 |

### 查询选项

| 方法 | 说明 |
|---|---|
| `Select(fields...)` | 只返回指定字段 |
| `Omit(fields...)` | 排除指定字段 |
| `Sort(fields...)` | 排序，`-field` 表示降序 |
| `Skip(n)` | 跳过数量 |
| `Limit(n)` | 限制数量 |
| `Page(page, pageSize)` | 分页 |
| `Hint(hint)` | 指定索引 |
| `MaxTime(duration)` | 最大执行时间 |
| `AllowDiskUse(bool)` | 是否允许使用磁盘 |
| `Collation(collation)` | 排序规则 |

### 执行方法

```go
query := gmongo.Coll("users").QueryWithContext(ctx).
    WhereEq("status", "active")

var users []User
err := query.Find(&users)

var user User
err = query.FindOne(&user)

count, err := query.Count()
exists, err := query.Exists()

updateResult, err := query.Update(bson.M{"$set": bson.M{"status": "disabled"}})
updateManyResult, err := query.UpdateMany(bson.M{"$set": bson.M{"status": "disabled"}})

deleteResult, err := query.Delete()
deleteManyResult, err := query.DeleteMany()

values, err := query.Distinct("city")
cursor, err := query.Cursor()
```

## 聚合管道

### 构建聚合

```go
ctx := context.Background()
users := gmongo.Coll("users")

type AgeGroup struct {
    ID     string  `bson:"_id"`
    Count  int     `bson:"count"`
    AvgAge float64 `bson:"avgAge"`
}

var results []AgeGroup
err := users.NewAggregationWithContext(ctx).
    Match(bson.M{"status": "active"}).
    Group("$city", bson.M{
        "count":  gmongo.Sum(1),
        "avgAge": gmongo.Avg("$age"),
    }).
    Sort(bson.D{{Key: "count", Value: -1}}).
    Limit(10).
    Execute(&results)
```

### 创建方式

推荐使用集合方法：

```go
aggregation := users.NewAggregationWithContext(ctx)
```

也可以使用包级函数：

```go
aggregation := gmongo.NewAggregationWithContext(ctx, users)
```

如果使用不带 context 的创建方式，需要显式设置 context：

```go
aggregation := gmongo.NewAggregation(users).Context(ctx)
```

### 常用 stage

| 方法 | MongoDB Stage |
|---|---|
| `Match` | `$match` |
| `Project` | `$project` |
| `Group` | `$group` |
| `Sort` | `$sort` |
| `Limit` | `$limit` |
| `Skip` | `$skip` |
| `Unwind` | `$unwind` |
| `Lookup` / `LookupPipeline` | `$lookup` |
| `AddFields` | `$addFields` |
| `ReplaceRoot` | `$replaceRoot` |
| `Sample` | `$sample` |
| `Count` | `$count` |
| `Bucket` / `BucketAuto` | `$bucket` / `$bucketAuto` |
| `Facet` | `$facet` |
| `GraphLookup` | `$graphLookup` |
| `Out` | `$out` |
| `Merge` | `$merge` |
| `Redact` | `$redact` |
| `GeoNear` | `$geoNear` |

### 聚合选项

```go
err := users.NewAggregationWithContext(ctx).
    Match(bson.M{"status": "active"}).
    AllowDiskUse(true).
    BatchSize(500).
    MaxTime(5 * time.Second).
    Comment("active-users-report").
    Execute(&results)
```

## 事务

### 简单事务

```go
ctx := context.Background()

err := gmongo.Tx(ctx, func(sessCtx mongo.SessionContext) error {
    users := gmongo.Coll("users")
    accounts := gmongo.Coll("accounts")

    if _, err := users.InsertOne(sessCtx, user); err != nil {
        return err
    }
    if _, err := accounts.InsertOne(sessCtx, account); err != nil {
        return err
    }
    return nil
})
```

### 带事务选项

```go
opts := gmongo.NewTxOptions().
    ReadConcernMajority().
    WriteConcernMajority().
    Build()

err := gmongo.TxWithOptions(ctx, opts, func(sessCtx mongo.SessionContext) error {
    return nil
})
```

### 独立事务管理器

```go
tm := gmongo.NewTransactionManager(client)

err := tm.Execute(ctx, func(sessCtx mongo.SessionContext) error {
    return nil
})

result, err := tm.ExecuteWithResult(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
    return "ok", nil
})
```

### 隔离级别辅助方法

```go
err := gmongo.WithIsolationLevel(ctx, gmongo.IsolationLevelSnapshot, func(sessCtx mongo.SessionContext) error {
    return nil
})

err = gmongo.WithReadConcernMajority(ctx, func(sessCtx mongo.SessionContext) error {
    return nil
})

err = gmongo.WithWriteConcernMajority(ctx, func(sessCtx mongo.SessionContext) error {
    return nil
})
```

事务中的 MongoDB command 会继续使用 `sessCtx` 传递到 MongoDB Driver，因此也会被 trace instrumentation 记录为 command span。

## 索引管理

```go
ctx := context.Background()
users := gmongo.Coll("users")
indexes := users.Indexes()

indexName, err := indexes.CreateUniqueIndex(ctx, bson.D{{Key: "email", Value: 1}})

indexName, err = indexes.CreateCompoundIndex(ctx, map[string]int{
    "name": 1,
    "age":  -1,
}, false)

indexName, err = indexes.CreateTextIndex(ctx, "description")
indexName, err = indexes.CreateTTLIndex(ctx, "created_at", 24*time.Hour)
indexName, err = indexes.CreateGeoIndex(ctx, "location", "2dsphere")

indexName, err = gmongo.NewIndexBuilder().
    AddAscending("name").
    AddDescending("age").
    Unique(true).
    Name("idx_name_age").
    Sparse(true).
    Create(ctx, indexes)

specs, err := indexes.ListSpecifications(ctx)
err = indexes.DropByName(ctx, "idx_name_age")
```

## 钩子系统

钩子系统用于业务层操作前后扩展，不等同于 MongoDB Driver 级 trace。Driver 级 trace 请使用 [OpenTelemetry Trace](#opentelemetry-trace)。

```go
gmongo.RegisterBeforeInsert(func(ctx context.Context, operation string, args ...interface{}) error {
    return nil
})

gmongo.RegisterAfterInsert(func(ctx context.Context, operation string, args ...interface{}) error {
    return nil
})

users := gmongo.Coll("users")
registry := gmongo.NewHookRegistry()
usersWithHooks := gmongo.NewCollectionWithHooks(users, registry)

usersWithHooks.RegisterBefore(gmongo.OpInsert, func(ctx context.Context, operation string, args ...interface{}) error {
    return nil
})

usersWithHooks.RegisterAfter(gmongo.OpUpdate, func(ctx context.Context, operation string, args ...interface{}) error {
    return nil
})
```

支持的操作常量：

| 常量 | 值 |
|---|---|
| `OpInsert` | `insert` |
| `OpUpdate` | `update` |
| `OpDelete` | `delete` |
| `OpFind` | `find` |
| `OpAggregate` | `aggregate` |
| `OpCount` | `count` |
| `OpDistinct` | `distinct` |

## OpenTelemetry Trace

GMongo 内置 MongoDB command tracing instrumentation。使用者只需要初始化 OpenTelemetry `TracerProvider`，再将 `mongotrace.New(...)` 配置到 `gmongo.Config.Tracing`。

### 工作机制

1. GMongo 在创建 MongoDB client 时将 tracing instrumentation 注入 `mongo/options.ClientOptions`。
2. MongoDB Driver 发出 command started 事件时，`tracing/trace` 基于当前操作的 `context.Context` 创建 client span。
3. MongoDB Driver 发出 command succeeded 事件时，instrumentation 标记 span 成功并结束 span。
4. MongoDB Driver 发出 command failed 事件时，instrumentation 记录错误、标记 span 失败并结束 span。
5. started / succeeded / failed 事件通过 `RequestID + ConnectionID` 关联。

### 快速接入

```go
package main

import (
    "context"
    "log"

    "github.com/nicexiaonie/gmongo"
    mongotrace "github.com/nicexiaonie/gmongo/tracing/trace"
    "go.opentelemetry.io/otel"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
    tracerProvider := sdktrace.NewTracerProvider()
    defer tracerProvider.Shutdown(context.Background())
    otel.SetTracerProvider(tracerProvider)

    config := gmongo.DefaultConfig()
    config.Database = "myapp"
    config.Tracing = mongotrace.New(mongotrace.Config{
        Enabled:        true,
        TracerProvider: tracerProvider,
        TracerName:     "myapp/mongodb",
    })

    if err := gmongo.Connect(config); err != nil {
        log.Fatal(err)
    }
    defer gmongo.Close()
}
```

如果项目已经在应用入口统一设置过全局 OpenTelemetry provider，也可以省略 `TracerProvider`，默认使用 `otel.GetTracerProvider()`。

```go
config.Tracing = mongotrace.New(mongotrace.DefaultConfig())
```

### 使用业务 context 关联 command span

```go
tracer := otel.Tracer("myapp/service")
ctx, span := tracer.Start(context.Background(), "CreateUser")
defer span.End()

users := gmongo.Coll("users")

_, err := users.InsertOne(ctx, User{Name: "张三"})

var list []User
err = users.QueryWithContext(ctx).
    WhereEq("status", "active").
    Find(&list)

err = users.NewAggregationWithContext(ctx).
    Match(bson.M{"status": "active"}).
    Execute(&list)
```

只要传入带 parent span 的 `ctx`，MongoDB command span 就会作为子 span 记录。若传入 nil 或使用默认 `context.Background()`，MongoDB command span 仍会产生，但不会挂到上游请求链路。

### Trace 配置项

| 配置项 | 类型 | 说明 |
|---|---|---|
| `Enabled` | `bool` | 是否启用 command span 创建 |
| `TracerProvider` | `trace.TracerProvider` | OpenTelemetry tracer provider |
| `TracerName` | `string` | tracer 名称 |
| `SpanNameFormatter` | `func(CommandStartedInfo) string` | 自定义 span name |
| `Attributes` | `func(CommandStartedInfo) []attribute.KeyValue` | 追加 span attributes |
| `RecordCommand` | `bool` | 是否启用 command 白名单属性记录 |
| `CommandSanitizer` | `func(CommandStartedInfo) []attribute.KeyValue` | command 脱敏和白名单属性函数 |
| `PoolEventHandler` | `func(PoolEventInfo)` | 连接池诊断事件处理 |
| `ServerEventHandler` | `func(ServerEventInfo)` | server/topology/heartbeat 诊断事件处理 |

### Span 命名

默认 span name：

```text
MongoDB <command> <collection>
```

示例：

```text
MongoDB find users
MongoDB insert users
MongoDB aggregate orders
```

自定义 span name：

```go
config.Tracing = mongotrace.New(mongotrace.Config{
    Enabled:        true,
    TracerProvider: tracerProvider,
    SpanNameFormatter: func(info mongotrace.CommandStartedInfo) string {
        return "mongo." + info.CommandName
    },
})
```

### 默认 Span Attributes

默认记录低风险属性：

| Attribute | 说明 |
|---|---|
| `db.system` | 固定为 `mongodb` |
| `db.name` | 数据库名 |
| `db.operation` | MongoDB command name |
| `db.mongodb.collection` | 集合名，可提取时记录 |
| `db.mongodb.request_id` | MongoDB driver request id |
| `db.mongodb.connection_id` | MongoDB driver connection id |
| `db.mongodb.duration_ms` | command 完成耗时 |

追加自定义属性：

```go
config.Tracing = mongotrace.New(mongotrace.Config{
    Enabled:        true,
    TracerProvider: tracerProvider,
    Attributes: func(info mongotrace.CommandStartedInfo) []attribute.KeyValue {
        return []attribute.KeyValue{
            attribute.String("app.mongo.collection", info.Collection),
        }
    },
})
```

### 敏感数据与脱敏

默认不会把完整 Mongo command、filter、document、update、reply 写入 span。查询条件、插入文档、更新内容可能包含手机号、邮箱、token、身份证号等敏感数据，不建议直接进入 trace 系统。

如确实需要记录部分 command 信息，应只记录白名单字段：

```go
config.Tracing = mongotrace.New(mongotrace.Config{
    Enabled:        true,
    TracerProvider: tracerProvider,
    RecordCommand:  true,
    CommandSanitizer: func(info mongotrace.CommandStartedInfo) []attribute.KeyValue {
        return []attribute.KeyValue{
            attribute.String("mongo.command", info.CommandName),
            attribute.String("mongo.collection", info.Collection),
        }
    },
})
```

### Pool / Server Diagnostics

MongoDB Driver 的 PoolMonitor 和 ServerMonitor 不携带业务 `context.Context`，因此默认不创建请求子 span。它们适合用于 metrics、日志和诊断。

```go
config.Tracing = mongotrace.New(mongotrace.Config{
    Enabled:        true,
    TracerProvider: tracerProvider,
    PoolEventHandler: func(info mongotrace.PoolEventInfo) {
        log.Printf("mongo pool event: type=%s address=%s reason=%s", info.Type, info.Address, info.Reason)
    },
    ServerEventHandler: func(info mongotrace.ServerEventInfo) {
        log.Printf("mongo server event: type=%s connectionID=%s", info.Type, info.ConnectionID)
    },
})
```

## 低层 Monitor 接入

如果不使用 OpenTelemetry，也可以直接配置 MongoDB Driver monitor。

```go
config.CommandMonitor = customCommandMonitor
config.PoolMonitor = customPoolMonitor
config.ServerMonitor = customServerMonitor
```

如果同时配置 `config.Tracing` 和手动 monitor，GMongo 会组合执行：

1. tracing monitor 先执行。
2. 手动 monitor 后执行。
3. `ClientOptionsHook` 最后执行，可覆盖前面所有配置。

```go
config.ClientOptionsHook = func(opts *options.ClientOptions) {
    opts.SetAppName("myapp")
}
```

保留低层 callback adapter：

```go
config.CommandMonitor = mongotrace.NewCommandMonitor(mongotrace.CommandMonitorOptions{
    Started: func(ctx context.Context, e *event.CommandStartedEvent) {},
    Succeeded: func(ctx context.Context, e *event.CommandSucceededEvent) {},
    Failed: func(ctx context.Context, e *event.CommandFailedEvent) {},
})
```

## 批量操作

```go
ctx := context.Background()
users := gmongo.Coll("users")

documents := []interface{}{
    User{Name: "用户1", Age: 20},
    User{Name: "用户2", Age: 25},
    User{Name: "用户3", Age: 30},
}
result, err := users.InsertMany(ctx, documents)

models := []mongo.WriteModel{
    mongo.NewInsertOneModel().SetDocument(User{Name: "新用户", Age: 22}),
    mongo.NewUpdateOneModel().
        SetFilter(bson.M{"name": "用户1"}).
        SetUpdate(bson.M{"$set": bson.M{"age": 21}}),
    mongo.NewDeleteOneModel().SetFilter(bson.M{"name": "用户3"}),
}
bulkResult, err := users.BulkWrite(ctx, models)
```

## Change Streams

```go
ctx := context.Background()
users := gmongo.Coll("users")

pipeline := []bson.M{
    {"$match": bson.M{"operationType": bson.M{"$in": []string{"insert", "update", "delete"}}}},
}

stream, err := users.Watch(ctx, pipeline)
if err != nil {
    log.Fatal(err)
}
defer stream.Close(ctx)

for stream.Next(ctx) {
    var changeEvent bson.M
    if err := stream.Decode(&changeEvent); err != nil {
        log.Fatal(err)
    }
}

if err := stream.Err(); err != nil {
    log.Fatal(err)
}
```

## 错误处理

```go
import (
    "errors"

    "go.mongodb.org/mongo-driver/mongo"
)

err := users.FindOne(ctx, filter).Decode(&user)
if errors.Is(err, mongo.ErrNoDocuments) {
    // 文档不存在
}

_, err = users.InsertOne(ctx, user)
if mongo.IsDuplicateKeyError(err) {
    // 唯一索引冲突
}
```

## 性能建议

1. 为高频查询字段创建索引。
2. 使用投影减少网络传输和反序列化成本。
3. 大批量写入优先使用 `InsertMany` 或 `BulkWrite`。
4. 根据服务并发和 MongoDB 容量合理配置连接池。
5. 聚合中尽量让 `$match`、`$project` 前置，减少后续 stage 数据量。
6. 大数据量聚合按需启用 `AllowDiskUse`。
7. 复杂查询可使用 `Hint` 明确索引。
8. 始终传递业务 `context.Context`，支持超时、取消和 trace 关联。
9. trace 中不要记录完整 command/filter/document，避免增加性能和合规风险。

```go
err := gmongo.Coll("users").QueryWithContext(ctx).
    WhereGt("age", 18).
    Hint(bson.D{{Key: "age", Value: 1}}).
    Find(&users)

err = gmongo.Coll("orders").NewAggregationWithContext(ctx).
    Match(bson.M{"status": "completed"}).
    AllowDiskUse(true).
    Group("$city", bson.M{"count": gmongo.Sum(1)}).
    Execute(&results)
```

## 最佳实践

1. **显式传递 context**：直接操作使用 `ctx` 参数，链式查询使用 `QueryWithContext(ctx)`，聚合使用 `NewAggregationWithContext(ctx)`。
2. **统一初始化 OpenTelemetry**：在应用入口初始化 `TracerProvider`，再配置 `gmongo.Config.Tracing`。
3. **不要泄露敏感数据**：trace、日志、metrics 中不要记录完整 command/filter/document。
4. **区分 trace 与 diagnostics**：command 事件用于 request-scoped trace；pool/server 事件用于 metrics 和诊断。
5. **按需配置连接池**：根据应用并发、MongoDB 规格和超时策略设置 pool 参数。
6. **事务保持短小**：事务内只放必须原子提交的操作，避免长事务。
7. **索引和查询同步设计**：新增高频查询前先评估索引。
8. **区分 hooks 与 monitor**：hooks 适合业务扩展；monitor 适合 driver 级 trace、metrics、diagnostics。
9. **关闭资源**：及时关闭 client、cursor、change stream。
10. **处理 cursor 错误**：遍历 cursor 或 change stream 后检查 `Err()`。

## 许可证

MIT License

## 相关链接

- [MongoDB 官方文档](https://www.mongodb.com/docs/)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)
