import { useState } from 'react'
import { Card, Input, Button, Tabs, message } from 'antd'
import { integrationApi } from '@/api/integration'

export default function MonitorDashboard() {
  const [promQuery, setPromQuery] = useState('')
  const [promResult, setPromResult] = useState('')
  const [elkQuery, setElkQuery] = useState('')
  const [elkResult, setElkResult] = useState('')
  const [emailTo, setEmailTo] = useState('')
  const [emailSubject, setEmailSubject] = useState('')
  const [emailBody, setEmailBody] = useState('')

  const handlePromQuery = async () => {
    try {
      const res: any = await integrationApi.prometheusQuery(promQuery)
      setPromResult(JSON.stringify(res.data, null, 2))
    } catch (err: any) {
      message.error(err.message || '查询失败')
    }
  }

  const handleElkSearch = async () => {
    try {
      const res: any = await integrationApi.elkSearch({ query: elkQuery })
      setElkResult(JSON.stringify(res.data, null, 2))
    } catch (err: any) {
      message.error(err.message || '搜索失败')
    }
  }

  const handleSendEmail = async () => {
    try {
      await integrationApi.sendTestEmail({ to: emailTo, subject: emailSubject, body: emailBody })
      message.success('邮件发送成功')
    } catch (err: any) {
      message.error(err.message || '发送失败')
    }
  }

  const items = [
    {
      key: 'prometheus',
      label: 'Prometheus',
      children: (
        <Card title="PromQL 查询">
          <Input.TextArea
            value={promQuery}
            onChange={(e) => setPromQuery(e.target.value)}
            placeholder="输入 PromQL，如: up"
            rows={3}
          />
          <Button type="primary" onClick={handlePromQuery} style={{ marginTop: 8 }}>查询</Button>
          <pre style={{ marginTop: 16, background: '#f5f5f5', padding: 16 }}>{promResult || '结果将显示在这里'}</pre>
        </Card>
      ),
    },
    {
      key: 'elk',
      label: 'ELK',
      children: (
        <Card title="日志搜索">
          <Input.TextArea
            value={elkQuery}
            onChange={(e) => setElkQuery(e.target.value)}
            placeholder="输入搜索条件"
            rows={3}
          />
          <Button type="primary" onClick={handleElkSearch} style={{ marginTop: 8 }}>搜索</Button>
          <pre style={{ marginTop: 16, background: '#f5f5f5', padding: 16 }}>{elkResult || '结果将显示在这里'}</pre>
        </Card>
      ),
    },
    {
      key: 'email',
      label: '邮件测试',
      children: (
        <Card title="发送测试邮件">
          <Input placeholder="收件人" value={emailTo} onChange={(e) => setEmailTo(e.target.value)} style={{ marginBottom: 8 }} />
          <Input placeholder="主题" value={emailSubject} onChange={(e) => setEmailSubject(e.target.value)} style={{ marginBottom: 8 }} />
          <Input.TextArea placeholder="内容" value={emailBody} onChange={(e) => setEmailBody(e.target.value)} rows={4} style={{ marginBottom: 8 }} />
          <Button type="primary" onClick={handleSendEmail}>发送</Button>
        </Card>
      ),
    },
  ]

  return (
    <Card title="监控集成">
      <Tabs items={items} />
    </Card>
  )
}
