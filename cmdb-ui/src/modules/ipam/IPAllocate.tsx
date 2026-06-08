import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Card, Table, Button, Tag, message, Space, Select } from 'antd'
import { ArrowLeftOutlined, ThunderboltOutlined } from '@ant-design/icons'
import { ipamApi } from '@/api/ipam'

export default function IPAllocate() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const initialSubnetId = Number(searchParams.get('subnet_id') || 0)

  const [subnets, setSubnets] = useState<any[]>([])
  const [ips, setIPs] = useState<any[]>([])
  const [selectedSubnet, setSelectedSubnet] = useState<number | undefined>(
    initialSubnetId || undefined
  )
  const [loading, setLoading] = useState(false)
  const [allocating, setAllocating] = useState<number | null>(null)

  const fetchSubnets = async () => {
    try {
      const res: any = await ipamApi.listSubnets()
      if (res.code === 0) setSubnets(res.data || [])
    } catch {
      message.error('加载子网失败')
    }
  }

  const fetchIPs = async (subnetId?: number) => {
    if (!subnetId) {
      setIPs([])
      return
    }
    setLoading(true)
    try {
      const res: any = await ipamApi.listIPsBySubnet(subnetId, { status: 'free' })
      if (res.code === 0) setIPs(res.data || [])
    } catch {
      message.error('加载 IP 失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSubnets()
    if (selectedSubnet) fetchIPs(selectedSubnet)
  }, [])

  useEffect(() => {
    fetchIPs(selectedSubnet)
  }, [selectedSubnet])

  const handleAllocateOne = async (ip: any) => {
    setAllocating(ip.id)
    try {
      const res: any = await ipamApi.allocateIPByID(ip.id)
      if (res.code === 0) {
        message.success(`IP ${ip.ip} 分配成功`)
        fetchIPs(selectedSubnet)
      } else {
        message.error(res.message || '分配失败')
      }
    } catch {
      message.error('分配失败')
    } finally {
      setAllocating(null)
    }
  }

  const handleAutoAllocate = async () => {
    if (!selectedSubnet) return
    setAllocating(-1)
    try {
      const res: any = await ipamApi.allocateIP({ subnet_id: selectedSubnet })
      if (res.code === 0) {
        message.success(`成功分配 IP: ${res.data?.ip}`)
        fetchIPs(selectedSubnet)
      } else {
        message.error(res.message || '自动分配失败')
      }
    } catch {
      message.error('自动分配失败')
    } finally {
      setAllocating(null)
    }
  }

  const columns = [
    { title: 'IP 地址', dataIndex: 'ip', key: 'ip' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (_status: string) => <Tag color="green">空闲</Tag>,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Button
          type="primary"
          size="small"
          loading={allocating === record.id}
          icon={<ThunderboltOutlined />}
          onClick={() => handleAllocateOne(record)}
        >
          分配此 IP
        </Button>
      ),
    },
  ]

  const subnetOptions = subnets.map((s: any) => ({
    label: `${s.name} (${s.cidr})`,
    value: s.id,
  }))

  return (
    <Card
      title="IP 分配"
      extra={
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/ipam')}>
          返回
        </Button>
      }
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <Select
          placeholder="选择子网"
          style={{ width: 320 }}
          options={subnetOptions}
          value={selectedSubnet}
          onChange={(v) => setSelectedSubnet(v)}
          allowClear
        />
        {selectedSubnet && (
          <Button
            type="primary"
            icon={<ThunderboltOutlined />}
            loading={allocating === -1}
            onClick={handleAutoAllocate}
          >
            自动分配
          </Button>
        )}
      </Space>

      {!selectedSubnet ? (
        <p>请选择一个子网查看可用 IP。</p>
      ) : (
        <Table
          dataSource={ips}
          columns={columns}
          loading={loading}
          rowKey="id"
          pagination={{ pageSize: 20 }}
          summary={() => (
            <Table.Summary.Row>
              <Table.Summary.Cell index={0} colSpan={3}>
                当前子网空闲 IP 数量：<strong>{ips.length}</strong>
              </Table.Summary.Cell>
            </Table.Summary.Row>
          )}
        />
      )}
    </Card>
  )
}
