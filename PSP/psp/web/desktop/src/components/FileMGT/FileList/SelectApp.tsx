/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Tooltip } from 'antd'
import { Icon, Modal } from '@/components'
import { cloudAppList, env } from '@/domain'
import { useObserver, useLocalStore } from 'mobx-react-lite'
import { BaseDirectory, BaseFile } from '@/utils/FileSystem'
import { openVisualApp } from '@/utils'
import { getToken, VisualHttp } from '@/domain/VisualHttp'
import CloudApp from '@/domain/Visualization/CloudApp'

type Node = BaseDirectory | BaseFile
type Props = {
  node: Node
}

const Styled = styled.div`
  display: inline-block;
  font-size: 14px;
  cursor: pointer;
  margin-left: 8px;

  .anticon {
    vertical-align: middle;
  }
`
const AppContainer = styled.div`
  padding: 25px;
  display: flex;
  flex-wrap: wrap;
  .app-item {
    margin: 10px;
    cursor: pointer;
    width: 90px;
    .img-box {
      width: 90px;
      height: 90px;
      border: 1px solid rgba(0, 0, 0, 0.15);
      border-radius: 2px;
      display: flex;
      justify-content: center;
      align-items: center;
      position: relative;
      &:hover {
        border-color: #3182ff;
      }
      img {
        height: 70px;
        width: 70px;
      }
    }

    &.selected {
      .img-box {
        border-color: #3182ff;
      }
      .check-container {
        display: flex;
      }
    }
    .title {
      margin-top: 5px;
      text-align: center;
    }
  }
  .check-container {
    display: none;
    top: 5px;
    right: 5px;
    height: 16px;
    width: 16px;
    position: absolute;
    background-color: #3182ff;
    border-radius: 8px;
    justify-content: center;
    align-items: center;
  }
  .check {
    display: inline-block;
    transform: rotate(45deg);
    height: 8px;
    width: 4px;
    border-bottom: 1px solid white;
    border-right: 1px solid white;
  }
`

export const SelectApp = observer(function SelecctApp({ node }: Props) {
  const state = useLocalStore(() => ({
    app: null,
    setApp(app) {
      this.app = app
    },
  }))
  const openVirtualApp = async (app: CloudApp, filePath: string) => {
    VisualHttp.post(
      '/worktask',
      {
        app_id: +app.id,
        template_name: app.name,
        project_id: env.project?.id,
        from: 'public_cloud',
        app_param: app.app_param,
        app_param_paths: [filePath],
      },
      { baseURL: '' }
    ).then(res => {
      const { link, user_id, id } = res.data
      const url = `/visualization/?link=${link}&worktask_id=${id}&user_id=${user_id}&access_token=${getToken()}`
      openVisualApp(url)
    })
  }
  async function open() {
    await cloudAppList.fetch()
    if (!cloudAppList.list.length) {
      return
    }
    state.setApp(null)
    const Content = function () {
      function select(app) {
        state.setApp(app)
      }
      const Apps = cloudAppList.list.map(function CloudApp(item) {
        return useObserver(() => {
          const selected = state.app?.id === item.id
          return (
            <div
              className={`app-item ${selected ? 'selected' : ''}`}
              key={item.id}
              onClick={() => select(item)}>
              <div className='img-box'>
                <div className='check-container'>
                  <div className='check'></div>
                </div>
                <img src={item.icon_data} />
              </div>
              <div className='title'>{item.name}</div>
            </div>
          )
        })
      })
      return <AppContainer>{Apps}</AppContainer>
    }
    await Modal.show({
      title: '打开方式',
      width: 600,
      className: '__open_vis__',
      style: {
        padding: 0,
      },
      bodyStyle: {
        padding: 0,
      },
      content: ({ onCancel, onOk }) => <Content />,
    })
    if (state.app) {
      openVirtualApp(state.app, node.path)
    }
  }

  return (
    <Styled>
      <Tooltip title='打开文件' className='btn'>
        <Icon type='copy' onClick={open} />
      </Tooltip>
    </Styled>
  )
})
