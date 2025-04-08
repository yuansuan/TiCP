/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

// single promise by id
const singleCache = {}
export function single<T>(id: string, resolver: () => Promise<T>): Promise<T> {
  if (!singleCache[id]) {
    singleCache[id] = resolver().finally(() => {
      singleCache[id] = null
    })
    return singleCache[id]
  }

  return Promise.reject()
}
