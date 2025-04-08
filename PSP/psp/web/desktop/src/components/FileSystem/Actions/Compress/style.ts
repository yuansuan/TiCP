import styled from 'styled-components'

export const StyledCompressInfo = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;

  .body {
    flex: 1;

    .module {
      margin-bottom: 20px;

      .name {
        display: inline-block;
        width: 60px;
        text-align: right;
        margin-right: 5px;
      }

      .widget {
        width: 260px;
      }
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
