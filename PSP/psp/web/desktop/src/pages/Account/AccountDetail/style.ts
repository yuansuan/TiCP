/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const InfoBlock = styled.div`
  background: #fff;
  margin-bottom: 20px;

  .header {
    font-size: 14px;
    border-bottom: 1px solid ${props => props.theme.borderColorBase};
    padding: 16px 20px;

    .title {
      font-weight: 500;
      position: relative;
      font-family: PingFangSC-Semibold;
      font-size: 16px;
      color: #333333;
    }
  }

  .info {
    display: grid;
    grid-template-columns: repeat(auto-fill, 246px);
    grid-row-gap: 12px;
    overflow-x: hidden;
    padding: 16px 20px;
    font-size: 14px;
  }

  .list {
    display: flex;
    padding: 0px 20px 16px;
    overflow-x: scroll;

    .combo {
      height: 150px;
      width: 250px;
      padding: 10px;
      margin-right: 6px;
      border: 1px solid #dbe3e4;
      border-radius: 2px;
    }
  }
`

export const InfoItem = styled.div`
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  align-items: center;
  max-width: 240px;
`
