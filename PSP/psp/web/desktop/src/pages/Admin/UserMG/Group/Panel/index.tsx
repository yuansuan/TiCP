import { computed, reaction } from 'mobx'
import { observer, disposeOnUnmount } from 'mobx-react'
import * as React from 'react'

import { sysConfig } from '@/domain'
import { Group, GroupList, RoleList, UserList } from '@/domain/UserMG'
import { Table, PartialList, Modal, Button } from '@/components'
import { createMobxStream, DataDashPlugin } from '@/utils'
import Operators from './Operators'
import { PanelWrapper } from './style'
import GroupPreview from '../GroupPreview'
import { ListQuery } from '../../utils'
import { untilDestroyed } from '@/utils/operators'
import ApproveStatus from '../../components/ApproveStatus'

interface IProps {
  width: number
  height: number
  loading: boolean
  store: typeof GroupList

  listQuery: ListQuery
  updateTotal: (number) => void
}

@observer
export default class GroupPanel extends React.Component<IProps> {
  state = {
    roleColumnWidth: 0,
  }

  updateRoleColumnWidth = width => {
    this.setState({
      roleColumnWidth: width,
    })
  }

  componentDidMount() {
    createMobxStream(() => this.props.width)
      .pipe(untilDestroyed(this))
      .subscribe(width => {
        this.updateRoleColumnWidth(width * 0.45)
      })
  }

  @computed
  get filteredList() {
    const {
      store,
      listQuery: { query },
    } = this.props
    if (query) {
      return store.groupList.filter(group => group.name.includes(query)) || []
    }
    return store.groupList
  }

  @computed
  get dataSource() {
    const { page, pageSize } = this.props.listQuery

    return this.filteredList
      .slice((page - 1) * pageSize, page * pageSize)
      .map(item => ({
        ...item,
        approve_status_str: item.approve_status_str,
        roles: item.roles.map(r =>
          RoleList.list.get(r) ? RoleList.list.get(r).name : ''
        ),
        users: item.users.map(r =>
          UserList.get(r) ? UserList.get(r).name : ''
        ),
      }))
  }

  @disposeOnUnmount
  disposer = reaction(
    () => this.filteredList,
    () => {
      this.props.updateTotal(this.filteredList.length)
    }
  )

  get columns() {
    const { width } = this.props
    const { roleColumnWidth } = this.state

    const _cols = [
      {
        props: {
          flexGrow: 1,
          minWidth: 200,
        },
        header: '用户组名称',
        cell: {
          props: {
            dataKey: 'name',
          },
          render: ({ rowData, dataKey }) => (
            <div
              className='groupName'
              title={rowData[dataKey]}
              onClick={() =>
                Modal.show({
                  title: '预览用户组',
                  bodyStyle: {
                    height: 710,
                    background: '#F0F5FD',
                    overflow: 'auto',
                  },
                  width: 1130,
                  footer: null,
                  content: () => (
                    <GroupPreview group={new Group(rowData.toRequest())} />
                  ),
                })
              }>
              <span className='name' title={rowData[dataKey]}>
                {rowData[dataKey]}
              </span>
            </div>
          ),
        },
      },
      {
        props: {
          width: width * 0.15,
        },
        header: '成员数',
        cell: {
          props: {
            dataKey: 'users',
          },
          render: ({ rowData, dataKey }) => (
            <span>{rowData[dataKey].length}</span>
          ),
        },
      },
      {
        props: {
          width: roleColumnWidth,
          resizable: true,
          onResize: this.updateRoleColumnWidth,
        },
        header: '角色',
        cell: {
          props: {
            dataKey: 'roles',
          },
          render: ({ rowData }) => {
            return (
              <PartialList
                maxWidth={roleColumnWidth}
                items={rowData.roles.length === 0 ? ['--'] : rowData.roles}
              />
            )
          },
        },
      },
      {
        props: {
          width: width * 0.15,
        },
        header: '审批状态',
        cell: {
          props: {
            dataKey: 'approve_status_str',
          },
          render: ({ rowData }) => {
            const { approve_status_str } = rowData
            return (
              <ApproveStatus
                targetId={rowData.id}
                targetType={'GROUP'}
                data={rowData}
                callback={() => GroupList.fetch()}>
                {approve_status_str === '--' ? (
                  <span style={{ display: 'flex', justifyContent: 'center' }}>
                    --
                  </span>
                ) : (
                  <Button type='link'>{approve_status_str}</Button>
                )}
              </ApproveStatus>
            )
          },
        },
      },
      {
        props: {
          width: width * 0.15,
        },
        header: '操作',
        cell: {
          render: ({ rowData }) => <Operators rowData={rowData} />,
        },
      },
    ]

    if (!sysConfig.enableThreeMembers) {
      _cols.splice(3, 1)
    }

    return _cols
  }

  public render() {
    const { loading } = this.props

    return (
      <PanelWrapper>
        <Table
          props={{
            height: this.props.height,
            data: this.dataSource,
            rowKey: 'id',
            loading,
          }}
          plugins={[new DataDashPlugin()]}
          columns={this.columns}
        />
      </PanelWrapper>
    )
  }
}
