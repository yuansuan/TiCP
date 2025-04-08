/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;
  .ant-page-header-heading-title {
    max-width: calc(100% - 100px);
    overflow: hidden;
    text-overflow: ellipsis;
  }
`

export const StyledHeader = styled.div`
  background-color: white;
  padding: 16px 24px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
`

export const StyledContent = styled.div`
  flex: 1;
  margin: 0;
  background-color: white;
  padding: 5px 20px;
  .ant-tabs-nav-wrap {
    background: white;
    border-bottom: 1px solid #eee;
  }
  .ant-tabs.ant-tabs-card .ant-tabs-card-bar {
    .ant-tabs-tab {
      border-radius: 6px 6px 0 0;
      width: 88px;
      height: 44px;
      text-align: center;
    }
  }

  .ant-tabs-bar {
    margin: 0;
  }
`
