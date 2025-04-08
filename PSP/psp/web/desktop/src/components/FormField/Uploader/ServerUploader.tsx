import { action, observable } from 'mobx'
import { observer } from 'mobx-react'
import * as React from 'react'
import { Subject } from 'rxjs'
import styled from 'styled-components'

// import  FileSystem  from '@/components/FileSystem'
import { ActionType } from '@/components/FileSystem/Toolbar'
import { Store } from '@/domain/FileSystem'
import { untilDestroyed } from '@/utils/operators'
import { Modal } from '@/components'
import { currentUser } from '@/domain'
import { RootPoint } from '@/domain/FileSystem'

export const ServerFileChooser = styled.div`
  height: 100%;
  display: flex;
  flex-direction: column;

  > .body {
    height: calc(100% - 42px);
  }

  > .footer {
    padding: 0 10px 10px 0;
    background-color: white;
  }
`

interface IProps {
  onUpload: (files: any[]) => void
  children: any
}

@observer
export default class ServerUploaderUI extends React.Component<IProps> {
  @observable public selectedKeys = []

  public selectedKeys$ = new Subject<string[]>()
  @action
  public updateSelectedKeys = keys => (this.selectedKeys = keys)

  public componentDidMount() {
    this.selectedKeys$
      .pipe(untilDestroyed(this))
      .subscribe(this.updateSelectedKeys)
  }

  public render() {
    const { children } = this.props
    return <>{children(this.uploadServerFiles)}</>
  }

  private onOk = async keys => {
    this.props.onUpload(keys.map(key => Store.get(key)))
  }

  private uploadServerFiles = () => {
    Modal.show({
      title: '请选择服务器文件',
      width: 1000,
      footer: null,
      bodyStyle: { height: 600, padding: 0 },
      content: ({ onCancel, onOk }) => (
        <ServerFileChooser>
          <div key='body' className='body'>
            {/* <FileSystem
              selectedKeys$={this.selectedKeys$}
              points={currentUser.mountList
                .filter(item => item.id !== 'favorites')
                .map(
                  point =>
                    new RootPoint({
                      pointId: point.id,
                      path: point.path,
                      name: point.name,
                    })
                )}
              toolbar={{
                includes: [ActionType.newFolder],
              }}
            /> */}
          </div>
          <div key='footer' className='footer'>
            <Modal.Footer
              onCancel={() => {
                this.updateSelectedKeys([])
                onCancel()
              }}
              onOk={() => {
                this.onOk(this.selectedKeys)
                this.updateSelectedKeys([])
                onOk()
              }}
            />
          </div>
        </ServerFileChooser>
      )
    })
  }
}
