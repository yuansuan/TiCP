/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useMemo } from 'react'
import { observer } from 'mobx-react-lite'
import * as FormField from '@/components/FormField'
import { env, currentUser } from '@/domain'
import { FieldType } from '@/domain/Applications/App/Field'
import { history } from '@/utils'
import { getUrlParams } from '@/utils/Validator'
import { PageWrapper } from './style'
import { BottomAction, InputFiles, Section } from './components'
import { Provider, useStore } from './store'
import styled from 'styled-components'
import { AppList } from './AppList'
const JobCreator = observer((props: Props) => {
  const { id, type, mode } = useMemo(
    () => getUrlParams(),
    [window.location.hash]
  )

  const store = useStore()
  const { data, fileTree, params } = store

  useEffect(() => {
    if (['redeploy', 'continuous'].includes(mode as string)) {
      store.setJobBuilderMode(mode as any, {
        id: id as string,
        type: type as any
      })
      // store.setUnblock(
      //   history.block(() => '离开页面将不会保留当前工作，确认要离开页面吗？')
      // )

      // store.setUnlisten(
      //   history.listen(location => {
      //     // 仍需跳转到当前页面时reload
      //     location.pathname === '/new-job-creator' &&
      //       window.location.replace('/')
      //   })
      // )
    } else {
      store.setJobBuilderMode('default')
    }

    // TODO PSP作业提交页面初始化
    store.init()

    return () => {
      store.removeHistoryBlock()
      store.reset()
    }
  }, [id, mode, store, type])

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
        [FieldType.node_selector]: 'NodeSelector'
      }[type]
    ]
  }

  return (
    <PageWrapper id='job_creator'>
      <div className='input-content'>
        <AppList action={props.action} is_trial={false} />
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

export default (props: { action?: any }) => (
  <Provider>
    <StyledDiv>
      <JobCreator {...props} pushHistoryUrl='/new-jobs?tab=jobs' />
    </StyledDiv>
  </Provider>
)

export const StyledInDrawerJobCreatorDiv = styled.div``

type Props = {
  onOk?: () => void
  pushHistoryUrl?: string
  action?: any
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
