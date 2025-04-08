import * as React from 'react'
import { observable, action, computed } from 'mobx'
import { observer } from 'mobx-react'
import remove from 'lodash/remove'
import { Checkbox } from 'antd'
import { ExclamationCircleOutlined } from '@ant-design/icons'
import { sysConfig } from '@/domain'

import { SelectPanel, RadiusItem } from '..'
import { SelectEditorWrapper } from './style'

interface ListDto {
  id: number
  name: string
}

enum SelectType {
  ALL,
  INDETERMINATE,
  NONE
}

interface IProps {
  selectedLeftKeys: number[]
  updateSelectedLeftKeys: (keys: number[], items?: any[]) => void
  leftDisabledCondition?: (item) => boolean
  rightDisabledCondition?: (item) => boolean
  leftUnVisibleCondition?: (item) => boolean
  rightUnVisibleCondition?: (item) => boolean
  LeftList: {
    fetch: () => Promise<any>
    list: ListDto[]
  }
  RightList?: {
    fetch: () => Promise<any>
    list: ListDto[]
  }
  title: {
    leftTab: string
  }
}

@observer
export default class SelectEditor extends React.Component<IProps> {
  @observable activeTab = 'left'
  @observable allLeftKeyList = []
  @observable allRightKeyList = []
  @observable leftInitialItems = []
  @observable rightInitialItems = []
  @action
  updateActiveTab = key => (this.activeTab = key)
  @action
  updateAllLeftKeyList = keys => (this.allLeftKeyList = keys)
  @action
  updateAllRightKeyList = keys => (this.allRightKeyList = keys)
  @action
  updateLeftInitialItems = items => (this.leftInitialItems = items)

  private tabs = [
    {
      title: this.props.title.leftTab,
      key: 'left'
    }
  ]

  private allHas(visibleKeys = [], allKeys = []) {
    const res = allKeys.filter(al => visibleKeys.includes(al))

    if (!res.length) return SelectType.NONE
    if (res.length === visibleKeys.length) return SelectType.ALL
    return SelectType.INDETERMINATE
  }

  @computed
  get leftInitialKeys() {
    return this.leftInitialItems.map(it => it.id)
  }

  @computed
  get rightInitialKeys() {
    return this.rightInitialItems.map(it => it.id)
  }

  @computed
  get rightAllKeys() {
    return this.allRightKeyList.map(k => k.id)
  }

  @computed
  get leftAllKeys() {
    return this.allLeftKeyList.map(k => k.id)
  }

  @computed
  get leftAllSelected() {
    const { selectedLeftKeys } = this.props
    return this.allHas(this.leftAllKeys, selectedLeftKeys) === SelectType.ALL
  }

  @computed
  get leftIndeterminate() {
    const { selectedLeftKeys } = this.props
    return (
      this.allHas(this.leftAllKeys, selectedLeftKeys) ===
      SelectType.INDETERMINATE
    )
  }

  private selectLeftAll = e => {
    const { selectedLeftKeys: originKeys } = this.props
    let { checked } = e.target

    if (
      originKeys.length >=
      this.allLeftKeyList.filter(k => k['approve_status'] !== 0).length
    ) {
      checked = false
    }

    const disabledSelectedKeys = this.allLeftKeyList
      .filter(k => k['approve_status'] === 0 && originKeys.includes(k.id))
      .map(k => k.id)

    const disabledUnSelectedKeys = this.allLeftKeyList
      .filter(k => k['approve_status'] === 0 && !originKeys.includes(k.id))
      .map(k => k.id)

    let keys = []
    if (checked) {
      keys = [...new Set([...this.leftAllKeys, ...originKeys])]

      if (sysConfig.enableThreeMembers) {
        keys = keys.filter(k => !disabledUnSelectedKeys.includes(k))
      }
    } else {
      remove(
        originKeys,
        iKey =>
          this.leftAllKeys.includes(iKey) &&
          !this.leftInitialKeys.includes(iKey)
      )
      keys = [...new Set([...originKeys])]

      if (sysConfig.enableThreeMembers) {
        keys = [...new Set([...disabledSelectedKeys, ...keys])]
      }
    }

    const selectedLeftList = this.props.LeftList.list.filter(item =>
      keys.includes(item.id)
    )

    this.props.updateSelectedLeftKeys(keys, selectedLeftList)
  }

  @computed
  get selectedLeftList() {
    const { LeftList, selectedLeftKeys } = this.props
    return LeftList.list.filter(item => selectedLeftKeys.includes(item.id))
  }

  render() {
    const {
      title: { leftTab },
      LeftList,
      selectedLeftKeys
    } = this.props

    return (
      <SelectEditorWrapper>
        <div className='editorBody'>
          <div className='left'>
            <div className='module'>
              <header>
                <span className='label'>{leftTab}</span>
              </header>
              <div className='body'>
                <RadiusItem
                  itemList={this.selectedLeftList.map(li => li.name)}
                />
              </div>
            </div>
          </div>
          <div className='right'>
            <div className='module'>
              <header>
                <span className='label'>
                  选择{leftTab}
                  {sysConfig.enableThreeMembers && (
                    <span style={{ fontSize: 12 }}>
                      「选项中的
                      <ExclamationCircleOutlined style={{ color: '#1D398B' }} />
                      表示有未审批信息」
                    </span>
                  )}
                </span>
              </header>
              <div className='body'>
                {/* <div className='header'>
                  {this.tabs.map(tab => (
                    <div
                      key={tab.key}
                      className={`tab ${
                        tab.key === this.activeTab ? 'active' : ''
                      }
                      ${tab.key === 'left' ? 'left_tab' : ''}
                      `}
                      onClick={() => this.updateActiveTab(tab.key)}>
                      <Checkbox
                        checked={this.leftAllSelected}
                        indeterminate={this.leftIndeterminate}
                        onClick={this.selectLeftAll}
                      />
                      <span className='title'>
                        {tab.title} 已选 {selectedLeftKeys?.length}
                      </span>
                    </div>
                  ))}
                </div> */}
                {this.tabs.map(tab => (
                  <div
                    key={tab.key}
                    className={`tabPanel ${
                      tab.key === this.activeTab ? 'active' : ''
                    }`}>
                    {tab.key === 'left' && (
                      <SelectPanel
                        fetchItems={LeftList.fetch}
                        placeholderName={tab.title}
                        selectedKeys={selectedLeftKeys}
                        updateSelectedKeys={this.props.updateSelectedLeftKeys}
                        disabledCondition={
                          this.props.leftDisabledCondition &&
                          this.props.leftDisabledCondition
                        }
                        unVisibleCondition={
                          this.props.leftUnVisibleCondition &&
                          this.props.leftUnVisibleCondition
                        }
                        updateList={this.updateAllLeftKeyList}
                        updateInitialItems={this.updateLeftInitialItems}
                      />
                    )}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </SelectEditorWrapper>
    )
  }
}
