import client from './client'

export const integrationApi = {
  prometheusQuery: (query: string) =>
    client.get('/integration/prometheus/query', { params: { query } }),
  elkSearch: (data: { query: string; size?: number }) =>
    client.post('/integration/elk/search', data),
  sendTestEmail: (data: { to: string; subject: string; body: string }) =>
    client.post('/integration/email/test', data),
}
