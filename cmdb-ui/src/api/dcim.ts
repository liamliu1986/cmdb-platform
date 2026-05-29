import client from './client'

export const dcimApi = {
  createIDC: (data: any) => client.post('/dcim/idcs', data),
  listIDCs: () => client.get('/dcim/idcs'),
  createRoom: (data: any) => client.post('/dcim/rooms', data),
  createRack: (data: any) => client.post('/dcim/racks', data),
}
