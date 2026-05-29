import { useEffect, useState } from 'react'
import { Table, Button, Card, message } from 'antd'
import { dcimApi } from '@/api/dcim'

export default function IDCList() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await dcimApi.listIDCs()
      if (res.code === 0) setData(res.data)
    } catch {
      message.error('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  return (
    <Card title="数据中心" extra={<Button type="primary">新建IDC</Button>}>
      <Table dataSource={data} columns={[
        { title: '名称', dataIndex: 'name' },
        { title: '地址', dataIndex: 'address' },
        { title: '联系人', dataIndex: 'contact' },
      ]} loading={loading} rowKey="id" />
    </Card>
  )
}
