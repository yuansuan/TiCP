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

export const BasicInfoWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  padding-right: 150px;
  font-size: 16px;
  margin-bottom: 10px;
  color: rgba(0, 0, 0, 0.85);

  & > div {
    display: flex;
    & > label {
    }
    & > span {
      margin-right: 30px;
      max-width: 240px;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
    }

    > .cert {
      cursor: pointer;
      color: #3182ff;
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
