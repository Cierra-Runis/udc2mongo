package model

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// ParseUCDXML 解析UCD XML数据
func ParseUCDXML(xmlData []byte) (*UCD, error) {
	var ucd UCD

	// 解析XML
	err := xml.Unmarshal(xmlData, &ucd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	fmt.Printf("Parsed UCD with %d repertoire items\n",
		len(getCodePointsFromRepertoire(ucd.Repertoire)))

	if ucd.Blocks != nil {
		fmt.Printf("Found %d blocks\n", len(ucd.Blocks.Blocks))
	}

	return &ucd, nil
}

// getCodePointsFromRepertoire 从 repertoire 中提取所有字符点
func getCodePointsFromRepertoire(repertoire *Repertoire) []CodePoint {
	if repertoire == nil {
		return nil
	}

	var allCodePoints []CodePoint

	// 添加普通字符
	fmt.Printf("  - Regular characters: %d\n", len(repertoire.CodePoints))
	allCodePoints = append(allCodePoints, repertoire.CodePoints...)

	// 添加保留字符
	fmt.Printf("  - Reserved characters: %d\n", len(repertoire.Reserved))
	for _, cp := range repertoire.Reserved {
		cp.Deprecated = true // 标记为保留
		allCodePoints = append(allCodePoints, cp)
	}

	// 添加非字符
	fmt.Printf("  - Noncharacters: %d\n", len(repertoire.Noncharacter))
	for _, cp := range repertoire.Noncharacter {
		cp.Noncharacter = true
		allCodePoints = append(allCodePoints, cp)
	}

	// 添加代理对
	fmt.Printf("  - Surrogate characters: %d\n", len(repertoire.Surrogate))
	allCodePoints = append(allCodePoints, repertoire.Surrogate...)

	return allCodePoints
}

// ExtractAllCodePoints 从UCD中提取所有字符点用于单独存储
func ExtractAllCodePoints(ucd *UCD) []CodePoint {
	var allCodePoints []CodePoint

	// 从repertoire提取字符点
	if ucd.Repertoire != nil {
		allCodePoints = append(allCodePoints, getCodePointsFromRepertoire(ucd.Repertoire)...)
	}

	fmt.Printf("Extracted %d total code points\n", len(allCodePoints))
	return allCodePoints
}

// ExtractBlocks 从UCD中提取块信息
func ExtractBlocks(ucd *UCD) []Block {
	if ucd.Blocks == nil {
		return nil
	}

	return ucd.Blocks.Blocks
}

// ValidateCodePoint 验证字符点数据
func ValidateCodePoint(cp *CodePoint) error {
	if cp.CP == "" && cp.FirstCP == "" {
		return fmt.Errorf("code point must have either cp or first-cp")
	}

	// 如果有范围，检查last-cp
	if cp.FirstCP != "" && cp.LastCP == "" {
		return fmt.Errorf("code point with first-cp must also have last-cp")
	}

	return nil
}

// NormalizeCodePoint 标准化字符点数据
func NormalizeCodePoint(cp *CodePoint) {
	// 规范化字符串字段，移除多余空格
	cp.Name = strings.TrimSpace(cp.Name)
	cp.Name1 = strings.TrimSpace(cp.Name1)
	cp.Block = strings.TrimSpace(cp.Block)
	cp.Script = strings.TrimSpace(cp.Script)

	// 规范化空字符串为空
	if cp.BidiMirroringGlyph == "#" {
		cp.BidiMirroringGlyph = ""
	}
	if cp.DecompositionMapping == "#" {
		cp.DecompositionMapping = ""
	}
	if cp.SimpleUppercase == "#" {
		cp.SimpleUppercase = ""
	}
	if cp.SimpleLowercase == "#" {
		cp.SimpleLowercase = ""
	}
	if cp.SimpleTitlecase == "#" {
		cp.SimpleTitlecase = ""
	}
	if cp.UppercaseMapping == "#" {
		cp.UppercaseMapping = ""
	}
	if cp.LowercaseMapping == "#" {
		cp.LowercaseMapping = ""
	}
	if cp.TitlecaseMapping == "#" {
		cp.TitlecaseMapping = ""
	}
	if cp.SimpleCaseFolding == "#" {
		cp.SimpleCaseFolding = ""
	}
	if cp.CaseFolding == "#" {
		cp.CaseFolding = ""
	}
}

// ProcessUCDForMongoDB 处理UCD数据准备保存到MongoDB
func ProcessUCDForMongoDB(ucd *UCD) ([]CodePoint, []Block, error) {
	// 提取所有字符点
	codePoints := ExtractAllCodePoints(ucd)

	// 验证和标准化字符点
	validCodePoints := make([]CodePoint, 0, len(codePoints))
	for i := range codePoints {
		cp := &codePoints[i]

		// 验证
		if err := ValidateCodePoint(cp); err != nil {
			fmt.Printf("Warning: skipping invalid code point: %v\n", err)
			continue
		}

		// 标准化
		NormalizeCodePoint(cp)

		validCodePoints = append(validCodePoints, *cp)
	}

	// 提取块
	blocks := ExtractBlocks(ucd)

	fmt.Printf("Processed %d valid code points and %d blocks\n",
		len(validCodePoints), len(blocks))

	return validCodePoints, blocks, nil
}
