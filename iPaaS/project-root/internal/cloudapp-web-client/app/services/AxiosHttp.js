import axios from 'axios'

axios.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response.status == 401) {
      location.href = '/#/login'
    }
    return Promise.reject(error)
  }
)

export default {
  get(url, params, isPublic = false) {
    params = params || {}
    const headers = {}
    if (!isPublic) {
      headers.Authorization = `Bearer ${sessionStorage.access_token}`
    }
    return axios.get(url, { params, headers })
  },
  post(url, data, isPublic = false) {
    const headers = {
      Accept: '*/*',
      'Content-Type': 'application/json'
    }
    if (!isPublic) {
      headers.Authorization = `Bearer ${sessionStorage.access_token}`
    }
    return axios.post(url, data, { headers })
  },
  put(url, data) {
    const headers = {
      Accept: '*/*',
      'Content-Type': 'application/json'
    }
    return axios.put(url, data, { headers })
  },
  postAsForm(url, data) {
    const param = new URLSearchParams(data)
    return axios.post(url, param, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    })
  }
}
