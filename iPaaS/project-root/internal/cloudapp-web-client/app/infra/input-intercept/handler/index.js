import BaseHandler from './basehandler'
import KeyHandler from './key'
import MouseClickHandler from './mouseclick'
import MouseMoveHandler from './mousemove'
import ScrollHandler from './scroll'

export default BaseHandler
export { KeyHandler, MouseClickHandler, MouseMoveHandler, ScrollHandler }

export function createAllHandlers(sender, idGen) {
  return [
    new KeyHandler(sender, idGen),
    new MouseClickHandler(sender, idGen),
    new MouseMoveHandler(sender, idGen),
    new ScrollHandler(sender, idGen),
  ]
}
