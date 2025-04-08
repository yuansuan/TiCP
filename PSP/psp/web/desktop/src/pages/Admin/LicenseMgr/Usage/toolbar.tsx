import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Input } from 'antd'
import { Button } from '@/components'
import { lmUsageList } from '@/domain'

const StyledLayout = styled.div`
  display: flex;
  justify-content: space-between;
  > .left {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    > .company {
      width: 100px;
    }
  }
  > .right {
    display: flex;
    width: 130px;
    justify-content: space-between;
  }
`

export const Toolbar = observer(function toolbar(refresh: any) {
  return (
    <StyledLayout>
      <div className='left'>
        <div className='company'>企业名称：</div>
        <Input
          value={lmUsageList?.company}
          placeholder={'请输入企业名称'}
          onChange={e => lmUsageList.setFilterParams(e.target.value)}
        />
      </div>
      <div className='right'>
        <Button
          onClick={() => {
            refresh
          }}
          type='primary'>
          查询
        </Button>
        <Button
          onClick={() => {
            lmUsageList.setFilterParams('')
            refresh
          }}>
          重置
        </Button>
      </div>
    </StyledLayout>
  )
})
