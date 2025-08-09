package database

import (
	"context"
	"fmt"
	"time"
	"udc2mongo/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client     *mongo.Client
	database   *mongo.Database
	ucd        *mongo.Collection
	CodePoints *mongo.Collection
	blocks     *mongo.Collection
}

func NewMongoClient(uri, dbName string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(dbName)

	return &MongoClient{
		client:     client,
		database:   database,
		ucd:        database.Collection("ucd"),
		CodePoints: database.Collection("code_points"),
		blocks:     database.Collection("blocks"),
	}, nil
}

func (mc *MongoClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return mc.client.Disconnect(ctx)
}

func (mc *MongoClient) SaveUCD(ucd *model.UCD) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ucdMetadata := &model.UCD{
		Xmlns:       ucd.Xmlns,
		Description: ucd.Description,
		Version:     ucd.Version,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	fmt.Println("Clearing existing UCD metadata...")
	err := mc.ucd.Drop(ctx)
	if err != nil {
		return fmt.Errorf("failed to clear existing UCD data: %w", err)
	}

	fmt.Println("Saving UCD metadata...")
	result, err := mc.ucd.InsertOne(ctx, ucdMetadata)
	if err != nil {
		return fmt.Errorf("failed to save UCD: %w", err)
	}

	fmt.Printf("UCD metadata saved with ID: %v\n", result.InsertedID)
	return nil
}

func (mc *MongoClient) SaveCodePoints(codePoints []model.CodePoint) error {
	if len(codePoints) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("Clearing existing code points...")
	_, err := mc.CodePoints.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to clear existing code points: %w", err)
	}

	now := time.Now()
	documents := make([]interface{}, len(codePoints))
	for i := range codePoints {
		codePoints[i].ID = primitive.NewObjectID()
		codePoints[i].CreatedAt = now
		codePoints[i].UpdatedAt = now
		documents[i] = codePoints[i]
	}

	fmt.Printf("Inserting %d code points...\n", len(documents))

	batchSize := 1000
	for i := 0; i < len(documents); i += batchSize {
		end := i + batchSize
		if end > len(documents) {
			end = len(documents)
		}

		batch := documents[i:end]
		_, err := mc.CodePoints.InsertMany(ctx, batch)
		if err != nil {
			return fmt.Errorf("failed to insert code points batch %d-%d: %w", i, end, err)
		}

		fmt.Printf("Inserted batch %d-%d\n", i, end)
	}

	fmt.Printf("Successfully saved %d code points\n", len(codePoints))
	return nil
}

func (mc *MongoClient) SaveBlocks(blocks []model.Block) error {
	if len(blocks) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Clearing existing blocks...")
	_, err := mc.blocks.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to clear existing blocks: %w", err)
	}

	now := time.Now()
	documents := make([]interface{}, len(blocks))
	for i := range blocks {
		blocks[i].ID = primitive.NewObjectID()
		blocks[i].CreatedAt = now
		blocks[i].UpdatedAt = now
		documents[i] = blocks[i]
	}

	fmt.Printf("Inserting %d blocks...\n", len(documents))
	_, err = mc.blocks.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert blocks: %w", err)
	}

	fmt.Printf("Successfully saved %d blocks\n", len(blocks))
	return nil
}

func (mc *MongoClient) CreateIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Creating indexes...")

	fmt.Println("Dropping existing indexes...")
	_, err := mc.CodePoints.Indexes().DropAll(ctx)
	if err != nil {
		if !isNamespaceNotFoundError(err) {
			return fmt.Errorf("failed to drop existing indexes: %w", err)
		}
		fmt.Println("Code points collection doesn't exist yet, skipping index drop")
	}

	_, err = mc.blocks.Indexes().DropAll(ctx)
	if err != nil {
		if !isNamespaceNotFoundError(err) {
			return fmt.Errorf("failed to drop existing block indexes: %w", err)
		}
		fmt.Println("Blocks collection doesn't exist yet, skipping index drop")
	}

	codePointIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "cp", Value: 1}},
			Options: options.Index().SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "block", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "general_category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "script", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "age", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "first_cp", Value: 1},
				{Key: "last_cp", Value: 1},
			},
		},
	}

	_, err = mc.CodePoints.Indexes().CreateMany(ctx, codePointIndexes)
	if err != nil {
		return fmt.Errorf("failed to create code points indexes: %w", err)
	}

	blockIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "first_cp", Value: 1},
				{Key: "last_cp", Value: 1},
			},
		},
	}

	_, err = mc.blocks.Indexes().CreateMany(ctx, blockIndexes)
	if err != nil {
		return fmt.Errorf("failed to create blocks indexes: %w", err)
	}

	fmt.Println("Indexes created successfully")
	return nil
}

func (mc *MongoClient) GetCodePointByCP(cp string) (*model.CodePoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var codePoint model.CodePoint
	err := mc.CodePoints.FindOne(ctx, bson.M{"cp": cp}).Decode(&codePoint)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find code point %s: %w", cp, err)
	}

	return &codePoint, nil
}

func (mc *MongoClient) GetCodePointsByBlock(blockName string) ([]model.CodePoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := mc.CodePoints.Find(ctx, bson.M{"block": blockName})
	if err != nil {
		return nil, fmt.Errorf("failed to find code points in block %s: %w", blockName, err)
	}
	defer cursor.Close(ctx)

	var codePoints []model.CodePoint
	err = cursor.All(ctx, &codePoints)
	if err != nil {
		return nil, fmt.Errorf("failed to decode code points: %w", err)
	}

	return codePoints, nil
}

func (mc *MongoClient) GetStats() (*DatabaseStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats := &DatabaseStats{}

	codePointCount, err := mc.CodePoints.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to count code points: %w", err)
	}
	stats.CodePointCount = codePointCount

	blockCount, err := mc.blocks.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to count blocks: %w", err)
	}
	stats.BlockCount = blockCount

	ucdCount, err := mc.ucd.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to count UCD documents: %w", err)
	}
	stats.UCDCount = ucdCount

	pipeline := []bson.M{
		{"$group": bson.M{
			"_id":   "$script",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	cursor, err := mc.CodePoints.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate by script: %w", err)
	}
	defer cursor.Close(ctx)

	var scriptStats []ScriptStat
	err = cursor.All(ctx, &scriptStats)
	if err != nil {
		return nil, fmt.Errorf("failed to decode script stats: %w", err)
	}
	stats.TopScripts = scriptStats

	return stats, nil
}

// DatabaseStats 数据库统计信息
type DatabaseStats struct {
	CodePointCount int64        `json:"code_point_count"`
	BlockCount     int64        `json:"block_count"`
	UCDCount       int64        `json:"ucd_count"`
	TopScripts     []ScriptStat `json:"top_scripts"`
}

// ScriptStat 脚本统计
type ScriptStat struct {
	Script string `bson:"_id" json:"script"`
	Count  int64  `bson:"count" json:"count"`
}

// isNamespaceNotFoundError 检查错误是否为命名空间未找到错误
func isNamespaceNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// 检查错误消息中是否包含 NamespaceNotFound 或相关的关键词
	errMsg := err.Error()
	return contains(errMsg, "NamespaceNotFound") || contains(errMsg, "ns not found")
}

// contains 检查字符串是否包含子字符串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && indexIgnoreCase(s, substr) >= 0))
}

// indexIgnoreCase 不区分大小写地查找子字符串位置
func indexIgnoreCase(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(s) < len(substr) {
		return -1
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := s[i+j]
			c2 := substr[j]
			if c1 >= 'A' && c1 <= 'Z' {
				c1 = c1 + 32 // 转换为小写
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 = c2 + 32 // 转换为小写
			}
			if c1 != c2 {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
