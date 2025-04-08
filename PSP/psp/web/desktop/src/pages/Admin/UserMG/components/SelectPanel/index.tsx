import * as React from 'react'
import { observable, action, computed, reaction } from 'mobx'
import { observer, disposeOnUnmount } from 'mobx-react'
import { Spin, Input, Checkbox } from 'antd'
import { sysConfig } from '@/domain'
import { ExclamationCircleOutlined } from '@ant-design/icons'

import ApproveStatus from '../ApproveStatus'

import { PanelWrapper } from './style'

interface IProps {
  fetchItems: () => Promise<any>
  placeholderName: string
  selectedKeys: number[]
  updateSelectedKeys: (keys: number[], items?: any[]) => void
  disabledCondition?: (item) => boolean
  unVisibleCondition?: (item) => boolean
  updateList: (keys: number[]) => void
  updateInitialItems: (item) => boolean
}

@observer
export default class SelectPanel extends React.Component<IProps> {
  @observable items = []
  @observable keyword = ''
  @observable fetching = false
  @action
  updateItems = items => (this.items = items)
  @action
  updateFetching = flag => (this.fetching = flag)
  @action
  updateKeyword = keyword => (this.keyword = keyword)

  // @disposeOnUnmount
  // disposer1 = reaction(
  //   () => this.finalVisibleItems,
  //   () => {
  //     this.props.updateList(this.finalVisibleItems)
  //   }
  // )

  // @disposeOnUnmount
  // disposer2 = reaction(
  //   () => this.initialItems,
  //   () => {
  //     this.props.updateInitialItems(this.initialItems)
  //   }
  // )

  componentDidMount() {
    this.updateFetching(true)
    this.props
      ?.fetchItems()
      .then(this.updateItems)
      .finally(() => {
        this.updateFetching(false)
        let keys = new Set([...this.props.selectedKeys, ...this.initialKeys])
        let selectedItems = this.items.filter(i => keys.has(i.id))

        this.props.updateSelectedKeys([...keys], selectedItems)
      })
  }

  @computed
  get initialItems() {
    return this.props.disabledCondition
      ? this.items.filter(this.props.disabledCondition)
      : []
  }

  @computed
  get initialKeys() {
    return this.initialItems.map(i => i.id)
  }

  private onSelect = (checked, id) => {
    let keys = [...this.props.selectedKeys]

    if (checked) {
      keys = [...new Set([...keys, id])]
    } else {
      keys = keys.filter(key => key !== id)
    }

    let keysSet = new Set([...keys])
    let selectedItems = this.items.filter(i => keysSet.has(i.id))

    this.props.updateSelectedKeys([...keysSet], selectedItems)
  }

  private onChange = e => this.updateKeyword(e.target.value)

  @computed
  get visibleItems() {
    return this.items.filter(
      item =>
        !this.keyword ||
        item.name.toLowerCase().includes(this.keyword.toLowerCase())
    )
  }

  @computed
  get finalVisibleItems() {
    const filterItems = this.visibleItems.filter(item => item.type !== 1)
    return this.props.unVisibleCondition
      ? filterItems.filter(item => !this.props.unVisibleCondition(item))
      : filterItems
  }

  render() {
    const { fetching, keyword } = this
    const { selectedKeys } = this.props

    return (
      <PanelWrapper>
        {fetching ? (
          <Spin style={{ margin: 'auto' }} />
        ) : (
          <>
            {/* <div className='filter'>
              搜索：
              <Input.Search
                size='small'
                maxLength={64}
                placeholder={`请输入${this.props.placeholderName}名称`}
                value={keyword}
                onChange={this.onChange}
              />
            </div> */}
            <div className='itemList'>
              {this.finalVisibleItems.map(item => (
                <div className='item' key={item.id}>
                  <Checkbox
                    checked={selectedKeys.includes(item.id)}
                    disabled={
                      this.initialItems.includes(item) ||
                      (sysConfig.enableThreeMembers &&
                        item.approve_status === 0)
                    }
                    onChange={e => this.onSelect(e.target.checked, item.id)}
                  />
                  <span title={item.name}>{item.name}</span>
                  {sysConfig.enableThreeMembers && item.approve_status === 0 && (
                    <ApproveStatus
                      title='未审批信息'
                      data={item}
                      callback={() => {}}
                      targetType={'USER'}
                      targetId={item.id}>
                      <ExclamationCircleOutlined style={{ color: '#1D398B' }} />
                    </ApproveStatus>
                  )}
                </div>
              ))}
            </div>
          </>
        )}
      </PanelWrapper>
    )
  }
}
