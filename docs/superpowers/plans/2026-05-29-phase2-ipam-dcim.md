# 阶段二：IPAM + DCIM Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development

**Goal:** 实现 IPAM（IP 地址管理）和 DCIM（数据中心基础设施管理）模块，包括子网树、IP 分配、IDC/机房/机柜管理。

**Architecture:** 基于阶段一的 Core 模块扩展，IPAM/DCIM 使用 `cmdb_ipam` 和 `cmdb_dcim` Schema，同时初始化内置 CIType。

**Tech Stack:** Go 1.22 + Gin + GORM + PostgreSQL + React/TS + Ant Design

---

## 文件结构

```
cmdb-api/
├── modules/
│   ├── ipam/
│   │   ├── model.go
│   │   ├── migrate.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── dto.go
│   └── dcim/
│       ├── model.go
│       ├── migrate.go
│       ├── repository.go
│       ├── service.go
│       ├── handler.go
│       └── dto.go
└── router/
    └── router.go (update)

cmdb-ui/src/
├── modules/
│   ├── ipam/
│   │   ├── SubnetTree.tsx
│   │   ├── IPList.tsx
│   │   └── IPAllocate.tsx
│   └── dcim/
│       ├── IDCList.tsx
│       ├── ServerRoom.tsx
│       └── RackView.tsx
```

---

## Task 1: 内置 CIType 初始化

**Files:**
- Modify: `cmdb-api/modules/core/migrate.go`

- [ ] **Step 1: 在 migrate 中添加内置 CIType 初始化**

系统启动时自动创建以下内置 CIType：
- `Region`（地域）
- `IDC`（数据中心）
- `ServerRoom`（机房）
- `Rack`（机柜）
- `Subnet`（子网）
- `IPAddress`（IP 地址）

---

## Task 2: IPAM 数据库模型与迁移

**Files:**
- Create: `cmdb-api/modules/ipam/model.go`
- Create: `cmdb-api/modules/ipam/migrate.go`

- [ ] **Step 1: 创建 Subnet 模型**
- [ ] **Step 2: 创建 IPAddress 模型**
- [ ] **Step 3: 创建迁移函数**

---

## Task 3: IPAM Service + Handler

**Files:**
- Create: `cmdb-api/modules/ipam/repository.go`
- Create: `cmdb-api/modules/ipam/service.go`
- Create: `cmdb-api/modules/ipam/handler.go`
- Create: `cmdb-api/modules/ipam/dto.go`

- [ ] **Step 1: Subnet CRUD + 树形查询**
- [ ] **Step 2: IP 分配/释放/预留**
- [ ] **Step 3: CIDR 校验和无重叠检查**

---

## Task 4: DCIM 数据库模型与迁移

**Files:**
- Create: `cmdb-api/modules/dcim/model.go`
- Create: `cmdb-api/modules/dcim/migrate.go`

- [ ] **Step 1: 创建 IDC/ServerRoom/Rack 模型**
- [ ] **Step 2: 创建 RackLayout 模型（U 位占用）**

---

## Task 5: DCIM Service + Handler

**Files:**
- Create: `cmdb-api/modules/dcim/repository.go`
- Create: `cmdb-api/modules/dcim/service.go`
- Create: `cmdb-api/modules/dcim/handler.go`
- Create: `cmdb-api/modules/dcim/dto.go`

- [ ] **Step 1: IDC/机房/机柜 CRUD**
- [ ] **Step 2: 机柜设备上架/下架**
- [ ] **Step 3: 机柜容量计算**

---

## Task 6: 路由注册

**Files:**
- Modify: `cmdb-api/router/router.go`

- [ ] **Step 1: 注册 IPAM 路由**
- [ ] **Step 2: 注册 DCIM 路由**

---

## Task 7: 前端 IPAM 页面

**Files:**
- Create: `cmdb-ui/src/modules/ipam/SubnetTree.tsx`
- Create: `cmdb-ui/src/modules/ipam/IPList.tsx`
- Modify: `cmdb-ui/src/router/index.tsx`

- [ ] **Step 1: 子网树形展示**
- [ ] **Step 2: IP 列表管理**

---

## Task 8: 前端 DCIM 页面

**Files:**
- Create: `cmdb-ui/src/modules/dcim/IDCList.tsx`
- Create: `cmdb-ui/src/modules/dcim/RackView.tsx`
- Modify: `cmdb-ui/src/router/index.tsx`

- [ ] **Step 1: IDC 列表**
- [ ] **Step 2: 机柜视图 1.0**
