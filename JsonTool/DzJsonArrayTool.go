package JsonTool

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

func (dja DzJsonArray) GetNumberArray() []float64 {
	numberArray := make([]float64, 0)
	for _, ja := range dja {
		switch v := ja.(type) {
		case float64:
			numberArray = append(numberArray, v)
		}
	}
	return numberArray
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
