# 关于项目

## 项目来源

这个项目是在 [gocron](https://github.com/ouqiang/gocron) 的基础上开发改造而来。

我们非常喜欢原项目的设计理念，但由于原作者不再维护，我在此基础上做了大量的重构工作，旨在打造一个更现代、更安全、更易用的分布式定时任务管理系统。

### 主要改进

我们在原项目基础上进行了以下重大改进：

- **国际化支持**：完整的中文/英文界面切换支持
- **安全性增强**：
  - 新增双因素认证 (2FA)
  - 支持 TLS 双向认证
- **日志系统**：新增日志自动清理功能，支持数据库和文件日志
- **前端重构**：采用现代化的 Vue3 + Element Plus + Vite 架构
- **后端重构**：
  - Web 框架升级到 [Gin](https://github.com/gin-gonic/gin)
  - ORM 升级到 [GORM](https://github.com/go-gorm/gorm)
  - 实现了基于 GORM 的数据库自动迁移
- **数据库支持**：新增 SQLite 支持（纯 Go 驱动，无需 CGO），方便轻量级部署
- **Kubernetes 部署**：提供 Helm Chart，支持一键部署到 K8s 集群
- **自动化发布**：基于 GoReleaser + GitHub Actions，tag 触发自动构建全平台二进制并生成 changelog
- **Bug 修复**：修复了 Shell 任务停止功能、系统重启后孤立任务处理等关键问题
- **体验优化**：全新的 UI 设计和交互体验

## 技术栈

### 后端 (Go)

- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能的 HTTP Web 框架
- **ORM 框架**: [GORM](https://github.com/go-gorm/gorm) - 功能强大的 ORM 库
- **RPC 框架**: [gRPC](https://github.com/grpc/grpc-go) - 高性能、开源的通用 RPC 框架
- **Cron 库**: [cron](https://github.com/gocronx-team/cron) - 强大的定时任务库 (Forked)
- **JWT**: [golang-jwt](https://github.com/golang-jwt/jwt) - JSON Web Token 实现
- **2FA**: [otp](https://github.com/pquerna/otp) - 一次性密码库，用于双因素认证
- **CLI**: [cli](https://github.com/urfave/cli) - 简单、快速、有趣的命令行包
- **日志**: [logrus](https://github.com/sirupsen/logrus) - 结构化日志记录器
- **静态资源**: [statik](https://github.com/rakyll/statik) - 将静态文件嵌入 Go 二进制文件

### 前端 (Vue 3)

- **核心框架**: [Vue 3](https://github.com/vuejs/core) - 渐进式 JavaScript 框架
- **构建工具**: [Vite](https://github.com/vitejs/vite) - 下一代前端开发与构建工具
- **UI 组件库**: [Element Plus](https://github.com/element-plus/element-plus) - 基于 Vue 3 的组件库
- **路由管理**: [Vue Router](https://github.com/vuejs/router) - Vue.js 官方路由
- **状态管理**: [Pinia](https://github.com/vuejs/pinia) - Vue 的专属状态管理库
- **HTTP 客户端**: [Axios](https://github.com/axios/axios) - 基于 Promise 的 HTTP 客户端
- **国际化**: [Vue I18n](https://github.com/intlify/vue-i18n-next) - Vue.js 国际化插件
- **工具库**: [VueUse](https://github.com/vueuse/vueuse) - 必要的 Vue Composition API 工具集

## 致谢与支持

gocron 的发展离不开开源社区的支持。

如果您觉得这个项目对您有帮助，请给项目点个 **Star** ⭐，这不仅是对我们工作的认可，也能让更多人发现这个项目。

同时，我们非常欢迎您提交 **Issue** 反馈问题或建议，或者提交 **Pull Request** 参与代码贡献。让我们一起努力，打造一个更优秀的分布式定时任务管理系统！


