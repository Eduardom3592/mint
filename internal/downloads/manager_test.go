package downloads

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager(t.TempDir(), 3, nil)
	defer m.Close()

	if m.maxWorkers != 3 {
		t.Errorf("expected maxWorkers 3, got %d", m.maxWorkers)
	}
}

func TestNewManagerDefaultsToAtLeastOneWorker(t *testing.T) {
	m := NewManager(t.TempDir(), 0, nil)
	defer m.Close()

	if m.maxWorkers != 1 {
		t.Errorf("expected maxWorkers 1, got %d", m.maxWorkers)
	}
}

func TestEnqueue(t *testing.T) {
	m := NewManager(t.TempDir(), 3, nil)
	defer m.Close()

	id := m.Enqueue("proj-1", "test project", "ver-1", "1.0.0", "https://example.com/file.jar", "file.jar", 1024, nil)
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}

	items := m.List()
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}

	if items[0].Status != StatusQueued {
		t.Errorf("expected queued status, got %s", items[0].Status)
	}
}

func TestEnqueueUsesPersistedID(t *testing.T) {
	p := &fakePersister{nextID: 41}
	m := NewManager(t.TempDir(), 3, p)
	defer m.Close()

	id := m.Enqueue("proj-1", "test project", "ver-1", "1.0.0", "https://example.com/file.jar", "file.jar", 1024, nil)
	if id != 42 {
		t.Errorf("expected persisted id 42, got %d", id)
	}

	items := m.List()
	if len(items) != 1 || items[0].ID != 42 {
		t.Fatalf("expected item id 42, got %#v", items)
	}
}

func TestEnqueueSanitizesFilename(t *testing.T) {
	m := NewManager(t.TempDir(), 3, nil)
	defer m.Close()

	m.Enqueue("proj-1", "test project", "ver-1", "1.0.0", "https://example.com/file.jar", "../file.jar", 1024, nil)
	items := m.List()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Filename != "file.jar" {
		t.Errorf("expected sanitized filename file.jar, got %q", items[0].Filename)
	}
}

func TestDownloadCompletesAfterEOF(t *testing.T) {
	content := []byte("download body")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(content)
	}))
	defer server.Close()

	dir := t.TempDir()
	m := NewManager(dir, 1, nil)
	defer m.Close()

	item := &Item{
		FileURL:   server.URL,
		Filename:  "file.jar",
		TotalSize: int64(len(content)),
	}

	if err := m.download(context.Background(), item); err != nil {
		t.Fatalf("download failed: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dir, "file.jar"))
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}
	if string(got) != string(content) {
		t.Fatalf("downloaded content mismatch: got %q want %q", got, content)
	}
}

func TestCancel(t *testing.T) {
	m := NewManager(t.TempDir(), 3, nil)
	defer m.Close()

	id := m.Enqueue("proj-1", "test project", "ver-1", "1.0.0", "https://example.com/file.jar", "file.jar", 1024, nil)

	if !m.Cancel(id) {
		t.Error("expected cancel to return true")
	}

	items := m.List()
	if items[0].Status != StatusCancelled {
		t.Errorf("expected cancelled status, got %s", items[0].Status)
	}
}

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		s    Status
		want string
	}{
		{StatusQueued, "queued"},
		{StatusPreparing, "preparing"},
		{StatusDownloading, "downloading"},
		{StatusVerifying, "verifying"},
		{StatusCompleted, "completed"},
		{StatusFailed, "failed"},
		{StatusCancelled, "cancelled"},
	}
	for _, tt := range tests {
		if string(tt.s) != tt.want {
			t.Errorf("expected %s, got %s", tt.want, string(tt.s))
		}
	}
}

type fakePersister struct {
	mu     sync.Mutex
	nextID int64
}

func (p *fakePersister) SaveDownload(_, _, _, _, _, _, _ string, _ int64, _ string, _ bool) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.nextID++
	return p.nextID, nil
}

func (p *fakePersister) UpdateDownloadStatus(int64, string) error {
	return nil
}

func (p *fakePersister) ListDownloads() ([]DownloadRecord, error) {
	return nil, nil
}

func (p *fakePersister) DeleteDownload(int64) error {
	return nil
}

func (p *fakePersister) DeleteAllDownloads() error {
	return nil
}

func (p *fakePersister) MarkInstalled(int64) error {
	return nil
}
