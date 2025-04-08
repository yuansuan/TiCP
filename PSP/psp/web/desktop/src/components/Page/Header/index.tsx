/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { PageHeader } from 'antd'

import { StyledLayout } from './style'

export class Header extends React.Component {
  render() {
    return (
      <StyledLayout>
        <PageHeader
          title='Title'
          breadcrumb={{
            routes: [
              {
                path: 'index',
                breadcrumbName: 'First-level Menu',
              },
              {
                path: 'first',
                breadcrumbName: 'Second-level Menu',
              },
              {
                path: 'second',
                breadcrumbName: 'Third-level Menu',
              },
            ],
          }}
          subTitle='This is a subtitle'
        />
      </StyledLayout>
    )
  }
}
