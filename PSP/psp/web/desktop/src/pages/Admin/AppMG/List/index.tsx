import * as React from 'react'
import styled from 'styled-components'
import { Tabs, Dropdown, Menu } from 'antd'
import { observable } from 'mobx'
import { sysConfig } from '@/domain'
import { ArrowLeftOutlined, CaretDownOutlined } from '@ant-design/icons'
import Normal from './Normal'
import { RemoteAppList, AppList } from '@/domain/Applications'
import { LabelMap } from './ThreeDMgr'
import { INSTALL_TYPE } from '@/utils/const'
import { getUrlParams } from '@/utils'
const StyledWrapper = styled.div`
  width: 100%;
  height: 100%;
  background-color: white;

  .body {
    margin: 0 20px;

    .ant-tabs-bar {
      border-bottom: none;
    }

    .tabLayout {
      height: calc(100vh - 210px);
    }
  }
`

type Tab = {
  title: string
  key: string
  Component: any
}

const { TabPane } = Tabs

export default class ApplicationList extends React.Component {
  @observable appList = new AppList()
  @observable remoteAppList = new RemoteAppList()
  state = {
    threeDMgrType: 'WorkStation',
    defaultActiveKey: 'normal'
  }

  get isAIO() {
    return sysConfig.installType === INSTALL_TYPE.aio
  }

  get tabs() {
      return [
        { title: '本地应用模版', key: 'normal', Component: Normal }
      ] as Array<Tab>
  }

  componentDidMount() {
    const query = getUrlParams(window.location.href)
    if (query.tab) {
      this.setState({
        defaultActiveKey: query.tab
      })
    }
  }
  handleDropdownMenuClick = ({ key }) => {
    this.setState({
      threeDMgrType: key
    })
  }
  onTabClick = key => {
    this.setState({ defaultActiveKey: key })
    if (key === 'remote') {
      this.remoteAppList.fetchTemplates()
    } else if (key === 'normal') {
      this.appList.fetchTemplates()
    }
  }

  renderDropdownMenu = text => {
    const menu = (
      <Menu onClick={this.handleDropdownMenuClick}>
        {Object.keys(LabelMap).map(key => (
          <Menu.Item key={key}>
            {key === this.state.threeDMgrType ? (
              <div style={{ color: '#194E8B' }}>
                {LabelMap[key]} <ArrowLeftOutlined />
              </div>
            ) : (
              LabelMap[key]
            )}
          </Menu.Item>
        ))}
      </Menu>
    )

    return (
      <>
        {text}
        <Dropdown overlay={menu}>
          <CaretDownOutlined />
        </Dropdown>
      </>
    )
  }

  render() {
    return (
      <StyledWrapper style={{ height: '100%' }}>
        <div className='body'>
          <Tabs
            defaultActiveKey={this.state.defaultActiveKey}
            activeKey={this.state.defaultActiveKey}
            onTabClick={this.onTabClick}
            animated={false}>
            {this.tabs.map(({ key, title, Component }) => (
              <TabPane
                key={key}
                tab={
                  key === 'visualmgr' ? this.renderDropdownMenu(title) : title
                }>
                <div className='tabLayout'>
                  {key === 'visualmgr' ? (
                    <Component type={this.state.threeDMgrType} />
                  ) : (
                    <Component />
                  )}
                </div>
              </TabPane>
            ))}
          </Tabs>
        </div>
      </StyledWrapper>
    )
  }
}
