import * as React from 'react'
import { ToolbarWrapper } from './style'
import { Search } from '@/components'
import { ListQuery } from '@/pages/Admin/UserMG/utils'

interface IProps {
  operators?: React.ReactNode
  updateListQuery: (listQuery: ListQuery) => void
  listQuery: ListQuery
  placeholder?: string
}

export default class Toolbar extends React.Component<IProps> {
  onPressSearch = (value: string) => {
    const { listQuery: newQuery, updateListQuery } = this.props
    newQuery.query = value
    newQuery.page = 1
    updateListQuery(newQuery)
  }

  public render() {
    const { operators, placeholder } = this.props

    return (
      <ToolbarWrapper>
        {operators}

        <div className='filter'>
          <Search
            onSearch={this.onPressSearch}
            debounceWait={300}
            placeholder={placeholder}
          />
        </div>
      </ToolbarWrapper>
    )
  }
}
