//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nicexiaonie/gmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User 测试用户结构
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Age       int                `bson:"age"`
	Tags      []string           `bson:"tags,omitempty"`
	Address   *Address           `bson:"address,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	ExpireAt  *time.Time         `bson:"expire_at,omitempty"`
}

// Address 地址结构
type Address struct {
	City     string  `bson:"city"`
	Province string  `bson:"province"`
	Location GeoJSON `bson:"location,omitempty"`
}

// GeoJSON 地理位置
type GeoJSON struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

// Product 测试产品结构
type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Price    float64            `bson:"price"`
	Stock    int                `bson:"stock"`
	Category string             `bson:"category"`
	Tags     []string           `bson:"tags,omitempty"`
}

// Order 订单结构
type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ProductID primitive.ObjectID `bson:"product_id"`
	Quantity  int                `bson:"quantity"`
	Total     float64            `bson:"total"`
	Status    string             `bson:"status"`
	CreatedAt time.Time          `bson:"created_at"`
}

/*

 docker run -it --name mongo -p 27017:27017 -v /Users/nieyuanpei/Server/mongo/data:/data/db -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=123456 -d 9e52bf8768ae

mongosh "mongodb://root:123456@127.0.0.1:27017"

*/

func main() {
	ctx := context.Background()

	// ========== 基础功能测试 ==========
	log.Println("\n========== 基础功能测试 ==========")
	testConnection(ctx)
	testBasicCRUD(ctx)
	testQueryBuilder(ctx)

	// ========== 高级查询测试 ==========
	log.Println("\n========== 高级查询测试 ==========")
	testAdvancedQueries(ctx)
	testProjectionAndSort(ctx)
	testArrayOperators(ctx)
	testFindOneAndModify(ctx)

	// ========== 聚合与索引测试 ==========
	log.Println("\n========== 聚合与索引测试 ==========")
	testAggregation(ctx)
	testComplexAggregation(ctx)
	testIndexManagement(ctx)
	testSpecialIndexes(ctx)

	// ========== 事务测试 ==========
	log.Println("\n========== 事务测试 ==========")
	testTransaction(ctx)
	testTransactionRollback(ctx)
	testTransactionIsolation(ctx)

	// ========== 批量操作测试 ==========
	log.Println("\n========== 批量操作测试 ==========")
	testBulkOperations(ctx)

	// ========== 并发与性能测试 ==========
	log.Println("\n========== 并发与性能测试 ==========")
	testConcurrency(ctx)
	testPerformance(ctx)

	// ========== 边界条件测试 ==========
	log.Println("\n========== 边界条件测试 ==========")
	testEdgeCases(ctx)

	// ========== 错误处理测试 ==========
	log.Println("\n========== 错误处理测试 ==========")
	testErrorHandling(ctx)

	// ========== 工具函数测试 ==========
	log.Println("\n========== 工具函数测试 ==========")
	testUtilityFunctions()

	// ========== 变更流测试 ==========
	log.Println("\n========== 变更流测试 ==========")
	testChangeStream(ctx)

	// 关闭连接
	if err := gmongo.Close(); err != nil {
		log.Printf("❌ 关闭连接失败: %v", err)
	}
	log.Println("\n✅ 所有测试完成")
}

// testConnection 测试连接管理
func testConnection(ctx context.Context) {
	// Connect
	config := gmongo.DefaultConfig()
	config.Host = "127.0.0.1"
	config.Port = 27017
	config.Database = "test_db"
	config.Username = "root"
	config.Password = "123456"
	config.AuthSource = "admin"

	if err := gmongo.Connect(config); err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	fmt.Println("✓ Connect 成功")

	// Ping
	if err := gmongo.Ping(ctx); err != nil {
		log.Printf("Ping 失败: %v", err)
	} else {
		fmt.Println("✓ Ping 成功")
	}

	// HealthCheck
	if err := gmongo.HealthCheck(); err != nil {
		log.Printf("健康检查失败: %v", err)
	} else {
		fmt.Println("✓ HealthCheck 成功")
	}

	// Stats
	stats := gmongo.Stats()
	fmt.Printf("✓ Stats: %s\n", stats)

	// ListDatabaseNames
	dbNames, err := gmongo.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Printf("列出数据库失败: %v", err)
	} else {
		fmt.Printf("✓ ListDatabaseNames: %v\n", dbNames)
	}
}

// testBasicCRUD 测试基本 CRUD 操作
func testBasicCRUD(ctx context.Context) {
	coll := gmongo.Coll("users")

	// 清空集合
	coll.DeleteMany(ctx, bson.M{})

	// InsertOne
	user := User{
		Name:      "张三",
		Email:     "zhangsan@example.com",
		Age:       25,
		Tags:      []string{"developer", "golang"},
		CreatedAt: time.Now(),
	}
	result, err := coll.InsertOne(ctx, user)
	if err != nil {
		log.Printf("InsertOne 失败: %v", err)
	} else {
		fmt.Printf("✓ InsertOne 成功, ID: %v\n", result.InsertedID)
	}

	// InsertMany
	users := []interface{}{
		User{Name: "李四", Email: "lisi@example.com", Age: 30, CreatedAt: time.Now()},
		User{Name: "王五", Email: "wangwu@example.com", Age: 28, CreatedAt: time.Now()},
		User{Name: "赵六", Email: "zhaoliu@example.com", Age: 35, CreatedAt: time.Now()},
	}
	manyResult, err := coll.InsertMany(ctx, users)
	if err != nil {
		log.Printf("InsertMany 失败: %v", err)
	} else {
		fmt.Printf("✓ InsertMany 成功, 插入 %d 条记录\n", len(manyResult.InsertedIDs))
	}

	// FindOne
	var foundUser User
	err = coll.FindOne(ctx, bson.M{"name": "张三"}).Decode(&foundUser)
	if err != nil {
		log.Printf("FindOne 失败: %v", err)
	} else {
		fmt.Printf("✓ FindOne 成功: %s\n", foundUser.Name)
	}

	// Find
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Find 失败: %v", err)
	} else {
		var allUsers []User
		if err = cursor.All(ctx, &allUsers); err != nil {
			log.Printf("解码失败: %v", err)
		} else {
			fmt.Printf("✓ Find 成功, 找到 %d 条记录\n", len(allUsers))
		}
	}

	// FindAll
	var allUsers []User
	err = coll.FindAll(ctx, bson.M{}, &allUsers)
	if err != nil {
		log.Printf("FindAll 失败: %v", err)
	} else {
		fmt.Printf("✓ FindAll 成功, 找到 %d 条记录\n", len(allUsers))
	}

	// UpdateOne
	updateResult, err := coll.UpdateOne(ctx,
		bson.M{"name": "张三"},
		bson.M{"$set": bson.M{"age": 26, "updated_at": time.Now()}})
	if err != nil {
		log.Printf("UpdateOne 失败: %v", err)
	} else {
		fmt.Printf("✓ UpdateOne 成功, 修改 %d 条记录\n", updateResult.ModifiedCount)
	}

	// UpdateMany
	updateManyResult, err := coll.UpdateMany(ctx,
		bson.M{"age": bson.M{"$gte": 30}},
		bson.M{"$set": bson.M{"updated_at": time.Now()}})
	if err != nil {
		log.Printf("UpdateMany 失败: %v", err)
	} else {
		fmt.Printf("✓ UpdateMany 成功, 修改 %d 条记录\n", updateManyResult.ModifiedCount)
	}

	// CountDocuments
	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("CountDocuments 失败: %v", err)
	} else {
		fmt.Printf("✓ CountDocuments: %d\n", count)
	}

	// Distinct
	distinctAges, err := coll.Distinct(ctx, "age", bson.M{})
	if err != nil {
		log.Printf("Distinct 失败: %v", err)
	} else {
		fmt.Printf("✓ Distinct ages: %v\n", distinctAges)
	}

	// FindOneAndUpdate
	var updatedUser User
	err = coll.FindOneAndUpdate(ctx,
		bson.M{"name": "李四"},
		bson.M{"$set": bson.M{"age": 31}}).Decode(&updatedUser)
	if err != nil {
		log.Printf("FindOneAndUpdate 失败: %v", err)
	} else {
		fmt.Printf("✓ FindOneAndUpdate 成功: %s\n", updatedUser.Name)
	}

	// DeleteOne
	deleteResult, err := coll.DeleteOne(ctx, bson.M{"name": "赵六"})
	if err != nil {
		log.Printf("DeleteOne 失败: %v", err)
	} else {
		fmt.Printf("✓ DeleteOne 成功, 删除 %d 条记录\n", deleteResult.DeletedCount)
	}
}

// testQueryBuilder 测试查询构建器
func testQueryBuilder(ctx context.Context) {
	coll := gmongo.Coll("users")

	// 链式查询
	query := coll.Query().
		Context(ctx).
		WhereGte("age", 25).
		WhereLte("age", 35).
		Select("name", "email", "age").
		Sort("-age").
		Limit(10)

	var users []User
	err := query.Find(&users)
	if err != nil {
		log.Printf("查询构建器 Find 失败: %v", err)
	} else {
		fmt.Printf("✓ 查询构建器 Find 成功, 找到 %d 条记录\n", len(users))
	}

	// WhereIn
	query2 := coll.Query().Context(ctx).WhereIn("age", []interface{}{25, 30, 35})
	count, err := query2.Count()
	if err != nil {
		log.Printf("WhereIn Count 失败: %v", err)
	} else {
		fmt.Printf("✓ WhereIn Count: %d\n", count)
	}

	// WhereBetween
	query3 := coll.Query().Context(ctx).WhereBetween("age", 25, 35)
	exists, err := query3.Exists()
	if err != nil {
		log.Printf("WhereBetween Exists 失败: %v", err)
	} else {
		fmt.Printf("✓ WhereBetween Exists: %v\n", exists)
	}

	// WhereRegex
	query4 := coll.Query().Context(ctx).WhereRegex("email", ".*@example.com", "i")
	var emailUsers []User
	err = query4.Find(&emailUsers)
	if err != nil {
		log.Printf("WhereRegex 失败: %v", err)
	} else {
		fmt.Printf("✓ WhereRegex 成功, 找到 %d 条记录\n", len(emailUsers))
	}

	// Or 查询
	query5 := coll.Query().Context(ctx).Or(
		bson.M{"age": bson.M{"$lt": 26}},
		bson.M{"age": bson.M{"$gt": 30}},
	)
	orCount, err := query5.Count()
	if err != nil {
		log.Printf("Or Count 失败: %v", err)
	} else {
		fmt.Printf("✓ Or Count: %d\n", orCount)
	}

	// 分页查询 - 使用 Collection 的 Paginate 方法
	var pageUsers []User
	total, err := coll.Paginate(ctx, bson.M{}, 1, 2, &pageUsers)
	if err != nil {
		log.Printf("Paginate 失败: %v", err)
	} else {
		fmt.Printf("✓ Paginate 成功, 总数: %d, 当前页数据: %d 条\n", total, len(pageUsers))
	}
}

// testAggregation 测试聚合操作
func testAggregation(ctx context.Context) {
	coll := gmongo.Coll("users")

	// 基本聚合
	pipeline := gmongo.NewAggregation(coll).
		Match(bson.M{"age": bson.M{"$gte": 25}}).
		Group("$age", bson.M{
			"count":  bson.M{"$sum": 1},
			"avgAge": gmongo.Avg("$age"),
			"names":  gmongo.Push("$name"),
		}).
		SortBy("-count").
		Limit(10)

	var results []bson.M
	err := pipeline.Execute(&results)
	if err != nil {
		log.Printf("聚合 Execute 失败: %v", err)
	} else {
		fmt.Printf("✓ 聚合 Execute 成功, 结果数: %d\n", len(results))
	}

	// Project 聚合
	pipeline2 := gmongo.NewAggregation(coll).
		Project(bson.M{
			"name": 1,
			"age":  1,
			"ageGroup": bson.M{
				"$cond": bson.A{
					bson.M{"$gte": bson.A{"$age", 30}},
					"senior",
					"junior",
				},
			},
		})

	var projectResults []bson.M
	err = pipeline2.Execute(&projectResults)
	if err != nil {
		log.Printf("Project 聚合失败: %v", err)
	} else {
		fmt.Printf("✓ Project 聚合成功, 结果数: %d\n", len(projectResults))
	}

	// AddFields 聚合
	pipeline3 := gmongo.NewAggregation(coll).
		AddFields(bson.M{
			"fullInfo": bson.M{"$concat": bson.A{"$name", " - ", "$email"}},
		})

	var addFieldsResults []bson.M
	err = pipeline3.Execute(&addFieldsResults)
	if err != nil {
		log.Printf("AddFields 聚合失败: %v", err)
	} else {
		fmt.Printf("✓ AddFields 聚合成功, 结果数: %d\n", len(addFieldsResults))
	}

	// Count 聚合
	pipeline4 := gmongo.NewAggregation(coll).
		Match(bson.M{"age": bson.M{"$gte": 25}}).
		Count("total")

	var countResult bson.M
	err = pipeline4.One(&countResult)
	if err != nil {
		log.Printf("Count 聚合失败: %v", err)
	} else {
		fmt.Printf("✓ Count 聚合成功: %v\n", countResult)
	}
}

// testIndexManagement 测试索引管理
func testIndexManagement(ctx context.Context) {
	coll := gmongo.Coll("users")
	indexMgr := coll.Indexes()

	// CreateIndex - 单字段索引
	_, err := indexMgr.CreateIndex(ctx, bson.D{{Key: "email", Value: 1}}, false)
	if err != nil {
		log.Printf("CreateIndex 失败: %v", err)
	} else {
		fmt.Println("✓ CreateIndex (email) 成功")
	}

	// CreateUniqueIndex - 唯一索引
	_, err = indexMgr.CreateUniqueIndex(ctx, bson.D{{Key: "email", Value: 1}})
	if err != nil {
		log.Printf("CreateUniqueIndex 失败: %v", err)
	} else {
		fmt.Println("✓ CreateUniqueIndex (email) 成功")
	}

	// CreateCompoundIndex - 复合索引
	_, err = indexMgr.CreateCompoundIndex(ctx, map[string]int{"name": 1, "age": -1}, false)
	if err != nil {
		log.Printf("CreateCompoundIndex 失败: %v", err)
	} else {
		fmt.Println("✓ CreateCompoundIndex (name, age) 成功")
	}

	// List - 列出所有索引
	cursor, err := indexMgr.List(ctx)
	if err != nil {
		log.Printf("List 索引失败: %v", err)
	} else {
		var indexes []bson.M
		if err = cursor.All(ctx, &indexes); err != nil {
			log.Printf("解析索引失败: %v", err)
		} else {
			fmt.Printf("✓ List 索引成功, 共 %d 个索引\n", len(indexes))
		}
	}

	// IndexExists - 检查索引是否存在
	exists, err := indexMgr.IndexExists(ctx, "email_1")
	if err != nil {
		log.Printf("IndexExists 失败: %v", err)
	} else {
		fmt.Printf("✓ IndexExists (email_1): %v\n", exists)
	}

	// IndexBuilder - 使用构建器创建索引
	builder := gmongo.NewIndexBuilder().
		AddField("created_at", -1).
		Unique(false).
		Background(true).
		Name("created_at_desc")

	_, err = indexMgr.CreateOne(ctx, builder.Build())
	if err != nil {
		log.Printf("IndexBuilder CreateOne 失败: %v", err)
	} else {
		fmt.Println("✓ IndexBuilder CreateOne 成功")
	}
}

// testTransaction 测试事务操作
func testTransaction(ctx context.Context) {
	// 简单事务
	err := gmongo.Tx(ctx, func(sessCtx mongo.SessionContext) error {
		coll := gmongo.Coll("users")

		// 在事务中插入
		_, err := coll.InsertOne(sessCtx, User{
			Name:      "事务用户",
			Email:     "tx@example.com",
			Age:       40,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		// 在事务中更新
		_, err = coll.UpdateOne(sessCtx,
			bson.M{"email": "tx@example.com"},
			bson.M{"$set": bson.M{"age": 41}})

		return err
	})
	if err != nil {
		log.Printf("Tx 失败: %v", err)
	} else {
		fmt.Println("✓ Tx 成功")
	}

	// WithTransaction - 带返回值的事务
	result, err := gmongo.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		coll := gmongo.Coll("users")
		count, err := coll.CountDocuments(sessCtx, bson.M{})
		return count, err
	})
	if err != nil {
		log.Printf("WithTransaction 失败: %v", err)
	} else {
		fmt.Printf("✓ WithTransaction 成功, 结果: %v\n", result)
	}

	// WithReadConcernMajority - 使用 majority 读关注
	err = gmongo.WithReadConcernMajority(ctx, func(sessCtx mongo.SessionContext) error {
		coll := gmongo.Coll("users")
		var user User
		return coll.FindOne(sessCtx, bson.M{}).Decode(&user)
	})
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("WithReadConcernMajority 失败: %v", err)
	} else {
		fmt.Println("✓ WithReadConcernMajority 成功")
	}

	// TxOptionsBuilder - 使用事务选项构建器
	txOpts := gmongo.NewTxOptions().
		ReadConcernMajority().
		WriteConcernMajority().
		Build()

	err = gmongo.TxWithOptions(ctx, txOpts, func(sessCtx mongo.SessionContext) error {
		coll := gmongo.Coll("users")
		_, err := coll.UpdateMany(sessCtx,
			bson.M{"age": bson.M{"$gte": 40}},
			bson.M{"$set": bson.M{"updated_at": time.Now()}})
		return err
	})
	if err != nil {
		log.Printf("TxWithOptions 失败: %v", err)
	} else {
		fmt.Println("✓ TxWithOptions 成功")
	}

	// WithIsolationLevel - 使用隔离级别
	err = gmongo.WithIsolationLevel(ctx, gmongo.IsolationLevelSnapshot, func(sessCtx mongo.SessionContext) error {
		coll := gmongo.Coll("users")
		var users []User
		return coll.FindAll(sessCtx, bson.M{}, &users)
	})
	if err != nil {
		log.Printf("WithIsolationLevel 失败: %v", err)
	} else {
		fmt.Println("✓ WithIsolationLevel 成功")
	}
}

// testUtilityFunctions 测试工具函数
func testUtilityFunctions() {
	// NewObjectID
	oid := gmongo.NewObjectID()
	fmt.Printf("✓ NewObjectID: %s\n", oid.Hex())

	// ObjectIDFromHex
	oid2, err := gmongo.ObjectIDFromHex(oid.Hex())
	if err != nil {
		log.Printf("ObjectIDFromHex 失败: %v", err)
	} else {
		fmt.Printf("✓ ObjectIDFromHex: %s\n", oid2.Hex())
	}

	// IsValidObjectID
	valid := gmongo.IsValidObjectID(oid.Hex())
	fmt.Printf("✓ IsValidObjectID: %v\n", valid)

	// BuildFilter
	filter := gmongo.BuildFilter(bson.M{
		"age":  gmongo.Gte(25),
		"name": gmongo.Regex("张.*", "i"),
	})
	fmt.Printf("✓ BuildFilter: %v\n", filter)

	// BuildUpdate
	update := gmongo.BuildUpdate(bson.M{"age": 30, "login_count": bson.M{"$inc": 1}})
	fmt.Printf("✓ BuildUpdate: %v\n", update)

	// BuildSort
	sort := gmongo.BuildSort("-age", "name")
	fmt.Printf("✓ BuildSort: %v\n", sort)

	// TimeRange - Today
	today := gmongo.Today()
	fmt.Printf("✓ Today: %v - %v\n", today.Start, today.End)

	// TimeRange - ThisWeek
	thisWeek := gmongo.ThisWeek()
	fmt.Printf("✓ ThisWeek: %v - %v\n", thisWeek.Start, thisWeek.End)

	// ToDoc / FromDoc
	user := User{Name: "测试", Email: "test@example.com", Age: 25}
	doc, err := gmongo.ToDoc(user)
	if err != nil {
		log.Printf("ToDoc 失败: %v", err)
	} else {
		fmt.Println("✓ ToDoc 成功")

		var user2 User
		err = gmongo.FromDoc(doc, &user2)
		if err != nil {
			log.Printf("FromDoc 失败: %v", err)
		} else {
			fmt.Printf("✓ FromDoc 成功: %s\n", user2.Name)
		}
	}
}

// testAdvancedQueries 测试高级查询
func testAdvancedQueries(ctx context.Context) {
	coll := gmongo.Coll("users")

	// FindOneByID
	var user User
	err := coll.FindAll(ctx, bson.M{}, &user)
	if err == nil && len(user.Name) > 0 {
		var foundUser User
		err = coll.FindOneByID(ctx, user.ID, &foundUser)
		if err != nil {
			log.Printf("FindOneByID 失败: %v", err)
		} else {
			fmt.Printf("✓ FindOneByID 成功: %s\n", foundUser.Name)
		}
	}

	// Upsert
	upsertResult, err := coll.Upsert(ctx,
		bson.M{"email": "upsert@example.com"},
		bson.M{"$set": bson.M{
			"name":       "Upsert用户",
			"age":        50,
			"created_at": time.Now(),
		}})
	if err != nil {
		log.Printf("Upsert 失败: %v", err)
	} else {
		fmt.Printf("✓ Upsert 成功, UpsertedID: %v\n", upsertResult.UpsertedID)
	}

	// Paginate
	var pageUsers []User
	total, err := coll.Paginate(ctx, bson.M{}, 1, 5, &pageUsers)
	if err != nil {
		log.Printf("Paginate 失败: %v", err)
	} else {
		fmt.Printf("✓ Paginate 成功, 总数: %d, 当前页数据: %d 条\n", total, len(pageUsers))
	}

	// WhereExists
	query := coll.Query().Context(ctx).WhereExists("tags", true)
	count, err := query.Count()
	if err != nil {
		log.Printf("WhereExists 失败: %v", err)
	} else {
		fmt.Printf("✓ WhereExists Count: %d\n", count)
	}

	// WhereSize
	query2 := coll.Query().Context(ctx).WhereSize("tags", 2)
	count2, err := query2.Count()
	if err != nil {
		log.Printf("WhereSize 失败: %v", err)
	} else {
		fmt.Printf("✓ WhereSize Count: %d\n", count2)
	}
}

// testBulkOperations 测试批量操作
func testBulkOperations(ctx context.Context) {
	coll := gmongo.Coll("products")

	// 清空集合
	coll.DeleteMany(ctx, bson.M{})

	// 准备批量操作
	operations := []mongo.WriteModel{
		mongo.NewInsertOneModel().SetDocument(Product{
			Name: "产品A", Price: 100.0, Stock: 10, Category: "电子",
		}),
		mongo.NewInsertOneModel().SetDocument(Product{
			Name: "产品B", Price: 200.0, Stock: 20, Category: "电子",
		}),
		mongo.NewInsertOneModel().SetDocument(Product{
			Name: "产品C", Price: 150.0, Stock: 15, Category: "家居",
		}),
		mongo.NewUpdateOneModel().
			SetFilter(bson.M{"name": "产品A"}).
			SetUpdate(bson.M{"$set": bson.M{"price": 110.0}}),
		mongo.NewDeleteOneModel().
			SetFilter(bson.M{"name": "产品C"}),
	}

	// BulkWrite
	bulkResult, err := coll.BulkWrite(ctx, operations)
	if err != nil {
		log.Printf("BulkWrite 失败: %v", err)
	} else {
		fmt.Printf("✓ BulkWrite 成功, 插入: %d, 更新: %d, 删除: %d\n",
			bulkResult.InsertedCount, bulkResult.ModifiedCount, bulkResult.DeletedCount)
	}
}

// testChangeStream 测试变更流
func testChangeStream(ctx context.Context) {
	coll := gmongo.Coll("users")

	// 创建变更流（仅演示，不实际监听）
	pipeline := mongo.Pipeline{}

	// Watch
	changeStream, err := coll.Watch(ctx, pipeline)
	if err != nil {
		log.Printf("Watch 失败: %v", err)
		return
	}
	defer changeStream.Close(ctx)

	fmt.Println("✓ Watch 创建成功")

	// 注意：实际使用时需要在 goroutine 中监听变更
	// 这里只是演示 API 的使用
}

// testProjectionAndSort 测试投影和排序
func testProjectionAndSort(ctx context.Context) {
	coll := gmongo.Coll("users")

	// 投影：只返回指定字段
	var users []User
	err := coll.Query().Context(ctx).
		Select("name", "email", "age").
		Sort("-age").
		Limit(5).
		Find(&users)
	if err != nil {
		log.Printf("投影查询失败: %v", err)
	} else {
		fmt.Printf("✓ 投影查询成功，返回 %d 条记录\n", len(users))
	}

	// 排除字段
	var users2 []User
	err = coll.Query().Context(ctx).
		Omit("tags", "created_at").
		Sort("name").
		Limit(3).
		Find(&users2)
	if err != nil {
		log.Printf("排除字段查询失败: %v", err)
	} else {
		fmt.Printf("✓ 排除字段查询成功，返回 %d 条记录\n", len(users2))
	}

	// 多字段排序
	var users3 []User
	err = coll.Query().Context(ctx).
		Sort("-age", "name").
		Limit(5).
		Find(&users3)
	if err != nil {
		log.Printf("多字段排序失败: %v", err)
	} else {
		fmt.Printf("✓ 多字段排序成功，返回 %d 条记录\n", len(users3))
	}
}

// testArrayOperators 测试数组操作符
func testArrayOperators(ctx context.Context) {
	coll := gmongo.Coll("users")

	// 插入测试数据
	coll.InsertOne(ctx, User{
		Name:      "数组测试用户1",
		Email:     "array1@example.com",
		Age:       25,
		Tags:      []string{"golang", "mongodb", "redis"},
		CreatedAt: time.Now(),
	})
	coll.InsertOne(ctx, User{
		Name:      "数组测试用户2",
		Email:     "array2@example.com",
		Age:       26,
		Tags:      []string{"python", "mongodb", "mysql"},
		CreatedAt: time.Now(),
	})

	// $all - 包含所有指定元素
	var users1 []User
	err := coll.FindAll(ctx, bson.M{"tags": bson.M{"$all": []string{"mongodb", "golang"}}}, &users1)
	if err != nil {
		log.Printf("$all 查询失败: %v", err)
	} else {
		fmt.Printf("✓ $all 查询成功，找到 %d 条记录\n", len(users1))
	}

	// $elemMatch - 数组元素匹配
	var users2 []User
	err = coll.FindAll(ctx, bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "mongodb"}}}, &users2)
	if err != nil {
		log.Printf("$elemMatch 查询失败: %v", err)
	} else {
		fmt.Printf("✓ $elemMatch 查询成功，找到 %d 条记录\n", len(users2))
	}

	// $size - 数组大小
	var users3 []User
	err = coll.FindAll(ctx, bson.M{"tags": bson.M{"$size": 3}}, &users3)
	if err != nil {
		log.Printf("$size 查询失败: %v", err)
	} else {
		fmt.Printf("✓ $size 查询成功，找到 %d 条记录\n", len(users3))
	}

	// $in - 数组包含任一元素
	var users4 []User
	err = coll.FindAll(ctx, bson.M{"tags": bson.M{"$in": []string{"python", "java"}}}, &users4)
	if err != nil {
		log.Printf("$in 查询失败: %v", err)
	} else {
		fmt.Printf("✓ $in 查询成功，找到 %d 条记录\n", len(users4))
	}
}

// testFindOneAndModify 测试原子操作
func testFindOneAndModify(ctx context.Context) {
	coll := gmongo.Coll("users")

	// FindOneAndUpdate - 查找并更新
	var updatedUser User
	result := coll.FindOneAndUpdate(ctx,
		bson.M{"email": "array1@example.com"},
		bson.M{"$set": bson.M{"age": 30}})
	err := result.Decode(&updatedUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("FindOneAndUpdate 失败: %v", err)
	} else if err == nil {
		fmt.Printf("✓ FindOneAndUpdate 成功: %s\n", updatedUser.Name)
	}

	// FindOneAndDelete - 查找并删除
	var deletedUser User
	result = coll.FindOneAndDelete(ctx, bson.M{"email": "array2@example.com"})
	err = result.Decode(&deletedUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("FindOneAndDelete 失败: %v", err)
	} else if err == nil {
		fmt.Printf("✓ FindOneAndDelete 成功: %s\n", deletedUser.Name)
	}

	// FindOneAndReplace - 查找并替换
	var replacedUser User
	result = coll.FindOneAndReplace(ctx,
		bson.M{"email": "array1@example.com"},
		User{
			Name:      "替换后的用户",
			Email:     "array1@example.com",
			Age:       35,
			CreatedAt: time.Now(),
		})
	err = result.Decode(&replacedUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("FindOneAndReplace 失败: %v", err)
	} else if err == nil {
		fmt.Printf("✓ FindOneAndReplace 成功: %s\n", replacedUser.Name)
	}
}

// testComplexAggregation 测试复杂聚合
func testComplexAggregation(ctx context.Context) {
	userColl := gmongo.Coll("users")
	orderColl := gmongo.Coll("orders")

	// 清空并准备测试数据
	orderColl.DeleteMany(ctx, bson.M{})

	// 获取一些用户ID
	var users []User
	userColl.Query().Context(ctx).Limit(3).Find(&users)
	if len(users) == 0 {
		log.Println("没有用户数据，跳过复杂聚合测试")
		return
	}

	// 插入订单数据
	for i, user := range users {
		orderColl.InsertOne(ctx, Order{
			UserID:    user.ID,
			ProductID: primitive.NewObjectID(),
			Quantity:  i + 1,
			Total:     float64((i + 1) * 100),
			Status:    "completed",
			CreatedAt: time.Now(),
		})
	}

	// $lookup - 关联查询
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"status": "completed"}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user_info",
		}}},
		{{Key: "$unwind", Value: "$user_info"}},
		{{Key: "$project", Value: bson.M{
			"quantity":  1,
			"total":     1,
			"user_name": "$user_info.name",
		}}},
	}

	cursor, err := orderColl.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("$lookup 聚合失败: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		log.Printf("读取聚合结果失败: %v", err)
	} else {
		fmt.Printf("✓ $lookup 聚合成功，返回 %d 条记录\n", len(results))
	}

	// $group + $sort - 分组统计
	pipeline2 := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":          "$status",
			"total_count":  bson.M{"$sum": 1},
			"total_amount": bson.M{"$sum": "$total"},
			"avg_amount":   bson.M{"$avg": "$total"},
		}}},
		{{Key: "$sort", Value: bson.M{"total_amount": -1}}},
	}

	cursor2, err := orderColl.Aggregate(ctx, pipeline2)
	if err != nil {
		log.Printf("$group 聚合失败: %v", err)
		return
	}
	defer cursor2.Close(ctx)

	var results2 []bson.M
	if err = cursor2.All(ctx, &results2); err != nil {
		log.Printf("读取聚合结果失败: %v", err)
	} else {
		fmt.Printf("✓ $group 聚合成功，返回 %d 条记录\n", len(results2))
	}
}

// testSpecialIndexes 测试特殊索引
func testSpecialIndexes(ctx context.Context) {
	coll := gmongo.Coll("users")

	// TTL 索引 - 自动过期
	ttlIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "expire_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(3600),
	}
	_, err := coll.Indexes().CreateMany(ctx, []mongo.IndexModel{ttlIndex})
	if err != nil {
		log.Printf("创建 TTL 索引失败: %v", err)
	} else {
		fmt.Println("✓ 创建 TTL 索引成功")
	}

	// 文本索引 - 全文搜索
	textIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}, {Key: "email", Value: "text"}},
	}
	_, err = coll.Indexes().CreateMany(ctx, []mongo.IndexModel{textIndex})
	if err != nil {
		log.Printf("创建文本索引失败: %v", err)
	} else {
		fmt.Println("✓ 创建文本索引成功")
	}

	// 使用文本搜索
	var users []User
	err = coll.FindAll(ctx, bson.M{"$text": bson.M{"$search": "test"}}, &users)
	if err != nil {
		log.Printf("文本搜索失败: %v", err)
	} else {
		fmt.Printf("✓ 文本搜索成功，找到 %d 条记录\n", len(users))
	}

	// 地理空间索引
	locationColl := gmongo.Coll("locations")
	geoIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "address.location", Value: "2dsphere"}},
	}
	_, err = locationColl.Indexes().CreateMany(ctx, []mongo.IndexModel{geoIndex})
	if err != nil {
		log.Printf("创建地理空间索引失败: %v", err)
	} else {
		fmt.Println("✓ 创建地理空间索引成功")
	}

	// 插入带地理位置的数据
	locationColl.InsertOne(ctx, User{
		Name:  "位置测试用户",
		Email: "location@example.com",
		Age:   30,
		Address: &Address{
			City:     "北京",
			Province: "北京",
			Location: GeoJSON{
				Type:        "Point",
				Coordinates: []float64{116.404, 39.915}, // 经度, 纬度
			},
		},
		CreatedAt: time.Now(),
	})

	// 地理位置查询 - 查找附近的点
	var nearUsers []User
	err = locationColl.FindAll(ctx, bson.M{
		"address.location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{116.404, 39.915},
				},
				"$maxDistance": 1000, // 1000米内
			},
		},
	}, &nearUsers)
	if err != nil {
		log.Printf("地理位置查询失败: %v", err)
	} else {
		fmt.Printf("✓ 地理位置查询成功，找到 %d 条记录\n", len(nearUsers))
	}
}

// testTransactionRollback 测试事务回滚
func testTransactionRollback(ctx context.Context) {
	coll := gmongo.Coll("users")

	// 记录初始数量
	initialCount, _ := coll.CountDocuments(ctx, bson.M{})

	// 执行会失败的事务
	err := gmongo.Tx(ctx, func(sessCtx mongo.SessionContext) error {
		// 插入一条记录
		_, err := coll.InsertOne(sessCtx, User{
			Name:      "回滚测试用户",
			Email:     "rollback@example.com",
			Age:       40,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}

		// 故意返回错误，触发回滚
		return fmt.Errorf("模拟错误，触发回滚")
	})

	if err != nil {
		fmt.Printf("✓ 事务回滚测试成功，错误: %v\n", err)
	}

	// 验证数据未被插入
	finalCount, _ := coll.CountDocuments(ctx, bson.M{})
	if initialCount == finalCount {
		fmt.Println("✓ 事务回滚验证成功，数据未被插入")
	} else {
		log.Printf("❌ 事务回滚验证失败，数据被插入了")
	}
}

// testTransactionIsolation 测试事务隔离级别
func testTransactionIsolation(ctx context.Context) {
	// Snapshot 隔离级别
	err := gmongo.WithIsolationLevel(ctx, gmongo.IsolationLevelSnapshot, func(sessCtx mongo.SessionContext) error {
		coll := gmongo.Coll("users")
		var users []User
		return coll.FindAll(sessCtx, bson.M{}, &users)
	})
	if err != nil {
		log.Printf("Snapshot 隔离级别测试失败: %v", err)
	} else {
		fmt.Println("✓ Snapshot 隔离级别测试成功")
	}

	// Majority 隔离级别
	err = gmongo.WithIsolationLevel(ctx, gmongo.IsolationLevelMajority, func(sessCtx mongo.SessionContext) error {
		coll := gmongo.Coll("users")
		var users []User
		return coll.FindAll(sessCtx, bson.M{}, &users)
	})
	if err != nil {
		log.Printf("Majority 隔离级别测试失败: %v", err)
	} else {
		fmt.Println("✓ Majority 隔离级别测试成功")
	}

	// Local 隔离级别
	err = gmongo.WithIsolationLevel(ctx, gmongo.IsolationLevelLocal, func(sessCtx mongo.SessionContext) error {
		coll := gmongo.Coll("users")
		var users []User
		return coll.FindAll(sessCtx, bson.M{}, &users)
	})
	if err != nil {
		log.Printf("Local 隔离级别测试失败: %v", err)
	} else {
		fmt.Println("✓ Local 隔离级别测试成功")
	}
}

// testConcurrency 测试并发安全
func testConcurrency(ctx context.Context) {
	coll := gmongo.Coll("users")
	var wg sync.WaitGroup
	concurrency := 10

	// 并发插入
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := coll.InsertOne(ctx, User{
				Name:      fmt.Sprintf("并发用户%d", index),
				Email:     fmt.Sprintf("concurrent%d@example.com", index),
				Age:       20 + index,
				CreatedAt: time.Now(),
			})
			if err != nil {
				log.Printf("并发插入失败 [%d]: %v", index, err)
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("✓ 并发插入测试完成，插入 %d 条记录\n", concurrency)

	// 并发查询
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			var users []User
			err := coll.Query().Context(ctx).
				WhereGte("age", 20).
				Limit(5).
				Find(&users)
			if err != nil {
				log.Printf("并发查询失败 [%d]: %v", index, err)
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("✓ 并发查询测试完成，执行 %d 次查询\n", concurrency)

	// 并发更新
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			_, err := coll.UpdateMany(ctx,
				bson.M{"email": fmt.Sprintf("concurrent%d@example.com", index)},
				bson.M{"$set": bson.M{"age": 30 + index}})
			if err != nil {
				log.Printf("并发更新失败 [%d]: %v", index, err)
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("✓ 并发更新测试完成，更新 %d 条记录\n", concurrency)
}

// testPerformance 测试性能
func testPerformance(ctx context.Context) {
	coll := gmongo.Coll("performance_test")
	coll.DeleteMany(ctx, bson.M{})

	// 批量插入性能测试
	batchSize := 1000
	start := time.Now()

	var docs []interface{}
	for i := 0; i < batchSize; i++ {
		docs = append(docs, User{
			Name:      fmt.Sprintf("性能测试用户%d", i),
			Email:     fmt.Sprintf("perf%d@example.com", i),
			Age:       20 + (i % 50),
			Tags:      []string{"tag1", "tag2", "tag3"},
			CreatedAt: time.Now(),
		})
	}

	_, err := coll.InsertMany(ctx, docs)
	elapsed := time.Since(start)
	if err != nil {
		log.Printf("批量插入失败: %v", err)
	} else {
		fmt.Printf("✓ 批量插入 %d 条记录，耗时: %v\n", batchSize, elapsed)
	}

	// 查询性能测试
	start = time.Now()
	var users []User
	err = coll.Query().Context(ctx).
		WhereGte("age", 30).
		Limit(100).
		Find(&users)
	elapsed = time.Since(start)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("✓ 查询 %d 条记录，耗时: %v\n", len(users), elapsed)
	}

	// 聚合性能测试
	start = time.Now()
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"age": bson.M{"$gte": 30}}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$age",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
	}
	cursor, err := coll.Aggregate(ctx, pipeline)
	elapsed = time.Since(start)
	if err != nil {
		log.Printf("聚合失败: %v", err)
	} else {
		var results []bson.M
		cursor.All(ctx, &results)
		cursor.Close(ctx)
		fmt.Printf("✓ 聚合查询返回 %d 条记录，耗时: %v\n", len(results), elapsed)
	}

	// 清理测试数据
	coll.DeleteMany(ctx, bson.M{})
}

// testEdgeCases 测试边界条件
func testEdgeCases(ctx context.Context) {
	coll := gmongo.Coll("edge_cases")
	coll.DeleteMany(ctx, bson.M{})

	// 空值测试
	_, err := coll.InsertOne(ctx, User{
		Name:      "",
		Email:     "empty@example.com",
		Age:       0,
		Tags:      []string{},
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("插入空值失败: %v", err)
	} else {
		fmt.Println("✓ 插入空值成功")
	}

	// nil 地址测试
	_, err = coll.InsertOne(ctx, User{
		Name:      "无地址用户",
		Email:     "noaddress@example.com",
		Age:       25,
		Address:   nil,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("插入 nil 地址失败: %v", err)
	} else {
		fmt.Println("✓ 插入 nil 地址成功")
	}

	// 查询不存在的文档
	var user User
	result := coll.FindOne(ctx, bson.M{"email": "notexist@example.com"})
	err = result.Decode(&user)
	if err == mongo.ErrNoDocuments {
		fmt.Println("✓ 查询不存在的文档返回正确错误")
	} else if err != nil {
		log.Printf("查询失败: %v", err)
	}

	// 更新不存在的文档
	updateResult, err := coll.UpdateOne(ctx,
		bson.M{"email": "notexist@example.com"},
		bson.M{"$set": bson.M{"age": 100}})
	if err != nil {
		log.Printf("更新不存在的文档失败: %v", err)
	} else if updateResult.MatchedCount == 0 {
		fmt.Println("✓ 更新不存在的文档，MatchedCount = 0")
	}

	// 删除不存在的文档
	delResult, err := coll.DeleteOne(ctx, bson.M{"email": "notexist@example.com"})
	if err != nil {
		log.Printf("删除不存在的文档失败: %v", err)
	} else if delResult.DeletedCount == 0 {
		fmt.Println("✓ 删除不存在的文档，DeletedCount = 0")
	}

	// 极大值测试
	_, err = coll.InsertOne(ctx, User{
		Name:      "极大值用户",
		Email:     "maxvalue@example.com",
		Age:       2147483647, // int32 最大值
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("插入极大值失败: %v", err)
	} else {
		fmt.Println("✓ 插入极大值成功")
	}

	// 空查询条件
	var users []User
	err = coll.Query().Context(ctx).Limit(5).Find(&users)
	if err != nil {
		log.Printf("空查询条件失败: %v", err)
	} else {
		fmt.Printf("✓ 空查询条件成功，返回 %d 条记录\n", len(users))
	}
}

// testErrorHandling 测试错误处理
func testErrorHandling(ctx context.Context) {
	// 测试无效的 ObjectID
	_, err := gmongo.ObjectIDFromHex("invalid_id")
	if err != nil {
		fmt.Println("✓ 无效 ObjectID 返回错误")
	}

	// 测试连接到不存在的数据库（使用默认客户端）
	coll := gmongo.Coll("test_collection")

	// 测试插入无效数据
	_, err = coll.InsertOne(ctx, nil)
	if err != nil {
		fmt.Println("✓ 插入 nil 数据返回错误")
	}

	// 测试无效的查询条件
	var user User
	result := coll.FindOne(ctx, "invalid_filter")
	err = result.Decode(&user)
	if err != nil {
		fmt.Println("✓ 无效查询条件返回错误")
	}

	// 测试超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
	defer cancel()
	time.Sleep(10 * time.Millisecond) // 确保超时

	var users []User
	err = coll.Query().Context(timeoutCtx).Find(&users)
	if err != nil {
		fmt.Println("✓ 超时上下文返回错误")
	}

	// 测试重复键错误（需要唯一索引）
	coll2 := gmongo.Coll("unique_test")
	coll2.DeleteMany(ctx, bson.M{})

	// 创建唯一索引
	uniqueIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	coll2.Indexes().CreateMany(ctx, []mongo.IndexModel{uniqueIndex})

	// 插入第一条记录
	coll2.InsertOne(ctx, User{
		Name:      "唯一测试1",
		Email:     "unique@example.com",
		Age:       30,
		CreatedAt: time.Now(),
	})

	// 尝试插入重复记录
	_, err = coll2.InsertOne(ctx, User{
		Name:      "唯一测试2",
		Email:     "unique@example.com",
		Age:       31,
		CreatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println("✓ 重复键错误被正确捕获")
	}
}
