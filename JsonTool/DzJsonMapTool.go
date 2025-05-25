package JsonTool

import (
	"encoding/json"
	"github.com/PWND0U/dztool/StringTool"
	"maps"
	"slices"
)

type DzJsonMap map[string]interface{}

func NewDzJsonMap(jsonData []byte) DzJsonMap {
	djm := make(DzJsonMap)
	err := json.Unmarshal(jsonData, &djm)
	if err != nil {
		//log.Fatal(err)
		return nil
	}
	return djm
}

func ParseDzJsonMap(jsonData map[string]interface{}) DzJsonMap {
	return jsonData
}

func (djm DzJsonMap) getMap(key string) DzJsonMap {
	if djm == nil {
		return nil
	}
	if slices.Contains(slices.Collect(maps.Keys(djm)), key) {
		switch v := djm[key].(type) {
		case map[string]interface{}:
			return v
		default:
			return nil
		}
	}
	return nil
}

func (djm DzJsonMap) GetMap(key string) DzJsonMap {
	dzStrings := StringTool.NewDzString(key).Spilt(".")
	tempDjm := djm
	for _, dzString := range dzStrings {
		if !dzString.IsEmpty() {
			tempDjm = tempDjm.getMap(dzString.ToString())
			if tempDjm == nil {
				return tempDjm
			}
		}
	}
	return tempDjm
}

func (djm DzJsonMap) GetMapArray(key string) []DzJsonMap {
	return djm.GetArray(key).GetMapArray()
}

func (djm DzJsonMap) getArray(key string) DzJsonArray {
	if djm == nil {
		return make(DzJsonArray, 0)
	}
	if slices.Contains(slices.Collect(maps.Keys(djm)), key) {
		switch v := djm[key].(type) {
		case []interface{}:
			return v
		default:
			return make(DzJsonArray, 0)
		}
	}
	return make(DzJsonArray, 0)
}

func (djm DzJsonMap) GetArray(key string) DzJsonArray {
	dzString := StringTool.NewDzString(key)
	if dzString.IsContains(".") {
		dzStrings := dzString.Spilt(".")
		getMap := djm.GetMap(dzStrings[0 : len(dzStrings)-1].Join(".").ToString())
		if getMap != nil {
			return getMap.getArray(dzStrings[len(dzStrings)-1].ToString())
		}
	} else {
		return djm.getArray(key)
	}
	return make(DzJsonArray, 0)
}
func (djm DzJsonMap) getBool(key string) bool {
	if djm == nil {
		return false
	}
	if slices.Contains(slices.Collect(maps.Keys(djm)), key) {
		switch v := djm[key].(type) {
		case bool:
			return v
		default:
			return false
		}
	}
	return false
}
func (djm DzJsonMap) GetBool(key string) bool {
	dzString := StringTool.NewDzString(key)
	if dzString.IsContains(".") {
		dzStrings := StringTool.NewDzString(key).Spilt(".")
		getMap := djm.GetMap(dzStrings[0 : len(dzStrings)-1].Join(".").ToString())
		if getMap != nil {
			return getMap.getBool(dzStrings[len(dzStrings)-1].ToString())
		}
	} else {
		return djm.getBool(key)
	}
	return false
}

func (djm DzJsonMap) getNumber(key string) float64 {
	if djm == nil {
		return 0
	}
	if slices.Contains(slices.Collect(maps.Keys(djm)), key) {
		switch v := djm[key].(type) {
		case float64:
			return v
		default:
			return 0
		}
	}
	return 0
}

func (djm DzJsonMap) GetNumber(key string) float64 {
	dzString := StringTool.NewDzString(key)
	if dzString.IsContains(".") {
		dzStrings := StringTool.NewDzString(key).Spilt(".")
		getMap := djm.GetMap(dzStrings[0 : len(dzStrings)-1].Join(".").ToString())
		if getMap != nil {
			return getMap.getNumber(dzStrings[len(dzStrings)-1].ToString())
		}
	} else {
		return djm.getNumber(key)
	}
	return 0
}

func (djm DzJsonMap) getString(key string) string {
	if djm == nil {
		return "make(DzJsonArray, 0)"
	}
	if slices.Contains(slices.Collect(maps.Keys(djm)), key) {
		switch v := djm[key].(type) {
		case string:
			return v
		default:
			return ""
		}
	}
	return ""
}

func (djm DzJsonMap) GetString(key string) string {
	dzString := StringTool.NewDzString(key)
	if dzString.IsContains(".") {
		dzStrings := StringTool.NewDzString(key).Spilt(".")
		getMap := djm.GetMap(dzStrings[0 : len(dzStrings)-1].Join(".").ToString())
		if getMap != nil {
			return getMap.getString(dzStrings[len(dzStrings)-1].ToString())
		}
	} else {
		return djm.getString(key)
	}
	return ""
}
