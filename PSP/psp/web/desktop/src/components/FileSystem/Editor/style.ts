import styled from 'styled-components'

export const StyledEditor = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 10px;

  .header {
  }

  .editorMain {
    flex: 1;
    border: 1px solid ${props => props.theme.borderColor};
    position: relative;

    .loading {
      display: flex;
      position: absolute;
      left: 0;
      right: 0;
      top: 0;
      bottom: 0;
      z-index: 999;

      > * {
        margin: auto;
      }
    }
  }

  .footer {
    display: flex;

    .footerMain {
      margin: 10px 0;
      margin-left: auto;

      button {
        margin: 0 10px;
      }
    }
  }
`
