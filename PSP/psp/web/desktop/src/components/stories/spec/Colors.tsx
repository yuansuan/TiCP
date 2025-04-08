/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { theme } from '@/utils'
import Table from '../../Table'

const StyledLayout = styled.div``

const StyledColor = styled.div`
  position: relative;
  margin-left: 30px;

  &::before {
    content: '';
    position: absolute;
    left: -25px;
    top: 25%;
    width: 20px;
    height: 20px;
    background-color: ${props => props['data-color']};
  }
`

export const Colors = function Colors() {
  const dataSource = [
    ['primaryColor', '系统主色', '用于主按钮颜色'],
    ['linkColor', '链接颜色', '用于链接'],
    ['errorColor', '异常提示色', '用于错误、警告'],
    ['warningColor', '异常提示色', '表示严重程度较低'],
    ['successColor', '成功提示色', '表示正确'],
    ['infoColor', '正常提示色', '表示正常、可行'],
    ['borderColorBase', '外部边框色', '用于组件外部边框'],
    ['borderColorSplit', '内部边框色', '用于组件内部边框分割线'],
    ['backgroundColorBase', '背景色', '用于背景'],
    ['disabledColor', '禁用色', '用于禁用按钮、链接颜色'],
    ['backgroundColorHover', '背景悬浮色', '用于鼠标悬浮时的背景颜色'],
    ['secondaryColor', '次色', '用于次按钮颜色'],
    ['cancelColor', '取消色', '用于取消按钮颜色'],
    ['cancelHighlightColor', '取消高亮色', '用于取消按钮悬浮高亮的颜色'],
  ].map(([key, name, desc]) => ({
    key,
    name,
    desc,
    value: theme[key],
  }))

  return (
    <StyledLayout>
      <Table
        props={{ data: dataSource, autoHeight: true }}
        columns={[
          {
            header: '名称',
            dataKey: 'name',
            props: {
              width: 150,
            },
          },
          {
            header: '键值',
            dataKey: 'key',
            props: {
              width: 200,
            },
          },
          {
            header: '颜色',
            props: {
              width: 150,
            },
            cell: {
              props: { dataKey: 'value' },
              render({ rowData, dataKey }) {
                const value = rowData[dataKey]
                return <StyledColor data-color={value}>{value}</StyledColor>
              },
            },
          },
          {
            header: '描述',
            dataKey: 'desc',
            props: {
              width: 300,
            },
          },
        ]}
      />
    </StyledLayout>
  )
}
