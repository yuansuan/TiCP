import * as React from 'react'
import { Input, message } from 'antd'
import { observer } from 'mobx-react'
import { observable, action } from 'mobx'

import { Http } from '@/utils'
import { Button } from '@/components'
import { FavoritePoint } from '@/domain/FileSystem'
import { StyledFavoritePopover } from './style'

interface IProps {
  favoritePoint: FavoritePoint
  path: string
  hide: () => void
}

@observer
export default class FavoritePopover extends React.Component<IProps> {
  @observable name = ''
  @action
  updateName = name => (this.name = name)

  private onCancel = () => {
    const { hide } = this.props

    hide()
    this.updateName('')
  }

  private onSave = () => {
    const { hide, path, favoritePoint } = this.props

    if (!this.name) {
      message.error('名称不能为空')
    } else {
      Http.post('/file/favorite', {
        name: this.name,
        path,
      }).then(() => {
        message.success('添加成功')
        favoritePoint.fetch()
      })

      hide()
      this.updateName('')
    }
  }

  private onKeyDown = e => {
    if (e.keyCode === 13) {
      this.onSave()
    }
  }

  render() {
    return (
      <StyledFavoritePopover>
        <div>名称</div>
        <div>
          <Input
            value={this.name}
            onChange={e => this.updateName(e.target.value)}
            onKeyDown={this.onKeyDown}
          />
        </div>
        <div className='toolbar'>
          <div className='main'>
            <Button onClick={this.onCancel}>取消</Button>
            <Button type='primary' onClick={this.onSave}>
              保存
            </Button>
          </div>
        </div>
      </StyledFavoritePopover>
    )
  }
}
