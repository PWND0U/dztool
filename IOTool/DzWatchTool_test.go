package IOTool

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	dir, err := os.MkdirTemp("", "dzwatch_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	w, err := NewWatcher(dir)
	if err != nil {
		t.Fatalf("NewWatcher 失败: %v", err)
	}
	defer w.Stop()

	if len(w.paths) != 1 {
		t.Fatalf("期望 paths 长度 1, 实际 %d", len(w.paths))
	}
	if w.depth != -1 {
		t.Fatalf("默认深度应为 -1, 实际 %d", w.depth)
	}
	if w.delay != 300*time.Millisecond {
		t.Fatalf("默认延迟应为 300ms, 实际 %v", w.delay)
	}
}

func TestSetDepth(t *testing.T) {
	w, err := NewWatcher(".")
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	w.SetDepth(2)
	if w.depth != 2 {
		t.Fatalf("期望深度 2, 实际 %d", w.depth)
	}

	w.SetDepth(0)
	if w.depth != 0 {
		t.Fatalf("期望深度 0, 实际 %d", w.depth)
	}
}

func TestSetDelay(t *testing.T) {
	w, err := NewWatcher(".")
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	w.SetDelay(100 * time.Millisecond)
	if w.delay != 100*time.Millisecond {
		t.Fatalf("期望延迟 100ms, 实际 %v", w.delay)
	}
}

func TestAddRemoveObserver(t *testing.T) {
	w, err := NewWatcher(".")
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	obs := &testObserver{}
	w.AddObserver(obs)
	if len(w.observers) != 1 {
		t.Fatalf("添加后期望 1 个观察者, 实际 %d", len(w.observers))
	}

	w.RemoveObserver(obs)
	if len(w.observers) != 0 {
		t.Fatalf("移除后期望 0 个观察者, 实际 %d", len(w.observers))
	}
}

func TestStartStop(t *testing.T) {
	dir, err := os.MkdirTemp("", "dzwatch_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	w, err := NewWatcher(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := w.Start(); err != nil {
		t.Fatalf("Start 失败: %v", err)
	}
	w.Stop()
}

func TestFileCreateEvent(t *testing.T) {
	dir, err := os.MkdirTemp("", "dzwatch_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	w, err := NewWatcher(dir)
	if err != nil {
		t.Fatal(err)
	}
	w.SetDelay(100 * time.Millisecond)

	obs := &testObserver{ch: make(chan WatchEvent, 10)}
	w.AddObserver(obs)

	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	testFile := filepath.Join(dir, "test_create.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0666); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-obs.ch:
		if event.Type != WatchEventCreate {
			t.Fatalf("期望 Create 事件, 实际 %v", event.Type)
		}
		if event.Path != testFile {
			t.Fatalf("路径不匹配: 期望 %s, 实际 %s", testFile, event.Path)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("等待事件超时")
	}
}

func TestDelayMerge(t *testing.T) {
	dir, err := os.MkdirTemp("", "dzwatch_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	w, err := NewWatcher(dir)
	if err != nil {
		t.Fatal(err)
	}
	w.SetDelay(200 * time.Millisecond)

	obs := &testObserver{ch: make(chan WatchEvent, 10)}
	w.AddObserver(obs)

	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	testFile := filepath.Join(dir, "test_merge.txt")
	if err := os.WriteFile(testFile, []byte("first"), 0666); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
	if err := os.WriteFile(testFile, []byte("second"), 0666); err != nil {
		t.Fatal(err)
	}

	var events []WatchEvent
	timeout := time.After(2 * time.Second)
collect:
	for {
		select {
		case e := <-obs.ch:
			events = append(events, e)
		case <-timeout:
			break collect
		}
	}

	if len(events) != 1 {
		t.Fatalf("延迟合并后期望 1 个事件, 实际收到 %d 个", len(events))
	}
	if events[0].Type != WatchEventCreate {
		t.Fatalf("Create+Write 合并后应为 Create, 实际 %v", events[0].Type)
	}
}

func TestSubDirWatch(t *testing.T) {
	dir, err := os.MkdirTemp("", "dzwatch_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	subDir := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(subDir, 0777); err != nil {
		t.Fatal(err)
	}

	w, err := NewWatcher(dir)
	if err != nil {
		t.Fatal(err)
	}
	w.SetDelay(100 * time.Millisecond)

	obs := &testObserver{ch: make(chan WatchEvent, 10)}
	w.AddObserver(obs)

	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	testFile := filepath.Join(subDir, "test_sub.txt")
	if err := os.WriteFile(testFile, []byte("sub"), 0666); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-obs.ch:
		if event.Path != testFile {
			t.Fatalf("路径不匹配: 期望 %s, 实际 %s", testFile, event.Path)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("子目录事件等待超时")
	}
}

func TestWatchSimpleAPI(t *testing.T) {
	dir, err := os.MkdirTemp("", "dzwatch_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ch := make(chan WatchEvent, 10)
	w, err := Watch(dir, func(event WatchEvent) {
		ch <- event
	})
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	testFile := filepath.Join(dir, "simple.txt")
	if err := os.WriteFile(testFile, []byte("simple"), 0666); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-ch:
		if event.Type != WatchEventCreate {
			t.Fatalf("期望 Create 事件, 实际 %v", event.Type)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Watch API 事件等待超时")
	}
}

func TestFollow(t *testing.T) {
	dir, err := os.MkdirTemp("", "dzwatch_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "follow.txt")
	if err := os.WriteFile(testFile, []byte("init"), 0666); err != nil {
		t.Fatal(err)
	}

	w, err := NewWatcher(dir)
	if err != nil {
		t.Fatal(err)
	}
	w.SetDelay(100 * time.Millisecond)

	obs := &testObserver{ch: make(chan WatchEvent, 20)}
	w.Follow(testFile, obs)

	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	if err := os.WriteFile(testFile, []byte("modified"), 0666); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-obs.ch:
		if event.Type != WatchEventWrite && event.Type != WatchEventCreate {
			t.Fatalf("期望 Write/Create 事件, 实际 %v", event.Type)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Follow 事件等待超时")
	}
}

type testObserver struct {
	mu     sync.Mutex
	events []WatchEvent
	ch     chan WatchEvent
}

func (o *testObserver) OnEvent(event WatchEvent) {
	o.mu.Lock()
	o.events = append(o.events, event)
	o.mu.Unlock()
	if o.ch != nil {
		select {
		case o.ch <- event:
		default:
		}
	}
}
