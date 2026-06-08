# 阶段五：外部集成 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development

**Goal:** 实现 Prometheus 监控指标查询、ELK 日志查询、邮件告警通知，将 CMDB 与运维系统打通。

**Architecture:** 后端代理层 + 前端 CI 详情页集成 Tab + 邮件 SMTP 配置。

---

## Task 1: Prometheus 集成 API

- 代理端点 `GET /api/v1/integration/prometheus/query`，转发查询到 Prometheus
- 按 CI IP 自动关联指标
- CI 详情页集成监控 Tab

## Task 2: ELK 日志查询 API

- 代理端点 `POST /api/v1/integration/elk/search`，转发查询到 Elasticsearch
- 按 CI hostname/IP 自动查询日志
- CI 详情页集成日志 Tab

## Task 3: 邮件告警配置

- SMTP 配置管理 API
- CI 变更 Webhook 触发邮件通知
- 告警规则配置页面

## Task 4: 前端 CI 详情页 + 集成 Tab

- CI 详情页（监控 tab + 日志 tab）
- 告警规则配置页
- 路由 + 菜单更新
