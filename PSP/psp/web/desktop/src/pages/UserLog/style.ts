import styled from 'styled-components'

export const Wrapper = styled.div`
  width: 100%;
  background-color: rgba(240, 242, 245, 1);

  .body {
    padding: 20px;
    height: 100%;
    overflow: hidden;
    background: #fff;
  }
`

export const ListWrapper = styled.div`
  .actions {
    display: flex;
    justify-content: flex-start;
    margin-bottom: 16px;

    .item {
      padding-right: 10px;
      display: flex;
      align-items: center;
      justify-content: flex-start;
      .label {
        flex: 0 0 80px;
      }
    }
  }

  .list {
    height: calc(100vh - 285px);

    .logTimeline {
      padding: 10px;
    }
  }

  .footer {
    padding: 10px;
    display: flex;
    justify-content: center;
    align-self: center;
  }
`
