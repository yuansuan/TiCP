/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState, useRef } from 'react'
import { Switch, InputNumber } from 'antd'
import { observer } from 'mobx-react'
import { message } from 'antd'
import { machineSetting } from '@/domain/Visual'
import { ConfigWrapper } from './style'

const UserTaskConfig = observer(() => {
  const { user_task, default_vm_task_number } = machineSetting
  const [sessionNumber, setSessionNumber] = useState(0)
  const [taskNumber, setTaskNumber] = useState(0)
  const [checked, setChecked] = useState(user_task.enable)
  const checkedRef = useRef(checked)

  useEffect(() => {
    machineSetting.fetch()
  }, [])

  useEffect(() => {
    setSessionNumber(user_task.number)
    setTaskNumber(default_vm_task_number)
  }, [default_vm_task_number, user_task])

  const updateMachineSetting = async () => {
    try {
      await machineSetting.update({
        user_task: {
          enable: checkedRef.current,
          number: +sessionNumber,
        },
        default_vm_task_number: +taskNumber,
      })
    } finally {
      await machineSetting.fetch()
    }
  }
  const onChecked = async value => {
    setChecked(value)
    checkedRef.current = value
    await updateMachineSetting()
  }

  const onBlur = async value => {
    if (
      checked === machineSetting.user_task.enable &&
      sessionNumber === machineSetting.user_task.number &&
      taskNumber === machineSetting.default_vm_task_number
    )
      return
    await updateMachineSetting()
    message.success('操作成功！')
  }
  return (
    <ConfigWrapper>
      <div className='item-row'>
        <div>
          <Switch
            style={{ marginLeft: 30 }}
            checked={checked}
            checkedChildren='开启用户任务数'
            unCheckedChildren='关闭用户任务数'
            onChange={value => onChecked(value)}
          />
        </div>
        <InputNumber
          style={{ margin: '0 20px' }}
          size='small'
          disabled={!checked}
          min={1}
          step={1}
          value={sessionNumber || machineSetting.user_task.number}
          precision={0}
          parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
          onChange={value => setSessionNumber(value)}
          onBlur={e => onBlur(e.target.value)}
        />
      </div>
      <div className='item-row'>
        <span className='label'>虚拟机任务数</span>
        <InputNumber
          style={{ margin: '0 20px' }}
          size='small'
          min={1}
          step={1}
          value={taskNumber || machineSetting.default_vm_task_number}
          precision={0}
          parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
          onChange={value => setTaskNumber(value)}
          onBlur={e => onBlur(e.target.value)}
        />
      </div>
    </ConfigWrapper>
  )
})

export default UserTaskConfig
