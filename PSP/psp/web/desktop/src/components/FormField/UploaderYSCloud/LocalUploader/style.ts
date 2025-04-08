import styled from 'styled-components'

export const UploadMenu = styled.div`
  display: flex;
  flex-direction: column;
  background-color: white;
  border: 1px solid #d1d1d1;

  .upload {
    width: 100%;

    &:first-child {
      border-bottom: 1px dashed #d1d1d1;
    }

    .ant-upload {
      width: 100%;

      .uploadItem {
        display: inline-block;
        width: 100%;
        text-align: center;
        padding: 5px 10px;
        cursor: pointer;

        &:hover {
          background-color: #d1d1d1;
        }
      }
    }
  }
`
