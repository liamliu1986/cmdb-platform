import client from './client'

export const dcimApi = {
  createIDC: (data: any) => client.post('/dcim/idcs', data),
  listIDCs: () => client.get('/dcim/idcs'),
  createRoom: (data: any) => client.post('/dcim/rooms', data),
  listRooms: (idcId: number) => client.get('/dcim/rooms', { params: { idc_id: idcId } }),
  createRack: (data: any) => client.post('/dcim/racks', data),
  listRacks: (roomId: number) => client.get('/dcim/racks', { params: { room_id: roomId } }),
  getRack: (id: number) => client.get(`/dcim/racks/${id}`),
  getRackLayout: (id: number) => client.get(`/dcim/racks/${id}/layout`),
  getRackCapacity: (id: number) => client.get(`/dcim/racks/${id}/capacity`),
  mountDevice: (data: any) => client.post('/dcim/racks/mount', data),
  unmountDevice: (rackId: number, uPosition: number) => client.delete(`/dcim/racks/${rackId}/devices/${uPosition}`),
}
