/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { useDispatch } from 'react-redux'
import { buryPoint, history } from '@/utils'
import { Tooltip } from 'antd'
import { Modal } from '@/components'
import { Monitors } from './Monitors'
import { LineChartOutlined } from '@ant-design/icons'
import { useStore } from '../../store'

const StyledLayout = styled.div`
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;

  .name {
    max-width: calc(100% - 40px);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: ${({ theme }) => theme.primaryColor};
  }

  .icons {
    display: flex;
    flex: 0 0 48px;
    align-items: center;
    justify-content: flex-start;
    .icon {
    }
  }

  .anticon {
    margin-left: 10px;
  }
`

type Props = {
  id: string
  name: string
  isCloud: boolean
  display_state: number
  residualVisible: boolean
  monitorVisible: boolean
  cloudGraphicVisible?: boolean
  projectId: string
  userId: string
  jobRuntimeId: string
  onClick?: (id: string) => void
}

export const JobName = observer(function JobName({
  id,
  name,
  isCloud,
  display_state,
  residualVisible,
  monitorVisible,
  cloudGraphicVisible,
  projectId,
  userId,
  jobRuntimeId,
  onClick
}: Props) {
  const title = '可视化分析'
  const store = useStore()
  const dispatch = useDispatch()

  function onJobClick() {
    buryPoint({
      category: '作业管理',
      action: '作业名称'
    })
    if (onClick) {
      onClick(id)
    } else {
      window.localStorage.setItem(
        'CURRENTROUTERPATH',
        `/new-job?jobId=${id}&isCloud=${isCloud}`
      )
    }
    dispatch({
      type: 'NEWJODETAIL',
      payload: 'close'
    })

    setTimeout(() => {
      dispatch({
        type: 'JOBMANAGE',
        payload: 'full'
      })
      dispatch({
        type: 'NEWJODETAIL',
        payload: 'full'
      })
    }, 100)
  }

  function showMonitor() {
    Modal.show({
      title,
      width: 1000,
      bodyStyle: {
        paddingTop: 0
      },
      destroyOnClose: display_state === 1, // 运行中 job，关闭对话框, 需要 destory modal
      footer: null,
      content: () => {
        return (
          <Monitors
            id={id}
            userId={userId}
            jobRuntimeId={jobRuntimeId}
            jobState={display_state}
            residualVisible={residualVisible}
            // monitorVisible={monitorVisible}
            cloudGraphicVisible={cloudGraphicVisible}
          />
        )
      }
    })
  }

  return (
    <StyledLayout>
      <div className='name' title={name} onClick={onJobClick}>
        {name}
      </div>
      <div className='icons'>
        {(residualVisible || cloudGraphicVisible || monitorVisible) && (
          <Tooltip title={title} className='residual-plot'>
            <LineChartOutlined
              style={{ fontSize: 18, flex: '0 0 24px' }}
              onClick={showMonitor}
            />
          </Tooltip>
        )}
      </div>
    </StyledLayout>
  )
})
