import styled from 'styled-components'

export const UserCreatorWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 50px 10px;

  .header {
    display: flex;
    margin: 15px 0;
    .upload-wrap {
      position: relative;
      display: inline-block;
      overflow: hidden;
      border: 1px solid #d8d8d8;
      border-radius: 4px;
      margin-right: 20px;
    }
    .upload-wrap .file-ele {
      position: absolute;
      top: 0;
      right: 0;
      opacity: 0;
      height: 100%;
      width: 100%;
      cursor: pointer;
    }
    .upload-wrap .file-open {
      width: 450px;
      height: 30px;
      line-height: 30px;
      text-align: center;
      background-color: white;
    }
    .create {
      width: 324px;
    }
  }

  .editorMain {
    display: flex;
    height: 490px;

    .content {
      width: 450px;
      overflow-y: auto;
      background-color: white;
      margin-right: 20px;
      border: 1px solid #d8d8d8;
      border-radius: 4px;
      .contentList {
        padding: 0 10px;
        display: flex;
        justify-content: space-between;
      }
    }

    .result {
      width: 324px;
      background-color: white;
      min-height: 0;
      border: 1px solid #d8d8d8;
      border-radius: 4px;
      word-break: keep-all;
    }
  }
`
export const FooterWrapper = styled.div`
  position: absolute;
  display: flex;
  bottom: 0px;
  right: 0;
  width: 100%;
  line-height: 70px;
  height: 70px;
  background: white;

  .footerMain {
    margin-left: auto;

    button {
      width: 120px;
      height: 40px;
      margin: 0 20px;
    }
  }
`
