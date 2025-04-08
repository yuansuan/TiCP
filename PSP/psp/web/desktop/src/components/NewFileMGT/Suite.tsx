/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useEffect } from 'react'
import styled from 'styled-components'
import {
  Menu,
  Toolbar,
  FileList,
  Context,
  useModel,
  History,
  useStore
} from '.'
import { observer } from 'mobx-react-lite'
import { Resizable } from 're-resizable'
import { useResize, env } from '@/domain'
import { Icon } from '@/components'
import { useLayoutRect } from '@/utils/hooks'

const StyledLayout = styled.div`
  height: 100%;
  > .areaSelectWrap {
    display: none;
    padding: 10px 20px;
    border-bottom: 6px solid #f5f5f5;
    > div {
      display: flex;
      align-items: center;
      /* h3 {
        margin-bottom: 0;
        font-weight: normal;
      } */
    }
  }

  > .file_header {
    padding: 20px;
  }

  > .file_body {
    display: flex;
    width: 100%;
    border-top: 4px solid ${({ theme }) => theme.backgroundColorBase};

    > .file_menu {
      .resizeBar {
        z-index: 5;
      }
      height: 100%;
    }

    > .mockbar {
      position: relative;
      width: 1px;
      background-color: ${({ theme }) => theme.borderColorBase};
      z-index: 1;

      > .wrapper {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 10px;
        height: 26px;
        color: #c9c9c9;
        background-color: #eee;
        border-radius: 5px;
        > .icon {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
        }
      }
    }

    > .files {
      flex: 1;
      display: flex;
      flex-direction: column;

      > .file_toolbar {
        padding: 20px;
      }

      > .file_list {
        flex: 1;
      }
    }
  }
`

type Props = {
  model?: ReturnType<typeof useModel>
  jobManger?: boolean // 是否是作业管理模块渲染
}

export const BaseSuite = observer(function Suite(props: any) {
  const { dirTree, getWidget, initDirTree, refresh } = useStore()
  const [width, setWidth] = useState(250)
  const [headerRect, headerRef, headerResize] = useLayoutRect()
  const [rect, ref, resize] = useResize()

  // hack: harmony file_header size
  useEffect(() => {
    headerResize()
    setTimeout(resize, 100)
  }, [])

  useEffect(() => {
    initDirTree()
    // 作业管理无需请求
    !props.jobManger && refresh()
  }, [])

  return (
    <StyledLayout>
      {getWidget('history') || (
        <div className='file_header' ref={headerRef}>
          <History />
        </div>
      )}
      <div
        className='file_body'
        ref={ref}
        style={{ height: `calc(100% - ${headerRect.height}px)` }}>
        <div className='file_menu'>
          <Resizable
            handleClasses={{ right: 'resizeBar' }}
            minWidth={120}
            enable={{ right: true }}
            size={{ width, height: '100%' }}
            onResizeStop={(e, direction, ref, d) => {
              setWidth(width + d.width)
              resize()
            }}>
            <Menu />
          </Resizable>
        </div>
        <div className='mockbar'>
          <div className='wrapper'>
            <Icon className='icon' type='drag' />
          </div>
        </div>
        <div className='files'>
          <div className='file_list'>
            {dirTree.children.length > 0 && (
              <FileList
                // width={rect.width - width - 43}
                height={rect.height - 127}
              />
            )}
          </div>
        </div>
      </div>
    </StyledLayout>
  )
})

export function Suite({ model, ...props }: Props) {
  const defaultModel = useModel()
  const finalModel = model || defaultModel

  return (
    <Context.Provider value={finalModel}>
      <BaseSuite {...props} />
    </Context.Provider>
  )
}
