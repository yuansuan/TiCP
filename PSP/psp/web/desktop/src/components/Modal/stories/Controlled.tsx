/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Modal, Button } from '../..'
import styled from 'styled-components'
import { observer, useObserver, useLocalStore } from 'mobx-react-lite'

const StyledContainer = styled.div`
  margin: 10px;

  button {
    margin: 5px;
  }
`

const ControlledFooter = observer(function ControlledFooter() {
  const state = useLocalStore(() => ({
    loading: false,
    setLoading(loading) {
      this.loading = loading
    },
  }))

  return (
    <Button
      type='primary'
      onClick={() =>
        Modal.showConfirm({
          onOk: () => {
            state.setLoading(true)
            return new Promise(resolve => {
              setTimeout(() => {
                resolve()
                state.setLoading(false)
              }, 1000)
            })
          },
          OkButton: function OkButton({ onOk }) {
            return useObserver(() => (
              <Button onClick={onOk} loading={state.loading}>
                确定
              </Button>
            ))
          },
        })
      }>
      control footer
    </Button>
  )
})

const ControlledContent = observer(function ControlledContent() {
  const state = useLocalStore(() => ({
    loading: false,
    setLoading(loading) {
      this.loading = loading
    },
  }))

  return (
    <Button
      type='primary'
      onClick={() =>
        Modal.showConfirm({
          content: ({ onOk }) => (
            <Modal.Footer
              CancelButton={null}
              OkButton={observer(function OkButton() {
                return (
                  <Button
                    type='primary'
                    onClick={() => {
                      state.setLoading(true)
                      setTimeout(() => {
                        state.setLoading(false)
                        onOk()
                      }, 1000)
                    }}
                    loading={state.loading}>
                    确定
                  </Button>
                )
              })}
            />
          ),
          footer: null,
        })
      }>
      control content
    </Button>
  )
})

export const Controlled = () => (
  <StyledContainer>
    <ControlledFooter />
    <ControlledContent />
  </StyledContainer>
)
