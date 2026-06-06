import { useEffect, useState } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import { Card, Table, Button, Tag, message, Space, Select, Modal, Form, Input } from 'antd'
import { ArrowLeftOutlined, PlusOutlined } from '@ant-design/icons'
import { ipamApi } from '@/api/ipam'

const statusMap: Record<string, { label: string; color: string }> = {
  free: { label: '空闲', color: 'green' },
  allocated: { label: '已分配', color: 'blue' },
  reserved: { label: '预留', color: 'orange' },
  disabled: { label: '禁用', color: 'red' },
}

export default function IPList() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const subnetId = Number(searchParams.get('subnet_id') || 0)
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [subnetName, setSubnetName] = useState('')
  const [modalOpen, setModalOpen] = useState(false)
  const [allocating, setAllocating] = useState(false)
  const [form] = Form.useForm()

  const fetchData = async () => {
    if (!subnetId) return
    setLoading(true)
    try {
      const params: { status?: string } = {}
      if (statusFilter) params.status = statusFilter
      const res: any = await ipamApi.listIPsBySubnet(subnetId, params)
      if (res.code === 0) {
        setData(res.data || [])
        if (res.data?.[0]?.subnet) {
          setSubnetName(res.data[0].subnet.name)
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

  const fetchSubnetInfo = async () => {
    try {
      const res: any = await ipamApi.listSubnets()
      if (res.code === 0) {
        const subnet = (res.data || []).find((s: any) => s.id === subnetId)
        if (subnet) setSubnetName(subnet.name)
      }
    } catch {
      // ignore
    }
  }

  useEffect(() => {
    if (subnetId) {
      fetchData()
      fetchSubnetInfo()
    }
  }, [subnetId, statusFilter])

  const handleRelease = async (id: number) => {
    try {
      const res: any = await ipamApi.releaseIP(id)
      if (res.code === 0) {
        message.success('释放成功')
        fetchData()
      } else {
        message.error(res.message || '释放失败')
      }
    } catch {
      message.error('释放失败')
    }
  }

  const handleAllocate = async (record: any) => {
    setAllocating(true)
    try {
      const res: any = await ipamApi.allocateIPByID(record.id)
      if (res.code === 0) {
        message.success(`IP ${record.ip} 分配成功`)
        fetchData()
      } else {
        message.error(res.message || '分配失败')
      }
    } catch {
      message.error('分配失败')
    } finally {
      setAllocating(false)
    }
  }

  const handleBulkAllocate = async () => {
    try {
      await form.validateFields()
      const res: any = await ipamApi.allocateIP({ subnet_id: subnetId })
      if (res.code === 0) {
        message.success('分配成功')
        setModalOpen(false)
        form.resetFields()
        fetchData()
      } else {
        message.error(res.message || '分配失败')
      }
    } catch {
      // validation or request error
    }
  }

  const columns = [
    { title: 'IP 地址', dataIndex: 'ip', key: 'ip' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const s = statusMap[status] || { label: status, color: 'default' }
        return <Tag color={s.color}>{s.label}</Tag>
      },
    },
    { title: '分配人', dataIndex: 'allocated_by', key: 'allocated_by' },
    { title: 'CI ID', dataIndex: 'ci_id', key: 'ci_id', render: (v: any) => v || '-' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Space>
          {record.status === 'free' && (
            <Button
              type="primary"
              size="small"
              loading={allocating}
              onClick={() => handleAllocate(record)}
            >
              分配
            </Button>
          )}
          {record.status === 'allocated' && (
            <Button size="small" danger onClick={() => handleRelease(record.id)}>
              释放
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <Card
      title={subnetId ? `IP 管理 - ${subnetName || `子网 #${subnetId}`}` : 'IP 管理'}
      extra={
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/ipam')}>
            返回
          </Button>
          {subnetId > 0 && (
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
              自动分配 IP
            </Button>
          )}
        </Space>
      }
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <Select
          placeholder="状态筛选"
          allowClear
          style={{ width: 160 }}
          value={statusFilter || undefined}
          onChange={(v) => setStatusFilter(v || '')}
          options={[
            { label: '空闲', value: 'free' },
            { label: '已分配', value: 'allocated' },
            { label: '预留', value: 'reserved' },
            { label: '禁用', value: 'disabled' },
          ]}
        />
        {!subnetId && <p>请从子网列表或子网拓扑中选择一个子网。</p>}
      </Space>

      <Table
        dataSource={data}
        columns={columns}
        loading={loading}
        rowKey="id"
        pagination={{ pageSize: 20 }}
      />

      <Modal
        title="自动分配 IP"
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields() }}
        onOk={handleBulkAllocate}
        okText="分配"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="note" label="备注">
            <Input placeholder="可选备注" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
