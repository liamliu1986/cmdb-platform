import { useState } from 'react'
import { Card, Space, message } from 'antd'
import RelationGraph from '@/components/graph/RelationGraph'

export default function CIRelationGraph() {
  const [ciId, setCiId] = useState(1)
  const [ciTypeName, setCITypeName] = useState('Server')

  const handleNodeClick = (id: number, type: string) => {
    setCiId(id)
    setCITypeName(type)
    message.info(`切换到: ${type} (ID: ${id})`)
  }

  return (
    <Card title="关系图谱" extra={
      <Space>
        <span>当前中心节点: {ciTypeName}</span>
      </Space>
    }>
      <RelationGraph ciId={ciId} ciTypeName={ciTypeName} onNodeClick={handleNodeClick} />
    </Card>
  )
}
