/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { currentUser } from '@/domain'
import { Http } from '@/utils'
export default class NewResult {
  async linkToCommon(pairs: Record<string, string>, dstPath: string) {
    const res = await Http.post('/storage/link', {
      src_paths: Object.keys(pairs),
      dst_path: dstPath,
      user_name: currentUser.name
    })
    return res.data.files
  }
}
