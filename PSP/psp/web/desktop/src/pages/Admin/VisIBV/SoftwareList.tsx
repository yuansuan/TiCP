import React, { useState } from 'react'
import styled from 'styled-components'
import { Button, Icon, Modal } from '@/components'
import { Http } from '@/utils'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from './store'
import { Divider, message, Table, Tag } from 'antd'
import SoftwareEditor from './SoftwareEditor'
import PresetEditor from './PresetEditor'
import { RemoteAppEditor } from './RemoteAppEditor'
import { Popover } from 'antd'

interface StyledLayoutProp {
  height: number
}

const StyledLayout = styled.div<StyledLayoutProp>`
  height: ${props => props.height}px;
  .name {
    cursor: pointer;

    &:hover {
      color: ${({ theme }) => theme.primaryHighlightColor};
    }
  }

  &.rs-table {
    border: 0;
  }

  th,
  td {
    padding: 16px;
    border-bottom: 1px solid #f3f5f8;
  }

  thead {
    th {
      text-align: left;
      background: #f3f5f8;
    }
  }

  tbody > tr:hover > td {
    background: #e8e8e8;
  }
  &.rs-table-hover .rs-table-row {
    &:hover {
      background: ${({ theme }) => theme.backgroundColorHover};

      .rs-table-cell {
        background: ${({ theme }) => theme.backgroundColorHover};
      }

      .rs-table-cell-group {
        background: ${({ theme }) => theme.backgroundColorHover};
      }
    }
  }

  th.padding,
  td.padding {
    padding: 0;
  }
`

interface IProps {
  height: number
}

export const SoftwareList = observer(function List(props: IProps) {
  const store = useStore()
  const [visible, setVisible] = useState(false)
  const state = useLocalStore(() => ({
    get dataSource() {
      return store.software.softwareList?.map(item => ({
        ...item,
        gpu: item.gpu_desired ? '是' : '否',
        isPublish: item.state === 'published' ? '已发布' : '未发布',
        enabled: item.state === 'published' ? true : false
      }))
    }
  }))

  async function isPublish(rowData) {
    await Http.put('/vis/software/publish', {
      id: rowData.id,
      state: rowData.enabled ? 'unpublished' : 'published'
    })
    store.refreshSoftware()
    message.success(rowData.enabled ? '取消发布成功' : '发布成功')
  }

  function edit(rowData) {
    Modal.show({
      title: '编辑镜像',
      width: 600,
      bodyStyle: { padding: 0, height: 640 },
      footer: null,
      content: ({ onCancel, onOk }) => (
        <SoftwareEditor
          softwareItem={rowData}
          onCancel={onCancel}
          onOk={() => {
            onOk()
            store.refreshSoftware()
          }}
        />
      )
    })
  }

  function remoteAppSetting(rowData) {
    Modal.show({
      title: '远程应用设置',
      width: 800,
      bodyStyle: { padding: 0, height: 460 },
      footer: null,
      content: ({ onCancel, onOk }) => (
        <RemoteAppEditor
          remoteAppList={rowData.remote_apps}
          softwareId={rowData.id}
          onCancel={onCancel}
          onOk={() => {
            onOk()
            store.refreshSoftware()
          }}
        />
      )
    })
  }
  function association(rowData) {
    Modal.show({
      title: '预设',
      width: 800,
      bodyStyle: { padding: 0, height: 460 },
      footer: null,
      content: ({ onCancel, onOk }) => (
        <PresetEditor
          softwareItem={rowData}
          onCancel={onCancel}
          onOk={() => {
            onOk()
            store.refreshSoftware()
          }}
        />
      )
    })
  }

  function ellipsis(text, title) {
    const showScript = (
      <>
        {text.length > 50 ? (
          <>
            {text.substring(0, 47)}...
            <Button
              type='link'
              style={{ padding: 0 }}
              onClick={() => {
                return Modal.show({
                  title,
                  width: 800,
                  footer: null,
                  content: () => text
                })
              }}>
              展示更多
            </Button>
          </>
        ) : (
          text
        )}
      </>
    )

    return showScript
  }

  const deleteSoftware = ({ id }) => {
    Modal.confirm({
      title: '删除镜像',
      content: '确认删除！',
      okText: '确认',
      visible,
      cancelText: '取消',
      onOk: async () => {
        await Http.delete(`/vis/software?id=${id}`)
          .then(res => {
            if (res.success) {
            }
          })
          .finally(() => {
            setVisible(false)
          })
        store.refreshSoftware()
      }
    })
  }
  const columns = [
    {
      width: 180,
      title: '镜像名称',
      dataIndex: 'name',
      key: 'name',
      render: (text, record) => (
        <>
          {
            <Popover
              content={
                <img
                  style={{ width: 225, height: 142 }}
                  src={record.icon || 'img/icon/3dcloudApp.png'}
                />
              }>
              <img
                style={{ width: 20, height: 20 }}
                src={record.icon || 'img/icon/3dcloudApp.png'}
              />
            </Popover>
          }
          <span className='name' title={record.name}>
            {record.name}
          </span>
        </>
      )
    },
    {
      width: 180,
      title: '镜像描述',
      dataIndex: 'desc',
      key: 'desc',
      render: (text, record, index) => ellipsis(text, '镜像描述')
    },

    {
      width: 120,
      title: '操作平台',
      dataIndex: 'platform'
    },

    {
      width: 80,
      title: 'GPU',
      dataIndex: 'gpu'
    },
    {
      width: 180,
      title: '初始化脚本',
      dataIndex: 'init_script',
      render: (text, record, index) => ellipsis(text, '初始化脚本')
    },
    {
      width: 200,
      title: '镜像编号',
      dataIndex: 'image_id'
    },
    {
      width: 100,
      title: '状态',
      dataIndex: 'isPublish'
    },
    {
      width: 260,
      title: '操作',
      dataIndex: 'id',
      key: 'id',
      fixed: 'right',
      render: (text, record, index) => {
        return (
          <div className='rowData'>
            <Button
              type='link'
              style={{ padding: 0 }}
              onClick={() => {
                isPublish(record)
              }}>
              {record['enabled'] ? '取消发布' : '发布'}
            </Button>
            <Divider type='vertical' />
            <Button
              type='link'
              disabled={record['enabled']}
              style={{ padding: 0 }}
              onClick={() => edit(record)}>
              编辑
            </Button>
            <Divider type='vertical' />
            <Button
              type='link'
              style={{ padding: 0 }}
              onClick={() => association(record)}>
              预设
            </Button>
            <Divider type='vertical' />
            <Button
              type='link'
              disabled={record['enabled']}
              style={{ padding: 0 }}
              onClick={() => deleteSoftware(record)}>
              删除
            </Button>
          </div>
        )
      }
    }
  ] as any

  return (
    <StyledLayout height={props.height - 110 || 300}>
      <Table
        rowKey='id'
        dataSource={state.dataSource}
        columns={columns}
        loading={store.loading}
        scroll={{ y: props.height - 180 || 300 }}
        pagination={false}
        expandedRowRender={record => {
          return (
            <div style={{ margin: 0 }}>
              实例列表：
              {record?.presets &&
                record?.presets?.map(item => (
                  <Tag color='#EAEAEA' key={item.id}>
                    <span style={{ color: '#666666' }}>{item.name}</span>
                  </Tag>
                ))}
            </div>
          )
        }}
      />
    </StyledLayout>
  )
})
