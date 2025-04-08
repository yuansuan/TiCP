import * as React from 'react'
import styled from 'styled-components'

import WorkStation from '@/pages/Admin/VisualMgr/WorkStation'
import WorkTask from '@/pages/Admin/VisualMgr/WorkTask'
import VMMgr from '@/pages/Admin/VisualMgr/VMMgr'
import ImageEditing from '@/pages/Admin/VisualMgr/ImageEditing'

const Wrapper = styled.div`
  .title {
    padding-left: 20px;
  }
`

export const ComponentsMap = {
  WorkStation,
  WorkTask,
  VMMgr,
  ImageEditing
}

export const LabelMap = {
  WorkStation: '应用资源',
  WorkTask: '已打开应用',
  VMMgr: '虚拟机管理',
  ImageEditing: '镜像编辑'
}

export default function ThreeDMgr({ type }) {
  let Component = ComponentsMap[type || 'WorkStation']
  return (
    <Wrapper>
      <div className='title'>{LabelMap[type]}</div>
      <Component />
    </Wrapper>
  )
}
