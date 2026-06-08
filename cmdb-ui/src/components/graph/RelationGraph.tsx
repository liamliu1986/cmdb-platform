import { useEffect, useRef } from 'react'
import { Graph } from '@antv/g6'

interface Props {
  ciId: number
  ciTypeName: string
  onNodeClick?: (ciId: number, ciType: string) => void
}

const mockData = {
  nodes: [
    {
      id: 'center',
      data: { label: 'Web Server 01', ciType: 'Server' },
      style: { fill: '#1890ff' },
    },
    {
      id: 'db1',
      data: { label: 'MySQL DB 01', ciType: 'MySQL' },
      style: { fill: '#52c41a' },
    },
    {
      id: 'redis1',
      data: { label: 'Redis Cache 01', ciType: 'Redis' },
      style: { fill: '#fa8c16' },
    },
    {
      id: 'lb1',
      data: { label: 'Load Balancer', ciType: 'SLB' },
      style: { fill: '#722ed1' },
    },
  ],
  edges: [
    { source: 'center', target: 'db1', data: { label: 'depends_on' } },
    { source: 'center', target: 'redis1', data: { label: 'depends_on' } },
    { source: 'lb1', target: 'center', data: { label: 'routes_to' } },
  ],
}

export default function RelationGraph({ ciId: _ciId, ciTypeName: _ciTypeName, onNodeClick }: Props) {
  const containerRef = useRef<HTMLDivElement>(null)
  const graphRef = useRef<Graph | null>(null)

  useEffect(() => {
    if (!containerRef.current) return

    const container = containerRef.current
    const width = container.scrollWidth || 800
    const height = container.scrollHeight || 500

    const graph = new Graph({
      container,
      width,
      height,
      data: mockData,
      layout: {
        type: 'd3-force',
        preventOverlap: true,
        nodeStrength: -50,
      },
      behaviors: [
        { type: 'drag-canvas' },
        { type: 'zoom-canvas' },
        { type: 'drag-element' },
      ],
      node: {
        type: 'rect',
        style: {
          size: [140, 40],
          radius: 8,
          fill: '#1890ff',
          labelText: (d: { data?: { label?: string } }) => d.data?.label ?? d.id ?? '',
          labelFill: '#fff',
          labelFontSize: 12,
          labelPlacement: 'center',
        },
      },
      edge: {
        type: 'polyline',
        style: {
          stroke: '#ccc',
          endArrow: true,
          labelText: (d: { data?: { label?: string } }) => d.data?.label ?? '',
          labelFontSize: 10,
          labelBackground: true,
          labelBackgroundFill: '#fff',
        },
      },
    })

    graph.render()
    graphRef.current = graph

    graph.on('node:click', (evt: unknown) => {
      const e = evt as { target?: { id?: string; type?: string } }
      if (onNodeClick && e.target) {
        const nodeId = e.target.id
        if (nodeId && nodeId !== 'center') {
          const nodeData = graph.getNodeData(nodeId)
          const ciType = (nodeData?.data?.ciType as string) ?? 'Unknown'
          onNodeClick(0, ciType)
        }
      }
    })

    return () => {
      graph.destroy()
      graphRef.current = null
    }
  }, [_ciId, _ciTypeName, onNodeClick])

  return (
    <div style={{ width: '100%', height: '100%', minHeight: 500, background: '#fff', borderRadius: 8 }}>
      <div ref={containerRef} style={{ width: '100%', height: 500 }} />
    </div>
  )
}
