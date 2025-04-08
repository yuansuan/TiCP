import styled from 'styled-components'

export const RoleEditorWrapper = styled.div`
  display: flex;
  flex-direction: column;
  font-size: 16px;
  height: 100%;

  .loading {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 999;
  }

  > .body {
    padding: 40px 50px;
    height: calc(100% - 70px);
    overflow: auto;
  }

  .Softwares {
    .perm {
      display: flex;
      line-height: 53px;
      align-items: center;
      border-bottom: 1px solid #d8d8d8;
    }

    .title {
      display: flex;
      line-height: 53px;
      align-items: center;
      border-bottom: 1px solid black;
      color: rgb(0, 0, 0, 0.85);
    }

    .sf {
      margin-left: 20px;
      width: 500px;
    }
  }
`

export const RoleBasicInfoWrapper = styled.div`
  display: flex;
  flex-direction: column;

  .module {
    display: flex;
    align-items: centeren;
  }

  .module-bottom {
    margin-top: 25px;
    margin-left: 12px;

    & > input {
      height: 80px;
    }
  }

  .warn {
    color: #e02020;
    margin-right: 5px;
  }

  .name {
    margin-left: 21px;
    margin-right: 5px;
    color: rgba(0, 0, 0, 0.85);
  }

  .widget {
    width: 300px;
  }
`

export const SysWrapper = styled.div`
  display: flex;
  margin-top: 30px;

  .name {
    width: 100px;
    margin-left: 32px;
    color: rgba(0, 0, 0, 0.85);
  }

  .permcheck {
    margin-left: 14px;
  }

  .body {
    width: 100%;
    display: grid;
    grid-template-columns: 33.3% 33.3% 33.3%;

    .sysperm {
      margin: 0px 75px 20px 0;
    }
  }
`
