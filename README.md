# dztool

一个 Go 工具库，提供字符串、字节、JSON、结构体、时间、IO、网络、SSE、JWT、OCR 等常用操作的链式封装类型。

## 安装

```bash
go get github.com/PWND0U/dztool
```

## 构建 & 测试

```bash
go build ./...
go test ./...
go vet ./...
```

---

## 包总览

| 包 | 核心类型 | 说明 |
|---|---|---|
| StringTool | `DzString` `DzStrings` | 字符串链式操作 |
| BytesTool | `DzBytes` | 字节切片链式操作 |
| JsonTool | `DzJsonMap` `DzJsonArray` `DzJsonMapArray` | JSON 点号路径导航 + JSON 修补 |
| StructTool | `DzStruct` | 结构体拷贝/克隆/比较/Map互转（类似 Java BeanUtils），全面支持指针类型 |
| DateTool | `DzDateTime` `DzLunarDate` `DzStopwatch` | Joda-Time 风格日期时间对象 + 农历 + 计时器 |
| TimeIntervalTool | `DzTimeInterval` | 命名计时器 |
| IOTool | `DzFastBuffer` `DzWatcher` | 流操作、文件读写、文件类型检测、目录监听、快速缓冲、路径工具 |
| NetTool | — | IP 转换、端口检测、网卡信息、Ping、CIDR、Cookie 解析 |
| ServerTool | `DzServerSentEvent` | SSE 编解码与 HTTP Flush |
| Algorithm | — | Levenshtein 编辑距离 |
| Algorithm/DzJWT | `JWTClaims` | JWT 创建、解析与验证（HS256） |
| dzAi/dzOcr | `DzPaddleOcrEngine` | PaddleOCR 文字识别引擎（CPU/CUDA/DirectML/TensorRT） |

---

## StringTool

```go
import "github.com/PWND0U/dztool/StringTool"
```

### DzString

字符串链式封装，底层类型 `string`。

```go
// 构造
ds := StringTool.NewDzString("hello world")

// 转换
ds.ToString()          // string
ds.ToBytes()           // []byte
ds.ToInt()             // int（失败返回 0）

// 替换
ds.ReplaceAll("o", "0")      // "hell0 w0rld"
ds.ReplaceN("o", "0", 1)     // "hell0 world"
ds.RexReplaceAll(`\d+`, "#")  // 正则替换

// 大小写
ds.Upper()             // "HELLO WORLD"
ds.Lower()             // "hello world"
ds.Title()             // "Hello World"

// 查找与分割
ds.Find("world")              // 6
ds.Split(" ")                 // DzStrings{"hello", "world"}
ds.SplitN(" ", 2)             // DzStrings{"hello", "world"}
ds.IsContains("world")        // true
ds.IsEmpty()                  // false

// 去除
ds.Strip("hd")         // "ello worl"
ds.LStrip("h")         // "ello world"
ds.RStrip("d")         // "hello worl"

// 相似度（基于 Levenshtein）
ds.SimilarText("hello")  // 0.69...

// 模板格式化
StringTool.NewDzString("{name} is {age}").FStringFormat(map[string]any{"name": "Alice", "age": "30"})
// "Alice is 30"

// 静态 Join
StringTool.Join([]string{"a", "b"}, "-")  // DzString("a-b")
```

### DzStrings

字符串切片封装，底层类型 `[]DzString`。

```go
dss := StringTool.NewDzStrings([]string{"a", "b", "c"})
dss.Join("-")           // DzString("a-b-c")
dss.IsContain("b")      // true
```

---

## BytesTool

```go
import "github.com/PWND0U/dztool/BytesTool"
```

### DzBytes

字节切片链式封装，底层类型 `[]byte`。

```go
// 构造
db := BytesTool.NewDzBytes("hello")          // 从字符串
db2 := BytesTool.NewDzBytesByBytes([]byte{1,2,3})  // 从 []byte

// 转换
db.ToString()           // "hello"
db.ToDzString()         // StringTool.DzString

// 替换
db.ReplaceAll([]byte("l"), []byte("L"))  // "heLLo"
db.ReplaceN([]byte("l"), []byte("L"), 1) // "heLlo"

// 查找与分割
db.Find([]byte("ll"))           // 2
db.Split([]byte("l"))           // []DzBytes{"he", "", "o"}
db.SplitN([]byte("l"), 2)       // []DzBytes{"he", "o"}
db.IsContains([]byte("ell"))     // true

// 静态 Join
BytesTool.Join([][]byte{[]byte("a"), []byte("b")}, []byte("-"))  // DzBytes("a-b")
```

---

## JsonTool

```go
import "github.com/PWND0U/dztool/JsonTool"
```

### DzJsonMap

JSON 对象封装，底层类型 `map[string]interface{}`，支持点号路径访问嵌套字段。

```go
// 构造
djm := JsonTool.NewDzJsonMap([]byte(`{"commodity":{"data":[{"name":"苹果"}]}}`))
djm2 := JsonTool.ParseDzJsonMap(map[string]interface{}{"key": "val"})

// 点号路径访问
djm.GetString("commodity.data[0].name")  // 不支持索引，需配合 GetArray
djm.GetArray("commodity.data").GetMapArray()[0].GetString("name")  // "苹果"
djm.GetMap("commodity")           // DzJsonMap
djm.GetMapArray("commodity.data") // []DzJsonMap
djm.GetInt("count")               // int（缺失返回 0）
djm.GetFloat("price")             // float64（缺失返回 0）
djm.GetBool("active")             // bool（缺失返回 false）
```

### DzJsonArray

JSON 数组封装，底层类型 `[]interface{}`，提供类型化提取。

```go
arr := djm.GetArray("items")
arr.GetStringArray()   // []string
arr.GetIntArray()      // []int
arr.GetFloatArray()    // []float64
arr.GetBoolArray()     // []bool
arr.GetMapArray()      // []DzJsonMap
```

### DzJsonMapArray

`DzJsonMap` 切片封装，底层类型 `[]DzJsonMap`。

```go
djma := JsonTool.NewDzJsonMapArray([]byte(`[{"name":"a"},{"name":"b"}]`))
for _, m := range djma {
    fmt.Println(m.GetString("name"))
}
```

### JSON 修补

使用 [jsonrepair](https://github.com/kaptinlin/jsonrepair) 修复不合法的 JSON 字符串，适用于 LLM 输出、JavaScript 代码片段等场景。

```go
// 核心修补函数
repaired, err := JsonTool.RepairJSON(`{name: John, age: 30}`)
// repaired = `{"name": "John", "age": 30}`

// 修补失败时 panic
result := JsonTool.MustRepairJSON(`{'a':'foo'}`)

// 修补并解析为 DzJsonMap（失败返回 nil）
dzm := JsonTool.RepairToDzJsonMap(`{name: "Alice", active: True}`)

// 修补并解析为原生 map
m, err := JsonTool.RepairToJsonMap(`{name: "Bob"}`)

// 修补并解析到结构体
var p Person
err = JsonTool.RepairToStruct(`{name: "Charlie", age: 30, active: True}`, &p)

// 判断是否可修补
JsonTool.IsRepairable(`{a:1}`)    // true
JsonTool.IsRepairable(`{{{`)      // false

// 尝试修补：返回结果 + 是否实际修补了 + 错误
result, repaired, err := JsonTool.TryRepair(text)
// 如果本身是合法 JSON，repaired=false，直接返回原文
```

**可修复的常见问题**：

| 问题 | 输入 | 输出 |
|------|------|------|
| 缺少引号 | `{name: John}` | `{"name": "John"}` |
| 单引号 | `{'a':'foo'}` | `{"a":"foo"}` |
| 尾部逗号 | `{"a":1,}` | `{"a":1}` |
| 缺少闭合括号 | `{"items":[1,2,3` | `{"items":[1,2,3]}` |
| JS 注释 | `{"a":1/*c*/}` | `{"a":1}` |
| Python 常量 | `{active: True}` | `{"active":true}` |
| JSONP 包装 | `callback({"ok":1});` | `{"ok":1}` |
| 截断 JSON | `{"msg":"hello` | `{"msg":"hello"}` |
| 省略号 | `[1,2,...]` | `[1,2]` |
| NDJSON | `{"id":1}\n{"id":2}` | `[{"id":1},{"id":2}]` |

---

## StructTool

```go
import "github.com/PWND0U/dztool/StructTool"
```

类似 Java BeanUtils 的结构体工具，支持相同/不同结构体间拷贝、合并、比较，以及结构体与 Map 的互转。

### DzStruct

```go
type Person struct {
    Name string
    Age  int
}
```

#### 构造与状态

```go
ds := StructTool.NewDzStruct(Person{Name: "Alice", Age: 30})
ds.IsValid()     // true
ds.IsZero()      // false
ds.Err()         // nil

// 传入指针同样支持
ds2 := StructTool.NewDzStruct(&Person{Name: "Bob", Age: 25})

// 错误情况
StructTool.NewDzStruct(nil)        // err != nil
StructTool.NewDzStruct("string")   // err != nil（非结构体）
```

#### 值提取

```go
ds.ToInterface()                    // Person{Name:"Alice", Age:30}
var p Person
ds.ToIntf(&p)                       // 将值填充到 p
```

#### 同结构体操作

```go
// ============ 深克隆 ============
// 递归复制所有层级数据，修改克隆体不影响原始结构体

// DeepClone — JSON 方式（适用于可 JSON 序列化的字段）
cloned := ds.DeepClone()

// DeepCloneReflect — 反射方式（不依赖 JSON，支持不可序列化类型）
cloned2 := ds.DeepCloneReflect()

// ============ 浅克隆 ============
// 仅复制值类型字段，引用类型（slice, map, pointer）与原始共享底层数据
// 修改克隆体的切片或 map 元素会影响原始结构体

shallow := ds.ShallowClone()

// ============ 字段操作 ============
ds.Fields()                          // ["Name", "Age"]
ds.Field("Name")                     // "Alice"
ds.SetField("Name", "Bob")           // 修改字段（支持链式）
ds.SetField("Age", "30")             // cast 类型转换："30" → 30
ds.SetField("Ptr", nil)              // 指针字段设为 nil
ds.SetField("Inner", otherStruct)    // 指针结构体字段自动递归复制
ds.Zero()                            // 所有字段置零值
```

#### 跨结构体操作

```go
type Employee struct {
    Name   string
    Age    int
    Salary float64
}

emp := Employee{Name: "Charlie", Age: 28, Salary: 5000}

// CopyFrom — 按字段名匹配复制（不同结构体间，cast 类型转换）
ds.CopyFrom(emp)                     // Name="Charlie", Age=28

// CopyFromByTag — 按 struct tag 匹配复制
type UserDTO struct {
    UserName string `json:"user_name"`
    UserAge  int    `json:"user_age"`
}
type User struct {
    UserName string `json:"user_name"`
    UserAge  int    `json:"user_age"`
    Address  string `json:"address"`
}
dto := UserDTO{UserName: "dto_name", UserAge: 25}
u := StructTool.NewDzStruct(User{}).CopyFromByTag(dto, "json")
// UserName="dto_name", UserAge=25, Address=""

// MergeFrom — 仅覆盖源非零字段（零值字段不覆盖目标）
p1 := StructTool.NewDzStruct(Person{Name: "Alice", Age: 30})
p2 := Person{Name: "Bob", Age: 0}    // Age 为零值
p1.MergeFrom(p2)                     // Name="Bob", Age=30（Age 未被覆盖）

// MergeFromByTag — 按 tag 合并非零字段

// 比较
diff := ds.CompareTo(other)          // map[Name:[Alice Bob]] — 仅包含差异字段
diffFields := ds.DiffFields(other)   // ["Name"]
equal := ds.EqualTo(other)           // false
```

#### 结构体与 Map 互转

```go
// 结构体 → Map
ds.ToMap()                           // map[Name:Alice Age:30]
ds.ToMapByTag("json")                // 按 json tag 作为 key

// 嵌套结构体递归展开
type Nested struct {
    Name    string
    Address struct { City string }
}
Nested{Name:"Alice", Address:struct{City string}{"Beijing"}}.ToMap()
// map[Name:Alice Address:map[City:Beijing]]

// 嵌入字段扁平化
type Inner struct { InnerName string }
type Outer struct { OuterName string; Inner }
Outer{OuterName:"O", Inner:Inner{InnerName:"I"}}.ToMap()
// map[OuterName:O InnerName:I]

// Map → 结构体
ds2 := StructTool.NewDzStruct(Person{}).FromMap(map[string]interface{}{
    "Name": "Charlie",
    "Age":  "28",                     // 字符串 "28" 自动 cast 为 int
})
ds2.FromMapByTag(data, "json")        // 按 tag 匹配

// Round-trip
ds.ToMap() → FromMap() → EqualTo(ds)  // true
```

#### 静态工具函数

无需创建 `DzStruct`，直接操作：

```go
// 拷贝
StructTool.CopyStruct(src, &dst)                         // 按名复制
StructTool.CopyStructByTag(src, &dst, "json")            // 按 tag 复制

// 结构体 ↔ Map
StructTool.StructToMap(src)                               // (map[string]interface{}, error)
StructTool.StructToMapByTag(src, "json")                  // (map[string]interface{}, error)
StructTool.MapToStruct(data, &output)                     // error
StructTool.MapToStructByTag(data, &output, "json")        // error

// 深克隆（递归复制所有层级，修改克隆体不影响原始）
StructTool.DeepCloneStruct(src)                            // JSON 方式
StructTool.DeepCloneStructReflect(src)                     // 反射方式

// 浅克隆（引用类型字段与原始共享底层数据）
StructTool.ShallowCloneStruct(src)

// 比较
StructTool.CompareStruct(a, b)                            // ([]string, error)
```

#### 指针类型支持

`DzStruct` 全面支持指针类型字段，包括指针结构体、双重指针、指针匿名嵌入等：

```go
type Inner struct { Name string; ID int }
type WithPtr struct {
    Name  string
    Inner *Inner     // 指针结构体字段
}

// SetField 支持多种赋值方式
ds := StructTool.NewDzStruct(WithPtr{})
ds.SetField("Inner", &Inner{Name: "test", ID: 1})  // 传指针
ds.SetField("Inner", Inner{Name: "test", ID: 2})   // 传值（自动创建指针）
ds.SetField("Inner", nil)                           // 设为 nil

// ToMap / FromMap 支持指针结构体
m := ds.ToMap()          // map[Name: Inner:map[Name:test ID:1]]
ds2.FromMap(map[string]interface{}{
    "Inner": map[string]interface{}{"Name": "from-map", "ID": 3},
})
ds2.FromMap(map[string]interface{}{"Inner": nil})  // 指针设为 nil

// CopyFrom / MergeFrom 支持不同类型的指针结构体（递归匹配字段名）
type InnerDTO struct { Name string `json:"name"`; ID int `json:"id"` }
ds.CopyFrom(WithPtrDTO{Inner: &InnerDTO{Name: "cross", ID: 5}})

// 深克隆正确复制指针结构体（修改克隆体不影响原始）
cloned := ds.DeepClone()

// 往返一致：ToMap → FromMap → EqualTo 原始
```

---

## DateTool

```go
import "github.com/PWND0U/dztool/DateTool"
```

Joda-Time 风格的日期时间操作库，提供三大模块：

- **DzDateTime** — 不可变链式日期时间对象
- **DzLunarDate** — 农历日期（天干地支、生肖、公历农历互转）
- **DzStopwatch** — 代码执行计时器

### DzDateTime

值类型，所有链式方法返回新实例，天然协程安全。内嵌 `err` 字段支持链式错误短路。

```go
// ============ 构造 ============
dt := DateTool.Now()                                        // 当前时间
dt = DateTool.DzDate(2024, 6, 15)                          // 仅日期
dt = DateTool.DzDateTimeOf(2024, 6, 15, 10, 30, 0)        // 日期+时间
dt = DateTool.Parse("2006-01-02", "2024-06-15")            // 解析字符串
dt = DateTool.FromTimestamp(1718439045)                     // Unix 时间戳
dt = DateTool.FromTimestampMilli(1718439045000)             // 毫秒时间戳

// ============ 链式 Setter（返回新值，不修改原对象）============
dt.WithYear(2020).WithMonth(1).WithDay(1).WithHour(12)

// ============ 链式算术 ============
dt.AddYears(1).AddMonths(3).AddDays(10)
dt.SubHours(5).SubMinutes(30).SubSeconds(15)
dt.AddDuration(2 * time.Hour)

// ============ 边界方法 ============
dt.StartOfDay()     // 当天 00:00:00
dt.EndOfDay()       // 当天 23:59:59.999999999
dt.StartOfMonth()   // 当月1日 00:00:00
dt.EndOfMonth()     // 当月最后一天 23:59:59
dt.StartOfYear()    // 当年1月1日
dt.EndOfYear()      // 当年12月31日 23:59:59
dt.StartOfWeek()    // 本周一 00:00:00
dt.EndOfWeek()      // 本周日 23:59:59

// ============ 比较 ============
dt.IsBefore(other)      // bool
dt.IsAfter(other)       // bool
dt.IsEqual(other)       // bool
dt.IsBetween(start, end) // bool（包含边界）
dt.IsToday()            // bool
dt.IsLeapYear()         // bool
dt.IsWeekend()          // bool
dt.IsSameDay(other)     // bool

// ============ 差量 ============
dt.DiffInYears(other)    // int
dt.DiffInMonths(other)   // int
dt.DiffInDays(other)     // int
dt.DiffInHours(other)    // float64

// ============ 格式化 ============
dt.Format("2006-01-02 15:04:05")
dt.ToDateString()           // "2024-06-15"
dt.ToDateTimeString()       // "2024-06-15 10:30:00"
dt.ToRfc3339String()        // RFC3339 格式

// ============ 时区 ============
dt.UTC()                    // UTC
dt.ToTimezone("Asia/Shanghai")

// ============ Getter ============
dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), dt.Second()
dt.Timestamp(), dt.TimestampMilli(), dt.DaysInMonth(), dt.DaysInYear()

// ============ 错误检查 ============
dt.Err()        // error（链式调用中累积）
dt.IsValid()    // bool

// ============ 转换 ============
dt.ToTime()     // time.Time
dt.ToLunar()    // DzLunarDate（公历转农历）
```

#### 静态工具函数

```go
DateTool.IsLeapYear(2024)                    // true
DateTool.DaysInMonth(2024, 2)                // 29
DateTool.DaysBetween(t1, t2)                 // int
DateTool.BeginOfDay(t)                       // time.Time
DateTool.EndOfMonth(t)                       // time.Time
DateTool.FormatSafe(t, "2006-01-02")         // panic-free 格式化
```

#### DzDateFormatter — 协程安全格式化器

```go
f := DateTool.NewDzDateFormatter("2006-01-02 15:04:05")
f.Format(time.Now())                  // string
f.Parse("2024-06-15 10:30:00")        // (time.Time, error)
f.FormatDzDateTime(dt)                // string
```

### DzLunarDate

农历日期封装，覆盖 1900-2100 年，支持公历农历互转、天干地支、生肖、中文命名。

```go
// ============ 构造 ============
lunar := DateTool.NewDzLunarDate(2024, 6, 15)             // 农历日期
lunar = DateTool.NewDzLunarDate(2023, 2, 15, true)        // 闰月
lunar = DateTool.LunarFromSolar(dt)                        // 公历转农历

// ============ 公历↔农历互转 ============
solar := lunar.ToSolar()          // DzDateTime（农历转公历）
lunar := dt.ToLunar()             // DzLunarDate（公历转农历）

// 往返一致
orig := DateTool.DzDate(2024, 6, 15)
back := DateTool.LunarFromSolar(orig).ToSolar()
orig.ToDateString() == back.ToDateString()  // true

// ============ 天干地支 / 生肖 ============
lunar.YearGanZhi()    // "甲辰"
lunar.MonthGanZhi()   // "庚午"
lunar.DayGanZhi()     // "丙寅"
lunar.Zodiac()        // "龙"

// ============ 中文命名 ============
lunar.YearName()      // "甲辰年"
lunar.MonthName()     // "六月" / "闰二月"
lunar.DayName()       // "十五" / "初一" / "三十"
lunar.String()        // "甲辰年六月十五"

// ============ 格式化（占位符）============
lunar.Format("YYYY年MM月DD日")       // "2024年06月15日"
lunar.Format("GZ年ZODIAC")           // "甲辰年龙"
lunar.Format("MNAMEDNAME")            // "六月十五"

// ============ 信息查询 ============
lunar.LeapMonth()     // 该年闰月月份（0=无）
lunar.DaysInMonth()   // 当月天数
lunar.DaysInYear()    // 全年天数

// ============ 比较 ============
lunar.IsBefore(other)
lunar.IsAfter(other)
lunar.IsEqual(other)
```

#### 静态工具函数

```go
DateTool.SolarToLunar(2024, 6, 15)            // (年,月,日,是否闰月)
DateTool.LunarToSolar(2024, 1, 1, false)      // (年,月,日)
DateTool.LeapMonthOfYear(2023)                 // 2（闰二月）
DateTool.GanZhiYear(2024)                      // "甲辰"
DateTool.ZodiacOfYear(2024)                    // "龙"
```

### DzStopwatch

代码执行计时器，支持暂停/继续累积，提供多种时间单位。

```go
sw := DateTool.NewDzStopwatch()    // 创建并自动启动

// 链式控制
sw.Stop()      // 暂停
sw.Start()     // 继续（不丢失已累积时间）
sw.Reset()     // 清零并停止
sw.Restart()   // 清零并重新启动

// 读取已用时间
sw.Elapsed()           // time.Duration
sw.ElapsedMs()         // int64（毫秒）
sw.ElapsedSeconds()    // float64（秒）
sw.ElapsedMinutes()    // float64（分钟）
sw.ElapsedHours()      // float64（小时）
sw.ElapsedDays()       // float64（天）
sw.ElapsedWeeks()      // float64（周）
sw.IsRunning()         // bool

// 格式化
sw.String()            // "1.234s"
sw.Format("ms")        // "1234ms"
sw.Format("sec")       // "1.234s"
sw.Format("full")      // "00:00:01.234"

// 暂停/继续累积正确
sw := DateTool.NewDzStopwatch()
time.Sleep(50 * time.Millisecond)
sw.Stop()                         // 第一段 ~50ms
time.Sleep(50 * time.Millisecond) // 等待期间不计时
sw.Start()                        // 继续
time.Sleep(50 * time.Millisecond)
sw.Stop()                         // 总计 ~100ms
```

---

## TimeIntervalTool

```go
import "github.com/PWND0U/dztool/TimeIntervalTool"
```

命名计时器，底层使用 `map[string]time.Time`，支持多个并发计时。

```go
timer := TimeIntervalTool.NewDzTimeInterval()

timer.Start("request")            // 启动计时器
timer.Start("db", "cache")        // 批量启动

timer.IntervalMs("request")       // 毫秒
timer.IntervalSecond("request")   // 秒
timer.IntervalMinute("request")   // 分钟
timer.IntervalHour("request")     // 小时
timer.IntervalDay("request")      // 天
timer.IntervalWeek("request")     // 周

timer.ReStart("request")          // 重置计时器
```

---

## IOTool

```go
import "github.com/PWND0U/dztool/IOTool"
```

IO 操作工具库，涵盖流操作、文件读写、文件类型检测、目录监听、快速缓冲、路径工具六大模块。

### DzStreamTool — 流操作

封装 `io.Reader` / `io.Writer` 的常用操作。

```go
// 读取
data, err := IOTool.ReadAll(reader)           // []byte
str, err := IOTool.ReadAllString(reader)      // string
lines, err := IOTool.ReadLines(path)          // []string（按行读取文件）
lines, err := IOTool.ReadLinesFromReader(r)   // []string（按行读取 Reader）

// 写入
n, err := IOTool.WriteAll(writer, data)       // int, error
n, err := IOTool.WriteString(writer, "hello")  // int, error

// 流拷贝
n, err := IOTool.Copy(dst, src)               // int64
n, err := IOTool.CopyBuffer(dst, src, buf)    // 带缓冲区

// 组合
r := IOTool.MultiReader(r1, r2, r3)          // 合并多个 Reader
w := IOTool.MultiWriter(w1, w2)              // 合并多个 Writer
r = IOTool.LimitReader(r, 1024)              // 限制读取字节数

// 其他
data, err := IOTool.ReadAt(readerAt, off, n) // 带偏移读取
n, err := IOTool.Pipe(src, dst)              // 管道操作
```

### DzFileTool — 文件读写操作

```go
// 读取
data, err := IOTool.ReadFile("test.txt")            // []byte
str, err := IOTool.ReadFileAsString("test.txt")     // string
lines, err := IOTool.ReadLines("test.txt")          // []string

// 写入
err := IOTool.WriteFile("out.txt", data)             // 覆盖写入
err := IOTool.WriteFileString("out.txt", "hello")
err := IOTool.WriteFileSync("out.txt", data)        // 写入并 Sync 刷盘

// 追加
err := IOTool.AppendFile("log.txt", data)
err := IOTool.AppendFileString("log.txt", "line\n")

// 操作
err := IOTool.CopyFile("src.txt", "dst.txt")         // 复制（自动创建目录）
err := IOTool.MoveFile("old.txt", "new.txt")          // 移动

// 状态查询
ok := IOTool.FileExists("test.txt")                  // bool
ok := IOTool.DirExists("/tmp/mydir")
size, err := IOTool.FileSize("test.txt")              // int64
modTime, err := IOTool.LastModified("test.txt")       // time.Time

// 目录
err := IOTool.MkdirAll("/tmp/a/b/c")                 // 递归创建
err := IOTool.Remove("/tmp/old")                      // 递归删除
infos, err := IOTool.ListDir(".")                     // 直接子项
err := IOTool.WalkDir(".", fn)                        // 遍历

// 临时文件
f, err := IOTool.CreateTempFile("", "dztool_*.tmp")
dir, err := IOTool.CreateTempDir("", "dztool_*")
```

### DzFileTypeTool — 文件类型判断

基于魔数(Magic Number)的文件类型检测，支持自定义注册和 30+ 内置常见类型。

```go
// 检测
typeName, err := IOTool.DetectFileType("photo.jpg")     // "JPEG"
typeName := IOTool.DetectFileTypeByBytes(data)          // string
typeName, err := IOTool.DetectFileTypeFromReader(r)
ok, err := IOTool.IsFileType("doc.pdf", "PDF")          // bool

// 自定义注册
IOTool.RegisterFileType("MyType",
    [][]byte{[]byte{0x89, 'M', 'Y', 'F'}},  // 魔数
    []string{".myf"},                        // 扩展名
    0)                                        // 偏移量

types := IOTool.GetRegisteredTypes()          // 所有已注册类型
```

**内置支持的文件类型**：

| 类别 | 支持格式 |
|------|---------|
| 图片 | JPEG、PNG、GIF、BMP、WebP、ICO、TIFF |
| 文档 | PDF、DOC/DOCX、XLS/XLSX、PPT/PPTX |
| 压缩 | ZIP、RAR、7Z、GZIP、BZ2、TAR |
| 音频 | MP3、WAV、FLAC、OGG、AAC |
| 视频 | MP4、AVI、MKV、MOV、WMV |
| 可执行 | EXE、ELF、Mach-O |
| 文本 | XML、HTML、JSON |

### DzWatchTool — 目录/文件监听

支持多级目录递归监听、延迟合并触发、多观察者模式、文件跟随。

```go
// 简洁 API — 一行代码搞定监听
watcher, err := IOTool.Watch("./config", func(event IOTool.WatchEvent) {
    fmt.Printf("[%s] %s (dir=%v)\n", event.Type, event.Path, event.IsDir)
})

// 完整 API
w, _ := IOTool.NewWatcher("./logs", "./data")

// 配置
w.SetDepth(3)                                    // 监听深度（-1=无限）
w.SetDelay(500 * time.Millisecond)                // 延迟合并时间窗口
w.SetFilter("*.log")                             // glob 过滤
w.SetFilterFunc(func(p string) bool {             // 自定义过滤
    return !strings.HasSuffix(p, ".tmp")
})

// 观察者
type MyObserver struct{}
func (o *MyObserver) OnEvent(e IOTool.WatchEvent) {
    log.Printf("事件: %s → %s", e.Type, e.Path)
}
w.AddObserver(&MyObserver{})

// 生命周期
w.Start()                                         // 开始监听
defer w.Stop()                                    // 停止释放

// 文件跟随 — 删除重建后自动重新跟随
w.Follow("app.log", &MyObserver{})
```

**延迟合并触发**：同一文件在延迟窗口内的多次修改事件合并为一个通知，避免频繁触发回调。

**多级目录**：新建子目录时自动递归注册监听，受 `SetDepth` 深度限制。

**事件类型**：`Create`、`Write`、`Remove`、`Rename`、`Chmod`

### DzFastBuffer — 快速缓冲

基于分块存储策略的高性能缓冲区，内部维护 `[][]byte` 缓冲集，避免单一数组扩容的大块内存拷贝。

```go
buf := IOTool.NewFastBuffer()                    // 默认块大小 4096
buf2 := IOTool.NewFastBufferWithSize(8192)       // 自定义块大小

// 写入
n := buf.Write([]byte("hello"))                  // 跨块自动扩容
buf.WriteFrom(strings.NewReader(" world"))       // io.ReaderFrom 接口

// 读取
chunk := buf.Read(5)                              // []byte
all := buf.ReadAll()                              // 全部数据并清空

// 信息
len := buf.Len()                                  // 未读长度
buf.Reset()                                       // 清空保留内存

// io.WriterTo 接口
buf.WriteTo(os.Stdout)                            // 写入 Writer
```

### DzPathTool — 路径工具

```go
// 解析
IOTool.BaseName("/path/to/file.txt")              // "file"
IOTool.Ext("/path/to/file.txt")                   // ".txt"
IOTool.FileName("/path/to/file.txt")              // "file.txt"
IOTool.Parent("/path/to/file.txt")                // "/path/to"

// 操作
IOTool.Join("a", "b", "c")                       // "a/b/c"
IOTool.Normalize(".././foo/../bar")               // 规范化路径
IOTool.Rel("/a/b/c", "/a/b/c/d/e")               // "d/e"
IOTool.ReplaceExt("/path/test.txt", ".md")        // "/path/test.md"

// 安全
IOTool.IsSubPath("/home/user", "/home/user/doc")  // true
IOTool.IsSafePath("../../etc/passwd")              // false
IOTool.IsAbs("/usr/local/bin")                    // true

// 遍历
IOTool.Depth("/a/b/c/d")                          // 4
paths := IOTool.PathIterator("/a/b/c/d")          // ["/", "/a", "/a/b", ...]

// 环境
home, _ := IOTool.HomeDir()
exeDir, _ := IOTool.ExecutableDir()
hostName := IOTool.GetLocalHostName()
pType := IOTool.PathType("/etc/passwd")            // "file" / "dir" / "not_exist"
```

---

## NetTool

```go
import "github.com/PWND0U/dztool/NetTool"
```

网络工具库，参考 Hutool NetUtil 实现，涵盖 IP 转换、端口检测、网卡信息、Ping、CIDR、Cookie 解析等 39 个函数。

### IP 地址转换

```go
// IPv4 ↔ uint32
NetTool.Ipv4ToLong("192.168.1.1")                 // 3232235777
NetTool.LongToIpv4(3232235777)                     // "192.168.1.1"

// IPv6 ↔ big.Int
bigInt, _ := NetTool.Ipv6ToBigInt("::1")           // *big.Int
NetTool.BigIntToIpv6(bigInt)                       // "::1"
```

### 端口检测

```go
NetTool.IsValidPort(8080)                          // true
NetTool.IsValidPort(-1)                            // false
NetTool.IsValidPort(70000)                         // false

NetTool.IsUsableLocalPort(8080)                    // 是否可用（尝试监听）

port := NetTool.GetUsableLocalPort()               // 随机可用端口（1024~65535）
port = NetTool.GetUsableLocalPortRange(20000)      // 从指定端口开始查找
port, _ = NetTool.GetUsableLocalPortBetween(30000, 40000)

ports, _ := NetTool.GetUsableLocalPorts(3, 5000, 60000) // 批量获取
```

### IP 判断与隐藏

```go
NetTool.IsInnerIP("10.0.0.1")                     // true（A 类内网）
NetTool.IsInnerIP("172.16.0.1")                   // true（B 类内网）
NetTool.IsInnerIP("192.168.1.1")                  // true（C 类内网）
NetTool.IsInnerIP("127.0.0.1")                    // true（回环）
NetTool.IsInnerIP("8.8.8.8")                      // false

NetTool.IsInRange("192.168.1.100", "192.168.1.0/24") // true

NetTool.HideIpPart("192.168.1.1")                 // "192.168.1.*"
NetTool.HideIpPartFromLong(3232235777)             // "192.168.1.*"
```

### 网络地址构建 & DNS

```go
addr, _ := NetTool.BuildTCPAddr("example.com:8080", 9090)  // 使用 host 中端口
addr, _ := NetTool.BuildTCPAddr("example.com", 9090)        // 使用默认端口
addr, _ := NetTool.CreateTCPAddr("localhost", 3306)

ip := NetTool.GetIpByHost("localhost")               // DNS 解析
```

### 网卡信息

```go
// 网卡
ifaces, _ := NetTool.GetNetworkInterfaces()          // 所有网卡
iface, _ := NetTool.GetNetworkInterface("eth0")      // 指定网卡

// IP 地址
ipv4s, _ := NetTool.LocalIpv4s()                    // 本机 IPv4 列表
ipv6s, _ := NetTool.LocalIpv6s()                    // 本机 IPv6 列表
ips, _ := NetTool.LocalIps()                         // 本机所有 IP

ipStr := NetTool.GetLocalhostStr()                   // 第一个非回环 IP
ip, _ := NetTool.GetLocalhost()

addrs, _ := NetTool.LocalAddressList(nil)            // 过滤后的地址列表
ipList := NetTool.ToIpList(addrs)                    // 地址→字符串

// MAC
mac, _ := NetTool.GetLocalMacAddress()               // 本机 MAC（如 "AA-BB-CC-DD-EE-FF"）
mac, _ := NetTool.GetMacAddress(iface)
mac, _ := NetTool.GetMacAddressWithSeparator(iface, ":")  // 自定义分隔符

// 主机名
name := NetTool.GetLocalHostName()
```

### Ping & 远程端口

```go
NetTool.Ping("127.0.0.1")                           // 默认 3 秒超时
NetTool.PingWithTimeout("8.8.8.8", 1000)             // 指定超时（毫秒）

NetTool.IsOpen("github.com", 443, 3*time.Second)     // 远程端口是否开启
```

### Cookie & URL & IDN

```go
cookies := NetTool.ParseCookies("session=abc123; lang=zh-CN")  // []*http.Cookie

url := NetTool.ToAbsoluteUrl("https://example.com/base/", "../page.html")
// "https://example.com/page.html"

punycode, _ := NetTool.IdnToASCII("中文.com")          // "xn--fiq228c.com"
```

### 反向代理 & 工具函数

```go
// 从多级反向代理提取真实 IP
ip := NetTool.GetMultistageReverseProxyIp("unknown, 10.0.0.1, 192.168.1.1")
// "10.0.0.1"（第一个非 unknown 的 IP）

NetTool.IsUnknown("")                                 // true
NetTool.IsUnknown("unknown")                          // true
NetTool.IsUnknown("NULL")                             // true
NetTool.IsUnknown("192.168.1.1")                      // false

// Socket 发送
NetTool.NetCat("example.com", 80, []byte("GET / HTTP/1.0\r\n\r\n"))
```

---

## ServerTool

```go
import "github.com/PWND0U/dztool/ServerTool"
```

Server-Sent Events (SSE) 编解码与 HTTP Flush。

```go
// 构造 SSE 事件
event := ServerTool.NewDzServerSentEvent(
    []byte("Hello,World"), // data
    "update",              // event
    "1",                   // id
    "注释",                 // comment
    1,                     // retry（秒）
)

// 编码
encoded := event.Encode()                           // []byte（SSE 文本格式）
encoded2 := ServerTool.EncodeDzServerSentEvent(event) // 同上

// 解码
decoded := ServerTool.DecodeDzServerSentEvent(encoded)

// HTTP Flush（写入响应并刷新）
http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
    event.SSEDataFlush(w)
})
```

---

## Algorithm

```go
import "github.com/PWND0U/dztool/Algorithm"
```

### Levenshtein 编辑距离

```go
// DzLevenshtein(str1, str2, weightInsert, weightReplace, weightDelete)
dist := Algorithm.DzLevenshtein("kitten", "sitting", 1, 1, 1)  // 3
```

---

## Algorithm/DzJWT

```go
import "github.com/PWND0U/dztool/Algorithm/DzJWT"
```

JWT 创建、解析与验证，使用 HS256 签名。

```go
// 创建 Claims
claims := DzJWT.NewJWTClaims()
claims.SetIssuer("my-app")
claims.SetSubject("user-123")
claims.SetAudience([]string{"web", "mobile"})
claims.SetExpiresAt(time.Now().Add(24 * time.Hour))
claims.SetNotBefore(time.Now())
claims.SetIssuedAt(time.Now())
claims.SetID("token-001")
claims.SetExtraData(map[string]interface{}{"role": "admin"})
claims.AddExtraDataByKey("level", 5)

// 生成 JWT
token, err := claims.GenJWT("my-secret")

// 验证
valid := DzJWT.VerifyJWT(token, "my-secret")  // true/false

// 解析为 map
claimsMap, err := DzJWT.ParseJWT(token, "my-secret")

// 解析为 JWTClaims 结构体
parsedClaims, err := DzJWT.ParseJWTToEntity(token, "my-secret")

// 从 token 直接构造 JWTClaims
jc := DzJWT.NewJWTClaimsByToken(token, "my-secret")

// 静态生成
token, err = DzJWT.GenJWT(*claims, "my-secret")
```

---

## dzAi/dzOcr

```go
import "github.com/PWND0U/dztool/dzAi/dzOcr"
```

基于 PaddleOCR + ONNX Runtime 的文字识别引擎，支持 CPU、CUDA、DirectML、TensorRT 后端。

> **前提**：需自行准备 PaddleOCR 的 ONNX 模型文件（det.onnx、rec.onnx、dict.txt）和 ONNX Runtime 动态库。

```go
// 配置
config := dzOcr.Config{
    OnnxRuntimeLibPath:    "path/to/onnxruntime.dll",
    GPUOnnxRuntimeLibPath: "path/to/onnxruntime_gpu.dll",
    DetModelPath:          "path/to/det.onnx",
    RecModelPath:          "path/to/rec.onnx",
    DictPath:              "path/to/dict.txt",
    // 可选参数
    UseCuda:             false,
    UseDirectML:         false,
    UseTensorrt:         false,
    DetMaxSideLen:       960,
    DetOutsideExpandPix: 10,
    RecHeight:           48,
    HeatmapThreshold:    0.3,
}

// 创建引擎
engine, err := dzOcr.NewDzPaddleOcrEngine(config)
if err != nil { log.Fatal(err) }
defer engine.Destroy()

// 全流程 OCR
results, err := engine.RunOCRByFile("test.png")
for _, r := range results {
    fmt.Printf("区域: %v, 文本: %s, 置信度: %.2f\n", r.Box, r.Text, r.Score)
}

// 仅检测
boxes, err := engine.RunDetect(img)

// 仅识别
result, err := engine.RunRecognize(img, boxes[0], session)

// 绘制检测框
resultImg := dzOcr.DrawBoxes(img, boxes, nil)  // nil = 红色
```

---

## 依赖

| 依赖 | 版本 | 用途 |
|---|---|---|
| `bytedance/sonic` | v1.15.1 | 高性能 JSON 序列化/反序列化（替代标准库 `encoding/json`） |
| `golang-jwt/jwt/v5` | v5.3.1 | JWT 签发与验证 |
| `spf13/cast` | v1.10.0 | 类型转换（JsonTool / StructTool） |
| `kaptinlin/jsonrepair` | v0.4.5 | JSON 字符串修补（修复不合法 JSON） |
| `yalue/onnxruntime_go` | v1.28.0 | ONNX Runtime Go 绑定（OCR 推理） |
| `up-zero/gotool` | v0.0.0-20260402 | 图像裁剪/旋转/绘制（OCR 后处理） |
| `golang.org/x/image` | v0.39.0 | 扩展图像处理 |
| `golang.org/x/text` | v0.36.0 | 大小写转换 |
| `fsnotify/fsnotify` | v1.8.0 | 跨平台文件系统监听（IOTool 目录监听） |
| `golang.org/x/net` | v0.39.0 | IDN 域名 Punycode 转换（NetTool） |

## 许可证

Apache License 2.0
