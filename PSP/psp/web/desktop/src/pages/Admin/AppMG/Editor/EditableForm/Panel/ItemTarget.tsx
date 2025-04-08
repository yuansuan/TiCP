import * as React from 'react'
import { observer } from 'mobx-react'
import styled from 'styled-components'
import { DropTarget, DropTargetMonitor, DropTargetConnector } from 'react-dnd'
import flow from 'lodash/flow'

import { inject } from '@/pages/context'
import Section from '@/domain/Applications/App/Section'
import Field from '@/domain/Applications/App/Field'
import { DndTypes } from '@/components/FormField'

const target = {
  drop: (props, monitor: DropTargetMonitor) => {
    const item = monitor.getItem()
    const type = monitor.getItemType()
    const { index, parentIndex, sections } = props

    const section = sections[parentIndex]
    // new field
    if (type === DndTypes.EMPTY_FORM_ITEM) {
      const field = new Field({ type: item.type })
      section.add(field, index)
    } else {
      // toggle fields in section
      if (item.parentIndex === parentIndex) {
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
        section.toggle(item.index, targetIndex)
      } else {
        // toggle fields between sections
        const sourceSection = sections[item.parentIndex]
        const targetSection = section
        const sourceField = sourceSection.fields[item.index]

        sourceSection.fields.splice(item.index, 1)
        targetSection.fields.splice(index, 0, sourceField)
      }
    }
  },
}

const targetCollect = (
  connect: DropTargetConnector,
  monitor: DropTargetMonitor
) => ({
  connectDropTarget: connect.dropTarget(),
  isOver: monitor.isOver(),
  canDrop: monitor.canDrop(),
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
`

interface IProps {
  connectDropTarget?: any
  isOver?: any
  canDrop?: any
  sections?: Section[]
  parentIndex: number
  index: number
}

@flow(
  observer,
  DropTarget(
    [DndTypes.FORM_ITEM, DndTypes.EMPTY_FORM_ITEM],
    target,
    targetCollect
  ),
  inject(({ app: { subForm: { sections } } }) => ({ sections }))
)
export default class ItemTarget extends React.Component<IProps> {
  render() {
    const {
      connectDropTarget,
      isOver,
      canDrop,
      sections,
      parentIndex,
    } = this.props
    const isActive = isOver && canDrop
    const section = sections[parentIndex]

    return connectDropTarget(
      section.fields.length === 0 ? (
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
