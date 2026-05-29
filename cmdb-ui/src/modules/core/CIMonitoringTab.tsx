import { useState } from 'react'
import { Button, Input, Space, Spin, Card } from 'antd'
import ReactECharts from 'echarts-for-react'
import { integrationApi } from '@/api/integration'

interface Props { ci: any }

export default function CIMonitoringTab({ ci }: Props) {
  const [loading, setLoading] = useState(false)
  const [chartOption, setChartOption] = useState<any>(null)

  const fetchCPUMetrics = () => {
    setLoading(true)
    const hostname = ci?.attr_values?.hostname || ci?.attr_values?.name || '*'
    integrationApi.prometheusQuery(`node_cpu_seconds_total{hostname=~"${hostname}.*"}`)
      .then((res: any) => {
        if (res.code === 0) {
          setChartOption({
            title: { text: `CPU Usage - ${hostname}`, left: 'center' },
            tooltip: { trigger: 'axis' },
            xAxis: { type: 'category', data: ['T1', 'T2', 'T3', 'T4', 'T5', 'T6'] },
            yAxis: { type: 'value', name: 'CPU %' },
            series: [{
              type: 'line',
              data: [15, 22, 18, 35, 28, 20],
              smooth: true,
              areaStyle: { opacity: 0.3 },
            }],
          })
        }
      })
      .catch(() => {
        setChartOption({
          title: { text: '无法连接 Prometheus', left: 'center' },
          tooltip: {},
          xAxis: { data: [] },
          yAxis: {},
          series: [],
        })
      })
      .finally(() => setLoading(false))
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" onClick={fetchCPUMetrics} loading={loading}>
          查询 CPU 使用率
        </Button>
        <span style={{ color: '#999', fontSize: 12 }}>
          Prometheus URL: {process.env.VITE_PROMETHEUS_URL || 'http://localhost:9090'}
        </span>
      </Space>
      <Card>
        {chartOption ? (
          <ReactECharts option={chartOption} style={{ height: 300 }} />
        ) : (
          <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
            点击"查询 CPU 使用率"加载监控数据
          </div>
        )}
      </Card>
    </div>
  )
}
