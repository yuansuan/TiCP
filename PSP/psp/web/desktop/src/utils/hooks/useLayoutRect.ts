/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useEffect, useRef } from 'react'

export function useLayoutRect(): [
  ClientRect,
  React.MutableRefObject<any>,
  () => void
] {
  const [rect, setRect] = useState({
    width: 0,
    height: 0,
    left: 0,
    top: 0,
    bottom: 0,
    right: 0,
  })
  const ref = useRef(null)

  const getClientRect = () => {
    if (ref.current) {
      setRect(ref.current.getBoundingClientRect())
    }
  }

  useEffect(() => {
    window.addEventListener('resize', getClientRect)

    return () => {
      window.removeEventListener('resize', getClientRect)
    }
  }, [])

  return [rect, ref, getClientRect]
}
