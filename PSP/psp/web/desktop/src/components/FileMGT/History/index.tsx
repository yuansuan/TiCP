/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { useStore } from '../store'
import { Button } from '@/components'
import { observer } from 'mobx-react-lite'
import { useResize } from '@/domain'
import { Bar } from './Bar'
import { LeftOutlined, RightOutlined } from '@ant-design/icons'

const StyledLayout = styled.div`
  display: flex;
  height: 32px;
  line-height: 32px;
  overflow: hidden;

  > .tip {
    color: rgba(0, 0, 0, 0.45);
    margin-right: 10px;
  }

  > .history-jump {
    button:first-child {
      margin-right: 8px;
    }
  }
`

export const History = observer(function History() {
  const { history, currentNode, setNodeId, dirTree } = useStore()
  const { current, prevable, nextable } = history
  const [rect, ref] = useResize()

  useEffect(() => {
    if (!currentNode) {
      return
    }

    if (current?.path !== currentNode.path) {
      history.push({
        path: currentNode.path,
      })
    }
  }, [currentNode])

  function prev() {
    const record = history.prev()
    if (record?.path) {
      const node = dirTree.filterFirstNode(item => item.path === record.path)
      setNodeId(node?.id)

      // record is depracted
      if (!node) {
        history.delete(history.cursor)
        if (history.prevable) {
          prev()
        } else {
          next()
        }
      }
    }
  }

  function next() {
    const record = history.next()
    if (record?.path) {
      const node = dirTree.filterFirstNode(item => item.path === record.path)
      setNodeId(node?.id)

      // record is depracted
      if (!node) {
        history.delete(history.cursor)
        if (history.nextable) {
          next()
        } else {
          prev()
        }
      }
    }
  }

  return (
    <StyledLayout ref={ref}>
      <div className='tip'>文件位置:</div>
      <Bar width={rect.width - 198} />
      <div className='history-jump'>
        <Button disabled={!prevable} onClick={prev}>
          <LeftOutlined />
        </Button>
        <Button disabled={!nextable} onClick={next}>
          <RightOutlined />
        </Button>
      </div>
    </StyledLayout>
  )
})
