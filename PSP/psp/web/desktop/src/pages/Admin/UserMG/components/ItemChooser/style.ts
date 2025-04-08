import styled from 'styled-components'

export const RoleChooserWrapper = styled.div`
  display: flex;
  margin-top: 10px;

  .label {
    font-family: 'PingFangSC-Medium';
    font-size: 14px;
    color: rgba(0, 0, 0, 0.85);
  }

  .module {
    .header {
      display: flex;
      align-items: center;
      height: 32px;
      margin-bottom: 10px;
    }

    .body {
      border: 1px solid #d8d8d8;
      border-radius: 4px;
      height: 400px;
      display: flex;
      flex-direction: column;
    }
  }

  .left {
    width: 300px;
    margin-right: 20px;

    .body {
      padding: 5px;
      overflow: auto;

      .itemList {
        display: flex;
        flex-wrap: wrap;

        .item {
          border: 1px solid #000000;
          border-radius: 12.5px;
          margin: 5px;
          padding: 3px 5px;
        }
      }
    }
  }

  .right {
    flex: 1;

    .filter {
      width: 120px;
      margin-left: auto;
    }

    .body {
      padding: 5px 10px;

      .all {
        padding: 8px 10px;
        border-bottom: 1px solid #b6b6b6;

        .title {
          font-family: 'PingFangSC-Regular';
          font-size: 14px;
          color: rgba(0, 0, 0, 0.85);
          letter-spacing: 0;
        }
      }

      .itemList {
        flex: 1;
        overflow: auto;

        .item {
          padding: 8px 10px;

          &:hover {
            background: #f2f6ff;
          }

          & > span {
            display: inline-block;
            max-width: 360px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            vertical-align: bottom;
          }
        }
      }
    }
  }
`
