/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { env } from '@/domain'
import { Http } from '@/utils'
export enum DraftType {
  Default = '',
  Vis = 'vis',
  Resubmit = 'resubmit',
  Continue = 'continue'
}

export default class Draft {
  type: DraftType

  constructor(type: DraftType = DraftType.Default) {
    this.type = type
  }

  async deleteFile(path: string) {
    return await Http.delete('/filemanager/rm', {
      params: {
        path
      }
    })
  }

  async deleteFiles(paths: string[]) {
    return await Promise.all(paths.map(path => this.deleteFile(path)))
  }

  async listFile(recursion: boolean = true) {
    try {
      const res = await Http.get('/filemanager/ls', {
        params: {
          path: '.',
          bucket_keys: { draft_type: this.type }
        }
      })
      return res.data ? res.data.files : []
    } catch (e) {
      return []
    }
  }

  async getFileContent(path: string) {
    try {
      const res = await Http.get('/filemanager/cat', {
        params: {
          path: window.encodeURIComponent(path),
          offset: 0,
          length: 10000
        }
      })
      return res.data ? res.data.content : ''
    } catch (e) {
      return ''
    }
  }

  moveFile(srcPath = '', destPath = '') {
    return Http.post('/filemanager/mv', {
      bucket: 'draft',
      src_path: srcPath,
      dest_path: destPath || '.',
      bucket_keys: { draft_type: this.type }
    })
  }

  mkdir(path = '') {
    return Http.post('/filemanager/mkdir', {
      bucket: 'draft',
      path,
      bucket_keys: { draft_type: this.type }
    })
  }

  clean() {
    return Http.post('/filemanager/draft/clean', {
      bucket_keys: { draft_type: this.type }
    })
  }

  submit() {
    return Http.post('/filemanager/draft/submit_to_input', {
      bucket_keys: { draft_type: this.type }
    })
  }

  back(
    id: string,
    config?: {
      disableErrorMessage?: boolean
      formatErrorMessage?: () => string
    }
  ) {
    return Http.post(
      '/filemanager/draft/input_back_to_draft',
      {
        input_folder_uuid: id,
        bucket_keys: { draft_type: this.type }
      },
      config
    )
  }

  async putResultToDraft(
    job_id,
    file_list?,
    config?: {
      disableErrorMessage?: boolean
      formatErrorMessage?: () => string
    }
  ) {
    return await Http.post(
      '/filemanager/draft/result_back_to_draft',
      {
        job_id,
        bucket_keys: { draft_type: this.type },
        file_list
      },
      config
    )
  }

  async linkFromCommon(pairs: Record<string, string>) {
    const res = await Http.post('/filemanager/common/link', {
      link_pairs: pairs,
      direction: 'to'
    })
    return res.data.infos
  }

  async upload({ filename, file, data }) {
    // 构造要上传的数据
    const formData = new FormData()
    formData.append('file', file)
    formData.append('filename', filename)
    Object.keys(data).forEach(key => formData.append(key, data[key]))
    await Http.post('/filemanager/upload', formData)
  }
}
