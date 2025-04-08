import React, { useCallback } from 'react'
import styled from 'styled-components'
import { Button, Modal } from '@/components'
import Search from '@/components/Search'
import { observer } from 'mobx-react-lite'
import { useStore } from './store'
import SoftwareEditor from './SoftwareEditor'
import HardwareEditor from './HardwareEditor'

const Wrapper = styled.div`
  width: 100%;
  padding-bottom: 14px;

  .toolbar {
    display: flex;
    justify-content: space-between;

    .button-group {
      button {
        margin: 0 10px;
      }
    }

    .search-input {
      display: border-box;
      width: 150px;
      margin-left: auto;
      margin-right: 10px;
    }
  }
`

export const Toolbar = observer(function Toolbar() {
  const store = useStore()

  function add() {
    Modal.show({
      title: `添加${store.tabType === '1' ? '镜像' : '实例'}`,
      width: 600,
      bodyStyle: { padding: 0, height: 640 },
      footer: null,
      content: ({ onCancel, onOk }) => {
        return (
          <>
            {store.tabType === '1' && (
              <SoftwareEditor
                onCancel={onCancel}
                onOk={() => {
                  onOk()
                  store.refreshSoftware()
                }}
              />
            )}
            {store.tabType === '2' && (
              <HardwareEditor
                onCancel={onCancel}
                onOk={() => {
                  onOk()
                  store.refreshHardware()
                }}
              />
            )}
          </>
        )
      }
    })
  }

  return (
    <Wrapper>
      <div className='toolbar'>
        <div className='button-group'>
          <Button icon='add' type='primary' onClick={add}>
            添加
          </Button>
        </div>

        <Search
          className='search-input'
          placeholder='请输入名称来搜索'
          onSearch={name => store.setName(name)}
        />
      </div>
    </Wrapper>
  )
})
