import styled from 'styled-components'

export const StyledOverlay = styled.div`
  display: flex;
  flex-direction: column;
  background-color: white;
  border: 1px solid ${props => props.theme.primaryColor};

  > * {
    display: inline-block;
    width: 100%;
    padding: 5px 10px;
    cursor: pointer;

    &:first-child {
      border-bottom: 1px dashed ${props => props.theme.primaryColor};
    }

    &:last-child {
      border-bottom: 0;
    }

    &:hover {
      background-color: ${props => props.theme.primaryColor};
      color: white;
    }
  }
`
