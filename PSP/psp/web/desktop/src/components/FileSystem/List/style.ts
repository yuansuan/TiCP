import styled from 'styled-components'

export const StyledPanel = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;

  .list {
    flex: 1;

    .tableRow {
      .nameCell {
        display: flex;
        height: 100%;
      }

      .typeCell {
        width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .operators {
        .item {
          cursor: pointer;

          &:hover {
            color: ${props => props.theme.primaryHighlightColor};
          }
        }

        .disabled {
          color: #bfbfbf;
        }
      }
    }
  }

  .footer {
    display: flex;
    height: 50px;
    border-top: 1px solid ${props => props.theme.borderColor};

    .right {
      display: flex;
      height: 100%;
      margin-left: auto;
      margin-right: 0;

      .info {
        margin: auto;
        margin-right: 20px;
        color: ${props => props.theme.primaryColor};
      }

      .refresh {
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 0 10px;
        font-size: 16px;
        color: ${props => props.theme.primaryColor};
        border-left: 1px solid ${props => props.theme.borderColor};
        cursor: pointer;

        &:hover {
          color: ${props => props.theme.primaryHighlightColor};
        }
      }
    }
  }
`
