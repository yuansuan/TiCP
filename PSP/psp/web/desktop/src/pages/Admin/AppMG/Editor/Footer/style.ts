import styled from 'styled-components'

export const Wrapper = styled.div`
  position: absolute;
  bottom: 0;
  right: 0;
  padding: 10px 20px;
  background-color: white;
  display: flex;
  z-index: 999999;

  .main {
    margin-left: auto;
    button {
      margin: 0px 10px;
      min-width: 50px;
    }
  }
`
