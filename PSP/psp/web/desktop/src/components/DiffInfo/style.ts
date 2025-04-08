import styled from 'styled-components'

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
`

export const StyledSection = styled.div`
  border-bottom: 1px solid ${props => props.theme.borderColor};
  margin-bottom: 16px;

  .tag {
    display: inline-block;
    padding: 2px 8px;
    background: #1b65f5;
    border-radius: 2px;
    color: white;
    margin-bottom: 15px;
  }

  .label {
    font-family: 'PingFangSC-Medium';
    font-size: 14px;
    color: rgba(0, 0, 0, 0.65);
    margin-right: 17px;
    width: 100px;
    text-align: right;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
`
