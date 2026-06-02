package dzUtils

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// ==================== UUID 测试 ====================

func TestNewUUID(t *testing.T) {
	id := NewUUID()
	if id == "" {
		t.Fatal("NewUUID() 返回空字符串")
	}
	if len(id) != 36 {
		t.Errorf("NewUUID() 长度 = %d, 期望 36", len(id))
	}
	if strings.Count(id, "-") != 4 {
		t.Errorf("NewUUID() 连字符数量 = %d, 期望 4", strings.Count(id, "-"))
	}
	fmt.Println("NewUUID:", id)
}

func TestNewUUIDv1(t *testing.T) {
	id, err := NewUUIDv1()
	if err != nil {
		t.Fatalf("NewUUIDv1() 失败: %v", err)
	}
	if len(id) != 36 {
		t.Errorf("NewUUIDv1() 长度 = %d, 期望 36", len(id))
	}
	fmt.Println("NewUUIDv1:", id)
}

func TestNewUUIDv4(t *testing.T) {
	id := NewUUIDv4()
	if len(id) != 36 {
		t.Errorf("NewUUIDv4() 长度 = %d, 期望 36", len(id))
	}
	fmt.Println("NewUUIDv4:", id)
}

func TestNewUUIDv7(t *testing.T) {
	id, err := NewUUIDv7()
	if err != nil {
		t.Fatalf("NewUUIDv7() 失败: %v", err)
	}
	if len(id) != 36 {
		t.Errorf("NewUUIDv7() 长度 = %d, 期望 36", len(id))
	}
	fmt.Println("NewUUIDv7:", id)
}

func TestNewUUIDWithoutDash(t *testing.T) {
	id := NewUUIDWithoutDash()
	if len(id) != 32 {
		t.Errorf("NewUUIDWithoutDash() 长度 = %d, 期望 32", len(id))
	}
	if strings.Contains(id, "-") {
		t.Error("NewUUIDWithoutDash() 不应包含连字符")
	}
	fmt.Println("NewUUIDWithoutDash:", id)
}

func TestUUIDUniqueness(t *testing.T) {
	const count = 1000
	seen := make(map[string]struct{}, count)
	for i := 0; i < count; i++ {
		id := NewUUID()
		if _, exists := seen[id]; exists {
			t.Errorf("UUID 重复: %s", id)
		}
		seen[id] = struct{}{}
	}
	fmt.Printf("UUID 唯一性测试通过 (%d 个无重复)\n", count)
}

func TestUUIDv7TimeOrdering(t *testing.T) {
	var ids []string
	for i := 0; i < 10; i++ {
		id, err := NewUUIDv7()
		if err != nil {
			t.Fatalf("NewUUIDv7() 失败: %v", err)
		}
		ids = append(ids, id)
	}
	// v7 UUID 应该是时间有序的（前48位包含毫秒时间戳）
	for i := 1; i < len(ids); i++ {
		if ids[i] < ids[i-1] {
			t.Errorf("UUID v7 时间排序异常: %s < %s", ids[i], ids[i-1])
		}
	}
	fmt.Println("UUID v7 时间排序测试通过")
}

func TestParseUUID(t *testing.T) {
	original := NewUUID()
	parsed, err := ParseUUID(original)
	if err != nil {
		t.Fatalf("ParseUUID() 失败: %v", err)
	}
	if parsed.String() != original {
		t.Errorf("ParseUUID(%q) = %q, 期望 %q", original, parsed.String(), original)
	}

	// 无效UUID
	_, err = ParseUUID("not-a-uuid")
	if err == nil {
		t.Error("ParseUUID(\"not-a-uuid\") 期望返回错误")
	}
	fmt.Println("ParseUUID 测试通过")
}

func TestUUIDIsValid(t *testing.T) {
	validUUID := NewUUID()
	if !UUIDIsValid(validUUID) {
		t.Errorf("UUIDIsValid(%q) = false, 期望 true", validUUID)
	}

	if UUIDIsValid("invalid") {
		t.Error("UUIDIsValid(\"invalid\") = true, 期望 false")
	}

	if UUIDIsValid("") {
		t.Error("UUIDIsValid(\"\") = true, 期望 false")
	}
	fmt.Println("UUIDIsValid 测试通过")
}

// ==================== ObjectId 测试 ====================

func TestNewObjectId(t *testing.T) {
	id := NewObjectId()
	if len(id) != 24 {
		t.Errorf("NewObjectId() 长度 = %d, 期望 24", len(id))
	}
	fmt.Println("NewObjectId:", id)
}

func TestNewObjectIdFromTime(t *testing.T) {
	now := time.Now()
	id := NewObjectIdFromTime(now)
	if len(id) != 24 {
		t.Errorf("NewObjectIdFromTime() 长度 = %d, 期望 24", len(id))
	}

	// 提取的时间戳应该与输入时间一致（秒级精度）
	ts, err := ObjectIdTimestamp(id)
	if err != nil {
		t.Fatalf("ObjectIdTimestamp() 失败: %v", err)
	}
	diff := ts.Unix() - now.Unix()
	if diff < -1 || diff > 1 {
		t.Errorf("ObjectId时间戳差异过大: 输入 %v, 提取 %v", now.Unix(), ts.Unix())
	}
	fmt.Println("NewObjectIdFromTime:", id, "时间戳:", ts)
}

func TestObjectIdTimestamp(t *testing.T) {
	knownTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	id := NewObjectIdFromTime(knownTime)

	ts, err := ObjectIdTimestamp(id)
	if err != nil {
		t.Fatalf("ObjectIdTimestamp() 失败: %v", err)
	}
	if ts.Unix() != knownTime.Unix() {
		t.Errorf("ObjectIdTimestamp() = %v, 期望 %v", ts.Unix(), knownTime.Unix())
	}
	fmt.Println("ObjectIdTimestamp 测试通过")
}

func TestObjectIdIsValid(t *testing.T) {
	validID := NewObjectId()
	if !ObjectIdIsValid(validID) {
		t.Errorf("ObjectIdIsValid(%q) = false, 期望 true", validID)
	}

	// 长度不对
	if ObjectIdIsValid("abc") {
		t.Error("ObjectIdIsValid(\"abc\") = true, 期望 false")
	}

	// 非十六进制字符
	if ObjectIdIsValid("zzzzzzzzzzzzzzzzzzzzzzzz") {
		t.Error("ObjectIdIsValid(非十六进制) = true, 期望 false")
	}

	// 空字符串
	if ObjectIdIsValid("") {
		t.Error("ObjectIdIsValid(\"\") = true, 期望 false")
	}
	fmt.Println("ObjectIdIsValid 测试通过")
}

func TestObjectIdTimestampError(t *testing.T) {
	// 长度不对
	_, err := ObjectIdTimestamp("abc")
	if err == nil {
		t.Error("ObjectIdTimestamp(\"abc\") 期望返回错误")
	}

	// 非十六进制
	_, err = ObjectIdTimestamp("zzzzzzzzzzzzzzzzzzzzzzzz")
	if err == nil {
		t.Error("ObjectIdTimestamp(非十六进制) 期望返回错误")
	}
	fmt.Println("ObjectIdTimestamp 错误处理测试通过")
}

func TestObjectIdUniqueness(t *testing.T) {
	const count = 1000
	seen := make(map[string]struct{}, count)
	for i := 0; i < count; i++ {
		id := NewObjectId()
		if _, exists := seen[id]; exists {
			t.Errorf("ObjectId 重复: %s", id)
		}
		seen[id] = struct{}{}
	}
	fmt.Printf("ObjectId 唯一性测试通过 (%d 个无重复)\n", count)
}

func TestObjectIdMonotonic(t *testing.T) {
	// 同一秒内生成的ObjectId应该单调递增
	var prev string
	for i := 0; i < 100; i++ {
		id := NewObjectId()
		if prev != "" && id <= prev {
			t.Errorf("ObjectId 非单调递增: %q <= %q", id, prev)
		}
		prev = id
	}
	fmt.Println("ObjectId 单调递增测试通过")
}

// ==================== Snowflake 测试 ====================

func TestInitSnowflake(t *testing.T) {
	err := InitSnowflake(1)
	if err != nil {
		t.Fatalf("InitSnowflake(1) 失败: %v", err)
	}

	// 无效节点ID（超过1023）
	err = InitSnowflake(1024)
	if err == nil {
		t.Error("InitSnowflake(1024) 期望返回错误")
	}
	fmt.Println("InitSnowflake 测试通过")
}

func TestNewSnowflakeID(t *testing.T) {
	_ = InitSnowflake(1)

	id := NewSnowflakeID()
	if id <= 0 {
		t.Errorf("NewSnowflakeID() = %d, 期望 > 0", id)
	}
	fmt.Println("NewSnowflakeID:", id)
}

func TestNewSnowflakeString(t *testing.T) {
	_ = InitSnowflake(1)

	s := NewSnowflakeString()
	if s == "" {
		t.Error("NewSnowflakeString() 返回空字符串")
	}
	if s == "0" {
		t.Error("NewSnowflakeString() 不应返回 \"0\"")
	}
	fmt.Println("NewSnowflakeString:", s)
}

func TestSnowflakeTime(t *testing.T) {
	_ = InitSnowflake(1)

	before := time.Now()
	id := NewSnowflakeID()
	after := time.Now()

	ts := SnowflakeTime(id)
	// Snowflake时间应在生成ID前后之间（允许1ms误差）
	if ts.Before(before.Add(-2*time.Millisecond)) || ts.After(after.Add(2*time.Millisecond)) {
		t.Errorf("SnowflakeTime() = %v, 期望在 [%v, %v] 范围内", ts, before, after)
	}
	fmt.Println("SnowflakeTime:", ts)
}

func TestSnowflakeNodeID(t *testing.T) {
	const testNode int64 = 42
	_ = InitSnowflake(testNode)

	id := NewSnowflakeID()
	nodeID := SnowflakeNodeID(id)
	if nodeID != testNode {
		t.Errorf("SnowflakeNodeID() = %d, 期望 %d", nodeID, testNode)
	}
	fmt.Println("SnowflakeNodeID:", nodeID)
}

func TestSnowflakeStep(t *testing.T) {
	_ = InitSnowflake(1)

	id := NewSnowflakeID()
	step := SnowflakeStep(id)
	// Step应在 [0, 4095] 范围内
	if step < 0 || step > 4095 {
		t.Errorf("SnowflakeStep() = %d, 超出有效范围 [0, 4095]", step)
	}
	fmt.Println("SnowflakeStep:", step)
}

func TestDecomposeSnowflake(t *testing.T) {
	const testNode int64 = 100
	_ = InitSnowflake(testNode)

	id := NewSnowflakeID()
	parts := DecomposeSnowflake(id)

	if parts.ID != id {
		t.Errorf("DecomposeSnowflake().ID = %d, 期望 %d", parts.ID, id)
	}
	if parts.NodeID != testNode {
		t.Errorf("DecomposeSnowflake().NodeID = %d, 期望 %d", parts.NodeID, testNode)
	}
	if parts.Time.IsZero() {
		t.Error("DecomposeSnowflake().Time 不应为零值")
	}
	fmt.Printf("DecomposeSnowflake: ID=%d, NodeID=%d, Step=%d, Time=%v\n",
		parts.ID, parts.NodeID, parts.Step, parts.Time)
}

func TestSnowflakeUniqueness(t *testing.T) {
	_ = InitSnowflake(1)

	const count = 5000
	seen := make(map[int64]struct{}, count)
	for i := 0; i < count; i++ {
		id := NewSnowflakeID()
		if _, exists := seen[id]; exists {
			t.Errorf("Snowflake ID 重复: %d", id)
		}
		seen[id] = struct{}{}
	}
	fmt.Printf("Snowflake 唯一性测试通过 (%d 个无重复)\n", count)
}

func TestSnowflakeMonotonic(t *testing.T) {
	_ = InitSnowflake(1)

	var prev int64
	for i := 0; i < 1000; i++ {
		id := NewSnowflakeID()
		if id <= prev {
			t.Errorf("Snowflake ID 非单调递增: %d <= %d", id, prev)
		}
		prev = id
	}
	fmt.Println("Snowflake 单调递增测试通过")
}

func TestSetSnowflakeEpoch(t *testing.T) {
	// 注意：SetSnowflakeEpoch 只影响后续创建的Node
	// 由于 Snowflake 的 Epoch 是包级变量，此测试仅验证函数不报错
	customEpoch := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	SetSnowflakeEpoch(customEpoch)
	// 重新初始化节点以使用新Epoch
	_ = InitSnowflake(1)

	id := NewSnowflakeID()
	ts := SnowflakeTime(id)
	// 时间应基于自定义Epoch计算，应该接近当前时间
	now := time.Now()
	diff := now.Sub(ts)
	if diff < 0 {
		diff = -diff
	}
	// 允许较大误差，因为Epoch偏移会影响时间计算
	if diff > 5*time.Second {
		t.Errorf("SnowflakeTime 与当前时间差异过大: %v", diff)
	}
	fmt.Println("SetSnowflakeEpoch 测试通过, 时间:", ts)

	// 恢复默认Epoch
	SetSnowflakeEpoch(time.Date(2010, 11, 4, 1, 42, 54, 0, time.UTC))
	_ = InitSnowflake(1)
}

// ==================== 并发安全测试 ====================

func TestUUIDConcurrent(t *testing.T) {
	const goroutines = 10
	const count = 100

	var wg sync.WaitGroup
	seen := make(map[string]struct{})
	var mu sync.Mutex

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < count; i++ {
				id := NewUUID()
				mu.Lock()
				if _, exists := seen[id]; exists {
					t.Errorf("并发UUID重复: %s", id)
				}
				seen[id] = struct{}{}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("UUID 并发测试通过 (%d goroutines × %d = %d 无重复)\n", goroutines, count, goroutines*count)
}

func TestObjectIdConcurrent(t *testing.T) {
	const goroutines = 10
	const count = 100

	var wg sync.WaitGroup
	seen := make(map[string]struct{})
	var mu sync.Mutex

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < count; i++ {
				id := NewObjectId()
				mu.Lock()
				if _, exists := seen[id]; exists {
					t.Errorf("并发ObjectId重复: %s", id)
				}
				seen[id] = struct{}{}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("ObjectId 并发测试通过 (%d goroutines × %d = %d 无重复)\n", goroutines, count, goroutines*count)
}

func TestSnowflakeConcurrent(t *testing.T) {
	_ = InitSnowflake(1)

	const goroutines = 10
	const count = 100

	var wg sync.WaitGroup
	seen := make(map[int64]struct{})
	var mu sync.Mutex

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < count; i++ {
				id := NewSnowflakeID()
				mu.Lock()
				if _, exists := seen[id]; exists {
					t.Errorf("并发Snowflake ID重复: %d", id)
				}
				seen[id] = struct{}{}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Snowflake 并发测试通过 (%d goroutines × %d = %d 无重复)\n", goroutines, count, goroutines*count)
}
