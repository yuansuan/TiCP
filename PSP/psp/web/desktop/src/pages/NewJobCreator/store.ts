/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { runInAction } from 'mobx'
import { JobDraftEnum } from '@/constant'
import { appList, currentUser, env, uploader } from '@/domain'
import { FileTree } from '@/domain/JobBuilder/FileTree'
import { getFilenameByPath, history, Http, parseUrlParam } from '@/utils'
import App from '@/domain/_Application/App'
import { DraftType, NewDraft } from '@/domain/Box'
import { message } from 'antd'
import { companyServer } from '@/server'
import {
  getContinuousRedeployInfo,
  getRedeployInfo
} from '@/domain/JobBuilder/NewJobBuilder'

class IJobBuilderData {
  name: string
  currentApp: App
  paramsModel: any
  scIds: string[]
  numCores: number
  currentAppId?: string
  mainFilePaths?: string[]
}

const initData = {
  name: '',
  compute_type: '',
  currentApp: null,
  paramsModel: { isTyping: false },
  numCores: null,
  mainFilePaths: undefined
}

export function useModel() {
  return useLocalStore(() => {
    return {
      tabKey: 'formal',
      setTabKey(key) {
        this.tabKey = key
      },
      get is_trial() {
        return this.tabKey === 'trial'
      },
      get is_cloud() {
        return this.data.compute_type === 'cloud'
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
      refresh() {},
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
          case 'redeploy':
            return new NewDraft(DraftType.Redeploy)
          case 'continuous':
            return new NewDraft(DraftType.Continuous)
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

      restoreData: undefined,
      redeployJobId: '',
      continuousJobId: '',
      redeployType: 'job',
      tempDirPath: '',
      setTempDirPath(tempDirPath) {
        this.tempDirPath = tempDirPath
      },
      async fetchTempDirPath() {
        const res = await Http.post('job/createTempDir', {
          user_name: currentUser.name,
          compute_type: this.data.compute_type
        }).then(res => {
          this.setTempDirPath(res.data.path)
        })
        return res.data?.path
      },
      setJobBuilderMode(
        key: 'default' | 'redeploy' | 'continuous',
        param?: { id: string; type: 'job' | 'jobset' }
      ) {
        this.jobBuildMode = key
        if (key === 'redeploy') {
          this.redeployJobId = param.id
          this.redeployType = param.type
        } else if (key === 'continuous') {
          this.continuousJobId = param.id
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
        return appList.publishedAppList
      },

      get displayExpectCoresScope() {
        return '范围(1-200)核'
      },
      // 填写期望核数，新的核数范围固定为 1-200
      get displayExpectCoreValid() {
        if (this.data.numCores) {
          return (
            this.data.numCores >= 1 &&
            this.data.numCores <= 200 &&
            typeof this.data.numCores === 'number'
          )
        } else {
          if (this.data.numCores === 0) {
            return false
          }
          return true
        }
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
        this.params.forEach(item => {
          item?.field.forEach(field => {
            const cache = cacheParams[field.id]
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
       * 创建作业：
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
          user_id: currentUser.id,
          main_file: this.tempDirPath + '/' + this.mainFilePaths,
          work_dir: {
            path: this.tempDirPath,
            is_temp: true
          },
          fields
        }

        // const store = {
        //   name: data.name,
        //   paramsModel,
        //   scIds: this.data.scIds,
        //   numCores: this.data.numCores,
        //   currentAppId: currentApp.id,
        //   mainFilePaths: this.mainFilePaths,
        //   input_folder_uuid: res.data.input_folder_uuid
        // }
        await Http.post('job/submit', { ...data }, {})

        try {
          this.isJobSubmitting = true
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

        const jobName = this.data.name || this.defaultJobName
        if (!jobName) {
          message.error('请填写作业名称')
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

      // 清理文件树
      reset() {
        this.fileTree = new FileTree({
          name: ''
        })
        // 重置时保留app
        this.updateData({ ...initData, currentApp: this.data.currentApp })
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
      },

      async restoreFileList(mainFilePaths?: string[]) {
        if (!mainFilePaths) {
          mainFilePaths = this.data.mainFilePaths || []
        }
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
        this.updateData({
          currentApp:
            this.apps.find(app => app.id === cache.currentAppId) || this.apps[0]
        })
        delete cache.currentAppId
        await this.restoreFileList(cache.mainFilePaths)

        if (this.data.currentApp) {
          // 恢复作业参数
          await this.initSection()
          this.initParamsData(cache?.paramsModel)

          // 恢复算力资源，去除不存在的算力资源
          cache.scIds = (cache.scIds || []).filter(scId =>
            this.data.currentApp.sc_ids.includes(scId)
          )
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

        delete cache['project_id']
        delete cache['user_id']

        // 恢复app，如果找不到，重置为第一个
        this.updateData({
          currentApp:
            this.apps.find(app => app.id === cache.currentAppId) || this.apps[0]
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
        this.updateData({
          currentApp:
            this.apps.find(app => app.id === cache.currentAppId) || this.apps[0]
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
