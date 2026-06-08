# 已知测试问题修复 - 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复联调测试中发现的 5 个前后端问题，使 Dashboard、CIList、登出、监控路由、子网创建功能正常工作。

**Architecture:** 后端新增统计 API 和 CIType 属性预加载；前端修复硬编码数据、添加缺失交互和路由。

**Tech Stack:** Go + Gin + GORM + PostgreSQL / React + TypeScript + Ant Design + Vite + ECharts

---

## 文件映射

### 后端修改
| 文件 | 职责 |
|------|------|
| `cmdb-api/modules/core/model.go` | 为 CIType 添加 Attributes 关联字段 |
| `cmdb-api/modules/core/repository.go` | 新增 GetCITypeWithAttributes、GetDashboardStats 查询 |
| `cmdb-api/modules/core/service.go` | 新增 GetCITypeWithAttributes、GetDashboardStats 业务逻辑 |
| `cmdb-api/modules/core/handler.go` | 修改 GetCIType 返回属性，新增 DashboardStats handler |
| `cmdb-api/router/router.go` | 注册 /stats/dashboard 路由 |

### 前端修改
| 文件 | 职责 |
|------|------|
| `cmdb-ui/src/api/core.ts` | 新增 getCIType(id)、getDashboardStats() |
| `cmdb-ui/src/api/stats.ts` | 新建 stats API 模块 |
| `cmdb-ui/src/modules/dashboard/Dashboard.tsx` | 调用 stats API 渲染真实数据 |
| `cmdb-ui/src/modules/core/CIList.tsx` | 选择 CIType 时动态获取属性 |
| `cmdb-ui/src/layouts/AppLayout.tsx` | 添加登出按钮 |
| `cmdb-ui/src/modules/integration/MonitorDashboard.tsx` | 新建监控集成页面 |
| `cmdb-ui/src/router/index.tsx` | 添加 /monitor 路由 |
| `cmdb-ui/src/modules/ipam/SubnetList.tsx` | 添加新建子网弹窗 |

---

## Task 1: 后端 - CIType 模型添加 Attributes 关联

**Files:**
- Modify: `cmdb-api/modules/core/model.go:9-23`

- [ ] **Step 1: 修改 CIType 结构体**

```go
type CIType struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:32;uniqueIndex;not null" json:"name"`
	Alias        string         `gorm:"size:32;not null" json:"alias"`
	UniqueAttrID uint           `gorm:"not null" json:"unique_attr_id"`
	Icon         string         `gorm:"size:255" json:"icon"`
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	IsBuiltin    bool           `gorm:"default:false" json:"is_builtin"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Attributes   []Attribute    `gorm:"many2many:ci_type_attributes;" json:"attributes,omitempty"`
}
```

- [ ] **Step 2: Commit**

```bash
git add cmdb-api/modules/core/model.go
git commit -m "feat: add Attributes association to CIType model"
```

---

## Task 2: 后端 - Repository 新增查询方法

**Files:**
- Modify: `cmdb-api/modules/core/repository.go`

- [ ] **Step 1: 新增 GetCITypeWithAttributes 方法**

在 repository.go 中 `GetCITypeByID` 方法之后添加：

```go
func (r *CoreRepository) GetCITypeWithAttributes(id uint) (*CIType, error) {
	var ct CIType
	err := database.DB.Preload("Attributes").First(&ct, id).Error
	return &ct, err
}
```

- [ ] **Step 2: 新增 Dashboard 统计查询方法**

在 repository.go 末尾添加：

```go
func (r *CoreRepository) CountCIsByType() ([]struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}, error) {
	var results []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	err := database.DB.Raw(`
		SELECT ct.name, COUNT(c.id) as value
		FROM cmdb_core.ci_types ct
		LEFT JOIN cmdb_core.cis c ON c.ci_type_id = ct.id
		GROUP BY ct.id, ct.name
		ORDER BY value DESC
	`).Scan(&results).Error
	return results, err
}

func (r *CoreRepository) CountCIsByStatus() ([]struct {
	Status string `json:"status"`
	Value  int64  `json:"value"`
}, error) {
	var results []struct {
		Status string `json:"status"`
		Value  int64  `json:"value"`
	}
	err := database.DB.Raw(`
		SELECT status, COUNT(*) as value
		FROM cmdb_core.cis
		GROUP BY status
		ORDER BY value DESC
	`).Scan(&results).Error
	return results, err
}

func (r *CoreRepository) CountTotal(table string) (int64, error) {
	var count int64
	err := database.DB.Table(table).Count(&count).Error
	return count, err
}
```

- [ ] **Step 3: Commit**

```bash
git add cmdb-api/modules/core/repository.go
git commit -m "feat: add GetCITypeWithAttributes and dashboard stats queries"
```

---

## Task 3: 后端 - Service 层新增业务方法

**Files:**
- Modify: `cmdb-api/modules/core/service.go`

- [ ] **Step 1: 新增 GetCITypeWithAttributes 方法**

在 service.go 中 `GetCIType` 方法之后添加：

```go
func (s *CoreService) GetCITypeWithAttributes(id uint) (*CIType, error) {
	return s.repo.GetCITypeWithAttributes(id)
}
```

- [ ] **Step 2: 新增 Dashboard 统计服务方法**

在 service.go 末尾添加：

```go
type DashboardStats struct {
	TotalCI      int64 `json:"total_ci"`
	TotalCIType  int64 `json:"total_citype"`
	TotalRule    int64 `json:"total_rule"`
	TotalAgent   int64 `json:"total_agent"`
	CIByType     []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	} `json:"ci_by_type"`
	CIByStatus   []struct {
		Status string `json:"status"`
		Value  int64  `json:"value"`
	} `json:"ci_by_status"`
}

func (s *CoreService) GetDashboardStats() (*DashboardStats, error) {
	totalCI, _ := s.repo.CountTotal("cmdb_core.cis")
	totalCIType, _ := s.repo.CountTotal("cmdb_core.ci_types")
	totalRule, _ := s.repo.CountTotal("cmdb_discovery.rules")
	totalAgent, _ := s.repo.CountTotal("cmdb_discovery.agents")
	ciByType, _ := s.repo.CountCIsByType()
	ciByStatus, _ := s.repo.CountCIsByStatus()

	return &DashboardStats{
		TotalCI:     totalCI,
		TotalCIType: totalCIType,
		TotalRule:   totalRule,
		TotalAgent:  totalAgent,
		CIByType:    ciByType,
		CIByStatus:  ciByStatus,
	}, nil
}
```

- [ ] **Step 3: Commit**

```bash
git add cmdb-api/modules/core/service.go
git commit -m "feat: add GetCITypeWithAttributes and GetDashboardStats services"
```

---

## Task 4: 后端 - Handler 修改和新增

**Files:**
- Modify: `cmdb-api/modules/core/handler.go`

- [ ] **Step 1: 修改 GetCIType 返回带属性的数据**

替换现有的 `GetCIType` 方法：

```go
func (h *CoreHandler) GetCIType(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ct, err := h.svc.GetCITypeWithAttributes(uint(id))
	if err != nil {
		response.Error(c, 404, "ci_type not found")
		return
	}
	response.Success(c, ct)
}
```

- [ ] **Step 2: 新增 DashboardStats handler 方法**

在 handler.go 末尾添加：

```go
func (h *CoreHandler) DashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}
```

- [ ] **Step 3: Commit**

```bash
git add cmdb-api/modules/core/handler.go
git commit -m "feat: update GetCIType to preload attributes, add DashboardStats handler"
```

---

## Task 5: 后端 - 注册 Dashboard 路由

**Files:**
- Modify: `cmdb-api/router/router.go`

- [ ] **Step 1: 在 authorized 路由组中添加 stats 路由**

在 router.go 的 authorized 组中，在 CIType 路由之前添加：

```go
// Stats
authorized.GET("/stats/dashboard", coreHandler.DashboardStats)
```

- [ ] **Step 2: Commit**

```bash
git add cmdb-api/router/router.go
git commit -m "feat: register GET /api/v1/stats/dashboard route"
```

---

## Task 6: 后端 - 重启 API 验证

- [ ] **Step 1: 重启 API**

```bash
# 停止旧进程
powershell -Command "Stop-Process -Id (Get-Content /tmp/cmdb-api.pid) -Force"
sleep 2
# 启动新进程
cd /e/Projects/cmdb/cmdb-api && DB_PASSWORD=cmdb DB_USER=cmdb DB_NAME=cmdb DB_HOST=localhost DB_PORT=5432 REDIS_HOST=localhost REDIS_PORT=6379 go run main.go > /tmp/cmdb-api.log 2>&1 &
echo $! > /tmp/cmdb-api.pid
sleep 5
```

- [ ] **Step 2: 验证 Dashboard Stats API**

```bash
curl -s http://localhost:8080/api/v1/stats/dashboard \
  -H "Authorization: Bearer $(curl -s -X POST http://localhost:8080/api/v1/auth/login -H 'Content-Type: application/json' -d '{"username":"admin","password":"admin123"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)"
```

**预期输出：**
```json
{"code":0,"data":{"total_ci":0,"total_citype":0,"total_rule":0,"total_agent":0,"ci_by_type":[],"ci_by_status":[]}}
```

- [ ] **Step 3: 验证 CIType 详情带属性**

先创建一个 CIType 再测试：
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login -H 'Content-Type: application/json' -d '{"username":"admin","password":"admin123"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
curl -s http://localhost:8080/api/v1/citypes/1 -H "Authorization: Bearer $TOKEN"
# 预期返回 code=0, data.attributes 字段存在（即使为空数组）
```

---

## Task 7: 前端 - 新增 stats API 模块

**Files:**
- Create: `cmdb-ui/src/api/stats.ts`

- [ ] **Step 1: 创建 stats API**

```typescript
import client from './client'

export const statsApi = {
  getDashboardStats: () => client.get('/stats/dashboard'),
}
```

- [ ] **Step 2: Commit**

```bash
git add cmdb-ui/src/api/stats.ts
git commit -m "feat: add stats API module"
```

---

## Task 8: 前端 - 修改 Dashboard 调用真实数据

**Files:**
- Modify: `cmdb-ui/src/modules/dashboard/Dashboard.tsx`

- [ ] **Step 1: 导入 stats API 并替换数据获取逻辑**

```tsx
import { useEffect, useState } from 'react'
import { Card, Col, Row, Statistic, Empty } from 'antd'
import { DatabaseOutlined, CloudServerOutlined, WifiOutlined, TeamOutlined } from '@ant-design/icons'
import ReactECharts from 'echarts-for-react'
import { statsApi } from '@/api/stats'

export default function Dashboard() {
  const [stats, setStats] = useState<any>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const fetchStats = async () => {
      setLoading(true)
      try {
        const res: any = await statsApi.getDashboardStats()
        if (res.code === 0) setStats(res.data)
      } catch {
        // ignore
      } finally {
        setLoading(false)
      }
    }
    fetchStats()
  }, [])

  const ciByTypeData = stats?.ci_by_type?.length
    ? stats.ci_by_type.map((item: any) => ({ value: item.value, name: item.name }))
    : []

  const ciByStatusData = stats?.ci_by_status?.length
    ? stats.ci_by_status.map((item: any) => ({ value: item.value, name: item.status }))
    : []

  const pieOption = ciByTypeData.length ? {
    title: { text: 'CI 类型分布', left: 'center' },
    tooltip: { trigger: 'item' },
    series: [{ type: 'pie', radius: ['40%', '70%'], data: ciByTypeData }],
  } : null

  const barOption = ciByStatusData.length ? {
    title: { text: '设备状态分布', left: 'center' },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ciByStatusData.map((d: any) => d.name) },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: ciByStatusData.map((d: any) => ({ value: d.value })) }],
  } : null

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic title="CI 总数" value={stats?.total_ci ?? 0} prefix={<DatabaseOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic title="CIType" value={stats?.total_citype ?? 0} prefix={<CloudServerOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic title="发现规则" value={stats?.total_rule ?? 0} prefix={<WifiOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic title="Agent 在线" value={stats?.total_agent ?? 0} prefix={<TeamOutlined />} />
          </Card>
        </Col>
      </Row>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Card>
            {pieOption ? (
              <ReactECharts option={pieOption} style={{ height: 350 }} />
            ) : (
              <Empty description="暂无 CI 类型分布数据" style={{ height: 350, display: 'flex', flexDirection: 'column', justifyContent: 'center' }} />
            )}
          </Card>
        </Col>
        <Col span={12}>
          <Card>
            {barOption ? (
              <ReactECharts option={barOption} style={{ height: 350 }} />
            ) : (
              <Empty description="暂无设备状态分布数据" style={{ height: 350, display: 'flex', flexDirection: 'column', justifyContent: 'center' }} />
            )}
          </Card>
        </Col>
      </Row>
    </div>
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add cmdb-ui/src/modules/dashboard/Dashboard.tsx
git commit -m "feat: dashboard uses real stats API instead of hardcoded data"
```

---

## Task 9: 前端 - CIList 动态获取属性

**Files:**
- Modify: `cmdb-ui/src/modules/core/CIList.tsx`
- Modify: `cmdb-ui/src/api/core.ts`

- [ ] **Step 1: core.ts 添加 getCIType**

```typescript
export const coreApi = {
  createCIType: (data: any) => client.post('/citypes', data),
  listCITypes: () => client.get('/citypes'),
  getCIType: (id: number) => client.get(`/citypes/${id}`),
  createCI: (data: any) => client.post('/ci', data),
  getCI: (id: number) => client.get(`/ci/${id}`),
  searchCI: (params: { q?: string; page?: number; page_size?: number; sort?: string }) =>
    client.get('/ci/s', { params }),
}
```

- [ ] **Step 2: 修改 CIList.tsx 中 handleCITypeChange**

替换 `handleCITypeChange` 方法：

```tsx
const handleCITypeChange = async (typeId: number) => {
  const ct = ciTypes.find((t: any) => t.id === typeId)
  setSelectedCIType(ct)
  try {
    const res: any = await coreApi.getCIType(typeId)
    if (res.code === 0 && res.data?.attributes) {
      setCITypeAttrs(res.data.attributes.map((attr: any) => ({
        name: attr.name,
        alias: attr.alias,
        value_type: attr.value_type,
        is_required: attr.is_required ?? false,
      })))
    } else {
      setCITypeAttrs([])
    }
  } catch {
    setCITypeAttrs([])
  }
  form.resetFields()
}
```

- [ ] **Step 3: Commit**

```bash
git add cmdb-ui/src/api/core.ts cmdb-ui/src/modules/core/CIList.tsx
git commit -m "feat: CIList dynamically fetches attributes from selected CIType"
```

---

## Task 10: 前端 - 添加登出功能

**Files:**
- Modify: `cmdb-ui/src/layouts/AppLayout.tsx`

- [ ] **Step 1: 导入 LogoutOutlined 并添加登出按钮**

```tsx
import { Layout, Menu, Button } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { HomeOutlined, BuildOutlined, DatabaseOutlined, AppstoreOutlined, ApartmentOutlined, BankOutlined, RadarChartOutlined, RobotOutlined, ClusterOutlined, EnvironmentOutlined, DashboardOutlined, LogoutOutlined } from '@ant-design/icons'

const { Header, Sider, Content } = Layout

export default function AppLayout() {
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    localStorage.removeItem('token')
    window.location.href = '/login'
  }

  // ... menuItems 不变 ...

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#fff', padding: '0 24px', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0 }}>CMDB</h2>
        <Button icon={<LogoutOutlined />} onClick={handleLogout}>登出</Button>
      </Header>
      <Layout>
        <Sider theme="light" style={{ borderRight: '1px solid #f0f0f0' }}>
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => navigate(key)}
          />
        </Sider>
        <Content style={{ padding: 24, background: '#f5f5f5' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add cmdb-ui/src/layouts/AppLayout.tsx
git commit -m "feat: add logout button to AppLayout"
```

---

## Task 11: 前端 - 新建监控集成页面

**Files:**
- Create: `cmdb-ui/src/modules/integration/MonitorDashboard.tsx`
- Modify: `cmdb-ui/src/router/index.tsx`

- [ ] **Step 1: 创建 MonitorDashboard.tsx**

```tsx
import { useState } from 'react'
import { Card, Input, Button, Tabs, message } from 'antd'
import { integrationApi } from '@/api/integration'

export default function MonitorDashboard() {
  const [promQuery, setPromQuery] = useState('')
  const [promResult, setPromResult] = useState('')
  const [elkQuery, setElkQuery] = useState('')
  const [elkResult, setElkResult] = useState('')
  const [emailTo, setEmailTo] = useState('')
  const [emailSubject, setEmailSubject] = useState('')
  const [emailBody, setEmailBody] = useState('')

  const handlePromQuery = async () => {
    try {
      const res: any = await integrationApi.prometheusQuery(promQuery)
      setPromResult(JSON.stringify(res.data, null, 2))
    } catch (err: any) {
      message.error(err.message || '查询失败')
    }
  }

  const handleElkSearch = async () => {
    try {
      const res: any = await integrationApi.elkSearch({ query: elkQuery })
      setElkResult(JSON.stringify(res.data, null, 2))
    } catch (err: any) {
      message.error(err.message || '搜索失败')
    }
  }

  const handleSendEmail = async () => {
    try {
      await integrationApi.sendTestEmail({ to: emailTo, subject: emailSubject, body: emailBody })
      message.success('邮件发送成功')
    } catch (err: any) {
      message.error(err.message || '发送失败')
    }
  }

  const items = [
    {
      key: 'prometheus',
      label: 'Prometheus',
      children: (
        <Card title="PromQL 查询">
          <Input.TextArea
            value={promQuery}
            onChange={(e) => setPromQuery(e.target.value)}
            placeholder="输入 PromQL，如: up"
            rows={3}
          />
          <Button type="primary" onClick={handlePromQuery} style={{ marginTop: 8 }}>查询</Button>
          <pre style={{ marginTop: 16, background: '#f5f5f5', padding: 16 }}>{promResult || '结果将显示在这里'}</pre>
        </Card>
      ),
    },
    {
      key: 'elk',
      label: 'ELK',
      children: (
        <Card title="日志搜索">
          <Input.TextArea
            value={elkQuery}
            onChange={(e) => setElkQuery(e.target.value)}
            placeholder="输入搜索条件"
            rows={3}
          />
          <Button type="primary" onClick={handleElkSearch} style={{ marginTop: 8 }}>搜索</Button>
          <pre style={{ marginTop: 16, background: '#f5f5f5', padding: 16 }}>{elkResult || '结果将显示在这里'}</pre>
        </Card>
      ),
    },
    {
      key: 'email',
      label: '邮件测试',
      children: (
        <Card title="发送测试邮件">
          <Input placeholder="收件人" value={emailTo} onChange={(e) => setEmailTo(e.target.value)} style={{ marginBottom: 8 }} />
          <Input placeholder="主题" value={emailSubject} onChange={(e) => setEmailSubject(e.target.value)} style={{ marginBottom: 8 }} />
          <Input.TextArea placeholder="内容" value={emailBody} onChange={(e) => setEmailBody(e.target.value)} rows={4} style={{ marginBottom: 8 }} />
          <Button type="primary" onClick={handleSendEmail}>发送</Button>
        </Card>
      ),
    },
  ]

  return (
    <Card title="监控集成">
      <Tabs items={items} />
    </Card>
  )
}
```

- [ ] **Step 2: router/index.tsx 添加 /monitor 路由**

```tsx
import MonitorDashboard from '@/modules/integration/MonitorDashboard'

// 在 Route path="/" element={<AppLayout />} 的子路由中添加：
<Route path="monitor" element={<MonitorDashboard />} />
```

- [ ] **Step 3: Commit**

```bash
git add cmdb-ui/src/modules/integration/MonitorDashboard.tsx cmdb-ui/src/router/index.tsx
git commit -m "feat: add MonitorDashboard page and route"
```

---

## Task 12: 前端 - SubnetList 添加新建子网弹窗

**Files:**
- Modify: `cmdb-ui/src/modules/ipam/SubnetList.tsx`

- [ ] **Step 1: 修改 SubnetList.tsx**

```tsx
import { useEffect, useState } from 'react'
import { Table, Button, Card, message, Modal, Form, Input, Select } from 'antd'
import { ipamApi } from '@/api/ipam'

export default function SubnetList() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await ipamApi.listSubnets()
      if (res.code === 0) setData(res.data)
    } catch {
      message.error('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const handleCreate = async () => {
    try {
      const values = await form.validateFields()
      const res: any = await ipamApi.createSubnet(values)
      if (res.code === 0) {
        message.success('创建成功')
        setModalOpen(false)
        form.resetFields()
        fetchData()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('创建失败')
    }
  }

  return (
    <Card
      title="子网管理"
      extra={<Button type="primary" onClick={() => setModalOpen(true)}>新建子网</Button>}
    >
      <Table dataSource={data} columns={[
        { title: '名称', dataIndex: 'name' },
        { title: 'CIDR', dataIndex: 'cidr' },
        { title: '状态', dataIndex: 'status' },
      ]} loading={loading} rowKey="id" />

      <Modal
        title="新建子网"
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields() }}
        onOk={handleCreate}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入子网名称' }]}>
            <Input placeholder="如：办公网" />
          </Form.Item>
          <Form.Item name="cidr" label="CIDR" rules={[{ required: true, message: '请输入 CIDR' }]}>
            <Input placeholder="如：192.168.1.0/24" />
          </Form.Item>
          <Form.Item name="vlan_id" label="VLAN ID">
            <Input placeholder="可选" />
          </Form.Item>
          <Form.Item name="parent_id" label="父网">
            <Select
              placeholder="可选"
              allowClear
              options={data.map((s: any) => ({ label: `${s.name} (${s.cidr})`, value: s.id }))}
            />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add cmdb-ui/src/modules/ipam/SubnetList.tsx
git commit -m "feat: add create subnet modal to SubnetList"
```

---

## Task 13: 验证所有修复

- [ ] **Step 1: 确认 UI 编译通过**

```bash
cd /e/Projects/cmdb/cmdb-ui && npm run build 2>&1 | tail -5
# 预期：无 TypeScript 编译错误
```

- [ ] **Step 2: 前端功能联调验证清单**

| # | 验证项 | 前端操作 | 预期结果 |
|---|--------|---------|---------|
| 1 | Dashboard 真实数据 | 登录后观察 Dashboard | 统计卡片显示 0，图表显示 Empty 占位 |
| 2 | CIType 属性预加载 | 创建 CIType 带属性 → 资源管理 → 新建 CI → 选择该 CIType | 表单字段与 CIType 定义的属性一致 |
| 3 | 登出 | 点击 Header 登出按钮 | Token 清除，跳转 /login |
| 4 | 监控集成路由 | 点击左侧"监控集成" | 页面显示 Tabs（Prometheus/ELK/邮件） |
| 5 | 新建子网 | IPAM 页面 → 新建子网 → 填写信息 → 确定 | 列表刷新，显示新子网 |

---

## 自检清单

- [x] **Spec coverage**: 5 个问题均有对应 Task
- [x] **Placeholder scan**: 无 TBD/TODO/模糊描述
- [x] **Type consistency**: 前后端字段名一致（ci_by_type/ci_by_status/attributes）
- [x] **File paths**: 使用绝对路径，与代码库一致
