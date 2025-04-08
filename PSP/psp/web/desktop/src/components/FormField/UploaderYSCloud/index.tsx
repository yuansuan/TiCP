/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import {
  action,
  computed,
  observable,
  runInAction,
  transaction,
  when,
} from 'mobx'
import { observer } from 'mobx-react'
import * as React from 'react'

import { Button } from '@/components'
import Uploader, { Task } from '@/domain/Uploader'
import { createMobxStream, Http } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import Container from '../Container'
import LocalUploader from './LocalUploader'
import ServerUploader from './ServerUploader'
import VirtualDirectory from './VirtualDirectory'
import { FileList } from './FileList'

interface IProps {
  model
  formModel: any
  showId?: boolean
  fetchUploadPath: () => Promise<any>
  win?: any
}

interface IFileExtendSion {
  _isMain: boolean
  _master: string
  _from: 'local' | 'server'
  content?: string
  slaveFiles?: Map<string, IFile>
}

interface IDoneFile {
  path: string
  name: string
  size: number
  status: string
  percent: number
  is_dir: boolean
  _isMain: boolean
  _master: string
  _from: 'local' | 'server'
  // cache valid slave files from master
  _validSlaveFiles?: string[]
  slaveFiles?: Map<string, IFile>
  content?: string
}

type IFile = ((Task | VirtualDirectory) & IFileExtendSion) | IDoneFile

export interface IHarmonyFile {
  path: string
  name: string
  size: number
  status: string
  // percent: number
  isDir: boolean
  isMain: boolean
  master: string
  from: 'local' | 'server'
  slaveFiles?: IHarmonyFile[]
}

@observer
export default class LocalUploaderItem extends React.Component<IProps> {
  // harmony Task/VirtualDirectory/Server File/Server Directory
  private harmony = (file): IHarmonyFile => {
    // Task
    if (file instanceof Task) {
      return {
        path: file.target.path,
        name: file.target.name,
        size: file.target.size,
        status: file.status,
        // percent: file.percent,
        isMain: false,
        isDir: false,
        master: file.target._master,
        from: file.target._from,
      }
    } else if (file instanceof VirtualDirectory) {
      // VirtualDirectory
      return {
        path: file.path,
        name: file.name,
        size: file.size,
        status: file.status,
        // percent: file.percent,
        isMain: false,
        isDir: file.is_dir,
        master: file._master,
        from: file._from,
      }
    } else {
      // done file/directory
      const size = Number.isNaN(parseFloat(file.size)) ? '--' : file.size

      return {
        path: file.path,
        name: file.name,
        size,
        status: 'done',
        // percent: 100,
        isMain: !!file._isMain,
        isDir: file.is_dir || file.mode[0] === 'd',
        master: file._master,
        from: file._from,
        // recurse harmony slaveFiles
        slaveFiles: file.slaveFiles
          ? [...file.slaveFiles.values()].map(this.harmony)
          : undefined,
      }
    }
  }
  // harmonyFiles just be used in FileTree
  @computed
  get harmonyFiles(): IHarmonyFile[] {
    return [...this.files.values()].map(this.harmony)
  }

  @observable public files: Map<string, IFile> = new Map()

  @action
  setMainFileKeysAction(keys: string[]) {
    const { formModel, model } = this.props
    formModel[model.id]['mainFileKeys'] = keys
  }

  constructor(props) {
    super(props)

    const { formModel, model } = props
    formModel[model.id] = {
      ...model,
      value: model.value || model.defaultValue,
      values: model.values.length > 0 ? model.values : model.defaultValues,
      masterSlave: model.masterSlave || '',
      // hack: let form to know the internal files
      _files: [],
    }
  }

  @action
  public uploadFile = file => {
    // 上传后可能已覆盖主文件，由于文件树非递归一次性展开，无法确认是否覆盖，统一取消主文件选择
    this.setMainFileKeysAction([])
    this.files.set(this.getFilePath(file), file)
  }

  @action
  public deleteFile = filePath => {
    this.files.delete(filePath)
  }

  async componentDidMount() {
    createMobxStream(() => this.harmonyFiles, false)
      .pipe(untilDestroyed(this))
      .subscribe(files => {
        const { formModel, model } = this.props
        const field = formModel[model.id]
        field._files = files
        field.values = files.map(item => item.path)
      })
  }

  public render() {
    const { model, showId } = this.props
    const uploadTypes = model.fileFromType.split('_').filter(item => !!item)
    return (
      <>
        <Container model={model} showId={showId}>
          {uploadTypes.map(type =>
            type === 'local' ? (
              <LocalUploader key={type} upload={this.uploadLocalFile}>
                <Button icon='upload' style={{ marginRight: '5px' }}>
                  本地文件
                </Button>
              </LocalUploader>
            ) : (
              <ServerUploader key={type} onUpload={this.uploadServerFile}>
                {upload => (
                  <Button icon='upload' onClick={upload}>
                    远程文件
                  </Button>
                )}
              </ServerUploader>
            )
          )}
        </Container>
        <div>
          {this.harmonyFiles.length > 0 ? (
            <FileList
              fileList={this.harmonyFiles}
              deleteAction={this.deleteAction}
              setMainFileKeysAction={(keys: string[]) => {
                this.setMainFileKeysAction(keys)
              }}
            />
          ) : null}
        </div>
      </>
    )
  }

  private getFilePath = file => {
    if (file instanceof Task) {
      return file.target.path
    } else {
      return file.path
    }
  }

  @action
  private uploadLocalFile = async (params, isDir) => {
    params.data.dirPath = await this.props.fetchUploadPath()
    const task = Uploader.upload(params, isDir)

    // upload local directory
    const { customPath, uploadPath } = task.target
    if (/[\\/]/.test(customPath)) {
      const dirPath = customPath.match(/^[^\\/]*/)[0]

      let directory
      directory = this.files.get(dirPath)
      if (!directory) {
        directory = new VirtualDirectory({
          path: dirPath,
          name: dirPath,
          _from: 'local',
        })

        // upload directory to master's slaveFiles
        this.uploadFile(directory)

        // when upload complete, replace the local directory with server directory
        const disposer = when(
          () => directory.isDone,
          () => {
            Http.get('/file/detail', {
              params: {
                paths: uploadPath.replace(/[\\/][^\\/]*$/, ''),
              },
            }).then(res => {
              // replace local dirctory with server directory
              runInAction(() => {
                // delete local directory
                this.deleteFile(directory.path)
                // upload server directory
                this.uploadFile({
                  ...res.data.files[0],
                  _from: 'local',
                })
              })
            })
          }
        )
        // if the uploadingDirectory is aborted, dispose the when monitor
        directory.hooks.aborted.tap('dispose monitor', () => {
          disposer && disposer()
        })
      }
      directory.addTask(task)
    } else {
      // upload local file
      // set master/from
      task.target._from = 'local'

      // upload file to master's slaveFiles
      this.uploadFile(task)

      // when upload complete, replace the local file with server file
      const subscription = task.status$.subscribe(status => {
        if (['done', 'error', 'aborted'].includes(status)) {
          subscription.unsubscribe()
        }

        if (status === 'done') {
          Http.get('/file/detail', {
            params: {
              paths: task.target.uploadPath,
            },
          }).then(res => {
            // replace local file with server file
            runInAction(() => {
              // delete local file
              this.deleteFile(task.target.path)
              // upload server file
              this.uploadFile({
                ...res.data.files[0],
                _from: 'local',
              })
            })
          })
        }
      })
    }
  }

  @action
  private uploadServerFile = files => {
    transaction(() => {
      files.forEach(file => {
        this.uploadFile({ ...file, _from: 'server' })
      })
    })
  }

  @action
  private deleteAction = path => {
    const targetFile = this.files.get(path)

    // local-upload file/directory
    if (targetFile instanceof Task || targetFile instanceof VirtualDirectory) {
      targetFile.abort()
    }

    // delete file from files
    this.deleteFile(path)
  }
}
