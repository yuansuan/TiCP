import throwError from '@/infra/throw-error'
import { createAllHandlers } from './handler'

class IdGenerator {
  constructor() {
    this._id = 0
  }

  /**
   * @return {number}
   */
  gen() {
    this._id += 1
    return this._id
  }
}

export default class InputInterceptor {
  /**
   * options for controlling input
   * @param {HTMLElement[]} elements
   * @param {handler.BaseHandler[]]} handlers -actually the type is abstract class InputInterceptHandler
   */
  constructor(elements, handlers) {
    if (elements.length === 0) {
      throwError('at least one element should be monitored')
    }

    this._handlers = handlers
    this._elems = elements
  }

  start() {
    this._intercept()
  }

  stop() {
    this._deintercept()
  }

  _intercept() {
    this._elems.forEach(elem => {
      this._handlers.forEach(h => {
        h.intercept(elem)
      })
    })
  }

  _deintercept() {
    this._elems.forEach(elem => {
      this._handlers.forEach(h => {
        h.deintercept(elem)
      })
    })
  }
}

export function makeInputInterceptor(sender, elements) {
  const idGenerator = new IdGenerator()
  const hs = createAllHandlers(sender, idGenerator.gen.bind(idGenerator))
  return new InputInterceptor(elements, hs)
}
