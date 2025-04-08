import React, { useEffect, useState } from 'react'
import { Pagination, DatePicker, Tabs, Select, message } from 'antd'
import styled from 'styled-components'
import { overviewDataUser as overviewData, 
  overviewPagingDataUser as overviewPagingData,
  detailDataUser as detailData, 
  detailPagingDataUser as detailPagingData,
  DATE_FORMAT,
} from './data'
import { Table, Button } from '@/components'
import { useResize } from '@/domain'
import { Http } from '@/utils'
import { observer } from 'mobx-react'
import { runInAction } from 'mobx'
import { eventEmitter } from '@/utils'
import moment from 'moment'
// import { userWorker as worker, _eventName } from './MainThread'
import qs from 'qs'
import { download } from './exportData'
import { formatISODateStr } from '@/utils/formatter'
import { DatePicker_FORMAT, DatePicker_SHOWTIME_FORMAT, GeneralDatePickerRange } from '@/constant'
import { ProjectSelector } from './common'

const { RangePicker } = DatePicker
const { TabPane } = Tabs
const { Option } = Select

const Wrapper = styled.div`
  padding-top: 0px;
  background: #fff;
  height: calc(100vh - 220px);
  width: 100%;
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

  .pagination {
    display: flex;
    justify-content: center;
    width: 100%;
    margin-top: 15px;
  }
`

export default observer(function ReportByUser(props) {
  const computeTypes = props.computeTypes || []
  const computeTypesMap = props.computeTypesMap || {}
  const [rect, ref] = useResize()
  const [overviewLoading, setOverviewLoading] = useState(false)
  const [detailLoading, setDetailLoading] = useState(false)
  const [appList, setAppList] = useState([])
  const [exporting, setExporting] = useState<boolean | string>(false)
  const [tabKey, setTabKey] = useState('overview')
  const [total, setTotal] = useState(0)

  const [queryString, setQueryString] = useState({
    projectIds: [],
    appIds: [],
    type: 'all',
    time: [moment().subtract(30, 'days'), moment()]
  })

  const getParams = (paging = null) => {
    const { appIds, time, type, projectIds } = queryString
    const params = {
      compute_type: type === 'all' ? null : type,
      start_time: time[0] ? Math.round(time[0]?.valueOf()/1000) : null,
      end_time: time[1] ? Math.round(time[1]?.valueOf()/1000) : null,
      names: appIds.length === 0 ? [] : appIds,
      project_ids: projectIds.length === 0 ? [] : projectIds,
      query_type: 'user',
    }

    return {
      ...params,
      ...(paging ? {
        page_size: paging.size ?? 10,
        page_index: paging.index ?? 1,
      } : {})
    }
  }

  const getTotal = async () => {
    const res = await Http.get('/job/statistics/totalCPUTime', {
      params: getParams()
    })    
    setTotal(res?.data?.cpu_time || 0)
  }

  const getList = async () => {
    setDetailLoading(true)
    try {
      const res = await Http.get('/job/statistics/detail', {
        params: getParams(detailPagingData)
      })    
      runInAction(() => {
        detailData.list = res.data.job_details || []
        detailData.total = res.data.total
      })
    } finally {
      setDetailLoading(false)
    }
  } 

  const getStatistics = async () => {
    setOverviewLoading(true)

    try {
      const res = await Http.get('/job/statistics/overview', {
        params: getParams(overviewPagingData)
      })

      runInAction(() => {
        overviewData.list = res?.data?.overviews || []
        overviewData.total = res.data.total
      })
    } finally {
      setOverviewLoading(false)
    }
  }  

  async function fetchUserList(value) {
    if (value === 'all') {
      const res = await Http.get('/job/userNames')
      setAppList(res?.data || [])
    } else {
      const res = await Http.get('/job/userNames', {
        params: {
          compute_type: value
        },
      })
      setAppList(res?.data || [])
    } 
  }
  
  useEffect(() => {
    fetchUserList('all')
    getStatistics()
    getList()
    getTotal()
  }, []) 

  useEffect(() => {
    eventEmitter.on('user_export_execl_success', ({message}) => {
      setExporting(false)
    })

    eventEmitter.on('user_export_execl_running', ({message}) => {
      setExporting('正在导出，请耐心等待')
    })

    eventEmitter.on('user_export_execl_error', ({message}) => {
      setExporting(false)
    })

    return () => {
      eventEmitter.off('user_export_execl_success')
      eventEmitter.off('user_export_execl_running')
      eventEmitter.off('user_export_execl_error')
    }
  }, [])


  const columns0 = [
    {
      props: {
        flexGrow: 1,
      },
      header: '用户编号',
      dataKey: 'id',
    },
    {
      props: {
        flexGrow: 2,
      },
      header: '用户名称',
      cell: {
        props: {
          dataKey: 'name',
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        },
      },
    },
    {
      props: {
        flexGrow: 1,
      },
      header: '计算类型',
      cell: {
        props: {
          dataKey: 'computeType',
        },
        render: ({ rowData, dataKey }) => {
          return computeTypesMap[rowData[dataKey]] || '--'
        },
      },
    },
    {
      props: {
        width: 200,
      },
      header: '项目名称',
      cell: {
        props: {
          dataKey: 'project_name',
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        },
      },
    },
    {
      props: {
        flexGrow: 2,
      },
      header: '核时(小时)',
      cell: {
        props: {
          dataKey: 'cpu_time',
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        },
      },
    },
  ]

  const columns1 = [
    {
      props: {
        width: 200,
      },
      header: '作业编号',
      dataKey: 'id',
    },
    {
      props: {
        width: 200,
      },
      header: '作业名称',
      dataKey: 'name',
    },
    {
      props: {
        width: 160,
      },
      header: '计算类型',
      cell: {
        props: {
          dataKey: 'type',
        },
        render: ({ rowData, dataKey }) => {
          return computeTypesMap[rowData[dataKey]] || '--'
        },
      },
    },
    {
      props: {
        width: 200,
      },
      header: '应用名称',
      cell: {
        props: {
          dataKey: 'app_name',
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        },
      },
    },
    {
      props: {
        width: 200,
      },
      header: '项目名称',
      cell: {
        props: {
          dataKey: 'project_name',
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        },
      },
    },
    {
      props: {
        width: 200,
      },
      header: '用户名称',
      cell: {
        props: {
          dataKey: 'user_name',
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        },
      },
    },
    {
      props: {
        width: 200,
      },
      header: '提交时间',
      cell: {
        props: {
          dataKey: 'submit_time',
        },
        render: ({ rowData, dataKey }) => {
          return formatISODateStr(rowData[dataKey], DATE_FORMAT)
        },
      },
    },
    {
      props: {
        width: 200,
      },
      header: '开始时间',
      cell: {
        props: {
          dataKey: 'start_time',
        },
        render: ({ rowData, dataKey }) => {
          return formatISODateStr(rowData[dataKey], DATE_FORMAT)
        },
      },
    },
    {
      props: {
        width: 200,
      },
      header: '结束时间',
      cell: {
        props: {
          dataKey: 'end_time',
        },
        render: ({ rowData, dataKey }) => {
          return formatISODateStr(rowData[dataKey], DATE_FORMAT)
        },
      },
    },
    {
      props: {
        width: 180,
      },
      header: '核时(小时)',
      cell: {
        props: {
          dataKey: 'cpu_time',
        },
        render: ({ rowData, dataKey }) => {
          return rowData[dataKey]
        },
      },
    },
  ]

  const exportExcelFile = (type) => {
    setExporting('正在导出，请耐心等待')
    message.info('正在生成导出文件，请耐心等待')
    
    const params = getParams()
    
    let url = '/api/v1/job/statistics/export?' + qs.stringify({...params, show_type: type}, { arrayFormat: 'repeat'})
    console.debug(url)
    download(url, type)

    setTimeout(() => {
      setExporting(false)
    }, 2000)

    // worker.postMessage({
    //   eventName: _eventName,
    //   eventData: {
    //     params:  {...params, show_type: type},
    //     type,
    //     computeTypesMap,
    //   },
    // })
  }

  const query = (tabKey) => {
    if ( tabKey === 'overview') {
      runInAction(() => {
        overviewPagingData.index = 1
        getStatistics()
        getTotal()
      })
    } else {
      runInAction(() => {
        detailPagingData.index = 1
        getList()
        getTotal()
      })
    }
  }

  return <Wrapper ref={ref}>
    <div className='action'>
     <div className='filter'>
      <div className='item'>
        <span className="label">计算类型: </span>
        <Select
          style={{width: 160}}
          value={queryString.type} 
          onSelect={value => {
            setQueryString({
              ...queryString,
              appIds: [],
              type: value
            })
            fetchUserList(value)
          }}
          >
          <Option key={'-1'} value="all">全部</Option>
          {
            computeTypes.map(t => 
              <Option 
                key={t.compute_type} 
                value={t.compute_type}> 
                {t.show_name} 
              </Option>
          )}
        </Select>
      </div>
      <div className='item'>
        <span className='label'>用户名称: </span>
        <Select 
          style={{width: 200}}
          value={queryString.appIds} 
          maxTagCount={4}
          mode={'multiple'}
          onSelect={(value) => {
            setQueryString({
              ...queryString,
              appIds: [...queryString.appIds, value]
            })
          }}
          onDeselect={(value) => {
            let { appIds } = queryString
            let index = appIds.indexOf(value)
            appIds.splice(index, 1)

            setQueryString({
              ...queryString,
              appIds: [...appIds]
            })
          }}>
          {
            appList.map(app => 
            <Option 
              key={app} 
              value={app}> 
              {app} 
            </Option>
            )
          }
        </Select>  
      </div>
      <div className='item'>
        <span className='label'>项目名称: </span>
        <ProjectSelector value={queryString.projectIds} 
          onSelect={(value) => {
            setQueryString({
              ...queryString,
              projectIds: [...queryString.projectIds, value]
            })
          }}
          onDeselect={(value) => {
            let { projectIds } = queryString
            let index = projectIds.indexOf(value)
            projectIds.splice(index, 1)

            setQueryString({
              ...queryString,
              projectIds: [...projectIds]
            })
          }}
        />
      </div>
      <div className='item'>
        <span className='label'>提交时间: </span>
        <RangePicker
          value={queryString.time as any}
          ranges={GeneralDatePickerRange}
          showTime={{ format: DatePicker_SHOWTIME_FORMAT }}
          format={DatePicker_FORMAT}
          onChange={(dates) => {
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
      <Button className="btn" onClick={() => query(tabKey)}>查询</Button>
     </div>
    </div>
    <Tabs activeKey={tabKey} 
      onChange={key => { 
        setTabKey(key)
        query(key)
      }}
      tabBarExtraContent={
        <>
          <span>总核时(小时): {total}</span>
          <Button 
            disabled={exporting} 
            loading={Boolean(exporting)}
            type="link" 
            onClick={() => exportExcelFile(tabKey)}>
            数据导出
          </Button>
        </>
    }>
      <TabPane tab="统计总览" key={'overview'}>
        <div className='firstTabWrapper'>
          <div className="table">
            <Table
              props={{
                height: rect.height - 160,
                loading: overviewLoading,
                data: overviewData.list,
                rowKey: 'u_id',
              }}
              columns={columns0 as any}
            />
          </div>
        </div>
        <div className='pagination'>
          <Pagination
            showSizeChanger
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
        </div>
      </TabPane>
      <TabPane tab="统计明细" key={'detail'}>
        <Table
          props={{
            height: rect.height - 160,
            loading: detailLoading,
            data: detailData.list,
            rowKey: 'id',
          }}
          columns={columns1 as any}
        />
        <div className='pagination'>
          <Pagination
            showSizeChanger
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
})