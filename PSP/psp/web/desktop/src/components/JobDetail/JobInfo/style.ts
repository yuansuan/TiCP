import styled from 'styled-components'

export const Wrapper = styled.div`
  height: 100%;
  padding: 20px 0;

  > header > .title {
    margin-bottom: 10px;
    font-family: PingFangSC-Semibold;
    font-size: 16px;
    line-height: 22px;
    color: #262626;
  }

  > header > .subTitle {
    display: flex;
    margin: 0 18px;

    > div {
      width: 25%;
      font-family: PingFangSC-Regular;
      font-size: 12px;
      line-height: 17px;
      .key {
        color: rgba(0, 0, 0, 0.65);
      }
      .val {
        color: rgba(0, 0, 0, 0.8);
      }
    }
  }

  > .buttonWrapper {
    margin: 14px 0;
  }
`
