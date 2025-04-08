/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Directory } from './Directory'
import { File } from './File'
import { BaseDirectory } from '@/utils/FileSystem'
import { formatRequest } from './common'
import { FileServer } from '@/server'
import { currentUser } from '@/domain'

type GetContentParams = { path: string; offset: number; length: number }

const formatPath = path => path.replace(/\/+/, '/').replace(/^\//, '')

export const serverFactory = (fileServer: FileServer) => {
  const server = {
    getUserCompressStatus: () => {
      return fileServer.getUserCompressStatus()
    },

    compress: (paths: string[], target_path: string) => {
      return fileServer.compress({
        paths,
        target_path,
      })
    },

    delete: (paths: string[], cross = false, is_cloud = false) => {
      return fileServer.delete({
        paths: [paths],
        cross,
        is_cloud
      })
    },
    // 获取文件列表
    fetch: async (
      path: string,
      cross = false,
      is_cloud = false,
      filter_regex?:string[]
      ): Promise<BaseDirectory> => {
      const dir = new BaseDirectory({ path })
      try {
        const res = await fileServer.list({
          path: path || '.',
          cross,
          is_cloud,
          user_name: currentUser.name,
          filter_regexp_list:filter_regex
        })
        const files = (res.data || [])
          .sort((x, y) => x.name.split('/').length - y.name.split('/').length)
          .map(item => ({
            ...item,
            name: formatPath(`${path}/${item.name}`)
          }))

        files.forEach(item => {
          let parentPath = item.name.split('/')
          parentPath.pop()
          parentPath = parentPath.join('/')
          let parent = dir

          if (parentPath) {
            parent = dir.filterFirstNode(item => item.path === parentPath)
          }
          parent.push(item.is_dir ? new Directory(item) : new File(item))
        })

        return dir
      } catch (e) {
        return dir
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
      return await fileServer.rename({ path, newName, cross, is_cloud })
    },
    move: async ({
      srcPaths,
      destPath,
      cross = false,
      overwrite = false,
      is_cloud = false
    }: {
      srcPaths: string
      destPath: string
      cross?: boolean
      overwrite?: boolean
      is_cloud?: boolean
    }) => {
      await fileServer.move({ srcPaths, destPath, cross, is_cloud ,overwrite})
    },
    // 创建文件夹
    mkdir: async (path = '', cross = false, is_cloud = false) => {
      await fileServer.mkdir({ path, cross, is_cloud })
    },
    sync: async (oldNode: BaseDirectory) => {
      const { path } = oldNode
      const newNode = await server.fetch(path)
      const oldChildren = [...oldNode.children]
      const newChildren = [...newNode.children]
      const oldPaths = oldChildren.map(item => item.path)
      const newPaths = newChildren.map(item => item.path)

      newChildren.forEach(item => {
        // add new node
        if (!oldPaths.includes(item.path)) {
          oldNode.unshift(item)

          if (!item.isFile) {
            server.fetch(item.path).then(dir => {
              ;(item as any).children = dir.children
            })
          }
        }
      })

      oldChildren.forEach(item => {
        // remove dir
        if (!newPaths.includes(item.path)) {
          item.parent.removeFirstNode(n => n.id === item.id)
        } else {
          // update
          const node = newNode.filterFirstNode(n => n.path === item.path)
          item.update({
            name: node.name,
            path: node.path,
            size: node.size,
            mtime: node.m_date
          })
        }
      })
    },
    getContent: (params: GetContentParams) => {
      return fileServer.getContent({ ...params })
    },
    stat: async (path: string): Promise<ReturnType<typeof formatRequest>> => {
      const { data } = await fileServer.stat({ path })

      return formatRequest({
        ...data,
        name: formatPath(`${path}/${data.name}`)
      })
    },
    getFileUrl: async (
      paths: string[],
      types?: boolean[],
      sizes?: string[],
      isImage?: boolean
    ) => {
      return fileServer.getFileUrl({
        base: '.',
        paths,
        types,
        sizes,
        isImage
      })
    },
    download: async (paths: string[], types?: boolean[], sizes?: string[],cross?:boolean) => {
      return fileServer.download({
        paths,
        types,
        sizes,
        cross
      })
    }
  }

  return server
}
