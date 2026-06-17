# Introduction

gocron is a lightweight scheduled task centralized scheduling and management system developed in Go language, designed to replace Linux-crontab.

## Key Features

- **Web Interface Management**: Provides a user-friendly web management interface for easy task creation, editing, and monitoring
- **Scheduled Tasks**: Supports crontab time expressions with second-level precision
- **Task Execution Methods**:
  - Shell command execution
  - HTTP request invocation
- **Task Logs**: Detailed task execution logs with search capabilities
- **Email Notifications**: Email notifications for task execution failures
- **Permission Management**: Multi-user support with different users managing different tasks
- **Security Authentication**:
  - Two-Factor Authentication (2FA)
  - TLS mutual authentication
- **Task Dependencies**: Support for task dependencies
- **Failure Retry**: Automatic retry on task execution failure
- **Single Instance Execution**: Prevents duplicate task execution

## Quick Start

Please refer to the following documentation to get started with gocron:

- [Configuration](./configuration)
- [Scheduled Tasks](./scheduled-tasks)
- [Log Management](./log-management)
- [API Documentation](./api)
