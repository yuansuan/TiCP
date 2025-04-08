/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import { useState, useEffect } from 'react'
import { Dropdown } from 'antd'

import { Http } from '@/utils'
import { Button } from '@/components'
import RemoteItem from './RemoteItem'
import { StyledOverlay } from './style'
import { RootPoint } from '@/domain/FileSystem'

export default function OpenRemote({
  selectedKeys = [],
  point,
  parentPath,
}: {
  selectedKeys: Array<string>
  point: RootPoint
  parentPath: string
}) {
  const [apps, setApps] = useState([])

  useEffect(() => {
    Http.get('/remote_app/list', {
      params: {
        state: 'published',
      },
    }).then(res => {
      setApps(res.data.applications)
    })
  }, [])

  if (apps.length === 0) {
    return null
  }

  const file: string = selectedKeys.length > 0 ? selectedKeys[0] : undefined
  // const ext = file && file.substring(file.lastIndexOf('.') + 1)
  const currentNode = point.filterFirstNode(node => node.path === parentPath)
  const fileObj =
    file &&
    currentNode &&
    currentNode.children.find(p => {
      return p.path === file
    })

  const matchedApps = apps.filter(app => {
    const ftypes: string = app.app_file_types
    if (!ftypes || !fileObj) {
      return false
    }
    const reList = ftypes.split(',').map(ftype => {
      let reStr = ftype.trim().replace(/[.+^${}()|[\]\\]/g, '\\$&') // exclude * and ?
      return reStr.replace(/\*/g, '(.*)').replace(/\?/g, '(.)')
    })
    const matched = reList.find(re => {
      const rexp = new RegExp(re)
      return rexp.test(fileObj.name)
    })
    return !!matched
  })
  return (
    <Dropdown
      disabled={
        selectedKeys.length !== 1 ||
        (fileObj && !fileObj.isFile) ||
        !matchedApps.length
      }
      overlay={
        <StyledOverlay>
          {matchedApps.map(app => (
            <RemoteItem key={app.app_id} app={app} files={selectedKeys} />
          ))}
        </StyledOverlay>
      }>
      <Button type='primary' ghost>
        打开方式
      </Button>
    </Dropdown>
  )
}
