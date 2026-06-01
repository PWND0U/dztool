package DateTool

import (
	"fmt"
	"testing"
	"time"
)

// ==================== 构造函数测试 ====================

func TestNewDzDateTime(t *testing.T) {
	now := time.Now()
	dt := NewDzDateTime(now)
	fmt.Println("NewDzDateTime IsValid:", dt.IsValid())
	fmt.Println("NewDzDateTime Year:", dt.Year())
	fmt.Println("NewDzDateTime String:", dt.String())
}

func TestNow(t *testing.T) {
	dt := Now()
	fmt.Println("Now IsValid:", dt.IsValid())
	fmt.Println("Now String:", dt.String())
}

func TestDzDate(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	fmt.Println("DzDate String:", dt.String())
	fmt.Println("DzDate Hour:", dt.Hour())
	fmt.Println("DzDate Minute:", dt.Minute())
	fmt.Println("DzDate Second:", dt.Second())
}

func TestDzDateTimeOf(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("DzDateTimeOf String:", dt.String())
	fmt.Println("DzDateTimeOf Hour:", dt.Hour())
	fmt.Println("DzDateTimeOf Minute:", dt.Minute())
	fmt.Println("DzDateTimeOf Second:", dt.Second())
}

func TestParse(t *testing.T) {
	dt := Parse("2006-01-02 15:04:05", "2024-06-15 10:30:45")
	fmt.Println("Parse IsValid:", dt.IsValid())
	fmt.Println("Parse String:", dt.String())

	dtErr := Parse("2006-01-02", "invalid")
	fmt.Println("Parse invalid IsValid:", dtErr.IsValid())
	fmt.Println("Parse invalid Err:", dtErr.Err())
}

func TestFromTimestamp(t *testing.T) {
	ts := int64(1718439045)
	dt := FromTimestamp(ts)
	fmt.Println("FromTimestamp String:", dt.String())
	fmt.Println("FromTimestamp Timestamp:", dt.Timestamp())
}

func TestFromTimestampMilli(t *testing.T) {
	ts := int64(1718439045000)
	dt := FromTimestampMilli(ts)
	fmt.Println("FromTimestampMilli String:", dt.String())
	fmt.Println("FromTimestampMilli TimestampMilli:", dt.TimestampMilli())
}

// ==================== 链式 Setter 测试 ====================

func TestWithYear(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.WithYear(2020)
	fmt.Println("WithYear original:", dt.String())
	fmt.Println("WithYear new:", dt2.String())
	fmt.Println("WithYear Year:", dt2.Year())
}

func TestWithMonth(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.WithMonth(1)
	fmt.Println("WithMonth new:", dt2.String())
	fmt.Println("WithMonth Month:", dt2.Month())
}

func TestWithDay(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.WithDay(1)
	fmt.Println("WithDay new:", dt2.String())
	fmt.Println("WithDay Day:", dt2.Day())
}

func TestWithHour(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	dt2 := dt.WithHour(20)
	fmt.Println("WithHour new:", dt2.String())
	fmt.Println("WithHour Hour:", dt2.Hour())
}

func TestWithMinute(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	dt2 := dt.WithMinute(0)
	fmt.Println("WithMinute new:", dt2.String())
	fmt.Println("WithMinute Minute:", dt2.Minute())
}

func TestWithSecond(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	dt2 := dt.WithSecond(0)
	fmt.Println("WithSecond new:", dt2.String())
	fmt.Println("WithSecond Second:", dt2.Second())
}

func TestChainedWith(t *testing.T) {
	dt := DzDate(2024, 6, 15).
		WithYear(2020).
		WithMonth(1).
		WithDay(1).
		WithHour(12).
		WithMinute(0).
		WithSecond(0)
	fmt.Println("ChainedWith:", dt.String())
}

// ==================== 链式算术测试 ====================

func TestAddYears(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.AddYears(1)
	fmt.Println("AddYears:", dt2.String())
}

func TestAddMonths(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.AddMonths(3)
	fmt.Println("AddMonths:", dt2.String())
}

func TestAddDays(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.AddDays(10)
	fmt.Println("AddDays:", dt2.String())
}

func TestSubYears(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.SubYears(1)
	fmt.Println("SubYears:", dt2.String())
}

func TestSubMonths(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.SubMonths(3)
	fmt.Println("SubMonths:", dt2.String())
}

func TestSubDays(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	dt2 := dt.SubDays(10)
	fmt.Println("SubDays:", dt2.String())
}

func TestAddHours(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 0, 0)
	dt2 := dt.AddHours(5)
	fmt.Println("AddHours:", dt2.String())
}

func TestAddMinutes(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 0, 0)
	dt2 := dt.AddMinutes(30)
	fmt.Println("AddMinutes:", dt2.String())
}

func TestAddSeconds(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 0, 0)
	dt2 := dt.AddSeconds(45)
	fmt.Println("AddSeconds:", dt2.String())
}

func TestAddDuration(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 0, 0)
	dt2 := dt.AddDuration(2*time.Hour + 30*time.Minute)
	fmt.Println("AddDuration:", dt2.String())
}

func TestChainedArithmetic(t *testing.T) {
	dt := DzDate(2024, 6, 15).
		AddYears(1).
		AddMonths(2).
		AddDays(10).
		AddHours(5).
		AddMinutes(30)
	fmt.Println("ChainedArithmetic:", dt.String())
}

// ==================== 边界方法测试 ====================

func TestStartEndOfDay(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("StartOfDay:", dt.StartOfDay().String())
	fmt.Println("EndOfDay:", dt.EndOfDay().String())
}

func TestStartEndOfHour(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("StartOfHour:", dt.StartOfHour().String())
	fmt.Println("EndOfHour:", dt.EndOfHour().String())
}

func TestStartEndOfMinute(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("StartOfMinute:", dt.StartOfMinute().String())
	fmt.Println("EndOfMinute:", dt.EndOfMinute().String())
}

func TestStartEndOfMonth(t *testing.T) {
	dt := DzDateTimeOf(2024, 2, 15, 10, 30, 45)
	fmt.Println("StartOfMonth:", dt.StartOfMonth().String())
	fmt.Println("EndOfMonth:", dt.EndOfMonth().String())
	fmt.Println("Feb DaysInMonth:", dt.DaysInMonth())

	dt2 := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("Jun EndOfMonth:", dt2.EndOfMonth().String())
}

func TestStartEndOfYear(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("StartOfYear:", dt.StartOfYear().String())
	fmt.Println("EndOfYear:", dt.EndOfYear().String())
	fmt.Println("DaysInYear:", dt.DaysInYear())
}

func TestStartEndOfWeek(t *testing.T) {
	// 2024-06-15 是周六
	dt := DzDate(2024, 6, 15)
	fmt.Println("Weekday:", dt.Weekday())
	fmt.Println("StartOfWeek:", dt.StartOfWeek().String())
	fmt.Println("EndOfWeek:", dt.EndOfWeek().String())
}

// ==================== 比较方法测试 ====================

func TestIsBefore(t *testing.T) {
	dt1 := DzDate(2024, 6, 15)
	dt2 := DzDate(2024, 6, 16)
	fmt.Println("IsBefore:", dt1.IsBefore(dt2))
	fmt.Println("IsBefore (reverse):", dt2.IsBefore(dt1))
}

func TestIsAfter(t *testing.T) {
	dt1 := DzDate(2024, 6, 15)
	dt2 := DzDate(2024, 6, 16)
	fmt.Println("IsAfter:", dt2.IsAfter(dt1))
}

func TestIsEqual(t *testing.T) {
	dt1 := DzDate(2024, 6, 15)
	dt2 := DzDate(2024, 6, 15)
	fmt.Println("IsEqual (same):", dt1.IsEqual(dt2))

	dt3 := DzDate(2024, 6, 16)
	fmt.Println("IsEqual (diff):", dt1.IsEqual(dt3))
}

func TestIsBetween(t *testing.T) {
	dt := DzDate(2024, 6, 15)
	start := DzDate(2024, 6, 1)
	end := DzDate(2024, 6, 30)
	fmt.Println("IsBetween:", dt.IsBetween(start, end))

	outside := DzDate(2024, 7, 15)
	fmt.Println("IsBetween (outside):", outside.IsBetween(start, end))
}

func TestIsToday(t *testing.T) {
	dt := Now()
	fmt.Println("IsToday (now):", dt.IsToday())
}

func TestIsLeapYear(t *testing.T) {
	dt1 := DzDate(2024, 1, 1)
	fmt.Println("2024 IsLeapYear:", dt1.IsLeapYear())

	dt2 := DzDate(2023, 1, 1)
	fmt.Println("2023 IsLeapYear:", dt2.IsLeapYear())
}

func TestIsWeekend(t *testing.T) {
	// 2024-06-15 是周六
	dt := DzDate(2024, 6, 15)
	fmt.Println("2024-06-15 IsWeekend:", dt.IsWeekend())
	fmt.Println("2024-06-15 IsWeekday:", dt.IsWeekday())

	// 2024-06-17 是周一
	dt2 := DzDate(2024, 6, 17)
	fmt.Println("2024-06-17 IsWeekend:", dt2.IsWeekend())
	fmt.Println("2024-06-17 IsWeekday:", dt2.IsWeekday())
}

func TestIsSameDayMonthYear(t *testing.T) {
	dt1 := DzDateTimeOf(2024, 6, 15, 10, 0, 0)
	dt2 := DzDateTimeOf(2024, 6, 15, 20, 0, 0)
	dt3 := DzDateTimeOf(2024, 7, 15, 10, 0, 0)
	dt4 := DzDateTimeOf(2025, 6, 15, 10, 0, 0)

	fmt.Println("IsSameDay:", dt1.IsSameDay(dt2))
	fmt.Println("IsSameDay (diff):", dt1.IsSameDay(dt3))
	fmt.Println("IsSameMonth:", dt1.IsSameMonth(dt2))
	fmt.Println("IsSameMonth (diff month):", dt1.IsSameMonth(dt3))
	fmt.Println("IsSameYear:", dt1.IsSameYear(dt3))
	fmt.Println("IsSameYear (diff year):", dt1.IsSameYear(dt4))
}

// ==================== Getter 测试 ====================

func TestGetter(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("Year:", dt.Year())
	fmt.Println("Month:", dt.Month())
	fmt.Println("Day:", dt.Day())
	fmt.Println("Hour:", dt.Hour())
	fmt.Println("Minute:", dt.Minute())
	fmt.Println("Second:", dt.Second())
	fmt.Println("Weekday:", dt.Weekday())
	fmt.Println("YearDay:", dt.YearDay())
	isoYear, isoWeek := dt.ISOWeek()
	fmt.Println("ISOWeek:", isoYear, isoWeek)
	fmt.Println("Timestamp:", dt.Timestamp())
	fmt.Println("TimestampMilli:", dt.TimestampMilli())
	fmt.Println("DaysInMonth:", dt.DaysInMonth())
	fmt.Println("DaysInYear:", dt.DaysInYear())
}

// ==================== 差量测试 ====================

func TestDiffInYears(t *testing.T) {
	dt1 := DzDate(2024, 6, 15)
	dt2 := DzDate(2020, 6, 15)
	fmt.Println("DiffInYears:", dt1.DiffInYears(dt2))
}

func TestDiffInMonths(t *testing.T) {
	dt1 := DzDate(2024, 6, 15)
	dt2 := DzDate(2023, 1, 15)
	fmt.Println("DiffInMonths:", dt1.DiffInMonths(dt2))
}

func TestDiffInDays(t *testing.T) {
	dt1 := DzDate(2024, 6, 15)
	dt2 := DzDate(2024, 6, 1)
	fmt.Println("DiffInDays:", dt1.DiffInDays(dt2))
}

func TestDiffInHours(t *testing.T) {
	dt1 := DzDateTimeOf(2024, 6, 15, 20, 0, 0)
	dt2 := DzDateTimeOf(2024, 6, 15, 10, 0, 0)
	fmt.Println("DiffInHours:", dt1.DiffInHours(dt2))
}

func TestDiffInMinutes(t *testing.T) {
	dt1 := DzDateTimeOf(2024, 6, 15, 10, 30, 0)
	dt2 := DzDateTimeOf(2024, 6, 15, 10, 0, 0)
	fmt.Println("DiffInMinutes:", dt1.DiffInMinutes(dt2))
}

func TestDiffInSeconds(t *testing.T) {
	dt1 := DzDateTimeOf(2024, 6, 15, 10, 0, 45)
	dt2 := DzDateTimeOf(2024, 6, 15, 10, 0, 0)
	fmt.Println("DiffInSeconds:", dt1.DiffInSeconds(dt2))
}

// ==================== 格式化测试 ====================

func TestFormat(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("Format:", dt.Format("2006-01-02 15:04:05"))
	fmt.Println("ToDateString:", dt.ToDateString())
	fmt.Println("ToTimeString:", dt.ToTimeString())
	fmt.Println("ToDateTimeString:", dt.ToDateTimeString())
	fmt.Println("ToRfc3339String:", dt.ToRfc3339String())
	fmt.Println("ToRfc1123String:", dt.ToRfc1123String())
	fmt.Println("ToAtomString:", dt.ToAtomString())
}

// ==================== 时区测试 ====================

func TestTimezone(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	fmt.Println("Original:", dt.String())
	fmt.Println("UTC:", dt.UTC().String())
	fmt.Println("Local:", dt.Local().String())

	dtShanghai := dt.ToTimezone("Asia/Shanghai")
	fmt.Println("Shanghai:", dtShanghai.String())
	fmt.Println("Shanghai IsValid:", dtShanghai.IsValid())

	dtInvalid := dt.ToTimezone("Invalid/Zone")
	fmt.Println("Invalid Timezone IsValid:", dtInvalid.IsValid())
	fmt.Println("Invalid Timezone Err:", dtInvalid.Err())
}

// ==================== 转换测试 ====================

func TestToTime(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	tm := dt.ToTime()
	fmt.Println("ToTime:", tm.Format("2006-01-02 15:04:05"))
}

func TestToInterface(t *testing.T) {
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)
	v := dt.ToInterface()
	fmt.Println("ToInterface type:", fmt.Sprintf("%T", v))
}

// ==================== 错误短路测试 ====================

func TestErrShortCircuit(t *testing.T) {
	dt := Parse("invalid-layout", "not-a-date")
	fmt.Println("ErrShortCircuit IsValid:", dt.IsValid())
	// 链式调用应在错误上短路
	dt2 := dt.AddYears(1).AddMonths(2).AddDays(3)
	fmt.Println("ErrShortCircuit after chain IsValid:", dt2.IsValid())
	fmt.Println("ErrShortCircuit Err:", dt2.Err())
}

// ==================== 静态工具函数测试 ====================

func TestStaticIsLeapYear(t *testing.T) {
	fmt.Println("IsLeapYear(2024):", IsLeapYear(2024))
	fmt.Println("IsLeapYear(2023):", IsLeapYear(2023))
	fmt.Println("IsLeapYear(2000):", IsLeapYear(2000))
	fmt.Println("IsLeapYear(1900):", IsLeapYear(1900))
}

func TestStaticDaysInMonth(t *testing.T) {
	fmt.Println("DaysInMonth(2024, 2):", DaysInMonth(2024, 2))
	fmt.Println("DaysInMonth(2023, 2):", DaysInMonth(2023, 2))
	fmt.Println("DaysInMonth(2024, 6):", DaysInMonth(2024, 6))
}

func TestStaticDaysInYear(t *testing.T) {
	fmt.Println("DaysInYear(2024):", DaysInYear(2024))
	fmt.Println("DaysInYear(2023):", DaysInYear(2023))
}

func TestStaticDaysBetween(t *testing.T) {
	t1 := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)
	fmt.Println("DaysBetween:", DaysBetween(t1, t2))
}

func TestStaticBeginEndOfDay(t *testing.T) {
	tm := time.Date(2024, 6, 15, 10, 30, 45, 0, time.Local)
	fmt.Println("BeginOfDay:", BeginOfDay(tm).Format("2006-01-02 15:04:05"))
	fmt.Println("EndOfDay:", EndOfDay(tm).Format("2006-01-02 15:04:05"))
}

func TestStaticBeginEndOfMonth(t *testing.T) {
	tm := time.Date(2024, 2, 15, 10, 30, 45, 0, time.Local)
	fmt.Println("BeginOfMonth:", BeginOfMonth(tm).Format("2006-01-02 15:04:05"))
	fmt.Println("EndOfMonth:", EndOfMonth(tm).Format("2006-01-02 15:04:05"))
}

func TestStaticBeginEndOfYear(t *testing.T) {
	tm := time.Date(2024, 6, 15, 10, 30, 45, 0, time.Local)
	fmt.Println("BeginOfYear:", BeginOfYear(tm).Format("2006-01-02 15:04:05"))
	fmt.Println("EndOfYear:", EndOfYear(tm).Format("2006-01-02 15:04:05"))
}

func TestStaticBeginEndOfWeek(t *testing.T) {
	tm := time.Date(2024, 6, 15, 10, 30, 45, 0, time.Local) // 周六
	fmt.Println("BeginOfWeek:", BeginOfWeek(tm).Format("2006-01-02 15:04:05"))
	fmt.Println("EndOfWeek:", EndOfWeek(tm).Format("2006-01-02 15:04:05"))
}

// ==================== DzDateFormatter 测试 ====================

func TestDzDateFormatter(t *testing.T) {
	f := NewDzDateFormatter("2006-01-02 15:04:05")
	dt := DzDateTimeOf(2024, 6, 15, 10, 30, 45)

	fmt.Println("DzDateFormatter Format:", f.Format(dt.ToTime()))
	fmt.Println("DzDateFormatter FormatDzDateTime:", f.FormatDzDateTime(dt))

	parsed, err := f.Parse("2024-06-15 10:30:45")
	fmt.Println("DzDateFormatter Parse:", parsed.Format("2006-01-02 15:04:05"), "err:", err)

	fmt.Println("DzDateFormatter Layout:", f.Layout())
}

func TestFormatSafe(t *testing.T) {
	tm := time.Date(2024, 6, 15, 10, 30, 45, 0, time.UTC)
	fmt.Println("FormatSafe:", FormatSafe(tm, "2006-01-02 15:04:05"))
}

func TestParseSafe(t *testing.T) {
	tm, err := ParseSafe("2006-01-02", "2024-06-15")
	fmt.Println("ParseSafe:", tm.Format("2006-01-02"), "err:", err)
}
