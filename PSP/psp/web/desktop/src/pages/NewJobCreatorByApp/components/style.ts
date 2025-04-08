/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const SectionStyle = styled.div`
  background-color: white;
  padding: 0 20px;

  .section-header {
    display: flex;
    align-items: center;
    height: 46px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);

    > .left {
      display: flex;
      align-items: center;

      .section-icon {
        color: #3182ff;
        font-size: 24px;
        margin-right: 8px;
      }

      h1 {
        font-size: 14px;
        margin: 0;
        color: rgba(0, 0, 0, 0.85);
        font-weight: 600;
        line-height: 45px;
        padding: 0 10px;
      }
    }

    > .right {
      margin-left: auto;
    }
  }

  .section-content {
    padding: 16px 0 30px 0;
    color: #000;
  }
`

export const BottomActionStyle = styled.div`
  position: absolute;
  width: 100%;
  bottom: 0;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.2);
  padding: 0 48px;
  transition: all 0.2s;
  z-index: 99;
  line-height: 22px;
  > div.formCompleteness {
    display: flex;
    align-items: center;
    color: #000;
    > span.label {
      margin-right: 10px;
    }
    > span.completeness {
      display: flex;
      align-items: center;
      margin-right: 30px;
      > span {
        margin-right: 10px;
      }
    }
  }
  > div.actions {
    button {
      margin-right: 6px;
    }
  }
`

interface IStyledProps {
  marks: Object
}

export const CoreSelectorStyle = styled.div<IStyledProps>`
  .ant-slider {
    width: ${props => (Object.keys(props.marks).length > 1 ? '400px' : 0)};
  }
`

export const SummaryStyle = styled.div`
  display: flex;
  flex-direction: column;
  height: 600px;
  > .jobName {
    margin-bottom: 20px;
    color: rgba(0, 0, 0, 0.85);
    > div.label {
      display: inline-block;
      > span.star {
        color: red;
        margin-right: 5px;
      }
      margin-right: 10px;
    }
    > input {
      width: 224px;
    }
  }

  > .main {
    flex: 1 1 auto;
    display: flex;
    flex-direction: column;
    padding-bottom: 40px;
    > div.table {
      flex: 1;
    }
    > span.note {
      font-size: 14px;
      color: #999999;
    }
  }

  > .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 10px 17px 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`

const CompletenessStyle = styled.span`
  display: inline-block;
  box-sizing: border-box;
  width: 16px;
  height: 16px;
  border: 1px solid;
  border-radius: 50%;
`

export const CompletedStyle = styled(CompletenessStyle)`
  border-color: #1890ff;
  line-height: 12px;
  font-size: 12px;
  color: #1890ff;
  background-color: #ddf0ff;
  padding: 1px;
`

export const UnCompletedStyle = styled(CompletenessStyle)`
  border-color: #d9d9d9;
`

export const AlertCheckboxAlertStyle = styled.div`
  height: 54px;
  display: flex;
  align-items: center;
  > label {
    margin: 10px;
  }
`
