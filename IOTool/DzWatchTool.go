package IOTool

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"path"
)

// WatchEventType 事件类型枚举
type WatchEventType int

const (
	WatchEventCreate WatchEventType = iota
	WatchEventWrite
	WatchEventRemove
	WatchEventRename
	WatchEventChmod
)

// WatchEvent 监听事件
type WatchEvent struct {
	Type  WatchEventType
	Path  string
	IsDir bool
}

// WatchObserver 观察者接口
type WatchObserver interface {
	OnEvent(event WatchEvent)
}

type pendingEvent struct {
	eventType WatchEventType
	isDir     bool
	timer     *time.Timer
}

// DzWatcher 文件监听器
type DzWatcher struct {
	paths     []string
	depth     int
	delay     time.Duration
	filter    func(path string) bool
	observers []WatchObserver
	mu        sync.RWMutex

	fw     *fsnotify.Watcher
	stopCh chan struct{}
	wg     sync.WaitGroup

	pendingMu sync.Mutex
	pending   map[string]*pendingEvent

	followMu    sync.Mutex
	followPaths map[string]bool
}

// NewWatcher 创建监听器
func NewWatcher(paths ...string) (*DzWatcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	absPaths := make([]string, 0, len(paths))
	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			fw.Close()
			return nil, err
		}
		absPaths = append(absPaths, abs)
	}
	return &DzWatcher{
		paths:       absPaths,
		depth:       -1,
		delay:       300 * time.Millisecond,
		filter:      nil,
		observers:   make([]WatchObserver, 0),
		fw:          fw,
		stopCh:      make(chan struct{}),
		pending:     make(map[string]*pendingEvent),
		followPaths: make(map[string]bool),
	}, nil
}

// SetDepth 设置监听深度，0=仅当前目录，-1=无限
func (w *DzWatcher) SetDepth(depth int) {
	w.depth = depth
}

// SetDelay 设置延迟合并时间
func (w *DzWatcher) SetDelay(duration time.Duration) {
	w.delay = duration
}

// SetFilter 设置 glob 过滤规则
func (w *DzWatcher) SetFilter(pattern string) {
	w.filter = func(p string) bool {
		matched, err := path.Match(pattern, filepath.Base(p))
		if err != nil {
			return false
		}
		return !matched
	}
}

// SetFilterFunc 设置自定义过滤函数，返回 true 表示忽略该路径
func (w *DzWatcher) SetFilterFunc(fn func(path string) bool) {
	w.filter = fn
}

// AddObserver 添加观察者
func (w *DzWatcher) AddObserver(observer WatchObserver) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.observers = append(w.observers, observer)
}

// RemoveObserver 移除观察者
func (w *DzWatcher) RemoveObserver(observer WatchObserver) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i, obs := range w.observers {
		if obs == observer {
			w.observers = append(w.observers[:i], w.observers[i+1:]...)
			return
		}
	}
}

// Start 启动监听
func (w *DzWatcher) Start() error {
	for _, p := range w.paths {
		if err := w.registerRecursive(p, 0); err != nil {
			return err
		}
	}
	w.wg.Add(1)
	go w.eventLoop()
	return nil
}

// Stop 停止监听
func (w *DzWatcher) Stop() {
	close(w.stopCh)
	w.fw.Close()
	w.wg.Wait()

	w.pendingMu.Lock()
	for _, pe := range w.pending {
		pe.timer.Stop()
	}
	w.pending = make(map[string]*pendingEvent)
	w.pendingMu.Unlock()
}

func (w *DzWatcher) registerRecursive(root string, currentDepth int) error {
	if w.depth >= 0 && currentDepth > w.depth {
		return nil
	}
	if err := w.fw.Add(root); err != nil {
		return err
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			fullPath := filepath.Join(root, entry.Name())
			if err := w.registerRecursive(fullPath, currentDepth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *DzWatcher) eventLoop() {
	defer w.wg.Done()
	for {
		select {
		case <-w.stopCh:
			return
		case event, ok := <-w.fw.Events:
			if !ok {
				return
			}
			w.handleEvent(event)
		case _, ok := <-w.fw.Errors:
			if !ok {
				return
			}
		}
	}
}

func (w *DzWatcher) handleEvent(event fsnotify.Event) {
	p := event.Name
	if w.filter != nil && w.filter(p) {
		return
	}
	eventType := convertEventType(event.Op)
	isDir := isDirPath(p)

	if eventType == WatchEventCreate && isDir {
		depth := w.pathDepth(p)
		w.registerRecursive(p, depth)
	}

	w.pendingMu.Lock()
	defer w.pendingMu.Unlock()

	if pe, ok := w.pending[p]; ok {
		pe.timer.Stop()
		eventType = mergeEventType(pe.eventType, eventType)
	}

	if eventType == WatchEventRemove && isDir {
		w.fw.Remove(p)
	}

	if eventType == WatchEventType(-1) {
		delete(w.pending, p)
		return
	}

	timer := time.AfterFunc(w.delay, func() {
		w.pendingMu.Lock()
		delete(w.pending, p)
		w.pendingMu.Unlock()
		w.notifyObservers(WatchEvent{Type: eventType, Path: p, IsDir: isDir})

		if eventType == WatchEventCreate || eventType == WatchEventWrite {
			w.followMu.Lock()
			if w.followPaths[p] {
				w.fw.Add(p)
			}
			w.followMu.Unlock()
		}
	})

	w.pending[p] = &pendingEvent{eventType: eventType, isDir: isDir, timer: timer}
}

func (w *DzWatcher) notifyObservers(event WatchEvent) {
	w.mu.RLock()
	observers := make([]WatchObserver, len(w.observers))
	copy(observers, w.observers)
	w.mu.RUnlock()
	for _, obs := range observers {
		obs.OnEvent(event)
	}
}

// Follow 文件跟随
func (w *DzWatcher) Follow(path string, observer WatchObserver) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}
	w.followMu.Lock()
	w.followPaths[abs] = true
	w.followMu.Unlock()
	w.AddObserver(observer)
	w.fw.Add(abs)
}

// Watch 一行代码监听
func Watch(path string, onEvent func(event WatchEvent)) (*DzWatcher, error) {
	w, err := NewWatcher(path)
	if err != nil {
		return nil, err
	}
	w.AddObserver(&funcObserver{fn: onEvent})
	if err := w.Start(); err != nil {
		w.Stop()
		return nil, err
	}
	return w, nil
}

type funcObserver struct {
	fn func(event WatchEvent)
}

func (f *funcObserver) OnEvent(event WatchEvent) {
	f.fn(event)
}

func convertEventType(op fsnotify.Op) WatchEventType {
	switch {
	case op&fsnotify.Create != 0:
		return WatchEventCreate
	case op&fsnotify.Write != 0:
		return WatchEventWrite
	case op&fsnotify.Remove != 0:
		return WatchEventRemove
	case op&fsnotify.Rename != 0:
		return WatchEventRename
	case op&fsnotify.Chmod != 0:
		return WatchEventChmod
	default:
		return WatchEventWrite
	}
}

func isDirPath(p string) bool {
	info, err := os.Lstat(p)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func mergeEventType(old, new WatchEventType) WatchEventType {
	if old == WatchEventCreate && new == WatchEventWrite {
		return WatchEventCreate
	}
	if old == WatchEventCreate && new == WatchEventRemove {
		return -1
	}
	if old == WatchEventRemove && new == WatchEventCreate {
		return WatchEventWrite
	}
	return new
}

func (w *DzWatcher) pathDepth(p string) int {
	for _, root := range w.paths {
		rel, err := filepath.Rel(root, p)
		if err != nil {
			continue
		}
		if rel == "." {
			return 0
		}
		if !strings.HasPrefix(rel, "..") {
			count := 0
			for _, c := range rel {
				if c == os.PathSeparator {
					count++
				}
			}
			return count + 1
		}
	}
	return 0
}


