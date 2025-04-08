import * as React from 'react'
import { inject, observer } from 'mobx-react'
import { observable, computed } from 'mobx'
import { Input, Select, Tooltip } from 'antd'
import debounce from 'lodash.debounce'
import sysConfig from '@/domain/SysConfig'

import { Job } from '@/domain/JobList/Job'
// import {
//   JobStatusFilterList,
//   jobStatusColumnFields,
//   HistoryJobStatusFilterList,
// } from '@/domain/JobList'
// import { JobActionButtonGroup, Search } from '@/components'

import {
  Wrapper,
  SearchWrapper,
  ActionWrapper,
  MoreFilterWrapper
} from './style'

const { Option } = Select
const genFilterItemTip = (v1, v2, mapping = v => v) =>
  v1 == v2 ? '无' : mapping(v1)

const daysMap = {
  7: '过去7天之内',
  30: '过去30天之内',
  365: '过去365天之内'
}

interface IButtonsOption {
  reSubmit?: boolean
}
interface ListActionsProps {
  isHistory?: boolean
  isOneBurst?: boolean
  hasUserFilter?: boolean
  isDropDown?: boolean
  queueList?: string[]
  selectedRowKeys?: string[]
  updateSelectedRowKeys?: (keys: string[]) => void
  jobList: Job[]
  filterByType?: (type, value, currentPage: number) => void
  hasActions: {
    buttons: boolean
    buttonsOption?: IButtonsOption
    search: boolean
    moreFilter: boolean
  }
  updateJobList: () => void
}

@inject((stores: any) => {
  const { store } = stores
  const { selectedRowKeys, updateSelectedRowKeys, filterByType } = store
  return {
    selectedRowKeys,
    updateSelectedRowKeys,
    filterByType
  }
})
@observer
export default class ListActions extends React.Component<ListActionsProps> {
  @observable moreFilter = {
    stateFilter: sysConfig.jobConfig?.job?.list_default_filter?.states || [], // 筛选出的状态
    userFilter: '', // 过滤用户
    appFilter: '', // 过滤应用
    queueFilter: '', // 过滤队列
    submitPastTime: null, // 按作业提交时间过滤
    endPastTime: null // 按作业结束时间过滤
  }

  @computed
  get filterTips() {
    if (
      this.moreFilter.stateFilter.length === 0 &&
      this.moreFilter.userFilter == '' &&
      this.moreFilter.appFilter == '' &&
      this.moreFilter.queueFilter == '' &&
      this.moreFilter.submitPastTime === null &&
      this.moreFilter.endPastTime === null
    ) {
      return '无过滤'
    } else {
      return `${
        this.props.hasUserFilter
          ? `用户: ${genFilterItemTip(this.moreFilter.userFilter, '')},`
          : ''
      } 
应用: ${genFilterItemTip(this.moreFilter.appFilter, '')}, 
队列: ${genFilterItemTip(this.moreFilter.queueFilter, '')},
作业状态: ${
        this.moreFilter.stateFilter.length === 0
          ? '无'
          : this.moreFilter.stateFilter
              .map(s => {
                return jobStatusColumnFields.filter(o => o.status === s)[0]
                  .label
              })
              .join(',')
      },
作业提交时间: ${genFilterItemTip(
        this.moreFilter.submitPastTime,
        null,
        v => daysMap[v]
      )},
作业结束时间: ${genFilterItemTip(
        this.moreFilter.endPastTime,
        null,
        v => daysMap[v]
      )}`
    }
  }

  @observable showMoreFilter = true

  debounceFilter = debounce((type, value) => {
    if (type === 'submitPastTime' || type === 'endPastTime') {
      this.props.filterByType(
        type,
        value ? Math.floor(Date.now() / 1000) - value * 24 * 60 * 60 : -1,
        1
      )
    } else {
      this.props.filterByType(type, value, 1)
    }

    // 清除 checkbox 状态
    this.props.updateSelectedRowKeys([])
  }, 300)

  onMoreFilterChange = (value, type) => {
    this.moreFilter[type] = value
    this.debounceFilter(type, value)
  }

  get selectedItems() {
    const { jobList, selectedRowKeys } = this.props
    return selectedRowKeys
      ? jobList.filter(item => selectedRowKeys.includes(item.id))
      : []
  }

  handleSearch = (value: string) => {
    this.props.filterByType('fuzzy', value, 1)
  }

  render() {
    const { hasActions, isHistory, isDropDown, isOneBurst } = this.props

    const jobStatusList = isHistory
      ? HistoryJobStatusFilterList
      : JobStatusFilterList

    return (
      <Wrapper>
        <ActionWrapper>
          {hasActions.buttons ? (
            <JobActionButtonGroup
              buttonsOption={hasActions.buttonsOption}
              isHistory={isHistory}
              isOneBurst={isOneBurst}
              isDropDown={isDropDown}
              selectedItems={this.selectedItems}
              operateCallback={() => {
                this.props.updateJobList()
                // 清除 checkbox 状态
                this.props.updateSelectedRowKeys([])
              }}
            />
          ) : (
            <span />
          )}
          {hasActions.search ? (
            <SearchWrapper>
              <Search
                style={{ width: 260, height: 32 }}
                onSearch={this.handleSearch}
                debounceWait={500}
                placeholder='按作业ID或作业名称搜索'
              />
              {hasActions.moreFilter && (
                <Tooltip
                  placement='topLeft'
                  title={`更多过滤项: ${this.filterTips}`}>
                  <Icon
                    style={{
                      margin: 10,
                      fontSize: 20,
                      color: this.filterTips !== '无过滤' ? '#1A6DBA' : 'inhert'
                    }}
                    type={this.showMoreFilter ? 'down-square' : 'up-square'}
                    onClick={e => (this.showMoreFilter = !this.showMoreFilter)}
                  />
                </Tooltip>
              )}
            </SearchWrapper>
          ) : (
            <span />
          )}
        </ActionWrapper>
        {hasActions.search && this.showMoreFilter && (
          <MoreFilterWrapper>
            {this.props.hasUserFilter && (
              <label className='item'>
                <div className='label'>用户:</div>
                <Input
                  style={{ width: 180 }}
                  placeholder='请输入提交作业用户'
                  value={this.moreFilter.userFilter}
                  onChange={e =>
                    this.onMoreFilterChange(e.target.value, 'userFilter')
                  }
                />
              </label>
            )}
            <label className='item'>
              <div className='label'>应用:</div>
              <Input
                style={{ width: 180 }}
                placeholder='请输入应用名称'
                value={this.moreFilter.appFilter}
                onChange={e =>
                  this.onMoreFilterChange(e.target.value, 'appFilter')
                }
              />
            </label>
            <label className='item'>
              <div className='label'> 队列:</div>
              <Select
                style={{ width: 180 }}
                placeholder='请选择队列'
                value={this.moreFilter.queueFilter}
                onChange={value =>
                  this.onMoreFilterChange(value, 'queueFilter')
                }>
                <Option key={''} value={''}>
                  请选择
                </Option>
                {this.props.queueList &&
                  this.props.queueList.map(q => (
                    <Option key={q} value={q}>
                      {q}
                    </Option>
                  ))}
              </Select>
            </label>
            <label className='item'>
              <div className='label'>作业状态:</div>
              <Select
                style={{ width: 180 }}
                mode='multiple'
                placeholder='请选择作业状态'
                value={this.moreFilter.stateFilter}
                onChange={value =>
                  this.onMoreFilterChange(value, 'stateFilter')
                }>
                {jobStatusList.map(s => (
                  <Option key={s} value={s}>
                    {
                      jobStatusColumnFields.filter(o => o.status === s)[0]
                        ?.label
                    }
                  </Option>
                ))}
              </Select>
            </label>
            <label className='item'>
              <div className='label'>作业提交时间:</div>
              <Select
                style={{ width: 180 }}
                value={this.moreFilter.submitPastTime}
                onChange={value =>
                  this.onMoreFilterChange(value, 'submitPastTime')
                }>
                <Option key={'-1'} value={null}>
                  请选择
                </Option>
                <Option key={'7d'} value={7}>
                  过去7天内
                </Option>
                <Option key={'30d'} value={30}>
                  过去30天内
                </Option>
                <Option key={'365d'} value={365}>
                  过去365天内
                </Option>
              </Select>
            </label>
            <label className='item'>
              <div className='label'>作业结束时间:</div>
              <Select
                style={{ width: 180 }}
                value={this.moreFilter.endPastTime}
                onChange={value =>
                  this.onMoreFilterChange(value, 'endPastTime')
                }>
                <Option key={'-1'} value={null}>
                  请选择
                </Option>
                <Option key={'7d'} value={7}>
                  过去7天内
                </Option>
                <Option key={'30d'} value={30}>
                  过去30天内
                </Option>
                <Option key={'365d'} value={365}>
                  过去365天内
                </Option>
              </Select>
            </label>
          </MoreFilterWrapper>
        )}
      </Wrapper>
    )
  }
}
