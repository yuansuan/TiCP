import { observable } from 'mobx'
import { untilDestroyed } from '@/utils/operators'
import { filter } from 'rxjs/operators'
import currentUser from '@/domain/User'
import { eventEmitter, IEventData } from '@/utils'
import { Fetch } from '@/utils'
import { PackageCheckError } from '@/domain/LocalPackage'
import { message } from 'antd'
import { Task } from '@/domain/Uploader'

interface ICheckPackageTask {
  checking: boolean
  checkingTips: string
  uploadTask: Task
}

class CheckPackageTask implements ICheckPackageTask {
  @observable public checking = false
  @observable public checkingTips = ''
  @observable public uploadTask = null
  @observable public packageId = null

  constructor(packageId, checking, checkingTips, uploadTask) {
    this.packageId = packageId
    this.checking = checking
    this.checkingTips = checkingTips
    this.uploadTask = uploadTask
  }
}

export default class CheckPackageTasks {
  // one by one 一个一个 check
  packageId = null
  @observable public checking = false
  @observable public checkingTips = ''
  @observable tasks: Map<number, ICheckPackageTask> = new Map()

  public uploader = null

  constructor(uploader) {
    this.uploader = uploader
  }

  upload = (props, dirPath, packageId, afterCheckSuccess) => {
    this.packageId = packageId
    const task = this.uploader.upload(props)

    this.tasks.set(packageId, new CheckPackageTask(packageId, false, '', task))
    const currPackageTask = this.tasks.get(packageId)

    currPackageTask.uploadTask.status$
      .pipe(
        untilDestroyed(this),
        filter(status => status === 'done' || status === 'aborted')
      )
      .subscribe(async status => {
        if (status === 'aborted') {
          this.tasks.delete(this.packageId)
          return
        }
        // 上传完成
        console.debug('upload finished', props.file, dirPath)

        // 发请求, 检验文件是否符合要求
        this.checking = true
        this.checkingTips = `包 ${props.file.name} 正在检查中，请稍候... ...`

        try {
          await Fetch.get(
            `/api/v3/localpackage/check?packageName=${encodeURIComponent(
              props.file.name
            )}&packageId=${packageId}&userName=${currentUser.name}`,
            { timeout: 0 }
          )

          // 事件监听，包检测
          eventEmitter.once(
            `PACKAGE_CHECK_${props.file.name}`,
            (data: IEventData) => {
              const msg = data.message
              console.debug(`PACKAGE_CHECK_${props.file.name}:`, msg)

              if (msg.success) {
                message.success(`安装包 ${props.file.name} 检测成功`)
              } else {
                const error = new PackageCheckError(msg.errCode)
                message.error(
                  `安装包 ${props.file.name} 检测失败, 原因: ${error.message}`
                )
              }
              afterCheckSuccess()
              this.tasks.delete(this.packageId)
              this.checking = false
              this.checkingTips = ''
            }
          )

          // 超时处理
          setTimeout(() => {
            message.error(
              `安装包 ${props.file.name} 检测失败, 原因: 检测服务超时，请联系管理员`
            )
            eventEmitter.off(`PACKAGE_CHECK_${props.file.name}`)

            this.tasks.delete(this.packageId)
            this.checking = false
            this.checkingTips = ''
          }, 1000 * 60 * 30)
        } catch (e) {
          this.tasks.delete(this.packageId)
          this.checking = false
          this.checkingTips = ''
        }
      })
  }
}
