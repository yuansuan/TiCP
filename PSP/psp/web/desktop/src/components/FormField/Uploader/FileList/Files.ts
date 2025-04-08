import { action, computed, observable } from 'mobx'
import { formatByte } from '@/utils/Validator'
import { statusMap } from '@/domain/Uploader/Task'

export class FileList {
  @observable files: FileTreeFile[] = []

  setFiles(files) {
    this.files = files
  }
}

export class FileTreeFile {
  @observable name: string = undefined
  @observable parent: FileTreeFile = undefined
  @observable children: FileTreeFile[] = undefined
  @observable size: number = undefined
  @observable path: string = undefined
  @observable status: string = ''
  @observable from: string = ''
  @observable isMain: boolean = false
  @observable isDir: boolean = false
  @observable isRenaming: boolean = false
  @observable is_text: boolean = false

  @computed
  get displaySize() {
    return formatByte(this.size)
  }

  @computed
  get displayStatus() {
    return statusMap[this.status]
  }

  @computed
  get displayFrom() {
    let file = this
    while (file.parent) {
      file = file.parent as any
    }
    return file.from === 'local' ? '本地' : '服务器'
  }

  constructor(props: Partial<FileTreeFile>) {
    this.update(props)
  }

  @action
  update(props: Partial<FileTreeFile>) {
    Object.assign(this, props)
  }
}

export const currentFileList = new FileList()