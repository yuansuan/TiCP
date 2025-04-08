import React, { useEffect, useRef, useState } from 'react'
import { observer } from 'mobx-react-lite'
import styled from 'styled-components'
import { useStore } from './store'
import { Pagination } from 'antd'
import { Toolbar } from './Toolbar'
import { SoftwareList } from './SoftwareList'
import { HardwareList } from './HardwareList'
import SessionListPage from './3DSessions'
import SoftwareReport from '../VisualMgr/SoftwareReport'

const Wrapper = styled.div`
  height: calc(100vh - 220px);

  .title {
    font-size: 20px;
    margin-bottom: 10px;
  }

  .main {
    padding: 14px 0;

    > .footer {
      display: flex;
      flex: 1;
      justify-content: center;
      align-items: center;
      margin: 10px;
    }
  }
`

export const Content = observer(function Content() {
  const store = useStore()
  const ref= useRef(null)
  const [height, setHeight] = useState(800)
  const { total } =
    store.tabType === '1' ? store.software.page_ctx : store.hardware.page_ctx
  const hardwarePage = store.hardware.page_ctx


  useEffect(() => {
    const resizeObserver = new ResizeObserver((entries) => {
      for (let entry of entries) {
          setHeight(entry.contentRect.height)
      }
    })
    
    resizeObserver.observe(ref.current)

    return () => resizeObserver.disconnect()
  }, [])

  function onPageChange(index, size) {
    store.setPageIndex(index)
    store.setPageSize(size)
  }

  function onSessionPageSizeChange(index, size) {
    store.setSessionPageIndex(index)
    store.setSessionPageSize(size)
  }

  function onSessionPageChange(index, size) {
    store.setSessionPageIndex(index)
    store.setSessionPageSize(size)
  }

  return (
    <Wrapper ref={ref}>
      {(store.tabType == '1' || store.tabType == '2') && <Toolbar />}
      <div className='main'>
        {store.tabType === '1' && (
          <>
            <SoftwareList height={height} />
            {total > 0 && (
              <div className='footer'>
                <Pagination
                  disabled={store.loading}
                  showSizeChanger
                  pageSize={store.pageSize}
                  current={store.pageIndex}
                  total={total}
                  onChange={onPageChange}
                />
              </div>
            )}
          </>
        )}
        {store.tabType === '2' && (
          <>
            <HardwareList height={height} />
            {hardwarePage.total > 0 && (
              <div className='footer'>
                <Pagination
                  disabled={store.loading}
                  showSizeChanger
                  pageSize={store.pageSize}
                  current={store.pageIndex}
                  total={total}
                  onChange={onPageChange}
                />
              </div>
            )}
          </>
        )}
        {store.tabType === '3' && (
          <>
            <SessionListPage height={height} />
            {store.model.total > 0 && (
              <div className='footer'>
                <Pagination
                  showSizeChanger
                  onChange={onSessionPageChange}
                  onShowSizeChange={onSessionPageSizeChange}
                  pageSize={store.page_size}
                  current={store.page_index}
                  total={store.model.total}
                />
              </div>
            )}
          </>
        )}

        {store.tabType === '4' && <SoftwareReport />}
      </div>
    </Wrapper>
  )
})
