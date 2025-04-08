import styled from 'styled-components'

export const ActionWrapper = styled.div`
  display: flex;

  .item {
    padding: 10px;

    .time {
      padding: 0 5px;
      display: inline-block;
    }
  }
`

export const PieWrapper = styled.div`
  display: flex;
  .data {
    padding: 5px;
    width: 50%;
  }

  .chart {
    width: 50%;
  }
`

export const SinglePieWrapper = styled.div`
  display: flex;
  justify-content: center;

  .data {
    padding: 5px;
    width: 50%;
  }

  .chart {
    width: 50%;
  }
`
