# 安全 - TLS 双向认证

TLS 双向认证为 gocron 提供了更高级别的安全保护，确保客户端和服务器之间的通信安全。

## 什么是 TLS 双向认证

TLS 双向认证（Mutual TLS Authentication）是一种安全机制，要求客户端和服务器都提供证书进行身份验证：
- **服务器认证**：客户端验证服务器的身份
- **客户端认证**：服务器验证客户端的身份

## 配置 TLS

### 1. 生成证书

首先需要生成 CA 证书、服务器证书和客户端证书。

**生成 CA 证书**：
```bash
# 生成 CA 私钥
openssl genrsa -out ca.key 2048

# 生成 CA 证书
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt
```

**生成服务器证书**：
```bash
# 生成服务器私钥
openssl genrsa -out server.key 2048

# 生成服务器证书签名请求
openssl req -new -key server.key -out server.csr

# 使用 CA 签名服务器证书
openssl x509 -req -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
```

**生成客户端证书**：
```bash
# 生成客户端私钥
openssl genrsa -out client.key 2048

# 生成客户端证书签名请求
openssl req -new -key client.key -out client.csr

# 使用 CA 签名客户端证书
openssl x509 -req -days 3650 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt
```

### 2. 配置 gocron

在配置文件 `.gocron/conf/app.ini` 中添加以下配置：

```ini
[tls]
enable_tls = true
ca_file = /path/to/ca.crt
cert_file = /path/to/server.crt
key_file = /path/to/server.key
```

### 3. 重启服务

修改配置后，重启 gocron 服务使配置生效。

## 客户端配置

使用 TLS 双向认证时，客户端也需要配置证书：

```bash
curl --cacert ca.crt --cert client.crt --key client.key https://gocron-server:5920
```

## 验证配置

可以使用以下命令验证 TLS 配置是否正确：

```bash
openssl s_client -connect localhost:5920 -CAfile ca.crt -cert client.crt -key client.key
```

## 故障排除

### 常见问题

**Q: 启用 TLS 后无法访问**
- 检查证书路径是否正确
- 确认证书文件权限
- 查看 gocron 日志了解详细错误信息

**Q: 证书验证失败**
- 确认 CA 证书、服务器证书和客户端证书是否由同一个 CA 签发
- 检查证书是否过期
- 验证证书的 Common Name (CN) 是否正确

## 最佳实践

- **定期更新证书**：证书应定期更新，避免过期
- **安全存储私钥**：私钥文件应设置适当的权限，避免泄露
- **使用强加密**：使用 2048 位或更高的 RSA 密钥
- **监控证书有效期**：设置提醒，在证书过期前及时更新

## 相关文档

- [安全 - 双因素认证 (2FA)](./security-2fa)
- [配置文件](./configuration)
