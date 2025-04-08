import styled from 'styled-components'

interface WrapperProps {
  bgImg: string
  headBgImg: string
  lineImg: string
  mapImg: string
  maskImg: string
}

export const ScreenFullMonitorWrapper = styled.div<WrapperProps>`
  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  li {
    list-style: none;
  }

  .dataVis {
    width: 100vw;
    height: 100vh;

    background: url(${props => props.bgImg}) no-repeat top center;
    background-size: 100% 100%;
  }

  header {
    position: relative;
    // 100px
    height: 100px;
    background: url(${props => props.headBgImg}) no-repeat top center;
    background-size: 100% 100%;

    h1 {
      // 38px
      font-size: 38px;
      color: #fff;
      text-align: center;
      // 100px
      line-height: 100px;
    }

    .showTime {
      position: absolute;
      top: 50%;
      // 80px
      right: 80px;
      transform: translateY(-50%);
      // 20px
      font-size: 20px;
      color: rgba(255, 255, 255, 0.7);
    }

    .fullScreen {
      position: absolute;
      top: 50%;
      // 20px
      right: 20px;
      transform: translateY(-50%);
      // 25px
      font-size: 25px;
      color: rgba(255, 255, 255, 0.7);
    }
  }

  .content {
    display: flex;
    margin: 0 auto;
    // 8px 10px
    padding: 8px 10px 0;

    .column {
      flex: 3;
    }

    .column:nth-child(2) {
      flex: 5;
      // 10px
      margin: 0 10px;
    }

    .panel {
      position: relative;
      // 310px
      height: 310px;
      // 15px 40px
      padding: 0 15px 40px;
      // 15px
      margin-bottom: 15px;
      border: 1px solid rgba(25, 186, 139, 0.17);
      background: url(${props => props.lineImg}) no-repeat
        rgba(255, 255, 255, 0.04);

      &:nth-child(3) {
        margin-bottom: 0;
      }

      h2 {
        color: rgba(255, 255, 255, 0.7);
      }

      &::before {
        position: absolute;
        top: -1px;
        left: -1px;
        content: '';
        // 10px
        width: 10px;
        height: 10px;
        border-top: 2px solid #02a6b5;
        border-left: 2px solid #02a6b5;
      }

      &::after {
        position: absolute;
        top: -1px;
        right: -1px;
        content: '';
        // 10px
        width: 10px;
        height: 10px;
        border-top: 2px solid #02a6b5;
        border-right: 2px solid #02a6b5;
      }

      .panelFooter {
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;

        &::before {
          position: absolute;
          bottom: -1px;
          left: -1px;
          content: '';
          // 10px
          width: 10px;
          height: 10px;
          border-bottom: 2px solid #02a6b5;
          border-left: 2px solid #02a6b5;
        }

        &::after {
          position: absolute;
          bottom: -1px;
          right: -1px;
          content: '';
          // 10px
          width: 10px;
          height: 10px;
          border-bottom: 2px solid #02a6b5;
          border-right: 2px solid #02a6b5;
        }
      }

      h2 {
        position: relative;
        // 48px
        height: 48px;
        // 20px
        font-size: 20px;
        // 48px
        line-height: 48px;
        font-weight: 400;
        text-align: center;

        ul {
          position: absolute;
          top: 0;
          // 10px
          right: 10px;
          // 13px
          font-size: 13px;

          li {
            float: right;
            .circle {
              // 10px
              width: 10px;
              height: 10px;
              border-radius: 50%;
              display: inline-block;
              // 5px 10px
              margin: 0 5px 0 10px;
            }

            .ok {
              background: #98c46d;
            }

            .notOk {
              background: #999;
            }
          }
        }
      }

      .chart {
        // 240px
        height: 240px;
      }

      .chart > div {
        max-width: 100%;
      }
    }

    .summary {
      background-color: rgba(101, 132, 226, 0.1);
      // 15px
      padding: 15px;

      .summaryHd {
        position: relative;
        // 82px
        height: 82px;
        border: 1px solid rgba(25, 186, 139, 0.17);

        ul {
          display: flex;
          li {
            position: relative;
            flex: 1;
            // 80px
            line-height: 80px;
            text-align: center;
            color: #ffeb7b;
            // 70px
            font-size: 70px;

            &:nth-child(-n + 3):after {
              position: absolute;
              top: 20%;
              right: 0;
              height: 60%;
              content: '';
              border: 1px solid rgba(255, 255, 255, 0.2);
            }
          }
        }

        &::before {
          position: absolute;
          top: -1px;
          left: -1px;
          content: '';
          // 30px
          width: 30px;
          // 10px
          height: 10px;
          border-top: 2px solid #02a6b5;
          border-left: 2px solid #02a6b5;
        }

        &::after {
          position: absolute;
          bottom: -1px;
          right: -1px;
          content: '';
          // 30px
          width: 30px;
          // 10px
          height: 10px;
          border-bottom: 2px solid #02a6b5;
          border-right: 2px solid #02a6b5;
        }
      }

      .summaryBd {
        ul {
          display: flex;
          li {
            flex: 1;
            // 40px
            height: 40px;
            line-height: 40px;
            text-align: center;
            color: rgba(255, 255, 255, 0.7);
            // 20px
            font-size: 20px;
            // 10px
            padding-top: 10px;
          }
        }
      }
    }

    .cluster {
      position: relative;
      // 468px
      height: 468px;
      // 15px
      margin: 15px 0;

      .sphere {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        // 388px
        width: 388px;
        height: 388px;
        background: url(${props => props.mapImg});
        background-size: 100% 100%;
        opacity: 0;
      }

      .mask {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        // 468px
        width: 468px;
        height: 468px;
        background: url(${props => props.maskImg});
        background-size: 100% 100%;
        opacity: 0.5;
        animation: rotate1 20s linear infinite;
      }

      .smallBox {
        position: absolute;
        // 230px
        width: 230px;
        // 100px
        height: 100px;
        // 35px
        padding: 0 35px;

        h3 {
          // 35px
          height: 35px;
          // 20px
          font-size: 20px;
          // 35px
          line-height: 35px;
          text-align: center;
          color: rgba(255, 255, 255, 0.7);
          // 2px
          border-bottom: 2px solid #02a6b5;
        }

        ul {
          // 5px
          margin-top: 5px;

          li {
            // 25px
            height: 25px;
            // 16px
            font-size: 16px;
            // 25px
            line-height: 25px;
            color: rgba(255, 255, 255, 0.7);

            span {
              // 18px
              font-size: 18px;
              color: #ffeb7b;
            }
          }
        }
      }

      .leftBottom {
        left: 0;
        bottom: 0;
      }

      .rightBottom {
        right: 0;
        bottom: 0;
      }

      @keyframes rotate1 {
        from {
          transform: translate(-50%, -50%) rotate(0deg);
        }
        to {
          transform: translate(-50%, -50%) rotate(360deg);
        }
      }
    }
  }
`
