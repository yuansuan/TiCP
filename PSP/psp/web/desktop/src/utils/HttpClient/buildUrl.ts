// @ts-ignore
import buildURL from 'axios/lib/helpers/buildURL'

export default function buildStoreURL(...args: any[]) {
  const builtURL = buildURL(...args)

  const [urlPath, queryString] = builtURL.split('?')

  if (queryString) {
    const paramsPair = queryString.split('&')
    return `${urlPath}?${paramsPair.sort().join('&')}`
  }

  return builtURL
}
