/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { observable, action } from 'mobx'
import { showNewComerModal } from './NewComerModal'
import { Hover } from '@/components'

const StyledDiv = styled.div`
  position: fixed;
  top: 160px;
  right: -70px;
  width: 90px;
  height: 60px;
  transition: transform 0.1s ease-in;
  transition-delay: 0.8s;

  .tip-pane-expanded {
    margin-top: 20px;
    width: 90px;
    height: 40px;
    background: #ffffff;
    box-shadow: 0 4px 12px 0 rgba(0, 0, 0, 0.2);
    border-radius: 100px 0 0 100px;
    padding-left: 24px;
    z-index: 1;
    line-height: 40px;
    font-size: 14px;
    color: #3182ff;
    cursor: pointer;
    visibility: hidden;
    transition-delay: 0.9s;
  }

  .tip-pane-light {
    position: absolute;
    left: -8px;
    top: -1px;
  }

  &.hovered {
    transform: translateX(-70px);
    transition-delay: 0s;

    .tip-pane-expanded {
      visibility: visible;
      transition-delay: 0s;
    }
  }
`

class LightStore {
  @observable light_expanded: boolean = false
  @action
  setLightExpanded(bool) {
    this.light_expanded = bool
  }
}

export const lightStore = new LightStore()

export const LightComp = observer(function LightComp() {
  const { light_expanded } = lightStore

  return (
    <Hover
      render={hovered => (
        <StyledDiv
          className={`light-tip ${hovered || light_expanded ? 'hovered' : ''}`}
          onClick={() => {
            lightStore.setLightExpanded(true)
            showNewComerModal({ nextBtnText: '确定' })
          }}>
          <div className='tip-pane-light'>
            <img src={require('@/assets/images/light_default.svg')} alt='tip' />
          </div>
          <div className='tip-pane-expanded'>教学视频</div>
        </StyledDiv>
      )}
    />
  )
})
