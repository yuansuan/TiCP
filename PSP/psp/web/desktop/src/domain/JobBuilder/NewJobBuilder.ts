/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { Http } from '@/utils'
import Draft, { DraftType } from '../Box/NewDraft'
import { newBoxServer } from '@/server'
import { message } from 'antd'
import { env } from '@/domain'

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
  const redeployDraft = new Draft(DraftType.Redeploy)
  const { input_folder_uuid } = data
  try {
    await redeployDraft.back(input_folder_uuid, {
      formatErrorMessage: () => '仅可重提交近30天内提交的作业'
    })
  } catch (e) {
    return Promise.reject()
  }

  clean && (await redeployDraft.clean())
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
  } = await newBoxServer.list({
    path: id,
    bucket: 'result',
    url: ''
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
  const continuousDraft = new Draft(DraftType.Continuous)
  // 传回续传草稿
  await continuousDraft.putResultToDraft(
    undefined,
    id,
    [datFile?.name, casFile?.name, jouFile?.name].filter(Boolean)
  )

  // 获取文件内容
  const content =
    (await newBoxServer.getContent({
      path: `${id}/${jouFile?.name}`,
      offset: 0,
      length: jouFile?.size || 0,
      bucket: 'result',
      url: ''
    })) || ''

  const jouName = jouFile?.name || 'default.jou'
  // 修改文件内容
  const newContent = content?.replace(
    '/solve/initialize/initialize-flow',
    `/file read-data "${datFile?.name}"`
  )
  // 上传文件
  await newBoxServer.upload({
    file: new window.File([newContent], jouName),
    path: `./${jouName}`,
    bucket: 'draft',
    bucket_keys: JSON.stringify({ draft_type: continuousDraft.type }),
    url: ''
  })

  clean && (await continuousDraft.clean())
  return data
}
