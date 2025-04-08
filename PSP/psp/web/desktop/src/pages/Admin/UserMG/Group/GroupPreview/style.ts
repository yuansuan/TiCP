import styled from 'styled-components'

export const GroupDetail = styled.div`
  display: flex;
  flex-direction: column;
  padding: 20px 50px;
  font-size: 16px;

  .subInfo {
    opacity: 0.65;
    font-size: 14px;
    color: rgba(0, 0, 0, 0.85);
    letter-spacing: 0.88px;
  }

  .action {
    margin-left: 10px;
  }

  .groupName {
    color: rgba(0, 0, 0, 0.85);
    margin-bottom: 15px;
  }

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
