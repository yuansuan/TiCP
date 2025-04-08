import PackageList from './PackageList'
import CheckPackageTasks from './CheckPackageTasks'
import { Uploader } from '@/domain/Uploader'

const packageUploader = new Uploader()

export { default as Package } from './Package'

export const packageList = new PackageList()
export const checkPackageTasks = new CheckPackageTasks(packageUploader)

export class PackageCheckError extends Error {
  message: string

  constructor(code: number) {
    super()
    this.message = PackageCheckError.errorMap[code] || '服务出错，请联系管理员'
  }

  static errorMap = {
    260005: '安装包解压失败',
    260006: '安装包中 manifest.json 不存在',
    260007: '安装包中 manifest.json 读取失败',
    260008: '安装包中 manifest.json 格式不正确，请参考开发者文档',
    260009: '安装包中 manifest.json 中不能存在空字段，请检查',
    260010: '服务出错，请联系管理员',
    260011: '服务出错，请联系管理员',
    260012: '服务出错，请联系管理员',
    260013: '服务出错，请联系管理员',
    260014: '服务出错，请联系管理员',
    260015: '安装包已经存在',
    260016: '安装包中的应用名和版本信息已重复, 请修改 manifest.json',
    260017: '服务器错误，请联系管理员',
    260018: '删除安装包失败',
    260019: '服务出错，请联系管理员',
    260020: '服务出错，请联系管理员',
    260026: '服务出错，请联系管理员',
    260027: '应用图标不存在，请检查 manifest.json 的 icon 字段',
    260028: '服务出错，请联系管理员',
    260029: '不能删除正在安装应用的安装包',
    260031: '安装包检测失败',
    260034: '应用图标大小不能大于 1 MB',
  }
}
