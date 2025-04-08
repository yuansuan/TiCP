/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createActionCreator } from 'deox'

export type Dispatch<T extends (...args: any[]) => any> = React.Dispatch<
  Parameters<T>[1]
>

export function createAction<K extends string>(name: K) {
  return function _createAction<T = void, M = void>() {
    return createActionCreator(name, resolve => (payload: T, meta?: M) =>
      resolve(payload, meta)
    )
  }
}

export { createReducer } from 'deox'

export { createStore } from './store'
export { usePrevious } from './store'
