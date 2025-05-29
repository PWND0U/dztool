package Algorithm

func DzLevenshtein(str1, str2 string, weightInsert, weightReplace, weightDelete int) int {
	tempStr1 := []rune(str1)
	tempStr2 := []rune(str2)
	str1Len := len(tempStr1)
	str2Len := len(tempStr2)
	d := make([][]int, str1Len+1)
	for i := 0; i <= str1Len; i++ {
		temp := make([]int, str2Len+1)
		temp[0] = i
		d[i] = temp
	}
	for i := 0; i <= str2Len; i++ {
		d[0][i] = i
	}
	for i := 1; i <= str1Len; i++ {
		for j := 1; j <= str2Len; j++ {
			if tempStr1[i-1] == tempStr2[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				// 这里需要看添加/删除/替换不同的操作
				// [i-1,j-1]->[i,j] 替换
				// [i,j-1]->[i,j] 插入
				// [i-1,j]->[i,j] 删除
				if d[i-1][j-1] <= d[i][j-1] && d[i-1][j-1] <= d[i-1][j] {
					// 表示是替换操作
					d[i][j] = weightReplace + d[i-1][j-1]
				} else if d[i][j-1] <= d[i-1][j-1] && d[i][j-1] <= d[i-1][j] {
					// 表示是插入操作
					d[i][j] = weightInsert + d[i][j-1]
				} else {
					//表示是删除
					d[i][j] = weightDelete + d[i-1][j]
				}
			}
		}

	}
	return d[str1Len][str2Len]
}
