/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { ReactElement, useRef } from 'react'
import { DndProvider, createDndContext } from 'react-dnd'
import Backend from 'react-dnd-html5-backend'

const RNDContext = createDndContext(Backend)

function useDNDProviderElement(props) {
  const manager = useRef(RNDContext)

  if (!props.children) return null

  return (
    <DndProvider manager={manager.current.dragDropManager}>
      {props.children}
    </DndProvider>
  )
}

export default function DragAndDrop(props) {
  const DNDElement = useDNDProviderElement(props)
  return <React.Fragment>{DNDElement}</React.Fragment>
}

interface IProps {
  children: ReactElement
}

export function DraggableWarpper(props: IProps) {
  return <DragAndDrop>{props.children}</DragAndDrop>
}
