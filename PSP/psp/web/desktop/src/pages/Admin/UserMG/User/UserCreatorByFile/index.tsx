import * as React from 'react'
import { message, Input, Button } from 'antd'
import { observer } from 'mobx-react'
import { observable, action, computed } from 'mobx'

import { UserList, User } from '@/domain/UserMG'
import { UserCreatorWrapper, FooterWrapper } from './style'
import { Modal } from '@/components'
import TextArea from 'antd/lib/input/TextArea'

interface IProps {
  user: User
  onOk?: () => void
}

@observer
export default class UserCreatorByFile extends React.Component<IProps> {
  @observable content = ''
  @observable file
  @observable fileList = []
  @observable adding = false
  @observable isShow = false
  @observable isCreate = false
  @observable successList = []
  @observable isAddedList = []
  @observable notExistList = []

  @action
  updateIsCreate = flag => (this.isCreate = flag)
  @action
  updateAdding = flag => (this.adding = flag)
  @action
  updateContent = content => (this.content = content)
  @action
  updateFile = file => (this.file = file)
  @action
  updateIsShow = flag => (this.isShow = flag)

  public async componentDidMount() {
    await UserList.fetchDisabled()
  }

  private onOk = () => {
    const { onOk } = this.props
    if (this.file && this.isCreate) {
      this.updateAdding(true)

      onOk && onOk()
      this.updateAdding(false)
    } else {
      Modal.showConfirm({
        title: '关闭',
        content: '请选择文件后点击开始创建按钮以添加用户',
      })
    }
  }
  delete(index) {
    this.fileList.splice(index, 1)
  }
  onCreate = () => {
    this.successList = []
    this.isAddedList = []
    this.notExistList = []

    if (this.file === undefined) {
      message.error('请选择文件')
      return
    }
    const sysUsers = UserList.disabledUsers
    const allUsers = sysUsers.concat(UserList.enabledUsers)
    const newlist = this.fileList.map(n => {
      const userName = n.split(/\s+/)[0]
      const adduser = allUsers.filter(n => userName === n.name)[0]

      const id = adduser ? adduser.id : -1
      let isadded
      if (id !== -1) {
        isadded = sysUsers.filter(n => n === adduser)[0] ? false : true
      }

      return { id: id, name: userName, isAdded: isadded }
    })

    UserList.addAll(newlist).then(res => {
      //通过result结果筛选出true或false对应的User
      const getUserByUserStatus = status =>
        res.data.filter(n => n.result === status)

      getUserByUserStatus(true).map(n => {
        this.successList.push(n.userName)
      })

      getUserByUserStatus(false).map(n => {
        if (n.isAdded) {
          this.isAddedList.push(n.userName)
        }
        if (n.reason) {
          this.notExistList.push(n.userName)
        }
      })
    })

    this.updateIsShow(true)
    this.updateIsCreate(true)
  }

  submitFile = async e => {
    const fileContent = await this.readFile(e.target.files[0])
    this.updateContent(fileContent)
    this.fileList = this.content.split(/\r\n|\r|\n/).filter(s => s && s.trim())
    this.updateIsCreate(false)
    this.updateIsShow(false)
  }

  readFile(file) {
    if (file) {
      return new Promise((resolve, reject) => {
        var reader = new FileReader()
        reader.readAsText(file, 'UTF-8')
        this.updateFile(file)
        reader.onload = evt => {
          resolve(evt.target.result)
        }
        reader.onerror = () => {
          reject('failed')
        }
      })
    } else {
      this.updateFile(file)
      return ' '
    }
  }

  @computed
  get result() {
    if (this.isShow) {
      let res
      if (this.successList.length !== 0) {
        const success = '用户' + this.successList + '添加成功'
        res = success
      }
      if (this.isAddedList.length !== 0) {
        const failure = '用户' + this.isAddedList + '已经存在，添加失败'

        res = res ? res + '\n' + failure : failure
      }
      if (this.notExistList.length !== 0) {
        const exist = '用户' + this.notExistList + '不存在，无法添加'
        res = res ? res + '\n' + exist : exist
      }
      return res
    } else {
      return ''
    }
  }

  render() {
    return (
      <UserCreatorWrapper>
        <div className='header'>
          <div className='upload-wrap anticon' nv-file-drop=''>
            <Input
              type='file'
              accept='.csv'
              onChange={this.submitFile}
              className='file-ele'
            />
            <div className='file-open'>
              <em className='icon icon-upload'></em>&nbsp;请选择csv文件
            </div>
          </div>

          <Button onClick={this.onCreate} className='create'>
            开始创建
          </Button>
        </div>

        <div className='editorMain'>
          <div className='content'>
            {this.fileList &&
              this.fileList.map((item, index) => (
                <li key={`key${index}`} className='contentList'>
                  {item}
                  <a
                    onClick={() => {
                      this.delete(index)
                    }}>
                    删除
                  </a>
                </li>
              ))}
          </div>

          <TextArea className='result' value={this.result} />
        </div>

        <FooterWrapper>
          <div className='footerMain'>
            <Button disabled={this.adding} type='primary' onClick={this.onOk}>
              {this.adding ? '关闭中...' : '关闭'}
            </Button>
          </div>
        </FooterWrapper>
      </UserCreatorWrapper>
    )
  }
}
