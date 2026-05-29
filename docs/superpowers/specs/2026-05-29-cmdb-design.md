# CMDB 系统设计文档

## 1. 项目概述

### 1.1 背景与目标

构建一套企业级 CMDB（配置管理数据库），作为集团 IT 资产的权威记录系统，对接智能运维体系（监控、日志采集、日志分析），覆盖 IT 资产全生命周期管理。

### 1.2 核心功能需求

1. **用户管理与权限**：认证、授权（AWS IAM 模式）、角色管理
2. **CIType 定义**：灵活定义配置项类型，支持属性、关系、继承、触发器
3. **CI 实例管理**：为每个 CIType 自动生成 CRUD API，支持搜索、关系、审计
4. **自动发现**：Agent 插件化架构，支持阿里云、腾讯云、华为云、AWS 及服务器/网络设备
5. **IPAM**：IP 地址与子网管理
6. **DCIM**：企业组织、数据中心、机房、机柜管理，支持地图展示与机柜可视化

### 1.3 规模与部署

- CI 实例：10 万
- 并发用户：100
- 部署方式：Kubernetes + Helm
- 组织架构：中心化部署，管理涉及分子公司

## 2. 参考项目分析（veops/cmdb）

### 2.1 项目概况

- 开源 CMDB，Python/Flask + Vue.js + Ant Design Vue
- 数据库：MySQL + Redis，可选 Elasticsearch
- 核心能力：自定义 CIType、自动发现、多维视图、IPAM、DCIM、ACL 权限

### 2.2 功能对标

| 功能 | veops 开源版 | veops 企业版 | 本项目目标 |
|------|-------------|-------------|-----------|
| CIType 自定义 | 支持 | 支持 | 完整实现 |
| CI 关系管理 | 支持 | 支持 | 完整实现 |
| 自动发现（公有云） | 阿里云/腾讯云/华为云/AWS | + vCenter | 完整实现 |
| 自动发现（私有云） | 无 | OpenStack / Proxmox VE | 完整实现 |
| 自动发现（网络） | SNMP 基础 | + 更多设备 | 完整实现 |
| IPAM | 支持 | 支持 | 完整实现 |
| DCIM | 支持 | 支持 | 完整实现 |
| 地图展示 | 基础 | 增强 | 完整实现 |
| 仪表盘 | 基础 | 增强 | 完整实现 |
| 权限控制 | ACL | ACL | AWS IAM 模式 |

## 3. 技术选型

| 层级 | 技术 |
|------|------|
| 后端框架 | Go 1.22 + Gin |
| ORM | GORM |
| 数据库迁移 | golang-migrate |
| 数据库 | PostgreSQL 15+ |
| 缓存 | Redis 7+ |
| 搜索引擎 | Elasticsearch 8+ |
| 消息队列 | NATS / Redis Streams |
| 对象存储 | MinIO |
| 前端框架 | React 18 + TypeScript + Vite |
| UI 组件 | Ant Design 5.x |
| 状态管理 | Zustand + React Query |
| 图表 | ECharts + @antv/g6 |
| 地图 | 高德地图 JS API |
| Agent | Go 1.22（跨平台编译） |
| 部署 | Kubernetes + Helm 3 |
| 监控 | Prometheus + Grafana |

## 4. 系统整体架构

### 4.1 架构模式：模块化单体（Modular Monolith）

一个 Go HTTP 服务，内部按模块严格分层，数据库按模块分 Schema。代码边界清晰，未来可平滑拆分为微服务。

### 4.2 模块划分

```
cmdb-api (Go/Gin)
├── modules/auth/         # 用户、角色、权限（AWS IAM 模式）
├── modules/core/         # CIType 引擎、CI 实例、搜索、关系
├── modules/ipam/         # IP 地址管理
├── modules/dcim/         # 数据中心基础设施管理
├── modules/discovery/    # 自动发现规则与任务调度
├── modules/notify/       # Webhook、邮件、订阅
└── modules/integration/  # Prometheus、ELK 对接
```

### 4.3 模块依赖

```
auth (无依赖，基础层)
  │
  ▼
core (依赖 auth 做权限校验)
  │
  ├──► ipam  (依赖 core)
  ├──► dcim  (依赖 core)
  ├──► discovery (依赖 core 写入 CI)
  └──► notify (依赖 core 的事件)
```

模块间禁止直接访问对方的数据库表，必须通过 Service 接口交互。

### 4.4 部署拓扑

```
用户
  │ HTTPS
  ▼
Ingress (Nginx)
  ├── /api/* → cmdb-api-service (ClusterIP) → cmdb-api Pod x2
  └── /*     → cmdb-ui-service  (ClusterIP) → cmdb-ui Pod x2 (Nginx)

内部依赖：
  PostgreSQL (StatefulSet)  - 主数据
  Redis      (StatefulSet)  - 缓存/会话/锁
  ES         (StatefulSet)  - 搜索/聚合
  NATS       - 异步消息
  MinIO      - 文件存储
```

## 5. 数据库设计

### 5.1 Schema 划分

| Schema | 用途 |
|--------|------|
| cmdb_auth | 用户、角色、权限 |
| cmdb_core | CIType、CI、关系、属性、触发器、日志 |
| cmdb_ipam | 子网、IP 地址 |
| cmdb_dcim | 地图坐标、机柜布局 |
| cmdb_discovery | 发现规则、任务、Agent |

### 5.2 核心设计决策：CI 属性值存储

采用 PostgreSQL JSONB 列存储 CI 属性值。

```sql
CREATE TABLE cmdb_core.cis (
    id BIGSERIAL PRIMARY KEY,
    type_id INT NOT NULL,
    status VARCHAR(16) DEFAULT 'active',
    attr_values JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    updated_by VARCHAR(64),
    is_auto_discovery BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_cis_type ON cmdb_core.cis(type_id);
CREATE INDEX idx_cis_attr_gin ON cmdb_core.cis USING GIN (attr_values jsonb_path_ops);
```

对需要唯一或索引的属性，创建表达式索引：
```sql
CREATE UNIQUE INDEX idx_cis_sn ON cmdb_core.cis ((attr_values->>'sn'))
    WHERE attr_values ? 'sn';
```

### 5.3 cmdb_auth 核心表

| 表名 | 说明 |
|------|------|
| users | 用户基础信息 |
| roles | 角色定义 |
| role_relations | 角色继承 |
| resource_types | 资源类型（CIType、CI、Subnet 等） |
| resources | 具体资源实例 |
| resource_groups | 资源组（批量授权） |
| permissions | 权限类型（create/read/update/delete/execute） |
| role_permissions | 角色资源权限绑定 |
| user_roles | 用户角色绑定 |

### 5.4 cmdb_core 核心表

| 表名 | 说明 |
|------|------|
| relation_types | 关系类型（contains、depends_on、deployed_on） |
| ci_types | CI 类型定义 |
| ci_type_inheritance | 继承关系 |
| attributes | 属性定义（类型、校验、默认值） |
| ci_type_attributes | CIType 与属性关联 |
| ci_type_relations | CIType 间关系模板 |
| ci_type_triggers | 触发器配置 |
| cis | CI 实例（JSONB 存储属性值） |
| ci_relations | CI 实例间关系 |
| operation_logs | 操作审计日志 |

### 5.5 cmdb_ipam 核心表

| 表名 | 说明 |
|------|------|
| subnets | 子网树（parent_id + path） |
| ip_addresses | IP 地址状态（free/allocated/reserved） |
| ipam_histories | IP 操作历史 |

### 5.6 cmdb_dcim 核心表

| 表名 | 说明 |
|------|------|
| location_coords | 地图坐标（lat, lng） |
| rack_layouts | 机柜 U 位占用 |
| idc_floors | 机房平面图 |

## 6. API 设计规范

### 6.1 统一响应格式

```json
{ "code": 0, "message": "success", "data": {} }
```

### 6.2 核心 API

**认证：**
- POST `/api/v1/auth/login` — 登录
- POST `/api/v1/auth/logout` — 退出
- GET `/api/v1/auth/me` — 当前用户

**用户权限：**
- CRUD `/api/v1/users`
- CRUD `/api/v1/roles`
- POST `/api/v1/roles/:id/permissions`

**CIType：**
- CRUD `/api/v1/citypes`
- GET `/api/v1/citypes/:id/attributes`
- GET `/api/v1/citypes/:id/relations`

**CI 实例：**
- GET `/api/v1/ci/s?q=...` — 搜索（核心）
- POST `/api/v1/ci` — 创建
- PUT `/api/v1/ci/:id` — 更新
- DELETE `/api/v1/ci/:id` — 删除
- GET `/api/v1/ci/:id` — 详情
- GET `/api/v1/ci/:id/relations` — 关系
- GET `/api/v1/ci/:id/history` — 历史

**搜索查询语法：**
- `_type:Server` — 指定 CIType
- `attr:value` — 属性过滤
- 逗号分隔 — AND
- `-attr:value` — OR
- `~attr:value` — NOT
- `attr:(v1;v2)` — IN
- `attr:[a _TO_ b]` — 范围
- `attr:>5` — 比较

**IPAM：**
- CRUD `/api/v1/ipam/subnets`
- GET `/api/v1/ipam/subnets/:id/ips`
- POST `/api/v1/ipam/subnets/:id/allocate`
- POST `/api/v1/ipam/ips/:id/release`

**DCIM：**
- CRUD `/api/v1/dcim/locations`
- GET `/api/v1/dcim/racks/:id/layout`
- POST `/api/v1/dcim/racks/:id/devices`

**发现：**
- CRUD `/api/v1/discovery/rules`
- POST `/api/v1/discovery/rules/:id/execute`
- GET `/api/v1/discovery/tasks`
- GET `/api/v1/discovery/agents`

## 7. 前端架构

### 7.1 目录结构

```
src/
├── api/           # 按模块封装的 API
├── components/    # 公共组件
│   ├── ci/        # CIForm、CITable、CITypeTree
│   ├── graph/     # RelationGraph、TopologyTree、IDCMap、RackView
│   └── search/    # CIQueryBuilder
├── hooks/         # useAuth、useCIType、usePermission
├── layouts/       # MainLayout、AuthLayout
├── modules/       # 按模块的页面
│   ├── auth/
│   ├── citype/
│   ├── ci/
│   ├── ipam/
│   ├── dcim/
│   ├── discovery/
│   ├── dashboard/
│   └── settings/
├── router/        # React Router v6
├── stores/        # Zustand（auth、app、citype）
├── types/         # TypeScript 类型定义
└── utils/         # 工具函数
```

### 7.2 核心组件

| 组件 | 用途 |
|------|------|
| CITypeDesigner | 可视化模型设计器（拖拽配置属性、关系） |
| CIForm | 动态表单（根据 CIType 属性自动生成） |
| CITable | 动态表格（列的显示/隐藏/顺序可配置） |
| RelationGraph | 关系图谱（@antv/g6） |
| TopologyTree | 层级拓扑树（集团→地域→IDC→机房→机柜→设备） |
| IDCMap | 高德地图展示 IDC 地理分布 |
| RackView | SVG 机柜 2D 平面图（U 位管理） |

### 7.3 状态管理

- **Zustand**：全局 UI 状态（主题、语言、用户信息）、CIType 元数据缓存
- **React Query**：服务器状态（列表、详情），自带缓存、去重、自动刷新

## 8. Agent 架构

### 8.1 插件化设计

```go
type Discoverer interface {
    Name() string
    Type() string  // cloud/server/network/custom
    Init(config map[string]interface{}) error
    Discover(ctx context.Context) ([]Resource, error)
}
```

### 8.2 内置插件

| 插件 | 发现内容 |
|------|----------|
| aliyun_ecs/rds/slb | 阿里云资源 |
| tencent_cvm | 腾讯云资源 |
| huawei_ecs | 华为云资源 |
| aws_ec2 | AWS 资源 |
| server_local/ssh | 物理机/虚拟机硬件信息 |
| network_snmp | 交换机、路由器、端口 |
| vmware_vcenter | VM、ESXi、Datastore |
| openstack | OpenStack 实例、Flavor、Image、Network |
| proxmox_ve | Proxmox VE VM、LXC、Node、Storage |
| custom_script | 用户自定义脚本 |

### 8.3 通信协议

默认 HTTPS，可选 gRPC。Agent 启动时注册获取 Token，定期心跳，接收任务，执行后上报结果。

### 8.4 安全

- Token 认证（`X-Agent-Token`）
- HTTPS 强制 TLS 1.2+
- 敏感配置 AES 加密
- **只上报数据，不接收执行远程命令**
- IP 白名单

## 9. K8s 部署架构

### 9.1 组件资源

| 组件 | 类型 | 副本 | CPU | 内存 | 存储 |
|------|------|------|-----|------|------|
| cmdb-api | Deployment | 2 | 500m | 512Mi | 无 |
| cmdb-ui | Deployment | 2 | 100m | 128Mi | 无 |
| postgresql | StatefulSet | 1 | 1000m | 2Gi | 50Gi |
| redis | StatefulSet | 1 | 500m | 512Mi | 10Gi |
| elasticsearch | StatefulSet | 1 | 1000m | 2Gi | 50Gi |

### 9.2 Helm Chart 核心设计

- `values.yaml`：默认配置
- `values-production.yaml`：生产覆盖（外接数据库）
- `templates/job-migrate.yaml`：Post-Install/Post-Upgrade Hook 执行数据库迁移
- `templates/ingress.yaml`：`/api/*` → API, `/*` → UI
- `templates/servicemonitor.yaml`：Prometheus 监控采集

### 9.3 部署命令

```bash
helm install cmdb ./cmdb -n cmdb --create-namespace -f values-production.yaml
helm upgrade cmdb ./cmdb -n cmdb
```

## 10. 分阶段实施路线图

### 阶段一：核心底座

**目标**：用户权限 + CIType 引擎 + CI 管理

- **后端**：auth 模块（AWS IAM）、core 模块（CIType/CI/搜索/日志）
- **前端**：登录页、用户/角色/权限管理、CIType 设计器、动态 CI 管理页
- **验收**：能创建 CIType，能增删改查 CI，权限控制生效，搜索语法正确

### 阶段二：基础设施管理

**目标**：IPAM + DCIM

- **后端**：ipam 模块（子树/分配算法）、dcim 模块、内置 CIType 初始化
- **前端**：IPAM 子网树、IP 管理、DCIM 资源列表、机柜视图 1.0
- **验收**：能划分子网、自动分配 IP，能管理机柜和设备上架

### 阶段三：自动发现

**目标**：Agent + 云资源/服务器自动采集

- **后端**：discovery 模块（规则/调度/结果入库）、云厂商 SDK 集成
- **前端**：发现规则配置、任务列表、Agent 管理
- **Agent**：OneAgent 框架 + 云厂商/服务器/网络插件
- **验收**：配置 AK/SK 后自动拉取云资源，Agent 上报硬件信息，重复执行不创建重复 CI

### 阶段四：可视化增强

**目标**：地图 + 拓扑 + 仪表盘

- **后端**：地图坐标管理、自定义仪表盘 API、订阅功能
- **前端**：IDC 地图、机房平面图、机柜视图 2.0（拖拽）、关系图谱、层级拓扑、仪表盘
- **验收**：地图展示 IDC 分布，机柜支持拖拽上架，关系图谱展示 3 层上下游

### 阶段五：外部集成

**目标**：Prometheus + ELK + 邮件告警

- **后端**：Prometheus 代理查询、ELK 日志代理查询、SMTP 邮件告警
- **前端**：CI 详情页监控 tab、日志 tab、告警规则配置
- **验收**：CI 详情页能看到实时 CPU/内存图表和最近日志，收到变更邮件通知

## 11. 安全设计

- JWT Token 认证，支持刷新
- AWS IAM 模式细粒度授权（角色→资源类型→资源→权限）
- 密码 bcrypt 加密存储
- 操作审计日志全量记录（操作人、时间、旧值、新值）
- SQL 注入防护：GORM + 参数化查询
- XSS 防护：前端转义 + Content-Security-Policy
- Agent 只上报不执行命令

## 12. 监控与可观测性

- **指标**：Prometheus `/metrics` 暴露 HTTP 请求数/耗时、CI 总量、数据库连接数
- **日志**：JSON 格式输出到 stdout，Fluent Bit 采集到 ELK
- **链路**：请求头携带 `X-Request-ID`，全链路追踪
- **告警**：Pod 崩溃、API 错误率升高、数据库连接池耗尽

---

*文档版本：v1.0*
*日期：2026-05-29*
