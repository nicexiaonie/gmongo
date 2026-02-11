# GMongo 快速开始指南

## 5 分钟快速上手

### 1. 安装

```bash
go get github.com/yourusername/go-sophon/gmongo
```

### 2. 最简单的例子

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/yourusername/go-sophon/gmongo"
    "go.mongodb.org/mongo-driver/bson"
)

func main() {
    // 连接数据库
    config := gmongo.DefaultConfig()
    config.Database = "myapp"

    if err := gmongo.Connect(config); err != nil {
        log.Fatal(err)
    }
    defer gmongo.Close()

    ctx := context.Background()
    users := gmongo.Collection("users")

    // 插入
    user := bson.M{"name": "张三", "age": 25}
    result, _ := users.InsertOne(ctx, user)
    fmt.Println("插入成功:", result.InsertedID)

    // 查询
    var found bson.M
    users.FindOne(ctx, bson.M{"name": "张三"}).Decode(&found)
    fmt.Println("查询结果:", found)
}
```

### 3. 使用链式查询

```go
type User struct {
    Name  string `bson:"name"`
    Age   int    `bson:"age"`
    Email string `bson:"email"`
}

var users []User

// 链式查询
err := gmongo.Collection("users").Query().
    WhereGt("age", 18).              // 年龄 > 18
    WhereLt("age", 60).              // 年龄 < 60
    Sort("-age").                    // 按年龄降序
    Limit(10).                       // 限制 10 条
    Find(&users)                     // 执行查询

if err != nil {
    log.Fatal(err)
}

for _, user := range users {
    fmt.Printf("%s: %d 岁\n", user.Name, user.Age)
}
```

### 4. 聚合查询

```go
type CityStats struct {
    City  string  `bson:"_id"`
    Count int     `bson:"count"`
    Avg   float64 `bson:"avgAge"`
}

var stats []CityStats

// 按城市分组统计
err := gmongo.NewAggregation(users).
    Group("$city", bson.M{
        "count":  gmongo.Sum(1),
        "avgAge": gmongo.Avg("$age"),
    }).
    Sort(bson.D{{Key: "count", Value: -1}}).
    Execute(&stats)

for _, stat := range stats {
    fmt.Printf("%s: %d 人, 平均 %.1f 岁\n",
        stat.City, stat.Count, stat.Avg)
}
```

### 5. 事务操作

```go
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
```

## 常用操作速查

### 插入

```go
// 插入单个
users.InsertOne(ctx, user)

// 插入多个
users.InsertMany(ctx, []interface{}{user1, user2, user3})

// Upsert（更新或插入）
users.Upsert(ctx, filter, update)
```

### 查询

```go
// 查询单个
var user User
users.FindOne(ctx, filter).Decode(&user)

// 查询多个
var users []User
users.Find(ctx, filter).All(ctx, &users)

// 链式查询
users.Query().WhereEq("status", "active").Find(&users)

// 分页查询
total, _ := users.Paginate(ctx, filter, page, pageSize, &users)

// 统计数量
count, _ := users.Count(ctx, filter)

// 检查存在
exists, _ := users.Exists(ctx, filter)
```

### 更新

```go
// 更新单个
update := bson.M{"$set": bson.M{"age": 26}}
users.UpdateOne(ctx, filter, update)

// 更新多个
users.UpdateMany(ctx, filter, update)

// 根据 ID 更新
users.UpdateByID(ctx, id, update)

// 使用查询构建器
users.Query().WhereEq("name", "张三").Update(update)
```

### 删除

```go
// 删除单个
users.DeleteOne(ctx, filter)

// 删除多个
users.DeleteMany(ctx, filter)

// 根据 ID 删除
users.DeleteByID(ctx, id)

// 使用查询构建器
users.Query().WhereEq("status", "inactive").DeleteMany()
```

### 索引

```go
indexes := users.Indexes()

// 创建唯一索引
indexes.CreateUniqueIndex(ctx, bson.D{{Key: "email", Value: 1}})

// 创建复合索引
indexes.CreateCompoundIndex(ctx, map[string]int{
    "city": 1,
    "age": -1,
}, false)

// 创建 TTL 索引（自动过期）
indexes.CreateTTLIndex(ctx, "created_at", 24*time.Hour)
```

## 查询条件速查

```go
// 比较操作
.WhereEq("age", 25)           // age == 25
.WhereNe("age", 25)           // age != 25
.WhereGt("age", 18)           // age > 18
.WhereGte("age", 18)          // age >= 18
.WhereLt("age", 60)           // age < 60
.WhereLte("age", 60)          // age <= 60

// 范围操作
.WhereIn("status", []string{"active", "pending"})
.WhereNin("status", []string{"deleted"})
.WhereBetween("age", 18, 60)

// 逻辑操作
.Or(bson.M{"age": 18}, bson.M{"age": 25})
.And(bson.M{"city": "北京"}, bson.M{"age": bson.M{"$gt": 18}})

// 字段操作
.WhereExists("email", true)   // email 字段存在
.WhereRegex("name", "^张")    // 正则匹配
.WhereType("age", "int")      // 类型匹配

// 数组操作
.WhereSize("tags", 3)         // 数组长度为 3
.WhereAll("tags", []string{"golang", "mongodb"})
.WhereElemMatch("scores", bson.M{"$gte": 80})
```

## 聚合函数速查

```go
// 统计函数
gmongo.Sum(1)              // 计数
gmongo.Sum("$amount")      // 求和
gmongo.Avg("$age")         // 平均值
gmongo.Min("$price")       // 最小值
gmongo.Max("$price")       // 最大值

// 数组函数
gmongo.Push("$item")       // 添加到数组
gmongo.AddToSet("$tag")    // 添加到集合（去重）
gmongo.First("$name")      // 第一个值
gmongo.Last("$name")       // 最后一个值
```

## 配置选项速查

```go
config := &gmongo.Config{
    Host:     "localhost",
    Port:     27017,
    Database: "myapp",
    Username: "admin",
    Password: "password",

    // 连接池
    MaxPoolSize:     100,
    MinPoolSize:     10,
    MaxConnIdleTime: 10 * time.Minute,

    // 超时
    ConnectTimeout: 10 * time.Second,
    SocketTimeout:  30 * time.Second,

    // 读写配置
    ReadPreference: "primary",
    WriteConcern:   "majority",

    // 压缩
    Compressors: []string{"snappy", "zlib"},
}
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

1. **使用 Context** - 所有操作都传递 context，支持超时和取消
2. **关闭资源** - 使用 defer 关闭游标和客户端
3. **错误处理** - 妥善处理所有错误
4. **使用索引** - 为常用查询创建索引
5. **批量操作** - 使用 InsertMany 和 BulkWrite 提高性能
6. **投影** - 只查询需要的字段，减少网络传输

## 下一步

- 阅读完整文档: [README.md](./README.md)
- 查看设计文档: [DESIGN.md](./DESIGN.md)
- 运行示例代码: [examples/basic_example.go](./examples/basic_example.go)
- 了解高级特性: 聚合管道、事务、钩子系统

## 获取帮助

- GitHub Issues: 提交问题和建议
- 文档: 查看完整文档
- 示例: 查看更多示例代码

祝你使用愉快！🎉
