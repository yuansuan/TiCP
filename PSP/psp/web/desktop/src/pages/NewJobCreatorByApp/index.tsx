/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useMemo, useRef } from 'react'
import { observer } from 'mobx-react-lite'
import * as FormField from '@/components/FormField'
import { account, currentUser } from '@/domain'
import AppParam, { FieldType } from '@/domain/_Application/AppParam'
import { clusterCores } from '@/domain/ClusterCores'
import { getUrlParams } from '@/utils/Validator'
import { PageWrapper } from './style'
import { BottomAction, InputFiles, Section } from './components'
import { Provider, useStore } from './store'
import styled from 'styled-components'
import { AppList } from './AppList'
import { Select } from 'antd'
import { extractPathAndParamsFromURL } from '@/utils'

const Option = Select.Option

const JobCreator = observer((props: Props) => {
  const currentPath = window.localStorage.getItem('CURRENTROUTERPATH')
  const { id, type, mode, upload_id, submit_param } = useMemo(
    () => extractPathAndParamsFromURL(currentPath),
    [currentPath]
  )

  const ref = useRef(null)
  const store = useStore()
  const { data, fileTree, params } = store

  useEffect(() => {
    if (['resubmit', 'continue'].includes(mode as string)) {
      store.setJobBuilderMode(mode as any, {
        id: id as string,
        type: type as any
      })
      store.resubmitParam = submit_param
      store.restoreForResubmit(upload_id)
    } else if (store.resubmitParam !== '') {
      store.restoreForResubmit()
    } else {
      store.setJobBuilderMode('default')
    }

    store.init()
    // 核数显示
    let timer = null
    timer = setInterval(async () => {
      if (store.isCloud) {
        await account.fetch()
      } else {
        await clusterCores.getClusterCoreInfo()
      }
    }, 30 * 1000)

    return () => {
      store.removeHistoryBlock()
      store.reset()
      clearInterval(timer)
    }
  }, [id, mode, type, store])

  useEffect(() => {
    // 点击重置按钮后 此时app为空 自动选择第一个app
    if (!store.data.currentApp) {
      store.updateData({
        currentApp: store.apps[0]
      })
      return
    }
    store.fetchParams()
  }, [store.data.currentApp])

  // useEffect(() => {
  //   if (props.isTop) {
  //     if (!store.isTempDirPath) {
  //       store.fetchJobTree()
  //     }
  //   }
  // }, [props.isTop])

  // useEffect(() => {
  //   store.fetchProjectList(true)
  // }, [])

  useEffect(() => {
    if (store.data.currentApp?.compute_type !== '') {
      store.reset()
      if (store.jobBuildMode !== 'resubmit') {
        store.setTempDirPath('')
      }
      if (store.data.currentApp?.compute_type === 'cloud') {
        account.fetch()
      } else {
        clusterCores.getClusterCoreInfo()
      }
    }
  }, [store.data.currentApp?.compute_type])

  useEffect(() => {
    if (store.data.paramsModel.isTyping) {
      if (store.isInRedeployMode) return

      localStorage.setItem(
        store.draftKey,
        JSON.stringify({
          ...data,
          user_id: currentUser.id
        })
      )
    }
  }, [
    store.data,
    store.mainFilePaths,
    store.currentAppId,
    store.data.paramsModel.isTyping
  ])

  const getComponentByType = (type: string) => {
    return FormField[
      {
        [FieldType.text]: 'Input',
        [FieldType.list]: 'Select',
        [FieldType.multiple]: 'MultiSelect',
        [FieldType.checkbox]: 'Checkbox',
        [FieldType.radio]: 'Radio',
        [FieldType.label]: 'Label',
        [FieldType.date]: 'Date',
        [FieldType.lsfile]: 'Input',
        [FieldType.texarea]: 'TextArea',
        [FieldType.node_selector]: 'NodeSelector',
        [FieldType.cascade_selector]: 'CascadeSelector'
      }[type]
    ]
  }

  return (
    <PageWrapper id='job_creator' ref={ref}>
      <div className='input-content'>
        <AppList action={props.action} is_trial={false} />
        {/* <Section title='项目选择' className='moduleProjectSelector'>
          <Select style={{ width: 300}}
            size='middle'
            placeholder="请选择提交作业所属项目"
            value={store.projectId}
            onDropdownVisibleChange={open => {
              if (open) store.fetchProjectList(false)
            }}
            onChange={store.setProjectId}>
            {
              store.projectList.map(item => {
                return <Option key={item.id} value={item.id}>{item.name}</Option>
              })
            }
          </Select>
        </Section> */}
        <Section title='上传模型' className='moduleUpload'>
          <InputFiles fileTree={fileTree} />
        </Section>

        {params?.length > 0 &&
          params.map(param => (
            <Section
              className='paramSettings'
              title={param?.name}
              key={param?.name}>
              {param?.field?.map((item, i) => {
                const Field = getComponentByType(item.type)
                return (
                  <Field
                    key={item?.id}
                    formModel={data.paramsModel}
                    model={item}
                    appId={store.data.currentApp.id}
                  />
                )
              })}
            </Section>
          ))}
      </div>

      <BottomAction {...props} />
    </PageWrapper>
  )
})

export const StyledDiv = styled.div`
  height: calc(100vh - 110px);
  max-height: calc(100vh - 110px);
`

export default (props: { action?: any; isTop: boolean }) => (
  <Provider>
    <StyledDiv>
      <JobCreator {...props} pushHistoryUrl='/new-jobs?tab=jobs' />
    </StyledDiv>
  </Provider>
)

export const StyledInDrawerJobCreatorDiv = styled.div`
  height: calc(100% - 20px);
`

type Props = {
  onOk?: () => void
  pushHistoryUrl?: string
  action?: any
  isTop?: boolean
}

export function InDrawerJobCreator(props: Props) {
  return (
    <Provider>
      <StyledInDrawerJobCreatorDiv>
        <JobCreator {...props} />
      </StyledInDrawerJobCreatorDiv>
    </Provider>
  )
}
