import { Icon } from 'antd'
import React from 'react'
import { Modal } from '@/components'
import { Scrollbars } from 'react-custom-scrollbars'

export const enlarge = (content, title) => {
  return (
    <Icon
      type='arrows-alt'
      rotate={90}
      style={{ fontSize: 18, paddingLeft: 8 }}
      onClick={() => {
        Modal.showConfirm({
          title,
          footer: null,
          width: 900,
          maskClosable: true,
          bodyStyle: {
            height: 710,
            overflow: 'auto',
          },
          content: (
            <Scrollbars style={{ height: '100%' }}>{content} </Scrollbars>
          ),
        })
      }}
    />
  )
}
