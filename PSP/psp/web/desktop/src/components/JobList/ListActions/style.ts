import styled from 'styled-components'

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  margin-bottom: 11px;
`

export const ActionWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
  height: 32px;
`
export const SearchWrapper = styled.div`
  display: flex;
  align-items: center;
`

export const MoreFilterWrapper = styled.div`
  margin: 5px 0;
  background: #eee;
  display: flex;
  flex-wrap: wrap;
  .item {
    display: flex;
    justify-content: flex-end;
    margin: 0px 5px;
    padding: 5px;
    .label {
      line-height: 32px;
      margin: 0px 5px;
    }
  }
`
