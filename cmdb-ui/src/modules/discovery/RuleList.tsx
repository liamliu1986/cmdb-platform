import { useEffect, useState } from 'react'
import { Table, Button, Card, message } from 'antd'
import { discoveryApi } from '@/api/discovery'

export default function RuleList() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await discoveryApi.listRules()
      if (res.code === 0) setData(res.data)
    } catch {
      message.error('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  return (
    <Card title="发现规则" extra={<Button type="primary">新建规则</Button>}>
      <Table dataSource={data} columns={[
        { title: '名称', dataIndex: 'name' },
        { title: '类型', dataIndex: 'type' },
        { title: '调度', dataIndex: 'schedule' },
        { title: '状态', dataIndex: 'enabled', render: (v: boolean) => v ? '启用' : '禁用' },
      ]} loading={loading} rowKey="id" />
    </Card>
  )
}
