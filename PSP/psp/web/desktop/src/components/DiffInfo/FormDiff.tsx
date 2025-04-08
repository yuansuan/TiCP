import React from 'react'

import { InlineDiff } from './Diff'
import SectionDiff from './SectionDiff'
import { StyledSection } from './style'

interface IProps {
  formDiff: any
}

export default class FormDiff extends React.Component<IProps> {
  render() {
    const { formDiff } = this.props
    const showAdd = formDiff.add.length > 0
    const showDelete = formDiff.delete.length > 0
    const showUpdate = formDiff.update.length > 0

    if (!showAdd && !showDelete && !showUpdate) {
      return null
    }

    return (
      <StyledSection>
        {(showAdd || showDelete || showUpdate) && (
          <div className='tag'>配置信息</div>
        )}
        {showDelete && (
          <>
            {formDiff.delete.map(item => (
              <InlineDiff
                key={item.name}
                Old={<span>Section（{item.name}）</span>}
              />
            ))}
          </>
        )}

        {showAdd && (
          <>
            {formDiff.add.map(item => (
              <InlineDiff
                key={item.name}
                New={<span>Section（{item.name}）</span>}
              />
            ))}
          </>
        )}

        {showUpdate && (
          <>
            {formDiff.update.map(item => (
              <div key={item.key}>
                <div className='tag'>更新Section（{item.key}）</div>
                <SectionDiff sectionDiff={item.field} />
              </div>
            ))}
          </>
        )}
      </StyledSection>
    )
  }
}
