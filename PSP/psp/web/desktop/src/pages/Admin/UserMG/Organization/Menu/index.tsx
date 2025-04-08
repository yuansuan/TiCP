import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Tree } from 'antd'
import { useStore } from '../store'
import { observer } from 'mobx-react'

const { DirectoryTree } = Tree
const { TreeNode } = Tree
const nodeHeight = '40px'

const StyledLayout = styled.div`
  overflow: auto;
  height: 100%;

  .ant-tree-child-tree {
    > li:first-child {
      padding-top: 0;
    }
  }

  .ant-tree {
    li {
      padding: 0;

      ul {
        height: auto !important;
      }

      &:last-child {
        padding-bottom: 0;
      }

      &:first-child {
        padding-top: 0;
      }

      span.ant-tree-node-content-wrapper {
        height: ${nodeHeight};
        line-height: ${nodeHeight};
        width: calc(100% - 24px);

        &:hover {
          background-color: ${props => props.theme.backgroundColor};
        }
      }

      span.ant-tree-switcher {
        height: ${nodeHeight};
        line-height: ${nodeHeight};
      }
    }

    &.ant-tree-directory {
      > li {
        &.ant-tree-treenode-selected {
          > span.ant-tree-node-content-wrapper::before {
            background-color: #cbdcf5;
            height: ${nodeHeight};
          }

          > span.ant-tree-switcher {
            color: black;
          }
        }

        .ant-tree-child-tree {
          > li span.ant-tree-node-content-wrapper {
            &:hover {
              &::before {
                background-color: ${props => props.theme.backgroundColor};
              }
            }
          }
        }

        span.ant-tree-node-content-wrapper {
          width: calc(100% - 24px);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;

          &::before {
            height: ${nodeHeight};
          }

          &:hover {
            &::before {
              background-color: ${props => props.theme.backgroundColor};
            }

            background-color: ${props => props.theme.backgroundColor};
          }

          > span {
            color: #252525;
          }
        }

        span.ant-tree-iconEle,
        > span.ant-tree-switcher {
          height: ${nodeHeight};
          line-height: ${nodeHeight};
        }
      }

      .ant-tree-child-tree {
        > li {
          span.ant-tree-node-content-wrapper {
            &::before {
              height: ${nodeHeight};
            }

            > span {
              color: #252525;
            }
          }

          &.ant-tree-treenode-selected {
            > span.ant-tree-node-content-wrapper::before {
              background-color: #cbdcf5;
            }

            > span.ant-tree-switcher {
              color: black;
            }
          }
        }
      }
    }
  }
`

type Props = {}

export const Menu = observer(function Menu({}: Props) {
  const store = useStore()
  const [fetch, loading] = store.getOrganization()
  const { org } = store

  const firstNode = [org?.data.org]

  useEffect(() => {
    fetch()
  }, [])

  function updateNodeId(keys) {
    store.setPage(1, 10)
    store.setOrder(true, 'name')
    store.setSearchKey('')
    store.setNodeId(keys)
  }

  function renderTreeNodes(data) {
    return data.map(item => {
      if (item?.children) {
        return (
          <TreeNode title={item.org_name_cn} key={item.id}>
            {renderTreeNodes(item.children)}
          </TreeNode>
        )
      }
      return <TreeNode title={item?.org_name_cn} key={item?.id}></TreeNode>
    })
  }
  return (
    <StyledLayout>
      <DirectoryTree
        onSelect={updateNodeId}
        showIcon={false}
        selectedKeys={store.nodeId}>
        {renderTreeNodes(firstNode)}
      </DirectoryTree>
    </StyledLayout>
  )
})
