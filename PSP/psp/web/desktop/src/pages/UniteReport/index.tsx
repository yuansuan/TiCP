// Copyright (C) 2019 LambdaCal Inc.

import { observer } from 'mobx-react'
import * as React from 'react'
import { Wrapper } from './style'
import Report from './Report'
import { Scrollbars } from 'react-custom-scrollbars'

@observer
export default class ReportPage extends React.Component<any> {
  render() {
    return (
      <Wrapper>
        <Scrollbars
          autoHide
          autoHideTimeout={1000}
          autoHideDuration={200}
          style={{ width: '100%', height: 'calc(100vh - 150px)' }}>
          <div className='body'>
            <Report />
          </div>
        </Scrollbars>
      </Wrapper>
    )
  }
}
