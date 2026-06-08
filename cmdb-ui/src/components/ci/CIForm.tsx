import { Form, Input, InputNumber, Select, DatePicker, Switch } from 'antd'

interface Props {
  attributes: any[]
  form: any
}

export default function CIForm({ attributes, form }: Props) {
  if (!attributes || attributes.length === 0) {
    return <div>该模型暂无属性定义</div>
  }

  const renderField = (attr: any) => {
    switch (attr.value_type) {
      case 'integer':
        return <InputNumber style={{ width: '100%' }} />
      case 'float':
        return <InputNumber step={0.1} style={{ width: '100%' }} />
      case 'bool':
        return <Switch />
      case 'choice':
        return (
          <Select
            options={[]}
            placeholder="请选择"
          />
        )
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
