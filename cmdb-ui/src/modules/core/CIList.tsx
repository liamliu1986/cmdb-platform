import { useEffect, useState } from 'react'
import { Table, Button, Card, Input, message, Modal, Form, Select } from 'antd'
import { coreApi } from '@/api/core'
import CIForm from '@/components/ci/CIForm'

export default function CIList() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [ciTypes, setCITypes] = useState<any[]>([])
  const [selectedCIType, setSelectedCIType] = useState<any>(null)
  const [ciTypeAttrs, setCITypeAttrs] = useState<any[]>([])
  const [form] = Form.useForm()
  const [searchQ, setSearchQ] = useState('')

  const fetchCITypes = async () => {
    try {
      const res: any = await coreApi.listCITypes()
      if (res.code === 0) setCITypes(res.data)
    } catch {
      // ignore
    }
  }

  const fetchData = async () => {
    setLoading(true)
    try {
      const res: any = await coreApi.searchCI({ q: searchQ, page: 1, page_size: 25 })
      if (res.code === 0) setData(res.data.list)
    } catch {
      message.error('加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchCITypes(); fetchData() }, [])

  const handleCITypeChange = (typeId: number) => {
    const ct = ciTypes.find((t: any) => t.id === typeId)
    setSelectedCIType(ct)
    // For now, use mock attributes. In a real implementation,
    // you would fetch attributes from the backend.
    setCITypeAttrs([
      { name: 'name', alias: '名称', value_type: 'text', is_required: true },
      { name: 'description', alias: '描述', value_type: 'text', is_required: false },
    ])
    form.resetFields()
  }

  const handleCreate = async () => {
    if (!selectedCIType) {
      message.error('请选择模型')
      return
    }
    const values = await form.validateFields()
    try {
      const res: any = await coreApi.createCI({
        ci_type_id: selectedCIType.id,
        attr_values: values,
      })
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

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '类型ID', dataIndex: 'ci_type_id', key: 'ci_type_id' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '更新人', dataIndex: 'updated_by', key: 'updated_by' },
    { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at' },
  ]

  return (
    <Card
      title="资源管理"
      extra={
        <div style={{ display: 'flex', gap: 8 }}>
          <Input.Search
            placeholder="搜索..."
            value={searchQ}
            onChange={(e) => setSearchQ(e.target.value)}
            onSearch={fetchData}
            style={{ width: 300 }}
          />
          <Button type="primary" onClick={() => setModalOpen(true)}>新建资源</Button>
        </div>
      }
    >
      <Table dataSource={data} columns={columns} loading={loading} rowKey="id" />

      <Modal
        title="新建 CI"
        open={modalOpen}
        onCancel={() => { setModalOpen(false); form.resetFields() }}
        onOk={handleCreate}
        width={600}
      >
        <Select
          placeholder="选择模型"
          style={{ width: '100%', marginBottom: 16 }}
          options={ciTypes.map((ct: any) => ({ label: ct.alias || ct.name, value: ct.id }))}
          onChange={handleCITypeChange}
        />
        {selectedCIType && (
          <CIForm attributes={ciTypeAttrs} form={form} />
        )}
      </Modal>
    </Card>
  )
}
