import { Button, Modal, Table } from '@/components'
import { formatUnixTime } from '@/utils'
import { formatByte } from '@/utils/Validator'
import { Input, message } from 'antd'
import { observer } from 'mobx-react-lite'
import React, { useCallback } from 'react'
import { useLocalStore } from 'mobx-react'
import styled from 'styled-components'
import { debounce } from 'lodash'
import { NewBoxHttp } from '@/domain/Box/NewBoxHttp'
import { env } from '@/domain'

const StyledLayout = styled.div`
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-right: 20px;
  }
  p {
    font-size: 16px;
    padding-top: 1em;
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
const MAX_SIZE = 10737418240 //10G

type Props = {
  onOk: (keys?: any[]) => void
  selectedJob: any
}
export const SearchFiles = observer(({ onOk, selectedJob }: Props) => {
  const state = useLocalStore(() => ({
    loading: false,
    name: '',
    files: [],
    setName(name) {
      this.name = name
    },
    setLoading(loading) {
      this.loading = loading
    },
    setFiles(list) {
      this.files = list
    },
    selectedKeys: [],
    setSelectedKeys(keys) {
      this.selectedKeys = keys
    },
    get searchDisabled() {
      if (this.name === '') {
        return '请输入关键字进行搜索'
      }
      return false
    },
    get downloadDisabled() {
      if (this.selectedKeys.length === 0) {
        return '当前无选中文件'
      }
      return false
    }
  }))

  const debouncedSetName = useCallback(
    debounce(name => {
      state.setName(name)
    }, 300),
    []
  )

  async function getFileList() {
    state.setLoading(true)
    try {
      await NewBoxHttp()
        .post('/filemanager/search', {
          keyword: state.name,
          job_ids: selectedJob
        })
        .then(res => {
          state.setFiles(res?.data || [])
        })
    } catch (e) {
      if (e['response']?.data?.code === 110018) {
        message.error(
          '您所要检索的文件数量过大，请输入更精确的关键字或减少选中的作业'
        )
      } else {
        message.error(`${e['response']?.data?.message}`)
      }
    }
    state.setLoading(false)
  }

  return (
    <StyledLayout>
      <div className='toolbar'>
        <div>
          关键字：
          <Input
            allowClear
            style={{ width: 280 }}
            placeholder='请输入文件名称关键字'
            onChange={e => {
              debouncedSetName(e.target.value)
            }}
          />
        </div>
        <Button onClick={getFileList} disabled={state.searchDisabled}>
          搜索
        </Button>
      </div>
      <p>在选中的作业文件中进行搜索</p>
      <Table
        props={{
          data: state.files,
          height: 400,
          rowKey: 'rel_path',
          loading: state.loading
        }}
        rowSelection={{
          selectedRowKeys: state.selectedKeys,
          onChange: keys => {
            const selectedRow = state.files
              .filter(file => keys?.find(path => path === file.rel_path))
              .map(file => file['size'])

            const totalSize =
              selectedRow.length !== 0 &&
              selectedRow.reduce((pre, current) => pre + current)

            if (totalSize >= MAX_SIZE) {
              message.info('当前所选文件已超过10GB,请重新勾选')

              state.setSelectedKeys([...state.selectedKeys])
            } else {
              state.setSelectedKeys(keys)
            }
          }
        }}
        columns={[
          {
            header: '文件名称',
            props: {
              fixed: true,
              flexGrow: 3
            },
            dataKey: 'name'
          },
          {
            header: '作业编号',
            props: {
              fixed: true,
              flexGrow: 2
            },
            dataKey: 'job_id'
          },
          {
            header: '创建时间',
            props: {
              fixed: true,
              flexGrow: 2
            },
            dataKey: 'mod_time',
            cell: {
              render: ({ rowData }) => (
                <div>{formatUnixTime(rowData.mod_time)}</div>
              )
            }
          },
          {
            header: '大小',
            props: {
              fixed: true,
              flexGrow: 1
            },
            dataKey: 'size',
            cell: {
              render: ({ rowData }) => <div>{formatByte(rowData.size)}</div>
            }
          }
        ]}
      />
      <Modal.Footer
        className='footer'
        CancelButton={null}
        OkButton={
          <Button
            type='primary'
            onClick={() => {
              onOk(state.selectedKeys)
              state.setSelectedKeys([])
            }}
            disabled={state.downloadDisabled}>
            下载
          </Button>
        }
      />
    </StyledLayout>
  )
})
