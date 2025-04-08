/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Icon } from '@/components'
import { RouterType } from '@/components/PageLayout/typing'
import {
  DesktopOutlined,
  AppstoreOutlined,
  FormOutlined,
  ClusterOutlined,
  SettingOutlined,
  ProfileOutlined,
  FundOutlined,
  CreditCardOutlined
} from '@ant-design/icons'

const iconStyle = {
  position: 'absolute',
  right: '-35px',
  top: '-2px'
}
const customName = name => (
  <span style={{ position: 'relative' }}>
    <span>{name}</span>
    <Icon type='nys-new' style={iconStyle as any} />
  </span>
)

export const NORMAL_PAGES: RouterType[] = []

export const SYS_PAGES: RouterType[] = [
  {
    name: '3D可视化',
    path: '/sys/visualmgr',
    visible: false,
    key: '/sys/visualmgr',
    icon: <AppstoreOutlined rev={'none'} />,
    component: () => import('./pages/Admin/VisualMgr')
  },
  {
    name: '计算应用',
    path: '/sys/template',
    visible: true,
    key: '/sys/template',
    icon: <FormOutlined rev={'none'} />,
    component: () => import('./pages/Admin/AppMG')
  },
  {
    name: '作业统计',
    path: '/sys/jobReport',
    visible: true,
    key: '/sys/jobReport',
    icon: <ProfileOutlined rev={'none'} />,
    component: () => import('./pages/Admin/JobReport')
  },
  {
    name: '项目管理',
    path: '/sys/projectMG',
    visible: true,
    key: '/sys/projectMG',
    // perms: ['sys_manager'],
    icon: <ProfileOutlined rev={'none'} />,
    component: () => import('./pages/Admin/ProjectMG')
  },
  {
    name: '报表管理',
    path: '/sys/Report',
    visible: true,
    key: '/sys/Report',
    icon: <FundOutlined rev={'none'} />,
    component: () => import('./pages/UniteReport')
  },
  {
    name: '许可证管理',
    path: '/sys/license_mgr',
    key: '/sys/license_mgr',
    visible: true,
    exact: true,
    icon: <CreditCardOutlined />,
    component: () => import('./pages/Admin/LicenseMgr')
  },
  {
    name: '集群管理',
    path: '/sys/node',
    key: '/sys/node',
    visible: true,
    icon: <ClusterOutlined rev={'none'} />,
    component: () => import('./pages/Admin/NodeMG/NodeList')
  },
  {
    component: () =>
      import(
         './pages/Admin/NodeMG/NodeDetail'
      ),
    path: '/node/:name'
  },
  {
    path: '/sys/template-edit',
    component: () => import('./pages/Admin/AppMG/Editor')
  },
  {
    name: '许可证详情',
    visible: false,
    path: '/sys/license_mgr/:id',
    exact: true,
    component: () => import('./pages/Admin/LicenseMgr/LicenseDetail')
  },
  {
    name: '许可证添加',
    visible: false,
    exact: true,
    path: '/sys/license_mgr-add',
    component: () => import('./pages/Admin/LicenseMgr/LicenseAdd')
  },
  {
    name: '许可证编辑',
    visible: false,
    exact: true,
    path: '/sys/license_mgr/license/:id',
    component: () => import('./pages/Admin/LicenseMgr/LicenseAdd')
  },
  {
    name: '在线用户详情',
    visible: false,
    exact: true,
    path: '/sys/onlineUser/:name',
    component: () => import('./pages/Admin/SysOnlineUserMG/OnlineUserDetail')
  },
  {
    name: '系统设置',
    icon: <SettingOutlined rev={'none'} />,
    children: [
      {
        name: '用户管理',
        path: '/sys/user',
        visible: true,
        key: '/sys/user',
        component: () => import('./pages/Admin/UserMG')
      },
      {
        name: '在线用户',
        path: '/sys/onlineUser',
        visible: true,
        key: '/sys/onlineUser',
        component: () => import('./pages/Admin/SysOnlineUserMG/OnlineUserList')
      },
      {
        name: '全局设置',
        path: '/sys/global',
        visible: true,
        key: '/sys/global',
        component: () => import('./pages/Admin/SysSetting')
      }
    ]
  }
]

export const PROJECT_PAGES: RouterType[] = [
  // 新版作业管理
  {
    component: () => import('./pages/NewJobManager'),
    path: '/new-jobs',
    name: '作业管理',
    customName: () => customName('作业管理'),
    icon: <Icon style={{ fontSize: 16 }} type='job_mgt_default' />,
    children: [
      {
        path: '/new-job/:id',
        name: '作业详情',
        component: () => import('@/pages/NewJobDetail'),
        isMenu: false
      },
      {
        path: '/new-job-set/:id',
        name: '作业集详情',
        component: () => import('@/pages/NewJobSetDetail'),
        isMenu: false
      }
    ]
  },
  {
    component: () => import('./pages/VisList/List'),
    path: '/vis-session',
    name: '3D云应用',
    icon: <DesktopOutlined rev={'none'} />,
    customName: () => customName('3D云应用')
  },
  {
    path: '/company/bill_user',
    icon: <Icon style={{ fontSize: 16 }} type='account_settings_hover' />,
    name: '个人账单',
    component: () => import('./pages/BillForUser')
  },

  {
    component: () => import('./pages/PersonalSetting'),
    path: '/personal-setting',
    name: '个人设置',
    isMenu: false
  },
  {
    component: () => import('./pages/MessageMGT'),
    path: '/messages',
    name: '消息中心',
    isMenu: false
  }
]
