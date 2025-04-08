/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'

import { Http, openVisual } from '@/utils'
import { workspaceList } from '@/domain/Workspace'
function openRemoteApp({ app, files }) {
  Http.post(
    '/visual/worktask',
    {
      from: 'computer',
      app_id: parseInt(app.app_id),
      template_name: app.name,
      app_param: app.app_params_with_file,
      app_param_paths: files,
      workspace_type: workspaceList.currentWorkspace?.isSharedWorkspace
        ? 'workspace'
        : 'home',
      storage_path: workspaceList.currentWorkspace?.storage_path,
    },
    { baseURL: '' }
  ).then(res => {
    const { link, id, user_id } = res.data
    openVisual(link, user_id, id)
  })
}

export default function RemoteItem({ app, files }) {
  return (
    <div key={app.app_id} onClick={() => openRemoteApp({ app, files })}>
      {app.name}
    </div>
  )
}
