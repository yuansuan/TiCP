import * as React from 'react'
import { message } from 'antd'

import { Point, RootPoint } from '@/domain/FileSystem'
import { Modal } from '@/components'
import CompressInfo from './CompressInfo'
import { ZipType } from '@/utils/const'
import { eventEmitter, IEventData } from '@/utils'

interface IProps {
  points: RootPoint[]
  point: Point
  onClick?: any
  targets: string[]
}

const isFileExist = async (point, fileName) => {
  return point.service.exist([fileName]).then(res => {
    if (res[0]) {
      return Modal.showConfirm({
        content: `${fileName} 已存在，是否覆盖？`,
      })
    } else {
      return true
    }
  })
}

export default class CompressAction extends React.Component<IProps> {
  render() {
    const { onClick, targets, point, points } = this.props
    const parentPath =
      targets.length > 0 && targets[0]?.replace(/[\\/][^\\/]+$/, '')

    const children = React.Children.map(this.props.children, child =>
      React.cloneElement(child as React.ReactElement, {
        onClick: () => {
          onClick && onClick()
          Modal.show({
            title: '压缩',
            bodyStyle: { height: 240 },
            content: ({ onCancel, onOk }) => (
              <CompressInfo
                points={points}
                parentPath={parentPath}
                zipName={
                  targets.length > 1 ? null : targets[0].split(/[\\/]/).pop()
                }
                onCancel={onCancel}
                onOk={onOk}
              />
            ),
            footer: null,
          }).then(({ zipName, compressType, destPath }) => {
            let _zipName = `${destPath}/${zipName}.${ZipType[compressType]}`
            let names = targets.map(path => path.split(/[\\/]/).pop())

            isFileExist(point, _zipName).then(res => {
              point.service
                .compress({
                  path: parentPath,
                  names,
                  compressType,
                  zipName: _zipName,
                })
                .then(() => {
                  // 监听消息
                  eventEmitter.once(
                    `COMPRESS_FILE_${_zipName}`,
                    (obj: IEventData) => {
                      if (obj.message.success) {
                        point.service.fetch(destPath).then(() => {})
                      }
                    }
                  )

                  message.success(`开始压缩文件 ${_zipName}，请稍后...`)
                })
                .catch(err => message.error(err.message))
            })
          })
        },
      })
    )

    return <>{children}</>
  }
}
