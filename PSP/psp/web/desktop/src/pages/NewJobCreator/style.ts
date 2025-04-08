/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const PageWrapper = styled.div`
  .areaSelectWrap {
    display: none;
    padding: 10px 20px;
    border-bottom: 6px solid #f5f5f5;
    background: #fff;
    > div {
      display: flex;
      align-items: center;
      /* h3 {
      margin-bottom: 0;
      font-weight: normal;
    } */
    }
  }
  position: relative;
  padding: 0 0 60px 0;
  margin: 20px 20px 0 20px;
  min-width: 1000px;
  height: 100%;
  // overflow: auto;

  .input-content {
    height: 100%;
    overflow: auto;
    padding-bottom: 50px;
    > .moduleUpload {
      > .section-header {
        > .right {
          margin-left: 0;
          > .areaSelectWrap {
            padding: 0;
            border-bottom: none;
          }
        }
      }
    }
  }

  .calculate_selection {
    padding: 3px 0;

    > .ant-checkbox-wrapper {
      margin-right: 12px;
    }
  }

  .apps {
    display: grid;
    grid-template-columns: repeat(auto-fill, 300px);
    grid-row-gap: 20px;
    grid-column-gap: 20px;
  }

  .field {
    display: flex;
    align-items: center;
    padding: 12px 0;
    position: relative;

    .label {
      > .required {
        color: red;
      }

      width: 250px;
      text-align: right;
      margin-right: 14px;
    }

    .ant-input-number-handler-wrap {
      display: none;
    }

    > span {
      position: absolute;
      left: 440px;
      font-family: PingFangSC-Regular;
      font-size: 14px;
      color: rgba(0, 0, 0, 0.45);
      line-height: 22px;
    }

    > .valid-core {
      left: 265px;
      bottom: -8px;
      position: absolute;
      font-size: 12px;
      font-family: PingFangSC-Regular;
      color: #f5222d;
      letter-spacing: 0;
    }
    > .suggested-cores {
      color: rgba(0, 0, 0, 0.45);
      position: relative;
      left: 1em;
    }
  }

  .remaining-resource {
    display: none;
    width: 200px;
    margin-left: 262px;
  }

  .field-description {
    > .post-text {
      display: flex;
      padding: 5px 11px;
      align-items: center;
      min-height: 32px;
      font-size: 13px;
      color: rgba(0, 0, 0, 0.65);
      font-family: PingFangSC-Regular;
      opacity: 0.7;
      background: rgba(0, 0, 0, 0.05);
      border: 1px solid rgba(0, 0, 0, 0.1);
      border-radius: 2px;
      margin-left: 264px;
      transform: translateY(-3px);
    }
  }
`
