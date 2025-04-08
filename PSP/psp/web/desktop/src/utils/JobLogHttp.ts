import { createHttp, interceptError } from '@/utils'

export const Http = createHttp({ baseURL: '/api/v1' })

interceptError(Http)

Http.interceptors.response.use(
  response => response,
  error => {
    const { status } = error.response
    if (status === 401) {
      window.location.replace(`/api/sso/login${window.location.hash}`)
    }

    return Promise.reject(error)
  }
)

Http.interceptors.response.use(response => {
  return response.data
})
