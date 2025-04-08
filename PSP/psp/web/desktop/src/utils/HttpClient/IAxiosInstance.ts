import {
  AxiosRequestConfig,
  AxiosPromise,
  AxiosInterceptorManager,
} from 'axios'

interface CustomRequestConfig extends AxiosRequestConfig {
  disableErrorMessage?: boolean
  formatErrorMessage?: (message: string) => string
}

interface AxiosResponse<T = any> {
  data: T
  status: number
  statusText: string
  headers: any
  config: CustomRequestConfig
  request?: any
}

export default interface AxiosInstance {
  (config: CustomRequestConfig): AxiosPromise
  (url: string, config?: CustomRequestConfig): AxiosPromise
  defaults: CustomRequestConfig
  interceptors: {
    request: AxiosInterceptorManager<CustomRequestConfig>
    response: AxiosInterceptorManager<AxiosResponse>
  }
  request<T = any>(config: CustomRequestConfig): Promise<T>
  get<T = any>(url: string, config?: CustomRequestConfig): Promise<T>
  delete(url: string, config?: CustomRequestConfig): Promise<any>
  head(url: string, config?: CustomRequestConfig): Promise<any>
  post<T = any>(
    url: string,
    data?: any,
    config?: CustomRequestConfig
  ): Promise<T>
  put<T = any>(
    url: string,
    data?: any,
    config?: CustomRequestConfig
  ): Promise<T>
  patch<T = any>(
    url: string,
    data?: any,
    config?: CustomRequestConfig
  ): Promise<T>
}
