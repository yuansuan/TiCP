import styled from 'styled-components'

export const GroupEditorWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 50px 10px;
  font-size: 16px;

  .groupName {
    margin-bottom: 20px;

    & > span {
      color: #e02020;
      margin-right: 5px;
    }

    & > input {
      width: 300px;
      height: 40px;
    }
  }
`
