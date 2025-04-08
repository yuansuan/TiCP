import React from 'react'
import { Page } from '@/components'
export default () => {
  return (
    <Page header={null}>
      <div
        style={{
          display: 'flex',
          height: '75vh',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center'
        }}>
        <div>
          <img src={require('@/assets/images/deving.svg')} />
        </div>
        <div>该页面正在开发中...</div>
      </div>
    </Page>
  )
}
