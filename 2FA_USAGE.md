# 2FA (双因素认证) 使用说明

## 功能概述
系统已集成基于TOTP的双因素认证功能，提升账户安全性。

## 数据库迁移
首次使用需要执行数据库迁移：

```bash
# PostgreSQL
psql -U gocron -d gocron -f ADD_2FA_FIELDS.sql

# MySQL
mysql -u root -p gocron < ADD_2FA_FIELDS.sql
```

或手动执行：
```sql
ALTER TABLE user ADD COLUMN two_factor_key VARCHAR(100) DEFAULT '' COMMENT '2FA密钥';
ALTER TABLE user ADD COLUMN two_factor_on TINYINT NOT NULL DEFAULT 0 COMMENT '2FA开关 1:开启 0:关闭';
```

## API接口

### 1. 获取2FA状态
```
GET /api/user/2fa/status
```

响应：
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "enabled": false
  }
}
```

### 2. 设置2FA（获取密钥和二维码）
```
GET /api/user/2fa/setup
```

响应：
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_code": "data:image/png;base64,iVBORw0KG..."
  }
}
```

### 3. 启用2FA
```
POST /api/user/2fa/enable
Content-Type: application/x-www-form-urlencoded

secret=JBSWY3DPEHPK3PXP&code=123456
```

或JSON格式：
```json
{
  "secret": "JBSWY3DPEHPK3PXP",
  "code": "123456"
}
```

### 4. 禁用2FA
```
POST /api/user/2fa/disable
Content-Type: application/x-www-form-urlencoded

code=123456
```

### 5. 登录（启用2FA后）
```
POST /api/user/login
Content-Type: application/x-www-form-urlencoded

username=admin&password=123456&two_factor_code=123456
```

如果用户启用了2FA但未提供验证码，会返回：
```json
{
  "code": 1,
  "message": "需要输入2FA验证码",
  "data": {
    "require_2fa": true
  }
}
```

## 使用流程

### 用户启用2FA
1. 用户登录系统
2. 访问 `/api/user/2fa/setup` 获取密钥和二维码
3. 使用Google Authenticator、Microsoft Authenticator等APP扫描二维码
4. 输入APP显示的6位验证码，调用 `/api/user/2fa/enable` 启用2FA
5. 启用成功后，下次登录需要输入验证码

### 用户登录（已启用2FA）
1. 输入用户名和密码
2. 如果启用了2FA，系统会提示需要验证码
3. 打开认证APP获取当前6位验证码
4. 输入验证码完成登录

### 用户禁用2FA
1. 用户登录系统
2. 打开认证APP获取当前6位验证码
3. 调用 `/api/user/2fa/disable` 并提供验证码
4. 禁用成功后，登录不再需要验证码

## 推荐的认证APP
- Google Authenticator (iOS/Android)
- Microsoft Authenticator (iOS/Android)
- Authy (iOS/Android/Desktop)
- 1Password (支持TOTP)

## 安全建议
1. 启用2FA后，请妥善保管密钥（secret），建议备份到安全的地方
2. 如果丢失认证设备，可以联系管理员重置2FA
3. 管理员账户强烈建议启用2FA
4. 验证码有效期为30秒，过期后会自动刷新

## 技术实现
- 使用TOTP (Time-based One-Time Password) 算法
- 基于RFC 6238标准
- 验证码长度：6位数字
- 时间步长：30秒
- 使用库：github.com/pquerna/otp
