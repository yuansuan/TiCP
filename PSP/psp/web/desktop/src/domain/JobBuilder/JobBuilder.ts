/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Http } from '@/utils'
import { boxServer } from '@/server'
import { message } from 'antd'
import { env } from '@/domain'
import { INNER_COMPANY_IDS } from '@/constant'

/**
 * 获取重提交信息
 * @param id 重提交作业/集 id
 * @param type 重提交类型
 * @param clean 是否要清理draft 一般预先验证时需要清理(true)，正式进入重提交无需清理(false)
 */
export async function getRedeployInfo({
  id,
  type,
  clean
}: {
  id: string
  type: 'job' | 'jobset'
  clean: boolean
}) {
  const { data } = await Http.get(`job/restore/${id}`, {
    params: { type }
  })

  return data
}

export async function getContinuousRedeployInfo({
  id,
  type,
  clean
}: {
  id: string
  type: 'job' | 'jobset'
  clean: boolean
}) {
  // 续算提交 获取参数data
  const { data } = await Http.get(`job/restore/${id}`, {
    params: { type }
  })

  const {
    data: { files }
  } = await boxServer.list({
    path: id
  })

  const datFile = files
    .filter(info => /.dat$/.test(info.name))
    .sort((a, z) => {
      return (z.name.match(/\d*$/)[0] || 0) - (a.name.match(/\d*$/)[0] || 0)
    })
    .pop()

  const casFile = files
    .filter(info => /.cas$/.test(info.name))
    .sort((a, b) => b.name.length - a.name.length)
    .pop()
  const jouFile = files
    .filter(info => /.jou$/.test(info.name))
    .sort((a, b) => {
      if (a.mod_time === b.mod_time) {
        return a.name.localeCompare(b.name)
      } else {
        return a.mod_time - b.mod_time
      }
    })
    .pop()

  // 获取文件内容
  const content =
    (await boxServer.getContent({
      path: `${id}/${jouFile?.name}`,
      offset: 0,
      length: jouFile?.size || 0,
      bucket: 'result'
    })) || ''

  const jouName = jouFile?.name || 'default.jou'
  // 修改文件内容
  const newContent = content?.replace(
    '/solve/initialize/initialize-flow',
    `/file read-data "${datFile?.name}"`
  )
  // 上传文件
  await boxServer.upload({
    file: new window.File([newContent], jouName),
    path: `./${jouName}`
  })

  return data
}
