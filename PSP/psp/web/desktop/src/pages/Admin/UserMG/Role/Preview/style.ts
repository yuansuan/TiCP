import styled from 'styled-components'

export const RolePreviewWrapper = styled.div`
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

  .roleName {
    margin-bottom: 30px;
  }
`

export const BasicInfoWrapper = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: flex-start;
  margin-bottom: 10px;

  .header {
    font-size: 16px;
    color: rgba(0, 0, 0, 0.85);
    width: 80px;
    text-align: right;
  }

  .roleName {
    margin-bottom: 30px;
  }

  .body {
    width: 800px;
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
