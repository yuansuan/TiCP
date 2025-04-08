import styled from 'styled-components'

const nodeHeight = '40px'
export default styled.div`
  .ant-tree-child-tree {
    > li:first-child {
      padding-top: 0;
    }
  }

  .ant-tree {
    li {
      padding: 0;

      ul {
        /* hack: fix the firefox bug that height is zero */
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
