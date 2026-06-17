# Agent 自动注册

gocron 支持通过 Web 界面一键生成安装命令，在目标服务器上执行即可自动安装并注册 Agent 节点。这是部署任务节点最简单快捷的方式。

## 功能特性

* ✅ **一键安装** - 复制命令，一行搞定
* ✅ **自动注册** - 无需手动配置，自动连接到 gocron 服务器
* ✅ **安全Token** - 一次性 Token，3小时有效期
* ✅ **批量部署** - Token 可在有效期内重复使用
* ✅ **系统服务** - 自动创建系统服务，开机自启
* ✅ **离线支持** - 支持内网环境，优先使用本地安装包
* ✅ **跨平台** - 支持 Linux、macOS、Windows

## 使用方法

::: tip 💡 提示：离线安装（可选）
如果你的环境无法访问 GitHub，或者希望加快安装速度，可以先准备离线安装包：

1. 在 gocron 服务器的所在目录上创建目录：`mkdir -p gocron-node-package`
2. 从 [GitHub Releases](https://github.com/gocronx-team/gocron/releases) 下载所需平台的 `gocron-node-*.tar.gz` 或 `gocron-node-*.zip` 文件
3. 将下载的文件放入 `gocron-node-package/` 目录

这样安装脚本会优先从本地服务器下载，无需访问 GitHub。详见下方[离线安装](#离线安装)章节。
:::

### 步骤 1：生成安装命令

1. 登录 gocron Web 界面
2. 进入「任务节点」页面
3. 点击「自动注册」按钮
4. 选择目标平台（Linux/macOS/Windows）
5. 复制生成的安装命令

### 步骤 2：在目标服务器执行

根据目标服务器的操作系统，选择对应的安装方式：

#### Linux / macOS

```bash
curl -fsSL http://your-server:5920/api/agent/install.sh | bash -s -- <token>
```

**示例**：
```bash
curl -fsSL http://192.168.1.100:5920/api/agent/install.sh | bash -s -- eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Windows

::: warning 注意
出于安全考虑，Windows 暂不支持自动安装脚本。请使用手动方式安装。
:::

**手动安装步骤**：

1. 下载 Windows 安装包：[gocron-node-windows-amd64.zip](https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-windows-amd64.zip)
2. 解压到指定目录（例如 `D:\gocron-node`）
3. 打开 PowerShell
4. 进入解压目录
5. 运行 `.\gocron-node.exe` 启动节点
6. 在 Web 界面手动添加该节点（IP 和端口）

### 步骤 3：验证安装

安装完成后，返回 Web 界面的「任务节点」页面，应该能看到新注册的节点。

## 安装过程

对于 Linux/macOS，安装脚本会自动完成以下操作：

1. **检测系统** - 自动识别操作系统和架构
2. **下载安装包** - 从服务器或 GitHub 下载对应版本
3. **解压安装** - 将文件解压到指定目录
4. **注册节点** - 使用 Token 向服务器注册
5. **创建服务** - 创建系统服务（systemd / launchd）
6. **启动服务** - 自动启动 gocron-node 服务

### 安装目录

| 操作系统 | 安装目录 |
|---------|---------|
| Linux | `/opt/gocron-node` |
| macOS | `/usr/local/gocron-node` |

### 配置文件

| 操作系统 | 配置文件路径 |
|---------|-------------|
| Linux/macOS | `~/.gocron-node/node.ini` |

## Token 机制

### Token 特性

- **一次性生成** - 每次点击「自动注册」生成新 Token
- **有效期** - 3小时内有效
- **可重复使用** - 同一 Token 可用于批量安装多个节点
- **自动失效** - 超过有效期后自动失效

### Token 安全性

- Token 采用 JWT 格式，包含加密签名
- Token 仅用于节点注册，不包含敏感信息
- 注册完成后，Token 即可丢弃
- 节点与服务器之间采用双向通信：
  - **注册时**：节点主动通过 HTTP 连接服务器
  - **执行任务时**：服务器主动通过 gRPC 连接节点（需确保服务器能访问节点的 IP:Port）

## 离线安装

对于内网环境或无法访问 GitHub 的场景，gocron 支持离线安装。

### 配置离线安装包

#### 1. 创建安装包目录

在 gocron 服务器的所在目录上创建目录：

```bash
mkdir -p gocron-node-package
```

#### 2. 下载安装包

访问 [GitHub Releases](https://github.com/gocronx-team/gocron/releases) 下载所需平台的安装包。

**方式一：手动下载**

下载以下文件并放入 `gocron-node-package/` 目录：
- `gocron-node-linux-amd64.tar.gz`
- `gocron-node-linux-arm64.tar.gz`
- `gocron-node-darwin-amd64.tar.gz`
- `gocron-node-darwin-arm64.tar.gz`
- `gocron-node-windows-amd64.zip`

**方式二：使用 curl 下载**

```bash
# Linux amd64
curl -L -o gocron-node-package/gocron-node-linux-amd64.tar.gz \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-linux-amd64.tar.gz

# Linux arm64
curl -L -o gocron-node-package/gocron-node-linux-arm64.tar.gz \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-linux-arm64.tar.gz

# macOS amd64
curl -L -o gocron-node-package/gocron-node-darwin-amd64.tar.gz \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-darwin-amd64.tar.gz

# macOS arm64 (Apple Silicon)
curl -L -o gocron-node-package/gocron-node-darwin-arm64.tar.gz \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-darwin-arm64.tar.gz

# Windows amd64
curl -L -o gocron-node-package/gocron-node-windows-amd64.zip \
  https://github.com/gocronx-team/gocron/releases/latest/download/gocron-node-windows-amd64.zip
```

#### 3. 目录结构

```
gocron/
├── gocron                          # gocron 主程序
├── gocron-node-package/            # 离线安装包目录
│   ├── gocron-node-linux-amd64.tar.gz
│   ├── gocron-node-linux-arm64.tar.gz
│   ├── gocron-node-darwin-amd64.tar.gz
│   ├── gocron-node-darwin-arm64.tar.gz
│   └── gocron-node-windows-amd64.zip
```

### 工作原理

安装脚本会自动检测本地是否有安装包：

- ✅ **本地有安装包** - 直接从内网服务器下载，速度快，无需外网
- ⚠️ **本地无安装包** - 自动跳转到 GitHub 下载（需要外网访问）

用户在执行安装命令时会看到清晰的提示：
```
正在从本地服务器下载安装包...
或
正在从 GitHub 下载安装包...
```

::: tip 提示
- 文件名必须严格按照格式：`gocron-node-{os}-{arch}.{ext}`
- 此配置为可选项，不配置时会自动从 GitHub 下载
- 建议定期更新到最新版本
:::

## 系统服务管理

### Linux (systemd)

安装脚本会自动创建 systemd 服务：

```bash
# 查看服务状态
sudo systemctl status gocron-node

# 启动服务
sudo systemctl start gocron-node

# 停止服务
sudo systemctl stop gocron-node

# 重启服务
sudo systemctl restart gocron-node

# 查看日志
sudo journalctl -u gocron-node -f
```

### macOS (后台进程)

macOS 安装脚本使用 `nohup` 在后台运行 gocron-node，而不是创建 launchd 服务：

```bash
# 查看进程状态
ps aux | grep gocron-node | grep -v grep

# 停止进程
pkill -f gocron-node

# 启动进程
nohup /usr/local/gocron-node/gocron-node > /tmp/gocron-node.log 2>&1 &

# 查看日志
tail -f /tmp/gocron-node.log
```

::: tip 提示
如果需要开机自启，可以手动创建 launchd 配置文件。参考 [macOS launchd 文档](https://www.launchd.info/)。
:::

### Windows (Windows Service)

如果是手动安装，默认不会创建 Windows 服务。你可以直接运行 `gocron-node.exe`，或者使用 `sc` 命令手动创建服务：

```powershell
# 创建服务
sc.exe create gocron-node binPath= "D:\gocron-node\gocron-node.exe" start= auto

# 启动服务
Start-Service gocron-node
```

如果已创建服务，可以使用以下命令管理：

```powershell
# 查看服务状态
Get-Service gocron-node

# 停止服务
Stop-Service gocron-node

# 重启服务
Restart-Service gocron-node
```

## 卸载节点

### Linux / macOS

```bash
# 停止服务/进程
sudo systemctl stop gocron-node  # Linux
pkill -f gocron-node  # macOS

# 删除服务文件
sudo rm /etc/systemd/system/gocron-node.service  # Linux

# 删除安装目录
sudo rm -rf /opt/gocron-node  # Linux
sudo rm -rf /usr/local/gocron-node  # macOS

# 删除配置文件
rm -rf ~/.gocron-node

# 删除日志文件（macOS）
rm /tmp/gocron-node.log  # macOS
```

### Windows

```powershell
# 停止并删除服务
Stop-Service gocron-node
sc.exe delete gocron-node

# 删除安装目录
Remove-Item -Recurse -Force "C:\Program Files\gocron-node"

# 删除配置文件
Remove-Item -Recurse -Force "$env:USERPROFILE\.gocron-node"
```

## 故障排除

### 安装失败

**问题**：安装脚本执行失败

**解决方案**：
1. 检查网络连接
2. 确认 Token 未过期（3小时有效期）
3. 检查是否有管理员权限
4. 查看安装日志获取详细错误信息

### 节点未显示

**问题**：安装完成但节点列表中看不到

**解决方案**：
1. 检查节点服务是否正常运行
2. 检查防火墙设置
3. 确认服务器地址配置正确
4. 查看节点日志：
   - Linux/macOS: `sudo journalctl -u gocron-node -f`
   - Windows: 事件查看器

### Token 过期

**问题**：Token 已过期无法使用

**解决方案**：
重新生成 Token。在 Web 界面点击「自动注册」生成新的 Token。

### 权限不足

**问题**：Linux/macOS 提示权限不足

**解决方案**：
使用具有 `sudo` 权限的普通用户执行安装命令。**注意：不要直接使用 root 用户运行安装脚本**，脚本会自动使用 `sudo` 提权执行需要管理员权限的操作。

**问题**：Windows 提示权限不足

**解决方案**：
右键点击 PowerShell，选择「以管理员身份运行」。

## 相关文档

- [快速开始](./quick-start) - 了解其他部署方式
- [配置文件](./configuration) - 节点配置说明
- [定时任务](./scheduled-tasks) - 创建和管理任务
