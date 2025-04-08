import styled from 'styled-components'

export const PanelHeaderWrapper = styled.div`
  display: inline-block;
`

export const Wrapper = styled.div`
  width: 100%;

  .body {
    position: relative;
    padding: 10px;
    height: calc(100vh - 130px);
    overflow-y: auto;

    .search {
      position: fixed;
      right: 30px;
      top: 130px;
      z-index: 10;
    }

    mark {
      background: yellow;
      color: black;
    }
    .ant-collapse {
      background-color: #fff !important;
    }
    .ant-collapse-header {
      .ant-tooltip {
        color: rgba(255, 255, 0, 1) !important;

        .ant-tooltip-arrow::before {
          background-color: rgba(255, 255, 0, 1) !important;
        }
        .ant-tooltip-inner {
          color: black !important;
          background-color: rgba(255, 255, 0, 1) !important;
          box-shadow: 0 2px 8px rgb(255 255 0 / 15%) !important;
        }
      }
    }
  }

  .loading {
    text-align: center;
    border-radius: 4px;
    margin-bottom: 20px;
    padding: 30px 50px;
    margin: 20% 0;
  }
`

export const ConfigWrapper = styled.div`
  display: flex;
  flex-direction: column;

  .item {
    padding: 5px;
    display:flex;
    > p {
      margin-bottom: 0px;
    }
    .left {
      margin-right: 15px;
      display: flex;
      .value {
        display: flex;
        flex-direction: column;
      }
      .unit {
        padding: 5px;
      }
    }
    .right {
      display: flex;
    }
    .field {
      width: 300px;
      margin-left: 30px;
    }

    .textField {
      margin-left: 30px;
    }

    .msg {
      padding-left: 20px;
      margin: 0;
      color: red;
    }
  }

  .animate-charcter {
    text-transform: uppercase;
    background-image: linear-gradient(
      -225deg,
      #231557 0%,
      #44107a 29%,
      #ff1361 67%,
      #fff800 100%
    );
    background-size: auto auto;
    background-clip: border-box;
    background-size: 200% auto;
    color: #fff;
    background-clip: text;
    text-fill-color: transparent;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    animation: textclip 2s linear infinite;
    display: inline-block;
    font-size: 190px;
  }

  @keyframes textclip {
    to {
      background-position: 200% center;
    }
  }
`

export const StepWrapper = styled.div`
  width: 740px;
  margin: 0 auto;

  .text {
    margin-top: 10px;
  }

  .tips {
    line-height: 30px;
    border-radius: 3px;
    background: #e3f4ff;
    border: 1px solid #0090fa;
    color: #0090fa;
    padding: 6px;
    display: flex;

    .machineId {
      margin-left: 5px;
      width: 570px;
      margin-bottom: 0;
    }

    .btn {
      height: 30px;
    }
  }

  .license {
    margin: 15px 0;
  }

  .upload {
    display: flex;
    margin: 15px 0;

    .file {
      height: 38px;
    }
  }
`
