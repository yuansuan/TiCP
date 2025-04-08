import { Pagination } from 'antd'
import { inject, observer } from 'mobx-react'
import * as React from 'react'
import { Wrapper } from './style'

interface ListPaginationProps {
  currentIndex?: number
  totalItems: number
  pageSize?: number
  updateCurrentIndex?: (currentIndex: number) => void
  updatePageSize?: (pageSize: number) => void
  updateSelectedRowKeys?: (keys: string[]) => void
}

@inject((stores: any) => {
  const {
    currentIndex,
    pageSize,
    updateCurrentIndex,
    updatePageSize,
    updateSelectedRowKeys,
  } = stores.store
  return {
    currentIndex,
    pageSize,
    updateCurrentIndex,
    updatePageSize,
    updateSelectedRowKeys,
  }
})
@observer
export default class ListPagination extends React.Component<
  ListPaginationProps
> {
  updatePageSize = (current: number, size: number) => {
    this.props.updateSelectedRowKeys([])
    this.props.updatePageSize(size)
    this.props.updateCurrentIndex(1)
  }

  updateCurrentIndex = (current: number) => {
    this.props.updateSelectedRowKeys([])
    this.props.updateCurrentIndex(current)
  }

  render() {
    const { currentIndex, pageSize, totalItems } = this.props
    return (
      <Wrapper>
        <Pagination
          showQuickJumper
          showSizeChanger
          pageSize={pageSize}
          current={currentIndex}
          total={totalItems}
          onChange={this.updateCurrentIndex}
          onShowSizeChange={this.updatePageSize}
        />
      </Wrapper>
    )
  }
}
