/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { JobDraftEnum } from '@/constant'
import { appList, currentUser, uploader, account } from '@/domain'
import { FileTree } from '@/domain/JobBuilder/FileTree'
import { getFilenameByPath, history, Http, parseUrlParam } from '@/utils'
import { fromJSON2Tree, fromJSON2Tree2 } from '@/domain/JobBuilder/utils'
import App from '@/domain/_Application/App'
import { message } from 'antd'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { clusterCores } from '@/domain/ClusterCores'
import { newBoxServer } from '@/server'
import {
  getContinuousRedeployInfo,
  getRedeployInfo
} from '@/domain/JobBuilder/NewJobBuilder'
import { DraftType, NewDraft } from '@/domain/Box'

class IJobBuilderData {
  name: string
  currentApp: App
  paramsModel: any
  scIds: string[]
  numCores: number
  currentAppId?: string
  mainFilePaths?: string[]
}

class ResubmitParam {
  app_id: string
  project_id: string
  user_id: string
  user_name: string
  queue_name: string
  main_files: string[]
  work_dir: {
    path: string
    is_temp: boolean
  }
  fields: {
    id: string
    type: string
    value: string
    values: string[]
  }[]
}

const initData = {
  name: '',
  currentApp: null,
  paramsModel: { isTyping: false },
  numCores: null,
  mainFilePaths: undefined
}

export function useModel() {
  return useLocalStore(() => {
    const currentPath = window.localStorage.getItem('CURRENTROUTERPATH')
    const params = parseUrlParam(currentPath)
    return {
      tabKey: 'formal',
      setTabKey(key) {
        this.tabKey = key
      },
      get currentJobHasParams() {
        // 校验模型文件
        if (this.fileTree?._children?.length > 0) {
          return true
        }

        // // 校验作业模版参数
        // const { paramsModel } = this.data
        // for (const key of Object.keys(paramsModel)) {
        //   const param = paramsModel[key]
        //   if (
        //     param &&
        //     param.required &&
        //     !param.hidden &&
        //     (param.value || param.values?.length)
        //   ) {
        //     return true
        //   }
        // }
        return false
      },
      get is_trial() {
        return this.tabKey === 'trial'
      },

      get isCross() {
        return this.tempDirPath ? true : false
      },
      get isCloud() {
        return this.data.currentApp?.compute_type === 'cloud' ? true : false
      },
      unblock: undefined,
      setUnblock(v) {
        this.unblock = v
      },
      unlisten: undefined,
      setUnlisten(v) {
        this.unlisten = v
      },

      removeHistoryBlock() {
        this.unblock && this.unblock()
        this.unlisten && this.unlisten()
      },

      async refresh() {},

      jobBuildMode: 'default',
      get draftKey(): string {
        switch (this.jobBuildMode) {
          case 'default':
            return JobDraftEnum.JOB_DRAFT_STORE_KEY
          case 'redeploy':
            return JobDraftEnum.JOB_REDEPLOY_DRAFT_STORE_KEY
          case 'continuous':
            return JobDraftEnum.JOB_CONTINUOUS_DRAFT_STORE_KEY
        }
        return null
      },
      get draft(): NewDraft {
        switch (this.jobBuildMode) {
          case 'default':
            return new NewDraft(DraftType.Default)
          case 'resubmit':
            return new NewDraft(DraftType.Resubmit)
          case 'continue':
            return new NewDraft(DraftType.Continue)
        }
        return null
      },

      fileTree: new FileTree({
        name: ''
      }),
      data: initData,

      params: [],
      isJobSubmitting: false,

      get isInRedeployMode() {
        return this.jobBuildMode === 'redeploy'
      },

      get isInConMode() {
        return this.jobBuildMode === 'continuous'
      },

      get isDefault() {
        return this.jobBuildMode === 'default'
      },

      expandKeys: [],
      expandFlag: false,
      resubmitParam: '',
      getResubmitParam() {
        return (
          this.resubmitParam !== '' &&
          (JSON.parse(this.resubmitParam || '{}') as ResubmitParam)
        )
      },

      restoreData: undefined,
      redeployJobId: '',
      continuousJobId: '',
      redeployType: 'job',
      jobQueue: '',
      setJobQueue(queue) {
        this.jobQueue = queue
      },
      tempDirPath: '',
      setTempDirPath(tempDirPath) {
        this.tempDirPath = tempDirPath
      },
      isTempDirPath: true, // 默认进来是临时工作目录
      setIsTempDirPath(isTempDirPath) {
        this.isTempDirPath = isTempDirPath
      },
      projectId: '',
      setProjectId(value) {
        this.projectId = value
      },
      projectList: [],
      setProjectList(values) {
        this.projectList = values || []
      },
      async fetchProjectList(isInitProjectId = false) {
        const res = await Http.get('project/list/current', {
          params: {
            state: 'Running'
          }
        })
        this.setProjectList(res?.data?.projects || [])
        if (isInitProjectId && this.projectId === '') {
          this.setProjectId(res?.data?.projects[0]?.id)
        }
      },
      async fetchTempDirPath() {
        await Http.post('job/createTempDir', {
          compute_type: this.data.currentApp.compute_type
        }).then(res => {
          this.setTempDirPath(res.data.path)
        })
      },
      setJobBuilderMode(
        key: 'default' | 'resubmit' | 'continue',
        param?: { id: string; type: 'job' | 'jobset' }
      ) {
        this.jobBuildMode = key
        if (key === 'resubmit') {
          this.redeployJobId = param.id
          this.redeployType = param.type
          this.jobBuildMode = 'resubmit'
        } else if (key === 'continue') {
          this.continuousJobId = param.id
        }
      },

      restoreForResubmit(upload_id = '') {
        if (
          this.jobBuildMode === 'resubmit' &&
          Object.keys(this.getResubmitParam()).length > 0
        ) {
          this.setProjectId(this.getResubmitParam().project_id)
          this.setJobQueue(this.getResubmitParam().queue_name)
          this.setTempDirPath(this.getResubmitParam().work_dir?.path)
          this.setIsTempDirPath(this.getResubmitParam().work_dir?.is_temp)

          this.apps.forEach(app => {
            if (app.id === this.getResubmitParam().app_id) {
              this.data.currentApp = app
              return
            }
          })

          if (upload_id !== '') {
            EE.emit(EE_CUSTOM_EVENT.SUPERCOMPUTING_TASKKEY, {
              taskKey: upload_id
            })
            EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, { visible: true })
            EE.once(
              EE_CUSTOM_EVENT.SERVER_FILE_TO_SUPERCOMPUTING,
              async ({ file_status }) => {
                console.log('file_status=========+>: ', file_status)
                if (file_status === 'success') {
                  const fileData = await this.draft.getFileList({
                    path: this.tempDirPath,
                    cross: this.isTempDirPath,
                    is_cloud: this.isCloud
                  })
                  let selectFileSet = new Set<String>()
                  if (this.jobBuildMode === 'resubmit') {
                    selectFileSet = new Set<String>(
                      this.getResubmitParam()?.main_files
                    )
                  }
                  this.fileTree.uploadCommonFiles(
                    fileData,
                    this.tempDirPath,
                    selectFileSet,
                    this.isTempDirPath,
                    null
                  )
                }
                this.expandFlag = true
              }
            )
            return
          }
          ;(async tempDirPath => {
            const fileData = await this.draft.getFileList({
              path: this.getResubmitParam().work_dir?.path,
              cross: this.isTempDirPath,
              is_cloud: this.isCloud
            })
            this.fileTree.uploadCommonFiles(
              fileData,
              tempDirPath,
              new Set<String>(this.getResubmitParam().main_files),
              this.isTempDirPath
            )
            this.expandFlag = true
          })(this.tempDirPath)
        }
      },

      async restore() {
        // reset draft
        const draftCache = localStorage.getItem(this.draftKey)
        if (draftCache) {
          try {
            const draft = JSON.parse(draftCache)
            if (draft.user_id !== currentUser.id) {
              localStorage.removeItem(this.draftKey)
            }
          } catch (e) {
            localStorage.removeItem(this.draftKey)
          }
        }
        switch (this.jobBuildMode) {
          case 'default':
            await this.restoreFromDraft()
            break
          case 'redeploy':
            try {
              this.restoreData = await getRedeployInfo({
                id: this.redeployJobId,
                type: this.redeployType as any,
                clean: false
              })
            } catch (e) {
              setTimeout(() => history.replace('/new-jobs'), 300)
              return
            }

            await this.restoreFromRedeploy()
            break
          case 'continuous':
            try {
              this.restoreData = await getContinuousRedeployInfo({
                id: this.continuousJobId,
                type: this.redeployType as any,
                clean: false
              })
            } catch (e) {
              setTimeout(() => history.replace('/new-jobs'), 300)
              return
            }

            await this.restoreContinuousFromResult()
            break
        }
      },

      get apps() {
        return appList.list.filter(app => app.state === 'published')
      },

      // only can be used for constructor function
      // to restore data from localstorage when page refresh
      async init() {
        await this.fetchApps()
        // 从草稿恢复作业
        await this.restore()
      },

      async fetchApps() {
        await Promise.all([appList.fetch()])
      },

      async fetchParams() {
        await Promise.all([
          // 获取应用参数
          this.initSection()
        ])

        // 初始化应用参数
        this.initParamsData()
      },

      updateData(data: Partial<IJobBuilderData>) {
        Object.assign(this.data, data)
      },

      get defaultJobName() {
        return this.mainFilePaths.length === 1
          ? getFilenameByPath(this.mainFilePaths[0])
          : ''
      },

      get mainFiles() {
        return this.fileTree.flatten().filter(node => {
          return node.isFile && node.isMain
        })
      },

      get mainFilePaths() {
        return this.mainFiles.map(node => node.path)
      },

      get currentAppId() {
        return this.data.currentApp ? this.data.currentApp.id : ''
      },

      // 获取app参数
      async initSection() {
        if (!this.data.currentApp) return
        this.params = []
        this.params = await this.data.currentApp.getParams()
      },

      // 初始化应用参数，合并缓存中的数据
      initParamsData(cacheParams: any = {}) {
        if (!this.data.currentApp) {
          return
        }

        const data = {}
        const resubmitParamFiledMap = new Map(
          (this.getResubmitParam()?.fields ?? []).map(field => [
            field.id,
            field
          ])
        )

        this.params.forEach(item => {
          item?.field.forEach(field => {
            let cache = cacheParams[field.id]
            if (
              this.jobBuildMode === 'resubmit' &&
              resubmitParamFiledMap.has(field.id)
            ) {
              cache = resubmitParamFiledMap.get(field.id)
            }

            const value =
              cache?.value === '' || cache?.value === undefined
                ? field.value === ''
                  ? field.defaultValue
                  : field.value
                : cache?.value
            const values =
              Array.isArray(cache?.values) && cache?.values?.length > 0
                ? cache?.values
                : field?.values?.length > 0
                ? field.values
                : field.defaultValues

            data[field.id] = {
              ...field,
              value,
              values
            }
          })
        })
        this.updateData({ paramsModel: data })
      },

      async jobSetNameCheck(name): Promise<boolean> {
        const { data } = await Http.get('job/set/name_check', {
          params: {
            name
          }
        })
        return data.exists
      },

      /**
       * 创建作业：调用盒子api将draft目录移到input目录
       * 将返回的参数传给创建作业模块
       */
      async create(entries) {
        const { paramsModel, currentApp } = this.data
        if (!this.beforeCreate()) return false
        entries.forEach(async e => {
          if (e.mainFile.endsWith('.jobs')) {
            // fetch file content
            // e.content = await this.draft.getFileContent(e.mainFile).then(a => a)
          }
        })

        try {
          this.isJobSubmitting = true

          const fields = Object.keys(paramsModel).map(key => {
            return {
              id: key,
              value: paramsModel[key].value,
              values: paramsModel[key].values,
              type: paramsModel[key].type
            }
          })

          const data = {
            app_id: currentApp.id,
            main_files:
              this.mainFilePaths?.map(mainFile =>
                mainFile.replace(/^\.\//, '')
              ) || [], // /^\.\// 匹配 ./ 开头的文件
            queue_name: this.jobQueue,
            work_dir: {
              path: this.tempDirPath,
              is_temp: this.isTempDirPath
            },
            fields,
            project_id: this.projectId
          }

          await Http.post('job/submit', { ...data })
        } catch (e) {
          return false
        } finally {
          this.isJobSubmitting = false
        }
        return true
      },

      beforeCreate() {
        // 校验作业参数
        if (!this.mainFilePaths.length) {
          message.error('请至少选择一个主文件')
          return false
        }

        if (!this.data.currentApp) {
          message.error('请选择应用')
          return false
        }

        // 校验作业模版
        const { paramsModel } = this.data
        for (const key of Object.keys(paramsModel)) {
          const param = paramsModel[key]
          if (param && param.required && !param.value && param.values === 0) {
            message.error(`请填写${param.label}`)
            return false
          }
        }
        return true
      },

      reset() {
        this.fileTree = new FileTree({
          name: ''
        })
        // 重置时保留app
        this.updateData({ ...initData, currentApp: this.data.currentApp })
      },
      resetParams() {
        this.isTempDirPath = true
        this.tempDirPath = ''
      },
      resetFileTree() {
        this.fileTree = new FileTree({ name: '' })
      },

      // 清理draft、取消上传、清理文件树
      async clean(refetch = false) {
        localStorage.removeItem(this.draftKey)

        // cancel uploader
        this.fileTree.tapNodes(
          () => true,
          node => {
            uploader.remove(node.uid)
          }
        )

        this.reset()

        if (this.isInRedeployMode || this.isInConMode) {
          // 重提交下重置会通过api恢复作业参数
          await this.init()
        } else if (refetch) {
          // 非重提交下 点击重置按钮仅需重新获取参数
          await this.fetchParams()
        }

        // 清理指定工作目录
        if (!this.isTempDirPath) {
          this.isTempDirPath = true
          this.tempDirPath = ''
        }
      },

      async fetchJobTree(mainFilePaths?: string[]) {
        if (!mainFilePaths) {
          mainFilePaths = this.mainFilePaths || []
        }

        // 获取目录的文件列表

        let res = await newBoxServer.list({
          path: this.tempDirPath,
          cross: this.isTempDirPath,
          recursive: true
        })

        this.fileTree = fromJSON2Tree2(
          res.data || [],
          this.tempDirPath,
          this.isTempDirPath
        )
        // 恢复主文件，从文件树中找到对应的主文件 id
        this.fileTree.tapNodes(
          node => mainFilePaths.includes(node.path),
          node => {
            node.isMain = true
          }
        )
      },

      async restoreFileList(mainFilePaths?: string[]) {
        if (!mainFilePaths) {
          mainFilePaths = this.data.mainFilePaths || []
        }

        // 恢复上传的文件，获取盒子中draft的文件
        // let files = await this.draft.listFile()
        let files = []

        this.fileTree = fromJSON2Tree(files)
        // 恢复主文件，从文件树中找到对应的主文件 id
        this.fileTree.tapNodes(
          node => mainFilePaths.includes(node.path),
          node => {
            node.isMain = true
          }
        )
      },

      // 从草稿中恢复到工作区
      async restoreFromDraft() {
        if (!this.apps.length) return

        let cache: IJobBuilderData, isCacheEmpty

        // parse cache
        const cacheStr = localStorage.getItem(this.draftKey)

        if (!cacheStr) {
          cache = {} as IJobBuilderData
          isCacheEmpty = true
        } else {
          cache = JSON.parse(cacheStr)
          isCacheEmpty = false
        }

        delete cache['user_id']

        // 恢复app，如果找不到，重置为第一个
        const firstApp = this.apps?.filter(app => app.id === params?.id)[0]
        this.updateData({
          currentApp:
            firstApp ||
            this.apps.find(app => app.id === cache.currentAppId) ||
            this.apps[0]
        })

        delete cache.currentAppId
        await this.restoreFileList(cache.mainFilePaths)

        if (this.data.currentApp) {
          // 恢复作业参数
          await this.initSection()
          this.initParamsData(cache?.paramsModel)
        }

        !isCacheEmpty &&
          this.updateData({
            name: cache.name,
            numCores: cache.numCores
          })
      },

      // 从重提交远端中恢复到工作区
      async restoreFromRedeploy() {
        if (!this.apps.length) return

        let cache: IJobBuilderData, isCacheEmpty
        cache = { ...this.restoreData }
        cache.name += '_re'
        // await this.draft.back(cache['input_folder_uuid']).catch(() => {
        //   setTimeout(() => (location.href = '/new-jobs'), 300)
        // })

        delete cache['user_id']

        const firstApp = this.apps?.filter(app => app.id === params?.id)[0]
        // 恢复app，如果找不到，重置为第一个
        this.updateData({
          currentApp:
            firstApp ||
            this.apps.find(app => app.id === cache.currentAppId) ||
            this.apps[0]
        })
        delete cache.currentAppId

        await this.restoreFileList(cache.mainFilePaths)

        if (this.data.currentApp) {
          // 恢复作业参数
          await this.initSection()
          this.initParamsData(cache?.paramsModel)
        }

        !isCacheEmpty && this.updateData(cache)
      },

      // 从作业结果恢复到提交界面
      async restoreContinuousFromResult() {
        if (!this.apps.length) return

        let cache: IJobBuilderData, isCacheEmpty
        cache = { ...this.restoreData }
        cache.name += '_continuous'

        delete cache['user_id']

        // 恢复app，如果找不到，重置为第一个
        const firstApp = this.apps?.filter(app => app.id === params?.id)[0]
        this.updateData({
          currentApp:
            firstApp ||
            this.apps.find(app => app.id === cache.currentAppId) ||
            this.apps[0]
        })
        delete cache.currentAppId

        await this.restoreFileList(cache.mainFilePaths)

        if (this.data.currentApp) {
          // 恢复作业参数
          await this.initSection()
          this.initParamsData(cache?.paramsModel)
        }

        !isCacheEmpty && this.updateData(cache)
      }
    }
  })
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
