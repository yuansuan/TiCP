/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { JobDirectory } from './JobDirectory'
import { JobFile } from './JobFile'
import { v4 as uuidv4 } from 'uuid'
import { escapeRegExp } from '@/utils'

interface IRemoteCommonFile {
  name: string
  size: number
  path: string
  mod_time: number
  is_dir: boolean
  common_path?: string
}

export class FileTree extends JobDirectory {
  uploadLocalFile = (id, file, isTempDir = true, tempDirPath) => {
    const node = this.filterFirstNode(item => item.id === id)
    if (node && !node.isFile) {
      const { webkitRelativePath = '' } = file.originFileObj || {}
      if (!webkitRelativePath) {
        node.unshift(
          new JobFile({
            name: file.name,
            uid: file.uid,
            size: file.size,
            percent: file.percent,
            status: file.status
          })
        )
      } else {
        const paths = webkitRelativePath.split('/')
        const filename = paths.pop()
        const dirPath = paths.join('/')
        let ensureDir = `${node.path}/${dirPath}`

        const dir = this.ensureDir(
          isTempDir
            ? ensureDir
            : ensureDir.replace(new RegExp(escapeRegExp(tempDirPath)), '')
        )
        dir.unshift(
          new JobFile({
            name: filename,
            uid: file.uid,
            size: file.size,
            percent: file.percent,
            status: file.status
          })
        )
      }
    }
  }

  uploadCommonFiles(
    files: IRemoteCommonFile[],
    tempDir: String,
    selectFileSet = new Set(),
    isTempDir = true,
    id = null,
    currentDirPathSet = new Set()
  ) {
    let expandKeys = new Set()
    if (true) {
      // 处理目录
      files
        .filter(i => i.is_dir)
        .forEach(file => {
          const paths = file.path
            .replace(tempDir + '/', '')
            .split('/')
            .filter(i => !!i)

          const dirPath = paths.join('/')
          const common_path_prefix = file.path.replace(tempDir + '/', '')

          // if (file?.common_path) {
          //   let lastIndex = file?.common_path.lastIndexOf(file.name)
          //   common_path_prefix = file?.common_path.substring(0, lastIndex)
          // }

          const dir = this.ensureDir(dirPath, common_path_prefix)
          if (
            currentDirPathSet.has(common_path_prefix) &&
            !expandKeys.has(dir.id)
          ) {
            expandKeys.add(dir.id)
          }
        })

      // 处理文件
      files
        .filter(i => !i.is_dir)
        .forEach(file => {
          const paths = file.path
            .replace(tempDir + '/', '')
            .split('/')
            .filter(i => !!i)
          const filename = paths.pop()
          const dirPath = paths.join('/')

          let common_path_prefix = file.path.replace(tempDir + '/', '')

          // if (file?.common_path) {
          //   let lastIndex = file?.common_path.lastIndexOf(file.name)
          //   common_path_prefix = file?.common_path.substring(0, lastIndex)
          // }

          const dir = this.ensureDir(dirPath, common_path_prefix)
          if (!expandKeys.has(dir.id)) {
            expandKeys.add(dir.id)
          }

          let jobFile = new JobFile({
            name: filename,
            uid: uuidv4(),
            size: file.size,
            realCommonPath: file?.common_path
          })

          let selectedFilePath = filename
          if (dirPath !== '') {
            selectedFilePath = dirPath + '/' + filename
          }
          if (selectFileSet !== null && selectFileSet.has(selectedFilePath)) {
            jobFile.isMain = true
          }

          dir.unshift(jobFile)
        })
    } else {
      const node = this.filterFirstNode(item => item.id === id)
      if (node && !node.isFile) {
        files.forEach(({ name, is_dir, path, ...rest }) => {
          path = './' + path

          if (is_dir) {
            const tmp = path.replace(new RegExp(escapeRegExp(node.path)), '')
            node.ensureDir(tmp)
          } else {
            const paths = path.split('/')
            const parentPath = paths.slice(0, paths.length - 1).join('/')

            const tmp = parentPath.replace(
              new RegExp(escapeRegExp(node.path)),
              ''
            )
            const parent = tmp ? node.ensureDir(tmp) : node
            const file = new JobFile({ name, is_dir, ...rest })
            parent.push(file)
          }
        })
      }
    }
    return expandKeys
  }
}
