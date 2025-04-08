/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { NewBoxHttp } from '@/domain/Box/NewBoxHttp'
import qs from 'querystring'
import { Http, getFilenameByPath } from '@/utils'
import { UPLOAD_CHUNK_SIZE } from '@/constant'
import { FileServer } from '.'
import { currentUser } from '@/domain'
import { downloadFile } from '@/utils/FileDownload'
import { Buffer } from 'buffer'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'

type DownloadParams = {
  paths: string[]
  types?: boolean[]
  sizes?: string[]
  base?: string
  path_rewrite?: any
  isImage?: boolean
  cross?: boolean
  is_cloud?: boolean
  user_name?: string
}

let is_show = false
EE.on(EE_CUSTOM_EVENT.SHOW_HIDE_FILE, ({ show }) => {
  is_show = show
})
// NOTE: PSP 文件管理
export const newBoxServer: FileServer & {} = {
  getUserCompressStatus: async () => {
    return Http.get('/storage/compressTasks', {
      params: {
        is_cloud: false
      },
      disableErrorMessage: true })
  },

  compress: async ({
    paths,
    target_path,
  }: {
    paths: string[], target_path: string
  }) => {
    console.log(paths, target_path)


    const names = paths.map(p => getFilenameByPath(p))
    let newName = names.slice(0, 2).join('_')

    if (names.length > 2) {
      newName += '等'
    }

    let timeStr = (new Date()).toISOString().substr(0, 19)
   
    return Http.post('/storage/compress', {
      src_paths: paths,
      dst_path: target_path + '/' + `${newName}_${timeStr}.zip`,
      base_path: '.',
      is_cloud: false
    })
  },

  delete: ({
    paths,
    cross = false,
    is_cloud = false
  }: {
    paths: string[]
    cross?: boolean
    is_cloud?: boolean
  }) => {
    const replacePath = paths.map(innerArr =>
      innerArr?.map(str => str?.replace(/^\.\//, '/'))
    )
    return Promise.all(
      replacePath.map(path =>
        Http.post('/storage/remove', {
          paths: path,
          cross,
          is_cloud,
          user_name: currentUser.name
        })
      )
    )
  },

  list: async ({
    path = '.',
    cross = false,
    is_cloud = false,
    user_name = '', //作业详情需要传递当前作业的user_name,不传递就是当前用户
    filter_regexp_list= [],
    recursive = false
  }: {
    path: string
    user_name?: string
    cross?: boolean
    is_cloud?: boolean
    filter_regexp_list?: string[] 
    recursive?: boolean
  }) => {
    if (recursive) {
      return Http.post(`/storage/listOfRecur`, {
        paths: [path],
        cross,
        is_cloud,
        user_name: user_name || currentUser.name,
        filter_regexp_list,
        show_hide_file: is_show
      })
    } else {
      return Http.post(`/storage/list`, {
        path,
        cross,
        is_cloud,
        user_name: user_name || currentUser.name,
        filter_regexp_list,
        show_hide_file: is_show
      })
    }
  },
  rename: async ({
    path,
    newName,
    cross = false,
    is_cloud = false
  }: {
    path: string
    newName: string
    cross?: boolean
    is_cloud?: boolean
  }) => {
    return Http.put(`/storage/rename`, {
      path: path?.replace(/^\.\//, '/'),
      newpath: newName?.replace(/^\.\//, '/'),
      user_name: currentUser.name,
      cross,
      is_cloud
    })
  },
  move: async ({
    srcPaths,
    destPath,
    cross = false,
    is_cloud = false,
    overwrite =false
  }: {
    srcPaths: string
    destPath: string
    cross?: boolean
    is_cloud?: boolean
    overwrite?: boolean
  }) => {
    await Http.put('/storage/move', {
      src_paths: srcPaths,
      dst_path: destPath,
      cross,
      is_cloud,
      overwrite,
      user_name: currentUser.name
    })
  },

  mkdir: async ({
    path,
    cross = false,
    is_cloud = false
  }: {
    path: string
    cross?: boolean
    is_cloud?: boolean
  }) => {
    await Http.post('/storage/createDir', {
      path,
      cross,
      is_cloud,
      user_name: currentUser.name
    })
  },

  getContent: async ({
    user_name,
    ...params
  }: {
    path: string
    offset: number
    len: number
    cross?: boolean
    is_cloud?: boolean
    user_name?: string
  }) => {
    const { data } = await Http.post(`/storage/read`, {
      ...params,
      user_name: user_name || currentUser.name
    })

    return Buffer.from(data, 'base64').toString()
  },

  stat: ({
    sync_id,
    ...params
  }: {
    path: string
    bucket?: string
    sync_id?: string
    url: string
  }) => {
    return NewBoxHttp().get(`/filemanager${sync_id ? '/remote' : ''}/stat`, {
      params: {
        ...params
      }
    })
  },

  getFileUrl: async ({ paths, types, ...params }: DownloadParams) => {},

  download: async (params: DownloadParams) => {
    await downloadFile(
      params.paths,
      params.cross,
      params.is_cloud,
      params.types,
      params.user_name
    )
  },

  linkToCommon: async ({
    current_path,
    dest_dir_path,
    src_dir_paths,
    src_file_paths,
    user_name
  }: {
    current_path: string
    user_name: string
    dest_dir_path: string
    src_dir_paths: string[]
    src_file_paths: string[]
  }) => {
    const res = await Http.post('/storage/link', {
      current_path,
      dest_dir_path,
      // overwrite: true,
      src_dir_paths,
      src_file_paths,
      user_name
    })
    return res
  },
  upload: async ({
    file,
    path,
    url,
    cross = false,
    is_cloud = false,
    ...params
  }: {
    file: File
    path: string
    url: string
    cross?: boolean
    is_cloud?: boolean
  }) => {
    const {
      data: { upload_id }
    } = await Http.post(
      `/storage/pre-upload?${qs.stringify({
        path,
        user_name: currentUser.name,
        file_size: file.size,
        ...params
      })}`
    )

    async function uploadChunk(index) {
      const finish = UPLOAD_CHUNK_SIZE * (index + 1) > file.size
      const query = {
        upload_id,
        path,
        file_size: file.size,
        offset: UPLOAD_CHUNK_SIZE * index,
        slice_size: finish
          ? file.size - UPLOAD_CHUNK_SIZE * index
          : UPLOAD_CHUNK_SIZE,
        finish,
        ...params
      }
      const formData = new FormData()
      formData.append(
        'slice',
        file.slice(UPLOAD_CHUNK_SIZE * index, UPLOAD_CHUNK_SIZE * (index + 1))
      )
      Object.keys(query).forEach(key => {
        formData.append(key, query[key])
      })
      await Http.post(`/file/upload`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      if (!finish) {
        await uploadChunk(index + 1)
      }
    }

    await uploadChunk(0)
  }
}
