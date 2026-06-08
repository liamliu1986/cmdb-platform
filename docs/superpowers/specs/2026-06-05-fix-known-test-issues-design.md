# 已知测试问题修复设计

## 背景

前后端联调测试中发现 5 个阻碍测试进行的已知问题，需一并修复。

## 问题清单

| # | 问题 | 影响 | 修复范围 |
|---|------|------|---------|
| 1 | Dashboard 图表数据硬编码 | 显示假数据，无意义 | 后端新增统计 API + 前端调用 |
| 2 | CIList 新建 CI 时属性硬编码 | 无论选择什么 CIType，表单字段固定为 name/description | 后端 CIType 详情预加载属性 + 前端动态获取 |
| 3 | 无登出功能 | 用户无法退出登录 | 前端添加登出按钮 |
| 4 | 监控集成菜单无对应路由 | 点击菜单 404 | 前端添加 MonitorDashboard 页面和路由 |
| 5 | SubnetList "新建子网"按钮无 onClick | 无法创建子网 | 前端添加新建子网弹窗 |

---

## 问题 1：Dashboard 统计图表

### 后端 API

**新增接口**：`GET /api/v1/stats/dashboard`

**响应格式**：
```json
{
  "code": 0,
  "data": {
    "total_ci": 0,
    "total_citype": 0,
    "total_rule": 0,
    "total_agent": 0,
    "ci_by_type": [
      { "name": "Server", "value": 10 }
    ],
    "ci_by_status": [
      { "status": "online", "value": 8 },
      { "status": "offline", "value": 2 }
    ]
  }
}
```

**实现**：在 `modules/core` 下新增 `stats.go`（handler + service + repository），通过 GORM Count 聚合查询。

### 前端修改

- Dashboard.tsx：加载时调用 `statsApi.getDashboardStats()`
- 统计卡片使用 `total_ci` 等字段
- ECharts 饼图使用 `ci_by_type`、柱状图使用 `ci_by_status`
- 当数据为空时显示空状态或占位文字

---

## 问题 2：CIList 动态属性

### 后端修改

**修改接口**：`GET /api/v1/citypes/:id`

**当前问题**：`GetCITypeByID` 只查询 `ci_types` 表，未预加载关联属性。

**修复**：Repository 中使用 GORM Preload 加载关联：
```go
database.DB.Preload("Attributes").First(&ct, id)
```

**注意**：`CIType` 模型需确认是否有 `Attributes` 关联字段。若无，需添加：
```go
type CIType struct {
    // ... 现有字段 ...
    Attributes []Attribute `gorm:"many2many:ci_type_attributes;" json:"attributes"`
}
```

### 前端修改

- CIList.tsx：选择 CIType 时调用 `coreApi.getCIType(typeId)`
- 从返回的 `attributes` 数组提取属性定义
- 传递给 CIForm 动态渲染表单字段

---

## 问题 3：登出功能

### 前端修改

- AppLayout.tsx Header 右侧添加"登出"按钮
- 点击执行：
  ```ts
  localStorage.removeItem('token')
  window.location.href = '/login'
  ```

---

## 问题 4：监控集成路由

### 前端修改

- 新建 `src/modules/integration/MonitorDashboard.tsx`
- 包含三个功能卡片：
  - Prometheus 查询：输入 PromQL，调用 `integrationApi.prometheusQuery()`
  - ELK 搜索：输入查询条件，调用 `integrationApi.elkSearch()`
  - 邮件测试：填写收件人/主题/内容，调用 `integrationApi.sendTestEmail()`
- router/index.tsx 添加 `/monitor` 路由

---

## 问题 5：子网创建弹窗

### 前端修改

- SubnetList.tsx 中"新建子网"按钮绑定 `onClick={() => setModalOpen(true)}`
- 新增弹窗表单字段：
  - 名称（required）
  - CIDR（required，如 192.168.1.0/24）
  - 描述
  - 父网 ID（Select，可选）
- 提交调用 `ipamApi.createSubnet(data)`，成功后刷新列表

---

## 测试策略

每修复一个问题，联调验证：
1. 前端操作 → 观察 UI 表现
2. curl 调用后端 API → 验证数据一致性
