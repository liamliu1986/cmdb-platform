import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Card, message, Modal, Form, Input, Select, Space } from 'antd'
import { ClusterOutlined, EyeOutlined, ThunderboltOutlined } from '@ant-design/icons'
import { ipamApi } from '@/api/ipam'

export default function SubnetList() {
  const navigate = useNavigate()
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await ipamApi.listSubnets()
      if (res.code === 0) setData(res.data)
    } catch {
      message.error('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const handleCreate = async () => {
    try {
      const values = await form.validateFields()
      const res: any = await ipamApi.createSubnet(values)
      if (res.code === 0) {
        message.success('创建成功')
        setModalOpen(false)
        form.resetFields()
        fetchData()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('创建失败')
    }
  }

  return (
    <Card
      title="子网管理"
      extra={
        <Space>
          <Button icon={<ClusterOutlined />} onClick={() => navigate('/ipam/tree')}>
            拓扑视图
          </Button>
          <Button type="primary" onClick={() => setModalOpen(true)}>新建子网</Button>
        </Space>
      }
    >
      <Table dataSource={data} columns={[
        { title: '名称', dataIndex: 'name' },
        { title: 'CIDR', dataIndex: 'cidr' },
        { title: '状态', dataIndex: 'status' },
        {
          title: '操作',
          key: 'action',
          render: (_: any, record: any) => (
            <Space>
              <Button size="small" icon={<EyeOutlined />} onClick={() => navigate(`/ipam/ips?subnet_id=${record.id}`)}>
                查看 IP
              </Button>
              <Button size="small" icon={<ThunderboltOutlined />} onClick={() => navigate(`/ipam/allocate?subnet_id=${record.id}`)}>
                分配 IP
              </Button>
            </Space>
          ),
        },
      ]} loading={loading} rowKey="id" />

      <Modal
        title="新建子网"
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields() }}
        onOk={handleCreate}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入子网名称' }]}>
            <Input placeholder="如：办公网" />
          </Form.Item>
          <Form.Item name="cidr" label="CIDR" rules={[{ required: true, message: '请输入 CIDR' }]}>
            <Input placeholder="如：192.168.1.0/24" />
          </Form.Item>
          <Form.Item name="vlan_id" label="VLAN ID">
            <Input placeholder="可选" />
          </Form.Item>
          <Form.Item name="parent_id" label="父网">
            <Select
              placeholder="可选"
              allowClear
              options={data.map((s: any) => ({ label: `${s.name} (${s.cidr})`, value: s.id }))}
            />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
