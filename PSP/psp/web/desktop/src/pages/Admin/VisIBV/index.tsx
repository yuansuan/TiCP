import React, { useEffect, useState } from 'react'
import { runInAction } from 'mobx'
import { Page } from '@/components/Page'
import { observer } from 'mobx-react-lite'
import { useStore, Provider } from './store'
import { Content } from './Content'
import { Http, formatDate } from '@/utils'

import {
  overviewData,
  detailData,
  detailPagingData,
  DATE_FORMAT
} from '../VisualMgr/SoftwareReport/data'
import moment from 'moment'

const VisIBV = observer(function VisIBV() {
  const store = useStore()

  useEffect(() => {
    if (store.tabType === '2') {
      store.setPageSize(10)
      store.refreshHardware()
    } else if (store.tabType === '1') {
      store.setPageSize(10)
      store.refreshSoftware()
    } else if (store.tabType === '3') {
      store.refreshSoftware(0)
      store.refreshHardware(0)
      store.fetchProjects()
      store.fetchSessionList()
    } else if (store.tabType === '4') {
      const queryString = {
        appIds: [],
        time: [moment().subtract(7, 'days'), moment()]
      }
      const fetchData0 = async () => {
        const { appIds, time } = queryString

        return Http.post('/vis/history/statistics', {
          start_time: time ? time[0]?.format(DATE_FORMAT) : '',
          end_time: time ? time[1]?.format(DATE_FORMAT) : '',
          app_ids: appIds.length === 0 ? null : appIds
        })
      }

      const fetchData1 = async () => {
        const { appIds, time } = queryString
        const { index, size } = detailPagingData
        return Http.post('/vis/historyList', {
          start_time: time ? time[0]?.format(DATE_FORMAT) : '',
          end_time: time ? time[1]?.format(DATE_FORMAT) : '',
          app_ids: appIds.length === 0 ? null : appIds,
          page: {
            size: size ?? 10,
            index: index ?? 1
          }
        })
      }

      const getList = async () => {
        const res = await fetchData1()
        let list = res?.data?.historiesDuration.map(app => ({
          ...app,
          start_time: formatDate(app.start_time?.seconds) || '--',
          end_time: formatDate(app.end_time?.seconds) || '--'
        }))

        runInAction(() => {
          detailData.list = list || []
          detailData.total = res.data.total
        })
      }

      const getStatistics = async () => {
        const res = await fetchData0()
        let list = res?.data?.statistics.map(app => ({
          ...app
        }))

        runInAction(() => {
          overviewData.list = list || []
          overviewData.total = res.data.total
        })
      }

      getList()
      getStatistics()
    }
  }, [store.pageIndex, store.pageSize, store.name, store.tabType, store.page_index, store.page_size])
  return (
    <Page
      header={null}
      tabConfig={{
        tabContentList: [
          {
            tabName: '镜像管理',
            tabKey: '1',
            content: <Content />
          },
          {
            tabName: '实例规格',
            tabKey: '2',
            content: <Content />
          },
          {
            tabName: '会话管理',
            tabKey: '3',
            content: <Content />
          }
          // {
          //   tabName: '应用统计',
          //   tabKey: '4',
          //   content: <Content />
          // }
        ].filter(Boolean),
        defaultActiveKey: store.tabType,
        onChange: activeKey => {
          store.setTabType(activeKey)
          store.setName('')
          store.setPageIndex(1)
          store.setSessionPageIndex(1)
        }
      }}
    />
  )
})

export default function VisIBVWithStore() {
  return (
    <Provider>
      <VisIBV />
    </Provider>
  )
}
