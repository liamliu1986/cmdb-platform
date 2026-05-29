import { useEffect, useState } from 'react'
import { Table, Button, Card, message } from 'antd'
import { ipamApi } from '@/api/ipam'

export default function SubnetList() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

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

  return (
    <Card title="子网管理" extra={<Button type="primary">新建子网</Button>}>
      <Table dataSource={data} columns={[
        { title: '名称', dataIndex: 'name' },
        { title: 'CIDR', dataIndex: 'cidr' },
        { title: '状态', dataIndex: 'status' },
      ]} loading={loading} rowKey="id" />
    </Card>
  )
}
