import styled from 'styled-components'

export const Wrapper = styled.div`
  width: 100%;

  .body {
    padding: 20px;
    height: calc(100vh - 158px);
    overflow-y: auto;
  }

  .loading {
    text-align: center;
    border-radius: 4px;
    margin-bottom: 20px;
    padding: 30px 50px;
    margin: 20% 0;
  }
`

export const ConfigWrapper = styled.div`
  display: flex;
  flex-direction: column;

  .item-row {
    padding: 5px;
    width: 300px;
    justify-content: space-between;
    display: flex;
    .label {
      margin-left: 28px;
    }
  }
`
