/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'

export function isClassComponent(component) {
  return (typeof component === 'function' &&
    !!component.prototype?.isReactComponent) ||
    (typeof component === 'object' &&
      component &&
      component.$$typeof === Symbol.for('react.memo'))
    ? true
    : false
}

export function isFunctionComponent(component) {
  return typeof component === 'function' ? true : false
}

export function isReactComponent(component) {
  return isClassComponent(component) || isFunctionComponent(component)
    ? true
    : false
}

export function isElement(element) {
  return React.isValidElement(element)
}

export function isDOMTypeElement(element) {
  return isElement(element) && typeof element.type === 'string'
}

export function isCompositeTypeElement(element) {
  return isElement(element) && typeof element.type === 'function'
}
