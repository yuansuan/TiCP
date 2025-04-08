import * as React from 'react'

import { Point } from '@/domain/FileSystem'
import { currentUser } from '@/domain'
import { message } from 'antd'

interface IProps {
  point: Point
  onClick?: any
  targets: string[]
  beforeDownload?: () => boolean | Promise<boolean>
}

export default class DownloadAction extends React.Component<IProps> {
  private onDownload = () => {
    const { point, targets } = this.props

    point.service
      .preDownload({ paths: targets, userId: currentUser.id })
      .then(res => res.data.token)
      .then(token =>
        point.service.download({
          token,
          userId: currentUser.id,
        })
      )
      .catch(err => {
        message.error(err.message)
      })
  }

  render() {
    const { onClick, beforeDownload } = this.props

    const children = React.Children.map(this.props.children, child =>
      React.cloneElement(child as React.ReactElement, {
        onClick: () => {
          onClick && onClick()
          if (beforeDownload) {
            const res = beforeDownload()
            if (res instanceof Promise) {
              res.then(this.onDownload)
            } else if (res) {
              this.onDownload()
            }
          } else {
            this.onDownload()
          }
        },
      })
    )

    return <>{children}</>
  }
}
