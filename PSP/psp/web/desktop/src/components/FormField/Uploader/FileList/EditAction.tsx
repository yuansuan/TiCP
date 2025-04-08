import * as React from 'react'
import { message } from 'antd'
// import { Editor } from '@/components/FileSystem'
import { Modal } from '@/components'
import { observable } from 'mobx'
import { Http } from '@/utils'

interface IProps {
  title: string
  node: any
  readOnly?: boolean
  onClick?: any
}

export const isTextFile = file => file.isFile && file.isText

export default class EditAction extends React.Component<IProps> {
  @observable fileSize = null

  render() {
    const { title, readOnly, node, onClick } = this.props
    const children = React.Children.map(this.props.children, child =>
      React.cloneElement(child as React.ReactElement, {
        onClick: () => {
          onClick && onClick()

          if (!node) {
            return
          }

          const operation = readOnly ? '查看' : '编辑'
          if (
            !isTextFile({
              ...node,
              isFile: !node.isDir,
              isText: node.is_text
            })
          ) {
            message.error(`只能${operation}文本文件`)
            return
          }

          const service = {
            view({ path, offset, len }) {
              return Http.get('/file/content', {
                params: { path, offset, len }
              }).then(res => res.data.content)
            },
            edit({ path, content }) {
              return Http.put(
                '/file/edit',
                { path, content },
                { formatErrorMessage: msg => `编辑失败：${msg}` }
              )
            },
            get(path: string) {
              return Http.get('/file/detail', {
                params: {
                  paths: path
                },
                formatErrorMessage: msg => `获取文件信息失败`
              })
            }
          }

          this.fileSize = node.size

          this.openEditor({
            title,
            readOnly,
            viewContent: service.view.bind(service),
            saveContent: service.edit.bind(service),
            file: node,
            getFile: service.get.bind(service)
          })
        }
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
    getFile
  }) => {
    Modal.show({
      title,
      width: 800,
      bodyStyle: {
        height: 600,
        padding: 10
      },
      onCancel: closeModal => {
        if (!readOnly) {
          Modal.showConfirm({
            content: '确认要取消保存吗？（取消保存会丢失已编辑的内容）'
          }).then(() => {
            closeModal()
          })
        } else {
          closeModal()
        }
      },
      content: ({ onCancel, onOk }) => {
        // return (<Editor
        //   viewContent={viewContent}
        //   saveContent={saveContent}
        //   path={file.path}
        //   fileSize={this.fileSize || file.size}
        //   readOnly={readOnly}
        //   onCancel={onCancel}
        //   onOk={onOk}
        //   onBeforeRefresh={async () => {
        //     if (getFile) {
        //       const res = await getFile(file.path)
        //       this.fileSize = file.size
        //       return res.data.files[0].size
        //     } else {
        //       this.fileSize = file.size
        //       return file.size
        //     }
        //   }}
        // />)
      },
      footer: null
    })
  }
}
