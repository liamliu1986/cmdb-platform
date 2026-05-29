import { useEffect, useRef } from 'react'
import { Graph } from '@antv/g6'

// Demo data: hierarchical organization structure
const treeData = {
  id: 'group',
  value: { label: 'XX集团', type: 'Group' },
  children: [
    {
      id: 'region-east',
      value: { label: '华东地区', type: 'Region' },
      children: [
        {
          id: 'idc-sh',
          value: { label: '上海IDC', type: 'IDC' },
          children: [
            {
              id: 'room-a1',
              value: { label: 'A1机房', type: 'ServerRoom' },
              children: [
                { id: 'rack-01', value: { label: 'RACK-A01', type: 'Rack' } },
                { id: 'rack-02', value: { label: 'RACK-A02', type: 'Rack' } },
              ],
            },
            {
              id: 'room-b1',
              value: { label: 'B1机房', type: 'ServerRoom' },
              children: [
                { id: 'rack-03', value: { label: 'RACK-B01', type: 'Rack' } },
              ],
            },
          ],
        },
      ],
    },
    {
      id: 'region-north',
      value: { label: '华北地区', type: 'Region' },
      children: [
        {
          id: 'idc-bj',
          value: { label: '北京IDC', type: 'IDC' },
          children: [
            {
              id: 'room-c1',
              value: { label: 'C1机房', type: 'ServerRoom' },
              children: [
                { id: 'rack-04', value: { label: 'RACK-C01', type: 'Rack' } },
                { id: 'rack-05', value: { label: 'RACK-C02', type: 'Rack' } },
              ],
            },
          ],
        },
      ],
    },
  ],
}

export default function TopologyTree() {
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!containerRef.current) return

    const graph = new Graph({
      container: containerRef.current,
      width: 800,
      height: 500,
      data: treeData,
      layout: {
        type: 'compact-box',
        direction: 'LR',
        getWidth: () => 80,
        getHeight: () => 36,
        getVGap: () => 10,
        getHGap: () => 50,
      },
      behaviors: ['drag-canvas', 'zoom-canvas', 'collapse-expand'],
      node: {
        type: 'rect',
        style: {
          size: [80, 36],
          radius: 4,
          fill: '#f0f5ff',
          stroke: '#1890ff',
          labelText: (d: any) => d.value?.label || '',
          labelFontSize: 11,
          labelFill: '#333',
        },
        state: {
          collapsed: {
            size: [80, 36],
          },
        },
      },
      edge: {
        type: 'polyline',
        style: {
          stroke: '#aaa',
          lineWidth: 1,
        },
      },
    })

    graph.render()

    return () => {
      graph.destroy()
    }
  }, [])

  return (
    <div ref={containerRef} style={{ width: '100%', height: 500, background: '#fff', borderRadius: 8 }} />
  )
}
