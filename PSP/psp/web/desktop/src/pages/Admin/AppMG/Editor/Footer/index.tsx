import * as React from 'react'
import { observer } from 'mobx-react'
import { message } from 'antd'
import { Button } from '@/components'
import { App } from '@/domain/Applications'
import { history, getUrlParams } from '@/utils'
import { inject } from '@/pages/context'
import { DiffInfo } from '@/components'
import { Modal } from '@/components'
import checkForm from '@/utils/checkForm'
import TestResult from './TestResult'
import SaveAs from './SaveAs'
import { Wrapper } from './style'

interface IProps {
  app?: App
  formModel: any
  className?: string
  disabled: boolean
  isRemote?: boolean
  updateDisabled: (disabled: boolean) => void
  uploadToken?: string
  getScriptData: () => string
}

@inject(({ app, uploadToken }) => ({
  app,
  uploadToken
}))
@observer
export default class Footer extends React.Component<IProps> {
  unblockHistory = null
  blockHistoryFlag = true
  tab = ''
  componentDidMount() {
    // history change block

    const query = getUrlParams()
    this.tab = query?.isRemote === 'true' ? 'remote' : 'normal'

    this.unblockHistory = history.block(() => {
      if (this.blockHistoryFlag) {
        return '离开页面将不会保留当前工作，确认要离开页面吗？'
      }
      return undefined
    })
  }

  componentWillUnmount() {
    this.unblockHistory && this.unblockHistory()
  }

  private pushHistory(path) {
    this.blockHistoryFlag = false
    history.push(path)
    window.localStorage.setItem('CURRENTROUTERPATH', path)
  }

  private onCancel = () => {
    const { app } = this.props
    const diff = app.diff()

    let modalPromise = null

    if (diff) {
      modalPromise = Modal.show({
        title: '确认取消变更',
        width: 600,
        bodyStyle: { height: 250 },
        content: (
          <DiffInfo
            diff={diff}
            script={{
              old: app.script,
              new: app.scriptData
            }}
          />
        )
      })
    } else {
      modalPromise = Modal.showConfirm({
        content: '确定要取消模版编辑吗？'
      })
    }

    modalPromise.then(() => {
      app.reset()

      this.pushHistory(`/sys/template?tab=${this.tab}`)
    })
  }

  private beforeSave = () =>
    new Promise((resolve, reject) => {
      const { app } = this.props

      const invalidSection = app.subForm.sections.find(
        section =>
          !section.name ||
          section.fields.length === 0 ||
          section.computedEditing
      )

      let supportWorkdirFieldOfNums = 0

      block_out: for (let i = 0; i < app.subForm.sections.length; i++) {
        let section = app.subForm.sections[i]
        let fields = section.fields

        for (let j = 0; j < fields.length; j++) {
          if (fields[j].isSupportWorkdir) {
            supportWorkdirFieldOfNums++
            if (supportWorkdirFieldOfNums >= 2) {
              break block_out
            }
          }
        }
      }

      if (invalidSection) {
        if (!invalidSection.name) {
          message.error('存在未命名的 section')
        } else if (invalidSection.fields.length === 0) {
          message.error(`${invalidSection.name} section 必须至少存在一个组件`)
        } else {
          message.error(`请完成 ${invalidSection.name} section 的编辑`)
        }
        reject()
      } else if (supportWorkdirFieldOfNums >= 2) {
        message.error('1个以上 Field 开启支持选择工作目录')
        reject()
      } else {
        resolve(null)
      }
    })

  private onSave = ({ publish = false } = {}) =>
    this.beforeSave().then(() => {
      const { app, getScriptData } = this.props

      // 设置 scriptData
      app.setScriptData(getScriptData())

      const diff = app.diff()

      if (diff) {
        Modal.show({
          title: '请确认变更',
          width: 800,
          bodyStyle: { height: 400, padding: 10, overflowY: 'auto' },
          content: (
            <DiffInfo
              diff={diff}
              script={{
                old: app.script,
                new: app.scriptData
              }}
            />
          )
        }).then(() => {
          if (publish) {
            this.saveAndPublish()
          } else {
            this.save()
          }
        })
      } else {
        message.warn('暂无可保存变更')
      }
    })

  private save = () => {
    const { app, updateDisabled } = this.props

    updateDisabled(true)
    return app
      .save()
      .then(() => {
        message.success('保存成功')
        this.pushHistory(`/sys/template?tab=${this.tab}`)
      })
      .finally(() => updateDisabled(false))
  }

  private saveAndPublish = () => {
    const { app, updateDisabled } = this.props

    updateDisabled(true)
    return app
      .save()
      .then(app.publish)
      .then(() => {
        message.success('保存并发布成功')

        this.pushHistory(`/sys/template?tab=${this.tab}`)
      })
      .finally(() => updateDisabled(false))
  }

  private saveAs = () =>
    this.beforeSave().then((resolve: any) => {
      const { app } = this.props
      Modal.show({
        title: '另存为模版',
        width: 300,
        bodyStyle: { height: 170 },
        content: ({ onOk, onCancel }) => (
          <SaveAs
            app={app}
            getScriptData={this.props.getScriptData}
            defaultName={`${app.version}.x`}
            onOk={onOk}
            onCancel={onCancel}
          />
        ),
        footer: null
      }).then(() => {
        this.pushHistory(`/sys/template?tab=${this.tab}`)
      })
    })

  private test = () => {
    const { formModel, app, updateDisabled, uploadToken, getScriptData } =
      this.props

    if (!checkForm(formModel)) {
      return
    }

    updateDisabled(true)
    app
      .test({
        req_fields: Object.keys(formModel).map(key => {
          const item = formModel[key]
          return {
            id: key,
            type: item.type,
            value: item.value || '',
            values: item.values || [],
            master_slave: item.masterSlave || '',
            required: item.required,
            custom_json_value_string: item.customJSONValueString || '{}',
            is_support_master: item.isSupportMaster ?? true, // 只保留支持主文件模式
            master_file: item.masterFile || '',
            is_support_workdir: item.isSupportWorkdir ?? false,
            workdir: item.workdir
          }
        }),
        upload_sub_token: uploadToken,
        name: app.name,
        state: app.state,
        script: getScriptData() || app.scriptData,
        scheduler: ''
      })
      .then(res => {
        const { debug_info, job_id } = res.data

        Modal.show({
          title: '测试报告',
          width: 800,
          bodyStyle: { height: 600, padding: '0 10px' },
          content: ({ onCancel, onOk }) => (
            <TestResult
              debugInfo={debug_info}
              jobId={job_id}
              reSubmitCallback={onOk}
            />
          ),
          footer: null,
          onCancel: async () => {
            // kill job
            console.log('kill job...')
          }
        })
      })
      .finally(() => updateDisabled(false))
  }

  render() {
    const { disabled, isRemote } = this.props

    return (
      <Wrapper>
        <div className='main'>
          {!isRemote && (
            <Button
              type='secondary'
              disabled={disabled}
              onClick={() => this.saveAs()}>
              另存为模版
            </Button>
          )}
          <Button
            type='secondary'
            disabled={disabled}
            onClick={() => this.onSave()}>
            仅保存
          </Button>
          <Button
            type='primary'
            disabled={disabled}
            onClick={() => this.onSave({ publish: true })}>
            保存并发布
          </Button>
          <Button
            type='primary'
            ghost
            disabled={disabled}
            onClick={this.onCancel}>
            取消
          </Button>
        </div>
      </Wrapper>
    )
  }
}
