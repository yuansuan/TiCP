import { stopEvent } from '@/infra/event'
import BaseHandler from './basehandler'

// the which denotes which mouse key is operated
// 0: no button is operated when mouse move
// 1: left mouse button
// 2: middle mouse button
// 3: right mouse button
const whichMap = new Map([
  [1, 0],
  [2, 2],
  [3, 1]
])
const typeMap = new Map([
  ['mousedown', 0],
  ['mouseup', 1]
])
/**
 * In windows server, action bits are as follows:
 * left-down     0x02
 * left-up       0x04
 * right-down    0x08
 * right-up      0x10
 * middle-down   0x20
 * middle-up     0x40
 * @param {number} which
 * @param {string} type
 */
const getActionBit = (which, type) =>
  2 << (whichMap.get(which) * 2 + typeMap.get(type)) // eslint-disable-line

/**
 * MouseClickEvent handles the mouse click event including 'mousedown' 'mouseup'
 */
export default class MouseClickHandler extends BaseHandler {
  /**
   * @implements
   */
  get _eventsMap() {
    if (!this.$m) {
      this.$m = new Map([
        ['mousedown', 'mouseclick'],
        ['mouseup', 'mouseclick']
      ])
    }
    return this.$m
  }

  _blockContexMenu(e) {
    stopEvent(e)
    return false
  }

  /**
   * @implements
   */
  _genMsg(e) {
    const actionBit = getActionBit(e.which, e.type)
    return Promise.resolve({ action_bit: actionBit })
  }

  /**
   * @override
   * @callback
   */
  _onIntercept(elem) {
    elem.addEventListener('contextmenu', this._blockContexMenu.bind(this))
  }

  /**
   * @override
   * @callback
   */
  _onDeintercept(elem) {
    elem.removeEventListener('contextmenu', this._blockContexMenu)
  }
}
