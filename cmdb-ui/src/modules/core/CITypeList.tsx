import { useEffect, useState } from 'react'
import { Table, Button, Card, message } from 'antd'
import { coreApi } from '@/api/core'
import CITypeDesigner from './CITypeDesigner'

export default function CITypeList() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await coreApi.listCITypes()
      if (res.code === 0) setData(res.data)
    } catch {
      message.error('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '别名', dataIndex: 'alias', key: 'alias' },
    { title: '状态', dataIndex: 'enabled', key: 'enabled', render: (v: boolean) => (v ? '启用' : '禁用') },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  ]

  return (
    <Card
      title="CIType 管理"
      extra={<Button type="primary" onClick={() => setModalOpen(true)}>新建模型</Button>}
    >
      <Table dataSource={data} columns={columns} loading={loading} rowKey="id" />
      <CITypeDesigner
        open={modalOpen}
        onClose={() => { setModalOpen(false); fetchData() }}
      />
    </Card>
  )
}
