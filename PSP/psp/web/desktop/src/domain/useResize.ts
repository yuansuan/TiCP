/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { useLayoutRect, useDidUpdate } from '@/utils/hooks'
import PageLayout from '@/components/PageLayout'
import { useLayoutEffect } from 'react'

export function useResize(): [
  ClientRect,
  React.MutableRefObject<any>,
  () => void
] {
  const [rect, ref, resize] = useLayoutRect()
  const store = PageLayout.useStore()

  useLayoutEffect(() => {
    resize()
  }, [])

  useDidUpdate(() => {
    setTimeout(resize, 300)
  }, [store?.menuExpanded?.[0]])

  return [rect, ref, resize]
}
