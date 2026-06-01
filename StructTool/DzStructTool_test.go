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

func TestClone_Basic(t *testing.T) {
	p := Person{Name: "Frank", Age: 40}
	ds := NewDzStruct(p)
	cloned := ds.Clone()
	fmt.Println("Original:", ds.ToInterface())
	fmt.Println("Cloned:", cloned.ToInterface())
	fmt.Println("Equal:", ds.EqualTo(cloned))
	// 修改克隆体不影响原始
	cloned.SetField("Name", "FrankClone")
	fmt.Println("After modify - Original Name:", ds.Field("Name"))
	fmt.Println("After modify - Cloned Name:", cloned.Field("Name"))
}

func TestClone_Nested(t *testing.T) {
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
	cloned := ds.Clone()
	fmt.Println("Nested Equal:", ds.EqualTo(cloned))
}

func TestClone_Slice(t *testing.T) {
	ws := WithSlice{Name: "Hank", Scores: []int{90, 85, 95}}
	ds := NewDzStruct(ws)
	cloned := ds.Clone()
	fmt.Println("Slice Equal:", ds.EqualTo(cloned))
}

func TestCloneReflect_Basic(t *testing.T) {
	p := Person{Name: "Ivy", Age: 27}
	ds := NewDzStruct(p)
	cloned := ds.CloneReflect()
	fmt.Println("CloneReflect Equal:", ds.EqualTo(cloned))
}

func TestCloneReflect_Slice(t *testing.T) {
	ws := WithSlice{Name: "Jack", Scores: []int{88, 92, 76}}
	ds := NewDzStruct(ws)
	cloned := ds.CloneReflect()
	fmt.Println("CloneReflect Slice Equal:", ds.EqualTo(cloned))
}

func TestCloneReflect_Map(t *testing.T) {
	wm := WithMap{Name: "Kate", Meta: map[string]string{"key": "value"}}
	ds := NewDzStruct(wm)
	cloned := ds.CloneReflect()
	fmt.Println("CloneReflect Map Equal:", ds.EqualTo(cloned))
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

func TestCloneStruct(t *testing.T) {
	p := Person{Name: "CloneSrc", Age: 50}
	cloned, err := CloneStruct(p)
	fmt.Println("CloneStruct err:", err)
	fmt.Println("CloneStruct result:", cloned)
	fmt.Println("CloneStruct equal:", reflect.DeepEqual(p, cloned))
}

func TestCloneStructReflect(t *testing.T) {
	ws := WithSlice{Name: "CloneSlice", Scores: []int{100, 99, 98}}
	cloned, err := CloneStructReflect(ws)
	fmt.Println("CloneStructReflect err:", err)
	fmt.Println("CloneStructReflect equal:", reflect.DeepEqual(ws, cloned))
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

	cloned := ds.Clone()
	fmt.Println("PointerField Clone Equal:", ds.EqualTo(cloned))
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

	cloned := ds.Clone()
	fmt.Println("AllZero Clone IsZero:", cloned.IsZero())
}

func TestMapWithNilPointer(t *testing.T) {
	wp := WithPointer{Name: nil, Age: 25, Score: 0}
	ds := NewDzStruct(wp)
	m := ds.ToMap()
	fmt.Println("WithNilPointer ToMap:", m)
}
