import { Modal, Table, Button } from '@/components'
import { observer } from 'mobx-react-lite'
import React, { useEffect ,useRef} from 'react'
import { useLocalStore } from 'mobx-react'
import styled from 'styled-components'
import { QuantityAdd } from './QuantityAdd'
import { Http } from '@/utils'

const StyledLayout = styled.div`
  height: 100%;
  overflow:auto;
  .top {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
  }
  > .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 10px 17px 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`
type ModuleConf = {
  module_name: string
  id: string
  free_num: number
  total: number // 实时总数量（监控统计的）
}

type Props = {
  moduleInfos: ModuleConf
  licenseId: string
  usedPercent: string
  onOk?: () => void
  onCancel?: () => void
}
export const QuantityDetail = observer(
  ({ moduleInfos, usedPercent, licenseId, onOk }: Props) => {
    const state = useLocalStore(() => ({
      loading: false,
      height: 500,
      list: [],
      visible: false,
      moduleConf: null,
      setHeight(h) {
        this.height = h
      },
      setModuleConf(config) {
        this.moduleConf = config
      },
      setLoading(loading) {
        this.loading = loading
      },
      setVisible(bool) {
        this.visible = bool
      },
      setList(list) {
        this.list = list
      },
      usedPercent: '',
      setUsedPercent(usedPercent) {
        this.usedPercent = usedPercent
      }
    }))
    const ref = useRef<HTMLDivElement>()
    const interval = useRef(null)
    

    useEffect(() => {
      const resizeObserver = new ResizeObserver(entries => {
        for (let entry of entries) {
          state.setHeight(entry.contentRect.height)
        }
      })
  
      resizeObserver.observe(ref.current)
  
      setTimeout(() => {
        ref.current.style.paddingRight = 1 + 'px'
      }, 3000)
  
      return () => {
        resizeObserver && resizeObserver.disconnect()
      }
    }, [])
    
    const fetchModuleConfigs = async id => {
      const { data } = await Http.get(`/licenseInfos/${id}/moduleConfigs`)
      state.setList(data?.module_config_infos || [])
      state.setUsedPercent(data?.used_percent || '')
    }

    useEffect(() => {
      interval.current = setInterval(async () => {
         fetchModuleConfigs(licenseId);
      }, 10000);

      return () => {
        interval.current && clearInterval(interval.current)
      }
    },[])
    useEffect(() => {
      try {
        state.setLoading(true)
        state.setList(moduleInfos)
        state.setUsedPercent(usedPercent)
      } finally {
        state.setLoading(false)
      }
    }, [])

    const deleteModule = ({ id, module_name }) => {
      Modal.confirm({
        title: '删除模块配置',
        content: `确认删除「${module_name}」！`,
        okText: '确认',
        visible: state.visible,
        cancelText: '取消',
        onOk: async () => {
          await Http.delete(`moduleConfigs/${id}`)
            .then(res => {
              if (res.success) {
              }
            })
            .finally(() => {
              state.setVisible(false)
            })
          fetchModuleConfigs(licenseId)
        }
      })
    }

    const moduleConfigModal = () => {
      return Modal.show({
        title: `${
          Object.keys(state.moduleConf).length > 0 ? '编辑' : '添加'
        }模块`,
        width: 800,
        footer: null,
        content: ({ onCancel }) => (
          <QuantityAdd
            moduleConf={[state.moduleConf]}
            refresh={fetchModuleConfigs}
            licenseId={licenseId}
            onCancel={onCancel}
          />
        )
      })
    }
    return (
      <StyledLayout  ref={ref}>
        <div className='top'>
          <Button
            icon='add'
            type='primary'
            onClick={() => {
              state.setModuleConf({})
              moduleConfigModal()
            }}>
            添加
          </Button>
          {state.usedPercent && <span>总使用率: {state.usedPercent}</span>}
        </div>
        <Table
          props={{
            data: state.list,
            height: state.height,
            rowKey: 'rel_path',
            loading: state.loading
          }}
          columns={[
            {
              header: '模块',
              props: {
                fixed: true,
                flexGrow: 3
              },
              dataKey: 'module_name'
            },
            {
              header: '总数量',
              props: {
                fixed: true,
                flexGrow: 2
              },
              dataKey: 'total'
            },
            {
              header: '使用数量',
              props: {
                fixed: true,
                flexGrow: 2
              },
              dataKey: 'used_num'
            },
            {
              header: '剩余数量',
              props: {
                fixed: true,
                flexGrow: 2
              },
              dataKey: 'free_num'
            },
            {
              header: '操作',
              props: {
                fixed: true,
                flexGrow: 2
              },
              dataKey: 'operation',
              cell: {
                render: ({ rowData }) => {
                  return (
                    <div>
                      <Button
                        type='link'
                        onClick={() => {
                          state.setModuleConf(rowData)
                          moduleConfigModal()
                        }}>
                        编辑
                      </Button>
                      <Button type='link' onClick={() => deleteModule(rowData)}>
                        删除
                      </Button>
                    </div>
                  )
                }
              }
            }
          ]}
        />
        <Modal.Footer className='footer' CancelButton={null} onOk={onOk} />
      </StyledLayout>
    )
  }
)
