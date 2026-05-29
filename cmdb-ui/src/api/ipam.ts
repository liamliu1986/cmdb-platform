import client from './client'

export const ipamApi = {
  createSubnet: (data: any) => client.post('/ipam/subnets', data),
  listSubnets: (params?: { parent_id?: number }) => client.get('/ipam/subnets', { params }),
  allocateIP: (data: { subnet_id: number }) => client.post('/ipam/ips/allocate', data),
  releaseIP: (id: number) => client.post(`/ipam/ips/${id}/release`),
}
