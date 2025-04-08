import styled from 'styled-components'

export const Wrapper = styled.div`
  position: relative;
  width: 100%;
  padding: 0;
  background: #fff;
  overflow: auto;
  height: calc(100vh - 155px);

  .title {
    font-family: PingFangSC-Medium;
    font-size: 16px;
    color: currentColor;
    line-height: 22px;
  }

  .name {
    font-size: 18px;
    color: currentColor;
    padding: 5px;
  }
`

export const SummaryWrapper = styled.div`
  padding: 20px 20px 20px 40px;
  border-bottom: 1px solid #e8e8e8;
`

export const ContentWarpper = styled.div`
  padding: 20px;
`
export const TopWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
`
