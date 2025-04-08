import styled from 'styled-components'

export const StyledTimeMachine = styled.div`
  display: flex;
  align-items: center;

  .skip {
    margin-right: 8px;

    button {
      border: 1px solid #bfbfbf;
      color: #595959;
      width: 32px;
      height: 32px;
      line-height: 20px;
      text-align: center;
      margin: 0 3px;
      padding: 0;
    }
  }

  .controller {
    display: flex;
    align-items: center;
    background: #ffffff;
    border: 1px solid #bfbfbf;
    border-radius: 4px;
    flex: 1;
    height: 32px;
    overflow: hidden;
    margin-right: 34px;

    .breadcrumb {
      margin: auto 10px;
      font-size: 12px;
      display: flex;
      align-items: center;

      .backList {
        border-right: 1px solid ${props => props.theme.borderColor};
        margin-right: 5px;
        padding-right: 5px;

        .listIcon {
          cursor: pointer;

          &:hover {
            color: ${props => props.theme.primaryHighlightColor};
          }
        }
      }
    }

    .collect {
      margin: auto;
      margin-right: 10px;
      cursor: pointer;
    }
  }

  .filter {
    margin-left: auto;
    margin-right: 0;
  }
`

export const StyledBackList = styled.div`
  ul {
    list-style: none;
    max-height: 200px;
    overflow: auto;

    li {
      padding: 4px 14px;
      line-height: 20px;
      cursor: pointer;

      &:hover {
        background-color: ${props => props.theme.backgroundColor};
      }
    }
  }
`
