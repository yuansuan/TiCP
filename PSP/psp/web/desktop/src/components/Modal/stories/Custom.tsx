/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Modal, Button, Icon } from '../..'
import styled from 'styled-components'

const StyledExample = styled.div`
  margin: 10px;

  button {
    margin: 5px;
  }
`

export const Custom = () => (
  <StyledExample>
    <Button
      type='primary'
      onClick={() =>
        Modal.show({
          CancelButton: null,
          onOk: () =>
            new Promise(resolve => {
              setTimeout(resolve, 1000)
            }),
          OkButton: ({ onOk, loading }) => (
            <Button onClick={onOk} loading={loading}>
              确认
            </Button>
          ),
        })
      }>
      custom OkButton
    </Button>
    <Button
      type='primary'
      onClick={() =>
        Modal.show({
          footer: ({ onCancel }) => (
            <Modal.Footer onCancel={onCancel} OkButton={null} />
          ),
        })
      }>
      自定义 footer
    </Button>
    <Button
      type='primary'
      onClick={() =>
        Modal.show({
          content: ({ onCancel }) => (
            <>
              <Modal.Toolbar
                actions={[
                  {
                    tip: 'cancel',
                    slot: <Icon type='cancel' />,
                  },
                  {
                    tip: 'rename',
                    slot: <Icon type='rename' />,
                  },
                  {
                    tip: 'define',
                    slot: <Icon type='define' />,
                  },
                ]}
              />
              <div style={{ height: 200 }}>content</div>
              <Modal.Footer onCancel={onCancel} OkButton={null} />
            </>
          ),
          footer: null,
        })
      }>
      弹窗内 footer
    </Button>
    <Button
      type='primary'
      onClick={() =>
        Modal.show({
          content: ({ onCancel }) => (
            <>
              <Modal.Toolbar
                actions={[
                  {
                    tip: 'cancel',
                    slot: <Icon type='cancel' />,
                  },
                  {
                    tip: 'rename',
                    slot: <Icon type='rename' />,
                  },
                  {
                    tip: 'define',
                    slot: <Icon type='define' />,
                  },
                ]}
              />
              <div style={{ height: 200 }}>content</div>
              <Modal.Footer onCancel={onCancel} OkButton={null} />
            </>
          ),
          footer: null,
          showHeader: false,
        })
      }>
      隐藏弹窗的header
    </Button>
  </StyledExample>
)
