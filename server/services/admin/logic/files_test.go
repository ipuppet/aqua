package logic_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ipuppet/aqua/server/pkg/utils/crypto"
	"github.com/ipuppet/aqua/server/services/admin/logic"
	"github.com/ipuppet/aqua/server/services/admin/types"
)

type stubFileConfigProvider struct {
	fileCfg types.FileStoreConfig
}

func (s *stubFileConfigProvider) FileStoreConfig() (types.FileStoreConfig, error) {
	return s.fileCfg, nil
}

func TestFileService_UploadAndDelete(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()
	cfg := &stubFileConfigProvider{
		fileCfg: types.FileStoreConfig{URLPrefix: "/static", SavePath: tempDir},
	}
	svc := logic.NewFileService(cfg, nil)

	module := "blog"
	original := "image.png"
	data := strings.NewReader("hello")
	info, err := svc.UploadImage(ctx, module, original, data)
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	imageURL := info.URL

	today := time.Now().Format("2006-01-02")
	expectedName := crypto.MD5(original) + ".png"
	expectedURL := "/static/" + module + "/" + today + "/" + expectedName
	if imageURL != expectedURL {
		t.Fatalf("unexpected url path %s", imageURL)
	}

	diskPath := filepath.Join(tempDir, module, today, expectedName)
	b, err := os.ReadFile(diskPath)
	if err != nil {
		t.Fatalf("read uploaded file: %v", err)
	}
	if string(b) != "hello" {
		t.Fatalf("file content mismatch: %q", b)
	}

	delPath := info.Path
	if err := svc.DeleteImage(ctx, delPath); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := os.Stat(diskPath); !os.IsNotExist(err) {
		t.Fatalf("file still exists after delete")
	}
}
