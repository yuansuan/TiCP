import React, { useEffect, useState } from 'react'
import { Pagination, DatePicker, Tabs, Select } from 'antd'
import styled from 'styled-components'
import {
  overviewData,
  overviewPagingData,
  detailData,
  detailPagingData,
  DATE_FORMAT
} from './data'
import { formatDate } from '@/utils'
import { Table, Button } from '@/components'
import { useResize } from '@/domain'
import PieChart from '@/components/Chart/PieChart'
import { Http } from '@/utils'
import { eventEmitter } from '@/utils'
import { useStore } from '../../VisIBV/store'
import { startExport } from './export'
import { exportExcel } from './exportData'
import { observer } from 'mobx-react'
import { runInAction } from 'mobx'
import moment from 'moment'
const { RangePicker } = DatePicker

const { TabPane } = Tabs
const { Option } = Select

const Wrapper = styled.div`
  padding: 20px;
  .action {
    display: flex;
    justify-content: space-between;
    margin-bottom: 5px;
    align-items: center;

    .filter {
      display: flex;
      justify-content: flex-start;
      align-items: center;

      .item {
        margin: 5px;
        display: flex;
        align-items: center;

        .label {
          flex: 1 0 66px;
          text-align: right;
          padding: 5px;
        }
      }
    }

    .btn {
      margin 0 5px;
    }
  }

  .firstTabWrapper {
    display: flex;
    .table {
      flex: 1;
    }
    .pie {
      margin-left: 10px;
      flex: 1.3;
      background: #eee;
    }
  }

  .pagination {
    display: flex;
    justify-content: center;
    width: 100%;
    margin-top: 30px;
  }
`

const execlInfo = [
  {
    execlName: '统计总览',
    columnNames: ['ID', '应用名称', '使用时长(小时)'],
    columnKeys: ['app_id', 'app_name', 'duration'],
    formatter: [val => val, val => val, val => val],
    url: '/visual/worktask/history/statistics'
  },
  {
    execlName: '统计明细',
    columnNames: [
      'ID',
      '应用名称',
      '工作站',
      '开始时间',
      '结束时间',
      '使用时长(小时)'
    ],
    columnKeys: [
      'id',
      'app_name',
      'workstation_name',
      'start_time',
      'end_time',
      'duration'
    ],
    formatter: [
      val => val,
      val => val,
      val => val,
      val => val,
      val => val,
      val => val
    ],
    url: '/visual/worktask/history/list'
  }
]

export default observer(function SoftwareReport() {
  const store = useStore()
  const [rect, ref] = useResize()
  const [resize, setResize] = useState(false)
  const [loading, setLoading] = useState(false)
  const [exporting, setExporting] = useState(false)
  const [tabKey, setTabKey] = useState('0')

  const [queryString, setQueryString] = useState({
    appIds: [],
    time: [moment().subtract(7, 'days'), moment()]
  })

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

  useEffect(() => {
    eventEmitter.on('LEFT_TREE_COLLAPSE', ({ message }) => {
      setResize(message.collapsed)
    })

    eventEmitter.on('export_execl_success', ({ message }) => {
      setExporting(false)
    })

    eventEmitter.on('export_execl_error', ({ message }) => {
      setExporting(false)
    })

    return () => {
      eventEmitter.off('LEFT_TREE_COLLAPSE')
      eventEmitter.off('export_execl_success')
      eventEmitter.off('export_execl_error')
    }
  }, [])

  const columns0 = [
    {
      props: {
        flexGrow: 1
      },
      header: '应用ID',
      dataKey: 'app_id'
    },
    {
      props: {
        flexGrow: 2
      },
      header: '应用名称',
      cell: {
        props: {
          dataKey: 'app_name'
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        }
      }
    },
    {
      props: {
        flexGrow: 2
      },
      header: '使用时长（小时）',
      cell: {
        props: {
          dataKey: 'duration'
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        }
      }
    }
  ]

  const columns1 = [
    {
      props: {
        flexGrow: 0.5
      },
      header: 'ID',
      dataKey: 'id'
    },
    {
      props: {
        flexGrow: 1
      },
      header: '应用名称',
      cell: {
        props: {
          dataKey: 'app_name'
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        }
      }
    },
    {
      props: {
        flexGrow: 2
      },
      header: '开始时间',
      cell: {
        props: {
          dataKey: 'start_time'
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        }
      }
    },
    {
      props: {
        flexGrow: 2
      },
      header: '结束时间',
      cell: {
        props: {
          dataKey: 'end_time'
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        }
      }
    },
    {
      props: {
        flexGrow: 2
      },
      header: '使用时长（小时）',
      cell: {
        props: {
          dataKey: 'duration'
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        }
      }
    }
  ]

  const exportExcelFile = () => {
    setExporting(true)

    const { url, execlName, columnKeys, columnNames, formatter } =
      execlInfo[tabKey]

    const { appIds, time } = queryString

    startExport(url, time, appIds, allData => {
      const sheetNameMap = {
        name: execlName
      }

      let sheetData = allData.flat().map(d => {
        let row = []
        columnKeys.forEach((key, index) => {
          row.push(formatter[index](d[key]))
        })

        return row
      })

      sheetData.unshift(columnNames)

      const sheets = []

      sheets.push({
        sheetName: sheetNameMap.name,
        data: sheetData
      })

      exportExcel({
        excelName: `${execlName}报表`,
        sheets
      })
    })
  }

  const query = tabKey => {
    if (tabKey === '0') {
      runInAction(() => {
        overviewPagingData.index = 1
        getStatistics()
      })
    } else {
      runInAction(() => {
        detailPagingData.index = 1
        getList()
      })
    }
  }

  return (
    <Wrapper ref={ref} style={{ height: 'calc(100vh - 180px)', width: '100%' }}>
      <div className='action'>
        <div className='filter'>
          <div className='item'>
            <span className='label'>应用名称: </span>
            <Select
              style={{ width: 300 }}
              value={queryString.appIds}
              maxTagCount={4}
              mode={'multiple'}
              onSelect={value => {
                setQueryString({
                  ...queryString,
                  appIds: [...queryString.appIds, value]
                })
              }}
              onDeselect={value => {
                let { appIds } = queryString
                let index = appIds.indexOf(value)
                appIds.splice(index, 1)

                setQueryString({
                  ...queryString,
                  appIds: [...appIds]
                })
              }}>
              {store.software.softwareList.map(app => (
                <Option key={app.id} value={app.id}>
                  {app.name}
                </Option>
              ))}
            </Select>
          </div>
          <div className='item'>
            <span className='label'>时间: </span>
            <RangePicker
              value={queryString.time}
              showTime={{ format: 'HH:mm' }}
              format='YYYY-MM-DD HH:mm'
              onChange={dates => {
                setQueryString({
                  ...queryString,
                  time: dates
                })
              }}
              allowClear={true}
            />
          </div>
        </div>
        <div>
          <Button className='btn' onClick={() => query(tabKey)}>
            查询
          </Button>
        </div>
      </div>
      <Tabs
        activeKey={tabKey}
        onChange={key => {
          setTabKey(key)
          query(key)
        }}
        // tabBarExtraContent={
        //   <Button
        //     disabled={exporting}
        //     type='link'
        //     onClick={() => exportExcelFile()}>
        //     数据导出
        //   </Button>
        // }
      >
        <TabPane tab='统计总览' key={'0'}>
          <div className='firstTabWrapper'>
            <div className='table'>
              <Table
                props={{
                  height: rect.height - 180,
                  loading: loading,
                  data: overviewData.list,
                  rowKey: 'app_id'
                }}
                columns={columns0 as any}
              />
            </div>
            <div
              className='pie'
              style={{ width: resize ? rect.width / 2 : rect.width / 2 - 200 }}>
              <PieChart
                data={overviewData.list.map(app => ({
                  key: app.app_name,
                  value: +app.duration
                }))}
                unit={'小时'}
                title={''}
              />
            </div>
          </div>
          {/* <div className='pagination'>
            <Pagination
              pageSize={overviewPagingData.size}
              current={overviewPagingData.index}
              total={overviewData.total}
              onChange={(index, size) => {
                runInAction(() => {
                  overviewPagingData.index = index
                  overviewPagingData.size = size
                  getStatistics()
                })
              }}
            />
          </div> */}
        </TabPane>
        <TabPane tab='统计明细' key={'1'}>
          <Table
            props={{
              height: rect.height - 180,
              loading: loading,
              data: detailData.list,
              rowKey: 'id'
            }}
            columns={columns1 as any}
          />
          <div className='pagination'>
            <Pagination
              pageSize={detailPagingData.size}
              current={detailPagingData.index}
              total={detailData.total}
              onChange={(index, size) => {
                runInAction(() => {
                  detailPagingData.index = index
                  detailPagingData.size = size
                  getList()
                })
              }}
            />
          </div>
        </TabPane>
      </Tabs>
    </Wrapper>
  )
})
