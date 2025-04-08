import * as React from 'react'
import { message } from 'antd'

import { Point } from '@/domain/FileSystem'
import { Modal } from '@/components'

interface IProps {
  point: Point
  onConfirm?: any
  targets: string[]
}

export default class DeleteAction extends React.Component<IProps> {
  render() {
    const { onConfirm, targets, point } = this.props

    const children = React.Children.map(this.props.children, child =>
      React.cloneElement(child as React.ReactElement, {
        onClick: () => {
          Modal.showConfirm({
            title: '删除确认弹窗',
            content: '确认要删除选中的文件吗？（删除以后文件不可恢复）',
          }).then(() => {
            const names = targets.map(path => path.split(/[\\/]/).pop())
            if (names.length > 0) {
              const promise = point.service
                .delete({
                  paths: targets,
                })
                .then(() => {
                  message.success('删除成功')
                })
              onConfirm && onConfirm(promise, point)
            }
          })
        },
      })
    )

    return <>{children}</>
  }
}
