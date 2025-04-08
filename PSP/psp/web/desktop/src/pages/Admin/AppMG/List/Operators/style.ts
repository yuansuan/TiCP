import styled from 'styled-components'

export const StyledOperators = styled.div`
  .item {
    display: inline-block;
    margin-right: 20px;
    cursor: pointer;

    &.large {
      width: 56px;
      text-align: left;
    }

    &.disabled {
      color: #bfbfbf;

      &:hover {
        color: #bfbfbf;
      }
    }

    &:hover {
      color: ${props => props.theme.primaryColor};
    }
  }
`
