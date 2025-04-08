import styled from 'styled-components'

export const EditableNameWrapper = styled.div`
  display: inline-block;

  input {
    width: 130px;
  }

  .editor {
    display: flex;
    align-items: center;

    .value {
      width: 80%;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      cursor: pointer;

      &.isLink {
        &:hover {
          color: $app-primary-color;
        }
      }
    }

    .edit {
      cursor: pointer;
      margin-left: 5px;

      &:hover {
        color: $app-primary-color;
      }
    }
  }
`
