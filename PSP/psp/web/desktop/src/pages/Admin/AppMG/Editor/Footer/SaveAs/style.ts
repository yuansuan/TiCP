import styled from 'styled-components'

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;

  .body {
    flex: 1;

    .error {
      color: red;
    }
  }

  .footer {
    display: flex;

    .footerMain {
      margin-left: auto;

      button {
        margin: 0 4px;
      }
    }
  }
`
