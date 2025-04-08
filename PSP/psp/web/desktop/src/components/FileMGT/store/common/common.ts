/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { IRequest } from './typing'

export * from './typing'
export function formatRequest({
  m_date,
  is_dir,
  name,
  ...rest
}: Partial<IRequest>) {
  const lastName = name.split('/').pop()
  return {
    ...rest,
    path: name,
    name: lastName,
    mtime: m_date,
    ...(!is_dir && {
      type: lastName.split('.').slice(1).pop()
    })
  }
}
