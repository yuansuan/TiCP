/**
 * @param {Event} e
 * TODO: better to check `cancellable` of the event
 * Ref: https://developer.mozilla.org/en-US/docs/Web/API/Event/cancelable
 */
export function stopEvent(e) {
  if (e.preventDefault) {
    e.preventDefault()
  }
  if (e.stopPropagation) {
    e.stopPropagation()
  }
}
