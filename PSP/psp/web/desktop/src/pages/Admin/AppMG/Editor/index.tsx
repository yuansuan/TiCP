/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import nanoid from 'nanoid'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import { observable, action } from 'mobx'
import { Spin, Tabs, Empty } from 'antd'
import { LeftOutlined } from '@ant-design/icons'
import { App } from '@/domain/Applications'
import GlobalContext from '@/pages/context'
import { history, Http, getUrlParams } from '@/utils'
import EditableForm from './EditableForm'
import Script from './Script'
import TemplateInfo from './BaseInfo'
import Footer from './Footer'
import Document from './Document'
import { sysConfig } from '@/domain'
import { DeployMode } from '@/constant'

const Wrapper = styled.div`
  height: 100%;
  background-color: white;
  padding-left: 15px;
  .link {
    cursor: pointer;
    &:hover {
      color: ${props => props.theme.primaryHighlightColor};
    }
  }

  .arrow {
    font-size: 15px;
    position: relative;
    top: 1px;
    right: 5px;
  }

  .ant-tabs-bar {
    margin: 0;
  }

  .tabLayout {
    height: calc(100vh - 280px);
  }
`

const EmptyContainer = styled.div`
  padding-bottom: 44px;
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
  app: any
  context?: any
}

@observer
export default class Editor extends React.Component<IProps> {
  // uploadToken is used to identify unque job submission
  uploadToken = nanoid()
  // The createDirPromise is singleton of request for uploadPath
  // It's used to avoid create directory repeatly when upload multiple files
  createDirPromise = null
  uploadPath = ''

  @observable loading = false
  @observable fetching = false
  @observable app
  @observable version = ''
  @observable private tabs = []
  @observable isRemote = false
  @action
  updateLoading = loading => (this.loading = loading)
  @action
  updateFetching = flag => (this.fetching = flag)

  formModel = observable({})

  editorRef = null

  async componentDidMount() {
    const query = getUrlParams()

    if (!query || !query.app) {
      return
    }
    try {
      this.isRemote = query.isRemote === 'true'
      this.updateFetching(true)

      this.app = await App.fetch({
        name: query.app,
        state: 'unpublished',
        version: query.version,
        compute_type: query.isRemote === 'true' ? 'cloud' : 'local'
      })
    } finally {
      this.updateFetching(false)
    }

    // render Script before fetchScript to correct height
    this.updateLoading(true)
    this.app.fetchScript().finally(() => this.updateLoading(false))
    this.tabs = [
      {
        title: '模版信息',
        key: 'INFO',
        icon: 'tag',
        component: <TemplateInfo isRemote={this.isRemote} />
      },
      {
        title: '表单',
        key: 'FORM',
        icon: 'form',
        component: <EditableForm formModel={this.formModel} />
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
        component: <Document helpDoc={this.app.helpDoc} />
      }
    ]
  }

  goBack = () => {
    const tab = this.isRemote ? 'remote' : 'normal'
    history.push(`/sys/template?tab=${tab}`)
  }

  // create upload directory
  fetchUploadPath = () => {
    // return the cache of uploadPath
    if (this.uploadPath) {
      return Promise.resolve(this.uploadPath)
    }

    // singleton promise
    if (this.createDirPromise) {
      return this.createDirPromise
    }
    this.createDirPromise = Http.post('/application/create_dir', {
      upload_sub_token: this.uploadToken
    }).then(({ data: { job_file_path } }) => {
      this.uploadPath = job_file_path
      return this.uploadPath
    })
    return this.createDirPromise
  }

  render() {
    const { app, uploadToken, fetchUploadPath, fetching } = this
    const appName = app && app.name ? app.name : ''

    return (
      <GlobalContext.Provider value={{ app, uploadToken, fetchUploadPath }}>
        <Wrapper>
          <div className='link' onClick={this.goBack}>
            <span className='arrow'>
              <LeftOutlined />
            </span>
            编辑模版{appName && `（${appName}）`}
          </div>

          {this.tabs.length > 0 ? (
            <>
              <Tabs defaultActiveKey='FORM' animated={false}>
                {this.tabs.map(tab => (
                  <TabPane key={tab.key} tab={tab.title}>
                    <div className='tabLayout'>{tab.component}</div>
                  </TabPane>
                ))}
              </Tabs>
              <Footer
                formModel={this.formModel}
                disabled={this.loading}
                isRemote={!!this.isRemote}
                updateDisabled={this.updateLoading}
                getScriptData={() => {
                  if (this.editorRef) {
                    return this.editorRef.getValue()
                  } else {
                    return app.scriptData
                  }
                }}
              />
            </>
          ) : (
            <EmptyContainer>
              <Empty
                description={
                  fetching ? '模版查询中...' : '未找到匹配的应用模版'
                }
              />
            </EmptyContainer>
          )}

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
