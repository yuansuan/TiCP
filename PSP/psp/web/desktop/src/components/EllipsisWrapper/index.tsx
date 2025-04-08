import styled from 'styled-components'

interface IProps {
  title?: string
  width?: number
}

const EllipsisWrapper = styled.div.attrs(props => ({
  title: props.title || props.children,
}))<IProps>`
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: ${props => (props.width ? props.width + 'px' : 'auto')};
`

export default EllipsisWrapper
