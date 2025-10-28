# gocron 升级指南

## 升级方式

gocron 支持**自动升级**，用户只需要替换可执行文件并重启服务即可。

## 自动升级机制

系统启动时会自动检测版本并执行数据库升级：

1. 系统读取 `~/.gocron/conf/.version` 文件获取当前版本号
2. 对比代码中的版本号（当前为 150）
3. 如果代码版本号更高，自动执行数据库迁移
4. 升级完成后更新 `.version` 文件

## 升级步骤

### 1. 停止服务

```bash
# 停止 gocron 服务
# 如果使用 systemd
sudo systemctl stop gocron

# 或者直接 kill 进程
pkill gocron
```

### 2. 备份数据（重要！）

```bash
# 备份数据库
mysqldump -u root -p gocron > gocron_backup_$(date +%Y%m%d).sql

# 备份配置文件
cp -r ~/.gocron/conf ~/.gocron/conf_backup_$(date +%Y%m%d)
```

### 3. 替换可执行文件

```bash
# 下载新版本
wget https://github.com/gocronx-team/gocron/releases/download/vX.X.X/gocron-vX.X.X-linux-amd64.tar.gz

# 解压
tar -xzf gocron-vX.X.X-linux-amd64.tar.gz

# 替换旧文件
cp gocron /path/to/your/gocron
chmod +x /path/to/your/gocron
```

### 4. 启动服务

```bash
# 启动 gocron
./gocron web

# 或使用 systemd
sudo systemctl start gocron
```

### 5. 查看升级日志

启动时会在日志中看到升级信息：

```
版本升级开始, 当前版本号150
开始升级到v1.5.1 - 添加2FA支持
已升级到v1.5.1
已升级到最新版本151
```

## 版本历史

### v1.5.1 (版本号: 151)
- 新增：双因素认证(2FA)功能
- 数据库变更：
  - user 表新增 `two_factor_key` 字段
  - user 表新增 `two_factor_on` 字段

### v1.5 (版本号: 150)
- 前端使用Vue+ElementUI重构
- 任务通知功能增强
- 数据库变更：
  - task 表新增 `notify_keyword` 字段
  - setting 表新增通知模板配置

### v1.4 (版本号: 140)
- HTTP任务支持POST请求
- 后台手动停止运行中的shell任务
- 数据库变更：
  - task 表新增 `retry_interval` 字段
  - task 表新增 `http_method` 字段

### v1.3 (版本号: 130)
- 支持多用户登录
- 增加用户权限控制
- 数据库变更：
  - user 表删除 `deleted` 字段

### v1.2.2 (版本号: 122)
- 任务批量操作
- 调度器与任务节点支持HTTPS
- 数据库变更：
  - task 表新增 `tag` 字段

### v1.1 (版本号: 110)
- 任务可同时在多个节点上运行
- 数据库变更：
  - 新增 `task_host` 表
  - task 表删除 `host_id` 字段

## 常见问题

### Q: 升级失败怎么办？

A: 如果升级失败，系统会在日志中显示错误信息。可以：
1. 查看日志文件 `~/.gocron/log/cron.log`
2. 恢复数据库备份
3. 使用旧版本可执行文件

### Q: 可以跨版本升级吗？

A: 可以。系统会自动执行所有中间版本的升级脚本。例如从 v1.3 直接升级到 v1.5.1，会依次执行 v1.4、v1.5、v1.5.1 的升级脚本。

### Q: 如何查看当前版本？

A: 
```bash
# 查看程序版本
./gocron -v

# 查看数据库版本
cat ~/.gocron/conf/.version
```

### Q: 首次安装需要手动执行 SQL 吗？

A: 不需要。首次安装时，系统会自动创建所有表和字段。只有升级时才会执行迁移脚本。

### Q: 升级后需要重新配置吗？

A: 不需要。配置文件 `~/.gocron/conf/app.ini` 会保持不变。

## 注意事项

1. **务必备份数据库**：升级前一定要备份数据库，以防万一
2. **停止服务后再升级**：避免数据不一致
3. **查看升级日志**：确认升级成功
4. **测试功能**：升级后测试关键功能是否正常
5. **v1.2 版本特殊说明**：v1.2 版本不支持自动升级，需要手动处理

## 回滚操作

如果升级后出现问题，可以回滚：

```bash
# 1. 停止服务
sudo systemctl stop gocron

# 2. 恢复数据库
mysql -u root -p gocron < gocron_backup_YYYYMMDD.sql

# 3. 恢复配置文件
cp -r ~/.gocron/conf_backup_YYYYMMDD/* ~/.gocron/conf/

# 4. 使用旧版本可执行文件
cp /path/to/old/gocron /path/to/your/gocron

# 5. 启动服务
sudo systemctl start gocron
```

## 技术细节

升级逻辑位于：
- 主程序：`cmd/gocron/gocron.go` 中的 `upgradeIfNeed()` 函数
- 迁移脚本：`internal/models/migration.go` 中的各个 `upgradeForXXX()` 函数
- 版本管理：`internal/modules/app/app.go` 中的版本号处理

每个升级函数都在事务中执行，确保原子性。如果某个步骤失败，整个升级会回滚。
