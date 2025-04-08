/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo } from 'react'
import { Checkbox } from 'antd'
import { StyledColumnManager } from './style'

interface IProps {
  Icon: any
  columns: any[]
  config: Array<{ key: string; active: boolean }>
  setConfig: (config: any) => void
}

export function ColumnManager(props: IProps) {
  const columns = useMemo(() => {
    const { columns: originColumns, config } = props
    const columns = config.map(item => {
      const column = originColumns.find(column => column.dataKey === item.key)
      return {
        ...item,
        name: column ? column.name : '',
        fixable: column?.props?.fixed || '',
      }
    })

    return columns
  }, [])
  const activeColumns = useMemo(() => columns.filter(column => column.active), [
    columns,
  ])

  function order(key: string, operator: string) {
    const { setConfig } = props
    const index = columns.findIndex(item => item.key === key)

    switch (operator) {
      case 'up': {
        if (index <= 0) {
          return
        }
        const prev = columns[index - 1]
        if (prev.fixable) return
        columns[index - 1] = columns[index]
        columns[index] = prev
        break
      }

      case 'down': {
        if (index >= columns.length - 1) {
          return
        }
        const next = columns[index + 1]
        if (next.fixable) return
        columns[index + 1] = columns[index]
        columns[index] = next
        break
      }
    }

    setConfig(
      columns.map(item => ({
        key: item.key,
        active: item.active,
      }))
    )
  }

  function activate(key: string, checked: boolean) {
    const { setConfig } = props

    const column = columns.find(item => item.key === key)
    column.active = checked

    setConfig(
      columns.map(item => ({
        key: item.key,
        active: item.active,
      }))
    )
  }

  const { Icon } = props

  return (
    <StyledColumnManager>
      {columns.map((column, index) => (
        <div key={index} className='item'>
          <Checkbox
            disabled={column.active && activeColumns.length <= 1}
            checked={column.active}
            onChange={e => activate(column.key, e.target.checked)}
          />
          <div className='name'>{column.name}</div>
          <div className='move'>
            {!column.fixable && (
              <>
                <Icon
                  type='arrow_up'
                  style={{ fontSize: 16 }}
                  onClick={() => order(column.key, 'up')}
                />
                <Icon
                  type='arrow_down'
                  style={{ fontSize: 16 }}
                  onClick={() => order(column.key, 'down')}
                />
              </>
            )}
          </div>
        </div>
      ))}
    </StyledColumnManager>
  )
}
