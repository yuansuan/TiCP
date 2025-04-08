import React, { useEffect, useState } from 'react'
import { Modal } from '@/components'
import { Button, Form, Input, message, Radio, Space } from 'antd'
import { observer } from 'mobx-react'
import { showDirSelector,showFailure } from '@/components/NewFileMGT'
import { useLocalStore } from 'mobx-react-lite'
import styled from 'styled-components'
import { Http } from '@/utils'
import { serverFactory } from '@/components/NewFileMGT/store/common'
import { newBoxServer } from '@/server'
import { currentUser } from '@/domain'

const StyledLayout = styled.div`
  .body {
    width: 450px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    align-item: center;
    border-bottom: 1px solid #ccc;
  }
  > .footer {
    margin-top: 10px;
  }
`

export const ShareFileContent = observer(({ onCancel, onOk, ...options }) => {
  const { isdir, name, path, share_type, size, type } = options
  const server = serverFactory(newBoxServer)

  const state = useLocalStore(() => ({
    userList: [],
    setUserList(list) {
      this.userList = list
    },
    checkUsers: [],
    setCheckedUsers(users) {
      this.checkUsers = users
    },
    get okDisable() {
      return state.setCheckedUsers.length > 0
    }
  }))
  const [value, setValue] = useState(1)
  const onChange = async e => {
    setValue(e.target.value)
  }

  const handlerFile = async () => {
    if (value === 1) {
      server.download([path], [!isdir], [size],true)
    } else {
      const dst_path = await showDirSelector()
      const current_path = path.substring(0, path.lastIndexOf('/'))
      
      if (!path) {
        return message.error('请选择保存目标路径！')
      } else {
        const url = share_type === 1 ? '/storage/copy' : '/storage/link'
        const HttpMethod = share_type === 1 ? Http.put : Http.post
        let  isOverwrite = false
        // check duplicate
        const targetDir = await server.fetch(dst_path)
        const rejectedNodes = []
        const resolvedNodes = []
        if (targetDir.getDuplicate({ id: undefined, name: name })) {
          rejectedNodes.push(options)
        } else {
          resolvedNodes.push(options)
        }
        if (rejectedNodes.length > 0) {
          const  coverNodes = await showFailure({
            actionName: '保存',
            items: rejectedNodes
          })
          if (coverNodes.length > 0) {
            isOverwrite = true
          } else {
            return onOk()
          }
        }
        await HttpMethod(url, {
            cross: true,
            current_path,
            dst_path: dst_path.replace(/^\./, currentUser.name),
            is_cloud: false,
            overwrite: isOverwrite,
            src_dir_paths: isdir ? [path] : [],
            src_file_paths: isdir ? [] : [path]
          }).then(res => {
            if (res.success) {
              if(share_type === 1){
                message.info('文件保存中，如文件较大可能耗时较长，请耐心等待')
              }else {
                message.success('文件保存完成')
              }
              onOk()
            }
          })
       
      }
    }
  }

  return (
    <StyledLayout>
      <div className='body'>
        <Radio.Group onChange={onChange} value={value}>
          <Space>
            <Radio value={1}>下载至本地</Radio>
            <Radio value={2}>保存至我的文件</Radio>
          </Space>
        </Radio.Group>
      </div>
      <Modal.Footer
        className='footer'
        onCancel={onCancel}
        OkButton={
          <Button
            type='primary'
            disabled={state.okDisable}
            onClick={() => handlerFile()}>
            确认
          </Button>
        }
      />
    </StyledLayout>
  )
})
