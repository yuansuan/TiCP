/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, transaction } from 'mobx'
import * as qs from 'qs'

import { Http, Fetch, Validator, formatRegExpStr } from '@/utils'
import { BaseDirectory, BaseFile } from '@/utils/FileSystem'
import Store from '../Store'
import { IPoint } from './Point'

export default class RootPoint extends BaseDirectory implements IPoint {
  @observable pointId
  @observable rootPath = ''
  @observable name = ''

  service = {
    point: this,
    // get file info
    get(path: string) {
      return Http.get('/file/detail', {
        params: {
          paths: path
        },
        formatErrorMessage: msg => `获取文件信息失败`
      })
    },
    // fetch files by path
    fetch(path) {
      return Http.get('/file/ls', {
        params: { path }
      })
        .then((res: any) => {
          if (res.data && res.data.files) {
            return res.data.files
          } else {
            return []
          }
        })
        .then(files => {
          const childPaths = this.getChildPaths(path)

          // delete files
          const filePaths = files.map(item => item.path)
          const deletedFiles = childPaths.filter(path => {
            return !filePaths.includes(path)
          })

          if (deletedFiles.length > 0) {
            Store.delete(deletedFiles)
          }

          // update files
          Store.update(files)

          return files
        })
        .catch(err => {
          // Clear the table when has error
          const childPaths = this.getChildPaths(path)

          if (childPaths.length > 0) {
            Store.delete(childPaths)
          }
        })
    },

    // Get the child file paths
    getChildPaths(path) {
      const parentNode = this.point.filterFirstNode(item => item.path === path)
      const childPaths =
        parentNode && parentNode.children
          ? parentNode.children.map(item => item.path)
          : []

      return childPaths
    },

    // fetch file tree by path
    fetchTree({ path, rootPath }) {
      return Http.get('/file/tree', {
        params: { path, root_path: rootPath }
      })
        .then((res: any) => res.data)
        .then(data => {
          if (!data) {
            return []
          }

          const extractFiles = node => {
            let files = [node]
            if (node.sub_files) {
              node.sub_files.forEach(item => {
                files = [...files, ...extractFiles(item)]
              })
              Reflect.deleteProperty(node, 'sub_files')
            }

            return files
          }
          const files = data.reduce((arr, node) => {
            return [...arr, ...extractFiles(node)]
          }, [])

          // update files
          if (files.length > 0) {
            Store.update(files)
          }

          return files
        })
    },

    // delete file
    delete({ paths }) {
      return Http.post(
        '/file/delete',
        { paths },
        { formatErrorMessage: msg => `删除失败：${msg}` }
      ).then(() => {
        const dirPath = paths[0].replace(/[\\\/][^\\\/]*$/, '')
        // freshen
        this.fetch(dirPath)

        return true
      })
    },

    // rename file/directory
    rename({ path, newName }) {
      const { error } = Validator.filename(newName)
      if (error) {
        return Promise.reject(error)
      }

      return Http.put(
        '/file/rename',
        {
          path,
          new_name: newName
        },
        { formatErrorMessage: msg => `重命名失败：${msg}` }
      ).then(({ data: { file } }: any) => {
        Store.rename(path, newName)
        Store.update([file])

        return true
      })
    },

    // create new directory
    createDir({ path }) {
      const { error } = Validator.filename(path.split(/[\\\/]/).pop())
      if (error) {
        return Promise.reject(error)
      }
      return Http.post(
        '/file/create_dir',
        {
          path
        },
        { formatErrorMessage: msg => `新建失败：${msg}` }
      ).then(() => {
        this.fetch(path.replace(/[\\/][^\\/]*$/, ''))

        return true
      })
    },

    // move to
    move({ srcpaths, dstpath, overwrite }) {
      return Http.post(
        '/file/move',
        {
          srcpaths,
          dstpath,
          overwrite
        },
        { formatErrorMessage: msg => `移动失败：${msg}` }
      ).then(() => {
        const srcDirPath = srcpaths[0].replace(/[\\/][^\\/]*$/, '')
        this.fetch(srcDirPath)
        this.fetch(dstpath)

        return true
      })
    },

    // copy to
    copy({ path, names, toPath }) {
      return Http.post('/file/copy', {
        path,
        names,
        to_path: toPath
      }).then(() => {
        this.fetch(toPath)

        return true
      })
    },

    preDownload({ paths, userId }) {
      return Fetch.post('/api/v3/file/download', {
        paths,
        user_id: userId
      }).then(res => {
        if (!res.success) {
          throw Error(
            {
              40005: '下载失败：非法的文件名格式',
              40060: '下载失败：没有下载权限'
            }[res.code] || '下载失败'
          )
        }

        return res
      })
    },

    download({ token, userId }) {
      const a = document.createElement('a')
      a.href = `/api/v3/file/download?${qs.stringify(
        { token, user_id: userId },
        { indices: false }
      )}`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
    },

    compress({ path, names, compressType, zipName }) {
      return Http.post('/file/compress', {
        parentPath: path,
        srcPaths: names,
        compressType: compressType,
        compressFile: zipName
      })
    },

    extract({ file, toPath, fileType }) {
      // const { error } = Validator.filename(file)
      // if (error) {
      //   return Promise.reject(error)
      // }

      return Http.post('/file/decompress', {
        compressFile: file,
        compressType: fileType,
        destPath: toPath
      })
    },

    view({ path, offset, len }) {
      return Http.get('/file/content', {
        params: { path, offset, len }
      }).then(res => res.data.content)
    },

    edit({ path, content }) {
      return Http.put(
        '/file/edit',
        { path, content },
        { formatErrorMessage: msg => `编辑失败：${msg}` }
      ).then(() => {
        this.fetch(path)
      })
    },

    exist(paths) {
      return Http.post('/file/exist', { paths }).then(res => res.data.isExist)
    }
  }

  constructor(props) {
    super(props)

    this.pointId = props.pointId
    this.rootPath = props.path
    this.name = props.name

    this.init(this, [...Store.nodeMap.values()])

    Store.hooks.afterDelete.tapAsync(
      'Store.delete: sync Point',
      this.onAfterDelete
    )

    Store.hooks.afterUpdate.tapAsync(
      'Store.update: sync Point',
      this.onAfterUpdate
    )

    Store.hooks.afterAdd.tapAsync('Store.add: sync Point', this.onAfterAdd)
  }

  private onAfterDelete = paths => {
    this.removeNodes(item => paths.includes(item.path))
  }

  private onAfterUpdate = nodes => {
    // use transaction to prevent frequent update
    transaction(() => {
      nodes.forEach(({ path, newProps }) => {
        const targetNode = this.filterFirstNode(item => item.path === path)
        targetNode && targetNode.update(newProps)
      })
    })
  }

  private onAfterAdd = nodes => {
    // use transaction to prevent frequent update
    transaction(() => {
      nodes.forEach(props => {
        const { path } = props
        const parentPath = path.replace(/[\\/][^\\/]+$/, '')
        const parentNode = this.filterFirstNode(
          item => item.path === parentPath
        )

        if (parentNode) {
          if (!props.isFile) {
            parentNode.add(new BaseDirectory(props))
          } else {
            parentNode.add(new BaseFile(props))
          }
        }
      })
    })
  }

  init = (parentNode, nodes) => {
    const parentPath = formatRegExpStr(parentNode.path)
    const descendantReg = new RegExp(`^${parentPath}[\\/].+$`)
    const childReg = new RegExp(`^${parentPath}[\\/][^\\/]+$`)

    let descendantNodes = []
    let childNodes = []

    nodes.forEach(node => {
      if (childReg.test(node.path)) {
        childNodes.push(node)
      } else if (descendantReg.test(node.path)) {
        descendantNodes.push(node)
      }
    })

    childNodes.forEach(props => {
      const childNode = !props.isFile
        ? new BaseDirectory(props)
        : new BaseFile(props)
      parentNode.add(childNode)

      this.init(childNode, descendantNodes)
    })
  }
}
