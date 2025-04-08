/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Switch, message } from 'antd'
import { Http } from '@/utils'
import IEditableText from './IEditableText'

import { StyledLayout } from './style'

const editConfig = (config: { key: string; value: number | string }) => {
  return Http.put('/company/config', config)
}

const SystemSetting = observer(function SystemSetting() {
  const state = useLocalStore(() => ({
    autoBoost: null,
    setAutoBoost(value: boolean) {
      this.autoBoost = value
    },
    timeBeforeBoost: null,
    setTimeBeforeBoost(value: number) {
      this.timeBeforeBoost = value
      if (value === null) this.setAutoBoost(false)
    },
    downloadSpeedLimit: null,
    setDownloadSpeedLimit(value: number) {
      this.downloadSpeedLimit = value
    },
    uploadSpeedLimit: null,
    setUploadSpeedLimit(value: number) {
      this.uploadSpeedLimit = value
    },
    cloudType: null,
    setCloudType(value: string) {
      this.cloudType = value
    },
  }))

  const autoBoostOnChange = (value: boolean) => {
    if (value === true) {
      editConfig({
        key: 'timeBeforeBoost',
        value: 1,
      }).then(data => {
        if (data.success) {
          state.setAutoBoost(true)
          state.setTimeBeforeBoost(1)
          message.success('启用自动爆发')
        }
      })
    } else if (value === false) {
      editConfig({
        key: 'timeBeforeBoost',
        value: null,
      }).then(data => {
        if (data.success) {
          state.setAutoBoost(false)
          state.setTimeBeforeBoost(null)
          message.success('取消自动爆发')
        }
      })
    }
  }

  const timeBeforeBoostSetValue = (value: number) => {
    editConfig({
      key: 'timeBeforeBoost',
      value: value,
    }).then(data => {
      if (data.success) {
        state.setTimeBeforeBoost(value)
        if (value === null) {
          state.setAutoBoost(false)
          message.success('取消自动爆发')
        } else {
          message.success(`设置自动爆发前等待时间为 ${value} 小时`)
        }
      }
    })
  }

  const downloadSpeedLimitSetValue = (value: number) => {
    editConfig({
      key: 'downloadSpeedLimit',
      value: value,
    }).then(data => {
      if (data?.success) {
        state.setDownloadSpeedLimit(value)
        message.success(
          value === null
            ? '取消云端任务下载限速'
            : `云端任务下载限速为 ${value} MB/s`
        )
      }
    })
  }

  const uploadSpeedLimitSetValue = (value: number) => {
    editConfig({
      key: 'uploadSpeedLimit',
      value: value,
    }).then(data => {
      if (data.success) {
        state.setUploadSpeedLimit(value)
        message.success(
          value === null
            ? '取消云端任务上传限速'
            : `云端任务上传限速为 ${value} MB/s`
        )
      }
    })
  }


  return (
    <StyledLayout>
      {state.cloudType === 'mixed' && (
        <section>
          <h1>爆发设置</h1>
          <div className='section-bottom'>
            <div className='row'>
              自动爆发：
              <Switch checked={state.autoBoost} onChange={autoBoostOnChange} />
            </div>
            {state.autoBoost && (
              <div className='row'>
                等待时间：
                <IEditableText
                  unit='小时'
                  value={state.timeBeforeBoost}
                  setValue={timeBeforeBoostSetValue}
                />
              </div>
            )}
          </div>
        </section>
      )}

      <section>
        <h1>网络限速</h1>
        <div className='section-bottom'>
          <div className='row'>
            云端任务下载限速：
            <IEditableText
              unit='MB/s'
              value={state.downloadSpeedLimit}
              setValue={downloadSpeedLimitSetValue}
            />
          </div>
          <div className='row'>
            云端任务上传限速：
            <IEditableText
              unit='MB/s'
              value={state.uploadSpeedLimit}
              setValue={uploadSpeedLimitSetValue}
            />
          </div>
        </div>
      </section>
    </StyledLayout>
  )
})

export default SystemSetting
