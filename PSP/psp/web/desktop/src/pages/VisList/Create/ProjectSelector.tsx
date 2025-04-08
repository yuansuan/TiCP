import React, { useState, useEffect, useCallback } from 'react'
import { ListActionWrapper } from '../style'
import { Project } from '@/domain/Vis'
import { Spin, Select, Tooltip, Checkbox, Row, Col, Descriptions, message } from 'antd'
import { observer } from 'mobx-react-lite'
import { Content, Layout, Sider, SiderTitle } from './styles'
import { Modal, Button } from '@/components'
import styled from 'styled-components'
import { Http } from '@/utils'
const tjsb = require('@/assets/images/tjsb.png')

const selectStyles = { width: '200px' }

const DeviceMountWrapper = styled.div`
  .body {
    .label {
      display: inline-block;
      font-size: 14px;
      font-family: PingFangSC-Regular;
      font-weight: 700;
    }
  }
  .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`
interface IProps {
  loading: boolean
  projects: Array<Project>
  onSelect?: (projectId: string) => void
  onMountSelect?: (projectNames: string[]) => void
}
type ISelectMount = {
  id: string
  name: string
}
type IMountInfo = {
  default_mounts: string[]
  select_mount:ISelectMount[]
  select_limit: number
}
type IDeviceMountProps = {
  onCancel: () => void
  onOk: () => void
  data: IMountInfo
  hasCheckedProjects: ISelectMount[]
  onMountSelect?: (projectNames: ISelectMount[]) => void
}

const DeviceMount = observer(
  ({
    onCancel,
    onOk,
    data,
    onMountSelect,
    hasCheckedProjects
  }: IDeviceMountProps) => {
    const [checkedProject, setCheckedProject] = useState(hasCheckedProjects)

    const onChange = checkedValues => {
      if(checkedValues.length > data.select_limit) {
        return message.warn('最多选择' + data.select_limit + '个项目存储！')
      }
      const checkedProject = data.select_mount.filter(item => checkedValues.includes(item.id))
      setCheckedProject(checkedProject)
    }

    const onSelectProject = useCallback(() => {
      onMountSelect(checkedProject)
      onOk()
    }, [checkedProject])


    return (
      <DeviceMountWrapper>
        <div className='body'>
          <Col span={6}>
            <div className='label'>
              项目存储：
            </div>
          </Col>
          <Checkbox.Group
            style={{ width: '100%' }}
            onChange={onChange}
            defaultValue={checkedProject.map(item => item.id)}
            value={checkedProject.map(item => item.id)}>
            <Row>
              {data.select_mount.map(item => (
                <Col span={5} key={item.id}>
                  <Checkbox value={item.id}>{item.name}</Checkbox>
                </Col>
              ))}
            </Row>
          </Checkbox.Group>
        </div>
        <Modal.Footer
          className='footer'
          onCancel={onCancel}
          OkButton={
            <Button type='primary' onClick={onSelectProject}>
              确认
            </Button>
          }
        />
      </DeviceMountWrapper>
    )
  }
)
export const ProjectSelector = observer(
  ({ projects, loading, onSelect, onMountSelect }: IProps) => {
    const [projectId, setProjectId] = useState(projects?.[0]?.id || '')
    const [checkedProjects, setCheckedProjects] = useState([])

    useEffect(() => {
      onSelect(projectId)
      // 切换所属项目，清空已选存储（项目名称）
      onMountSelect([])
      setCheckedProjects([])
    }, [projectId])

    const hasSelectProjectNames = checkProject=> {
      setCheckedProjects(checkProject)
      onMountSelect(checkProject.map(item => item.id))
    }
    const showSider = async () => {
      const res = await Http.get('/vis/session/getMountInfo', {
        params: {
          project_id: projectId
        }
      })
      Modal.show({
        title: '存储挂载',
        width: 600,
        bodyStyle: { overflow: 'auto', height: 500 },
        footer: null,
        content: ({ onCancel, onOk }) => (
          <DeviceMount
            onCancel={onCancel}
            onOk={onOk}
            data={res.data}
            hasCheckedProjects={checkedProjects}
            onMountSelect={hasSelectProjectNames}
          />
        )
      })
    }
    
    return (
      <Layout>
        <Sider>
          <SiderTitle>项目名称</SiderTitle>
        </Sider>
        <Content>
          <Spin spinning={loading}>
            <ListActionWrapper style={{ padding: 0 }}>
              <div className='item'>
                <Select
                  className={'status'}
                  value={projectId}
                  onChange={setProjectId}
                  style={selectStyles}>
                  {projects?.map(v => (
                    <Select.Option value={v.id} key={v.id}>
                      <Tooltip placement='topLeft' title={v.name}>
                        {v.name}
                      </Tooltip>
                    </Select.Option>
                  ))}
                </Select>
              </div>
              <div className='mountDevice'>
                <Tooltip title={'存储挂载'}>
                  <img src={tjsb} alt='存储挂载' onClick={showSider} />
                </Tooltip>
                <div className='showNames'>
                  {checkedProjects.length > 0 && (
                    <Descriptions>
                      <Descriptions.Item label='已选存储'>
                        {checkedProjects.map(item => item.name).join(', ')}
                      </Descriptions.Item>
                    </Descriptions>
                  )}
                </div>
              </div>
            </ListActionWrapper>
          </Spin>
          {!projectId && (
            <div className='validate_tip validate_soft_tip '>
              {'请选择项目名称'}
            </div>
          )}
        </Content>
      </Layout>
    )
  }
)
