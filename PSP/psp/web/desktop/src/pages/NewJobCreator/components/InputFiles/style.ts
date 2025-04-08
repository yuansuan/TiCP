/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const InputFilesStyle = styled.div`
  height: 100%;

  .toolbar {
    > * {
      margin: 0 5px;
    }
  }
`

export const ToolbarStyle = styled.div`
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 4px 0;

  .right-search {
    display: flex;
    align-items: center;

    span.label {
      font-weight: 600;
      margin-right: 10px;
      color: rgba(0, 0, 0, 0.85);
    }

    .ant-input-search, .ant-input-group-wrapper {
      width: 220px;
      margin-right: 12px;
    }
  }
`

export const FileSearchResultListStyle = styled.div`
  width: 220px;
  background: #ffffff;
  box-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.15);
  max-height: 280px;
  overflow: auto;

  .ant-checkbox-group {
    width: 100%;
  }

  li {
    height: 24px;
    line-height: 24px;
    padding: 4px 0 4px 12px;
    box-sizing: content-box;

    label > span:nth-child(2) {
      width: 155px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  .actions {
    text-align: right;
    margin-right: 12px;
    margin-bottom: 12px;
  }

  .dropdown-content {
    overflow: auto;
    max-height: 192px;
  }
`
