import React,{useEffect} from 'react'
import { Modal } from '@/components'
import { Button, message} from 'antd'
import { observer } from 'mobx-react'
import {useLocalStore} from 'mobx-react-lite'
import styled from 'styled-components'
import OrganizationTree from './OrganizationTree'
import {Http} from '@/utils'
import {currentUser} from '@/domain'
const StyledLayout = styled.div`
 
  .body {
    width: 450px;
    min-height: 300px;
    max-height: 500px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    align-item:center;
    border-bottom: 1px solid #ccc;
    .title {
      display: block;
      margin-bottom: 5px;
      font-size: 15px;
      font-weight: 500;
    }
  }
  >.footer{
    margin-top:10px;
  }
`

const ShareFileContent = observer(({ onCancel, onOk,actType,selectedNodes }) => {
  const state = useLocalStore(() => ({
    userList: [],
    setUserList(list) {
      this.userList= list
    },
    checkUsers: [],
    setCheckedUsers(users) {
      this.checkUsers = users
    },
    get okDisable() {
      return state.setCheckedUsers.length >0
    },
    get actName() {
      return actType === 'send' ? '发送': '分享'
    }
  }))
  
  async function getUserOptionList(name) {
    const {data} =await Http.get('/user/optionList',{
      params:{
        filterName:name
      }
    })
    const newData = data?.filter(item => item.title !==  currentUser.name)
    state.setUserList(newData || [])
  }
  useEffect(() => {
    getUserOptionList('')
  },[])



  const handleCheck = (keys) => {
    state.setCheckedUsers(keys)
  };

  const onShareOrSendFile =() => {
    if(state.checkUsers.length && selectedNodes.length){
      Http.post('/storage/share/send',{
        share_file_path: selectedNodes[0].path,
        share_type: actType === 'send' ? 1 : 2,
        share_user_list:state.checkUsers
      }).then(res => {
        if(res.success){
          message.success(`${state.actName}成功`)
          onOk()
        }
      })
    }else { 
      message.warn('请至少选择一名成员')
    }
  }
  return (
    <StyledLayout>
      <div className='body'>
        <div className='title'>*{state.actName}人员</div>
        <OrganizationTree treeData={state.userList} onCheck={handleCheck} refresh={getUserOptionList} />
      </div>
      <Modal.Footer
        className='footer'
        onCancel={onCancel}
        OkButton={
          <Button type='primary' disabled={state.okDisable} onClick={onShareOrSendFile}>
            确认
          </Button>
        }
      />
    </StyledLayout>
  )
})

export const showShareFile = options => {
  const selectedFileName = options?.selectedNodes[0].name
  return Modal.show({
    title: `${options?.actType ==='send' ? "发送" : '分享'}文件（${selectedFileName}）`,
    footer: null,
    width: 500,
    content: ({ onCancel, onOk }) => (
      <ShareFileContent {...options} onCancel={onCancel} onOk={onOk} />
    )
  })
}
