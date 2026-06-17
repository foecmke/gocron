# 快速开始

本指南将帮助你快速部署和运行 gocron。

## 环境要求

- **Go**: 1.23+
- **数据库**: MySQL / PostgreSQL / SQLite（见下方说明）
- **Node.js**: 20+（仅前端开发需要）

## 数据库支持说明

| 部署方式 | MySQL | PostgreSQL | SQLite |
|---------|-------|------------|--------|
| Docker 部署 | ✅ 支持 | ✅ 支持 | ✅ 支持 |
| Kubernetes 部署 | ✅ 支持 | ✅ 支持 | ✅ 支持 |
| 二进制部署 | ✅ 支持 | ✅ 支持 | ✅ 支持 |
| 开发环境 | ✅ 支持 | ✅ 支持 | ✅ 支持 |

::: tip 提示
- 所有部署方式均支持三种数据库，gocron 使用纯 Go SQLite 驱动，无需 CGO
- **生产环境推荐**：使用 MySQL 或 PostgreSQL，性能更好，支持分布式部署
:::

## Docker Compose 部署（推荐）

使用 Docker 部署最简单快捷，适合快速体验和测试环境。

### 步骤

```bash
# 1. 克隆项目
git clone https://github.com/gocronx-team/gocron.git
cd gocron

# 2. 启动服务（自动构建镜像）
docker-compose up -d

# 3. 访问 Web 界面
# http://localhost:5920
```

### 默认账号

- 用户名：`admin`
- 密码：`admin123`

::: tip 提示
- Docker Compose 仅部署 gocron 管理端
- 任务节点（gocron-node）需要单独安装部署
- 请参考 [Agent 自动注册](./agent-registration) 章节安装任务节点
:::

## Kubernetes 部署（Helm）

使用 Helm Chart 一键部署到 Kubernetes 集群。

### 添加 Helm 仓库

```bash
helm repo add gocron https://gocronx-team.github.io/gocron
helm repo update
```

### 部署

```bash
# 使用 SQLite（默认，最简单）
helm install gocron gocron/gocron

# 使用 MySQL
helm install gocron gocron/gocron \
  --set db.engine=mysql \
  --set db.host=mysql.default \
  --set db.port=3306 \
  --set db.user=gocron \
  --set db.password=your_password \
  --set db.database=gocron

# 使用 PostgreSQL
helm install gocron gocron/gocron \
  --set db.engine=postgres \
  --set db.host=pg.default \
  --set db.port=5432 \
  --set db.user=gocron \
  --set db.password=your_password \
  --set db.database=gocron
```

### 配置 Ingress

```bash
helm install gocron gocron/gocron \
  --set ingress.enabled=true \
  --set 'ingress.hosts[0].host=gocron.example.com' \
  --set 'ingress.hosts[0].paths[0].path=/' \
  --set 'ingress.hosts[0].paths[0].pathType=Prefix'
```

::: tip 提示
完整的 Helm 配置项请参考 [Kubernetes 部署](./kubernetes) 章节。
:::

## 二进制部署（生产环境推荐）

适合生产环境，支持所有数据库（包括 SQLite）。

### 下载安装包

访问 [GitHub Releases](https://github.com/gocronx-team/gocron/releases) 下载最新版本的安装包。

选择对应平台的包：
- Linux: `gocron-linux-amd64.tar.gz` 或 `gocron-linux-arm64.tar.gz`
- macOS: `gocron-darwin-amd64.tar.gz` 或 `gocron-darwin-arm64.tar.gz`
- Windows: `gocron-windows-amd64.zip` 或 `gocron-windows-arm64.zip`

### 安装步骤

```bash
# 1. 解压安装包
tar -xzf gocron-linux-amd64.tar.gz
cd gocron-linux-amd64

# 2. 配置数据库（可选）
# 编辑 .gocron/conf/app.ini
# 默认使用 SQLite，无需配置

# 3. 启动服务
./gocron web

# 4. 访问 Web 界面
# http://localhost:5920
```

### 配置数据库

gocron 支持三种数据库，根据需要选择：

#### SQLite（默认）

无需配置，开箱即用。适合小型部署和测试环境。

#### MySQL

编辑 `.gocron/conf/app.ini`：

```ini
[db]
engine = mysql
host = 127.0.0.1
port = 3306
user = root
password = your_password
database = gocron
charset = utf8mb4
```

#### PostgreSQL

编辑 `.gocron/conf/app.ini`：

```ini
[db]
engine = postgres
host = 127.0.0.1
port = 5432
user = postgres
password = your_password
database = gocron
```

## 开发环境

适合开发和调试。

### 前置要求

- Go 1.23+
- Node.js 20+
- Yarn

### 步骤

```bash
# 1. 克隆项目
git clone https://github.com/gocronx-team/gocron.git
cd gocron

# 2. 安装 Go 依赖
go mod download

# 3. 配置数据库
# 编辑 .gocron/conf/app.ini
# 或复制 app.ini.sqlite.example 使用 SQLite

# 4. 安装开发工具
go install github.com/air-verse/air@latest

# 5. 启动后端（热更新）
air

# 6. 启动前端（另一个终端）
cd web/vue
yarn install
yarn run dev
```

访问 http://localhost:8080

## 安装任务节点

gocron 采用分布式架构，任务在独立的节点上执行。

### 方式一：自动注册（推荐）

使用 Web 界面一键生成安装命令，详见 [Agent 自动注册](./agent-registration)。

### 方式二：手动安装

```bash
# 1. 下载 gocron-node 安装包
# 访问 GitHub Releases 下载对应平台的 gocron-node 包

# 2. 解压
tar -xzf gocron-node-linux-amd64.tar.gz
cd gocron-node-linux-amd64

# 3. 启动节点
./gocron-node

# 4. 在 Web 界面添加节点
# 进入「任务节点」页面，点击「添加节点」
```

## 验证安装

1. 访问 Web 界面：http://localhost:5920
2. 使用默认账号登录（admin / admin123）
3. 进入「任务节点」页面，确认节点已连接
4. 创建一个测试任务，验证执行是否正常

## 下一步

- [配置文件](./configuration) - 了解详细的配置选项
- [Kubernetes 部署](./kubernetes) - Helm Chart 详细配置
- [定时任务](./scheduled-tasks) - 学习如何创建和管理任务
- [Agent 自动注册](./agent-registration) - 快速部署任务节点
- [API 文档](./api) - 使用 API 管理任务

## 常见问题

### Docker 部署可以使用 SQLite 吗？

可以。gocron v1.5.9+ 使用纯 Go SQLite 驱动，Docker 镜像已支持 SQLite。

### 如何修改默认端口？

```bash
# 方式一：命令行参数
./gocron web -p 8080

# 方式二：配置文件
# 编辑 .gocron/conf/app.ini
[server]
port = 8080
```

### 如何启用 HTTPS？

参考 [安全 - TLS 双向认证](./security-tls) 章节。

### 任务节点无法连接？

1. 检查节点是否正常启动
2. 检查防火墙设置
3. 确认节点地址配置正确
4. 查看节点日志：`log/cron.log`
