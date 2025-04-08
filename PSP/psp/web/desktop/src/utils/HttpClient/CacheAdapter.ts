import axios, { AxiosAdapter } from 'axios'
import { cache } from '../cache'
import buildStoreURL from './buildUrl'

declare module 'axios' {
  interface AxiosRequestConfig {
    useCache?: boolean
  }
}

const adapter = axios.defaults.adapter

interface Options {
  // TODO 以后扩展配置
  moreConfig: any
}

export default function cacheAdapter(
  options: Options = { moreConfig: false }
): AxiosAdapter {
  // const { moreConfig } = options

  return config => {
    const { url, method, params, paramsSerializer, useCache } = config

    if (method === 'get' && useCache) {
      const cloneParams = JSON.parse(JSON.stringify(params))
      delete cloneParams.__timestamp__

      const index = buildStoreURL(url, cloneParams, paramsSerializer)

      let response = null

      let responsePromise = (async () => {
        try {
          response = cache.get(index)
          if (!response) {
            response = await adapter(config)
            cache.set(index, response)
            return response
          } else {
            return response
          }
        } catch (reason) {
          cache.delete(index)
          throw reason
        }
      })()

      return responsePromise
    }

    return adapter(config)
  }
}
