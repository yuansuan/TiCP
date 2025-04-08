import * as React from 'react'
import { Wrapper, TopWrapper, ContentWarpper } from './style'
import { BackButton } from '@/components'
import { Table, Button, Modal } from '@/components'
import { history } from '@/utils'
import { observer } from 'mobx-react'
import { computed, observable, action } from 'mobx'
import userSessionList from '@/domain/SysUserMG/UserSessionList'
import { formatDateFromMilliSec } from '@/utils/formatter'
import { message, Pagination } from 'antd'
import sysUserManager from '@/domain/SysUserMG/SysUserList'

@observer
export default class SysUserDetail extends React.Component<any> {
  wrapperRef = null

  constructor(props) {
    super(props);
    this.wrapperRef = React.createRef();
  }

  @observable sortType = ''
  @observable selectedRowKeys = []
  @observable loading = false
  @observable height = 400

  @action
  updateSelectedRowKeys = keys => {
    this.selectedRowKeys = keys
  }
  @action
  updateSortType = type => (this.sortType = type)

  componentDidMount() {
    this.updateList()

    const rect = this.wrapperRef.current.getBoundingClientRect();
    this.height = rect.height
  }

  updateList = async () => {
    this.loading = true
    try {
      await Promise.all([userSessionList.getUserSessionList(this.user_name)])
    } finally {
      this.loading = false
    }
  }

  get user_name() {
    const { name } = this.props.match.params
    return name
  }

  backSysUserListPage = () => {
    history.push('/sys/onlineUser')
  }

  @computed
  get columns() {
    return [
      {
        props: {
          flexGrow: 1
        },
        header: '会话ID',
        dataKey: 'jti'
      },
      {
        props: {
          flexGrow: 1
        },
        header: 'IP地址',
        dataKey: 'ip'
      },
      {
        props: {
          flexGrow: 1
        },
        header: '会话过期时间',
        cell: {
          props: {
            dataKey: 'expire_time'
          },
          render: ({ rowData, dataKey }) =>
            formatDateFromMilliSec(rowData[dataKey])
        }
      }
    ]
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

  private onSelectAll = keys => {
    this.updateSelectedRowKeys(keys)
  }
  private onSelectInvert = () => {
    this.updateSelectedRowKeys([])
  }

  @computed
  get tableData() {
    return userSessionList.sessionList
  }

  private logout = () => {
    Modal.showConfirm({
      title: '终止登录会话',
      content:
        '如果被选定的会话在进程中，当会话终止时该进程将停止，你确定要终止被选的登录用户会话吗？',
      onOk: async () => {
        await userSessionList.logoutByJwtToken(
          this.selectedRowKeys,
          this.user_name
        )
        this.updateSelectedRowKeys([])
        message.success('此会话已终止')
        userSessionList.getUserSessionList(this.user_name)
      }
    })
  }

  private onPageChange = current => {
    userSessionList.page_index = current
    userSessionList.getUserSessionList(this.user_name)
  }

  private onPageSizeChange = (current, size) => {
    userSessionList.page_index = current
    userSessionList.page_size = size
    userSessionList.getUserSessionList(this.user_name)
  }

  @computed
  get user() {
    return sysUserManager.sysuserList.filter(n => n.name === this.user_name)[0]
  }
  render() {
    return (
      <Wrapper ref={this.wrapperRef}>
        <div>
          <BackButton
            title='返回登录用户管理'
            onClick={this.backSysUserListPage}
            style={{
              fontSize: 20
            }}>
            <>
              <span className='title'>会话详情</span>
              <span className='name'>{this.user?.name}</span>
            </>
          </BackButton>
        </div>
        <ContentWarpper>
          <TopWrapper>
            <Button
              ghost
              type='primary'
              disabled={this.selectedRowKeys.length === 0}
              onClick={this.logout}>
              终止
            </Button>
          </TopWrapper>
          <Table
            props={{
              data: this.tableData,
              height: this.height - 160,
              rowKey: 'jti',
              loading: this.loading,
              locale: {
                emptyMessage: '没有用户数据',
                loading: '数据加载中...'
              }
            }}
            columns={this.columns as any}
            rowSelection={{
              selectedRowKeys: this.selectedRowKeys,
              onSelect: this.onSelect,
              onSelectAll: this.onSelectAll,
              onSelectInvert: this.onSelectInvert
            }}
          />
          <div style={{ textAlign: 'center', padding: 10 }}>
            <Pagination
              showSizeChanger
              onChange={this.onPageChange}
              current={userSessionList.page_index}
              total={userSessionList.total}
              onShowSizeChange={this.onPageSizeChange}
            />
          </div>
        </ContentWarpper>
      </Wrapper>
    )
  }
}
