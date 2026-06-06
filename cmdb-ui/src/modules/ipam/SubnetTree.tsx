import { useEffect, useState } from 'react'
import { Card, Tree, Spin, Empty, Button, Space, message } from 'antd'
import { ClusterOutlined, ArrowLeftOutlined } from '@ant-design/icons'
import { ipamApi } from '@/api/ipam'
import { useNavigate } from 'react-router-dom'

interface SubnetNode {
  id: number
  parent_id?: number
  name: string
  cidr: string
  children?: SubnetNode[]
}

export default function SubnetTree() {
  const [treeData, setTreeData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [expandedKeys, setExpandedKeys] = useState<string[]>([])
  const navigate = useNavigate()

  const buildTree = (subnets: SubnetNode[]): any[] => {
    const nodeMap = new Map<number, SubnetNode & { children?: SubnetNode[] }>()
    subnets.forEach((s) => nodeMap.set(s.id, { ...s, children: [] }))

    const roots: any[] = []
    subnets.forEach((s) => {
      const node = nodeMap.get(s.id)!
      if (s.parent_id && nodeMap.has(s.parent_id)) {
        const parent = nodeMap.get(s.parent_id)!
        if (!parent.children) parent.children = []
        parent.children.push(node)
      } else {
        roots.push(node)
      }
    })

    const toTreeNode = (n: SubnetNode & { children?: SubnetNode[] }): any => ({
      key: String(n.id),
      title: (
        <Space>
          <ClusterOutlined />
          <span>{n.name}</span>
          <span style={{ color: '#888', fontSize: 12 }}>{n.cidr}</span>
        </Space>
      ),
      children: n.children?.map(toTreeNode) || [],
    })

    return roots.map(toTreeNode)
  }

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await ipamApi.listSubnets()
      if (res.code === 0) {
        const data = buildTree(res.data || [])
        setTreeData(data)
        if (data.length > 0) {
          setExpandedKeys(data.map((d) => d.key))
        }
      } else {
        message.error(res.message || '加载失败')
      }
    } catch {
      message.error('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const handleSelect = (selectedKeys: any[]) => {
    if (selectedKeys.length > 0) {
      navigate(`/ipam/ips?subnet_id=${selectedKeys[0]}`)
    }
  }

  return (
    <Card
      title="子网拓扑"
      extra={
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/ipam')}>
          返回列表
        </Button>
      }
    >
      <Spin spinning={loading}>
        {treeData.length === 0 ? (
          <Empty description="暂无子网数据" />
        ) : (
          <Tree
            treeData={treeData}
            expandedKeys={expandedKeys}
            onExpand={(keys) => setExpandedKeys(keys as string[])}
            onSelect={handleSelect}
            showLine
            defaultExpandAll
          />
        )}
      </Spin>
    </Card>
  )
}
