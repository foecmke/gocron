//go:build !windows
// +build !windows

package utils

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Result struct {
	output string
	err    error
}

// isTermux 通过检查 Termux 标志文件判断当前是否运行在 Termux (Android) 环境下。
// 该文件是 Termux 安装后必定存在的，比检查 PREFIX 环境变量更可靠。
func isTermux() bool {
	_, err := os.Stat("/data/data/com.termux/files/usr/etc/termux")
	return err == nil
}

// termuxPrefix 返回 Termux 的安装前缀路径。
func termuxPrefix() string {
	return "/data/data/com.termux/files/usr"
}

// init 在包加载时处理 Termux/Android 的 DNS 和 TLS CA 证书适配。
func init() {
	if !isTermux() {
		return
	}

	// DNS：Termux 没有 /etc/resolv.conf，通过自定义 Resolver 接管 DNS 解析。
	// 如果用户安装了 resolv.conf 则读取其中的 nameserver，否则使用 114.114.114.114。
	resolvPath := termuxPrefix() + "/etc/resolv.conf"
	dnsAddr := readResolvConfDNS(resolvPath)
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, dnsAddr)
		},
	}

	// TLS CA：Termux 下没有标准 CA 证书路径，设置 SSL_CERT_FILE 指向 Termux 证书
	setTermuxCACerts()
}

// readResolvConfDNS 解析 resolv.conf 文件，提取第一个 nameserver 并返回 ip:53。
// 文件不存在或解析失败则回退到 114.114.114.114:53。
func readResolvConfDNS(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "114.114.114.114:53"
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if after, ok := strings.CutPrefix(line, "nameserver"); ok {
			addr := strings.TrimSpace(after)
			if addr != "" {
				return net.JoinHostPort(addr, "53")
			}
		}
	}
	return "114.114.114.114:53"
}

// setTermuxCACerts 查找 Termux 的 CA 证书文件并设置 SSL_CERT_FILE 环境变量。
// Go 的 crypto/x509 包读取该变量以加载 CA 证书池。
func setTermuxCACerts() {
	prefix := termuxPrefix()
	paths := []string{
		prefix + "/etc/tls/cert.pem",
		prefix + "/etc/tls/ca-bundle.pem",
		prefix + "/etc/ssl/certs/ca-certificates.crt",
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			_ = os.Setenv("SSL_CERT_FILE", p)
			return
		}
	}
}

// detectBashPath 根据当前系统环境检测 bash 路径。
// Termux (Android) 下 bash 位于 termuxPrefix/bin/bash，标准 Linux / macOS 下为 /bin/bash。
func detectBashPath() string {
	if isTermux() {
		return termuxPrefix() + "/bin/bash"
	}
	return "/bin/bash"
}

// 执行shell命令，可设置执行超时时间
// 改进：将命令写入临时脚本执行，即使超时或被取消，也会返回已产生的输出
func ExecShell(ctx context.Context, command string) (string, error) {
	// 清理可能存在的 HTML 实体编码
	command = CleanHTMLEntities(command)
	// 将换行符统一替换为Unix风格的\n
	command = strings.ReplaceAll(command, "\r\n", "\n")

	// 使用系统临时目录
	tmpDir := os.TempDir()
	timestamp := time.Now().Format("20060102150405")
	scriptPattern := fmt.Sprintf("gocron_%s_*.sh", timestamp)

	tmpFile, err := os.CreateTemp(tmpDir, scriptPattern)
	if err != nil {
		return "", fmt.Errorf("创建临时脚本文件失败: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // 执行完毕后删除临时文件
	defer tmpFile.Close()

	// 将命令写入临时文件
	_, err = tmpFile.WriteString(command)
	if err != nil {
		return "", fmt.Errorf("写入脚本内容失败: %w", err)
	}

	// 确保文件写入磁盘
	err = tmpFile.Sync()
	if err != nil {
		return "", fmt.Errorf("同步文件失败: %w", err)
	}

	// 给脚本文件添加执行权限
	err = os.Chmod(tmpFile.Name(), 0700)
	if err != nil {
		return "", fmt.Errorf("设置脚本执行权限失败: %w", err)
	}

	// 根据当前系统环境检测 bash 路径，执行脚本文件
	scriptPath := tmpFile.Name()
	bashPath := detectBashPath()
	cmd := exec.Command(bashPath, scriptPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	// 设置工作目录为用户家目录，避免 getcwd 错误
	if homeDir, err := os.UserHomeDir(); err == nil {
		cmd.Dir = homeDir
	} else {
		cmd.Dir = tmpDir
	}

	// 使用管道实时捕获输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	// 用于收集输出
	var outputBuffer bytes.Buffer
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 启动命令
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// 实时读取 stdout 和 stderr
	wg.Add(2)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				mu.Lock()
				outputBuffer.Write(buf[:n])
				mu.Unlock()
			}
			if err != nil {
				break
			}
		}
	}()
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				mu.Lock()
				outputBuffer.Write(buf[:n])
				mu.Unlock()
			}
			if err != nil {
				break
			}
		}
	}()

	// 等待命令完成或超时
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// 超时或被取消，尝试优雅终止
		if cmd.Process != nil && cmd.Process.Pid > 0 {
			// 先发送 SIGTERM，给进程清理的机会
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)

			// 等待 2 秒，看进程是否自行退出
			timer := time.NewTimer(2 * time.Second)
			select {
			case <-done:
				timer.Stop()
			case <-timer.C:
				// 进程仍未退出，强制杀死
				_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
				<-done // 等待 Wait() 返回
			}
		}

		// 等待 IO 读取完成
		wg.Wait()

		// 返回已捕获的输出和错误信息
		mu.Lock()
		output := outputBuffer.String()
		mu.Unlock()
		return output, errors.New("timeout killed")

	case err := <-done:
		// 命令正常完成
		wg.Wait()
		mu.Lock()
		output := outputBuffer.String()
		mu.Unlock()
		return output, err
	}
}
