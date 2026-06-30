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

func TestIsTermux(t *testing.T) {
	result := isTermux()
	if result {
		// isTermux 返回 true 时，标志文件必须存在
		if _, err := os.Stat("/data/data/com.termux/files/usr/etc/termux"); err != nil {
			t.Fatal("isTermux returns true but termux marker doesn't exist")
		}
	}
}

func TestDetectBashPathTermux(t *testing.T) {
	if !isTermux() {
		t.Skip("Not running on Termux, skipping")
	}

	path := detectBashPath()
	expected := termuxPrefix() + "/bin/bash"
	if path != expected {
		t.Fatalf("On Termux, expected %s, got: %s", expected, path)
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

func TestReadResolvConfDNS(t *testing.T) {
	f, err := os.CreateTemp("", "resolv.conf.*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	content := "# comment\nnameserver 1.2.3.4\nnameserver 5.6.7.8\n"
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	addr := readResolvConfDNS(f.Name())
	if addr != "1.2.3.4:53" {
		t.Fatalf("Expected 1.2.3.4:53, got: %s", addr)
	}
}

func TestReadResolvConfDNSEmpty(t *testing.T) {
	f, err := os.CreateTemp("", "resolv.conf.*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	// 文件为空时回退到默认 DNS
	addr := readResolvConfDNS(f.Name())
	if addr != "114.114.114.114:53" {
		t.Fatalf("Expected fallback 114.114.114.114:53, got: %s", addr)
	}
}

func TestTermuxPrefix(t *testing.T) {
	prefix := termuxPrefix()
	if prefix != "/data/data/com.termux/files/usr" {
		t.Fatalf("Expected fixed Termux prefix, got: %s", prefix)
	}
	if isTermux() {
		if _, err := os.Stat(prefix); err != nil {
			t.Fatalf("Termux prefix %s should exist when isTermux is true: %v", prefix, err)
		}
	}
}

func TestSetTermuxCACerts(t *testing.T) {
	if !isTermux() {
		t.Skip("Not running on Termux, skipping")
	}

	// 手动触发以确保 init 之后的状态一致
	setTermuxCACerts()
	certFile := os.Getenv("SSL_CERT_FILE")
	if certFile == "" {
		t.Fatal("Expected SSL_CERT_FILE to be set on Termux")
	}
	if _, err := os.Stat(certFile); err != nil {
		t.Fatalf("SSL_CERT_FILE %s does not exist: %v", certFile, err)
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
