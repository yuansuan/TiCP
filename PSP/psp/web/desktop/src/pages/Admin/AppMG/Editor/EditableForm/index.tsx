import * as React from 'react'
import styled from 'styled-components'
import HTML5Backend from 'react-dnd-html5-backend'
import { DndProvider } from 'react-dnd'
import { App } from '@/domain/Applications'
import Menu from './Menu'
import Panel from './Panel'

const Wrapper = styled.div`
  display: flex;
  height: 100%;
  background-color: white;
`
const backend = HTML5Backend

interface IProps {
  win?: any
  app?: App
  formModel: any
  showMenu?:boolean
}

export default class Form extends React.Component<IProps> {
  render() {
    const { showMenu = true,formModel } = this.props

    return (
      <Wrapper>
        <DndProvider backend={backend}>
          {showMenu && <Menu />}
          <Panel formModel={formModel} />
        </DndProvider>
      </Wrapper>
    )
  }
}
