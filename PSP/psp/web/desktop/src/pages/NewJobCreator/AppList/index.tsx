/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Section } from '../components'
import { Software } from './Software'
import { Input, Empty } from 'antd'
import { appList } from '@/domain'
import { useStore } from '../store'

const upFlod = require('@/assets/images/up-flod.svg')
const downFlod = require('@/assets/images/down-flod.svg')

const StyledToolbar = styled.div`
  display: flex;
  align-items: center;

  > .appSearch {
    display: flex;
    align-items: center;
  }

  > .unfold,
  > .fold {
    cursor: pointer;
    margin-left: 20px;
    display: flex;
    align-items: center;
    user-select: none;
    .ysicon {
      font-size: 38px;
      margin-right: 8px;
    }
  }

  > .fold {
    .ysicon {
      transform: rotate(90deg);
    }
  }

  > .unfold {
    .ysicon {
      transform: rotate(-90deg);
    }
  }
`
const StyledAppList = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  .appListSection {
    div.header {
      margin-top: 15px;
      display: flex;
      flex-direction: row;
      h2 {
        font-size: 14px;
        margin: 0;
        color: rgba(0, 0, 0, 0.85);
        font-weight: 600;
        line-height: 22px;
        padding-left: 10px;
      }
      p {
        font-size: 14px;
        color: #666666;
        line-height: 22px;
        font-weight: 400;
        margin: 0;
      }
    }
  }
`

type App = {
  name: string
  icon: string
  type: string
  versions: [string, string][]
  is_trial: boolean
}

type Props = {
  is_trial: boolean
  action?: any
}

export const AppList = observer(function AppList(props: Props) {
  const store = useStore()
  const state = useLocalStore(() => ({
    searchKey: '',
    setSearchKey(key) {
      this.searchKey = key
    },
    active: true,
    setActive(active) {
      this.active = active
    },
    get apps(): App[] {
      const { searchKey } = this
      let list =
        localStorage.getItem('FLAG_ENTERTAINMENT') === 'undefined' ||
        !localStorage.getItem('FLAG_ENTERTAINMENT')
          ? appList.publishedAppList
          : JSON.parse(localStorage.getItem('FLAG_ENTERTAINMENT') || '[]')

      const key = searchKey.toLocaleLowerCase()
      if (searchKey) {
        list = list.filter(item => item.name.toLowerCase().includes(key))
      }

      const apps = list.reduce((res, item) => {
        res[item.type] = res[item.type] || {
          name: item.type,
          icon: item.icon,
          type: item.type,
          versions: []
        }
        res[item.type].versions.push([item.id, item.version])

        return res
      }, {})
      return Object.values(apps)
    }
  }))

  const AppVersionSelect = () => {
    return (
      state.apps.length && (
        <StyledAppList>
          <div className='appListSection'>
            <div className='header'>
              {state.apps.map((app: App) => {
                return (
                  app.versions.length > 0 &&
                  app.versions.map(v => {
                    if (v[0] === store.currentAppId) {
                      return (
                        <Software
                          action={props.action}
                          key={app.type}
                          {...app}
                          icon={app.icon}
                        />
                      )
                    } else {
                      return ''
                    }
                  })
                )
              })}
            </div>
          </div>
        </StyledAppList>
      )
    )
  }

  return (
    <Section
      className='appList'
      title={AppVersionSelect()}
      // toolbar={
      //   <StyledToolbar>
      //     <div className='appSearch'>
      //       <Input.Search
      //         placeholder='请输入关键字'
      //         onChange={e => state.setSearchKey(e.target.value)}
      //       />
      //     </div>
      //   </StyledToolbar>
      // }
    ></Section>
  )
})
