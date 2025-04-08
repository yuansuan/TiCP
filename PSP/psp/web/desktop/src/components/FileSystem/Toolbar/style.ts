import styled from 'styled-components'

export const StyledToolbar = styled.div`
  display: flex;

  .operator {
    display: flex;
    margin: auto 0;

    > * {
      margin-right: 14px;
    }

    .operatorGroup {
      display: none;

      &.active {
        display: inline-block;
      }
    }
  }
`
