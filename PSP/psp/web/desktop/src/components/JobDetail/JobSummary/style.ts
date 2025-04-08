import styled from 'styled-components'

export const Wrapper = styled.div`
  padding: 33px 20px 0 20px;
  border-bottom: 1px solid #e8e8e8;
`

export const Title = styled.span`
  font-family: PingFangSC-Medium;
  font-size: 16px;
  color: currentColor;
  line-height: 22px;
`

export const JobNameWrapper = styled.span`
  font-size: 18px;
  color: currentColor;
  padding: 10px;
  line-height: 25px;

  .input {
    width: 300px;
  }

  .icon {
    margin-left: 15px;
  }
`
export const ActionButtonWrapper = styled.div``

export const JobBaseInfoWrapper = styled.div`
  padding: 24px 0;
  display: flex;
  align-items: center;
  .field {
    display: flex;
    align-items: center;
    font-size: 16px;
    min-width: 200px;
    line-height: 22px;
    color: #4c4c4c;
    .label {
      &::after {
        content: ':';
        padding-right: 5px;
      }
    }
    .value {
      overflow: hidden;
      text-overflow: ellipsis;
      word-break: break-all;
      white-space: nowrap;
      width: 150px;
    }
  }
`
