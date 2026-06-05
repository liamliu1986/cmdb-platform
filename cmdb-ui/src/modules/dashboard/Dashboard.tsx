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
