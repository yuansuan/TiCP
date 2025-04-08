/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Button, Modal } from '@/components'
import { useStore, Context } from '../store'
import { LicenseModal } from './LicenseModal'

const StyledLayout = styled.div`
  > * {
    margin: 0 4px;
  }
`

type Props = {
  merchandiseId: string
}

export const Actions = observer(function Actions({ merchandiseId }: Props) {
  const store = useStore()

  async function showModal() {
    await Modal.show({
      title: '许可证设置',
      centered: true,
      footer: null,
      bodyStyle: {
        padding: 0,
        height: 400,
      },
      content: ({ onCancel, onOk }) => (
        <Context.Provider value={store}>
          <LicenseModal
            merchandiseId={merchandiseId}
            onOk={onOk}
            onCancel={onCancel}
          />
        </Context.Provider>
      ),
    })
  }

  return (
    <StyledLayout>
      <Button type='link' onClick={showModal}>
        设置
      </Button>
    </StyledLayout>
  )
})
