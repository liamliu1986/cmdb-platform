import { useEffect, useState } from 'react'
import { Card, Table, message } from 'antd'
import { dcimApi } from '@/api/dcim'

// Mock IDC map display using a simple visual representation
// In production, replace with Gaode Map JS API (AMap)
export default function IDCMap() {
  const [idcs, setIDCs] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedIDC, setSelectedIDC] = useState<any>(null)

  const fetchIDCs = async () => {
    setLoading(true)
    try {
      const res: any = await dcimApi.listIDCs()
      if (res.code === 0) setIDCs(res.data)
    } catch {
      message.error('加载IDC列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchIDCs() }, [])

  return (
    <Card title="IDC 地理分布">
      <div style={{ marginBottom: 16 }}>
        <p style={{ color: '#666' }}>
          数据来源：后端 DCIM 模块。在生产环境中，此页面将集成高德地图 JS API，
          在地图上以 Marker 形式展示各 IDC 的地理位置，支持聚合显示和点击查看详情。
        </p>
      </div>
      <Table
        dataSource={idcs}
        loading={loading}
        rowKey="id"
        columns={[
          { title: '名称', dataIndex: 'name' },
          { title: '地址', dataIndex: 'address' },
          { title: '联系人', dataIndex: 'contact' },
        ]}
        onRow={(record) => ({
          onClick: () => setSelectedIDC(record),
          style: {
            background: selectedIDC?.id === record.id ? '#e6f7ff' : undefined,
            cursor: 'pointer',
          },
        })}
      />
      {selectedIDC && (
        <div style={{ marginTop: 16, padding: 16, background: '#fafafa', borderRadius: 8 }}>
          <h3>{selectedIDC.name} 详情</h3>
          <p>地址: {selectedIDC.address || '-'}</p>
          <p>联系人: {selectedIDC.contact || '-'}</p>
        </div>
      )}
    </Card>
  )
}
