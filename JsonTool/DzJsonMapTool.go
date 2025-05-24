package JsonTool

import (
	"encoding/json"
	"log"
)

type DzJsonMap map[string]interface{}

func NewDzJsonMap(jsonData []byte) DzJsonMap {
	djm := make(DzJsonMap)
	err := json.Unmarshal(jsonData, &djm)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return djm
}

func ParseDzJsonMap(jsonData map[string]interface{}) DzJsonMap {
	return jsonData
}

func (djm DzJsonMap) GetMap(key string) DzJsonMap {
	switch v := djm[key].(type) {
	case map[string]interface{}:
		return v
	default:
		return nil
	}
}

func (djm DzJsonMap) GetMapArray(key string) DzJsonMapArray {
	switch v := djm[key].(type) {
	case []map[string]interface{}:
		return v
	default:
		return nil
	}
}

func (djm DzJsonMap) GetArray(key string) DzJsonArray {
	switch v := djm[key].(type) {
	case []interface{}:
		return v
	default:
		return nil
	}
}

func (djm DzJsonMap) GetBool(key string) bool {
	switch v := djm[key].(type) {
	case bool:
		return v
	default:
		return false
	}
}

func (djm DzJsonMap) GetNumber(key string) float64 {
	switch v := djm[key].(type) {
	case float64:
		return v
	default:
		return 0
	}
}
