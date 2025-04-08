import styled from 'styled-components'
const select = require('@/assets/images/select.png')

export const Layout = styled.div`
  display: flex;
  width: 100%;
  padding: 20px;
`

export const BillingLayout = styled.div`
  display: flex;
  width: 100%;
  padding: 20px 20px 0;
  align-items: center;
`

export const ArrearsWrapper = styled.div`
  border-radius: 2px;
  padding: 16px 24px;
  box-shadow: unset;
  background-color: #fffbe6;
  border: 1px solid #ffe58f;

  .icon {
    color: #f9bf02;
    font-size: 16px;
    padding-right: 6px;
  }
`

export const Sider = styled.div`
  width: 15%;
`

export const SiderTitle = styled.span`
  font-weight: bold;
  font-size: 14px;
`

export const Content = styled.div<{ chargeType?: number }>`
  flex-grow: 1;
  width: 85%;

  .cost {
    font-family: PingFangSC-Medium;
    font-size: 16px;
    color: #e37a41;
    padding-right: 16px;
  }

  .bottomAffix {
    width: 100%;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    height: 80px;
    line-height: 80px;
    border-top: 1px solid #dbe3e4;
    background-color: #fff;
    box-shadow: 0 -4px 4px -2px #e4e9f0;
  }

  .validate_tip {
    margin-top: 14px;
    color: red;
    font-size: 14px;
    font-family: PingFangSC;
    line-height: 14px;
  }
  .errtips {
    animation: errtips 0.8s linear;
  }
  @keyframes errtips {
    10% {
      transform: translateY(1px);
    }
    20% {
      transform: translateY(-2px);
    }
    30% {
      transform: translateY(2px);
    }
    40% {
      transform: translateY(-2px);
    }
    50% {
      transform: translateY(2px);
    }
    60% {
      transform: translateY(-2px);
    }
    70% {
      transform: translateY(2px);
    }
    80% {
      transform: translateY(-2px);
    }
    90% {
      transform: translateY(2px);
    }
    100% {
      transform: translateY(0);
    }
  }

  .ant-radio-group-solid
    .ant-radio-button-wrapper-checked:not(.ant-radio-button-wrapper-disabled) {
    border: 1px solid #005dfc !important;
    background: #fff !important;
    &:after {
      content: '';
      position: absolute;
      width: 22px;
      height: 22px;
      background: url(${select}) no-repeat;
      background-size: cover;
      top: 8px;
      right: 8px;
    }
  }

  .shop-checkblock-wrap {
    display: flex;

    .shop-checkblock-item {
      ${props =>
        props.chargeType === 2
          ? ' width: 260px; height: 320px; '
          : ' width: 288px;  height: 330px;'};

      padding: 0;
      margin-right: 15px;

      &:hover {
        border: 1px solid #005dfc !important;
        box-shadow: 2px 0 8px 0 rgba(33, 150, 243, 0.11);
      }

      span {
        color: black;
        display: flex;
        flex-wrap: wrap;
        height: 100%;
        align-content: space-between;

        .main {
          padding: 10px 15px 5px;

          h3 {
            margin-bottom: 0;
            font-weight: 600;
          }

          .configure {
            .title {
              color: #777;
              line-height: 34px;
            }
            .desc {
              color: #555;
              font-size: 13px;
            }
            .info {
              line-height: 23px;
              .model {
                font-size: 13px;
                color: rgba(70, 70, 70);
              }
            }
          }

          .time {
            color: #f9924e;
            line-height: 20px;
            padding: 8px 0 5px;

            p {
              margin-bottom: 6px;
            }
          }
        }

        .usage {
          background-color: #1479ce;
          line-height: 54px;
          display: flex;
          justify-content: space-between;
          padding: 0 16px;
          width: 100%;
          color: #fff;
        }
      }
    }
  }
`

export const ComboCardWrapper = styled.div`
  overflow-x: scroll;
`
