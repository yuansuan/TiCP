import * as React from 'react'
import { observer } from 'mobx-react'
import styled from 'styled-components'
import { DropTarget, DropTargetMonitor, DropTargetConnector } from 'react-dnd'
import flow from 'lodash/flow'

import { inject } from '@/pages/context'
import { App } from '@/domain/Applications'
import Section from '@/domain/Applications/App/Section'
import { DndTypes } from '@/components/FormField'

const target = {
  drop: (props, monitor: DropTargetMonitor) => {
    const item = monitor.getItem()
    const type = monitor.getItemType()
    const {
      index,
      app: { subForm }
    } = props

    // new section
    if (type === DndTypes.EMPTY_SECTION) {
      const section = new Section()
      subForm.add(section, index)
    } else {
      // toggle two section
      // ignore adjacent toggle
      if (index === item.index || index === item.index + 1) {
        return
      }
      let targetIndex = -1
      // insert before
      if (index < item.index) {
        targetIndex = index
      } else {
        // insert after
        targetIndex = index - 1
      }
      subForm.toggle(item.index, targetIndex)
    }
  }
}

const targetCollect = (
  connect: DropTargetConnector,
  monitor: DropTargetMonitor
) => ({
  connectDropTarget: connect.dropTarget(),
  isOver: monitor.isOver(),
  canDrop: monitor.canDrop()
})

const Wrapper: any = styled.div`
  width: 100%;
  height: 10px;
  background: ${(props: any) => (props.isActive ? '#368EFF' : 'transparent')};
`

const EmptyWrapper: any = styled.div`
  width: 100%;
  height: 200px;
  background: ${(props: any) => (props.isActive ? '#aaa' : '')};
  position: absolute;
`

interface IProps {
  connectDropTarget?: any
  isOver?: any
  canDrop?: any
  app?: App
  index: number
}

@flow(
  observer,
  DropTarget([DndTypes.SECTION, DndTypes.EMPTY_SECTION], target, targetCollect),
  inject(({ app }) => ({ app }))
)
export default class SectionTarget extends React.Component<IProps> {
  render() {
    const {
      connectDropTarget,
      isOver,
      canDrop,
      app: {
        subForm: { sections }
      }
    } = this.props
    const isActive = isOver && canDrop

    return connectDropTarget(
      sections.length === 0 ? (
        <div>
          <EmptyWrapper isActive={isActive} />
        </div>
      ) : (
        <div>
          <Wrapper isActive={isActive} />
        </div>
      )
    )
  }
}
