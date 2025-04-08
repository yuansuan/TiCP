/* Copyright (C) 2016-present, Yuansuan.cn */
const defState = {
  hide: true,
  top: 80,
  left: 360,
  opts: 'desk',
  attr: null,
  dataset: null,
  data: {
    desk: {
      width: '310px',
      secwid: '200px'
    },
    task: {
      width: '220px',
      secwid: '120px',
      ispace: false // show the space for icons in menu
    },
    app: {
      width: '310px',
      secwid: '200px'
    },
    fileManager: {
      width: '310px',
      secwid: '200px'
    },
    cloudMenuClose: {
      width: '310px',
      secwid: '200px'
    },
    cloudMenuOpen: {
      width: '310px',
      secwid: '200px'
    },
    fileMenu: {
      width: '310px',
      secwid: '200px'
    },
    rightMenu: {
      width: '310px',
      secwid: '200px'
    },
    multipleFiles:{
      width: '310px',
      secwid: '500px'
    },
    singleFile:{
      width: '310px',
      secwid: '500px'
    },
    noFile: {
      width: '310px',
      secwid: '200px'
    }
  },
  menus: {
    desk: [
      {
        name: '查看',
        icon: 'view',
        type: 'svg',
        opts: [
          {
            name: '大图标',
            action: 'changeIconSize',
            payload: 'large'
          },
          {
            name: '中等图标',
            action: 'changeIconSize',
            payload: 'medium'
          },
          {
            name: '小图标',
            action: 'changeIconSize',
            payload: 'small',
            dot: true
          },
          {
            type: 'hr'
          },
          {
            name: '展示桌面图标',
            action: 'deskHide',
            check: true
          }
        ]
      },
      {
        name: '刷新',
        action: 'refresh',
        type: 'svg',
        icon: 'refresh'
      }
    ],
    task: [
      {
        name: 'Align icons',
        opts: [
          {
            name: 'Left',
            action: 'changeTaskAlign',
            payload: 'left'
          },
          {
            name: 'Center',
            action: 'changeTaskAlign',
            payload: 'center',
            dot: true
          }
        ]
      },
      {
        type: 'hr'
      },
      {
        name: 'Search',
        opts: [
          {
            name: 'Show',
            action: 'TASKSRCH',
            payload: true
          },
          {
            name: 'Hide',
            action: 'TASKSRCH',
            payload: false
          }
        ]
      },
      {
        name: 'Widgets',
        opts: [
          {
            name: 'Show',
            action: 'TASKWIDG',
            payload: true
          },
          {
            name: 'Hide',
            action: 'TASKWIDG',
            payload: false
          }
        ]
      },
      {
        type: 'hr'
      },
      {
        name: 'Show Desktop',
        action: 'SHOWDSK'
      }
    ],
    app: [
      {
        name: '打开',
        action: 'performApp',
        payload: 'open'
      }
    ],
    fileManager: [
      {
        name: '新建文件夹',
        action: 'fileAction',
        payload: 'createFile'
      },
      {
        name: '上传文件',
        action: 'fileAction',
        payload: 'uploadFile'
      },
      {
        name: '上传文件夹',
        action: 'fileAction',
        payload: 'uploadFiles'
      }
    ],
    cloudMenuOpen: [
      {
        name: '打开',
        action: 'performApp',
        payload: 'open'
      },
      {
        name: '关闭',
        action: 'performApp',
        payload: 'delshort'
      }
    ],
    cloudMenuClose: [
      {
        name: '创建',
        action: 'performApp',
        payload: 'open',
        key: 'create'
      }
    ],
    fileMenu: [
      {
        name: '重新命名',
        action: 'fileMenu',
        payload: 'rename'
      },
      {
        name: '移动',
        action: 'fileMenu',
        payload: 'move'
      },
      {
        name: '删除',
        action: 'fileMenu',
        payload: 'remove'
      },
      {
        name: '删除',
        action: 'fileMenu',
        payload: 'remove'
      },
    ],
    // 文件多选菜单
    multipleFiles: [
      {
        name: '下载',
        action: 'multipleFiles',
        payload: 'download'
      },
      {
        name: '新建文件夹',
        action: 'fileAction',
        payload: 'createFile'
      },
      {
        name: '移动',
        action: 'multipleFiles',
        payload: 'move'
      },
      {
        name: '删除',
        action: 'multipleFiles',
        payload: 'remove'
      },
      {
        name: '压缩',
        action: 'multipleFiles',
        payload: 'compress'
      },
    ],
    singleFile: [
      {
        name: '下载',
        action: 'singleFile',
        payload: 'download'
      },
      {
        name: '新建文件夹',
        action: 'singleFile',
        payload: 'createFile'
      },
      {
        name: '移动',
        action: 'singleFile',
        payload: 'move'
      },
      {
        name: '删除',
        action: 'singleFile',
        payload: 'remove'
      },
      {
        name: '重新命名',
        action: 'singleFile',
        payload: 'rename'
      },
      {
        name: '分享',
        action: 'singleFile',
        payload: 'share'
      },
      {
        name: '发送',
        action: 'singleFile',
        payload: 'send'
      },
      {
        name: '压缩',
        action: 'singleFile',
        payload: 'compress'
      },
    ],
    noFile: [
      {
        name: '新建文件夹',
        action: 'noFile',
        payload: 'createFile'
      },
      {
        name: '上传文件',
        action: 'noFile',
        payload: 'uploadFile'
      },
      {
        name: '上传文件夹',
        action: 'noFile',
        payload: 'uploadFiles'
      }
    ],
    rightMenu: [
      {
        name: '移动',
        action: 'rightMenu',
        payload: 'move'
      },
      {
        name: '删除',
        action: 'rightMenu',
        payload: 'remove'
      },
      {
        name: '下载',
        action: 'rightMenu',
        payload: 'download'
      }
    ]
  }
}

const menusReducer = (state = defState, action) => {
  let tmpState = {
    ...state
  }
  if (action.type == 'MENUHIDE') {
    tmpState.hide = true
  } else if (action.type == 'MENUSHOW') {
    tmpState.hide = false
    tmpState.top = (action.payload && action.payload.top) || 272
    tmpState.left = (action.payload && action.payload.left) || 430
    tmpState.opts = (action.payload && action.payload.menu) || 'desk'
    tmpState.attr = action.payload && action.payload.attr
    tmpState.dataset = action.payload && action.payload.dataset
  } else if (action.type == 'MENUCHNG') {
    tmpState = {
      ...action.payload
    }
  } else if (action.type == 'MENUITEMUPDATE') {
    const { menu, menuItemPayload, menuItemAttr } = action.payload || {}
    const currenMenu = tmpState.menus[menu]
    if (currenMenu) {
      for (let i=0; i<currenMenu.length; i++) {
        if (currenMenu[i].payload === menuItemPayload) {
          currenMenu[i] = {
            ...currenMenu[i],
            ...menuItemAttr,
          }
          break
        }
      }
    }
  }

  return tmpState
}

export default menusReducer
