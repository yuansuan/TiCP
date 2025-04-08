/**
 * @module RoleChooser
 * multi choose roles
 */
import * as React from 'react'
import { observable, action, computed, reaction } from 'mobx'
import { observer, disposeOnUnmount } from 'mobx-react'
import { Input, Spin, Checkbox } from 'antd'
import remove from 'lodash/remove'

import { createMobxStream, fromStream } from '@/utils'
import { RoleChooserWrapper } from './style'

interface IProps {
  className?: string
  title: string
  fetchItems: () => Promise<any>
  selectedKeys: { current: string[] }
  disabledCondition: (item) => {}
}

enum SelectType {
  ALL,
  INDETERMINATE,
  NONE
}

@observer
export default class ItemChooser extends React.Component<IProps> {
  @observable selectedKeys = []
  @observable items = []
  @observable name = ''
  @observable keyword = ''
  @observable fetching = false
  @action
  updateSelectedKeys = keys => (this.selectedKeys = keys)
  @action
  updateName = name => (this.name = name)
  @action
  updateKeyword = keyword => (this.keyword = keyword)
  @action
  updateItems = items => (this.items = items)
  @action
  updateFetching = flag => (this.fetching = flag)

  constructor(props) {
    super(props)

    this.updateSelectedKeys([
      ...props.selectedKeys.current,
      ...this.initialKeys
    ])
  }

  @disposeOnUnmount
  disposer1 = reaction(
    () => this.initialKeys,
    () => {
      this.updateSelectedKeys(this.initialKeys)
    }
  )

  componentDidMount() {
    this.fetchItems()

    // export selected keys
    fromStream(
      createMobxStream(() => this.selectedKeys),
      this.props.selectedKeys
    )
  }

  private allHas(visibleKeys = [], allKeys = []) {
    const res = allKeys.filter(al => visibleKeys.includes(al))

    if (!res.length) return SelectType.NONE
    if (res.length === visibleKeys.length) return SelectType.ALL
    return SelectType.INDETERMINATE
  }

  @computed
  get selectedItems() {
    return this.items.filter(item => this.selectedKeys.includes(item.id))
  }

  @computed
  get initialItems() {
    return this.items.filter(item => this.props.disabledCondition(item))
  }

  @computed
  get initialKeys() {
    return this.initialItems.map(it => it.id)
  }

  @computed
  get visibleRoles() {
    return this.items.filter(
      item => !this.keyword || item.name.includes(this.keyword)
    )
  }

  @computed
  get visibleRoleIds() {
    return this.visibleRoles.map(r => r.id)
  }

  @computed
  get indeterminate() {
    return (
      this.allHas(this.visibleRoleIds, this.selectedKeys) ===
      SelectType.INDETERMINATE
    )
  }

  @computed
  get allSelected() {
    return (
      this.allHas(this.visibleRoleIds, this.selectedKeys) === SelectType.ALL
    )
  }

  private fetchItems = () => {
    this.updateFetching(true)
    this.props
      .fetchItems()
      .then(this.updateItems)
      .finally(() => this.updateFetching(false))
  }

  private onChangeKeyword = e => this.updateKeyword(e.target.value)

  private onSelectAll = e => {
    const { checked } = e.target
    let oldKeys = [...this.selectedKeys]
    let keys

    // select all
    if (checked) {
      keys = [...new Set([...this.visibleRoleIds, ...oldKeys])]
    } else {
      remove(
        oldKeys,
        iKey =>
          this.visibleRoleIds.includes(iKey) && !this.initialKeys.includes(iKey)
      )

      keys = [...new Set([...oldKeys])]
    }

    this.updateSelectedKeys(keys)
  }

  private onSelect = (checked, id) => {
    let keys = [...this.selectedKeys]

    if (checked) {
      keys = [...new Set([...keys, id])]
    } else {
      keys = keys.filter(key => key !== id)
    }

    this.updateSelectedKeys(keys)
  }

  render() {
    const {
      selectedItems,
      visibleRoles,
      selectedKeys,
      allSelected,
      indeterminate,
      keyword,
      fetching
    } = this
    const { className = '', title = '' } = this.props

    return (
      <RoleChooserWrapper className={className}>
        <div className='module left'>
          <div className='header'>
            <span className='label'>添加新{title}</span>
          </div>

          <div className='body'>
            <div className='itemList'>
              {selectedItems.map(item => (
                <div className='item' key={item.id}>
                  {item.name}
                </div>
              ))}
            </div>
          </div>
        </div>
        <div className='module right'>
          <div className='header'>
            <span className='label'>选择{title}</span>
            <Input.Search
              className='filter'
              maxLength={64}
              value={keyword}
              onChange={this.onChangeKeyword}
            />
          </div>
          <div className='body'>
            {fetching ? (
              <Spin style={{ margin: 'auto' }} />
            ) : (
              <>
                <div className='all'>
                  <Checkbox
                    checked={allSelected}
                    indeterminate={indeterminate}
                    onChange={this.onSelectAll}
                  />
                  <span className='title'>全部{title}</span>
                </div>
                <div className='itemList'>
                  {visibleRoles.map(item => (
                    <div className='item' key={item.id}>
                      <Checkbox
                        disabled={this.initialKeys.includes(item.id)}
                        checked={selectedKeys.includes(item.id)}
                        onChange={e => this.onSelect(e.target.checked, item.id)}
                      />
                      <span title={item.name}>{item.name}</span>
                    </div>
                  ))}
                </div>
              </>
            )}
          </div>
        </div>
      </RoleChooserWrapper>
    )
  }
}
