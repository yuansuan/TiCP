/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { message, Tabs, Tooltip } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Table, Button } from '@/components'
import { formatByte } from '@/utils/Validator'
import { formatUnixTime } from '@/utils'
import { JobFileTypeEnum } from '@/constant'
import { useStore } from '../store'
import { env } from '@/domain'
import { showTextEditor } from '@/components'
import { boxServer } from '@/server'
import { Toolbar } from './Toolbar'
import { Monitors } from '@/pages/JobManager/JobList/JobName/Monitors'
import { CodeOutlined } from '@ant-design/icons'
import { showLogsMonitor } from './LogsMonitor'
import { copy2clipboard } from '@/utils'
import { CopyOutlined } from '@ant-design/icons'

const IconStyle: React.CSSProperties = {
  fontSize: '14px',
  fontWeight: 'bold',
  fontFamily: 'PingFangSC-Regular'
}
const StyledLayout = styled.div`
  display: flex;
  padding: 20px 24px;
  background: #fff;
  margin-top: 10px;

  min-height: calc(100vh - 302px);

  .ant-tabs-nav .ant-tabs-tab {
    > div {
      min-width: 88px;
      text-align: center;
    }
  }

  .item {
    width: 100%;

    .action {
      display: flex;
      margin: 4px 0 20px 0;

      .search-input {
        margin-left: auto;
        width: 200px;
      }
    }
  }
`

const { TabPane } = Tabs
const MONITOR_CHARTS_TAB_KEY = -1

export const FileList = observer(function FileList() {
  const store = useStore()
  const { job, jobFile } = store

  const state = useLocalStore(() => ({
    activeTab: 0,
    setActiveTab(key) {
      this.activeTab = key
    },
    get logFileName() {
      return (
        jobFile?.list?.filter(file => /.log$/.test(file.name))[0]?.name ||
        'out.log'
      )
    },
    get displayFiles() {
      return jobFile?.list?.filter(
        file =>
          (this.activeTab === file.type ||
            this.activeTab === JobFileTypeEnum.all) &&
          file.name.includes(store.searchKey)
      )
    },
    get residualVisible() {
      return !!job.runtime?.have_residual
    },
    get monitorVisible() {
      return (
        job?.runtime?.server_params?.map?.MONITOR_CHART_ENABLE?.values[0] ===
          'yes' && job.display_state !== 7
      )
    }, // 排队中的作业不显示监控
    get cloudGraphicVisible() {
      // TODO 先暂时复用餐差图
      return !!job.runtime?.have_residual
    }
  }))

  useEffect(() => {
    store.setSearchKey('')
    store.setSelectedKeys([])
  }, [state.activeTab])

  const preview = data => {
    showTextEditor({
      ...(job?.useRemote
        ? {
            path: data.name,
            sync_id: job?.runtime?.download_task_id
          }
        : {
            path: `${job.id}/${data.name}`,
            bucket: 'result'
          }),
      readonly: true,
      boxServerUtil: boxServer
    })
  }

  // 实时日志查看
  const openLogsTextEditor = () => {
    const path = `${job.id}/${state.logFileName}`
    showLogsMonitor({
      path: path,
      readonly: true,
      ...job
    })
  }

  const LogBtn = (
    <Button
      disabled={job.display_state !== 1}
      icon={<CodeOutlined />}
      onClick={openLogsTextEditor}
      style={IconStyle}>
      {'\u00A0'}实时日志
    </Button>
  )

  const copy = (titleText) => {
    copy2clipboard(titleText)
    message.success('复制文件名成功')
  }

  const renderTitle = (titleText) => {
    return (<div>
      {titleText} <CopyOutlined onClick={e => {copy(titleText); e.stopPropagation()}}/>
    </div>)
  }
  return (
    <StyledLayout>
      <div className='item'>
        <Tabs
          defaultActiveKey='0'
          onChange={key => state.setActiveTab(+key)}
          tabBarExtraContent={LogBtn}>
          <TabPane tab='全部' key={JobFileTypeEnum.all} />
          <TabPane tab='结果' key={JobFileTypeEnum.result} />
          <TabPane tab='模型' key={JobFileTypeEnum.model} />
          <TabPane tab='日志' key={JobFileTypeEnum.log} />
          <TabPane tab='中间文件' key={JobFileTypeEnum.middle} />
          <TabPane tab='其他' key={JobFileTypeEnum.others} />
          {(state.residualVisible ||
            state.monitorVisible ||
            state.cloudGraphicVisible) && (
            <TabPane tab='可视化分析' key={MONITOR_CHARTS_TAB_KEY}>
              <Monitors
                id={job.id}
                jobRuntimeId={job.runtime.id}
                projectId={job.project_id}
                userId={job.user_id}
                jobState={job.display_state}
                residualVisible={state.residualVisible}
                monitorVisible={state.monitorVisible}
                cloudGraphicVisible={state.cloudGraphicVisible}
              />
            </TabPane>
          )}
        </Tabs>
        {state.activeTab !== MONITOR_CHARTS_TAB_KEY && (
          <>
            <Toolbar />
            <Table
              props={{
                data: state.displayFiles,
                virtualized: true,
                rowKey: 'name',
                height: 800,
                shouldUpdateScroll: false
              }}
              rowSelection={{
                selectedRowKeys: store.selectedKeys,
                onChange: keys => store.setSelectedKeys(keys)
              }}
              columns={[
                {
                  header: '名称',
                  props: {
                    fixed: true,
                    flexGrow: 3
                  },
                  dataKey: 'name',
                  cell: {
                    render: ({ rowData }) => (
                      <Button
                        type='link'
                        disabled={
                          env.accountHasDebt && '账户余额不足，请充值后使用'
                        }
                        onClick={() => preview(rowData)}>
                        <Tooltip title={renderTitle(rowData.name)}>
                          {rowData.name}
                        </Tooltip>
                      </Button>
                    )
                  }
                },
                {
                  header: '创建时间',
                  props: {
                    fixed: true,
                    flexGrow: 2
                  },
                  dataKey: 'lastModified',
                  cell: {
                    render: ({ rowData }) => (
                      <div>{formatUnixTime(rowData.mod_time)}</div>
                    )
                  }
                },
                {
                  header: '大小',
                  props: {
                    fixed: true,
                    flexGrow: 1
                  },
                  dataKey: 'size',
                  cell: {
                    render: ({ rowData }) => (
                      <div>{formatByte(rowData.size)}</div>
                    )
                  }
                }
              ]}
            />
          </>
        )}
      </div>
    </StyledLayout>
  )
})
