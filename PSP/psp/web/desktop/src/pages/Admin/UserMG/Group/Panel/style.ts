import styled from 'styled-components'

export const PanelWrapper = styled.div`
  height: 100%;

  .groupName {
    display: flex;
    align-items: center;
    .name {
      margin-left: 6px;
      width: 95%;
      display: inline-block;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    color: ${props => props.theme.primaryHighlightColor};
    cursor: pointer;
  }
`

export const OperatorsWrapper = styled.div`
  display: flex;
  align-items: center;

  .action {
    margin-right: 40px;
    cursor: pointer;

    > span {
      margin-left: 3px;
    }

    &:hover {
      color: ${props => props.theme.primaryHighlightColor};
    }
  }
`
