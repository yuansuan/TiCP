import styled from 'styled-components'

interface WrapperProps {
  bigBg: string
  smallBg: string
  logo: string
}

export const Wrapper = styled.div<WrapperProps>`
  position: relative;
  min-width: 1000px;
  min-height: 700px;
  height: 100vh;
  display: flex;
  background: url(${props => props.bigBg}) no-repeat;
  background-size: cover;
  overflow: hidden;

  .footerBox {
    position: absolute;
    bottom: 30px;
    left: 50%;
    transform: translateX(-50%);
    .ysLogo {
      display: flex;
      justify-content: center;
      margin-bottom: 10px;
      .logo {
        width: 60px;
        height: 20px;
        background: url(${props => props.logo}) no-repeat;
        background-size: contain;
      }
    }
    .copyRight {
      font-family: Helvetica;
      font-size: 12px;
      color: #ffffff;
    }
  }

  .centerBox {
    width: 900px;
    display: flex;
    height: auto;
    margin: auto;
    background: #fff;
    box-sizing: content-box;
    border-radius: 4px;
    box-shadow: 0 0 6px #ccc;

    .loginBox {
      padding: 20px 20px 60px 20px;
      margin: 30px 20px;
      flex: 1.25;
      flex-direction: column;

      .title {
        text-align: left;
        p {
          color: #1b4b92;
          font-size: 36px;
          line-height: 36px;
          font-family: PingFangSC-Light;
          margin-bottom: 0;
        }
        h2 {
          text-align: right;
          padding-right: 45px;
          color: #1b4b92;
          font-size: 45px;
          margin-bottom: 0;
        }
      }

      .form {
        padding: 20px 0;
        & > .field {
          display: flex;
          padding: 15px 0;

          label {
            width: 70px;
            text-align: right;
            padding-right: 10px;
            line-height: 40px;
            align-self: flex-end;
          }
          .ant-input-password{
            width: 340px;
          }
          input {
            width: 340px;
          }
        }

        & > div {
          .remember {
            margin-left: 75px;
            display: none;
          }

          .loginBtn {
            width: 340px;
            height: 40px;
            margin-top: 15px;
            margin-left: 70px;
          }
        }
      }
    }

    .bgBox {
      flex: 1;
      background: url(${props => props.smallBg}) no-repeat;
      background-size: cover;
      border-top-left-radius: 4px;
      border-bottom-left-radius: 4px;
    }
  }
`
