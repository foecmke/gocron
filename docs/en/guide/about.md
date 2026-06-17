# About Project

## Project Origin

This project is developed and refactored based on [gocron](https://github.com/ouqiang/gocron).

We really appreciate the design philosophy of the original project, but since the original author is no longer maintaining it, we have done extensive refactoring work on this foundation to build a more modern, secure, and easy-to-use distributed scheduled task management system.

### Key Improvements

We have made significant improvements based on the original project:

- **Internationalization**: Complete Chinese/English interface switching support
- **Security Enhancements**:
  - Added Two-Factor Authentication (2FA)
  - Added TLS Mutual Authentication support
- **Log System**: Added automatic log cleanup, supporting both database and file logs
- **Frontend Refactoring**: Adopted modern Vue3 + Element Plus + Vite architecture
- **Backend Refactoring**:
  - Upgraded Web framework to [Gin](https://github.com/gin-gonic/gin)
  - Upgraded ORM to [GORM](https://github.com/go-gorm/gorm)
  - Implemented GORM-based automatic database migration
- **Database Support**: Added SQLite support (pure Go driver, no CGO required) for lightweight deployment
- **Kubernetes Deployment**: Helm Chart for one-click deployment to K8s clusters
- **Automated Release**: GoReleaser + GitHub Actions, tag-triggered multi-platform builds with auto-generated changelog
- **Bug Fixes**: Fixed critical issues such as Shell task stop functionality and orphaned task handling after restart
- **Experience Optimization**: Brand new UI design and interaction experience
 
## Tech Stack

### Backend (Go)

- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP Web framework
- **ORM Framework**: [GORM](https://github.com/go-gorm/gorm) - The fantastic ORM library for Golang
- **RPC Framework**: [gRPC](https://github.com/grpc/grpc-go) - High performance, open source universal RPC framework
- **Cron Library**: [cron](https://github.com/gocronx-team/cron) - Robust cron library (Forked)
- **JWT**: [golang-jwt](https://github.com/golang-jwt/jwt) - JSON Web Token implementation
- **2FA**: [otp](https://github.com/pquerna/otp) - One-time password library for 2FA
- **CLI**: [cli](https://github.com/urfave/cli) - Simple, fast, and fun package for building command line apps
- **Logger**: [logrus](https://github.com/sirupsen/logrus) - Structured logger for Golang
- **Static Assets**: [statik](https://github.com/rakyll/statik) - Embeds static files into a Go executable

### Frontend (Vue 3)

- **Core Framework**: [Vue 3](https://github.com/vuejs/core) - The Progressive JavaScript Framework
- **Build Tool**: [Vite](https://github.com/vitejs/vite) - Next Generation Frontend Tooling
- **UI Library**: [Element Plus](https://github.com/element-plus/element-plus) - A Vue 3 based component library
- **Router**: [Vue Router](https://github.com/vuejs/router) - The official router for Vue.js
- **State Management**: [Pinia](https://github.com/vuejs/pinia) - The Vue Store that you will enjoy using
- **HTTP Client**: [Axios](https://github.com/axios/axios) - Promise based HTTP client
- **I18n**: [Vue I18n](https://github.com/intlify/vue-i18n-next) - Internationalization plugin for Vue.js
- **Utilities**: [VueUse](https://github.com/vueuse/vueuse) - Collection of Essential Vue Composition Utilities

## Acknowledgments & Support

The development of gocron relies on the support of the open source community.

If you find this project helpful, please give it a **Star** ⭐. This is not only a recognition of our work but also helps more people discover this project.

We also warmly welcome you to submit **Issues** for feedback or suggestions, or submit **Pull Requests** to contribute code. Let's work together to build an even better distributed scheduled task management system!


