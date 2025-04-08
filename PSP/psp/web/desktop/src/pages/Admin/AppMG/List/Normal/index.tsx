import * as React from 'react'
import styled from 'styled-components'
import { observer, inject } from 'mobx-react'
import { observable, action, computed } from 'mobx'
import { Popover } from 'antd'

import {
  AppList,
  FavoriteList,
  RemoteAppList,
  App
} from '@/domain/Applications'
import { Table, Modal } from '@/components'
import Header from '../Header'
import Operators from '../Operators'
import Previewer from './Previewer'
import { StatsBall } from '@/components'
import { statusMap } from '@/domain/Applications/App/utils'

const StyledContent = styled.div`
  height: 100%;
  overflow: auto;

  .name {
    margin-left: 5px;
    color: ${props => props.theme.primaryHighlightColor};
    cursor: pointer;
  }
`
const StyledTemplates = styled.div`
  height: calc(100% - 54px);
  width: 100%;
`

enum StarFilter {
  star,
  noStar
}

enum CloudFilter {
  local,
  cloud
}

enum StateFilter {
  published = 'published',
  unpublished = 'unpublished'
}

interface IProps {
  appList?: AppList
  remoteAppList?: RemoteAppList
  favoriteList?: FavoriteList
}

@inject(({ appList, favoriteList, remoteAppList }) => ({
  appList,
  favoriteList,
  remoteAppList
}))
@observer
export default class ApplicationList extends React.Component<IProps> {
  resizeObserver = null
  @observable loading = true
  @observable _keyword = { current: '' }
  @observable selectedRowKeys = []
  @observable width = 0
  @observable height = 0
  @observable starFilter: StarFilter[] = []
  @observable cloudFilter: CloudFilter[] = []
  @observable stateFilter: StateFilter[] = []
  @action
  updateWidth = width => (this.width = width)
  @action
  updateHeight = height => (this.height = height)
  @action
  updateSelectedRowKeys = keys => (this.selectedRowKeys = keys)
  @action
  updateStarFilter = types => (this.starFilter = types)
  @action
  updateCloudFilter = types => (this.cloudFilter = types)
  @action
  updateStateFilter = types => (this.stateFilter = types)

  tableContainerRef = null

  @computed
  get keyword() {
    return this._keyword.current
  }

  @computed
  get visibleApps() {
    const { appList, favoriteList } = this.props
    let apps = [...appList]

    return (
      apps
        .filter(item => {
          if (
            this.keyword &&
            !item.name.toLowerCase().includes(this.keyword.toLowerCase())
          ) {
            return false
          }

          // match star filter
          const { starFilter } = this
          if (starFilter.length !== 0) {
            if (
              starFilter.includes(StarFilter.star) &&
              starFilter.includes(StarFilter.noStar)
            ) {
              return true
            }

            if (starFilter.includes(StarFilter.star)) {
              return !![...favoriteList].find(
                favorite => favorite.name === item.name
              )
            }

            if (starFilter.includes(StarFilter.noStar)) {
              return ![...favoriteList].find(
                favorite => favorite.name === item.name
              )
            }
          }

          // match cloudFilter
          const { cloudFilter } = this
          if (cloudFilter.length !== 0) {
            if (
              cloudFilter.includes(CloudFilter.local) &&
              cloudFilter.includes(CloudFilter.cloud)
            ) {
              return true
            }

            if (cloudFilter.includes(CloudFilter.cloud)) {
              return !!item.cloudTarget
            }

            if (cloudFilter.includes(CloudFilter.local)) {
              return !item.cloudTarget
            }
          }

          // match stateFilter
          const { stateFilter } = this
          if (stateFilter.length !== 0) {
            if (
              stateFilter.includes(StateFilter.published) &&
              stateFilter.includes(StateFilter.unpublished)
            ) {
              return true
            }

            if (stateFilter.includes(StateFilter.published)) {
              return item.state === StateFilter.published
            }

            if (stateFilter.includes(StateFilter.unpublished)) {
              return item.state === StateFilter.unpublished
            }
          }

          return true
        })
        // specify the properties to activate observer
        .map(item => ({
          name: item.name,
          state: item.state,
          iconData: item.iconData,
          version: item.version,
          cloudOutAppId: item.cloud_out_app_id,
          cloudOutAppName: item.cloud_out_app_name,
          description: item.description,
          status: statusMap[item.state] && statusMap[item.state].text,
          statusColor: statusMap[item.state] && statusMap[item.state].color,
          isInternal: item.isInternal
        }))
    )
  }

  async componentDidMount() {
    const { appList } = this.props
    await appList.fetchTemplates().finally(() => {
      this.loading = false
    })

    this.resizeObserver = new ResizeObserver(entries => {
      for (let entry of entries) {
        this.updateWidth(entry.contentRect.width)
        this.updateHeight(entry.contentRect.height)
      }
    })

    this.resizeObserver.observe(this.tableContainerRef)

    // hack: 处理Table首次加载 bug
    setTimeout(() => {
      this.tableContainerRef.style.paddingRight = '1px'
    }, 3000)
  }

  componentWillUnmount() {
    this.resizeObserver && this.resizeObserver.disconnect()
  }

  @computed
  get localColumns() {
    return [
      {
        props: {
          flexGrow: 1,
          minWidth: 200
        },
        header: '模版名称',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey }) => (
            <>
              {
                <Popover
                  content={
                    <img
                      style={{ width: 225, height: 142 }}
                      src={rowData.iconData || 'img/asset/defaultApp.svg'}
                    />
                  }>
                  <img
                    style={{ width: 20, height: 20 }}
                    src={rowData.iconData || 'img/asset/defaultApp.svg'}
                  />
                </Popover>
              }
              <span
                className='name'
                title={rowData[dataKey]}
                onClick={() => this.preview(rowData[dataKey])}>
                {rowData[dataKey]}
              </span>
            </>
          )
        }
      },
      {
        props: {
          flexGrow: 1,
          width: 100
        },
        header: '模版类型',
        cell: {
          props: {
            dataKey: 'isInternal'
          },
          render: ({ rowData, dataKey }) =>
            rowData[dataKey] ? '内置' : '自定义'
        }
      },
      {
        props: {
          flexGrow: 1,
          width: 100
        },
        header: '版本',
        cell: {
          props: {
            dataKey: 'version'
          }
        }
      },
      {
        props: {
          flexGrow: 1,
          width: 100
        },
        header: '模版描述',
        cell: {
          props: {
            dataKey: 'description'
          }
        }
      },
      {
        props: {
          width: 100
        },
        header: '状态',
        filter: {
          onChange: types => this.updateStateFilter(types),
          items: [
            {
              key: StateFilter.published,
              name: '已发布'
            },
            {
              key: StateFilter.unpublished,
              name: '未发布'
            }
          ]
        },
        cell: {
          props: {
            flexGrow: 1,
            dataKey: 'state'
          },
          render: ({ rowData }) => (
            <StatsBall color={rowData.statusColor}>{rowData.status}</StatsBall>
          )
        }
      },
      {
        props: {
          minWidth: 150,
          flexGrow: 1.1
        },
        header: '操作',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey }) => (
            <Operators model={this.props.appList.get(rowData[dataKey])} />
          )
        }
      }
    ]
  }

  @computed
  get columns() {
    return [
      {
        props: {
          flexGrow: 1,
          minWidth: 200
        },
        header: '模版名称',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey }) => (
            <>
              {
                <Popover
                  content={
                    <img
                      style={{ width: 225, height: 142 }}
                      src={rowData.iconData || 'img/asset/defaultApp.svg'}
                    />
                  }>
                  <img
                    style={{ width: 20, height: 20 }}
                    src={rowData.iconData || 'img/asset/defaultApp.svg'}
                  />
                </Popover>
              }
              <span
                className='name'
                title={rowData[dataKey]}
                onClick={() => this.preview(rowData[dataKey])}>
                {rowData[dataKey]}
              </span>
            </>
          )
        }
      },
      {
        props: {
          flexGrow: 1,
          width: 100
        },
        header: '关联云应用',
        cell: {
          props: {
            dataKey: 'cloudOutAppName'
          }
        }
      },

      {
        props: {
          flexGrow: 1,
          width: 100
        },
        header: '模版类型',
        cell: {
          props: {
            dataKey: 'isInternal'
          },
          render: ({ rowData, dataKey }) =>
            rowData[dataKey] ? '内置' : '自定义'
        }
      },
      {
        props: {
          flexGrow: 1,
          width: 100
        },
        header: '版本',
        cell: {
          props: {
            dataKey: 'version'
          }
        }
      },
      {
        props: {
          flexGrow: 1,
          width: 100
        },
        header: '模版描述',
        cell: {
          props: {
            dataKey: 'description'
          }
        }
      },
      {
        props: {
          width: 100
        },
        header: '状态',
        filter: {
          onChange: types => this.updateStateFilter(types),
          items: [
            {
              key: StateFilter.published,
              name: '已发布'
            },
            {
              key: StateFilter.unpublished,
              name: '未发布'
            }
          ]
        },
        cell: {
          props: {
            flexGrow: 1,
            dataKey: 'state'
          },
          render: ({ rowData }) => (
            <StatsBall color={rowData.statusColor}>{rowData.status}</StatsBall>
          )
        }
      },
      {
        props: {
          minWidth: 150,
          flexGrow: 1.1
        },
        header: '操作',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey }) => (
            <Operators model={this.props.appList.get(rowData[dataKey])} />
          )
        }
      }
    ]
  }

  preview(name) {
    const { application, icon_data } = this.props.appList.get(name).toRequest()
    application.icon = icon_data

    Modal.show({
      title: '模版预览',
      width: 800,
      bodyStyle: { padding: 0, overflow: 'auto' },
      cancelButtonProps: { style: { display: 'none' } },
      content: <Previewer app={new App(application)} />
    })
  }

  private onSelectAll = keys => {
    this.updateSelectedRowKeys(keys)
  }

  private onSelectInvert = () => {
    this.updateSelectedRowKeys([])
  }

  private onSelect = (rowKey, checked) => {
    let keys = this.selectedRowKeys

    if (checked) {
      keys = [...keys, rowKey]
    } else {
      const index = keys.findIndex(item => item === rowKey)
      keys.splice(index, 1)
    }

    this.updateSelectedRowKeys(keys)
  }

  render() {
    return (
      <StyledContent>
        <Header
          keyword={this._keyword}
          selectedRowKeys={this.selectedRowKeys}
        />
        <StyledTemplates ref={ref => (this.tableContainerRef = ref)}>
          <Table
            props={{
              height: this.height,
              data: this.visibleApps,
              rowKey: 'name',
              loading: this.loading,
              locale: {
                emptyMessage: '没有数据',
                loading: '数据加载中...'
              }
            }}
            columns={ this.localColumns as any[] }
            rowSelection={{
              selectedRowKeys: this.selectedRowKeys,
              onSelect: this.onSelect,
              onSelectAll: this.onSelectAll,
              onSelectInvert: this.onSelectInvert
            }}
          />
        </StyledTemplates>
      </StyledContent>
    )
  }
}
