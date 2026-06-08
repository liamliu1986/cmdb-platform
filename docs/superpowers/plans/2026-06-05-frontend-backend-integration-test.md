# CMDB 前后端联调测试计划

> 测试环境：从零数据开始，仅保留 admin / admin123 管理员账户
> 测试方式：每一步前端操作后，用 curl 验证后端数据状态

---

## 阶段一：基础框架验证

### 测试 1：登录认证

**前置条件**：数据库仅含 admin 用户，API/UI 已启动

**前端操作**：
1. 访问 http://localhost:3000
2. 观察是否自动跳转 `/login`
3. 输入错误密码 `admin` / `wrong`，点击登录
4. 输入正确密码 `admin` / `admin123`，点击登录

**验证点**：
- [ ] 错误密码提示 "invalid username or password"
- [ ] 正确密码显示 "登录成功" Toast
- [ ] 页面跳转到 `/` (Dashboard)
- [ ] F12 → Application → Local Storage → `token` 存在且有值
- [ ] 清除 token 刷新页面 → 自动跳回 `/login`

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/auth/login -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
# 预期返回 code=0, data.token 非空
```

---

### 测试 2：Dashboard 与菜单导航

**前置条件**：已登录

**前端操作**：
1. 观察 Dashboard 页面
2. 逐一点击左侧 12 个菜单项

**验证点**：
- [ ] 4 个统计卡片显示（CI 总数=0, CIType=0, 发现规则=0, Agent=0）
- [ ] 饼图/柱状图正常渲染（注意：数据为前端硬编码）
- [ ] 点击每个菜单项 → 页面切换无报错、无 404
- [ ] **监控集成** (`/monitor`) 菜单 → 观察是否为 404 或空白

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/ci/s?page=1&page_size=1 \
  -H "Authorization: Bearer <token>"
# 预期返回 code=0, data.pagination.total=0

curl -s http://localhost:8080/api/v1/citypes \
  -H "Authorization: Bearer <token>"
# 预期返回 code=0, data=[]
```

---

## 阶段二：核心功能（CIType → CI）

### 测试 3：CIType 模型管理 - 创建

**前置条件**：已登录，当前无 CIType 数据

**前端操作**：
1. 点击左侧 **模型管理**
2. 点击右上角 **"新建模型"**
3. 在弹窗中输入：
   - 模型名称：`Server`
   - 别名：`服务器`
4. 在属性定义区域添加属性：
   - 属性名：`hostname`，别名：`主机名`，类型：`text` → 点击"添加"
   - 属性名：`ip`，别名：`IP地址`，类型：`text` → 点击"添加"
   - 属性名：`os`，别名：`操作系统`，类型：`choice` → 点击"添加"
5. 点击 **"确定"**

**验证点**：
- [ ] 弹窗关闭
- [ ] 列表自动刷新，出现 "Server" 行
- [ ] 显示 "创建成功" Toast

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/citypes \
  -H "Authorization: Bearer <token>"
# 预期返回 code=0, data 包含 Server 对象，attributes 包含 3 个属性
```

**补充测试**：重复创建同名 CIType `Server`
- [ ] 预期报错：名称已存在或后端唯一性约束错误

---

### 测试 4：CIType 模型管理 - 列表与详情

**前置条件**：已创建 Server CIType

**前端操作**：
1. 刷新 **模型管理** 页面
2. 观察列表列：名称、别名、状态、创建时间
3. 点击某行查看是否有详情入口

**验证点**：
- [ ] 列表正常加载，数据与创建时一致
- [ ] 状态列显示 "启用"

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/citypes/1 \
  -H "Authorization: Bearer <token>"
# 预期返回 code=0, data 包含完整 CIType 信息
```

---

### 测试 5：CI 实例管理 - 创建

**前置条件**：已创建 Server CIType

**前端操作**：
1. 点击左侧 **资源管理**
2. 点击右上角 **"新建资源"**
3. 在弹窗中选择模型：`Server`
4. 填写表单（名称、描述等）
5. 点击 **"确定"**

**验证点**：
- [ ] 弹窗关闭
- [ ] 列表自动刷新，出现新 CI 行
- [ ] 显示 "创建成功" Toast

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/ci/s?page=1&page_size=25 \
  -H "Authorization: Bearer <token>"
# 预期返回 code=0, data.list 非空

curl -s http://localhost:8080/api/v1/ci/1 \
  -H "Authorization: Bearer <token>"
# 预期返回 code=0, data 包含 attr_values
```

**已知问题观察**：
- CIList 前端 `handleCITypeChange` 中属性为硬编码（name/description），实际应根据选择的 CIType 动态获取属性。观察创建时表单字段是否正确。

---

### 测试 6：CI 实例管理 - 搜索

**前置条件**：已创建至少 1 个 CI

**前端操作**：
1. 在 **资源管理** 页面上方搜索框输入关键字
2. 点击搜索或回车

**验证点**：
- [ ] 列表根据关键字过滤
- [ ] 无结果时显示空状态

---

### 测试 7：CI 实例管理 - 详情页

**前置条件**：已创建 CI

**前端操作**：
1. 在 **资源管理** 列表中点击某行
2. 观察是否跳转到 `/ci/:id`

**验证点**：
- [ ] 页面正常加载，显示 CI 详情
- [ ] 包含 日志/监控 Tab（CILogTab、CIMonitoringTab）

---

## 阶段三：IPAM 功能

### 测试 8：子网管理 - 创建

**前置条件**：已登录

**前端操作**：
1. 点击左侧 **IPAM**
2. 点击右上角 **"新建子网"**
3. 填写信息（CIDR 如 `192.168.1.0/24`）

**验证点**：
- [ ] 列表出现新子网

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/ipam/subnets \
  -H "Authorization: Bearer <token>"
# 预期返回 code=0, data 包含子网
```

---

### 测试 9：IP 分配与释放

**前置条件**：已创建子网

**前端操作**：
1. 在子网列表中找到操作按钮
2. 点击 **"分配 IP"**
3. 观察分配的 IP 地址
4. 点击 **"释放"**

**验证点**：
- [ ] 分配成功返回 IP 地址
- [ ] 释放后状态更新

**后端验证**：
```bash
curl -s -X POST http://localhost:8080/api/v1/ipam/ips/allocate \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"subnet_id":1}'

curl -s -X POST http://localhost:8080/api/v1/ipam/ips/1/release \
  -H "Authorization: Bearer <token>"
```

---

## 阶段四：DCIM 功能

### 测试 10：IDC 管理 - 创建与列表

**前端操作**：
1. 点击左侧 **DCIM**
2. 创建 IDC：`北京数据中心`

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/dcim/idcs \
  -H "Authorization: Bearer <token>"
```

---

### 测试 11：机房/机柜管理

**前端操作**：
1. 在 IDC 下创建机房
2. 在机房下创建机柜

---

### 测试 12：机柜视图

**前端操作**：
1. 点击左侧 **机柜视图**
2. 观察机柜布局渲染

---

### 测试 13：IDC 地图

**前端操作**：
1. 点击左侧 **IDC地图**
2. 观察地图组件渲染

---

## 阶段五：发现与 Agent

### 测试 14：发现规则管理

**前端操作**：
1. 点击左侧 **自动发现**
2. 创建发现规则（名称、类型、配置）
3. 点击 **"执行规则"**

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/discovery/rules \
  -H "Authorization: Bearer <token>"
```

---

### 测试 15：Agent 管理

**前端操作**：
1. 点击左侧 **Agent管理**
2. 观察 Agent 列表

**后端验证**：
```bash
curl -s http://localhost:8080/api/v1/discovery/agents \
  -H "Authorization: Bearer <token>"
# 预期空列表（无 Agent 注册）
```

---

## 阶段六：可视化

### 测试 16：关系图谱

**前置条件**：已创建多个 CI 且有关系数据

**前端操作**：
1. 点击左侧 **关系图谱**
2. 观察 G6 图渲染

---

### 测试 17：层级拓扑

**前端操作**：
1. 点击左侧 **层级拓扑**
2. 观察拓扑树渲染

---

## 阶段七：监控集成

### 测试 18：Prometheus 查询

**前端操作**：
1. 点击左侧 **监控集成**
2. 输入 PromQL 查询

**注意**：需要配置 Prometheus 地址才能真实调用

---

## 问题清单（测试前已知）

| # | 问题 | 位置 | 严重程度 |
|---|------|------|---------|
| 1 | Dashboard 图表数据硬编码 | Dashboard.tsx:45-71 | 中 |
| 2 | CIList 属性硬编码（不根据 CIType 动态获取） | CIList.tsx:44-47 | 高 |
| 3 | 无登出功能 | AppLayout.tsx | 中 |
| 4 | 监控集成菜单无路由 | AppLayout.tsx 有菜单，router 无匹配 | 低 |
| 5 | SubnetList "新建子网" 按钮无 onClick 绑定 | SubnetList.tsx:24 | 高 |

---

## 执行顺序

```
测试 1 → 测试 2 → 测试 3 → 测试 4 → 测试 5 → 测试 6 → 测试 7
  ↓
测试 8 → 测试 9
  ↓
测试 10 → 测试 11 → 测试 12 → 测试 13
  ↓
测试 14 → 测试 15
  ↓
测试 16 → 测试 17
  ↓
测试 18
```

每完成一项测试，用户在前端验证后回复结果，再进行下一项。
