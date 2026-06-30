//go:build !windows
// +build !windows

package utils

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDetectBashPath(t *testing.T) {
	path := detectBashPath()
	if path == "" {
		t.Fatal("Expected non-empty bash path")
	}

	// 验证返回的路径是 Termux 或标准 Linux/macOS 的合法 bash 路径
	validPaths := []string{
		"/data/data/com.termux/files/usr/bin/bash",
		"/bin/bash",
	}
	found := false
	for _, p := range validPaths {
		if path == p {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Expected bash path to be one of %v, got: %s", validPaths, path)
	}

	// 验证返回的 bash 路径文件存在
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Bash path %s does not exist: %v", path, err)
	}
}

func TestDetectBashPathTermux(t *testing.T) {
	// 若 Termux bash 文件存在，则 detectBashPath 应返回 Termux 路径
	termuxBash := "/data/data/com.termux/files/usr/bin/bash"
	if _, err := os.Stat(termuxBash); err != nil {
		t.Skip("Not running on Termux, skipping")
	}

	path := detectBashPath()
	if path != termuxBash {
		t.Fatalf("On Termux, expected %s, got: %s", termuxBash, path)
	}
}

func TestExecShellUsesTempDir(t *testing.T) {
	ctx := context.Background()
	output, err := ExecShell(ctx, "echo 'hello world'")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !strings.Contains(output, "hello world") {
		t.Fatalf("Expected output to contain 'hello world', got: %s", output)
	}

	// 检查临时目录中没有遗留的脚本文件
	tempDir := os.TempDir()
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "gocron_") && strings.HasSuffix(entry.Name(), ".sh") {
			t.Fatalf("Temporary script file was not cleaned up: %s", filepath.Join(tempDir, entry.Name()))
		}
	}
}

func TestExecShellSuccess(t *testing.T) {
	ctx := context.Background()
	output, err := ExecShell(ctx, "echo 'hello world'")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !strings.Contains(output, "hello world") {
		t.Fatalf("Expected output to contain 'hello world', got: %s", output)
	}
}

func TestExecShellTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// 运行一个会产生输出然后睡眠的命令
	output, err := ExecShell(ctx, "echo 'partial output'; sleep 1; echo 'should not see this'")

	if err == nil {
		t.Fatal("Expected timeout error")
	}
	if err.Error() != "timeout killed" {
		t.Fatalf("Expected 'timeout killed' error, got: %v", err)
	}
	if !strings.Contains(output, "partial output") {
		t.Fatalf("Expected partial output to contain 'partial output', got: %s", output)
	}
	if strings.Contains(output, "should not see this") {
		t.Fatalf("Should not contain output after timeout, got: %s", output)
	}
}

func TestExecShellCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// 启动一个长时间运行的命令
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel() // 手动取消
	}()

	output, err := ExecShell(ctx, "echo 'before cancel'; sleep 1; echo 'after cancel'")

	if err == nil {
		t.Fatal("Expected cancel error")
	}
	if err.Error() != "timeout killed" {
		t.Fatalf("Expected 'timeout killed' error, got: %v", err)
	}
	if !strings.Contains(output, "before cancel") {
		t.Fatalf("Expected partial output to contain 'before cancel', got: %s", output)
	}
}

func TestExecShellCommandError(t *testing.T) {
	ctx := context.Background()
	output, err := ExecShell(ctx, "nonexistentcommand")

	if err == nil {
		t.Fatal("Expected command error")
	}
	// 应该有错误输出
	if output == "" {
		t.Fatal("Expected some error output")
	}
}
