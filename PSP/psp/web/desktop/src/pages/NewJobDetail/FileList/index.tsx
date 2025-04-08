/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Table, message, Tabs, Tooltip } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Button, Table as YSTable } from '@/components'
import { formatByte } from '@/utils/Validator'
import { formatUnixTime } from '@/utils'
import { JobFileTypeEnum } from '@/constant'
import { useStore } from '../store'
import { env } from '@/domain'
import { showTextEditor } from '@/components'
import { newBoxServer } from '@/server'
import { Toolbar } from './Toolbar'
import { copy2clipboard } from '@/utils'
import { CopyOutlined, FolderOutlined, FileOutlined } from '@ant-design/icons'

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

export const FileList = observer(function FileList() {
  const store = useStore()
  const { job, jobFile, expandedRowKeys, setExpandedRowKeys } = store
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
    // get displayFiles() {
    //   console.log('jobFile?.list: ', jobFile?.list);
    //   return jobFile?.list?.filter(
    //     file =>
    //       (this.activeTab === file.type ||
    //         this.activeTab === JobFileTypeEnum.all) &&
    //       file.name.includes(store.searchKey)
    //   )
    // }
    get displayFiles() {
      return jobFile?.list.filter(file => file.name.includes(store.searchKey))
    },
    data: jobFile?.list || [],
    prepareData(serverData) {
      // 将服务器数据进行结构调整，添加children属性
      return serverData.map(item => ({
        ...item,
        key: item.path,
        children: item.is_dir ? [] : null
      }))
    },

    async handleExpand(expanded, record) {
      if (expanded && record.is_dir && !record.children.length) {
        // 仅当展开目录且尚未加载子数据时，进行异步请求加载
        const childrenData = await jobFile.fetch({
          path: record.path,
          is_cloud: job.isCloud,
          user_name: job.user_name
        })

        const newData = jobFile?.list.map(item =>
          item.path === record.path
            ? { ...item, children: this.prepareData(childrenData) }
            : item
        )
      }
    },

    setData(serverData) {
      this.data = this.prepareData(serverData)
    }
  }))

  useEffect(() => {
    store.setSearchKey('')
    // store.setSelectedKeys([])
  }, [state.activeTab])

  const preview = data => {
    const { size, name, path, type } = data
    showTextEditor({
      path,
      fileInfo: {
        size,
        name,
        path,
        type,
        is_cloud: job.isCloud,
        user_name: job.user_name
      },
      readonly: true,
      boxServerUtil: newBoxServer
    })
  }

  const copy = titleText => {
    copy2clipboard(titleText)
    message.success('复制文件名成功')
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

  useEffect(() => {
    state.setData(jobFile.list)
  }, [state])

  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      render: (text, record) => (
        <Button type='link'>
          <Tooltip title={renderTitle(record.name)}>
            {record.is_dir ? (
              <>
                <FolderOutlined /> {record.name}
              </>
            ) : (
              <div onClick={() => preview(record)}>
                <FileOutlined />
                {record.name}
              </div>
            )}
          </Tooltip>
        </Button>
      )
    },
    {
      title: '创建时间',
      dataIndex: 'lastModified',
      render: (text, record) => <div>{formatUnixTime(record.m_date)}</div>
    },
    {
      title: '大小',
      dataIndex: 'size',
      render: (text, record) => <div>{formatByte(record.size)}</div>
    }
  ]

  return (
    <StyledLayout>
      <div className='item'>
        <Tabs defaultActiveKey='0' onChange={key => state.setActiveTab(+key)}>
          <TabPane tab='全部' key={JobFileTypeEnum.all} />
        </Tabs>
        {
          <>
            <Toolbar />
            <Table
              columns={columns}
              rowKey={record => record.name}
              rowSelection={{
                selectedRowKeys: store.selectedKeys,
                onChange: keys => store.setSelectedKeys(keys)
              }}
              dataSource={state.displayFiles}
              expandable={{
                onExpand: state.handleExpand
              }}
            />
          </>
        }
      </div>
    </StyledLayout>
  )
})
