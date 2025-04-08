/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { ReactElement } from 'react'
import { DndProvider } from 'react-dnd'
import Backend from 'react-dnd-html5-backend'

interface IProps {
  children: ReactElement
}

export function DraggableWarpper(props: IProps) {
  return <DndProvider backend={Backend}>{props.children}</DndProvider>
}
