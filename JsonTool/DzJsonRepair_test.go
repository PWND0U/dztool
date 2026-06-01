package JsonTool

import (
	"fmt"

	"github.com/bytedance/sonic"
	"testing"
)

// ==================== RepairJSON 测试 ====================

func TestRepairJSON_MissingQuotes(t *testing.T) {
	input := `{name: John, age: 30}`
	result, err := RepairJSON(input)
	fmt.Println("MissingQuotes result:", result)
	fmt.Println("MissingQuotes err:", err)
	fmt.Println("MissingQuotes valid:", sonic.Valid([]byte(result)))
}

func TestRepairJSON_SingleQuotes(t *testing.T) {
	input := `{'a': 'foo', 'b': 'bar'}`
	result, err := RepairJSON(input)
	fmt.Println("SingleQuotes result:", result)
	fmt.Println("SingleQuotes err:", err)
}

func TestRepairJSON_TrailingComma(t *testing.T) {
	input := `{"a":1,"b":2,}`
	result, err := RepairJSON(input)
	fmt.Println("TrailingComma result:", result)
	fmt.Println("TrailingComma err:", err)
}

func TestRepairJSON_MissingBracket(t *testing.T) {
	input := `{"items":[1,2,3`
	result, err := RepairJSON(input)
	fmt.Println("MissingBracket result:", result)
	fmt.Println("MissingBracket err:", err)
}

func TestRepairJSON_JavaScriptComments(t *testing.T) {
	input := `{"a": 1 /* comment */, "b": 2 // line comment\n}`
	result, err := RepairJSON(input)
	fmt.Println("JSComments result:", result)
	fmt.Println("JSComments err:", err)
}

func TestRepairJSON_PythonConstants(t *testing.T) {
	input := `{"active": True, "value": None, "flag": False}`
	result, err := RepairJSON(input)
	fmt.Println("PythonConstants result:", result)
	fmt.Println("PythonConstants err:", err)
}

func TestRepairJSON_JSONP(t *testing.T) {
	input := `callback({"status": "ok"});`
	result, err := RepairJSON(input)
	fmt.Println("JSONP result:", result)
	fmt.Println("JSONP err:", err)
}

func TestRepairJSON_Truncated(t *testing.T) {
	input := `{"message": "truncated string`
	result, err := RepairJSON(input)
	fmt.Println("Truncated result:", result)
	fmt.Println("Truncated err:", err)
}

func TestRepairJSON_Ellipsis(t *testing.T) {
	input := `[1, 2, 3, ...]`
	result, err := RepairJSON(input)
	fmt.Println("Ellipsis result:", result)
	fmt.Println("Ellipsis err:", err)
}

func TestRepairJSON_NDJSON(t *testing.T) {
	input := "{\"id\":1}\n{\"id\":2}\n{\"id\":3}"
	result, err := RepairJSON(input)
	fmt.Println("NDJSON result:", result)
	fmt.Println("NDJSON err:", err)
}

func TestRepairJSON_ValidJSON(t *testing.T) {
	input := `{"name": "Alice", "age": 30}`
	result, err := RepairJSON(input)
	fmt.Println("ValidJSON result:", result)
	fmt.Println("ValidJSON err:", err)
}

func TestRepairJSON_InvalidInput(t *testing.T) {
	input := `totally not json at all [[[`
	result, err := RepairJSON(input)
	fmt.Println("InvalidInput result:", result)
	fmt.Println("InvalidInput err:", err)
}

func TestRepairJSON_EmptyInput(t *testing.T) {
	result, err := RepairJSON("")
	fmt.Println("EmptyInput result:", result)
	fmt.Println("EmptyInput err:", err)
}

// ==================== MustRepairJSON 测试 ====================

func TestMustRepairJSON(t *testing.T) {
	input := `{name: MustRepair}`
	result := MustRepairJSON(input)
	fmt.Println("MustRepairJSON result:", result)
}

// ==================== RepairToDzJsonMap 测试 ====================

func TestRepairToDzJsonMap(t *testing.T) {
	input := `{name: "Alice", age: 25}`
	result := RepairToDzJsonMap(input)
	fmt.Println("RepairToDzJsonMap:", result)
	if result != nil {
		fmt.Println("RepairToDzJsonMap name:", result.GetString("name"))
		fmt.Println("RepairToDzJsonMap age:", result.GetInt("age"))
	}
}

func TestRepairToDzJsonMap_Invalid(t *testing.T) {
	input := `not json at all`
	result := RepairToDzJsonMap(input)
	fmt.Println("RepairToDzJsonMap invalid:", result)
}

// ==================== RepairToJsonMap 测试 ====================

func TestRepairToJsonMap(t *testing.T) {
	input := `{name: "Bob", active: True}`
	m, err := RepairToJsonMap(input)
	fmt.Println("RepairToJsonMap:", m)
	fmt.Println("RepairToJsonMap err:", err)
}

// ==================== RepairToJsonSlice 测试 ====================

func TestRepairToJsonSlice(t *testing.T) {
	input := `[1, 2, 3, ...]`
	s, err := RepairToJsonSlice(input)
	fmt.Println("RepairToJsonSlice:", s)
	fmt.Println("RepairToJsonSlice err:", err)
}

// ==================== RepairToStruct 测试 ====================

func TestRepairToStruct(t *testing.T) {
	type Person struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Active bool   `json:"active"`
	}

	input := `{name: "Charlie", age: 30, active: True}`
	var p Person
	err := RepairToStruct(input, &p)
	fmt.Println("RepairToStruct:", p)
	fmt.Println("RepairToStruct err:", err)
}

// ==================== IsRepairable 测试 ====================

func TestIsRepairable(t *testing.T) {
	fmt.Println("IsRepairable (valid):", IsRepairable(`{"a":1}`))
	fmt.Println("IsRepairable (fixable):", IsRepairable(`{a:1}`))
	fmt.Println("IsRepairable (broken):", IsRepairable(`{{{`))
}

// ==================== TryRepair 测试 ====================

func TestTryRepair_Valid(t *testing.T) {
	input := `{"name": "Alice"}`
	result, repaired, err := TryRepair(input)
	fmt.Println("TryRepair valid result:", result)
	fmt.Println("TryRepair valid repaired:", repaired)
	fmt.Println("TryRepair valid err:", err)
}

func TestTryRepair_NeedsRepair(t *testing.T) {
	input := `{name: "Alice"}`
	result, repaired, err := TryRepair(input)
	fmt.Println("TryRepair fixable result:", result)
	fmt.Println("TryRepair fixable repaired:", repaired)
	fmt.Println("TryRepair fixable err:", err)
}

func TestTryRepair_Unrepairable(t *testing.T) {
	input := `{{{`
	result, repaired, err := TryRepair(input)
	fmt.Println("TryRepair broken result:", result)
	fmt.Println("TryRepair broken repaired:", repaired)
	fmt.Println("TryRepair broken err:", err)
}
