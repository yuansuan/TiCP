/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Modal } from '@/components'
import { Job } from '@/domain/JobList/Job'
import { Descriptions } from 'antd'
import { scList } from '@/domain'
import { jobServer } from '@/server'
import { observer, useLocalStore } from 'mobx-react-lite'
import { buryPoint } from '@/utils'

const StyledLayout = styled.div`
  > h2 {
    padding-left: 20px;
    font-size: 14px;
    line-height: 54px;
    font-weight: bold;
    color: #666666;
    background: #f3f5f8;
    margin: 0 0 6px;
  }

  > table {
    width: 100%;
    word-break: break-all;
    tbody {
      border-bottom: 1px solid #e8e8e8;
      vertical-align: baseline;
      tr {
        &.date {
          td {
            padding-top: 6px;
            color: #999999;
            font-size: 10px;
            line-height: 14px;
          }
        }
        &.message {
          font-size: 14px;
          line-height: 20px;
          color: #666666;
          td {
            padding-bottom: 6px;
            &:first-of-type {
              width: 62px;
              text-align: right;
              &::after {
                content: '：';
              }
            }
          }
        }
      }
    }
  }
`
interface Props {
  jobID: string
}

const HelpHistory = observer(function HelpHistory({ jobID }: Props) {
  const state = useLocalStore(() => ({
    job: new Job(),
    helpList: [],
    setHelpList(list) {
      this.helpList = list
    }
  }))
  const { job } = state

  useEffect(() => {
    jobServer.get(jobID).then(({ data }) => {
      job.update(data)
    })
  }, [jobID])

  return (
    <StyledLayout>
      <Descriptions>
        <Descriptions.Item label='作业名称'>
          {job?.name || '--'}
        </Descriptions.Item>
        <Descriptions.Item label='作业编号'>
          {job?.id || '--'}
        </Descriptions.Item>
       
       
        <Descriptions.Item label='软件'>
          {job?.app_name || '--'}
        </Descriptions.Item>
        <Descriptions.Item label='版本'>
          {job?.app_version || '--'}
        </Descriptions.Item>
      </Descriptions>
    </StyledLayout>
  )
})

export const showHelpHistory = (jobID: string) => {
  buryPoint({
    category: '作业管理',
    action: '帮助记录'
  })
  Modal.show({
    title: '帮助记录',
    footer: null,
    content: <HelpHistory jobID={jobID} />,
    width: 850,
    bodyStyle: {
      height: 585
    }
  })
}
