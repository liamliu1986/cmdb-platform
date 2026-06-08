import { useEffect, useState } from 'react'
import { Table, Button, Card, message, Modal, Form, Input, Space } from 'antd'
import { EyeOutlined, PlusOutlined, EnvironmentOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { dcimApi } from '@/api/dcim'

export default function IDCList() {
  const navigate = useNavigate()
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await dcimApi.listIDCs()
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
      const res: any = await dcimApi.createIDC(values)
      if (res.code === 0) {
        message.success('创建成功')
        setModalOpen(false)
        form.resetFields()
        fetchData()
      } else {
        message.error(res.message || '创建失败')
      }
    } catch {
      message.error('创建失败')
    }
  }

  const columns = [
    { title: '名称', dataIndex: 'name' },
    { title: '地址', dataIndex: 'address' },
    { title: '联系人', dataIndex: 'contact' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button size="small" icon={<EyeOutlined />} onClick={() => navigate(`/dcim?idc_id=${record.id}`)}>
            机房
          </Button>
          <Button size="small" icon={<EnvironmentOutlined />} onClick={() => navigate(`/idcmap?idc_id=${record.id}`)}>
            地图
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <Card
      title="数据中心"
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
          新建IDC
        </Button>
      }
    >
      <Table dataSource={data} columns={columns} loading={loading} rowKey="id" />

      <Modal
        title="新建 IDC"
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields() }}
        onOk={handleCreate}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入 IDC 名称' }]}>
            <Input placeholder="如：北京亦庄机房" />
          </Form.Item>
          <Form.Item name="address" label="地址">
            <Input placeholder="如：北京市大兴区..." />
          </Form.Item>
          <Form.Item name="contact" label="联系人">
            <Input placeholder="可选" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
