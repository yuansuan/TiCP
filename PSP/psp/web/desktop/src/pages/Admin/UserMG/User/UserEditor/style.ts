import styled from 'styled-components'

export const UserEditorWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 50px 10px;
  font-size: 16px;

  .email {
    display: inline-block;
    vertical-align: bottom;
    max-width: 200px;
    min-width: 100px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
`

export const BasicInfoWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 30px;
  height: 32px;
  color: rgba(0, 0, 0, 0.85);

  span {
    margin-right: 10px;
  }

  svg {
    margin-left: 10px;
  }

  input {
    width: 150px;
  }
`

export const EditWrapper = styled.div`
  width: 300px;

  .mgRight {
    width: 200px;
  }
`
