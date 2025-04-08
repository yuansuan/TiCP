import styled from 'styled-components'

interface IProps {
  isFullScreen: boolean
  ref: any
}

export const Wrapper = styled.div<IProps>`
  --base-size: calc(100vw / 80);
  --base-height: calc(100vh - var(--base-size) * 3);
  --base-row: calc(var(--base-height) * 0.33 - var(--base-size));
  --first-row: calc(var(--base-height) * 0.4 - var(--base-size));
  --other-row: calc(var(--base-height) * 0.3 - var(--base-size));

  width: 100%;
  padding: ${props =>
    props.isFullScreen
      ? '0 var(--base-size) var(--base-size) var(--base-size)'
      : '0px 20px 15px 20px'};
  background-color: rgba(240, 242, 245, 1);

  .pageTitle {
    display: flex;
    justify-content: space-between;

    font-size: ${props => (props.isFullScreen ? 'var(--base-size)' : '16px')};
    line-height: ${props =>
      props.isFullScreen ? 'calc(var(--base-size)*3)' : '48px'};

    .right {
      display: flex;
      align-items: center;
      justify-content: space-between;
      .btn {
        padding: ${props =>
          props.isFullScreen ? '0 calc(var(--base-size)*0.6)' : '0 10px'};
        font-size: ${props =>
          props.isFullScreen ? 'var(--base-size)' : '16px'};
      }
      .online {
        font-size: ${props =>
          props.isFullScreen ? 'calc(var(--base-size)*0.8)' : '14px'};
        color: #999999;
        font-weight: 500;

        .num {
          font-size: ${props =>
            props.isFullScreen ? 'calc(var(--base-size)*0.8)' : '14px'};
          color: #1a6eba;
          cursor: pointer;
        }
      }
    }
  }

  .gridBody {
    display: grid;
    width: ${props =>
      props.isFullScreen
        ? 'calc(100% - var(--base-size)*2)'
        : 'calc(100% - 20px)'};
    grid-template-columns: repeat(3, 33.33%);
    grid-template-rows: ${props =>
      (props as any).isFullScreen
        ? 'calc(var(--first-row)) calc(var(--other-row)) calc(var(--other-row))'
        : '290px 250px 250px 250px'};
    grid-gap: ${props =>
      props.isFullScreen ? 'calc(var(--base-size))' : '12px'};

    .item {
      background-color: rgba(255, 255, 255, 1);
      box-sizing: border-box;
      border-width: 1px;
      border-style: solid;
      border-color: rgba(242, 242, 242, 1);
      padding: ${props =>
        props.isFullScreen ? 'calc(var(--base-size)*1.1)' : '20px'};

      .title {
        font-size: ${props =>
          props.isFullScreen ? 'var(--base-size)' : '18px'};
        color: rgba(0, 0, 0, 0.85);
        font-weight: 500;
        display: flex;
        align-items: center;
        .detail {
          font-weight: 400;
          font-size: ${props =>
            props.isFullScreen ? 'calc(var(--base-size)*0.8)' : '12px'};
          color: #999999;
          margin: 0 auto;
        }
      }
      .util {
        display: grid;
        grid-template-columns: 30% 50%;
      }
    }

    .head {
      grid-column-start: 1;
      grid-column-end: 4;
    }
    .footer {
      grid-column-start: 1;
      grid-column-end: 3;
    }
  }
`
