import * as React from 'react'
import nanoid from 'nanoid'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import { observable, action } from 'mobx'
import { Spin, Tabs } from 'antd'
import EditableForm from '../../../Editor/EditableForm'
import { App } from '@/domain/Applications'
import GlobalContext from '@/pages/context'
import Script from './Script'
import TemplateInfo from './BaseInfo'
import Document from './Document'

const Wrapper = styled.div`
  height: 100%;
  padding:0 10px 0 10px;
  .ant-tabs-bar {
    background-color: white;
    margin: 0;
    border-bottom: none;
  }

  .tabLayout {
    height: calc(100vh - 280px);
    overflow: auto;
  }
`

const Loading = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  right: 0;
`

const { TabPane } = Tabs

interface IProps {
  app: App
  remoteAppList?: any
  context?: any
}

@observer
export default class Editor extends React.Component<IProps> {
  uploadToken = nanoid()
  @observable loading = false
  @observable app
  @observable private tabs = []
  @action
  updateLoading = loading => (this.loading = loading)

  formModel = observable({})

  editorRef = null

  async componentDidMount() {
    const { app } = this.props

    // render Script before fetchScript to correct height
    this.updateLoading(true)
    await app.fetchScript().finally(() => this.updateLoading(false))

    this.tabs = [
      {
        title: '模版信息',
        key: 'INFO',
        icon: 'tag',
        component: (
          <TemplateInfo
            app={app}
            isRemote={false}
          />
        )
      },
      {
        title: '表单',
        key: 'FORM',
        icon: 'form',
        component: <EditableForm formModel={this.formModel} showMenu={false} />
      },
      {
        title: '脚本',
        key: 'SCRIPT',
        icon: 'script',
        component: <Script ref={ref => (this.editorRef = ref)} />
      },
      {
        title: '说明文档',
        key: 'DOC',
        icon: 'readme',
        component: <Document helpDoc={app.helpDoc} />
      }
    ]
  }

  render() {
    const { app } = this.props

    return (
      <GlobalContext.Provider value={{ app }}>
        <Wrapper>
          <>
            <Tabs defaultActiveKey='FORM' animated={false}>
              {this.tabs.map(tab => (
                <TabPane key={tab.key} tab={tab.title}>
                  <div className='tabLayout'>{tab.component}</div>
                </TabPane>
              ))}
            </Tabs>
          </>

          {this.loading ? (
            <Loading>
              <Spin />
            </Loading>
          ) : null}
        </Wrapper>
      </GlobalContext.Provider>
    )
  }
}
