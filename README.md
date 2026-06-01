# dztool

一个 Go 工具库，提供字符串、字节、JSON、结构体、时间、SSE、JWT、OCR 等常用操作的链式封装类型。

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
| JsonTool | `DzJsonMap` `DzJsonArray` `DzJsonMapArray` | JSON 点号路径导航 |
| StructTool | `DzStruct` | 结构体拷贝/克隆/比较/Map互转（类似 Java BeanUtils） |
| TimeIntervalTool | `DzTimeInterval` | 命名计时器 |
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
| `golang-jwt/jwt/v5` | v5.3.1 | JWT 签发与验证 |
| `spf13/cast` | v1.10.0 | 类型转换（JsonTool / StructTool） |
| `yalue/onnxruntime_go` | v1.28.0 | ONNX Runtime Go 绑定（OCR 推理） |
| `up-zero/gotool` | v0.0.0-20260402 | 图像裁剪/旋转/绘制（OCR 后处理） |
| `golang.org/x/image` | v0.39.0 | 扩展图像处理 |
| `golang.org/x/text` | v0.36.0 | 大小写转换 |

## 许可证

Apache License 2.0
