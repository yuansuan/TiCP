export function addClassName(elem, className) {
  const names = elem.className.split(/\s+/)
  if (names.indexOf(className) === -1) {
    elem.className += ` ${className}`
  }
}

export function removeClassName(elem, className) {
  const names = elem.className.split(/\s+/)
  if (names.indexOf(className) > -1) {
    names.splice(names.indexOf(className), 1)
    elem.className = names.join(' ')
  }
}
