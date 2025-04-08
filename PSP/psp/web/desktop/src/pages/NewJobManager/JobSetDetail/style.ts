/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const InfoItem = styled.div`
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  align-items: center;
  max-width: 240px;
`

export const InfoBlock = styled.div`
  padding: 0 20px;
  background: #fff;

  .header {
    height: 46px;
    line-height: 46px;
    font-size: 14px;
    color: ${props => props.theme.primaryColor};
    border-bottom: 1px solid ${props => props.theme.borderColorBase};

    .title {
      width: 120px;
      text-align: center;
      font-weight: 500;
      position: relative;

      ::after {
        content: '';
        display: block;
        width: 100%;
        height: 2px;
        background: ${props => props.theme.primaryColor};
        position: absolute;
        bottom: 0;
        left: 0;
      }
    }
  }

  .content {
    display: grid;
    grid-template-columns: repeat(auto-fill, 246px);
    grid-row-gap: 12px;
    overflow-x: hidden;
    padding: 16px;

    > .job-status {
      display: flex;
    }
  }
`

export const Wrapper = styled.div`
  padding: 20px;
`

export const JobListWrapper = styled.div`
  display: flex;
  background: #fff;
  margin-top: 10px;
  min-height: calc(100vh - 302px);

  .ant-tabs-nav .ant-tabs-tab {
    width: 88px;
    text-align: center;
  }

  .item {
    flex: 1;

    .action {
      margin: 4px 0 20px 0;
    }
  }
`
