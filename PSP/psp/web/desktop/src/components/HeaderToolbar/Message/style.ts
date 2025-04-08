/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 10px;
`

export const StyledOverlay = styled.div`
  width: 480px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

  > .body {
    padding: 20px;

    .tabName {
      padding: 0 12px;
    }

    .ant-tabs-bar {
      margin-bottom: 0;
    }

    .ant-list-bordered {
      border: none;
      border-radius: 0;
    }
  }

  .footer {
    display: flex;
    padding: 8px 0;
    align-items: center;
    justify-content: center;
    background-color: ${({ theme }) => theme.backgroundColorBase};
    > .link {
      flex: 1;
      text-align: center;
      cursor: pointer;
      color: ${({ theme }) => theme.linkColor};

      &:hover {
        text-decoration: underline;
      }
    }
    > .read {
      flex: 1;
      text-align: center;
      cursor: pointer;
      color: ${({ theme }) => theme.linkColor};

      &:hover {
        text-decoration: underline;
      }
    }
  }
`

export const StyledPanel = styled.div`
  .ant-list-item.item {
    position: relative;

    .title {
      display: flex;
      align-items: center;
      font-size: 14px;
      margin-bottom: 4px;

      .time {
        font-size: 12px;
        font-weight: normal;
        color: #ccc;
        padding-left: 8px;
      }

      .actions {
        margin-right: 0px;
        margin-left: auto;
        position: relative;
        font-weight: normal;

        .ant-btn {
          margin: 0;
          padding: 0;
          height: 1em;
          line-height: 1em;
        }

        > .ant-btn {
          padding: 0 6px;
        }

        > * {
          position: relative;
          padding: 0 6px;

          &:not(:last-child) {
            border-radius: 0;

            &::after {
              content: '';
              position: absolute;
              right: -1px;
              top: 0px;
              height: 100%;
              border-right: 1px solid ${({ theme }) => theme.borderColorBase};
            }
          }
        }
      }
    }

    .ant-list-item-meta-avatar {
      font-size: 20px;
      color: #1890ff;
      margin-top: auto;
      margin-bottom: auto;
    }

    .ant-list-item-meta-content {
      .ant-list-item-meta-description {
        line-height: 1em;
        word-break: break-all;
        padding-right: 8px;
      }
    }

    .isRead {
      cursor: default;
      color: #ccc;
    }

    .notRead {
      cursor: pointer;
      color: #ff4d4f;
    }
  }

  .ant-list-item {
    padding: 16px 4px;
  }

  .ant-list-bordered {
    .ant-list-item {
      padding-left: 12px;
      padding-right: 4px;
    }
  }

  .ant-list-bordered {
    > .ant-list-footer {
      padding: 0;
    }
  }
`
