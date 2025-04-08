import * as React from 'react'
import { observer } from 'mobx-react'
import { Menu, Dropdown, message } from 'antd'
import { Icon } from '@/components'

interface IProps {
  node: any
}

interface IState {
  visible: boolean
}

@observer
export default class StatusMenu extends React.Component<IProps, IState> {
  state = {
    visible: false,
  }
  handleClick = async (clickParam, action: string) => {
    clickParam.domEvent.stopPropagation()
    const { node } = this.props
    const res = await node.updateState(action)

    if (res && res.errorCode === 0) {
      message.success(`更新节点 ${node.name} 状态成功`)
    }
  }

  handleVisibleChange = (visible: boolean) => {
    this.setState({ visible })
  }

  render() {
    const menu = (
      <Menu>
        <Menu.Item
          disabled={this.props.node.state === 'free'}
          onClick={clickParam => this.handleClick(clickParam, 'free')}>
          free
        </Menu.Item>
        <Menu.Item
          disabled={this.props.node.state === 'offline'}
          onClick={clickParam => this.handleClick(clickParam, 'offline')}>
          offline
        </Menu.Item>
      </Menu>
    )

    return (
      <Dropdown overlay={menu} onVisibleChange={this.handleVisibleChange}>
        <a href='#' onClick={e => e.preventDefault()}>
          <Icon
            type={this.state.visible ? 'caret-up' : 'caret-down'}
            style={{ color: '#3E75CB', fontSize: 8 }}
          />
        </a>
      </Dropdown>
    )
  }
}
