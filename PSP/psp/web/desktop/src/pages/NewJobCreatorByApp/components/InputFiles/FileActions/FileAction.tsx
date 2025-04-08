/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Icon } from '@/components'
import { Spin, Tooltip } from 'antd'
import { WarningOutlined } from '@ant-design/icons'
import { FileActionStyle } from './style'

interface Props {
  children: string | React.ReactNode
  icon?: string | (() => React.ReactNode)
  loading?: boolean
  onClick: () => void
  warningTip?: string
}

export const FileAction = ({
  children,
  icon,
  onClick,
  loading,
  warningTip
}: Props) => {
  const isString = typeof icon === 'string'

  return (<FileActionStyle onClick={loading ? undefined : onClick}>
    {icon && isString && (
      <Spin spinning={loading}>
        <Icon type={icon} className={icon} />
        <Icon type={icon + '_active'} className={icon} />
        {
          warningTip && <Tooltip title={warningTip}>
            <WarningOutlined style={{verticalAlign: 'baseline', color: '#ec942c'}} rev={''} />
          </Tooltip>
        }
      </Spin>
    )}
     {icon && !isString && (
      <Spin spinning={loading}>
        { icon() }
        {
          warningTip && <Tooltip title={warningTip}>
            <WarningOutlined style={{color: '#ec942c'}} rev={''} />
          </Tooltip>
        }
      </Spin>
    )}
    {children}   
  </FileActionStyle>)
}
