/* Copyright (C) 2016-present, Yuansuan.cn */
import { currentUser } from '@/domain'
import { DeployMode } from '@/constant'

export const gene_name = () =>
  Math.random().toString(36).substring(2, 10).toUpperCase()

let installed = JSON.parse(localStorage.getItem('installed') || '[]')

export function generateCalcApp(data) {
  if (!Array.isArray(data)) return
  const iconTypeMap = {}
  const appType = [
    ...data.map(item => {
      return {
        type: item.type,
        id: item.id,
        icon: item.icon,
        name: item.name,
        computeType: item.compute_type
      }
    })
  ]

  // data.forEach(item => {
  //   if (!iconTypeMap.hasOwnProperty(item.type)) {
  //     iconTypeMap[item.type] = item.icon
  //   }
  // })

  const handleIcon = item => {
    // TODO 桌面icon
    switch (true) {
      // case /^Star-?CCM\+?\d*/gi.test(item.type):
      //   return 'Telemac'
      // case item === 'Telemac':
      //   return 'Telemac'
      // case item === 'Code Saturne':
      //   return 'Code Saturne'
      // case item === 'Code Aster':
      //   return 'Code Aster'
      // case item === 'Fluent':
      // return 'fluent'
      default:
        return item.icon || item.type
    }
  }
  const generateAppInfo = appType.map(item => {
    const icon = handleIcon(item)
    return {
      name: item.type,
      title: item.type,
      icon: icon,
      routerPath: `/new-job-creator?action=${item.type}&appType=${item.type}&id=${item.id}`,
      type: 'app',
      action: item.type,
      computeType: item.computeType,
      appType: 'generate',
      className: '',
      renderType: 'calcApp'
      // action: 'CALCUAPP',
    }
  })

  return {
    appType,
    generateAppInfo
  }
}

function handleIcon(name, icon, id) {
  switch (true) {
    default:
      return icon || id + 'CLOUDAPP'
  }
}

export function generateCloudApp(cloudSoftware) {
  const softWareType = cloudSoftware?.map(item => item.name)
  const generateCloudApp = cloudSoftware?.map(item => {
    const icon = handleIcon(item.name, item.icon, item.id)

    const result = {
      name: item.name,
      title: item.name,
      icon,
      id: item.session_id,
      routerPath: `/vis-session?actionApp=${icon}`,
      type: 'app',
      action: icon,
      appType: 'generate',
      className:
        item.status === 'STARTED' ? 'CloudAppWrap_open' : 'CloudAppWrap_close',
      menu: item.status === 'STARTED' ? 'cloudMenuOpen' : 'cloudMenuClose',
      renderType: 'cloudApp',
      url: window.atob(item.stream_url)
    }

    return result
  })

  window.localStorage.setItem(
    'GENERATECLOUDAPP',
    JSON.stringify(generateCloudApp)
  )
  return { softWareType, CloudApp: generateCloudApp }
}

const apps = [
  {
    name: 'EnterpriseManage', // 应用名称
    title: '系统管理',
    perm: ['sys_manager'],
    routerPath: '/sys',
    icon: 'enterpriseManage', // 应用icon
    type: 'app', // 应用类型
    action: 'ENTERPRISEMANAGE' // 触发的action，一般为name的全大写
  },
  {
    name: 'FileManage',
    title: '文件管理',
    icon: 'explorer',
    perm: ['file_manager'],
    routerPath: '/files',
    type: 'app',
    action: 'FILEMANAGE'
  },
  {
    name: 'JobManage',
    title: '作业管理',
    icon: 'jobLists',
    perm: ['job_manager', 'personal_job_manager'],
    routerPath: '/new-jobs',
    type: 'app',
    action: 'JOBMANAGE'
  },
  {
    name: 'FileExplorer',
    title: 'File Explorer',
    icon: 'FileExplorer',
    routerPath: '/files',
    type: 'app',
    action: 'FILEXOLORER'
  },
  {
    name: 'Calculator',
    title: '作业提交',
    icon: 'calculator',
    routerPath: '/new-job-creator',
    type: 'app',
    action: 'CALCUAPP'
  },
  {
    name: 'Messages', // 应用名称
    title: '消息中心',
    routerPath: '/messages?tab=messages',
    icon: 'mail', // 应用icon
    type: 'app', // 应用类型
    action: 'MESSAGES' // 触发的action，一般为name的全大写
  },
  {
    name: '3dcloudApp', // 应用名称
    title: '3D云应用',
    icon: '3dcloudApp', // 应用icon
    routerPath: '/vis-session',
    type: 'app', // 应用类型
    action: '3DCLOUDAPP' // 触发的action，一般为name的全大写
  },
  {
    name: 'NewJobDetail', // 应用名称
    title: '作业详情',
    icon: 'jobDetail', // 应用icon
    type: 'app', // 应用类型
    action: 'NEWJODETAIL' // 触发的action，一般为name的全大写
  },
  {
    name: 'NewJobSetDetail', // 应用名称
    title: '作业集详情',
    icon: 'jobSetDetail', // 应用icon
    type: 'app', // 应用类型
    action: 'NEWJOBSETDETAIL' // 触发的action，一般为name的全大写
  },
  {
    name: 'Template', // 应用名称
    title: '模版',
    icon: 'template', // 应用icon
    type: 'app', // 应用类型
    action: 'TEMPLATE' // 触发的action，一般为name的全大写
  },
  {
    name: 'Dashboard',
    title: '集群监控',
    perm: ['cluster_monitor'],
    icon: 'dashboard',
    type: 'app',
    action: 'DASHBOARD'
  },
  {
    name: 'UserLog',
    title: '用户日志',
    icon: 'userlog',
    perm: ['job_manager', 'personal_job_manager'],
    type: 'app',
    action: 'USERLOG'
  },
  {
    name: 'AuditLog',
    title: '审计日志',
    icon: 'auditLog',
    perm: [
      'normal_audit_log',
      'system_admin_audit_log',
      'security_admin_audit_log'
    ],
    type: 'app',
    action: 'AUDITLOG'
  },
  {
    name: 'SecurityApproval',
    title: '安全审批', // 审批申请
    icon: 'securityApproval',
    perm: ['security_approval', 'sys_manager'], // no_perm means always show
    type: 'app',
    action: 'SECURITYAPPROVAL'
  },
  {
    name: 'Report',
    title: '报表管理',
    icon: 'template',
    type: 'app',
    action: 'REPORT'
  },
  {
    name: 'ProjectMG',
    title: '项目管理',
    perm: ['project_manager', 'personal_project_manager'],
    icon: 'projectMG',
    type: 'app',
    action: 'PROJECTMG'
  }
  // ...generateAppInfo,
  // ...generateCloudApp
]

for (let i = 0; i < installed.length; i++) {
  installed[i].action = gene_name()
  apps.push(installed[i])
}

// 注册应用

let mainApp = []
// 云端和本地的求解应用相同，只显示本地的
function uniqueArrayWithPriority(arr, propertyName) {
  const uniqueMap = new Map()

  for (const item of arr) {
    if (!uniqueMap.has(item[propertyName])) {
      uniqueMap.set(item[propertyName], item)
    } else if (
      item.computeType === 'local' &&
      uniqueMap.get(item[propertyName]).computeType !== 'local'
    ) {
      uniqueMap.set(item[propertyName], item)
    }
  }

  return Array.from(uniqueMap.values())
}

export const GenerateDesktopApps = (apps, generateDesktop) => {
  const globalConfig = JSON.parse(localStorage.getItem('GlobalConfig') || '{}')
  const systemPerm = JSON.parse(localStorage.getItem('SystemPerm') || '[]')
  const cloudApps = apps.filter(app => app.renderType === 'cloudApp')

  const filterCloudApps = apps?.filter(app => app.renderType !== 'cloudApp')
  const filteredApps = uniqueArrayWithPriority(filterCloudApps, 'name')

  const localDesktopApp =
    apps
      .filter(app =>
        systemPerm?.some(
          p => app?.perm?.includes(p?.key) || app?.perm?.includes('no_perm')
        )
      ).map(app => app.name) || []

  // 全局控制可视化应用显示，以及权限共同判断
  if (
    globalConfig?.enable_visual &&
    systemPerm?.some(p => ['visual'].includes(p?.key))
  ) {
    localDesktopApp.push('3dcloudApp')
    // filteredApps.push(...cloudApps)
  }
  mainApp = [...localDesktopApp]

  const desktopType = [...mainApp, ...generateDesktop]

  return filteredApps
    .filter(x => desktopType.includes(x.name))
    .sort((a, b) => {
      return desktopType.indexOf(a.name) > desktopType.indexOf(b.name) ? 1 : -1
    })
}

const { taskbar, desktop, pinned, recent } = {
  // 设置底部taskbar的数组
  taskbar: mainApp,
  // 设置桌面应用
  desktop: [...mainApp],
  // 设置便捷连接
  pinned: [],

  recent: []
}
export const taskApps = apps.filter(x => taskbar.includes(x.name))

export const desktopApps = apps
  .filter(x => desktop.includes(x.name))
  .sort((a, b) => {
    return desktop.indexOf(a.name) > desktop.indexOf(b.name) ? 1 : -1
  })

export const pinnedApps = apps
  .filter(x => pinned.includes(x.name))
  .sort((a, b) => {
    return pinned.indexOf(a.name) > pinned.indexOf(b.name) ? 1 : -1
  })

export const recentApps = apps
  .filter(x => recent.includes(x.name))
  .sort((a, b) => {
    return recent.indexOf(a.name) > recent.indexOf(b.name) ? 1 : -1
  })

export const allApps = apps.filter(app => {
  return app.type === 'app'
})

export const dfApps = {
  taskbar,
  desktop,
  pinned,
  recent
}
export default apps
