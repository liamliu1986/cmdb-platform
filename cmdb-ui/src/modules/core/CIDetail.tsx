import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Card, Descriptions, Tabs, Spin, message, Alert, Table } from 'antd'
import { coreApi } from '@/api/core'
import CIMonitoringTab from './CIMonitoringTab'
import CILogTab from './CILogTab'

export default function CIDetail() {
  const { id } = useParams<{ id: string }>()
  const [ci, setCI] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!id) return
    setLoading(true)
    coreApi.getCI(parseInt(id))
      .then((res: any) => {
        if (res.code === 0) setCI(res.data)
        else message.error(res.message)
      })
      .catch(() => message.error('加载CI详情失败'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) return <Spin size="large" style={{ display: 'block', margin: '100px auto' }} />
  if (!ci) return <Alert type="error" message="CI not found" />

  const attrEntries = ci.attr_values
    ? Object.entries(ci.attr_values)
    : []

  return (
    <div>
      <Card title={`${ci.ci_type || 'CI'} 详情`} style={{ marginBottom: 16 }}>
        <Descriptions bordered column={2} size="small">
          <Descriptions.Item label="CI ID">{ci.id}</Descriptions.Item>
          <Descriptions.Item label="类型">{ci.ci_type}</Descriptions.Item>
          <Descriptions.Item label="状态">{ci.status || 'active'}</Descriptions.Item>
          <Descriptions.Item label="更新人">{ci.updated_by || '-'}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{ci.updated_at || '-'}</Descriptions.Item>
        </Descriptions>
        {attrEntries.length > 0 && (
          <div style={{ marginTop: 16 }}>
            <h4>属性值</h4>
            <Table
              dataSource={attrEntries.map(([k, v]) => ({ key: k, value: String(v) }))}
              columns={[
                { title: '属性', dataIndex: 'key', width: 200 },
                { title: '值', dataIndex: 'value' },
              ]}
              pagination={false}
              size="small"
            />
          </div>
        )}
      </Card>

      <Card>
        <Tabs items={[
          {
            key: 'monitoring',
            label: '监控指标',
            children: <CIMonitoringTab ci={ci} />,
          },
          {
            key: 'logs',
            label: '日志查询',
            children: <CILogTab ci={ci} />,
          },
        ]} />
      </Card>
    </div>
  )
}
