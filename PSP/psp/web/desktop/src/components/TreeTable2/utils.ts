/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { IData } from './types'

export enum ItemTypes {
  TableRow = 'TableRow',
}

export function getAllParentKey(data: IData, keyField = 'key') {
  return data.reduce((ids: string[], curr: IData) => {
    if (!curr.children) {
      return ids
    }
    ids.push(curr[keyField])
    return [...ids, ...getAllParentKey(curr.children, keyField)]
  }, [])
}
