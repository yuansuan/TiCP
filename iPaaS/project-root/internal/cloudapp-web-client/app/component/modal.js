import { addClassName, removeClassName } from './util'

const ACTIVATE_MODAL = 'is-active'

export default class Modal {
  constructor(elem) {
    this._elem = elem
    this._lastActiveElem = null
  }

  show() {
    // FIXME: hacking way to make the container focused again
    this._lastActiveElem = document.activeElement
    addClassName(this._elem, ACTIVATE_MODAL)
  }

  hide() {
    removeClassName(this._elem, ACTIVATE_MODAL)
    if (this._lastActiveElem) {
      this._lastActiveElem.focus()
      this._lastActiveElem = null
    }
  }
}
