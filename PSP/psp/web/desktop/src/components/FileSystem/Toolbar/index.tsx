import * as React from 'react'
import { observer } from 'mobx-react'
import { observable, action, computed, toJS } from 'mobx'
import { message } from 'antd'
import { Subject } from 'rxjs'

import {
  DeleteAction,
  ChooseFileAction,
  DownloadAction,
  CompressAction
} from '../Actions'
import { untilDestroyed } from '@/utils/operators'
import { Button, Modal } from '@/components'
import { RootPoint, FavoritePoint, PathHistory } from '@/domain/FileSystem'
import LocalUploader from './LocalUploader'
// import OpenRemote from './OpenRemote'
import { StyledToolbar } from './style'
import { eventEmitter, IEventData } from '@/utils'
import { ZIP_TYPE } from '@/utils/const'

export enum ActionType {
  upload = 'upload',
  download = 'download',
  newFolder = 'newFolder',
  edit = 'edit',
  delete = 'delete',
  moveTo = 'moveTo',
  copyTo = 'copyTo',
  rename = 'rename',
  compress = 'compress',
  decompress = 'decompress',
  customActions = 'customActions'
}

export interface IToolbarConfig {
  visible: boolean // show toolbar
  excludes: ActionType[] // exclude function
  includes: ActionType[] // include function
}

interface IProps {
  hasPerm?: boolean
  selectedKeys$: Subject<any>
  newDirectories$: any
  sourcePoint: RootPoint | FavoritePoint
  currentPoint: RootPoint
  points: RootPoint[]
  parentPath: string
  config: Partial<IToolbarConfig>
  showLoading?: any
  hideLoading?: any
  selectedMenuKeys: Object
  history: PathHistory
}

@observer
export default class Toolbar extends React.Component<IProps> {
  @observable selectedKeys: string[] = []
  selectedMenuKey: string

  @action
  updateKeys = keys => (this.selectedKeys = keys)

  config: IToolbarConfig = {
    visible: true,
    includes: [
      ActionType.upload,
      ActionType.download,
      ActionType.newFolder,
      ActionType.edit,
      ActionType.delete,
      ActionType.moveTo,
      ActionType.copyTo,
      ActionType.rename,
      ActionType.compress,
      ActionType.decompress,
      ActionType.customActions
    ],
    excludes: []
  }

  constructor(props) {
    super(props)

    // config toolbar
    if (props.config) {
      this.config = Object.assign({}, this.config, props.config)
    }
  }

  @computed
  get noneSelected() {
    const { selectedMenuKeys, points } = this.props
    //根据选中的文件夹在那个根目录下进行过滤，找到当前文件夹所在路径
    const point = points
      .map(point => selectedMenuKeys[point.pointId])
      .filter(n => n)[0]

    this.selectedMenuKey = point ? toJS(point)[0] : null
    const rootPath = points.map(point => point.rootPath)

    return (
      this.selectedKeys.length === 0 &&
      rootPath.filter(path => path === this.selectedMenuKey).length !== 0
    )
  }

  @computed
  get newDirectoryable() {
    const { sourcePoint } = this.props

    return sourcePoint && !(sourcePoint instanceof FavoritePoint)
  }

  get uploadable() {
    const { sourcePoint, parentPath } = this.props

    return parentPath && !(sourcePoint instanceof FavoritePoint)
  }

  get targets() {
    return this.selectedKeys.length !== 0
      ? this.selectedKeys
      : [this.selectedMenuKey]
  }

  componentDidMount() {
    this.props.selectedKeys$.pipe(untilDestroyed(this)).subscribe(keys => {
      this.updateKeys(keys)
    })
  }

  private newDirectory = () => {
    return this.props.newDirectories$
      .filter(n => n.pointId === this.props.currentPoint.pointId)[0]
      ['newDirectory$'].next()
  }

  private onDecompress = async (keys: string[], next, currPoint) => {
    await Modal.showConfirm({
      title: '警告',
      content:
        '系统不会检测同名文件，解压操作会强制覆盖目标目录中的同名文件，请谨慎操作？'
    })

    const { currentPoint } = this.props

    const dstpath = keys[0]
    const file = this.targets[0]
    const fileType = file.endsWith('.tar.gz')
      ? 'tar.gz'
      : file.substring(file.lastIndexOf('.') + 1)

    currentPoint.service
      .extract({
        file,
        toPath: dstpath,
        fileType: ZIP_TYPE[fileType]
      })
      .then(() => {
        // this.props.history.push({ source: currPoint, path: dstpath })
        // 监听消息
        eventEmitter.once(`DECOMPRESS_FILE_${file}`, (obj: IEventData) => {
          if (obj.message.success) {
            currentPoint.service.fetch(dstpath).then(() => {})
          }
        })
        message.success(`开始解压文件${file}, 请稍后...`)
      })

    next()
  }

  private onMove = async (keys: string[], next, currPoint) => {
    const { currentPoint, parentPath } = this.props

    //find parentPath
    let _parentPath =
      this.selectedKeys.length !== 0
        ? parentPath
        : this.selectedMenuKey.replace(/[\\\/][^\\\/]*$/, '')

    if (keys.length > 0) {
      const dstpath = keys[0]
      // can't move file to same directory
      if (dstpath === _parentPath) {
        message.info('不能在同一目录下移动文件或目录')
        return
      }

      // check exist files
      const dstFiles = this.selectedKeys.map(
        item => `${dstpath}/${item.split(/[\\/]/).pop()}`
      )
      const existArr = await currentPoint.service.exist(dstFiles)
      const conflictArr = dstFiles.filter((item, index) => existArr[index])
      let promise = Promise.resolve()
      if (conflictArr.length > 0) {
        promise = Modal.showConfirm({
          content: `下列文件已存在，是否覆盖：${conflictArr.join(', ')}`
        })
      }

      promise.then(() =>
        currentPoint.service
          .move({
            srcpaths: this.targets,
            dstpath,
            overwrite: true
          })
          .then(() => {
            this.props.history.push({ source: currPoint, path: dstpath })
            message.success('文件移动成功')
          })
      )
    }
    next()
  }

  private onConfirmDelete = async (promise, currPoint) => {
    const { showLoading, hideLoading } = this.props
    try {
      showLoading && showLoading('文件删除中...')
      await promise
      if (this.selectedKeys.length === 0) {
        this.props.history.push({
          source: currPoint,
          path: this.selectedMenuKey.replace(/[\\\/][^\\\/]*$/, '')
        })
      }
      eventEmitter.emit('FILE_SYSTEM_FILE_DELETE_EVNET')
    } finally {
      hideLoading && hideLoading()
    }
  }

  private hasFilePerm = () => {
    const { hasPerm } = this.props
    return !(typeof hasPerm !== undefined ? hasPerm : true)
  }

  private isCompressFile = fileName => {
    return ['.zip', '.tar', '.tar.gz', '.gz'].some(fix =>
      fileName.endsWith(fix)
    )
  }

  private canDecompress = () => {
    return this.selectedKeys.length === 1
      ? this.isCompressFile(this.selectedKeys[0])
      : false
  }

  private beforeDownload = async () => {
    return Modal.showConfirm({
      title: '下载确认弹窗',
      content: '确认要下载选中的文件吗？'
    })
  }

  render() {
    const { noneSelected, selectedKeys, newDirectoryable, uploadable, config } =
      this

    const { currentPoint: point, points, parentPath } = this.props

    const actions = config.includes.filter(
      item => !config.excludes.includes(item)
    )

    if (config.visible === false) {
      return null
    }

    return (
      <StyledToolbar>
        <div className='operator'>
          {actions.includes(ActionType.newFolder) ? (
            <Button
              type='primary'
              icon='folder-add'
              disabled={!newDirectoryable || this.hasFilePerm()}
              onClick={this.newDirectory}
              ghost>
              新建
            </Button>
          ) : null}
          {actions.includes(ActionType.download) ? (
            <DownloadAction
              point={point}
              targets={this.targets}
              beforeDownload={this.beforeDownload}>
              <Button
                disabled={noneSelected || this.hasFilePerm()}
                type='primary'
                icon='download'
                ghost>
                下载
              </Button>
            </DownloadAction>
          ) : null}
          {actions.includes(ActionType.upload) ? (
            <LocalUploader
              disabled={!uploadable || this.hasFilePerm()}
              point={point}
              path={parentPath}
              showLoading={this.props.showLoading}
              hideLoading={this.props.hideLoading}
            />
          ) : null}
          {actions.includes(ActionType.moveTo) ? (
            <ChooseFileAction
              points={points}
              disabledKeys={selectedKeys}
              path={parentPath}
              onOk={this.onMove}
              hasPerm={this.props.hasPerm}
              onCancel={(keys, next) => next()}>
              <Button
                disabled={noneSelected || this.hasFilePerm()}
                type='primary'
                icon='move'
                ghost>
                移动
              </Button>
            </ChooseFileAction>
          ) : null}
          {actions.includes(ActionType.delete) ? (
            <DeleteAction
              point={point}
              targets={this.targets}
              onConfirm={this.onConfirmDelete}>
              <Button
                disabled={noneSelected || this.hasFilePerm()}
                type='primary'
                icon='delete'
                ghost>
                删除
              </Button>
            </DeleteAction>
          ) : null}
          {actions.includes(ActionType.compress) ? (
            <CompressAction
              points={points}
              point={point}
              targets={this.targets}
              onClick={() => {}}>
              <Button
                disabled={noneSelected || this.hasFilePerm()}
                type='primary'
                icon='yasuo'
                ghost>
                压缩
              </Button>
            </CompressAction>
          ) : null}

          {actions.includes(ActionType.decompress) ? (
            <ChooseFileAction
              points={points}
              disabledKeys={selectedKeys}
              path={parentPath}
              onOk={this.onDecompress}
              hasPerm={this.props.hasPerm}
              onCancel={(keys, next) => next()}>
              <Button
                disabled={!this.canDecompress() || this.hasFilePerm()}
                type='primary'
                icon='jieyasuo'
                ghost>
                解压缩
              </Button>
            </ChooseFileAction>
          ) : null}

          {/* <OpenRemote
            selectedKeys={selectedKeys}
            point={point}
            parentPath={parentPath}
          /> */}
        </div>
      </StyledToolbar>
    )
  }
}
