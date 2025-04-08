/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useDispatch } from 'react-redux'
import { useStore } from './store'
import { message, Pagination, Tooltip } from 'antd'
import { Table } from '@/components'
import {
  buryPoint,
  getDisplayRunTime,
  copy2clipboard,
  getComputeType
} from '@/utils'
import { JobName } from './JobName'
import { JobStatus, Status } from '@/components'
import { Operator } from './Operator'
import { Toolbar } from './Toolbar'
import { runInAction } from 'mobx'
import { CopyOutlined } from '@ant-design/icons'
import { ALL_JOB_STATES, DATASTATEMAP } from '@/constant'
export { Context, useModel } from './store'

const StyledLayout = styled.div`
  > .body {
    padding: 15px 20px;
    display: flex;
    flex-direction: column;

    .pagination {
      margin: 10px auto;
    }
  }
`

const NameStyled = styled.div`
  max-width: calc(100%);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  color: ${({ theme }) => theme.primaryColor};
`

type Props = {
  showJobSetName?: boolean
  hidePagination?: boolean
  operatable?: boolean
  onNameClick?: (id: string) => void
  height?: number
  width?: number
}

export const JobList = observer(function JobList({
  showJobSetName = true,
  hidePagination = false,
  operatable = true,
  onNameClick,
  height,
  width
}: Props) {
  const store = useStore()
  const dispatch = useDispatch()
  const { model } = store
  const state = useLocalStore(() => ({
    get dataSource() {
      return model.list.map(item => {
        return {
          ...item,
          terminalable: item.terminalable,
          resubmittable: item.resubmittable,
          displayCreateTime: item.submit_time,
          displayStartTime: item.start_time,
          displayEndTime: item?.end_time,
          displayState: ALL_JOB_STATES[item.state],
          residualVisible: item.enable_residual,
          snapshotVisible: item.enable_snapshot
        }
      })
    }
  }))

  function onPageChange(index, size) {
    runInAction(() => {
      store.setPageIndex(index)
      store.setPageSize(size)
      store.setSelectedKeys([])
    })
  }

  function onJobIdClick(rowData) {
    const { id, isCloud } = rowData
    buryPoint({
      category: '作业管理',
      action: '作业编号'
    })
    window.localStorage.setItem(
      'CURRENTROUTERPATH',
      `/new-job?jobId=${id}&isCloud=${isCloud}`
    )

    dispatch({
      type: 'NEWJODETAIL',
      payload: 'close'
    })

    setTimeout(() => {
      dispatch({
        type: 'JOBMANAGE',
        payload: 'full'
      })
      dispatch({
        type: 'NEWJODETAIL',
        payload: 'full'
      })
    }, 500)
  }

  function onJobSetIdOrNameClick(rowData) {
    const { job_set_id, isCloud } = rowData
    buryPoint({
      category: '作业管理',
      action: '作业集编号'
    })
    window.localStorage.setItem(
      'CURRENTROUTERPATH',
      `/new-job-set?jobSetId=${job_set_id}&isCloud=${isCloud}`
    )

    dispatch({
      type: 'NEWJOBSETDETAIL',
      payload: 'close'
    })

    setTimeout(() => {
      dispatch({
        type: 'JOBMANAGE',
        payload: 'full'
      })
      dispatch({
        type: 'NEWJOBSETDETAIL',
        payload: 'full'
      })
    }, 500)
  }

  const copy = titleText => {
    copy2clipboard(titleText)
    message.success('复制成功')
  }

  const renderTitle = titleText => {
    return (
      <div>
        {titleText}{' '}
        <CopyOutlined
          onClick={e => {
            copy(titleText)
            e.stopPropagation()
          }}
        />
      </div>
    )
  }

  return (
    <StyledLayout>
      <Toolbar />
      <div className='body'>
        <Table
          tableId={operatable ? 'job_list_table' : undefined}
          defaultConfig={[
            { key: 'id', active: true },
            { key: 'name', active: true },
            showJobSetName && { key: 'job_set_id', active: true },
            showJobSetName && { key: 'job_set_name', active: true },
            { key: 'user_name', active: true },
            { key: 'app_name', active: true },
            { key: 'queue', active: true },
            { key: 'project_name', active: true },
            { key: 'type', active: true },
            { key: 'cpus_alloc', active: true },
            { key: 'exec_duration', active: true },
            { key: 'displayCreateTime', active: true },
            { key: 'displayStartTime', active: true },
            { key: 'displayEndTime', active: true },
            { key: 'displayState', active: true },
            { key: 'data_state', active: true },
            { key: 'options', active: true }
          ].filter(Boolean)}
          props={{
            data: state.dataSource,
            rowKey: 'out_job_id',
            height: height,
            shouldUpdateScroll: false
          }}
          rowSelection={{
            selectedKeys: store.selectedKeys,
            onChange: keys => {
              store.setSelectedKeys(keys)
            },
            props: {
              fixed: true
            }
          }}
          columns={[
            {
              header: '作业编号',
              props: {
                width: 130,
                fixed: 'left' as 'left'
              },
              dataKey: 'id',
              cell: {
                props: {
                  dataKey: 'id'
                },
                render: ({ rowData }) => (
                  <NameStyled
                    title={rowData['id']}
                    onClick={() => onJobIdClick(rowData)}>
                    <Tooltip title={renderTitle(rowData['id'])}>
                      {rowData['id']}
                    </Tooltip>
                  </NameStyled>
                )
              }
            },
            {
              header: '作业名称',
              headerClassName: 'table-job-name-header',
              props: {
                resizable: true,
                width: 200
              },
              dataKey: 'name',
              cell: {
                props: {
                  dataKey: 'name'
                },
                render: ({ rowData }) => (
                  <JobName
                    onClick={onNameClick}
                    {...rowData}
                    userId={rowData.user_id}
                    cloudGraphicVisible={rowData.snapshotVisible}
                    residualVisible={rowData.residualVisible}
                  />
                )
              }
            },
            {
              header: '作业集编号',
              props: {
                width: 130
              },
              dataKey: 'job_set_id',
              cell: {
                props: {
                  dataKey: 'job_set_id'
                },
                render: ({ rowData }) => (
                  <>
                    {rowData['job_set_id'] === '' && '--'}
                    <NameStyled
                      title={rowData['job_set_id']}
                      onClick={() => onJobSetIdOrNameClick(rowData)}>
                      <Tooltip title={renderTitle(rowData['job_set_id'])}>
                        {rowData['job_set_id']}
                      </Tooltip>
                    </NameStyled>
                  </>
                )
              }
            },
            {
              header: '作业集名称',
              props: {
                resizable: true,
                width: 200
              },
              dataKey: 'job_set_name',
              cell: {
                props: {
                  dataKey: 'job_set_name'
                },
                render: ({ rowData }) => (
                  <>
                    {rowData['job_set_name'] === '' && '--'}
                    <NameStyled
                      title={rowData['job_set_name']}
                      onClick={() => onJobSetIdOrNameClick(rowData)}>
                      {rowData['job_set_name']}
                    </NameStyled>
                  </>
                )
              }
            },
            {
              header: '用户名称',
              props: {
                width: 120
              },
              dataKey: 'user_name',
              cell: {
                render: ({ rowData }) => {
                  return <div>{rowData.user_name}</div>
                }
              }
            },

            {
              header: '应用名称',
              props: {
                width: 200
              },
              dataKey: 'app_name'
            },
            {
              header: '队列名称',
              props: {
                width: 120
              },
              dataKey: 'queue'
            },
            {
              header: '项目名称',
              props: {
                width: 200
              },
              dataKey: 'project_name'
            },
            {
              header: '作业类型',
              props: {
                width: 120
              },
              dataKey: 'type',
              cell: {
                props: {
                  dataKey: 'type'
                },
                render: ({ rowData, dataKey }) => {
                  return <span>{getComputeType(rowData['type'])}</span>
                }
              }
            },
            {
              header: '核数',
              props: {
                width: 100
              },
              dataKey: 'cpus_alloc'
            },
            {
              header: '提交时间',
              props: {
                width: 200
              },
              dataKey: 'displayCreateTime'
            },
            {
              header: '开始时间',
              props: {
                width: 200
              },
              dataKey: 'displayStartTime'
            },
            {
              header: '结束时间',
              props: {
                width: 200
              },
              dataKey: 'displayEndTime'
            },
            {
              header: '运行时长',
              props: {
                width: 140
              },
              cell: {
                props: {
                  dataKey: 'exec_duration'
                },
                render: ({ rowData, dataKey }) =>
                  getDisplayRunTime(rowData[dataKey])
              }
            },
            {
              header: '计算状态',
              props: {
                width: 130,
                fixed: 'right' as 'right'
              },
              cell: {
                props: {
                  dataKey: 'displayState'
                },
                render: ({ rowData }) => {
                  return <JobStatus showDropDown={true} job={rowData} />
                }
              }
            },
            {
              header: '数据状态',
              props: {
                width: 130,
                fixed: 'right' as 'right'
              },
              cell: {
                props: {
                  dataKey: 'data_state'
                },
                render: ({ rowData }) => {
                  return (
                    <Status
                      text={DATASTATEMAP[rowData?.data_state]?.text}
                      type={DATASTATEMAP[rowData?.data_state]?.type}
                    />
                  )
                }
              }
            },

            operatable && {
              header: '操作',
              props: {
                width: 130,
                fixed: 'right' as 'right'
              },
              cell: {
                props: {
                  dataKey: 'options'
                },
                render: ({ rowData }) => <Operator {...rowData} />
              }
            }
          ].filter(Boolean)}
        />
        {!hidePagination && (
          <Pagination
            className='pagination'
            showSizeChanger
            pageSize={store.pageSize}
            pageSizeOptions={['10', '20', '50', '100']}
            current={store.pageIndex}
            total={model.page_ctx.total}
            onChange={onPageChange}
          />
        )}
      </div>
    </StyledLayout>
  )
})
