import styled from 'styled-components'

export const ChartTitle = styled.div`
  text-align: center;
  font-size: 14px;
  color: #000000;
  font-weight: normal;
  margin-bottom: 20px;
`

export const ChartActionWrapper = styled.div`
  display: flex;
  justify-content: flex-end;
  align-items: center;
`

export const SwitchKeyWrapper = styled.div`
  display: flex;
  width: calc(100% - 20px);
  margin-left: 40px;
  flex-wrap: wrap;
  justify-content: flex-start;
  align-items: center;

  span {
    padding: 2px 5px;
    cursor: pointer;
    color: #888;
    &.active {
      color: blue;
    }
  }
`
