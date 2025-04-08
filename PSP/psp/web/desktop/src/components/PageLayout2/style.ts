/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'
import { Layout } from 'antd'
import { COLLAPSED_WIDTH } from './constant'

const { Header, Content } = Layout

export const StyledLayout = styled(Layout)`
  aside {
    box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
    z-index: 999;
    user-select: none;
    height: calc(100vh - 64px);
    position: fixed;
  }

  .ys-pageLayout-title {
    display: flex;
    margin: 16px 0;
    justify-content: center;
    align-items: center;

    img {
      height: 32px;
    }

    .text {
      color: #fff;
    }
  }

  .ys-pageLayout-menu-container {
    overflow: auto;
    position: relative;
    height: 100%;
    display: flex;
    flex-flow: column nowrap;
    &::-webkit-scrollbar  {
      display: none;
    }
  }

  .ys-pageLayout-sider-footer {
    width: 100%;
  }
`

export const StyledBody = styled(Layout)``

export const StyledHeader = styled(Header)`
  &.ant-layout-header {
    justify-content: left;
    background: #fff;
    display: flex;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    padding: 0 24px 0 0;
    z-index: 1;
  }

  .anticon {
    color: rgba(0, 0, 0, 0.65);
  }

  > .ys-pageLayout-breadcrumb {
    overflow: hidden;
  }

  > .ys-pageLayout-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 64px;
    background: #f5f5f5;

    > .anticon {
      font-size: 20px;
    }
  }

  > .ys-pageLayout-toolbar {
    flex: 1;
  }
`

export const StyledContent = styled(Content)`
  padding-left: ${COLLAPSED_WIDTH}px;

  &.ant-layout-content {
    min-height: calc(100vh - 64px);
  }

  &.expanded {
    padding-left: 200px;
  }

  &::-webkit-scrollbar-track {
    box-shadow: inset 0 0 6px rgba(255, 255, 255, 0.3);
    border-radius: 10px;
    background-color: white;
  }

  &::-webkit-scrollbar {
    width: 4px;
    background-color: white;
  }

  &::-webkit-scrollbar-thumb {
    border-radius: 10px;
    box-shadow: inset 0 0 6px rgba(255, 255, 255, 0.3);
    background-color: gray;
  }
`
