/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const FileListStyle = styled.div`
  width: 100%;
  max-height: 460px;
  padding: 0;
  overflow: auto;
  background-color: white;

  span.disabled {
    color: rgba(0, 0, 0, 0.25);
  }

  tr:hover td .actions {
    display: flex;
  }

  a.delete {
    color: #f5222d;
    &:hover {
      color: #f5222d;
    }
  }
`

export const FileNameStyle = styled.div`
  width: 100%;
  display: flex;
  align-items: center;

  > .anticon {
    font-size: 24px;
    margin-right: 4px;
    cursor: pointer;
    color: ${props => props.theme.primaryColor};
  }

  .icon-prefix {
    width: 24px;
  }

  .filename-text {
    width: calc(100% - 24px);
    display: flex;
    align-items: center;

    .visibleName {
      display: block;
    }

    .actions {
      margin-left: 10px;
      display: none;
      align-items: center;
    }
  }
`

export const FileProgressStyle = styled.div`
  display: flex;
  align-items: center;
  width: 100%;

  .ant-progress {
    width: 100%;
  }

  .status-icon {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    margin-right: 4px;
    position: relative;
    box-sizing: content-box;
  }
`

export const EditableNameStyle = styled.div`
  width: 100%;

  > div {
    min-height: 0;
    height: 22px;
  }
`
