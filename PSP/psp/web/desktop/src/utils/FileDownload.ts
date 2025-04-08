import * as qs from 'qs'
import { currentUser } from '@/domain'
import { Fetch } from '@/utils'

interface IPreDownload {
  file_paths: string[]
  cross: boolean
  file_name: string
  is_cloud: boolean
  user_name: string
  is_compress: boolean
}
function preDownload({ file_paths, cross = false, is_cloud = false,file_name,user_name,is_compress }:IPreDownload) {
  return Fetch.post('/storage/batchDownloadPre', {
    file_paths,
    file_name,
    cross,
    is_cloud,
    user_name,
    is_compress
  }).then(res => {
    if (!res.success) {
      throw Error(
        res.message || '下载失败'
      )
    }
    return res
  })
}
function download({ token, is_cloud }) {
  const a = document.createElement('a')
  a.href = `/api/v1/storage/batchDownload?${qs.stringify(
    { token,is_cloud },
    { indices: false }
  )}`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

export const downloadFile = async (
  paths,
  cross = false,
  is_cloud = false,
  types,
  user_name
) => {
  let isCompress = true
  const names = paths.map(path => path.split('/').pop())
  let name = ''
  if (names.length === 1) {
    // 是单个文件
    if (types[0]) {
      isCompress = false
      name = names[0]
    } else {
      // 是单个文件夹
      name = names[0] + '.zip'
    }
  } else {
    name = `[批量下载]${
      names.length > 2
        ? names.slice(0, 2).join('、') + '等.zip'
        : names.join('、') + '.zip'
    }`
  }
  
  const bodyParams =  {
    file_paths: paths,
    file_name: name,
    cross,
    is_cloud,
    is_compress: isCompress, //是否压缩
    user_name: user_name || currentUser.name
  }
  const res = await preDownload({...bodyParams})

  download({token: res.data.token, is_cloud})
}