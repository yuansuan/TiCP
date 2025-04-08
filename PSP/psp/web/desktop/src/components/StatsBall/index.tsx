import styled from 'styled-components'

const StatsBall = styled.div`
  display: inline-block;
  &:before {
    display: inline-block;
    content: '';
    height: 10px;
    width: 10px;
    border-radius: 50%;
    background: ${props => props.color};
    margin-right: 10px;
    margin-bottom: 0px;
  }
`

export default StatsBall
