/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Icon } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { message, Form, Input } from 'antd'
import { StyledLayout } from './style'
import { currentUser } from '@/domain'
import { validateName } from '@/utils/Validator'
import { userServer } from '@/server'
import { runInAction } from 'mobx'

export const Name = observer(function Name() {
  const store = useLocalStore(() => ({
    editMode: false,
    setEditMode(flag) {
      this.editMode = flag
    },
    realName: currentUser.real_name,
    setRealName(name) {
      this.realName = name
    }
  }))
  const [formRef] = Form.useForm()

  const onFinish = async values => {
    const real_name = values['realname'].trim()
    await userServer.updateRealName(real_name)
    runInAction(() => {
      currentUser.update({
        real_name
      })
    })
    store.setEditMode(false)
    message.success('修改成功')
  }

  const onSave = async () => {
    formRef.submit()
  }

  const onEdit = () => {
    store.setEditMode(true)
  }
  return (
    <StyledLayout>
      <Form
        hideRequiredMark
        form={formRef}
        onFinish={onFinish}
        initialValues={{ realname: currentUser.name }}>
        <Form.Item label='用户名'>
          {!store.editMode ? (
            <>
              <span className='text'>{currentUser.name}</span>
              {/* <div className='right-edit'>
                <Icon className='edit' type='rename' onClick={onEdit} />
              </div> */}
            </>
          ) : (
            <>
              <Form.Item
                name='realname'
                rules={[
                  {
                    required: true,
                    transform: value => value.trim(),
                    message: '姓名不能为空'
                  },
                  {
                    pattern: validateName.reg,
                    message: "姓名只能包含中文、英文、空格和 , . ' -"
                  }
                ]}>
                <Input
                  placeholder='请输入姓名'
                  onFocus={e => e.target.select()}
                />
              </Form.Item>
              <div className='right-confirm'>
                <Icon className='ok' type='define' onClick={onSave} />
                <Icon
                  className='cancel'
                  type='cancel'
                  onClick={() => {
                    store.setEditMode(false)
                    formRef.setFieldsValue({ realname: currentUser.name })
                  }}
                />
              </div>
            </>
          )}
        </Form.Item>
      </Form>
    </StyledLayout>
  )
})
