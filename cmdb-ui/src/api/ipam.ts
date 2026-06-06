import client from './client'

export const ipamApi = {
  createSubnet: (data: any) => client.post('/ipam/subnets', data),
  listSubnets: (params?: { parent_id?: number }) => client.get('/ipam/subnets', { params }),
  listIPsBySubnet: (subnetId: number, params?: { status?: string }) =>
    client.get(`/ipam/subnets/${subnetId}/ips`, { params }),
  allocateIP: (data: { subnet_id: number }) => client.post('/ipam/ips/allocate', data),
  allocateIPByID: (id: number) => client.post(`/ipam/ips/${id}/allocate-by-id`),
  releaseIP: (id: number) => client.post(`/ipam/ips/${id}/release`),
}
