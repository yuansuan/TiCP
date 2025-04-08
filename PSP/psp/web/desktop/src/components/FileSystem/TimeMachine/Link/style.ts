import styled from 'styled-components'

export const StyledLink = styled.div`
  display: inline-block;
  position: relative;
  top: 2px;
  cursor: pointer;
  overflow: hidden;
  max-width: 100px;
  text-overflow: ellipsis;
  white-space: nowrap;

  &.active,
  &:hover {
    color: ${props => props.theme.primaryHighlightColor};
  }
`
