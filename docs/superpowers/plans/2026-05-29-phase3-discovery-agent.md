# 阶段三：自动发现 + Agent Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development

**Goal:** 实现 OneAgent 自动发现框架和 Master 端发现管理，支持云资源（阿里云/腾讯云/华为云/AWS）和服务器自动采集。

**Architecture:** 
- Master: Go/Gin + 发现规则管理 + 任务调度 + Agent 心跳管理
- Agent: Go 插件化架构，支持云厂商 SDK / SSH / SNMP
- 通信: HTTPS（注册/心跳/任务下发/结果上报）

**Tech Stack:** Go 1.22 + Gin + GORM + PostgreSQL + Redis

---

## 文件结构

```
cmdb-api/
├── modules/
│   └── discovery/
│       ├── model.go          # DiscoveryRule, DiscoveryTask, Agent
│       ├── migrate.go
│       ├── repository.go
│       ├── service.go        # 任务调度、规则执行
│       ├── handler.go
│       └── dto.go
cmdb-agent/                     # 独立 Go 项目
├── go.mod
├── main.go                     # Agent 入口
├── core/
│   ├── agent.go                # Agent 核心（注册、心跳、任务调度）
│   ├── config.go               # Agent 配置
│   └── reporter.go             # 结果上报
└── plugins/
    ├── interfaces.go           # 插件接口
    ├── aliyun/
    │   └── ecs.go              # 阿里云 ECS 插件
    ├── tencent/
    │   └── cvm.go              # 腾讯云 CVM 插件
    ├── huawei/
    │   └── ecs.go              # 华为云 ECS 插件
    ├── aws/
    │   └── ec2.go              # AWS EC2 插件
    └── server/
        └── local.go            # 服务器本地采集插件
cmdb-ui/src/
└── modules/
    └── discovery/
        ├── RuleList.tsx
        └── AgentList.tsx
```

---

## Task 1: Master 端发现模块数据库模型与迁移

**Files:**
- Create: `cmdb-api/modules/discovery/model.go`
- Create: `cmdb-api/modules/discovery/migrate.go`

模型：
- `DiscoveryRule` — 发现规则（类型、配置、调度 cron）
- `DiscoveryTask` — 发现任务（规则ID、状态、结果摘要）
- `DiscoveryResult` — 发现结果（任务ID、原始数据、匹配CIID）
- `Agent` — Agent 注册信息（名称、Token、IP、OS、最后心跳）

---

## Task 2: Master 端发现 Service + Handler

**Files:**
- Create: `cmdb-api/modules/discovery/repository.go`
- Create: `cmdb-api/modules/discovery/service.go`
- Create: `cmdb-api/modules/discovery/handler.go`
- Create: `cmdb-api/modules/discovery/dto.go`

- [ ] 发现规则 CRUD
- [ ] 手动执行任务
- [ ] Agent 注册/心跳
- [ ] 结果入库（调用 core 模块 CI 创建/更新 API）

---

## Task 3: 路由注册 + 前端发现页面

**Files:**
- Modify: `cmdb-api/router/router.go`
- Create: `cmdb-ui/src/api/discovery.ts`
- Create: `cmdb-ui/src/modules/discovery/RuleList.tsx`
- Create: `cmdb-ui/src/modules/discovery/AgentList.tsx`
- Modify: `cmdb-ui/src/router/index.tsx`
- Modify: `cmdb-ui/src/layouts/AppLayout.tsx`

---

## Task 4: OneAgent 核心框架

**Files:**
- Create: `cmdb-agent/go.mod`
- Create: `cmdb-agent/main.go`
- Create: `cmdb-agent/core/config.go`
- Create: `cmdb-agent/core/agent.go`
- Create: `cmdb-agent/core/reporter.go`

- [ ] Agent 注册（启动时向 Master 注册获取 Token）
- [ ] 心跳（每 30 秒上报状态）
- [ ] 任务调度（接收任务、调用插件、上报结果）
- [ ] HTTPS 通信 + Token 认证

---

## Task 5: Agent 插件接口 + 云厂商插件

**Files:**
- Create: `cmdb-agent/plugins/interfaces.go`
- Create: `cmdb-agent/plugins/aliyun/ecs.go`
- Create: `cmdb-agent/plugins/tencent/cvm.go`
- Create: `cmdb-agent/plugins/huawei/ecs.go`
- Create: `cmdb-agent/plugins/aws/ec2.go`

插件接口：
```go
type Discoverer interface {
    Name() string
    Type() string
    Init(config map[string]interface{}) error
    Discover(ctx context.Context) ([]Resource, error)
}
```

---

## Task 6: 服务器本地采集插件

**Files:**
- Create: `cmdb-agent/plugins/server/local.go`

采集内容：
- CPU 信息（型号、核心数）
- 内存信息（总大小）
- 磁盘信息（挂载点、大小）
- 网卡信息（IP、MAC）
- OS 版本

---

## Task 7: Agent 构建验证

**Files:**
- Create: `cmdb-agent/Makefile`

- [ ] 支持跨平台编译（Linux/Windows/macOS, x86_64/arm64）
- [ ] 验证 `go build` 成功
