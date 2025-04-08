import styled from 'styled-components'

export const Wrapper = styled.div`
  display: flex;
  height: 100%;
  flex-direction: column;
  overflow-y: auto;
  padding-bottom: 50px;

  .form {
    padding: 10px;
  }

  .footer {
    position: absolute;
    width: 100%;
    padding: 10px;
    bottom: 0;
    left: 0;
  }
`
