import styled from 'styled-components'

export const PermPreviewWrapper = styled.div`
  font-size: 16px;

  .ant-tabs-tab {
    font-size: 16px;
  }

  .ant-tabs-tabpane,
  .body {
    font-size: 16px;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    padding-left: 0;
    width: 100%;
  }

  .ant-tabs {
    width: 100%;
  }

  .special {
    height: 100%;
    width: 100%;
  }

  .header > .icon {
    color: #10398b;
  }
`

export const AuthorityWrapper = styled.div`
  display: flex;
  flex-wrap: wrap;
  padding: 0 25px;

  .appAuthority {
    width: 100%;
    display: flex;
    align-items: center;
    font-size: 16px;
    color: rgba(0, 0, 0, 0.85);
    margin-top: 30px;
    margin-bottom: 20px;
    line-height: 22px;
    height: 22px;
  }

  .bubble {
    border-radius: 15px;
    margin-right: 12px;
    margin-bottom: 15px;
    padding: 0 25px;
    height: 30px;
    font-size: 16px;
    line-height: 30px;
  }

  .enable {
    border: 1px solid #10398b;
    background-color: #10398b;
    color: #fff;
  }

  .disable {
    border: 1px solid #bfbfbf;
    background-color: #fff;
    color: #bfbfbf;
  }
`
