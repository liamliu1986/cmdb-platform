import { useState } from 'react'
import { Modal, Form, Input, Button, Table, message, Select } from 'antd'
import { coreApi } from '@/api/core'

interface Props {
  open: boolean
  onClose: () => void
}

const VALUE_TYPES = [
  'text', 'integer', 'float', 'date', 'bool',
  'choice', 'list', 'password', 'link', 'reference', 'computed',
]

export default function CITypeDesigner({ open, onClose }: Props) {
  const [form] = Form.useForm()
  const [attrs, setAttrs] = useState<any[]>([])
  const [attrForm] = Form.useForm()

  const addAttr = async () => {
    const values = await attrForm.validateFields()
    setAttrs([...attrs, { ...values, id: Date.now() }])
    attrForm.resetFields()
  }

  const handleSubmit = async () => {
    const values = await form.validateFields()
    try {
      const res: any = await coreApi.createCIType({
        ...values,
        attributes: attrs,
      })
      if (res.code === 0) {
        message.success('创建成功')
        form.resetFields()
        setAttrs([])
        onClose()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('创建失败')
    }
  }

  return (
    <Modal
      title="新建 CIType"
      open={open}
      onCancel={onClose}
      onOk={handleSubmit}
      width={700}
    >
      <Form form={form} layout="vertical">
        <Form.Item name="name" label="模型名称" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="alias" label="别名">
          <Input />
        </Form.Item>
      </Form>

      <div style={{ marginTop: 16 }}>
        <h4>属性定义</h4>
        <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
          <Form form={attrForm} layout="inline" style={{ flex: 1 }}>
            <Form.Item name="name" rules={[{ required: true }]} style={{ flex: 1 }}>
              <Input placeholder="属性名" />
            </Form.Item>
            <Form.Item name="alias" style={{ flex: 1 }}>
              <Input placeholder="别名" />
            </Form.Item>
            <Form.Item name="value_type" initialValue="text" style={{ width: 120 }}>
              <Select options={VALUE_TYPES.map(t => ({ label: t, value: t }))} />
            </Form.Item>
          </Form>
          <Button onClick={addAttr}>添加</Button>
        </div>
        <Table
          dataSource={attrs}
          columns={[
            { title: '属性名', dataIndex: 'name' },
            { title: '别名', dataIndex: 'alias' },
            { title: '类型', dataIndex: 'value_type' },
          ]}
          pagination={false}
          size="small"
          rowKey="id"
        />
      </div>
    </Modal>
  )
}
