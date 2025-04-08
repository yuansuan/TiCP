import styled from 'styled-components'

export const PanelWrapper = styled.div`
  height: 100%;

  .comment,
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

  .roleName {
    color: ${props => props.theme.primaryHighlightColor};
    cursor: pointer;
  }
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
