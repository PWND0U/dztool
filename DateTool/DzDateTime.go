package DateTool

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// DzDateTime 封装 time.Time，提供 Joda-Time 风格的链式日期时间操作。
// 值类型，所有链式方法返回新的 DzDateTime 实例，不修改原始对象，天然协程安全。
type DzDateTime struct {
	t   time.Time
	err error
}

// ==================== 构造函数 ====================

// NewDzDateTime 从 time.Time 构造 DzDateTime
func NewDzDateTime(t time.Time) DzDateTime {
	return DzDateTime{t: t}
}

// Now 返回当前本地时间的 DzDateTime
func Now() DzDateTime {
	return DzDateTime{t: time.Now()}
}

// DzDate 从年月日构造 DzDateTime（时间部分为零值）
func DzDate(y, m, d int) DzDateTime {
	return DzDateTime{t: time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)}
}

// DzDateTimeOf 从年月日时分秒构造 DzDateTime
func DzDateTimeOf(y, m, d, h, mi, s int) DzDateTime {
	return DzDateTime{t: time.Date(y, time.Month(m), d, h, mi, s, 0, time.Local)}
}

// Parse 使用 Go 标准 layout 解析日期时间字符串
func Parse(layout, value string) DzDateTime {
	t, err := time.Parse(layout, value)
	if err != nil {
		return DzDateTime{err: fmt.Errorf("DateTool: parse %q with layout %q: %w", value, layout, err)}
	}
	return DzDateTime{t: t}
}

// ParseInLocation 使用指定时区解析日期时间字符串
func ParseInLocation(layout, value string, loc *time.Location) DzDateTime {
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return DzDateTime{err: fmt.Errorf("DateTool: parse %q with layout %q in location: %w", value, layout, err)}
	}
	return DzDateTime{t: t}
}

// FromTimestamp 从 Unix 时间戳（秒）构造 DzDateTime
func FromTimestamp(ts int64) DzDateTime {
	return DzDateTime{t: time.Unix(ts, 0)}
}

// FromTimestampMilli 从毫秒时间戳构造 DzDateTime
func FromTimestampMilli(ts int64) DzDateTime {
	return DzDateTime{t: time.UnixMilli(ts)}
}

// FromTimestampMicro 从微秒时间戳构造 DzDateTime
func FromTimestampMicro(ts int64) DzDateTime {
	return DzDateTime{t: time.UnixMicro(ts)}
}

// ==================== 错误检查 ====================

// Err 返回链式调用中累积的错误
func (d DzDateTime) Err() error {
	return d.err
}

// IsValid 返回 DzDateTime 是否有效（无错误）
func (d DzDateTime) IsValid() bool {
	return d.err == nil
}

// ==================== 链式 Setter（返回新 DzDateTime）====================

// WithYear 返回设置年份后的新 DzDateTime
func (d DzDateTime) WithYear(y int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.AddDate(y-d.t.Year(), 0, 0)}
}

// WithMonth 返回设置月份后的新 DzDateTime
func (d DzDateTime) WithMonth(m int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.AddDate(0, m-int(d.t.Month()), 0)}
}

// WithDay 返回设置日期后的新 DzDateTime
func (d DzDateTime) WithDay(dd int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.AddDate(0, 0, dd-d.t.Day())}
}

// WithHour 返回设置小时后的新 DzDateTime
func (d DzDateTime) WithHour(h int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), h, d.t.Minute(), d.t.Second(), d.t.Nanosecond(), d.t.Location())}
}

// WithMinute 返回设置分钟后的新 DzDateTime
func (d DzDateTime) WithMinute(m int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), d.t.Hour(), m, d.t.Second(), d.t.Nanosecond(), d.t.Location())}
}

// WithSecond 返回设置秒后的新 DzDateTime
func (d DzDateTime) WithSecond(s int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), d.t.Hour(), d.t.Minute(), s, d.t.Nanosecond(), d.t.Location())}
}

// ==================== 链式算术（返回新 DzDateTime）====================

// AddYears 增加指定年数
func (d DzDateTime) AddYears(n int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.AddDate(n, 0, 0)}
}

// AddMonths 增加指定月数
func (d DzDateTime) AddMonths(n int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.AddDate(0, n, 0)}
}

// AddDays 增加指定天数
func (d DzDateTime) AddDays(n int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.AddDate(0, 0, n)}
}

// AddHours 增加指定小时数
func (d DzDateTime) AddHours(n int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.Add(time.Duration(n) * time.Hour)}
}

// AddMinutes 增加指定分钟数
func (d DzDateTime) AddMinutes(n int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.Add(time.Duration(n) * time.Minute)}
}

// AddSeconds 增加指定秒数
func (d DzDateTime) AddSeconds(n int) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.Add(time.Duration(n) * time.Second)}
}

// AddDuration 增加指定时间段
func (d DzDateTime) AddDuration(dur time.Duration) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.Add(dur)}
}

// SubYears 减少指定年数
func (d DzDateTime) SubYears(n int) DzDateTime {
	return d.AddYears(-n)
}

// SubMonths 减少指定月数
func (d DzDateTime) SubMonths(n int) DzDateTime {
	return d.AddMonths(-n)
}

// SubDays 减少指定天数
func (d DzDateTime) SubDays(n int) DzDateTime {
	return d.AddDays(-n)
}

// SubHours 减少指定小时数
func (d DzDateTime) SubHours(n int) DzDateTime {
	return d.AddHours(-n)
}

// SubMinutes 减少指定分钟数
func (d DzDateTime) SubMinutes(n int) DzDateTime {
	return d.AddMinutes(-n)
}

// SubSeconds 减少指定秒数
func (d DzDateTime) SubSeconds(n int) DzDateTime {
	return d.AddSeconds(-n)
}

// ==================== 边界方法 ====================

// StartOfDay 返回当天开始时间 00:00:00
func (d DzDateTime) StartOfDay() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), 0, 0, 0, 0, d.t.Location())}
}

// EndOfDay 返回当天结束时间 23:59:59.999999999
func (d DzDateTime) EndOfDay() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), 23, 59, 59, 999999999, d.t.Location())}
}

// StartOfHour 返回当前小时的开始时间
func (d DzDateTime) StartOfHour() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), d.t.Hour(), 0, 0, 0, d.t.Location())}
}

// EndOfHour 返回当前小时的结束时间
func (d DzDateTime) EndOfHour() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), d.t.Hour(), 59, 59, 999999999, d.t.Location())}
}

// StartOfMinute 返回当前分钟的开始时间
func (d DzDateTime) StartOfMinute() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), d.t.Hour(), d.t.Minute(), 0, 0, d.t.Location())}
}

// EndOfMinute 返回当前分钟的结束时间
func (d DzDateTime) EndOfMinute() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day(), d.t.Hour(), d.t.Minute(), 59, 999999999, d.t.Location())}
}

// StartOfMonth 返回当月开始时间
func (d DzDateTime) StartOfMonth() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), 1, 0, 0, 0, 0, d.t.Location())}
}

// EndOfMonth 返回当月结束时间
func (d DzDateTime) EndOfMonth() DzDateTime {
	if d.err != nil {
		return d
	}
	days := daysInMonth(d.t.Year(), d.t.Month())
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), days, 23, 59, 59, 999999999, d.t.Location())}
}

// StartOfYear 返回当年开始时间
func (d DzDateTime) StartOfYear() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), 1, 1, 0, 0, 0, 0, d.t.Location())}
}

// EndOfYear 返回当年结束时间
func (d DzDateTime) EndOfYear() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: time.Date(d.t.Year(), 12, 31, 23, 59, 59, 999999999, d.t.Location())}
}

// StartOfWeek 返回本周开始时间（周一开始）
func (d DzDateTime) StartOfWeek() DzDateTime {
	if d.err != nil {
		return d
	}
	wd := int(d.t.Weekday())
	if wd == 0 {
		wd = 7 // Sunday -> 7
	}
	offset := 1 - wd // Monday = 1
	return DzDateTime{t: time.Date(d.t.Year(), d.t.Month(), d.t.Day()+offset, 0, 0, 0, 0, d.t.Location())}
}

// EndOfWeek 返回本周结束时间（周日结束）
func (d DzDateTime) EndOfWeek() DzDateTime {
	if d.err != nil {
		return d
	}
	return d.StartOfWeek().AddDays(6).EndOfDay()
}

// ==================== 比较方法 ====================

// IsBefore 判断是否在 other 之前
func (d DzDateTime) IsBefore(other DzDateTime) bool {
	if d.err != nil || other.err != nil {
		return false
	}
	return d.t.Before(other.t)
}

// IsAfter 判断是否在 other 之后
func (d DzDateTime) IsAfter(other DzDateTime) bool {
	if d.err != nil || other.err != nil {
		return false
	}
	return d.t.After(other.t)
}

// IsEqual 判断是否与 other 相等
func (d DzDateTime) IsEqual(other DzDateTime) bool {
	if d.err != nil || other.err != nil {
		return false
	}
	return d.t.Equal(other.t)
}

// IsBetween 判断是否在 start 和 end 之间（包含边界）
func (d DzDateTime) IsBetween(start, end DzDateTime) bool {
	if d.err != nil || start.err != nil || end.err != nil {
		return false
	}
	return (d.t.After(start.t) || d.t.Equal(start.t)) && (d.t.Before(end.t) || d.t.Equal(end.t))
}

// IsToday 判断是否是今天
func (d DzDateTime) IsToday() bool {
	if d.err != nil {
		return false
	}
	now := time.Now()
	return d.t.Year() == now.Year() && d.t.Month() == now.Month() && d.t.Day() == now.Day()
}

// IsLeapYear 判断是否闰年
func (d DzDateTime) IsLeapYear() bool {
	if d.err != nil {
		return false
	}
	return isLeapYear(d.t.Year())
}

// IsWeekend 判断是否周末
func (d DzDateTime) IsWeekend() bool {
	if d.err != nil {
		return false
	}
	wd := d.t.Weekday()
	return wd == time.Saturday || wd == time.Sunday
}

// IsWeekday 判断是否工作日
func (d DzDateTime) IsWeekday() bool {
	return !d.IsWeekend()
}

// IsPast 判断是否是过去的时间
func (d DzDateTime) IsPast() bool {
	if d.err != nil {
		return false
	}
	return d.t.Before(time.Now())
}

// IsFuture 判断是否是未来的时间
func (d DzDateTime) IsFuture() bool {
	if d.err != nil {
		return false
	}
	return d.t.After(time.Now())
}

// IsSameDay 判断是否与 other 同一天
func (d DzDateTime) IsSameDay(other DzDateTime) bool {
	if d.err != nil || other.err != nil {
		return false
	}
	return d.t.Year() == other.t.Year() && d.t.Month() == other.t.Month() && d.t.Day() == other.t.Day()
}

// IsSameMonth 判断是否与 other 同一月
func (d DzDateTime) IsSameMonth(other DzDateTime) bool {
	if d.err != nil || other.err != nil {
		return false
	}
	return d.t.Year() == other.t.Year() && d.t.Month() == other.t.Month()
}

// IsSameYear 判断是否与 other 同一年
func (d DzDateTime) IsSameYear(other DzDateTime) bool {
	if d.err != nil || other.err != nil {
		return false
	}
	return d.t.Year() == other.t.Year()
}

// ==================== Getter ====================

// Year 返回年份
func (d DzDateTime) Year() int {
	if d.err != nil {
		return 0
	}
	return d.t.Year()
}

// Month 返回月份（1-12）
func (d DzDateTime) Month() int {
	if d.err != nil {
		return 0
	}
	return int(d.t.Month())
}

// Day 返回日
func (d DzDateTime) Day() int {
	if d.err != nil {
		return 0
	}
	return d.t.Day()
}

// Hour 返回小时
func (d DzDateTime) Hour() int {
	if d.err != nil {
		return 0
	}
	return d.t.Hour()
}

// Minute 返回分钟
func (d DzDateTime) Minute() int {
	if d.err != nil {
		return 0
	}
	return d.t.Minute()
}

// Second 返回秒
func (d DzDateTime) Second() int {
	if d.err != nil {
		return 0
	}
	return d.t.Second()
}

// Nanosecond 返回纳秒
func (d DzDateTime) Nanosecond() int {
	if d.err != nil {
		return 0
	}
	return d.t.Nanosecond()
}

// Weekday 返回星期几
func (d DzDateTime) Weekday() time.Weekday {
	if d.err != nil {
		return time.Sunday
	}
	return d.t.Weekday()
}

// YearDay 返回一年中的第几天
func (d DzDateTime) YearDay() int {
	if d.err != nil {
		return 0
	}
	return d.t.YearDay()
}

// ISOWeek 返回 ISO 周号
func (d DzDateTime) ISOWeek() (int, int) {
	if d.err != nil {
		return 0, 0
	}
	return d.t.ISOWeek()
}

// Timestamp 返回 Unix 时间戳（秒）
func (d DzDateTime) Timestamp() int64 {
	if d.err != nil {
		return 0
	}
	return d.t.Unix()
}

// TimestampMilli 返回毫秒时间戳
func (d DzDateTime) TimestampMilli() int64 {
	if d.err != nil {
		return 0
	}
	return d.t.UnixMilli()
}

// TimestampMicro 返回微秒时间戳
func (d DzDateTime) TimestampMicro() int64 {
	if d.err != nil {
		return 0
	}
	return d.t.UnixMicro()
}

// TimestampNano 返回纳秒时间戳
func (d DzDateTime) TimestampNano() int64 {
	if d.err != nil {
		return 0
	}
	return d.t.UnixNano()
}

// DaysInMonth 返回当前月份的天数
func (d DzDateTime) DaysInMonth() int {
	if d.err != nil {
		return 0
	}
	return daysInMonth(d.t.Year(), d.t.Month())
}

// DaysInYear 返回当前年份的天数
func (d DzDateTime) DaysInYear() int {
	if d.err != nil {
		return 0
	}
	if isLeapYear(d.t.Year()) {
		return 366
	}
	return 365
}

// ==================== 差量 ====================

// DiffInYears 计算与 other 相差的年数
func (d DzDateTime) DiffInYears(other DzDateTime) int {
	if d.err != nil || other.err != nil {
		return 0
	}
	return d.t.Year() - other.t.Year()
}

// DiffInMonths 计算与 other 相差的月数
func (d DzDateTime) DiffInMonths(other DzDateTime) int {
	if d.err != nil || other.err != nil {
		return 0
	}
	return (d.t.Year()-other.t.Year())*12 + int(d.t.Month()) - int(other.t.Month())
}

// DiffInDays 计算与 other 相差的天数
func (d DzDateTime) DiffInDays(other DzDateTime) int {
	if d.err != nil || other.err != nil {
		return 0
	}
	return int(d.t.Sub(other.t).Hours() / 24)
}

// DiffInHours 计算与 other 相差的小时数
func (d DzDateTime) DiffInHours(other DzDateTime) float64 {
	if d.err != nil || other.err != nil {
		return 0
	}
	return d.t.Sub(other.t).Hours()
}

// DiffInMinutes 计算与 other 相差的分钟数
func (d DzDateTime) DiffInMinutes(other DzDateTime) float64 {
	if d.err != nil || other.err != nil {
		return 0
	}
	return d.t.Sub(other.t).Minutes()
}

// DiffInSeconds 计算与 other 相差的秒数
func (d DzDateTime) DiffInSeconds(other DzDateTime) float64 {
	if d.err != nil || other.err != nil {
		return 0
	}
	return d.t.Sub(other.t).Seconds()
}

// ==================== 格式化 ====================

// Format 使用 Go 标准 layout 格式化日期时间
func (d DzDateTime) Format(layout string) string {
	if d.err != nil {
		return ""
	}
	return d.t.Format(layout)
}

// ToDateString 返回日期字符串 "2006-01-02"
func (d DzDateTime) ToDateString() string {
	return d.Format("2006-01-02")
}

// ToTimeString 返回时间字符串 "15:04:05"
func (d DzDateTime) ToTimeString() string {
	return d.Format("15:04:05")
}

// ToDateTimeString 返回日期时间字符串 "2006-01-02 15:04:05"
func (d DzDateTime) ToDateTimeString() string {
	return d.Format("2006-01-02 15:04:05")
}

// ToAtomString 返回 Atom 格式字符串 "2006-01-02T15:04:05Z07:00"
func (d DzDateTime) ToAtomString() string {
	return d.Format(time.RFC3339)
}

// ToRfc3339String 返回 RFC3339 格式字符串
func (d DzDateTime) ToRfc3339String() string {
	return d.Format(time.RFC3339)
}

// ToRfc1123String 返回 RFC1123 格式字符串
func (d DzDateTime) ToRfc1123String() string {
	return d.Format(time.RFC1123)
}

// ==================== 时区 ====================

// In 转换到指定时区
func (d DzDateTime) In(loc *time.Location) DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.In(loc)}
}

// UTC 转换到 UTC 时区
func (d DzDateTime) UTC() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.UTC()}
}

// Local 转换到本地时区
func (d DzDateTime) Local() DzDateTime {
	if d.err != nil {
		return d
	}
	return DzDateTime{t: d.t.Local()}
}

// ToTimezone 按时区名称转换时区
func (d DzDateTime) ToTimezone(name string) DzDateTime {
	if d.err != nil {
		return d
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return DzDateTime{err: fmt.Errorf("DateTool: unknown timezone %q: %w", name, err)}
	}
	return DzDateTime{t: d.t.In(loc)}
}

// ==================== 转换 ====================

// ToTime 返回底层的 time.Time
func (d DzDateTime) ToTime() time.Time {
	return d.t
}

// ToLunar 转换为农历日期
func (d DzDateTime) ToLunar() DzLunarDate {
	if d.err != nil {
		return DzLunarDate{err: d.err}
	}
	return LunarFromSolar(d)
}

// ToInterface 返回底层 time.Time 的 interface{} 形式
func (d DzDateTime) ToInterface() interface{} {
	if d.err != nil {
		return nil
	}
	return d.t
}

// String 实现 fmt.Stringer 接口
func (d DzDateTime) String() string {
	if d.err != nil {
		return fmt.Sprintf("DzDateTime{err: %v}", d.err)
	}
	return d.t.Format("2006-01-02 15:04:05")
}

// ==================== 内部辅助函数 ====================

// isLeapYear 判断是否闰年
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// daysInMonth 返回指定年月的天数
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// ==================== 静态工具函数 ====================

// DaysBetween 计算两个 time.Time 之间的天数差
func DaysBetween(t1, t2 time.Time) int {
	d := t2.Sub(t1).Hours() / 24
	return int(math.Round(d))
}

// IsLeapYear 判断指定年份是否闰年
func IsLeapYear(year int) bool {
	return isLeapYear(year)
}

// DaysInMonth 返回指定年月的天数
func DaysInMonth(year, month int) int {
	return daysInMonth(year, time.Month(month))
}

// DaysInYear 返回指定年份的天数
func DaysInYear(year int) int {
	if isLeapYear(year) {
		return 366
	}
	return 365
}

// BeginOfDay 返回当天开始时间
func BeginOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay 返回当天结束时间
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// BeginOfMonth 返回当月开始时间
func BeginOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth 返回当月结束时间
func EndOfMonth(t time.Time) time.Time {
	days := daysInMonth(t.Year(), t.Month())
	return time.Date(t.Year(), t.Month(), days, 23, 59, 59, 999999999, t.Location())
}

// BeginOfYear 返回当年开始时间
func BeginOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear 返回当年结束时间
func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}

// BeginOfWeek 返回本周开始时间（周一开始）
func BeginOfWeek(t time.Time) time.Time {
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	offset := 1 - wd
	return time.Date(t.Year(), t.Month(), t.Day()+offset, 0, 0, 0, 0, t.Location())
}

// EndOfWeek 返回本周结束时间（周日结束）
func EndOfWeek(t time.Time) time.Time {
	start := BeginOfWeek(t)
	return time.Date(start.Year(), start.Month(), start.Day()+6, 23, 59, 59, 999999999, start.Location())
}

// ==================== 协程安全的格式化/解析 ====================

// DzDateFormatter 提供协程安全的日期格式化和解析。
// 内部使用 sync.Pool 复用 buffer，避免频繁内存分配。
// 虽然 Go 的 time.Format/Parse 天然协程安全，DzDateFormatter 额外提供
// 预编译 layout 缓存和错误恢复能力。
type DzDateFormatter struct {
	layout string
	pool   sync.Pool
}

// NewDzDateFormatter 创建一个协程安全的日期格式化器
func NewDzDateFormatter(layout string) *DzDateFormatter {
	return &DzDateFormatter{
		layout: layout,
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 64)
			},
		},
	}
}

// Format 协程安全地格式化 time.Time
func (f *DzDateFormatter) Format(t time.Time) string {
	return t.Format(f.layout)
}

// Parse 协程安全地解析日期字符串
func (f *DzDateFormatter) Parse(value string) (time.Time, error) {
	return time.Parse(f.layout, value)
}

// ParseInLocation 协程安全地按指定时区解析日期字符串
func (f *DzDateFormatter) ParseInLocation(value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(f.layout, value, loc)
}

// FormatDzDateTime 协程安全地格式化 DzDateTime
func (f *DzDateFormatter) FormatDzDateTime(d DzDateTime) string {
	if d.err != nil {
		return ""
	}
	return d.t.Format(f.layout)
}

// Layout 返回格式化器的 layout
func (f *DzDateFormatter) Layout() string {
	return f.layout
}

// FormatSafe 安全格式化，不会 panic，出错返回空字符串
func FormatSafe(t time.Time, layout string) string {
	defer func() {
		recover()
	}()
	return t.Format(layout)
}

// ParseSafe 安全解析，不会 panic
func ParseSafe(layout, value string) (time.Time, error) {
	defer func() {
		recover()
	}()
	return time.Parse(layout, value)
}
