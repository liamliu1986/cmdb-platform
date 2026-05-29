import client from './client'

export const discoveryApi = {
  createRule: (data: any) => client.post('/discovery/rules', data),
  listRules: () => client.get('/discovery/rules'),
  executeRule: (data: { rule_id: number }) => client.post('/discovery/rules/execute', data),
  listAgents: () => client.get('/discovery/agents'),
}
