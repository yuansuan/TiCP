/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useEffect, useState } from 'react'
import { ProtalWrapper } from '../style'
import { Hardware, Software, StartSessionResponse,Project } from '@/domain/Vis'
import { Button, Modal } from '@/components'
import { observer } from 'mobx-react-lite'
import { HardwareSelector } from './HardwareSelector'
import { SoftwareSelector } from './SoftwareSelector'
import { ProjectSelector } from './ProjectSelector'
import { useStore } from '../List/store'
const centerStyle = {
  zIndex: 99999,
  position: 'absolute',
  left: '50%',
  top: '30%',
  transform: 'translate(-50%, 0)'
}

interface IProps {
  vis: any
  loading: boolean
  hardware: Array<Hardware>
  software: Array<Software>
  projects: Array<Project>
  onCreated?: (response: StartSessionResponse) => void
}

export const PortableCreate = observer(
  ({ vis, hardware, software, projects, loading, onCreated }: IProps) => {
    const [hardwareId, setHardwareId] = useState(null)
    const [softwareId, setSoftwareId] = useState(null)
    const [projectId, setProjectId] = useState(null)
    const [projectNames, onMountSelect] = useState(null)
    const [errMsg, setErrMsg] = useState('')
    const [preset, setPreset] = useState({ presets: [] })

    // 当修改 softwareID 时获取该软件的预设清单
    useEffect(() => {
      if (softwareId) {
        vis.autoChoseHardware(softwareId).then(res => {
          setPreset(res)
        })
      }
    }, [softwareId])

    useEffect(() => {
      setHardwareId(preset?.presets?.find(item => item.default_preset)?.id)
    }, [preset])
    useEffect(() => {
      if (errMsg) {
        setTimeout(() => {
          setErrMsg('')
        }, 5 * 1000)
      }
    }, [errMsg])

    const onStartSession = () => {
      if (!softwareId || !hardwareId) {
        let allAnimateDiv: NodeListOf<Element> =
          document.querySelectorAll('.validate_tip')
        allAnimateDiv &&
          Array.from(allAnimateDiv).map(sub => sub.classList.add('errtips'))

        setTimeout(() => {
          allAnimateDiv &&
            Array.from(allAnimateDiv).map(sub =>
              sub.classList.remove('errtips')
            )
        }, 2000)
      } else {
        const selectedSoftware = software.find(item => item.id === softwareId)
        const selectedHardtware = hardware?.find(item => item.id === hardwareId)
        const selectedProject = projects?.find(item => item.id === projectId)

        const pass = selectedSoftware.gpu_desired && selectedHardtware.gpu > 0

        Modal.showConfirm({
          title: '确认创建会话',
          content: pass
            ? '您当前所选镜像需要GPU才能正常工作，您是否确定使用当前配置创建会话？'
            : '确认按所选配置创建会话吗？',

          onOk() {
            return new Promise((resolve, reject) => {
              // check 是否可以创建新会话
              vis
                .startSession({
                  hardware_id: hardwareId,
                  software_id: softwareId,
                  project_id: selectedProject?.id,
                  mounts: projectNames
                })
                .then(res => {
                  resolve(res)
                  onCreated(res)
                  window.localStorage.setItem(
                    'CURRENTROUTERPATH',
                    `/vis-session`
                  )
                })
                .catch(() => {
                  resolve('创建会话失败')
                })
            }).catch((err) => {
              setErrMsg('创建会话失败')
            })
          }
        })
      }
    }

    const disabledBtn = () => {
      if (projectId && softwareId && hardwareId) {
        return false
      }
      return true
    }

    return (
      <ProtalWrapper>
        {errMsg && (
          <div className='ant-message-notice' style={centerStyle as any}>
            <div className='ant-message-notice-content'>
              <div className='ant-message-custom-content ant-message-error'>
                <span
                  role='img'
                  aria-label='check-circle'
                  className='anticon anticon-check-circle'>
                  <svg
                    viewBox='64 64 896 896'
                    focusable='false'
                    data-icon='close-circle'
                    width='1em'
                    height='1em'
                    fill='currentColor'
                    aria-hidden='true'>
                    <path d='M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm165.4 618.2l-66-.3L512 563.4l-99.3 118.4-66.1.3c-4.4 0-8-3.5-8-8 0-1.9.7-3.7 1.9-5.2l130.1-155L340.5 359a8.32 8.32 0 01-1.9-5.2c0-4.4 3.6-8 8-8l66.1.3L512 464.6l99.3-118.4 66-.3c4.4 0 8 3.5 8 8 0 1.9-.7 3.7-1.9 5.2L553.5 514l130 155c1.2 1.5 1.9 3.3 1.9 5.2 0 4.4-3.6 8-8 8z'></path>
                  </svg>
                </span>
                <span style={{ padding: '0 10px' }}>{errMsg}</span>
              </div>
            </div>
          </div>
        )}
        <ProjectSelector
          projects={projects}
          loading={loading}
          onSelect={setProjectId}
          onMountSelect={onMountSelect}
         />
        <SoftwareSelector
          software={software}
          loading={loading}
          onSelect={setSoftwareId}
        />
        <HardwareSelector
          hardware={hardware}
          loading={loading}
          defaultId={hardwareId}
          onSelect={setHardwareId}
        />
        <Modal.Footer
          className='footer'
          OkButton={
            <Button
              type='primary'
              style={{ marginRight: 10 }}
              disabled={disabledBtn()}
              onClick={onStartSession}>
              创建会话
            </Button>
          }
          CancelButton={null}
        />
      </ProtalWrapper>
    )
  }
)
