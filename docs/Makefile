# Makefile：本地构建 / 预览 / 发布 gocron 文档到 Cloudflare Pages
# 在 docs/ 目录下执行（docs/ 是独立 pnpm 工作区）。

PROJECT_NAME := gocron-docs
OUTPUT_DIR := .vitepress/dist
PACKAGE_MANAGER := pnpm

.PHONY: help install dev build preview deploy clean

help:
	@echo "Available commands:"
	@echo "  make install  - 安装依赖"
	@echo "  make dev      - 启动本地开发服务器 (http://localhost:5173)"
	@echo "  make build    - 构建文档站点"
	@echo "  make preview  - 本地预览构建产物"
	@echo "  make deploy   - 构建并发布到 Cloudflare Pages"
	@echo "  make clean    - 清理构建产物"

install:
	$(PACKAGE_MANAGER) install

dev:
	$(PACKAGE_MANAGER) run docs:dev

build:
	$(PACKAGE_MANAGER) run docs:build

preview:
	$(PACKAGE_MANAGER) run docs:preview

# 构建并部署到 Cloudflare Pages（项目不存在时会创建）
deploy: build
	npx wrangler pages deploy $(OUTPUT_DIR) --project-name $(PROJECT_NAME) --branch production --no-bundle

clean:
	rm -rf $(OUTPUT_DIR)
