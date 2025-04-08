import React from 'react'
import { InlineDiff } from './Diff'

import { StyledSection } from './style'

interface IProps {
  binPathDiff: any
  title: ''
}

export default class BinPathDiff extends React.Component<IProps> {
  render() {
    const { binPathDiff,title } = this.props
    const showAdd = binPathDiff.add.length > 0
    const showDelete = binPathDiff.delete.length > 0
    const showUpdate = binPathDiff.update.length > 0

    if (!showAdd && !showDelete && !showUpdate) {
      return null
    }

    return (
      <StyledSection>
        {(showAdd || showDelete || showUpdate) && title &&  (
        <div className='tag'>{title}</div>
        )}
        {showDelete && (
          <>
            {binPathDiff.delete.map(item => (
              <InlineDiff
                name={item.key}
                key={item.key}
                Old={<span>{item.value}</span>}
              />
            ))}
          </>
        )}

        {showAdd && (
          <>
            {binPathDiff.add.map(item => (
              <InlineDiff
                name={item.key}
                key={item.key}
                New={<span>{item.value}</span>}
              />
            ))}
          </>
        )}

        {showUpdate && (
          <>
            {binPathDiff.update.map(item => (
              <div key={item.old.key}>
                <InlineDiff
                  key={item.old.key}
                  Old={item.old.key}
                  New={item.new.value}
                />
              </div>
            ))}
          </>
        )}
      </StyledSection>
    )
  }
}
