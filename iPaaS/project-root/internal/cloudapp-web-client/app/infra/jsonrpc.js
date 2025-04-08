const latestVersion = '2.0'

/**
 *
 * @typedef {Object} JSONRPCReq
 *  @prop {number} id
 *  @prop {string} method
 *  @prop {Object} params
 *  @prop {string} jsonrpc
 *
 * @typedef {Object} JSONRPCResp
 *  @prop {number} id
 *  @prop {Object} result
 *  @prop {string} jsonrpc
 */
export default class Protool {
  constructor(idGen) {
    this._idGen = idGen
  }

  /**
   * @param {string} method
   * @param {Object} msg
   * @return {string}
   */
  encode(method, msg) {
    /** @type {JSONRPCReq} */
    const req = {
      id: 0,
      method,
      params: msg,
      jsonrpc: latestVersion,
    }
    return JSON.stringify(req)
  }

  /**
   * @param {string} resp
   * @return {Array}
   */
  decode(resp) {
    /** @type {JSONRPCResp} */
    const res = JSON.parse(resp)
    return [res.method || '', res.params || {}]
  }
}
