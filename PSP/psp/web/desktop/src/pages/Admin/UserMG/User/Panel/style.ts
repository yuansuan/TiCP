import styled from 'styled-components'

export const PanelWrapper = styled.div`
  height: 100%;

  .userName {
    display: flex;

    .name {
      width: 95%;
      margin-left: 6px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    color: ${props => props.theme.primaryHighlightColor};
    cursor: pointer;
  }

  .roleName {
    display: flex;
    align-items: center;

    .name {
      width: 95%;
      margin-left: 6px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
`

export const StateWrapper = styled.div`
  margin: 0 10px;
`

export const OperatorsWrapper = styled.div`
  display: flex;
  align-items: center;

  .action {
    margin-right: 10px;
    cursor: pointer;

    > span {
      margin-left: 3px;
    }

    &:hover {
      color: ${props => props.theme.primaryHighlightColor};
    }
  }

  .disabled {
    color: #ccc;
    cursor: not-allowed;

    &:hover {
      color: #ccc;
    }
  }
`
