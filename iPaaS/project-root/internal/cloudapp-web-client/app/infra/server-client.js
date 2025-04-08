import log from '@/infra/log'
import config from '@/config'
import path from 'path'

const emptyF = () => {}
/**
 * @class
 */
export class ServerClient {
  /**
   * @param {string} host
   * @param {number} port
   */
  constructor(host = window.location.hostname, port = 80) {
    this._host = host
    this._port = port
  }

  // TODO: implement it if necessary
  getHTTPClient() {
    log.error('getHTTPClient is not implemented yet')
  }

  /**
   * @typedef Options
   *  @prop {string} subUrl
   *  @prop {number} port
   *  @prop {function} onOpen
   *  @prop {function} onMessage
   *  @prop {function} onError
   *  @prop {function} onClose
   * @param {Options} options
   */
  createWSConn(options) {
    const { subUrl, onOpen, onMessage, onError, onClose } = options
    const port = options.port || this._port

    const ishttps = window.location.protocol === 'https:'
    let url = `${ishttps ? `wss` : `ws`}://${this._host}:${port}`
    if (subUrl) {
      url = path.join(url, subUrl)
    }
    const ws = new WebSocket(url)
    ws.onopen = onOpen || emptyF
    ws.onmessage = onMessage || emptyF
    ws.onerror = onError || emptyF
    ws.onclose = onClose || emptyF
    return ws
  }
}

export default new ServerClient(window.location.hostname, config.port)
