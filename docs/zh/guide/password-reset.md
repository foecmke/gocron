# 密码重置 (CLI)

当管理员忘记密码，或因丢失双因素认证（2FA）设备被锁在门外时，可以使用内置的
`resetpwd` 命令，直接在命令行重置密码。

## 概述

`resetpwd` 命令会直接在数据库中更新用户密码。它只初始化配置和数据库连接，**不会**
启动 Web 服务、调度器或 leader 选举，因此在 gocron 运行期间执行也是安全的。

::: tip
`resetpwd` 的配置文件（`conf/app.ini`）是相对 gocron **二进制**定位的，而非当前
目录；并且会自动把相对的 SQLite 数据库路径锚定到 gocron 的基准目录。因此你可以在
任意工作目录下运行它——只要用**已安装的同一个二进制**（服务端运行的那个）即可。
对于 SQLite，它会打印所连接的数据库文件路径，确认前请核对该行与服务端数据库一致。
:::

## 用法

```bash
gocron resetpwd [username]
```

| 参数 / 选项 | 说明 |
| --- | --- |
| `[username]` | 要重置的用户。可选，默认为 `admin`。按用户名匹配。 |
| `--disable-2fa` | 同时关闭该用户的双因素认证。 |

## 交互流程

1. 显示二次确认提示（`y/N`），输入非 `y`/`yes` 即取消。
2. 提示输入新密码。输入**不回显**到终端。
3. 若手动输入密码，需要再输入一次以确认；两次不一致则中止操作，不做任何改动。
4. 若**留空**直接回车，则自动生成一个 12 位随机密码并打印到屏幕上。

## 示例

重置默认的 `admin` 账户：

```bash
./gocron resetpwd
```

重置指定用户：

```bash
./gocron resetpwd alice
```

重置密码**并**关闭 2FA（当管理员因丢失验证器被锁定时使用）：

```bash
./gocron resetpwd admin --disable-2fa
```

### 在 Docker 中

如果 gocron 运行在容器里，进入容器对容器内的二进制执行命令：

```bash
docker exec -it gocron ./gocron resetpwd admin
```

同时清除 2FA：

```bash
docker exec -it gocron ./gocron resetpwd admin --disable-2fa
```

## 故障排查

- **提示 `gocron is not installed; cannot reset password`** —— 二进制旁边找不到
  完成安装的标志（`conf/app.ini` 和安装锁）。请使用已安装的 gocron 二进制——即服务端
  运行的那个（Docker 环境下进入正在运行的容器执行）。
- **提示 `user [xxx] not found`** —— 该用户名不存在，或你连接的数据库与预期不符。
  请核对用户名以及 `conf/app.ini` 中配置的数据库。

::: warning
出于安全考虑，建议让工具自动生成随机密码（输入时留空）。重新登录后，请通过 Web 界面
修改密码；如果之前关闭了 2FA，也请重新启用。
:::

## 相关文档

- [安全 - 双因素认证 (2FA)](./security-2fa)
- [配置文件](./configuration)
