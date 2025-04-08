import BaseHandler from './basehandler'

/**
 * MousemoveHandler handles the mouse move event
 */
export default class MouseMoveHandler extends BaseHandler {
  /**
   * @implements
   */
  get _eventsMap() {
    if (!this.$m) {
      this.$m = new Map([['mousemove', 'mousemove']])
    }
    return this.$m
  }

  /**
   * @implements
   */
  _genMsg(e) {
    return Promise.resolve({
      x: e.offsetX / e.target.clientWidth,
      y: e.offsetY / e.target.clientHeight
    })
  }
}
