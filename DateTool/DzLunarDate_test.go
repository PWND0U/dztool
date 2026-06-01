package DateTool

import (
	"fmt"
	"testing"
)

// ==================== 构造函数测试 ====================

func TestNewDzLunarDate(t *testing.T) {
	l := NewDzLunarDate(2024, 6, 15)
	fmt.Println("NewDzLunarDate IsValid:", l.IsValid())
	fmt.Println("NewDzLunarDate Year:", l.Year())
	fmt.Println("NewDzLunarDate Month:", l.Month())
	fmt.Println("NewDzLunarDate Day:", l.Day())
	fmt.Println("NewDzLunarDate IsLeap:", l.IsLeap())
	fmt.Println("NewDzLunarDate String:", l.String())
}

func TestNewDzLunarDate_InvalidYear(t *testing.T) {
	l := NewDzLunarDate(1800, 1, 1)
	fmt.Println("InvalidYear IsValid:", l.IsValid())
	fmt.Println("InvalidYear Err:", l.Err())
}

func TestNewDzLunarDate_InvalidMonth(t *testing.T) {
	l := NewDzLunarDate(2024, 13, 1)
	fmt.Println("InvalidMonth IsValid:", l.IsValid())
	fmt.Println("InvalidMonth Err:", l.Err())
}

func TestNewDzLunarDate_InvalidLeap(t *testing.T) {
	// 2024 年闰月不是5月
	l := NewDzLunarDate(2024, 5, 1, true)
	fmt.Println("InvalidLeap IsValid:", l.IsValid())
	fmt.Println("InvalidLeap Err:", l.Err())
}

func TestParseLunar(t *testing.T) {
	l := ParseLunar(2024, 1, 1)
	fmt.Println("ParseLunar String:", l.String())
}

// ==================== 公历转农历测试 ====================

func TestLunarFromSolar(t *testing.T) {
	// 2024-06-15 公历
	dt := DzDate(2024, 6, 15)
	l := LunarFromSolar(dt)
	fmt.Println("LunarFromSolar (2024-06-15):", l.String())
	fmt.Println("  Year:", l.Year())
	fmt.Println("  Month:", l.Month())
	fmt.Println("  Day:", l.Day())
	fmt.Println("  IsLeap:", l.IsLeap())
}

func TestLunarFromSolar_SpringFestival(t *testing.T) {
	// 2024-02-10 是农历甲辰年正月初一（2024年春节）
	dt := DzDate(2024, 2, 10)
	l := LunarFromSolar(dt)
	fmt.Println("SpringFestival (2024-02-10):", l.String())
	fmt.Println("  Year:", l.Year())
	fmt.Println("  Month:", l.Month())
	fmt.Println("  Day:", l.Day())
}

func TestLunarFromSolar_LeapMonth(t *testing.T) {
	// 2023年有闰二月
	// 2023-03-22 附近应该是农历闰二月
	dt := DzDate(2023, 3, 22)
	l := LunarFromSolar(dt)
	fmt.Println("2023-03-22 Lunar:", l.String())
	fmt.Println("  IsLeap:", l.IsLeap())
}

// ==================== 农历转公历测试 ====================

func TestToSolar(t *testing.T) {
	// 农历 2024年正月初一 -> 公历
	l := NewDzLunarDate(2024, 1, 1)
	dt := l.ToSolar()
	fmt.Println("ToSolar (2024-01-01 lunar):", dt.ToDateString())
}

func TestToSolar_LeapMonth(t *testing.T) {
	// 2023年闰二月十五
	l := NewDzLunarDate(2023, 2, 15, true)
	dt := l.ToSolar()
	fmt.Println("ToSolar (2023 leap 2-15):", dt.ToDateString())
	fmt.Println("  IsLeap:", l.IsLeap())
	fmt.Println("  LeapMonth:", l.LeapMonth())
}

// ==================== 往返测试 ====================

func TestRoundTrip(t *testing.T) {
	// 公历 -> 农历 -> 公历
	orig := DzDate(2024, 6, 15)
	l := LunarFromSolar(orig)
	back := l.ToSolar()
	fmt.Println("RoundTrip original:", orig.ToDateString())
	fmt.Println("RoundTrip lunar:", l.String())
	fmt.Println("RoundTrip solar:", back.ToDateString())
	fmt.Println("RoundTrip equal:", orig.ToDateString() == back.ToDateString())
}

func TestRoundTrip_LeapMonth(t *testing.T) {
	// 农历 -> 公历 -> 农历
	orig := NewDzLunarDate(2023, 2, 15, true)
	solar := orig.ToSolar()
	back := LunarFromSolar(solar)
	fmt.Println("RoundTrip leap original:", orig.String())
	fmt.Println("RoundTrip leap solar:", solar.ToDateString())
	fmt.Println("RoundTrip leap back:", back.String())
	fmt.Println("RoundTrip leap IsEqual:", orig.IsEqual(back))
}

// ==================== 天干地支/生肖测试 ====================

func TestGanZhi(t *testing.T) {
	// 2024 = 甲辰年
	l := NewDzLunarDate(2024, 1, 1)
	fmt.Println("2024 YearGanZhi:", l.YearGanZhi())
	fmt.Println("2024 Zodiac:", l.Zodiac())

	// 1900 = 庚子年
	l2 := NewDzLunarDate(1900, 1, 1)
	fmt.Println("1900 YearGanZhi:", l2.YearGanZhi())
	fmt.Println("1900 Zodiac:", l2.Zodiac())
}

func TestMonthName(t *testing.T) {
	l := NewDzLunarDate(2024, 1, 1)
	fmt.Println("Month 1 Name:", l.MonthName())

	l2 := NewDzLunarDate(2024, 12, 1)
	fmt.Println("Month 12 Name:", l2.MonthName())
}

func TestDayName(t *testing.T) {
	l1 := NewDzLunarDate(2024, 1, 1)
	fmt.Println("Day 1 Name:", l1.DayName())

	l15 := NewDzLunarDate(2024, 1, 15)
	fmt.Println("Day 15 Name:", l15.DayName())

	l10 := NewDzLunarDate(2024, 1, 10)
	fmt.Println("Day 10 Name:", l10.DayName())

	l20 := NewDzLunarDate(2024, 1, 20)
	fmt.Println("Day 20 Name:", l20.DayName())

	// 找一个有30天的月来测试三十
	// 2024年农历各月天数
	for m := 1; m <= 12; m++ {
		if lunarMonthDays(2024, m) == 30 {
			l30 := NewDzLunarDate(2024, m, 30)
			fmt.Println("Day 30 Name (month", m, "):", l30.DayName())
			break
		}
	}
}

func TestYearName(t *testing.T) {
	l := NewDzLunarDate(2024, 1, 1)
	fmt.Println("YearName:", l.YearName())
}

func TestDayGanZhi(t *testing.T) {
	l := NewDzLunarDate(2024, 1, 1)
	fmt.Println("DayGanZhi:", l.DayGanZhi())
}

// ==================== 天数信息测试 ====================

func TestDaysInfo(t *testing.T) {
	l := NewDzLunarDate(2024, 1, 1)
	fmt.Println("2024 LeapMonth:", l.LeapMonth())
	fmt.Println("2024 DaysInMonth (1月):", l.DaysInMonth())
	fmt.Println("2024 DaysInYear:", l.DaysInYear())
}

// ==================== 格式化测试 ====================

func TestLunarFormat(t *testing.T) {
	l := NewDzLunarDate(2024, 6, 15)
	fmt.Println("Format YYYY-MM-DD:", l.Format("YYYY年MM月DD日"))
	fmt.Println("Format GZ ZODIAC:", l.Format("GZ年ZODIAC"))
	fmt.Println("Format MNAME DNAME:", l.Format("MNAMEDNAME"))
}

// ==================== 比较测试 ====================

func TestLunarIsBefore(t *testing.T) {
	l1 := NewDzLunarDate(2024, 1, 1)
	l2 := NewDzLunarDate(2024, 1, 15)
	fmt.Println("IsBefore:", l1.IsBefore(l2))
	fmt.Println("IsBefore (reverse):", l2.IsBefore(l1))
}

func TestLunarIsAfter(t *testing.T) {
	l1 := NewDzLunarDate(2024, 1, 15)
	l2 := NewDzLunarDate(2024, 1, 1)
	fmt.Println("IsAfter:", l1.IsAfter(l2))
}

func TestLunarIsEqual(t *testing.T) {
	l1 := NewDzLunarDate(2024, 1, 1)
	l2 := NewDzLunarDate(2024, 1, 1)
	fmt.Println("IsEqual (same):", l1.IsEqual(l2))

	l3 := NewDzLunarDate(2024, 1, 2)
	fmt.Println("IsEqual (diff):", l1.IsEqual(l3))
}

// ==================== DzDateTime.ToLunar 测试 ====================

func TestDzDateTime_ToLunar(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	l := dt.ToLunar()
	fmt.Println("ToLunar:", l.String())
}

// ==================== 静态工具函数测试 ====================

func TestStaticSolarToLunar(t *testing.T) {
	ly, lm, ld, isLeap := SolarToLunar(2024, 6, 15)
	fmt.Println("SolarToLunar(2024,6,15):", ly, lm, ld, isLeap)
}

func TestStaticLunarToSolar(t *testing.T) {
	sy, sm, sd := LunarToSolar(2024, 1, 1, false)
	fmt.Println("LunarToSolar(2024,1,1):", sy, sm, sd)
}

func TestStaticLeapMonthOfYear(t *testing.T) {
	fmt.Println("LeapMonthOfYear(2023):", LeapMonthOfYear(2023))
	fmt.Println("LeapMonthOfYear(2024):", LeapMonthOfYear(2024))
	fmt.Println("LeapMonthOfYear(2025):", LeapMonthOfYear(2025))
}

func TestStaticGanZhiYear(t *testing.T) {
	fmt.Println("GanZhiYear(2024):", GanZhiYear(2024))
	fmt.Println("GanZhiYear(1900):", GanZhiYear(1900))
	fmt.Println("GanZhiYear(2000):", GanZhiYear(2000))
}

func TestStaticZodiacOfYear(t *testing.T) {
	fmt.Println("ZodiacOfYear(2024):", ZodiacOfYear(2024))
	fmt.Println("ZodiacOfYear(2023):", ZodiacOfYear(2023))
	fmt.Println("ZodiacOfYear(2000):", ZodiacOfYear(2000))
}

// ==================== 闰月完整验证 ====================

func TestLeapMonth2023(t *testing.T) {
	// 2023年闰二月，验证从公历3月22日开始进入闰二月
	leap := LeapMonthOfYear(2023)
	fmt.Println("2023 LeapMonth:", leap)

	// 3月22日应该落在闰二月
	dt1 := DzDate(2023, 3, 22)
	l1 := LunarFromSolar(dt1)
	fmt.Println("2023-03-22 Lunar:", l1.String(), "IsLeap:", l1.IsLeap())

	// 从农历闰二月初一转回公历
	l2 := NewDzLunarDate(2023, 2, 1, true)
	dt2 := l2.ToSolar()
	fmt.Println("Lunar 2023 leap 2-1 -> Solar:", dt2.ToDateString())
}
