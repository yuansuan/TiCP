/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { message, Spin } from 'antd'
import {
  action,
  computed,
  observable,
  runInAction,
  transaction,
  when
} from 'mobx'
import { observer } from 'mobx-react'
import * as React from 'react'
import { currentUser } from '@/domain'
import { Button, Modal } from '@/components'
import Uploader, { Task } from '@/domain/Uploader'
import { createMobxStream, Http, formatRegExpStr } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import Container from '../Container'
import Editor from './Editor'
import FileTree from './FileTree'
import { FileList } from './FileList'
import LocalUploader from './LocalUploader'
import ServerUploader from './ServerUploader'
import VirtualDirectory from './VirtualDirectory'
import SysConfig from '@/domain/SysConfig'
import { currentFileList } from './FileList/Files'
import eventEmitter from '@/utils/EventEmitter'
import WorkdirAction from './WorkdirAction'
import debounce from 'lodash.debounce'

const workdirTextStyle = {
  padding: 10,
  lineHeight: '20px',
  fontSize: 14,
  fontWeight: 600,
  background: '#F3F5F8',
  borderBottom: '1px solid #ccc'
}

interface IProps {
  model
  formModel: any
  showId?: boolean
  fetchUploadPath: (isReSubmitUseSameJobDir?) => Promise<any>
  win?: any
  isResubmit?: boolean
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
  is_text?: boolean
}

let topFiles = {}
// path: []
@observer
export default class LocalUploaderItem extends React.Component<IProps> {
  @observable mainFilePath = null
  @observable jobDirPath = null
  @observable fileLoading = false
  @observable fileLoadingText = '数据加载中'
  @observable workdir = '' // 610定制需求，自定义workdir
  @observable rootPath = ''
  public static Editor = Editor

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
        is_text: file.target.is_text
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
        is_text: false
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
        isDir: file.is_dir,
        master: file._master,
        from: file._from,
        // recurse harmony slaveFiles
        slaveFiles: file.slaveFiles
          ? [...file.slaveFiles.values()].map(this.harmony)
          : undefined,
        is_text: file?.is_text
      }
    }
  }
  // harmonyFiles just be used in FileTree
  @computed
  get harmonyFiles(): IHarmonyFile[] {
    return [...this.files.values()].map(this.harmony)
  }

  @computed
  get _files() {
    return [...currentFileList.files]
  }

  @observable public files: Map<string, IFile> = new Map()
  @observable public expandedKeys: string[] = []

  constructor(props) {
    super(props)

    const { formModel, model } = props
    formModel[model.id] = {
      ...model,
      value: model.value || model.defaultValue,
      values: model.values.length > 0 ? model.values : model.defaultValues,
      masterSlave: model.masterSlave || '',
      // hack: let form to know the internal files
      _files: []
    }
  }

  @action
  public updateFiles = files => (this.files = files)

  @action
  public uploadFile = (file, master?) => {
    const masterFile = this.getMasterFile(master)
    if (masterFile) {
      masterFile.slaveFiles = masterFile.slaveFiles || new Map()
      masterFile.slaveFiles.set(this.getFilePath(file), file)
    } else {
      this.files.set(this.getFilePath(file), file)
    }
  }
  @action
  public deleteFile = (filePath, master?) => {
    const masterFile = this.getMasterFile(master)
    if (masterFile) {
      masterFile.slaveFiles.delete(filePath)
    } else {
      this.files.delete(filePath)
    }
  }

  async componentDidMount() {
    // resubmit/draft
    const { formModel, model, isResubmit } = this.props
    const { isSupportWorkdir, workdir } = formModel[model.id]
    const { homedir } = SysConfig.userConfig
    if (homedir) this.rootPath = homedir + '/' + currentUser?.name

    if (isSupportWorkdir) this.workdir = workdir

    if (model.isSupportMaster) {
      this.fileLoading = true
      try {
        const { values: files, masterFile } = formModel[model.id]

        let tmpPath = ''

        if (
          isResubmit &&
          (SysConfig.jobConfig?.job?.resubmit?.use_identical_job_dir ||
            isSupportWorkdir)
        ) {
          if (!files?.length) return
          // 主文件模式，使用原目录
          tmpPath = await this.props.fetchUploadPath(true)
        } else {
          // 主文件模式，copy 文件进入，临时作业目录
          if (!files?.length) return
          // create tmp dir
          tmpPath = await this.props.fetchUploadPath()

          // copy
          await Http.post('/file/copy', {
            srcpaths: files,
            dstpath: tmpPath
          })
        }
        if (tmpPath) {
          // list
          const { data } = await Http.get('/file/ls', {
            params: { path: tmpPath }
          })

          transaction(() => {
            data.files
              .filter(f => !!f?.path)
              .forEach(f => {
                this.uploadFile({
                  ...f,
                  _isMain: null,
                  _from: 'local'
                })
              })
          })
        }

        if (masterFile) {
          let tmpMainPath = masterFile.split(/\.tmp_[^\/]+/)?.[1]

          if (!tmpMainPath) {
            tmpMainPath = masterFile.split(tmpPath)?.[1]
          }

          if (!tmpMainPath) {
            tmpMainPath = masterFile.split(workdir)?.[1]
          }

          this.jobDirPath = tmpPath
          this.mainFilePath = `${tmpPath}${tmpMainPath}`
        }
      } catch (e) {
        message.error('作业数据文件失效, 请重新上传')
      } finally {
        this.fileLoading = false
      }
    } else {
      const { masterSlave } = formModel[model.id]

      if (masterSlave) {
        const relations = JSON.parse(masterSlave)
        const jobFiles = new Set()
        // flatten job files
        Object.keys(relations).forEach(key => {
          const slaves = relations[key]
          jobFiles.add(key)
          slaves && slaves.forEach(item => jobFiles.add(item))
        })
        if (jobFiles.size > 0) {
          // fetch files and upload files
          Http.get('/file/detail', {
            params: { paths: [...jobFiles].join(',') },
            disableErrorMessage: true
          })
            .then(res => {
              const {
                data: { files }
              } = res
              Object.keys(relations).forEach(key => {
                const slaves = relations[key]
                const masterFile = files.find(file => file.path === key)
                if (masterFile) {
                  this.uploadFile({
                    ...masterFile,
                    _isMain: slaves !== null,
                    _from: 'server'
                  })
                  slaves &&
                    slaves.forEach(item => {
                      const slaveFile = files.find(file => file.path === item)
                      this.uploadFile(
                        {
                          ...slaveFile,
                          _from: 'server',
                          _master: key
                        },
                        key
                      )
                    })
                }
              })
            })
            .catch(() => {
              message.error('作业数据文件失效, 请重新上传')
            })
        }
      }

      createMobxStream(() => this.harmonyFiles, false)
        .pipe(untilDestroyed(this))
        .subscribe(files => {
          const { formModel, model } = this.props
          if (model.isSupportMaster) return
          const field = formModel[model.id]
          field._files = files
          // only upload file which is't slave file
          field.values = files
            .filter(item => !item.master)
            .map(item => item.path)
          // update file relations
          const relations = files.reduce((tree, file) => {
            tree[file.path] = file.isMain
              ? (file.slaveFiles || []).map(item => item.path)
              : null
            return tree
          }, {})
          field.masterSlave = JSON.stringify(relations)
        })
    }

    createMobxStream(() => this._files, false)
      .pipe(untilDestroyed(this))
      .subscribe(files => {
        this.onFilesChange(files)
      })
  }

  onFilesChange(files) {
    const { formModel, model } = this.props
    if (model.isSupportMaster) {
      const field = formModel[model.id]
      field._files = files

      let values = []

      const visitServerFiles = serverFiles => {
        serverFiles.forEach(file => {
          if (file?.children?.length) {
            visitServerFiles(file.children)
          } else {
            values.push(file.path)
          }
        })
      }

      files.forEach(f => {
        // local 只包含顶层文件夹，因为删除为 物理删除
        if (f.from == 'local') {
          values.push(f.path)
        } else if (f.from == 'server') {
          if (f?.children?.length) {
            // 递归获取，如果用户展开文件夹，
            // 可能存在删除文件或文件夹的可能，这里的删除是逻辑删除
            visitServerFiles(f.children)
          } else {
            // 没有 children，说明包含整个文件夹里的所有文件
            values.push(f.path)
          }
        }
      })

      field.values = values
      // update file relations
      field.masterSlave = ''
    }
  }

  putIntoData = debounce(() => {
    if (!topFiles[this.workdir] || topFiles[this.workdir].length === 0) return
    for (let i = 0; i < 5; i++) {
      let f = topFiles[this.workdir][i]
      if (f?.path) {
        this.files.set(f.path, {
          ...f,
          _isMain: null,
          _from: 'local'
        })
      }
    }

    topFiles[this.workdir].splice(0, 5)
  }, 300)

  onScroll = e => {
    if (
      e.target.clientHeight + e.target.scrollTop + 50 >=
      e.target.scrollHeight
    ) {
      this.putIntoData()
    }
  }

  onWorkdirChange = async workdir => {
    this.fileLoading = true
    // list
    try {
      const { data } = await Http.get('/file/ls', {
        params: { path: workdir }
      })

      topFiles[workdir] = data.files || []

      let files = new Map()

      for (let i = 0; i < 2000; i++) {
        let f = topFiles[workdir][i]
        if (f?.path) {
          files.set(f.path, {
            ...f,
            _isMain: null,
            _from: 'local'
          })
        }
      }

      topFiles[workdir].splice(0, 2000)

      this.updateFiles(files)

      this.workdir = workdir
      const { formModel, model } = this.props
      const field = formModel[model.id]
      field.workdir = workdir
    } catch (e) {
      message.error('工作目录切换后，文件数据加载失败')
    } finally {
      this.fileLoading = false
    }
  }

  public render() {
    const { model, showId, isResubmit } = this.props
    const { isSupportMaster, isSupportWorkdir } = model

    const uploadTypes = model.fileFromType.split('_').filter(item => !!item)

    return (
      <>
        <Spin tip={this.fileLoadingText} spinning={this.fileLoading}>
          <Container model={model} showId={showId}>
            {isSupportWorkdir && (
              <WorkdirAction
                disabled={isResubmit}
                onSelect={this.onWorkdirChange}
                workdir={this.workdir}
                rootPath={this.rootPath}>
                <Button disabled={isResubmit} type='link'>
                  选择工作目录
                </Button>
              </WorkdirAction>
            )}
            {uploadTypes.map(type =>
              type === 'local' ? (
                <LocalUploader
                  disabled={isSupportWorkdir ? !this.workdir : false}
                  key={type}
                  upload={this.uploadLocalFile}
                  beforeUpload={this.beforeUpload}>
                  <Button icon='upload' style={{ marginRight: '5px' }}>
                    本地文件
                  </Button>
                </LocalUploader>
              ) : (
                <ServerUploader key={type} onUpload={this.uploadServerFile}>
                  {upload => (
                    <Button
                      disabled={
                        (isSupportWorkdir ? !this.workdir : false) ||
                        this.fileLoading
                      }
                      icon='upload'
                      onClick={upload}>
                      远程文件
                    </Button>
                  )}
                </ServerUploader>
              )
            )}
          </Container>
          <div>
            {isSupportWorkdir && this.workdir && (
              <div style={workdirTextStyle}>当前工作目录: {this.workdir}</div>
            )}
            {this.harmonyFiles.length > 0 ? (
              isSupportMaster ? (
                <FileList
                  onScroll={e => {
                    this.onScroll(e)
                  }}
                  expandedKeys={this.expandedKeys}
                  onExpandedKeysChange={keys => (this.expandedKeys = keys)}
                  fileList={this.harmonyFiles}
                  deleteAction={this.deleteAction}
                  onFileListChange={this.onFilesChange.bind(this)}
                  jobDirPath={this.jobDirPath}
                  mainFilePath={this.mainFilePath}
                  setMainFileKeysAction={(keys: any[]) => {
                    const { formModel, model } = this.props
                    if (model.isSupportMaster) {
                      const field = formModel[model.id]
                      field.masterFile = keys[0]?.path
                    }
                  }}
                  beforeUploadLocalFile={this.beforeUploadLocalSubFile}
                  uploadLocalFile={this.uploadLocalSubFile}
                  uploadServerFile={this.uploadServerFile}
                />
              ) : (
                <FileTree
                  model={model}
                  files={this.harmonyFiles}
                  deleteAction={this.deleteAction}
                  setMain={this.setMainAction}
                  beforeUploadLocalFile={this.beforeUploadLocalFile}
                  uploadLocalFile={this.uploadLocalFile}
                  uploadServerFile={this.uploadServerFile}
                />
              )
            ) : null}
          </div>
        </Spin>
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

  private getMasterFile = (masterPath): IDoneFile => {
    if (!masterPath) {
      return null
    }

    const masterFile = this.files.get(masterPath)

    if (!masterFile) {
      throw new Error(`master file: ${masterPath} is not exist`)
    }

    return masterFile as IDoneFile
  }

  private filterSlaveFilesFromMaster = async masterPath => {
    if (!masterPath) {
      return Promise.resolve()
    }

    const masterFile = this.getMasterFile(masterPath)
    // if master file is larger than 5M, skip filter
    if (masterFile.size > 5 * 1024 * 1024) {
      return Promise.resolve()
    }
    // serverFile: fetch masterFile's content
    if (masterFile.content === undefined) {
      const content = await Http.get('/file/content', {
        params: {
          path: masterFile.path,
          offset: 0,
          len: masterFile.size
        }
      }).then(res => res.data.content)

      masterFile.content = content
    }

    // filter slave files by keywords
    if (masterFile._validSlaveFiles === undefined) {
      const { formModel, model } = this.props
      const field = formModel[model.id]
      const { masterIncludeKeywords } = field
      const keywords = masterIncludeKeywords.split(';').filter(item => !!item)
      const prefix = `^(${formatRegExpStr(
        keywords.join('|')
      )})(\\r\\n|\\n|\\u0020)`
      const validReg = new RegExp(`${prefix}.+$`, 'gm')
      const validSlaveFiles = (masterFile.content.match(validReg) || []).map(
        item => item.replace(new RegExp(prefix), '')
      )
      masterFile._validSlaveFiles = validSlaveFiles
    }

    return Promise.resolve(masterFile._validSlaveFiles)
  }

  private beforeUpload = params => {
    let isShowConfirm = false
    let fileName = params.file.name
    let needDeleteFile = null

    if (params.isDir) {
      fileName = params.file.webkitRelativePath.split('/')?.[0]
    }

    this.files.forEach(value => {
      if (value?.['name'] === fileName) {
        isShowConfirm = true
        needDeleteFile = value
      }
    })

    if (isShowConfirm) {
      if (params.isDir) {
        if (params.isLast) {
          return Modal.showConfirm({
            title: '确认',
            content: `是否覆盖文件或文件夹${fileName}`
          })
            .then(() => {
              return Http.post(
                '/file/delete',
                { paths: [needDeleteFile.path] },
                { formatErrorMessage: msg => `删除失败：${msg}` }
              )
                .then(() => {
                  this.deleteFile(needDeleteFile.path, needDeleteFile._master)
                  return Promise.resolve()
                })
                .catch(e => {
                  console.error(`删除失败：${e}`)
                  return Promise.reject()
                })
            })
            .catch(() => {
              return Promise.reject()
            })
        }
      } else {
        return Modal.showConfirm({
          title: '确认',
          content: `是否覆盖文件或文件夹${fileName}`
        })
          .then(() => {
            return Http.post(
              '/file/delete',
              { paths: [needDeleteFile.path] },
              { formatErrorMessage: msg => `删除失败：${msg}` }
            )
              .then(() => {
                this.deleteFile(needDeleteFile.path, needDeleteFile._master)
                return Promise.resolve()
              })
              .catch(e => {
                console.error(`删除失败：${e}`)
                return Promise.reject()
              })
          })
          .catch(() => {
            return Promise.reject()
          })
      }
    }

    return Promise.resolve()
  }

  @action
  private beforeUploadLocalSubFile = async params => {
    const master = params.data ? params.data.master : undefined

    let isShowConfirm = false
    let fileName = params.file.name
    let needDeleteFile = null

    if (params.isDir) {
      fileName = params.file.webkitRelativePath.split('/')?.[0]
    }

    return Http.get('/file/list', {
      params: { path: master }
    })
      .then(res => {
        res.data.files.forEach(file => {
          if (file?.['name'] === fileName) {
            isShowConfirm = true
            needDeleteFile = file
          }
        })

        if (isShowConfirm) {
          if (params.isDir) {
            if (params.isLast) {
              return Modal.showConfirm({
                title: '确认',
                content: `是否覆盖文件或文件夹${fileName}`
              })
                .then(() => {
                  return Http.post(
                    '/file/delete',
                    { paths: [needDeleteFile.path] },
                    { formatErrorMessage: msg => `删除失败：${msg}` }
                  )
                    .then(() => {
                      this.deleteFile(needDeleteFile.path, null)
                      return Promise.resolve()
                    })
                    .catch(e => {
                      console.error(`删除失败：${e}`)
                      return Promise.reject()
                    })
                })
                .catch(() => {
                  return Promise.reject()
                })
            }
          } else {
            return Modal.showConfirm({
              title: '确认',
              content: `是否覆盖文件或文件夹${fileName}`
            })
              .then(() => {
                return Http.post(
                  '/file/delete',
                  { paths: [needDeleteFile.path] },
                  { formatErrorMessage: msg => `删除失败：${msg}` }
                )
                  .then(() => {
                    this.deleteFile(needDeleteFile.path, null)
                    return Promise.resolve()
                  })
                  .catch(e => {
                    console.error(`删除失败：${e}`)
                    return Promise.reject()
                  })
              })
              .catch(() => {
                return Promise.reject()
              })
          }
        }

        return Promise.resolve()
      })
      .catch(e => {
        return Promise.reject()
      })
  }

  @action
  private uploadLocalSubFile = async (params, isDir) => {
    // dirPath
    const master = params.data ? params.data.master : undefined

    if (
      !this.props.isResubmit &&
      this.props.model.isSupportWorkdir &&
      this.workdir
    ) {
      params.data.dirPath = this.workdir
    } else {
      params.data.dirPath = master ? master : await this.props.fetchUploadPath()
    }

    const task = Uploader.upload(params, isDir)

    // upload local directory
    const { customPath } = task.target
    if (/[\\/]/.test(customPath)) {
      const dirPath = customPath.match(/^[^\\/]*/)[0]

      let directory = new VirtualDirectory({
        path: dirPath,
        name: dirPath,
        _master: master,
        _from: 'local'
      })

      // when upload complete, replace the local directory with server directory
      const disposer = when(
        () => directory.isDone,
        () => {
          eventEmitter.emitEvent('AFTER_UPLOAD_REFRESH')
        }
      )
      // if the uploadingDirectory is aborted, dispose the when monitor
      directory.hooks.aborted.tap('dispose monitor', () => {
        disposer && disposer()
      })
      directory.addTask(task)
    } else {
      // upload local file
      // set master/from
      task.target._master = master
      task.target._from = 'local'

      // when upload complete, replace the local file with server file
      const subscription = task.status$.subscribe(status => {
        if (['done', 'error', 'aborted'].includes(status)) {
          eventEmitter.emitEvent('AFTER_UPLOAD_REFRESH')
          subscription.unsubscribe()
        }
      })
    }
  }

  @action
  private beforeUploadLocalFile = async params => {
    // master flag
    const master = params.data ? params.data.master : undefined
    // filter slave files
    if (master) {
      const filePath = params.file.webkitRelativePath || params.file.name
      const validSlaveFiles = await this.filterSlaveFilesFromMaster(master)
      if (!validSlaveFiles || validSlaveFiles.length === 0) {
        message.error('当前主文件无从文件依赖')
        return Promise.reject()
      }

      const isDir = /[\\/]/.test(filePath)
      const localPath = `./${filePath}`

      const validateExtensions = path => {
        // validate file extension
        const { formModel, model } = this.props
        const field = formModel[model.id]
        const { masterIncludeExtensions } = field
        const extensions = masterIncludeExtensions
          .split(/\s*;\s*/)
          .filter(item => !!item)
        const ext = path.includes('.') ? path.split('.').pop() : ''
        if (!extensions.includes(ext)) {
          message.error(`${filePath} 上传失败（不支持 ${ext} 格式的文件依赖）`)
          return false
        }

        return true
      }

      // validate directory import
      if (isDir) {
        let dirPath = localPath
        while (dirPath.search(/[\\/]/) > -1) {
          dirPath = dirPath.replace(/[\\/][^\\/]*$/, '')
          if (validSlaveFiles.includes(dirPath)) {
            // validate file extension
            if (!validateExtensions(filePath)) {
              return Promise.reject()
            }
            return Promise.resolve()
          }
        }
      }

      // validate file import
      if (!validSlaveFiles.includes(localPath)) {
        // validate declaration
        message.error(`${filePath} 上传失败（未声明的文件依赖）`)
        return Promise.reject()
      } else {
        // validate file extension
        if (!validateExtensions(filePath)) {
          return Promise.reject()
        }
      }
    }

    return Promise.resolve()
  }

  @action
  private uploadLocalFile = async (params, isDir) => {
    // master flag
    const master = params.data ? params.data.master : undefined
    if (
      !this.props.isResubmit &&
      this.props.model.isSupportWorkdir &&
      this.workdir
    ) {
      params.data.dirPath = this.workdir
    } else {
      params.data.dirPath = await this.props.fetchUploadPath()
    }
    const task = Uploader.upload(params, isDir)

    // upload local directory
    const { customPath } = task.target
    if (/[\\/]/.test(customPath)) {
      const dirPath = customPath.match(/^[^\\/]*/)[0]

      let directory
      directory = this.files.get(dirPath)
      if (!directory) {
        directory = new VirtualDirectory({
          path: dirPath,
          name: dirPath,
          _master: master,
          _from: 'local'
        })

        // upload directory to master's slaveFiles
        this.uploadFile(directory, master)

        // when upload complete, replace the local directory with server directory
        const disposer = when(
          () => directory.isDone,
          () => {
            Http.get('/file/detail', {
              params: {
                paths: params?.data?.dirPath + '/' + dirPath
              }
            }).then(res => {
              // replace local dirctory with server directory
              runInAction(() => {
                // delete local directory
                this.deleteFile(directory.path, master)
                // upload server directory
                this.uploadFile(
                  {
                    ...res.data.files[0],
                    _master: master,
                    _from: 'local'
                  },
                  master
                )
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
      task.target._master = master
      task.target._from = 'local'

      // upload file to master's slaveFiles
      this.uploadFile(task, master)

      // when upload complete, replace the local file with server file
      const subscription = task.status$.subscribe(status => {
        if (['done', 'error', 'aborted'].includes(status)) {
          subscription.unsubscribe()
        }

        if (status === 'done') {
          Http.get('/file/detail', {
            params: {
              paths: task.target.uploadPath
            }
          }).then(res => {
            // replace local file with server file
            runInAction(() => {
              // delete local file
              this.deleteFile(task.target.path, master)
              // upload server file
              this.uploadFile(
                {
                  ...res.data.files[0],
                  _master: master,
                  _from: 'local'
                },
                master
              )
            })
          })
        }
      })
    }
  }

  @action
  private uploadServerFile = async (files, master?, targetPath?) => {
    const { model } = this.props

    if (model.isSupportMaster) {
      // 新主文件模式，copy 文件到upload目录下
      this.fileLoading = true

      let path = null

      if (!this.props.isResubmit && model.isSupportWorkdir && this.workdir) {
        path = targetPath ? targetPath : this.workdir
      } else {
        path = targetPath ? targetPath : await this.props.fetchUploadPath()
      }
      try {
        // 检测同名文件是否需要覆盖

        const { data } = await Http.get('/file/list', {
          params: { path: path }
        })

        const { files: fileList } = data

        let fileNames = files.map(file => file.name)

        let confirmFileList = []

        fileList.forEach(value => {
          if (fileNames.includes(value?.['name'])) {
            confirmFileList.push(value)
          }
        })

        if (confirmFileList.length > 0) {
          await Modal.showConfirm({
            title: '确认',
            content: `是否覆盖文件或文件夹${confirmFileList
              .map(f => f.name)
              .join(',')}`
          })

          await Promise.all(
            confirmFileList.map(async f => {
              this.deleteFile(f.path, f._master)
              return Http.post(
                '/file/delete',
                { paths: [f.path] },
                { formatErrorMessage: msg => `删除失败：${msg}` }
              )
            })
          )
        }

        // copy 操作
        await Http.post('/file/copy', {
          srcpaths: files.map(f => f.path),
          dstpath: path
        })

        if (targetPath) {
          eventEmitter.emitEvent('AFTER_SERVER_UPLOAD_REFRESH')
        } else {
          // 因为文件已近被 copy 到 local,  _from: 'local'
          transaction(() => {
            files.forEach(file => {
              this.uploadFile(
                {
                  ...file,
                  path: `${path}/${file.name}`,
                  _master: null,
                  _from: 'local',
                  is_dir: !file.isFile
                },
                null
              )
            })
          })
        }
      } catch (e) {
        if (e) message.error('上传远程文件失败')
      } finally {
        this.fileLoading = false
      }
    } else {
      let fileNames = files.map(file => file.name)

      let confirmFileList = []

      this.files.forEach(value => {
        if (fileNames.includes(value?.['name'])) {
          confirmFileList.push(value)
        }
      })

      if (confirmFileList.length > 0) {
        Modal.showConfirm({
          title: '确认',
          content: `是否覆盖文件或文件夹${confirmFileList
            .map(f => f.name)
            .join(',')}`
        })
          .then(() => {
            confirmFileList.forEach(f => {
              this.deleteFile(f.path, f._master)
              // 删除本地文件
              if (f._from === 'local') {
                Http.post(
                  '/file/delete',
                  { paths: [f.path] },
                  { formatErrorMessage: msg => `删除失败：${msg}` }
                ).catch(e => {
                  console.error(`删除失败：${e}`)
                })
              }
            })

            transaction(() => {
              files.forEach(file => {
                this.uploadFile(
                  {
                    ...file,
                    is_dir: !file.isFile,
                    _master: master,
                    _from: 'server'
                  },
                  master
                )
              })
            })
            return Promise.resolve()
          })
          .catch(() => {
            return Promise.reject()
          })
      }

      transaction(() => {
        files.forEach(file => {
          this.uploadFile(
            { ...file, _master: master, _from: 'server', is_dir: !file.isFile },
            master
          )
        })
      })
    }
  }

  @action
  private deleteAction = (path, master?) => {
    const masterFile = this.getMasterFile(master)

    const targetFile = masterFile
      ? masterFile.slaveFiles.get(path)
      : this.files.get(path)

    // local-upload file/directory
    if (targetFile instanceof Task || targetFile instanceof VirtualDirectory) {
      targetFile.abort()
    }

    // delete file from files
    this.deleteFile(path, master)
  }

  @action
  private setMainAction = (path, isMain) => {
    const file = this.files.get(path)

    // hack: prefetch file content for uploading slave files validation
    // This will activate the filterSlaveFilesFromMaster cache and beforeUploadLocalFile will trigger in a micro task
    // and uploadFile will trigger before the Uploader is unmounted
    this.filterSlaveFilesFromMaster((file as IDoneFile).path).then(() => {
      // update file by ref
      file._isMain = isMain
    })
  }
}
