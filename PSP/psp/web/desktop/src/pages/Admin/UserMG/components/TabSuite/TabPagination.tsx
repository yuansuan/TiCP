import * as React from 'react'
import { Pagination } from 'antd'
import { observer } from 'mobx-react'

import { ListQuery } from '../../utils'

interface IProps {
  updateListQuery: (newQuery: ListQuery) => void
  listQuery: ListQuery
  total: number
}

@observer
export default class TabPagination extends React.Component<IProps> {
  onShowSizeChange = (_, pageSize) => {
    const { listQuery: newQuery, updateListQuery } = this.props
    if (newQuery.pageSize !== pageSize) {
      newQuery.pageSize = pageSize
      newQuery.page = 1
      updateListQuery(newQuery)
    }
  }

  changePagination = (page, pageSize) => {
    const { listQuery: newQuery, updateListQuery } = this.props
    if (newQuery.page !== page) {
      newQuery.page = page
      updateListQuery(newQuery)
    }
    if (newQuery.pageSize !== pageSize) {
      newQuery.pageSize = pageSize
      newQuery.page = 1
      updateListQuery(newQuery)
    }
  }

  public render() {
    const {
      total,
      listQuery: { page },
    } = this.props
    return (
      <Pagination
        current={page}
        total={total}
        onChange={this.changePagination}
        onShowSizeChange={this.onShowSizeChange}
        showSizeChanger
      />
    )
  }
}
