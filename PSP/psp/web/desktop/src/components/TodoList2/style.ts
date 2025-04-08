import styled from 'styled-components'

export const StyleLayout = styled.div`
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: space-between;
  font-family: 'PingFangSC-Regular';
  font-size: 12px;
  color: rgba(38, 38, 38, 0.65);
  .todo-item {
    display: flex;
    align-items: center;
    margin-bottom: 10px;
  }

  .ant-select,
  .ant-input {
    margin-right: 10px !important;
    font-size: 12px !important;
    font-family: PingFangSC-Regular !important;
    color: rgba(38, 38, 38, 0.65) !important;
  }

  .todo-input {
    width: 200px;
  }
  .delete-button {
    border: none;
    padding: 5px 10px;
    cursor: pointer;
  }

  .add-todo {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .add-item-action {
    color: #0064ff;
    margin-top: 10px;
    padding: 5px 10px;
    cursor: pointer;
    margin: 0 auto;
  }
`
