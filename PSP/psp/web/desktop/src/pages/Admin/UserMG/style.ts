import styled from 'styled-components'

export const UserMGWrapper = styled.div`
  background: white;
  width: 100%;
  height: calc(100vh - 155px);

  .loading {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    right: 0;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .body {
    padding: 0 20px;
    .ant-btn-sm {
      font-size: 14px;
    }
  }
`
