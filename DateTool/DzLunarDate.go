package DateTool

import (
	"fmt"
	"time"
)

// DzLunarDate 封装农历日期，提供公历-农历互转、天干地支、生肖等功能。
// 值类型，不可变，天然协程安全。
// 覆盖范围：1900-2100 年。
type DzLunarDate struct {
	year   int
	month  int  // 1-12
	day    int  // 1-29/30
	isLeap bool // 是否闰月
	err    error
}

// ==================== 农历查表数据 ====================

// lunarInfo 存储 1900-2100 年的农历信息，每个 uint32 编码格式：
//   bit 0-3:   闰月月份（0=无闰月，1-12=闰X月）
//   bit 4-15:  1-12月大小月标志（bit4=1月, ..., bit15=12月; 1=30天大月, 0=29天小月）
//   bit 16:    闰月天数（0=29天, 1=30天）
//   bit 17-31: 保留
//
// 基准日期：1900-01-31（公历）= 农历庚子年正月初一
var lunarInfo = []uint32{
	0x04bd8, 0x04ae0, 0x0a570, 0x054d5, 0x0d260, 0x0d950, 0x16554, 0x056a0, 0x09ad0, 0x055d2, // 1900-1909
	0x04ae0, 0x0a5b6, 0x0a4d0, 0x0d250, 0x1d255, 0x0b540, 0x0d6a0, 0x0ada2, 0x095b0, 0x14977, // 1910-1919
	0x04970, 0x0a4b0, 0x0b4b5, 0x06a50, 0x06d40, 0x1ab54, 0x02b60, 0x09570, 0x052f2, 0x04970, // 1920-1929
	0x06566, 0x0d4a0, 0x0ea50, 0x16a95, 0x05ad0, 0x02b60, 0x186e3, 0x092e0, 0x1c8d7, 0x0c950, // 1930-1939
	0x0d4a0, 0x1d8a6, 0x0b550, 0x056a0, 0x1a5b4, 0x025d0, 0x092d0, 0x0d2b2, 0x0a950, 0x0b557, // 1940-1949
	0x06ca0, 0x0b550, 0x15355, 0x04da0, 0x0a5b0, 0x14573, 0x052b0, 0x0a9a8, 0x0e950, 0x06aa0, // 1950-1959
	0x0aea6, 0x0ab50, 0x04b60, 0x0aae4, 0x0a570, 0x05260, 0x0f263, 0x0d950, 0x05b57, 0x056a0, // 1960-1969
	0x096d0, 0x04dd5, 0x04ad0, 0x0a4d0, 0x0d4d4, 0x0d250, 0x0d558, 0x0b540, 0x0b6a0, 0x195a6, // 1970-1979
	0x095b0, 0x049b0, 0x0a974, 0x0a4b0, 0x0b27a, 0x06a50, 0x06d40, 0x0af46, 0x0ab60, 0x09570, // 1980-1989
	0x04af5, 0x04970, 0x064b0, 0x074a3, 0x0ea50, 0x06b58, 0x05ac0, 0x0ab60, 0x096d5, 0x092e0, // 1990-1999
	0x0c960, 0x0d954, 0x0d4a0, 0x0da50, 0x07552, 0x056a0, 0x0abb7, 0x025d0, 0x092d0, 0x0cab5, // 2000-2009
	0x0a950, 0x0b4a0, 0x0baa4, 0x0ad50, 0x055d9, 0x04ba0, 0x0a5b0, 0x15176, 0x052b0, 0x0a930, // 2010-2019
	0x07954, 0x06aa0, 0x0ad50, 0x05b52, 0x04b60, 0x0a6e6, 0x0a4e0, 0x0d260, 0x0ea65, 0x0d530, // 2020-2029
	0x05aa0, 0x076a3, 0x096d0, 0x04afb, 0x04ad0, 0x0a4d0, 0x1d0b6, 0x0d250, 0x0d520, 0x0dd45, // 2030-2039
	0x0b5a0, 0x056d0, 0x055b2, 0x049b0, 0x0a577, 0x0a4b0, 0x0aa50, 0x1b255, 0x06d20, 0x0ada0, // 2040-2049
	0x14b63, 0x09370, 0x049f8, 0x04970, 0x064b0, 0x168a6, 0x0ea50, 0x06aa0, 0x1a6c4, 0x0aae0, // 2050-2059
	0x092e0, 0x0d2e3, 0x0c960, 0x0d557, 0x0d4a0, 0x0da50, 0x05d55, 0x056a0, 0x0a6d0, 0x055d4, // 2060-2069
	0x052d0, 0x0a9b8, 0x0a950, 0x0b4a0, 0x0b6a6, 0x0ad50, 0x055a0, 0x0aba4, 0x0a5b0, 0x052b0, // 2070-2079
	0x0b273, 0x06930, 0x07337, 0x06aa0, 0x0ad50, 0x14b55, 0x04b60, 0x0a570, 0x054e4, 0x0d160, // 2080-2089
	0x0e968, 0x0d520, 0x0daa0, 0x16aa6, 0x056d0, 0x04ae0, 0x0a9d4, 0x0a4d0, 0x0d150, 0x0f252, // 2090-2099
	0x0d520, // 2100
}

// ==================== 农历数据查表辅助 ====================

// lunarLeapMonth 返回农历 year 的闰月月份，0 表示无闰月
func lunarLeapMonth(year int) int {
	if year < 1900 || year > 2100 {
		return 0
	}
	return int(lunarInfo[year-1900] & 0xf)
}

// lunarLeapDays 返回农历 year 闰月的天数，0 表示无闰月
func lunarLeapDays(year int) int {
	if lunarLeapMonth(year) == 0 {
		return 0
	}
	if lunarInfo[year-1900]&0x10000 != 0 {
		return 30
	}
	return 29
}

// lunarMonthDays 返回农历 year 月 month（1-12）的天数
func lunarMonthDays(year, month int) int {
	if month < 1 || month > 12 || year < 1900 || year > 2100 {
		return 0
	}
	if lunarInfo[year-1900]&(0x10000>>uint(month)) != 0 {
		return 30
	}
	return 29
}

// lunarYearDays 返回农历 year 全年的总天数
func lunarYearDays(year int) int {
	if year < 1900 || year > 2100 {
		return 0
	}
	total := 0
	for i := 1; i <= 12; i++ {
		total += lunarMonthDays(year, i)
	}
	total += lunarLeapDays(year)
	return total
}

// ==================== 公历↔农历核心转换 ====================

// solarToLunar 将公历日期转换为农历日期
// 返回 (农历年, 农历月, 农历日, 是否闰月)
func solarToLunar(y, m, d int) (int, int, int, bool) {
	// 基准日期: 1900-01-31 = 农历1900年正月初一
	baseDate := time.Date(1900, 1, 31, 0, 0, 0, 0, time.UTC)
	target := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)

	offset := int(target.Sub(baseDate).Hours() / 24)
	if offset < 0 {
		return 0, 0, 0, false
	}

	// 定位农历年份
	lunarYear := 1900
	for lunarYear <= 2100 {
		daysInYear := lunarYearDays(lunarYear)
		if offset < daysInYear {
			break
		}
		offset -= daysInYear
		lunarYear++
	}

	if lunarYear > 2100 {
		return 0, 0, 0, false
	}

	// 定位农历月份
	// 农历月份遍历顺序：1月, 2月, ..., leap月(正常), 闰leap月, leap+1月, ..., 12月
	leap := lunarLeapMonth(lunarYear)
	isLeap := false
	lunarMonth := 1

	for lunarMonth <= 12 {
		// 先检查正常月份
		daysInMonthVal := lunarMonthDays(lunarYear, lunarMonth)
		if offset < daysInMonthVal {
			break
		}
		offset -= daysInMonthVal

		// 如果该月有闰月，接着检查闰月
		if leap > 0 && lunarMonth == leap {
			daysInLeap := lunarLeapDays(lunarYear)
			if offset < daysInLeap {
				isLeap = true
				break
			}
			offset -= daysInLeap
		}

		lunarMonth++
	}

	lunarDay := offset + 1
	return lunarYear, lunarMonth, lunarDay, isLeap
}

// lunarToSolar 将农历日期转换为公历日期
// 返回 (公历年, 公历月, 公历日)
func lunarToSolar(ly, lm, ld int, isLeap bool) (int, int, int) {
	if ly < 1900 || ly > 2100 || lm < 1 || lm > 12 || ld < 1 {
		return 0, 0, 0
	}

	// 从基准日期 1900-01-31 开始计算偏移天数
	offset := 0

	// 累加完整年份的天数
	for y := 1900; y < ly; y++ {
		offset += lunarYearDays(y)
	}

	// 累加当年已过月份的天数
	leap := lunarLeapMonth(ly)
	for m := 1; m < lm; m++ {
		offset += lunarMonthDays(ly, m)
		// 如果闰月在当前月份之后，需要加上闰月天数
		if leap > 0 && m == leap {
			offset += lunarLeapDays(ly)
		}
	}

	// 如果目标是闰月
	if isLeap && leap == lm {
		offset += lunarMonthDays(ly, lm)
	}

	// 加上日偏移
	offset += ld - 1

	// 基准日期 + 偏移
	baseDate := time.Date(1900, 1, 31, 0, 0, 0, 0, time.UTC)
	result := baseDate.AddDate(0, 0, offset)
	return result.Year(), int(result.Month()), result.Day()
}

// ==================== 天干地支/生肖 ====================

var (
	gan  = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhi  = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	shengxiao = []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	monthName = []string{"正月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "冬月", "腊月"}
	dayName1  = []string{"初", "十", "廿", "卅"}
	dayName2  = []string{"", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
)

// ganZhiYear 返回农历年的天干地支
func ganZhiYear(year int) string {
	// 1900 年 = 庚子年
	idx := (year - 4) % 60 // 4年是甲子年
	if idx < 0 {
		idx += 60
	}
	return gan[idx%10] + zhi[idx%12]
}

// ganZhiMonth 返回农历月的天干地支
func ganZhiMonth(year, month int) string {
	// 月干支以年干为基础推算
	// 寅月（正月）天干 = (年干序号 * 2 + 2) % 10，地支固定从寅开始
	yearGanIdx := (year - 4) % 10
	if yearGanIdx < 0 {
		yearGanIdx += 10
	}
	monthGanIdx := (yearGanIdx*2 + 2 + month - 1) % 10
	monthZhiIdx := (month + 1) % 12 // 正月=寅(2), 二月=卯(3), ...
	return gan[monthGanIdx] + zhi[monthZhiIdx]
}

// ganZhiDay 返回农历日的天干地支
// 基于 1900-01-01（公历）= 甲戌日 计算
func ganZhiDay(y, m, d int) string {
	base := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	target := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	days := int(target.Sub(base).Hours() / 24)
	// 1900-01-01 = 甲戌 -> 干0支10
	ganIdx := (days + 0) % 10
	zhiIdx := (days + 10) % 12
	if ganIdx < 0 {
		ganIdx += 10
	}
	if zhiIdx < 0 {
		zhiIdx += 12
	}
	return gan[ganIdx] + zhi[zhiIdx]
}

// zodiac 返回农历年的生肖
func zodiac(year int) string {
	idx := (year - 4) % 12
	if idx < 0 {
		idx += 12
	}
	return shengxiao[idx]
}

// lunarDayName 返回农历日的中文表示（如 "初一"、"十五"、"三十"）
func lunarDayName(day int) string {
	if day < 1 || day > 30 {
		return ""
	}
	switch day {
	case 10:
		return "初十"
	case 20:
		return "二十"
	case 30:
		return "三十"
	default:
		return dayName1[day/10] + dayName2[day%10]
	}
}

// ==================== 构造函数 ====================

// NewDzLunarDate 从农历年月日构造 DzLunarDate
// isLeap 为可选参数，表示是否闰月（默认 false）
func NewDzLunarDate(year, month, day int, isLeap ...bool) DzLunarDate {
	leap := false
	if len(isLeap) > 0 {
		leap = isLeap[0]
	}

	// 校验范围
	if year < 1900 || year > 2100 {
		return DzLunarDate{err: fmt.Errorf("DateTool: lunar year %d out of range [1900, 2100]", year)}
	}
	if month < 1 || month > 12 {
		return DzLunarDate{err: fmt.Errorf("DateTool: lunar month %d out of range [1, 12]", month)}
	}

	// 校验闰月合法性
	leapMonth := lunarLeapMonth(year)
	if leap && leapMonth != month {
		return DzLunarDate{err: fmt.Errorf("DateTool: year %d has no leap month %d", year, month)}
	}

	// 校验日
	maxDay := lunarMonthDays(year, month)
	if leap {
		maxDay = lunarLeapDays(year)
	}
	if day < 1 || day > maxDay {
		return DzLunarDate{err: fmt.Errorf("DateTool: lunar day %d out of range [1, %d]", day, maxDay)}
	}

	return DzLunarDate{year: year, month: month, day: day, isLeap: leap}
}

// LunarFromSolar 从公历 DzDateTime 转换为农历 DzLunarDate
func LunarFromSolar(dt DzDateTime) DzLunarDate {
	if dt.err != nil {
		return DzLunarDate{err: dt.err}
	}
	ly, lm, ld, isLeap := solarToLunar(dt.t.Year(), int(dt.t.Month()), dt.t.Day())
	if ly == 0 {
		return DzLunarDate{err: fmt.Errorf("DateTool: failed to convert solar to lunar")}
	}
	return DzLunarDate{year: ly, month: lm, day: ld, isLeap: isLeap}
}

// ParseLunar 从农历年月日构造 DzLunarDate（同 NewDzLunarDate）
func ParseLunar(year, month, day int, isLeap ...bool) DzLunarDate {
	return NewDzLunarDate(year, month, day, isLeap...)
}

// ==================== 错误检查 ====================

// Err 返回链式调用中累积的错误
func (l DzLunarDate) Err() error {
	return l.err
}

// IsValid 返回 DzLunarDate 是否有效
func (l DzLunarDate) IsValid() bool {
	return l.err == nil
}

// ==================== Getter ====================

// Year 返回农历年
func (l DzLunarDate) Year() int {
	return l.year
}

// Month 返回农历月（1-12）
func (l DzLunarDate) Month() int {
	return l.month
}

// Day 返回农历日
func (l DzLunarDate) Day() int {
	return l.day
}

// IsLeap 返回是否闰月
func (l DzLunarDate) IsLeap() bool {
	return l.isLeap
}

// YearGanZhi 返回天干地支年名（如 "甲子"）
func (l DzLunarDate) YearGanZhi() string {
	if l.err != nil {
		return ""
	}
	return ganZhiYear(l.year)
}

// MonthGanZhi 返回天干地支月名（如 "丙寅"）
func (l DzLunarDate) MonthGanZhi() string {
	if l.err != nil {
		return ""
	}
	return ganZhiMonth(l.year, l.month)
}

// DayGanZhi 返回天干地支日名（如 "庚午"）
func (l DzLunarDate) DayGanZhi() string {
	if l.err != nil {
		return ""
	}
	// 需要公历日期来计算日干支
	sy, sm, sd := lunarToSolar(l.year, l.month, l.day, l.isLeap)
	return ganZhiDay(sy, sm, sd)
}

// Zodiac 返回生肖（如 "鼠"）
func (l DzLunarDate) Zodiac() string {
	if l.err != nil {
		return ""
	}
	return zodiac(l.year)
}

// MonthName 返回农历月名（如 "正月"、"腊月"）
func (l DzLunarDate) MonthName() string {
	if l.err != nil {
		return ""
	}
	if l.isLeap {
		return "闰" + monthName[l.month-1]
	}
	return monthName[l.month-1]
}

// DayName 返回农历日名（如 "初一"、"十五"）
func (l DzLunarDate) DayName() string {
	if l.err != nil {
		return ""
	}
	return lunarDayName(l.day)
}

// YearName 返回农历年份全名（如 "庚子年"）
func (l DzLunarDate) YearName() string {
	if l.err != nil {
		return ""
	}
	return ganZhiYear(l.year) + "年"
}

// LeapMonth 返回该年闰月月份（0=无闰月）
func (l DzLunarDate) LeapMonth() int {
	if l.err != nil {
		return 0
	}
	return lunarLeapMonth(l.year)
}

// DaysInMonth 返回当前农历月的天数
func (l DzLunarDate) DaysInMonth() int {
	if l.err != nil {
		return 0
	}
	if l.isLeap {
		return lunarLeapDays(l.year)
	}
	return lunarMonthDays(l.year, l.month)
}

// DaysInYear 返回当前农历年的总天数
func (l DzLunarDate) DaysInYear() int {
	if l.err != nil {
		return 0
	}
	return lunarYearDays(l.year)
}

// ==================== 转换 ====================

// ToSolar 转换为公历 DzDateTime
func (l DzLunarDate) ToSolar() DzDateTime {
	if l.err != nil {
		return DzDateTime{err: l.err}
	}
	sy, sm, sd := lunarToSolar(l.year, l.month, l.day, l.isLeap)
	if sy == 0 {
		return DzDateTime{err: fmt.Errorf("DateTool: failed to convert lunar to solar")}
	}
	return DzDateTime{t: time.Date(sy, time.Month(sm), sd, 0, 0, 0, 0, time.Local)}
}

// ==================== 格式化 ====================

// String 实现 fmt.Stringer，输出如 "庚子年闰四月十五"
func (l DzLunarDate) String() string {
	if l.err != nil {
		return fmt.Sprintf("DzLunarDate{err: %v}", l.err)
	}
	return l.YearGanZhi() + "年" + l.MonthName() + l.DayName()
}

// Format 格式化农历日期，支持占位符：
//   YYYY  - 农历年（数字）
//   YY    - 农历年（后两位）
//   MM    - 农历月（数字，1-12）
//   DD    - 农历日（数字，1-30）
//   GZ    - 天干地支年
//   ZODIAC - 生肖
//   MNAME - 月名（正月、二月…）
//   DNAME - 日名（初一、初二…）
func (l DzLunarDate) Format(layout string) string {
	if l.err != nil {
		return ""
	}
	result := layout
	result = replaceOnce(result, "YYYY", fmt.Sprintf("%04d", l.year))
	result = replaceOnce(result, "YY", fmt.Sprintf("%02d", l.year%100))
	result = replaceOnce(result, "MM", fmt.Sprintf("%02d", l.month))
	result = replaceOnce(result, "DD", fmt.Sprintf("%02d", l.day))
	result = replaceOnce(result, "GZ", l.YearGanZhi())
	result = replaceOnce(result, "ZODIAC", l.Zodiac())
	result = replaceOnce(result, "MNAME", l.MonthName())
	result = replaceOnce(result, "DNAME", l.DayName())
	return result
}

// replaceOnce 替换字符串中第一个匹配
func replaceOnce(s, old, new string) string {
	idx := indexOf(s, old)
	if idx < 0 {
		return s
	}
	return s[:idx] + new + s[idx+len(old):]
}

// indexOf 查找子串位置
func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// ==================== 比较 ====================

// IsBefore 判断是否在 other 之前
func (l DzLunarDate) IsBefore(other DzLunarDate) bool {
	if l.err != nil || other.err != nil {
		return false
	}
	lTime := l.ToSolar()
	oTime := other.ToSolar()
	return lTime.IsBefore(oTime)
}

// IsAfter 判断是否在 other 之后
func (l DzLunarDate) IsAfter(other DzLunarDate) bool {
	if l.err != nil || other.err != nil {
		return false
	}
	lTime := l.ToSolar()
	oTime := other.ToSolar()
	return lTime.IsAfter(oTime)
}

// IsEqual 判断是否与 other 相等
func (l DzLunarDate) IsEqual(other DzLunarDate) bool {
	if l.err != nil || other.err != nil {
		return false
	}
	return l.year == other.year && l.month == other.month && l.day == other.day && l.isLeap == other.isLeap
}

// ==================== 静态工具函数 ====================

// LunarToSolar 将农历日期转换为公历日期（纯数据转换）
func LunarToSolar(year, month, day int, isLeap bool) (int, int, int) {
	return lunarToSolar(year, month, day, isLeap)
}

// SolarToLunar 将公历日期转换为农历日期（纯数据转换）
// 返回 (农历年, 农历月, 农历日, 是否闰月)
func SolarToLunar(year, month, day int) (int, int, int, bool) {
	return solarToLunar(year, month, day)
}

// LeapMonthOfYear 返回指定农历年的闰月月份（0=无闰月）
func LeapMonthOfYear(year int) int {
	return lunarLeapMonth(year)
}

// GanZhiYear 返回指定年份的天干地支
func GanZhiYear(year int) string {
	return ganZhiYear(year)
}

// ZodiacOfYear 返回指定年份的生肖
func ZodiacOfYear(year int) string {
	return zodiac(year)
}
