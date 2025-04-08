/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { theme } from '@/utils'
import Table from '../../Table'

const StyledLayout = styled.div``

export const Fonts = function Fonts() {
  const dataSource = [
    'Caption',
    'Body',
    'Head4',
    'Display',
    'Head3',
    'Head2',
    'Head1',
  ].map((desc, index) => ({
    name: `${index + 1}号字体`,
    value: theme.fontSize[index],
    desc,
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
              width: 100,
            },
          },
          {
            header: '值',
            dataKey: 'value',
            props: {
              width: 100,
            },
          },
          {
            header: '应用场景',
            props: {
              width: 200,
            },
            cell: {
              props: { dataKey: 'desc' },
              render({ rowData, dataKey }) {
                return (
                  <div style={{ fontSize: rowData.value }}>
                    {rowData[dataKey]}
                  </div>
                )
              },
            },
          },
        ]}
      />
    </StyledLayout>
  )
}
