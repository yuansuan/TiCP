/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useEffect } from 'react'
import styled from 'styled-components'
import { Menu, Toolbar, FileList, Context, useModel, useStore } from '.'
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

      > .toolbar {
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
  job?: any
}

export const BaseSuite = observer(function Suite(props: any) {
  const { job } = props
  const { dirTree, getWidget, initDirTree } = useStore()
  const [width, setWidth] = useState(250)
  const [headerRect, headerRef, headerResize] = useLayoutRect()
  const [rect, ref, resize] = useResize()

  // hack: harmony file_header size
  useEffect(() => {
    headerResize()
    setTimeout(resize, 100)
    
  }, [])

  useEffect(() => {
    if (job?.work_dir && job?.work_dir.endsWith('/')) {
      let rootPath = job?.work_dir.slice(0, -1)
      initDirTree(rootPath)
    }
  }, [job?.work_dir])

  return (
    <StyledLayout>
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
            <Menu isSyncToLocal={job?.isSyncToLocal} />
          </Resizable>
        </div>
        <div className='mockbar'>
          <div className='wrapper'>
            <Icon className='icon' type='drag' />
          </div>
        </div>
        <div className='files'>
          <div className='toolbar'>
            <Toolbar isSyncToLocal={job?.isSyncToLocal} userName={job?.user_name}  />
          </div>
          <div className='file_list'>
            {dirTree.children.length > 0 && (
              <FileList
                // width={rect.width - width - 43}
                height={rect.height - 127}
                isCloud={job?.isCloud}
                isSyncToLocal={job?.isSyncToLocal}
                noNeedLoading={job?.state === 'Running'}
                userName={job?.user_name}
              />
            )}
          </div>
        </div>
      </div>
    </StyledLayout>
  )
})

export function Suite({ ...props }: Props) {
  const defaultModel = useModel(props)

  return (
    <Context.Provider value={defaultModel}>
      <BaseSuite {...props} />
    </Context.Provider>
  )
}
