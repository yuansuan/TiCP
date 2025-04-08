/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { serverFactory } from './store/common/server'
import { newBoxServer } from '@/server'

// 右侧表格区域+当前目录...
export { FileList } from './FileList'
// 左侧搜索+新建文件夹+Home
export { Menu } from './Menu'
// 文件位置这一行
export { History } from './History'
// 上传+移动这一行
export { Toolbar } from './Toolbar'
export * from './store'
export { Suite } from './Suite'
export { showFileSelector } from './FileSelector'
export { showDirSelector } from './DirSelector'
export { showFailure } from './Failure'

export const boxFileServer = serverFactory(newBoxServer)
