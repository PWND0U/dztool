package JsonTool

import "github.com/spf13/cast"

type DzJsonArray []interface{}

func (dja DzJsonArray) GetStringArray() []string {
	strArray := make([]string, 0)
	for _, ja := range dja {
		switch v := ja.(type) {
		case string:
			strArray = append(strArray, v)
		}
	}
	return strArray
}

func (dja DzJsonArray) GetFloatArray() []float64 {
	numberArray := make([]float64, 0)
	for _, ja := range dja {
		numberArray = append(numberArray, cast.ToFloat64(ja))
		//switch v := ja.(type) {
		//case float64:
		//	numberArray = append(numberArray, v)
		//}
	}
	return numberArray
}

func (dja DzJsonArray) GetIntArray() []int {
	numberArray := make([]int, 0)
	for _, ja := range dja {
		numberArray = append(numberArray, cast.ToInt(ja))
		//switch v := ja.(type) {
		//case float64:
		//	numberArray = append(numberArray, v)
		//}
	}
	return numberArray
}

func (dja DzJsonArray) GetMapArray() []DzJsonMap {
	mapArray := make([]DzJsonMap, 0)
	for _, ja := range dja {
		switch v := ja.(type) {
		case map[string]interface{}:
			mapArray = append(mapArray, v)
		}
	}
	return mapArray
}

func (dja DzJsonArray) GetBoolArray() []bool {
	boolArray := make([]bool, 0)
	for _, ja := range dja {
		switch v := ja.(type) {
		case bool:
			boolArray = append(boolArray, v)
		}
	}
	return boolArray
}
