/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { useDispatch } from 'react-redux'
import { BottomActionStyle } from './style'
import { Button, Modal } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Completeness } from './Completeness'
import { useStore } from '../store'
import { buryPoint, getFilenameByPath, history } from '@/utils'
import { message } from 'antd'

type Props = {
  onOk?: () => void
  pushHistoryUrl?: string
}

export const BottomAction = observer(({ onOk, pushHistoryUrl }: Props) => {
  const store = useStore()

  const dispatch = useDispatch()
  const state = useLocalStore(() => ({
    get paramsSettingsOk() {
      if (Object.keys(store.data.paramsModel).length === 0) return false
      for (const key of Object.keys(store.data.paramsModel)) {
        const param = store.data.paramsModel[key]
        if (param && param.required && !param.value) {
          return false
        }
      }
      return true
    },
    get prerequisites() {
      return [
        {
          label: '项目选择',
          isCompleted: store.projectId
        },
        {
          label: '上传模型',
          isCompleted:
            store.fileTree.flatten().filter(node => node.status === 'done')
              .length > 0 &&
            store.fileTree.filterNodes(
              file => file.isFile && file.status !== 'done'
            ).length === 0
        },
        {
          label: '设置主文件',
          isCompleted: store.mainFiles?.length > 0
        },
        {
          label: '参数配置',
          isCompleted: state.paramsSettingsOk
        }
      ]
    },
    async onSubmit() {
      try {
        store.isJobSubmitting = true
        buryPoint({
          category: '作业提交',
          action: '提交'
        })
        const map = new Map()

        const newJobs = store.mainFiles.map(file => {
          let jobName = getFilenameByPath(file.path)

          if (map.has(jobName)) {
            let num = map.get(jobName)
            map.set(jobName, num + 1)
            jobName = `${jobName}_${num}`
          } else {
            map.set(jobName, 1)
          }

          return {
            mainFile: file.path,
            jobName,
            content: ''
          }
        })
        const isSuccess = await store.create(newJobs)
        if (isSuccess) {
          store.removeHistoryBlock()
          store.resetParams()
          store.resetFileTree()
          store.resubmitParam = ''
          store.jobBuildMode = 'default'
          message.success('作业创建成功')

          // dispatch({
          //   type: store.data.currentApp?.type,
          //   payload: 'close'
          // })
          dispatch({
            type: 'JOBMANAGE',
            payload: 'winRefresh'
          })
          dispatch({
            type: 'JOBMANAGE',
            payload: 'full'
          })

          onOk && onOk()
          if (pushHistoryUrl) {
            window.localStorage.setItem('CURRENTROUTERPATH', pushHistoryUrl)
          }
          store.clean(true)
        }
      } catch (e) {
        console.error(e)
      } finally {
        store.isJobSubmitting = false
      }
    },
    async onCancel() {
      const content = {
        default: '重置将会清空工作区的文件和已填写的参数，是否确认？',
        redeploy: '重置将会恢复重提交作业时的参数，是否确认？',
        continuous: '重置将会恢复续算提交作业时的参数，是否确认？'
      }[store.jobBuildMode]

      await Modal.showConfirm({
        title: '重置',
        content
      })

      buryPoint({
        category: '作业提交',
        action: '重置'
      })

      await store.clean(true)
      message.success('已重置')
    },
    get createBtnDisabled() {
      const { data, isJobSubmitting, fileTree, projectId } = store
      const { paramsModel } = data

      if (!projectId) {
        return '请选择提交作业所属的项目'
      }

      const successFiles = store.fileTree
        .flatten()
        .filter(node => node.status === 'done')

      if (successFiles.length === 0) {
        return '请上传模型至作业目录'
      }

      if (!store.mainFilePaths.length) {
        return '请至少选择一个主文件'
      }

      if (
        store.fileTree.filterNodes(
          file => file.isFile && file.status !== 'done'
        ).length > 0
      ) {
        return '请等待全部文件上传完成'
      }

      for (const key of Object.keys(paramsModel)) {
        const param = paramsModel[key]
        if (param && param.required && !param.value) {
          return `请填写${param.label}`
        }
      }

      const uploadingFile = fileTree
        .flatten()
        .find(node => node.isFile && node.status === 'uploading')
      if (uploadingFile) {
        return `请等待${uploadingFile.name}上传完成`
      }

      return isJobSubmitting
    }
  }))

  return (
    <BottomActionStyle className='bottomToolbar'>
      <div className='formCompleteness'>
        <span className='label'>作业完整情况：</span>
        {state.prerequisites.map((prerequisite, index) => (
          <span className='completeness' key={index}>
            <Completeness isCompleted={prerequisite.isCompleted} />
            {prerequisite.label}
          </span>
        ))}
      </div>

      <div className='actions'>
        <Button onClick={state.onCancel}>重置</Button>
        <Button
          type='primary'
          className='submitJob'
          onClick={state.onSubmit}
          disabled={state.createBtnDisabled}
          loading={store.isJobSubmitting}>
          提交
        </Button>
      </div>
    </BottomActionStyle>
  )
})
