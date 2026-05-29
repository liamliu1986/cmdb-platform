import { useEffect, useState } from 'react'
import { Table, Card, message } from 'antd'
import { discoveryApi } from '@/api/discovery'

export default function AgentList() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await discoveryApi.listAgents()
      if (res.code === 0) setData(res.data)
    } catch {
      message.error('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  return (
    <Card title="Agent 管理">
      <Table dataSource={data} columns={[
        { title: '名称', dataIndex: 'name' },
        { title: 'IP', dataIndex: 'ip' },
        { title: '系统', dataIndex: 'os' },
        { title: '架构', dataIndex: 'arch' },
        { title: '状态', dataIndex: 'status' },
        { title: '最后心跳', dataIndex: 'last_heartbeat' },
      ]} loading={loading} rowKey="id" />
    </Card>
  )
}
