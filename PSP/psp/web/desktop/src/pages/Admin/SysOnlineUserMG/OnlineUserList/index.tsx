import * as React from 'react'
import { observer } from 'mobx-react'
import { history } from '@/utils'
import { Wrapper, TopWrapper, TableWrapper } from './style'
import { Search } from '@/components'
import { Table, Button, Modal } from '@/components'
import { computed, observable, action } from 'mobx'
import sysUserManager from '@/domain/SysUserMG/SysUserList'
import { message } from 'antd'
import { Pagination } from 'antd'
import { sortByStartTime } from '../common'
@observer
export default class SysUserList extends React.Component<any> {
  wrapperRef = null
  resizeObserver = null

  constructor(props) {
    super(props);
    this.wrapperRef = React.createRef();
  }

  @observable searchStr = ''
  @observable sortType = ''
  @observable selectedRowKeys = []
  @observable loading = true
  @observable current = 1
  @observable pageSize = 10
  @observable height = 400

  @action
  updateSelectedRowKeys = keys => (this.selectedRowKeys = keys)

  @computed
  get userNames() {
    return sysUserManager.sysuserList.map(user => user.name)
  }

  @action
  updateSortType = type => (this.sortType = type)

  componentDidMount() {
    sysUserManager.getSysUserList().finally(() => {
      this.loading = false
    })

    this.resizeObserver = new ResizeObserver((entries) => {
      for (let entry of entries) {
          this.height = entry.contentRect.height
      }
    })
    
    this.resizeObserver.observe(this.wrapperRef.current)

    // hack: 处理Table首次加载 bug
    setTimeout(() => {
      this.wrapperRef.current.style.paddingRight = '1px'
    }, 3000)
  }

  componentWillUnmount(): void {
    this.resizeObserver && this.resizeObserver.disconnect()
  }

  @computed
  get tableData() {
    const filterData = this.filterSearchStr(sysUserManager.sysuserList)
    //sort user_name by con_starttime(会话开始时间)
    return sortByStartTime(filterData, this.sortType)
  }

  @computed
  get pagingTableData() {
    const res = this.tableData.slice(
      (this.current - 1) * this.pageSize,
      this.current * this.pageSize
    )
    return res
  }

  private filterSearchStr(list) {
    return !this.searchStr
      ? list
      : list.filter(user => user.user_name.includes(this.searchStr))
  }

  @computed
  get columns(): any {
    const cols = [
      {
        props: {
          flexGrow: 1
        },
        header: '用户',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData }) => {
            return (
              <a
                title={rowData.name}
                style={{ color: '#1458E0' }}
                onClick={() => {
                  this.onRowClick && this.onRowClick({ rowData })
                }}>
                {rowData.name}
              </a>
            )
          }
        },
      },
      {
        props: {
          flexGrow: 1
        },
        header: '会话数',
        dataKey: 'count'
      }
    ]
    return cols
  }

  private logout = () => {
    Modal.showConfirm({
      title: '注销用户',
      content: '此操作将结束所选的用户会话和相关操作。您确定要继续吗？',
      onOk: async () => {
        await sysUserManager.logoutByUserName(this.selectedRowKeys)
        this.updateSelectedRowKeys([])
        sysUserManager.getSysUserList()
        message.success('注销成功')
      }
    })
  }

  private onSearch = (value: string) => {
    sysUserManager.getSysUserList(value)
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

  onRowClick = ({ rowData }) => {
    history.push(`/sys/onlineUser/${rowData.name}`)
  }

  private onPageChange = current => {
    sysUserManager.page_index = current
    sysUserManager.getSysUserList()
  }

  private onPageSizeChange = (current, size) => {
    sysUserManager.page_index = current
    sysUserManager.page_size = size
    sysUserManager.getSysUserList()
  }

  render() {
    return (
      <Wrapper ref={this.wrapperRef}>
        <TableWrapper>
          <TopWrapper>
            <Button
              ghost
              type='primary'
              disabled={this.selectedRowKeys.length === 0}
              onClick={this.logout}>
              注销
            </Button>
            <Search onSearch={this.onSearch} placeholder={'请输入用户名称'} />
          </TopWrapper>
          <Table
            columns={this.columns}
            props={{
              data: this.pagingTableData,
              rowKey: 'name',
              height: this.height - 140,
              loading: this.loading,
              locale: {
                emptyMessage: '没有用户数据',
                loading: '数据加载中...'
              }
            }}
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
              current={sysUserManager.page_index}
              total={sysUserManager.total}
              onShowSizeChange={this.onPageSizeChange}
            />
          </div>
        </TableWrapper>
      </Wrapper>
    )
  }
}
