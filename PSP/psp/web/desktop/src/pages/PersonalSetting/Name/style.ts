/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  display: flex;
  align-items: center;

  > form {
    .ant-row {
      margin: 0;
      width: fit-content;

      .ant-form-item-control-input {
        width: 200px;
      }

      .ant-form-item-control-input-content {
        display: flex;
        align-items: center;
      }

      .ant-form-item-explain {
        width: fit-content;
      }

      label {
        display: inline-block;
        width: 120px;
        text-align: right;
        color: rgba(0, 0, 0, 0.65);
        font-size: 14px;
        line-height: 30px;
        margin-right: 10px;
      }
    }

    .text {
      display: inline-block;
    }

    .right-edit {
      display: inline-block;
      margin-left: 20px;

      &:hover {
        color: ${({ theme }) => theme.primaryColor};
      }
    }

    .right-confirm {
      top: 5px;
      right: -55px;
      position: absolute;

      > .ok {
        color: #72c14b;
      }
      > .cancel {
        color: #f5222d;
      }
    }
  }
`
