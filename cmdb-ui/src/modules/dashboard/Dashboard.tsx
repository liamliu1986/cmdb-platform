import { useEffect, useState } from 'react'
import { Card, Col, Row, Statistic } from 'antd'
import { DatabaseOutlined, CloudServerOutlined, WifiOutlined, TeamOutlined } from '@ant-design/icons'
import ReactECharts from 'echarts-for-react'
import { coreApi } from '@/api/core'
import { discoveryApi } from '@/api/discovery'

export default function Dashboard() {
  const [ciCount, setCiCount] = useState(0)
  const [citypeCount, setCITypeCount] = useState(0)
  const [ruleCount, setRuleCount] = useState(0)

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const ciRes: any = await coreApi.searchCI({ page: 1, page_size: 1 })
        if (ciRes.code === 0) setCiCount(ciRes.data?.pagination?.total || 0)
      } catch {
        // API not yet available, keep default 0
      }
      try {
        const ctRes: any = await coreApi.listCITypes()
        if (ctRes.code === 0) setCITypeCount(ctRes.data?.length || 0)
      } catch {
        // API not yet available, keep default 0
      }
      try {
        const ruleRes: any = await discoveryApi.listRules()
        if (ruleRes.code === 0) setRuleCount(ruleRes.data?.length || 0)
      } catch {
        // API not yet available, keep default 0
      }
    }
    fetchStats()
  }, [])

  // Device type distribution pie chart
  const pieOption = {
    title: { text: 'CI 类型分布', left: 'center' },
    tooltip: { trigger: 'item' },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        data: [
          { value: 120, name: 'Server' },
          { value: 45, name: 'VM' },
          { value: 30, name: 'Switch' },
          { value: 25, name: 'MySQL' },
          { value: 15, name: 'Redis' },
          { value: 10, name: 'Other' },
        ],
      },
    ],
  }

  // Status distribution bar chart
  const barOption = {
    title: { text: '设备状态分布', left: 'center' },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: ['Online', 'Offline', 'Maintenance', 'Pending'] },
    yAxis: { type: 'value' },
    series: [
      {
        type: 'bar',
        data: [
          { value: 180, itemStyle: { color: '#52c41a' } },
          { value: 25, itemStyle: { color: '#ff4d4f' } },
          { value: 15, itemStyle: { color: '#fa8c16' } },
          { value: 10, itemStyle: { color: '#d9d9d9' } },
        ],
      },
    ],
  }

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col span={6}>
          <Card>
            <Statistic title="CI 总数" value={ciCount} prefix={<DatabaseOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="CIType" value={citypeCount} prefix={<CloudServerOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="发现规则" value={ruleCount} prefix={<WifiOutlined />} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="Agent 在线" value={0} prefix={<TeamOutlined />} />
          </Card>
        </Col>
      </Row>
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Card>
            <ReactECharts option={pieOption} style={{ height: 350 }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card>
            <ReactECharts option={barOption} style={{ height: 350 }} />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
