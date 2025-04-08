import styled from 'styled-components'

export const PanelWrapper = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;

  .filter {
    height: 60px;
    line-height: 60px;
    margin: auto;

    .ant-input-search {
      height: 32px;
      width: 235px;
    }
  }

  .itemList {
    flex: 1;
    overflow: auto;
    border-top: 1px solid #d8d8d8;

    .item {
      line-height: 45px;
      padding-left: 15px;
      padding-right: 15px;

      &:hover {
        background: #f2f6ff;
      }

      & > span {
        width: 230px;
        margin-left: 25px;
        display: inline-block;
        max-width: 360px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        vertical-align: bottom;
      }
    }
  }
`
