import React, { useState, useEffect } from 'react'
import { message, Tree, Button, Spin, Icon, Descriptions } from 'antd'
import { Modal } from '@/components'
import { Http } from '@/utils'
import { StyledTreeSelector } from './style'
import { formatByte } from '@/utils/Validator'

const { TreeNode, DirectoryTree } = Tree

interface IProps {
  title?: string
  rootPath: string
  workdir: string
  onClick?: any
  onSelect?: any
  disabled?: boolean
}

function bfsSearch<T>(
  list: Array<T>,
  pickFunc: (item: T) => boolean
): Array<T> {
  const queue = []
  queue.push(...list)
  const results = []
  
  while (queue.length > 0) {
    const item = queue.shift()
    if (pickFunc.apply(item, [item])) {
      results.push(item)
    } else {
      if (Array.isArray(item.children)) {
        queue.push(...item.children)
      }
    }
  }

  return results
}

function DirTreeSelector({rootPath, workdir, onOk, onCancel}) {

  const [treeData, setTreeData] = useState([])
  const [dir, setDir] = useState(workdir)
  const [expandedKeys, setExpandedKeys] = useState([])
  const [selectedKeys, setSelectedKeys] = useState([])
  const [selectedNode, setSelectedNode] = useState(null)
  const [loading, setLoading] = useState(false)


  useEffect(() => {
    (async () => {
      setLoading(true)

      try {
        const { data } = await Http.get('/file/ls', {
          params: { path: rootPath },
        })

        let _tmpData = data?.files?.filter(f => f.is_dir) ?? []

        setTreeData(_tmpData)

        if (workdir) { // 存在workdir目录继续下探

          let pathes = workdir.split(rootPath)
      
          let subFiles = pathes[1].split('/')
      
          subFiles.shift()
      
          let expandedKeys = []
      
          let prefix = rootPath + '/'
    
          while (subFiles.length) {
            let dirName = subFiles.shift()
            expandedKeys.push(`${prefix}${dirName}`)
            prefix = `${prefix}${dirName}/`
          }
    
          let promises = []
          let keys = []
    
          expandedKeys.forEach(key => {
            promises.push(() => {
              const currentNodeData = bfsSearch(_tmpData, item => item.path === key)?.[0]
              if (currentNodeData) {
                keys.push(key)
                return new Promise(((solve, reject) =>
                  Http.get('/file/ls', {
                    params: { path: currentNodeData.path },
                  }).then(({data}) => {
                    currentNodeData.children = data?.files?.filter(f => f.is_dir) ?? []
                    solve(true)
                  }).catch((e) => {
                    reject(false)
                    setLoading(false)
                  })
                ))
              } else {
                return Promise.resolve()
              }  
            })
          })

          promises.push(() => {
            return new Promise((resolve, reject) => {
              setTreeData([..._tmpData])
              setExpandedKeys(expandedKeys)
              setSelectedKeys([workdir])
              setLoading(false)
              const currentNodeData = bfsSearch(_tmpData, item => item.path === workdir)?.[0]
              setSelectedNode(currentNodeData)
              resolve(true)
            })
          })
          
          promises.reduce(
            (previousPromise, nextPromise) => previousPromise.then(() => nextPromise()),
            Promise.resolve()
          )
        } else {
          setLoading(false)
        }
      } catch(e) {
        message.error('数据加载异常，请确认选择的工作目录是否存在')
        setLoading(false)
      }
    })()

  }, [])


  const onLoadData = async (treeNode) => {
    if (treeNode.props.children) {
      return;
    }

    const { data } = await Http.get('/file/ls', {
      params: { path: treeNode.props.dataRef.path },
    })

    treeNode.props.dataRef.children = data?.files?.filter(f => f.is_dir) ?? []

    setTreeData([...treeData])
  }

  const onSelect = (selectedKeys, e) => {
    if (e.selected) {
      setDir(e.node?.props?.dataRef.path)
      setSelectedNode(e.node?.props?.dataRef) 
    } else {
      setDir('')
      setSelectedNode(null) 
    }
    setSelectedKeys(selectedKeys)
  }

  const renderTreeNodes = data =>
    data.map(item => {
      if (item.children) {
        return (
          <TreeNode title={item.name} key={item.path} dataRef={item} isLeaf={false}>
            {renderTreeNodes(item.children)}
          </TreeNode>
        )
      }
      return <TreeNode title={item.name} key={item.path} dataRef={item} isLeaf={false}/>
    })


  return (
    <StyledTreeSelector>
      <Spin tip="数据加载中..." spinning={loading}>
      <div className='header'>已选工作目录：{dir}</div>
      <div className='main'>
        <div className='tree'>
          <DirectoryTree 
            expandedKeys={expandedKeys}
            onExpand={(expandedKeys) => setExpandedKeys(expandedKeys)}
            selectedKeys={selectedKeys}
            onSelect={onSelect} 
            loadData={onLoadData}>
              {renderTreeNodes(treeData)}
          </DirectoryTree>
        </div>
        {selectedNode && 
          (<div className='preview'>
            <div className='name'>
              {selectedNode.name}
            </div>
            <div className='pic'>
              <Icon type="folder" style={{fontSize: 80}}/>
            </div>
            <div className='desc'>
              <Descriptions title="" size={'small'} column={1}>
                <Descriptions.Item label="大小">{formatByte(selectedNode.size)}</Descriptions.Item>
                <Descriptions.Item label="模式">{selectedNode.mode}</Descriptions.Item>
                <Descriptions.Item label="所有者">{selectedNode.owner.user}</Descriptions.Item>
                <Descriptions.Item label="所属组">{selectedNode.owner.group}</Descriptions.Item>
                <Descriptions.Item label="修改日期">{new Date(selectedNode.m_date*1000).toDateString()}</Descriptions.Item>
              </Descriptions>
            </div>
          </div>)
        }
      </div>
      <div className='footer'>
        <div className='footerMain'>
          <Button onClick={onCancel}>取消</Button>
          <Button
            disabled={!dir}
            type='primary'
            onClick={() => onOk(dir)}>
            确定
          </Button>
        </div>
      </div>
      </Spin>
    </StyledTreeSelector>
  )
}


export default class WorkdirAction extends React.Component<IProps> {

  render() {
    const { title, rootPath, workdir, onClick, onSelect, disabled } = this.props
    const children = React.Children.map(this.props.children, child =>
      React.cloneElement(child as React.ReactElement, {
        onClick: () => {
          if (disabled) return
          onClick && onClick()

          this.openDialog({
            title,
            rootPath,
            workdir,
            onSelect
          })
        },
      })
    )

    return <>{children}</>
  }

  private openDialog = ({
    title="选择工作目录",
    rootPath,
    workdir,
    onSelect
  }) => {
    Modal.show({
      title,
      width: 800,
      bodyStyle: {
        height: 600,
        padding: 10,
        background: '#fff'
      },
      content: ({ onCancel, onOk }) => {

        const _onOk = (newWorkdir) => {
          if (workdir) {
            Modal.showConfirm({
              content: '确认要切换工作目录?',
            }).then(() => {
              onSelect(newWorkdir)
              onOk()
            })
          } else {
            onSelect(newWorkdir)
            onOk()
          }
        } 

        return (<DirTreeSelector
          rootPath={rootPath}
          workdir={workdir}
          onCancel={onCancel}
          onOk={_onOk}
        />)
      },
      footer: null,
    })
  }
}
