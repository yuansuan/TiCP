import styled from 'styled-components'

export const UserEditorWrapper = styled.div`
  position: relative;
  padding: 20px 50px;
  display: flex;
  flex-direction: column;
  font-size: 16px;

  .Softwares {
    .special {
      .rs-table-cell,
      .rs-table-row-header .rs-table-cell {
        background-color: #f0f5fd;
      }

      .rs-table-row {
        border-bottom-color: rgba(109, 114, 120, 0.25);
      }

      .rs-table-row-header {
        border-bottom-color: rgba(109, 114, 120, 0.85);
      }
    }
  }
`

export const BaseInfoWraper = styled.div`
  .ant-descriptions-item-label {
    font-size: 16px;
  }
  .ant-descriptions-item-content {
    font-size: 16px;
  }
`

export const StyledLoading = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
`
