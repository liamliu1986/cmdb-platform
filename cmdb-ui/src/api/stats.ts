import client from './client'

export const statsApi = {
  getDashboardStats: () => client.get('/stats/dashboard'),
}
