import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Modal, Button } from '@/components'
import { Http } from '@/utils'
import { observer, useLocalStore } from 'mobx-react-lite'
import { message, Form, Input, Transfer } from 'antd'
import { Software, Preset } from '@/domain/VIsIBV/Software'
import { useStore } from './store'

const colProps = { labelCol: { span: 3 } }
const StyledLayout = styled.div`
  padding: 20px;

  .footer {
    position: absolute;
    bottom: 20px;
    right: 10px;
    padding: 10px;
  }
`

type Props = {
  softwareItem?: Software
  onCancel?: () => void
  onOk?: () => void
}

interface TransferItem {
  key: string
  title: string
  defaulted: boolean
  selected: boolean
}

type TransferItemMap = { [key: string]: TransferItem }

export default observer(function PresetEditor({
  softwareItem,
  onCancel,
  onOk
}: Props) {
  const [form] = Form.useForm()

  const state = useLocalStore(() => ({
    loading: false,
    setLoading(loading: boolean) {
      this.loading = loading
    },

    selectedHardware: [],
    setSelectedHardware(hardwareList: string[]) {
      Object.values(this.publishedHardware).forEach(
        (v: any) => (v.selected = false)
      )
      hardwareList?.forEach((key: string) => {
        if (this.publishedHardware[key]) {
          this.publishedHardware[key].selected = true
        }
      })
      this.selectedHardware = hardwareList
    },

    publishedHardware: {},
    setPublishedHardware(hardware: TransferItemMap) {
      this.publishedHardware = hardware
    },
    get hardwareDataSource() {
      return Object.values(this.publishedHardware)
    },
    setHardwareAsDefault(key: string) {
      const copy = {}
      Object.keys(this.publishedHardware).forEach(k => {
        copy[k] = this.publishedHardware[k]
        copy[k].defaulted = copy[k].key === key
      })
      this.setPublishedHardware(copy)
    }
  }))

  useEffect(() => {
    state.setLoading(true)

    // 获取所有可用的实例列表
    Http.get('/vis/hardware', {
      params: {
        page_size: 1000,
        page_index: 1
      }
    })
      .then(res => {
        const hardware: TransferItemMap = {}
        res.data.hardwares?.map((item: any) => {
          hardware[item.id] = {
            key: item.id,
            title: item.name,
            selected: false,
            defaulted: false
          }
        })

        softwareItem?.presets?.forEach((item: Preset) => {
          if (hardware[item.id]) {
            hardware[item.id].defaulted = item.default_preset
          }
        })

        // 设置可用实例列表状态
        state.setPublishedHardware(hardware)
        // 设置已选中的实例列表
        state.setSelectedHardware(
          softwareItem?.presets?.map((item: Preset) => item.id)
        )
      })
      .finally(() => state.setLoading(false))
  }, [])

  async function onFinish(values: any) {
    try {
      state.setLoading(true)

      await Http.post(
        '/vis/software/preset',
        {
          software_id: softwareItem.id,
          presets: state.selectedHardware?.map((k: string) => {
            return {
              hardware_id: k,
              default: state.publishedHardware[k]?.defaulted || false
            }
          })
        },
        {}
      )

      onOk()
      message.success('镜像预设成功')
    } finally {
      state.setLoading(false)
    }
  }

  return (
    <StyledLayout>
      <Form
        form={form}
        onFinish={onFinish}
        {...colProps}
        initialValues={{ ...softwareItem }}>
        <Form.Item
          label='镜像名称'
          name='software_name'
          initialValue={softwareItem.name}>
          <Input maxLength={64} disabled style={{ width: 200 }} />
        </Form.Item>
        <Form.Item
          label='预设'
          name='preset_hardware'
          initialValue={state.selectedHardware}>
          <Transfer
            dataSource={state.hardwareDataSource}
            targetKeys={state.selectedHardware}
            locale={{ itemUnit: '项', itemsUnit: '项' }}
            listStyle={{
              width: 310,
              height: 300
            }}
            titles={['全部实例', '当前预设']}
            onChange={state.setSelectedHardware}
            render={(item: any) => (
              <div>
                <span>{item.title}</span>
                {item.selected && (
                  <Button
                    type='link'
                    style={{ float: 'right', zIndex: 9999 }}
                    onClick={e => {
                      e.stopPropagation()
                      state.setHardwareAsDefault(item.key)
                    }}>
                    {item.defaulted ? '默认实例' : '设为默认'}
                  </Button>
                )}
              </div>
            )}></Transfer>
        </Form.Item>
        <Modal.Footer
          className='footer'
          onCancel={onCancel}
          OkButton={
            <Button
              type='primary'
              loading={state.loading}
              onClick={form.submit}>
              确认
            </Button>
          }
        />
      </Form>
    </StyledLayout>
  )
})
