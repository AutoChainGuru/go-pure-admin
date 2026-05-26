import request from './request'

export function login(data) {
  return request.post('/auth/login', data)
}

export function getAuthInfo() {
  return request.get('/auth/info')
}

export function getProfile() {
  return request.get('/auth/profile')
}

export function updateProfile(data) {
  return request.put('/auth/profile', data)
}

export function changePassword(data) {
  return request.put('/auth/password', data)
}
