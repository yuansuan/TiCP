/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { env, BoxHttp } from '@/domain'

export default class Result {
  async linkToCommon(pairs: Record<string, string>) {
    const res = await BoxHttp.post('/filemanager/common/link', {
      link_pairs: pairs,
      direction: 'from'
    })
    return res.data.infos
  }
}
