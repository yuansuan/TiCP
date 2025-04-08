import IdGen from '@/infra/id-gen'
import throwError from '@/infra/throw-error'
import JSONRPC from '@/infra/jsonrpc'
import { stopEvent } from '@/infra/event'
import UAParser from 'ua-parser-js'

const INPUT_METHOD = 'input'

/**
 * @class
 * @abstract
 *  @method _genMsg
 *  @method _eventsMap
 *
 * @callback _onEvent
 * @callback _onIntercept
 * @callback _onDeintercept
 *
 * @public
 *  @method intercept
 *  @method deintercept
 */
export default class BaseHandler {
  /**
   * @param {Object} sender a class with send method
   * @param {function} idGen
   */
  constructor(sender, idGen) {
    if (this.constructor === BaseHandler) {
      throwError('Abstract InputMonitorHandler cannot be instantiated')
    }
    this._sender = sender
    this._protocol = new JSONRPC(IdGen)
    this._idGen = idGen
  }

  /**
   * @protected
   * @param {Object} msg
   * @return {Promise}
   */
  _sendMessage(msg) {
    const o = this._protocol.encode(INPUT_METHOD, msg)
    return new Promise((resolve, reject) => {
      try {
        this._sender.send(o)
        resolve(true)
      } catch (err) {
        reject(err)
      }
    })
  }

  /**
   * @param {Event} e
   * @param {number} seq
   * @param {Object} extra
   * @return {Promise}
   * NOTE: no send if msg is null
   */
  // TODO: send multiple messages for efficiency
  // TODO: compression for multiple messages
  async _hEvent(e, extra) {
    stopEvent(e)
    const msg = await this._genMsg(e, extra)
    if (!msg) {
      return Promise.resolve()
    }
    // Caps Lock on Mac OS
    const parser = new UAParser()
    const osType = parser.getOS().name
    if (e.keyCode == 20 && osType === 'Mac OS') {
      if (e.keyCode == 20) {
        await this._sendMessage({
          seq: this._idGen(),
          type: 'keydown',
          timestamp: e.timeStamp,
          msg,
        })
        return this._sendMessage({
          seq: this._idGen(),
          type: 'keyup',
          timestamp: e.timeStamp,
          msg,
        })
      }
    }
    return this._sendMessage({
      seq: this._idGen(),
      type: this._getDomainEventType(e.type),
      timestamp: e.timeStamp,
      msg,
    })
  }

  /**
   * @param {Event} e
   */
  _wraplistener(e) {
    const extra = this._onEvent(e)
    this._hEvent(e, extra)
  }

  /**
   * @param {Element} elem
   */
  intercept(elem) {
    this._onIntercept(elem)
    // store the listener func to _lis for remvoval
    this._lis = this._wraplistener.bind(this)
    this._getHandleableEventTypes().forEach((eventName) => {
      elem.addEventListener(eventName, this._lis)
    })
  }

  /**
   * @param {Element} elem
   */
  deintercept(elem) {
    this._onDeintercept(elem)
    this._getHandleableEventTypes().forEach((eventName) => {
      elem.removeEventListener(eventName, this._lis)
    })
  }

  /**
   * raw html event type => domain event type
   * @protected
   * @param {string} type
   */
  _getDomainEventType(type) {
    return this._eventsMap.get(type)
  }

  /**
   * @return {string[]}
   */
  _getHandleableEventTypes() {
    return Array.from(this._eventsMap.keys())
  }

  /**
   * @abstract
   * @return {Map<string, string>}
   */
  get _eventsMap() {} // eslint-disable-line

  /**
   * @abstract
   * @param {Event} e
   * @param {Object} extra
   * @return {Promise}
   */
  _genMsg(e, seq) {} // eslint-disable-line

  /**
   * @callback
   * @param {Event} e
   * @return {Object} {extra}
   */
  _onEvent(e) {} // eslint-disable-line

  /**
   * @callback
   * @param {Element} elem
   */
  _onIntercept(elem) {} // eslint-disable-line

  /**
   * @callback
   * @param {Element} elem
   */
  _onDeintercept(elem) {} // eslint-disable-line
}
