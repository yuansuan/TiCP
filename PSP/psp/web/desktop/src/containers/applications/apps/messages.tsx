/* Copyright (C) 2016-present, Yuansuan.cn */
import React from 'react'
import { Route, Router, Switch } from 'react-router-dom'
import history from '@/utils/history'
import Loadable from 'react-loadable'
import { useSelector } from 'react-redux'
import { RouterType } from '@/components/PageLayout/typing'
import { ToolBar } from '@/utils/general'

import { NORMAL_PAGES } from '@/router'

export const Messages = () => {
  const wnapp = useSelector(state => state.apps.mail)

  const Loading = ({ error }) => {
    if (error) {
      throw error
    }

    return <div>Loading...</div>
  }

  const createRoute = CustomRoute => page => {
    const { visible } = page
    const LoadableComponent = Loadable({
      loader: page.component,
      loading: Loading
    })
    const NoPerm = Loadable({
      loader: () => import('@/pages/403'),
      loading: Loading
    })

    const jobPathMap = {
      '^/standard-jobs': {
        destPath: '/jobs'
      },
      '^/jobs': {
        destPath: '/standard-jobs'
      },
      '^/job-creator': {
        destPath: '/standard-job-creator'
      },
      '^/standard-job-creator': {
        destPath: '/job-creator'
      },
      '^/standard-job/': {
        destPath: '/jobs'
      },
      '^/job/': {
        destPath: '/standard-jobs'
      },
      '^/jobset/': {
        destPath: '/standard-jobs'
      }
      // '^/files': {
      //   destPath: '/files'
      // }
    }

    return (
      <CustomRoute
        exact={page.exact}
        path={page.path}
        key='ys'
        render={props => {
          if (visible !== undefined && page.path.includes('job')) {
            if ((typeof visible === 'function' && !visible()) || !visible) {
              let pathKey = Object.keys(jobPathMap).filter(key => {
                const pathReg = new RegExp(key)
                return pathReg.test(page.path)
              })[0]
              const jobPath = jobPathMap[pathKey]
              if (jobPath) {
                history.push(jobPath.destPath)
              } else {
                return <NoPerm {...props} />
              }
            }
          }
          return <LoadableComponent {...props} />
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

  return wnapp ? (
    <div
      className='calcApp floatTab dpShad'
      data-menu='application'
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
          <Router history={history}>
            {/* <WebConfigComponent /> */}

            <Switch>
              {NORMAL_PAGES.map(createRoute(Route))}

              <Route
                // exact
                // path='/'
                // render={() => <Redirect to='/company/workspaces' />}
                component={Loadable({
                  loader: () => import('@/pages/MessageMGT'),
                  loading: Loading
                })}
              />
              <Route
                component={Loadable({
                  loader: () => import('@/pages/404'),
                  loading: Loading
                })}
              />
            </Switch>
          </Router>
        )}
      </div>
    </div>
  ) : null
}
