import log from '@/infra/log'
import { isMacOS } from '@/infra/os-detector'

import BaseHandler from './basehandler'

/**
 * @param {Event} e
 * @return {Object}
 */
function extractDeltaXY(e) {
  const dx = e.deltaX || 0
  const dy = e.deltaY || 0
  let scale = 1
  switch (e.deltaMode) {
    /**
     * 0 = pixels
        1 = lines
        2 = pages
     */
    case 1:
      log.info('wont fix the such mode')
      break
    case 2:
      scale = window.innerHeight
      break
    default:
  }
  return {
    dx: dx * scale,
    dy: dy * scale
  }
}

export default class WheelHandler extends BaseHandler {
  /**
   * @implements
   */
  get _eventsMap() {
    if (!this.$m) {
      this.$m = new Map([['wheel', 'scroll']])
    }
    return this.$m
  }

  /**
   * @override
   * @callback
   */
  _onEvent(e) {
    return extractDeltaXY(e)
  }

  /**
   * @implements
   */
  _genMsg(e, extra) {
    return Promise.resolve({
      x: e.x / e.target.clientWidth,
      y: e.y / e.target.clientHeight,
      dx: extra.dx,
      dy: extra.dy,
    })
  }
}
