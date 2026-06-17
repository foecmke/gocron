# 配置文件

::: tip 配置文件路径
配置文件路径：`.gocron/conf/app.ini`

系统安装后会自动创建该配置文件。
:::

## 配置项说明

### 基础配置

- **allow_ips**  
  允许访问的客户端 IP，多个 IP 用逗号分隔，默认为空（不限制）

- **app.name**  
  应用名称

### 数据库配置

- **db.engine**  
  数据库引擎，支持 `mysql` 和 `postgres`

- **db.host**  
  数据库主机名

- **db.port**  
  数据库端口

- **db.user**  
  数据库用户名

- **db.password**  
  数据库密码

- **db.charset**  
  数据库字符集

- **db.prefix**  
  表前缀

- **db.database**  
  数据库名

### API 配置

- **api.key**  
  API 接口 key，未配置则不能使用接口

- **api.secret**  
  API 接口秘钥，未配置则不能使用接口

- **api.sign.enable**  
  是否启用签名验证

### TLS 配置

- **enable_tls**  
  开启 TLS

- **ca_file**  
  CA 证书文件路径

- **cert_file**  
  客户端证书路径

- **key_file**  
  客户端私钥路径

## 配置示例

```ini
[app]
name = gocron

[server]
allow_ips = 

[db]
engine = mysql
host = 127.0.0.1
port = 3306
user = root
password = 
database = gocron
charset = utf8mb4
prefix = 

[api]
key = your_api_key
secret = your_api_secret
sign.enable = true

[tls]
enable_tls = false
ca_file = 
cert_file = 
key_file = 
```

## 相关文档

- [安全 - TLS 双向认证](./security-tls)
- [API 文档](./api)
