/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { RouterType } from './typing'

function traverse({
  routers,
  tap,
  paths,
}: {
  routers: RouterType[]
  tap: (item: RouterType) => boolean
  paths: RouterType[]
}) {
  routers.every((item, index) => {
    if (paths.length > 0 && tap(paths[paths.length - 1])) {
      return false
    }
    if (index !== 0) {
      paths.pop()
    }
    paths.push(item)
    if (tap(item)) {
      return false
    }
    if (item.children && item.children.length > 0) {
      traverse({
        routers: item.children,
        paths,
        tap,
      })
    }

    return true
  })

  if (paths.length > 0 && !tap(paths[paths.length - 1])) {
    paths.pop()
  }

  return paths
}

export function getBreadCrumb(pathname, routers) {
  return traverse({
    tap: item => {
      if (!item.path) {
        return false
      }
      const targetPaths = item.path.split('/').filter(p => !!p)
      const sourcePaths = pathname.split('/').filter(p => !!p)

      const wildcard =
        targetPaths[targetPaths.length - 1] &&
        targetPaths[targetPaths.length - 1].startsWith(':')
      if (wildcard) {
        targetPaths.pop()
        sourcePaths.pop()
      }
      return targetPaths.join('/') === sourcePaths.join('/')
    },
    routers,
    paths: [],
  })
}
