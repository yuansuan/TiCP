import styled from 'styled-components'

export const Wrapper = styled.div`
  width: 100%;
  background-color: rgba(240, 242, 245, 1);

  .body {
    padding: 10px 20px;
    height: 100%;
    overflow: hidden;
    background: #fff;
  }
`

export const ListWrapper = styled.div`
  .actions {
    display: flex;
    justify-content: space-between;
    margin-bottom: 16px;

    .btnArea {
      display: flex;
      .btn {
        margin: 0 5px;
      }
      .btn:first-child {
        margin: 0;
      }
    }

    .filterArea {
      display: flex;
      flex-wrap: wrap;

      .item {
        padding: 5px;
        display: flex;
        align-items: center;
        justify-content: flex-start;

        .label {
          padding: 0 10px;
        }
      }
    }
  }

  .list {
    height: calc(100vh - 285px);
  }

  .footer {
    padding: 10px;
    display: flex;
    justify-content: center;
    align-self: center;
  }
`

export const FormWrapper = styled.div`
  padding: 10px;

  .item {
    display: flex;
    align-self: flex-end;
    flex-direction: column;
    width: 350px;
  }

  .formItem {
    width: 300px;
  }

  .ant-descriptions-item {
    display: flex;
  }

  .ant-descriptions-item-label {
    padding-top: 5px;
  }

  .footer {
    position: absolute;
    display: flex;
    bottom: 0px;
    right: 0;
    width: 100%;
    line-height: 64px;
    height: 64px;
    background: white;

    .footerMain {
      margin-left: auto;
      margin-right: 8px;

      button {
        margin: 0 8px;
      }
    }
  }
`
