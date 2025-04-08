import * as React from 'react'
import { Tree} from 'antd'
import { Icon } from '@/components'
import { Wrapper } from './style'

const { TreeNode } = Tree

interface IProps {
  sectionDiff: any
}

export default class SectionDiff extends React.Component<IProps> {
  render() {
    const { sectionDiff } = this.props
    const showAdd = sectionDiff.add.length > 0
    const showDelete = sectionDiff.delete.length > 0
    const showUpdate = sectionDiff.update.length > 0

    return (
      <Wrapper>
        {showAdd ? (
          <>
            <div className='sub-title'>新增 Section：</div>

            <Tree
              showIcon
              selectable={false}
              switcherIcon={<Icon type='down' /}>
              {sectionDiff.add.map((section, index) => (
                <TreeNode
                  key={index}
                  icon={<Icon type='tree' />}
                  title={section.name}>
                  {section.field.map(field => (
                    <TreeNode
                      key={field.id}
                      icon={<Icon type='setting' />}
                      title={field.label}
                    />
                  ))}
                </TreeNode>
              ))}
            </Tree>
          </>
        ) : null}

        {showDelete ? (
          <>
            <div className='sub-title'>删除 Section：</div>

            <Tree
              showIcon
              selectable={false}
              switcherIcon={<Icon type='down' />}>
              {sectionDiff.delete.map((section, index) => (
                <TreeNode
                  key={index}
                  icon={<Icon type='tree' />}
                  title={section.name}>
                  {section.field.map(field => (
                    <TreeNode
                      key={field.id}
                      icon={<Icon type='setting' />}
                      title={field.label}
                    />
                  ))}
                </TreeNode>
              ))}
            </Tree>
          </>
        ) : null}

        {showUpdate ? (
          <>
            <div className='sub-title'>更新 Section：</div>
            <Tree
              showIcon
              selectable={false}
              switcherIcon={<Icon type='down' />}>
              {sectionDiff.update.map((section, index) => {
                const { field } = section

                return (
                  <TreeNode
                    key={index}
                    icon={<Icon type='tree' />}
                    title={section.key}>
                    {field.add.length > 0 ? (
                      <TreeNode
                        key='add'
                        icon={<Icon type='tree' />}
                        title='新增控件'>
                        {field.add.map(field => (
                          <TreeNode
                            key={field.id}
                            icon={<Icon type='setting' />}
                            title={field.label}
                          />
                        ))}
                      </TreeNode>
                    ) : null}
                    {field.delete.length > 0 ? (
                      <TreeNode
                        key='delete'
                        icon={<Icon type='tree' />}
                        title='删除控件'>
                        {field.delete.map(field => (
                          <TreeNode
                            key={field.id}
                            icon={<Icon type='setting' />}
                            title={field.label}
                          />
                        ))}
                      </TreeNode>
                    ) : null}
                    {field.update.length > 0 ? (
                      <TreeNode
                        key='update'
                        icon={<Icon type='tree' />}
                        title='变更控件'>
                        {field.update.map(item => (
                          <TreeNode
                            key={item.key}
                            icon={<Icon type='setting' />}
                            title={item.key}>
                            {Object.keys(item.props).map(key => {
                              const prop = item.props[key]
                              return (
                                <TreeNode
                                  key={key}
                                  icon={<Icon type='setting' />}
                                  title={
                                    <>
                                      {Array.isArray(prop) ? (
                                        <span>{key}</span>
                                      ) : (
                                        <span>
                                          {key}: {prop.old || `""`} {'->'}{' '}
                                          {prop.new}
                                        </span>
                                      )}
                                    </>
                                  }
                                />
                              )
                            })}
                          </TreeNode>
                        ))}
                      </TreeNode>
                    ) : null}
                  </TreeNode>
                )
              })}
            </Tree>
          </>
        ) : null}
      </Wrapper>
    )
  }
}
