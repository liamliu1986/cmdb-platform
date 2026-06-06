import { useEffect, useState } from 'react'
import { Form, Input, InputNumber, Select, DatePicker, Switch, Tag, Space } from 'antd'
import { ipamApi } from '@/api/ipam'

interface Props {
  attributes: any[]
  form: any
}

// IP Reference Selector — shows IPs assigned to the current user
function IPReferenceSelector({ value, onChange }: any) {
  const [ips, setIPs] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    setLoading(true)
    // Get current user ID from localStorage (set at login)
    const userId = localStorage.getItem('user_id')
    if (!userId) {
      setLoading(false)
      return
    }
    ipamApi.getUserAssignedIPs(parseInt(userId))
      .then((res: any) => {
        if (res.code === 0) setIPs(res.data || [])
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const selectedIP = ips.find((ip: any) => ip.id === value)

  return (
    <Space direction="vertical" style={{ width: '100%' }}>
      <Select
        placeholder="选择已分配的 IP"
        style={{ width: '100%' }}
        value={value}
        onChange={onChange}
        loading={loading}
        allowClear
        showSearch
        options={ips.map((ip: any) => ({
          label: `${ip.ip} (${ip.status})`,
          value: ip.id,
        }))}
      />
      {selectedIP && (
        <Tag color="blue">{selectedIP.ip}</Tag>
      )}
    </Space>
  )
}

export default function CIForm({ attributes, form }: Props) {
  if (!attributes || attributes.length === 0) {
    return <div>该模型暂无属性定义</div>
  }

  const renderField = (attr: any) => {
    if (attr.is_reference && attr.ref_table === 'cmdb_ipam.ip_addresses') {
      return <IPReferenceSelector />
    }
    switch (attr.value_type) {
      case 'integer':
        return <InputNumber style={{ width: '100%' }} />
      case 'float':
        return <InputNumber step={0.1} style={{ width: '100%' }} />
      case 'bool':
        return <Switch />
      case 'choice':
        return <Select options={[]} placeholder="请选择" />
      case 'date':
        return <DatePicker style={{ width: '100%' }} />
      case 'password':
        return <Input.Password />
      default:
        return <Input />
    }
  }

  return (
    <Form form={form} layout="vertical">
      {attributes.map((attr: any) => (
        <Form.Item
          key={attr.id || attr.name}
          name={attr.name}
          label={attr.alias || attr.name}
          rules={attr.is_required ? [{ required: true, message: `请输入${attr.alias || attr.name}` }] : []}
        >
          {renderField(attr)}
        </Form.Item>
      ))}
    </Form>
  )
}
