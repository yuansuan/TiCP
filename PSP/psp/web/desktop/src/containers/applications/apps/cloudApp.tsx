/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useState } from 'react'
import { useSelector } from 'react-redux'
import { ToolBar } from '@/utils/general'
import VisList from '@/pages/VisList/List'
import { env } from '@/domain'
import styled from 'styled-components'

const Wrapper  = styled.div` 
  height: calc(100vh - 100px);
` 

export const CloudApp = () => {
  const wnapp = useSelector(state => state.apps['3dcloudApp'])
  const [ID, setID] = useState(1)

  if (wnapp.winRefresh) {
    setID(ID+1)
    wnapp.winRefresh = false
  } 

  return wnapp ? (
    <div
      className='calcApp floatTab dpShad'
      data-size={wnapp.size}
      id={wnapp.icon + 'App'}
      data-max={wnapp.max}
      style={{
        ...(wnapp.size == 'cstm' ? wnapp.dim : null),
        zIndex: wnapp.z
      }}
      data-hide={wnapp.hide}>
      <ToolBar
        app={wnapp.action}
        icon={env.isPersonal ? 'search' : wnapp.icon}
        size={wnapp.size}
        name={env.isPersonal ? '404' : wnapp.title}
        hasRefresh={true}
      />
      <div className='windowScreen flex flex-col' data-dock='true'>
        {!wnapp.hide && 
          <Wrapper>
            <VisList key={ID} /> 
          </Wrapper>
        }
      </div>
    </div>
  ) : null
}
