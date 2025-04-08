import * as React from 'react'
import styled from 'styled-components'
import { computed, observable } from 'mobx'
import { observer } from 'mobx-react'
import { sofwareAppList } from '@/domain/Visual'
import { Table, Button, Modal } from '@/components'
import { Tag, message,  Descriptions} from 'antd'
import { Search } from '@/components'
import { SoftwareAppForm } from './SoftwareAppForm'
import { BindForm } from './BindForm'

const Wrapper = styled.div`
  padding: 20px;
  .action {
    display: flex;
    justify-content: space-between;
    margin-bottom: 20px;
  }

  .table {
  }
`

@observer
export default class SoftwareMgr extends React.Component<any> {
  @observable searchKey = ''

  fetch = async () => {
    await sofwareAppList.fetch()
  }

  componentDidMount = () => {
    this.fetch()
  }

  previewSoftwareInfo = (rowData) => {
    Modal.show({
      title: '预览软件信息',
      content: () => {
        return <>
          <Descriptions title="软件" column={1}>
            <Descriptions.Item label="软件名称">{rowData['name']}</Descriptions.Item>
            <Descriptions.Item label="软件版本">{rowData['version'] || '--'}</Descriptions.Item>
            <Descriptions.Item label="软件图标">
              {rowData['icon_data'] ? <img height='44px' src={rowData['icon_data']} /> : '--'}
            </Descriptions.Item>
            <Descriptions.Item label="是否显示桌面">{rowData['only_show_desktop'] ? '是' : '否'}</Descriptions.Item>
            <Descriptions.Item label="软件路径">{rowData['path'] || '--'}</Descriptions.Item>
            <Descriptions.Item label="操作系统">{rowData['os_type'] || '--'}</Descriptions.Item>
            <Descriptions.Item label="是否支持GPU">{rowData['gpu_support'] ? '是' : '否'}</Descriptions.Item>
          </Descriptions>
          <Descriptions title="参数" column={1}>
            <Descriptions.Item label="参数">{rowData['app_param'] || '--' }</Descriptions.Item>
            <Descriptions.Item label="参数路径">{rowData['app_param_paths'] || '--'}</Descriptions.Item>
          </Descriptions>
          <Descriptions title="关联工作站" column={1}>
            <Descriptions.Item label="工作站">
              {
                rowData['WS_list'].length !== 0 ? rowData['WS_list'].map(s => {
                  return <Tag key={s.name}>{s.name}</Tag>
                }) : '--'
               }
            </Descriptions.Item>
          </Descriptions>
        </>
      }
    })
  }

  openAddForm = () => {
    Modal.show({
      title: '添加软件',
      footer: null,
      content: ({ onCancel, onOk }) => {
        const ok = async data => {
          const res = await sofwareAppList.add(data)
          if (res.success) {
            message.success('软件添加成功')
            this.fetch()
            onOk()
          } else {
            message.error('软件添加失败')
          }

          return res
        }
        return <SoftwareAppForm onOk={ok} onCancel={onCancel} />
      },
      width: 700,
    })
  }

  openEditForm = (rowData) => {
    Modal.show({
      title: '编辑软件',
      footer: null,
      content: ({ onCancel, onOk }) => {
        const ok = async data => {
          const res = await sofwareAppList.edit(data)
          if (res.success) {
            message.success('软件编辑成功')
            this.fetch()
            onOk()
          } else {
            message.error('软件编辑失败')
          }

          return res
        }
        return <SoftwareAppForm onOk={ok} onCancel={onCancel} data={rowData}/>
      },
      width: 700,
    })
  }

  openBindForm = (rowData) => {
    Modal.show({
      title: '编辑关联工作站',
      footer: null,
      content: ({ onCancel, onOk }) => {
        const ok = async data => {
          const res = await sofwareAppList.bindOrUnBind(data)
          
          if (res.every(r => r.success)) {
            message.success('编辑关联工作站成功')
            this.fetch()
            onOk()
          } else {
            message.error('编辑关联工作站失败')
          }

          return res
        }
        return <BindForm onOk={ok} onCancel={onCancel} data={rowData}/>
      },
      width: 700,
    })
  }

  @computed
  get columns() {
    return [
      {
        props: {
          resizable: true,
          fixed: 'left',
        },
        header: 'ID',
        dataKey: 'id',
      },
      {
        props: {
          resizable: true,
          fixed: 'left',
        },
        header: '软件名称',
        cell: {
          props: {
            dataKey: 'name',
          },
          render: ({ rowData, dataKey }) => {
            return <a onClick={() => this.previewSoftwareInfo(rowData)}> {rowData[dataKey] }</a>
          },
        },
      },
      {
        props: {
          resizable: true,
          fixed: 'left',
        },
        header: '软件版本',
        dataKey: 'version',
      },
      {
        props: {
          resizable: true,
        },
        header: '软件图标',
        cell: {
          props: {
            dataKey: 'icon_data',
          },
          render: ({ rowData, dataKey }) => {
            return rowData[dataKey] ? <img height='44px' src={rowData[dataKey]} /> : '--'
          },
        },
      },
      {
        props: {
          resizable: true,
        },
        header: '是否显示桌面',
        cell: {
          props: {
            dataKey: 'only_show_desktop',
          },
          render: ({ rowData, dataKey }) => {
            return <div>{rowData[dataKey] ? '是' : '否'}</div>
          },
        },
      },
      {
        props: {
          resizable: true,
        },
        header: '软件路径',
        dataKey: 'path',
      },
      {
        props: {
          resizable: true,
        },
        header: '操作系统',
        dataKey: 'os_type',
      },
      {
        props: {
          resizable: true,
        },
        header: '是否支持GPU',
        cell: {
          props: {
            dataKey: 'gpu_support',
          },
          render: ({ rowData, dataKey }) => {
            return <div>{rowData[dataKey] ? '是' : '否'}</div>
          },
        },
      },
      {
        props: {
          resizable: true,
        },
        header: '参数',
        dataKey: 'app_param',
      },
      {
        props: {
          resizable: true,
        },
        header: '参数路径',
        dataKey: 'app_param_paths',
      },
      {
        props: {
          resizable: true,
        },
        header: '工作站',
        cell: {
          props: {
            dataKey: 'WS_list',
          },
          render: ({ rowData, dataKey }) => {
            return rowData[dataKey].map(s => {
              return <Tag key={s.name}>{s.name}</Tag>
            }) || '--'
          },
        },
      },
      {
        header: '操作',
        props: {
          flexGrow: 1,
          minWidth: 240,
          fixed: 'right',
        },
        cell: {
          render: ({ rowData }) => {
            return (
              <>
              <Button
                type='link'
                onClick={async () => {
                  this.openBindForm(rowData)
                }}>
                关联工作站
              </Button>
              <Button
                type='link'
                disabled={rowData['WS_list'].length !== 0 ? '无法进行编辑操作，请取消已关联的工作站' : false}
                onClick={async () => {
                  this.openEditForm(rowData)
                }}>
                编辑
              </Button>
              <Button
                type='link'
                disabled={rowData['WS_list'].length !== 0 ? '无法进行删除操作，请取消已关联的工作站' : false}
                onClick={async () => {
                  await Modal.showConfirm({
                    title: '确认',
                    content: `确认删除软件${rowData['name']}?`,
                  })
                  try {
                    await sofwareAppList.delete(rowData)
                    message.success('删除软件成功')
                  } finally {
                    this.fetch()
                  }
                }}>
                删除
              </Button>
              </>
            )
          },
        },
      },
    ]
  }

  @computed
  get filteredAppList() {
    return this.searchKey ? [...sofwareAppList].filter(app => app.name.includes(this.searchKey)) : [...sofwareAppList]
  }

  render() {
    return (
      <Wrapper>
        <div className='action'>
          <Button onClick={() => {this.openAddForm()}}>添加软件</Button>
          <Search
            placeholder={'输入软件名称搜索'}
            debounceWait={300}
            onSearch={value => {
              this.searchKey = value
            }}
          />
        </div>
        <Table
          props={{
            autoHeight: true,
            loading: sofwareAppList.loading,
            data: this.filteredAppList,
            rowKey: 'id',
          }}
          columns={this.columns as any}
        />
      </Wrapper>
    )
  }
}
