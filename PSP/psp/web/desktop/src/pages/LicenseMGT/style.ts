/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  .header {
    padding: 16px 0;
    > h3 {
      margin-bottom: 0;
    }
  }
  .main {
    padding: 14px 20px;
    .toolbar {
      display: flex;
      margin-bottom: 18px;
      .left {
        > button {
          margin-right: 8px;
        }
      }
      .right {
        margin-left: auto;
      }
    }
    .body {
      .ant-form-item {
        margin-bottom: 0;
      }
    }
  }
`

export const FormWrapper = styled.div`
  .ant-input-number {
    width: 100%;
  }
`
