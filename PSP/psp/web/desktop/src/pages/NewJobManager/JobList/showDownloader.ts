/* Copyright (C) 2016-present, Yuansuan.cn */

import {
  boxFileServer,
  showDirSelector,
  showFailure
} from '@/components/NewFileMGT'

export async function showDownloader(jobs: { id: string; name: string }[]) {
  const countMap = {}
  const targetDir = (await showDirSelector()).replace(/^\//, '')
  const existNodes = []
  const dir = await boxFileServer.fetch(targetDir)
  const allDirPaths = dir.flatten().map(item => item.path)
  const pathsObj = jobs.reduce((o, p) => {
    const srcPath = p.id
    let dstPath = targetDir ? `${targetDir}/${p.name}` : p.name
    dstPath = dstPath.replace(/^\//, '')

    if (!countMap[dstPath]) {
      countMap[dstPath] = 0
    }
    countMap[dstPath]++

    dstPath =
      countMap[dstPath] > 1 ? `${dstPath}(${countMap[dstPath] - 1})` : dstPath

    if (allDirPaths.includes(dstPath)) {
      existNodes.push({ path: dstPath, srcPath, name: p.name, isFile: false })
    } else {
      o[srcPath] = dstPath
    }
    return o
  }, {})

  if (existNodes.length > 0) {
    const coverNodes = await showFailure({
      actionName: '下载',
      items: existNodes
    })
    if (!!coverNodes.length) {
      await boxFileServer.delete(coverNodes.map(n => n.path))

      coverNodes.reduce((o, n) => {
        o[n.srcPath] = n.path
      }, pathsObj)
    }
  }

  return pathsObj
}
