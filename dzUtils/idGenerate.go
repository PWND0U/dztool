package dzUtils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

// ==================== UUID ====================

// NewUUID 生成UUID v4字符串（最常用，基于随机数）
func NewUUID() string {
	return uuid.New().String()
}

// NewUUIDv1 生成UUID v1字符串（基于时间戳和MAC地址）
func NewUUIDv1() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("dzUtils: 生成UUID v1失败: %w", err)
	}
	return id.String(), nil
}

// NewUUIDv4 生成UUID v4字符串（基于随机数）
func NewUUIDv4() string {
	return uuid.New().String()
}

// NewUUIDv7 生成UUID v7字符串（基于时间戳，时间有序）
func NewUUIDv7() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("dzUtils: 生成UUID v7失败: %w", err)
	}
	return id.String(), nil
}

// NewUUIDWithoutDash 生成无连字符的UUID v4字符串（32位十六进制）
func NewUUIDWithoutDash() string {
	id := uuid.New()
	return hex.EncodeToString(id[:])
}

// ParseUUID 解析UUID字符串
func ParseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dzUtils: 解析UUID失败 %q: %w", s, err)
	}
	return id, nil
}

// UUIDIsValid 验证UUID字符串是否有效
func UUIDIsValid(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// ==================== ObjectId ====================

// ObjectId MongoDB兼容的ObjectId（12字节 = 24位十六进制字符串）
// 结构：4字节时间戳 + 5字节随机进程标识 + 3字节递增计数器

var (
	objectIdRand [5]byte // 5字节随机进程标识
	objectIdCnt  uint32  // 3字节递增计数器
	objectIdOnce sync.Once
)

// initObjectId 延迟初始化ObjectId的随机进程标识和计数器
func initObjectId() {
	objectIdOnce.Do(func() {
		_, _ = rand.Read(objectIdRand[:])
		var buf [4]byte
		_, _ = rand.Read(buf[:])
		objectIdCnt = uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	})
}

// objectIdEncode 将12字节编码为24位十六进制字符串
func objectIdEncode(ts uint32, randBytes [5]byte, cnt uint32) string {
	var buf [12]byte
	// 4字节时间戳（大端序）
	buf[0] = byte(ts >> 24)
	buf[1] = byte(ts >> 16)
	buf[2] = byte(ts >> 8)
	buf[3] = byte(ts)
	// 5字节随机进程标识
	copy(buf[4:9], randBytes[:])
	// 3字节递增计数器（仅取低24位）
	cnt &= 0x00FFFFFF
	buf[9] = byte(cnt >> 16)
	buf[10] = byte(cnt >> 8)
	buf[11] = byte(cnt)
	return hex.EncodeToString(buf[:])
}

// NewObjectId 生成MongoDB兼容的ObjectId字符串（24位十六进制）
func NewObjectId() string {
	initObjectId()
	ts := uint32(time.Now().Unix())
	cnt := atomic.AddUint32(&objectIdCnt, 1)
	return objectIdEncode(ts, objectIdRand, cnt)
}

// NewObjectIdFromTime 根据指定时间生成ObjectId
func NewObjectIdFromTime(t time.Time) string {
	initObjectId()
	ts := uint32(t.Unix())
	cnt := atomic.AddUint32(&objectIdCnt, 1)
	return objectIdEncode(ts, objectIdRand, cnt)
}

// ObjectIdTimestamp 从ObjectId中提取时间戳
func ObjectIdTimestamp(id string) (time.Time, error) {
	if len(id) != 24 {
		return time.Time{}, fmt.Errorf("dzUtils: ObjectId长度无效，期望24，实际 %d", len(id))
	}
	b, err := hex.DecodeString(id)
	if err != nil {
		return time.Time{}, fmt.Errorf("dzUtils: 解码ObjectId失败: %w", err)
	}
	if len(b) != 12 {
		return time.Time{}, fmt.Errorf("dzUtils: ObjectId字节数无效，期望12，实际 %d", len(b))
	}
	ts := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return time.Unix(int64(ts), 0), nil
}

// ObjectIdIsValid 验证ObjectId字符串是否有效
func ObjectIdIsValid(id string) bool {
	if len(id) != 24 {
		return false
	}
	_, err := hex.DecodeString(id)
	return err == nil
}

// ==================== Snowflake ====================

var (
	sfNode   *snowflake.Node
	sfNodeMu sync.RWMutex
	sfOnce   sync.Once
)

// InitSnowflake 初始化Snowflake节点
// 参数 nodeID: 节点ID，范围 [0, 1023]
// 必须在使用其他Snowflake函数之前调用，否则将使用默认节点1
func InitSnowflake(nodeID int64) error {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return fmt.Errorf("dzUtils: 初始化Snowflake节点失败: %w", err)
	}
	sfNodeMu.Lock()
	sfNode = node
	sfNodeMu.Unlock()
	return nil
}

// SetSnowflakeEpoch 设置Snowflake纪元时间（必须在InitSnowflake之前调用）
func SetSnowflakeEpoch(t time.Time) {
	snowflake.Epoch = t.UnixMilli()
}

// ensureSnowflakeNode 确保Snowflake节点已初始化（默认使用节点1）
func ensureSnowflakeNode() *snowflake.Node {
	sfOnce.Do(func() {
		sfNodeMu.RLock()
		node := sfNode
		sfNodeMu.RUnlock()
		if node != nil {
			return
		}
		n, err := snowflake.NewNode(1)
		if err != nil {
			panic(fmt.Sprintf("dzUtils: 初始化默认Snowflake节点失败: %v", err))
		}
		sfNodeMu.Lock()
		sfNode = n
		sfNodeMu.Unlock()
	})
	sfNodeMu.RLock()
	defer sfNodeMu.RUnlock()
	return sfNode
}

// NewSnowflakeID 生成Snowflake ID（int64）
func NewSnowflakeID() int64 {
	return ensureSnowflakeNode().Generate().Int64()
}

// NewSnowflakeString 生成Snowflake ID字符串
func NewSnowflakeString() string {
	return ensureSnowflakeNode().Generate().String()
}

// SnowflakeTime 从Snowflake ID中提取生成时间
func SnowflakeTime(id int64) time.Time {
	ms := snowflake.ID(id).Time()
	return time.UnixMilli(ms)
}

// SnowflakeNodeID 从Snowflake ID中提取节点ID
func SnowflakeNodeID(id int64) int64 {
	return snowflake.ID(id).Node()
}

// SnowflakeStep 从Snowflake ID中提取序列号
func SnowflakeStep(id int64) int64 {
	return snowflake.ID(id).Step()
}

// SnowflakeParts Snowflake ID分解后的各组成部分
type SnowflakeParts struct {
	ID     int64     // 完整ID
	NodeID int64     // 节点ID
	Step   int64     // 序列号
	Time   time.Time // 生成时间
}

// DecomposeSnowflake 分解Snowflake ID为各组成部分
func DecomposeSnowflake(id int64) *SnowflakeParts {
	sfID := snowflake.ID(id)
	return &SnowflakeParts{
		ID:     id,
		NodeID: sfID.Node(),
		Step:   sfID.Step(),
		Time:   time.UnixMilli(sfID.Time()),
	}
}
