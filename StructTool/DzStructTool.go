package StructTool

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

// DzStruct 封装任意结构体类型，提供反射驱动的拷贝/克隆/比较/转换功能
type DzStruct struct {
	data interface{} // 持有的结构体值（总是非指针形式的副本）
	err  error       // 链式调用中累积的错误
}

// ==================== 构造函数 ====================

// NewDzStruct 包装任意结构体为 DzStruct。
// src 可以是结构体值或结构体指针。
// 内部始终存储值副本（若传入指针则解引用）。
func NewDzStruct(src interface{}) *DzStruct {
	if src == nil {
		return &DzStruct{err: errors.New("StructTool: src is nil")}
	}
	val := reflect.ValueOf(src)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return &DzStruct{err: errors.New("StructTool: src is nil pointer")}
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return &DzStruct{err: fmt.Errorf("StructTool: src is %s, not struct", val.Kind())}
	}
	return &DzStruct{data: val.Interface()}
}

// ==================== 状态检查 ====================

// Err 返回链式调用过程中累积的第一个错误
func (d *DzStruct) Err() error {
	return d.err
}

// IsValid 返回底层数据是否为有效的结构体
func (d *DzStruct) IsValid() bool {
	return d.err == nil && d.data != nil
}

// IsZero 返回底层结构体是否为零值（所有字段为零）
func (d *DzStruct) IsZero() bool {
	if d.err != nil {
		return true
	}
	val := reflect.ValueOf(d.data)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true
		}
		val = val.Elem()
	}
	return val.IsZero()
}

// ToInterface 返回底层存储的结构体值（以 interface{} 形式）
func (d *DzStruct) ToInterface() interface{} {
	if d.err != nil {
		return nil
	}
	return d.data
}

// ToIntf 将底层结构体值填充到 target 指向的结构体中。
// target 必须为结构体指针。
func (d *DzStruct) ToIntf(target interface{}) error {
	if d.err != nil {
		return d.err
	}
	if target == nil {
		return errors.New("StructTool: target is nil")
	}
	dstVal := reflect.ValueOf(target)
	if dstVal.Kind() != reflect.Ptr {
		return errors.New("StructTool: target must be a pointer")
	}
	if dstVal.IsNil() {
		return errors.New("StructTool: target is nil pointer")
	}
	srcVal := reflect.ValueOf(d.data)
	dstElem := dstVal.Elem()
	if dstElem.Kind() != reflect.Struct {
		return fmt.Errorf("StructTool: target must point to struct, got %s", dstElem.Kind())
	}
	if srcVal.Type() != dstElem.Type() {
		return fmt.Errorf("StructTool: type mismatch, src=%s, dst=%s", srcVal.Type(), dstElem.Type())
	}
	dstElem.Set(srcVal)
	return nil
}

// ==================== 内部辅助函数 ====================

// resolveValue 将 interface{} 解析为 reflect.Value（解引用指针到结构体值）
func resolveValue(src interface{}) (reflect.Value, error) {
	if src == nil {
		return reflect.Value{}, errors.New("StructTool: src is nil")
	}
	if ds, ok := src.(*DzStruct); ok {
		if ds.err != nil {
			return reflect.Value{}, ds.err
		}
		src = ds.data
	}
	val := reflect.ValueOf(src)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return reflect.Value{}, errors.New("StructTool: src is nil pointer")
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("StructTool: src is %s, not struct", val.Kind())
	}
	return val, nil
}

// getTagName 从 struct tag 中获取指定 tag 的名称部分（去掉 omitempty 等选项）
func getTagName(tag reflect.StructTag, tagName string) string {
	if tagName == "" {
		return ""
	}
	tagVal := tag.Get(tagName)
	if tagVal == "" {
		return ""
	}
	return strings.Split(tagVal, ",")[0]
}

// isExported 判断字段是否可导出
func isExported(field reflect.StructField) bool {
	return field.IsExported()
}

// fieldKey 确定字段在 map 中使用的 key。
// 如果指定了 tagName 且字段有对应 tag，则使用 tag 值；否则使用字段名。
func fieldKey(field reflect.StructField, tagName string) string {
	if tagName != "" {
		if tag := getTagName(field.Tag, tagName); tag != "" {
			return tag
		}
	}
	return field.Name
}

// setFieldWithCast 使用 cast 库进行类型转换后设置字段值
func setFieldWithCast(dstField reflect.Value, srcValue interface{}) bool {
	if !dstField.CanSet() {
		return false
	}
	srcVal := reflect.ValueOf(srcValue)
	// 如果类型相同且可直接赋值
	if srcVal.IsValid() && srcVal.Type() == dstField.Type() {
		dstField.Set(srcVal)
		return true
	}
	// 解引用指针
	for srcVal.IsValid() && srcVal.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return false
		}
		srcVal = srcVal.Elem()
	}
	// 处理目标为指针类型
	if dstField.Kind() == reflect.Ptr {
		if srcVal.IsValid() && srcVal.Type() == dstField.Type().Elem() {
			ptr := reflect.New(dstField.Type().Elem())
			ptr.Elem().Set(srcVal)
			dstField.Set(ptr)
			return true
		}
		return false
	}
	// 使用 cast 进行基本类型转换
	switch dstField.Kind() {
	case reflect.String:
		dstField.SetString(cast.ToString(srcValue))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dstField.SetInt(cast.ToInt64(srcValue))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dstField.SetUint(cast.ToUint64(srcValue))
	case reflect.Float32, reflect.Float64:
		dstField.SetFloat(cast.ToFloat64(srcValue))
	case reflect.Bool:
		dstField.SetBool(cast.ToBool(srcValue))
	default:
		// 其他类型（struct, slice, map 等）尝试直接赋值
		if srcVal.IsValid() && srcVal.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcVal)
			return true
		}
		return false
	}
	return true
}

// getAllFields 获取所有可导出字段（包含嵌入字段扁平化展开）
func getAllFields(val reflect.Value, tagName string) map[string]int {
	result := make(map[string]int)
	collectFields(val, "", tagName, result)
	return result
}

// collectFields 递归收集字段（嵌入字段扁平化）
func collectFields(val reflect.Value, prefix string, tagName string, result map[string]int) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if !isExported(field) {
			continue
		}
		fieldVal := val.Field(i)
		// 匿名嵌入字段：扁平化展开
		if field.Anonymous {
			// 解引用指针
			fv := fieldVal
			for fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					break
				}
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				collectFields(fv, prefix, tagName, result)
				continue
			}
		}
		key := fieldKey(field, tagName)
		if prefix != "" {
			key = prefix + key
		}
		result[key] = i
	}
}

// getFlatFields 获取扁平化的字段信息（字段名 -> 字段索引），返回字段名和字段值的映射
func getFlatFieldsIter(val reflect.Value, tagName string, callback func(key string, fieldVal reflect.Value, fieldTyp reflect.StructField)) {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if !isExported(field) {
			continue
		}
		fieldVal := val.Field(i)
		// 匿名嵌入字段：扁平化展开
		if field.Anonymous {
			fv := fieldVal
			for fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					break
				}
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				getFlatFieldsIter(fv, tagName, callback)
				continue
			}
		}
		key := fieldKey(field, tagName)
		callback(key, fieldVal, field)
	}
}

// deepCopyReflect 使用反射进行深拷贝
func deepCopyReflect(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Struct:
		result := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() {
				result.Field(i).Set(deepCopyReflect(v.Field(i)))
			}
		}
		return result
	case reflect.Ptr:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		cp := reflect.New(v.Elem().Type())
		cp.Elem().Set(deepCopyReflect(v.Elem()))
		return cp
	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		cp := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			cp.Index(i).Set(deepCopyReflect(v.Index(i)))
		}
		return cp
	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		cp := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			cp.SetMapIndex(deepCopyReflect(key), deepCopyReflect(v.MapIndex(key)))
		}
		return cp
	case reflect.Array:
		cp := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			cp.Index(i).Set(deepCopyReflect(v.Index(i)))
		}
		return cp
	case reflect.Interface:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		return deepCopyReflect(v.Elem())
	default:
		return v
	}
}

// structToMapRecursive 将结构体递归转换为 map（嵌入字段扁平化）
func structToMapRecursive(val reflect.Value, tagName string) map[string]interface{} {
	result := make(map[string]interface{})
	getFlatFieldsIter(val, tagName, func(key string, fieldVal reflect.Value, _ reflect.StructField) {
		// 解引用指针
		fv := fieldVal
		for fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				result[key] = nil
				return
			}
			fv = fv.Elem()
		}
		if fv.Kind() == reflect.Struct {
			// 对 time.Time 等特殊类型不做递归展开
			if fv.Type().String() == "time.Time" {
				result[key] = fv.Interface()
				return
			}
			result[key] = structToMapRecursive(fv, tagName)
		} else {
			result[key] = fv.Interface()
		}
	})
	return result
}

// mapToStructRecursive 将 map 填充到结构体（嵌入字段扁平化匹配）
func mapToStructRecursive(data map[string]interface{}, dstVal reflect.Value, tagName string) {
	if !dstVal.CanSet() {
		return
	}
	typ := dstVal.Type()
	for i := 0; i < dstVal.NumField(); i++ {
		field := typ.Field(i)
		if !isExported(field) {
			continue
		}
		fieldVal := dstVal.Field(i)
		// 匿名嵌入字段：递归填充
		if field.Anonymous {
			fv := fieldVal
			for fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					// 自动初始化指针
					fv.Set(reflect.New(fv.Type().Elem()))
					fv = fv.Elem()
				} else {
					fv = fv.Elem()
				}
			}
			if fv.Kind() == reflect.Struct {
				mapToStructRecursive(data, fv, tagName)
				continue
			}
		}
		key := fieldKey(field, tagName)
		mapVal, ok := data[key]
		if !ok {
			continue
		}
		// 处理嵌套结构体
		fv := fieldVal
		for fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
				fv = fv.Elem()
			} else {
				fv = fv.Elem()
			}
		}
		if fv.Kind() == reflect.Struct {
			if subMap, ok := mapVal.(map[string]interface{}); ok {
				mapToStructRecursive(subMap, fv, tagName)
				continue
			}
		}
		setFieldWithCast(fv, mapVal)
	}
}

// ==================== 同结构体操作 ====================

// Clone 深拷贝当前结构体。
// 使用 JSON 序列化/反序列化实现，能正确处理所有可 JSON 序列化的字段。
func (d *DzStruct) Clone() *DzStruct {
	if d.err != nil {
		return d
	}
	data, err := json.Marshal(d.data)
	if err != nil {
		return &DzStruct{err: fmt.Errorf("StructTool: Clone marshal: %w", err)}
	}
	srcType := reflect.TypeOf(d.data)
	target := reflect.New(srcType).Interface()
	if err := json.Unmarshal(data, target); err != nil {
		return &DzStruct{err: fmt.Errorf("StructTool: Clone unmarshal: %w", err)}
	}
	return &DzStruct{data: reflect.ValueOf(target).Elem().Interface()}
}

// CloneReflect 使用反射进行深拷贝，不依赖 JSON 序列化。
// 支持嵌套结构体、指针、切片、map。
func (d *DzStruct) CloneReflect() *DzStruct {
	if d.err != nil {
		return d
	}
	val := reflect.ValueOf(d.data)
	copied := deepCopyReflect(val)
	return &DzStruct{data: copied.Interface()}
}

// Fields 返回所有可导出字段的名称列表（嵌入字段扁平化）
func (d *DzStruct) Fields() []string {
	if d.err != nil {
		return nil
	}
	var result []string
	val := reflect.ValueOf(d.data)
	getFlatFieldsIter(val, "", func(key string, _ reflect.Value, _ reflect.StructField) {
		result = append(result, key)
	})
	return result
}

// Field 按名称获取指定字段的值（以 interface{} 形式）
func (d *DzStruct) Field(name string) interface{} {
	if d.err != nil {
		return nil
	}
	val := reflect.ValueOf(d.data)
	result := fieldByNameFlat(val, name)
	if !result.IsValid() {
		d.err = fmt.Errorf("StructTool: field %s not found", name)
		return nil
	}
	return result.Interface()
}

// fieldByNameFlat 按名称查找字段（支持嵌入字段扁平化）
func fieldByNameFlat(val reflect.Value, name string) reflect.Value {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if !isExported(field) {
			continue
		}
		if field.Anonymous {
			fv := val.Field(i)
			for fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					break
				}
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				if found := fieldByNameFlat(fv, name); found.IsValid() {
					return found
				}
			}
		}
		if field.Name == name {
			return val.Field(i)
		}
	}
	return reflect.Value{}
}

// SetField 按名称设置字段值。使用 cast 进行类型转换。
func (d *DzStruct) SetField(fieldName string, value interface{}) *DzStruct {
	if d.err != nil {
		return d
	}
	// 创建可修改的副本
	val := reflect.ValueOf(d.data)
	newVal := reflect.New(val.Type()).Elem()
	newVal.Set(val)
	field := fieldByNameFlat(newVal, fieldName)
	if !field.IsValid() {
		d.err = fmt.Errorf("StructTool: field %s not found or unexported", fieldName)
		return d
	}
	if !field.CanSet() {
		d.err = fmt.Errorf("StructTool: field %s cannot be set", fieldName)
		return d
	}
	if !setFieldWithCast(field, value) {
		d.err = fmt.Errorf("StructTool: cannot set field %s with value type %T", fieldName, value)
		return d
	}
	d.data = newVal.Interface()
	return d
}

// Zero 将结构体所有字段置为零值
func (d *DzStruct) Zero() *DzStruct {
	if d.err != nil {
		return d
	}
	srcType := reflect.TypeOf(d.data)
	d.data = reflect.New(srcType).Elem().Interface()
	return d
}

// ==================== 跨结构体操作 ====================

// CopyFrom 从 src 复制匹配字段到当前结构体。
// 匹配规则：字段名精确匹配（大小写敏感），使用 cast.ToXxx 进行类型转换。
// src 可以是结构体值/指针，也可以是 *DzStruct。
func (d *DzStruct) CopyFrom(src interface{}) *DzStruct {
	return d.copyFromInternal(src, "")
}

// CopyFromByTag 使用指定的 tag 名称匹配字段进行复制。
func (d *DzStruct) CopyFromByTag(src interface{}, tagName string) *DzStruct {
	return d.copyFromInternal(src, tagName)
}

// copyFromInternal 复制的内部实现
func (d *DzStruct) copyFromInternal(src interface{}, tagName string) *DzStruct {
	if d.err != nil {
		return d
	}
	srcVal, err := resolveValue(src)
	if err != nil {
		d.err = err
		return d
	}
	dstVal := reflect.ValueOf(d.data)
	// 创建可修改的副本
	newDst := reflect.New(dstVal.Type()).Elem()
	newDst.Set(dstVal)
	copyFieldsBetweenStructs(newDst, srcVal, tagName, false)
	d.data = newDst.Interface()
	return d
}

// MergeFrom 从 src 合并非零字段到当前结构体。
// 仅当 src 字段为非零值时才会覆盖目标字段。
func (d *DzStruct) MergeFrom(src interface{}) *DzStruct {
	return d.mergeFromInternal(src, "")
}

// MergeFromByTag 按 tag 名称匹配进行合并（非零字段覆盖）。
func (d *DzStruct) MergeFromByTag(src interface{}, tagName string) *DzStruct {
	return d.mergeFromInternal(src, tagName)
}

// mergeFromInternal 合并的内部实现
func (d *DzStruct) mergeFromInternal(src interface{}, tagName string) *DzStruct {
	if d.err != nil {
		return d
	}
	srcVal, err := resolveValue(src)
	if err != nil {
		d.err = err
		return d
	}
	dstVal := reflect.ValueOf(d.data)
	newDst := reflect.New(dstVal.Type()).Elem()
	newDst.Set(dstVal)
	copyFieldsBetweenStructs(newDst, srcVal, tagName, true)
	d.data = newDst.Interface()
	return d
}

// copyFieldsBetweenStructs 在两个不同结构体之间复制字段
// mergeOnlyNonZero=true 时仅复制源非零字段
func copyFieldsBetweenStructs(dstVal reflect.Value, srcVal reflect.Value, tagName string, mergeOnlyNonZero bool) {
	// 构建源字段映射（字段名/tag名 -> 字段索引）
	srcFieldMap := make(map[string]int)
	srcType := srcVal.Type()
	for i := 0; i < srcVal.NumField(); i++ {
		field := srcType.Field(i)
		if !isExported(field) {
			continue
		}
		// 匿名嵌入字段：扁平化
		if field.Anonymous {
			fv := srcVal.Field(i)
			for fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					break
				}
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				copyFieldsBetweenStructs(dstVal, fv, tagName, mergeOnlyNonZero)
				continue
			}
		}
		key := fieldKey(field, tagName)
		srcFieldMap[key] = i
	}

	// 遍历目标字段进行匹配复制
	dstType := dstVal.Type()
	for i := 0; i < dstVal.NumField(); i++ {
		field := dstType.Field(i)
		if !isExported(field) {
			continue
		}
		// 匿名嵌入字段：递归处理
		if field.Anonymous {
			fv := dstVal.Field(i)
			for fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					fv.Set(reflect.New(fv.Type().Elem()))
					fv = fv.Elem()
				} else {
					fv = fv.Elem()
				}
			}
			if fv.Kind() == reflect.Struct {
				copyFieldsBetweenStructs(fv, srcVal, tagName, mergeOnlyNonZero)
				continue
			}
		}
		dstKey := fieldKey(field, tagName)
		srcIdx, ok := srcFieldMap[dstKey]
		if !ok {
			continue
		}
		srcField := srcVal.Field(srcIdx)
		// mergeOnlyNonZero 时仅复制非零值
		if mergeOnlyNonZero && srcField.IsZero() {
			continue
		}
		dstField := dstVal.Field(i)
		setFieldWithCast(dstField, srcField.Interface())
	}
}

// CompareTo 与另一个结构体逐字段比较。
// 返回 map[string][2]interface{}，key 为字段名，
// value[0] 为当前结构体该字段的值，value[1] 为比较对象该字段的值。
// 仅包含存在差异的字段（嵌入字段扁平化）。
func (d *DzStruct) CompareTo(other interface{}) map[string][2]interface{} {
	if d.err != nil {
		return nil
	}
	otherVal, err := resolveValue(other)
	if err != nil {
		d.err = err
		return nil
	}
	result := make(map[string][2]interface{})
	dVal := reflect.ValueOf(d.data)
	compareStructsFlat(dVal, otherVal, result)
	return result
}

// compareStructsFlat 递归比较两个结构体（嵌入字段扁平化）
func compareStructsFlat(a, b reflect.Value, result map[string][2]interface{}) {
	aType := a.Type()
	for i := 0; i < a.NumField(); i++ {
		field := aType.Field(i)
		if !isExported(field) {
			continue
		}
		if field.Anonymous {
			af := a.Field(i)
			bf := b.Field(i)
			for af.Kind() == reflect.Ptr {
				if af.IsNil() {
					break
				}
				af = af.Elem()
			}
			for bf.Kind() == reflect.Ptr {
				if bf.IsNil() {
					break
				}
				bf = bf.Elem()
			}
			if af.Kind() == reflect.Struct && bf.Kind() == reflect.Struct {
				compareStructsFlat(af, bf, result)
				continue
			}
		}
		aField := a.Field(i)
		bField := b.Field(i)
		if !reflect.DeepEqual(aField.Interface(), bField.Interface()) {
			result[field.Name] = [2]interface{}{aField.Interface(), bField.Interface()}
		}
	}
}

// DiffFields 返回与另一个结构体差异字段的名称列表
func (d *DzStruct) DiffFields(other interface{}) []string {
	diff := d.CompareTo(other)
	var result []string
	for k := range diff {
		result = append(result, k)
	}
	return result
}

// EqualTo 判断两个结构体是否深度相等（所有可导出字段值相同）
func (d *DzStruct) EqualTo(other interface{}) bool {
	if d.err != nil {
		return false
	}
	otherVal, err := resolveValue(other)
	if err != nil {
		return false
	}
	dVal := reflect.ValueOf(d.data)
	return reflect.DeepEqual(dVal.Interface(), otherVal.Interface())
}

// ==================== 结构体-Map 互转 ====================

// ToMap 将结构体转换为 map[string]interface{}。
// map key 使用结构体字段名，嵌入字段扁平化展开。
// 嵌套结构体会递归展开为嵌套 map。
func (d *DzStruct) ToMap() map[string]interface{} {
	return d.toMapInternal("")
}

// ToMapByTag 使用指定 tag 名称的值作为 map key。
func (d *DzStruct) ToMapByTag(tagName string) map[string]interface{} {
	return d.toMapInternal(tagName)
}

// toMapInternal ToMap 的内部实现
func (d *DzStruct) toMapInternal(tagName string) map[string]interface{} {
	if d.err != nil {
		return nil
	}
	val := reflect.ValueOf(d.data)
	return structToMapRecursive(val, tagName)
}

// FromMap 从 map[string]interface{} 填充结构体字段。
// 优先匹配字段名，然后匹配 json tag。
// 值通过 cast 库进行类型转换。
func (d *DzStruct) FromMap(data map[string]interface{}) *DzStruct {
	return d.fromMapInternal(data, "")
}

// FromMapByTag 使用指定 tag 匹配 map key 到结构体字段。
func (d *DzStruct) FromMapByTag(data map[string]interface{}, tagName string) *DzStruct {
	return d.fromMapInternal(data, tagName)
}

// fromMapInternal FromMap 的内部实现
func (d *DzStruct) fromMapInternal(data map[string]interface{}, tagName string) *DzStruct {
	if d.err != nil {
		return d
	}
	val := reflect.ValueOf(d.data)
	newVal := reflect.New(val.Type()).Elem()
	newVal.Set(val)
	mapToStructRecursive(data, newVal, tagName)
	d.data = newVal.Interface()
	return d
}

// ==================== 静态工具函数 ====================

// CopyStruct 将 src 结构体的匹配字段复制到 dst 结构体。
// src, dst 必须为结构体指针。
func CopyStruct(src, dst interface{}) error {
	return copyStructByTag(src, dst, "")
}

// CopyStructByTag 按 tag 匹配规则复制字段。
func CopyStructByTag(src, dst interface{}, tagName string) error {
	return copyStructByTag(src, dst, tagName)
}

func copyStructByTag(src, dst interface{}, tagName string) error {
	srcVal, err := resolveValue(src)
	if err != nil {
		return err
	}
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return errors.New("StructTool: dst must be a pointer")
	}
	if dstVal.IsNil() {
		return errors.New("StructTool: dst is nil pointer")
	}
	dstElem := dstVal.Elem()
	if dstElem.Kind() != reflect.Struct {
		return fmt.Errorf("StructTool: dst must point to struct, got %s", dstElem.Kind())
	}
	copyFieldsBetweenStructs(dstElem, srcVal, tagName, false)
	return nil
}

// StructToMap 将任意结构体转换为 map。
func StructToMap(src interface{}) (map[string]interface{}, error) {
	val, err := resolveValue(src)
	if err != nil {
		return nil, err
	}
	return structToMapRecursive(val, ""), nil
}

// StructToMapByTag 按指定 tag 转换为 map。
func StructToMapByTag(src interface{}, tagName string) (map[string]interface{}, error) {
	val, err := resolveValue(src)
	if err != nil {
		return nil, err
	}
	return structToMapRecursive(val, tagName), nil
}

// MapToStruct 将 map 填充到结构体。
func MapToStruct(data map[string]interface{}, output interface{}) error {
	return mapToStructByTag(data, output, "")
}

// MapToStructByTag 按 tag 匹配填充。
func MapToStructByTag(data map[string]interface{}, output interface{}, tagName string) error {
	if output == nil {
		return errors.New("StructTool: output is nil")
	}
	outVal := reflect.ValueOf(output)
	if outVal.Kind() != reflect.Ptr {
		return errors.New("StructTool: output must be a pointer")
	}
	if outVal.IsNil() {
		return errors.New("StructTool: output is nil pointer")
	}
	dstElem := outVal.Elem()
	if dstElem.Kind() != reflect.Struct {
		return fmt.Errorf("StructTool: output must point to struct, got %s", dstElem.Kind())
	}
	mapToStructRecursive(data, dstElem, tagName)
	return nil
}

// mapToStructByTag 是 MapToStructByTag 的内部调用入口
func mapToStructByTag(data map[string]interface{}, output interface{}, tagName string) error {
	return MapToStructByTag(data, output, tagName)
}

// CloneStruct 深拷贝结构体（JSON 方式）。
func CloneStruct(src interface{}) (interface{}, error) {
	val, err := resolveValue(src)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(val.Interface())
	if err != nil {
		return nil, fmt.Errorf("StructTool: CloneStruct marshal: %w", err)
	}
	target := reflect.New(val.Type()).Interface()
	if err := json.Unmarshal(data, target); err != nil {
		return nil, fmt.Errorf("StructTool: CloneStruct unmarshal: %w", err)
	}
	return reflect.ValueOf(target).Elem().Interface(), nil
}

// CloneStructReflect 深拷贝结构体（反射方式）。
func CloneStructReflect(src interface{}) (interface{}, error) {
	val, err := resolveValue(src)
	if err != nil {
		return nil, err
	}
	copied := deepCopyReflect(val)
	return copied.Interface(), nil
}

// CompareStruct 比较两个结构体，返回差异字段名列表。
func CompareStruct(a, b interface{}) ([]string, error) {
	aVal, err := resolveValue(a)
	if err != nil {
		return nil, err
	}
	bVal, err := resolveValue(b)
	if err != nil {
		return nil, err
	}
	result := make(map[string][2]interface{})
	compareStructsFlat(aVal, bVal, result)
	var diffs []string
	for k := range result {
		diffs = append(diffs, k)
	}
	return diffs, nil
}
