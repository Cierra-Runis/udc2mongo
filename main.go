package main

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"udc2mongo/database"
	"udc2mongo/model"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		log.Println("Continuing with system environment variables...")
	} else {
		fmt.Println("✓ .env file loaded successfully")
	}

	fmt.Println("Unicode Data to MongoDB Processor")
	fmt.Println("==================================")

	baseUrl := "https://www.unicode.org/Public/16.0.0/ucdxml/"

	// 获取XML内容（带缓存）
	fmt.Println("1. Fetching Unicode data...")
	content, err := fetchUcdXmlContentWithCache(baseUrl)
	if err != nil {
		fmt.Printf("Error fetching UCD XML content: %v\n", err)
		return
	}
	fmt.Printf("Successfully fetched %d bytes of XML data\n", len(content))

	// 解析XML
	fmt.Println("\n2. Parsing XML data...")
	ucd, err := model.ParseUCDXML(content)
	if err != nil {
		fmt.Printf("Error parsing XML: %v\n", err)
		return
	}

	// 处理数据
	fmt.Println("\n3. Processing data for MongoDB...")
	codePoints, blocks, err := model.ProcessUCDForMongoDB(ucd)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return
	}

	// 连接MongoDB
	fmt.Println("\n4. Connecting to MongoDB...")
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // 默认本地连接
	}

	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "unicode_db" // 默认数据库名
	}

	mongoClient, err := database.NewMongoClient(mongoURI, dbName)
	if err != nil {
		fmt.Printf("Error connecting to MongoDB: %v\n", err)
		return
	}
	defer mongoClient.Close()

	fmt.Printf("Connected to MongoDB at %s, database: %s\n", mongoURI, dbName)

	// 创建索引
	fmt.Println("\n5. Creating database indexes...")
	err = mongoClient.CreateIndexes()
	if err != nil {
		fmt.Printf("Error creating indexes: %v\n", err)
		return
	}

	// 保存数据到MongoDB
	fmt.Println("\n6. Saving data to MongoDB...")

	// 保存UCD主文档
	ucd.Version = "16.0.0"
	err = mongoClient.SaveUCD(ucd)
	if err != nil {
		fmt.Printf("Error saving UCD: %v\n", err)
		return
	}

	// 保存字符点
	err = mongoClient.SaveCodePoints(codePoints)
	if err != nil {
		fmt.Printf("Error saving code points: %v\n", err)
		return
	}

	// 保存块
	err = mongoClient.SaveBlocks(blocks)
	if err != nil {
		fmt.Printf("Error saving blocks: %v\n", err)
		return
	}

	// 获取统计信息
	fmt.Println("\n7. Database Statistics:")
	stats, err := mongoClient.GetStats()
	if err != nil {
		fmt.Printf("Error getting stats: %v\n", err)
		return
	}

	fmt.Printf("✓ Total Code Points: %d\n", stats.CodePointCount)
	fmt.Printf("✓ Total Blocks: %d\n", stats.BlockCount)
	fmt.Printf("✓ UCD Documents: %d\n", stats.UCDCount)

	// 详细字符类型统计
	fmt.Println("\n8. Detailed Character Type Analysis:")
	err = analyzeCharacterTypes(mongoClient)
	if err != nil {
		fmt.Printf("Error analyzing character types: %v\n", err)
		return
	}

	if len(stats.TopScripts) > 0 {
		fmt.Println("\nTop Scripts by Character Count:")
		for i, script := range stats.TopScripts {
			if i >= 5 { // 只显示前5个
				break
			}
			scriptName := script.Script
			if scriptName == "" {
				scriptName = "(No Script)"
			}
			fmt.Printf("  %s: %d characters\n", scriptName, script.Count)
		}
	}

	fmt.Println("\n✅ Data successfully imported to MongoDB!")
	fmt.Println("\nExample queries you can run:")
	fmt.Printf("  - Find character by code point: db.code_points.findOne({\"cp\": \"0041\"})\n")
	fmt.Printf("  - Find characters in Latin block: db.code_points.find({\"block\": \"ASCII\"})\n")
	fmt.Printf("  - Find Chinese characters: db.code_points.find({\"script\": \"Hani\"})\n")
}

const (
	cacheFile    = "ucd.all.flat.xml"
	cacheZipFile = "ucd.all.flat.zip"
)

// fetchUcdXmlContentWithCache 带缓存的数据获取函数
func fetchUcdXmlContentWithCache(baseUrl string) ([]byte, error) {
	cacheDir := os.TempDir()

	// 确保缓存目录存在
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	cacheFilePath := filepath.Join(cacheDir, cacheFile)
	cacheZipPath := filepath.Join(cacheDir, cacheZipFile)

	// 检查XML缓存是否存在
	if isCacheValid(cacheFilePath) {
		fmt.Println("Using cached XML data...")
		return os.ReadFile(cacheFilePath)
	}

	// 检查ZIP缓存是否存在，如果存在就解压
	if isCacheValid(cacheZipPath) {
		fmt.Println("Found cached ZIP file, extracting XML...")
		content, err := extractXmlFromZipFile(cacheZipPath)
		if err != nil {
			fmt.Printf("Failed to extract from cached ZIP: %v, downloading fresh copy...\n", err)
		} else {
			// 保存解压后的XML到缓存
			if err := os.WriteFile(cacheFilePath, content, 0644); err != nil {
				fmt.Printf("Warning: failed to save XML cache: %v\n", err)
			} else {
				fmt.Println("XML cached successfully.")
			}
			return content, nil
		}
	}

	fmt.Println("No cache found, downloading from network...")

	// 从网络获取数据
	content, err := fetchUcdXmlContent(baseUrl)
	if err != nil {
		return nil, err
	}

	// 保存XML到缓存
	if err := os.WriteFile(cacheFilePath, content, 0644); err != nil {
		fmt.Printf("Warning: failed to save XML cache: %v\n", err)
		// 即使缓存保存失败，也返回获取到的内容
	} else {
		fmt.Println("XML data cached successfully.")
	}

	return content, nil
}

// isCacheValid 检查缓存文件是否存在
func isCacheValid(cacheFilePath string) bool {
	_, err := os.Stat(cacheFilePath)
	return err == nil // 文件存在就返回true
}

// extractXmlFromZipFile 从本地ZIP文件中提取XML内容
func extractXmlFromZipFile(zipFilePath string) ([]byte, error) {
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name != "ucd.all.flat.xml" {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer rc.Close()

		content, err := io.ReadAll(rc)
		if err != nil {
			return nil, fmt.Errorf("failed to read XML content: %w", err)
		}

		return content, nil
	}

	return nil, fmt.Errorf("ucd.all.flat.xml not found in zip file")
}

func fetchUcdXmlContent(baseUrl string) ([]byte, error) {
	fileName := "ucd.all.flat.zip"

	filePath, err := url.JoinPath(baseUrl, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to construct file URL: %w", err)
	}

	resp, err := http.Get(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	var xmlContent []byte
	for _, file := range zipReader.File {
		if file.Name != "ucd.all.flat.xml" {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer rc.Close()

		xmlContent, err = io.ReadAll(rc)
		if err != nil {
			return nil, fmt.Errorf("failed to read XML content: %w", err)
		}
	}

	if len(xmlContent) == 0 {
		return nil, fmt.Errorf("ucd.all.flat.xml not found in zip file")
	}

	return xmlContent, nil
}

// analyzeCharacterTypes 分析字符类型统计
func analyzeCharacterTypes(mongoClient *database.MongoClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 总字符数
	total, err := mongoClient.CodePoints.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("error counting total: %w", err)
	}
	fmt.Printf("总字符数: %d\n", total)

	// 有名称的字符
	withNames, err := mongoClient.CodePoints.CountDocuments(ctx, bson.M{
		"name": bson.M{"$exists": true, "$ne": ""},
	})
	if err != nil {
		return fmt.Errorf("error counting with names: %w", err)
	}
	fmt.Printf("有名称的字符: %d\n", withNames)

	// 保留字符
	deprecated, err := mongoClient.CodePoints.CountDocuments(ctx, bson.M{
		"deprecated": true,
	})
	if err != nil {
		return fmt.Errorf("error counting deprecated: %w", err)
	}
	fmt.Printf("保留字符: %d\n", deprecated)

	// 非字符
	nonchar, err := mongoClient.CodePoints.CountDocuments(ctx, bson.M{
		"noncharacter": true,
	})
	if err != nil {
		return fmt.Errorf("error counting noncharacter: %w", err)
	}
	fmt.Printf("非字符: %d\n", nonchar)

	// 有CP字段的字符
	withCP, err := mongoClient.CodePoints.CountDocuments(ctx, bson.M{
		"cp": bson.M{"$exists": true, "$ne": ""},
	})
	if err != nil {
		return fmt.Errorf("error counting with CP: %w", err)
	}
	fmt.Printf("有CP字段的字符: %d\n", withCP)

	return nil
}
