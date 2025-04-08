/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import PageLayout from '@/components/PageLayout'
import styled from 'styled-components'
import { env } from '@/domain'

const StyledDiv = styled.div`
  width: 100%;

  &:hover {
    color: white;
  }

  .project-switcher {
    background: #001b43;
    padding-left: 26px;
    display: flex;
    align-items: center;
    margin-top: 0;
    margin-bottom: 8px;

    &.ant-menu-item:not(.ant-menu-item-selected):hover {
      background: #001b43;
    }

    .anticon.ysicon {
      font-size: 14px;
      margin-right: 12px;
    }
  }

  .foot-version {
    display: flex;
    flex-direction: column;
    background: #001529;
    color: white;  text-align:center;
    .foot-version-main{
    
      font-size:16px;
      font-weight:800
    }
}
  }
`

export const SiderFooter = () => {
  const {
    menuExpanded: [menuExpanded]
  } = PageLayout.useStore()

  return (
    <StyledDiv>
      {env.isKaiwu && (
        <div className='foot-version'>
          <div className='foot-version-main'>开物 v1.0.0</div>
          <div>Copyright © 2016 - 2022</div>
          <div>浙江远算科技有限公司</div>
        </div>
      )}
    </StyledDiv>
  )
}
