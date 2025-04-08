import * as React from 'react'
import { message } from 'antd'
import { Editor } from '@/components/FileSystem'
import { Point } from '@/domain/FileSystem'
import { Modal } from '@/components'
import { observable } from 'mobx'

interface IProps {
  title: string
  point: Point
  path: string
  readOnly?: boolean
  onClick?: any
}

export const isTextFile = file => file.isFile && file.isText

export default class EditAction extends React.Component<IProps> {
  @observable fileSize = null

  render() {
    const { title, readOnly, point, path, onClick } = this.props
    const children = React.Children.map(this.props.children, child =>
      React.cloneElement(child as React.ReactElement, {
        onClick: () => {
          onClick && onClick()

          const node = point.filterFirstNode(item => item.path === path)
          if (!node) {
            return
          }

          const operation = readOnly ? '查看' : '编辑'
          if (
            !isTextFile({
              ...node,
              isText: node.is_text,
            })
          ) {
            message.error(`只能${operation}文本文件`)
            return
          }

          const { service } = point

          this.fileSize = node.size

          this.openEditor({
            title,
            readOnly,
            viewContent: service.view.bind(service),
            saveContent: service.edit.bind(service),
            file: node,
            getFile: service.get.bind(service),
          })
        },
      })
    )

    return <>{children}</>
  }

  private openEditor = ({
    title,
    readOnly = false,
    viewContent,
    saveContent,
    file,
    getFile,
  }) => {
    Modal.show({
      title,
      width: 800,
      bodyStyle: {
        height: 600,
        padding: 10,
      },
      onCancel: closeModal => {
        if (!readOnly) {
          Modal.showConfirm({
            content: '确认要取消保存吗？（取消保存会丢失已编辑的内容）',
          }).then(() => {
            closeModal()
          })
        } else {
          closeModal()
        }
      },
      content: ({ onCancel, onOk }) => {
        return (<Editor
          viewContent={viewContent}
          saveContent={saveContent}
          path={file.path}
          fileSize={this.fileSize || file.size}
          readOnly={readOnly}
          onCancel={onCancel}
          onOk={onOk}
          onBeforeRefresh={async () => {
            if (getFile) {
              const res = await getFile(file.path)
              this.fileSize = file.size
              return res.data.files[0].size
            } else {
              this.fileSize = file.size
              return file.size
            }
          }}
        />)
      },
      footer: null,
    })
  }
}
