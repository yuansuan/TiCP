/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  width: 100%;
  min-height: 28px;

  &:hover {
    > .operator {
      .edit.hoverAction {
        visibility: visible;
      }
    }
  }

  > .operator {
    padding: 0 5px;
    box-sizing: border-box;
    display: inline-block;
    vertical-align: middle;

    > * {
      vertical-align: middle;
    }

    > .confirm {
      color: ${({ theme }) => theme.successColor};
    }

    > .cancel {
      color: ${({ theme }) => theme.errorColor};
    }

    > .help {
      color: #ccc;
    }

    .edit {
      cursor: pointer;
      margin-left: 8px;

      .anticon {
        vertical-align: middle;
      }

      &:hover {
        color: ${({ theme }) => theme.linkColor};
      }

      &.hoverAction {
        visibility: hidden;
      }
    }

    svg {
      margin: 0 3px;
      cursor: pointer;

      &:hover {
        transform: scale(1.2);
      }
    }
  }

  > .unit {
    display: inline-block;
    vertical-align: middle;
  }

  > .main {
    display: inline-block;
    vertical-align: middle;

    .ant-input-affix-wrapper {
      vertical-align: middle;

      &.error {
        border-color: ${({ theme }) => theme.errorColor};

        .ant-input-suffix {
          color: ${({ theme }) => theme.errorColor};
        }
      }
    }

    input {
      width: 100%;
      vertical-align: middle;
    }

    .text {
      width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      cursor: pointer;
      vertical-align: middle;

      &.isLink {
        color: ${({ theme }) => theme.linkColor};
      }
    }
  }
`
