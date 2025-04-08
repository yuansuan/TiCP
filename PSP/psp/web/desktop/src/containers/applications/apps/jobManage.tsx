/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useState } from 'react'
import { useSelector } from 'react-redux'
import { ToolBar } from '../../../utils/general'
import NewJobManager from '@/pages/NewJobManager'
import styled from 'styled-components'

const Wrapper  = styled.div` 
  height: calc(100vh - 100px);
` 

export const JobManage = () => {
  const wnapp = useSelector((state) => state.apps.jobLists)
  const [ID, setID] = useState(1)

  if (wnapp.winRefresh) {
    setID(ID+1)
    wnapp.winRefresh = false
  } 

  return wnapp ? (
    <div
      className="notepad floatTab dpShad"
      data-size={wnapp.size}
      data-max={wnapp.max}
      style={{
        ...(wnapp.size == 'cstm' ? wnapp.dim : null),
        zIndex: wnapp.z,
      }}
      data-hide={wnapp.hide}
      id={wnapp.icon + 'App'}
    >
      <ToolBar
        app={wnapp.action}
        icon={wnapp.icon}
        size={wnapp.size}
        name={wnapp.title}
        hasRefresh={true}
      />
      <div className="windowScreen flex flex-col" data-dock="true">
        { !wnapp.hide && 
          <Wrapper>
            <NewJobManager key={ID}/>
          </Wrapper>
        }
      </div>
    </div>
  ) : null
}
