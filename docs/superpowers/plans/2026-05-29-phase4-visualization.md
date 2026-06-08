# 阶段四：可视化增强 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development

**Goal:** 实现地图展示、拓扑关系、层级树、机柜视图和自定义仪表盘。

**Architecture:** 前端为主（React + AntV/G6 + 高德地图 + ECharts），后端补充坐标和布局数据 API。

**Tech Stack:** React 18 + TypeScript + AntV/G6 + 高德地图 JS API + ECharts

---

## Task 1: 后端地图坐标与机柜布局 API

**Files:**
- Create: `cmdb-api/modules/dcim/handler.go` (append) — 地图坐标 CRUD
- Create: `cmdb-api/modules/dcim/dto.go` (append) — 坐标 DTO
- Modify: `cmdb-api/router/router.go` — 添加坐标路由

- [ ] IDC 地图坐标保存/查询
- [ ] 机柜布局查询（设备列表 + U位）

---

## Task 2: 关系图谱（@antv/g6）

**Files:**
- Create: `cmdb-ui/src/components/graph/RelationGraph.tsx`
- Create: `cmdb-ui/src/modules/core/CIRelationGraph.tsx`

- [ ] 以指定 CI 为中心，展开上下游关系
- [ ] 不同 CIType 用不同颜色/图标
- [ ] 支持力导向布局和层次布局切换
- [ ] 点击节点查看详情

---

## Task 3: 层级拓扑树

**Files:**
- Create: `cmdb-ui/src/components/graph/TopologyTree.tsx`
- Create: `cmdb-ui/src/modules/core/TopologyView.tsx`

- [ ] 集团 → 地域 → IDC → 机房 → 机柜 → 设备 层级展示
- [ ] 左侧折叠树，右侧详情面板
- [ ] 支持点击展开/折叠

---

## Task 4: IDC 地理分布地图

**Files:**
- Create: `cmdb-ui/src/modules/dcim/IDCMap.tsx`

- [ ] 使用高德地图 API 展示 IDC 坐标
- [ ] Marker 点击弹出 InfoWindow（IDC 名称、地址、设备数量）
- [ ] 聚合模式（缩放级别低时聚合显示数量）

---

## Task 5: 机柜视图 2.0

**Files:**
- Create: `cmdb-ui/src/modules/dcim/RackView.tsx`

- [ ] SVG 绘制机柜 U 位
- [ ] 设备用颜色区分状态（在线=绿、离线=红、维护=黄）
- [ ] 显示设备名称和U位范围
- [ ] 点击设备查看详情

---

## Task 6: 自定义仪表盘

**Files:**
- Create: `cmdb-ui/src/modules/dashboard/Dashboard.tsx`
- Create: `cmdb-ui/src/components/dashboard/StatCard.tsx`

- [ ] 各 CIType 数量统计（饼图/柱状图）
- [ ] 设备状态分布
- [ ] IP 地址使用趋势
- [ ] 最近发现任务状态
