/* Copyright (C) 2016-present, Yuansuan.cn */
import styled from 'styled-components'

export const ListWrapper = styled.div`
  background: #fff;
  width: 100%;
  height: calc(100vh - 100px);

  .footer {
    padding: 10px 17px 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};

    .main {
      display: flex;
      justify-content: space-between;
      width: 100%;
    }
  }
`

export const ProtalWrapper = styled.div`
  background: #fff;
  width: 100%;

  .footer {
    padding: 10px 17px 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`
export const CloseTimeWrapper = styled.div`
  display: flex;
  width: 80%;
  align-items: center;
  padding-left: 20px;
`

export const ZoneWrapper = styled.div`
  display: flex;
  line-height: 32px;
  padding: 20px;
  margin-bottom: 3px;
  background: #fff;

  display: none;
`

export const ListActionWrapper = styled.div`
  display: flex;
  flex-wrap: nowrap;
  padding: 20px;
  justify-content: flex-start;

  .item {
    padding: 5px;
    display: flex;
    justify-content: flex-start;
    align-items: flex-start;
    
    & > div {
      padding: 0 0 10px 0;
    }

    .label {
      flex: 0 0 auto;
      width: 80px;
      margin-left: 20px;
    }

    .field {
      width: 200px;
    }

    .btn {
      margin: 0 5px;
    }
  }
  .mountDevice {
    line-height: 30px;
    display:flex;
    > img {
      cursor: pointer;
      width: 30px;
      height: 30px;
      margin-top:6px;
    }
    .showNames {
      margin-top: 10px;
    }
  }
  
  
`

export const ListDataWrapper = styled.div`
  width: 100%;
  padding: 0 20px;
  .main {
    padding: 0px 20px;
    display: flex;
    flex-direction: column;

    .pagination {
      margin: 20px auto;
    }
  }
`

export const ModalListDataWrapper = styled.div`
  width: 100%;
  max-height: 220px;
  overflow-y: scroll;
`

export const StatusWrapper = styled.div`
  display: flex;
  align-items: center;

  .icon {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    margin-right: 4px;
    position: relative;
    box-sizing: content-box;
  }

  .text {
    margin-left: 4px;
  }

  .icon-right {
    margin-left: 4px;

    .anticon {
      height: 12px;
      width: 12px;
      position: absolute;
      top: 21px;
    }
  }
`
export const UpdateWapper = styled.div`
  display: flex;
  align-items: center;
  padding: 20px;
`
