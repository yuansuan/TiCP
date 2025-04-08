import { message } from 'antd'
import { FieldType } from '@/domain/Applications/App/Field'

/**
 * check formModel when test or submit
 * @param formModel
 */
export default function checkForm(formModel) {
  return Object.values(formModel).every((item: any) => {
    // check file uploader
    if (item.type === FieldType.lsfile) {
      const files = item._files || []
      // required
      if (item.required && !files.find(file => file.status === 'done')) {
        message.error(`${item.label} 至少上传一个文件`)
        return false
      }

      if (item.required) {
        const movingFile = files.find(
          file => file.status === 'done' && !file.path
        )

        if (movingFile) {
          message.error(`文件 ${movingFile.name} 正在处理中，请稍等`)
          return false
        }
      }

      // check uploading file
      const uploadingFile = files.find(file => file.status === 'uploading')
      if (uploadingFile) {
        message.error(`文件 ${uploadingFile.name} 正在上传`)
        return false
      }

      // check merge and md5 file
      const mergeMd5File = files.find(file => file.status === 'mergeAndMd5')
      if (mergeMd5File) {
        message.error(`文件 ${uploadingFile.name} 正在合并与校验`)
        return false
      }

      if (item.required && item.isSupportWorkdir) {
        if (!item.workdir) {
          message.error(`${item.label} 未选择工作目录`)
          return false
        }
      }

      if (item.required && item.isSupportMaster) {
        if (!item.masterFile) {
          message.error(`${item.label} 未选择主文件`)
          return false
        }
      }

    } else if (item.required && !item.value && item.values.length === 0) {
      // required
      message.error(`请填写 ${item.label}`)
      return false
    }

    return true
  })
}
