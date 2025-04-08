import * as React from 'react'
import { observable, action, computed } from 'mobx'
import { observer } from 'mobx-react'
import remove from 'lodash/remove'
import { Checkbox } from 'antd'

import { SelectPanel, RadiusItem } from '../../components'
import { SelectEditorWrapper } from '../../components/SelectEditor/style'
import { RoleList } from '@/domain/UserMG'

enum SelectType {
  ALL,
  INDETERMINATE,
  NONE,
}

interface IProps {
  selectedKeys: number[]
  updateSelectedKeys: (keys: number[], items?: any[]) => void
  disabledCondition?: (item) => boolean
  title: string
}

@observer
export default class SelectEditor extends React.Component<IProps> {
  @observable allKeyList = []
  @observable initialItems = []
  @action
  updateAllKeyList = keys => (this.allKeyList = keys)
  @action
  updateInitialItems = items => (this.initialItems = items)

  private allHas(visibleKeys = [], allKeys = []) {
    const res = allKeys.filter(al => visibleKeys.includes(al))

    if (!res.length) return SelectType.NONE
    if (res.length === visibleKeys.length) return SelectType.ALL
    return SelectType.INDETERMINATE
  }

  @computed
  get initialKeys() {
    return this.initialItems.map(it => it.id)
  }

  @computed
  get allKeys() {
    return this.allKeyList.map(k => k.id)
  }

  @computed
  get allSelected() {
    const { selectedKeys } = this.props
    return this.allHas(this.allKeys, selectedKeys) === SelectType.ALL
  }

  @computed
  get indeterminate() {
    const { selectedKeys } = this.props
    return this.allHas(this.allKeys, selectedKeys) === SelectType.INDETERMINATE
  }

  private selectAll = e => {
    const { selectedKeys: originKeys } = this.props
    const { checked } = e.target

    let keys = []
    if (checked) {
      keys = [...new Set([...this.allKeys, ...originKeys])]
    } else {
      remove(
        originKeys,
        iKey => this.allKeys.includes(iKey) && !this.initialKeys.includes(iKey)
      )
      keys = [...new Set([...originKeys])]
    }
    this.props.updateSelectedKeys(keys, this.selectedList)
  }

  @computed
  get selectedList() {
    const { selectedKeys } = this.props
    return RoleList.roleList?.filter(item => selectedKeys.includes(item.id))
  }

  render() {
    const { title, selectedKeys } = this.props

    return (
      <SelectEditorWrapper isSync={true}>
        <div className='editorBody'>
          <div className='left'>
            <div className='module'>
              <header>
                <span className='label'>{title}</span>
              </header>
              <div className='body'>
                <RadiusItem itemList={this.selectedList?.map(i => i.name)} />
              </div>
            </div>
          </div>
          <div className='right'>
            <div className='module'>
              <header>
                <span className='label'>选择{title}</span>
              </header>
              <div className='body'>
                <div className='header'>
                  <div className={'tab active'}>
                    <Checkbox
                      checked={this.allSelected}
                      indeterminate={this.indeterminate}
                      onClick={this.selectAll}
                    />
                    <span className='title'>
                      {title} 已选 {selectedKeys.length}
                    </span>
                  </div>
                </div>
                <div className={'tabPanel active'}>
                  <SelectPanel
                    fetchItems={() =>
                      RoleList.fetch().then(res => res.data.roles)
                    }
                    placeholderName={title}
                    selectedKeys={selectedKeys}
                    updateSelectedKeys={this.props.updateSelectedKeys}
                    disabledCondition={
                      this.props.disabledCondition &&
                      this.props.disabledCondition
                    }
                    updateList={this.updateAllKeyList}
                    updateInitialItems={this.updateInitialItems}
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </SelectEditorWrapper>
    )
  }
}
