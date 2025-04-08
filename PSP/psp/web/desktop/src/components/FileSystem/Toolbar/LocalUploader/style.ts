import styled from 'styled-components'

export const StyledUploadMenu = styled.div`
  display: flex;
  flex-direction: column;
  background-color: white;
  border: 1px solid ${props => props.theme.primaryColor};

  .ant-upload.ant-upload-select {
    width: 100%;
  }

  .uploadItem {
    display: inline-block;
    width: 100%;
    padding: 5px 10px;
    cursor: pointer;

    &.first {
      border-bottom: 1px dashed ${props => props.theme.primaryColor};
    }

    &:hover {
      background-color: ${props => props.theme.primaryColor};
      color: white;
    }
  }
`
