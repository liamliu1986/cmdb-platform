import { useState } from 'react'
import { Button, Input, Space, Table, message } from 'antd'
import { integrationApi } from '@/api/integration'

interface Props { ci: any }

export default function CILogTab({ ci }: Props) {
  const [loading, setLoading] = useState(false)
  const [logs, setLogs] = useState<any[]>([])
  const [keyword, setKeyword] = useState('')

  const fetchLogs = () => {
    setLoading(true)
    const hostname = ci?.attr_values?.hostname || '*'
    const q = keyword ? `host:${hostname} AND ${keyword}` : `host:${hostname}`
    integrationApi.elkSearch({ query: q, size: 50 })
      .then((res: any) => {
        if (res.code === 0 && res.data?.hits?.hits) {
          setLogs(res.data.hits.hits.map((h: any) => ({
            timestamp: h._source?.['@timestamp'] || '-',
            message: h._source?.message || h._source?.log || JSON.stringify(h._source),
          })))
        } else {
          setLogs([])
        }
      })
      .catch(() => message.info('无法连接 Elasticsearch'))
      .finally(() => setLoading(false))
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Input.Search
          placeholder="搜索关键字..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onSearch={fetchLogs}
          style={{ width: 300 }}
        />
        <Button onClick={fetchLogs} loading={loading}>查询日志</Button>
        <span style={{ color: '#999', fontSize: 12 }}>
          ES URL: {process.env.VITE_ES_URL || 'http://localhost:9200'}
        </span>
      </Space>
      <Table
        dataSource={logs}
        loading={loading}
        columns={[
          { title: '时间', dataIndex: 'timestamp', width: 200 },
          { title: '消息', dataIndex: 'message', ellipsis: true },
        ]}
        pagination={{ pageSize: 20 }}
        size="small"
        rowKey={(_, i) => String(i)}
        locale={{ emptyText: '暂无日志数据 - 请先查询' }}
      />
    </div>
  )
}
