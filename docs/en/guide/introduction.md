# Introduction

gocron is a lightweight scheduled task centralized scheduling and management system developed in Go language, designed to replace Linux-crontab.

## Key Features

- **Web Interface**: Intuitive web UI for creating, editing, and monitoring tasks
- **Second-level Precision**: Crontab expressions with second-level precision
- **Task Execution Methods**: Shell command execution and HTTP request invocation
- **High Availability**: Database-lock-based leader election, automatic failover in seconds
- **Task Dependency**: Supports task dependency configuration
- **Task Retry**: Configurable retry policies for failed tasks
- **Single Instance Execution**: Prevents duplicate task execution
- **Agent Auto-Registration**: One-click task node installation for Linux/macOS
- **Access Control**: Comprehensive multi-user and permission management
- **Security Authentication**: Two-Factor Authentication (2FA) and TLS mutual authentication
- **AI Assist**: Natural-language to cron, AI failure-log diagnosis, and an AI ops chat assistant — backed by any OpenAI-compatible model (cloud or self-hosted/local)
- **MCP Support**: Remote management by AI clients (Claude Desktop, Cursor, etc.) via the Model Context Protocol, secured with web-managed access tokens
- **Multi-Database**: MySQL / PostgreSQL / SQLite support
- **Log Management**: Complete execution logs with search and auto-cleanup
- **Notifications**: Email, Slack, and Webhook support

## Quick Start

Please refer to the following documentation to get started with gocron:

- [Configuration](./configuration)
- [Scheduled Tasks](./scheduled-tasks)
- [Log Management](./log-management)
- [AI Features](./ai-features)
- [API Documentation](./api)
