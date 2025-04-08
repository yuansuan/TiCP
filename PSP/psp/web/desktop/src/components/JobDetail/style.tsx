import styled from 'styled-components'

export const TabAreaWrapper = styled.div`
  padding: 0 20px;
  overflow: auto;

  .ant-tabs-bar {
    border-bottom: 0;
  }

  .tabArea {
    height: calc(100vh - 380px);
    overflow-y: auto;
  }
`
