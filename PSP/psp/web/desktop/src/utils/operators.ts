/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Observable, Subject } from 'rxjs'
import { takeUntil } from 'rxjs/operators'

export const untilDestroyed = (
  componentInstance: object,
  destroyMethodName = 'componentWillUnmount'
) => <T>(source: Observable<T>) => {
  const originalDestroy = componentInstance[destroyMethodName]

  const key = '__takeUntilDestroy'
  if (!componentInstance[key]) {
    componentInstance[key] = new Subject()

    componentInstance[destroyMethodName] = function() {
      componentInstance[key].next(true)
      componentInstance[key].complete()
      if (typeof originalDestroy === 'function') {
        originalDestroy.apply(this, arguments)
      }
    }
  }

  return source.pipe(takeUntil(componentInstance[key]))
}
