import client from './client'

export const ipamApi = {
  createSubnet: (data: any) => client.post('/ipam/subnets', data),
  listSubnets: (params?: { parent_id?: number }) => client.get('/ipam/subnets', { params }),
  listIPsBySubnet: (subnetId: number, params?: { status?: string }) =>
    client.get(`/ipam/subnets/${subnetId}/ips`, { params }),
  getIP: (id: number) => client.get(`/ipam/ips/${id}`),
  allocateIP: (data: { subnet_id: number }) => client.post('/ipam/ips/allocate', data),
  allocateIPByID: (id: number) => client.post(`/ipam/ips/${id}/allocate-by-id`),
  releaseIP: (id: number) => client.post(`/ipam/ips/${id}/release`),
  listAvailableIPs: (subnetId: number) => client.get(`/ipam/subnets/${subnetId}/ips/available`),
  assignIPToUser: (userId: number, data: { ip_address_id: number }) => client.post(`/ipam/users/${userId}/ips`, data),
  unassignIPFromUser: (userId: number, ipAddressId: number) => client.delete(`/ipam/users/${userId}/ips/${ipAddressId}`),
  getUserAssignedIPs: (userId: number) => client.get(`/ipam/users/${userId}/ips`),
}
