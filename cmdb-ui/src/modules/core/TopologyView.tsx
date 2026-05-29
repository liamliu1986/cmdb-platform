import { Card } from 'antd'
import TopologyTree from '@/components/graph/TopologyTree'

export default function TopologyView() {
  return (
    <Card title="层级拓扑">
      <TopologyTree />
    </Card>
  )
}
