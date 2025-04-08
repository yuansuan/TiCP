import styled from 'styled-components'

export const Wrapper = styled.div`
  position: relative;
  width: 100%;
  padding: 0;
  background: #fff;
  overflow: auto;

  .title {
    font-family: PingFangSC-Medium;
    font-size: 16px;
    color: currentColor;
    line-height: 22px;
  }

  .name {
    font-size: 18px;
    color: currentColor;
    padding: 5px 0;
  }
`

export const TabAreaWrapper = styled.div`
  padding: 0 20px;

  .ant-tabs-bar {
    border-bottom: 0;
  }
`

export const SummaryWrapper = styled.div`
  padding: 20px;
  border-bottom: 1px solid #e8e8e8;

  .header-line {
    font-size: 16px;
    display: flex;
    padding: 10px 0;
    > div {
      display: flex;
      align-items: center;
    }
  }
`
