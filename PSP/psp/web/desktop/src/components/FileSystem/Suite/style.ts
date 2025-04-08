import styled from 'styled-components'

export const StyledSuite = styled.div`
  padding: 20px;
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: inherit;
  max-height: inherit;
  background-color: white;

  .header {
    .timeMachine {
      margin-bottom: 16px;
    }

    .toolbar {
      margin-bottom: 20px;
    }
  }

  .main {
    flex: 1;
    display: flex;
    border: 1px solid ${props => props.theme.borderColor};
    min-height: 0;

    .menu {
      position: relative;

      .menuWrapper {
        width: 300px;
        height: 100%;
        overflow: scroll;
      }

      .resizeBar {
        position: absolute;
        display: flex;
        width: 2px;
        height: 100%;
        right: -2px;
        top: 0;
        z-index: 3;
        cursor: ew-resize;
        background-color: ${props => props.theme.borderColor};

        .dragIcon {
          position: absolute;
          top: 48%;
          left: -6px;
          background-color: white;
          padding: 5px 0;
        }
      }
    }

    .panel {
      flex: 4;
    }
  }

  .loading {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    right: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    background: (255, 255, 255, 0.1);
  }
`
