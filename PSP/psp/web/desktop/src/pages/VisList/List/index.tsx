/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useState, useEffect, useRef } from 'react'
import { Context, useModel, useStore } from './store'
import { observer } from 'mobx-react-lite'
import { useAsync } from '@/utils/hooks'
import { Modal } from 'antd'
import { ListWrapper } from '../style'
import {
  ListSessionRequest,
  ListSessionResponse,
  StartSessionResponse
} from '@/domain/Vis'
import { Action } from './Action'
import { DataTable } from './DataTable'
import { PortableCreate } from '../Create/Portable'
import { Page } from '@/components'

type IProps = {
  isRefresh?: boolean
}
const List = observer((props: IProps) => {
  const store = useStore()
  const { vis } = store
  const [height, setHeight] = useState(800)
  const [hardware, setHardware] = useState([])
  const [software, setSoftware] = useState([])
  const [projects, setProjects] = useState([])

  const [hasPermHardware, setHasPermHardware] = useState([])
  const [hasPermSoftware, setHasPermSoftware] = useState([])

  const [isModalVisible, setIsModalVisible] = useState(false)
  const [response, setResponse] = useState(new ListSessionResponse())
  const [request, setRequest] = useState(
    new ListSessionRequest({ ...vis.filterQuery })
  )
  const { execute: listSession, loading } = useAsync(async () => {
    return await vis.listSession({ ...vis.filterQuery })
  }, true)
  const [isStopFetch, setStopFetch] = useState(false)

  const timerRef = useRef(null)
  const domRef = useRef(null)

  useEffect(() => {
    checkVisRinningOrClosing()
    // 防止实例在后台关闭时
    let timer = setInterval(() => {
      checkVisRinningOrClosing()
    }, 60 * 1000)

    return () => clearInterval(timer)
  }, [])

  useEffect(() => {
    listSession().then((res: any) => setResponse(res))

    vis
      .listAllHardware()
      .then(list => {
        setHardware(list)
        return vis.listAllSoftware()
      })
      .then(list => setSoftware(list))
    return () => {
      vis.setFilterParams({
        statuses: [],
        hardware_ids: [],
        software_ids: [],
        project_ids: [],
        page_index: 1,
        page_size: 20
      })
      setStopFetch(true)
      clearInterval(timerRef.current)
      timerRef.current = null
    }
  }, [])

  useEffect(() => {
    if (isStopFetch) {
      setStopFetch(true)
      clearInterval(timerRef.current)
      timerRef.current = null
    }
  }, [isStopFetch])

  useEffect(() => {
    const resizeObserver = new ResizeObserver(entries => {
      for (let entry of entries) {
        setHeight(entry.contentRect.height)
      }
    })

    resizeObserver.observe(domRef.current)

    // hack: 处理Table首次加载 bug
    setTimeout(() => {
      if (domRef.current) domRef.current.style.paddingRight = 1 + 'px'
    }, 3000)
    return () => resizeObserver.disconnect()
  }, [])

  const onPagination = (index: number, size: number) => {
    vis.setFilterParams({
      ...vis.filterQuery,
      page_index: index,
      page_size: size
    })
    setRequest({ ...vis.filterQuery })
    listSession().then((res: any) => setResponse(res))
  }

  const onQuerySession = async (values?: any) => {
    if (
      (values?.software_ids && values?.software_ids?.length !== 0) ||
      (values?.hardware_ids && values?.hardware_ids?.length !== 0) ||
      (values?.project_ids && values?.project_ids?.length !== 0) ||
      (values?.statuses && values?.statuses?.length !== 0) ||
      values?.user_id
    ) {
      await vis.listSession(values).then((res: any) => setResponse(res))
    } else {
      listSession().then((res: any) => setResponse(res))
    }
  }

  const checkVisRinningOrClosing = () => {
    clearInterval(timerRef.current)
    if (!isStopFetch || !timerRef.current) {
      timerRef.current = setInterval(() => {
        visualIsOpening()
      }, 10 * 1000)
    }
  }

  const onRestartSession = (id: string) => {
    return new Promise((resolve, reject) => {
      vis
        .restartSession(id)
        .then(res => {
          if (res.success) {
            return true
          } else {
            throw new Error('重启会话失败')
          }
        })
        .then(() => setStopFetch(false))
        .then(() => checkVisRinningOrClosing())
        .then(() => listSession())
        .then(res => {
          setResponse(res)
          resolve(id)
        })
        .catch(() => {
          setStopFetch(false)
          reject('重启会话失败')
        })
    })
  }

  const onCloseSession = (id: string) => {
    return new Promise((resolve, reject) => {
      vis
        .closeSession(id)
        .then(() => setStopFetch(false))
        .then(() => checkVisRinningOrClosing())
        .then(() => listSession())
        .then(res => {
          setResponse(res)
          resolve(id)
        })

        .catch(() => {
          setStopFetch(false)
          reject('删除会话失败')
        })
    })
  }

  const onPowerOnSession = (id: string) => {
    return new Promise((resolve, reject) => {
      vis
        .powerOnSession(id)
        .then(() => {
          setStopFetch(false)
        })
        .then(() => checkVisRinningOrClosing())
        .then(() => {
          return listSession()
        })
        .then(res => {
          setResponse(res)
          resolve(id)
        })
        .catch(() => {
          listSession().then(res => setResponse(res))
          setStopFetch(true)
          reject('开启会话失败')
        })
    })
  }

  const onDeleteSession = (id: string) => {
    return new Promise((resolve, reject) => {
      vis
        .deleteSession(id)
        .then(() => setStopFetch(false))
        .then(() => listSession())
        .then(res => {
          setResponse(res)
          resolve(id)
        })
        .catch(() => {
          setStopFetch(true)
          reject('删除会话失败')
        })
    })
  }

  const onUpdateSession = (id: string, autoClose: boolean, time?: string) => {
    return new Promise((resolve, reject) => {
      vis
        .updateSession(id, autoClose, time)
        .then(() => {
          setStopFetch(false)
        })
        .then(() => {
          return listSession()
        })
        .then(res => {
          setResponse(res)
          resolve(id)
        })
        .catch(() => {
          setStopFetch(true)
          reject('更新会话关闭时间失败')
        })
    })
  }

  // TODO 优化会话轮询
  const visualIsOpening = async () => {
    const lists = await vis.listSession({ ...request })

    setResponse(lists)
    const { sessions } = JSON.parse(JSON.stringify(lists))
    const findStartingVis = sessions.find(
      ({ session }) => session.status !== 'CLOSED'
    )
    const findClosingVis = sessions.find(
      ({ session }) => session.status === 'CLOSING'
    )
    const findRestartingVis = sessions.find(
      ({ session }) => session.status === 'REBOOTING'
    )
    const findPoweringVis = sessions.find(
      ({ session }) =>
        session.status === 'POWERING ON' || session.status === 'POWERING OFF'
    )

    if (
      findStartingVis?.loading ||
      findClosingVis ||
      findPoweringVis ||
      findRestartingVis
    ) {
      // 应用还在启动中或者关闭中
    } else {
      clearInterval(timerRef.current)
      timerRef.current = null
      onQuerySession(vis.filterQuery)
      return
    }
  }

  const createdSession = (response: StartSessionResponse) => {
    onQuerySession()
    setIsModalVisible(false)
    setStopFetch(false)
    checkVisRinningOrClosing()
  }

  const onOpenSession = async (id: string, row: any) => {
    // 后端是base64 加密的前端需要解密
    return window.atob(row.session.stream_url)
  }

  const onOpenRemoteApp = async (id: string, row: any) => {
    const app_name = row.session?.software?.remote_apps[0]?.name
    if (app_name) {
      return vis.getRemoteAppUrl(id, app_name)
    } else {
      return ''
    }
  }

  const setModalVisible = async () => {
    setStopFetch(false)
    const projects = await vis.getCurrentUserProjects()
    setProjects(projects)
    vis
      .listPermHardware()
      .then(list => {
        setHasPermHardware(list)
        return vis.listPermSoftware()
      })
      .then(list => setHasPermSoftware(list))
      .then(() => setIsModalVisible(true))
  }
  return (
    <>
      <ListWrapper ref={domRef}>
        {/* <Page header={null}> */}
        <Action
          hardware={hardware}
          software={software}
          value={request}
          onCreate={setModalVisible}
          onSubmit={value => {
            Object.assign(request, value)
            setRequest(request)
            vis.setFilterParams(request)
            onQuerySession(value)
          }}
        />
        <DataTable
          request={request}
          response={response}
          loading={loading}
          height={height - 160}
          onPagination={onPagination}
          deleteSession={onDeleteSession}
          closeSession={onCloseSession}
          restartSession={onRestartSession}
          powerOnSession={onPowerOnSession}
          updateSession={onUpdateSession}
          openSession={onOpenSession}
          openRemoteApp={onOpenRemoteApp}
        />
        {/* </Page> */}
      </ListWrapper>
      <Modal
        width={'80%'}
        footer={null}
        title='创建可视化会话'
        maskClosable={false}
        centered={true}
        visible={isModalVisible}
        destroyOnClose={true}
        onOk={() => setStopFetch(false)}
        onCancel={() => setIsModalVisible(false)}>
        <PortableCreate
          loading={loading}
          projects={projects}
          hardware={hasPermHardware}
          software={hasPermSoftware}
          vis={vis}
          onCreated={createdSession}
        />
      </Modal>
    </>
  )
})

export default observer((props: IProps) => {
  const model = useModel()

  return (
    <Context.Provider value={model}>
      <List {...props} />
    </Context.Provider>
  )
})
