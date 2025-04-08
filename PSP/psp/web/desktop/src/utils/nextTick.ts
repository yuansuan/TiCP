/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

export function nextTick<T>(cb: () => T): Promise<T> {
  // use microtask to mock nextTick
  return Promise.resolve().then(cb)
}
