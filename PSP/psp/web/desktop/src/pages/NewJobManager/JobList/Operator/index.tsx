/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import qs from 'qs'
import { useDispatch } from 'react-redux'
import { Button, Modal } from '@/components'
import { Divider, message } from 'antd'
import { env } from '@/domain'
import { Dropdown } from '@/components/Dropdown'
import { jobServer } from '@/server'
import { useStore } from '../store'
import { Http, history } from '@/utils'
import { showDownloader } from '../showDownloader'
import { buryPoint } from '@/utils'
import {
  getContinuousRedeployInfo,
  getRedeployInfo,
  getResubmit
} from '@/domain/JobBuilder/NewJobBuilder'
import { getUrlParams } from '@/utils/Validator'

const StyledLayout = styled.div``

type Props = {
  id: string
  out_job_id: string
  type: string
  raw_state: string
  name: string
  terminalable: boolean
  resubmittable: boolean
  downloadable: boolean
  deleteable: boolean
  showContinuousRedeploy: boolean
  display_back_state: number
}

export const Operator = observer(function Operator({
  id,
  out_job_id,
  raw_state,
  name,
  type,
  terminalable,
  resubmittable,
  downloadable,
  deleteable,
  showContinuousRedeploy,
  display_back_state
}: Props) {
  const store = useStore()
  const dispatch = useDispatch()
  const { model } = store

  async function cancelJob() {
    if (raw_state === 'Terminating') {
      message.warn(`【${name}】已经在终止中，请稍后操作！`)
    } else {
      await Modal.showConfirm({
        title: '确认终止',
        content: `终止作业【${name}】，是否确认？`
      })

      buryPoint({
        category: '作业管理',
        action: '终止作业'
      })
      await jobServer.terminate({ out_job_id, compute_type: type })

      message.success(`作业【${name}】终止中，请耐心等待`)
      store.refresh()
    }
  }

  async function resubmitJob(jobId) {
    buryPoint({
      category: '作业管理',
      action: '重新提交'
    })

    const {
      data: {
        param,
        extension: { app_type, upload_id }
      }
    } = await jobServer.resubmit(jobId)

    window.localStorage.setItem(
      'CURRENTROUTERPATH',
      `/new-job-creator?id=${jobId}&type=jobs&app_type=${app_type}&mode=resubmit&upload_id=${upload_id}&submit_param=${JSON.stringify(
        param
      )}`
    )

    JSON.parse(window.localStorage.getItem('FLAG_ENTERTAINMENT') || '[]')?.map(
      app => {
        dispatch({
          type: app?.type,
          payload: 'close'
        })
      }
    )

    dispatch({
      type: 'DESKTOP',
      payload: 'winRefresh'
    })
    dispatch({
      type: 'JOBMANAGE',
      payload: 'togg'
    })
    setTimeout(() => {
      dispatch({
        type: app_type,
        payload: 'full'
      })
    }, 500)
  }

  async function goToContinuousRedeployPage(id) {
    buryPoint({
      category: '作业管理',
      action: '续算提交'
    })

    await getContinuousRedeployInfo({
      id,
      type: 'job',
      clean: false
    })

    history.push(
      `/new-job-creator?${qs.stringify({
        id: id,
        type: 'job',
        mode: 'continuous'
      })}`
    )
    dispatch({
      type: 'CALCUAPP',
      payload: 'close'
    })
    dispatch({
      type: 'JOBMANAGE',
      payload: 'close'
    })
    setTimeout(() => {
      dispatch({
        type: 'CALCUAPP',
        payload: 'togg'
      })
    })
  }

  function downloadJob() {
    buryPoint({
      category: '作业管理',
      action: '下载至本地'
    })
  }

  async function downloadJobToCommon() {
    buryPoint({
      category: '作业管理',
      action: '下载至我的文件'
    })
    const job = model.list.find(job => job.id === id)
    const resolvedJobs = await showDownloader([
      {
        id: job.id,
        name: job.name
      }
    ])
    if (Object.keys(resolvedJobs).length > 0) {
      message.success('下载完成')
      store.setSelectedKeys([])
    }
  }

  async function deleteJob() {
    buryPoint({
      category: '作业管理',
      action: '删除'
    })
    await Modal.showConfirm({
      title: '确认删除',
      content: '删除作业同时会删除作业产生的文件，是否确认？'
    })
    await jobServer.delete([id])
    store.refresh()
    message.success('删除成功')
  }

  return (
    <StyledLayout>
      <Button
        style={{ padding: 0 }}
        disabled={!terminalable}
        onClick={() => cancelJob()}
        type='link'>
        终止作业
      </Button>

      <Divider type='vertical' style={{ margin: 2 }} />

      <Button
        style={{ padding: 0 }}
        disabled={!resubmittable}
        onClick={() => resubmitJob(id)}
        type='link'>
        重新提交
      </Button>

      {/* <Divider type='vertical' style={{ margin: 2 }} /> */}

      {/* <Dropdown
        menuContentList={[
          showContinuousRedeploy && {
            children: (
              <Button
                disabled={
                  display_back_state !== 2
                    ? '文件未回传完成无法续算提交'
                    : false
                }
                style={{
                  color: 'rgba(0,0,0,0.65)'
                }}
                onClick={() => goToContinuousRedeployPage(id)}
                type='link'>
                续算提交
              </Button>
            )
          },
          {
            children: (
              <Button
                style={{ color: 'rgba(0,0,0,0.65)' }}
                onClick={() => downloadJob()}
                disabled={!downloadable}
                type='link'>
                下载至本地
              </Button>
            )
          },
          {
            children: (
              <Button
                style={{ color: 'rgba(0,0,0,0.65)' }}
                onClick={() => downloadJobToCommon()}
                disabled={!downloadable}
                type='link'>
                下载至我的文件
              </Button>
            )
          },

          {
            children: (
              <Button
                style={{ color: 'red' }}
                onClick={() => deleteJob()}
                disabled={!deleteable}
                type='link'>
                删除
              </Button>
            )
          }
        ].filter(Boolean)}
      /> */}
    </StyledLayout>
  )
})
