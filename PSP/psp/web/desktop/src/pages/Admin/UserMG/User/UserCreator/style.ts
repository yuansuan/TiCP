import styled from 'styled-components'

export const UserCreatorWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 50px 10px;
  font-size: 16px;

  .userName {
    display: flex;
    align-items: center;

    .widget {
      width: 224px;
      margin-left: 20px;
    }
  }

  .roleChooser {
    flex: 1;
  }

  .ant-select-selection--single {
    width: 200px;
  }

  .warn {
    color: #e02020;
    margin-right: 5px;
  }

  .nameSelect {
    margin-bottom: 20px;
  }
`
