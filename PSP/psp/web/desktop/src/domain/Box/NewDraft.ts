/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { currentUser, NewBoxHttp } from '@/domain'
import { Http } from '@/utils'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'

export enum DraftType {
  Default = '',
  Vis = 'vis',
  Resubmit = 'resubmit',
  Continue = 'continue'
}
type File = {
  id: string
  isFile: boolean
  mtime: number
  name: string
  parent: File
  children: File[]
  path: string
  size: number
  type: string
}

export default class Draft {
  type: DraftType

  constructor(type: DraftType = DraftType.Default) {
    this.type = type
  }

  async deleteFile(path: string, cross = false, is_cloud = false) {
    return await Http.post('/storage/remove', {
      paths: [path],
      user_name: currentUser.name,
      cross,
      is_cloud
    })
  }

  async deleteFiles(paths: string[], cross = false, is_cloud = false) {
    return await Promise.all(
      paths.map(path => this.deleteFile(path, cross, is_cloud))
    )
  }

  async listFile(recursion: boolean = true) {
    try {
      const res = await NewBoxHttp().get('/filemanager/ls', {
        params: {
          path: '.'
        }
      })
      return res.data ? res.data : []
    } catch (e) {
      return []
    }
  }

  async getFileContent(path: string) {
    try {
      const res = await NewBoxHttp().get('/filemanager/cat', {
        params: {
          path: window.encodeURIComponent(path),
          offset: 0,
          length: 10000,
          bucket_keys: { draft_type: this.type }
        }
      })
      return res.data ? res.data.content : ''
    } catch (e) {
      return ''
    }
  }

  moveFile(srcPath = '', destPath = '', cross = false, is_cloud = false) {
    return Http.put('/storage/move', {
      params: {
        src_paths: srcPath,
        dst_path: destPath || '.',
        cross,
        is_cloud,
        user_name: currentUser.name
      }
    })
  }

  mkdir(path = '') {
    return NewBoxHttp().post('/filemanager/mkdir', {
      path,
      bucket_keys: { draft_type: this.type }
    })
  }

  clean() {
    return NewBoxHttp().post('/filemanager/draft/clean', {
      bucket_keys: { draft_type: this.type }
    })
  }

  submit() {
    return NewBoxHttp().post('/filemanager/draft/submit_to_input', {
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
    return NewBoxHttp().post(
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
    return await NewBoxHttp().post(
      '/filemanager/draft/result_back_to_draft',
      {
        job_id,
        bucket_keys: { draft_type: this.type },
        file_list
      },
      config
    )
  }
  async getFiles(paths, destPath, currentUser) {
    return new Promise((resolve, reject) => {
      const promises = []

      paths.forEach(path => {
        const res = Http.get('/file/ls', {
          params: {
            path: path.replace(/^\./, destPath),

            user_name: currentUser.name
          }
        })

        promises.push(res)
      })

      Promise.all(promises)
        .then(results => {
          const files = results.map(
            result => (result.data && result.data.files) || []
          )

          resolve(files)
        })
        .catch(error => reject(error))
    })
  }

  async getFileList({
    path,
    cross = false,
    is_cloud = false,
    show_hide_file = false
  }) {
    const res = await Http.post('/storage/list', {
      path,
      cross,
      is_cloud,
      user_name: currentUser.name,
      show_hide_file
    })

    return res.data ? res.data : []
  }

  async getFileListOfRecur({
    paths,
    cross = false,
    is_cloud = false,
    show_hide_file = false
  }) {
    const res = await Http.post('/storage/listOfRecur', {
      paths,
      cross,
      is_cloud,
      user_name: currentUser.name,
      show_hide_file
    })

    return res.data ? res.data : []
  }

  async serverFilesToSpuerComputing({
    destPath,
    overwrite,
    selectedFiles,
    cross
  }: {
    destPath: string
    overwrite: boolean
    selectedFiles: File[]
    cross: boolean
  }) {
    const srcDirPaths = selectedFiles
      .filter(file => !file?.isFile)
      .map(item => item?.path)

    const srcFilePaths = selectedFiles
      .filter(file => file?.isFile)
      .map(item => item?.path)

    const res = await Http.post('/storage/hpcUpload/submitTask', {
      overwrite,
      current_path: selectedFiles[0]?.parent?.path,
      dest_dir_path: destPath,
      src_dir_paths: srcDirPaths,
      src_file_paths: srcFilePaths,
      user_name: currentUser.name,
      cross
    })

    EE.emit(EE_CUSTOM_EVENT.SUPERCOMPUTING_TASKKEY, { taskKey: res.data })
    EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, {
      visible: true
    })

    return []
  }
  async copyOrUploadFromCommon(
    pairs: Record<string, string>,
    destPath: string,
    cross = false,
    is_cloud = false,
    selectedFiles: File[]
  ) {
    const src_file_paths = selectedFiles
      .filter(item => item.isFile === true)
      .map(item => item.path.replace(/^\./, currentUser.name))
    const src_dir_paths = selectedFiles
      .filter(item => item.isFile === false)
      .map(item => item.path.replace(/^\./, currentUser.name))

    if (is_cloud) {
      return await this.serverFilesToSpuerComputing({
        selectedFiles,
        destPath,
        overwrite: true,
        cross
      })
    } else {
      if (cross) {
        await Http.post('/storage/link', {
          cross: true,
          is_cloud,
          src_dir_paths,
          src_file_paths,
          dst_path: destPath,
          current_path: selectedFiles[0]?.parent?.path?.replace(
            /^\./,
            currentUser.name
          ),
          user_name: currentUser.name
        })
      }
      return await this.getFileList({
        path: destPath,
        cross,
        is_cloud
      })
    }
  }

  async upload({ filename, file, data }) {
    // 构造要上传的数据
    const formData = new FormData()
    formData.append('file', file)
    formData.append('filename', filename)
    Object.keys(data).forEach(key => formData.append(key, data[key]))
    await NewBoxHttp().post('/filemanager/upload', formData)
  }
}
