package StructTool

import (
	"fmt"
	"reflect"
	"testing"
)

// ==================== 测试用结构体定义 ====================

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Name   string
	Age    int
	Salary float64
}

type Student struct {
	Name  string
	Age   int
	Score float64
}

type UserDTO struct {
	UserName string `json:"user_name"`
	UserAge  int    `json:"user_age"`
	Email    string `json:"email"`
}

type User struct {
	UserName string `json:"user_name"`
	UserAge  int    `json:"user_age"`
	Address  string `json:"address"`
}

type Address struct {
	City    string
	Street  string
	ZipCode string
}

type NestedPerson struct {
	Name    string
	Age     int
	Address Address
}

type Inner struct {
	InnerName string
	InnerID   int
}

type Outer struct {
	OuterName string
	Inner     // 匿名嵌入字段
}

type WithPointer struct {
	Name  *string
	Age   int
	Score float64
}

type WithSlice struct {
	Name   string
	Scores []int
}

type WithMap struct {
	Name    string
	Meta    map[string]string
}

type EmptyStruct struct{}

// ==================== 指针结构体测试类型 ====================

type PtrInner struct {
	Name string
	ID   int
}

type WithPtrStruct struct {
	Name  string
	Inner *PtrInner
}

type WithDoublePtrStruct struct {
	Name  string
	Inner **PtrInner
}

type PtrEmbedded struct {
	*PtrInner // 指针匿名嵌入字段
	Label     string
}

type PtrInnerDTO struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type WithPtrStructDTO struct {
	Name  string        `json:"name"`
	Inner *PtrInnerDTO  `json:"inner"`
}

// ==================== 构造与状态测试 ====================

func TestNewDzStruct_ValidStruct(t *testing.T) {
	p := Person{Name: "Alice", Age: 30}
	ds := NewDzStruct(p)
	fmt.Println("IsValid:", ds.IsValid())
	fmt.Println("IsZero:", ds.IsZero())
	fmt.Println("Name:", ds.Field("Name"))
	fmt.Println("Age:", ds.Field("Age"))
}

func TestNewDzStruct_PointerInput(t *testing.T) {
	p := &Person{Name: "Bob", Age: 25}
	ds := NewDzStruct(p)
	fmt.Println("IsValid:", ds.IsValid())
	fmt.Println("Name:", ds.Field("Name"))
}

func TestNewDzStruct_NilInput(t *testing.T) {
	ds := NewDzStruct(nil)
	fmt.Println("IsValid:", ds.IsValid())
	fmt.Println("Err:", ds.Err())
}

func TestNewDzStruct_NonStructInput(t *testing.T) {
	ds := NewDzStruct("hello")
	fmt.Println("IsValid:", ds.IsValid())
	fmt.Println("Err:", ds.Err())

	ds2 := NewDzStruct(42)
	fmt.Println("IsValid:", ds2.IsValid())
	fmt.Println("Err:", ds2.Err())
}

func TestNewDzStruct_NilPointer(t *testing.T) {
	var p *Person = nil
	ds := NewDzStruct(p)
	fmt.Println("IsValid:", ds.IsValid())
	fmt.Println("Err:", ds.Err())
}

// ==================== ToInterface / ToIntf 测试 ====================

func TestToInterface(t *testing.T) {
	p := Person{Name: "Charlie", Age: 28}
	ds := NewDzStruct(p)
	result := ds.ToInterface()
	fmt.Println("ToInterface:", result)
	fmt.Println("Type:", reflect.TypeOf(result))
}

func TestToIntf(t *testing.T) {
	p := Person{Name: "David", Age: 35}
	ds := NewDzStruct(p)
	var target Person
	err := ds.ToIntf(&target)
	fmt.Println("Err:", err)
	fmt.Println("Target:", target)
}

func TestToIntf_TypeMismatch(t *testing.T) {
	p := Person{Name: "Eve", Age: 22}
	ds := NewDzStruct(p)
	var target Employee
	err := ds.ToIntf(&target)
	fmt.Println("Err:", err)
}

// ==================== 同结构体操作测试 ====================

func TestDeepClone_Basic(t *testing.T) {
	p := Person{Name: "Frank", Age: 40}
	ds := NewDzStruct(p)
	cloned := ds.DeepClone()
	fmt.Println("Original:", ds.ToInterface())
	fmt.Println("Cloned:", cloned.ToInterface())
	fmt.Println("Equal:", ds.EqualTo(cloned))
	// 修改克隆体不影响原始
	cloned.SetField("Name", "FrankClone")
	fmt.Println("After modify - Original Name:", ds.Field("Name"))
	fmt.Println("After modify - Cloned Name:", cloned.Field("Name"))
}

func TestDeepClone_Nested(t *testing.T) {
	np := NestedPerson{
		Name: "Grace",
		Age:  30,
		Address: Address{
			City:    "Beijing",
			Street:  "Chaoyang",
			ZipCode: "100000",
		},
	}
	ds := NewDzStruct(np)
	cloned := ds.DeepClone()
	fmt.Println("Nested Equal:", ds.EqualTo(cloned))
}

func TestDeepClone_Slice(t *testing.T) {
	ws := WithSlice{Name: "Hank", Scores: []int{90, 85, 95}}
	ds := NewDzStruct(ws)
	cloned := ds.DeepClone()
	fmt.Println("Slice Equal:", ds.EqualTo(cloned))
}

func TestDeepClone_SliceIndependence(t *testing.T) {
	// 深克隆：修改克隆体的切片不影响原始
	ws := WithSlice{Name: "Hank", Scores: []int{90, 85, 95}}
	ds := NewDzStruct(ws)
	cloned := ds.DeepClone()
	cloned.SetField("Scores", []int{1, 2, 3})
	fmt.Println("DeepClone - Original Scores:", ds.Field("Scores"))
	fmt.Println("DeepClone - Cloned Scores:", cloned.Field("Scores"))
	fmt.Println("DeepClone - Slices independent:", !reflect.DeepEqual(ds.Field("Scores"), cloned.Field("Scores")))
}

func TestDeepCloneReflect_Basic(t *testing.T) {
	p := Person{Name: "Ivy", Age: 27}
	ds := NewDzStruct(p)
	cloned := ds.DeepCloneReflect()
	fmt.Println("DeepCloneReflect Equal:", ds.EqualTo(cloned))
}

func TestDeepCloneReflect_Slice(t *testing.T) {
	ws := WithSlice{Name: "Jack", Scores: []int{88, 92, 76}}
	ds := NewDzStruct(ws)
	cloned := ds.DeepCloneReflect()
	fmt.Println("DeepCloneReflect Slice Equal:", ds.EqualTo(cloned))
}

func TestDeepCloneReflect_Map(t *testing.T) {
	wm := WithMap{Name: "Kate", Meta: map[string]string{"key": "value"}}
	ds := NewDzStruct(wm)
	cloned := ds.DeepCloneReflect()
	fmt.Println("DeepCloneReflect Map Equal:", ds.EqualTo(cloned))
}

func TestDeepCloneReflect_SliceIndependence(t *testing.T) {
	// 深克隆（反射）：修改克隆体的切片不影响原始
	ws := WithSlice{Name: "Jack", Scores: []int{88, 92, 76}}
	ds := NewDzStruct(ws)
	cloned := ds.DeepCloneReflect()
	cloned.SetField("Scores", []int{1, 2})
	fmt.Println("DeepCloneReflect - Original Scores:", ds.Field("Scores"))
	fmt.Println("DeepCloneReflect - Cloned Scores:", cloned.Field("Scores"))
}

func TestShallowClone_Basic(t *testing.T) {
	p := Person{Name: "ShallowAlice", Age: 30}
	ds := NewDzStruct(p)
	cloned := ds.ShallowClone()
	fmt.Println("ShallowClone Equal:", ds.EqualTo(cloned))
	// 修改值类型字段不影响原始
	cloned.SetField("Name", "ShallowBob")
	fmt.Println("ShallowClone - Original Name:", ds.Field("Name"))
	fmt.Println("ShallowClone - Cloned Name:", cloned.Field("Name"))
}

func TestShallowClone_SliceShared(t *testing.T) {
	// 浅克隆：切片等引用类型字段与原始共享底层数据
	ws := WithSlice{Name: "SharedSlice", Scores: []int{1, 2, 3}}
	ds := NewDzStruct(ws)
	cloned := ds.ShallowClone()

	// 获取原始和克隆体的 Scores 字段
	origScores := ds.Field("Scores").([]int)
	clonedScores := cloned.Field("Scores").([]int)

	// 修改克隆体的切片元素应影响原始（因为共享底层数组）
	clonedScores[0] = 999
	fmt.Println("ShallowClone - Original Scores after modify cloned:", origScores)
	fmt.Println("ShallowClone - Cloned Scores:", clonedScores)
	fmt.Println("ShallowClone - Slices share underlying data:", origScores[0] == 999)
}

func TestShallowClone_MapShared(t *testing.T) {
	// 浅克隆：map 字段与原始共享底层数据
	wm := WithMap{Name: "SharedMap", Meta: map[string]string{"key": "value"}}
	ds := NewDzStruct(wm)
	cloned := ds.ShallowClone()

	// 获取 map 并修改
	clonedMeta := cloned.Field("Meta").(map[string]string)
	clonedMeta["new_key"] = "new_value"

	origMeta := ds.Field("Meta").(map[string]string)
	fmt.Println("ShallowClone - Original Meta after modify cloned:", origMeta)
	fmt.Println("ShallowClone - Cloned Meta:", clonedMeta)
	fmt.Println("ShallowClone - Maps share underlying data:", origMeta["new_key"] == "new_value")
}

func TestFields(t *testing.T) {
	p := Person{Name: "Leo", Age: 33}
	ds := NewDzStruct(p)
	fmt.Println("Fields:", ds.Fields())
}

func TestFields_Embedded(t *testing.T) {
	o := Outer{OuterName: "Outer1", Inner: Inner{InnerName: "Inner1", InnerID: 1}}
	ds := NewDzStruct(o)
	fmt.Println("Fields (embedded flat):", ds.Fields())
}

func TestField(t *testing.T) {
	p := Person{Name: "Mia", Age: 29}
	ds := NewDzStruct(p)
	fmt.Println("Field Name:", ds.Field("Name"))
	fmt.Println("Field Age:", ds.Field("Age"))
	fmt.Println("Field NotFound:", ds.Field("NotExist"))
	fmt.Println("Err after not found:", ds.Err())
}

func TestSetField(t *testing.T) {
	p := Person{Name: "Nora", Age: 31}
	ds := NewDzStruct(p)
	ds.SetField("Name", "NoraUpdated").SetField("Age", 32)
	fmt.Println("Name after SetField:", ds.Field("Name"))
	fmt.Println("Age after SetField:", ds.Field("Age"))
}

func TestSetField_TypeConversion(t *testing.T) {
	p := Person{Name: "Oscar", Age: 25}
	ds := NewDzStruct(p)
	// 字符串 "30" 应通过 cast 转换为 int
	ds.SetField("Age", "30")
	fmt.Println("Age after SetField with string:", ds.Field("Age"))
}

func TestZero(t *testing.T) {
	p := Person{Name: "Peggy", Age: 45}
	ds := NewDzStruct(p)
	ds.Zero()
	fmt.Println("IsZero after Zero():", ds.IsZero())
	fmt.Println("Name after Zero():", ds.Field("Name"))
	fmt.Println("Age after Zero():", ds.Field("Age"))
}

func TestIsZero(t *testing.T) {
	p := Person{}
	ds := NewDzStruct(p)
	fmt.Println("IsZero (empty struct):", ds.IsZero())

	p2 := Person{Name: "Quinn"}
	ds2 := NewDzStruct(p2)
	fmt.Println("IsZero (partial):", ds2.IsZero())
}

// ==================== 跨结构体操作测试 ====================

func TestCopyFrom_SameType(t *testing.T) {
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}
	ds := NewDzStruct(p1)
	ds.CopyFrom(p2)
	fmt.Println("After CopyFrom (same type):")
	fmt.Println("  Name:", ds.Field("Name"))
	fmt.Println("  Age:", ds.Field("Age"))
}

func TestCopyFrom_DifferentType(t *testing.T) {
	p := Person{Name: "Alice", Age: 30}
	e := Employee{Name: "EmpAlice", Age: 28, Salary: 5000}
	ds := NewDzStruct(p)
	ds.CopyFrom(e)
	fmt.Println("After CopyFrom (different type):")
	fmt.Println("  Name:", ds.Field("Name"))
	fmt.Println("  Age:", ds.Field("Age"))
	// Person 没有 Salary 字段，不会被复制
}

func TestCopyFrom_PartialOverlap(t *testing.T) {
	s := Student{Name: "StuBob", Age: 20, Score: 95.5}
	p := Person{Name: "Original", Age: 99}
	ds := NewDzStruct(p)
	ds.CopyFrom(s)
	fmt.Println("After CopyFrom (partial overlap):")
	fmt.Println("  Name:", ds.Field("Name"))
	fmt.Println("  Age:", ds.Field("Age"))
}

func TestCopyFromByTag(t *testing.T) {
	dto := UserDTO{UserName: "dto_name", UserAge: 25, Email: "dto@test.com"}
	u := User{UserName: "original", UserAge: 0, Address: "default_addr"}
	ds := NewDzStruct(u)
	ds.CopyFromByTag(dto, "json")
	fmt.Println("After CopyFromByTag:")
	fmt.Println("  UserName:", ds.Field("UserName"))
	fmt.Println("  UserAge:", ds.Field("UserAge"))
	fmt.Println("  Address:", ds.Field("Address"))
}

func TestMergeFrom_NonZeroOverwrites(t *testing.T) {
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 0} // Age 为零值
	ds := NewDzStruct(p1)
	ds.MergeFrom(p2)
	fmt.Println("After MergeFrom (zero Age not overwritten):")
	fmt.Println("  Name:", ds.Field("Name"))
	fmt.Println("  Age:", ds.Field("Age"))
}

func TestMergeFrom_NonZeroOverwrites2(t *testing.T) {
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25} // 两个都非零
	ds := NewDzStruct(p1)
	ds.MergeFrom(p2)
	fmt.Println("After MergeFrom (both non-zero):")
	fmt.Println("  Name:", ds.Field("Name"))
	fmt.Println("  Age:", ds.Field("Age"))
}

func TestMergeFromByTag(t *testing.T) {
	dto := UserDTO{UserName: "merged_name", UserAge: 0, Email: "merged@test.com"}
	u := User{UserName: "original", UserAge: 30, Address: "some_addr"}
	ds := NewDzStruct(u)
	ds.MergeFromByTag(dto, "json")
	fmt.Println("After MergeFromByTag:")
	fmt.Println("  UserName:", ds.Field("UserName"))
	fmt.Println("  UserAge:", ds.Field("UserAge"))
	fmt.Println("  Address:", ds.Field("Address"))
}

func TestCompareTo_Identical(t *testing.T) {
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	ds := NewDzStruct(p1)
	diff := ds.CompareTo(p2)
	fmt.Println("Diff (identical):", diff)
}

func TestCompareTo_Different(t *testing.T) {
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}
	ds := NewDzStruct(p1)
	diff := ds.CompareTo(p2)
	fmt.Println("Diff (different):", diff)
}

func TestDiffFields(t *testing.T) {
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 30} // 只有 Name 不同
	ds := NewDzStruct(p1)
	diffs := ds.DiffFields(p2)
	fmt.Println("DiffFields:", diffs)
}

func TestEqualTo_True(t *testing.T) {
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	ds := NewDzStruct(p1)
	fmt.Println("EqualTo (same):", ds.EqualTo(p2))
}

func TestEqualTo_False(t *testing.T) {
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}
	ds := NewDzStruct(p1)
	fmt.Println("EqualTo (different):", ds.EqualTo(p2))
}

// ==================== 结构体-Map 互转测试 ====================

func TestToMap_Basic(t *testing.T) {
	p := Person{Name: "Alice", Age: 30}
	ds := NewDzStruct(p)
	m := ds.ToMap()
	fmt.Println("ToMap:", m)
}

func TestToMap_Nested(t *testing.T) {
	np := NestedPerson{
		Name: "Bob",
		Age:  28,
		Address: Address{
			City:    "Shanghai",
			Street:  "Nanjing Rd",
			ZipCode: "200000",
		},
	}
	ds := NewDzStruct(np)
	m := ds.ToMap()
	fmt.Println("ToMap (nested):", m)
}

func TestToMap_Embedded(t *testing.T) {
	o := Outer{OuterName: "Outer1", Inner: Inner{InnerName: "Inner1", InnerID: 42}}
	ds := NewDzStruct(o)
	m := ds.ToMap()
	fmt.Println("ToMap (embedded flat):", m)
}

func TestToMapByTag(t *testing.T) {
	u := UserDTO{UserName: "tag_user", UserAge: 25, Email: "tag@test.com"}
	ds := NewDzStruct(u)
	m := ds.ToMapByTag("json")
	fmt.Println("ToMapByTag (json):", m)
}

func TestFromMap_Basic(t *testing.T) {
	p := Person{}
	ds := NewDzStruct(p)
	ds.FromMap(map[string]interface{}{
		"Name": "Charlie",
		"Age":  35,
	})
	fmt.Println("After FromMap:")
	fmt.Println("  Name:", ds.Field("Name"))
	fmt.Println("  Age:", ds.Field("Age"))
}

func TestFromMap_TypeConversion(t *testing.T) {
	p := Person{}
	ds := NewDzStruct(p)
	// map 中 Age 为字符串 "28"，应通过 cast 转为 int
	ds.FromMap(map[string]interface{}{
		"Name": "Diana",
		"Age":  "28",
	})
	fmt.Println("After FromMap (type conversion):")
	fmt.Println("  Name:", ds.Field("Name"))
	fmt.Println("  Age:", ds.Field("Age"))
}

func TestFromMap_Nested(t *testing.T) {
	np := NestedPerson{}
	ds := NewDzStruct(np)
	ds.FromMap(map[string]interface{}{
		"Name": "Eve",
		"Age":  22,
		"Address": map[string]interface{}{
			"City":    "Guangzhou",
			"Street":  "Tianhe Rd",
			"ZipCode": "510000",
		},
	})
	fmt.Println("After FromMap (nested):")
	fmt.Println("  Name:", ds.Field("Name"))
	fmt.Println("  Age:", ds.Field("Age"))
	addr := ds.Field("Address")
	fmt.Println("  Address:", addr)
}

func TestFromMapByTag(t *testing.T) {
	u := UserDTO{}
	ds := NewDzStruct(u)
	ds.FromMapByTag(map[string]interface{}{
		"user_name": "TagUser",
		"user_age":  30,
		"email":     "tag@test.com",
	}, "json")
	fmt.Println("After FromMapByTag:")
	fmt.Println("  UserName:", ds.Field("UserName"))
	fmt.Println("  UserAge:", ds.Field("UserAge"))
	fmt.Println("  Email:", ds.Field("Email"))
}

func TestRoundTrip(t *testing.T) {
	p := Person{Name: "Frank", Age: 40}
	ds := NewDzStruct(p)
	m := ds.ToMap()
	p2 := Person{}
	ds2 := NewDzStruct(p2)
	ds2.FromMap(m)
	fmt.Println("RoundTrip Equal:", ds.EqualTo(ds2))
}

// ==================== 静态工具函数测试 ====================

func TestCopyStruct(t *testing.T) {
	src := Employee{Name: "EmpSrc", Age: 30, Salary: 8000}
	dst := Person{Name: "DstOrig", Age: 0}
	err := CopyStruct(&src, &dst)
	fmt.Println("CopyStruct err:", err)
	fmt.Println("CopyStruct dst:", dst)
}

func TestCopyStructByTag(t *testing.T) {
	src := UserDTO{UserName: "dto_src", UserAge: 25, Email: "src@test.com"}
	dst := User{UserName: "user_dst", UserAge: 0, Address: "addr"}
	err := CopyStructByTag(&src, &dst, "json")
	fmt.Println("CopyStructByTag err:", err)
	fmt.Println("CopyStructByTag dst:", dst)
}

func TestStructToMap(t *testing.T) {
	p := Person{Name: "StaticAlice", Age: 30}
	m, err := StructToMap(p)
	fmt.Println("StructToMap err:", err)
	fmt.Println("StructToMap result:", m)
}

func TestStructToMapByTag(t *testing.T) {
	u := UserDTO{UserName: "StaticTag", UserAge: 28, Email: "static@test.com"}
	m, err := StructToMapByTag(u, "json")
	fmt.Println("StructToMapByTag err:", err)
	fmt.Println("StructToMapByTag result:", m)
}

func TestMapToStruct(t *testing.T) {
	data := map[string]interface{}{
		"Name": "MapSrc",
		"Age":  33,
	}
	var p Person
	err := MapToStruct(data, &p)
	fmt.Println("MapToStruct err:", err)
	fmt.Println("MapToStruct result:", p)
}

func TestMapToStructByTag(t *testing.T) {
	data := map[string]interface{}{
		"user_name": "MapTagSrc",
		"user_age":  29,
		"email":     "maptag@test.com",
	}
	var u UserDTO
	err := MapToStructByTag(data, &u, "json")
	fmt.Println("MapToStructByTag err:", err)
	fmt.Println("MapToStructByTag result:", u)
}

func TestDeepCloneStruct(t *testing.T) {
	p := Person{Name: "CloneSrc", Age: 50}
	cloned, err := DeepCloneStruct(p)
	fmt.Println("DeepCloneStruct err:", err)
	fmt.Println("DeepCloneStruct result:", cloned)
	fmt.Println("DeepCloneStruct equal:", reflect.DeepEqual(p, cloned))
}

func TestDeepCloneStructReflect(t *testing.T) {
	ws := WithSlice{Name: "CloneSlice", Scores: []int{100, 99, 98}}
	cloned, err := DeepCloneStructReflect(ws)
	fmt.Println("DeepCloneStructReflect err:", err)
	fmt.Println("DeepCloneStructReflect equal:", reflect.DeepEqual(ws, cloned))
}

func TestShallowCloneStruct(t *testing.T) {
	ws := WithSlice{Name: "ShallowStatic", Scores: []int{10, 20}}
	cloned, err := ShallowCloneStruct(ws)
	fmt.Println("ShallowCloneStruct err:", err)
	fmt.Println("ShallowCloneStruct equal:", reflect.DeepEqual(ws, cloned))

	// 浅克隆的切片共享底层数据
	origScores := ws.Scores
	clonedTyped := cloned.(WithSlice)
	clonedTyped.Scores[0] = 999
	fmt.Println("ShallowCloneStruct - Original Scores[0] after modify cloned:", origScores[0])
	fmt.Println("ShallowCloneStruct - Slices share data:", origScores[0] == 999)
}

func TestCompareStruct(t *testing.T) {
	p1 := Person{Name: "Comp1", Age: 30}
	p2 := Person{Name: "Comp2", Age: 30}
	diffs, err := CompareStruct(p1, p2)
	fmt.Println("CompareStruct err:", err)
	fmt.Println("CompareStruct diffs:", diffs)
}

// ==================== 边界情况测试 ====================

func TestEmptyStruct(t *testing.T) {
	e := EmptyStruct{}
	ds := NewDzStruct(e)
	fmt.Println("EmptyStruct IsValid:", ds.IsValid())
	fmt.Println("EmptyStruct IsZero:", ds.IsZero())
	fmt.Println("EmptyStruct Fields:", ds.Fields())
	m := ds.ToMap()
	fmt.Println("EmptyStruct ToMap:", m)
}

func TestUnexportedField(t *testing.T) {
	type PrivateField struct {
		Name   string
		secret string // 非导出字段
	}
	pf := PrivateField{Name: "Visible", secret: "hidden"}
	ds := NewDzStruct(pf)
	fmt.Println("Fields (should skip secret):", ds.Fields())
	m := ds.ToMap()
	fmt.Println("ToMap (should skip secret):", m)
}

func TestPointerFields(t *testing.T) {
	name := "PointerName"
	wp := WithPointer{Name: &name, Age: 25, Score: 88.5}
	ds := NewDzStruct(wp)
	fmt.Println("PointerField Name:", ds.Field("Name"))
	fmt.Println("PointerField Age:", ds.Field("Age"))

	cloned := ds.DeepClone()
	fmt.Println("PointerField DeepClone Equal:", ds.EqualTo(cloned))
}

func TestWithSlice_Fields(t *testing.T) {
	ws := WithSlice{Name: "SliceTest", Scores: []int{1, 2, 3}}
	ds := NewDzStruct(ws)
	m := ds.ToMap()
	fmt.Println("WithSlice ToMap:", m)
}

func TestEmbeddedStruct_FlatCopy(t *testing.T) {
	o := Outer{OuterName: "OuterSrc", Inner: Inner{InnerName: "InnerSrc", InnerID: 100}}
	ds := NewDzStruct(o)
	fmt.Println("Embedded OuterName:", ds.Field("OuterName"))
	fmt.Println("Embedded InnerName (flat):", ds.Field("InnerName"))
	fmt.Println("Embedded InnerID (flat):", ds.Field("InnerID"))

	m := ds.ToMap()
	fmt.Println("Embedded ToMap (flat):", m)
}

func TestChainedOperations(t *testing.T) {
	p := Person{}
	ds := NewDzStruct(p).
		SetField("Name", "Chained").
		SetField("Age", 20)
	fmt.Println("Chained Name:", ds.Field("Name"))
	fmt.Println("Chained Age:", ds.Field("Age"))

	m := ds.ToMap()
	fmt.Println("Chained ToMap:", m)

	ds2 := NewDzStruct(Person{}).FromMap(m)
	fmt.Println("Chained round-trip Equal:", ds.EqualTo(ds2))
}

func TestCopyFrom_DzStructInput(t *testing.T) {
	p1 := Person{Name: "Src", Age: 30}
	ds1 := NewDzStruct(p1)
	p2 := Person{Name: "Dst", Age: 0}
	ds2 := NewDzStruct(p2)
	ds2.CopyFrom(ds1)
	fmt.Println("CopyFrom DzStruct Name:", ds2.Field("Name"))
	fmt.Println("CopyFrom DzStruct Age:", ds2.Field("Age"))
}

func TestAllZeroStruct(t *testing.T) {
	p := Person{}
	ds := NewDzStruct(p)
	fmt.Println("AllZero IsZero:", ds.IsZero())
	fmt.Println("AllZero IsValid:", ds.IsValid())

	cloned := ds.DeepClone()
	fmt.Println("AllZero DeepClone IsZero:", cloned.IsZero())
}

func TestMapWithNilPointer(t *testing.T) {
	wp := WithPointer{Name: nil, Age: 25, Score: 0}
	ds := NewDzStruct(wp)
	m := ds.ToMap()
	fmt.Println("WithNilPointer ToMap:", m)
}

// ==================== 指针结构体测试 ====================

func TestPtrStruct_ToMap(t *testing.T) {
	inner := PtrInner{Name: "Inner1", ID: 1}
	wps := WithPtrStruct{Name: "Outer", Inner: &inner}
	ds := NewDzStruct(wps)
	m := ds.ToMap()
	fmt.Println("PtrStruct ToMap (non-nil):", m)
}

func TestPtrStruct_ToMap_Nil(t *testing.T) {
	wps := WithPtrStruct{Name: "Outer", Inner: nil}
	ds := NewDzStruct(wps)
	m := ds.ToMap()
	fmt.Println("PtrStruct ToMap (nil):", m)
}

func TestPtrStruct_FromMap(t *testing.T) {
	wps := WithPtrStruct{}
	ds := NewDzStruct(wps)
	ds.FromMap(map[string]interface{}{
		"Name": "FromMap",
		"Inner": map[string]interface{}{
			"Name": "InnerValue",
			"ID":   42,
		},
	})
	fmt.Println("PtrStruct FromMap Name:", ds.Field("Name"))
	fmt.Println("PtrStruct FromMap Inner:", ds.Field("Inner"))
}

func TestPtrStruct_FromMap_NilValue(t *testing.T) {
	inner := PtrInner{Name: "Original", ID: 1}
	wps := WithPtrStruct{Name: "Outer", Inner: &inner}
	ds := NewDzStruct(wps)
	ds.FromMap(map[string]interface{}{
		"Inner": nil,
	})
	fmt.Println("PtrStruct FromMap nil Inner:", ds.Field("Inner"))
}

func TestPtrStruct_FromMap_KeepExistingOnNil(t *testing.T) {
	inner := PtrInner{Name: "KeepMe", ID: 99}
	wps := WithPtrStruct{Name: "Outer", Inner: &inner}
	ds := NewDzStruct(wps)
	// 非 nil 的 map 不包含 Inner key → 不应修改
	ds.FromMap(map[string]interface{}{
		"Name": "Updated",
	})
	fmt.Println("PtrStruct FromMap (no Inner key) Name:", ds.Field("Name"))
	fmt.Println("PtrStruct FromMap (no Inner key) Inner:", ds.Field("Inner"))
}

func TestPtrStruct_RoundTrip(t *testing.T) {
	inner := PtrInner{Name: "RoundTrip", ID: 7}
	wps := WithPtrStruct{Name: "Outer", Inner: &inner}
	ds := NewDzStruct(wps)
	m := ds.ToMap()
	wps2 := WithPtrStruct{}
	ds2 := NewDzStruct(wps2)
	ds2.FromMap(m)
	fmt.Println("PtrStruct RoundTrip Equal:", ds.EqualTo(ds2))
	fmt.Println("PtrStruct RoundTrip Name:", ds2.Field("Name"))
	fmt.Println("PtrStruct RoundTrip Inner:", ds2.Field("Inner"))
}

func TestPtrStruct_RoundTrip_Nil(t *testing.T) {
	wps := WithPtrStruct{Name: "Outer", Inner: nil}
	ds := NewDzStruct(wps)
	m := ds.ToMap()
	fmt.Println("PtrStruct RoundTrip ToMap (nil):", m)
	wps2 := WithPtrStruct{Inner: &PtrInner{Name: "ShouldBeCleared"}}
	ds2 := NewDzStruct(wps2)
	ds2.FromMap(m)
	fmt.Println("PtrStruct RoundTrip Inner after nil map:", ds2.Field("Inner"))
}

func TestPtrStruct_SetField_PointerValue(t *testing.T) {
	wps := WithPtrStruct{Name: "Test"}
	ds := NewDzStruct(wps)
	inner := PtrInner{Name: "SetByPointer", ID: 5}
	ds.SetField("Inner", &inner)
	fmt.Println("PtrStruct SetField (*Inner):", ds.Field("Inner"))
}

func TestPtrStruct_SetField_StructValue(t *testing.T) {
	wps := WithPtrStruct{Name: "Test"}
	ds := NewDzStruct(wps)
	ds.SetField("Inner", PtrInner{Name: "SetByValue", ID: 6})
	fmt.Println("PtrStruct SetField (Inner value):", ds.Field("Inner"))
}

func TestPtrStruct_SetField_Nil(t *testing.T) {
	inner := PtrInner{Name: "WillBeNil", ID: 1}
	wps := WithPtrStruct{Name: "Test", Inner: &inner}
	ds := NewDzStruct(wps)
	ds.SetField("Inner", nil)
	fmt.Println("PtrStruct SetField nil:", ds.Field("Inner"))
}

func TestPtrStruct_SetField_CastToPtr(t *testing.T) {
	wp := WithPointer{Age: 20}
	ds := NewDzStruct(wp)
	// "30" → *string 不适用此场景，但 "25" → *int 用 cast
	ds.SetField("Name", "HelloPtr")
	fmt.Println("PtrStruct SetField cast *string:", ds.Field("Name"))
}

func TestPtrStruct_DeepClone(t *testing.T) {
	inner := PtrInner{Name: "CloneMe", ID: 10}
	wps := WithPtrStruct{Name: "Outer", Inner: &inner}
	ds := NewDzStruct(wps)
	cloned := ds.DeepClone()
	fmt.Println("PtrStruct DeepClone Equal:", ds.EqualTo(cloned))
	// 修改克隆体的 Inner 不应影响原始
	cloned.SetField("Inner", &PtrInner{Name: "Modified", ID: 0})
	fmt.Println("PtrStruct DeepClone - Original Inner:", ds.Field("Inner"))
	fmt.Println("PtrStruct DeepClone - Cloned Inner:", cloned.Field("Inner"))
}

func TestPtrStruct_DeepCloneReflect(t *testing.T) {
	inner := PtrInner{Name: "ReflectClone", ID: 20}
	wps := WithPtrStruct{Name: "Outer", Inner: &inner}
	ds := NewDzStruct(wps)
	cloned := ds.DeepCloneReflect()
	fmt.Println("PtrStruct DeepCloneReflect Equal:", ds.EqualTo(cloned))
}

func TestPtrStruct_CompareTo(t *testing.T) {
	inner1 := PtrInner{Name: "Same", ID: 1}
	inner2 := PtrInner{Name: "Diff", ID: 2}
	p1 := WithPtrStruct{Name: "A", Inner: &inner1}
	p2 := WithPtrStruct{Name: "A", Inner: &inner2}
	ds := NewDzStruct(p1)
	diff := ds.CompareTo(p2)
	fmt.Println("PtrStruct CompareTo diff:", diff)
}

func TestPtrStruct_CompareTo_Identical(t *testing.T) {
	inner := PtrInner{Name: "Same", ID: 1}
	p1 := WithPtrStruct{Name: "A", Inner: &inner}
	p2 := WithPtrStruct{Name: "A", Inner: &PtrInner{Name: "Same", ID: 1}}
	ds := NewDzStruct(p1)
	diff := ds.CompareTo(p2)
	fmt.Println("PtrStruct CompareTo identical diff:", diff)
}

func TestPtrStruct_EqualTo(t *testing.T) {
	inner := PtrInner{Name: "Eq", ID: 1}
	p1 := WithPtrStruct{Name: "A", Inner: &inner}
	p2 := WithPtrStruct{Name: "A", Inner: &PtrInner{Name: "Eq", ID: 1}}
	ds := NewDzStruct(p1)
	fmt.Println("PtrStruct EqualTo (same values):", ds.EqualTo(p2))
}

func TestPtrStruct_EqualTo_Nil(t *testing.T) {
	p1 := WithPtrStruct{Name: "A", Inner: nil}
	p2 := WithPtrStruct{Name: "A", Inner: nil}
	ds := NewDzStruct(p1)
	fmt.Println("PtrStruct EqualTo (both nil):", ds.EqualTo(p2))
}

func TestPtrStruct_CopyFrom(t *testing.T) {
	inner := PtrInner{Name: "SrcInner", ID: 5}
	src := WithPtrStruct{Name: "Src", Inner: &inner}
	dst := WithPtrStruct{Name: "Dst"}
	ds := NewDzStruct(dst)
	ds.CopyFrom(src)
	fmt.Println("PtrStruct CopyFrom Name:", ds.Field("Name"))
	fmt.Println("PtrStruct CopyFrom Inner:", ds.Field("Inner"))
}

func TestPtrStruct_CopyFrom_DifferentType(t *testing.T) {
	inner := PtrInnerDTO{Name: "CrossType", ID: 3}
	src := WithPtrStructDTO{Name: "Src", Inner: &inner}
	dst := WithPtrStruct{Name: "Dst"}
	ds := NewDzStruct(dst)
	ds.CopyFrom(src)
	fmt.Println("PtrStruct CopyFrom (cross type) Name:", ds.Field("Name"))
	fmt.Println("PtrStruct CopyFrom (cross type) Inner:", ds.Field("Inner"))
}

func TestPtrStruct_MergeFrom(t *testing.T) {
	inner := PtrInner{Name: "MergeInner", ID: 8}
	src := WithPtrStruct{Name: "SrcName", Inner: &inner}
	dst := WithPtrStruct{Name: "DstName", Inner: &PtrInner{Name: "DstInner", ID: 0}}
	ds := NewDzStruct(dst)
	ds.MergeFrom(src)
	fmt.Println("PtrStruct MergeFrom Name:", ds.Field("Name"))
	fmt.Println("PtrStruct MergeFrom Inner:", ds.Field("Inner"))
}

func TestPtrStruct_MergeFrom_NilPointer(t *testing.T) {
	src := WithPtrStruct{Name: "SrcName", Inner: nil}
	dst := WithPtrStruct{Name: "DstName", Inner: &PtrInner{Name: "KeepMe", ID: 1}}
	ds := NewDzStruct(dst)
	ds.MergeFrom(src)
	fmt.Println("PtrStruct MergeFrom (nil src Inner) Name:", ds.Field("Name"))
	fmt.Println("PtrStruct MergeFrom (nil src Inner) Inner:", ds.Field("Inner"))
}

// ==================== 双重指针测试 ====================

func TestDoublePtr_ToMap(t *testing.T) {
	inner := PtrInner{Name: "Double", ID: 1}
	ptrInner := &inner
	wds := WithDoublePtrStruct{Name: "Outer", Inner: &ptrInner}
	ds := NewDzStruct(wds)
	m := ds.ToMap()
	fmt.Println("DoublePtr ToMap:", m)
}

func TestDoublePtr_ToMap_Nil(t *testing.T) {
	wds := WithDoublePtrStruct{Name: "Outer", Inner: nil}
	ds := NewDzStruct(wds)
	m := ds.ToMap()
	fmt.Println("DoublePtr ToMap (nil):", m)
}

func TestDoublePtr_FromMap(t *testing.T) {
	wds := WithDoublePtrStruct{}
	ds := NewDzStruct(wds)
	ds.FromMap(map[string]interface{}{
		"Name": "DoubleFromMap",
		"Inner": map[string]interface{}{
			"Name": "InnerVal",
			"ID":   55,
		},
	})
	fmt.Println("DoublePtr FromMap Name:", ds.Field("Name"))
	fmt.Println("DoublePtr FromMap Inner:", ds.Field("Inner"))
}

func TestDoublePtr_FromMap_NilValue(t *testing.T) {
	inner := PtrInner{Name: "Original", ID: 1}
	ptrInner := &inner
	wds := WithDoublePtrStruct{Name: "Outer", Inner: &ptrInner}
	ds := NewDzStruct(wds)
	ds.FromMap(map[string]interface{}{
		"Inner": nil,
	})
	fmt.Println("DoublePtr FromMap nil Inner:", ds.Field("Inner"))
}

// ==================== 指针匿名嵌入字段测试 ====================

func TestPtrEmbedded_Fields(t *testing.T) {
	inner := PtrInner{Name: "Embedded", ID: 10}
	pe := PtrEmbedded{PtrInner: &inner, Label: "Test"}
	ds := NewDzStruct(pe)
	fmt.Println("PtrEmbedded Fields:", ds.Fields())
	fmt.Println("PtrEmbedded Field Name:", ds.Field("Name"))
	fmt.Println("PtrEmbedded Field ID:", ds.Field("ID"))
	fmt.Println("PtrEmbedded Field Label:", ds.Field("Label"))
}

func TestPtrEmbedded_ToMap(t *testing.T) {
	inner := PtrInner{Name: "Embedded", ID: 10}
	pe := PtrEmbedded{PtrInner: &inner, Label: "Test"}
	ds := NewDzStruct(pe)
	m := ds.ToMap()
	fmt.Println("PtrEmbedded ToMap:", m)
}

func TestPtrEmbedded_NilField(t *testing.T) {
	pe := PtrEmbedded{PtrInner: nil, Label: "Test"}
	ds := NewDzStruct(pe)
	fmt.Println("PtrEmbedded NilField Fields:", ds.Fields())
	m := ds.ToMap()
	fmt.Println("PtrEmbedded NilField ToMap:", m)
}
