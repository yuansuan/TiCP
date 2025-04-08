/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { buryPoint, history } from '@/utils'

const StyledLayout = styled.div`
  cursor: pointer;
  display: flex;
  align-items: center;

  .name {
    max-width: calc(100% - 40px);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    cursor: pointer;
    color: ${({ theme }) => theme.primaryColor};
  }

  .anticon {
    margin-left: 10px;
  }
`

type Props = {
  id: string
  name: string
  project_id: string
}

export const JobName = observer(function JobName({
  id,
  name,
  project_id,
}: Props) {
  function checkDetail() {
    buryPoint({
      category: '作业中心',
      action: '作业名称',
    })
    history.push({
      pathname: `/company/job/${id}`,
      search: `?project_id=${project_id}`,
    })
  }

  return (
    <StyledLayout onClick={checkDetail}>
      <div className='name' title={name}>
        {name}
      </div>
    </StyledLayout>
  )
})
