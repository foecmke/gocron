# 简介

gocron 是一个使用 Go 语言开发的轻量级定时任务集中调度和管理系统，用于替代 Linux-crontab。

## 主要特性

- **Web 界面**：直观的 Web 管理界面，方便任务的创建、编辑和监控
- **秒级精度**：支持精确到秒的 crontab 表达式
- **任务执行方式**：Shell 命令执行、HTTP 请求调用
- **高可用**：基于数据库锁的 Leader 选举，秒级自动故障转移
- **任务依赖**：支持任务间的依赖配置
- **失败重试**：可配置的失败重试策略
- **单实例运行**：防止任务重复执行
- **Agent 自动注册**：Linux/macOS 一键部署任务节点
- **权限管理**：完善的多用户与权限管理
- **安全认证**：双因素认证（2FA）与 TLS 双向认证
- **AI 辅助**：自然语言转 cron、失败日志 AI 诊断、AI 运维助手对话——基于任意 OpenAI 兼容模型(云端或自托管/本地)
- **MCP 支持**：通过 Model Context Protocol 供 AI 客户端(Claude Desktop、Cursor 等)远程管理,以 Web 端管理的访问令牌鉴权
- **多数据库**：支持 MySQL / PostgreSQL / SQLite
- **日志管理**：完整的执行日志，支持搜索与自动清理
- **通知**：支持邮件、Slack、Webhook

## 快速开始

请参考以下文档开始使用 gocron：

- [配置文件](./configuration)
- [定时任务](./scheduled-tasks)
- [日志管理](./log-management)
- [AI 功能](./ai-features)
- [API 文档](./api)
