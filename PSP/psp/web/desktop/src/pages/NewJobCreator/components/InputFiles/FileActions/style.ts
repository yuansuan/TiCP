/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const FileActionStyle = styled.button`
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.15);
  padding: 5px 15px;
  border-radius: 2px;
  cursor: pointer;

  display: flex;
  align-items: center;

  font-size: 14px;
  color: rgba(0, 0, 0, 0.65);

  span.anticon {
    font-size: 18px;
    margin-right: 6px;
    &.folder_mine {
      font-size: 22px;
    }
  }

  &:focus {
    outline: none;
  }

  span:first-of-type {
    display: inline-block;
  }
  span:nth-of-type(2) {
    display: none;
  }

  &:hover {
    border: 1px solid #005dfc;

    span:first-of-type {
      display: none;
    }
    span:nth-of-type(2) {
      display: inline-block;
    }
  }
`

export const FileActionsStyle = styled.div`
  display: flex;
  flex-direction: row;
  background: #ffffff;

  button {
    margin-right: 10px;
  }
`

export const AreaSelectStyle = styled.div`
  padding: 0;
  border: none;
  /* width: 100px; */
  margin-right: 10px;
  > .AreaSelectContainer {
    height: 100%;
    > .ant-select {
      height: 100%;
      width: 140px!important;
    }
  }
`
