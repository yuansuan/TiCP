import log from './log'

export default function throwError(msg, ...kvs) {
  if (kvs && kvs.length > 0) {
    if (kvs.length % 2) {
      kvs.push('')
    }
    log.error('Extra error informations:', ...kvs)
  }
  throw new Error(msg)
}
