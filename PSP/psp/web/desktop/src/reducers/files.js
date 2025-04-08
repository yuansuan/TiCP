/* Copyright (C) 2016-present, Yuansuan.cn */
import { Bin, Item } from '../utils/bin'
import { message } from 'antd'
import fdata from './dir.json'
import { env, uploader, NewBoxHttp, sysConfig } from '@/domain'
import { serverFactory } from '@/components/NewFileMGT/store/common'
import { showDirSelector } from '@/components/NewFileMGT/DirSelector'
import { showShareFile } from '@/components/ShareFile'
import { showFailure } from '@/components/NewFileMGT'
import { newBoxServer } from '@/server'
import { v4 as uuid } from 'uuid'
import { globalSizes } from '@/domain/Box/states'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { formatByte, formatUnixTime } from '@/utils'
import { currentUser } from '@/domain'
const server = serverFactory(newBoxServer)
import { Http, Validator } from '@/utils'
import { hijackUploaderController } from '@/utils/uploader'

const defState = {
  cdir: '%homedrive%', // 当前所在的目录
  hist: [],
  hid: 0,
  view: 1,
  editing: false
}

const getDuplicate = item => {
  let duplicateNode = null

  item.data.every(child => {
    // check duplicate name
    if (child.name === item.name && item.id !== child.id) {
      duplicateNode = child
      return false
    }

    return true
  })

  return duplicateNode
}

const _compress = (itemPaths, currentDir, callback) => {
  server.compress(itemPaths, currentDir.path || '/').then((res) => {
    message.success('开始压缩')

    let packageId = res.data.package_id
    const checkPackageStatus = () => {
      server.getUserCompressStatus().then((res) => {
        if (res.data.length === 0) {
          callback(currentDir.path || '/', currentDir.id)
        } else {
          const hasRunning = res.data.some(l => l.Status === 2)
            
          if (hasRunning) {
            setTimeout(checkPackageStatus, 3000)
          } else {
            // TODO 目前是单任务，代码可以正常运行，支持多个任务，需要修改代码
            const hasFailed = res.data.some(l => l.Status === 1)

            if (hasFailed) {
              message.error('压缩任务失败')
            }

            setTimeout(checkPackageStatus, 3000)
          }         
        }
      }).catch((e) => {
        console.log(e)
        message.error('检查压缩任务失败')
      })
    }
    checkPackageStatus()
  }).catch(() => {
    message.error('压缩失败')
  })
}

const upload = props => {
  uploader.upload({
    action: '/storage/upload',
    httpAdapter: Http,
    ...props,
    origin: props.origin,
    data: {
      ...props.data
    }
  })
}

// 劫持上传控制器以实现异步并发上传
hijackUploaderController(uploader)
//  文件管理上传
const _upload = (directory = false, dir, callback) => {
  const UPLOAD_ID = uuid()

  const uniqueID = uuid()
  const dirFinal = `/${(dir.path || '/').replace(/^\//, '')}`
  let uploadFilesLength = 0
  let uploadSucceedCount = 0
  upload({
    origin: UPLOAD_ID,
    by: 'chunk',
    multiple: true,
    data: {
      directory,
      _uid: uniqueID,
      dir: dirFinal,
      user_name: currentUser.name
    },
    directory,
    beforeUpload: async files => {
      const resolvedFiles = []
      const rejectedFiles = []
      files.forEach(file => {
        const filePath = file.webkitRelativePath || file.name
        const fileName = filePath.split('/')[0]
        // check filename
        const { error } = Validator.filename(file.name)
        if (error) {
          message.error(error.message)
          return Promise.reject(error)
        }
        if (getDuplicate({ id: undefined, name: fileName, data: dir.data })) {
          rejectedFiles.push(file)
        } else {
          resolvedFiles.push(file)
        }
      })

      globalSizes[uniqueID] = 0
      files.forEach(file => {
        globalSizes[uniqueID] += file.size
      })

      if (rejectedFiles.length > 0) {
        if (directory) {
          const filePath =
            rejectedFiles[0].webkitRelativePath || rejectedFiles[0].name
          const topDirName = filePath.split('/')[0]

          const coverNodes = await showFailure({
            actionName: '上传',
            items: [
              {
                isFile: false,
                name: topDirName
              }
            ]
          })
          if (coverNodes.length > 0) {
            // remove dir
            await server.delete([`${dirFinal}/${topDirName}`])
            // should upload newFiles
            resolvedFiles.push(...rejectedFiles)
          }
        } else {
          const coverNodes = await showFailure({
            actionName: '上传',
            items: rejectedFiles.map(item => ({
              name: item.name,
              uid: item.uid,
              isFile: true
            }))
          })
          if (coverNodes.length > 0) {
            await server.delete(
              coverNodes.map(item => `${dirFinal}/${item.name}`)
            )
            resolvedFiles.push(...coverNodes)
          }
        }
      }

      if (resolvedFiles.length > 0) {
        uploadFilesLength = resolvedFiles.length
        message.success('文件开始上传')
        // 触发显示 dropdown
        EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, { visible: true })
      }

      return resolvedFiles.map(item => item.uid)
    },
    onChange: ({ file, origin }) => {
      if (origin !== UPLOAD_ID) {
        return
      }
      if (file.status === 'done') {
        uploadSucceedCount++
        // 有文件上传完成，check 是否要关闭 dropdown
        if (uploadSucceedCount === uploadFilesLength) {
          EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, { visible: false })
        }
        callback(dir.path || '.', dir.id)
      }
    }
  })
}

const onMove = (selectedNodes, callback) => {
  showDirSelector({
    disabledPaths: selectedNodes.map(item => item.path)
  }).then(path => {
    moveTo(path, selectedNodes, callback)
  })
}

const onShareFiles = (selectedNodes, actType) => {
  showShareFile({selectedNodes,actType})
}

function getCurrentDir(file_path) {
  if (!file_path) return '/'
  let last_slash_index = file_path.lastIndexOf('/')
  let current_dir = file_path.substring(0, last_slash_index).slice(0)
  return current_dir
}

async function moveTo(path, selectedNodes, callback) {
  const nodes = selectedNodes

  // check duplicate
  const targetDir = await server.fetch(path)
  const rejectedNodes = []
  const resolvedNodes = []

  nodes.forEach(item => {
    if (targetDir.getDuplicate({ id: undefined, name: item.name })) {
      rejectedNodes.push(item)
    } else {
      resolvedNodes.push(item)
    }
  })

  const destMoveNodes = [...resolvedNodes]
  if (rejectedNodes.length > 0) {
    const coverNodes = await showFailure({
      actionName: '移动',
      items: rejectedNodes
    })
    if (coverNodes.length > 0) {
      // coverNodes 中的要删除
      // await server.delete(coverNodes.map(item => `${path}/${item.name}`))
      destMoveNodes.push(...coverNodes)
    }
  }

  if (destMoveNodes.length > 0) {
    const srcPaths = destMoveNodes[0]?.selectedname
      ? destMoveNodes[0]?.selectedname
          .split(',')
          .map(name => `${getCurrentDir(destMoveNodes[0].path)}/${name}`)
      : [`${getCurrentDir(destMoveNodes[0].path)}/${destMoveNodes[0].name}`]

    const destPath = path ? path : '/'
    await server.move({ srcPaths, destPath, overwrite: true })

    await server.fetch(path)

    message.success('文件移动成功')
  }
  callback && callback()
}

defState.hist.push(defState.cdir)
defState.data = new Bin()
// defState.data.parse(fdata || '{}')

const handleData = data => {
  const collectData = {}
  data.forEach(item => {
    item.originId = item.id
    item.mtime = formatUnixTime(item.mtime)
    item.type = item.type
    item.size = item.size
    collectData[item.name] = item
  })
  return collectData
}

const fileReducer = (state = defState, action) => {
  let tmp = { ...state }
  let navHist = false

  // 生成第一层级的folder和file
  if (action.type === 'generateFiles') {
    // const HomeDirectory = currentUser.mountList[0]?.path || 'Home' + '>>'
    const HomeDirectory = 'Home' + '>>'
    fdata['Home:'].data = handleData(action.payload)
    tmp.data.parse(fdata)
  }

  // 第二层级的文件夹和文件的生成
  if (action.type === 'FILEDIR') {
    tmp.cdir = action.payload
    if (action?.data) {
      let itemById = tmp.data.getId(action.payload)
      const mapData = handleData(action.data)
      let drives = Object.keys(mapData)
      itemById.data = drives.map(item => {
        const data = mapData[item]
        return tmp.data.parseFolder(mapData[item], item, itemById)
      })
    }
    // tmp.data.parse(handleData(action.data), action.payload)
  } else if (action.type === 'CREATEFILEDIR') {
    // 新建文件夹
    let itemById = tmp.data.getId(tmp.cdir)
    // 判断当前目录下有多少个未命名文件夹，进行递增
    const findUnKnowFolders = itemById.data?.filter(item =>
      /^未命名文件夹\d?/.test(item?.name)
    )
    let name =
      findUnKnowFolders.length > 0
        ? `未命名文件夹${findUnKnowFolders.length}`
        : '未命名文件夹'
    // 如果递增的名字也存在于当前目录，修改为时间戳拼接
    const findResult = itemById.data?.find(item => item.name === name)
    if (findResult) {
      name = `未命名文件夹${Date.now()}`
    }
    const path = `${itemById.path ? itemById.path + '/' : ''}${name}`
    const createObj = {
      data: [],
      editFlag: false,
      host: null,
      id: Math.random().toString(36).substring(2, 10).toLowerCase(),
      info: { icon: 'folder' },
      name,
      path,
      isMkdir: true,
      _key: 'newFolder',
      size: 0,
      type: 'folder'
    }

    const createItem = tmp.data.parseFolder(createObj, name, itemById)
    createItem.editFlag = true
    createItem.isMkdir = true
    tmp.editing = true
    itemById.data.push(createItem)
  } else if (action.type === 'FILEOPERATE') {
    // 修改文件名称
    // 排他，将其他处于修改的关闭
    if (tmp.editing) {
      let currentFolder = tmp.data.getId(tmp.cdir)
      currentFolder.data.forEach(
        item => item.editFlag && (item.editFlag = false)
      )
      tmp.editing = false
    }
    if (action.payload === 'rename' && action.data) {
      // 新建文件夹
      let itemById = tmp.data.getId(tmp.cdir)
      // 判断当前目录下有多少个未命名文件夹，进行递增
      const findUnKnowFolders = itemById.data?.filter(item =>
        /^未命名文件夹\d?/.test(item?.name)
      )
      // 如果递增的名字也存在于当前目录，修改为时间戳拼接
      const findResult = itemById.data?.find(item => item?._key === 'newFolder')
    }
    // 将对应的name变成input
    if (action.data && action.payload !== 'handleRename') {
      let itemById = tmp.data.getId(action.data)
      itemById.editFlag = true
      tmp.editing = true
    }

    let currentFolder = tmp.data.getId(tmp.cdir)
    let targetItem = tmp.data.getId(action.data?.id)
    let itemById = tmp.data.getId(tmp.cdir)
    const newName =
      (currentFolder.path ? currentFolder.path + '/' : '/') +
      action.data?.targetPath
    const path = action.data?.originPath ? action.data?.originPath : '/'

    // 修改名称
    if (
      action.payload === 'handleRename' &&
      action.data?.targetPath !== action.data.name
    ) {
      const _copyTargetItem = targetItem
      const regex = /^[^\\,;'"`]+$/gi
      // 判断是否已经存在名称
      const findResult = currentFolder.data.find(
        item => item.name === action.data.targetPath
      )
      if (findResult) {
        message.error(`${action.data.targetPath}名称已存在！`)
      } else {
        try {
          if (regex.test(newName)) {
            if (action.data?.isMkdir) {
              const findAddingDir = currentFolder.data.find(
                item => item._key === 'newFolder'
              )
              if (newName.endsWith('/')) {
                message.warn('文件夹名称不能为空，请重新创建')
                itemById.editFlag = true
                tmp.editing = true
                return
              }
              server
                .mkdir(newName)
                .then(() => {
                  targetItem.isMkdir = false
                  itemById.editFlag = false
                  tmp.editing = false
                })
                .catch(e => {
                  itemById.editFlag = true
                  tmp.editing = true
                  message.error('操作失败')
                })
            } else {
              server.rename({ path, newName }).then(res => {
                if (res?.code === 0) {
                  message.success('操作成功')
                } else {
                  message.error('操作失败')
                  targetItem = _copyTargetItem
                }
              })
            }
            targetItem.name = action.data.targetPath
            targetItem.path = newName
          } else {
            message.warn(
              `名称不能包含反斜杠、逗号、分号、单引号、双引号和反引号`
            )
          }
        } catch (err) {
          message.error('操作失败', err)
          targetItem = _copyTargetItem
        }
      }
      // 是新建文件夹但是没有修改默认的文件夹名字
    } else {
      if (action.data?.isMkdir) {
        server
          .mkdir(newName)
          .then(() => {
            targetItem.isMkdir = false
            itemById.editFlag = false
            tmp.editing = false
          })
          .catch(e => {
            itemById.editFlag = true
            tmp.editing = true
            message.error('操作失败')
          })
      }
    }
  } else if (action.type === 'FILEREMOVE') {
    // 删除文件和文件夹
    // let itemById = tmp.data.getId(action.data.id)
    let currentItem = tmp.data.getId(tmp.cdir)
    let deletePath = [action.data.path]
    // 批量删除
    if (action.data?.selected.length > 0) {
      const deleteItems = currentItem.data?.filter(item =>
        action.data.selected?.includes(item.id)
      )
      deletePath = deleteItems.map(item => item.path)
    }
    // 调用删除接口
    server.delete(deletePath).then(res => {})
    currentItem.data = currentItem.data.filter(
      item => !action.data.selected.includes(item.id)
    )
  } else if (action.type === 'FILECOMPRESSION') {
    let currentDir = tmp.data.getId(tmp.cdir)
    let itemPaths = [action.data.path]
    let items = [action.data]
    // 多个文件压缩
    if (action.data?.selected.length > 0) {
      items = currentDir.data?.filter(item =>
        action.data.selected?.includes(item.id)
      )
      itemPaths = items.map(item => item.path)
    }
    // TODO 检查是否有压缩任务进行 新接口，如果有，提示用户，有压缩任务正在进行中，退出
    // 调用压缩接口
    _compress(itemPaths, currentDir, action.callback)
  } else if (action.type === 'UPLOAD') {
    const current = tmp.data.getId(tmp.cdir)
    _upload(action.payload, current, action.callback)
  } else if (action.type === 'DOWNLOAD') {
    let item = tmp.data.getId(action.payload)
    let path = [item.path]
    let isFile = [item.isFile]
    let size = [item.size]
    if (action.data) {
      path = action.data.map(item => item.path)
      isFile = action.data.map(item => item.isfile)
      size = action.data.map(item => item.size)
    }
    server.download(path, isFile, size)
  } else if (action.type === 'FIlESMOVE') {
    let item = tmp.data.getId(tmp.cdir)
    onMove(action.payload, () => action.callback(item.path || '/', item.id))
  } else if(action.type === 'FIlESSHARE') {
    let item = tmp.data.getId(tmp.cdir)
    const actType = action?.actType
    onShareFiles(action.payload, actType)
  } else if (action.type === 'FILEPATH') {
    let pathid = tmp.data.parsePath(action.payload)
    if (pathid) tmp.cdir = pathid
  } else if (action.type === 'FILEBACK') {
    let item = tmp.data.getId(tmp.cdir)
    if (item.host) {
      tmp.cdir = item.host.id
    }
  } else if (action.type === 'FILEVIEW') {
    tmp.view = action.payload
  } else if (action.type === 'FILEPREV') {
    tmp.hid--
    if (tmp.hid < 0) tmp.hid = 0
    navHist = true
  } else if (action.type === 'FILENEXT') {
    tmp.hid++
    if (tmp.hid > tmp.hist.length - 1) tmp.hid = tmp.hist.length - 1
    navHist = true
  }

  if (!navHist && tmp.cdir != tmp.hist[tmp.hid]) {
    tmp.hist.splice(tmp.hid + 1)
    tmp.hist.push(tmp.cdir)
    tmp.hid = tmp.hist.length - 1
  }

  tmp.cdir = tmp.hist[tmp.hid]
  if (tmp.cdir.includes('%')) {
    if (tmp.data.special[tmp.cdir] != null) {
      tmp.cdir = tmp.data.special[tmp.cdir]
      tmp[tmp.hid] = tmp.cdir
    }
  }

  tmp.cpath = tmp.data.getPath(tmp.cdir)
  return tmp
}

export default fileReducer
