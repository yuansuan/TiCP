/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import BoxHttp from '@/domain/Box/BoxHttp'
import { env, currentUser } from '@/domain'
import qs from 'querystring'
import { Http } from '@/utils'
import { UPLOAD_CHUNK_SIZE } from '@/constant'
import { FileServer } from '.'

type DownloadParams = {
  paths: string[]
  types?: boolean[]
  sizes?: string[]
  base?: string
  bucket?: string
  sync_id?: string
  path_rewrite?: any
  boxUrl?: any
  isImage?: boolean
}

export const boxServer: FileServer & {} = {
  delete: ({
    paths,
    bucket,
    project_id
  }: {
    paths: string[]
    bucket: string
    project_id?: string
  }) => {
    return Promise.all(
      paths.map(path =>
        Http.delete('/filemanager/rm', {
          params: {
            path
          }
        })
      )
    )
  },

  list: async ({
    path = '.',
    user_name = '',
    ...params
  }: {
    path?: string
    user_name?: string
  }) => {
    return Http.get(`file/ls`, {
      params: {
        path,
        user_name: currentUser.name,
        ...params
      }
    })
  },

  move: async ({ items }: { items: [string, string][] }) => {
    await Promise.all(
      items.map(([srcPath, destPath]) =>
        Http.post('/file/move', {
          src_paths: srcPath,
          dst_path: destPath,
          overwrite: true
        })
      )
    )
  },

  mkdir: ({ path }: { path: string; user_name: string }) => {
    return Http.post('/file/create_dir', {
      path,
      user_name: currentUser.name
    })
  },

  getContent: async ({
    sync_id,
    ...params
  }: {
    path: string
    offset: number
    length: number
    bucket?: string
    sync_id?: string
  }) => {
    const {
      data: { content }
    } = await BoxHttp.get(`/filemanager${sync_id ? '/remote' : ''}/cat`, {
      params: {
        sync_id,
        ...params
      }
    })

    return content
  },

  stat: ({ ...params }: { path: string }) => {
    return Http.get(`/filemanager/stat`, {
      params: {
        ...params
      }
    })
  },

  getFileUrl: async ({
    paths,
    sync_id,
    boxUrl,
    types,
    ...params
  }: DownloadParams) => {
    if (types && types.length === 1 && types[0]) {
      const names = paths.map(path => path.split('/').pop())

      // !params.isImage &&
      //   Http.post('/filerecord/record', {
      //     type: 2,
      //     info: {
      //       storage_size: params.sizes[0],
      //       file_name: names[0],
      //       file_type: 1 || 0
      //     }
      //   })

      return `${boxUrl}/api/filemanager${
        sync_id ? '/remote' : ''
      }/single/download?path=${encodeURIComponent(paths[0])}`
    } else {
      const {
        data: { token, total_size }
      } = await BoxHttp.post(
        `/filemanager${sync_id ? '/remote' : ''}/download`,
        {
          ...(sync_id
            ? {
                path: paths[0],
                sync_id
              }
            : {
                paths
              }),
          types: types,
          ...params
        }
      )

      return `${boxUrl}/api/filemanager${
        sync_id ? '/remote' : ''
      }/download?token=${token}&total_size=${total_size}`
    }
  },

  download: async (params: DownloadParams) => {
    const aEl = document.createElement('a')
    aEl.href = await boxServer.getFileUrl(params)
    document.body.appendChild(aEl)
    aEl.click()
    document.body.removeChild(aEl)
  },

  upload: async ({
    file,
    path,
    bucket,
    ...params
  }: {
    file: File
    path: string
    bucket: string
  }) => {
    const {
      data: { upload_id }
    } = await BoxHttp.post(
      `/filemanager/pre-upload?${qs.stringify({
        path,
        file_size: file.size,
        bucket,
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
        bucket,
        ...params
      }
      const formData = new FormData()
      formData.append(
        'slice',
        file.slice(UPLOAD_CHUNK_SIZE * index, UPLOAD_CHUNK_SIZE * (index + 1))
      )
      await BoxHttp.post(
        `/filemanager/upload?${qs.stringify(query)}`,
        formData,
        {
          headers: { 'Content-Type': 'multipart/form-data' }
        }
      )
      if (!finish) {
        await uploadChunk(index + 1)
      }
    }

    await uploadChunk(0)
  }
}
