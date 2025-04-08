import styled from 'styled-components'

export const FormItem = styled.div`
  width: 800px;
  display: flex;
  align-items: center;
  position: relative;
  padding: 6px 0;
  margin: 4px 0;

  input {
    width: 300px;
    height: 40px;
    font-size: 14px;
  }

  .ant-select {
    width: 300px;
    font-size: 14px;

    .ant-select-selection--multi {
      height: auto;
    }

    .ant-select-selection--single {
      height: 40px;

      .ant-select-selection__rendered {
        line-height: 38px;
      }
    }
  }

  .ant-input-number-sm input {
    height: 40px;
    line-height: 40px;
  }

  textarea {
    width: 350px;
    height: 80px;
    font-size: 14px;
    resize: none;
  }
`
export const Label: any = styled.div`
  width: 180px;
  font-size: 14px;
  padding-left: 25px;
  color: rgba(38, 38, 38, 0.65);
  text-align: right;
  margin-right: 8px;
  line-height: 24px;

  .required {
    color: red;
  }

  .info {
    display: flex;
    flex-direction: column;

    .help {
      color: ${props => props.theme.primaryColor};
    }

    > .label {
      width: 100%;
      display: flex;
      align-items: center;

      .text {
        display: inline-block;
        width: calc(100% - 16px);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    > .id {
      width: 100%;
      display: flex;
      align-items: center;

      .value {
        display: inline-block;
        width: calc(100% - 16px);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
`

export const InfoText = styled.div``

export const FormItemWrapper = styled.div`
  width: 800px;
  padding: 20px 0;
  background-color: #f5f9ff;
`
export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  padding-right: 30px;

  button {
    width: 67px;
    height: 24px;
    font-size: 14px;
    margin-left: 20px;
  }
`

export const Options = styled.div`
  display: flex;
  flex-direction: column;

  .input-wrapper {
    display: flex;
    align-items: center;

    input {
      width: 306px;
      height: 24px;
      line-height: 24px;
      margin-bottom: 6px;
    }

    .right-option {
      display: none;

      span {
        display: block;
        margin-left: 10px;
        font-size: 14px;
        color: #262626;
      }

      span.active {
        color: #368eff;
      }

      .icon {
        margin-left: 20px;
      }
    }

    &:hover > .right-option {
      display: flex;
    }
  }
  .new {
    width: 306px;
    height: 24px;
    line-height: 24px;
    padding-left: 8px;
    border: 1px dashed #d9d9d9;
    border-radius: 4px;
    font-size: 14px;
    color: #d9d9d9;
  }
`
