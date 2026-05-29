import client from './client'

export const coreApi = {
  createCIType: (data: any) => client.post('/citypes', data),
  listCITypes: () => client.get('/citypes'),
  getCIType: (id: number) => client.get(`/citypes/${id}`),
  createCI: (data: any) => client.post('/ci', data),
  getCI: (id: number) => client.get(`/ci/${id}`),
  searchCI: (params: { q?: string; page?: number; page_size?: number; sort?: string }) =>
    client.get('/ci/s', { params }),
}
