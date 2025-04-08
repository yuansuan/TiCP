/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import { Redirect, Route, Router, Switch } from 'react-router-dom'
import history from '@/utils/history'
import Loadable from 'react-loadable'
import { useSelector } from 'react-redux'
import { BackTop } from 'antd'
import { RouterType } from '@/components/PageLayout/typing'
import { ToolBar } from '@/utils/general'
import { SYS_PAGES } from '@/router'
import { CompanyLayout } from '@/components'
import styled from 'styled-components'
import { currentUser } from '@/domain'

const Wrapper  = styled.div` 
  height: calc(100vh - 48px);
` 
export const EnterpriseManage = () => {
  const apps = useSelector(state => state.apps)
  const wnapp = useSelector(state => state.apps.enterpriseManage)
  const globalConfig = JSON.parse(window.localStorage.getItem('GlobalConfig') || '{}')

  const New_SYS_PAGES = SYS_PAGES.map(route => {
    if (route.path === '/sys/visualmgr') {
      route.visible = globalConfig?.enable_visual
    }
    return route
  }).filter(route => {
    if (route.perms && route.perms.length !== 0) {
      return route.perms.some(perm => currentUser.permKeys.includes(perm)) 
    } else {
      return true
    }
  })


  const Loading = ({ error }) => {
    if (error) {
      throw error
    }

    return <div>Loading...</div>
  }

  const CompanyRoute = ({ component: Component, render, path, ...rest }) => {
    return (
      <Route
        {...rest}
        path={path}
        render={matchProps => (
          <CompanyLayout routers={New_SYS_PAGES}>
            <>
              {/* <Notice /> */}
              <BackTop style={{ right: 15, bottom: 130 }} />
              {Component && <Component {...matchProps} />}
              {render && render(matchProps)}
            </>
          </CompanyLayout>
        )}
      />
    )
  }

  const createRoute = CustomRoute => page => {
    const LoadableComponent = Loadable({
      loader: page.component,
      loading: Loading
    })

    return (
      <CustomRoute
        exact={page.exact}
        path={page.path}
        key='ys'
        render={props => {
          return <LoadableComponent {...props} refresh={apps.hz === apps.enterpriseManage.z}/>
        }}
      />
    )
  }

  function createRouters(routers: RouterType[], RouteType: any = Route) {
    return routers.map(item => {
      return [
        item.path && createRoute(RouteType)(item),
        item.children &&
          createRouters(
            item.children.map(child => ({
              ...child,
              visible: item.visible
            })),
            RouteType
          )
      ]
    })
  }

  const redirectPath = New_SYS_PAGES.filter(r => r.visible === true)[0]?.path

  return wnapp ? (
    <div
      className='calcApp floatTab dpShad systemCompany'
      data-size={wnapp.size}
      id={wnapp.icon + 'App'}
      data-max={wnapp.max}
      style={{
        ...(wnapp.size == 'cstm' ? wnapp.dim : null),
        zIndex: wnapp.z
      }}
      data-hide={wnapp.hide}>
      <ToolBar
        app={wnapp.action}
        icon={wnapp.icon}
        size={wnapp.size}
        name={wnapp.title}
      />
      <div className='windowScreen flex flex-col' data-dock='true'>
        {!wnapp.hide && (
          <Wrapper>
            <Router history={history}>
              {/* <WebConfigComponent /> */}

              <Switch>
                {createRouters(New_SYS_PAGES, CompanyRoute)}

                <Route
                  path='/'
                  render={() => <Redirect to={redirectPath} />}
                />
                <Route
                  component={Loadable({
                    loader: () => import('@/pages/404'),
                    loading: Loading
                  })}
                />
              </Switch>
            </Router>
          </Wrapper>
        )}
      </div>
    </div>
  ) : null
}
