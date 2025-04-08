/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'
import { Button } from 'antd'

export const StyledButton: typeof Button = styled(Button)`
  &.ant-btn {
    height: 32px;
    font-size: 14px;
    padding: 0 16px;

    > .container {
      display: flex;
      height: 100%;
      align-items: center;
      justify-content: center;
    }

    .ysicon {
      font-size: 18px;
    }

    .ysicon + span,
    span + .ysicon {
      margin-left: 8px;
    }

    &.loading {
      cursor: not-allowed;
      pointer-events: none;
    }

    &.ant-btn-sm {
      height: 24px;
      padding: 0 10px;
      font-size: 12px;

      .ysicon {
        font-size: 14px;
      }
    }

    &.ant-btn-lg {
      height: 40px;
      padding: 0 20px;
      font-size: 16px;

      .ysicon {
        font-size: 24px;
      }
    }

    &.ant-btn-background-ghost:disabled {
      background-color: ${({ theme }) => theme.backgroundColorBase}!important;
    }

    &.ant-btn-secondary:not(:disabled) {
      background-color: ${({ theme }) => theme.secondaryColor};
      color: white;

      &.ant-btn-background-ghost {
        color: ${({ theme }) => theme.secondaryColor};
        border-color: ${({ theme }) => theme.secondaryColor};

        &:hover,
        &:active {
          color: ${({ theme }) => theme.linkColor};
          border-color: ${({ theme }) => theme.linkColor};
        }
      }

      &:hover,
      &:active {
        border-color: ${({ theme }) => theme.linkColor};
      }
    }

    &.ant-btn-cancel:not(:disabled) {
      background-color: ${({ theme }) => theme.cancelColor};
      color: white;

      &.ant-btn-background-ghost {
        color: ${({ theme }) => theme.cancelColor};
        border-color: ${({ theme }) => theme.cancelColor};

        &:hover,
        &:active {
          color: ${({ theme }) => theme.cancelHighlightColor};
          border-color: ${({ theme }) => theme.cancelHighlightColor};
        }
      }

      &:hover,
      &:active {
        background-color: ${({ theme }) => theme.cancelHighlightColor};
        border-color: ${({ theme }) => theme.cancelHighlightColor};
        color: white;
      }
    }
  }
`
