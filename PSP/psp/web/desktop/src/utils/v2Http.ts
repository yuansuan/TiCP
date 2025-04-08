import { createHttp, interceptError } from './HttpClient2'

export const Http = createHttp({ baseURL: '/v2api' })

interceptError(Http)

Http.interceptors.response.use(
  response => response,
  error => {
    const { status } = error.response
    if (status === 401) {
      // window.location.replace(`/api/sso/login${window.location.hash}`)
    }

    return Promise.reject(error)
  }
)

Http.interceptors.response.use(response => {
  return response.data
})
