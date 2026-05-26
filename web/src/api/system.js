import request from './request'

export const userApi = {
  list: (params) => request.get('/system/users', { params }),
  create: (data) => request.post('/system/users', data),
  update: (id, data) => request.put(`/system/users/${id}`, data),
  remove: (id) => request.delete(`/system/users/${id}`),
}

export const roleApi = {
  list: (params) => request.get('/system/roles', { params }),
  all: () => request.get('/system/roles/all'),
  create: (data) => request.post('/system/roles', data),
  update: (id, data) => request.put(`/system/roles/${id}`, data),
  remove: (id) => request.delete(`/system/roles/${id}`),
  getMenus: (id) => request.get(`/system/roles/${id}/menus`),
  assignMenus: (id, menuIds) => request.put(`/system/roles/${id}/menus`, { menuIds }),
}

export const menuApi = {
  tree: () => request.get('/system/menus'),
  create: (data) => request.post('/system/menus', data),
  update: (id, data) => request.put(`/system/menus/${id}`, data),
  remove: (id) => request.delete(`/system/menus/${id}`),
}

export const deptApi = {
  tree: () => request.get('/system/depts'),
  create: (data) => request.post('/system/depts', data),
  update: (id, data) => request.put(`/system/depts/${id}`, data),
  remove: (id) => request.delete(`/system/depts/${id}`),
}

export const dictApi = {
  list: (params) => request.get('/system/dicts', { params }),
  create: (data) => request.post('/system/dicts', data),
  update: (id, data) => request.put(`/system/dicts/${id}`, data),
  remove: (id) => request.delete(`/system/dicts/${id}`),
  listDetails: (dictId) => request.get('/system/dict-details', { params: { dictId } }),
  createDetail: (data) => request.post('/system/dict-details', data),
  updateDetail: (id, data) => request.put(`/system/dict-details/${id}`, data),
  removeDetail: (id) => request.delete(`/system/dict-details/${id}`),
}

export const loginLogApi = {
  list: (params) => request.get('/system/login-logs', { params }),
}

export const operLogApi = {
  list: (params) => request.get('/system/oper-logs', { params }),
  get: (id) => request.get(`/system/oper-logs/${id}`),
}
