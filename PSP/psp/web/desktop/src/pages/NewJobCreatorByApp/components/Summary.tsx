/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Observer } from 'mobx-react'
import { EditableText, Modal, Table, Button } from '@/components'
import { observable } from 'mobx'
import { SummaryStyle, AlertCheckboxAlertStyle } from './style'
import { Input, Checkbox, message } from 'antd'
import { userServer } from '@/server'
import { Context, useStore } from '../store'

type Props = {
  jobs: { jobName: string; mainFile: string }[]
  store: any
  onCancel: any
  onOk: (props: { jobName: string; mainFile: string }[]) => void
}

class Job {
  @observable mainFile: string
  @observable jobName: string
  @observable willAlert: boolean = false
  @observable inputVal: string

  constructor({ jobName, mainFile }) {
    this.mainFile = mainFile
    this.jobName = jobName
    this.inputVal = jobName
  }
}

export const Summary = observer(({ jobs, onCancel, onOk }: Props) => {
  const store = useStore()

  const state = useLocalStore(() => ({
    jobs: jobs.map(item => new Job(item)),
    get alertAll() {
      return this.jobs.every(job => job.willAlert)
    },
    notification_activated: false,
    setNA(bool) {
      this.notification_activated = bool
      if (bool) {
        this.jobs.forEach(v => (v.willAlert = true))
      }
    }
  }))

  useEffect(() => {
    // 为作业集名称提供默认值，默认是第一个作业的主文件名称
    store.updateData({ name: state?.jobs[0]?.jobName || ''})
  }, [])

  useEffect(() => {
    const fetch = async () => {
      const {
        data: { notification_activated }
      } = await userServer.checkWxBind('job')
      state.setNA(!!notification_activated)
    }

    fetch()
  }, [])

  return (
    <SummaryStyle>
      <div className='jobName'>
        <div className='label'>
          <span className='star'>*</span>
          <span>作业集名称：</span>
        </div>
        <Input
          value={store.data.name}
          onChange={e => {
            store.updateData({ name: e.target.value })
          }}
          placeholder='请输入'
        />
      </div>
      <div className='main'>
        <div className='table'>
          <Table
            columns={[
              {
                header: '作业名',
                props: {
                  flexGrow: 1
                },
                cell: {
                  props: {
                    dataKey: 'jobName'
                  },
                  render({ rowData }) {
                    const currentMainFile = rowData['mainFile']
                    return (
                      <div style={{ paddingRight: 50 }}>
                        <EditableText
                          filter = {item => {rowData['inputVal'] = item; return item}}
                          defaultValue={rowData['jobName']}
                          beforeConfirm={input => {
                            if (input === '') {
                              return '作业名称不能为空'
                            } else if (/\s/.test(input)) {
                              return '作业名称不能包含空格'
                            }

                            if (
                              state.jobs.find(
                                ({ jobName, mainFile }) =>
                                  mainFile !== currentMainFile &&
                                  jobName === input
                              )
                            ) {
                              return '作业名称已存在'
                            }
                            return true
                          }}
                          onConfirm={input => {
                            rowData.jobName = input
                          }}
                          onCancel={() => {
                            rowData['inputVal'] = rowData.jobName
                          }}
                        />
                      </div>
                    )
                  }
                }
              },
              {
                header: '主文件',
                props: {
                  flexGrow: 1
                },
                dataKey: 'mainFile'
              },
              {
                header: () => (
                  <Observer>
                    {() => (
                      <AlertCheckboxAlertStyle>
                        <Checkbox
                          disabled={state.notification_activated === false}
                          checked={state.alertAll}
                          onChange={({ target: { checked } }) => {
                            state.jobs.forEach(job => {
                              job.willAlert = checked
                            })
                          }}
                        />
                        全部通知
                      </AlertCheckboxAlertStyle>
                    )}
                  </Observer>
                ),
                props: {
                  width: 180
                },
                cell: {
                  render({ rowData }) {
                    return (
                      <Observer>
                        {() => (
                          <AlertCheckboxAlertStyle>
                            <Checkbox
                              disabled={state.notification_activated === false}
                              checked={rowData.willAlert}
                              onChange={({ target: { checked } }) => {
                                rowData.willAlert = checked
                              }}
                            />
                            通知
                          </AlertCheckboxAlertStyle>
                        )}
                      </Observer>
                    )
                  }
                }
              }
            ]}
            props={{
              data: state.jobs,
              height: 450
            } as any}
          />
        </div>
        {state.notification_activated === false && (
          <span className='note'>
            注：请前往【个人设置-我的微信】绑定微信后使用通知功能
          </span>
        )}
      </div>
      <Modal.Footer
        className='footer'
        OkButton={() => (
          <Button
            disabled={!store.data.name && '请输入作业集名称'}
            type={'primary'}
            onClick={() => {
              if (state.jobs.some(job => job.jobName !== job.inputVal)) {
                message.error('作业名发生了变化，请完成保存操作')
              } else if (state.jobs.some(job => job.jobName.trim() === '')) { 
                message.error('作业名不能为空，请修改') 
              } else {
                onOk(state.jobs)
              }
            }}>
            确认
          </Button>
        )}
        onCancel={onCancel}
      />
    </SummaryStyle>
  )
})

export const showSummary = (props: Omit<Props, 'onCancel' | 'onOk'>) =>
  Modal.show({
    title: '确认提交',
    width: 860,
    footer: null,
    content: ({ onCancel, onOk }) => (
      <Context.Provider value={props.store}>
        <Summary {...props} onCancel={onCancel} onOk={onOk} />
      </Context.Provider>
    )
  })
