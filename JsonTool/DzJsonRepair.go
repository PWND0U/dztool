package JsonTool

import (
	"fmt"

	"github.com/bytedance/sonic"

	"github.com/kaptinlin/jsonrepair"
)

// ==================== JSON 修补工具函数 ====================

// RepairJSON 修补不合法的 JSON 字符串，返回修补后的合法 JSON 字符串。
// 可修复的常见问题包括：
//   - 缺少引号的键或值: {name: John} → {"name": "John"}
//   - 单引号替代双引号: {'a':'foo'} → {"a":"foo"}
//   - 尾部逗号: {"a":1,} → {"a":1}
//   - 缺少闭合括号: {"items":[1,2,3 → {"items":[1,2,3]}
//   - JavaScript 注释: {"a":1/*comment*/} → {"a":1}
//   - Python 常量: {active: True} → {"active":true}
//   - JSONP 包装: callback({"ok":1}); → {"ok":1}
//   - 截断的 JSON: {"msg":"hello → {"msg":"hello"}
//   - 省略号占位: [1,2,...] → [1,2]
//   - NDJSON 多行: {"id":1}\n{"id":2} → [{"id":1},{"id":2}]
func RepairJSON(text string) (string, error) {
	return jsonrepair.Repair(text)
}

// MustRepairJSON 修补不合法的 JSON 字符串，修补失败时 panic。
func MustRepairJSON(text string) string {
	result, err := jsonrepair.Repair(text)
	if err != nil {
		panic(fmt.Sprintf("JsonTool: repair failed: %v", err))
	}
	return result
}

// RepairToDzJsonMap 修补不合法的 JSON 字符串并解析为 DzJsonMap。
// 修补失败或解析失败时返回 nil。
func RepairToDzJsonMap(text string) DzJsonMap {
	repaired, err := jsonrepair.Repair(text)
	if err != nil {
		return nil
	}
	return NewDzJsonMap([]byte(repaired))
}

// RepairToDzJsonMapArray 修补不合法的 JSON 字符串并解析为 DzJsonMapArray。
// 修补失败或解析失败时返回空切片。
func RepairToDzJsonMapArray(text string) DzJsonMapArray {
	repaired, err := jsonrepair.Repair(text)
	if err != nil {
		return make(DzJsonMapArray, 0)
	}
	return NewDzJsonMapArray([]byte(repaired))
}

// RepairToJsonMap 修补不合法的 JSON 字符串并解析为 map[string]interface{}。
func RepairToJsonMap(text string) (map[string]interface{}, error) {
	repaired, err := jsonrepair.Repair(text)
	if err != nil {
		return nil, fmt.Errorf("JsonTool: repair failed: %w", err)
	}
	var result map[string]interface{}
	if err := sonic.Unmarshal([]byte(repaired), &result); err != nil {
		return nil, fmt.Errorf("JsonTool: unmarshal repaired json: %w", err)
	}
	return result, nil
}

// RepairToJsonSlice 修补不合法的 JSON 字符串并解析为 []interface{}。
func RepairToJsonSlice(text string) ([]interface{}, error) {
	repaired, err := jsonrepair.Repair(text)
	if err != nil {
		return nil, fmt.Errorf("JsonTool: repair failed: %w", err)
	}
	var result []interface{}
	if err := sonic.Unmarshal([]byte(repaired), &result); err != nil {
		return nil, fmt.Errorf("JsonTool: unmarshal repaired json: %w", err)
	}
	return result, nil
}

// RepairToStruct 修补不合法的 JSON 字符串并解析到目标结构体。
// target 必须为结构体指针。
func RepairToStruct(text string, target interface{}) error {
	repaired, err := jsonrepair.Repair(text)
	if err != nil {
		return fmt.Errorf("JsonTool: repair failed: %w", err)
	}
	if err := sonic.Unmarshal([]byte(repaired), target); err != nil {
		return fmt.Errorf("JsonTool: unmarshal repaired json: %w", err)
	}
	return nil
}

// IsRepairable 判断给定的字符串是否可以通过 jsonrepair 修补。
// 如果输入本身是合法 JSON，也返回 true（无需修补即可通过）。
func IsRepairable(text string) bool {
	_, err := jsonrepair.Repair(text)
	return err == nil
}

// TryRepair 尝试修补 JSON 字符串，返回修补结果和是否进行了修补。
// 如果输入本身就是合法 JSON，返回原文且 repaired=false。
func TryRepair(text string) (result string, repaired bool, err error) {
	// 先检查是否本身合法
	if sonic.Valid([]byte(text)) {
		return text, false, nil
	}
	// 尝试修补
	fixed, err := jsonrepair.Repair(text)
	if err != nil {
		return "", false, err
	}
	return fixed, true, nil
}
