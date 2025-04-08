import styled from 'styled-components'

export const SelectEditorWrapper = styled.div<{ isSync?: boolean }>`
  display: flex;
  flex-direction: column;
  height: 100%;

  .label {
    color: rgba(0, 0, 0, 0.85);
  }

  .module {
    display: flex;
    flex-direction: column;
    height: 500px;

    header {
      display: flex;
      align-items: center;
      height: 32px;
      margin-bottom: 10px;
    }

    .body {
      flex: 1;
      background: white;
      border: 1px solid #d8d8d8;
      border-radius: 4px;
      display: flex;
      flex-direction: column;
    }
  }

  .editorBody {
    display: flex;

    .left {
      width: 450px;
      margin-right: 20px;
      display: flex;
      flex-direction: column;

      .module {
        height: ${props => (props.isSync ? '500px' : '250px')};

        .body {
          overflow-y: auto;
          padding: 15px 10px;
        }
      }
    }

    .right {
      width: 324px;

      .module {
        .body {
          background-color: white;
          min-height: 0;

          .header {
            display: flex;
            border-bottom: 1px solid #d9d9d9;

            .title {
              margin-left: 10px;
              padding: 0;
            }

            .tab {
              flex: 1;
              cursor: pointer;
              position: relative;
              display: flex;
              justify-content: center;
              height: 53px;
              line-height: 53px;

              &.active {
                color: ${props => props.theme.primaryColor};
              }

              &.left_tab {
                border-right: 1px solid #d8d8d8;
              }
            }
          }

          .tabPanel {
            display: none;
            height: 100%;
            overflow: auto;

            &.active {
              display: flex;
              flex-direction: column;
            }
          }
        }
      }
    }
  }
`
