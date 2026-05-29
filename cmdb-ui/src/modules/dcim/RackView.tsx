import { useState } from 'react'
import { Card, Tag, message } from 'antd'
import { dcimApi } from '@/api/dcim'

const statusColors: Record<string, string> = {
  active: '#52c41a',
  maintenance: '#fa8c16',
  offline: '#ff4d4f',
}

interface RackDevice {
  id: number
  u_position: number
  device_ci_id: number
  device_name: string
  status: string
}

// Mock rack data for demo
const mockDevices: RackDevice[] = [
  { id: 1, u_position: 1, device_ci_id: 100, device_name: 'Server-A', status: 'active' },
  { id: 2, u_position: 4, device_ci_id: 101, device_name: 'Server-B', status: 'active' },
  { id: 3, u_position: 10, device_ci_id: 102, device_name: 'Switch-01', status: 'maintenance' },
  { id: 4, u_position: 20, device_ci_id: 103, device_name: 'Server-C', status: 'offline' },
  { id: 5, u_position: 40, device_ci_id: 104, device_name: 'Server-D', status: 'active' },
]

export default function RackView() {
  const [rackName] = useState('RACK-A01')
  const totalU = 42
  const uHeight = 20 // px per U
  const gap = 2
  const canvasHeight = totalU * (uHeight + gap) + 40

  const occupiedMap = new Map<number, RackDevice>()
  mockDevices.forEach((d) => occupiedMap.set(d.u_position, d))

  const handleDeviceClick = (device: RackDevice) => {
    message.info(`设备: ${device.device_name} (CI ID: ${device.device_ci_id})`)
  }

  return (
    <Card title={`机柜视图: ${rackName}`} extra={<Tag color="blue">{totalU}U</Tag>}>
      <div style={{ display: 'flex', gap: 16 }}>
        {/* U position labels */}
        <div style={{ display: 'flex', flexDirection: 'column-reverse', marginTop: uHeight / 2 }}>
          {Array.from({ length: totalU }, (_, i) => (
            <div
              key={i}
              style={{
                height: uHeight + gap,
                lineHeight: `${uHeight + gap}px`,
                fontSize: 11,
                color: '#999',
                textAlign: 'right',
                paddingRight: 8,
              }}
            >
              U{i + 1}
            </div>
          ))}
        </div>

        {/* SVG rack */}
        <svg width="200" height={canvasHeight} style={{ border: '1px solid #d9d9d9', borderRadius: 4 }}>
          {Array.from({ length: totalU }, (_, i) => {
            const uPos = totalU - i
            const y = i * (uHeight + gap) + 20
            const device = occupiedMap.get(uPos)
            const fill = device ? (statusColors[device.status] || '#d9d9d9') : '#f5f5f5'
            return (
              <g key={i}>
                <rect
                  x="10"
                  y={y}
                  width="180"
                  height={uHeight}
                  rx="2"
                  fill={fill}
                  stroke="#d9d9d9"
                  onClick={() => device && handleDeviceClick(device)}
                  style={{ cursor: device ? 'pointer' : 'default' }}
                />
                {device && (
                  <text
                    x="100"
                    y={y + uHeight / 2 + 4}
                    textAnchor="middle"
                    fontSize="11"
                    fill="#fff"
                    style={{ pointerEvents: 'none' }}
                  >
                    {device.device_name}
                  </text>
                )}
              </g>
            )
          })}
        </svg>
      </div>
    </Card>
  )
}
