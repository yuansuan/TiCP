import styled from 'styled-components'

export const StyledTreeSelector = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 10px;

  .header {
    font-size: 16px;
    padding: 0 10px 15px 5px;
  }

  .main {
    display: flex;

    .tree {
      width: 500px;
      height: 480px;
      overflow-y: auto;
      border: 1px solid #ddd;
      border-radius: 5px;
    }

    .preview {
      flex: 1 1 0;
      margin-left: 10px;

      .name {
        font-size: 16px;
      }

      .pic {
        margin: 10px 0;
        border-radius: 5px;
        background: #eee;
        height: 150px;
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .desc {

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
