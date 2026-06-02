package dzUtils

import (
	"fmt"
	"math/big"
	"testing"
)

// ==================== 十进制转其他进制 ====================

func TestDecToBin(t *testing.T) {
	tests := []struct {
		input  int64
		expect string
	}{
		{0, "0"},
		{1, "1"},
		{10, "1010"},
		{255, "11111111"},
		{256, "100000000"},
		{-1, "-1"},
		{-255, "-11111111"},
	}
	for _, tt := range tests {
		got := DecToBin(tt.input)
		if got != tt.expect {
			t.Errorf("DecToBin(%d) = %q, 期望 %q", tt.input, got, tt.expect)
		}
	}
	fmt.Println("DecToBin 测试通过")
}

func TestDecToOct(t *testing.T) {
	tests := []struct {
		input  int64
		expect string
	}{
		{0, "0"},
		{8, "10"},
		{255, "377"},
		{1024, "2000"},
		{-8, "-10"},
	}
	for _, tt := range tests {
		got := DecToOct(tt.input)
		if got != tt.expect {
			t.Errorf("DecToOct(%d) = %q, 期望 %q", tt.input, got, tt.expect)
		}
	}
	fmt.Println("DecToOct 测试通过")
}

func TestDecToHex(t *testing.T) {
	tests := []struct {
		input  int64
		expect string
	}{
		{0, "0"},
		{16, "10"},
		{255, "ff"},
		{4096, "1000"},
		{-1, "-1"},
		{0xdeadbeef, "deadbeef"},
	}
	for _, tt := range tests {
		got := DecToHex(tt.input)
		if got != tt.expect {
			t.Errorf("DecToHex(%d) = %q, 期望 %q", tt.input, got, tt.expect)
		}
	}
	fmt.Println("DecToHex 测试通过")
}

// ==================== 其他进制转十进制 ====================

func TestBinToDec(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
		hasErr bool
	}{
		{"0", 0, false},
		{"1", 1, false},
		{"1010", 10, false},
		{"11111111", 255, false},
		{"0b1010", 10, false},
		{"0b11111111", 255, false},
		{"2", 0, true}, // 无效二进制字符
		{"", 0, true},  // 空字符串
	}
	for _, tt := range tests {
		got, err := BinToDec(tt.input)
		if tt.hasErr {
			if err == nil {
				t.Errorf("BinToDec(%q) 期望返回错误，实际未返回", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("BinToDec(%q) 未期望错误: %v", tt.input, err)
			}
			if got != tt.expect {
				t.Errorf("BinToDec(%q) = %d, 期望 %d", tt.input, got, tt.expect)
			}
		}
	}
	fmt.Println("BinToDec 测试通过")
}

func TestOctToDec(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
		hasErr bool
	}{
		{"0", 0, false},
		{"10", 8, false},
		{"377", 255, false},
		{"0o377", 255, false},
		{"0377", 255, false},
		{"9", 0, true}, // 无效八进制字符
	}
	for _, tt := range tests {
		got, err := OctToDec(tt.input)
		if tt.hasErr {
			if err == nil {
				t.Errorf("OctToDec(%q) 期望返回错误，实际未返回", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("OctToDec(%q) 未期望错误: %v", tt.input, err)
			}
			if got != tt.expect {
				t.Errorf("OctToDec(%q) = %d, 期望 %d", tt.input, got, tt.expect)
			}
		}
	}
	fmt.Println("OctToDec 测试通过")
}

func TestHexToDec(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
		hasErr bool
	}{
		{"0", 0, false},
		{"10", 16, false},
		{"ff", 255, false},
		{"FF", 255, false},
		{"0xff", 255, false},
		{"0XFF", 255, false},
		{"deadbeef", 3735928559, false},
		{"gg", 0, true}, // 无效十六进制字符
	}
	for _, tt := range tests {
		got, err := HexToDec(tt.input)
		if tt.hasErr {
			if err == nil {
				t.Errorf("HexToDec(%q) 期望返回错误，实际未返回", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("HexToDec(%q) 未期望错误: %v", tt.input, err)
			}
			if got != tt.expect {
				t.Errorf("HexToDec(%q) = %d, 期望 %d", tt.input, got, tt.expect)
			}
		}
	}
	fmt.Println("HexToDec 测试通过")
}

// ==================== 非十进制之间的交叉转换 ====================

func TestBinToOct(t *testing.T) {
	got, err := BinToOct("11111111")
	if err != nil {
		t.Fatalf("BinToOct 失败: %v", err)
	}
	if got != "377" {
		t.Errorf("BinToOct(\"11111111\") = %q, 期望 \"377\"", got)
	}

	got, err = BinToOct("0b1010")
	if err != nil {
		t.Fatalf("BinToOct 失败: %v", err)
	}
	if got != "12" {
		t.Errorf("BinToOct(\"0b1010\") = %q, 期望 \"12\"", got)
	}

	_, err = BinToOct("2")
	if err == nil {
		t.Error("BinToOct(\"2\") 期望返回错误")
	}
	fmt.Println("BinToOct 测试通过")
}

func TestBinToHex(t *testing.T) {
	got, err := BinToHex("11111111")
	if err != nil {
		t.Fatalf("BinToHex 失败: %v", err)
	}
	if got != "ff" {
		t.Errorf("BinToHex(\"11111111\") = %q, 期望 \"ff\"", got)
	}

	got, err = BinToHex("0b1010")
	if err != nil {
		t.Fatalf("BinToHex 失败: %v", err)
	}
	if got != "a" {
		t.Errorf("BinToHex(\"0b1010\") = %q, 期望 \"a\"", got)
	}
	fmt.Println("BinToHex 测试通过")
}

func TestOctToBin(t *testing.T) {
	got, err := OctToBin("377")
	if err != nil {
		t.Fatalf("OctToBin 失败: %v", err)
	}
	if got != "11111111" {
		t.Errorf("OctToBin(\"377\") = %q, 期望 \"11111111\"", got)
	}

	got, err = OctToBin("0o12")
	if err != nil {
		t.Fatalf("OctToBin 失败: %v", err)
	}
	if got != "1010" {
		t.Errorf("OctToBin(\"0o12\") = %q, 期望 \"1010\"", got)
	}
	fmt.Println("OctToBin 测试通过")
}

func TestOctToHex(t *testing.T) {
	got, err := OctToHex("377")
	if err != nil {
		t.Fatalf("OctToHex 失败: %v", err)
	}
	if got != "ff" {
		t.Errorf("OctToHex(\"377\") = %q, 期望 \"ff\"", got)
	}
	fmt.Println("OctToHex 测试通过")
}

func TestHexToBin(t *testing.T) {
	got, err := HexToBin("ff")
	if err != nil {
		t.Fatalf("HexToBin 失败: %v", err)
	}
	if got != "11111111" {
		t.Errorf("HexToBin(\"ff\") = %q, 期望 \"11111111\"", got)
	}

	got, err = HexToBin("0xFF")
	if err != nil {
		t.Fatalf("HexToBin 失败: %v", err)
	}
	if got != "11111111" {
		t.Errorf("HexToBin(\"0xFF\") = %q, 期望 \"11111111\"", got)
	}
	fmt.Println("HexToBin 测试通过")
}

func TestHexToOct(t *testing.T) {
	got, err := HexToOct("ff")
	if err != nil {
		t.Fatalf("HexToOct 失败: %v", err)
	}
	if got != "377" {
		t.Errorf("HexToOct(\"ff\") = %q, 期望 \"377\"", got)
	}
	fmt.Println("HexToOct 测试通过")
}

// ==================== 往返转换测试 ====================

func TestRoundTripDecBin(t *testing.T) {
	values := []int64{0, 1, 42, 127, 255, 1024, 65535, 2147483647}
	for _, v := range values {
		bin := DecToBin(v)
		got, err := BinToDec(bin)
		if err != nil {
			t.Errorf("BinToDec(%q) 失败: %v", bin, err)
		}
		if got != v {
			t.Errorf("往返转换: %d → %q → %d, 不一致", v, bin, got)
		}
	}
	fmt.Println("往返转换 Dec↔Bin 测试通过")
}

func TestRoundTripDecOct(t *testing.T) {
	values := []int64{0, 1, 42, 127, 255, 1024, 65535}
	for _, v := range values {
		oct := DecToOct(v)
		got, err := OctToDec(oct)
		if err != nil {
			t.Errorf("OctToDec(%q) 失败: %v", oct, err)
		}
		if got != v {
			t.Errorf("往返转换: %d → %q → %d, 不一致", v, oct, got)
		}
	}
	fmt.Println("往返转换 Dec↔Oct 测试通过")
}

func TestRoundTripDecHex(t *testing.T) {
	values := []int64{0, 1, 42, 127, 255, 1024, 65535, 3735928559}
	for _, v := range values {
		hex := DecToHex(v)
		got, err := HexToDec(hex)
		if err != nil {
			t.Errorf("HexToDec(%q) 失败: %v", hex, err)
		}
		if got != v {
			t.Errorf("往返转换: %d → %q → %d, 不一致", v, hex, got)
		}
	}
	fmt.Println("往返转换 Dec↔Hex 测试通过")
}

// ==================== 大整数进制转换 ====================

func TestBigDecToBin(t *testing.T) {
	n, ok := new(big.Int).SetString("18446744073709551615", 10) // uint64 max
	if !ok {
		t.Fatal("big.Int.SetString 失败")
	}
	got := BigDecToBin(n)
	expected := "1111111111111111111111111111111111111111111111111111111111111111"
	if got != expected {
		t.Errorf("BigDecToBin = %q, 期望 %q", got, expected)
	}
	fmt.Println("BigDecToBin 测试通过")
}

func TestBigDecToOct(t *testing.T) {
	n := new(big.Int)
	n.SetString("18446744073709551615", 10)
	got := BigDecToOct(n)
	if got != "1777777777777777777777" {
		t.Errorf("BigDecToOct = %q, 期望 \"1777777777777777777777\"", got)
	}
	fmt.Println("BigDecToOct 测试通过")
}

func TestBigDecToHex(t *testing.T) {
	n := new(big.Int)
	n.SetString("18446744073709551615", 10)
	got := BigDecToHex(n)
	if got != "ffffffffffffffff" {
		t.Errorf("BigDecToHex = %q, 期望 \"ffffffffffffffff\"", got)
	}
	fmt.Println("BigDecToHex 测试通过")
}

func TestBigBinToDec(t *testing.T) {
	input := "1111111111111111111111111111111111111111111111111111111111111111"
	got, err := BigBinToDec(input)
	if err != nil {
		t.Fatalf("BigBinToDec 失败: %v", err)
	}
	expected := new(big.Int)
	expected.SetString("18446744073709551615", 10)
	if got.Cmp(expected) != 0 {
		t.Errorf("BigBinToDec = %s, 期望 %s", got.String(), expected.String())
	}

	// 带 0b 前缀
	got2, err := BigBinToDec("0b" + input)
	if err != nil {
		t.Fatalf("BigBinToDec(带0b前缀) 失败: %v", err)
	}
	if got2.Cmp(expected) != 0 {
		t.Errorf("BigBinToDec(带0b前缀) = %s, 期望 %s", got2.String(), expected.String())
	}

	// 无效输入
	_, err = BigBinToDec("102")
	if err == nil {
		t.Error("BigBinToDec(\"102\") 期望返回错误")
	}
	fmt.Println("BigBinToDec 测试通过")
}

func TestBigOctToDec(t *testing.T) {
	got, err := BigOctToDec("1777777777777777777777")
	if err != nil {
		t.Fatalf("BigOctToDec 失败: %v", err)
	}
	expected := new(big.Int)
	expected.SetString("18446744073709551615", 10)
	if got.Cmp(expected) != 0 {
		t.Errorf("BigOctToDec = %s, 期望 %s", got.String(), expected.String())
	}
	fmt.Println("BigOctToDec 测试通过")
}

func TestBigHexToDec(t *testing.T) {
	got, err := BigHexToDec("ffffffffffffffff")
	if err != nil {
		t.Fatalf("BigHexToDec 失败: %v", err)
	}
	expected := new(big.Int)
	expected.SetString("18446744073709551615", 10)
	if got.Cmp(expected) != 0 {
		t.Errorf("BigHexToDec = %s, 期望 %s", got.String(), expected.String())
	}

	// 带 0x 前缀
	got2, err := BigHexToDec("0xffffffffffffffff")
	if err != nil {
		t.Fatalf("BigHexToDec(带0x前缀) 失败: %v", err)
	}
	if got2.Cmp(expected) != 0 {
		t.Errorf("BigHexToDec(带0x前缀) = %s, 期望 %s", got2.String(), expected.String())
	}
	fmt.Println("BigHexToDec 测试通过")
}

func TestBigRoundTrip(t *testing.T) {
	bigValues := []string{
		"0",
		"1",
		"18446744073709551615",               // uint64 max
		"340282366920938463463374607431768211455", // 2^128 - 1
	}
	for _, dec := range bigValues {
		n := new(big.Int)
		n.SetString(dec, 10)

		bin := BigDecToBin(n)
		back, err := BigBinToDec(bin)
		if err != nil || back.Cmp(n) != 0 {
			t.Errorf("大整数往返 Bin 失败: %s → %q → %v (err=%v)", dec, bin, back, err)
		}

		oct := BigDecToOct(n)
		back, err = BigOctToDec(oct)
		if err != nil || back.Cmp(n) != 0 {
			t.Errorf("大整数往返 Oct 失败: %s → %q → %v (err=%v)", dec, oct, back, err)
		}

		hex := BigDecToHex(n)
		back, err = BigHexToDec(hex)
		if err != nil || back.Cmp(n) != 0 {
			t.Errorf("大整数往返 Hex 失败: %s → %q → %v (err=%v)", dec, hex, back, err)
		}
	}
	fmt.Println("大整数往返转换测试通过")
}

// ==================== 字符串与进制转换 ====================

func TestStrToBin(t *testing.T) {
	got := StrToBin("AB")
	if got != "01000001 01000010" {
		t.Errorf("StrToBin(\"AB\") = %q, 期望 \"01000001 01000010\"", got)
	}

	got = StrToBin("")
	if got != "" {
		t.Errorf("StrToBin(\"\") = %q, 期望 \"\"", got)
	}
	fmt.Println("StrToBin 测试通过")
}

func TestStrToOct(t *testing.T) {
	got := StrToOct("AB")
	if got != "101 102" {
		t.Errorf("StrToOct(\"AB\") = %q, 期望 \"101 102\"", got)
	}
	fmt.Println("StrToOct 测试通过")
}

func TestStrToHex(t *testing.T) {
	got := StrToHex("AB")
	if got != "4142" {
		t.Errorf("StrToHex(\"AB\") = %q, 期望 \"4142\"", got)
	}

	got = StrToHex("")
	if got != "" {
		t.Errorf("StrToHex(\"\") = %q, 期望 \"\"", got)
	}
	fmt.Println("StrToHex 测试通过")
}

func TestBinToStr(t *testing.T) {
	got, err := BinToStr("0100000101000010")
	if err != nil {
		t.Fatalf("BinToStr 失败: %v", err)
	}
	if got != "AB" {
		t.Errorf("BinToStr = %q, 期望 \"AB\"", got)
	}

	// 带空格
	got, err = BinToStr("01000001 01000010")
	if err != nil {
		t.Fatalf("BinToStr(带空格) 失败: %v", err)
	}
	if got != "AB" {
		t.Errorf("BinToStr(带空格) = %q, 期望 \"AB\"", got)
	}

	// 长度不是8的倍数
	_, err = BinToStr("010")
	if err == nil {
		t.Error("BinToStr(\"010\") 期望返回错误")
	}
	fmt.Println("BinToStr 测试通过")
}

func TestOctToStr(t *testing.T) {
	got, err := OctToStr("101102")
	if err != nil {
		t.Fatalf("OctToStr 失败: %v", err)
	}
	if got != "AB" {
		t.Errorf("OctToStr = %q, 期望 \"AB\"", got)
	}

	// 带空格
	got, err = OctToStr("101 102")
	if err != nil {
		t.Fatalf("OctToStr(带空格) 失败: %v", err)
	}
	if got != "AB" {
		t.Errorf("OctToStr(带空格) = %q, 期望 \"AB\"", got)
	}

	// 长度不是3的倍数
	_, err = OctToStr("10")
	if err == nil {
		t.Error("OctToStr(\"10\") 期望返回错误")
	}
	fmt.Println("OctToStr 测试通过")
}

func TestHexToStr(t *testing.T) {
	got, err := HexToStr("4142")
	if err != nil {
		t.Fatalf("HexToStr 失败: %v", err)
	}
	if got != "AB" {
		t.Errorf("HexToStr = %q, 期望 \"AB\"", got)
	}

	// 大写
	got, err = HexToStr("4142")
	if err != nil {
		t.Fatalf("HexToStr(大写) 失败: %v", err)
	}
	if got != "AB" {
		t.Errorf("HexToStr(大写) = %q, 期望 \"AB\"", got)
	}

	// 带 0x 前缀
	got, err = HexToStr("0x4142")
	if err != nil {
		t.Fatalf("HexToStr(带0x) 失败: %v", err)
	}
	if got != "AB" {
		t.Errorf("HexToStr(带0x) = %q, 期望 \"AB\"", got)
	}

	// 长度不是2的倍数
	_, err = HexToStr("414")
	if err == nil {
		t.Error("HexToStr(\"414\") 期望返回错误")
	}
	fmt.Println("HexToStr 测试通过")
}

func TestStrRoundTrip(t *testing.T) {
	inputs := []string{"Hello", "AB", "测试", "123!@#"}
	for _, s := range inputs {
		// Bin 往返
		bin := StrToBin(s)
		back, err := BinToStr(bin)
		if err != nil || back != s {
			t.Errorf("字符串 Bin 往返失败: %q → %q → %q (err=%v)", s, bin, back, err)
		}

		// Oct 往返
		oct := StrToOct(s)
		back, err = OctToStr(oct)
		if err != nil || back != s {
			t.Errorf("字符串 Oct 往返失败: %q → %q → %q (err=%v)", s, oct, back, err)
		}

		// Hex 往返
		hex := StrToHex(s)
		back, err = HexToStr(hex)
		if err != nil || back != s {
			t.Errorf("字符串 Hex 往返失败: %q → %q → %q (err=%v)", s, hex, back, err)
		}
	}
	fmt.Println("字符串往返转换测试通过")
}

// ==================== 字节数组与进制转换 ====================

func TestBytesToBin(t *testing.T) {
	got := BytesToBin([]byte{0x41, 0x42})
	if got != "01000001 01000010" {
		t.Errorf("BytesToBin = %q, 期望 \"01000001 01000010\"", got)
	}

	got = BytesToBin([]byte{})
	if got != "" {
		t.Errorf("BytesToBin(nil) = %q, 期望 \"\"", got)
	}
	fmt.Println("BytesToBin 测试通过")
}

func TestBytesToOct(t *testing.T) {
	got := BytesToOct([]byte{0x41, 0x42})
	if got != "101 102" {
		t.Errorf("BytesToOct = %q, 期望 \"101 102\"", got)
	}
	fmt.Println("BytesToOct 测试通过")
}

func TestBytesToHex(t *testing.T) {
	got := BytesToHex([]byte{0x41, 0x42})
	if got != "4142" {
		t.Errorf("BytesToHex = %q, 期望 \"4142\"", got)
	}

	got = BytesToHex([]byte{0xde, 0xad, 0xbe, 0xef})
	if got != "deadbeef" {
		t.Errorf("BytesToHex = %q, 期望 \"deadbeef\"", got)
	}
	fmt.Println("BytesToHex 测试通过")
}

func TestBinToBytes(t *testing.T) {
	got, err := BinToBytes("0100000101000010")
	if err != nil {
		t.Fatalf("BinToBytes 失败: %v", err)
	}
	if string(got) != "AB" {
		t.Errorf("BinToBytes = %q, 期望 \"AB\"", string(got))
	}

	// 带空格
	got, err = BinToBytes("01000001 01000010")
	if err != nil {
		t.Fatalf("BinToBytes(带空格) 失败: %v", err)
	}
	if string(got) != "AB" {
		t.Errorf("BinToBytes(带空格) = %q, 期望 \"AB\"", string(got))
	}

	// 长度不是8的倍数
	_, err = BinToBytes("010")
	if err == nil {
		t.Error("BinToBytes(\"010\") 期望返回错误")
	}
	fmt.Println("BinToBytes 测试通过")
}

func TestOctToBytes(t *testing.T) {
	got, err := OctToBytes("101102")
	if err != nil {
		t.Fatalf("OctToBytes 失败: %v", err)
	}
	if string(got) != "AB" {
		t.Errorf("OctToBytes = %q, 期望 \"AB\"", string(got))
	}

	// 带空格
	got, err = OctToBytes("101 102")
	if err != nil {
		t.Fatalf("OctToBytes(带空格) 失败: %v", err)
	}
	if string(got) != "AB" {
		t.Errorf("OctToBytes(带空格) = %q, 期望 \"AB\"", string(got))
	}

	// 长度不是3的倍数
	_, err = OctToBytes("10")
	if err == nil {
		t.Error("OctToBytes(\"10\") 期望返回错误")
	}
	fmt.Println("OctToBytes 测试通过")
}

func TestHexToBytes(t *testing.T) {
	got, err := HexToBytes("4142")
	if err != nil {
		t.Fatalf("HexToBytes 失败: %v", err)
	}
	if string(got) != "AB" {
		t.Errorf("HexToBytes = %q, 期望 \"AB\"", string(got))
	}

	// 大写
	got, err = HexToBytes("4142")
	if err != nil {
		t.Fatalf("HexToBytes(大写) 失败: %v", err)
	}
	if string(got) != "AB" {
		t.Errorf("HexToBytes(大写) = %q, 期望 \"AB\"", string(got))
	}

	// 带 0x 前缀
	got, err = HexToBytes("0x4142")
	if err != nil {
		t.Fatalf("HexToBytes(带0x) 失败: %v", err)
	}
	if string(got) != "AB" {
		t.Errorf("HexToBytes(带0x) = %q, 期望 \"AB\"", string(got))
	}

	// 长度不是2的倍数
	_, err = HexToBytes("414")
	if err == nil {
		t.Error("HexToBytes(\"414\") 期望返回错误")
	}

	// 无效字符
	_, err = HexToBytes("gg")
	if err == nil {
		t.Error("HexToBytes(\"gg\") 期望返回错误")
	}
	fmt.Println("HexToBytes 测试通过")
}

func TestBytesRoundTrip(t *testing.T) {
	inputs := [][]byte{
		{},
		{0x00},
		{0x41, 0x42},
		{0xde, 0xad, 0xbe, 0xef},
		{0xff, 0x00, 0x80},
	}
	for _, data := range inputs {
		// Bin 往返
		bin := BytesToBin(data)
		back, err := BinToBytes(bin)
		if err != nil {
			t.Errorf("字节 Bin 往返失败: %v → %q: %v", data, bin, err)
			continue
		}
		if !bytesEqual(back, data) {
			t.Errorf("字节 Bin 往返不一致: %v → %q → %v", data, bin, back)
		}

		// Oct 往返
		oct := BytesToOct(data)
		back, err = OctToBytes(oct)
		if err != nil {
			t.Errorf("字节 Oct 往返失败: %v → %q: %v", data, oct, err)
			continue
		}
		if !bytesEqual(back, data) {
			t.Errorf("字节 Oct 往返不一致: %v → %q → %v", data, oct, back)
		}

		// Hex 往返
		hex := BytesToHex(data)
		back, err = HexToBytes(hex)
		if err != nil {
			t.Errorf("字节 Hex 往返失败: %v → %q: %v", data, hex, err)
			continue
		}
		if !bytesEqual(back, data) {
			t.Errorf("字节 Hex 往返不一致: %v → %q → %v", data, hex, back)
		}
	}
	fmt.Println("字节数组往返转换测试通过")
}

// ==================== 通用进制转换 ====================

func TestConvertBase(t *testing.T) {
	// 二进制 → 十六进制
	got, err := ConvertBase("11111111", 2, 16)
	if err != nil {
		t.Fatalf("ConvertBase 失败: %v", err)
	}
	if got != "ff" {
		t.Errorf("ConvertBase(\"11111111\", 2, 16) = %q, 期望 \"ff\"", got)
	}

	// 十六进制 → 十进制
	got, err = ConvertBase("ff", 16, 10)
	if err != nil {
		t.Fatalf("ConvertBase 失败: %v", err)
	}
	if got != "255" {
		t.Errorf("ConvertBase(\"ff\", 16, 10) = %q, 期望 \"255\"", got)
	}

	// 八进制 → 二进制
	got, err = ConvertBase("377", 8, 2)
	if err != nil {
		t.Fatalf("ConvertBase 失败: %v", err)
	}
	if got != "11111111" {
		t.Errorf("ConvertBase(\"377\", 8, 2) = %q, 期望 \"11111111\"", got)
	}

	// 无效进制
	_, err = ConvertBase("10", 1, 10)
	if err == nil {
		t.Error("ConvertBase 进制=1 期望返回错误")
	}
	_, err = ConvertBase("10", 2, 37)
	if err == nil {
		t.Error("ConvertBase 进制=37 期望返回错误")
	}

	// 无效输入
	_, err = ConvertBase("2", 2, 10)
	if err == nil {
		t.Error("ConvertBase(\"2\", 2, 10) 期望返回错误")
	}
	fmt.Println("ConvertBase 测试通过")
}

func TestConvertBaseBig(t *testing.T) {
	// 大数：36进制 → 16进制
	got, err := ConvertBaseBig("zzzzzzzzzz", 36, 16)
	if err != nil {
		t.Fatalf("ConvertBaseBig 失败: %v", err)
	}
	n := new(big.Int)
	n.SetString("zzzzzzzzzz", 36)
	expected := n.Text(16)
	if got != expected {
		t.Errorf("ConvertBaseBig = %q, 期望 %q", got, expected)
	}

	// 无效进制
	_, err = ConvertBaseBig("10", 1, 10)
	if err == nil {
		t.Error("ConvertBaseBig 进制=1 期望返回错误")
	}
	_, err = ConvertBaseBig("10", 2, 63)
	if err == nil {
		t.Error("ConvertBaseBig 进制=63 期望返回错误")
	}
	fmt.Println("ConvertBaseBig 测试通过")
}

// ==================== 辅助函数 ====================

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
