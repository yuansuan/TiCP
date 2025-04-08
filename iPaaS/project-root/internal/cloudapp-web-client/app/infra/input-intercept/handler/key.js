/* eslint-disable prefer-destructuring */
import BaseHandler from './basehandler'
import UAParser from 'ua-parser-js'
import JSONRPC from '@/infra/jsonrpc'

// https://developer.mozilla.org/zh-CN/docs/Web/API/KeyboardEvent/keyCode
const conflictKeyCodes = { 59: 186, 173: 189, 61: 187 }

/**
 * KeyHandler handles the key event
 */
export default class KeyHandler extends BaseHandler {
  constructor(sender, idGen) {
    super(sender, idGen)
    this._sender = sender
    this._idGen = idGen
    this._protocol = new JSONRPC(idGen)
    this._registryKeyAllUpAction()
  }
  /**
   * @implements
   */
  get _eventsMap() {
    if (!this.$m) {
      this.$m = new Map([
        ['keydown', 'keydown'],
        ['keyup', 'keyup'],
      ])
    }
    return this.$m
  }

  /**
   * @implements
   */

  _genMsg(e) {
    let keyCode = e.keyCode
    let key = e.key
    // escape the Meta key
    if (keyCode === 91) {
      return null
    }

    const parser = new UAParser()
    const browser = parser.getBrowser().name
    //ignore command + key on MacOS
    //http://phabricator.intern.yuansuan.cn/T23193
    if (e.metaKey && parser.getOS() === 'Mac OS') {
      return null
    }

    if (key === 'Alt' || key === 'Shift' || key === 'Control') {
      key = e.code
    }

    if (browser === 'Firefox' && conflictKeyCodes[e.keyCode]) {
      keyCode = conflictKeyCodes[e.keyCode]
    }
    const ret = {
      key,
      code: e.code,
      key_code: keyCode,
      ctrlKey: e.ctrlKey,
      extraMsg: '',
    }

    return Promise.resolve(ret)
  }

  // 修复按住某个键，将页面隐藏/失焦, 重新回复到页面焦点的情况下，通知server将所有键弹起
  _registryKeyAllUpAction() {
    window.onfocus = () => {
      const msg = this._protocol.encode("input", {
        type: "allkeyup",
      })
      this._sender.send(msg)
    }
  }
}
