/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  padding: 40px;
  position: relative;

  > h3 {
    font-size: 14px;
  }

  .click-back-div {
    position: absolute;
    font-family: PingFangSC-Medium;
    font-size: 14px;
    color: rgba(0, 0, 0, 0.65);
    letter-spacing: 0;
    line-height: 22px;
    top: 20px;
    right: 20px;
    cursor: pointer;
    width: 30px;

    > span {
      position: absolute;
      top: -1px;
      right: 33px;
    }

    &:hover {
      color: ${({ theme }) => theme.primaryColor};
      transform: scale(1.2);
    }
  }
`

export const StyledItem = styled.div`
  display: flex;
  align-items: center;
  margin: 24px 0;
  font-size: 14px;

  > label {
    width: 120px;
    text-align: right;
    margin-right: 10px;
  }

  > input,
  > .text {
    width: fit-content;
    margin-right: 15px;
  }

  > .psd {
    cursor: pointer;
    color: #3182ff;
  }

  > .right {
    margin-left: 20px;
    transform: translateY(1px);

    > .edit {
      cursor: pointer;
      color: ${({ theme }) => theme.linkColor};
    }
  }
`
