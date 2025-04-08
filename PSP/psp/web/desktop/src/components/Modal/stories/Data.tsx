import React from 'react'
import { Modal, Button } from '../..'
import styled from 'styled-components'
import { Input as AntInput, message } from 'antd'
import { useObserver, useLocalStore } from 'mobx-react-lite'

const StyledContainer = styled.div`
  margin: 10px;

  button {
    margin: 5px;
  }
`

const ControlledInput = ({ data, onChange }) =>
  useObserver(() => <AntInput value={data.value} onChange={onChange} />)

function ControlledModal() {
  const state = useLocalStore(() => ({
    value: '',
    setValue(value) {
      this.value = value
    },
  }))

  function onChange(e) {
    state.setValue(e.target.value)
  }

  return (
    <Button
      type='primary'
      onClick={() =>
        Modal.show({
          content: <ControlledInput data={state} onChange={onChange} />,
        }).then(() => {
          message.success(`获取数据：${state.value}`)
        })
      }>
      受控弹窗
    </Button>
  )
}

export const Data = () => (
  <StyledContainer>
    <ControlledModal />

    <Button
      type='primary'
      onClick={() =>
        Modal.show({
          content: ({ onOk }) => (
            <Button onClick={() => onOk('valueFromModal')}>获取数据</Button>
          ),
          footer: null,
        }).then(data => {
          message.success(`获取数据：${data}`)
        })
      }>
      回调返回值
    </Button>
  </StyledContainer>
)
