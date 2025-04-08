import styled from 'styled-components'

export const Wrapper = styled.div`
  width: 100%;
  padding: 20px;
  background: #fff;
  height: calc(100vh - 155px);
`

export const ListWrapper = styled.div`
  width: 100%;
  background: #fff;
  height: calc(100vh - 320px);
`

export const TopWrapper = styled.div`
  .action {
    display: flex;
    justify-content: space-between;
    margin-bottom: 5px;
    align-items: center;

    .filter {
      display: flex;
      justify-content: flex-start;
      align-items: center;

      .item {
        margin: 5px;
        display: flex;
        align-items: center;

        .label {
          flex: 1 0 80px;
          text-align: right;
          padding: 5px;
        }
      }
    }

    .btn {
      margin 0 5px;
    }
  }
`

export const FormWrapper = styled.div`
  .ant-input-number {
    width: 100%;
  }
`
