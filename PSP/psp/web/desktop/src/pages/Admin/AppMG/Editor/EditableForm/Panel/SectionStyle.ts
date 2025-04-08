import styled from 'styled-components'

export const Wrapper: any = styled.div`
  .section-item {
    opacity: ${(props: any) => (props.isDragging ? 0.4 : 1)};
    position: relative;
    display: flex;
    flex-direction: column;
    background-color: ${(props: any) => (props.active ? '#cdddf7' : 'inherit')};

    .section-header {
      display: flex;
      align-items: center;
      padding: 6px 20px;
      padding-left: 40px;
      cursor: default;

      .editor {
        width: 100%;
        display: flex;
        align-items: center;

        input {
          width: 160px;
        }

        > .operators {
          flex: 1;
          margin-left: 20px;

          a {
            font-size: 14px;
            margin: 0 5px;
          }
        }
      }

      .name {
        font-size: 18px;
      }

      .drag-icon {
        display: none;
      }

      &:hover {
        background-color: #cdddf7;

        > .operators {
          display: flex;
        }

        .drag-icon {
          display: block;
          position: absolute;
          left: 10px;
          top: 10px;
          z-index: 99;
          cursor: move;
        }
      }

      > .operators {
        margin-left: auto;
        cursor: default;
        display: none;

        .anticon {
          font-size: 14px;
          margin-right: 20px;
          cursor: pointer;
        }
      }
    }

    .body {
      flex: 1;
      position: relative;
      padding: 0 20px;
    }
  }
`
