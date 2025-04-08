/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Section } from '../components'
import { Software } from './Software'
import { ProjectSelector } from './Project'
import { appList, account } from '@/domain'
import { clusterCores } from '@/domain/ClusterCores'
import { useStore } from '../store'
import { Balance } from '@/components/HeaderToolbar'

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
  versions: [string, string, string][]
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
    get isCloudApp() {
      return store?.data?.currentApp?.compute_type === 'cloud'
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
        res[item.type].versions.push([item.id, item.version, item.compute_type])

        return res
      }, {})
      return Object.values(apps)
    }
  }))

  const ProjectAndAppVersionSelector = () => {
    return (
      state.apps.length && (
        <StyledAppList>
          <div className='appListSection'>
            <div className='header'>
              <span style={{ paddingRight: 5 }}>
                <ProjectSelector />
              </span>
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
      title={ProjectAndAppVersionSelector()}
      toolbar={
        state.isCloudApp ? (
          <Balance />
        ) : (
          <div className='core'>
            总核数：{clusterCores.total_cores}
            {'\u00A0'} | {'\u00A0'}可用核数：{clusterCores.available_cores}
          </div>
        )
      }></Section>
  )
})
