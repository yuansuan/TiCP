export class MergeAndMd5Error extends Error {
  message: string

  constructor(code: number) {
    super()
    this.message = MergeAndMd5Error.errorMap[code] || '合并与校验失败'
  }

  static errorMap = {
    40019: '没有权限，上传失败',
    40012: '磁盘空间不足',
    40036: '未知错误',
    40047: '文件校验失败',
    40046: '文件合并失败',
  }
}

export default class UploadError extends Error {
  message: string

  constructor(code: number) {
    super()
    this.message = UploadError.errorMap[code] || '上传失败'
  }

  static errorMap = {
    40019: '没有权限，上传失败',
    40012: '磁盘空间不足',
    40036: '未知错误',
  }
}
